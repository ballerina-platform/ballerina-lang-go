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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
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

// listenerState is the Go-side state of an http:Listener object, stored on the
// object's "$state" field. The HTTP server is created lazily in Listener.start
// (driven by the module's $start lifecycle hook); attach only registers
// services. The program stays alive while the runtime is in its listening
// state — the runtime lifecycle owns signal handling and shutdown.
type listenerState struct {
	host        string
	port        int
	timeout     time.Duration
	httpVersion string
	tlsCfg      *tls.Config
	mu          sync.RWMutex
	services    []*serviceEntry
	server      *http.Server
}

type serviceEntry struct {
	basePath string
	svcObj   *values.Object
}

// registerListenerExterns registers the Listener class definition and its
// extern methods. Called from initHttpModule.
func registerListenerExterns(rt *runtime.Runtime) {
	listenerClassDef := &bir.BIRClassDef{
		Name:      model.Name("Listener"),
		LookupKey: "ballerina/http:Listener",
		Fields:    []bir.ObjectField{},
		VTable: map[string]*bir.BIRFunction{
			"initNative":    {FunctionLookupKey: "ballerina/http:Listener.initNative"},
			"attach":        {FunctionLookupKey: "ballerina/http:Listener.attach"},
			"detach":        {FunctionLookupKey: "ballerina/http:Listener.detach"},
			"start":         {FunctionLookupKey: "ballerina/http:Listener.start"},
			"gracefulStop":  {FunctionLookupKey: "ballerina/http:Listener.gracefulStop"},
			"immediateStop": {FunctionLookupKey: "ballerina/http:Listener.immediateStop"},
		},
	}
	runtime.RegisterExternClassDef(rt, listenerClassDef)

	// Default lambdas for the optional config/name parameters (both default to ()).
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Listener.init$default$1",
		func(_ *extern.Context, _ []values.BalValue) (values.BalValue, error) { return nil, nil })
	runtime.RegisterExternFunction(rt, orgName, moduleName, "$Listener.attach$default$1",
		func(_ *extern.Context, _ []values.BalValue) (values.BalValue, error) { return nil, nil })

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.initNative",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			port := int(args[1].(int64))
			state := &listenerState{
				host:        "0.0.0.0",
				port:        port,
				timeout:     60 * time.Second,
				httpVersion: "2.0",
			}
			if len(args) > 2 {
				if cfg, ok := args[2].(*values.Map); ok {
					if v, ok := cfg.Get("host"); ok {
						if s, ok := v.(string); ok && s != "" {
							state.host = s
						}
					}
					if v, ok := cfg.Get("timeout"); ok {
						if d, ok := v.(*decimal.Decimal); ok {
							state.timeout = decimalToDuration(d)
						}
					}
					if v, ok := cfg.Get("httpVersion"); ok {
						if s, ok := v.(string); ok && s != "" {
							state.httpVersion = s
						}
					}
					if v, ok := cfg.Get("secureSocket"); ok {
						if ssMap, ok := v.(*values.Map); ok {
							tlsCfg, err := buildListenerTLSConfig(ssMap, rt.Platform().FS)
							if err != nil {
								return values.NewErrorWithMessage("Listener.init secureSocket: " + err.Error()), nil
							}
							state.tlsCfg = tlsCfg
						}
					}
				}
			}
			self.Put("$state", state)
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.attach",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			svcObj, ok := args[1].(*values.Object)
			if !ok {
				return values.NewErrorWithMessage("Listener.attach: expected service object"), nil
			}
			var attachArg values.BalValue
			if len(args) > 2 {
				attachArg = args[2]
			}
			basePath := extractAttachPath(attachArg)
			stateVal, _ := self.Get("$state")
			state := stateVal.(*listenerState)

			if msg := validateServiceForHTTP(svcObj); msg != "" {
				return values.NewErrorWithMessage("Listener.attach: " + msg), nil
			}

			state.mu.Lock()
			entry := &serviceEntry{basePath: basePath, svcObj: svcObj}
			state.services = append(state.services, entry)
			// Longest base path first so the most specific service wins routing.
			sort.Slice(state.services, func(i, j int) bool {
				return len(state.services[i].basePath) > len(state.services[j].basePath)
			})
			state.mu.Unlock()
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.detach",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			svcObj := args[1].(*values.Object)
			stateVal, _ := self.Get("$state")
			state := stateVal.(*listenerState)
			state.mu.Lock()
			defer state.mu.Unlock()
			for i, e := range state.services {
				if e.svcObj == svcObj {
					state.services = append(state.services[:i], state.services[i+1:]...)
					break
				}
			}
			return nil, nil
		})

	// Listener.start creates and starts the HTTP server. It is invoked by the
	// module's $start lifecycle hook after all services have been attached.
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.start",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			stateVal, ok := self.Get("$state")
			if !ok {
				return values.NewErrorWithMessage("Listener.start: listener not initialised"), nil
			}
			state := stateVal.(*listenerState)
			state.mu.Lock()
			alreadyStarted := state.server != nil
			state.mu.Unlock()
			if alreadyStarted {
				return nil, nil
			}
			server, err := startHTTPServer(rt, state)
			if err != nil {
				return values.NewErrorWithMessage("Listener.start: " + err.Error()), nil
			}
			state.mu.Lock()
			state.server = server
			state.mu.Unlock()
			return nil, nil
		})

	// Listener.gracefulStop drains in-flight requests before closing the server.
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.gracefulStop",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			stateVal, ok := self.Get("$state")
			if !ok {
				return nil, nil
			}
			state := stateVal.(*listenerState)
			state.mu.RLock()
			server := state.server
			state.mu.RUnlock()
			if server != nil {
				_ = server.Shutdown(context.Background())
			}
			return nil, nil
		})

	// Listener.immediateStop closes the server and all active connections at once.
	runtime.RegisterExternFunction(rt, orgName, moduleName, "Listener.immediateStop",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			self := args[0].(*values.Object)
			stateVal, ok := self.Get("$state")
			if !ok {
				return nil, nil
			}
			state := stateVal.(*listenerState)
			state.mu.RLock()
			server := state.server
			state.mu.RUnlock()
			if server != nil {
				_ = server.Close()
			}
			return nil, nil
		})
}

