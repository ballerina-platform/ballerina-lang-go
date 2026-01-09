package semtypes

import (
	"ballerina-lang-go/common"
	"fmt"
	"math/big"
	"math/bits"
)

const (
	MAX_VALUE = int64(^uint(0) >> 1) // Platform max int (typically 2^63-1 on 64-bit systems)
	MIN_VALUE = -MAX_VALUE - 1       // Platform min int
)

func cellAtomType(atom Atom) CellAtomicType {
	ta := atom.(*TypeAtom)
	atomicType := ta.AtomicType
	if cellAtomicType, ok := atomicType.(*CellAtomicType); ok {
		return *cellAtomicType
	}
	panic("expected cell atomic type")
}

func Diff(t1, t2 SemType) SemType {
	var all1, all2, some1, some2 int
	if b1, ok := t1.(*BasicTypeBitSet); ok {
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			return common.ToPointer(BasicTypeBitSetFrom(b1.bitset & ^b2.bitset))
		} else {
			if b1.bitset == 0 {
				return t1
			}
			complexT2 := t2.(ComplexSemType)
			all2 = complexT2.All()
			some2 = complexT2.Some()
		}
		all1 = b1.bitset
		some1 = 0
	} else {
		c1 := t1.(ComplexSemType)
		all1 = c1.All()
		some1 = c1.Some()
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			if b2.bitset == VT_MASK {
				return common.ToPointer(BasicTypeBitSetFrom(0))
			}
			all2 = b2.bitset
			some2 = 0
		} else {
			c2 := t2.(ComplexSemType)
			all2 = c2.All()
			some2 = c2.Some()
		}
	}
	all := BasicTypeBitSetFrom(all1 & ^(all2 | some2))
	someBitSet := (all1 | some1) & ^all2
	someBitSet = someBitSet & ^all.bitset
	some := BasicTypeBitSetFrom(someBitSet)
	if some.bitset == 0 {
		return common.ToPointer(basicTypeUnion(all.bitset))
	}
	var subtypes []BasicSubtype
	for pair := range newSubtypePairs(t1, t2, some) {
		code := pair.BasicTypeCode
		data1 := pair.SubtypeData1
		data2 := pair.SubtypeData2
		var data SubtypeData
		if data1 == nil {
			data = ops[code.Code].Complement(data2)
		} else if data2 == nil {
			data = data1
		} else {
			data = ops[code.Code].Diff(data1, data2)
		}
		if allOrNothing, ok := data.(AllOrNothingSubtype); !ok {
			subtypes = append(subtypes, BasicSubtypeFrom(code, data.(ProperSubtypeData)))
		} else if allOrNothing.IsAllSubtype() {
			c := code.Code
			all = BasicTypeBitSetFrom(all.bitset | (1 << c))
		}
	}
	if len(subtypes) == 0 {
		return common.ToPointer(all)
	}
	return CreateComplexSemType(all.bitset, subtypes...)
}

func unpackComplexSemType(t ComplexSemType) []BasicSubtype {
	some := t.Some()
	var subtypes []BasicSubtype
	for _, data := range t.SubtypeDataList() {
		code := BasicTypeCodeFrom(int(bits.TrailingZeros(uint(some))))
		subtypes = append(subtypes, BasicSubtypeFrom(code, data))
		some = some & ^(1 << code.Code)
	}
	return subtypes
}

func getComplexSubtypeData(t ComplexSemType, code BasicTypeCode) SubtypeData {
	c := code.Code
	c = 1 << c
	if (t.All() & c) != 0 {
		return CreateAll()
	}
	if (t.Some() & c) == 0 {
		return CreateNothing()
	}
	loBits := t.Some() & (c - 1)
	var index int
	if loBits == 0 {
		index = 0
	} else {
		index = bits.OnesCount(uint(loBits))
	}
	return t.SubtypeDataList()[index]
}

