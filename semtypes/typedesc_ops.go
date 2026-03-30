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

type typedescOps struct {
	CommonOps
}

var _ BasicTypeOps = &typedescOps{}

func typedescSubtypeComplement(t SubtypeData) SubtypeData {
	// migrated from typedescOps.java:38:5
	return bddComplement(t.(Bdd))
}

func typedescSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from typedescOps.java:42:5
	b := t.(Bdd)
	if bddPosMaybeEmpty(b) {
		b = bddIntersect(b, BDD_SUBTYPE_RO)
	}
	return mappingSubtypeIsEmpty(cx, b)
}

func newTypedescOps() typedescOps {
	this := typedescOps{}
	return this
}

func (this *typedescOps) complement(d SubtypeData) SubtypeData {
	// migrated from typedescOps.java:51:5
	return typedescSubtypeComplement(d)
}

func (this *typedescOps) IsEmpty(cx Context, d SubtypeData) bool {
	// migrated from typedescOps.java:56:5
	return typedescSubtypeIsEmpty(cx, d)
}
