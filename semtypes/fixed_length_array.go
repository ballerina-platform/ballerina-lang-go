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

type FixedLengthArray struct {
	Initial     []CellSemType
	FixedLength int
}

func NewFixedLengthArrayFromInitialFixedLength(initial []CellSemType, fixedLength int) FixedLengthArray {
	this := FixedLengthArray{}
	copiedInitial := make([]CellSemType, len(initial))
	copy(copiedInitial, initial)
	common.Assert(fixedLength >= 0)
	this.Initial = copiedInitial
	this.FixedLength = fixedLength
	return this
}

func FixedLengthArrayFrom(initial []CellSemType, fixedLength int) FixedLengthArray {
	// migrated from FixedLengthArray.java:45:5
	return NewFixedLengthArrayFromInitialFixedLength(initial, fixedLength)
}

func FixedLengthArrayEmpty() FixedLengthArray {
	// migrated from FixedLengthArray.java:53:5
	return FixedLengthArrayFrom(nil, 0)
}