func Union(t1, t2 SemType) SemType {
	common.Assert(t1 != nil && t2 != nil)
	var all1, all2, some1, some2 int
	if b1, ok := t1.(*BasicTypeBitSet); ok {
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			return common.ToPointer(BasicTypeBitSetFrom(b1.bitset | b2.bitset))
		} else {
			complexT2 := t2.(ComplexSemType)
			all2 = complexT2.All()
			some2 = complexT2.Some()
		}
		all1 = b1.bitset
		some1 = 0
	} else {
		complexT1 := t1.(ComplexSemType)
		all1 = complexT1.All()
		some1 = complexT1.Some()
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			all2 = b2.bitset
			some2 = 0
		} else {
			complexT2 := t2.(ComplexSemType)
			all2 = complexT2.All()
			some2 = complexT2.Some()
		}
	}
	all := BasicTypeBitSetFrom(all1 | all2)
	some := BasicTypeBitSetFrom((some1 | some2) & ^all.bitset)
	if some.bitset == 0 {
		return common.ToPointer(basicTypeUnion(all.bitset))
	}
	var subtypes []BasicSubtype
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
			data = ops[code.Code].Union(data1, data2)
		}
		if allOrNothing, ok := data.(AllOrNothingSubtype); ok && allOrNothing.IsAllSubtype() {
			c := code.Code
			all = BasicTypeBitSetFrom(all.bitset | (1 << c))
		} else {
			subtypes = append(subtypes, BasicSubtypeFrom(code, data.(ProperSubtypeData)))
		}
	}
	if len(subtypes) == 0 {
		return common.ToPointer(all)
	}
	return CreateComplexSemType(all.bitset, subtypes...)
}

func Intersect(t1, t2 SemType) SemType {
	common.Assert(t1 != nil && t2 != nil)
	var all1, all2, some1, some2 int
	if b1, ok := t1.(*BasicTypeBitSet); ok {
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			return common.ToPointer(BasicTypeBitSetFrom(b1.bitset & b2.bitset))
		} else {
			if b1.bitset == 0 {
				return t1
			}
			if b1.bitset == VT_MASK {
				return t2
			}
			complexT2 := t2.(ComplexSemType)
			all2 = complexT2.All()
			some2 = complexT2.Some()
		}
		all1 = b1.bitset
		some1 = 0
	} else {
		complexT1 := t1.(ComplexSemType)
		all1 = complexT1.All()
		some1 = complexT1.Some()
		if b2, ok := t2.(*BasicTypeBitSet); ok {
			if b2.bitset == 0 {
				return t2
			}
			if b2.bitset == VT_MASK {
				return t1
			}
			all2 = b2.bitset
			some2 = 0
		} else {
			complexT2 := t2.(ComplexSemType)
			all2 = complexT2.All()
			some2 = complexT2.Some()
		}
	}
	all := BasicTypeBitSetFrom(all1 & all2)
	some := BasicTypeBitSetFrom((some1 | all1) & (some2 | all2))
	some = BasicTypeBitSetFrom(some.bitset & ^all.bitset)
	if some.bitset == 0 {
		return common.ToPointer(basicTypeUnion(all.bitset))
	}
	var subtypes []BasicSubtype
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
			data = ops[code.Code].Intersect(data1, data2)
		}
		if allOrNothing, ok := data.(AllOrNothingSubtype); !ok || !allOrNothing.IsAllSubtype() {
			subtypes = append(subtypes, BasicSubtypeFrom(code, data.(ProperSubtypeData)))
		}
	}
	if len(subtypes) == 0 {
		return common.ToPointer(all)
	}
	return CreateComplexSemType(all.bitset, subtypes...)
}

func intersectMemberSemTypes(env Env, t1, t2 CellSemType) CellSemType {
	c1 := cellAtomicType(t1)
	c2 := cellAtomicType(t2)
	common.Assert(c1 != nil && c2 != nil)
	atomicType := IntersectCellAtomicType(c1, c2)
	var mut CellMutability
	if atomicType.Ty == &UNDEF {
		mut = CellMutability_CELL_MUT_NONE
	} else {
		mut = atomicType.Mut
	}
	return CellContainingWithEnvSemTypeCellMutability(env, atomicType.Ty, mut)
}

func Complement(t SemType) SemType {
	return Diff(&VAL, t)
}

func IsNever(t SemType) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return b.bitset == 0
	}
	return false
}

