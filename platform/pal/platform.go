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

// Package pal provides the Platform Adaptation Layer (PAL).
//
// PAL abstracts away interactions with the underlying platform such that the
// runtime can be agnostic toward the underlying platform. All library functions
// that interact with the underlying platform (e.g. io, http) should use PAL.
// Each supported platform (e.g. native-cli, web-editor) should provide an
// implementation of PAL to the runtime.
package pal

import "time"

type (
	Platform struct {
		IO   IO
		FS   FS
		HTTP HTTP
	}
	IO struct {
		Stdout func(p []byte) (n int, err error)
		Stderr func(p []byte) (n int, err error)
	}
	FS struct {
		ReadFile func(path string) ([]byte, error)
	}
	HTTP struct {
		NewClient func(cfg ClientConfig) HTTPClient
	}
)

// HTTP
type (
	// TLSConfig carries TLS settings derived from Ballerina's secureSocket config.
	TLSConfig struct {
		InsecureSkipVerify    bool          // secureSocket.enable=false OR verifyHostName=false
		CACertPEM             []byte        // secureSocket.cert (string PEM file path) → file contents
		ClientCertPEM         []byte        // secureSocket.key.certFile → file contents
		ClientKeyPEM          []byte        // secureSocket.key.keyFile  → file contents
		ServerName            string        // secureSocket.serverName → tls.Config.ServerName (SNI)
		CipherSuiteNames      []string      // secureSocket.ciphers → IANA names; platform resolves IDs
		MinVersion            uint16        // secureSocket.protocol.versions min → tls.Config.MinVersion
		MaxVersion            uint16        // secureSocket.protocol.versions max → tls.Config.MaxVersion
		HandshakeTimeout      time.Duration // secureSocket.handshakeTimeout → transport.TLSHandshakeTimeout
		DisableSessionTickets bool          // secureSocket.shareSession=false → tls.Config.SessionTicketsDisabled
	}
	// FollowRedirects controls HTTP redirect behaviour, matching Ballerina's http:FollowRedirects.
	FollowRedirects struct {
		Enabled          bool // default false — no redirects by default (Ballerina spec)
		MaxCount         int  // 0 uses Ballerina default of 5; only used when Enabled=true
		AllowAuthHeaders bool // if true, forward Authorization/Proxy-Authorization on redirect
	}
	// ClientConfig bundles all static options for a new HTTP client instance.
	ClientConfig struct {
		Timeout         time.Duration
		FollowRedirects FollowRedirects
		HTTPVersion     string // "1.1" or "2.0"; defaults to "2.0"
		TLS             TLSConfig
	}
	// HTTPClient is an opaque handle to an HTTP client created by the platform.
	// It is created once per Ballerina http:Client init and reused across requests.
	HTTPClient interface {
		Execute(method, url string, body []byte, contentType string, reqHeaders map[string][]string) (statusCode int, respHeaders map[string][]string, respBody []byte, err error)
	}
)
