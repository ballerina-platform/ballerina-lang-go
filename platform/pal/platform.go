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

import (
	"context"
	"io"
	"time"
)

type (
	Platform struct {
		IO   IO
		FS   FS
		Time Time
		HTTP HTTP
	}
	IO struct {
		Stdout func(p []byte) (n int, err error)
		Stderr func(p []byte) (n int, err error)
	}
	FS struct {
		ReadFile func(path string) ([]byte, error)
		WriteFile     func(path string, data []byte) error
		AppendFile    func(path string, data []byte) error
	}
	Time struct {
		Now          func() time.Time
		MonotonicNow func() time.Duration
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
	// PoolConfig carries connection pool settings derived from Ballerina's http:PoolConfiguration.
	// Zero values trigger platform-chosen defaults that mirror jBallerina's pool config.
	PoolConfig struct {
		// MaxIdleConnsPerHost is the per-host idle connection pool size.
		// 0 → default (100); matches jBallerina maxIdleConnections=100.
		MaxIdleConnsPerHost int
		// MaxIdleConns is the global idle pool size across all hosts.
		// 0 → default (512).
		MaxIdleConns int
		// MaxConnsPerHost is the maximum total (active+idle) connections per host.
		// 0 = unlimited; matches jBallerina maxActiveConnections=-1.
		MaxConnsPerHost int
		// IdleConnTimeout is how long idle connections are kept in the pool before eviction.
		// 0 → default (300s); matches jBallerina minEvictableIdleTime=300.
		IdleConnTimeout time.Duration
		// DialTimeout is the maximum time to establish a new TCP connection.
		// 0 → default (15s); matches jBallerina socketConfig.connectTimeOut=15.
		DialTimeout time.Duration
		// ResponseHeaderTimeout limits the time waiting for the first response byte.
		// 0 = no limit (disabled by default; opt-in via poolConfig).
		ResponseHeaderTimeout time.Duration
		// WriteBufferSize and ReadBufferSize size the per-connection user-space I/O buffers.
		// 0 → default (32 KB each).
		WriteBufferSize int
		ReadBufferSize  int
		// DisableCompression prevents the transport from injecting Accept-Encoding: gzip.
		DisableCompression bool
	}
	// ClientConfig bundles all static options for a new HTTP client instance.
	ClientConfig struct {
		Timeout         time.Duration
		FollowRedirects FollowRedirects
		HTTPVersion     string // "1.1" or "2.0"; defaults to "2.0"
		TLS             TLSConfig
		Pool            PoolConfig
	}
	// HTTPClient is an opaque handle to an HTTP client created by the platform.
	// It is created once per Ballerina http:Client init and reused across requests.
	HTTPClient interface {
		Execute(ctx context.Context, method, url string, body io.Reader, contentLength int64, contentType string, reqHeaders map[string][]string) (statusCode int, respHeaders map[string][]string, respBody io.ReadCloser, err error)
	}
)
