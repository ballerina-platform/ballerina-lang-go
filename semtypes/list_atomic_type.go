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

type ListAtomicType struct {
	Members fixedLengthArray
	Rest    *ComplexSemType
}

var _ AtomicType = &ListAtomicType{}

func (this *ListAtomicType) equals(other AtomicType) bool {
	if other, ok := other.(*ListAtomicType); ok {
		if !this.Rest.equals(other.Rest) {
			return false
		}
		return other.Members.FixedLength == this.Members.FixedLength &&
			slices.EqualFunc(other.Members.Initial, this.Members.Initial, func(a, b ComplexSemType) bool { return a.equals(&b) })
	}
	return false
}

func newListAtomicTypeFromMembersRest(members fixedLengthArray, rest *ComplexSemType) ListAtomicType {
	this := ListAtomicType{}
	this.Members = members
	this.Rest = rest
	return this
}

func listAtomicTypeFrom(members fixedLengthArray, rest *ComplexSemType) ListAtomicType {

	return newListAtomicTypeFromMembersRest(members, rest)
}

func (this *ListAtomicType) AtomKind() Kind {
	return Kind_LIST_ATOM
}

func (atomic *ListAtomicType) MemberAtInnerVal(index int) SemType {
	return CellInnerVal(atomic.MemberAt(index))
}

func (atomic *ListAtomicType) MemberAt(index int) *ComplexSemType {
	return listMemberAt(atomic.Members, atomic.Rest, index)
}
