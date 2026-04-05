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

import "ballerina-lang-go/common"

type ListDefinition struct {
	rec     *recAtom
	semType SemType
}

var _ Definition = &ListDefinition{}

func NewListDefinition() ListDefinition {
	this := ListDefinition{}
	this.rec = nil
	this.semType = nil
	// Default field initializations

	return this
}

func (this *ListDefinition) GetSemType(env Env) SemType {
	s := this.semType
	if s == nil {
		rec := env.recListAtom()
		this.rec = &rec
		return this.createSemType(env, &rec)
	} else {
		return s
	}
}

func (this *ListDefinition) TupleTypeWrapped(env Env, members ...SemType) SemType {
	return this.DefineListTypeWrappedWithEnvSemTypesInt(env, members, len(members))
}

func (this *ListDefinition) TupleTypeWrappedRo(env Env, members ...SemType) SemType {
	return this.DefineListTypeWrapped(env, members, len(members), NEVER, CellMutability_CELL_MUT_NONE)
}

func (this *ListDefinition) DefineListTypeWrapped(env Env, initial []SemType, fixedLength int, rest SemType, mut CellMutability) SemType {
	common.Assert(rest != nil)
	var initialCells []ComplexSemType
	for _, member := range initial {
		initialCells = append(initialCells, *cellContainingWithEnvSemTypeCellMutability(env, member, mut))
	}
	var restMut CellMutability
	if IsNever(rest) {
		restMut = CellMutability_CELL_MUT_NONE
	} else {
		restMut = mut
	}
	restCell := cellContainingWithEnvSemTypeCellMutability(env, Union(rest, UNDEF), restMut)
	return this.define(env, initialCells, fixedLength, restCell)
}

func (this *ListDefinition) DefineListTypeWrappedWithEnvSemTypesInt(env Env, initial []SemType, size int) SemType {
	return this.DefineListTypeWrapped(env, initial, size, NEVER, CellMutability_CELL_MUT_LIMITED)
}

func (this *ListDefinition) DefineListTypeWrappedWithEnvSemTypesIntSemType(env Env, initial []SemType, fixedLength int, rest SemType) SemType {
	return this.DefineListTypeWrapped(env, initial, fixedLength, rest, CellMutability_CELL_MUT_LIMITED)
}

func (this *ListDefinition) DefineListTypeWrappedWithEnvSemType(env Env, rest SemType) SemType {
	return this.DefineListTypeWrappedWithEnvSemTypesIntSemType(env, nil, 0, rest)
}

func (this *ListDefinition) DefineListTypeWrappedWithEnvSemTypeCellMutability(env Env, rest SemType, mut CellMutability) SemType {
	return this.DefineListTypeWrapped(env, nil, 0, rest, mut)
}

func (this *ListDefinition) DefineListTypeWrappedWithEnvSemTypesSemType(env Env, initial []SemType, rest SemType) SemType {
	return this.DefineListTypeWrapped(env, initial, len(initial), rest, CellMutability_CELL_MUT_LIMITED)
}

func (this *ListDefinition) define(env Env, initial []ComplexSemType, fixedLength int, rest *ComplexSemType) *ComplexSemType {
	members := this.fixedLengthNormalize(fixedLengthArrayFrom(initial, fixedLength))
	atomicType := listAtomicTypeFrom(members, rest)
	var atom atom
	rec := this.rec
	if rec != nil {
		atom = rec
		env.setRecListAtomType(*rec, &atomicType)
	} else {
		atom = new(env.listAtom(&atomicType))
	}
	return this.createSemType(env, atom)
}

func (this *ListDefinition) fixedLengthNormalize(array fixedLengthArray) fixedLengthArray {
	initial := array.Initial
	i := (len(initial) - 1)
	if i <= 0 {
		return array
	}
	last := initial[i]
	i = (i - 1)
	for i >= 0 {
		if !last.equals(&initial[i]) {
			break
		}
		i = (i - 1)
	}
	return fixedLengthArrayFrom(initial[:i+2], array.FixedLength)
}

func (this *ListDefinition) createSemType(env Env, atom atom) *ComplexSemType {
	bdd := bddAtom(atom)
	complexSemType := getBasicSubtype(BTList, bdd)
	this.semType = complexSemType
	return complexSemType
}
