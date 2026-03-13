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

type FunctionDefinition struct {
	rec     *RecAtom
	semType SemType
}

var _ Definition = &FunctionDefinition{}

func NewFunctionDefinition() FunctionDefinition {
	this := FunctionDefinition{}
	return this
}

func (f *FunctionDefinition) GetSemType(env Env) SemType {
	// migrated from FunctionDefinition.java:43:5
	if f.semType != nil {
		return f.semType
	}
	rec := env.recFunctionAtom()
	f.rec = &rec
	return f.createSemType(&rec)
}

func (f *FunctionDefinition) createSemType(rec Atom) SemType {
	// migrated from FunctionDefinition.java:53:5
	bdd := BddAtom(rec)
	s := basicSubtype(BTFunction, bdd)
	f.semType = s
	return s
}

func (f *FunctionDefinition) Define(env Env, args SemType, ret SemType, qualifiers FunctionQualifiers) SemType {
	// migrated from FunctionDefinition.java:60:5
	atomicType := FunctionAtomicTypeFrom(args, ret, qualifiers.semType)
	return f.defineInternal(env, atomicType)
}

func (f *FunctionDefinition) DefineGeneric(env Env, args SemType, ret SemType, qualifiers FunctionQualifiers) SemType {
	// migrated from FunctionDefinition.java:65:5
	atomicType := FunctionAtomicTypeGenericFrom(args, ret, qualifiers.semType)
	return f.defineInternal(env, atomicType)
}

func (f *FunctionDefinition) defineInternal(env Env, atomicType FunctionAtomicType) SemType {
	// migrated from FunctionDefinition.java:70:5
	var atom Atom
	rec := f.rec
	if rec != nil {
		atom = rec
		env.setRecFunctionAtomType(*rec, &atomicType)
	} else {
		atom = new(env.functionAtom(&atomicType))
	}
	return f.createSemType(atom)
}