func IsEmpty(cx Context, t SemType) bool {
	common.Assert(t != nil && cx != nil)
	if b, ok := t.(*BasicTypeBitSet); ok {
		return b.bitset == 0
	} else {
		ct := t.(ComplexSemType)
		if ct.All() != 0 {
			return false
		}
		for _, st := range unpackComplexSemType(ct) {
			if !ops[st.BasicTypeCode.Code].IsEmpty(cx, st.SubtypeData) {
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
	var bits int
	if b1, ok := t1.(*BasicTypeBitSet); ok {
		bits = b1.bitset
	} else {
		complexT1 := t1.(ComplexSemType)
		bits = complexT1.All() | complexT1.Some()
	}
	return (bits & ^t2.bitset) == 0
}

func IsSameType(cx Context, t1, t2 SemType) bool {
	return IsSubtype(cx, t1, t2) && IsSubtype(cx, t2, t1)
}

func widenToBasicTypes(t SemType) BasicTypeBitSet {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return *b
	} else {
		complexSemType := t.(ComplexSemType)
		return BasicTypeBitSetFrom(complexSemType.All() | complexSemType.Some())
	}
}

func wideUnsigned(t SemType) SemType {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return b
	} else {
		if !IsSubtypeSimple(t, INT) {
			return t
		}
		data := IntSubtypeWidenUnsigned(subtypeData(t, BT_INT))
		if _, ok := data.(AllOrNothingSubtype); ok {
			return &INT
		} else {
			return basicSubtype(BT_INT, data.(ProperSubtypeData))
		}
	}
}

func booleanSubtype(t SemType) SubtypeData {
	return subtypeData(t, BT_BOOLEAN)
}

func intSubtype(t SemType) SubtypeData {
	return subtypeData(t, BT_INT)
}

func floatSubtype(t SemType) SubtypeData {
	return subtypeData(t, BT_FLOAT)
}

func decimalSubtype(t SemType) SubtypeData {
	return subtypeData(t, BT_DECIMAL)
}

func stringSubtype(t SemType) SubtypeData {
	return subtypeData(t, BT_STRING)
}

func listMemberTypeInnerVal(cx Context, t, k SemType) SemType {
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & LIST.bitset) != 0 {
			return &VAL
		} else {
			return &NEVER
		}
	} else {
		keyData := intSubtype(k)
		if isNothingSubtype(keyData) {
			return &NEVER
		}
		return BddListMemberTypeInnerVal(cx, getComplexSubtypeData(t.(ComplexSemType), BT_LIST).(Bdd), keyData, &VAL)
	}
}

var LIST_MEMBER_TYPES_ALL = ListMemberTypesFrom([]Range{RangeFrom(0, int64(MAX_VALUE))}, []SemType{&VAL})

var LIST_MEMBER_TYPES_NONE = ListMemberTypesFrom([]Range{}, []SemType{})

func ListAllMemberTypesInner(cx Context, t SemType) ListMemberTypes {
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & LIST.bitset) != 0 {
			return LIST_MEMBER_TYPES_ALL
		} else {
			return LIST_MEMBER_TYPES_NONE
		}
	}

	ct := t.(ComplexSemType)
	ranges := []Range{}
	types := []SemType{}

	allRanges := bddListAllRanges(cx, getComplexSubtypeData(ct, BT_LIST).(Bdd), []Range{})
	for _, r := range allRanges {
		m := listMemberTypeInnerVal(cx, t, IntConst(r.Min))
		if !IsNever(m) {
			ranges = append(ranges, r)
			types = append(types, m)
		}
	}
	return ListMemberTypesFrom(ranges, types)
}

func bddListAllRanges(cx Context, b Bdd, accum []Range) []Range {
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accum
		} else {
			return []Range{}
		}
	} else {
		bddNode := b.(BddNode)
		listMemberTypes := ListAtomicTypeAllMemberTypesInnerVal(cx.listAtomType(bddNode.Atom()))
		return distinctRanges(bddListAllRanges(cx, bddNode.Left(),
			distinctRanges(listMemberTypes.Ranges, accum)),
			distinctRanges(bddListAllRanges(cx, bddNode.Middle(), accum),
				bddListAllRanges(cx, bddNode.Right(), accum)))
	}
}

