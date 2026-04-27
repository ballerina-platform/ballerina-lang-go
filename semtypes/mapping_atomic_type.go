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
	"slices"
)

type MappingAtomicType struct {
	Names []string
	Types []ComplexSemType
	Rest  *ComplexSemType
}

var _ atomicType = &MappingAtomicType{}

func (m *MappingAtomicType) equals(other atomicType) bool {
	if other, ok := other.(*MappingAtomicType); ok {
		if !m.Rest.equals(other.Rest) {
			return false
		}
		return slices.Equal(other.Names, m.Names) && slices.EqualFunc(other.Types, m.Types, func(a, b ComplexSemType) bool { return a.equals(&b) })
	}
	return false
}

func mappingAtomicTypeFrom(names []string, types []ComplexSemType, rest *ComplexSemType) MappingAtomicType {
	return MappingAtomicType{
		Names: names,
		Types: types,
		Rest:  rest,
	}
}

func (m *MappingAtomicType) atomKind() kind {
	return kind_MAPPING_ATOM
}

func (m *MappingAtomicType) FieldInnerVal(name string) SemType {
	for i, n := range m.Names {
		if n == name {
			return cellInnerVal(&m.Types[i])
		}
	}
	return cellInnerVal(m.Rest)
}

func (m *MappingAtomicType) IsOptional(cx Context, name string) bool {
	for i, n := range m.Names {
		if n == name {
			return IsSubtype(cx, UNDEF, cellInnerVal(&m.Types[i]))
		}
	}
	return true
}