// extractAttachPath converts the Ballerina attach-point value to a base path string.
// () → "/", "foo" → "/foo", ["a","b"] → "/a/b"
func extractAttachPath(v values.BalValue) string {
	if v == nil {
		return "/"
	}
	switch val := v.(type) {
	case string:
		if val == "" {
			return "/"
		}
		if !strings.HasPrefix(val, "/") {
			return "/" + val
		}
		return val
	case *values.List:
		parts := make([]string, val.Len())
		for i := range val.Len() {
			if s, ok := val.Get(i).(string); ok {
				parts[i] = s
			}
		}
		return "/" + strings.Join(parts, "/")
	}
	return "/"
}

// buildListenerTLSConfig builds a *tls.Config from a ListenerSecureSocket map.
func buildListenerTLSConfig(ssMap *values.Map, fs pal.FS) (*tls.Config, error) {
	keyVal, ok := ssMap.Get("key")
	if !ok {
		return nil, fmt.Errorf("secureSocket.key is required")
	}
	keyMap, ok := keyVal.(*values.Map)
	if !ok {
		return nil, fmt.Errorf("secureSocket.key must be a CertKey record")
	}

	certFileVal, _ := keyMap.Get("certFile")
	keyFileVal, _ := keyMap.Get("keyFile")
	certFilePath, _ := certFileVal.(string)
	keyFilePath, _ := keyFileVal.(string)

	certPEM, err := fs.ReadFile(certFilePath)
	if err != nil {
		return nil, fmt.Errorf("key.certFile: %w", err)
	}
	keyPEM, err := fs.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("key.keyFile: %w", err)
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("X509KeyPair: %w", err)
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	// mTLS: client certificate verification.
	if v, ok := ssMap.Get("mutualSsl"); ok {
		if b, ok := v.(bool); ok && b {
			if caCertPathVal, ok := ssMap.Get("cert"); ok {
				if caCertPath, ok := caCertPathVal.(string); ok && caCertPath != "" {
					caCertPEM, err := fs.ReadFile(caCertPath)
					if err != nil {
						return nil, fmt.Errorf("secureSocket.cert (CA): %w", err)
					}
					pool := x509.NewCertPool()
					if !pool.AppendCertsFromPEM(caCertPEM) {
						return nil, fmt.Errorf("failed to parse CA certificate")
					}
					tlsCfg.ClientCAs = pool
					tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
				}
			}
		}
	}

	// TLS version bounds.
	if v, ok := ssMap.Get("protocol"); ok {
		if list, ok := v.(*values.List); ok {
			tlsVersionMap := map[string]uint16{
				"TLSv1.0": tls.VersionTLS10, "TLSv1.1": tls.VersionTLS11,
				"TLSv1.2": tls.VersionTLS12, "TLSv1.3": tls.VersionTLS13,
			}
			for i := range list.Len() {
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

	// Cipher suites.
	if v, ok := ssMap.Get("ciphers"); ok {
		if list, ok := v.(*values.List); ok {
			allSuites := append(tls.CipherSuites(), tls.InsecureCipherSuites()...)
			nameToID := make(map[string]uint16, len(allSuites))
			for _, cs := range allSuites {
				nameToID[cs.Name] = cs.ID
			}
			for i := range list.Len() {
				if s, ok := list.Get(i).(string); ok {
					if id, found := nameToID[s]; found {
						tlsCfg.CipherSuites = append(tlsCfg.CipherSuites, id)
					}
				}
			}
		}
	}

	// Session tickets.
	if v, ok := ssMap.Get("shareSession"); ok {
		if b, ok := v.(bool); ok && !b {
			tlsCfg.SessionTicketsDisabled = true
		}
	}

	return tlsCfg, nil
}

// validateServiceForHTTP rejects service objects that contain remote methods,
// which are not supported for HTTP dispatch. Normal and resource methods are
// allowed. Returns a non-empty error message if validation fails.
func validateServiceForHTTP(svcObj *values.Object) string {
	if svcObj.HasRemoteMethods() {
		return "service object must not have remote methods"
	}
	return ""
}

// startHTTPServer creates the net/http server, binds the listening socket
// (optionally wrapped in TLS), and serves on a background goroutine.
func startHTTPServer(rt *runtime.Runtime, state *listenerState) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, fmt.Sprintf("%v", rec), http.StatusInternalServerError)
			}
		}()
		dispatchRequest(rt, state, w, r)
	})

	addr := fmt.Sprintf("%s:%d", state.host, state.port)
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	if state.httpVersion == "2.0" {
		protocols.SetHTTP2(true)
		if state.tlsCfg == nil {
			protocols.SetUnencryptedHTTP2(true)
		}
	}

	writeTimeout := state.timeout
	if writeTimeout == 0 {
		writeTimeout = 60 * time.Second
	}
	server := &http.Server{
		Addr:      addr,
		Handler:   mux,
		Protocols: protocols,
		// ReadHeaderTimeout guards against slow-loris without aborting request bodies
		// mid-stream (ReadTimeout would do that and breaks proxies/uploads).
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      writeTimeout,
		// Evict idle HTTP/1.1 keep-alive connections that haven't been used.
		IdleTimeout: 300 * time.Second,
		// 16 KB max header size; the default 1 MB is wasteful for typical REST APIs.
		MaxHeaderBytes: 1 << 14,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	var serveLn net.Listener = ln
	if state.tlsCfg != nil {
		serveLn = tls.NewListener(ln, state.tlsCfg)
	}
	go func() {
		_ = server.Serve(serveLn)
	}()
	return server, nil
}

