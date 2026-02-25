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

type FunctionOps struct {
}

var _ BasicTypeOps = &FunctionOps{}

func (this *FunctionOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from FunctionOps.java:45:5
	return memoSubtypeIsEmpty(cx, cx.functionMemo(), func(cx Context, b Bdd) bool {
		return bddEvery(cx, b, nil, nil, functionFormulaIsEmpty)
	}, t.(Bdd))
}

func (this *FunctionOps) Complement(t SubtypeData) SubtypeData {
	// migrated from FunctionOps.java:49:5
	return BddComplement(t.(Bdd))
}

func (this *FunctionOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FunctionOps.java:51:5
	return BddDiff(t1.(Bdd), t2.(Bdd))
}

func (this *FunctionOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FunctionOps.java:53:5
	return BddIntersect(t1.(Bdd), t2.(Bdd))
}

func (this *FunctionOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FunctionOps.java:55:5
	return BddUnion(t1.(Bdd), t2.(Bdd))
}

func functionFormulaIsEmpty(cx Context, pos *Conjunction, neg *Conjunction) bool {
	// migrated from FunctionOps.java:51:5
	return functionPathIsEmpty(cx, functionIntersectRet(cx, pos), functionUnionParams(cx, pos), functionUnionQualifiers(cx, pos), pos, neg)
}

func functionPathIsEmpty(cx Context, rets SemType, params SemType, qualifiers SemType, pos *Conjunction, neg *Conjunction) bool {
	// migrated from FunctionOps.java:56:5
	if neg == nil {
		return false
	} else {
		t := cx.functionAtomType(neg.Atom)
		t0 := t.ParamType
		t1 := t.RetType
		t2 := t.Qualifiers
		if t.IsGeneric {
			return (((IsSubtype(cx, qualifiers, t2) && IsSubtype(cx, params, t0)) && IsSubtype(cx, rets, t1)) || functionPathIsEmpty(cx, rets, params, qualifiers, pos, neg.Next))
		}
		return (((IsSubtype(cx, qualifiers, t2) && IsSubtype(cx, t0, params)) && functionPhi(cx, t0, Complement(t1), pos)) || functionPathIsEmpty(cx, rets, params, qualifiers, pos, neg.Next))
	}
}

func functionPhi(cx Context, t0 SemType, t1 SemType, pos *Conjunction) bool {
	// migrated from FunctionOps.java:81:5
	if pos == nil {
		return ((!IsNever(t0)) && (IsEmpty(cx, t0) || IsEmpty(cx, t1)))
	}
	return functionPhiInner(cx, t0, t1, pos)
}

func functionPhiInner(cx Context, t0 SemType, t1 SemType, pos *Conjunction) bool {
	// migrated from FunctionOps.java:89:5
	if pos == nil {
		return (IsEmpty(cx, t0) || IsEmpty(cx, t1))
	} else {
		s := cx.functionAtomType(pos.Atom)
		s0 := s.ParamType
		s1 := s.RetType
		return (((IsSubtype(cx, t0, s0) || IsSubtype(cx, functionIntersectRet(cx, pos.Next), Complement(t1))) && functionPhiInner(cx, t0, Intersect(t1, s1), pos.Next)) && functionPhiInner(cx, Diff(t0, s0), t1, pos.Next))
	}
}

func functionUnionParams(cx Context, pos *Conjunction) SemType {
	// migrated from FunctionOps.java:104:5
	if pos == nil {
		return &NEVER
	}
	return Union(cx.functionAtomType(pos.Atom).ParamType, functionUnionParams(cx, pos.Next))
}

func functionUnionQualifiers(cx Context, pos *Conjunction) SemType {
	// migrated from FunctionOps.java:111:5
	if pos == nil {
		return &NEVER
	}
	return Union(cx.functionAtomType(pos.Atom).Qualifiers, functionUnionQualifiers(cx, pos.Next))
}

