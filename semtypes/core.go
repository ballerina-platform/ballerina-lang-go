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

package semtypes

import (
	"fmt"
	"math/big"
	"math/bits"

	"ballerina-lang-go/common"
)

const (
	MAX_VALUE = int64(^uint(0) >> 1) // Platform max int (typically 2^63-1 on 64-bit systems)
	MIN_VALUE = -MAX_VALUE - 1       // Platform min int
)

func bitCount(b BasicTypeBitSet) int {
	return bits.OnesCount(uint(b.all()))
}

func cellAtomType(atom atom) cellAtomicType {
	ta := atom.(*typeAtom)
	atomicType := ta.AtomicType
	if cellAtomicType, ok := atomicType.(*cellAtomicType); ok {
		return *cellAtomicType
	}
	panic("expected cell atomic type")
}

func Diff(t1, t2 SemType) SemType {
	var all1, all2, some1, some2 BasicTypeBitSet
	if b1, ok := t1.(BasicTypeBitSet); ok {
		if b2, ok := t2.(BasicTypeBitSet); ok {
			return b1.all() & ^b2.all()
		} else {
			if b1.all() == 0 {
				return t1
			}
			complexT2 := t2.(*ComplexSemType)
			all2 = complexT2.all()
			some2 = complexT2.some()
		}
		all1 = b1.all()
		some1 = 0
	} else {
		c1 := t1.(*ComplexSemType)
		all1 = c1.all()
		some1 = c1.some()
		if b2, ok := t2.(BasicTypeBitSet); ok {
			if b2.all() == ValueTypeMask {
				return BasicTypeBitSet(0)
			}
			all2 = b2.all()
			some2 = 0
		} else {
			c2 := t2.(*ComplexSemType)
			all2 = c2.all()
			some2 = c2.some()
		}
	}
	all := all1 & ^(all2 | some2)
	someBitSet := (all1 | some1) & ^all2
	someBitSet = someBitSet & ^all
	some := someBitSet
	if some == 0 {
		return basicTypeUnion(all)
	}
	var subtypes []basicSubtype
	for pair := range newSubtypePairs(t1, t2, some) {
		code := pair.BasicTypeCode
		data1 := pair.SubtypeData1
		data2 := pair.SubtypeData2
		var data SubtypeData
		if data1 == nil {
			data = ops[code.Code()].complement(data2)
		} else if data2 == nil {
			data = data1
		} else {
			data = ops[code.Code()].Diff(data1, data2)
		}
		if allOrNothing, ok := data.(allOrNothingSubtype); !ok {
			subtypes = append(subtypes, basicSubtypeFrom(code, data.(ProperSubtypeData)))
		} else if allOrNothing.IsAllSubtype() {
			c := code.Code()
			all = all | (1 << c)
		}
	}
	if len(subtypes) == 0 {
		return all
	}
	return createComplexSemType(all, subtypes...)
}

func getUnpackComplexSemType(t *ComplexSemType) []basicSubtype {
	some := t.some()
	subtypeDataList := t.subtypeDataList()
	subtypes := make([]basicSubtype, len(subtypeDataList))
	for i, data := range subtypeDataList {
		code := basicTypeCodeFrom(bits.TrailingZeros(uint(some)))
		subtypes[i] = basicSubtypeFrom(code, data)
		some = some & ^(1 << code.Code())
	}
	return subtypes
}

func getComplexSubtypeData(t *ComplexSemType, code BasicTypeCode) SubtypeData {
	c := BasicTypeBitSet(1 << code.Code())
	if (t.all() & c) != 0 {
		return createAll()
	}
	if (t.some() & c) == 0 {
		return createNothing()
	}
	loBits := t.some() & (c - 1)
	var index int
	if loBits == 0 {
		index = 0
	} else {
		index = bits.OnesCount(uint(loBits))
	}
	return t.subtypeDataList()[index]
}

