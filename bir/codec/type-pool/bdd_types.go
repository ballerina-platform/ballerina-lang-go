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

package typepool

// I am not very happy about the way this turned out. In an ideal world what we should be doing is serailize semtypes in to a form
// that is close as possible to ballerina syntax (with negation) then reparse it and re-resolve the type cleanly in a new type env.
// To do this not only we need to serilaize type but type operations more generally, which we currently do in a half-hearted manner for BDDs.
// When we have time we should think about a better way to do this that does not require these speical cases.

import (
	"bytes"

	"ballerina-lang-go/semtypes"
)

// BDD DNF representation: union of intersections

type atomEntry struct {
	isRec bool
	index int32
}

type conjunction struct {
	nPosAtoms uint32
	posAtoms  []atomEntry
	nNegAtoms uint32
	negAtoms  []atomEntry
}

type unionOfIntersections struct {
	nConjunctions uint32
	conjunctions  []conjunction
}

// Atomic type entries for each BDD kind

type listAtomicTypeEntry struct {
	fixedLength int32
	nInitial    int32
	initial     []Index
	rest        Index
	mut         uint8
}

type mappingAtomicTypeEntry struct {
	nFields int32
	names   []enumerableStringDataEntry
	types   []Index
	rest    Index
	mut     uint8
}

type functionAtomicTypeEntry struct {
	paramType  Index
	retType    Index
	qualifiers Index
	isGeneric  bool
}

type xmlAtomicTypeEntry struct {
	index int32
}

type xmlSubtypeEntry struct {
	primitives int32
	sequence   unionOfIntersections
}

// Serialization: decompose BDD into DNF

type bddSerializationContext struct {
	pool *TypePool
	cx   semtypes.Context

	listAtomMap     map[semtypes.Atom]int32
	mappingAtomMap  map[semtypes.Atom]int32
	functionAtomMap map[semtypes.Atom]int32
	xmlAtomMap      map[semtypes.Atom]int32

	bp *binaryPool
}

func newBddSerializationContext(pool *TypePool, cx semtypes.Context, bp *binaryPool) *bddSerializationContext {
	sc := &bddSerializationContext{
		pool:            pool,
		cx:              cx,
		listAtomMap:     make(map[semtypes.Atom]int32),
		mappingAtomMap:  make(map[semtypes.Atom]int32),
		functionAtomMap: make(map[semtypes.Atom]int32),
		xmlAtomMap:      make(map[semtypes.Atom]int32),
		bp:              bp,
	}
	// Reserve index 0 in list and mapping atom tables for BDD_REC_ATOM_READONLY
	bp.listAtomicTypes = append(bp.listAtomicTypes, listAtomicTypeEntry{})
	bp.mappingAtomicTypes = append(bp.mappingAtomicTypes, mappingAtomicTypeEntry{})
	return sc
}

func (sc *bddSerializationContext) serializeListBdd(bdd semtypes.Bdd) unionOfIntersections {
	return sc.serializeBdd(bdd, sc.listAtomMap, sc.serializeListAtom)
}

func (sc *bddSerializationContext) serializeMappingBdd(bdd semtypes.Bdd) unionOfIntersections {
	return sc.serializeBdd(bdd, sc.mappingAtomMap, sc.serializeMappingAtom)
}

func (sc *bddSerializationContext) serializeFunctionBdd(bdd semtypes.Bdd) unionOfIntersections {
	return sc.serializeBdd(bdd, sc.functionAtomMap, sc.serializeFunctionAtom)
}

func (sc *bddSerializationContext) serializeXmlBdd(bdd semtypes.Bdd) unionOfIntersections {
	return sc.serializeBdd(bdd, sc.xmlAtomMap, sc.serializeXMLAtom)
}

func (sc *bddSerializationContext) serializeXmlSubtype(xs semtypes.XmlSubtype) xmlSubtypeEntry {
	return xmlSubtypeEntry{
		primitives: int32(xs.Primitives),
		sequence:   sc.serializeXmlBdd(xs.Sequence),
	}
}

