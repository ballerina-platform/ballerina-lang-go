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

package palnative

import (
	"crypto/tls"
	"testing"
	"time"

	"ballerina-lang-go/platform/pal"
)

// ---------------------------------------------------------------------------
// tlsMatchCN
// ---------------------------------------------------------------------------

func TestTlsMatchCN(t *testing.T) {
	tests := []struct {
		pattern string
		host    string
		match   bool
		desc    string
	}{
		{"example.com", "example.com", true, "exact match"},
		{"example.com", "other.com", false, "different host"},
		{"example.com", "sub.example.com", false, "no wildcard on exact pattern"},
		{"*.example.com", "sub.example.com", true, "single-label wildcard match"},
		{"*.example.com", "example.com", false, "wildcard does not match apex"},
		{"*.example.com", "deep.sub.example.com", false, "wildcard matches only one label"},
		{"", "example.com", false, "empty pattern"},
		{"example.com", "", false, "empty host"},
		{"EXAMPLE.COM", "example.com", true, "case-insensitive exact"},
		{"*.EXAMPLE.COM", "sub.example.com", true, "case-insensitive wildcard"},
		{"example.com.", "example.com", true, "trailing dot stripped in pattern"},
		{"example.com", "example.com.", true, "trailing dot stripped in host"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tlsMatchCN(tc.pattern, tc.host)
			if tc.match && err != nil {
				t.Errorf("tlsMatchCN(%q, %q): expected match, got error: %v", tc.pattern, tc.host, err)
			}
			if !tc.match && err == nil {
				t.Errorf("tlsMatchCN(%q, %q): expected mismatch error, got nil", tc.pattern, tc.host)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// resolveCipherSuites
// ---------------------------------------------------------------------------

func TestResolveCipherSuites_Empty(t *testing.T) {
	result := resolveCipherSuites([]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %v", result)
	}
}

func TestResolveCipherSuites_UnknownName(t *testing.T) {
	result := resolveCipherSuites([]string{"NOT_A_REAL_CIPHER"})
	if len(result) != 0 {
		t.Errorf("expected empty result for unknown cipher, got %v", result)
	}
}

func TestResolveCipherSuites_KnownSecureName(t *testing.T) {
	suites := tls.CipherSuites()
	if len(suites) == 0 {
		t.Skip("no cipher suites available on this platform")
	}
	name := suites[0].Name
	result := resolveCipherSuites([]string{name})
	if len(result) != 1 {
		t.Errorf("expected 1 resolved cipher for %q, got %d", name, len(result))
	}
	if result[0] != suites[0].ID {
		t.Errorf("expected cipher ID %d for %q, got %d", suites[0].ID, name, result[0])
	}
}

func TestResolveCipherSuites_MixedKnownUnknown(t *testing.T) {
	suites := tls.CipherSuites()
	if len(suites) == 0 {
		t.Skip("no cipher suites available on this platform")
	}
	names := []string{suites[0].Name, "UNKNOWN_CIPHER", suites[len(suites)-1].Name}
	result := resolveCipherSuites(names)
	if len(result) != 2 {
		t.Errorf("expected 2 resolved ciphers for mixed input, got %d", len(result))
	}
}

func TestResolveCipherSuites_InsecureName(t *testing.T) {
	insecure := tls.InsecureCipherSuites()
	if len(insecure) == 0 {
		t.Skip("no insecure cipher suites available on this platform")
	}
	name := insecure[0].Name
	result := resolveCipherSuites([]string{name})
	if len(result) != 1 {
		t.Errorf("expected 1 resolved insecure cipher for %q, got %d", name, len(result))
	}
}

// ---------------------------------------------------------------------------
// NewHTTPClient
// ---------------------------------------------------------------------------

func TestNewHTTPClient_DefaultConfig(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{})
	if client == nil {
		t.Fatal("expected non-nil client for default config")
	}
}

func TestNewHTTPClient_WithTimeout(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{Timeout: 5 * time.Second})
	if client == nil {
		t.Fatal("expected non-nil client with timeout config")
	}
}

func TestNewHTTPClient_HTTP2(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{HTTPVersion: "2.0"})
	if client == nil {
		t.Fatal("expected non-nil client with HTTP/2 config")
	}
}

func TestNewHTTPClient_RedirectsDisabled(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{
		FollowRedirects: pal.FollowRedirects{Enabled: false},
	})
	if client == nil {
		t.Fatal("expected non-nil client with redirects disabled")
	}
}

func TestNewHTTPClient_RedirectsEnabled(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{
		FollowRedirects: pal.FollowRedirects{
			Enabled:          true,
			MaxCount:         3,
			AllowAuthHeaders: true,
		},
	})
	if client == nil {
		t.Fatal("expected non-nil client with redirects enabled")
	}
}

func TestNewHTTPClient_TLSVersions(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{
		TLS: pal.TLSConfig{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		},
	})
	if client == nil {
		t.Fatal("expected non-nil client with TLS version range")
	}
}

func TestNewHTTPClient_ValidCipherSuites(t *testing.T) {
	suites := tls.CipherSuites()
	if len(suites) == 0 {
		t.Skip("no cipher suites available")
	}
	client := NewHTTPClient(pal.ClientConfig{
		TLS: pal.TLSConfig{
			CipherSuiteNames: []string{suites[0].Name},
		},
	})
	if client == nil {
		t.Fatal("expected non-nil client with valid cipher suites")
	}
}

func TestNewHTTPClient_InvalidCipherSuites(t *testing.T) {
	// All unknown names — falls through to the warning path, but should not panic.
	client := NewHTTPClient(pal.ClientConfig{
		TLS: pal.TLSConfig{
			CipherSuiteNames: []string{"NOT_A_REAL_CIPHER"},
		},
	})
	if client == nil {
		t.Fatal("expected non-nil client even when cipher suite names are unresolvable")
	}
}

func TestNewHTTPClient_InsecureSkipVerify(t *testing.T) {
	client := NewHTTPClient(pal.ClientConfig{
		TLS: pal.TLSConfig{InsecureSkipVerify: true},
	})
	if client == nil {
		t.Fatal("expected non-nil client with InsecureSkipVerify")
	}
}
