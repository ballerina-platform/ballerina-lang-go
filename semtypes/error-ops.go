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

type ErrorOps struct {
	CommonOps
}

var _ BasicTypeOps = &ErrorOps{}

func errorSubtypeComplement(t SubtypeData) SubtypeData {
	// migrated from ErrorOps.java:39:5
	return bddSubtypeDiff(BDD_SUBTYPE_RO, t)
}

func errorSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from ErrorOps.java:43:5
	b := t.(Bdd)
	if bddPosMaybeEmpty(b) {
		b = BddIntersect(b, BDD_SUBTYPE_RO)
	}
	return memoSubtypeIsEmpty(cx, cx.mappingMemo(), errorBddIsEmpty, b)
}

func errorBddIsEmpty(cx Context, b Bdd) bool {
	// migrated from ErrorOps.java:52:5
	return bddEveryPositive(cx, b, nil, nil, mappingFormulaIsEmpty)
}

func (this *ErrorOps) Complement(d SubtypeData) SubtypeData {
	// migrated from ErrorOps.java:56:5
	return errorSubtypeComplement(d)
}

func (this *ErrorOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from ErrorOps.java:61:5
	return errorSubtypeIsEmpty(cx, t)
}