func Union(t1, t2 SemType) SemType {
	common.Assert(t1 != nil && t2 != nil)
	var all1, all2, some1, some2 BasicTypeBitSet
	if b1, ok := t1.(BasicTypeBitSet); ok {
		if b2, ok := t2.(BasicTypeBitSet); ok {
			return b1.all() | b2.all()
		} else {
			complexT2 := t2.(*ComplexSemType)
			all2 = complexT2.all()
			some2 = complexT2.some()
		}
		all1 = b1.all()
		some1 = 0
	} else {
		complexT1 := t1.(*ComplexSemType)
		all1 = complexT1.all()
		some1 = complexT1.some()
		if b2, ok := t2.(BasicTypeBitSet); ok {
			all2 = b2.all()
			some2 = 0
		} else {
			complexT2 := t2.(*ComplexSemType)
			all2 = complexT2.all()
			some2 = complexT2.some()
		}
	}
	all := all1 | all2
	some := (some1 | some2) & ^all
	if some == 0 {
		return all
	}
	var subtypes []basicSubtype
	for pair := range newSubtypePairs(t1, t2, some) {
		code := pair.BasicTypeCode
		data1 := pair.SubtypeData1
		data2 := pair.SubtypeData2
		var data SubtypeData
		if data1 == nil {
			data = data2
		} else if data2 == nil {
			data = data1
		} else {
			data = ops[code.Code()].Union(data1, data2)
		}
		if allOrNothing, ok := data.(allOrNothingSubtype); ok && allOrNothing.IsAllSubtype() {
			c := code.Code()
			all = all | (1 << c)
		} else {
			subtypes = append(subtypes, basicSubtypeFrom(code, data.(ProperSubtypeData)))
		}
	}
	if len(subtypes) == 0 {
		return all
	}
	return createComplexSemType(all, subtypes...)
}

func Intersect(t1, t2 SemType) SemType {
	common.Assert(t1 != nil && t2 != nil)
	var all1, all2, some1, some2 BasicTypeBitSet
	if b1, ok := t1.(BasicTypeBitSet); ok {
		if b2, ok := t2.(BasicTypeBitSet); ok {
			return b1.all() & b2.all()
		} else {
			if b1.all() == 0 {
				return t1
			}
			if b1.all() == ValueTypeMask {
				return t2
			}
			complexT2 := t2.(*ComplexSemType)
			all2 = complexT2.all()
			some2 = complexT2.some()
		}
		all1 = b1.all()
		some1 = 0
	} else {
		complexT1 := t1.(*ComplexSemType)
		all1 = complexT1.all()
		some1 = complexT1.some()
		if b2, ok := t2.(BasicTypeBitSet); ok {
			if b2.all() == 0 {
				return t2
			}
			if b2.all() == ValueTypeMask {
				return t1
			}
			all2 = b2.all()
			some2 = 0
		} else {
			complexT2 := t2.(*ComplexSemType)
			all2 = complexT2.all()
			some2 = complexT2.some()
		}
	}
	all := all1 & all2
	some := (some1 | all1) & (some2 | all2)
	some = some & ^all
	if some == 0 {
		return basicTypeUnion(all)
	}
	var subtypes []basicSubtype
	for pair := range newSubtypePairs(t1, t2, some) {
		code := pair.BasicTypeCode
		data1 := pair.SubtypeData1
		data2 := pair.SubtypeData2
		var data SubtypeData
		if data1 == nil {
			data = data2
		} else if data2 == nil {
			data = data1
		} else {
			data = ops[code.Code()].Intersect(data1, data2)
		}
		if allOrNothing, ok := data.(allOrNothingSubtype); !ok || allOrNothing.IsAllSubtype() {
			subtypes = append(subtypes, basicSubtypeFrom(code, data.(ProperSubtypeData)))
		}
	}
	if len(subtypes) == 0 {
		return all
	}
	return createComplexSemType(all, subtypes...)
}

func intersectMemberSemTypes(env Env, t1, t2 *ComplexSemType) *ComplexSemType {
	c1 := getCellAtomicType(t1)
	c2 := getCellAtomicType(t2)
	common.Assert(c1 != nil && c2 != nil)
	atomicType := intersectCellAtomicType(c1, c2)
	var mut CellMutability
	if atomicType.Ty == UNDEF {
		mut = CellMutability_CELL_MUT_NONE
	} else {
		mut = atomicType.Mut
	}
	return cellContainingWithEnvSemTypeCellMutability(env, atomicType.Ty, mut)
}

