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

type enumerableDecimal struct {
	value big.Rat
}

var _ enumerableType[big.Rat] = &enumerableDecimal{}

func (this *enumerableDecimal) Value() big.Rat {
	return this.value
}

func (t1 *enumerableDecimal) Compare(t2 enumerableType[big.Rat]) int {
	f1 := t1.Value()
	f2 := t2.Value()
	return f1.Cmp(&f2)
}

func newEnumerableDecimalFromBigDecimal(value big.Rat) enumerableDecimal {
	this := enumerableDecimal{}
	this.value = value
	return this
}

func enumerableDecimalFrom(d big.Rat) enumerableDecimal {
	// migrated from enumerableDecimal.java:34:5
	return newEnumerableDecimalFromBigDecimal(d)
}
