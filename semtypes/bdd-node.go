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

type BddNode interface {
	Bdd
	Atom() Atom
	Left() Bdd
	Middle() Bdd
	Right() Bdd
}

func BddNodeCreate(atom Atom, left Bdd, middle Bdd, right Bdd) BddNode {
	// migrated from BddNode.java:31:5
	if IsSimpleNode(left, middle, right) {
		return &BddNodeSimple{atom: atom}
	}
	return &BddNodeImpl{atom: atom, left: left, middle: middle, right: right}
}

func IsSimpleNode(left Bdd, middle Bdd, right Bdd) bool {
	var leftIsAll, middleIsNothing, rightIsNothing bool
	if leftNode, ok := left.(BddAllOrNothing); ok {
		leftIsAll = leftNode.IsAll()
	}
	if middleNode, ok := middle.(BddAllOrNothing); ok {
		middleIsNothing = middleNode.IsNothing()
	}
	if rightNode, ok := right.(BddAllOrNothing); ok {
		rightIsNothing = rightNode.IsNothing()
	}
	return leftIsAll && middleIsNothing && rightIsNothing
}
