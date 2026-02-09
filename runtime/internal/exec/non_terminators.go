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

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/semtypes"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var decimalStringRegex = regexp.MustCompile(`^[+-]?[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)

func execConstantLoad(constantLoad *bir.ConstantLoad, frame *Frame) {
	frame.SetOperand(constantLoad.LhsOp.Index, constantLoad.Value)
}

func execMove(moveIns *bir.Move, frame *Frame) {
	frame.SetOperand(moveIns.LhsOp.Index, frame.GetOperand(moveIns.RhsOp.Index))
}

func execNewArray(newArray *bir.NewArray, frame *Frame) {
	size := 0
	if newArray.SizeOp != nil {
		size = int(frame.GetOperand(newArray.SizeOp.Index).(int64))
	}
	arr := make([]any, size)
	if size > 0 {
		for i, value := range newArray.Values {
			arr[i] = frame.GetOperand(value.Index)
		}
	}
	frame.SetOperand(newArray.LhsOp.Index, &arr)
}

func execArrayStore(access *bir.FieldAccess, frame *Frame) {
	arrPtr := frame.GetOperand(access.LhsOp.Index).(*[]any)
	arr := *arrPtr
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	arr = resizeArrayIfNeeded(arrPtr, arr, idx)
	arr[idx] = frame.GetOperand(access.RhsOp.Index)
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame) {
	arr := *(frame.GetOperand(access.RhsOp.Index).(*[]any))
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 || idx >= len(arr) {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	frame.SetOperand(access.LhsOp.Index, arr[idx])
}

func execTypeCast(typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := frame.GetOperand(typeCast.RhsOp.Index)
	targetType := typeCast.Type
	result := castValue(sourceValue, targetType)
	frame.SetOperand(typeCast.LhsOp.Index, result)
}

func castValue(value any, targetType semtypes.SemType) any {
	b, ok := targetType.(*semtypes.BasicTypeBitSet)
	if !ok {
		panic(fmt.Sprintf("bad type cast: unsupported target type %T", targetType))
	}
	if b.All() == semtypes.ANY.All() {
		return value
	}
	bitsetValue := b.All()
	switch {
	case bitsetValue&semtypes.INT.All() != 0:
		return toInt(value)
	case bitsetValue&semtypes.FLOAT.All() != 0:
		return toFloat(value)
	case bitsetValue&semtypes.DECIMAL.All() != 0:
		return toDecimal(value)
	case bitsetValue&semtypes.BOOLEAN.All() != 0:
		if v, ok := value.(bool); ok {
			return v
		}
		panic(fmt.Sprintf("bad type cast: cannot cast %v to boolean", value))
	}
	panic(fmt.Sprintf("bad type cast: unsupported basic type %s", b.String()))
}

func resizeArrayIfNeeded(arrPtr *[]any, arr []any, idx int) []any {
	if idx >= len(arr) {
		newArr := make([]any, idx+1)
		copy(newArr, arr)
		*arrPtr = newArr
		return newArr
	}
	return arr
}

func toInt(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
		}
		if v < float64(math.MinInt64) || v > float64(math.MaxInt64) {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
		}
		return int64(v)
	case *big.Rat:
		if !v.IsInt() {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
		}
		num := v.Num()
		if num.BitLen() > 63 {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
		}
		return num.Int64()
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to int", value))
	}
}

func toFloat(value any) float64 {
	switch v := value.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case *big.Rat:
		f, _ := v.Float64()
		return f
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to float", value))
	}
}

func toDecimal(value any) *big.Rat {
	switch v := value.(type) {
	case int64:
		return big.NewRat(v, 1)
	case float64:
		r := new(big.Rat)
		if r.SetFloat64(v) == nil {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to decimal", v))
		}
		return r
	case *big.Rat:
		return v
	case string:
		if strings.Contains(v, "/") || !decimalStringRegex.MatchString(v) {
			panic(fmt.Sprintf("cannot cast %v to decimal", v))
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot cast %v to decimal", v))
		}
		r := new(big.Rat)
		r.SetFloat64(f)
		return r
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to decimal", value))
	}
}
