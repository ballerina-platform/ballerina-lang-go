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
	"reflect"
)

type CompareResult int8

const (
	CmpLT CompareResult = -1
	CmpEQ CompareResult = 0
	CmpGT CompareResult = 1
	CmpUN CompareResult = 2
)

func Compare(x, y BalValue) CompareResult {
	if x == y {
		return CmpEQ
	}
	if x == nil || y == nil {
		return CmpUN
	}
	switch v1 := x.(type) {
	case int64:
		if v1 < y.(int64) {
			return CmpLT
		}
		return CmpGT
	case float64:
		return compareFloat(v1, y.(float64))
	case bool:
		if !v1 && y.(bool) {
			return CmpLT
		}
		return CmpGT
	case string:
		if v1 < y.(string) {
			return CmpLT
		}
		return CmpGT
	case *big.Rat:
		switch v1.Cmp(y.(*big.Rat)) {
		case -1:
			return CmpLT
		case 0:
			return CmpEQ
		default:
			return CmpGT
		}
	case *List:
		return compareListValues(v1, y.(*List))
	default:
		panic(NewErrorWithMessage(fmt.Sprintf("unsupported type for comparison: %T", x)))
	}
}

func CompareRef(x, y BalValue) CompareResult {
	if x == nil && y == nil {
		return CmpEQ
	}
	if x == nil || y == nil {
		return CmpUN
	}
	switch xv := x.(type) {
	case *Function:
		yFn, ok := y.(*Function)
		if !ok || xv.LookupKey != yFn.LookupKey {
			return CmpUN
		}
		return CmpEQ
	case float64:
		yf, ok := y.(float64)
		if !ok || !compareFloatRef(xv, yf) {
			return CmpUN
		}
		return CmpEQ
	case *big.Rat:
		yDec, ok := y.(*big.Rat)
		if !ok {
			return CmpUN
		}
		if xv.Cmp(yDec) == 0 {
			return CmpEQ
		}
		return CmpUN
	}
	if x == y {
		return CmpEQ
	}
	return CmpUN
}

func Equal(x, y BalValue) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	if reflect.TypeOf(x) != reflect.TypeOf(y) {
		return false
	}
	if xf, ok := x.(float64); ok {
		yf := y.(float64)
		if math.IsNaN(xf) && math.IsNaN(yf) {
			return true
		}
	}
	return Compare(x, y) == CmpEQ
}

func compareFloatRef(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Float64bits(a) == math.Float64bits(b)
}

func compareFloat(a, b float64) CompareResult {
	if math.IsNaN(a) || math.IsNaN(b) {
		return CmpUN
	}
	if a == b {
		return CmpEQ
	}
	if a < b {
		return CmpLT
	}
	if a > b {
		return CmpGT
	}
	return CmpUN
}

func compareListValues(x, y *List) CompareResult {
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
