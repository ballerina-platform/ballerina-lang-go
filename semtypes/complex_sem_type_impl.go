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

import "slices"

type ComplexSemType struct {
	allBitSet  BasicTypeBitSet
	someBitSet BasicTypeBitSet
	dataList   []ProperSubtypeData
}

var _ SemType = &ComplexSemType{}

func (c *ComplexSemType) all() BasicTypeBitSet {
	return c.allBitSet
}

func (c *ComplexSemType) some() BasicTypeBitSet {
	return c.someBitSet
}

func (c *ComplexSemType) subtypeDataList() []ProperSubtypeData {
	return c.dataList
}

func (c *ComplexSemType) equals(other *ComplexSemType) bool {
	return c == other || (c.allBitSet == other.allBitSet && c.someBitSet == other.someBitSet &&
		slices.Equal(c.dataList, other.dataList))
}