// dispatchRequest routes an incoming HTTP request to the matching service and
// resource method.
func dispatchRequest(rt *runtime.Runtime, state *listenerState, w http.ResponseWriter, r *http.Request) {
	state.mu.RLock()
	var found *serviceEntry
	var subPath string
	for _, e := range state.services {
		if strings.HasPrefix(r.URL.Path, e.basePath) {
			found = e
			subPath = r.URL.Path[len(e.basePath):]
			break
		}
	}
	state.mu.RUnlock()

	if found == nil {
		writeErrorJSON(w, r, http.StatusNotFound, "no matching resource found for path")
		return
	}

	segments := splitURLPath(subPath)
	ctx := rt.NewExternContext()

	httpMethod := strings.ToLower(r.Method)
	for _, accessorKey := range []string{httpMethod, "default"} {
		candidates, ok := found.svcObj.ResourceEntries(accessorKey)
		if !ok {
			continue
		}
		for i := range candidates {
			coerced, ok := coercePathForCandidate(ctx.TypeCtx, &candidates[i], segments)
			if !ok {
				continue
			}
			handle, ok := ctx.LookupResourceMethod(found.svcObj, accessorKey, coerced)
			if !ok {
				continue
			}
			// Count non-literal path params to determine how many user args the method expects.
			nonLiteralCount := 0
			for _, seg := range candidates[i].PathSegments {
				if _, isLit := values.LiteralPathSegment(seg); !isLit {
					nonLiteralCount++
				}
			}
			totalParams := runtime.GetBIRFunctionParamCount(rt, candidates[i].FunctionLookupKey)
			extraArgCount := 0
			if totalParams >= 0 {
				extraArgCount = totalParams - nonLiteralCount
			}

			var invocationArgs []values.BalValue
			if extraArgCount > 0 {
				var bodyBuf []byte
				var bodyStream io.ReadCloser
				cl := r.ContentLength
				if r.Body == nil || cl == 0 {
					// no body or explicitly empty
				} else if cl >= 0 && cl <= eagerBufferThreshold {
					data, _ := io.ReadAll(r.Body)
					_ = r.Body.Close()
					bodyBuf = data
					cl = int64(len(data))
				} else {
					bodyStream = r.Body
				}
				reqObj := buildRequest(ctx.TypeCtx, r.Method, r.URL.Path, r.Proto, r.Header, bodyStream, cl, r.URL.RawQuery, bodyBuf)
				invocationArgs = []values.BalValue{reqObj}
			} else if r.Body != nil {
				// Resource method does not take a Request parameter; discard the body.
				_ = r.Body.Close()
			}
			result, err := ctx.InvokeMethod(handle, invocationArgs)
			if err != nil {
				writeErrorJSON(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			writeResult(ctx.TypeCtx, w, r, result)
			return
		}
	}
	// Path matched a service but no accessor+path combination worked. Check whether the
	// path would have matched under a different HTTP method and return 405 if so.
	for _, accessor := range found.svcObj.AllResourceMethodNames() {
		if accessor == httpMethod || accessor == "default" {
			continue
		}
		candidates, _ := found.svcObj.ResourceEntries(accessor)
		for i := range candidates {
			if _, ok := coercePathForCandidate(ctx.TypeCtx, &candidates[i], segments); ok {
				writeErrorJSON(w, r, http.StatusMethodNotAllowed, "method not allowed for path")
				return
			}
		}
	}
	writeErrorJSON(w, r, http.StatusNotFound, "no matching resource found for path")
}

// splitURLPath splits a URL sub-path into segments, stripping leading/trailing slashes.
func splitURLPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

// coercePathForCandidate tries to coerce URL string segments to the typed values expected
// by the candidate resource entry. Returns (nil, false) if the segments don't match.
func coercePathForCandidate(tc semtypes.Context, entry *values.ResourceEntry, segments []string) ([]values.BalValue, bool) {
	required := len(entry.PathSegments)
	hasRest := !semtypes.IsNever(entry.RestSegmentTy)
	if len(segments) < required {
		return nil, false
	}
	if len(segments) > required && !hasRest {
		return nil, false
	}

	result := make([]values.BalValue, len(segments))
	for i := range required {
		seg := entry.PathSegments[i]
		v, ok := coerceSegment(tc, seg.Ty, segments[i])
		if !ok {
			return nil, false
		}
		result[i] = v
	}
	for i := required; i < len(segments); i++ {
		v, ok := coerceSegment(tc, entry.RestSegmentTy, segments[i])
		if !ok {
			return nil, false
		}
		result[i] = v
	}
	return result, true
}

// decodeBalIdentifier converts a Ballerina identifier token text to its URL-path form:
// strips a leading quoted-identifier prefix (') and replaces backslash escapes (\X → X).
func decodeBalIdentifier(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '\'' {
		s = s[1:]
	}
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// coerceSegment coerces a URL path segment string to a typed value matching segTy.
func coerceSegment(tc semtypes.Context, segTy semtypes.SemType, s string) (values.BalValue, bool) {
	// Literal segment: must equal the expected string constant, after decoding any
	// Ballerina quoted-identifier prefix or backslash escapes from the stored literal.
	if shape := semtypes.SingleShape(segTy); shape.IsPresent() {
		if lit, ok := shape.Get().Value.(string); ok {
			if s != decodeBalIdentifier(lit) {
				return nil, false
			}
			// Return the raw stored literal so its singleton type matches the stored
			// entry type when LookupResourceMethod re-validates via resourcePathMatches.
			return lit, true
		}
	}
	// Parameter segment: coerce based on type.
	if semtypes.IsSubtype(tc, semtypes.INT, segTy) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, false
		}
		return n, true
	}
	if semtypes.IsSubtype(tc, semtypes.FLOAT, segTy) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, false
		}
		return f, true
	}
	if semtypes.IsSubtype(tc, semtypes.BOOLEAN, segTy) {
		switch s {
		case "true":
			return true, true
		case "false":
			return false, true
		}
		return nil, false
	}
	// STRING or any other type: accept as-is.
	return s, true
}