func functionIntersectRet(cx Context, pos *Conjunction) SemType {
	// migrated from FunctionOps.java:119:5
	if pos == nil {
		return &VAL
	}
	return Intersect(cx.functionAtomType(pos.Atom).RetType, functionIntersectRet(cx, pos.Next))
}

func NewFunctionOps() FunctionOps {
	this := FunctionOps{}
	return this
}

func (this *FunctionOps) functionTheta(cx Context, t0 SemType, t1 SemType, pos *Conjunction) bool {
	// migrated from FunctionOps.java:126:5
	if pos == nil {
		return (IsEmpty(cx, t0) || IsEmpty(cx, t1))
	} else {
		s := cx.functionAtomType(pos.Atom)
		s0 := s.ParamType
		s1 := s.RetType
		return ((IsSubtype(cx, t0, s0) || this.functionTheta(cx, Diff(s0, t0), s1, pos.Next)) && (IsSubtype(cx, t1, Complement(s1)) || this.functionTheta(cx, s0, Intersect(s1, t1), pos.Next)))
	}
}

// Corresponds to dom^? in AMK tutorial.
func FunctionParamListType(cx Context, fnTy SemType) SemType {
	if !IsSubtypeSimple(fnTy, FUNCTION) {
		return nil
	}
	switch ty := fnTy.(type) {
	case *BasicTypeBitSet:
		return &NEVER
	case ComplexSemType:
		bdd := getComplexSubtypeData(ty, BT_FUNCTION).(Bdd)
		return functionParamListTypeInner(cx, &NEVER, bdd)
	default:
		panic("impossible")
	}
}

func functionParamListTypeInner(cx Context, accumTy SemType, bdd Bdd) SemType {
	if allOrNothing, ok := bdd.(BddAllOrNothing); ok {
		if allOrNothing.IsAll() {
			return accumTy
		}
		return &ANY
	}
	bddNode := bdd.(BddNode)
	atomArgListTy := cx.functionAtomType(bddNode.Atom()).ParamType
	return Intersect(functionParamListTypeInner(cx, Union(accumTy, atomArgListTy), bddNode.Left()),
		Intersect(functionParamListTypeInner(cx, accumTy, bddNode.Middle()),
			functionParamListTypeInner(cx, accumTy, bddNode.Right())))
}

// Corresponds to apply^? in AMK tutorial.
func FunctionReturnType(cx Context, fnTy SemType, argList SemType) SemType {
	domain := FunctionParamListType(cx, fnTy)
	if domain == nil || !IsSubtype(cx, argList, domain) {
		return nil
	}
	switch ty := fnTy.(type) {
	case *BasicTypeBitSet:
		return &VAL
	case ComplexSemType:
		bdd := getComplexSubtypeData(ty, BT_FUNCTION).(Bdd)
		return functionReturnTypeInner(cx, argList, &VAL, bdd)
	default:
		panic("impossible")
	}
}

func functionReturnTypeInner(cx Context, accumArgList SemType, accumReturn SemType, bdd Bdd) SemType {
	if IsEmpty(cx, accumArgList) || IsEmpty(cx, accumReturn) {
		return &NEVER
	}
	switch b := bdd.(type) {
	case BddAllOrNothing:
		if b.IsAll() {
			return accumReturn
		}
		return &NEVER
	case BddNode:
		fnAtom := cx.functionAtomType(b.Atom())
		atomArgListTy := fnAtom.ParamType
		atomReturnTy := fnAtom.RetType
		return Union(functionReturnTypeInner(cx, accumArgList, Intersect(accumReturn, atomReturnTy), b.Left()),
			Union(functionReturnTypeInner(cx, Diff(accumArgList, atomArgListTy), accumReturn, b.Left()),
				Union(functionReturnTypeInner(cx, accumArgList, accumReturn, b.Middle()),
					functionReturnTypeInner(cx, accumArgList, accumReturn, b.Right()))))

	default:
		panic("impossible")
	}
}
