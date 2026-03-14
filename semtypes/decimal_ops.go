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

type DecimalOps struct {
	CommonOps
}

var _ BasicTypeOps = &DecimalOps{}

func NewDecimalOps() DecimalOps {
	this := DecimalOps{}
	return this
}

func (d *DecimalOps) Union(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from DecimalOps.java:36:5
	var values []EnumerableType[big.Rat]
	var v1 EnumerableSubtype[big.Rat] = new(t1.(DecimalSubtype))
	var v2 EnumerableSubtype[big.Rat] = new(t2.(DecimalSubtype))
	allowed := EnumerableSubtypeUnion(v1, v2, &values)
	return CreateDecimalSubtype(allowed, values)
}

func (d *DecimalOps) Intersect(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from DecimalOps.java:44:5
	var values []EnumerableType[big.Rat]
	var v1 EnumerableSubtype[big.Rat] = new(t1.(DecimalSubtype))
	var v2 EnumerableSubtype[big.Rat] = new(t2.(DecimalSubtype))
	allowed := EnumerableSubtypeIntersect(v1, v2, &values)
	return CreateDecimalSubtype(allowed, values)
}

func (d *DecimalOps) Diff(t1 SubtypeData, t2 SubtypeData) SubtypeData {
	// migrated from DecimalOps.java:52:5
	return d.Intersect(t1, d.Complement(t2))
}

func (d *DecimalOps) Complement(t SubtypeData) SubtypeData {
	// migrated from DecimalOps.java:57:5
	s := t.(DecimalSubtype)
	return CreateDecimalSubtype((!s.allowed), s.Values())
}

func (d *DecimalOps) IsEmpty(tc Context, t SubtypeData) bool {
	// migrated from DecimalOps.java:63:5
	return notIsEmpty(tc, t)
}
