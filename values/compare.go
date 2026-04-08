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
