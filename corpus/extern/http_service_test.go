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
	goruntime "runtime"
	"testing"

	"ballerina-lang-go/platform/palnative"
)

// skipIfNoLoopback skips on platforms without loopback TCP (js/wasm). Unlike
// skipIfNoNetwork these service tests only need localhost, not the internet.
func skipIfNoLoopback(t *testing.T) {
	t.Helper()
	if goruntime.GOOS == "js" {
		t.Skip("skipping loopback-dependent test on js/wasm")
	}
}

// TestHttpServiceBasic starts a Ballerina http service on a real port and
// drives it from testMain via a real http:Client, exercising the full
// listener → dispatch → resource → response path.
func TestHttpServiceBasic(t *testing.T) {
	skipIfNoLoopback(t)
	runExtern(t, fileCase("http-svc-basic-v"), newHTTPPal(palnative.NewHTTPClient), nil)
}

// TestHttpServicePathParam exercises a typed (int) path parameter.
func TestHttpServicePathParam(t *testing.T) {
	skipIfNoLoopback(t)
	runExtern(t, fileCase("http-svc-path-param-v"), newHTTPPal(palnative.NewHTTPClient), nil)
}

// TestHttpServiceRequest exercises Request injection and JSON body round-trip
// through a POST resource.
func TestHttpServiceRequest(t *testing.T) {
	skipIfNoLoopback(t)
	runExtern(t, fileCase("http-svc-request-v"), newHTTPPal(palnative.NewHTTPClient), nil)
}

// TestHttpServiceRouting exercises 200 / 404 (unknown path) / 405 (wrong
// method) dispatch outcomes.
func TestHttpServiceRouting(t *testing.T) {
	skipIfNoLoopback(t)
	runExtern(t, fileCase("http-svc-routing-v"), newHTTPPal(palnative.NewHTTPClient), nil)
}

// TestHttpServiceTLS exercises a TLS listener: the server loads its cert/key
// from disk (realFS) and the client connects over https with verification
// disabled for the self-signed cert.
func TestHttpServiceTLS(t *testing.T) {
	skipIfNoLoopback(t)
	runExtern(t, fileCase("http-svc-tls-v"), newHTTPPal(palnative.NewHTTPClient).withRealFS(), nil)
}