func (sc *bddSerializationContext) serializeBdd(
	bdd semtypes.Bdd,
	atomMap map[semtypes.Atom]int32,
	serializeAtom func(semtypes.Atom) int32,
) unionOfIntersections {
	var conjs []conjunction
	semtypes.BddEvery(sc.cx, bdd, nil, nil, func(_ semtypes.Context, pos *semtypes.Conjunction, neg *semtypes.Conjunction) bool {
		var posAtoms []atomEntry
		for c := pos; c != nil; c = c.Next {
			posAtoms = append(posAtoms, sc.resolveAtom(c.Atom, atomMap, serializeAtom))
		}
		var negAtoms []atomEntry
		for c := neg; c != nil; c = c.Next {
			negAtoms = append(negAtoms, sc.resolveAtom(c.Atom, atomMap, serializeAtom))
		}
		conjs = append(conjs, conjunction{
			nPosAtoms: uint32(len(posAtoms)),
			posAtoms:  posAtoms,
			nNegAtoms: uint32(len(negAtoms)),
			negAtoms:  negAtoms,
		})
		return true
	})
	return unionOfIntersections{
		nConjunctions: uint32(len(conjs)),
		conjunctions:  conjs,
	}
}

func (sc *bddSerializationContext) resolveAtom(
	atom semtypes.Atom,
	atomMap map[semtypes.Atom]int32,
	serializeAtom func(semtypes.Atom) int32,
) atomEntry {
	if recAtom, ok := atom.(*semtypes.RecAtom); ok && recAtom.Index() == semtypes.BDD_REC_ATOM_READONLY {
		return atomEntry{isRec: false, index: 0}
	}
	if idx, ok := atomMap[atom]; ok {
		return atomEntry{isRec: true, index: idx}
	}
	idx := serializeAtom(atom)
	return atomEntry{isRec: false, index: idx}
}

func (sc *bddSerializationContext) serializeListAtom(atom semtypes.Atom) int32 {
	idx := int32(len(sc.bp.listAtomicTypes))
	sc.listAtomMap[atom] = idx
	sc.bp.listAtomicTypes = append(sc.bp.listAtomicTypes, listAtomicTypeEntry{})

	at := sc.cx.ListAtomType(atom)
	initial := make([]Index, len(at.Members.Initial))
	var mut uint8
	for i, cell := range at.Members.Initial {
		initial[i] = sc.pool.Put(semtypes.CellInnerVal(cell))
		mut = uint8(cellMut(cell))
	}
	rest := sc.pool.Put(semtypes.CellInnerVal(at.Rest))
	if len(at.Members.Initial) == 0 {
		mut = uint8(cellMut(at.Rest))
	}

	sc.bp.listAtomicTypes[idx] = listAtomicTypeEntry{
		fixedLength: int32(at.Members.FixedLength),
		nInitial:    int32(len(initial)),
		initial:     initial,
		rest:        rest,
		mut:         mut,
	}
	return idx
}

func (sc *bddSerializationContext) serializeMappingAtom(atom semtypes.Atom) int32 {
	idx := int32(len(sc.bp.mappingAtomicTypes))
	sc.mappingAtomMap[atom] = idx
	sc.bp.mappingAtomicTypes = append(sc.bp.mappingAtomicTypes, mappingAtomicTypeEntry{})

	at := sc.cx.MappingAtomType(atom)
	names := make([]enumerableStringDataEntry, len(at.Names))
	types := make([]Index, len(at.Types))
	var mut uint8
	for i, name := range at.Names {
		b := []byte(name)
		names[i] = enumerableStringDataEntry{len: int32(len(b)), values: b}
		types[i] = sc.pool.Put(semtypes.CellInnerVal(at.Types[i]))
		mut = uint8(cellMut(at.Types[i]))
	}
	rest := sc.pool.Put(semtypes.CellInnerVal(at.Rest))
	if len(at.Types) == 0 {
		mut = uint8(cellMut(at.Rest))
	}

	sc.bp.mappingAtomicTypes[idx] = mappingAtomicTypeEntry{
		nFields: int32(len(names)),
		names:   names,
		types:   types,
		rest:    rest,
		mut:     mut,
	}
	return idx
}

