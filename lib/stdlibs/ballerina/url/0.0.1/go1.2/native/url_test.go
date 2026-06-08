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
	"testing"

	"ballerina-lang-go/values"
)

func TestResolveEncoding(t *testing.T) {
	t.Parallel()
	tests := []struct {
		charset string
		wantNil bool
		desc    string
	}{
		{"UTF-8", false, "UTF-8 is identity encoding"},
		{"utf-8", false, "lowercase utf-8"},
		{"UTF8", false, "UTF8 alias"},
		{"ISO-8859-1", false, "ISO-8859-1"},
		{"iso-8859-1", false, "lowercase ISO-8859-1"},
		{"ISO8859-1", false, "ISO8859-1 alias"},
		{"ISO_8859_1", false, "ISO_8859_1 alias"},
		{"LATIN-1", false, "LATIN-1 alias"},
		{"LATIN1", false, "LATIN1 alias"},
		{"US-ASCII", false, "ASCII"},
		{"ASCII", false, "ASCII alias"},
		{"UTF-16", false, "UTF-16"},
		{"UTF-16BE", false, "UTF-16BE"},
		{"UTF-16LE", false, "UTF-16LE"},
		{"UNKNOWN-CHARSET", true, "unknown charset returns nil"},
		{"", true, "empty charset returns nil"},
		{"windows-1252", true, "unsupported charset returns nil"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := resolveEncoding(tc.charset)
			if tc.wantNil && got != nil {
				t.Errorf("resolveEncoding(%q) = %v, want nil", tc.charset, got)
			}
			if !tc.wantNil && got == nil {
				t.Errorf("resolveEncoding(%q) = nil, want non-nil", tc.charset)
			}
		})
	}
}

func TestIsUnreserved(t *testing.T) {
	t.Parallel()
	unreserved := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~"
	for _, c := range []byte(unreserved) {
		if !isUnreserved(c) {
			t.Errorf("isUnreserved(%q) = false, want true", c)
		}
	}
	reserved := " !@#$%^&*()+=[]{}|;:',/<>?`\"\\"
	for _, c := range []byte(reserved) {
		if isUnreserved(c) {
			t.Errorf("isUnreserved(%q) = true, want false", c)
		}
	}
}

func TestFromHex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c     byte
		want  byte
		valid bool
		desc  string
	}{
		{'0', 0, true, "digit 0"},
		{'9', 9, true, "digit 9"},
		{'a', 10, true, "lowercase a"},
		{'f', 15, true, "lowercase f"},
		{'A', 10, true, "uppercase A"},
		{'F', 15, true, "uppercase F"},
		{'g', 0, false, "invalid g"},
		{'G', 0, false, "invalid G"},
		{'z', 0, false, "invalid z"},
		{'/', 0, false, "slash not hex"},
		{':', 0, false, "colon not hex"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, ok := fromHex(tc.c)
			if ok != tc.valid {
				t.Errorf("fromHex(%q) valid = %v, want %v", tc.c, ok, tc.valid)
			}
			if tc.valid && got != tc.want {
				t.Errorf("fromHex(%q) = %d, want %d", tc.c, got, tc.want)
			}
		})
	}
}

func TestEncodeBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		raw  []byte
		want string
		desc string
	}{
		{[]byte("hello"), "hello", "unreserved chars pass through"},
		{[]byte("Hello World"), "Hello%20World", "space becomes %20"},
		{[]byte("abc~-_."), "abc~-_.", "all unreserved chars"},
		{[]byte{0xFF}, "%FF", "high byte encoded"},
		{[]byte{0x00}, "%00", "null byte encoded"},
		{[]byte("a+b"), "a%2Bb", "plus sign encoded"},
		{[]byte("a=b&c=d"), "a%3Db%26c%3Dd", "equals and ampersand encoded"},
		{[]byte{}, "", "empty input"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := encodeBytes(tc.raw)
			if got != tc.want {
				t.Errorf("encodeBytes(%v) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestPercentDecodeToBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    []byte
		wantErr bool
		desc    string
	}{
		{"hello", []byte("hello"), false, "no encoding"},
		{"hello+world", []byte("hello world"), false, "plus becomes space"},
		{"hello%20world", []byte("hello world"), false, "%20 is space"},
		{"%41%42%43", []byte("ABC"), false, "hex uppercase"},
		{"%61%62%63", []byte("abc"), false, "hex lowercase"},
		{"a%2Bb", []byte("a+b"), false, "encoded plus"},
		{"%", []byte("%"), false, "lone percent passes through as literal"},
		{"%2", []byte("%2"), false, "partial percent sequence passes through as literals"},
		{"%GG", nil, true, "invalid hex chars"},
		{"", []byte{}, false, "empty input"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := percentDecodeToBytes(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("percentDecodeToBytes(%q): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("percentDecodeToBytes(%q): unexpected error: %v", tc.input, err)
				return
			}
			if string(got) != string(tc.want) {
				t.Errorf("percentDecodeToBytes(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestEncodeExtern(t *testing.T) {
	t.Parallel()
	fn := encodeExtern()

	tests := []struct {
		args    []values.BalValue
		want    string
		wantErr bool
		desc    string
	}{
		{[]values.BalValue{"hello world", "UTF-8"}, "hello%20world", false, "UTF-8 space encoding"},
		{[]values.BalValue{"abc~-_.", "UTF-8"}, "abc~-_.", false, "UTF-8 unreserved passthrough"},
		{[]values.BalValue{"a=b&c", "UTF-8"}, "a%3Db%26c", false, "UTF-8 special chars"},
		{[]values.BalValue{"hello", "US-ASCII"}, "hello", false, "ASCII passthrough"},
		{[]values.BalValue{"test", "UNKNOWN-CHARSET"}, "", true, "unknown charset is error"},
		{[]values.BalValue{42, "UTF-8"}, "", true, "non-string value is error"},
		{[]values.BalValue{"hello", 42}, "", true, "non-string charset is error"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := fn(nil, tc.args)
			if err != nil {
				t.Fatalf("encodeExtern: unexpected Go error: %v", err)
			}
			if tc.wantErr {
				if _, ok := got.(*values.Error); !ok {
					t.Errorf("expected BalError result, got %T: %v", got, got)
				}
				return
			}
			if got != tc.want {
				t.Errorf("encodeExtern(%v) = %q, want %q", tc.args, got, tc.want)
			}
		})
	}
}

func TestDecodeExtern(t *testing.T) {
	t.Parallel()
	fn := decodeExtern()

	tests := []struct {
		args    []values.BalValue
		want    string
		wantErr bool
		desc    string
	}{
		{[]values.BalValue{"hello%20world", "UTF-8"}, "hello world", false, "UTF-8 decode space"},
		{[]values.BalValue{"hello+world", "UTF-8"}, "hello world", false, "plus as space"},
		{[]values.BalValue{"abc~-_.", "UTF-8"}, "abc~-_.", false, "unreserved passthrough"},
		{[]values.BalValue{"hello", "UNKNOWN-CHARSET"}, "", true, "unknown charset is error"},
		{[]values.BalValue{42, "UTF-8"}, "", true, "non-string value is error"},
		{[]values.BalValue{"hello", 42}, "", true, "non-string charset is error"},
		{[]values.BalValue{"%GG", "UTF-8"}, "", true, "invalid percent encoding is error"},
		{[]values.BalValue{"%FF", "UTF-8"}, "", true, "invalid UTF-8 byte is error"},
		// Literal non-ASCII chars must pass through unchanged, not be re-interpreted
		// through the charset decoder (e.g. ISO-8859-1 would corrupt UTF-8 bytes).
		{[]values.BalValue{"caf\xc3\xa9", "ISO-8859-1"}, "caf\xc3\xa9", false, "ISO-8859-1: literal non-ASCII not charset-decoded"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := fn(nil, tc.args)
			if err != nil {
				t.Fatalf("decodeExtern: unexpected Go error: %v", err)
			}
			if tc.wantErr {
				if _, ok := got.(*values.Error); !ok {
					t.Errorf("expected BalError result, got %T: %v", got, got)
				}
				return
			}
			if got != tc.want {
				t.Errorf("decodeExtern(%v) = %q, want %q", tc.args, got, tc.want)
			}
		})
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	t.Parallel()
	encode := encodeExtern()
	decode := decodeExtern()

	inputs := []string{
		"hello world",
		"path/to/resource?key=value&other=test",
		"special: !@#$%^&*()",
		"unicode: café",
	}
	charsets := []string{"UTF-8", "ISO-8859-1", "US-ASCII"}
	for _, input := range inputs {
		for _, charset := range charsets {
			t.Run(input+"/"+charset, func(t *testing.T) {
				encoded, err := encode(nil, []values.BalValue{input, charset})
				if err != nil {
					t.Fatalf("encode error: %v", err)
				}
				if _, isErr := encoded.(*values.Error); isErr {
					return // skip if encoding fails (e.g. ASCII can't encode café)
				}
				decoded, err := decode(nil, []values.BalValue{encoded, charset})
				if err != nil {
					t.Fatalf("decode error: %v", err)
				}
				if _, isErr := decoded.(*values.Error); isErr {
					return // skip if decoding fails
				}
				if decoded != input {
					t.Errorf("roundtrip(%q, %s): got %q", input, charset, decoded)
				}
			})
		}
	}
}
