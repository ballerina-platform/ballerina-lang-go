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
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

func execConstantLoad(ctx *Context, constantLoad *bir.ConstantLoad, frame *Frame) {
	setOperandValue(ctx, constantLoad.LhsOp, frame, constantLoad.Value)
}

func execMove(ctx *Context, moveIns *bir.Move, frame *Frame) {
	setOperandValue(ctx, moveIns.LhsOp, frame, getOperandValue(ctx, moveIns.RhsOp, frame))
}

func execNewArray(ctx *Context, newArray *bir.NewArray, frame *Frame) {
	size := 0
	if newArray.SizeOp != nil {
		size = int(getOperandValue(ctx, newArray.SizeOp, frame).(int64))
	}
	list := values.NewList(size, newArray.Type, newArray.Filler)
	for i, value := range newArray.Values {
		list.FillingSet(i, getOperandValue(ctx, value, frame))
	}
	setOperandValue(ctx, newArray.LhsOp, frame, list)
}

func execNewMap(ctx *Context, newMap *bir.NewMap, frame *Frame) {
	m := values.NewMap(newMap.Type)
	for _, entry := range newMap.Values {
		kv := entry.(*bir.MappingConstructorKeyValueEntry)
		keyVal := getOperandValue(ctx, kv.KeyOp(), frame)
		keyStr := keyVal.(string)
		valueVal := getOperandValue(ctx, kv.ValueOp(), frame)
		m.Put(keyStr, valueVal)
	}
	for _, def := range newMap.Defaults {
		if _, exists := m.Get(def.FieldName); !exists {
			fn := ctx.GetBIRFunction(def.FunctionLookupKey)
			val := executeFunction(ctx, *fn, nil, frame)
			m.Put(def.FieldName, val)
		}
	}
	setOperandValue(ctx, newMap.GetLhsOperand(), frame, m)
}

func execNewError(ctx *Context, newError *bir.NewError, frame *Frame) {
	msgVal := getOperandValue(ctx, newError.MessageOp, frame)
	message := msgVal.(string)

	var cause values.BalValue
	if newError.CauseOp != nil {
		cause = getOperandValue(ctx, newError.CauseOp, frame)
	}

	var detailMap *values.Map
	if newError.DetailOp != nil {
		detailMap = getOperandValue(ctx, newError.DetailOp, frame).(*values.Map)
	}
	errVal := values.NewError(newError.Type, message, cause, newError.TypeName, detailMap)
	setOperandValue(ctx, newError.GetLhsOperand(), frame, errVal)
}

func execNewObject(ctx *Context, newObject *bir.NewObject, frame *Frame) {
	classDef := ctx.GetClassDef(newObject.ClassDefRef)
	fieldValues := make(map[string]values.BalValue, len(classDef.Fields))
	for _, field := range classDef.Fields {
		fieldValues[field.Name] = values.DefaultValueForType(field.Ty)
	}
	methodKeys := make(map[string]string, len(classDef.VTable))
	for methodName, method := range classDef.VTable {
		methodKeys[methodName] = method.FunctionLookupKey
	}
	objType := newObject.GetLhsOperand().VariableDcl.GetType()
	obj := values.NewObject(objType, fieldValues, methodKeys)
	setOperandValue(ctx, newObject.GetLhsOperand(), frame, obj)
}

func execArrayStore(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	list := getOperandValue(ctx, access.LhsOp, frame).(*values.List)
	idx := int(getOperandValue(ctx, access.KeyOp, frame).(int64))
	if idx < 0 {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid array index: %d", idx)))
	}
	list.FillingSet(idx, getOperandValue(ctx, access.RhsOp, frame))
}

func execArrayLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	list := getOperandValue(ctx, access.RhsOp, frame).(*values.List)
	idx := int(getOperandValue(ctx, access.KeyOp, frame).(int64))
	if idx < 0 || idx >= list.Len() {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid array index: %d", idx)))
	}
	setOperandValue(ctx, access.LhsOp, frame, list.Get(idx))
}

func execMapStore(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	m := getOperandValue(ctx, access.LhsOp, frame).(*values.Map)
	keyVal := getOperandValue(ctx, access.KeyOp, frame)
	keyStr := keyVal.(string)
	valueVal := getOperandValue(ctx, access.RhsOp, frame)
	m.Put(keyStr, valueVal)
}

func execMapLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	m := getOperandValue(ctx, access.RhsOp, frame).(*values.Map)
	key := getOperandValue(ctx, access.KeyOp, frame).(string)
	value, _ := m.Get(key)
	setOperandValue(ctx, access.LhsOp, frame, value)
}

func execObjectStore(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	obj := getOperandValue(ctx, access.LhsOp, frame).(*values.Object)
	field := getOperandValue(ctx, access.KeyOp, frame).(string)
	value := getOperandValue(ctx, access.RhsOp, frame)
	obj.Put(field, value)
}

func execObjectLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	obj := getOperandValue(ctx, access.RhsOp, frame).(*values.Object)
	field := getOperandValue(ctx, access.KeyOp, frame).(string)
	value, _ := obj.Get(field)
	setOperandValue(ctx, access.LhsOp, frame, value)
}

func execTypeCast(ctx *Context, typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := getOperandValue(ctx, typeCast.RhsOp, frame)
	result := castValue(ctx, sourceValue, typeCast.Type)
	setOperandValue(ctx, typeCast.LhsOp, frame, result)
}

func execFPLoad(ctx *Context, fpLoad *bir.FPLoad, frame *Frame) {
	fn := &values.Function{
		Type:      fpLoad.Type,
		LookupKey: fpLoad.FunctionLookupKey,
	}
	if fpLoad.IsClosure {
		fn.ParentFrame = frame
	}
	setOperandValue(ctx, fpLoad.LhsOp, frame, fn)
}

func execTypeTest(ctx *Context, typeTest *bir.TypeTest, frame *Frame) {
	sourceValue := getOperandValue(ctx, typeTest.RhsOp, frame)
	valueType := values.SemTypeForValue(sourceValue)
	typeCtx := ctx.TypeCheckContext()
	matches := semtypes.IsSubtype(typeCtx, valueType, typeTest.Type) != typeTest.IsNegation
	setOperandValue(ctx, typeTest.LhsOp, frame, matches)
}

func castValue(ctx *Context, value values.BalValue, targetType semtypes.SemType) values.BalValue {
	typeCtx := ctx.TypeCheckContext()
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
	panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast value of type %s to %s",
		semtypes.ToString(typeCtx, valueType), semtypes.ToString(typeCtx, targetType))))
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