func distinctRanges(range1, range2 []Range) []Range {
	combined := CombineRanges(range1, range2)
	rangeResult := make([]Range, len(combined))
	for i := 0; i < len(combined); i++ {
		rangeResult[i] = combined[i].Range
	}
	return rangeResult
}

func CombineRanges(ranges1, ranges2 []Range) []CombinedRange {
	combined := []CombinedRange{}
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

		var in1 *int64 = nil
		if i1 < len1 {
			r := ranges1[i1]
			if cur >= int64(r.Min) && max <= int64(r.Max) {
				temp := int64(i1)
				in1 = &temp
			}
		}

		var in2 *int64 = nil
		if i2 < len2 {
			r := ranges2[i2]
			if cur >= int64(r.Min) && max <= int64(r.Max) {
				temp := int64(i2)
				in2 = &temp
			}
		}

		if in1 != nil || in2 != nil {
			combined = append(combined, CombinedRangeFrom(RangeFrom(cur, max), *in1, *in2))
		}

		if next == nil {
			break
		}
		cur = *next
	}
	return combined
}

func nextBoundary(cur int64, r Range, next *int64) *int64 {
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

func ListAtomicTypeAllMemberTypesInnerVal(atomicType *ListAtomicType) ListMemberTypes {
	ranges := []Range{}
	types := []SemType{}

	cellInitial := atomicType.Members.Initial
	initialLength := int64(len(cellInitial))

	initial := make([]SemType, 0, initialLength)
	for i, c := range cellInitial {
		initial[i] = CellInnerVal(c)
	}

	fixedLength := int64(atomicType.Members.FixedLength)
	if initialLength != 0 {
		types = append(types, initial...)
		for i := int64(0); i < initialLength; i++ {
			ranges = append(ranges, RangeFrom(i, i))
		}
		if initialLength < fixedLength {
			ranges[initialLength-1] = RangeFrom(initialLength-1, fixedLength-1)
		}
	}

	rest := CellInnerVal(atomicType.Rest)
	if !IsNever(rest) {
		types = append(types, rest)
		ranges = append(ranges, RangeFrom(fixedLength, MAX_VALUE))
	}

	return ListMemberTypesFrom(ranges, types)
}

func mappingAtomicType(cx Context, t SemType) *MappingAtomicType {
	mappingAtomicInner := MAPPING_ATOMIC_INNER
	if b, ok := t.(*BasicTypeBitSet); ok {
		if b.bitset == MAPPING.bitset {
			return &mappingAtomicInner
		} else {
			return nil
		}
	} else {
		env := cx.env()
		if !IsSubtypeSimple(t, MAPPING) {
			return nil
		}
		return bddMappingAtomicType(env,
			getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd),
			mappingAtomicInner)
	}
}

func bddMappingAtomicType(env Env, bdd Bdd, top MappingAtomicType) *MappingAtomicType {
	if allOrNothing, ok := bdd.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bddNode := bdd.(BddNode)
	if bddNodeSimple, ok := bddNode.(*BddNodeSimple); ok {
		result := env.mappingAtomType(bddNodeSimple.Atom())
		return result
	}
	return nil
}

func MappingMemberTypeInnerVal(cx Context, t, k SemType) SemType {
	return Diff(MappingMemberTypeInner(cx, t, k), &UNDEF)
}

func MappingMemberTypeInner(cx Context, t, k SemType) SemType {
	if b, ok := t.(*BasicTypeBitSet); ok {
		if (b.bitset & MAPPING.bitset) != 0 {
			return &VAL
		} else {
			return &UNDEF
		}
	} else {
		keyData := stringSubtype(k)
		if isNothingSubtype(keyData) {
			return &UNDEF
		}
		return BddMappingMemberTypeInner(cx, getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd), keyData,
			&INNER)
	}
}

