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

package native

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/decimal"
	"ballerina-lang-go/model"
	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

var httpPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("http")},
	model.Name("0.0.1"),
)

const (
	orgName    = "ballerina"
	moduleName = "http"
)

func init() {
	runtime.RegisterModuleInitializer(initHttpModule)
}

// httpTypes holds the lazily-built semtypes used by the http runtime.
// All entries are computed once on first use (so concurrent reads shouldn't be a problem);
type httpTypes struct {
	byteArrTy  semtypes.SemType
	strArrTy   semtypes.SemType
	jsonListTy semtypes.SemType
	jsonMapTy  semtypes.SemType
}

// newMappingValue builds a fresh open `map<anydata|error>` value.
func newMappingValue(tc semtypes.Context) *values.Map {
	return values.NewMap(semtypes.MAPPING, semtypes.ToMappingAtomicType(tc, semtypes.MAPPING), false, nil)
}

// newListValue builds a fresh open `(anydata|error)[]` value seeded with the
// supplied items.
func newListValue(tc semtypes.Context, items []values.BalValue) *values.List {
	return values.NewList(semtypes.LIST, semtypes.ToListAtomicType(tc, semtypes.LIST), false, nil, 0, items)
}

// newTypedListValue builds a fresh list with the supplied inherent type seeded
// with the supplied items. The atomic representation must exist; callers pass
// concrete list types built via list defs.
func newTypedListValue(tc semtypes.Context, ty semtypes.SemType, items []values.BalValue) *values.List {
	return values.NewList(ty, semtypes.ToListAtomicType(tc, ty), false, nil, 0, items)
}

