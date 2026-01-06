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

type StringOps struct {
}

var _ BasicTypeOps = &StringOps{}

func NewStringOps() StringOps {
	this := StringOps{}
	return this
}

func (this *StringOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from StringOps.java:45:5
	sd1 := d1.(*StringSubtype)
	sd2 := d2.(*StringSubtype)
	var chars []EnumerableType[string]
	var nonChars []EnumerableType[string]
	charsAllowed := EnumerableSubtypeUnion(sd1.GetChar(), sd2.GetChar(), &chars)
	nonCharsAllowed := EnumerableSubtypeUnion(sd1.GetNonChar(), sd2.GetNonChar(), &nonChars)
	return CreateStringSubtype(CharStringSubtypeFrom(charsAllowed, chars), NonCharStringSubtypeFrom(nonCharsAllowed, nonChars))
}

func (this *StringOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from StringOps.java:64:5
	if allOrNothing1, ok := d1.(*AllOrNothingSubtype); ok {
		if allOrNothing1.IsAllSubtype() {
			return d2
		} else {
			return CreateNothing()
		}
	}
	if allOrNothing2, ok := d2.(*AllOrNothingSubtype); ok {
		if allOrNothing2.IsAllSubtype() {
			return d1
		} else {
			return CreateNothing()
		}
	}
	sd1 := d1.(*StringSubtype)
	sd2 := d2.(*StringSubtype)
	var chars []EnumerableType[string]
	var nonChars []EnumerableType[string]
	charsAllowed := EnumerableSubtypeIntersect(sd1.GetChar(), sd2.GetChar(), &chars)
	nonCharsAllowed := EnumerableSubtypeIntersect(sd1.GetNonChar(), sd2.GetNonChar(), &nonChars)
	return CreateStringSubtype(CharStringSubtypeFrom(charsAllowed, chars), NonCharStringSubtypeFrom(nonCharsAllowed, nonChars))
}

func (this *StringOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from StringOps.java:86:5
	return this.Intersect(d1, this.Complement(d2))
}

func (this *StringOps) Complement(d SubtypeData) SubtypeData {
	// migrated from StringOps.java:91:5
	st := d.(*StringSubtype)
	if len(st.GetChar().Values()) == 0 && len(st.GetNonChar().Values()) == 0 {
		if st.GetChar().Allowed() && st.GetNonChar().Allowed() {
			return CreateAll()
		} else if !st.GetChar().Allowed() && !st.GetNonChar().Allowed() {
			return CreateNothing()
		}
	}
	return CreateStringSubtype(CharStringSubtypeFrom(!st.GetChar().Allowed(), st.GetChar().Values()), NonCharStringSubtypeFrom(!st.GetNonChar().Allowed(), st.GetNonChar().Values()))
}

func (this *StringOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from StringOps.java:106:5
	return notIsEmpty(cx, t)
}

func stringSubtypeListCoverage(subtype StringSubtype, values []string) StringSubtypeListCoverage {
	// migrated from StringOps.java:113:5
	var indices []int
	ch := subtype.GetChar()
	nonChar := subtype.GetNonChar()
	stringConsts := 0
	if ch.Allowed() {
		StringListIntersect(values, toStringArray(ch.Values()), &indices)
		stringConsts = len(ch.Values())
	} else if len(ch.Values()) == 0 {
		for i := 0; i < len(values); i++ {
			if len(values[i]) == 1 {
				indices = append(indices, i)
			}
		}
	}
	if nonChar.Allowed() {
		StringListIntersect(values, toStringArray(nonChar.Values()), &indices)
		stringConsts += len(nonChar.Values())
	} else if len(nonChar.Values()) == 0 {
		for i := 0; i < len(values); i++ {
			if len(values[i]) != 1 {
				indices = append(indices, i)
			}
		}
	}
	return StringSubtypeListCoverageFrom(stringConsts == len(indices), indices)
}

func toStringArray(ar []EnumerableType[string]) []string {
	strings := make([]string, len(ar))
	for i, value := range ar {
		strings[i] = value.Value()
	}
	return strings
}

func StringListIntersect(values []string, target []string, indices *[]int) {
	// migrated from StringOps.java:158:5
	i1 := 0
	i2 := 0
	len1 := len(values)
	len2 := len(target)
	for {
		if i1 >= len1 || i2 >= len2 {
			break
		} else {
			comp := CompareEnumerable(EnumerableStringFrom(values[i1]), EnumerableStringFrom(target[i2]))
			switch comp {
			case EQ:
				*indices = append(*indices, i1)
				i1 = i1 + 1
				i2 = i2 + 1
			case LT:
				i1 = i1 + 1
			case GT:
				i2 = i2 + 1
			default:
				panic("Invalid comparison value!")
			}
		}
	}
}
