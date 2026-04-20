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

func Equal(x, y BalValue) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	if x == y {
		return true
	}
	return equalValue(x, y)
}

func RefEqual(x, y BalValue) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	switch xv := x.(type) {
	case *Function:
		yFn, ok := y.(*Function)
		return ok && xv.LookupKey == yFn.LookupKey
	case float64:
		yf, ok := y.(float64)
		return ok && compareFloatRef(xv, yf)
	case *big.Rat:
		yDec, ok := y.(*big.Rat)
		return ok && xv.Cmp(yDec) == 0
	default:
		return x == y
	}
}

func equalValue(x, y BalValue) bool {
	switch v1 := x.(type) {
	case int64:
		v2, ok := y.(int64)
		return ok && v1 == v2
	case float64:
		v2, ok := y.(float64)
		return ok && (v1 == v2 || (math.IsNaN(v1) && math.IsNaN(v2)))
	case bool:
		v2, ok := y.(bool)
		return ok && v1 == v2
	case string:
		v2, ok := y.(string)
		return ok && v1 == v2
	case *big.Rat:
		v2, ok := y.(*big.Rat)
		return ok && v1.Cmp(v2) == 0
	case *List:
		v2, ok := y.(*List)
		return ok && compareListValues(v1, v2) == CmpEQ
	default:
		panic(NewErrorWithMessage(fmt.Sprintf("unsupported type for comparison: %T", x)))
	}
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

func compareFloat(a, b float64) CompareResult {
	if math.IsNaN(a) || math.IsNaN(b) {
		return CmpUN
	}
	if a < b {
		return CmpLT
	}
	return CmpGT
}

func compareFloatRef(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Float64bits(a) == math.Float64bits(b)
}
