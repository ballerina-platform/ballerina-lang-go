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
	"fmt"
	"math"
	"math/big"
	"strconv"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

func execConstantLoad(constantLoad *bir.ConstantLoad, frame *Frame, reg *modules.Registry) {
	setOperandValue(constantLoad.LhsOp, frame, reg, constantLoad.Value)
}

func execMove(moveIns *bir.Move, frame *Frame, reg *modules.Registry) {
	setOperandValue(moveIns.LhsOp, frame, reg, getOperandValue(moveIns.RhsOp, frame, reg))
}

func execNewArray(newArray *bir.NewArray, frame *Frame, reg *modules.Registry) {
	size := 0
	if newArray.SizeOp != nil {
		size = int(getOperandValue(newArray.SizeOp, frame, reg).(int64))
	}
	list := values.NewList(size, newArray.Type, newArray.Filler)
	for i, value := range newArray.Values {
		list.FillingSet(i, getOperandValue(value, frame, reg))
	}
	setOperandValue(newArray.LhsOp, frame, reg, list)
}

func execNewMap(newMap *bir.NewMap, frame *Frame, reg *modules.Registry) {
	m := values.NewMap(newMap.Type)
	for _, entry := range newMap.Values {
		kv := entry.(*bir.MappingConstructorKeyValueEntry)
		keyVal := getOperandValue(kv.KeyOp(), frame, reg)
		keyStr := keyVal.(string)
		valueVal := getOperandValue(kv.ValueOp(), frame, reg)
		m.Put(keyStr, valueVal)
	}
	setOperandValue(newMap.GetLhsOperand(), frame, reg, m)
}

func execNewError(newError *bir.NewError, frame *Frame, reg *modules.Registry) {
	msgVal := getOperandValue(newError.MessageOp, frame, reg)
	message := msgVal.(string)

	var cause values.BalValue
	if newError.CauseOp != nil {
		cause = getOperandValue(newError.CauseOp, frame, reg)
	}

	var detailMap *values.Map
	if newError.DetailOp != nil {
		detailMap = getOperandValue(newError.DetailOp, frame, reg).(*values.Map)
	}
	errVal := values.NewError(newError.Type, message, cause, newError.TypeName, detailMap)
	setOperandValue(newError.GetLhsOperand(), frame, reg, errVal)
}

func execArrayStore(access *bir.FieldAccess, frame *Frame, reg *modules.Registry) {
	list := getOperandValue(access.LhsOp, frame, reg).(*values.List)
	idx := int(getOperandValue(access.KeyOp, frame, reg).(int64))
	if idx < 0 {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid array index: %d", idx)))
	}
	list.FillingSet(idx, getOperandValue(access.RhsOp, frame, reg))
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame, reg *modules.Registry) {
	list := getOperandValue(access.RhsOp, frame, reg).(*values.List)
	idx := int(getOperandValue(access.KeyOp, frame, reg).(int64))
	if idx < 0 || idx >= list.Len() {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid array index: %d", idx)))
	}
	setOperandValue(access.LhsOp, frame, reg, list.Get(idx))
}

func execMapStore(access *bir.FieldAccess, frame *Frame, reg *modules.Registry) {
	m := getOperandValue(access.LhsOp, frame, reg).(*values.Map)
	keyVal := getOperandValue(access.KeyOp, frame, reg)
	keyStr := keyVal.(string)
	valueVal := getOperandValue(access.RhsOp, frame, reg)
	m.Put(keyStr, valueVal)
}

func execMapLoad(access *bir.FieldAccess, frame *Frame, reg *modules.Registry) {
	m := getOperandValue(access.RhsOp, frame, reg).(*values.Map)
	key := getOperandValue(access.KeyOp, frame, reg).(string)
	value, _ := m.Get(key)
	setOperandValue(access.LhsOp, frame, reg, value)
}

func execTypeCast(typeCast *bir.TypeCast, frame *Frame, reg *modules.Registry) {
	sourceValue := getOperandValue(typeCast.RhsOp, frame, reg)
	result := castValue(sourceValue, typeCast.Type, reg)
	setOperandValue(typeCast.LhsOp, frame, reg, result)
}

func execFPLoad(fpLoad *bir.FPLoad, frame *Frame, reg *modules.Registry) {
	fn := &values.Function{
		Type:      fpLoad.Type,
		LookupKey: fpLoad.FunctionLookupKey,
	}
	setOperandValue(fpLoad.LhsOp, frame, reg, fn)
}

func execTypeTest(typeTest *bir.TypeTest, frame *Frame, reg *modules.Registry) {
	sourceValue := getOperandValue(typeTest.RhsOp, frame, reg)
	valueType := values.SemTypeForValue(sourceValue)
	typeEnv := reg.GetTypeEnv()
	typeCtx := semtypes.TypeCheckContext(typeEnv)
	matches := semtypes.IsSubtype(typeCtx, valueType, typeTest.Type) != typeTest.IsNegation
	setOperandValue(typeTest.LhsOp, frame, reg, matches)
}

func castValue(value values.BalValue, targetType semtypes.SemType, reg *modules.Registry) values.BalValue {
	typeEnv := reg.GetTypeEnv()
	typeCtx := semtypes.TypeCheckContext(typeEnv)
	valueType := values.SemTypeForValue(value)
	if semtypes.IsSubtype(typeCtx, valueType, targetType) {
		return value
	}
	switch {
	case semtypes.IsSubtypeSimple(targetType, semtypes.INT):
		return toInt(value)
	case semtypes.IsSubtypeSimple(targetType, semtypes.FLOAT):
		return toFloat(value)
	case semtypes.IsSubtypeSimple(targetType, semtypes.DECIMAL):
		return toDecimal(value)
	}
	panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: unsupported target type %s", semtypes.ToString(typeCtx, targetType))))
}

func toInt(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast non-finite value %v to int", v)))
		}
		if v < float64(math.MinInt64) || v > float64(math.MaxInt64) {
			panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast out-of-range value %v to int", v)))
		}
		return int64(v)
	case *big.Rat:
		return decimalToInt(v)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to int", value)))
	}
}

func decimalToInt(v *big.Rat) int64 {
	num := v.Num()
	denom := v.Denom()
	q, r := new(big.Int).QuoRem(num, denom, new(big.Int))
	if r.Sign() != 0 && new(big.Int).Mul(new(big.Int).Abs(r), big.NewInt(2)).Cmp(denom) >= 0 {
		if num.Sign() >= 0 {
			q.Add(q, big.NewInt(1))
		} else {
			q.Sub(q, big.NewInt(1))
		}
	}
	if !q.IsInt64() {
		panic(values.NewErrorWithMessage(fmt.Sprintf("cannot convert %v to int64: value out of range", v)))
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
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to float", value)))
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
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to decimal", value)))
	}
}

func toBoolean(value any) bool {
	if v, ok := value.(bool); ok {
		return v
	}
	panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to boolean", value)))
}
