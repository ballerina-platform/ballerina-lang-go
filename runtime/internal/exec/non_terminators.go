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

	"ballerina-lang-go/bir"
	"ballerina-lang-go/decimal"
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
		fv, _ := values.FillerValue(ctx.TypeCheckContext(), field.Ty)
		fieldValues[field.Name] = fv
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

func execArrayFillingLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	list := getOperandValue(ctx, access.RhsOp, frame).(*values.List)
	idx := int(getOperandValue(ctx, access.KeyOp, frame).(int64))
	if idx < 0 {
		panic(values.NewErrorWithMessage(fmt.Sprintf("invalid array index: %d", idx)))
	}
	setOperandValue(ctx, access.LhsOp, frame, list.FillingGet(idx))
}

func execMapStore(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	container := getOperandValue(ctx, access.LhsOp, frame)
	keyStr := getOperandValue(ctx, access.KeyOp, frame).(string)
	if container == nil {
		panic(values.NewErrorWithMessage(fmt.Sprintf("missing key: %q", keyStr)))
	}
	m := container.(*values.Map)
	valueVal := getOperandValue(ctx, access.RhsOp, frame)
	if valueVal == nil && m.ShouldDeleteOnNilStore(ctx.TypeCheckContext(), keyStr) {
		m.Delete(keyStr)
		return
	}
	m.Put(keyStr, valueVal)
}

func execMapFillingLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	container := getOperandValue(ctx, access.RhsOp, frame)
	key := getOperandValue(ctx, access.KeyOp, frame).(string)
	if container == nil {
		panic(values.NewErrorWithMessage(fmt.Sprintf("missing key: %q", key)))
	}
	setOperandValue(ctx, access.LhsOp, frame, container.(*values.Map).FillingGet(key, access.Filler))
}

func execMapLoad(ctx *Context, access *bir.FieldAccess, frame *Frame) {
	container := getOperandValue(ctx, access.RhsOp, frame)
	key := getOperandValue(ctx, access.KeyOp, frame).(string)
	if container == nil {
		setOperandValue(ctx, access.LhsOp, frame, nil)
		return
	}
	value, _ := container.(*values.Map).Get(key)
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
		return int64(math.RoundToEven(v))
	case *decimal.Decimal:
		return decimalToInt(v)
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to int", value)))
	}
}

func decimalToInt(v *decimal.Decimal) int64 {
	n, ok, err := v.Int64()
	if err != nil {
		panic(values.NewErrorWithMessage(fmt.Sprintf("cannot convert %v to int: %v", v, err)))
	}
	if !ok {
		panic(values.NewErrorWithMessage(fmt.Sprintf("cannot convert %v to int64: value out of range", v)))
	}
	return n
}

func toFloat(value any) float64 {
	switch v := value.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	case *decimal.Decimal:
		return v.Float64()
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to float", value)))
	}
}

// floatToDecimal converts an IEEE 754 float64 into a Ballerina decimal.
// Ballerina decimals do not support NaN, infinities, or subnormals, so any
// such input triggers a runtime panic with the spec-mandated message.
func floatToDecimal(v float64) *decimal.Decimal {
	d, err := decimal.FromFloat64(v)
	if err != nil {
		panic(values.NewErrorWithMessage(err.Error()))
	}
	return d
}

func toDecimal(value any) *decimal.Decimal {
	switch v := value.(type) {
	case int64:
		return decimal.FromInt64(v)
	case float64:
		return floatToDecimal(v)
	case *decimal.Decimal:
		return v
	default:
		panic(values.NewErrorWithMessage(fmt.Sprintf("bad type cast: cannot cast %v to decimal", value)))
	}
}

func execNewXMLText(ctx *Context, instr *bir.NewXMLText, frame *Frame) {
	body := getOperandValue(ctx, instr.BodyOp, frame).(string)
	setOperandValue(ctx, instr.LhsOp, frame, &values.XMLText{Body: body})
}

func execNewXMLComment(ctx *Context, instr *bir.NewXMLComment, frame *Frame) {
	body := getOperandValue(ctx, instr.BodyOp, frame).(string)
	setOperandValue(ctx, instr.LhsOp, frame, &values.XMLComment{Body: body})
}

func execNewXMLPI(ctx *Context, instr *bir.NewXMLPI, frame *Frame) {
	target := getOperandValue(ctx, instr.TargetOp, frame).(string)
	data := getOperandValue(ctx, instr.DataOp, frame).(string)
	setOperandValue(ctx, instr.LhsOp, frame, &values.XMLProcessingInstruction{Target: target, Data: data})
}

func execNewXMLElement(ctx *Context, instr *bir.NewXMLElement, frame *Frame) {
	name := getOperandValue(ctx, instr.NameOp, frame).(string)
	var children values.XMLValue
	if instr.ChildrenOp != nil {
		if v, ok := getOperandValue(ctx, instr.ChildrenOp, frame).(values.XMLValue); ok {
			children = v
		}
	}
	var attrs *values.Map
	if instr.AttrsOp != nil {
		attrs = getOperandValue(ctx, instr.AttrsOp, frame).(*values.Map)
	}
	var namespaces *values.Map
	if instr.NamespacesOp != nil {
		namespaces = getOperandValue(ctx, instr.NamespacesOp, frame).(*values.Map)
	}
	setOperandValue(ctx, instr.LhsOp, frame, &values.XMLElement{Name: name, Attributes: attrs, Namespaces: namespaces, Children: children})
}

func execNewXMLSequence(ctx *Context, instr *bir.NewXMLSequence, frame *Frame) {
	items := make([]values.XMLValue, 0, len(instr.Children))
	for _, op := range instr.Children {
		val := getOperandValue(ctx, op, frame)
		if x, ok := val.(values.XMLValue); ok {
			items = append(items, x)
		}
	}
	setOperandValue(ctx, instr.LhsOp, frame, values.NewXMLSequence(items))
}