func complement(t SemType) SemType {
	return Diff(VAL, t)
}

func IsNever(t SemType) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return b.all() == 0
	}
	return false
}

func IsEmpty(cx Context, t SemType) bool {
	common.Assert(t != nil && cx != nil)
	if b, ok := t.(BasicTypeBitSet); ok {
		return b.all() == 0
	} else {
		ct := t.(*ComplexSemType)
		if ct.all() != 0 {
			return false
		}
		for _, st := range getUnpackComplexSemType(ct) {
			if !ops[st.BasicTypeCode.Code()].IsEmpty(cx, st.SubtypeData) {
				return false
			}
		}
		return true
	}
}

func IsSubtype(cx Context, t1, t2 SemType) bool {
	return IsEmpty(cx, Diff(t1, t2))
}

func IsSubtypeSimple(t1 SemType, t2 BasicTypeBitSet) bool {
	var b BasicTypeBitSet
	if b1, ok := t1.(BasicTypeBitSet); ok {
		b = b1.all()
	} else {
		complexT1 := t1.(*ComplexSemType)
		b = complexT1.all() | complexT1.some()
	}
	return (b & ^t2.all()) == 0
}

func IsSameType(cx Context, t1, t2 SemType) bool {
	return IsSubtype(cx, t1, t2) && IsSubtype(cx, t2, t1)
}

// NBasicTypes returns the number of basic types to which the given type belongs to
func NBasicTypes(t SemType) int {
	return bitCount(WidenToBasicTypes(t))
}

func WidenToBasicTypes(t SemType) BasicTypeBitSet {
	if b, ok := t.(BasicTypeBitSet); ok {
		return b
	} else {
		complexSemType := t.(*ComplexSemType)
		return complexSemType.all() | complexSemType.some()
	}
}

func wideUnsigned(t SemType) SemType {
	if b, ok := t.(BasicTypeBitSet); ok {
		return b
	} else {
		if !IsSubtypeSimple(t, INT) {
			return t
		}
		data := intSubtypeWidenUnsigned(subtypeData(t, BTInt))
		if _, ok := data.(allOrNothingSubtype); ok {
			return INT
		} else {
			return getBasicSubtype(BTInt, data.(ProperSubtypeData))
		}
	}
}

func getBooleanSubtype(t SemType) SubtypeData {
	return subtypeData(t, BTBoolean)
}

func getIntSubtype(t SemType) SubtypeData {
	return subtypeData(t, BTInt)
}

func getFloatSubtype(t SemType) SubtypeData {
	return subtypeData(t, BTFloat)
}

func getDecimalSubtype(t SemType) SubtypeData {
	return subtypeData(t, BTDecimal)
}

func getStringSubtype(t SemType) SubtypeData {
	return subtypeData(t, BTString)
}

func ListMemberTypeInnerVal(cx Context, t, k SemType) SemType {
	if b, ok := t.(BasicTypeBitSet); ok {
		if (b.all() & LIST.all()) != 0 {
			return VAL
		} else {
			return NEVER
		}
	} else {
		keyData := getIntSubtype(k)
		if isNothingSubtype(keyData) {
			return NEVER
		}
		return bddListMemberTypeInnerVal(cx, getComplexSubtypeData(t.(*ComplexSemType), BTList).(Bdd), keyData, VAL)
	}
}

var LIST_MEMBER_TYPES_ALL = listMemberTypesFrom([]intRange{rangeFrom(0, int64(MAX_VALUE))}, []SemType{VAL})

var LIST_MEMBER_TYPES_NONE = listMemberTypesFrom([]intRange{}, []SemType{})

func ListAllMemberTypesInner(cx Context, t SemType) ListMemberTypes {
	if b, ok := t.(BasicTypeBitSet); ok {
		if (b.all() & LIST.all()) != 0 {
			return LIST_MEMBER_TYPES_ALL
		} else {
			return LIST_MEMBER_TYPES_NONE
		}
	}

	ct := t.(*ComplexSemType)
	ranges := []intRange{}
	types := []SemType{}

	allRanges := bddListAllRanges(cx, getComplexSubtypeData(ct, BTList).(Bdd), []intRange{})
	for _, r := range allRanges {
		m := ListMemberTypeInnerVal(cx, t, IntConst(r.Min))
		if !IsNever(m) {
			ranges = append(ranges, r)
			types = append(types, m)
		}
	}
	return listMemberTypesFrom(ranges, types)
}

