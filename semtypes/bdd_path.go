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

type bddPath struct {
	bdd Bdd
	pos []atom
	neg []atom
}

func newBddPathFromBddPath(src bddPath) bddPath {
	this := bddPath{}
	this.bdd = src.bdd
	this.pos = slices.Clone(src.pos)
	this.neg = slices.Clone(src.neg)
	return this
}

func newBddPath() bddPath {
	this := bddPath{}
	this.bdd = bddAll()
	this.pos = nil
	this.neg = nil
	return this
}

func bddPaths(b Bdd, paths *[]bddPath, accum bddPath) {
	allOrNothing, ok := b.(*bddAllOrNothing)
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
		left.bdd = bddIntersect(left.bdd, bddAtom(bn.Atom()))
		bddPaths(bn.Left(), paths, left)
		bddPaths(bn.Middle(), paths, accum)
		right.neg = append(right.neg, bn.Atom())
		right.bdd = bddDiff(right.bdd, bddAtom(bn.Atom()))
		bddPaths(bn.Right(), paths, right)
	}
}

func bddPathClone(path bddPath) bddPath {

	return newBddPathFromBddPath(path)
}

func bddPathFrom() bddPath {
	return newBddPath()
}
