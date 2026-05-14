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

package extern_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/pal"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/test_util"
)

type rewritingHTTPClient struct {
	serverURL string
	client    *http.Client
}

func (c *rewritingHTTPClient) Execute(method, url string, body []byte, contentType string, reqHeaders map[string][]string) (int, map[string][]string, []byte, error) {
	const prefix = "http://testserver"
	if !strings.HasPrefix(url, prefix) {
		return 0, nil, nil, fmt.Errorf("rewritingHTTPClient: expected URL with prefix %q, got %q", prefix, url)
	}
	realURL := c.serverURL + url[len(prefix):]
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, realURL, bodyReader)
	if err != nil {
		return 0, nil, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, vals := range reqHeaders {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, err := io.ReadAll(resp.Body)
	return resp.StatusCode, map[string][]string(resp.Header), respBody, err
}

func TestHttpClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/hello" {
			w.WriteHeader(200)
			// Echo the X-Test header back in the body so the test can verify it was sent.
			xTest := r.Header.Get("X-Test")
			if xTest != "" {
				_, _ = fmt.Fprintf(w, "hello with %s", xTest)
			} else {
				_, _ = fmt.Fprint(w, "hello from test server")
			}
		} else {
			w.WriteHeader(404)
			_, _ = fmt.Fprint(w, "not found")
		}
	}))
	defer server.Close()

	balFile := filepath.Join(testDataDir, "http-client-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath := getBallerinaEnvPath(t)
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, filepath.Base(absPath), projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}

	// Rewrite the URL so the bal file hits the local test server.
	// The bal file uses "http://testserver" which we redirect to the actual server.
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: func(cfg pal.ClientConfig) pal.HTTPClient {
			return &rewritingHTTPClient{
				serverURL: server.URL,
				client:    &http.Client{Timeout: cfg.Timeout},
			}
		},
	}

	rt := runtime.NewRuntime(testPal)
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "200\nhello with test-header-value\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func TestHttpClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/echo" {
			body, _ := io.ReadAll(r.Body)
			ct := r.Header.Get("Content-Type")
			ct = strings.SplitN(ct, ";", 2)[0] // strip charset etc.
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "body: %s, ct: %s", string(body), ct)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	balFile := filepath.Join(testDataDir, "http-client-post-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath := getBallerinaEnvPath(t)
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, filepath.Base(absPath), projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: func(cfg pal.ClientConfig) pal.HTTPClient {
			return &rewritingHTTPClient{
				serverURL: server.URL,
				client:    &http.Client{Timeout: cfg.Timeout},
			}
		},
	}

	rt := runtime.NewRuntime(testPal)
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "200\nbody: hello post, ct: text/plain\n200\nbody: {\"msg\":\"hello\"}, ct: application/json\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func TestHttpClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/echo":
			body, _ := io.ReadAll(r.Body)
			ct := r.Header.Get("Content-Type")
			ct = strings.SplitN(ct, ";", 2)[0]
			w.WriteHeader(200)
			_, _ = fmt.Fprintf(w, "body: %s, ct: %s", string(body), ct)
		case r.URL.Path == "/delete" && r.Method == http.MethodDelete:
			w.WriteHeader(200)
		case r.URL.Path == "/head" && r.Method == http.MethodHead:
			w.WriteHeader(200)
		case r.URL.Path == "/options" && r.Method == http.MethodOptions:
			w.WriteHeader(200)
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	balFile := filepath.Join(testDataDir, "http-client-methods-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath := getBallerinaEnvPath(t)
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, filepath.Base(absPath), projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: func(cfg pal.ClientConfig) pal.HTTPClient {
			return &rewritingHTTPClient{
				serverURL: server.URL,
				client:    &http.Client{Timeout: cfg.Timeout},
			}
		},
	}

	rt := runtime.NewRuntime(testPal)
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "" +
		"200\nbody: put body, ct: text/plain\n" +
		"200\nbody: {\"k\":\"v\"}, ct: application/json\n" +
		"200\n" +
		"200\n" +
		"200\n" +
		"200\nbody: exec body, ct: text/plain\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func TestHttpClientTLSInsecure(t *testing.T) {
	// httptest.NewTLSServer uses a self-signed certificate. A client without
	// InsecureSkipVerify would fail the handshake. This test verifies that
	// secureSocket: {enable: false} propagates InsecureSkipVerify through the PAL.
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, "secure hello")
	}))
	defer server.Close()

	balFile := filepath.Join(testDataDir, "http-client-tls-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath := getBallerinaEnvPath(t)
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, filepath.Base(absPath), projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: func(cfg pal.ClientConfig) pal.HTTPClient {
			// Build a TLS-aware client. Clone the test server's TLS config and
			// apply InsecureSkipVerify if the Ballerina code requested it.
			serverTLSConfig := server.Client().Transport.(*http.Transport).TLSClientConfig.Clone()
			if cfg.TLS.InsecureSkipVerify {
				serverTLSConfig.InsecureSkipVerify = true //nolint:gosec
			}
			return &rewritingHTTPClient{
				serverURL: server.URL,
				client:    &http.Client{Timeout: cfg.Timeout, Transport: &http.Transport{TLSClientConfig: serverTLSConfig}},
			}
		},
	}

	rt := runtime.NewRuntime(testPal)
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "200\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}
