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

type objectOps struct {
}

var _ BasicTypeOps = &objectOps{}

func (this *objectOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from objectOps.java:51:5
	return bddSubtypeDiff(t1, t2)
}

func (this *objectOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from objectOps.java:51:5
	return bddSubtypeIntersect(t1, t2)
}

func (this *objectOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from objectOps.java:51:5
	return bddSubtypeUnion(t1, t2)
}

func objectSubTypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from objectOps.java:43:5
	return memoSubtypeIsEmpty(cx, cx.mappingMemo(), objectBddIsEmpty, t.(Bdd))
}

func objectBddIsEmpty(cx Context, b Bdd) bool {
	// migrated from objectOps.java:47:5
	return bddEveryPositive(cx, b, conjunctionNil, conjunctionNil, mappingFormulaIsEmpty)
}

func newObjectOps() objectOps {
	this := objectOps{}
	return this
}

func (this *objectOps) complement(t SubtypeData) SubtypeData {
	// migrated from objectOps.java:33:5
	return this.objectSubTypeComplement(t)
}

func (this *objectOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from objectOps.java:38:5
	return objectSubTypeIsEmpty(cx, t)
}

func (this *objectOps) objectSubTypeComplement(t SubtypeData) SubtypeData {
	// migrated from objectOps.java:51:5
	return bddSubtypeDiff(MAPPING_SUBTYPE_OBJECT, t)
}
