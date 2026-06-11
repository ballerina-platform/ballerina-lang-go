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

type preallocatedTypeVals struct {
	none      basicCellTypeVals
	limited   basicCellTypeVals
	unlimited basicCellTypeVals
}

type basicCellTypeVals struct {
	never          ComplexSemType
	nil            ComplexSemType
	boolean        ComplexSemType
	int            ComplexSemType
	float          ComplexSemType
	decimal        ComplexSemType
	string         ComplexSemType
	error          ComplexSemType
	list           ComplexSemType
	mapping        ComplexSemType
	table          ComplexSemType
	undef          ComplexSemType
	regexp         ComplexSemType
	function       ComplexSemType
	typedesc       ComplexSemType
	handle         ComplexSemType
	xml            ComplexSemType
	object         ComplexSemType
	stream         ComplexSemType
	future         ComplexSemType
	val            ComplexSemType
	inner          ComplexSemType
	any            ComplexSemType
	simpleOrString ComplexSemType
	number         ComplexSemType
	simpleBasic    ComplexSemType
}

func newPreallocatedTypeVals(env Env) preallocatedTypeVals {
	return preallocatedTypeVals{
		none:      newBasicCellTypeVals(env, CellMutability_CELL_MUT_NONE),
		limited:   newBasicCellTypeVals(env, CellMutability_CELL_MUT_LIMITED),
		unlimited: newBasicCellTypeVals(env, CellMutability_CELL_MUT_UNLIMITED),
	}
}

func newBasicCellTypeVals(env Env, mut CellMutability) basicCellTypeVals {
	return basicCellTypeVals{
		never:          createCellTypeVal(env, NEVER, mut),
		nil:            createCellTypeVal(env, NIL, mut),
		boolean:        createCellTypeVal(env, BOOLEAN, mut),
		int:            createCellTypeVal(env, INT, mut),
		float:          createCellTypeVal(env, FLOAT, mut),
		decimal:        createCellTypeVal(env, DECIMAL, mut),
		string:         createCellTypeVal(env, STRING, mut),
		error:          createCellTypeVal(env, ERROR, mut),
		list:           createCellTypeVal(env, LIST, mut),
		mapping:        createCellTypeVal(env, MAPPING, mut),
		table:          createCellTypeVal(env, TABLE, mut),
		undef:          createCellTypeVal(env, UNDEF, mut),
		regexp:         createCellTypeVal(env, REGEXP, mut),
		function:       createCellTypeVal(env, FUNCTION, mut),
		typedesc:       createCellTypeVal(env, TYPEDESC, mut),
		handle:         createCellTypeVal(env, HANDLE, mut),
		xml:            createCellTypeVal(env, XML, mut),
		object:         createCellTypeVal(env, OBJECT, mut),
		stream:         createCellTypeVal(env, STREAM, mut),
		future:         createCellTypeVal(env, FUTURE, mut),
		val:            createCellTypeVal(env, VAL, mut),
		inner:          createCellTypeVal(env, INNER, mut),
		any:            createCellTypeVal(env, ANY, mut),
		simpleOrString: createCellTypeVal(env, SIMPLE_OR_STRING, mut),
		number:         createCellTypeVal(env, NUMBER, mut),
		simpleBasic:    createCellTypeVal(env, SIMPLE_BASIC, mut),
	}
}

func createCellTypeVal(env Env, ty BasicTypeBitSet, mut CellMutability) ComplexSemType {
	atomicCell := cellAtomicTypeFrom(ty, mut)
	atom := env.cellAtom(&atomicCell)
	bdd := bddAtom(&atom)
	return getBasicSubtype(BTCell, bdd)
}

func (p preallocatedTypeVals) basicTypeCell(ty BasicTypeBitSet, mut CellMutability) (ComplexSemType, bool) {
	switch mut {
	case CellMutability_CELL_MUT_NONE:
		return p.none.basicTypeCell(ty)
	case CellMutability_CELL_MUT_LIMITED:
		return p.limited.basicTypeCell(ty)
	case CellMutability_CELL_MUT_UNLIMITED:
		return p.unlimited.basicTypeCell(ty)
	default:
		return ComplexSemType{}, false
	}
}

func (b basicCellTypeVals) basicTypeCell(ty BasicTypeBitSet) (ComplexSemType, bool) {
	switch ty {
	case NEVER:
		return b.never, true
	case NIL:
		return b.nil, true
	case BOOLEAN:
		return b.boolean, true
	case INT:
		return b.int, true
	case FLOAT:
		return b.float, true
	case DECIMAL:
		return b.decimal, true
	case STRING:
		return b.string, true
	case ERROR:
		return b.error, true
	case LIST:
		return b.list, true
	case MAPPING:
		return b.mapping, true
	case TABLE:
		return b.table, true
	case UNDEF:
		return b.undef, true
	case REGEXP:
		return b.regexp, true
	case FUNCTION:
		return b.function, true
	case TYPEDESC:
		return b.typedesc, true
	case HANDLE:
		return b.handle, true
	case XML:
		return b.xml, true
	case OBJECT:
		return b.object, true
	case STREAM:
		return b.stream, true
	case FUTURE:
		return b.future, true
	case VAL:
		return b.val, true
	case INNER:
		return b.inner, true
	case ANY:
		return b.any, true
	case SIMPLE_OR_STRING:
		return b.simpleOrString, true
	case NUMBER:
		return b.number, true
	case SIMPLE_BASIC:
		return b.simpleBasic, true
	default:
		return ComplexSemType{}, false
	}
}
