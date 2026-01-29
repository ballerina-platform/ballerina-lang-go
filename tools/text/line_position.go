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

// LinePosition represents a line number and a character offset from the start of the line.
type LinePosition interface {
	Line() int
	Offset() int
	String() string
	LinePositionLookupKey() LinePositionLookupKey
}

// LinePositionLookupKey represents the comparable fields of LinePosition for equality/hashing.
type LinePositionLookupKey struct {
	Line   int
	Offset int
}

type linePositionImpl struct {
	line   int
	offset int
}

func LinePositionFromLineAndOffset(line, offset int) LinePosition {
	return &linePositionImpl{
		line:   line,
		offset: offset,
	}
}

func (lp linePositionImpl) Line() int {
	return lp.line
}

func (lp linePositionImpl) Offset() int {
	return lp.offset
}

func (lp linePositionImpl) String() string {
	return fmt.Sprintf("%d:%d", lp.line, lp.offset)
}

// LinePositionLookupKey returns the lookup key for equality comparisons.
func (lp linePositionImpl) LinePositionLookupKey() LinePositionLookupKey {
	return LinePositionLookupKey{
		Line:   lp.line,
		Offset: lp.offset,
	}
}
