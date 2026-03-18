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

// Package typepool provide the internal implementation on how to serialize and deserialize semtypes
package typepool

// TODO: I think we will eventually need to serialize symbol table as well and then move this package to either
// semtypes or to top level
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"ballerina-lang-go/semtypes"
)

type TypePool struct {
	tys  []semtypes.SemType
	memo map[semtypes.SemType]Index
}

func NewTypePool() *TypePool {
	return &TypePool{
		memo: make(map[semtypes.SemType]Index),
	}
}

// Index represent handle to the [TypePool]
// If first bit is set or value is 0 it is a inline type (used to represent simple basic types)
// else it is an Index + 1 to the typePool
type Index int32

const indexMask = (1 << 31) - 1

func (pool *TypePool) Get(i Index) semtypes.SemType {
	if i <= 0 {
		bits := i & indexMask
		return semtypes.BasicTypeBitSetFrom(int(bits))
	}
	return pool.tys[i-1]
}

func (pool *TypePool) Put(ty semtypes.SemType) Index {
	switch ty := ty.(type) {
	case semtypes.BasicTypeBitSet:
		return Index(ty.All() | 1<<31)
	case semtypes.ComplexSemType:
		if cached, ok := pool.memo[ty]; ok {
			return cached
		}
		id := len(pool.tys) + 1
		pool.tys = append(pool.tys, ty)
		pool.memo[ty] = Index(id)
		return Index(id)
	}
	panic("unreachable")
}

func fromTypePool(pool *TypePool) binaryPool {
	bp := binaryPool{}
	for _, ty := range pool.tys {
		cst := ty.(semtypes.ComplexSemType)
		start := uint32(len(bp.subtypeData))
		for _, sd := range cst.SubtypeDataList() {
			var entry subtypeDataEntry
			switch data := sd.(type) {
			case semtypes.IntSubtype:
				entry = subtypeDataEntry{kind: intSubtypeData, index: uint32(len(bp.intSubtypes))}
				bp.intSubtypes = append(bp.intSubtypes, fromIntSubtype(&data))
			case semtypes.BooleanSubtype:
				entry = subtypeDataEntry{kind: booleanSubtypeData, index: uint32(len(bp.booleanSubtype))}
				bp.booleanSubtype = append(bp.booleanSubtype, fromBooleanSubtype(&data))
			case semtypes.FloatSubtype:
				entry = subtypeDataEntry{kind: floatSubtypeData, index: uint32(len(bp.floatSubtype))}
				bp.floatSubtype = append(bp.floatSubtype, fromFloatSubtype(&data))
			case semtypes.DecimalSubtype:
				entry = subtypeDataEntry{kind: decimalSubtypeData, index: uint32(len(bp.decimalSubtype))}
				bp.decimalSubtype = append(bp.decimalSubtype, fromDecimalSubtype(&data))
			case semtypes.StringSubtype:
				entry = subtypeDataEntry{kind: stringSubtypeData, index: uint32(len(bp.stringSubtype))}
				bp.stringSubtype = append(bp.stringSubtype, fromStringSubtype(&data))
			case semtypes.Bdd:
				panic("unimplemented")
			default:
				panic("unexpected subtype data type")
			}
			bp.subtypeData = append(bp.subtypeData, entry)
		}
		end := uint32(len(bp.subtypeData))
		bp.types = append(bp.types, typeEntry{
			all:              uint32(cst.All()),
			some:             uint32(cst.Some()),
			subtypeDataStart: start,
			subtypeDataEnd:   end,
		})
	}
	bp.nIntSubtypes = uint32(len(bp.intSubtypes))
	bp.nBooleanSubtypes = uint32(len(bp.booleanSubtype))
	bp.nFloatSubtypes = uint32(len(bp.floatSubtype))
	bp.nDecimalSubtypes = uint32(len(bp.decimalSubtype))
	bp.nStringSubtypes = uint32(len(bp.stringSubtype))
	return bp
}

func toTypePool(bp binaryPool) *TypePool {
	pool := &TypePool{memo: make(map[semtypes.SemType]Index)}
	for _, te := range bp.types {
		var subtypeDataList []semtypes.ProperSubtypeData
		for _, sde := range bp.subtypeData[te.subtypeDataStart:te.subtypeDataEnd] {
			var data semtypes.ProperSubtypeData
			switch sde.kind {
			case intSubtypeData:
				data = toIntSubtype(bp.intSubtypes[sde.index])
			case booleanSubtypeData:
				data = toBooleanSubtype(bp.booleanSubtype[sde.index])
			case floatSubtypeData:
				data = toFloatSubtype(bp.floatSubtype[sde.index])
			case decimalSubtypeData:
				data = toDecimalSubtype(bp.decimalSubtype[sde.index])
			case stringSubtypeData:
				data = toStringSubtype(bp.stringSubtype[sde.index])
			}
			subtypeDataList = append(subtypeDataList, data)
		}
		ty := semtypes.CreateComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(
			int(te.all), int(te.some), subtypeDataList,
		)
		pool.tys = append(pool.tys, ty)
	}
	return pool
}