func initHttpModule(rt *runtime.Runtime) {
	// Register module-level constants so BIR global-variable loads of http:LEADING
	// and http:TRAILING resolve correctly. Keys use buildGlobalVarLookupKey format:
	// org/pkg:varName = "ballerina/http:LEADING".
	runtime.RegisterModuleGlobals(rt, httpPackageID, map[string]values.BalValue{
		"ballerina/http:LEADING":  "LEADING",
		"ballerina/http:TRAILING": "TRAILING",
	})

	var (
		once  sync.Once
		types httpTypes
	)
	ensureTypes := func() {
		once.Do(func() {
			env := rt.GetTypeEnv()
			bld := semtypes.NewListDefinition()
			types.byteArrTy = bld.DefineListTypeWrappedWithEnvSemType(env, semtypes.BYTE)
			sld := semtypes.NewListDefinition()
			types.strArrTy = sld.DefineListTypeWrappedWithEnvSemType(env, semtypes.STRING)
			typCtx := semtypes.ContextFrom(env)
			jsonTy := semtypes.CreateJSON(typCtx)
			jmd := semtypes.NewMappingDefinition()
			types.jsonMapTy = jmd.DefineMappingTypeWrapped(env, nil, jsonTy)
			jld := semtypes.NewListDefinition()
			types.jsonListTy = jld.DefineListTypeWrappedWithEnvSemType(env, jsonTy)
		})
	}

	msgToBody := func(tc semtypes.Context, msg values.BalValue) ([]byte, string) {
		ensureTypes()
		switch v := msg.(type) {
		case string:
			return []byte(v), "text/plain"
		case *values.List:
			if v.Type != nil && semtypes.IsSubtype(tc, v.Type, types.byteArrTy) {
				if b, ok := listToBytes(v); ok {
					return b, "application/octet-stream"
				}
			}
			b, err := toJSONBytes(v)
			if err != nil {
				return nil, "json_error"
			}
			return b, "application/json"
		default:
			b, err := toJSONBytes(v)
			if err != nil {
				return nil, "json_error"
			}
			return b, "application/json"
		}
	}
	execBody := func(ctx *extern.Context, verb string, args []values.BalValue) (values.BalValue, error) {
		self := args[0].(*values.Object)
		path := args[1].(string)
		var body []byte
		contentType := ""
		if len(args) > 2 && args[2] != nil {
			body, contentType = msgToBody(ctx.TypeCtx, args[2])
			if body == nil && contentType == "json_error" {
				return values.NewErrorWithMessage("failed to serialize body to JSON"), nil
			}
		}
		var reqHeaders map[string][]string
		if len(args) > 3 {
			reqHeaders = extractHeaders(args[3])
			for hdrKey, hdrVals := range reqHeaders {
				if strings.EqualFold(hdrKey, "content-type") && len(hdrVals) > 0 {
					contentType = hdrVals[0]
					break
				}
			}
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
		return buildResponse(ctx.TypeCtx, statusCode, respHeaders, respBody), nil
	}

	// Remote method name uses the "$remote$" prefix (model.RemoteMethodName).
	// The BIR gen emits `callInfo.Name = "$remote$get"` for c->get(...), which
	// resolveObjectMethod then looks up in the object's methodKeys map.
	clientClassDef := &bir.BIRClassDef{
		Name:      model.Name("Client"),
		LookupKey: "ballerina/http:Client",
		Fields: []bir.ObjectField{
			{Name: "url", Ty: semtypes.STRING},
			{Name: "timeout", Ty: semtypes.DECIMAL},
			{Name: "followRedirects", Ty: semtypes.Union(semtypes.MAPPING, semtypes.NIL)},
			{Name: "httpVersion", Ty: semtypes.STRING},
		},
		VTable: map[string]*bir.BIRFunction{
			"initExtern":      {FunctionLookupKey: "ballerina/http:Client.initExtern"},
			"$remote$get":     {FunctionLookupKey: "ballerina/http:Client.$remote$get"},
			"$remote$post":    {FunctionLookupKey: "ballerina/http:Client.$remote$post"},
			"$remote$head":    {FunctionLookupKey: "ballerina/http:Client.$remote$head"},
			"$remote$options": {FunctionLookupKey: "ballerina/http:Client.$remote$options"},
			"$remote$put":     {FunctionLookupKey: "ballerina/http:Client.$remote$put"},
			"$remote$patch":   {FunctionLookupKey: "ballerina/http:Client.$remote$patch"},
			"$remote$delete":  {FunctionLookupKey: "ballerina/http:Client.$remote$delete"},
			"$remote$execute": {FunctionLookupKey: "ballerina/http:Client.$remote$execute"},
		},
	}
	runtime.RegisterExternClassDef(rt, clientClassDef)

	runtime.RegisterExternFunction(rt, orgName, moduleName, "parseHeader",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			input, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("parseHeader: expected string argument")
			}
			result, err := parseHeader(ctx.TypeCtx, input)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return result, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.initExtern",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			url := args[1].(string)

			timeout := decimal.FromInt64(30)
			var followRedirects pal.FollowRedirects // Enabled=false by default (Ballerina spec)
			httpVersion := "2.0"

			var tlsCfg pal.TLSConfig
			if len(args) > 2 {
				if cfg, ok := args[2].(*values.Map); ok {
					if v, ok := cfg.Get("timeout"); ok {
						if d, ok := v.(*decimal.Decimal); ok {
							timeout = d
						}
					}
					if v, ok := cfg.Get("followRedirects"); ok {
						if frMap, ok := v.(*values.Map); ok {
							if ev, ok := frMap.Get("enabled"); ok {
								if b, ok := ev.(bool); ok {
									followRedirects.Enabled = b
								}
							}
							followRedirects.MaxCount = 5
							if mv, ok := frMap.Get("maxCount"); ok {
								if n, ok := mv.(int64); ok {
									followRedirects.MaxCount = int(n)
								}
							}
							if av, ok := frMap.Get("allowAuthHeaders"); ok {
								if b, ok := av.(bool); ok {
									followRedirects.AllowAuthHeaders = b
								}
							}
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
									data, err := rt.Platform().FS.ReadFile(certPath)
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
											data, err := rt.Platform().FS.ReadFile(p)
											if err != nil {
												return values.NewErrorWithMessage("secureSocket.key.certFile: " + err.Error()), nil
											}
											tlsCfg.ClientCertPEM = data
										}
									}
									if kv, ok := keyMap.Get("keyFile"); ok {
										if p, ok := kv.(string); ok && p != "" {
											data, err := rt.Platform().FS.ReadFile(p)
											if err != nil {
												return values.NewErrorWithMessage("secureSocket.key.keyFile: " + err.Error()), nil
											}
											tlsCfg.ClientKeyPEM = data
										}
									}
									// keyPassword: accepted at compile time, ignored at runtime
								}
							}
							if v, ok := ssMap.Get("serverName"); ok {
								if s, ok := v.(string); ok && s != "" {
									tlsCfg.ServerName = s
								}
							}
							if v, ok := ssMap.Get("shareSession"); ok {
								if b, ok := v.(bool); ok && !b {
									tlsCfg.DisableSessionTickets = true
								}
							}
							if v, ok := ssMap.Get("handshakeTimeout"); ok {
								if d, ok := v.(*decimal.Decimal); ok {
									tlsCfg.HandshakeTimeout = decimalToDuration(d)
								}
							}
							if v, ok := ssMap.Get("ciphers"); ok {
								if list, ok := v.(*values.List); ok {
									for i := 0; i < list.Len(); i++ {
										if name, ok := list.Get(i).(string); ok {
											tlsCfg.CipherSuiteNames = append(tlsCfg.CipherSuiteNames, name)
										}
									}
								}
							}
							if v, ok := ssMap.Get("protocol"); ok {
								if protoMap, ok := v.(*values.Map); ok {
									if vv, ok := protoMap.Get("versions"); ok {
										if list, ok := vv.(*values.List); ok {
											tlsVersionMap := map[string]uint16{
												"TLSv1.0": 0x0301,
												"TLSv1.1": 0x0302,
												"TLSv1.2": 0x0303,
												"TLSv1.3": 0x0304,
											}
											for i := 0; i < list.Len(); i++ {
												if s, ok := list.Get(i).(string); ok {
													if ver, found := tlsVersionMap[s]; found {
														if tlsCfg.MinVersion == 0 || ver < tlsCfg.MinVersion {
															tlsCfg.MinVersion = ver
														}
														if ver > tlsCfg.MaxVersion {
															tlsCfg.MaxVersion = ver
														}
													}
												}
											}
										}
									}
								}
							}
							// certValidation/sessionTimeout: accepted at compile time, not supported at runtime
						}
					}
				}
			}
			httpClient := rt.Platform().HTTP.NewClient(pal.ClientConfig{
				Timeout:         decimalToDuration(timeout),
				FollowRedirects: followRedirects,
				HTTPVersion:     httpVersion,
				TLS:             tlsCfg,
			})
			self.Put("url", url)
			self.Put("timeout", timeout)
			self.Put("followRedirects", nil)
			self.Put("httpVersion", httpVersion)
			self.Put("$httpClient", httpClient)
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$get",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
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
			return buildResponse(ctx.TypeCtx, statusCode, respHeaders, body), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$post",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return execBody(ctx, "POST", args)
		})

	// head: body-less, like get
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$head",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
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
			return buildResponse(ctx.TypeCtx, statusCode, respHeaders, body), nil
		})

	// options: body-less, like get
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$options",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
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
			return buildResponse(ctx.TypeCtx, statusCode, respHeaders, body), nil
		})

	// put: body required, like post
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$put",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return execBody(ctx, "PUT", args)
		})

	// patch: body required, like post
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$patch",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return execBody(ctx, "PATCH", args)
		})

	// delete: message is optional (defaults to nil = empty body)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$delete",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return execBody(ctx, "DELETE", args)
		})

	// execute: args = [self, httpVerb, path, message, headers?, mediaType?]
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Client.$remote$execute",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			verb := args[1].(string)
			path := args[2].(string)

			var body []byte
			contentType := ""
			if len(args) > 3 && args[3] != nil {
				body, contentType = msgToBody(ctx.TypeCtx, args[3])
				if body == nil && contentType == "json_error" {
					return values.NewErrorWithMessage("failed to serialize body to JSON"), nil
				}
			}

			var reqHeaders map[string][]string
			if len(args) > 4 {
				reqHeaders = extractHeaders(args[4])
				for hdrKey, hdrVals := range reqHeaders {
					if strings.EqualFold(hdrKey, "content-type") && len(hdrVals) > 0 {
						contentType = hdrVals[0]
						break
					}
				}
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
			return buildResponse(ctx.TypeCtx, statusCode, respHeaders, respBody), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getTextPayload",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			body, _ := self.Get("body")
			return body, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getJsonPayload",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			self := args[0].(*values.Object)
			bodyVal, _ := self.Get("body")
			body := bodyVal.(string)
			dec := json.NewDecoder(strings.NewReader(body))
			dec.UseNumber()
			var v interface{}
			if err := dec.Decode(&v); err != nil {
				return values.NewErrorWithMessage("failed to parse JSON payload: " + err.Error()), nil
			}
			return goToBalValue(ctx.TypeCtx, v, types.jsonListTy, types.jsonMapTy), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getBinaryPayload",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			self := args[0].(*values.Object)
			bodyVal, _ := self.Get("body")
			body := bodyVal.(string)
			raw := []byte(body)
			items := make([]values.BalValue, len(raw))
			for i, b := range raw {
				items[i] = int64(b)
			}
			return newTypedListValue(ctx.TypeCtx, types.byteArrTy, items), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.hasHeader",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			// args[2] is position — ignored
			_, ok := responseHeaders(self).Get(name)
			return ok, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeader",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			v, ok := responseHeaders(self).Get(name)
			if !ok {
				return values.NewErrorWithMessage("header not found: " + name), nil
			}
			list := v.(*values.List)
			if list.Len() == 0 {
				return values.NewErrorWithMessage("header has no values: " + name), nil
			}
			return list.Get(0), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeaders",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			name := strings.ToLower(args[1].(string))
			v, ok := responseHeaders(self).Get(name)
			if !ok {
				return values.NewErrorWithMessage("header not found: " + name), nil
			}
			return v.(*values.List), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Response.getHeaderNames",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			self := args[0].(*values.Object)
			keys := responseHeaders(self).Keys()
			items := make([]values.BalValue, len(keys))
			for i, k := range keys {
				items[i] = k
			}
			return newTypedListValue(ctx.TypeCtx, types.strArrTy, items), nil
		})
}

