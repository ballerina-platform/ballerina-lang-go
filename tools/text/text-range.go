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

import "fmt"

// TextRange describes a contiguous sequence of unicode code points in the TextDocument.
type TextRange interface {
	StartOffset() int
	EndOffset() int
	Length() int
	Contains(position int) bool
	IntersectionExists(textRange TextRange) bool
	String() string
	TextRangeLookupKey() TextRangeLookupKey
}

// TextRangeLookupKey represents the comparable fields of TextRange for equality/hashing.
type TextRangeLookupKey struct {
	StartOffset int
	EndOffset   int
}

type textRangeImpl struct {
	startOffset int
	endOffset   int
	length      int
}

func TextRangeFromStartOffsetAndLength(startOffset, length int) TextRange {
	return &textRangeImpl{
		startOffset: startOffset,
		length:      length,
		endOffset:   startOffset + length,
	}
}

func (tr textRangeImpl) StartOffset() int {
	return tr.startOffset
}

func (tr textRangeImpl) EndOffset() int {
	return tr.endOffset
}

func (tr textRangeImpl) Length() int {
	return tr.length
}

func (tr textRangeImpl) Contains(position int) bool {
	return tr.startOffset <= position && position < tr.endOffset
}

// IntersectionExists tests whether there exists an intersection of this range and the given range.
// The ranges R1(S1, E1) and R2(S2, E2) intersects if S1 is greater than or equal to E2 and
// S2 is less than or equal to E1.
func (tr textRangeImpl) IntersectionExists(textRange TextRange) bool {
	return tr.startOffset <= textRange.EndOffset() && textRange.StartOffset() <= tr.endOffset
}

func (tr textRangeImpl) String() string {
	return fmt.Sprintf("(%d,%d)", tr.startOffset, tr.endOffset)
}

// TextRangeLookupKey returns the lookup key for equality comparisons.
func (tr textRangeImpl) TextRangeLookupKey() TextRangeLookupKey {
	return TextRangeLookupKey{
		StartOffset: tr.startOffset,
		EndOffset:   tr.endOffset,
	}
}
