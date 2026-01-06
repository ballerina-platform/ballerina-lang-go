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

type ListMemberTypes struct {
	Ranges   []Range
	SemTypes []SemType
}

func NewListMemberTypesFromRangesSemTypes(ranges []Range, semTypes []SemType) ListMemberTypes {
	this := ListMemberTypes{}
	rangesCopy := make([]Range, len(ranges))
	copy(rangesCopy, ranges)
	semTypesCopy := make([]SemType, len(semTypes))
	copy(semTypesCopy, semTypes)
	this.Ranges = rangesCopy
	this.SemTypes = semTypesCopy
	return this
}

func ListMemberTypesFrom(ranges []Range, semTypes []SemType) ListMemberTypes {
	return NewListMemberTypesFromRangesSemTypes(ranges, semTypes)
}
