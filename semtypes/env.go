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

import "sync"

// Env is an opaque pointer to the type environment. All (potentially recursive) types are defined in a type environment.
// Performing type operations involving types defined in different type environments is an undefined behaviour. Therefore it
// is advisable to use the same type environment during the full execution of an interpreter process.
type Env = *env

func CreateTypeEnv() Env {
	env := &env{
		atomTable: make(map[atomicType]typeAtom),
		types:     make(map[string]SemType),
	}
	fillRecAtoms(predefTypeEnv, &env.recListAtoms, predefTypeEnv.initializedRecListAtoms)
	fillRecAtoms(predefTypeEnv, &env.recMappingAtoms, predefTypeEnv.initializedRecMappingAtoms)
	for _, each := range predefTypeEnv.initializedCellAtoms {
		env.cellAtom(each.atomicType)
	}
	for _, each := range predefTypeEnv.initializedListAtoms {
		env.listAtom(each.atomicType)
	}
	return env
}

type env struct {
	recListAtoms      []*ListAtomicType
	recListAtomsMutex sync.Mutex

	recMappingAtoms      []*MappingAtomicType
	recMappingAtomsMutex sync.Mutex

	recFunctionAtoms      []*functionAtomicType
	recFunctionAtomsMutex sync.Mutex

	distinctAtoms     int
	distinctAtomMutex sync.Mutex
	// migration-note: unlike java implementation this will leak memory. So be careful about adding atoms in an unbounded way.
	atomTableMutex sync.Mutex
	atomTable      map[atomicType]typeAtom

	types map[string]SemType
}

func (this *env) recListAtomCount() int {
	return len(this.recListAtoms)
}

func (this *env) recMappingAtomCount() int {
	return len(this.recMappingAtoms)
}

func (this *env) recFunctionAtomCount() int {
	return len(this.recFunctionAtoms)
}

func (this *env) distinctAtomCount() int {
	this.distinctAtomMutex.Lock()
	defer this.distinctAtomMutex.Unlock()
	return this.distinctAtoms
}

func (this *env) distinctAtomCountGetAndIncrement() int {
	this.distinctAtomMutex.Lock()
	defer this.distinctAtomMutex.Unlock()
	this.distinctAtoms++
	return this.distinctAtoms
}

func (this *env) recFunctionAtom() recAtom {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	result := len(this.recFunctionAtoms)
	this.recFunctionAtoms = append(this.recFunctionAtoms, nil)
	return createRecAtom(result)
}

func (this *env) setRecFunctionAtomType(rec recAtom, atomicType *functionAtomicType) {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	this.recFunctionAtoms[rec.index()] = atomicType
}

func (this *env) getRecFunctionAtomType(rec recAtom) *functionAtomicType {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	return this.recFunctionAtoms[rec.index()]
}

func (this *env) listAtom(atomicType *ListAtomicType) typeAtom {
	return this.typeAtom(atomicType)
}

func (this *env) mappingAtom(atomicType *MappingAtomicType) typeAtom {
	return this.typeAtom(atomicType)
}

func (this *env) functionAtom(atomicType *functionAtomicType) typeAtom {
	return this.typeAtom(atomicType)
}

func (this *env) cellAtom(atomicType *cellAtomicType) typeAtom {
	return this.typeAtom(atomicType)
}

func (this *env) typeAtom(atomicType atomicType) typeAtom {
	this.atomTableMutex.Lock()
	defer this.atomTableMutex.Unlock()
	ta, ok := this.atomTable[atomicType]
	if ok {
		return ta
	}
	ta = createTypeAtom(len(this.atomTable), atomicType)
	this.atomTable[atomicType] = ta
	return ta
}

func (this *env) listAtomType(atom atom) *ListAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return this.getRecListAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*ListAtomicType)
}

func (this *env) functionAtomType(atom atom) *functionAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return this.getRecFunctionAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*functionAtomicType)
}

func (this *env) mappingAtomType(atom atom) *MappingAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return this.getRecMappingAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*MappingAtomicType)
}

func (this *env) recListAtom() recAtom {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	result := len(this.recListAtoms)
	this.recListAtoms = append(this.recListAtoms, nil)
	return createRecAtom(result)
}

func (this *env) recMappingAtom() recAtom {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	result := len(this.recMappingAtoms)
	this.recMappingAtoms = append(this.recMappingAtoms, nil)
	return createRecAtom(result)
}

func (this *env) setRecListAtomType(rec recAtom, atomicType *ListAtomicType) {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	this.recListAtoms[rec.index()] = atomicType
}

func (this *env) setRecMappingAtomType(rec recAtom, atomicType *MappingAtomicType) {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	this.recMappingAtoms[rec.index()] = atomicType
}

func (this *env) getRecListAtomType(rec recAtom) *ListAtomicType {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	return this.recListAtoms[rec.index()]
}

func (this *env) getRecMappingAtomType(rec recAtom) *MappingAtomicType {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	return this.recMappingAtoms[rec.index()]
}

func (this *env) cellAtomType(atom atom) *cellAtomicType {
	return atom.(*typeAtom).AtomicType.(*cellAtomicType)
}

// Public/package methods

// initializeEnv populates the environment with predefined atoms
// func (this *predefinedTypeEnv) initializeEnv(env Env) {
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
func fillRecAtoms[E atomicType](env *predefinedTypeEnv, envRecAtomList *[]E, initializedRecAtoms []E) {
	count := env.ReservedRecAtomCount()
	for i := range count {
		if i < len(initializedRecAtoms) {
			*envRecAtomList = append(*envRecAtomList, initializedRecAtoms[i])
		} else {
			var zero E
			*envRecAtomList = append(*envRecAtomList, zero)
		}
	}
}
