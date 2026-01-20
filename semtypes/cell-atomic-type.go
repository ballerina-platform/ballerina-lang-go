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

type CellMutability uint

type CellAtomicType struct {
	Ty  SemType
	Mut CellMutability
}

const (
	CellMutability_CELL_MUT_NONE CellMutability = iota
	CellMutability_CELL_MUT_LIMITED
	CellMutability_CELL_MUT_UNLIMITED
)

func (this *CellAtomicType) equals(other AtomicType) bool {
	if other, ok := other.(*CellAtomicType); ok {
		return other.Ty == this.Ty && other.Mut == this.Mut
	}
	return false
}

var _ AtomicType = &CellAtomicType{}

func NewCellAtomicTypeFromTyMut(ty SemType, mut CellMutability) CellAtomicType {
	this := CellAtomicType{}
	common.Assert(ty != nil)
	this.Ty = ty
	this.Mut = mut
	return this
}

func CellAtomicTypeFrom(ty SemType, mut CellMutability) CellAtomicType {
	// migrated from CellAtomicType.java:33:5
	common.Assert(ty != nil)
	return NewCellAtomicTypeFromTyMut(ty, mut)
}

func (this *CellAtomicType) AtomKind() Kind {
	// migrated from CellAtomicType.java:39:5
	return Kind_CELL_ATOM
}
