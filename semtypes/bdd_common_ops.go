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

type BddCommonOpsData any

type BddCommonOps interface {
	BddCommonOpsData
}

type bddOpMemoKey struct {
	B1 Bdd
	B2 Bdd
}

type bddOpMemo struct {
	UnionMemo        map[bddOpMemoKey]Bdd
	IntersectionMemo map[bddOpMemoKey]Bdd
	DiffMemo         map[bddOpMemoKey]Bdd
}

type bddCommonOpsBase struct{}

type bddCommonOpsMethods struct {
	Self BddCommonOps
}

func bddAtom(atom Atom) BddNode {
	// migrated from BddCommonOps.java:36:5
	return bddNodeCreate(atom, bddAll(), bddNothing(), bddNothing())
}

func bddUnion(b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:43:5
	return bddUnionWithMemo(createBddOpMemo(), b1, b2)
}

func bddUnionWithMemo(memoTable *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:47:5
	key := bddOpMemoKey{B1: b1, B2: b2}
	memoized, ok := memoTable.UnionMemo[key]
	if ok {
		return memoized
	}
	memoized = bddUnionInner(memoTable, b1, b2)
	memoTable.UnionMemo[key] = memoized
	return memoized
}

func bddUnionInner(memo *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:58:5
	if b1 == b2 {
		return b1
	}

	if allOrNothing1, ok := b1.(*bddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return bddAll()
		}
		return b2
	}

	if allOrNothing2, ok := b2.(*bddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return bddAll()
		}
		return b1
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := atomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return bddCreate(b1Bdd.Atom(), b1Bdd.Left(), bddUnionWithMemo(memo, b1Bdd.Middle(), b2), b1Bdd.Right())
	} else if cmp > 0 {
		return bddCreate(b2Bdd.Atom(), b2Bdd.Left(), bddUnionWithMemo(memo, b1, b2Bdd.Middle()), b2Bdd.Right())
	} else {
		return bddCreate(b1Bdd.Atom(), bddUnionWithMemo(memo, b1Bdd.Left(), b2Bdd.Left()), bddUnionWithMemo(memo, b1Bdd.Middle(), b2Bdd.Middle()), bddUnionWithMemo(memo, b1Bdd.Right(), b2Bdd.Right()))
	}
}

func bddIntersect(b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:88:5
	return bddIntersectWithMemo(createBddOpMemo(), b1, b2)
}

func bddIntersectWithMemo(memo *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:92:5
	key := bddOpMemoKey{B1: b1, B2: b2}
	memoized, ok := memo.IntersectionMemo[key]
	if ok {
		return memoized
	}
	memoized = bddIntersectInner(memo, b1, b2)
	memo.IntersectionMemo[key] = memoized
	return memoized
}

func bddIntersectInner(memo *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:103:5
	if b1 == b2 {
		return b1
	}

	if allOrNothing1, ok := b1.(*bddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return b2
		}
		return bddNothing()
	}

	if allOrNothing2, ok := b2.(*bddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return b1
		}
		return bddNothing()
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := atomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return bddCreate(b1Bdd.Atom(), bddIntersectWithMemo(memo, b1Bdd.Left(), b2), bddIntersectWithMemo(memo, b1Bdd.Middle(), b2), bddIntersectWithMemo(memo, b1Bdd.Right(), b2))
	} else if cmp > 0 {
		return bddCreate(b2Bdd.Atom(), bddIntersectWithMemo(memo, b1, b2Bdd.Left()), bddIntersectWithMemo(memo, b1, b2Bdd.Middle()), bddIntersectWithMemo(memo, b1, b2Bdd.Right()))
	} else {
		return bddCreate(b1Bdd.Atom(), bddIntersectWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), bddNothing(), bddIntersectWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	}
}

func bddDiff(b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:137:5
	return bddDiffWithMemo(createBddOpMemo(), b1, b2)
}

func bddDiffWithMemo(memo *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:141:5
	key := bddOpMemoKey{B1: b1, B2: b2}
	memoized, ok := memo.DiffMemo[key]
	if ok {
		return memoized
	}
	memoized = bddDiffInner(memo, b1, b2)
	memo.DiffMemo[key] = memoized
	return memoized
}

