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

func (this *ComplexSemType) all() BasicTypeBitSet {
	return this.allBitSet
}

func (this *ComplexSemType) some() BasicTypeBitSet {
	return this.someBitSet
}

func (this *ComplexSemType) subtypeDataList() []ProperSubtypeData {
	return this.dataList
}

func (this *ComplexSemType) equals(other *ComplexSemType) bool {
	return this == other || (this.allBitSet == other.allBitSet && this.someBitSet == other.someBitSet &&
		slices.Equal(this.dataList, other.dataList))
}
