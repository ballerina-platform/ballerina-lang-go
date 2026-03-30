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

type booleanOps struct{}

var _ BasicTypeOps = &booleanOps{}

func (this *booleanOps) Union(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from booleanOps.java:33:5
	v1 := d1.(booleanSubtype)
	v2 := d2.(booleanSubtype)
	if v1.Value == v2.Value {
		return v1
	} else {
		return createAll()
	}
}

func (this *booleanOps) Intersect(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from booleanOps.java:40:5
	v1 := d1.(booleanSubtype)
	v2 := d2.(booleanSubtype)
	if v1.Value == v2.Value {
		return v1
	} else {
		return createNothing()
	}
}

func (this *booleanOps) Diff(d1 SubtypeData, d2 SubtypeData) SubtypeData {
	// migrated from booleanOps.java:47:5
	v1 := d1.(booleanSubtype)
	v2 := d2.(booleanSubtype)
	if v1.Value == v2.Value {
		return createNothing()
	} else {
		return v1
	}
}

func (this *booleanOps) complement(d SubtypeData) SubtypeData {
	// migrated from booleanOps.java:54:5
	v := d.(booleanSubtype)
	t := booleanSubtypeFrom(!v.Value)
	return t
}

func (this *booleanOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from booleanOps.java:61:5
	return notIsEmpty(cx, t)
}
