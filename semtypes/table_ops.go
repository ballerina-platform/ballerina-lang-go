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

type tableOps struct {
}

var _ BasicTypeOps = &tableOps{}

func tableSubtypeComplement(t SubtypeData) SubtypeData {
	return bddSubtypeDiff(LIST_SUBTYPE_THREE_ELEMENT, t)
}

func tableSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	b := t.(Bdd)
	if bddPosMaybeEmpty(b) {
		b = bddIntersect(b, LIST_SUBTYPE_THREE_ELEMENT)
	}
	return listSubtypeIsEmpty(cx, b)
}

func (t *tableOps) complement(d SubtypeData) SubtypeData {
	return tableSubtypeComplement(d)
}

func (t *tableOps) IsEmpty(cx Context, d SubtypeData) bool {
	return tableSubtypeIsEmpty(cx, d)
}

func (t *tableOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeUnion(d1, d2)
}

func (t *tableOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeIntersect(d1, d2)
}

func (t *tableOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	return bddSubtypeDiff(d1, d2)
}
