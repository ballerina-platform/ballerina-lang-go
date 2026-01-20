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

type MemoStatus uint

type BddMemo struct {
	isEmpty MemoStatus
}

const (
	MemoStatus_LOOP MemoStatus = iota
	MemoStatus_TRUE
	MemoStatus_FALSE
	MemoStatus_CYCLIC
	MemoStatus_PROVISIONAL
	MemoStatus_NULL
)

func NewBddMemo() BddMemo {
	this := BddMemo{}
	this.isEmpty = MemoStatus_NULL
	return this
}

func (this *BddMemo) SetIsEmpty(isEmpty bool) {
	// migrated from BddMemo.java:33:5
	if isEmpty {
		this.isEmpty = MemoStatus_TRUE
	} else {
		this.isEmpty = MemoStatus_FALSE
	}
}

func (this *BddMemo) IsEmpty() bool {
	// migrated from BddMemo.java:37:5
	return (this.isEmpty == MemoStatus_TRUE)
}
