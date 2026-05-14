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

package values

import (
	"fmt"
	"math"
	"math/big"
	"unsafe"
)

type CompareResult int8

const (
	CmpLT CompareResult = -1
	CmpEQ CompareResult = 0
	CmpGT CompareResult = 1
	CmpUN CompareResult = 2
)

func Compare(x, y BalValue) CompareResult {
	if x == nil && y == nil {
		return CmpEQ
	}
	if x == nil || y == nil {
		return CmpUN
	}
	switch v1 := x.(type) {
	case int64:
		v2 := y.(int64)
		if v1 < v2 {
			return CmpLT
		}
		if v1 > v2 {
			return CmpGT
		}
		return CmpEQ
	case float64:
		v2 := y.(float64)
		if v1 < v2 {
			return CmpLT
		}
		if v1 > v2 {
			return CmpGT
		}
		if v1 == v2 {
			return CmpEQ
		}
		return CmpUN
	case bool:
		v2 := y.(bool)
		if v1 == v2 {
			return CmpEQ
		}
		if !v1 && v2 {
			return CmpLT
		}
		return CmpGT
	case string:
		v2 := y.(string)
		if v1 < v2 {
			return CmpLT
		}
		if v1 > v2 {
			return CmpGT
		}
		return CmpEQ
	case *big.Rat:
		v2 := y.(*big.Rat)
		switch v1.Cmp(v2) {
		case -1:
			return CmpLT
		case 0:
			return CmpEQ
		default:
			return CmpGT
		}
	case *List:
		v2 := y.(*List)
		return compareList(v1, v2)
	default:
		panic(NewErrorWithMessage(fmt.Sprintf("unsupported type for comparison: %T", x)))
	}
}

func CompareA(x, y BalValue) CompareResult {
	return compareForSort(x, y, true)
}

func CompareD(x, y BalValue) CompareResult {
	return compareForSort(x, y, false)
}

func CompareK(x, y BalValue, ascending bool) CompareResult {
	if ascending {
		return CompareA(x, y)
	}
	return reverseCompareResult(CompareD(x, y))
}

func compareForSort(x, y BalValue, ascending bool) CompareResult {
	if x == nil {
		if y == nil {
			return CmpEQ
		}
		if ascending {
			return CmpGT
		}
		return CmpLT
	}
	if y == nil {
		if ascending {
			return CmpLT
		}
		return CmpGT
	}
	if xList, ok := x.(*List); ok {
		yList := y.(*List)
		return compareListForSort(xList, yList, ascending)
	}
	if xFloat, ok := x.(float64); ok {
		yFloat := y.(float64)
		xNaN := math.IsNaN(xFloat)
		yNaN := math.IsNaN(yFloat)
		switch {
		case xNaN && yNaN:
			return CmpEQ
		case xNaN:
			if ascending {
				return CmpGT
			}
			return CmpLT
		case yNaN:
			if ascending {
				return CmpLT
			}
			return CmpGT
		}
	}

	r := Compare(x, y)
	if r == CmpUN {
		panic(NewErrorWithMessage(fmt.Sprintf("unsupported type for comparison: %T and %T", x, y)))
	}
	return r
}

func compareList(x, y *List) CompareResult {
	xLen := x.Len()
	yLen := y.Len()
	minLen := min(yLen, xLen)
	for i := range minLen {
		r := Compare(x.Get(i), y.Get(i))
		if r != CmpEQ {
			return r
		}
	}
	if xLen < yLen {
		return CmpLT
	}
	if xLen > yLen {
		return CmpGT
	}
	return CmpEQ
}

func compareListForSort(x, y *List, ascending bool) CompareResult {
	xLen := x.Len()
	yLen := y.Len()
	minLen := min(yLen, xLen)
	for i := range minLen {
		r := compareForSort(x.Get(i), y.Get(i), ascending)
		if r != CmpEQ {
			return r
		}
	}
	if xLen < yLen {
		return CmpLT
	}
	if xLen > yLen {
		return CmpGT
	}
	return CmpEQ
}

func reverseCompareResult(cmp CompareResult) CompareResult {
	switch cmp {
	case CmpLT:
		return CmpGT
	case CmpGT:
		return CmpLT
	default:
		return cmp
	}
}

type deepEqualVisit struct {
	left  uintptr
	right uintptr
}

func DeepEqual(x, y BalValue) bool {
	return deepEqual(x, y, make(map[deepEqualVisit]bool))
}

func deepEqual(x, y BalValue, visited map[deepEqualVisit]bool) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	switch left := x.(type) {
	case bool:
		right, ok := y.(bool)
		return ok && left == right
	case int64:
		right, ok := y.(int64)
		return ok && left == right
	case float64:
		right, ok := y.(float64)
		return ok && (left == right || (math.IsNaN(left) && math.IsNaN(right)))
	case string:
		right, ok := y.(string)
		return ok && left == right
	case *big.Rat:
		right, ok := y.(*big.Rat)
		return ok && left.Cmp(right) == 0
	case *List:
		right, ok := y.(*List)
		if !ok {
			return false
		}
		return deepEqualLists(left, right, visited)
	case *Map:
		right, ok := y.(*Map)
		if !ok {
			return false
		}
		return deepEqualMaps(left, right, visited)
	default:
		return false
	}
}

func deepEqualLists(left, right *List, visited map[deepEqualVisit]bool) bool {
	if left == right {
		return true
	}
	visit := deepEqualVisit{
		left:  uintptr(unsafe.Pointer(left)),
		right: uintptr(unsafe.Pointer(right)),
	}
	if visited[visit] {
		return true
	}
	visited[visit] = true
	defer delete(visited, visit)

	if left.Len() != right.Len() {
		return false
	}
	for i := 0; i < left.Len(); i++ {
		if !deepEqual(left.Get(i), right.Get(i), visited) {
			return false
		}
	}
	return true
}

func deepEqualMaps(left, right *Map, visited map[deepEqualVisit]bool) bool {
	if left == right {
		return true
	}
	visit := deepEqualVisit{
		left:  uintptr(unsafe.Pointer(left)),
		right: uintptr(unsafe.Pointer(right)),
	}
	if visited[visit] {
		return true
	}
	visited[visit] = true
	defer delete(visited, visit)

	if left.Len() != right.Len() {
		return false
	}
	for _, key := range left.Keys() {
		leftValue, _ := left.Get(key)
		rightValue, ok := right.Get(key)
		if !ok || !deepEqual(leftValue, rightValue, visited) {
			return false
		}
	}
	return true
}