func listAtomicType(cx Context, t SemType) *ListAtomicType {
	listAtomicInner := LIST_ATOMIC_INNER
	if b, ok := t.(*BasicTypeBitSet); ok {
		if b.bitset == LIST.bitset {
			return &listAtomicInner
		} else {
			return nil
		}
	} else {
		env := cx.env()
		if !IsSubtypeSimple(t, LIST) {
			return nil
		}
		return bddListAtomicType(env,
			getComplexSubtypeData(t.(ComplexSemType), BT_LIST).(Bdd),
			listAtomicInner)
	}
}

func bddListAtomicType(env Env, bdd Bdd, top ListAtomicType) *ListAtomicType {
	if allOrNothing, ok := bdd.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bddNode := bdd.(BddNode)
	if bddNodeSimple, ok := bddNode.(*BddNodeSimple); ok {
		result := env.listAtomType(bddNodeSimple.Atom())
		return result
	}
	return nil
}

func CellInnerVal(t CellSemType) SemType {
	return Diff(CellInner(t), &UNDEF)
}

func CellInner(t CellSemType) SemType {
	cat := cellAtomicType(t)
	common.Assert(cat != nil)
	return cat.Ty
}

func CellContainingInnerVal(env Env, t CellSemType) CellSemType {
	cat := cellAtomicType(t)
	common.Assert(cat != nil)
	return CellContainingWithEnvSemTypeCellMutability(env, Diff(cat.Ty, &UNDEF), cat.Mut)
}

func cellAtomicType(t SemType) *CellAtomicType {
	if bt, ok := t.(*BasicTypeBitSet); ok {
		if *bt == CELL {
			return CELL_ATOMIC_VAL
		} else {
			return nil
		}
	} else {
		if !IsSubtypeSimple(t, CELL) {
			return nil
		}
		return bddCellAtomicType(getComplexSubtypeData(t.(ComplexSemType), BT_CELL).(Bdd), *CELL_ATOMIC_VAL)
	}
}

func bddCellAtomicType(bdd Bdd, top CellAtomicType) *CellAtomicType {
	if allOrNothing, ok := bdd.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return &top
		}
		return nil
	}
	bddNode := bdd.(BddNode)
	leftBdd := bddNode.Left()
	middleBdd := bddNode.Middle()
	rightBdd := bddNode.Right()

	if leftAll, ok := leftBdd.(BddAllOrNothing); ok && leftAll.IsAll() {
		if middleNothing, ok := middleBdd.(BddAllOrNothing); ok && middleNothing.IsNothing() {
			if rightNothing, ok := rightBdd.(BddAllOrNothing); ok && rightNothing.IsNothing() {
				result := cellAtomType(bddNode.Atom())
				return &result
			}
		}
	}
	return nil
}

func SingleShape(t SemType) common.Optional[Value] {
	if t == &NIL {
		return common.OptionalOf(ValueFrom(nil))
	} else if _, ok := t.(*BasicTypeBitSet); ok {
		return common.OptionalEmpty[Value]()
	} else if IsSubtypeSimple(t, INT) {
		sd := getComplexSubtypeData(t.(ComplexSemType), BT_INT)
		value := IntSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(ValueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, FLOAT) {
		sd := getComplexSubtypeData(t.(ComplexSemType), BT_FLOAT)
		value := FloatSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(ValueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, STRING) {
		sd := getComplexSubtypeData(t.(ComplexSemType), BT_STRING)
		value := StringSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(ValueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, BOOLEAN) {
		sd := getComplexSubtypeData(t.(ComplexSemType), BT_BOOLEAN)
		value := BooleanSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(ValueFrom(value.Get()))
		}
	} else if IsSubtypeSimple(t, DECIMAL) {
		sd := getComplexSubtypeData(t.(ComplexSemType), BT_DECIMAL)
		value := DecimalSubtypeSingleValue(sd)
		if value.IsEmpty() {
			return common.OptionalEmpty[Value]()
		} else {
			return common.OptionalOf(ValueFrom(fmt.Sprintf("%v", value.Get())))
		}
	}
	return common.OptionalEmpty[Value]()
}

func Singleton(v any) SemType {
	if v == nil {
		return &NIL
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

func ContainsConst(t SemType, v any) bool {
	if v == nil {
		return ContainsNil(t)
	} else if lng, ok := v.(int64); ok {
		return ContainsConstInt(t, lng)
	} else if d, ok := v.(float64); ok {
		return ContainsConstFloat(t, d)
	} else if s, ok := v.(string); ok {
		return ContainsConstString(t, s)
	} else if b, ok := v.(bool); ok {
		return ContainsConstBoolean(t, b)
	} else {
		// Assuming it's a BigDecimal (big.Rat in Go)
		return ContainsConstDecimal(t, v.(big.Rat))
	}
}

func ContainsNil(t SemType) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return (b.bitset & (1 << BT_NIL.Code)) != 0
	} else {
		// todo: Need to verify this behavior
		complexSubtypeData := getComplexSubtypeData(t.(ComplexSemType), BT_NIL).(AllOrNothingSubtype)
		return complexSubtypeData.IsAllSubtype()
	}
}

func ContainsConstString(t SemType, s string) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return (b.bitset & (1 << BT_STRING.Code)) != 0
	} else {
		return StringSubtypeContains(
			getComplexSubtypeData(t.(ComplexSemType), BT_STRING), s)
	}
}