// splitOutsideQuotes splits s on every occurrence of sep that is not inside a
// double-quoted string (RFC 7230 §3.2.6 quoted-string), honouring backslash escapes.
func splitOutsideQuotes(s string, sep byte) []string {
	var out []string
	inQuote := false
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\\' && inQuote && i+1 < len(s):
			i++ // skip the escaped character
		case c == '"':
			inQuote = !inQuote
		case c == sep && !inQuote:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

func parseHeader(tc semtypes.Context, input string) (*values.List, error) {
	segments := splitOutsideQuotes(input, ',')
	list := newListValue(tc, nil)
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			return nil, fmt.Errorf("invalid header value: empty segment")
		}
		parts := splitOutsideQuotes(seg, ';')
		headerVal := strings.TrimSpace(parts[0])
		if headerVal == "" {
			return nil, fmt.Errorf("invalid header value: missing value before parameters")
		}
		params := newMappingValue(tc)
		for _, param := range parts[1:] {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}
			eqIdx := strings.IndexByte(param, '=')
			if eqIdx < 0 {
				params.Put(tc, strings.ToLower(param), "")
				continue
			}
			key := strings.ToLower(strings.TrimSpace(param[:eqIdx]))
			val := strings.TrimSpace(param[eqIdx+1:])
			if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}
			params.Put(tc, key, val)
		}
		entry := newMappingValue(tc)
		entry.Put(tc, "value", headerVal)
		entry.Put(tc, "params", params)
		list.Append(tc, entry)
	}
	return list, nil
}