func (sc *bddSerializationContext) serializeFunctionAtom(atom semtypes.Atom) int32 {
	idx := int32(len(sc.bp.functionAtomicTypes))
	sc.functionAtomMap[atom] = idx
	sc.bp.functionAtomicTypes = append(sc.bp.functionAtomicTypes, functionAtomicTypeEntry{})

	at := sc.cx.FunctionAtomType(atom)
	entry := functionAtomicTypeEntry{
		paramType:  sc.pool.Put(at.ParamType),
		retType:    sc.pool.Put(at.RetType),
		qualifiers: sc.pool.Put(at.Qualifiers),
		isGeneric:  at.IsGeneric,
	}

	sc.bp.functionAtomicTypes[idx] = entry
	return idx
}

func (sc *bddSerializationContext) serializeXMLAtom(atom semtypes.Atom) int32 {
	idx := int32(len(sc.bp.xmlAtomicTypes))
	sc.xmlAtomMap[atom] = idx
	sc.bp.xmlAtomicTypes = append(sc.bp.xmlAtomicTypes, xmlAtomicTypeEntry{
		index: int32(atom.Index()),
	})
	return idx
}

func cellMut(cell semtypes.CellSemType) semtypes.CellMutability {
	bdd := cell.SubtypeDataList()[0].(semtypes.BddNode)
	cat := bdd.Atom().(*semtypes.TypeAtom).AtomicType.(*semtypes.CellAtomicType)
	return cat.Mut
}

// Deserialization: reconstruct BDDs from DNF

type bddDeserializationContext struct {
	pool *TypePool
	env  semtypes.Env
	bp   *binaryPool

	listAtomDefs []*semtypes.ListDefinition
	listAtoms    []semtypes.Atom

	mappingAtomDefs []*semtypes.MappingDefinition
	mappingAtoms    []semtypes.Atom

	functionAtomDefs []*semtypes.FunctionDefinition
	functionAtoms    []semtypes.Atom

	xmlAtoms []semtypes.Atom
}

func newBddDeserializationContext(pool *TypePool, env semtypes.Env, bp *binaryPool) *bddDeserializationContext {
	return &bddDeserializationContext{
		pool:             pool,
		env:              env,
		bp:               bp,
		listAtomDefs:     make([]*semtypes.ListDefinition, len(bp.listAtomicTypes)),
		listAtoms:        make([]semtypes.Atom, len(bp.listAtomicTypes)),
		mappingAtomDefs:  make([]*semtypes.MappingDefinition, len(bp.mappingAtomicTypes)),
		mappingAtoms:     make([]semtypes.Atom, len(bp.mappingAtomicTypes)),
		functionAtomDefs: make([]*semtypes.FunctionDefinition, len(bp.functionAtomicTypes)),
		functionAtoms:    make([]semtypes.Atom, len(bp.functionAtomicTypes)),
		xmlAtoms:         make([]semtypes.Atom, len(bp.xmlAtomicTypes)),
	}
}

