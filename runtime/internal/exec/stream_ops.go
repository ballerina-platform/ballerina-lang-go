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
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/values"
)

func execNewStream(ctx *extern.Context, instr *bir.NewStream, frame *Frame) {
	impl := getOperandValue(ctx, instr.ImplOp, frame).(*values.Object)
	stream := values.NewStream(instr.StreamType, impl)
	setOperandValue(ctx, instr.LhsOp, frame, stream)
}

func execStreamNext(ctx *extern.Context, instr *bir.StreamNext, frame *Frame) {
	stream := getOperandValue(ctx, instr.StreamOp, frame).(*values.Stream)
	handle := LookupObjectMethod(ctx, stream.Underlying, "next")
	if handle == nil {
		panic(values.NewErrorWithMessage("stream implementor missing 'next' method"))
	}
	result := InvokeObjectMethod(ctx, handle, []values.BalValue{stream.Underlying})
	setOperandValue(ctx, instr.LhsOp, frame, result)
}

func execStreamClose(ctx *extern.Context, instr *bir.StreamClose, frame *Frame) {
	stream := getOperandValue(ctx, instr.StreamOp, frame).(*values.Stream)
	handle := LookupObjectMethod(ctx, stream.Underlying, "close")
	if handle == nil {
		setOperandValue(ctx, instr.LhsOp, frame, nil)
		return
	}
	result := InvokeObjectMethod(ctx, handle, []values.BalValue{stream.Underlying})
	setOperandValue(ctx, instr.LhsOp, frame, result)
}
