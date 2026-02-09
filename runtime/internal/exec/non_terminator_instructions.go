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
	"math/big"
	"strconv"
)

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
		if size < 0 {
			size = 0
		}
	}
	arr := make([]any, size)
	frame.SetOperand(newArray.LhsOp.Index, &arr)
}

func execArrayStore(access *bir.FieldAccess, frame *Frame) {
	arrPtr := frame.GetOperand(access.LhsOp.Index).(*[]any)
	arr := *arrPtr
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	arr = resizeArrayIfNeeded(arrPtr, arr, idx)
	arr[idx] = frame.GetOperand(access.RhsOp.Index)
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame) {
	arr := *(frame.GetOperand(access.RhsOp.Index).(*[]any))
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	frame.SetOperand(access.LhsOp.Index, arr[idx])
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

func execTypeCast(typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := frame.GetOperand(typeCast.RhsOp.Index)
	targetType := typeCast.Type
	result := convertValue(sourceValue, targetType)
	frame.SetOperand(typeCast.LhsOp.Index, result)
}

func convertValue(value any, targetType semtypes.SemType) any {
	bitset := semtypes.WidenToBasicTypes(targetType)
	bitsetValue := (&bitset).All()
	switch {
	case bitsetValue&semtypes.INT.All() != 0:
		return convertToInt(value)
	case bitsetValue&semtypes.FLOAT.All() != 0:
		return convertToFloat(value)
	case bitsetValue&semtypes.DECIMAL.All() != 0:
		return convertToDecimal(value)
	}
	return value
}

func convertToInt(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case *big.Rat:
		f, _ := v.Float64()
		return int64(f)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("unsupported type conversion to int: value %v", v))
		}
		return int64(f)
	default:
		panic(fmt.Sprintf("unsupported type conversion to int: value %v", value))
	}
}

func convertToFloat(value any) float64 {
	switch v := value.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case *big.Rat:
		f, _ := v.Float64()
		return f
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("unsupported type conversion to float: value %v", v))
		}
		return f
	default:
		panic(fmt.Sprintf("unsupported type conversion to float: value %v", value))
	}
}

func convertToDecimal(value any) *big.Rat {
	switch v := value.(type) {
	case int64:
		return big.NewRat(v, 1)
	case float64:
		return new(big.Rat).SetFloat64(v)
	case *big.Rat:
		return v
	case string:
		r := new(big.Rat)
		if _, ok := r.SetString(v); !ok {
			panic(fmt.Sprintf("unsupported type conversion to decimal: value %v", v))
		}
		return r
	default:
		panic(fmt.Sprintf("unsupported type conversion to decimal: value %v", value))
	}
}
