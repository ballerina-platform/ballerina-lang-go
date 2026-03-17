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

// Migrated from io.ballerina.types.CellSemType

type CellSemType struct {
	// Holding on to the single value instead of the array with a single value is more memory efficient. However, if
	// this start to cause problems in the future, we can change this to an array.
	subtypeData ProperSubtypeData
}

var _ ComplexSemType = &CellSemType{}

func NewCellSemType(subtypeDataList []ProperSubtypeData) CellSemType {
	// migrated from CellSemType.java:31
	// assert subtypeDataList.length == 1;
	if len(subtypeDataList) != 1 {
		panic("subtypeDataList length must be 1")
	}
	return CellSemType{
		subtypeData: subtypeDataList[0],
	}
}

func CellSemTypeFrom(subtypeDataList []ProperSubtypeData) CellSemType {
	// migrated from CellSemType.java:36
	return NewCellSemType(subtypeDataList)
}

func (this CellSemType) All() int {
	// migrated from CellSemType.java:46
	return 0
}

func (this CellSemType) Some() int {
	// migrated from CellSemType.java:51
	return CELL.bitset
}

func (this CellSemType) SubtypeDataList() []ProperSubtypeData {
	// migrated from CellSemType.java:56
	return []ProperSubtypeData{this.subtypeData}
}
