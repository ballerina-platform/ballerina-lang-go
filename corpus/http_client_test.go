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

package corpus

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/platform/palnative"
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

	balFile := filepath.Join(externTestDataDir, "http-client-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
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

	balFile := filepath.Join(externTestDataDir, "http-client-post-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
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

	balFile := filepath.Join(externTestDataDir, "http-client-methods-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
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

	balFile := filepath.Join(externTestDataDir, "http-client-tls-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "200\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

// TestHttpClientPublicGet exercises palnative.NewHTTPClient against a real
// public endpoint, ensuring the full Ballerina → PAL → palnative path is
// covered. Skipped when CORPUS_SKIP_NETWORK=1 or no network is available.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientPublicGet(t *testing.T) {
	skipIfNoNetwork(t)

	balFile := filepath.Join(externTestDataDir, "http-client-public-get-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: palnative.NewHTTPClient,
	}

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	if strings.TrimSpace(stdoutBuf.String()) != "200" {
		t.Errorf("expected status 200, got: %q", stdoutBuf.String())
	}
}

// TestHttpClientRedirect exercises redirect-following through palnative.NewHTTPClient.
func TestHttpClientRedirect(t *testing.T) {
	skipIfNoNetwork(t)

	balFile := filepath.Join(externTestDataDir, "http-client-redirect-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{
		NewClient: palnative.NewHTTPClient,
	}

	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	if strings.TrimSpace(stdoutBuf.String()) != "200" {
		t.Errorf("expected final status 200 after redirect, got: %q", stdoutBuf.String())
	}
}

// skipIfNoNetwork skips the test when CORPUS_SKIP_NETWORK is set or when
// running under WASM (js/wasm), which has no outbound TCP access.
func skipIfNoNetwork(t *testing.T) {
	t.Helper()
	if os.Getenv("CORPUS_SKIP_NETWORK") != "" || goruntime.GOOS == "js" {
		t.Skip("skipping network-dependent test")
	}
}

// runNetworkBal compiles and interprets a static .bal file using
// palnative.NewHTTPClient. Returns trimmed stdout. Callers must guard
// with skipIfNoNetwork before calling.
func runNetworkBal(t *testing.T, balFile string) string {
	t.Helper()
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatal("compilation had errors")
	}
	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{NewClient: palnative.NewHTTPClient}
	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}
	return strings.TrimSpace(stdoutBuf.String())
}

// TestHttpClientJson exercises Response.getJsonPayload against httpbin /json.
// Asserts that the returned value satisfies `is json` (map<json> semtype).
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientJson(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-json-v.bal"))
	if out != "true" {
		t.Errorf("expected getJsonPayload() result to be json, got: %q", out)
	}
}

// TestHttpClientText exercises Response.getTextPayload against httpbin /html.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientText(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-text-v.bal"))
	if out != "true" {
		t.Errorf("expected getTextPayload() to return non-empty string, got: %q", out)
	}
}

// TestHttpClientBinary exercises Response.getBinaryPayload against httpbin /bytes/16.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientBinary(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-binary-v.bal"))
	if out != "true" {
		t.Errorf("expected getBinaryPayload() to return non-empty byte[], got: %q", out)
	}
}

// TestHttpClientPublicMethods exercises POST, PUT, DELETE, and PATCH against
// dedicated httpbin endpoints that each return 200 for their verb.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientPublicMethods(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-public-methods-v.bal"))
	if out != "200\n200\n200\n200" {
		t.Errorf("expected four 200 status codes, got: %q", out)
	}
}

// TestHttpClientTimeout verifies that a 1-second timeout fires before
// httpbin /delay/5 responds, and that the resulting error propagates to
// Ballerina as an error value.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientTimeout(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-timeout-v.bal"))
	if out != "true" {
		t.Errorf("expected timeout to produce an error value, got: %q", out)
	}
}

// TestHttpClientConnectionError verifies that a DNS resolution failure for an
// unreachable host propagates back to Ballerina as an error value.
// TODO: Replace with a Ballerina HTTP service once server support lands.
func TestHttpClientConnectionError(t *testing.T) {
	skipIfNoNetwork(t)
	out := runNetworkBal(t, filepath.Join(externTestDataDir, "http-client-connection-error-v.bal"))
	if out != "true" {
		t.Errorf("expected connection error to produce an error value, got: %q", out)
	}
}

