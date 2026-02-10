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

package exec

import (
	"ballerina-lang-go/bir"
	"math"
)

func execConstantLoad(constantLoad *bir.ConstantLoad, frame *Frame) {
	frame.SetOperand(constantLoad.LhsOp.Index, constantLoad.Value)
}

func execMove(moveIns *bir.Move, frame *Frame) {
	frame.SetOperand(moveIns.LhsOp.Index, frame.GetOperand(moveIns.RhsOp.Index))
}

func execNewArray(newArray *bir.NewArray, frame *Frame) {
	size := 0
	if newArray.SizeOp != nil {
		size = int(frame.GetOperand(newArray.SizeOp.Index).(int64))
	}
	arr := make([]any, size)
	if size > 0 {
		for i, value := range newArray.Values {
			arr[i] = frame.GetOperand(value.Index)
		}
	}
	lat := newArray.AtomicType
	for i := size; i < lat.Members.FixedLength; i++ {
		ty := lat.MemberAt(i)
		val := fillMember(ty)
		if val == NeverValue {
			panic("never value encountered")
		}
		arr = append(arr, val)
	}
	frame.SetOperand(newArray.LhsOp.Index, &arr)
}

const MaxArraySize = math.MaxInt32

func execArrayStore(access *bir.FieldAccess, frame *Frame) {
	arrPtr := frame.GetOperand(access.LhsOp.Index).(*[]any)
	arr := *arrPtr
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 {
		panic("index out of bounds")
	}

	if idx >= MaxArraySize {
		panic("list too long")
	}
	for idx > len(arr)-1 {
		// FIXME: properly do the filling read here instead
		arr = append(arr, nil)
	}
	arr[idx] = frame.GetOperand(access.RhsOp.Index)
	frame.SetOperand(access.LhsOp.Index, &arr)
	*arrPtr = arr
}

func execArrayLoad(access *bir.FieldAccess, frame *Frame) {
	arr := *(frame.GetOperand(access.RhsOp.Index).(*[]any))
	idx := int(frame.GetOperand(access.KeyOp.Index).(int64))
	if idx < 0 || idx >= len(arr) {
		panic("index out of bounds")
	}
	frame.SetOperand(access.LhsOp.Index, arr[idx])
}
