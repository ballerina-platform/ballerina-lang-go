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

type Env interface {
	cellAtom(atomicType CellAtomicType) TypeAtom
	recFunctionAtomType() *RecAtom
	recMappingAtomType() *RecAtom
	recListAtomType() *RecAtom
	setRecFunctionAtomType(rec *RecAtom, atomicType FunctionAtomicType)
	setRecMappingAtomType(rec *RecAtom, atomicType MappingAtomicType)
	setRecListAtomType(rec *RecAtom, atomicType ListAtomicType)
	functionAtom(atomicType FunctionAtomicType) *TypeAtom
	mappingAtom(atomicType MappingAtomicType) *TypeAtom
	listAtom(atomicType ListAtomicType) *TypeAtom
	recListAtom() *RecAtom
	recMappingAtom() *RecAtom
	initializeFromPredefinedTypeEnv(PredefinedTypeEnv *PredefinedTypeEnv)
	mappingAtomType(atom Atom) *MappingAtomicType
	functionAtomType(atom Atom) *FunctionAtomicType
	listAtomType(atom Atom) *ListAtomicType
}

// Public/package methods - migrated from PredefinedTypeEnv.java:606-644

// initializeEnv populates the environment with predefined atoms
// migrated from PredefinedTypeEnv.java:606-611
// func (this *PredefinedTypeEnv) initializeEnv(env Env) {
// 	fillRecAtoms(this, &env.recListAtoms, this.initializedRecListAtoms)
// 	fillRecAtoms(this, &env.recMappingAtoms, this.initializedRecMappingAtoms)
// 	for _, each := range this.initializedCellAtoms {
// 		env.cellAtom(each.atomicType)
// 	}
// 	for _, each := range this.initializedListAtoms {
// 		env.listAtom(each.atomicType)
// 	}
// }

// fillRecAtoms fills the environment rec atom list with initialized rec atoms
// migrated from PredefinedTypeEnv.java:613-624
func fillRecAtoms[E AtomicType](env *PredefinedTypeEnv, envRecAtomList *[]E, initializedRecAtoms []E) {
	count := env.ReservedRecAtomCount()
	for i := 0; i < count; i++ {
		if i < len(initializedRecAtoms) {
			*envRecAtomList = append(*envRecAtomList, initializedRecAtoms[i])
		} else {
			var zero E
			*envRecAtomList = append(*envRecAtomList, zero)
		}
	}
}
