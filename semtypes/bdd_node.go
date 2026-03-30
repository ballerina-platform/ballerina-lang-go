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

func bddNodeCreate(atom Atom, left Bdd, middle Bdd, right Bdd) BddNode {
	// migrated from BddNode.java:31:5
	if isSimpleNode(left, middle, right) {
		return newBddNodeSimple(atom)
	}
	return newBddNodeImpl(atom, left, middle, right)
}

func isSimpleNode(left Bdd, middle Bdd, right Bdd) bool {
	leftIsAll := isAll(left)
	middleIsNothing := isNothing(middle)
	rightIsNothing := isNothing(right)
	return leftIsAll && middleIsNothing && rightIsNothing
}

func isAll(bdd Bdd) bool {
	if allOrNothig, ok := bdd.(*bddAllOrNothing); ok {
		return allOrNothig.isAll
	}
	return false
}

func isNothing(bdd Bdd) bool {
	if allOrNothig, ok := bdd.(*bddAllOrNothing); ok {
		return !allOrNothig.isAll
	}
	return false
}
