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

// TextLine represents a single line in the TextDocument.
type TextLine interface {
	LineNo() int
	Text() string
	StartOffset() int
	EndOffset() int
	EndOffsetWithNewLines() int
	Length() int
	LengthWithNewLineChars() int
}

type textLineImpl struct {
	lineNo               int
	text                 string
	startOffset          int
	endOffset            int
	lengthOfNewLineChars int
}

func NewTextLine(lineNo int, text string, startOffset, endOffset, lengthOfNewLineChars int) TextLine {
	return &textLineImpl{
		lineNo:               lineNo,
		text:                 text,
		startOffset:          startOffset,
		endOffset:            endOffset,
		lengthOfNewLineChars: lengthOfNewLineChars,
	}
}

func (tl textLineImpl) LineNo() int {
	return tl.lineNo
}

func (tl textLineImpl) Text() string {
	return tl.text
}

func (tl textLineImpl) StartOffset() int {
	return tl.startOffset
}

func (tl textLineImpl) EndOffset() int {
	return tl.endOffset
}

func (tl textLineImpl) EndOffsetWithNewLines() int {
	return tl.endOffset + tl.lengthOfNewLineChars
}

func (tl textLineImpl) Length() int {
	return tl.endOffset - tl.startOffset
}

func (tl textLineImpl) LengthWithNewLineChars() int {
	return tl.endOffset - tl.startOffset + tl.lengthOfNewLineChars
}