func bddListAllRanges(cx Context, b Bdd, accum []intRange) []intRange {
	if allOrNothing, ok := b.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accum
		} else {
			return []intRange{}
		}
	} else {
		bn := b.(bddNode)
		listMemberTypes := listAtomicTypeAllMemberTypesInnerVal(cx.ListAtomType(bn.atom()))
		return distinctRanges(bddListAllRanges(cx, bn.left(),
			distinctRanges(listMemberTypes.Ranges, accum)),
			distinctRanges(bddListAllRanges(cx, bn.middle(), accum),
				bddListAllRanges(cx, bn.right(), accum)))
	}
}

func distinctRanges(range1, range2 []intRange) []intRange {
	combined := combineRanges(range1, range2)
	rangeResult := make([]intRange, len(combined))
	for i := range combined {
		rangeResult[i] = combined[i].Range
	}
	return rangeResult
}

func combineRanges(ranges1, ranges2 []intRange) []combinedRange {
	combined := []combinedRange{}
	i1 := 0
	i2 := 0
	len1 := len(ranges1)
	len2 := len(ranges2)
	cur := int64(MIN_VALUE)

	// This iterates over the boundaries between ranges
	for {
		for i1 < len1 && cur > int64(ranges1[i1].Max) {
			i1 += 1
		}
		for i2 < len2 && cur > int64(ranges2[i2].Max) {
			i2 += 1
		}

		var next *int64 = nil
		if i1 < len1 {
			next = nextBoundary(cur, ranges1[i1], next)
		}
		if i2 < len2 {
			next = nextBoundary(cur, ranges2[i2], next)
		}

		var max int64
		if next == nil {
			max = int64(MAX_VALUE)
		} else {
			max = *next - 1
		}

		var in1 int64 = -1
		if i1 < len1 {
			r := ranges1[i1]
			if cur >= int64(r.Min) && max <= int64(r.Max) {
				in1 = int64(i1)
			}
		}

		var in2 int64 = -1
		if i2 < len2 {
			r := ranges2[i2]
			if cur >= int64(r.Min) && max <= int64(r.Max) {
				in2 = int64(i2)
			}
		}

		if in1 != -1 || in2 != -1 {
			combined = append(combined, combinedRangeFrom(rangeFrom(cur, max), in1, in2))
		}

		if next == nil {
			break
		}
		cur = *next
	}
	return combined
}

func nextBoundary(cur int64, r intRange, next *int64) *int64 {
	if (int64(r.Min) > cur) && (next == nil || int64(r.Min) < *next) {
		result := int64(r.Min)
		return &result
	}
	if r.Max != int64(MAX_VALUE) {
		i := int64(r.Max) + 1
		if i > cur && (next == nil || i < *next) {
			return &i
		}
	}
	return next
}

func listAtomicTypeAllMemberTypesInnerVal(atomicType *ListAtomicType) ListMemberTypes {
	ranges := []intRange{}
	types := []SemType{}

	cellInitial := atomicType.Members.initial
	initialLength := int64(len(cellInitial))

	initial := make([]SemType, 0, initialLength)
	for i := range cellInitial {
		initial = append(initial, cellInnerVal(&cellInitial[i]))
	}

	fixedLength := int64(atomicType.Members.FixedLength)
	if initialLength != 0 {
		types = append(types, initial...)
		for i := range initialLength {
			ranges = append(ranges, rangeFrom(i, i))
		}
		if initialLength < fixedLength {
			ranges[initialLength-1] = rangeFrom(initialLength-1, fixedLength-1)
		}
	}

	rest := cellInnerVal(atomicType.rest)
	if !IsNever(rest) {
		types = append(types, rest)
		ranges = append(ranges, rangeFrom(fixedLength, MAX_VALUE))
	}

	return listMemberTypesFrom(ranges, types)
}

