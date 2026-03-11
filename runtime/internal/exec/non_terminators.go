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
	"regexp"
	"strconv"
	"strings"
)

var decimalStringRegex = regexp.MustCompile(`^[+-]?[0-9]+(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)

func execConstantLoad(constantLoad *bir.ConstantLoad, frame *Frame) {
	Store(frame, constantLoad.LhsOp.Address, constantLoad.Value)
}

func execMove(moveIns *bir.Move, frame *Frame) {
	Store(frame, moveIns.LhsOp.Address, Load(frame, moveIns.RhsOp.Address))
}

func execNewArray(newArray *bir.NewArray, frame *Frame) {
	size := 0
	if newArray.SizeOp != nil {
		size = int(Load(frame, newArray.SizeOp.Address).(int64))
	}
	list := values.NewList(size, newArray.Type, newArray.Filler)
	for i, value := range newArray.Values {
		list.FillingSet(i, Load(frame, value.Address))
	}
	Store(frame, newArray.LhsOp.Address, list)
}

func execNewMap(newMap *bir.NewMap, frame *Frame) {
	m := values.NewMap(newMap.Type)
	for _, entry := range newMap.Values {
		kv := entry.(*bir.MappingConstructorKeyValueEntry)
		keyVal := Load(frame, kv.KeyOp().Address)
		keyStr := keyVal.(string)
		valueVal := Load(frame, kv.ValueOp().Address)
		m.Put(keyStr, valueVal)
	}
	Store(frame, newMap.GetLhsOperand().Address, m)
}

func execNewError(newError *bir.NewError, frame *Frame) {
	msgVal := Load(frame, newError.MessageOp.Address)
	message := msgVal.(string)

	var cause values.BalValue
	if newError.CauseOp != nil {
		cause = Load(frame, newError.CauseOp.Address)
	}

	var detailMap *values.Map
	if newError.DetailOp != nil {
		detailMap = Load(frame, newError.DetailOp.Address).(*values.Map)
	}
	errVal := values.NewError(newError.Type, message, cause, newError.TypeName, detailMap)
	Store(frame, newError.GetLhsOperand().Address, errVal)
}

func execArrayStore(access *bir.FieldAccess, frame *Frame) {
	list := Load(frame, access.LhsOp.Address).(*values.List)
	idx := int(Load(frame, access.KeyOp.Address).(int64))
	if idx < 0 {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	list.FillingSet(idx, Load(frame, access.RhsOp.Address))
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame) {
	list := Load(frame, access.RhsOp.Address).(*values.List)
	idx := int(Load(frame, access.KeyOp.Address).(int64))
	if idx < 0 || idx >= list.Len() {
		panic(fmt.Sprintf("invalid array index: %d", idx))
	}
	Store(frame, access.LhsOp.Address, list.Get(idx))
}

func execMapStore(access *bir.FieldAccess, frame *Frame) {
	m := Load(frame, access.LhsOp.Address).(*values.Map)
	keyVal := Load(frame, access.KeyOp.Address)
	keyStr := keyVal.(string)
	valueVal := Load(frame, access.RhsOp.Address)
	m.Put(keyStr, valueVal)
}

func execMapLoad(access *bir.FieldAccess, frame *Frame) {
	m := Load(frame, access.RhsOp.Address).(*values.Map)
	key := Load(frame, access.KeyOp.Address).(string)
	value, _ := m.Get(key)
	Store(frame, access.LhsOp.Address, value)
}

func execTypeCast(typeCast *bir.TypeCast, frame *Frame) {
	sourceValue := Load(frame, typeCast.RhsOp.Address)
	result := castValue(sourceValue, typeCast.Type)
	Store(frame, typeCast.LhsOp.Address, result)
}

func execFPLoad(fpLoad *bir.FPLoad, frame *Frame) {
	fn := &values.Function{
		Type:      fpLoad.Type,
		LookupKey: fpLoad.FunctionLookupKey,
	}
	Store(frame, fpLoad.LhsOp.Address, fn)
}

func execTypeTest(typeTest *bir.TypeTest, frame *Frame, reg *modules.Registry) {
	sourceValue := Load(frame, typeTest.RhsOp.Address)
	valueType := values.SemTypeForValue(sourceValue)
	typeEnv := reg.GetTypeEnv()
	typeCtx := semtypes.TypeCheckContext(typeEnv)
	matches := semtypes.IsSubtype(typeCtx, valueType, typeTest.Type) != typeTest.IsNegation
	Store(frame, typeTest.LhsOp.Address, matches)
}

func castValue(value values.BalValue, targetType semtypes.SemType) values.BalValue {
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
		s := strconv.FormatFloat(v, 'g', -1, 64)
		r := new(big.Rat)
		if _, ok := r.SetString(s); !ok {
			panic(fmt.Sprintf("bad type cast: cannot cast %v to decimal", v))
		}
		return r
	case *big.Rat:
		return v
	case string:
		if strings.Contains(v, "/") || !decimalStringRegex.MatchString(v) {
			panic(fmt.Sprintf("cannot cast %v to decimal", v))
		}
		r := new(big.Rat)
		if _, ok := r.SetString(v); !ok {
			panic(fmt.Sprintf("cannot cast %v to decimal", v))
		}
		return r
	default:
		panic(fmt.Sprintf("bad type cast: cannot cast %v to decimal", value))
	}
}