func bddDiffInner(memo *bddOpMemo, b1 Bdd, b2 Bdd) Bdd {
	// migrated from BddCommonOps.java:152:5
	if b1 == b2 {
		return bddNothing()
	}

	if allOrNothing2, ok := b2.(*bddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return bddNothing()
		}
		return b1
	}

	if allOrNothing1, ok := b1.(*bddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return bddComplement(b2)
		}
		return bddNothing()
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := atomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return bddCreate(b1Bdd.Atom(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), b2), bddNothing(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), b2))
	} else if cmp > 0 {
		return bddCreate(b2Bdd.Atom(), bddDiffWithMemo(memo, b1, bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), bddNothing(), bddDiffWithMemo(memo, b1, bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	} else {
		return bddCreate(b1Bdd.Atom(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), bddNothing(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	}
}

func bddComplement(b Bdd) Bdd {
	// migrated from BddCommonOps.java:190:5
	if allOrNothing, ok := b.(*bddAllOrNothing); ok {
		return allOrNothing.complement()
	}
	return bddNodeComplement(b.(BddNode))
}

func bddNodeComplement(b BddNode) Bdd {
	// migrated from BddCommonOps.java:198:5
	bddNothing := bddNothing()
	if b.Right() == bddNothing {
		return bddCreate(b.Atom(), bddNothing, bddComplement(bddUnion(b.Left(), b.Middle())), bddComplement(b.Middle()))
	} else if b.Left() == bddNothing {
		return bddCreate(b.Atom(), bddComplement(b.Middle()), bddComplement(bddUnion(b.Right(), b.Middle())), bddNothing)
	} else if b.Middle() == bddNothing {
		return bddCreate(b.Atom(), bddComplement(b.Left()), bddComplement(bddUnion(b.Left(), b.Right())), bddComplement(b.Right()))
	} else {
		return bddCreate(b.Atom(), bddComplement(bddUnion(b.Left(), b.Middle())), bddNothing, bddComplement(bddUnion(b.Right(), b.Middle())))
	}
}

func bddCreate(atom Atom, left Bdd, middle Bdd, right Bdd) Bdd {
	// migrated from BddCommonOps.java:226:5
	if allOrNothing, ok := middle.(*bddAllOrNothing); ok && allOrNothing.IsAll() {
		return middle
	}
	if left == right {
		return bddUnion(left, right)
	}
	return bddNodeCreate(atom, left, middle, right)
}

func atomCmp(a1 Atom, a2 Atom) int {
	// migrated from BddCommonOps.java:238:5
	r1, ok1 := a1.(*recAtom)
	r2, ok2 := a2.(*recAtom)

	if ok1 {
		if ok2 {
			return r1.Index() - r2.Index()
		}
		return -1
	} else if ok2 {
		return 1
	}
	return a1.Index() - a2.Index()
}

func (this *bddCommonOpsMethods) BddToString(b Bdd, inner bool) string {
	// migrated from BddCommonOps.java:254:5
	if allOrNothing, ok := b.(*bddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return "1"
		}
		return "0"
	}

	var str string
	bdd := b.(BddNode)
	a := bdd.Atom()
	if recAtom, ok := a.(*recAtom); ok {
		str = "r" + string(rune(recAtom.Index()))
	} else {
		str = "a" + string(rune(a.Index()))
	}
	str = str + "?" + this.BddToString(bdd.Left(), true) + ":" + this.BddToString(bdd.Middle(), true) + ":" + this.BddToString(bdd.Right(), true)
	if inner {
		str = "(" + str + ")"
	}
	return str
}

func createBddOpMemo() *bddOpMemo {
	// migrated from BddCommonOps.java:283:9
	return &bddOpMemo{
		UnionMemo:        make(map[bddOpMemoKey]Bdd),
		IntersectionMemo: make(map[bddOpMemoKey]Bdd),
		DiffMemo:         make(map[bddOpMemoKey]Bdd),
	}
}
