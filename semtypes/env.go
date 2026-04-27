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

func (e *env) recListAtomCount() int {
	return len(e.recListAtoms)
}

func (e *env) recMappingAtomCount() int {
	return len(e.recMappingAtoms)
}

func (e *env) recFunctionAtomCount() int {
	return len(e.recFunctionAtoms)
}

func (e *env) distinctAtomCount() int {
	e.distinctAtomMutex.Lock()
	defer e.distinctAtomMutex.Unlock()
	return e.distinctAtoms
}

func (e *env) distinctAtomCountGetAndIncrement() int {
	e.distinctAtomMutex.Lock()
	defer e.distinctAtomMutex.Unlock()
	e.distinctAtoms++
	return e.distinctAtoms
}

func (e *env) recFunctionAtom() recAtom {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	result := len(e.recFunctionAtoms)
	e.recFunctionAtoms = append(e.recFunctionAtoms, nil)
	return createRecAtom(result)
}

func (e *env) setRecFunctionAtomType(rec recAtom, atomicType *functionAtomicType) {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	e.recFunctionAtoms[rec.index()] = atomicType
}

func (e *env) getRecFunctionAtomType(rec recAtom) *functionAtomicType {
	e.recFunctionAtomsMutex.Lock()
	defer e.recFunctionAtomsMutex.Unlock()
	return e.recFunctionAtoms[rec.index()]
}

func (e *env) listAtom(atomicType *ListAtomicType) typeAtom {
	return e.typeAtom(atomicType)
}

func (e *env) mappingAtom(atomicType *MappingAtomicType) typeAtom {
	return e.typeAtom(atomicType)
}

func (e *env) functionAtom(atomicType *functionAtomicType) typeAtom {
	return e.typeAtom(atomicType)
}

func (e *env) cellAtom(atomicType *cellAtomicType) typeAtom {
	return e.typeAtom(atomicType)
}

func (e *env) typeAtom(atomicType atomicType) typeAtom {
	e.atomTableMutex.Lock()
	defer e.atomTableMutex.Unlock()
	ta, ok := e.atomTable[atomicType]
	if ok {
		return ta
	}
	ta = createTypeAtom(len(e.atomTable), atomicType)
	e.atomTable[atomicType] = ta
	return ta
}

func (e *env) listAtomType(atom atom) *ListAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return e.getRecListAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*ListAtomicType)
}

func (e *env) functionAtomType(atom atom) *functionAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return e.getRecFunctionAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*functionAtomicType)
}

func (e *env) mappingAtomType(atom atom) *MappingAtomicType {
	if recAtom, ok := atom.(*recAtom); ok {
		return e.getRecMappingAtomType(*recAtom)
	}
	return atom.(*typeAtom).AtomicType.(*MappingAtomicType)
}

func (e *env) recListAtom() recAtom {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	result := len(e.recListAtoms)
	e.recListAtoms = append(e.recListAtoms, nil)
	return createRecAtom(result)
}

func (e *env) recMappingAtom() recAtom {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	result := len(e.recMappingAtoms)
	e.recMappingAtoms = append(e.recMappingAtoms, nil)
	return createRecAtom(result)
}

func (e *env) setRecListAtomType(rec recAtom, atomicType *ListAtomicType) {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	e.recListAtoms[rec.index()] = atomicType
}

func (e *env) setRecMappingAtomType(rec recAtom, atomicType *MappingAtomicType) {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	e.recMappingAtoms[rec.index()] = atomicType
}

func (e *env) getRecListAtomType(rec recAtom) *ListAtomicType {
	e.recListAtomsMutex.Lock()
	defer e.recListAtomsMutex.Unlock()
	return e.recListAtoms[rec.index()]
}

func (e *env) getRecMappingAtomType(rec recAtom) *MappingAtomicType {
	e.recMappingAtomsMutex.Lock()
	defer e.recMappingAtomsMutex.Unlock()
	return e.recMappingAtoms[rec.index()]
}

func (e *env) cellAtomType(atom atom) *cellAtomicType {
	return atom.(*typeAtom).AtomicType.(*cellAtomicType)
}

// Public/package methods

// initializeEnv populates the environment with predefined atoms
// func (p *predefinedTypeEnv) initializeEnv(env Env) {
// 	fillRecAtoms(p, &env.recListAtoms, p.initializedRecListAtoms)
// 	fillRecAtoms(p, &env.recMappingAtoms, p.initializedRecMappingAtoms)
// 	for _, each := range p.initializedCellAtoms {
// 		env.cellAtom(each.atomicType)
// 	}
// 	for _, each := range p.initializedListAtoms {
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
