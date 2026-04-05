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

import "fmt"

type recAtom struct {
	index      int
	targetKind Kind
}

var ZERO = newRecAtomFromInt(BDD_REC_ATOM_READONLY)
var _ Atom = &recAtom{}

func newRecAtomFromInt(index int) recAtom {
	this := recAtom{}

	this.index = index
	return this
}

func newRecAtomFromIntKind(index int, targetKind Kind) recAtom {
	this := recAtom{}

	this.index = index
	this.targetKind = targetKind
	return this
}

func createRecAtom(index int) recAtom {
	if index == BDD_REC_ATOM_READONLY {
		return ZERO
	}
	return newRecAtomFromInt(index)
}

func createXMLRecAtom(index int) recAtom {
	return newRecAtomFromIntKind(index, Kind_XML_ATOM)
}

func createDistinctRecAtom(index int) recAtom {
	return newRecAtomFromIntKind(index, Kind_DISTINCT_ATOM)
}

func (this *recAtom) setKind(targetKind Kind) {
	this.targetKind = targetKind
}

func (this *recAtom) Index() int {
	return this.index
}

func (this *recAtom) Kind() Kind {
	// if this.targetKind == 0 {
	// 	panic("Target kind is not set for the recursive type atom")
	// }
	return this.targetKind
}

func (this *recAtom) canonicalKey() string {
	return fmt.Sprintf("r%d", this.index)
}

func (this *recAtom) String() string {
	return fmt.Sprintf("r%d", this.index)
}
