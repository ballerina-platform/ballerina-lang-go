// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package identifierutil

import (
	"testing"
)

func TestEscapeSpecialCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test\\name", "test\\\\name"},
		{"test name", "test\\ name"},
		{"$&+,:;=?@#\\|/' []<>.\"^*{}~`()%!-", "\\$\\&\\+\\,\\:\\;\\=\\?\\@\\#\\\\\\|\\/\\'\\ \\[\\]\\<\\>\\.\\\"\\^\\*\\{\\}\\~\\`\\(\\)\\%\\!\\-"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EscapeSpecialCharacters(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeSpecialCharacters(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeJava(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test\\\\name", "test\\name"},
		{"pre\\n\\t\\rpost", "pre\n\t\rpost"},
		{"\\'\\\"", "'\""},
		{"\\u0041\\u0061", "Aa"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeJava(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeJava(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecodeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"&0047", "/"},
		{"&0092", "\\"},
		{"&0046&0058&0059", ".:;"},
		{"prefix&0091index&0093", "prefix[index]"},
		{"test&0046name&0058value", "test.name:value"},
		{"&0060init&0062", "<init>"},
		{"&0046&0060init&0062", ".<init>"},
		{"$gen$test", "test"},
		{"$gen$&0046&0060init&0062", ".<init>"},
		{"&abcd", "&abcd"},
		{"&12ab", "&12ab"},
		{"&", "&"},
		{"&123", "&123"},
		{"&12345", "Ӓ5"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DecodeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeBallerina(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test\\nname", "test\nname"},
		{"test\\u{41}", "testA"},
		{"\\u{0048}ello", "Hello"},
		{"test\\u{1F600}", "test😀"},
		{"test\\u{5C}", "test\\"},
		{"\\\\u{61}", "\\u{61}"},
		{"test\\\\\\u{61}", "test\\a"},
		{"Line1\\nLine2\\u{0009}Tab", "Line1\nLine2\tTab"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeBallerina(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeBallerina(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeUnicodeCodepoints(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"", ""},
		{"no unicode here", "no unicode here"},
		{"\\u{41}", "A"},
		{"\\u{0048}ello", "Hello"},
		{"test\\u{1F600}emoji", "test😀emoji"},
		{"\\u{1F44D}", "👍"},
		{"\\u{48}\\u{65}\\u{6C}\\u{6C}\\u{6F}", "Hello"},
		{"prefix\\u{41}middle\\u{42}suffix", "prefixAmiddleBsuffix"},
		{"\\\\u{41}", "\\\\u{41}"},
		{"\\\\\\u{41}", "\\\\A"},
		{"test\\u{5C}end", "test\\u005Cend"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeUnicodeCodepoints(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeUnicodeCodepoints(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsEscapedNumericEscape(t *testing.T) {
	tests := []struct {
		leadingSlashes     string
		expectedNotEscaped bool
	}{
		{"", true},
		{"\\", false},
		{"\\\\", true},
		{"\\\\\\", false},
		{"\\\\\\\\", true},
	}

	for _, tt := range tests {
		t.Run(tt.leadingSlashes, func(t *testing.T) {
			result := IsEscapedNumericEscape(tt.leadingSlashes)
			expected := !tt.expectedNotEscaped
			if result != expected {
				t.Errorf("IsEscapedNumericEscape(%q) = %v, want %v", tt.leadingSlashes, result, expected)
			}
		})
	}
}

func TestEncodeFunctionIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".<init>", "$gen$&0046&0060init&0062"},
		{"normalFunction", "normalFunction"},
		{"method_name", "method_name"},
		{"method123", "method123"},
		{"test.method", "$gen$test&0046method"},
		{"test:method", "$gen$test&0058method"},
		{"test[method", "$gen$test&0091method"},
		{"test/method", "$gen$test&0047method"},
		{"test>method", "$gen$test&0062method"},
		{"test.method:value", "$gen$test&0046method&0058value"},
		{"test\\method", "testmethod"},
		{"test\\$method", "test&0036method"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EncodeFunctionIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeFunctionIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeNonFunctionIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"simpleField", "simpleField"},
		{"field_name", "field_name"},
		{"field123", "field123"},
		{"test.field", "test&0046field"},
		{"test:field", "test&0058field"},
		{"test/field", "test&0047field"},
		{"test[field", "test&0091field"},
		{"test<field", "test&0060field"},
		{"test>field", "test&0062field"},
		{"test.field:value", "test&0046field&0058value"},
		{"test\\field", "testfield"},
		{"test\\$field", "test&0036field"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EncodeNonFunctionIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeNonFunctionIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
