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

// TextDocument is an abstract representation of a Ballerina source file (.bal).
type TextDocument interface {
	Apply(textDocumentChange TextDocumentChange) TextDocument
	ToCharArray() []rune
	Line(line int) (TextLine, error)
	LinePositionFromTextPosition(textPosition int) (LinePosition, error)
	TextPositionFromLinePosition(linePosition LinePosition) (int, error)
	TextLines() []string
	Lines() LineMap
	PopulateTextLineMap() LineMap
	String() string
}

type textDocumentBase struct {
	lineMap LineMap
}

func (td textDocumentBase) Line(line int) (TextLine, error) {
	return td.lineMap.TextLine(line)
}

func (td textDocumentBase) LinePositionFromTextPosition(textPosition int) (LinePosition, error) {
	return td.lineMap.LinePositionFromPosition(textPosition)
}

func (td textDocumentBase) TextPositionFromLinePosition(linePosition LinePosition) (int, error) {
	return td.lineMap.TextPositionFromLinePosition(linePosition)
}
