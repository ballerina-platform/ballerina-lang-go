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

import "ballerina-lang-go/common"

type FloatOps struct {
	CommonOps
}

var _ BasicTypeOps = &FloatOps{}

func NewFloatOps() FloatOps {
	this := FloatOps{}
	return this
}

func (this *FloatOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FloatOps.java:36:5
	var values []EnumerableType[float64]
	var v1 EnumerableSubtype[float64] = common.ToPointer(t1.(FloatSubtype))
	var v2 EnumerableSubtype[float64] = common.ToPointer(t2.(FloatSubtype))
	allowed := EnumerableSubtypeUnion(v1, v2, &values)
	return CreateFloatSubtype(allowed, values)
}

func (this *FloatOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FloatOps.java:44:5
	var values []EnumerableType[float64]
	var v1 EnumerableSubtype[float64] = common.ToPointer(t1.(FloatSubtype))
	var v2 EnumerableSubtype[float64] = common.ToPointer(t2.(FloatSubtype))
	allowed := EnumerableSubtypeIntersect(v1, v2, &values)
	return CreateFloatSubtype(allowed, values)
}

func (this *FloatOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from FloatOps.java:51:5
	return this.Intersect(t1, this.Complement(t2))
}

func (this *FloatOps) Complement(t SubtypeData) SubtypeData {
	// migrated from FloatOps.java:56:5
	s := t.(FloatSubtype)
	return CreateFloatSubtype((!s.allowed), s.Values())
}

func (this *FloatOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from FloatOps.java:62:5
	return notIsEmpty(cx, t)
}