func ToObjectAtomicType(cx Context, t SemType) *MappingAtomicType {
	mappingTy := convertObjectToMappingTy(cx, t)
	if mappingTy == nil {
		return nil
	}
	return ToMappingAtomicType(cx, mappingTy)
}

func ToMappingAtomicType(cx Context, t SemType) *MappingAtomicType {
	mappingAtomicInner := MAPPING_ATOMIC_INNER
	if b, ok := t.(BasicTypeBitSet); ok {
		if b.all() == MAPPING.all() {
			return &mappingAtomicInner
		} else {
			return nil
		}
	} else {
		env := cx.Env()
		if !IsSubtypeSimple(t, MAPPING) {
			return nil
		}
		return bddMappingAtomicType(env,
			getComplexSubtypeData(t.(*ComplexSemType), BTMapping).(Bdd),
			mappingAtomicInner)
	}
}

func bddMappingAtomicType(env Env, bdd Bdd, top MappingAtomicType) *MappingAtomicType {
	if allOrNothing, ok := bdd.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bn := bdd.(bddNode)
	if bddNodeSimple, ok := bn.(*bddNodeSimple); ok {
		result := env.mappingAtomType(bddNodeSimple.atom())
		return result
	}
	return nil
}

func MappingMemberTypeInnerVal(cx Context, t, k SemType) SemType {
	return Diff(MappingMemberTypeInner(cx, t, k), UNDEF)
}

func MappingMemberTypeInner(cx Context, t, k SemType) SemType {
	if b, ok := t.(BasicTypeBitSet); ok {
		if (b.all() & MAPPING.all()) != 0 {
			return VAL
		} else {
			return UNDEF
		}
	} else {
		keyData := getStringSubtype(k)
		if isNothingSubtype(keyData) {
			return UNDEF
		}
		return bddMappingMemberTypeInnerCore(cx, getComplexSubtypeData(t.(*ComplexSemType), BTMapping).(Bdd), keyData,
			INNER)
	}
}

func ToListAtomicType(cx Context, t SemType) *ListAtomicType {
	listAtomicInner := LIST_ATOMIC_INNER
	if b, ok := t.(BasicTypeBitSet); ok {
		if b.all() == LIST.all() {
			return &listAtomicInner
		} else {
			return nil
		}
	} else {
		env := cx.Env()
		if !IsSubtypeSimple(t, LIST) {
			return nil
		}
		return bddListAtomicType(env,
			getComplexSubtypeData(t.(*ComplexSemType), BTList).(Bdd),
			listAtomicInner)
	}
}

func bddListAtomicType(env Env, bdd Bdd, top ListAtomicType) *ListAtomicType {
	if allOrNothing, ok := bdd.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bn := bdd.(bddNode)
	if bddNodeSimple, ok := bn.(*bddNodeSimple); ok {
		result := env.listAtomType(bddNodeSimple.atom())
		return result
	}
	return nil
}

func cellInnerVal(t *ComplexSemType) SemType {
	return Diff(cellInner(t), UNDEF)
}

func cellInner(t *ComplexSemType) SemType {
	cat := getCellAtomicType(t)
	common.Assert(cat != nil)
	return cat.Ty
}

func cellContainingInnerVal(env Env, t *ComplexSemType) *ComplexSemType {
	cat := getCellAtomicType(t)
	common.Assert(cat != nil)
	return cellContainingWithEnvSemTypeCellMutability(env, Diff(cat.Ty, UNDEF), cat.Mut)
}

func getCellAtomicType(t SemType) *cellAtomicType {
	if bt, ok := t.(BasicTypeBitSet); ok {
		if bt == CELL {
			return CELL_ATOMIC_VAL
		} else {
			return nil
		}
	} else {
		if !IsSubtypeSimple(t, CELL) {
			return nil
		}
		return bddCellAtomicType(getComplexSubtypeData(t.(*ComplexSemType), BTCell).(Bdd), *CELL_ATOMIC_VAL)
	}
}

