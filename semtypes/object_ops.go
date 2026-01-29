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

type ObjectOps struct {
}

var _ BasicTypeOps = &ObjectOps{}

func (this *ObjectOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from ObjectOps.java:51:5
	return bddSubtypeDiff(MAPPING_SUBTYPE_OBJECT, t2)
}

func (this *ObjectOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from ObjectOps.java:51:5
	return bddSubtypeIntersect(MAPPING_SUBTYPE_OBJECT, t2)
}

func (this *ObjectOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from ObjectOps.java:51:5
	return bddSubtypeUnion(MAPPING_SUBTYPE_OBJECT, t2)
}

func objectSubTypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from ObjectOps.java:43:5
	return memoSubtypeIsEmpty(cx, cx.mappingMemo(), objectBddIsEmpty, Bdd(t))
}

func objectBddIsEmpty(cx Context, b Bdd) bool {
	// migrated from ObjectOps.java:47:5
	return bddEveryPositive(cx, b, nil, nil, mappingFormulaIsEmpty)
}

func NewObjectOps() ObjectOps {
	this := ObjectOps{}
	return this
}

func (this *ObjectOps) Complement(t SubtypeData) SubtypeData {
	// migrated from ObjectOps.java:33:5
	return this.objectSubTypeComplement(t)
}

func (this *ObjectOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from ObjectOps.java:38:5
	return objectSubTypeIsEmpty(cx, t)
}

func (this *ObjectOps) objectSubTypeComplement(t SubtypeData) SubtypeData {
	// migrated from ObjectOps.java:51:5
	return bddSubtypeDiff(MAPPING_SUBTYPE_OBJECT, t)
}
