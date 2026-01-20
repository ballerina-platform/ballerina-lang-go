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

type TableSubtype struct {
}

func newTableSubtype() TableSubtype {
	this := TableSubtype{}
	return this
}

func TableContainingKeyConstraint(cx Context, tableConstraint SemType, keyConstraint SemType) SemType {
	// migrated from TableSubtype.java:48:5
	var normalizedKc SemType
	lat := listAtomicType(cx, keyConstraint)
	if (lat != nil) && (CELL_ATOMIC_UNDEF == cellAtomicType(lat.Rest)) {
		members := lat.Members
		switch members.FixedLength {
		case 0:
			normalizedKc = &VAL
		case 1:
			normalizedKc = cellAtomicType(members.Initial[0]).Ty
		default:
			normalizedKc = keyConstraint
		}
	} else {
		normalizedKc = keyConstraint
	}
	return tableContainingWithEnvSemTypeSemTypeSemType(cx.env(), tableConstraint, normalizedKc, &VAL)
}

func TableContainingKeySpecifier(cx Context, tableConstraint SemType, fieldNames []string) SemType {
	// migrated from TableSubtype.java:64:5
	fieldNameSingletons := make([]SemType, len(fieldNames))
	fieldTypes := make([]SemType, len(fieldNames))
	for i := 0; i < len(fieldNames); i++ {
		key := StringConst(fieldNames[i])
		fieldNameSingletons[i] = key
		fieldTypes[i] = MappingMemberTypeInnerVal(cx, tableConstraint, key)
	}
	listDef1 := NewListDefinition()
	normalizedKs := listDef1.TupleTypeWrapped(cx.env(), fieldNameSingletons...)
	var normalizedKc SemType
	if len(fieldTypes) > 1 {
		ld := NewListDefinition()
		normalizedKc = ld.TupleTypeWrapped(cx.env(), fieldTypes...)
	} else {
		normalizedKc = fieldTypes[0]
	}
	return tableContainingWithEnvSemTypeSemTypeSemType(cx.env(), tableConstraint, normalizedKc, normalizedKs)
}

func TableContaining(env Env, tableConstraint SemType) SemType {
	// migrated from TableSubtype.java:85:5
	return TableContainingWithEnvSemTypeCellMutability(env, tableConstraint, CellMutability_CELL_MUT_LIMITED)
}

func TableContainingWithEnvSemTypeCellMutability(env Env, tableConstraint SemType, mut CellMutability) SemType {
	// migrated from TableSubtype.java:89:5
	var normalizedKc SemType = &VAL
	var normalizedKs SemType = &VAL
	return tableContaining(env, tableConstraint, normalizedKc, normalizedKs, mut)
}

func tableContaining(env Env, tableConstraint SemType, normalizedKc SemType, normalizedKs SemType, mut CellMutability) SemType {
	// migrated from TableSubtype.java:95:5
	if !IsSubtypeSimple(tableConstraint, MAPPING) {
		panic("assertion failed")
	}
	typeParamArrDef := NewListDefinition()
	typeParamArray := typeParamArrDef.DefineListTypeWrappedWithEnvSemTypeCellMutability(env, tableConstraint, mut)
	listDef := NewListDefinition()
	tupleType := listDef.TupleTypeWrapped(env, typeParamArray, normalizedKc, normalizedKs)
	bdd := subtypeData(tupleType, BT_LIST).(Bdd)
	return CreateBasicSemType(BT_TABLE, bdd)
}

func tableContainingWithEnvSemTypeSemTypeSemType(env Env, tableConstraint SemType, normalizedKc SemType, normalizedKs SemType) SemType {
	// migrated from TableSubtype.java:109:5
	return tableContaining(env, tableConstraint, normalizedKc, normalizedKs, CellMutability_CELL_MUT_LIMITED)
}
