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

import "math"

type EnumerableSubtypeData interface {
}

type EnumerableSubtype[T any] interface {
	EnumerableSubtypeData
	Allowed() bool
	Values() []EnumerableType[T]
}

type EnumerableSubtypeBase struct {
}

type EnumerableSubtypeMethods[T any] struct {
	Self EnumerableSubtype[T]
}

var LT = (-1)
var EQ = 0
var GT = 1

func EnumerableSubtypeUnion[T any](t1 EnumerableSubtype[T], t2 EnumerableSubtype[T], result *[]EnumerableType[T]) bool {
	// migrated from EnumerableSubtype.java:37:5
	b1 := t1.Allowed()
	b2 := t2.Allowed()
	var allowed bool
	if b1 && b2 {
		EnumerableListUnion(t1.Values(), t2.Values(), result)
		allowed = true
	} else if (!b1) && (!b2) {
		EnumerableListIntersect(t1.Values(), t2.Values(), result)
		allowed = false
	} else if b1 && (!b2) {
		EnumerableListDiff(t2.Values(), t1.Values(), result)
		allowed = false
	} else {
		EnumerableListDiff(t1.Values(), t2.Values(), result)
		allowed = false
	}
	return allowed
}

func EnumerableSubtypeIntersect[T any](t1 EnumerableSubtype[T], t2 EnumerableSubtype[T], result *[]EnumerableType[T]) bool {
	// migrated from EnumerableSubtype.java:59:5
	b1 := t1.Allowed()
	b2 := t2.Allowed()
	var allowed bool
	if b1 && b2 {
		EnumerableListIntersect(t1.Values(), t2.Values(), result)
		allowed = true
	} else if (!b1) && (!b2) {
		EnumerableListUnion(t1.Values(), t2.Values(), result)
		allowed = false
	} else if b1 && (!b2) {
		EnumerableListDiff(t1.Values(), t2.Values(), result)
		allowed = true
	} else {
		EnumerableListDiff(t2.Values(), t1.Values(), result)
		allowed = true
	}
	return allowed
}

func EnumerableListUnion[T any](v1 []EnumerableType[T], v2 []EnumerableType[T], result *[]EnumerableType[T]) {
	// migrated from EnumerableSubtype.java:81:5
	i1 := 0
	i2 := 0
	len1 := len(v1)
	len2 := len(v2)
	for true {
		if i1 >= len1 {
			if i2 >= len2 {
				break
			}
			*result = append(*result, v2[i2])
			i2 = (i2 + 1)
		} else if i2 >= len2 {
			*result = append(*result, v1[i1])
			i1 = (i1 + 1)
		} else {
			s1 := v1[i1]
			s2 := v2[i2]
			switch CompareEnumerable(s1, s2) {
			case EQ:
				*result = append(*result, s1)
				i1 = (i1 + 1)
				i2 = (i2 + 1)
				break
			case LT:
				*result = append(*result, s1)
				i1 = (i1 + 1)
				break
			case GT:
				*result = append(*result, s2)
				i2 = (i2 + 1)
				break
			}
		}
	}
}

func EnumerableListIntersect[T any](v1 []EnumerableType[T], v2 []EnumerableType[T], result *[]EnumerableType[T]) {
	// migrated from EnumerableSubtype.java:121:5
	i1 := 0
	i2 := 0
	len1 := len(v1)
	len2 := len(v2)
	for true {
		if (i1 >= len1) || (i2 >= len2) {
			break
		} else {
			s1 := v1[i1]
			s2 := v2[i2]
			switch CompareEnumerable(s1, s2) {
			case EQ:
				*result = append(*result, s1)
				i1 = (i1 + 1)
				i2 = (i2 + 1)
				break
			case LT:
				i1 = (i1 + 1)
				break
			case GT:
				i2 = (i2 + 1)
				break
			}
		}
	}
}

func EnumerableListDiff[T any](v1 []EnumerableType[T], v2 []EnumerableType[T], result *[]EnumerableType[T]) {
	// migrated from EnumerableSubtype.java:152:5
	i1 := 0
	i2 := 0
	len1 := len(v1)
	len2 := len(v2)
	for true {
		if i1 >= len1 {
			break
		}
		if i2 >= len2 {
			*result = append(*result, v1[i1])
			i1 = (i1 + 1)
		} else {
			s1 := v1[i1]
			s2 := v2[i2]
			switch CompareEnumerable(s1, s2) {
			case EQ:
				i1 = (i1 + 1)
				i2 = (i2 + 1)
				break
			case LT:
				*result = append(*result, s1)
				i1 = (i1 + 1)
				break
			case GT:
				i2 = (i2 + 1)
				break
			}
		}
	}
}

func CompareEnumerable[T any](v1 EnumerableType[T], v2 EnumerableType[T]) int {
	// migrated from EnumerableSubtype.java:187:5
	return v1.Compare(v2)
}

func bFloatEq(f1 float64, f2 float64) bool {
	// migrated from EnumerableSubtype.java:216:5
	if math.IsNaN(f1) {
		return math.IsNaN(f2)
	}
	return (f1 == f2)
}
