// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package http

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/lib/http/compile"
	"ballerina-lang-go/model"
	"ballerina-lang-go/pal"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "http"
)

func parseHeader(input string) (*values.List, error) {
	segments := strings.Split(input, ",")
	list := values.NewList(0, semtypes.LIST, nil)
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			return nil, fmt.Errorf("invalid header value: empty segment")
		}
		parts := strings.Split(seg, ";")
		headerVal := strings.TrimSpace(parts[0])
		if headerVal == "" {
			return nil, fmt.Errorf("invalid header value: missing value before parameters")
		}
		params := values.NewMap(semtypes.MAPPING)
		for _, param := range parts[1:] {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}
			eqIdx := strings.IndexByte(param, '=')
			if eqIdx < 0 {
				params.Put(strings.ToLower(param), "")
				continue
			}
			key := strings.ToLower(strings.TrimSpace(param[:eqIdx]))
			val := strings.TrimSpace(param[eqIdx+1:])
			if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}
			params.Put(key, val)
		}
		entry := values.NewMap(semtypes.MAPPING)
		entry.Put("value", headerVal)
		entry.Put("params", params)
		list.Append(entry)
	}
	return list, nil
}

func ratToDuration(r *big.Rat) time.Duration {
	f, _ := r.Float64()
	return time.Duration(f * float64(time.Second))
}

