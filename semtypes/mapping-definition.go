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
	"ballerina-lang-go/common"
	"sort"
)

type MappingDefinition struct {
	rec     *RecAtom
	semType SemType
}

var _ Definition = &MappingDefinition{}

func fieldName(f CellField) string {
	// migrated from MappingDefinition.java:135:5
	return f.Name
}

func NewMappingDefinition() MappingDefinition {
	this := MappingDefinition{}
	this.rec = nil
	this.semType = nil
	// Default field initializations

	return this
}

func (this *MappingDefinition) GetSemType(env Env) SemType {
	// migrated from MappingDefinition.java:56:5
	s := this.semType
	if s == nil {
		rec := env.recMappingAtom()
		this.rec = &rec
		return this.createSemType(env, &rec)
	} else {
		return s
	}
}

func (this *MappingDefinition) SetSemTypeToNever() {
	// migrated from MappingDefinition.java:72:5
	this.semType = &NEVER
}

func (this *MappingDefinition) Define(env Env, fields []CellField, rest CellSemType) SemType {
	// migrated from MappingDefinition.java:76:5
	sfh := this.splitFields(fields)
	atomicType := MappingAtomicTypeFrom(sfh.Names, sfh.Types, rest)
	var atom Atom
	rec := this.rec
	if rec != nil {
		atom = rec
		env.setRecMappingAtomType(*rec, &atomicType)
	} else {
		atom = common.ToPointer(env.mappingAtom(&atomicType))
	}
	return this.createSemType(env, atom)
}

func (this *MappingDefinition) DefineMappingTypeWrapped(env Env, fields []Field, rest SemType) SemType {
	// migrated from MappingDefinition.java:91:5
	return this.DefineMappingTypeWrappedWithEnvFieldsSemTypeCellMutability(env, fields, rest, CellMutability_CELL_MUT_LIMITED)
}

func (this *MappingDefinition) DefineMappingTypeWrappedWithEnvFieldsSemTypeCellMutability(env Env, fields []Field, rest SemType, mut CellMutability) SemType {
	// migrated from MappingDefinition.java:95:5
	var cellFields []CellField
	for _, field := range fields {
		ty := field.Ty
		var optTy SemType
		if field.Opt {
			optTy = Union(ty, &UNDEF)
		} else {
			optTy = ty
		}
		var ro CellMutability
		if field.Ro {
			ro = CellMutability_CELL_MUT_NONE
		} else {
			ro = mut
		}
		cellFields = append(cellFields, CellFieldFrom(field.Name, CellContainingWithEnvSemTypeCellMutability(env, optTy, ro)))
	}
	var restMut CellMutability
	if IsNever(rest) {
		restMut = CellMutability_CELL_MUT_NONE
	} else {
		restMut = mut
	}
	restCell := CellContainingWithEnvSemTypeCellMutability(env, Union(rest, &UNDEF), restMut)
	return this.Define(env, cellFields, restCell)
}

func (this *MappingDefinition) createSemType(env Env, atom Atom) SemType {
	// migrated from MappingDefinition.java:116:5
	bdd := BddAtom(atom)
	s := basicSubtype(BT_MAPPING, bdd)
	this.semType = s
	return s
}

func (this *MappingDefinition) splitFields(fields []CellField) SplitField {
	// migrated from MappingDefinition.java:123:5
	sortedFields := make([]CellField, len(fields))
	copy(sortedFields, fields)
	// Arrays.sort(sortedFields, Comparator.comparing(MappingDefinition::fieldName))
	sort.Slice(sortedFields, func(i, j int) bool {
		return fieldName(sortedFields[i]) < fieldName(sortedFields[j])
	})
	var names []string
	var types []CellSemType
	for _, field := range sortedFields {
		names = append(names, field.Name)
		types = append(types, field.Type)
	}
	return SplitFieldFrom(names, types)
}