func bddCellAtomicType(bdd Bdd, top cellAtomicType) *cellAtomicType {
	if allOrNothing, ok := bdd.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bn := bdd.(bddNode)
	leftBdd := bn.left()
	middleBdd := bn.middle()
	rightBdd := bn.right()

	if leftAll, ok := leftBdd.(*bddAllOrNothing); ok && leftAll.IsAll() {
		if middleNothing, ok := middleBdd.(*bddAllOrNothing); ok && middleNothing.IsNothing() {
			if rightNothing, ok := rightBdd.(*bddAllOrNothing); ok && rightNothing.IsNothing() {
				result := cellAtomType(bn.atom())
				return &result
			}
		}
	}
	return nil
}

func SingleShape(t SemType) common.Optional[Value] {
	if t == NIL {
		return common.OptionalOf(valueFrom(nil))
	} else if _, ok := t.(BasicTypeBitSet); ok {
		return common.OptionalEmpty[Value]()
	} else if IsSubtypeSimple(t, INT) {
		sd := getComplexSubtypeData(t.(*ComplexSemType), BTInt)
		value := intSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(valueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, FLOAT) {
		sd := getComplexSubtypeData(t.(*ComplexSemType), BTFloat)
		value := floatSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(valueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, STRING) {
		sd := getComplexSubtypeData(t.(*ComplexSemType), BTString)
		value := stringSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(valueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, BOOLEAN) {
		sd := getComplexSubtypeData(t.(*ComplexSemType), BTBoolean)
		value := booleanSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(valueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, DECIMAL) {
		sd := getComplexSubtypeData(t.(*ComplexSemType), BTDecimal)
		value := decimalSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(valueFrom(fmt.Sprintf("%v", value.Get())))
		}
	}
	return common.OptionalEmpty[Value]()
}

func singleton(v any) SemType {
	if v == nil {
		return NIL
	}

	if lng, ok := v.(int64); ok {
		return IntConst(lng)
	} else if d, ok := v.(float64); ok {
		return FloatConst(d)
	} else if s, ok := v.(string); ok {
		return StringConst(s)
	} else if b, ok := v.(bool); ok {
		return BooleanConst(b)
	} else {
		panic("Unsupported type: " + fmt.Sprintf("%T", v))
	}
}

func containsConst(t SemType, v any) bool {
	if v == nil {
		return containsNil(t)
	} else if lng, ok := v.(int64); ok {
		return containsConstInt(t, lng)
	} else if d, ok := v.(float64); ok {
		return containsConstFloat(t, d)
	} else if s, ok := v.(string); ok {
		return containsConstString(t, s)
	} else if b, ok := v.(bool); ok {
		return containsConstBoolean(t, b)
	} else {
		// Assuming it's a BigDecimal (big.Rat in Go)
		return containsConstDecimal(t, v.(big.Rat))
	}
}

func containsNil(t SemType) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return (b.all() & (1 << BTNil.Code())) != 0
	} else {
		// todo: Need to verify this behavior
		complexSubtypeData := getComplexSubtypeData(t.(*ComplexSemType), BTNil).(allOrNothingSubtype)
		return complexSubtypeData.IsAllSubtype()
	}
}

func containsConstString(t SemType, s string) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return (b.all() & (1 << BTString.Code())) != 0
	} else {
		return stringSubtypeContains(
			getComplexSubtypeData(t.(*ComplexSemType), BTString), s)
	}
}

func containsConstInt(t SemType, n int64) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return (b.all() & (1 << BTInt.Code())) != 0
	} else {
		return intSubtypeContains(
			getComplexSubtypeData(t.(*ComplexSemType), BTInt), n)
	}
}

func containsConstFloat(t SemType, n float64) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return (b.all() & (1 << BTFloat.Code())) != 0
	} else {
		return floatSubtypeContains(
			getComplexSubtypeData(t.(*ComplexSemType), BTFloat), enumerableFloatFrom(n))
	}
}

func containsConstDecimal(t SemType, n big.Rat) bool {
	if b, ok := t.(BasicTypeBitSet); ok {
		return (b.all() & (1 << BTDecimal.Code())) != 0
	} else {
		return decimalSubtypeContains(
			getComplexSubtypeData(t.(*ComplexSemType), BTDecimal), enumerableDecimalFrom(n))
	}
}