func initHttpModule(rt *runtime.Runtime) {
	// Register module-level constants so BIR global-variable loads of http:LEADING
	// and http:TRAILING resolve correctly. Keys use buildGlobalVarLookupKey format:
	// org/pkg:varName = "ballerina/http:LEADING".
	runtime.RegisterModuleGlobals(rt, compile.HttpPackageID, map[string]values.BalValue{
		"ballerina/http:LEADING":  "LEADING",
		"ballerina/http:TRAILING": "TRAILING",
	})

	// Remote method name uses the "$remote$" prefix (model.RemoteMethodName).
	// The BIR gen emits `callInfo.Name = "$remote$get"` for c->get(...), which
	// resolveObjectMethod then looks up in the object's methodKeys map.
	clientClassDef := &bir.BIRClassDef{
		Name:      model.Name("Client"),
		LookupKey: "ballerina/http:Client",
		Fields: []bir.ObjectField{
			{Name: "url", Ty: semtypes.STRING},
			{Name: "timeout", Ty: semtypes.DECIMAL},
			{Name: "followRedirects", Ty: semtypes.BOOLEAN},
			{Name: "httpVersion", Ty: semtypes.STRING},
		},
		VTable: map[string]*bir.BIRFunction{
			"init":             {FunctionLookupKey: "ballerina/http:Client.init"},
			"$remote$get":     {FunctionLookupKey: "ballerina/http:Client.get"},
			"$remote$post":    {FunctionLookupKey: "ballerina/http:Client.post"},
			"$remote$head":    {FunctionLookupKey: "ballerina/http:Client.head"},
			"$remote$options": {FunctionLookupKey: "ballerina/http:Client.options"},
			"$remote$put":     {FunctionLookupKey: "ballerina/http:Client.put"},
			"$remote$patch":   {FunctionLookupKey: "ballerina/http:Client.patch"},
			"$remote$delete":  {FunctionLookupKey: "ballerina/http:Client.delete"},
			"$remote$execute": {FunctionLookupKey: "ballerina/http:Client.execute"},
		},
	}
	runtime.RegisterExternClassDef(rt, clientClassDef)

	// Default lambda for config param: called as $Client.init$default$1(url) → returns {}.
	// Receives [url] (the preceding arg) and ignores it; the default is always {}.
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.init$default$1",
		func(args []values.BalValue) (values.BalValue, error) {
			return values.NewMap(semtypes.MAPPING), nil
		})

	// Default lambda for headers param: called as $Client.get$default$1(path) → returns () = nil.
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.get$default$1",
		func(args []values.BalValue) (values.BalValue, error) {
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "parseHeader",
		func(args []values.BalValue) (values.BalValue, error) {
			input, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("parseHeader: expected string argument")
			}
			result, err := parseHeader(input)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return result, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.init",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			url := args[1].(string)

			timeout := big.NewRat(30, 1)
			followRedirects := true
			httpVersion := "2.0"

			var tlsCfg pal.TLSConfig
			if len(args) > 2 {
				if cfg, ok := args[2].(*values.Map); ok {
					if v, ok := cfg.Get("timeout"); ok {
						if rat, ok := v.(*big.Rat); ok {
							timeout = rat
						}
					}
					if v, ok := cfg.Get("followRedirects"); ok {
						if b, ok := v.(bool); ok {
							followRedirects = b
						}
					}
					if v, ok := cfg.Get("httpVersion"); ok {
						if s, ok := v.(string); ok {
							httpVersion = s
						}
					}
					if ss, ok := cfg.Get("secureSocket"); ok {
						if ssMap, ok := ss.(*values.Map); ok {
							if v, ok := ssMap.Get("enable"); ok {
								if b, ok := v.(bool); ok && !b {
									tlsCfg.InsecureSkipVerify = true
								}
							}
							if v, ok := ssMap.Get("verifyHostName"); ok {
								if b, ok := v.(bool); ok && !b {
									tlsCfg.InsecureSkipVerify = true
								}
							}
							if v, ok := ssMap.Get("cert"); ok {
								if certPath, ok := v.(string); ok && certPath != "" {
									data, err := os.ReadFile(certPath)
									if err != nil {
										return values.NewErrorWithMessage("secureSocket.cert: " + err.Error()), nil
									}
									tlsCfg.CACertPEM = data
								}
							}
							if v, ok := ssMap.Get("key"); ok {
								if keyMap, ok := v.(*values.Map); ok {
									if cv, ok := keyMap.Get("certFile"); ok {
										if p, ok := cv.(string); ok && p != "" {
											data, err := os.ReadFile(p)
											if err != nil {
												return values.NewErrorWithMessage("secureSocket.key.certFile: " + err.Error()), nil
											}
											tlsCfg.ClientCertPEM = data
										}
									}
									if kv, ok := keyMap.Get("keyFile"); ok {
										if p, ok := kv.(string); ok && p != "" {
											data, err := os.ReadFile(p)
											if err != nil {
												return values.NewErrorWithMessage("secureSocket.key.keyFile: " + err.Error()), nil
											}
											tlsCfg.ClientKeyPEM = data
										}
									}
									// keyPassword: accepted at compile time, ignored at runtime
								}
							}
							// ciphers/shareSession/handshakeTimeout/sessionTimeout/serverName: accepted, ignored
						}
					}
				}
			}
			httpClient := rt.Platform().HTTP.NewClient(pal.ClientConfig{
				Timeout:         ratToDuration(timeout),
				FollowRedirects: followRedirects,
				HTTPVersion:     httpVersion,
				TLS:             tlsCfg,
			})
			self.Put("url", url)
			self.Put("timeout", timeout)
			self.Put("followRedirects", followRedirects)
			self.Put("httpVersion", httpVersion)
			self.Put("$httpClient", httpClient)
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.get",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			path := args[1].(string)

			var reqHeaders map[string][]string
			if len(args) > 2 {
				reqHeaders = extractHeaders(args[2])
			}

			urlVal, _ := self.Get("url")
			clientHandle, _ := self.Get("$httpClient")
			statusCode, respHeaders, body, err := clientHandle.(pal.HTTPClient).Execute("GET", urlVal.(string)+path, nil, "", reqHeaders)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return buildResponse(statusCode, respHeaders, body), nil
		})

	// Default lambdas for post optional params (both return nil = Ballerina ())
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.post$default$2",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.post$default$3",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.post",
		func(args []values.BalValue) (values.BalValue, error) {
			return execBodyMethod("POST", args)
		})

	// head: body-less, like get
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.head$default$1",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.head",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			path := args[1].(string)
			var reqHeaders map[string][]string
			if len(args) > 2 {
				reqHeaders = extractHeaders(args[2])
			}
			urlVal, _ := self.Get("url")
			clientHandle, _ := self.Get("$httpClient")
			statusCode, respHeaders, body, err := clientHandle.(pal.HTTPClient).Execute("HEAD", urlVal.(string)+path, nil, "", reqHeaders)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return buildResponse(statusCode, respHeaders, body), nil
		})

	// options: body-less, like get
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.options$default$1",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.options",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			path := args[1].(string)
			var reqHeaders map[string][]string
			if len(args) > 2 {
				reqHeaders = extractHeaders(args[2])
			}
			urlVal, _ := self.Get("url")
			clientHandle, _ := self.Get("$httpClient")
			statusCode, respHeaders, body, err := clientHandle.(pal.HTTPClient).Execute("OPTIONS", urlVal.(string)+path, nil, "", reqHeaders)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return buildResponse(statusCode, respHeaders, body), nil
		})

	// put: body required, like post
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.put$default$2",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.put$default$3",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.put",
		func(args []values.BalValue) (values.BalValue, error) {
			return execBodyMethod("PUT", args)
		})

	// patch: body required, like post
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.patch$default$2",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.patch$default$3",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.patch",
		func(args []values.BalValue) (values.BalValue, error) {
			return execBodyMethod("PATCH", args)
		})

	// delete: message is optional (defaults to nil = empty body)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.delete$default$1",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.delete$default$2",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.delete$default$3",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.delete",
		func(args []values.BalValue) (values.BalValue, error) {
			return execBodyMethod("DELETE", args)
		})

	// execute: args = [self, httpVerb, path, message, headers?, mediaType?]
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.execute$default$3",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Client.execute$default$4",
		func(args []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.execute",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			verb := args[1].(string)
			path := args[2].(string)

			var body []byte
			contentType := ""
			if len(args) > 3 && args[3] != nil {
				body, contentType = messageToBody(args[3])
				if body == nil && contentType == "json_error" {
					return values.NewErrorWithMessage("failed to serialize body to JSON"), nil
				}
			}

			var reqHeaders map[string][]string
			if len(args) > 4 {
				reqHeaders = extractHeaders(args[4])
			}
			if len(args) > 5 {
				if mt, ok := args[5].(string); ok && mt != "" {
					contentType = mt
				}
			}

			urlVal, _ := self.Get("url")
			clientHandle, _ := self.Get("$httpClient")
			statusCode, respHeaders, respBody, err := clientHandle.(pal.HTTPClient).Execute(verb, urlVal.(string)+path, body, contentType, reqHeaders)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return buildResponse(statusCode, respHeaders, respBody), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getTextPayload",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			body, _ := self.Get("body")
			return body, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getJsonPayload",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			bodyVal, _ := self.Get("body")
			body := bodyVal.(string)
			dec := json.NewDecoder(strings.NewReader(body))
			dec.UseNumber()
			var v interface{}
			if err := dec.Decode(&v); err != nil {
				return values.NewErrorWithMessage("failed to parse JSON payload: " + err.Error()), nil
			}
			return goToBalValue(v), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getBinaryPayload",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			bodyVal, _ := self.Get("body")
			body := bodyVal.(string)
			raw := []byte(body)
			list := values.NewList(len(raw), semtypes.LIST, nil)
			for i, b := range raw {
				list.FillingSet(i, int64(b))
			}
			return list, nil
		})

	// Default lambdas for position param (all return "LEADING")
	leading := values.BalValue("LEADING")
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Response.hasHeader$default$1",
		func(_ []values.BalValue) (values.BalValue, error) { return leading, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Response.getHeader$default$1",
		func(_ []values.BalValue) (values.BalValue, error) { return leading, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Response.getHeaders$default$1",
		func(_ []values.BalValue) (values.BalValue, error) { return leading, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Response.getHeaderNames$default$0",
		func(_ []values.BalValue) (values.BalValue, error) { return leading, nil })

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.hasHeader",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			// args[2] is position — ignored
			_, ok := responseHeaders(self).Get(name)
			return ok, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeader",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			v, ok := responseHeaders(self).Get(name)
			if !ok {
				return values.NewErrorWithMessage("header not found: " + name), nil
			}
			return v.(*values.List).Get(0), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeaders",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			v, ok := responseHeaders(self).Get(name)
			if !ok {
				return values.NewErrorWithMessage("header not found: " + name), nil
			}
			return v.(*values.List), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeaderNames",
		func(args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			keys := responseHeaders(self).Keys()
			list := values.NewList(len(keys), semtypes.LIST, nil)
			for i, k := range keys {
				list.FillingSet(i, k)
			}
			return list, nil
		})
}

