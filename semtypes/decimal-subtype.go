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
	"ballerina-lang-go/common"
	"math/big"
)

type DecimalSubtype struct {
	EnumerableSubtype[big.Rat]
	allowed bool
	values  []EnumerableDecimal
}

var _ ProperSubtypeData = &DecimalSubtype{}

func newDecimalSubtypeFromBoolEnumerableDecimal(allowed bool, value EnumerableDecimal) DecimalSubtype {
	this := DecimalSubtype{}
	this.allowed = allowed
	this.values = []EnumerableDecimal{value}
	return this
}

func newDecimalSubtypeFromBoolEnumerableDecimals(allowed bool, values []EnumerableType[big.Rat]) DecimalSubtype {
	this := DecimalSubtype{}
	this.allowed = allowed
	var decimals []EnumerableDecimal
	for _, value := range values {
		decimals = append(decimals, EnumerableDecimalFrom(value.Value()))
	}
	this.values = decimals
	return this
}

func DecimalConst(value big.Rat) SemType {
	return basicSubtype(BT_DECIMAL, newDecimalSubtypeFromBoolEnumerableDecimal(true, EnumerableDecimalFrom(value)))
}

func DecimalSubtypeSingleValue(d SubtypeData) common.Optional[big.Rat] {
	if _, ok := d.(AllOrNothingSubtype); ok {
		return common.OptionalEmpty[big.Rat]()
	}
	v := d.(DecimalSubtype)
	if !v.allowed {
		return common.OptionalEmpty[big.Rat]()
	}
	if len(v.values) != 1 {
		return common.OptionalEmpty[big.Rat]()
	}
	return common.OptionalOf(v.values[0].value)
}

func DecimalSubtypeContains(d SubtypeData, f EnumerableDecimal) bool {
	// migrated from DecimalSubtype.java:73:5
	if allOrNothingSubtype, ok := d.(AllOrNothingSubtype); ok {
		return allOrNothingSubtype.IsAllSubtype()
	}
	v := d.(DecimalSubtype)
	for _, val := range v.values {
		if val.Compare(&f) == 0 {
			return v.allowed
		}
	}
	return (!v.allowed)
}

func CreateDecimalSubtype(allowed bool, values []EnumerableType[big.Rat]) SubtypeData {
	// migrated from DecimalSubtype.java:87:5
	if len(values) == 0 {
		if allowed {
			return CreateNothing()
		} else {
			return CreateAll()
		}
	}
	return newDecimalSubtypeFromBoolEnumerableDecimals(allowed, values)
}

func (this *DecimalSubtype) Allowed() bool {
	// migrated from DecimalSubtype.java:94:5
	return this.allowed
}

func (this *DecimalSubtype) Values() []EnumerableType[big.Rat] {
	// migrated from DecimalSubtype.java:99:5
	var values []EnumerableType[big.Rat]
	for _, value := range this.values {
		values = append(values, &value)
	}
	return values
}