// buildRequest constructs a Ballerina Request object from HTTP request data.
// bodyStream is the raw request body; it is stored lazily in a requestBodyHolder
// so the body is only read from the network when a getPayload method is called.
// bodyBuf, when non-nil, is an already-read body; bodyStream must be nil in that case.
func buildRequest(tc semtypes.Context, method, rawPath, httpVersion string, headers map[string][]string, bodyStream io.ReadCloser, contentLength int64, rawQuery string, bodyBuf []byte) *values.Object {
	headersMap := newMappingValue(tc)
	for k, vals := range headers {
		items := make([]values.BalValue, len(vals))
		for i, v := range vals {
			items[i] = v
		}
		headersMap.Put(tc, strings.ToLower(k), newListValue(tc, items))
	}
	var holder *requestBodyHolder
	switch {
	case bodyBuf != nil:
		holder = &requestBodyHolder{buf: bodyBuf, contentLength: int64(len(bodyBuf))}
	case bodyStream != nil:
		holder = &requestBodyHolder{stream: bodyStream, contentLength: contentLength}
	default:
		holder = &requestBodyHolder{buf: []byte{}, contentLength: 0}
	}
	return values.NewObject(
		semtypes.OBJECT,
		map[string]values.BalValue{
			"rawPath":     rawPath,
			"method":      method,
			"httpVersion": httpVersion,
			"$headers":    headersMap,
			"$body":       holder,
			"$queryStr":   rawQuery,
		},
		map[string]string{
			"getTextPayload":     "ballerina/http:Request.getTextPayload",
			"getJsonPayload":     "ballerina/http:Request.getJsonPayload",
			"getBinaryPayload":   "ballerina/http:Request.getBinaryPayload",
			"getHeader":          "ballerina/http:Request.getHeader",
			"getHeaders":         "ballerina/http:Request.getHeaders",
			"hasHeader":          "ballerina/http:Request.hasHeader",
			"getQueryParams":     "ballerina/http:Request.getQueryParams",
			"getQueryParamValue": "ballerina/http:Request.getQueryParamValue",
		},
		nil,
	)
}