// extractHeaders converts a Ballerina map<string|string[]>? value to Go request headers.
func extractHeaders(arg values.BalValue) map[string][]string {
	if arg == nil {
		return nil
	}
	hdrMap, ok := arg.(*values.Map)
	if !ok {
		return nil
	}
	result := make(map[string][]string, hdrMap.Len())
	for _, key := range hdrMap.Keys() {
		val, _ := hdrMap.Get(key)
		switch v := val.(type) {
		case string:
			result[key] = []string{v}
		case *values.List:
			strs := make([]string, v.Len())
			for i := range v.Len() {
				if s, ok := v.Get(i).(string); ok {
					strs[i] = s
				}
			}
			result[key] = strs
		}
	}
	return result
}

// listToBytes converts a Ballerina byte[] (List of int64) to []byte.
func listToBytes(list *values.List) []byte {
	b := make([]byte, list.Len())
	for i := range list.Len() {
		if n, ok := list.Get(i).(int64); ok {
			b[i] = byte(n)
		}
	}
	return b
}

// goToBalValue converts a Go value (from json.Decoder with UseNumber) to a Ballerina BalValue.
// JSON null → nil, bool → bool, json.Number → int64 or float64, string → string,
// []interface{} → *values.List, map[string]interface{} → *values.Map.
func goToBalValue(v interface{}) values.BalValue {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return v
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		f, _ := v.Float64()
		return f
	case string:
		return v
	case []interface{}:
		list := values.NewList(len(v), semtypes.LIST, nil)
		for i, elem := range v {
			list.FillingSet(i, goToBalValue(elem))
		}
		return list
	case map[string]interface{}:
		m := values.NewMap(semtypes.MAPPING)
		for k, val := range v {
			m.Put(k, goToBalValue(val))
		}
		return m
	default:
		return nil
	}
}

