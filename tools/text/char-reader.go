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

package text

import (
	"unicode"
	"unicode/utf8"
)

// CharReader is a character reader utility used by the Ballerina lexer.
type CharReader interface {
	Reset(offset int)
	Peek() rune
	PeekN(k int) rune
	Advance()
	AdvanceN(k int)
	Mark()
	GetMarkedChars() string
	IsEOF() bool
}

// charReaderImpl is the concrete implementation of CharReader.
type charReaderImpl struct {
	charBuffer       string
	offset           int
	charBufferLength int
	lexemeStartPos   int
}

func CharReaderFromTextDocument(textDocument TextDocument) CharReader {
	return CharReaderFromText(textDocument.String())
}

func CharReaderFromText(text string) CharReader {
	return &charReaderImpl{
		charBuffer:       text,
		offset:           0,
		charBufferLength: len(text),
		lexemeStartPos:   0,
	}
}

func (cr *charReaderImpl) Reset(offset int) {
	cr.offset = offset
}

func (cr charReaderImpl) Peek() rune {
	if cr.offset < cr.charBufferLength {
		r, _ := utf8.DecodeRuneInString(cr.charBuffer[cr.offset:])
		return r
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

func (cr charReaderImpl) PeekN(k int) rune {
	n := cr.offset
	for range k {
		_, size := utf8.DecodeRuneInString(cr.charBuffer[n:])
		n = n + size
	}

	if n < cr.charBufferLength {
		r, _ := utf8.DecodeRuneInString(cr.charBuffer[n:])
		return r
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

func (cr *charReaderImpl) Advance() {
	_, size := utf8.DecodeRuneInString(cr.charBuffer[cr.offset:])
	cr.offset = cr.offset + size
}

func (cr *charReaderImpl) AdvanceN(k int) {
	for range k {
		if cr.offset < cr.charBufferLength {
			_, size := utf8.DecodeRuneInString(cr.charBuffer[cr.offset:])
			cr.offset = cr.offset + size
		}
	}
}

func (cr *charReaderImpl) Mark() {
	cr.lexemeStartPos = cr.offset
}

func (cr charReaderImpl) GetMarkedChars() string {
	return cr.charBuffer[cr.lexemeStartPos:cr.offset]
}

func (cr charReaderImpl) IsEOF() bool {
	return cr.offset >= cr.charBufferLength
}