func MarshalTypePool(pool *TypePool) []byte {
	bp := fromTypePool(pool)
	buf := &bytes.Buffer{}

	write(buf, bp.nIntSubtypes)
	write(buf, bp.nBooleanSubtypes)
	write(buf, bp.nFloatSubtypes)
	write(buf, bp.nDecimalSubtypes)
	write(buf, bp.nStringSubtypes)

	for _, entry := range bp.intSubtypes {
		marshalIntSubtype(buf, entry)
	}
	for _, entry := range bp.booleanSubtype {
		marshalBooleanSubtype(buf, entry)
	}
	for _, entry := range bp.floatSubtype {
		marshalFloatSubtype(buf, entry)
	}
	for _, entry := range bp.decimalSubtype {
		marshalDecimalSubtype(buf, entry)
	}
	for _, entry := range bp.stringSubtype {
		marshalStringSubtype(buf, entry)
	}

	marshalSubtypeData(buf, bp.subtypeData)
	marshalTypes(buf, bp.types)

	return buf.Bytes()
}

func UnmarshalTypePool(data []byte) *TypePool {
	r := bytes.NewReader(data)
	bp := binaryPool{}

	read(r, &bp.nIntSubtypes)
	read(r, &bp.nBooleanSubtypes)
	read(r, &bp.nFloatSubtypes)
	read(r, &bp.nDecimalSubtypes)
	read(r, &bp.nStringSubtypes)

	bp.intSubtypes = make([]intSubTypeEntry, bp.nIntSubtypes)
	for i := range bp.intSubtypes {
		bp.intSubtypes[i] = unmarshalIntSubtype(r)
	}
	bp.booleanSubtype = make([]booleanSubtypeEntry, bp.nBooleanSubtypes)
	for i := range bp.booleanSubtype {
		bp.booleanSubtype[i] = unmarshalBooleanSubtype(r)
	}
	bp.floatSubtype = make([]floatSubtypeEntry, bp.nFloatSubtypes)
	for i := range bp.floatSubtype {
		bp.floatSubtype[i] = unmarshalFloatSubtype(r)
	}
	bp.decimalSubtype = make([]decimalSubtypeEntry, bp.nDecimalSubtypes)
	for i := range bp.decimalSubtype {
		bp.decimalSubtype[i] = unmarshalDecimalSubtype(r)
	}
	bp.stringSubtype = make([]stringSubtypeEntry, bp.nStringSubtypes)
	for i := range bp.stringSubtype {
		bp.stringSubtype[i] = unmarshalStringSubtype(r)
	}

	bp.subtypeData = unmarshalSubtypeData(r)
	bp.types = unmarshalTypes(r)

	return toTypePool(bp)
}

// binaryPool represent how type pool is represented in memory
type binaryPool struct {
	nIntSubtypes     uint32
	nBooleanSubtypes uint32
	nFloatSubtypes   uint32
	nDecimalSubtypes uint32
	nStringSubtypes  uint32
	intSubtypes      []intSubTypeEntry
	booleanSubtype   []booleanSubtypeEntry
	floatSubtype     []floatSubtypeEntry
	decimalSubtype   []decimalSubtypeEntry
	stringSubtype    []stringSubtypeEntry
	subtypeData      []subtypeDataEntry
	types            []typeEntry
}

type typeEntry struct {
	all              uint32
	some             uint32
	subtypeDataStart uint32
	subtypeDataEnd   uint32
}

func marshalTypes(buf *bytes.Buffer, types []typeEntry) {
	write(buf, uint32(len(types)))
	for _, entry := range types {
		write(buf, entry.all)
		write(buf, entry.some)
		write(buf, entry.subtypeDataStart)
		write(buf, entry.subtypeDataEnd)
	}
}

func unmarshalTypes(r *bytes.Reader) []typeEntry {
	var count uint32
	read(r, &count)
	types := make([]typeEntry, count)
	for i := range types {
		read(r, &types[i].all)
		read(r, &types[i].some)
		read(r, &types[i].subtypeDataStart)
		read(r, &types[i].subtypeDataEnd)
	}
	return types
}

type subtypeDataEntry struct {
	kind  subtypeDataKind
	index uint32
}

type subtypeDataKind uint8

const (
	intSubtypeData subtypeDataKind = iota
	booleanSubtypeData
	floatSubtypeData
	decimalSubtypeData
	stringSubtypeData
)

