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
	"fmt"

	"ballerina-lang-go/common/errors"
)

// LineMap represents a collection of text lines in the TextDocument.
type LineMap interface {
	TextLine(line int) (TextLine, error)
	LinePositionFromPosition(position int) (LinePosition, error)
	TextPositionFromLinePosition(linePosition LinePosition) (int, error)
	TextLines() []string
}

type lineMapImpl struct {
	textLines []TextLine
	length    int
}

func NewLineMap(textLines []TextLine) LineMap {
	return &lineMapImpl{
		textLines: textLines,
		length:    len(textLines),
	}
}

func (lm lineMapImpl) TextLine(line int) (TextLine, error) {
	if err := lm.lineRangeCheck(line); err != nil {
		return nil, err
	}
	return lm.textLines[line], nil
}

func (lm lineMapImpl) LinePositionFromPosition(position int) (LinePosition, error) {
	if err := lm.positionRangeCheck(position); err != nil {
		return nil, err
	}
	textLine := lm.findLineFrom(position)
	return LinePositionFromLineAndOffset(textLine.LineNo(), position-textLine.StartOffset()), nil
}

func (lm lineMapImpl) TextPositionFromLinePosition(linePosition LinePosition) (int, error) {
	if err := lm.lineRangeCheck(linePosition.Line()); err != nil {
		return -1, err
	}
	textLine := lm.textLines[linePosition.Line()]
	if textLine.Length() < linePosition.Offset() {
		return -1, errors.NewIllegalArgumentError(fmt.Sprintf("Cannot find a line with the character offset '%d'", linePosition.Offset()))
	}
	// TODO: Lazy initialize and cache
	return textLine.StartOffset() + linePosition.Offset(), nil
}

func (lm lineMapImpl) TextLines() []string {
	lines := make([]string, len(lm.textLines))
	for i, textLine := range lm.textLines {
		lines[i] = textLine.Text()
	}
	return lines
}

func (lm lineMapImpl) positionRangeCheck(position int) error {
	if position < 0 || position > lm.textLines[lm.length-1].EndOffset() {
		return errors.NewIndexOutOfBoundsError(position, lm.textLines[lm.length-1].EndOffset())
	}
	return nil
}

func (lm lineMapImpl) lineRangeCheck(lineNo int) error {
	if lineNo < 0 || lineNo > lm.length {
		return errors.NewIndexOutOfBoundsError(lineNo, lm.length)
	}
	return nil
}

// findLineFrom returns the TextLine to which the given position belongs.
// Performs a binary search to find the matching text line.
func (lm lineMapImpl) findLineFrom(position int) TextLine {
	// Check boundary conditions
	if position == 0 {
		return lm.textLines[0]
	} else if position == lm.textLines[lm.length-1].EndOffset() {
		return lm.textLines[lm.length-1]
	}
	left := 0
	right := lm.length - 1
	for left <= right {
		lhs := left >> 1
		rhs := right >> 1
		middle := (lhs + rhs) + (left & right & 1)
		startOffset := lm.textLines[middle].StartOffset()
		endOffset := lm.textLines[middle].EndOffsetWithNewLines()
		if startOffset <= position && position < endOffset {
			return lm.textLines[middle]
		} else if endOffset <= position {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}
	// This should never happen given the boundary checks above
	panic("binary search failed to find matching text line")
}
