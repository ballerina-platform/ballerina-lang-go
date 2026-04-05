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

import "math"

type enumerableFloat struct {
	value float64
}

var _ enumerableType[float64] = &enumerableFloat{}

func (this *enumerableFloat) Value() float64 {
	return this.value
}

func (t1 *enumerableFloat) Compare(t2 enumerableType[float64]) int {
	f1 := t1.Value()
	f2 := t2.Value()
	if bFloatEq(f1, f2) {
		return eq
	} else if math.IsNaN(f1) {
		return lt
	} else if math.IsNaN(f2) {
		return gt
	} else if f1 < f2 {
		return lt
	}
	return gt
}

func newEnumerableFloatFromFloat64(value float64) enumerableFloat {
	this := enumerableFloat{}
	this.value = value
	return this
}

func enumerableFloatFrom(d float64) enumerableFloat {
	return newEnumerableFloatFromFloat64(d)
}
