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

type FunctionDefinition struct {
	rec     *RecAtom
	semType SemType
}

var _ Definition = &FunctionDefinition{}

func NewFunctionDefinition() FunctionDefinition {
	this := FunctionDefinition{}
	return this
}

func (this *FunctionDefinition) GetSemType(env Env) SemType {
	// migrated from FunctionDefinition.java:43:5
	if this.semType != nil {
		return this.semType
	}
	rec := env.recFunctionAtom()
	this.rec = &rec
	return this.createSemType(&rec)
}

func (this *FunctionDefinition) createSemType(rec Atom) SemType {
	// migrated from FunctionDefinition.java:53:5
	bdd := BddAtom(rec)
	s := basicSubtype(BT_FUNCTION, bdd)
	this.semType = s
	return s
}

func (this *FunctionDefinition) Define(env Env, args SemType, ret SemType, qualifiers FunctionQualifiers) SemType {
	// migrated from FunctionDefinition.java:60:5
	atomicType := FunctionAtomicTypeFrom(args, ret, qualifiers.semType)
	return this.defineInternal(env, atomicType)
}

func (this *FunctionDefinition) DefineGeneric(env Env, args SemType, ret SemType, qualifiers FunctionQualifiers) SemType {
	// migrated from FunctionDefinition.java:65:5
	atomicType := FunctionAtomicTypeGenericFrom(args, ret, qualifiers.semType)
	return this.defineInternal(env, atomicType)
}

func (this *FunctionDefinition) defineInternal(env Env, atomicType FunctionAtomicType) SemType {
	// migrated from FunctionDefinition.java:70:5
	var atom Atom
	rec := this.rec
	if rec != nil {
		atom = rec
		env.setRecFunctionAtomType(*rec, &atomicType)
	} else {
		atom = common.ToPointer(env.functionAtom(&atomicType))
	}
	return this.createSemType(atom)
}
