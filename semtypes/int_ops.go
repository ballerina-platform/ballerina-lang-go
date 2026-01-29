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

type IntOps struct {
}

var _ BasicTypeOps = &IntOps{}

func NewIntOps() IntOps {
	this := IntOps{}
	return this
}

var intOpsInstance = NewIntOps()

func (this *IntOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from IntOps.java:45:5
	v1 := d1.(IntSubtype)
	v2 := d2.(IntSubtype)
	v := rangeListUnion(v1.Ranges, v2.Ranges)
	if len(v) == 1 && v[0].Min == MIN_VALUE && v[0].Max == MAX_VALUE {
		return CreateAll()
	}
	return CreateIntSubtype(v...)
}

func (this *IntOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from IntOps.java:56:5
	v1 := d1.(IntSubtype)
	v2 := d2.(IntSubtype)
	v := rangeListIntersect(v1.Ranges, v2.Ranges)
	if len(v) == 0 {
		return CreateNothing()
	}
	return CreateIntSubtype(v...)
}

func (this *IntOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from IntOps.java:67:5
	v1 := d1.(IntSubtype)
	v2 := d2.(IntSubtype)
	v := rangeListIntersect(v1.Ranges, rangeListComplement(v2.Ranges))
	if len(v) == 0 {
		return CreateNothing()
	}
	return CreateIntSubtype(v...)
}

func (this *IntOps) Complement(d SubtypeData) SubtypeData {
	// migrated from IntOps.java:78:5
	v := d.(IntSubtype)
	return CreateIntSubtype(rangeListComplement(v.Ranges)...)
}

func intSubtypeOverlapRange(subtype IntSubtype, r Range) bool {
	// migrated from IntOps.java:83:5
	subtypeData := intOpsInstance.Intersect(subtype, CreateIntSubtype(r))
	if allOrNothingSubtype, ok := subtypeData.(AllOrNothingSubtype); ok {
		return !allOrNothingSubtype.IsNothingSubtype()
	}
	return true
}

func intSubtypeMax(subtype IntSubtype) int64 {
	// migrated from IntOps.java:89:5
	return subtype.Ranges[len(subtype.Ranges)-1].Max
}

func intSubtypeMin(subtype IntSubtype) int64 {
	// migrated from IntOps.java:93:5
	return subtype.Ranges[0].Min
}

func (this *IntOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from IntOps.java:98:5
	return notIsEmpty(cx, t)
}

func rangeListUnion(v1 []Range, v2 []Range) []Range {
	// migrated from IntOps.java:102:5
	var result []Range
	i1 := 0
	i2 := 0
	len1 := len(v1)
	len2 := len(v2)

	for {
		if i1 >= len1 {
			if i2 >= len2 {
				break
			}
			result = rangeUnionPush(result, v2[i2])
			i2 += 1
		} else if i2 >= len2 {
			result = rangeUnionPush(result, v1[i1])
			i1 += 1
		} else {
			r1 := v1[i1]
			r2 := v2[i2]
			combined := rangeUnion(r1, r2)
			if combined.Status == 0 {
				result = rangeUnionPush(result, *combined.Range)
				i1 += 1
				i2 += 1
			} else if combined.Status < 0 {
				result = rangeUnionPush(result, r1)
				i1 += 1
			} else {
				result = rangeUnionPush(result, r2)
				i2 += 1
			}
		}
	}
	return result
}

func rangeUnionPush(ranges []Range, next Range) []Range {
	// migrated from IntOps.java:140:5
	lastIndex := len(ranges) - 1
	if lastIndex < 0 {
		return append(ranges, next)
	}
	combined := rangeUnion(ranges[lastIndex], next)
	if combined.Status == 0 {
		ranges[lastIndex] = *combined.Range
		return ranges
	}
	return append(ranges, next)
}

func rangeUnion(r1 Range, r2 Range) RangeUnion {
	// migrated from IntOps.java:157:5
	if r1.Max < r2.Min {
		if r1.Max+1 != r2.Min {
			return RangeUnionFrom(-1)
		}
	}
	if r2.Max < r1.Min {
		if r2.Max+1 != r1.Min {
			return RangeUnionFrom(1)
		}
	}
	return FromWithRange(RangeFrom(min(r1.Min, r2.Min), max(r1.Max, r2.Max)))
}

func rangeListIntersect(v1 []Range, v2 []Range) []Range {
	// migrated from IntOps.java:171:5
	var result []Range
	i1 := 0
	i2 := 0
	len1 := len(v1)
	len2 := len(v2)
	for {
		if i1 >= len1 || i2 >= len2 {
			break
		} else {
			r1 := v1[i1]
			r2 := v2[i2]
			combined := rangeIntersect(r1, r2)
			if combined.Status == 0 {
				result = append(result, *combined.Range)
				i1 += 1
				i2 += 1
			} else if combined.Status < 0 {
				i1 += 1
			} else {
				i2 += 1
			}
		}
	}
	return result
}

func rangeIntersect(r1 Range, r2 Range) RangeUnion {
	// migrated from IntOps.java:202:5
	if r1.Max < r2.Min {
		return RangeUnionFrom(-1)
	}
	if r2.Max < r1.Min {
		return RangeUnionFrom(1)
	}
	return FromWithRange(RangeFrom(max(r1.Min, r2.Min), min(r1.Max, r2.Max)))
}

func rangeListComplement(v []Range) []Range {
	// migrated from IntOps.java:212:5
	var result []Range
	length := len(v)
	minVal := v[0].Min
	if minVal > MIN_VALUE {
		result = append(result, RangeFrom(MIN_VALUE, minVal-1))
	}
	for i := 1; i < length; i++ {
		result = append(result, RangeFrom(v[i-1].Max+1, v[i].Min-1))
	}
	maxVal := v[len(v)-1].Max
	if maxVal < MAX_VALUE {
		result = append(result, RangeFrom(maxVal+1, MAX_VALUE))
	}
	return result
}
