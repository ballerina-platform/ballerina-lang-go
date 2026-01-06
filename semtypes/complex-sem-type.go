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

type ComplexSemType interface {
	SemType
	All() int
	Some() int
	SubtypeDataList() []ProperSubtypeData
}

func CreateComplexSemType(allBitset int, subtypeList ...BasicSubtype) ComplexSemType {
	// migrated from ComplexSemType.java:33:5
	return CreateComplexSemTypeWithAllBitSetSubtypeList(allBitset, subtypeList)
}

func CreateComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(allBitset int, someBitset int, subtypeDataList []ProperSubtypeData) ComplexSemType {
	// migrated from ComplexSemType.java:37:5
	if (allBitset == 0) && (someBitset == (1 << BT_CELL.Code)) {
		return CellSemTypeFrom(subtypeDataList)
	}
	return &complexSemTypeImpl{
		all:             allBitset,
		some:            someBitset,
		subtypeDataList: subtypeDataList,
	}
}

func CreateComplexSemTypeWithAllBitSetSubtypeList(allBitset int, subtypeList []BasicSubtype) ComplexSemType {
	// migrated from ComplexSemType.java:44:5
	some := 0
    var dataList []ProperSubtypeData
	for _, basicSubtype := range subtypeList {
		dataList = append(dataList, basicSubtype.SubtypeData)
		c := basicSubtype.BasicTypeCode.Code
		some = (some | (1 << c))
	}
	return CreateComplexSemTypeWithAllBitSetSomeBitSetSubtypeDataList(allBitset, some, dataList)
}
