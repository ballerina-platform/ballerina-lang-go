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
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"reflect"
	"sort"
	"testing"
	"time"

	"ballerina-lang-go/decimal"
)

// ---------------------------------------------------------------------------
// removeHopByHopHeaders
// ---------------------------------------------------------------------------

func TestRemoveHopByHopHeaders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    map[string][]string
		wantGone []string
		wantKept []string
		desc     string
	}{
		{
			desc: "removes standard hop-by-hop headers",
			input: map[string][]string{
				"connection":        {"keep-alive"},
				"transfer-encoding": {"chunked"},
				"upgrade":           {"websocket"},
				"content-type":      {"application/json"},
				"authorization":     {"Bearer token"},
			},
			wantGone: []string{"connection", "transfer-encoding", "upgrade"},
			wantKept: []string{"content-type", "authorization"},
		},
		{
			desc: "honours Connection header token list",
			input: map[string][]string{
				"connection":    {"X-Custom-Hop, X-Another-Hop"},
				"x-custom-hop":  {"value1"},
				"x-another-hop": {"value2"},
				"x-safe-header": {"safe"},
			},
			wantGone: []string{"connection", "x-custom-hop", "x-another-hop"},
			wantKept: []string{"x-safe-header"},
		},
		{
			desc: "leaves non-hop-by-hop headers untouched",
			input: map[string][]string{
				"content-type":   {"text/plain"},
				"content-length": {"42"},
				"x-request-id":   {"abc123"},
			},
			wantGone: []string{},
			wantKept: []string{"content-type", "content-length", "x-request-id"},
		},
		{
			desc:     "empty map is a no-op",
			input:    map[string][]string{},
			wantGone: []string{},
			wantKept: []string{},
		},
		{
			desc: "removes keep-alive and te",
			input: map[string][]string{
				"keep-alive": {"timeout=5"},
				"te":         {"trailers"},
				"accept":     {"*/*"},
			},
			wantGone: []string{"keep-alive", "te"},
			wantKept: []string{"accept"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			removeHopByHopHeaders(tc.input)
			for _, k := range tc.wantGone {
				if _, ok := tc.input[k]; ok {
					t.Errorf("header %q should have been removed but is still present", k)
				}
			}
			for _, k := range tc.wantKept {
				if _, ok := tc.input[k]; !ok {
					t.Errorf("header %q should be present but was removed", k)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// splitOutsideQuotes
// ---------------------------------------------------------------------------

func TestSplitOutsideQuotes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		sep   byte
		want  []string
		desc  string
	}{
		{"a,b,c", ',', []string{"a", "b", "c"}, "basic comma split"},
		{"a", ',', []string{"a"}, "single element"},
		{"", ',', []string{""}, "empty string"},
		{`a,"b,c",d`, ',', []string{"a", `"b,c"`, "d"}, "comma inside quotes not split"},
		{`"a,b"`, ',', []string{`"a,b"`}, "fully quoted"},
		{`a,"b\"c",d`, ',', []string{"a", `"b\"c"`, "d"}, "escaped quote inside quotes"},
		{"a;b;c", ';', []string{"a", "b", "c"}, "semicolon separator"},
		{`"a;b";c`, ';', []string{`"a;b"`, "c"}, "semicolon inside quotes"},
		{"a,,b", ',', []string{"a", "", "b"}, "adjacent separators"},
		{`text/html; charset="utf-8"`, ';', []string{"text/html", ` charset="utf-8"`}, "content-type style"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := splitOutsideQuotes(tc.input, tc.sep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("splitOutsideQuotes(%q, %q) = %v, want %v",
					tc.input, tc.sep, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// applyCompressionHeaders
// ---------------------------------------------------------------------------

func TestApplyCompressionHeaders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		mode      string
		headers   map[string][]string
		wantHasAE bool
		wantAEVal string
		desc      string
	}{
		{
			mode:      "ALWAYS",
			headers:   map[string][]string{},
			wantHasAE: true,
			wantAEVal: "deflate, gzip",
			desc:      "ALWAYS adds Accept-Encoding to empty headers",
		},
		{
			mode:      "ALWAYS",
			headers:   map[string][]string{"content-type": {"application/json"}},
			wantHasAE: true,
			wantAEVal: "deflate, gzip",
			desc:      "ALWAYS adds Accept-Encoding when absent",
		},
		{
			mode:      "ALWAYS",
			headers:   map[string][]string{"accept-encoding": {"br"}},
			wantHasAE: true,
			wantAEVal: "br",
			desc:      "ALWAYS does not override existing Accept-Encoding",
		},
		{
			mode:      "NEVER",
			headers:   map[string][]string{"accept-encoding": {"gzip"}},
			wantHasAE: false,
			desc:      "NEVER removes Accept-Encoding",
		},
		{
			mode:      "NEVER",
			headers:   map[string][]string{"content-type": {"text/plain"}},
			wantHasAE: false,
			desc:      "NEVER does nothing when Accept-Encoding absent",
		},
		{
			mode:      "AUTO",
			headers:   map[string][]string{"accept-encoding": {"gzip"}},
			wantHasAE: true,
			wantAEVal: "gzip",
			desc:      "AUTO leaves headers untouched",
		},
		{
			mode:      "AUTO",
			headers:   map[string][]string{},
			wantHasAE: false,
			desc:      "AUTO does not add Accept-Encoding",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			result := applyCompressionHeaders(tc.mode, tc.headers)
			// Find Accept-Encoding case-insensitively
			var aeVals []string
			for k, v := range result {
				if k == "accept-encoding" || k == "Accept-Encoding" {
					aeVals = v
					break
				}
			}
			hasAE := len(aeVals) > 0
			if hasAE != tc.wantHasAE {
				t.Errorf("mode=%q: Accept-Encoding present = %v, want %v", tc.mode, hasAE, tc.wantHasAE)
			}
			if tc.wantHasAE && tc.wantAEVal != "" && len(aeVals) > 0 && aeVals[0] != tc.wantAEVal {
				t.Errorf("mode=%q: Accept-Encoding = %q, want %q", tc.mode, aeVals[0], tc.wantAEVal)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// decompressResponseBody
// ---------------------------------------------------------------------------

func gzipCompress(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

func deflateCompress(data []byte) []byte {
	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.DefaultCompression)
	_, _ = w.Write(data)
	_ = w.Close()
	return buf.Bytes()
}

func TestDecompressResponseBody(t *testing.T) {
	t.Parallel()
	original := []byte("Hello, compressed world!")

	tests := []struct {
		desc    string
		headers map[string][]string
		body    []byte
		want    []byte
	}{
		{
			desc:    "gzip decompressed and header removed",
			headers: map[string][]string{"content-encoding": {"gzip"}},
			body:    gzipCompress(original),
			want:    original,
		},
		{
			desc:    "deflate decompressed and header removed",
			headers: map[string][]string{"content-encoding": {"deflate"}},
			body:    deflateCompress(original),
			want:    original,
		},
		{
			desc:    "no content-encoding passes through unchanged",
			headers: map[string][]string{"content-type": {"text/plain"}},
			body:    original,
			want:    original,
		},
		{
			desc:    "unknown encoding passes through unchanged",
			headers: map[string][]string{"content-encoding": {"br"}},
			body:    original,
			want:    original,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			bodyRC := io.NopCloser(bytes.NewReader(tc.body))
			result := decompressResponseBody(tc.headers, bodyRC)
			got, err := io.ReadAll(result)
			if err != nil {
				t.Fatalf("reading decompressed body: %v", err)
			}
			_ = result.Close()
			if !bytes.Equal(got, tc.want) {
				t.Errorf("decompressResponseBody: got %q, want %q", got, tc.want)
			}
			// Verify header was removed after decompression for gzip/deflate
			for k := range tc.headers {
				if k == "content-encoding" {
					if tc.desc != "no content-encoding passes through unchanged" &&
						tc.desc != "unknown encoding passes through unchanged" {
						t.Errorf("content-encoding header should have been removed but is still present")
					}
				}
			}
		})
	}
}

func TestDecompressResponseBody_RemovesHeader(t *testing.T) {
	t.Parallel()
	original := []byte("test data")
	headers := map[string][]string{"content-encoding": {"gzip"}, "content-type": {"text/plain"}}
	body := gzipCompress(original)

	bodyRC := io.NopCloser(bytes.NewReader(body))
	result := decompressResponseBody(headers, bodyRC)
	_, _ = io.ReadAll(result)
	_ = result.Close()

	if _, ok := headers["content-encoding"]; ok {
		t.Error("content-encoding should have been removed from headers after gzip decompression")
	}
	if _, ok := headers["content-type"]; !ok {
		t.Error("content-type should still be present")
	}
}

// ---------------------------------------------------------------------------
// decimalToDuration
// ---------------------------------------------------------------------------

func TestDecimalToDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input *decimal.Decimal
		want  time.Duration
		desc  string
	}{
		{decimal.FromInt64(1), time.Second, "1 second"},
		{decimal.FromInt64(0), 0, "zero"},
		{decimal.FromInt64(60), 60 * time.Second, "60 seconds"},
	}
	// Build fractional values from string
	if d, err := decimal.FromString("1.5"); err == nil {
		tests = append(tests, struct {
			input *decimal.Decimal
			want  time.Duration
			desc  string
		}{d, 1500 * time.Millisecond, "1.5 seconds"})
	}
	if d, err := decimal.FromString("0.001"); err == nil {
		tests = append(tests, struct {
			input *decimal.Decimal
			want  time.Duration
			desc  string
		}{d, time.Millisecond, "1ms"})
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := decimalToDuration(tc.input)
			if got != tc.want {
				t.Errorf("decimalToDuration(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// eagerBufferResponse
// ---------------------------------------------------------------------------

func TestEagerBufferResponse_NilStream(t *testing.T) {
	t.Parallel()
	holder := eagerBufferResponse(map[string][]string{}, nil)
	data, err := holder.materialize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty buf for nil stream, got %d bytes", len(data))
	}
}

func TestEagerBufferResponse_SmallBody(t *testing.T) {
	t.Parallel()
	body := []byte("small body")
	headers := map[string][]string{"content-length": {"10"}}
	holder := eagerBufferResponse(headers, io.NopCloser(bytes.NewReader(body)))
	data, err := holder.materialize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, body) {
		t.Errorf("got %q, want %q", data, body)
	}
}

func TestEagerBufferResponse_LargeBody(t *testing.T) {
	t.Parallel()
	body := make([]byte, eagerBufferThreshold+1)
	for i := range body {
		body[i] = byte(i % 256)
	}
	headers := map[string][]string{}
	holder := eagerBufferResponse(headers, io.NopCloser(bytes.NewReader(body)))
	data, err := holder.materialize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, body) {
		t.Errorf("large body not correctly buffered")
	}
}

// ---------------------------------------------------------------------------
// requestBodyHolder
// ---------------------------------------------------------------------------

func TestRequestBodyHolder_MaterializeTwice(t *testing.T) {
	t.Parallel()
	body := []byte("hello world")
	h := &requestBodyHolder{stream: io.NopCloser(bytes.NewReader(body))}
	got1 := h.materialize()
	got2 := h.materialize()
	if !bytes.Equal(got1, body) {
		t.Errorf("first materialize: got %q, want %q", got1, body)
	}
	if !bytes.Equal(got2, body) {
		t.Errorf("second materialize: got %q, want %q", got2, body)
	}
}

func TestRequestBodyHolder_TakeStream(t *testing.T) {
	t.Parallel()
	body := []byte("stream data")
	h := &requestBodyHolder{stream: io.NopCloser(bytes.NewReader(body))}
	stream := h.takeStream()
	if stream == nil {
		t.Fatal("takeStream should return stream on first call")
	}
	data, _ := io.ReadAll(stream)
	_ = stream.Close()
	if !bytes.Equal(data, body) {
		t.Errorf("takeStream data: got %q, want %q", data, body)
	}
	// Second call should return nil (already consumed)
	if second := h.takeStream(); second != nil {
		t.Error("takeStream should return nil after stream was taken")
	}
}

func TestRequestBodyHolder_TakeStreamAfterMaterialize(t *testing.T) {
	t.Parallel()
	body := []byte("data")
	h := &requestBodyHolder{stream: io.NopCloser(bytes.NewReader(body))}
	_ = h.materialize() // consume stream
	stream := h.takeStream()
	if stream != nil {
		t.Error("takeStream should return nil when already materialized")
	}
}

// ---------------------------------------------------------------------------
// responseBodyHolder
// ---------------------------------------------------------------------------

func TestResponseBodyHolder_MaterializeNilStream(t *testing.T) {
	t.Parallel()
	holder := newResponseBodyHolder(nil)
	data, err := holder.materialize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("nil stream should produce empty body, got %d bytes", len(data))
	}
}

func TestResponseBodyHolder_MaterializeTwice(t *testing.T) {
	t.Parallel()
	body := []byte("response body")
	holder := newResponseBodyHolder(io.NopCloser(bytes.NewReader(body)))
	got1, err1 := holder.materialize()
	got2, err2 := holder.materialize()
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v, %v", err1, err2)
	}
	if !bytes.Equal(got1, body) || !bytes.Equal(got2, body) {
		t.Errorf("expected both materializations to return same data")
	}
}

// ---------------------------------------------------------------------------
// balToGoJSON
// ---------------------------------------------------------------------------

func TestBalToGoJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input interface{}
		want  interface{}
		desc  string
	}{
		{nil, nil, "nil"},
		{int64(42), int64(42), "int64 preserved"},
		{int64(-1), int64(-1), "negative int64 preserved"},
		{float64(3.14), float64(3.14), "float64"},
		{true, true, "bool true"},
		{false, false, "bool false"},
		{"hello", "hello", "string"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := balToGoJSON(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("balToGoJSON(%v) = %v (%T), want %v (%T)",
					tc.input, got, got, tc.want, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// goCtxOrBackground
// ---------------------------------------------------------------------------

func TestGoCtxOrBackground(t *testing.T) {
	t.Parallel()
	ctx := goCtxOrBackground(nil)
	if ctx == nil {
		t.Error("goCtxOrBackground(nil) should return non-nil context")
	}
}

// ---------------------------------------------------------------------------
// sortedKeys helper for test determinism
// ---------------------------------------------------------------------------

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