func containsConstBoolean(t SemType, b bool) bool {
	if bType, ok := t.(BasicTypeBitSet); ok {
		return (bType.all() & (1 << BTBoolean.Code())) != 0
	} else {
		return booleanSubtypeContains(
			getComplexSubtypeData(t.(*ComplexSemType), BTBoolean), b)
	}
}

func SingleNumericType(semType SemType) common.Optional[BasicTypeBitSet] {
	numType := Intersect(semType, NUMBER)
	if b, ok := numType.(BasicTypeBitSet); ok {
		if b.all() == NEVER.all() {
			return common.OptionalEmpty[BasicTypeBitSet]()
		}
	}
	if IsSubtypeSimple(numType, INT) {
		return common.OptionalOf(INT)
	}
	if IsSubtypeSimple(numType, FLOAT) {
		return common.OptionalOf(FLOAT)
	}
	if IsSubtypeSimple(numType, DECIMAL) {
		return common.OptionalOf(DECIMAL)
	}
	return common.OptionalEmpty[BasicTypeBitSet]()
}

func subtypeData(s SemType, code BasicTypeCode) SubtypeData {
	if b, ok := s.(BasicTypeBitSet); ok {
		if (b.all() & (1 << code.Code())) != 0 {
			return createAll()
		}
		return createNothing()
	} else {
		return getComplexSubtypeData(s.(*ComplexSemType), code)
	}
}

func TypeCheckContext(env Env) Context {
	return ContextFrom(env)
}

func createJson(context Context) SemType {
	memo := context.jsonMemo()
	env := context.Env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	j := Union(SIMPLE_OR_STRING, Union(listDef.GetSemType(env), mapDef.GetSemType(env)))
	listDef.DefineListTypeWrappedWithEnvSemType(env, j)
	mapDef.DefineMappingTypeWrapped(env, nil, j)
	context.setJsonMemo(j)
	return j
}

func CreateAnydata(context Context) SemType {
	memo := context.anydataMemo()
	env := context.Env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	tableTy := tableContainingDefault(env, mapDef.GetSemType(env))
	ad := Union(Union(SIMPLE_OR_STRING, Union(XML, Union(REGEXP, tableTy))),
		Union(listDef.GetSemType(env), mapDef.GetSemType(env)))
	listDef.DefineListTypeWrappedWithEnvSemType(env, ad)
	mapDef.DefineMappingTypeWrapped(env, nil, ad)
	context.setAnydataMemo(ad)
	return ad
}

func CreateCloneable(context Context) SemType {
	memo := context.cloneableMemo()
	env := context.Env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	tableTy := tableContainingDefault(env, mapDef.GetSemType(env))
	ad := Union(VAL_READONLY, Union(XML, Union(listDef.GetSemType(env), Union(tableTy,
		mapDef.GetSemType(env)))))
	listDef.DefineListTypeWrappedWithEnvSemType(env, ad)
	mapDef.DefineMappingTypeWrapped(env, []Field{}, ad)
	context.setCloneableMemo(ad)
	return ad
}

func createIsolatedObject(context Context) SemType {
	memo := context.isolatedObjectMemo()
	if memo != nil {
		return memo
	}

	quals := ObjectQualifiersFrom(true, false, NetworkQualifierNone)
	od := NewObjectDefinition()
	isolatedObj := od.Define(context.Env(), quals, []Member{})
	context.setIsolatedObjectMemo(isolatedObj)
	return isolatedObj
}

func createServiceObject(context Context) SemType {
	memo := context.serviceObjectMemo()
	if memo != nil {
		return memo
	}

	quals := ObjectQualifiersFrom(false, false, NetworkQualifierService)
	od := NewObjectDefinition()
	serviceObj := od.Define(context.Env(), quals, []Member{})
	context.setServiceObjectMemo(serviceObj)
	return serviceObj
}

func createBasicSemType(typeCode BasicTypeCode, subtypeData SubtypeData) SemType {
	if _, ok := subtypeData.(allOrNothingSubtype); ok {
		if isAllSubtype(subtypeData) {
			return basicTypeBitSetFrom(1 << typeCode.Code())
		} else {
			return basicTypeBitSetFrom(0)
		}
	} else {
		return createComplexSemType(0,
			basicSubtypeFrom(typeCode, subtypeData.(ProperSubtypeData)))
	}
}