func (dc *bddDeserializationContext) deserializeType(poolIndex int) semtypes.SemType {
	if dc.pool.tys[poolIndex] != nil {
		return dc.pool.tys[poolIndex]
	}
	te := dc.bp.types[poolIndex]
	var subtypeDataList []semtypes.ProperSubtypeData
	for _, sde := range dc.bp.subtypeData[te.subtypeDataStart:te.subtypeDataEnd] {
		var data semtypes.ProperSubtypeData
		switch sde.kind {
		case intSubtypeData:
			data = toIntSubtype(dc.bp.intSubtypes[sde.index])
		case booleanSubtypeData:
			data = toBooleanSubtype(dc.bp.booleanSubtype[sde.index])
		case floatSubtypeData:
			data = toFloatSubtype(dc.bp.floatSubtype[sde.index])
		case decimalSubtypeData:
			data = toDecimalSubtype(dc.bp.decimalSubtype[sde.index])
		case stringSubtypeData:
			data = toStringSubtype(dc.bp.stringSubtype[sde.index])
		case listBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.listBdds[sde.index], dc.deserializeListAtom)
		case mappingBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.mappingBdds[sde.index], dc.deserializeMappingAtom)
		case functionBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.functionBdds[sde.index], dc.deserializeFunctionAtom)
		case errorBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.errorBdds[sde.index], dc.deserializeMappingAtom)
		case tableBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.tableBdds[sde.index], dc.deserializeListAtom)
		case objectBddSubtypeData:
			data = dc.deserializeBddFromDnf(dc.bp.objectBdds[sde.index], dc.deserializeMappingAtom)
		case xmlSubtypeData:
			entry := dc.bp.xmlSubtypes[sde.index]
			sequence := dc.deserializeBddFromDnf(entry.sequence, dc.deserializeXmlAtom)
			data = semtypes.XmlSubtypeFrom(int(entry.primitives), sequence)
		}
		subtypeDataList = append(subtypeDataList, data)
	}
	ty := semtypes.CreateComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(
		int(te.all), int(te.some), subtypeDataList,
	)
	dc.pool.tys[poolIndex] = ty
	return ty
}

func (dc *bddDeserializationContext) deserializeBddFromDnf(
	dnf unionOfIntersections,
	deserializeAtom func(int32) semtypes.Atom,
) semtypes.Bdd {
	atoms := make(map[int32]semtypes.Atom)
	for _, conj := range dnf.conjunctions {
		for _, a := range conj.posAtoms {
			if _, ok := atoms[a.index]; !ok {
				atoms[a.index] = deserializeAtom(a.index)
			}
		}
		for _, a := range conj.negAtoms {
			if _, ok := atoms[a.index]; !ok {
				atoms[a.index] = deserializeAtom(a.index)
			}
		}
	}
	return buildBddFromDnf(dnf, atoms)
}

func (dc *bddDeserializationContext) deserializeListAtom(atomIndex int32) semtypes.Atom {
	if atomIndex == 0 {
		ro := semtypes.CreateRecAtom(semtypes.BDD_REC_ATOM_READONLY)
		return &ro
	}
	if dc.listAtoms[atomIndex] != nil {
		return dc.listAtoms[atomIndex]
	}
	if def := dc.listAtomDefs[atomIndex]; def != nil {
		placeholder := def.GetSemType(dc.env)
		atom := extractAtom(placeholder)
		dc.listAtoms[atomIndex] = atom
		return atom
	}

	def := semtypes.NewListDefinition()
	dc.listAtomDefs[atomIndex] = &def

	entry := dc.bp.listAtomicTypes[atomIndex]
	initial := make([]semtypes.SemType, entry.nInitial)
	for j := range initial {
		initial[j] = dc.resolvePoolType(entry.initial[j])
	}
	rest := dc.resolvePoolType(entry.rest)
	mut := semtypes.CellMutability(entry.mut)

	result := def.DefineListTypeWrapped(dc.env, initial, int(entry.fixedLength), rest, mut)
	atom := extractAtom(result)
	dc.listAtoms[atomIndex] = atom
	return atom
}

func (dc *bddDeserializationContext) deserializeMappingAtom(atomIndex int32) semtypes.Atom {
	if atomIndex == 0 {
		ro := semtypes.CreateRecAtom(semtypes.BDD_REC_ATOM_READONLY)
		return &ro
	}
	if dc.mappingAtoms[atomIndex] != nil {
		return dc.mappingAtoms[atomIndex]
	}
	if def := dc.mappingAtomDefs[atomIndex]; def != nil {
		placeholder := def.GetSemType(dc.env)
		atom := extractAtom(placeholder)
		dc.mappingAtoms[atomIndex] = atom
		return atom
	}

	def := semtypes.NewMappingDefinition()
	dc.mappingAtomDefs[atomIndex] = &def

	entry := dc.bp.mappingAtomicTypes[atomIndex]
	fields := make([]semtypes.Field, entry.nFields)
	for j := range fields {
		fields[j] = semtypes.Field{
			Name: string(entry.names[j].values),
			Ty:   dc.resolvePoolType(entry.types[j]),
		}
	}
	rest := dc.resolvePoolType(entry.rest)
	mut := semtypes.CellMutability(entry.mut)

	result := def.DefineMappingTypeWrappedWithEnvFieldsSemTypeCellMutability(dc.env, fields, rest, mut)
	atom := extractAtom(result)
	dc.mappingAtoms[atomIndex] = atom
	return atom
}

