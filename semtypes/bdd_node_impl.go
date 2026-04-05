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

type bddNodeImpl struct {
	atom      atom
	left      Bdd
	middle    Bdd
	right     Bdd
	canonical string
}

var _ BddNode = &bddNodeImpl{}

func (this *bddNodeImpl) Atom() atom {
	return this.atom
}

func (this *bddNodeImpl) Left() Bdd {
	return this.left
}

func (this *bddNodeImpl) Middle() Bdd {
	return this.middle
}

func (this *bddNodeImpl) Right() Bdd {
	return this.right
}

func newBddNodeImpl(atom atom, left, middle, right Bdd) *bddNodeImpl {
	return &bddNodeImpl{
		atom:      atom,
		left:      left,
		middle:    middle,
		right:     right,
		canonical: fmt.Sprintf("(%s (%s) (%s) (%s))", atom.canonicalKey(), left.canonicalKey(), middle.canonicalKey(), right.canonicalKey()),
	}
}

func (this *bddNodeImpl) canonicalKey() string {
	return this.canonical
}
