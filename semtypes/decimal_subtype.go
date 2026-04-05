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

	"ballerina-lang-go/common"
)

type decimalSubtype struct {
	allowed bool
	values  []enumerableDecimal
}

var _ ProperSubtypeData = &decimalSubtype{}

func newDecimalSubtypeFromBoolEnumerableDecimal(allowed bool, value enumerableDecimal) decimalSubtype {
	this := decimalSubtype{}
	this.allowed = allowed
	this.values = []enumerableDecimal{value}
	return this
}

func newDecimalSubtypeFromBoolEnumerableDecimals(allowed bool, values []enumerableType[big.Rat]) decimalSubtype {
	this := decimalSubtype{}
	this.allowed = allowed
	var decimals []enumerableDecimal
	for _, value := range values {
		decimals = append(decimals, enumerableDecimalFrom(value.Value()))
	}
	this.values = decimals
	return this
}

func DecimalConst(value big.Rat) SemType {
	return getBasicSubtype(BTDecimal, newDecimalSubtypeFromBoolEnumerableDecimal(true, enumerableDecimalFrom(value)))
}

func decimalSubtypeSingleValue(d SubtypeData) common.Optional[big.Rat] {
	if _, ok := d.(allOrNothingSubtype); ok {
		return common.OptionalEmpty[big.Rat]()
	}
	v := d.(decimalSubtype)
	if !v.allowed {
		return common.OptionalEmpty[big.Rat]()
	}
	if len(v.values) != 1 {
		return common.OptionalEmpty[big.Rat]()
	}
	return common.OptionalOf(v.values[0].value)
}

func decimalSubtypeContains(d SubtypeData, f enumerableDecimal) bool {
	if allOrNothingSubtype, ok := d.(allOrNothingSubtype); ok {
		return allOrNothingSubtype.IsAllSubtype()
	}
	v := d.(decimalSubtype)
	for _, val := range v.values {
		if val.Compare(&f) == 0 {
			return v.allowed
		}
	}
	return (!v.allowed)
}

func createDecimalSubtype(allowed bool, values []enumerableType[big.Rat]) ProperSubtypeData {
	if len(values) == 0 {
		if allowed {
			return createNothing()
		} else {
			return createAll()
		}
	}
	return newDecimalSubtypeFromBoolEnumerableDecimals(allowed, values)
}

func (this *decimalSubtype) Allowed() bool {
	return this.allowed
}

func (this *decimalSubtype) Values() []enumerableType[big.Rat] {
	var values []enumerableType[big.Rat]
	for _, value := range this.values {
		values = append(values, &value)
	}
	return values
}
