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

func execTypeCast(typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := frame.GetOperand(typeCast.RhsOp.Index)
	targetType := typeCast.Type
	result := castValue(sourceValue, targetType)
	frame.SetOperand(typeCast.LhsOp.Index, result)
}

func castValue(value any, targetType semtypes.SemType) any {
	b := targetType.(*semtypes.BasicTypeBitSet)
	// If casting to any, just return the value as-is
	if b.All() == semtypes.ANY.All() {
		return value
	}
	bitsetValue := b.All()
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

func resizeArrayIfNeeded(arrPtr *[]any, arr []any, idx int) []any {
	if idx < len(arr) {
		return arr
	}

	newLen := idx + 1
	newCap := newLen
	if cap(arr)*2 > newCap {
		newCap = cap(arr) * 2
	}

	newArr := make([]any, newLen, newCap)
	copy(newArr, arr)
	*arrPtr = newArr
	return newArr
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
			panic("bad type cast")
		}
		return int64(f)
	default:
		panic("bad type cast")
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
			panic("bad type cast")
		}
		return f
	default:
		panic("bad type cast")
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
			panic("bad type cast")
		}
		return r
	default:
		panic("bad type cast")
	}
}