func ContainsConstInt(t SemType, n int64) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return (b.bitset & (1 << BT_INT.Code)) != 0
	} else {
		return IntSubtypeContains(
			getComplexSubtypeData(t.(ComplexSemType), BT_INT), n)
	}
}

func ContainsConstFloat(t SemType, n float64) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return (b.bitset & (1 << BT_FLOAT.Code)) != 0
	} else {
		return FloatSubtypeContains(
			getComplexSubtypeData(t.(ComplexSemType), BT_FLOAT), EnumerableFloatFrom(n))
	}
}

func ContainsConstDecimal(t SemType, n big.Rat) bool {
	if b, ok := t.(*BasicTypeBitSet); ok {
		return (b.bitset & (1 << BT_DECIMAL.Code)) != 0
	} else {
		return DecimalSubtypeContains(
			getComplexSubtypeData(t.(ComplexSemType), BT_DECIMAL), EnumerableDecimalFrom(n))
	}
}

func ContainsConstBoolean(t SemType, b bool) bool {
	if bType, ok := t.(*BasicTypeBitSet); ok {
		return (bType.bitset & (1 << BT_BOOLEAN.Code)) != 0
	} else {
		return BooleanSubtypeContains(
			getComplexSubtypeData(t.(ComplexSemType), BT_BOOLEAN), b)
	}
}

