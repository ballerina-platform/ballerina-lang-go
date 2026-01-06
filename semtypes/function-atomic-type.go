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

type FunctionAtomicType struct {
	ParamType  SemType
	RetType    SemType
	Qualifiers SemType
	IsGeneric  bool
}

var _ AtomicType = &FunctionAtomicType{}

func FunctionAtomicTypeFrom(paramType SemType, rest SemType, qualifiers SemType) FunctionAtomicType {
	// migrated from FunctionAtomicType.java:32:5
	// FIXME: failed to find constructor for FunctionAtomicType

	return NewFunctionAtomicType(paramType, rest, qualifiers, false)
}

func FunctionAtomicTypeGenericFrom(paramType SemType, rest SemType, qualifiers SemType) FunctionAtomicType {
	// migrated from FunctionAtomicType.java:36:5
	// FIXME: failed to find constructor for FunctionAtomicType

	return NewFunctionAtomicType(paramType, rest, qualifiers, true)
}

func NewFunctionAtomicType(paramType SemType, retType SemType, qualifiers SemType, isGeneric bool) FunctionAtomicType {
	this := FunctionAtomicType{}
	this.ParamType = paramType
	this.RetType = retType
	this.Qualifiers = qualifiers
	this.IsGeneric = isGeneric
	return this
}

func (this *FunctionAtomicType) AtomKind() Kind {
	// migrated from FunctionAtomicType.java:40:5
	return Kind_FUNCTION_ATOM
}
