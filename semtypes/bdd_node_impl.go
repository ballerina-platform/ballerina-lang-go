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
	"sync"
)

type BddNodeImpl struct {
	atom             Atom
	left             Bdd
	middle           Bdd
	right            Bdd
	canonicalKeyFunc func() string
}

var _ BddNode = &BddNodeImpl{}

func (b *BddNodeImpl) Atom() Atom {
	return b.atom
}

func (b *BddNodeImpl) Left() Bdd {
	return b.left
}

func (b *BddNodeImpl) Middle() Bdd {
	return b.middle
}

func (b *BddNodeImpl) Right() Bdd {
	return b.right
}

func newBddNodeImpl(atom Atom, left, middle, right Bdd) *BddNodeImpl {
	node := &BddNodeImpl{atom: atom, left: left, middle: middle, right: right}
	node.canonicalKeyFunc = sync.OnceValue(func() string {
		return fmt.Sprintf("(%s (%s) (%s) (%s))", atom.canonicalKey(), left.canonicalKey(), middle.canonicalKey(), right.canonicalKey())
	})
	return node
}

func (b *BddNodeImpl) canonicalKey() string {
	return b.canonicalKeyFunc()
}
