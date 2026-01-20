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

package diagnostics

import "ballerina-lang-go/tools/text"

// Location represents the location in TextDocument.
// It is a combination of source file path, start and end line numbers, and start and end column numbers.
type Location interface {
	LineRange() text.LineRange
	TextRange() text.TextRange
}

// BLangDiagnosticLocation is a minimal implementation of Location for AST nodes
type BLangDiagnosticLocation struct {
	filePath    string
	startLine   int
	endLine     int
	startColumn int
	endColumn   int
	startOffset int
	length      int
}

func NewBLangDiagnosticLocation(
	filePath string,
	startLine, endLine, startColumn, endColumn, startOffset, length int,
) Location {
	return &BLangDiagnosticLocation{
		filePath:    filePath,
		startLine:   startLine,
		endLine:     endLine,
		startColumn: startColumn,
		endColumn:   endColumn,
		startOffset: startOffset,
		length:      length,
	}
}

var _ Location = &BLangDiagnosticLocation{}

func (loc *BLangDiagnosticLocation) LineRange() text.LineRange {
	startLinePos := text.LinePositionFromLineAndOffset(loc.startLine, loc.startColumn)
	endLinePos := text.LinePositionFromLineAndOffset(loc.endLine, loc.endColumn)
	return text.LineRangeFromLinePositions(loc.filePath, startLinePos, endLinePos)
}

func (loc *BLangDiagnosticLocation) TextRange() text.TextRange {
	return text.TextRangeFromStartOffsetAndLength(loc.startOffset, loc.length)
}
