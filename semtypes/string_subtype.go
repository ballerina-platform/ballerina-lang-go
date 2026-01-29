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

type StringSubtypeListCoverage struct {
	IsSubtype bool
	Indices   []int
}

type StringSubtype struct {
	charData    CharStringSubtype
	nonCharData NonCharStringSubtype
}

var EMPTY_STRING_ARR = []EnumerableType[string]{}
var EMPTY_CHAR_ARR = []EnumerableType[string]{}
var _ ProperSubtypeData = &StringSubtype{}

func NewStringSubtypeListCoverageFromBoolInts(isSubtype bool, indices []int) StringSubtypeListCoverage {
	this := StringSubtypeListCoverage{}
	this.IsSubtype = isSubtype
	this.Indices = indices
	return this
}

func StringSubtypeListCoverageFrom(isSubtype bool, indices []int) StringSubtypeListCoverage {
	// migrated from StringSubtype.java:140:9
	return NewStringSubtypeListCoverageFromBoolInts(isSubtype, indices)
}

func newStringSubtypeFromCharStringSubtypeNonCharStringSubtype(charData CharStringSubtype, nonCharData NonCharStringSubtype) StringSubtype {
	this := StringSubtype{}
	this.charData = charData
	this.nonCharData = nonCharData
	return this
}

func StringSubtypeFrom(chara CharStringSubtype, nonChar NonCharStringSubtype) StringSubtype {
	// migrated from StringSubtype.java:56:5
	return newStringSubtypeFromCharStringSubtypeNonCharStringSubtype(chara, nonChar)
}

func StringSubtypeContains(d SubtypeData, s string) bool {
	// migrated from StringSubtype.java:60:5
	if allOrNothingSubtype, ok := d.(AllOrNothingSubtype); ok {
		return allOrNothingSubtype.IsAllSubtype()
	}
	st := d.(StringSubtype)
	chara := st.charData
	nonChar := st.nonCharData
	if len(s) == 1 {
		charString := EnumerableCharStringFrom(s)
		for _, value := range chara.Values() {
			if value == charString {
				return chara.Allowed()
			}
		}
		return !nonChar.Allowed()
	}
	stringString := EnumerableStringFrom(s)
	for _, value := range nonChar.Values() {
		if value == stringString {
			return nonChar.Allowed()
		}
	}
	return !nonChar.Allowed()
}

func CreateStringSubtype(chara CharStringSubtype, nonChar NonCharStringSubtype) SubtypeData {
	// migrated from StringSubtype.java:73:5
	if len(chara.Values()) == 0 && len(nonChar.Values()) == 0 {
		if (!chara.allowed) && (!nonChar.allowed) {
			return CreateAll()
		} else if chara.allowed && nonChar.allowed {
			return CreateNothing()
		}
	}
	return StringSubtypeFrom(chara, nonChar)
}

func StringSubtypeSingleValue(d SubtypeData) common.Optional[string] {
	if _, ok := d.(AllOrNothingSubtype); ok {
		return common.OptionalEmpty[string]()
	}
	st := d.(StringSubtype)
	chara := st.charData
	nonChar := st.nonCharData
	var charCount int
	if chara.Allowed() {
		charCount = len(chara.Values())
	} else {
		charCount = 2
	}
	var nonCharCount int
	if nonChar.Allowed() {
		nonCharCount = len(nonChar.Values())
	} else {
		nonCharCount = 2
	}
	if charCount+nonCharCount == 1 {
		if charCount != 0 {
			return common.OptionalOf(chara.Values()[0].Value())
		}
		return common.OptionalOf(nonChar.Values()[0].Value())
	}
	return common.OptionalEmpty[string]()
}

func StringConst(value string) SemType {
	// migrated from StringSubtype.java:100:5
	var chara CharStringSubtype
	var nonChar NonCharStringSubtype
	if codePointCount(value, 0, len(value)) == 1 {
		chara = CharStringSubtypeFrom(true, []EnumerableType[string]{EnumerableCharStringFrom(value)})
		nonChar = NonCharStringSubtypeFrom(true, EMPTY_STRING_ARR)
	} else {
		chara = CharStringSubtypeFrom(true, EMPTY_CHAR_ARR)
		nonChar = NonCharStringSubtypeFrom(true, []EnumerableType[string]{EnumerableStringFrom(value)})
	}
	return basicSubtype(BT_STRING, newStringSubtypeFromCharStringSubtypeNonCharStringSubtype(chara, nonChar))
}

func codePointCount(s string, start, end int) int {
	return len([]rune(s[start:end]))
}

func (this *StringSubtype) GetChar() EnumerableSubtype[string] {
	// migrated from StringSubtype.java:43:5
	return &this.charData
}

func (this *StringSubtype) GetNonChar() EnumerableSubtype[string] {
	// migrated from StringSubtype.java:47:5
	return &this.nonCharData
}

func StringChar() SemType {
	st := newStringSubtypeFromCharStringSubtypeNonCharStringSubtype(
		CharStringSubtypeFrom(false, EMPTY_CHAR_ARR),
		NonCharStringSubtypeFrom(true, EMPTY_STRING_ARR),
	)
	return basicSubtype(BT_STRING, st)
}
