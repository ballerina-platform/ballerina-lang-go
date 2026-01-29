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

type BddCommonOpsData interface {
}

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

type BddCommonOpsBase struct {
}

type BddCommonOpsMethods struct {
	Self BddCommonOps
}

func BddAtom(atom Atom) BddNode {
	// migrated from BddCommonOps.java:36:5
	return BddNodeCreate(atom, BddAll(), BddNothing(), BddNothing())
}

func BddUnion(b1 Bdd, b2 Bdd) Bdd {
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

	if allOrNothing1, ok := b1.(BddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return BddAll()
		}
		return b2
	}

	if allOrNothing2, ok := b2.(BddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return BddAll()
		}
		return b1
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := AtomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return BddCreate(b1Bdd.Atom(), b1Bdd.Left(), bddUnionWithMemo(memo, b1Bdd.Middle(), b2), b1Bdd.Right())
	} else if cmp > 0 {
		return BddCreate(b2Bdd.Atom(), b2Bdd.Left(), bddUnionWithMemo(memo, b1, b2Bdd.Middle()), b2Bdd.Right())
	} else {
		return BddCreate(b1Bdd.Atom(), bddUnionWithMemo(memo, b1Bdd.Left(), b2Bdd.Left()), bddUnionWithMemo(memo, b1Bdd.Middle(), b2Bdd.Middle()), bddUnionWithMemo(memo, b1Bdd.Right(), b2Bdd.Right()))
	}
}

func BddIntersect(b1 Bdd, b2 Bdd) Bdd {
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

	if allOrNothing1, ok := b1.(BddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return b2
		}
		return BddNothing()
	}

	if allOrNothing2, ok := b2.(BddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return b1
		}
		return BddNothing()
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := AtomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return BddCreate(b1Bdd.Atom(), bddIntersectWithMemo(memo, b1Bdd.Left(), b2), bddIntersectWithMemo(memo, b1Bdd.Middle(), b2), bddIntersectWithMemo(memo, b1Bdd.Right(), b2))
	} else if cmp > 0 {
		return BddCreate(b2Bdd.Atom(), bddIntersectWithMemo(memo, b1, b2Bdd.Left()), bddIntersectWithMemo(memo, b1, b2Bdd.Middle()), bddIntersectWithMemo(memo, b1, b2Bdd.Right()))
	} else {
		return BddCreate(b1Bdd.Atom(), bddIntersectWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), BddNothing(), bddIntersectWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	}
}

func BddDiff(b1 Bdd, b2 Bdd) Bdd {
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
		return BddNothing()
	}

	if allOrNothing2, ok := b2.(BddAllOrNothing); ok {
		if allOrNothing2.IsAll() {
			return BddNothing()
		}
		return b1
	}

	if allOrNothing1, ok := b1.(BddAllOrNothing); ok {
		if allOrNothing1.IsAll() {
			return bddComplement(b2)
		}
		return BddNothing()
	}

	b1Bdd := b1.(BddNode)
	b2Bdd := b2.(BddNode)
	cmp := AtomCmp(b1Bdd.Atom(), b2Bdd.Atom())
	if cmp < 0 {
		return BddCreate(b1Bdd.Atom(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), b2), BddNothing(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), b2))
	} else if cmp > 0 {
		return BddCreate(b2Bdd.Atom(), bddDiffWithMemo(memo, b1, bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), BddNothing(), bddDiffWithMemo(memo, b1, bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	} else {
		return BddCreate(b1Bdd.Atom(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Left(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Left(), b2Bdd.Middle())), BddNothing(), bddDiffWithMemo(memo, bddUnionWithMemo(memo, b1Bdd.Right(), b1Bdd.Middle()), bddUnionWithMemo(memo, b2Bdd.Right(), b2Bdd.Middle())))
	}
}

func BddComplement(b Bdd) Bdd {
	// migrated from BddCommonOps.java:190:5
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		return allOrNothing.Complement()
	}
	return BddNodeComplement(b.(BddNode))
}

func bddComplement(b Bdd) Bdd {
	return BddComplement(b)
}

func BddNodeComplement(b BddNode) Bdd {
	// migrated from BddCommonOps.java:198:5
	bddNothing := BddNothing()
	if b.Right() == bddNothing {
		return BddCreate(b.Atom(), bddNothing, BddComplement(BddUnion(b.Left(), b.Middle())), BddComplement(b.Middle()))
	} else if b.Left() == bddNothing {
		return BddCreate(b.Atom(), BddComplement(b.Middle()), BddComplement(BddUnion(b.Right(), b.Middle())), bddNothing)
	} else if b.Middle() == bddNothing {
		return BddCreate(b.Atom(), BddComplement(b.Left()), BddComplement(BddUnion(b.Left(), b.Right())), BddComplement(b.Right()))
	} else {
		return BddCreate(b.Atom(), BddComplement(BddUnion(b.Left(), b.Middle())), bddNothing, BddComplement(BddUnion(b.Right(), b.Middle())))
	}
}

func BddCreate(atom Atom, left Bdd, middle Bdd, right Bdd) Bdd {
	// migrated from BddCommonOps.java:226:5
	if allOrNothing, ok := middle.(BddAllOrNothing); ok && allOrNothing.IsAll() {
		return middle
	}
	if left == right {
		return BddUnion(left, right)
	}
	return BddNodeCreate(atom, left, middle, right)
}

func AtomCmp(a1 Atom, a2 Atom) int {
	// migrated from BddCommonOps.java:238:5
	r1, ok1 := a1.(*RecAtom)
	r2, ok2 := a2.(*RecAtom)

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

func (this *BddCommonOpsMethods) BddToString(b Bdd, inner bool) string {
	// migrated from BddCommonOps.java:254:5
	if allOrNothing, ok := b.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return "1"
		}
		return "0"
	}

	var str string
	bdd := b.(BddNode)
	a := bdd.Atom()
	if recAtom, ok := a.(*RecAtom); ok {
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