// responseHeaders returns the internal header map stored on a Response object.
func responseHeaders(self *values.Object) *values.Map {
	h, _ := self.Get("$headers")
	return h.(*values.Map)
}

// buildResponse constructs a Ballerina Response object from HTTP response data.
// All header values are stored as *values.List under the internal "$headers" key.
func buildResponse(statusCode int, respHeaders map[string][]string, body []byte) *values.Object {
	headersMap := values.NewMap(semtypes.MAPPING)
	for k, vals := range respHeaders {
		list := values.NewList(len(vals), semtypes.LIST, nil)
		for i, v := range vals {
			list.FillingSet(i, v)
		}
		headersMap.Put(strings.ToLower(k), list)
	}
	return values.NewObject(
		semtypes.OBJECT,
		map[string]values.BalValue{
			"statusCode": int64(statusCode),
			"$headers":   headersMap,
			"body":       string(body),
		},
		map[string]string{
			"getTextPayload":   "ballerina/http:Response.getTextPayload",
			"getJsonPayload":   "ballerina/http:Response.getJsonPayload",
			"getBinaryPayload": "ballerina/http:Response.getBinaryPayload",
			"hasHeader":        "ballerina/http:Response.hasHeader",
			"getHeader":        "ballerina/http:Response.getHeader",
			"getHeaders":       "ballerina/http:Response.getHeaders",
			"getHeaderNames":   "ballerina/http:Response.getHeaderNames",
		},
	)
}

// messageToBody converts a Ballerina message value to a []byte body and a default Content-Type.
// Returns (nil, "json_error") on JSON serialization failure.
func messageToBody(msg values.BalValue) ([]byte, string) {
	switch v := msg.(type) {
	case string:
		return []byte(v), "text/plain"
	case *values.List:
		return listToBytes(v), "application/octet-stream"
	default:
		b, err := toJSONBytes(v)
		if err != nil {
			return nil, "json_error"
		}
		return b, "application/json"
	}
}

// execBodyMethod implements PUT/PATCH/DELETE externs.
// args = [self, path, message, headers?, mediaType?]
func execBodyMethod(verb string, args []values.BalValue) (values.BalValue, error) {
	self := args[0].(*values.Object)
	path := args[1].(string)

	var body []byte
	contentType := ""
	if len(args) > 2 && args[2] != nil {
		body, contentType = messageToBody(args[2])
		if body == nil && contentType == "json_error" {
			return values.NewErrorWithMessage("failed to serialize body to JSON"), nil
		}
	}

	var reqHeaders map[string][]string
	if len(args) > 3 {
		reqHeaders = extractHeaders(args[3])
	}
	if len(args) > 4 {
		if mt, ok := args[4].(string); ok && mt != "" {
			contentType = mt
		}
	}

	urlVal, _ := self.Get("url")
	clientHandle, _ := self.Get("$httpClient")
	statusCode, respHeaders, respBody, err := clientHandle.(pal.HTTPClient).Execute(verb, urlVal.(string)+path, body, contentType, reqHeaders)
	if err != nil {
		return values.NewErrorWithMessage(err.Error()), nil
	}
	return buildResponse(statusCode, respHeaders, respBody), nil
}

// balToGoJSON converts a Ballerina value to a Go value suitable for json.Marshal.
// Handles all Ballerina json-compatible types: nil, bool, int, float, decimal, string, map, list.
func balToGoJSON(v values.BalValue) any {
	switch t := v.(type) {
	case nil:
		return nil
	case bool:
		return t
	case int64:
		return t
	case float64:
		return t
	case *big.Rat:
		f, _ := t.Float64()
		return f
	case string:
		return t
	case *values.Map:
		m := make(map[string]any, t.Len())
		for _, k := range t.Keys() {
			val, _ := t.Get(k)
			m[k] = balToGoJSON(val)
		}
		return m
	case *values.List:
		s := make([]any, t.Len())
		for i := range t.Len() {
			s[i] = balToGoJSON(t.Get(i))
		}
		return s
	default:
		return nil
	}
}

// toJSONBytes serializes a Ballerina value to JSON bytes.
func toJSONBytes(v values.BalValue) ([]byte, error) {
	return json.Marshal(balToGoJSON(v))
}

func init() {
	runtime.RegisterModuleInitializer(initHttpModule)
}
