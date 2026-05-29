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
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "url"
)

// resolveEncoding maps a Java/IANA charset name to an x/text Encoding.
// Returns nil for unrecognised charsets.
func resolveEncoding(charset string) encoding.Encoding {
	switch strings.ToUpper(charset) {
	case "UTF-8", "UTF8":
		return encoding.Nop
	case "ISO-8859-1", "ISO8859-1", "ISO_8859_1", "LATIN-1", "LATIN1":
		return charmap.ISO8859_1
	case "US-ASCII", "ASCII":
		return charmap.ISO8859_1 // ASCII is a subset; ISO-8859-1 encoder errors on >U+00FF
	case "UTF-16":
		return unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	case "UTF-16BE":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "UTF-16LE":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	default:
		return nil
	}
}

// isUnreserved reports whether b is an RFC 3986 unreserved character that
// Java's URLEncoder also leaves unencoded: A-Z, a-z, 0-9, -, _, ., ~.
func isUnreserved(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '~'
}

// encodeBytes percent-encodes raw bytes, leaving unreserved bytes as-is.
// Space is encoded as %20 (not +) to match Ballerina/jBallerina behaviour.
func encodeBytes(raw []byte) string {
	var buf strings.Builder
	for _, b := range raw {
		if isUnreserved(b) {
			buf.WriteByte(b)
		} else {
			fmt.Fprintf(&buf, "%%%02X", b)
		}
	}
	return buf.String()
}

// percentDecodeToBytes decodes a percent-encoded string to raw bytes.
// Treats '+' as space, matching Java URLDecoder for all charsets.
func percentDecodeToBytes(s string) ([]byte, error) {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		switch {
		case s[i] == '+':
			buf = append(buf, ' ')
			i++
		case s[i] == '%' && i+2 < len(s):
			hi, ok1 := fromHex(s[i+1])
			lo, ok2 := fromHex(s[i+2])
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("invalid percent-encoding at position %d", i)
			}
			buf = append(buf, byte(hi<<4|lo))
			i += 3
		default:
			buf = append(buf, s[i])
			i++
		}
	}
	return buf, nil
}

func fromHex(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	default:
		return 0, false
	}
}

// encodeExtern replicates Java URLEncoder.encode() + Ballerina post-processing:
//   - space -> %20 (Java uses +, Ballerina converts to %20)
//   - * -> %2A (encoded, unlike RFC 3986 sub-delimiters)
//   - ~ -> ~ (unreserved, not encoded)
//
// For non-UTF-8 charsets the string is first converted to the target charset
// bytes, then those bytes are percent-encoded.
func encodeExtern() extern.NativeFunc {
	return func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
		value, ok := args[0].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while encoding. invalid string argument"), nil
		}
		charset, ok := args[1].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while encoding. invalid charset argument"), nil
		}
		enc := resolveEncoding(charset)
		if enc == nil {
			return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while encoding. %s", charset)), nil
		}

		var raw []byte
		if enc == encoding.Nop {
			raw = []byte(value)
		} else {
			converted, err := enc.NewEncoder().Bytes([]byte(value))
			if err != nil {
				return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while encoding. %s", err.Error())), nil
			}
			raw = converted
		}
		return encodeBytes(raw), nil
	}
}

// decodeExtern replicates Java URLDecoder.decode().
// For non-UTF-8 charsets the percent-decoded bytes are interpreted in the
// target charset and converted to UTF-8.
func decodeExtern() extern.NativeFunc {
	return func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
		value, ok := args[0].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while decoding. invalid string argument"), nil
		}
		charset, ok := args[1].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while decoding. invalid charset argument"), nil
		}
		enc := resolveEncoding(charset)
		if enc == nil {
			return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while decoding. %s", charset)), nil
		}

		raw, err := percentDecodeToBytes(value)
		if err != nil {
			return values.NewErrorWithMessage("Error occurred while decoding. " + err.Error()), nil
		}

		if enc == encoding.Nop {
			if !utf8.Valid(raw) {
				return values.NewErrorWithMessage("Error occurred while decoding. invalid UTF-8 sequence"), nil
			}
			return string(raw), nil
		}

		utf8Bytes, err := enc.NewDecoder().Bytes(raw)
		if err != nil {
			return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while decoding. %s", err.Error())), nil
		}
		return string(utf8Bytes), nil
	}
}

func initURLModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "encode", encodeExtern())
	runtime.RegisterExternFunction(rt, orgName, moduleName, "decode", decodeExtern())
}

func init() {
	runtime.RegisterModuleInitializer(initURLModule)
}
