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

func CreateTypeEnv() Env {
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
	return env
}

type envImpl struct {
	recListAtoms      []*ListAtomicType
	recListAtomsMutex sync.Mutex

	recMappingAtoms      []*MappingAtomicType
	recMappingAtomsMutex sync.Mutex

	recFunctionAtoms      []*FunctionAtomicType
	recFunctionAtomsMutex sync.Mutex

	distinctAtoms     int
	distinctAtomMutex sync.Mutex
	// migration-note: unlike java implementation this will leak memory. So be careful about adding atoms in an unbounded way.
	atomTableMutex sync.Mutex
	atomTable      map[AtomicType]TypeAtom

	types map[string]SemType
}

var _ Env = &envImpl{}

func (e *envImpl) recListAtomCount() int {
	return len(e.recListAtoms)
}

func (e *envImpl) recMappingAtomCount() int {
	return len(e.recMappingAtoms)
}

func (e *envImpl) recFunctionAtomCount() int {
	return len(e.recFunctionAtoms)
}

func (e *envImpl) distinctAtomCount() int {
	e.distinctAtomMutex.Lock()
	defer e.distinctAtomMutex.Unlock()
	return e.distinctAtoms
}

func (e *envImpl) distinctAtomCountGetAndIncrement() int {
	e.distinctAtomMutex.Lock()
	defer e.distinctAtomMutex.Unlock()
	e.distinctAtoms++
	return e.distinctAtoms
}

func (e *envImpl) recFunctionAtom() RecAtom {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	result := len(e.recFunctionAtoms)
	e.recFunctionAtoms = append(e.recFunctionAtoms, nil)
	return CreateRecAtom(result)
}

func (e *envImpl) setRecFunctionAtomType(rec RecAtom, atomicType *FunctionAtomicType) {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	rec.SetKind(Kind_FUNCTION_ATOM)
	e.recFunctionAtoms[rec.Index()] = atomicType
}

func (e *envImpl) getRecFunctionAtomType(rec RecAtom) *FunctionAtomicType {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	return e.recFunctionAtoms[rec.Index()]
}

func (e *envImpl) listAtom(atomicType *ListAtomicType) TypeAtom {
	return e.typeAtom(atomicType)
}

func (e *envImpl) mappingAtom(atomicType *MappingAtomicType) TypeAtom {
	return e.typeAtom(atomicType)
}

func (e *envImpl) functionAtom(atomicType *FunctionAtomicType) TypeAtom {
	return e.typeAtom(atomicType)
}

func (e *envImpl) cellAtom(atomicType *CellAtomicType) TypeAtom {
	return e.typeAtom(atomicType)
}

func (e *envImpl) typeAtom(atomicType AtomicType) TypeAtom {
	e.atomTableMutex.Lock()
	defer e.atomTableMutex.Unlock()
	ta, ok := e.atomTable[atomicType]
	if ok {
		return ta
	}
	ta = CreateTypeAtom(len(e.atomTable), atomicType)
	e.atomTable[atomicType] = ta
	return ta
}

func (e *envImpl) listAtomType(atom Atom) *ListAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return e.getRecListAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*ListAtomicType)
}

func (e *envImpl) functionAtomType(atom Atom) *FunctionAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return e.getRecFunctionAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*FunctionAtomicType)
}

func (e *envImpl) mappingAtomType(atom Atom) *MappingAtomicType {
	if recAtom, ok := atom.(*RecAtom); ok {
		return e.getRecMappingAtomType(*recAtom)
	}
	return atom.(*TypeAtom).AtomicType.(*MappingAtomicType)
}

func (e *envImpl) recListAtom() RecAtom {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	result := len(e.recListAtoms)
	e.recListAtoms = append(e.recListAtoms, nil)
	return CreateRecAtom(result)
}

func (e *envImpl) recMappingAtom() RecAtom {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	result := len(e.recMappingAtoms)
	e.recMappingAtoms = append(e.recMappingAtoms, nil)
	return CreateRecAtom(result)
}

func (e *envImpl) setRecListAtomType(rec RecAtom, atomicType *ListAtomicType) {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	rec.SetKind(Kind_LIST_ATOM)
	e.recListAtoms[rec.Index()] = atomicType
}

func (e *envImpl) setRecMappingAtomType(rec RecAtom, atomicType *MappingAtomicType) {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	rec.SetKind(Kind_MAPPING_ATOM)
	e.recMappingAtoms[rec.Index()] = atomicType
}

func (e *envImpl) getRecListAtomType(rec RecAtom) *ListAtomicType {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	return e.recListAtoms[rec.Index()]
}

func (e *envImpl) getRecMappingAtomType(rec RecAtom) *MappingAtomicType {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	return e.recMappingAtoms[rec.Index()]
}

func (e *envImpl) cellAtomType(atom Atom) *CellAtomicType {
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
	for i := range count {
		if i < len(initializedRecAtoms) {
			*envRecAtomList = append(*envRecAtomList, initializedRecAtoms[i])
		} else {
			var zero E
			*envRecAtomList = append(*envRecAtomList, zero)
		}
	}
}
