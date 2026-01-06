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

type BooleanOps struct {
}

var _ BasicTypeOps = &BooleanOps{}


func (this *BooleanOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from BooleanOps.java:33:5
	v1 := d1.(BooleanSubtype)
	v2 := d2.(BooleanSubtype)
	if v1.value == v2.value {
		return v1
	} else {
		return CreateAll()
	}
}

func (this *BooleanOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from BooleanOps.java:40:5
	v1 := d1.(BooleanSubtype)
	v2 := d2.(BooleanSubtype)
	if v1.value == v2.value {
		return v1
	} else {
		return CreateNothing()
	}
}

func (this *BooleanOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from BooleanOps.java:47:5
	v1 := d1.(BooleanSubtype)
	v2 := d2.(BooleanSubtype)
	if v1.value == v2.value {
		return CreateNothing()
	} else {
		return v1
	}
}

func (this *BooleanOps) Complement(d SubtypeData) SubtypeData {
	// migrated from BooleanOps.java:54:5
	v := d.(BooleanSubtype)
	t := BooleanSubtypeFrom(!v.value)
	return t
}

func (this *BooleanOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from BooleanOps.java:61:5
	return notIsEmpty(cx, t)
}
