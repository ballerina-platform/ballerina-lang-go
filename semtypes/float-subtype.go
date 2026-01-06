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

type FloatSubtype struct {
	EnumerableSubtype[float64]
	allowed bool
	values  []EnumerableFloat
}

var _ ProperSubtypeData = &FloatSubtype{}

func newFloatSubtypeFromBoolEnumerableFloat(allowed bool, value EnumerableFloat) FloatSubtype {
	this := FloatSubtype{}
	this.allowed = allowed
	this.values = []EnumerableFloat{value}
	return this
}

func newFloatSubtypeFromBoolEnumerableFloats(allowed bool, values []EnumerableType[float64]) FloatSubtype {
	this := FloatSubtype{}
	this.allowed = allowed
	var floats []EnumerableFloat
	for _, value := range values {
		floats = append(floats, EnumerableFloatFrom(value.Value()))
	}
	this.values = floats
	return this
}

func FloatConst(value float64) SemType {
	return basicSubtype(BT_FLOAT, newFloatSubtypeFromBoolEnumerableFloat(true, EnumerableFloatFrom(value)))
}

func FloatSubtypeSingleValue(d SubtypeData) common.Optional[float64] {
	if _, ok := d.(AllOrNothingSubtype); ok {
		return common.OptionalEmpty[float64]()
	}
	v := d.(FloatSubtype)
	if !v.allowed {
		return common.OptionalEmpty[float64]()
	}
	if len(v.values) != 1 {
		return common.OptionalEmpty[float64]()
	}
	return common.OptionalOf(v.values[0].value)
}

func FloatSubtypeContains(d SubtypeData, f EnumerableFloat) bool {
	// migrated from FloatSubtype.java:72:5
	if allOrNothingSubtype, ok := d.(AllOrNothingSubtype); ok {
		return allOrNothingSubtype.IsAllSubtype()
	}
	v := d.(FloatSubtype)
	for _, val := range v.values {
		if val == f {
			return v.allowed
		}
	}
	return (!v.allowed)
}

func CreateFloatSubtype(allowed bool, values []EnumerableType[float64]) SubtypeData {
	// migrated from FloatSubtype.java:86:5
	if len(values) == 0 {
		if allowed {
			return CreateNothing()
		} else {
			return CreateAll()
		}
	}
	return newFloatSubtypeFromBoolEnumerableFloats(allowed, values)
}

func (this *FloatSubtype) Allowed() bool {
	// migrated from FloatSubtype.java:93:5
	return this.allowed
}

func (this *FloatSubtype) Values() []EnumerableType[float64] {
	// migrated from FloatSubtype.java:98:5
	var values []EnumerableType[float64]
	for _, value := range this.values {
		values = append(values, &value)
	}
	return values
}