func (dc *bddDeserializationContext) deserializeFunctionAtom(atomIndex int32) semtypes.Atom {
	if dc.functionAtoms[atomIndex] != nil {
		return dc.functionAtoms[atomIndex]
	}
	if def := dc.functionAtomDefs[atomIndex]; def != nil {
		placeholder := def.GetSemType(dc.env)
		atom := extractAtom(placeholder)
		dc.functionAtoms[atomIndex] = atom
		return atom
	}

	def := semtypes.NewFunctionDefinition()
	dc.functionAtomDefs[atomIndex] = &def

	entry := dc.bp.functionAtomicTypes[atomIndex]
	paramType := dc.resolvePoolType(entry.paramType)
	retType := dc.resolvePoolType(entry.retType)
	qualifiers := dc.resolvePoolType(entry.qualifiers)

	var result semtypes.SemType
	if entry.isGeneric {
		result = def.DefineGeneric(dc.env, paramType, retType, semtypes.NewFunctionQualifiersFromSemType(qualifiers))
	} else {
		result = def.Define(dc.env, paramType, retType, semtypes.NewFunctionQualifiersFromSemType(qualifiers))
	}
	atom := extractAtom(result)
	dc.functionAtoms[atomIndex] = atom
	return atom
}

func (dc *bddDeserializationContext) deserializeXmlAtom(atomIndex int32) semtypes.Atom {
	if dc.xmlAtoms[atomIndex] != nil {
		return dc.xmlAtoms[atomIndex]
	}
	entry := dc.bp.xmlAtomicTypes[atomIndex]
	recAtom := semtypes.CreateXMLRecAtom(int(entry.index))
	dc.xmlAtoms[atomIndex] = &recAtom
	return &recAtom
}

func (dc *bddDeserializationContext) resolvePoolType(idx Index) semtypes.SemType {
	if idx <= 0 {
		bits := idx & indexMask
		return semtypes.BasicTypeBitSetFrom(int(bits))
	}
	return dc.deserializeType(int(idx - 1))
}

func extractAtom(ty semtypes.SemType) semtypes.Atom {
	cst := ty.(semtypes.ComplexSemType)
	return cst.SubtypeDataList()[0].(semtypes.BddNode).Atom()
}

func buildBddFromDnf(dnf unionOfIntersections, atoms map[int32]semtypes.Atom) semtypes.Bdd {
	var bdd semtypes.Bdd = semtypes.BddNothing()
	for _, conj := range dnf.conjunctions {
		var term semtypes.Bdd = semtypes.BddAll()
		for _, a := range conj.posAtoms {
			term = semtypes.BddIntersect(term, semtypes.BddAtom(atoms[a.index]))
		}
		for _, a := range conj.negAtoms {
			term = semtypes.BddDiff(term, semtypes.BddAtom(atoms[a.index]))
		}
		bdd = semtypes.BddUnion(bdd, term)
	}
	return bdd
}

// Marshal/unmarshal for BDD types

func marshalAtomEntry(buf *bytes.Buffer, entry atomEntry) {
	write(buf, entry.isRec)
	write(buf, entry.index)
}

func unmarshalAtomEntry(r *bytes.Reader) atomEntry {
	var entry atomEntry
	read(r, &entry.isRec)
	read(r, &entry.index)
	return entry
}

