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

import (
	"math/big"
)

type decimalOps struct {
	CommonOps
}

var _ BasicTypeOps = &decimalOps{}

func newDecimalOps() decimalOps {
	this := decimalOps{}
	return this
}

func (this *decimalOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from decimalOps.java:36:5
	var values []enumerableType[big.Rat]
	var v1 enumerableSubtype[big.Rat] = new(t1.(decimalSubtype))
	var v2 enumerableSubtype[big.Rat] = new(t2.(decimalSubtype))
	allowed := enumerableSubtypeUnion(v1, v2, &values)
	return createDecimalSubtype(allowed, values)
}

func (this *decimalOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from decimalOps.java:44:5
	var values []enumerableType[big.Rat]
	var v1 enumerableSubtype[big.Rat] = new(t1.(decimalSubtype))
	var v2 enumerableSubtype[big.Rat] = new(t2.(decimalSubtype))
	allowed := enumerableSubtypeIntersect(v1, v2, &values)
	return createDecimalSubtype(allowed, values)
}

func (this *decimalOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from decimalOps.java:52:5
	return this.Intersect(t1, this.complement(t2))
}

func (this *decimalOps) complement(t SubtypeData) SubtypeData {
	// migrated from decimalOps.java:57:5
	s := t.(decimalSubtype)
	return createDecimalSubtype((!s.allowed), s.Values())
}

func (this *decimalOps) IsEmpty(tc Context, t SubtypeData) bool {
	// migrated from decimalOps.java:63:5
	return notIsEmpty(tc, t)
}
