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

type BddPath struct {
	bdd Bdd
	pos []Atom
	neg []Atom
}

func newBddPathFromBddPath(bddPath BddPath) BddPath {
	this := BddPath{}
	this.bdd = bddPath.bdd
	this.pos = slices.Clone(bddPath.pos)
	this.neg = slices.Clone(bddPath.neg)
	return this
}

func NewBddPath() BddPath {
	this := BddPath{}
	this.bdd = BddAll()
	this.pos = nil
	this.neg = nil
	return this
}

func BddPaths(b Bdd, paths *[]BddPath, accum BddPath) {
	// migrated from BddPath.java:50:5
	allOrNothing, ok := b.(BddAllOrNothing)
	if ok {
		if allOrNothing.IsAll() {
			*paths = append(*paths, accum)
		}
	} else {
		left := bddPathClone(accum)
		right := bddPathClone(accum)
		bn, ok := b.(BddNode)
		if !ok {
			panic("b is not a BddNode")
		}
		left.pos = append(left.pos, bn.Atom())
		left.bdd = BddIntersect(left.bdd, BddAtom(bn.Atom()))
		BddPaths(bn.Left(), paths, left)
		BddPaths(bn.Middle(), paths, accum)
		right.neg = append(right.neg, bn.Atom())
		right.bdd = BddDiff(right.bdd, BddAtom(bn.Atom()))
		BddPaths(bn.Right(), paths, right)
	}
}

func bddPathClone(path BddPath) BddPath {
	// migrated from BddPath.java:69:5

	return newBddPathFromBddPath(path)
}

func BddPathFrom() BddPath {
	// migrated from BddPath.java:73:5
	return NewBddPath()
}
