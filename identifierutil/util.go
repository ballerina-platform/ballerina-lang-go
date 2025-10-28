/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package identifierutil

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	charPrefix            = "&"
	escapePrefix          = "\\"
	generatedMethodPrefix = "$gen$"
	unicodePointLen       = 5
)

var (
	unicodePattern              = regexp.MustCompile(`\\(\\*)u\{([a-fA-F0-9]+)\}`)
	unescapedSpecialCharPattern = regexp.MustCompile(`([$&+,:;=\?@#\\|/'\\ \[\}\]<>."^*{}~` + "`" + `()%!-])`)

	jvmReservedChars = map[rune]string{
		'\\': "0092",
		'.':  "0046",
		':':  "0058",
		';':  "0059",
		'[':  "0091",
		']':  "0093",
		'/':  "0047",
		'<':  "0060",
		'>':  "0062",
	}

	javaEscapes = map[byte]byte{
		'n':  '\n',
		't':  '\t',
		'r':  '\r',
		'b':  '\b',
		'f':  '\f',
		'\\': '\\',
		'"':  '"',
		'\'': '\'',
	}
)

type identifier struct {
	name      string
	isEncoded bool
}

func encodeSpecialCharacters(id string) string {
	var sb strings.Builder
	i := 0

	for i < len(id) {
		if id[i] == '\\' && i+1 < len(id) {
			if formatted := getFormattedStringForQuotedIdentifiers(rune(id[i+1])); formatted != "" {
				unicodePoint := charPrefix + formatted
				sb.WriteString(unicodePoint)
				i = i + 2
				continue
			}
			i = i + 1
		}
		sb.WriteByte(id[i])
		i = i + 1
	}
	return sb.String()
}

func EscapeSpecialCharacters(id string) string {
	return unescapedSpecialCharPattern.ReplaceAllString(id, "\\$1")
}

func encodeIdentifier(id string) string {
	if strings.Contains(id, escapePrefix) {
		return UnescapeJava(encodeSpecialCharacters(id))
	}
	return id
}

func UnescapeJava(str string) string {
	if str == "" {
		return str
	}

	var sb strings.Builder
	sb.Grow(len(str))

	for i := 0; i < len(str); i++ {
		if str[i] != '\\' || i+1 >= len(str) {
			sb.WriteByte(str[i])
			continue
		}

		next := str[i+1]

		if escaped, ok := javaEscapes[next]; ok {
			sb.WriteByte(escaped)
			i++
			continue
		}

		if next == 'u' && i+5 < len(str) {
			if codePoint, err := strconv.ParseInt(str[i+2:i+6], 16, 32); err == nil {
				sb.WriteRune(rune(codePoint))
				i += 5
				continue
			}
		}

		sb.WriteByte(str[i])
	}

	return sb.String()
}

func encodeGeneratedName(id string) identifier {
	var sb strings.Builder
	sb.Grow(len(id))

	isEncoded := false

	for _, ch := range id {
		if formatted, ok := jvmReservedChars[ch]; ok {
			sb.WriteString(charPrefix)
			sb.WriteString(formatted)
			isEncoded = true
		} else {
			sb.WriteRune(ch)
		}
	}
	return identifier{name: sb.String(), isEncoded: isEncoded}
}

func getFormattedStringForQuotedIdentifiers(c rune) string {
	if c == '$' {
		return "0036"
	}
	return jvmReservedChars[c]
}

func DecodeIdentifier(encodedId string) string {
	if encodedId == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(encodedId))

	for i := 0; i < len(encodedId); i++ {
		if encodedId[i] == '&' && i+4 < len(encodedId) && isUnicodePoint(encodedId, i) {
			codePoint, _ := strconv.ParseInt(encodedId[i+1:i+5], 10, 32)
			sb.WriteRune(rune(codePoint))
			i += 4
		} else {
			sb.WriteByte(encodedId[i])
		}
	}

	return decodeGeneratedMethodName(sb.String())
}

func decodeGeneratedMethodName(decodedName string) string {
	return strings.TrimPrefix(decodedName, generatedMethodPrefix)
}

func UnescapeBallerina(text string) string {
	return UnescapeJava(UnescapeUnicodeCodepoints(text))
}

func UnescapeUnicodeCodepoints(id string) string {
	return unicodePattern.ReplaceAllStringFunc(id, func(match string) string {
		submatch := unicodePattern.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		leadingSlashes := submatch[1]
		if IsEscapedNumericEscape(leadingSlashes) {
			return match
		}

		codePoint, err := strconv.ParseInt(submatch[2], 16, 32)
		if err != nil {
			return match
		}

		ch := rune(codePoint)
		if ch == '\\' {
			return leadingSlashes + "\\u005C"
		}

		return leadingSlashes + string(ch)
	})
}

func IsEscapedNumericEscape(leadingSlashes string) bool {
	return len(leadingSlashes)&1 != 0
}

func isUnicodePoint(encodedName string, index int) bool {
	if index+unicodePointLen > len(encodedName) {
		return false
	}
	return containsOnlyDigits(encodedName[index+1 : index+5])
}

func containsOnlyDigits(digitString string) bool {
	for _, ch := range digitString {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func EncodeFunctionIdentifier(functionName string) string {
	functionName = encodeIdentifier(functionName)

	specialCases := map[string]string{
		".<init>":     "$gen$&0046&0060init&0062",
		".<start>":    "$gen$&0046&0060start&0062",
		".<stop>":     "$gen$&0046&0060stop&0062",
		".<testinit>": "$gen$&0046&0060testinit&0062",
	}

	if encoded, ok := specialCases[functionName]; ok {
		return encoded
	}

	encodedName := encodeGeneratedName(functionName)
	if encodedName.isEncoded {
		return generatedMethodPrefix + encodedName.name
	}
	return functionName
}

func EncodeNonFunctionIdentifier(identifierString string) string {
	identifierString = encodeIdentifier(identifierString)
	return encodeGeneratedName(identifierString).name
}
