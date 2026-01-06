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

package semtypes

import "ballerina-lang-go/common"

type IntSubtype struct {
	Ranges []Range
}

var _ ProperSubtypeData = &IntSubtype{}

func NewIntSubtypeFromRanges(ranges []Range) IntSubtype {
	this := IntSubtype{}
	this.Ranges = ranges
	return this
}

func CreateIntSubtype(ranges ...Range) IntSubtype {
	// migrated from IntSubtype.java:43:5
	return NewIntSubtypeFromRanges(ranges)
}

func CreateSingleRangeSubtype(min, max int64) IntSubtype {
	return NewIntSubtypeFromRanges([]Range{RangeFrom(min, max)})
}

func IntConst(value int64) SemType {
	// migrated from IntSubtype.java:51:5
	return basicSubtype(BT_INT, CreateSingleRangeSubtype(value, value))
}

func validIntWidth(signed bool, bits int64) {
	if bits <= 0 {
		var message string
		if bits == 0 {
			message = "zero"
		} else {
			message = "negative"
		}
		message = message + " width in bits"
		panic(message)
	}
	if signed {
		if bits > 64 {
			panic("width of signed integers limited to 64")
		}
	} else {
		if bits > 63 {
			panic("width of unsigned integers limited to 63")
		}
	}
}

func ValidIntWidthSigned(bits int) {
	// migrated from IntSubtype.java:74:5
	validIntWidth(true, int64(bits))
}

func ValidIntWidthUnsigned(bits int) {
	// migrated from IntSubtype.java:78:5
	validIntWidth(false, int64(bits))
}

func IntWidthSigned(bits int64) SemType {
	// migrated from IntSubtype.java:82:5
	validIntWidth(true, bits)
	if bits == 64 {
		return &INT
	}
	t := CreateSingleRangeSubtype((-(int64(1) << (bits - int64(1)))), ((int64(1) << (bits - int64(1))) - int64(1)))
	return basicSubtype(BT_INT, t)
}

func IntWidthUnsigned(bits int) SemType {
	// migrated from IntSubtype.java:91:5
	validIntWidth(false, int64(bits))
	t := CreateSingleRangeSubtype(int64(0), ((int64(1) << bits) - int64(1)))
	return basicSubtype(BT_INT, t)
}

func IntSubtypeWidenUnsigned(d SubtypeData) SubtypeData {
	// migrated from IntSubtype.java:98:5
	if _, ok := d.(AllOrNothingSubtype); ok {
		return d
	}
	v := d.(IntSubtype)
	if v.Ranges[0].Min < int64(0) {
		return CreateAll()
	}
	r := v.Ranges[len(v.Ranges)-1]
	i := int64(8)
	for i <= int64(32) {
		if r.Max < (int64(1) << i) {
			w := CreateSingleRangeSubtype(int64(0), ((int64(1) << i) - 1))
			return w
		}
		i = (i * 2)
	}
	return CreateAll()
}

func IntSubtypeSingleValue(d SubtypeData) common.Optional[int64] {
	if _, ok := d.(AllOrNothingSubtype); ok {
		return common.OptionalEmpty[int64]()
	}
	v := d.(IntSubtype)
	if len(v.Ranges) != 1 {
		return common.OptionalEmpty[int64]()
	}
	min := v.Ranges[0].Min
	if min != v.Ranges[0].Max {
		return common.OptionalEmpty[int64]()
	}
	return common.OptionalOf(min)
}

func IntSubtypeContains(d SubtypeData, n int64) bool {
	// migrated from IntSubtype.java:137:5
	if allOrNothingSubtype, ok := d.(AllOrNothingSubtype); ok {
		return allOrNothingSubtype.IsAllSubtype()
	}
	v := d.(IntSubtype)
	for _, r := range v.Ranges {
		if (r.Min <= n) && (n <= r.Max) {
			return true
		}
	}
	return false
}
