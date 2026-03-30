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

type floatOps struct {
	CommonOps
}

var _ BasicTypeOps = &floatOps{}

func newFloatOps() floatOps {
	this := floatOps{}
	return this
}

func (this *floatOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from floatOps.java:36:5
	var values []enumerableType[float64]
	var v1 enumerableSubtype[float64] = new(t1.(floatSubtype))
	var v2 enumerableSubtype[float64] = new(t2.(floatSubtype))
	allowed := enumerableSubtypeUnion(v1, v2, &values)
	return createFloatSubtype(allowed, values)
}

func (this *floatOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from floatOps.java:44:5
	var values []enumerableType[float64]
	var v1 enumerableSubtype[float64] = new(t1.(floatSubtype))
	var v2 enumerableSubtype[float64] = new(t2.(floatSubtype))
	allowed := enumerableSubtypeIntersect(v1, v2, &values)
	return createFloatSubtype(allowed, values)
}

func (this *floatOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from floatOps.java:51:5
	return this.Intersect(t1, this.complement(t2))
}

func (this *floatOps) complement(t SubtypeData) SubtypeData {
	// migrated from floatOps.java:56:5
	s := t.(floatSubtype)
	return createFloatSubtype((!s.allowed), s.Values())
}

func (this *floatOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from floatOps.java:62:5
	return notIsEmpty(cx, t)
}
