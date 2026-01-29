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

type TypedescOps struct {
	CommonOps
}

var _ BasicTypeOps = &TypedescOps{}

func typedescSubtypeComplement(t SubtypeData) SubtypeData {
	// migrated from TypedescOps.java:38:5
	return bddComplement(Bdd(t))
}

func typedescSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from TypedescOps.java:42:5
	b := Bdd(t)
	if bddPosMaybeEmpty(b) {
		b = BddIntersect(b, BDD_SUBTYPE_RO)
	}
	return mappingSubtypeIsEmpty(cx, b)
}

func NewTypedescOps() TypedescOps {
	this := TypedescOps{}
	return this
}

func (this *TypedescOps) Complement(d SubtypeData) SubtypeData {
	// migrated from TypedescOps.java:51:5
	return typedescSubtypeComplement(d)
}

func (this *TypedescOps) IsEmpty(cx Context, d SubtypeData) bool {
	// migrated from TypedescOps.java:56:5
	return typedescSubtypeIsEmpty(cx, d)
}