func SingleNumericType(semType SemType) common.Optional[BasicTypeBitSet] {
	numType := Intersect(semType, &NUMBER)
	if b, ok := numType.(*BasicTypeBitSet); ok {
		if b.bitset == NEVER.bitset {
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
	if b, ok := s.(*BasicTypeBitSet); ok {
		if (b.bitset & (1 << code.Code)) != 0 {
			return CreateAll()
		}
		return CreateNothing()
	} else {
		return getComplexSubtypeData(s.(ComplexSemType), code)
	}
}

func TypeCheckContext(env Env) Context {
	return ContextFrom(env)
}

func CreateJson(context Context) SemType {
	memo := context.jsonMemo()
	env := context.env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	j := Union(&SIMPLE_OR_STRING, Union(listDef.GetSemType(env), mapDef.GetSemType(env)))
	listDef.DefineListTypeWrappedWithEnvSemType(env, j)
	mapDef.DefineMappingTypeWrapped(env, nil, j)
	context.setJsonMemo(j)
	return j
}

func CreateAnydata(context Context) SemType {
	memo := context.anydataMemo()
	env := context.env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	tableTy := TableContaining(env, mapDef.GetSemType(env))
	ad := Union(Union(&SIMPLE_OR_STRING, Union(&XML, Union(&REGEXP, tableTy))),
		Union(listDef.GetSemType(env), mapDef.GetSemType(env)))
	listDef.DefineListTypeWrappedWithEnvSemType(env, ad)
	mapDef.DefineMappingTypeWrapped(env, nil, ad)
	context.setAnydataMemo(ad)
	return ad
}

func CreateCloneable(context Context) SemType {
	memo := context.cloneableMemo()
	env := context.env()

	if memo != nil {
		return memo
	}
	listDef := &ListDefinition{}
	mapDef := &MappingDefinition{}
	tableTy := TableContaining(env, mapDef.GetSemType(env))
	ad := Union(VAL_READONLY, Union(&XML, Union(listDef.GetSemType(env), Union(tableTy,
		mapDef.GetSemType(env)))))
	listDef.DefineListTypeWrappedWithEnvSemType(env, ad)
	mapDef.DefineMappingTypeWrapped(env, []Field{}, ad)
	context.setCloneableMemo(ad)
	return ad
}

func CreateIsolatedObject(context Context) SemType {
	memo := context.isolatedObjectMemo()
	if memo != nil {
		return memo
	}

	quals := ObjectQualifiersFrom(true, false, NetworkQualifierNone)
	od := NewObjectDefinition()
	isolatedObj := od.Define(context.env(), quals, []Member{})
	context.setIsolatedObjectMemo(isolatedObj)
	return isolatedObj
}

func CreateServiceObject(context Context) SemType {
	memo := context.serviceObjectMemo()
	if memo != nil {
		return memo
	}

	quals := ObjectQualifiersFrom(false, false, NetworkQualifierService)
	od := NewObjectDefinition()
	serviceObj := od.Define(context.env(), quals, []Member{})
	context.setServiceObjectMemo(serviceObj)
	return serviceObj
}

func CreateBasicSemType(typeCode BasicTypeCode, subtypeData SubtypeData) SemType {
	if _, ok := subtypeData.(AllOrNothingSubtype); ok {
		if isAllSubtype(subtypeData) {
			return common.ToPointer(BasicTypeBitSetFrom(1 << typeCode.Code))
		} else {
			return common.ToPointer(BasicTypeBitSetFrom(0))
		}
	} else {
		return CreateComplexSemType(0,
			BasicSubtypeFrom(typeCode, subtypeData.(ProperSubtypeData)))
	}
}

func MappingAtomicTypesInUnion(cx Context, t SemType) common.Optional[[]MappingAtomicType] {
	matList := []MappingAtomicType{}
	mappingAtomicInner := MAPPING_ATOMIC_INNER
	if b, ok := t.(*BasicTypeBitSet); ok {
		if b.bitset == MAPPING.bitset {
			matList = append(matList, mappingAtomicInner)
			return common.OptionalOf(matList)
		}
		return common.OptionalEmpty[[]MappingAtomicType]()
	} else {
		env := cx.env()
		if !IsSubtypeSimple(t, MAPPING) {
			return common.OptionalEmpty[[]MappingAtomicType]()
		}
		if collectBddMappingAtomicTypesInUnion(env,
			getComplexSubtypeData(t.(ComplexSemType), BT_MAPPING).(Bdd),
			mappingAtomicInner, &matList) {
			return common.OptionalOf(matList)
		} else {
			return common.OptionalEmpty[[]MappingAtomicType]()
		}
	}
}

func collectBddMappingAtomicTypesInUnion(env Env, bdd Bdd, top MappingAtomicType, matList *[]MappingAtomicType) bool {
	if allOrNothing, ok := bdd.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			*matList = append(*matList, top)
			return true
		}
		return false
	}
	bddNode := bdd.(BddNode)
	if bddNodeSimple, ok := bddNode.(*BddNodeSimple); ok {
		*matList = append(*matList, *env.mappingAtomType(bddNodeSimple.Atom()))
		return true
	}

	bddNodeImpl := bdd.(BddNodeImpl)
	leftBdd := bddNodeImpl.Left()
	rightBdd := bddNodeImpl.Right()

	if leftNode, ok := leftBdd.(BddAllOrNothing); ok && leftNode.IsAll() {
		if rightNode, ok := rightBdd.(BddAllOrNothing); ok && rightNode.IsNothing() {
			*matList = append(*matList, *env.mappingAtomType(bddNodeImpl.Atom()))
			return collectBddMappingAtomicTypesInUnion(env, bddNodeImpl.Middle(), top, matList)
		}
	}

	return false
}