func marshalSubtypeData(buf *bytes.Buffer, entries []subtypeDataEntry) {
	write(buf, uint32(len(entries)))
	for _, entry := range entries {
		write(buf, uint8(entry.kind))
		write(buf, entry.index)
	}
}

func unmarshalSubtypeData(r *bytes.Reader) []subtypeDataEntry {
	var count uint32
	read(r, &count)
	entries := make([]subtypeDataEntry, count)
	for i := range entries {
		var kind uint8
		read(r, &kind)
		entries[i].kind = subtypeDataKind(kind)
		read(r, &entries[i].index)
	}
	return entries
}

// int subtype

type intSubTypeEntry struct {
	nRanges uint32
	ranges  []intRange
}

// currently this is fixed size but I don't want semtypes to have such a constraint so in the future
// if needed this should define it's own
type intRange semtypes.Range

func fromIntSubtype(st *semtypes.IntSubtype) intSubTypeEntry {
	ranges := make([]intRange, len(st.Ranges))
	for i, r := range st.Ranges {
		ranges[i] = intRange(r)
	}
	return intSubTypeEntry{nRanges: uint32(len(ranges)), ranges: ranges}
}

func toIntSubtype(entry intSubTypeEntry) semtypes.IntSubtype {
	ranges := make([]semtypes.Range, len(entry.ranges))
	for i, r := range entry.ranges {
		ranges[i] = semtypes.Range(r)
	}
	return semtypes.NewIntSubtypeFromRanges(ranges)
}

func marshalIntSubtype(buf *bytes.Buffer, entry intSubTypeEntry) {
	write(buf, entry.nRanges)
	for _, r := range entry.ranges {
		write(buf, r.Min)
		write(buf, r.Max)
	}
}

func unmarshalIntSubtype(r *bytes.Reader) intSubTypeEntry {
	var entry intSubTypeEntry
	read(r, &entry.nRanges)
	entry.ranges = make([]intRange, entry.nRanges)
	for j := range entry.ranges {
		read(r, &entry.ranges[j].Min)
		read(r, &entry.ranges[j].Max)
	}
	return entry
}

// boolean subtype

type booleanSubtypeEntry bool

func fromBooleanSubtype(st *semtypes.BooleanSubtype) booleanSubtypeEntry {
	return booleanSubtypeEntry(st.Value)
}

func toBooleanSubtype(entry booleanSubtypeEntry) semtypes.BooleanSubtype {
	return semtypes.BooleanSubtypeFrom(bool(entry))
}

func marshalBooleanSubtype(buf *bytes.Buffer, entry booleanSubtypeEntry) {
	write(buf, bool(entry))
}

func unmarshalBooleanSubtype(r *bytes.Reader) booleanSubtypeEntry {
	var v bool
	read(r, &v)
	return booleanSubtypeEntry(v)
}

// float subtype

// These have fixed sizes so binary can read them
type sized interface {
	int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | decimal
}
type enumerableSubtypeEntry[T sized] struct {
	allowed bool
	nValues int32
	values  []T
}

type floatSubtypeEntry enumerableSubtypeEntry[float64]

func fromFloatSubtype(st *semtypes.FloatSubtype) floatSubtypeEntry {
	nValues := len(st.Values())
	values := make([]float64, nValues)
	for i, v := range st.Values() {
		values[i] = v.Value()
	}
	return floatSubtypeEntry{allowed: st.Allowed(), nValues: int32(nValues), values: values}
}

func toFloatSubtype(entry floatSubtypeEntry) semtypes.ProperSubtypeData {
	values := make([]semtypes.EnumerableType[float64], entry.nValues)
	for i, v := range entry.values {
		f := semtypes.EnumerableFloatFrom(v)
		values[i] = &f
	}
	return semtypes.CreateFloatSubtype(entry.allowed, values)
}

func marshalFloatSubtype(buf *bytes.Buffer, entry floatSubtypeEntry) {
	write(buf, entry.allowed)
	write(buf, entry.nValues)
	for _, v := range entry.values {
		write(buf, v)
	}
}

func unmarshalFloatSubtype(r *bytes.Reader) floatSubtypeEntry {
	var entry floatSubtypeEntry
	read(r, &entry.allowed)
	read(r, &entry.nValues)
	entry.values = make([]float64, entry.nValues)
	for j := range entry.values {
		read(r, &entry.values[j])
	}
	return entry
}

// decimal subtype

type decimalSubtypeEntry enumerableSubtypeEntry[decimal]

// big.Rat is not fixed size
type decimal struct {
	num   int64
	denom int64
}

func fromDecimalSubtype(st *semtypes.DecimalSubtype) decimalSubtypeEntry {
	values := make([]decimal, len(st.Values()))
	for i, v := range st.Values() {
		r := v.Value()
		values[i] = decimal{num: r.Num().Int64(), denom: r.Denom().Int64()}
	}
	return decimalSubtypeEntry{allowed: st.Allowed(), nValues: int32(len(values)), values: values}
}

