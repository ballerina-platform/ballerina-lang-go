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
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
	"fmt"
	"math"
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
	}
	list := values.NewList(size, newArray.Type, newArray.Filler)
	for i, value := range newArray.Values {
		list.FillingSet(i, frame.GetOperand(value.Index))
	}
	frame.SetOperand(newArray.LhsOp.Index, list)
}

func execNewMap(newMap *bir.NewMap, frame *Frame) {
	m := values.NewMap(newMap.Type)
	for _, entry := range newMap.Values {
		kv := entry.(*bir.MappingConstructorKeyValueEntry)
		keyVal := frame.GetOperand(kv.KeyOp().Index)
		keyStr := keyVal.(string)
		valueVal := frame.GetOperand(kv.ValueOp().Index)
		m.Put(keyStr, valueVal)
	}
	frame.SetOperand(newMap.GetLhsOperand().Index, m)
}

func execNewError(newError *bir.NewError, frame *Frame) {
	msgVal := frame.GetOperand(newError.MessageOp.Index)
	message := msgVal.(string)

	var cause values.BalValue
	if newError.CauseOp != nil {
		cause = frame.GetOperand(newError.CauseOp.Index)
	}

	var detailMap *values.Map
	if newError.DetailOp != nil {
		detailMap = frame.GetOperand(newError.DetailOp.Index).(*values.Map)
	}
	errVal := values.NewError(newError.Type, message, cause, newError.TypeName, detailMap)
	frame.SetOperand(newError.GetLhsOperand().Index, errVal)
}

func execArrayStore(access *bir.FieldAccess, frame *Frame) {
	list := frame.GetOperand(access.LhsOp.Index).(*values.List)
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	list.FillingSet(idx, frame.GetOperand(access.RhsOp.Index))
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame) {
	list := frame.GetOperand(access.RhsOp.Index).(*values.List)
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 || idx >= list.Len() {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	frame.SetOperand(access.LhsOp.Index, list.Get(idx))
}

func execMapStore(access *bir.FieldAccess, frame *Frame) {
	m := frame.GetOperand(access.LhsOp.Index).(*values.Map)
	keyVal := frame.GetOperand(access.KeyOp.Index)
	keyStr := keyVal.(string)
	valueVal := frame.GetOperand(access.RhsOp.Index)
	m.Put(keyStr, valueVal)
}

func execMapLoad(access *bir.FieldAccess, frame *Frame) {
	m := frame.GetOperand(access.RhsOp.Index).(*values.Map)
	key := frame.GetOperand(access.KeyOp.Index).(string)
	value, _ := m.Get(key)
	frame.SetOperand(access.LhsOp.Index, value)
}

func execTypeCast(typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := frame.GetOperand(typeCast.RhsOp.Index)
	result := castValue(sourceValue, typeCast.Type)
	frame.SetOperand(typeCast.LhsOp.Index, result)
}

func execTypeTest(typeTest *bir.TypeTest, frame *Frame, reg *modules.Registry) {
	sourceValue := frame.GetOperand(typeTest.RhsOp.Index)
	valueType := values.SemTypeForValue(sourceValue)
	typeEnv := reg.GetTypeEnv()
	typeCtx := semtypes.TypeCheckContext(typeEnv)
	matches := semtypes.IsSubtype(typeCtx, valueType, typeTest.Type) != typeTest.IsNegation
	frame.SetOperand(typeTest.LhsOp.Index, matches)
}

func castValue(value values.BalValue, targetType semtypes.SemType) values.BalValue {
	b := targetType.(*semtypes.BasicTypeBitSet)
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
		return value.(bool)
	}
	panic(fmt.Sprintf("bad type cast: unsupported basic type %s", b.String()))
}

func toInt(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		if v < float64(math.MinInt64) || v > float64(math.MaxInt64) {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
		}
		return int64(v)
	case *big.Rat:
		return decimalToInt(v)
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to int", value))
	}
}

func decimalToInt(v *big.Rat) int64 {
	num := v.Num()
	denom := v.Denom()
	q := new(big.Int).Quo(num, denom)
	r := new(big.Int).Rem(new(big.Int).Set(num), denom)
	rAbs := new(big.Int).Abs(r)
	rAbs.Mul(rAbs, big.NewInt(2))
	cmp := rAbs.Cmp(denom)
	isHalf := cmp == 0
	roundUp := cmp > 0
	if isHalf {
		if q.Bit(0) == 0 {
			roundUp = false
		} else {
			roundUp = true
			if num.Sign() < 0 {
				q.Sub(q, big.NewInt(1))
			} else {
				q.Add(q, big.NewInt(1))
			}
			if q.BitLen() > 63 {
				panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
			}
			return q.Int64()
		}
	}
	if roundUp {
		if num.Sign() >= 0 {
			q.Add(q, big.NewInt(1))
		} else {
			q.Sub(q, big.NewInt(1))
		}
	}
	if q.BitLen() > 63 {
		panic(fmt.Sprintf("bad type cast: cannot cast %v to int", v))
	}
	return q.Int64()
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
		s := strconv.FormatFloat(v, 'g', -1, 64)
		r.SetString(s)
		return r
	case *big.Rat:
		return v
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to decimal", value))
	}
}
