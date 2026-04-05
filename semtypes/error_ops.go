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

type errorOps struct{}

var _ BasicTypeOps = &errorOps{}

func errorSubtypeComplement(t SubtypeData) SubtypeData {
	return bddSubtypeDiff(BDD_SUBTYPE_RO, t)
}

func errorSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	b := t.(Bdd)
	if bddPosMaybeEmpty(b) {
		b = bddIntersect(b, BDD_SUBTYPE_RO)
	}
	return memoSubtypeIsEmpty(cx, cx.mappingMemo(), errorBddIsEmpty, b)
}

func errorBddIsEmpty(cx Context, b Bdd) bool {
	return bddEveryPositive(cx, b, conjunctionNil, conjunctionNil, mappingFormulaIsEmpty)
}

func (this *errorOps) complement(d SubtypeData) SubtypeData {
	return errorSubtypeComplement(d)
}

func (this *errorOps) IsEmpty(cx Context, t SubtypeData) bool {
	return errorSubtypeIsEmpty(cx, t)
}

func (this *errorOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeUnion(d1, d2)
}

func (this *errorOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeIntersect(d1, d2)
}

func (this *errorOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeDiff(d1, d2)
}
