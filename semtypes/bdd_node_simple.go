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

type bddNodeSimple struct {
	atom      Atom
	canonical string
}

var _ BddNode = &bddNodeSimple{}

func (this *bddNodeSimple) Left() Bdd {
	return bddAll()
}

func (this *bddNodeSimple) Middle() Bdd {
	return bddNothing()
}

func (this *bddNodeSimple) Right() Bdd {
	return bddNothing()
}

func (this *bddNodeSimple) Atom() Atom {
	return this.atom
}

func newBddNodeSimple(atom Atom) *bddNodeSimple {
	return &bddNodeSimple{
		atom:      atom,
		canonical: fmt.Sprintf("(%s (true) (false) (false))", atom.canonicalKey()),
	}
}

func (this *bddNodeSimple) canonicalKey() string {
	return this.canonical
}
