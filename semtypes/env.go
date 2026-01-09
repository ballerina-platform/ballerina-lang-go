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

// migration-note: we can turning this to an interface to avoid accidentally copying the env
type Env interface {
	cellAtom(atomicType *CellAtomicType) TypeAtom
	recFunctionAtom() RecAtom
	recMappingAtom() RecAtom
	recListAtom() RecAtom
	setRecFunctionAtomType(rec RecAtom, atomicType *FunctionAtomicType)
	setRecMappingAtomType(rec RecAtom, atomicType *MappingAtomicType)
	setRecListAtomType(rec RecAtom, atomicType *ListAtomicType)
	functionAtom(atomicType *FunctionAtomicType) TypeAtom
	mappingAtom(atomicType *MappingAtomicType) TypeAtom
	listAtom(atomicType *ListAtomicType) TypeAtom
	mappingAtomType(atom Atom) *MappingAtomicType
	functionAtomType(atom Atom) *FunctionAtomicType
	listAtomType(atom Atom) *ListAtomicType
}

var typeEnv Env = nil
var typeEnvInitializer sync.Once

func GetTypeEnv() Env {
	typeEnvInitializer.Do(func() {
		env := &envImpl{
			atomTable: make(map[AtomicType]TypeAtom),
			types:     make(map[string]SemType),
		}
		fillRecAtoms(predefinedTypeEnv, &env.recListAtoms, predefinedTypeEnv.initializedRecListAtoms)
		fillRecAtoms(predefinedTypeEnv, &env.recMappingAtoms, predefinedTypeEnv.initializedRecMappingAtoms)
		for _, each := range predefinedTypeEnv.initializedCellAtoms {
			env.cellAtom(each.atomicType)
		}
		for _, each := range predefinedTypeEnv.initializedListAtoms {
			env.listAtom(each.atomicType)
		}
		typeEnv = env
	})
	return typeEnv
}

type envImpl struct {
	recListAtoms []*ListAtomicType
	recListAtomsMutex sync.Mutex

	recMappingAtoms []*MappingAtomicType
	recMappingAtomsMutex sync.Mutex

	recFunctionAtoms []*FunctionAtomicType
	recFunctionAtomsMutex sync.Mutex

	distinctAtoms int
	distinctAtomMutex sync.Mutex
	// migration-note: unlike java implementation this will leak memory. So be careful about adding atoms in an unbounded way.
	atomTableMutex sync.Mutex
	atomTable map[AtomicType]TypeAtom

	types map[string]SemType
}
var _ Env = &envImpl{}

func (this *envImpl) recListAtomCount() int {
	return len(this.recListAtoms)
}

func (this *envImpl) recMappingAtomCount() int {
	return len(this.recMappingAtoms)
}

func (this *envImpl) recFunctionAtomCount() int {
	return len(this.recFunctionAtoms)
}

func (this *envImpl) distinctAtomCount() int {
	this.distinctAtomMutex.Lock()
	defer this.distinctAtomMutex.Unlock()
	return this.distinctAtoms
}

func (this *envImpl) distinctAtomCountGetAndIncrement() int {
	this.distinctAtomMutex.Lock()
	defer this.distinctAtomMutex.Unlock()
	this.distinctAtoms++
	return this.distinctAtoms
}

func (this *envImpl) recFunctionAtom() RecAtom {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	result := len(this.recFunctionAtoms)
	this.recFunctionAtoms = append(this.recFunctionAtoms, nil)
	return CreateRecAtom(result)
}

func (this *envImpl) setRecFunctionAtomType(rec RecAtom, atomicType *FunctionAtomicType) {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	rec.SetKind(Kind_FUNCTION_ATOM)
	this.recFunctionAtoms[rec.Index()] = atomicType
}

func (this *envImpl) getRecFunctionAtomType(rec RecAtom) *FunctionAtomicType {
	this.recFunctionAtomsMutex.Lock()
	defer this.recFunctionAtomsMutex.Unlock()
	return this.recFunctionAtoms[rec.Index()]
}

func (this *envImpl) listAtom(atomicType *ListAtomicType) TypeAtom {
	return this.typeAtom(atomicType)
}

func (this *envImpl) mappingAtom(atomicType *MappingAtomicType) TypeAtom {
	return this.typeAtom(atomicType)
}

func (this *envImpl) functionAtom(atomicType *FunctionAtomicType) TypeAtom {
	return this.typeAtom(atomicType)
}

func (this *envImpl) cellAtom(atomicType *CellAtomicType) TypeAtom {
	return this.typeAtom(atomicType)
}

func (this *envImpl) typeAtom(atomicType AtomicType) TypeAtom {
	this.atomTableMutex.Lock()
	defer this.atomTableMutex.Unlock()
	ta, ok := this.atomTable[atomicType]
	if ok {
		return ta
	}
	ta = CreateTypeAtom(len(this.atomTable), atomicType)
	this.atomTable[atomicType] = ta
	return ta
}

func (this *envImpl) listAtomType(atom Atom) *ListAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return this.getRecListAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*ListAtomicType)
}

func (this *envImpl) functionAtomType(atom Atom) *FunctionAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return this.getRecFunctionAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*FunctionAtomicType)
}

func (this *envImpl) mappingAtomType(atom Atom) *MappingAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return this.getRecMappingAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*MappingAtomicType)
}

func (this *envImpl) recListAtom() RecAtom {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	result := len(this.recListAtoms)
	this.recListAtoms = append(this.recListAtoms, nil)
	return CreateRecAtom(result)
}

func (this *envImpl) recMappingAtom() RecAtom {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	result := len(this.recMappingAtoms)
	this.recMappingAtoms = append(this.recMappingAtoms, nil)
	return CreateRecAtom(result)
}

func (this *envImpl) setRecListAtomType(rec RecAtom, atomicType *ListAtomicType) {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	rec.SetKind(Kind_LIST_ATOM)
	this.recListAtoms[rec.Index()] = atomicType
}

func (this *envImpl) setRecMappingAtomType(rec RecAtom, atomicType *MappingAtomicType) {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	rec.SetKind(Kind_MAPPING_ATOM)
	this.recMappingAtoms[rec.Index()] = atomicType
}

func (this *envImpl) getRecListAtomType(rec RecAtom) *ListAtomicType {
	this.recListAtomsMutex.Lock()
	defer this.recListAtomsMutex.Unlock()
	return this.recListAtoms[rec.Index()]
}

func (this *envImpl) getRecMappingAtomType(rec RecAtom) *MappingAtomicType {
	this.recMappingAtomsMutex.Lock()
	defer this.recMappingAtomsMutex.Unlock()
	return this.recMappingAtoms[rec.Index()]
}

func (this *envImpl) cellAtomType(atom Atom) *CellAtomicType {
	return atom.(*TypeAtom).AtomicType.(*CellAtomicType)
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