func decimalToDuration(d *decimal.Decimal) time.Duration {
	return time.Duration(d.Float64() * float64(time.Second))
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

// buildResponse constructs a Ballerina Response object from HTTP response data.
// All header values are stored as *values.List under the internal "$headers" key.
func buildResponse(tc semtypes.Context, statusCode int, respHeaders map[string][]string, body []byte) *values.Object {
	headersMap := newMappingValue(tc)
	for k, vals := range respHeaders {
		items := make([]values.BalValue, len(vals))
		for i, v := range vals {
			items[i] = v
		}
		headersMap.Put(tc, strings.ToLower(k), newListValue(tc, items))
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

// responseHeaders returns the internal header map stored on a Response object.
func responseHeaders(self *values.Object) *values.Map {
	h, _ := self.Get("$headers")
	return h.(*values.Map)
}

// listToBytes converts a Ballerina byte[] (List of int64 in 0–255) to []byte.
// Returns (nil, false) if any element is not an integer in the byte range,
// indicating the list should be treated as a JSON array instead.
func listToBytes(list *values.List) ([]byte, bool) {
	b := make([]byte, list.Len())
	for i := range list.Len() {
		n, ok := list.Get(i).(int64)
		if !ok || n < 0 || n > 255 {
			return nil, false
		}
		b[i] = byte(n)
	}
	return b, true
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
	case *decimal.Decimal:
		// Emit the decimal128 string verbatim as a JSON number so the full
		// precision of the value is preserved — going through Float64() truncates
		// past ~17 significant digits.
		return json.RawMessage(t.String())
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

// goToBalValue converts a Go value (from json.Decoder with UseNumber) to a Ballerina BalValue.
// JSON null → nil, bool → bool, json.Number → int64 or float64, string → string,
// []interface{} → *values.List with json[] type, map[string]interface{} → *values.Map with map<json> type.
// jsonListTy and jsonMapTy must be the structural json[] and map<json> semtypes so that
// `value is json` type checks return true for the produced values.
func goToBalValue(tc semtypes.Context, v interface{}, jsonListTy, jsonMapTy semtypes.SemType) values.BalValue {
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
		items := make([]values.BalValue, len(v))
		for i, elem := range v {
			items[i] = goToBalValue(tc, elem, jsonListTy, jsonMapTy)
		}
		return newTypedListValue(tc, jsonListTy, items)
	case map[string]interface{}:
		m := values.NewMap(jsonMapTy, semtypes.ToMappingAtomicType(tc, jsonMapTy), false, nil)
		for k, val := range v {
			m.Put(tc, k, goToBalValue(tc, val, jsonListTy, jsonMapTy))
		}
		return m
	default:
		return nil
	}
}
