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

import "slices"

type ListAtomicType struct {
	Members FixedLengthArray
	Rest    CellSemType
}

var _ AtomicType = &ListAtomicType{}

func (this *ListAtomicType) equals(other AtomicType) bool {
	if other, ok := other.(*ListAtomicType); ok {
		if other.Rest != this.Rest {
			return false
		}
		return other.Members.FixedLength == this.Members.FixedLength &&
			slices.Equal(other.Members.Initial, this.Members.Initial)
	}
	return false
}

func NewListAtomicTypeFromMembersRest(members FixedLengthArray, rest CellSemType) ListAtomicType {
	this := ListAtomicType{}
	this.Members = members
	this.Rest = rest
	return this
}

func ListAtomicTypeFrom(members FixedLengthArray, rest CellSemType) ListAtomicType {
	// migrated from ListAtomicType.java:34:5

	return NewListAtomicTypeFromMembersRest(members, rest)
}

func (this *ListAtomicType) AtomKind() Kind {
	// migrated from ListAtomicType.java:38:5
	return Kind_LIST_ATOM
}