func toDecimalSubtype(entry decimalSubtypeEntry) semtypes.ProperSubtypeData {
	values := make([]semtypes.EnumerableType[big.Rat], len(entry.values))
	for i, v := range entry.values {
		r := new(big.Rat).SetFrac64(v.num, v.denom)
		d := semtypes.EnumerableDecimalFrom(*r)
		values[i] = &d
	}
	return semtypes.CreateDecimalSubtype(entry.allowed, values)
}

func marshalDecimalSubtype(buf *bytes.Buffer, entry decimalSubtypeEntry) {
	write(buf, entry.allowed)
	write(buf, entry.nValues)
	for _, v := range entry.values {
		write(buf, v.num)
		write(buf, v.denom)
	}
}

func unmarshalDecimalSubtype(r *bytes.Reader) decimalSubtypeEntry {
	var entry decimalSubtypeEntry
	read(r, &entry.allowed)
	read(r, &entry.nValues)
	entry.values = make([]decimal, entry.nValues)
	for j := range entry.values {
		read(r, &entry.values[j].num)
		read(r, &entry.values[j].denom)
	}
	return entry
}

// string subtype

type stringSubtypeEntry struct {
	charData    enumerableStringData
	nonCharData enumerableStringData
}

type enumerableStringData struct {
	allowed bool
	num     int64
	values  []enumerableStringDataEntry
}

type enumerableStringDataEntry struct {
	len    int32
	values []byte
}

func fromEnumerableStringSubtype(es semtypes.EnumerableSubtype[string]) enumerableStringData {
	entries := make([]enumerableStringDataEntry, len(es.Values()))
	for i, v := range es.Values() {
		b := []byte(v.Value())
		entries[i] = enumerableStringDataEntry{len: int32(len(b)), values: b}
	}
	return enumerableStringData{allowed: es.Allowed(), num: int64(len(entries)), values: entries}
}

func toEnumerableStrings(data enumerableStringData) (bool, []semtypes.EnumerableType[string]) {
	values := make([]semtypes.EnumerableType[string], len(data.values))
	for i, v := range data.values {
		values[i] = semtypes.EnumerableCharStringFrom(string(v.values))
	}
	return data.allowed, values
}

func fromStringSubtype(st *semtypes.StringSubtype) stringSubtypeEntry {
	return stringSubtypeEntry{
		charData:    fromEnumerableStringSubtype(st.GetChar()),
		nonCharData: fromEnumerableStringSubtype(st.GetNonChar()),
	}
}

func toStringSubtype(entry stringSubtypeEntry) semtypes.StringSubtype {
	charAllowed, charValues := toEnumerableStrings(entry.charData)
	nonCharAllowed, nonCharValues := toEnumerableStrings(entry.nonCharData)
	return semtypes.StringSubtypeFrom(
		semtypes.CharStringSubtypeFrom(charAllowed, charValues),
		semtypes.NonCharStringSubtypeFrom(nonCharAllowed, nonCharValues),
	)
}

func marshalStringSubtype(buf *bytes.Buffer, entry stringSubtypeEntry) {
	marshalEnumerableStringData(buf, entry.charData)
	marshalEnumerableStringData(buf, entry.nonCharData)
}

func unmarshalStringSubtype(r *bytes.Reader) stringSubtypeEntry {
	return stringSubtypeEntry{
		charData:    unmarshalEnumerableStringData(r),
		nonCharData: unmarshalEnumerableStringData(r),
	}
}

func marshalEnumerableStringData(buf *bytes.Buffer, data enumerableStringData) {
	write(buf, data.allowed)
	write(buf, data.num)
	for _, entry := range data.values {
		write(buf, entry.len)
		_, err := buf.Write(entry.values)
		if err != nil {
			panic(fmt.Sprintf("writing string entry bytes: %v", err))
		}
	}
}

func unmarshalEnumerableStringData(r *bytes.Reader) enumerableStringData {
	var data enumerableStringData
	read(r, &data.allowed)
	read(r, &data.num)
	data.values = make([]enumerableStringDataEntry, data.num)
	for i := range data.values {
		read(r, &data.values[i].len)
		data.values[i].values = make([]byte, data.values[i].len)
		_, err := r.Read(data.values[i].values)
		if err != nil {
			panic(fmt.Sprintf("reading string entry bytes: %v", err))
		}
	}
	return data
}

// binary helpers

func write(buf *bytes.Buffer, data any) {
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		panic(fmt.Sprintf("writing binary data: %v", err))
	}
}

func read(r *bytes.Reader, v any) {
	if err := binary.Read(r, binary.BigEndian, v); err != nil {
		panic(fmt.Sprintf("reading binary data: %v", err))
	}
}