func marshalBddDnf(buf *bytes.Buffer, dnf unionOfIntersections) {
	write(buf, dnf.nConjunctions)
	for _, conj := range dnf.conjunctions {
		write(buf, conj.nPosAtoms)
		for _, a := range conj.posAtoms {
			marshalAtomEntry(buf, a)
		}
		write(buf, conj.nNegAtoms)
		for _, a := range conj.negAtoms {
			marshalAtomEntry(buf, a)
		}
	}
}

func unmarshalBddDnf(r *bytes.Reader) unionOfIntersections {
	var dnf unionOfIntersections
	read(r, &dnf.nConjunctions)
	dnf.conjunctions = make([]conjunction, dnf.nConjunctions)
	for i := range dnf.conjunctions {
		read(r, &dnf.conjunctions[i].nPosAtoms)
		dnf.conjunctions[i].posAtoms = make([]atomEntry, dnf.conjunctions[i].nPosAtoms)
		for j := range dnf.conjunctions[i].posAtoms {
			dnf.conjunctions[i].posAtoms[j] = unmarshalAtomEntry(r)
		}
		read(r, &dnf.conjunctions[i].nNegAtoms)
		dnf.conjunctions[i].negAtoms = make([]atomEntry, dnf.conjunctions[i].nNegAtoms)
		for j := range dnf.conjunctions[i].negAtoms {
			dnf.conjunctions[i].negAtoms[j] = unmarshalAtomEntry(r)
		}
	}
	return dnf
}

func marshalListAtomicType(buf *bytes.Buffer, entry listAtomicTypeEntry) {
	write(buf, entry.fixedLength)
	write(buf, entry.nInitial)
	for _, idx := range entry.initial {
		write(buf, int32(idx))
	}
	write(buf, int32(entry.rest))
	write(buf, entry.mut)
}

func unmarshalListAtomicType(r *bytes.Reader) listAtomicTypeEntry {
	var entry listAtomicTypeEntry
	read(r, &entry.fixedLength)
	read(r, &entry.nInitial)
	entry.initial = make([]Index, entry.nInitial)
	for j := range entry.initial {
		var idx int32
		read(r, &idx)
		entry.initial[j] = Index(idx)
	}
	var rest int32
	read(r, &rest)
	entry.rest = Index(rest)
	read(r, &entry.mut)
	return entry
}

func marshalMappingAtomicType(buf *bytes.Buffer, entry mappingAtomicTypeEntry) {
	write(buf, entry.nFields)
	for _, name := range entry.names {
		write(buf, name.len)
		buf.Write(name.values)
	}
	for _, idx := range entry.types {
		write(buf, int32(idx))
	}
	rest := int32(entry.rest)
	write(buf, rest)
	write(buf, entry.mut)
}

func unmarshalMappingAtomicType(r *bytes.Reader) mappingAtomicTypeEntry {
	var entry mappingAtomicTypeEntry
	read(r, &entry.nFields)
	entry.names = make([]enumerableStringDataEntry, entry.nFields)
	for j := range entry.names {
		read(r, &entry.names[j].len)
		entry.names[j].values = make([]byte, entry.names[j].len)
		if _, err := r.Read(entry.names[j].values); err != nil {
			panic(err)
		}
	}
	entry.types = make([]Index, entry.nFields)
	for j := range entry.types {
		var idx int32
		read(r, &idx)
		entry.types[j] = Index(idx)
	}
	var rest int32
	read(r, &rest)
	entry.rest = Index(rest)
	read(r, &entry.mut)
	return entry
}

func marshalFunctionAtomicType(buf *bytes.Buffer, entry functionAtomicTypeEntry) {
	write(buf, int32(entry.paramType))
	write(buf, int32(entry.retType))
	write(buf, int32(entry.qualifiers))
	write(buf, entry.isGeneric)
}

func unmarshalFunctionAtomicType(r *bytes.Reader) functionAtomicTypeEntry {
	var entry functionAtomicTypeEntry
	var paramType, retType, qualifiers int32
	read(r, &paramType)
	read(r, &retType)
	read(r, &qualifiers)
	entry.paramType = Index(paramType)
	entry.retType = Index(retType)
	entry.qualifiers = Index(qualifiers)
	read(r, &entry.isGeneric)
	return entry
}