// writeErrorJSON writes a JSON error response in the standard Ballerina HTTP error format.
func writeErrorJSON(w http.ResponseWriter, r *http.Request, status int, message string) {
	type errorPayload struct {
		Timestamp string `json:"timestamp"`
		Status    int    `json:"status"`
		Reason    string `json:"reason"`
		Message   string `json:"message"`
		Path      string `json:"path"`
		Method    string `json:"method"`
	}
	payload := errorPayload{
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000") + "Z",
		Status:    status,
		Reason:    http.StatusText(status),
		Message:   message,
		Path:      r.URL.Path,
		Method:    r.Method,
	}
	body, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// writeResult writes a Ballerina resource method return value as an HTTP response.
func writeResult(_ semtypes.Context, w http.ResponseWriter, r *http.Request, result values.BalValue) {
	switch v := result.(type) {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	case *values.Error:
		writeErrorJSON(w, r, http.StatusInternalServerError, v.Message)
	case *values.Object:
		statusCodeVal, _ := v.Get("statusCode")
		statusCode := http.StatusOK
		if sc, ok := statusCodeVal.(int64); ok {
			statusCode = int(sc)
		}
		bodyVal, _ := v.Get("body")
		holder, _ := bodyVal.(*responseBodyHolder)

		// Emit headers from the response object, excluding hop-by-hop headers.
		// Forwarding hop-by-hop headers (e.g. Transfer-Encoding, Connection) from a
		// backend response to the downstream client violates RFC 7230 §6.1 and can
		// cause framing errors in HTTP/1.1 keep-alive connections.
		if hdrsVal, ok := v.Get("$headers"); ok {
			if hdrs, ok := hdrsVal.(*values.Map); ok {
				for _, k := range hdrs.Keys() {
					if _, skip := hopByHopHeaders[strings.ToLower(k)]; skip {
						continue
					}
					val, _ := hdrs.Get(k)
					list, ok := val.(*values.List)
					if !ok {
						continue
					}
					for i := range list.Len() {
						s, _ := list.Get(i).(string)
						if i == 0 {
							w.Header().Set(k, s)
						} else {
							w.Header().Add(k, s)
						}
					}
				}
			}
		}
		// WriteHeader must be called before writing the body; once body bytes
		// start flowing via writeStream, headers are already committed.
		w.WriteHeader(statusCode)
		if holder != nil {
			_ = holder.writeStream(w)
		}
	default:
		writeErrorJSON(w, r, http.StatusInternalServerError, "unexpected return type from resource method")
	}
}

// writeStream writes the body to w via io.Copy (streaming) or w.Write (buffered),
// then closes the stream. After this call the holder is exhausted.
func (h *responseBodyHolder) writeStream(w io.Writer) error {
	var (
		s   io.ReadCloser
		buf []byte
	)
	h.once.Do(func() {
		if h.stream != nil {
			s = h.stream
			h.stream = nil
			h.buf = []byte{}
		} else if len(h.buf) > 0 {
			buf = h.buf
			h.buf = []byte{}
		}
	})
	if s != nil {
		_, err := io.Copy(w, s)
		_ = s.Close()
		return err
	}
	if len(buf) > 0 {
		_, err := w.Write(buf)
		return err
	}
	// once was already fired by a prior materialize(); h.buf holds the materialized bytes.
	if len(h.buf) > 0 {
		_, err := w.Write(h.buf)
		return err
	}
	return nil
}
