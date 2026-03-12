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
	"fmt"
	"slices"
	"strings"
)

type MappingAtomicType struct {
	Names []string
	Types []CellSemType
	Rest  CellSemType
}

var _ AtomicType = &MappingAtomicType{}

func (this *MappingAtomicType) equals(other AtomicType) bool {
	if other, ok := other.(*MappingAtomicType); ok {
		if other.Rest != this.Rest {
			return false
		}
		return slices.Equal(other.Names, this.Names) && slices.Equal(other.Types, this.Types)
	}
	return false
}

func MappingAtomicTypeFrom(names []string, types []CellSemType, rest CellSemType) MappingAtomicType {
	// migrated from MappingAtomicType.java:52:5
	return MappingAtomicType{
		Names: names,
		Types: types,
		Rest:  rest,
	}
}

func (this *MappingAtomicType) String() string {
	var builder strings.Builder
	builder.WriteString("(mapping")
	for i, name := range this.Names {
		builder.WriteString(fmt.Sprintf(" (%s %s)", name, this.Types[i].String()))
	}
	builder.WriteString(" ")
	builder.WriteString(this.Rest.String())
	builder.WriteString(")")
	return builder.String()
}

func (this *MappingAtomicType) AtomKind() Kind {
	// migrated from MappingAtomicType.java:74:5
	return Kind_MAPPING_ATOM
}

func (this *MappingAtomicType) FieldInnerVal(name string) SemType {
	for i, n := range this.Names {
		if n == name {
			return CellInnerVal(this.Types[i])
		}
	}
	return CellInnerVal(this.Rest)
}