func mappingAtomicTypesInUnion(cx Context, t SemType) common.Optional[[]MappingAtomicType] {
	matList := []MappingAtomicType{}
	mappingAtomicInner := MAPPING_ATOMIC_INNER
	if b, ok := t.(BasicTypeBitSet); ok {
		if b.all() == MAPPING.all() {
			matList = append(matList, mappingAtomicInner)
			return common.OptionalOf(matList)
		}
		return common.OptionalEmpty[[]MappingAtomicType]()
	} else {
		env := cx.Env()
		if !IsSubtypeSimple(t, MAPPING) {
			return common.OptionalEmpty[[]MappingAtomicType]()
		}
		if collectBddMappingAtomicTypesInUnion(env,
			getComplexSubtypeData(t.(*ComplexSemType), BTMapping).(Bdd),
			mappingAtomicInner, &matList) {
			return common.OptionalOf(matList)
		} else {
			return common.OptionalEmpty[[]MappingAtomicType]()
		}
	}
}

func collectBddMappingAtomicTypesInUnion(env Env, bdd Bdd, top MappingAtomicType, matList *[]MappingAtomicType) bool {
	if allOrNothing, ok := bdd.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			*matList = append(*matList, top)
			return true
		}
		return false
	}
	bn := bdd.(bddNode)
	if bddNodeSimple, ok := bn.(*bddNodeSimple); ok {
		*matList = append(*matList, *env.mappingAtomType(bddNodeSimple.atom()))
		return true
	}

	bddNodeImpl := bdd.(*bddNodeImpl)
	leftBdd := bddNodeImpl.left()
	rightBdd := bddNodeImpl.right()

	if leftNode, ok := leftBdd.(*bddAllOrNothing); ok && leftNode.IsAll() {
		if rightNode, ok := rightBdd.(*bddAllOrNothing); ok && rightNode.IsNothing() {
			*matList = append(*matList, *env.mappingAtomType(bddNodeImpl.atom()))
			return collectBddMappingAtomicTypesInUnion(env, bddNodeImpl.middle(), top, matList)
		}
	}

	return false
}

func Comparable(cx Context, t1, t2 SemType) bool {
	semType := Diff(Union(t1, t2), NIL)
	if IsSubtypeSimple(semType, SIMPLE_OR_STRING) {
		nOrderings := bitCount(WidenToBasicTypes(semType))
		return nOrderings <= 1
	}
	if IsSubtypeSimple(semType, LIST) {
		return comparableNillableList(cx, t1, t2)
	}
	return false
}

// t1, t2 must be subtype of LIST|?
func comparableNillableList(cx Context, t1, t2 SemType) bool {
	memoized := cx.comparableMemo(t1, t2)
	if memoized != nil {
		return memoized.comparable
	}
	memo := comparableMemo{semType1: t1, semType2: t2}
	cx.setComparableMemo(t1, t2, &memo)
	listMemberTypes1 := ListAllMemberTypesInner(cx, t1)
	listMemberTypes2 := ListAllMemberTypesInner(cx, t2)
	ranges1 := listMemberTypes1.Ranges
	ranges2 := listMemberTypes2.Ranges
	memberTypes1 := listMemberTypes1.SemTypes
	memberTypes2 := listMemberTypes2.SemTypes
	for _, combinedRange := range combineRanges(ranges1, ranges2) {
		i1 := combinedRange.I1
		i2 := combinedRange.I2
		if i1 != -1 && i2 != -1 && !Comparable(cx, memberTypes1[i1], memberTypes2[i2]) {
			memo.comparable = false
			return false
		}
	}
	memo.comparable = true
	return true
}

func ContainsUndef(t SemType) bool {
	switch t := t.(type) {
	case BasicTypeBitSet:
		bitSet := t.all()
		return (bitSet & (1 << BTUndef.Code())) != 0
	case *ComplexSemType:
		switch data := getComplexSubtypeData(t, BTUndef).(type) {
		case allOrNothingSubtype:
			return data.isAll
		case *bool:
			return *data
		default:
			panic("unexpected subtype data")
		}
	default:
		panic("unexpected semtype")
	}
}