// TestHttpClientMTLS verifies mutual TLS by spinning up a local HTTPS server
// that requires a valid client certificate, then confirming the Ballerina
// secureSocket.key.certFile / keyFile config causes the client to send it.
// Uses locally generated certs so no external network is needed.
// Server cert verification is disabled (verifyHostName: false) because Go's
// TLS stack does not send SNI for IP-address targets — the mTLS focus is
// the client-certificate path, not server-cert verification.
func TestHttpClientMTLS(t *testing.T) {
	caCertPEM, serverCertPEM, serverKeyPEM, clientCertPEM, clientKeyPEM := generateTestCerts(t)

	serverTLSCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		t.Fatalf("creating server TLS cert: %v", err)
	}
	serverTLSCert.Leaf, err = x509.ParseCertificate(serverTLSCert.Certificate[0])
	if err != nil {
		t.Fatalf("parsing server TLS cert leaf: %v", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(caCertPEM)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		ClientCAs:    clientCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	server.StartTLS()
	defer server.Close()

	tmpDir := t.TempDir()
	clientCertFile := filepath.Join(tmpDir, "client.pem")
	clientKeyFile := filepath.Join(tmpDir, "client-key.pem")
	for _, pair := range []struct {
		path string
		data []byte
	}{
		{clientCertFile, clientCertPEM},
		{clientKeyFile, clientKeyPEM},
	} {
		if err := os.WriteFile(pair.path, pair.data, 0600); err != nil {
			t.Fatalf("writing %s: %v", pair.path, err)
		}
	}

	// verifyHostName: false skips server-cert verification (the IP-address
	// target has no SNI, causing cs.ServerName="" in the custom VerifyConnection
	// callback). The test focus is the client-cert path: the server still
	// validates the client cert against its CA pool and rejects without it.
	//
	// Use forward slashes in paths embedded in the Ballerina string literal so
	// that backslashes on Windows are not misinterpreted as escape sequences.
	certFileSlash := filepath.ToSlash(clientCertFile)
	keyFileSlash := filepath.ToSlash(clientKeyFile)
	balContent := fmt.Sprintf(`
import ballerina/http;
import ballerina/io;

public function main() returns error? {
    http:Client c = check new ("%s", {
        secureSocket: {
            verifyHostName: false,
            key: {certFile: "%s", keyFile: "%s"}
        }
    });
    http:Response r = check c->get("/");
    io:println(r.statusCode);
    return;
}
`, server.URL, certFileSlash, keyFileSlash)

	tmpBalFile := filepath.Join(tmpDir, "http-client-mtls-v.bal")
	if err := os.WriteFile(tmpBalFile, []byte(balContent), 0644); err != nil {
		t.Fatalf("writing bal file: %v", err)
	}

	absPath, err := filepath.Abs(tmpBalFile)
	if err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(filepath.Dir(absPath))
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	stdoutBuf := &bytes.Buffer{}
	testPal := test_util.TestPal(stdoutBuf, os.Stderr)
	testPal.HTTP = pal.HTTP{NewClient: palnative.NewHTTPClient}
	rt := runtime.NewRuntime(testPal, result.Project().Environment().TypeEnv())
	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	if strings.TrimSpace(stdoutBuf.String()) != "200" {
		t.Errorf("expected status 200 from mTLS server, got: %q", stdoutBuf.String())
	}
}

// generateTestCerts generates a self-signed CA, a server cert for 127.0.0.1,
// and a client cert. All leaves are signed by the CA. Returns PEM-encoded
// bytes ready for use with tls.X509KeyPair or file writes.
func generateTestCerts(t *testing.T) (caCertPEM, serverCertPEM, serverKeyPEM, clientCertPEM, clientKeyPEM []byte) {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	caSerial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	caTemplate := &x509.Certificate{
		SerialNumber:          caSerial,
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatal(err)
	}
	caCertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})

	newLeafCert := func(cn string, ips []net.IP, extUsage x509.ExtKeyUsage) (certPEM, keyPEM []byte) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
		tmpl := &x509.Certificate{
			SerialNumber: serial,
			Subject:      pkix.Name{CommonName: cn},
			NotBefore:    time.Now().Add(-time.Hour),
			NotAfter:     time.Now().Add(time.Hour),
			ExtKeyUsage:  []x509.ExtKeyUsage{extUsage},
			IPAddresses:  ips,
		}
		certDER, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
		if err != nil {
			t.Fatal(err)
		}
		keyDER, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}),
			pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	}

	serverCertPEM, serverKeyPEM = newLeafCert("127.0.0.1", []net.IP{net.ParseIP("127.0.0.1")}, x509.ExtKeyUsageServerAuth)
	clientCertPEM, clientKeyPEM = newLeafCert("client", nil, x509.ExtKeyUsageClientAuth)
	return
}
