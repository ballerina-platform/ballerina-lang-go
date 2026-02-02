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
	"ballerina-lang-go/runtime/internal/modules"
	"fmt"
)

func execBranch(branchTerm *bir.Branch, frame *Frame) *bir.BIRBasicBlock {
	opIndex := branchTerm.Op.Index
	value := frame.GetOperand(opIndex)
	cond, ok := value.(bool)
	if !ok {
		panic(fmt.Sprintf("invalid branch condition type at index %d: %T (expected bool)", opIndex, value))
	}
	if cond {
		return branchTerm.TrueBB
	}
	return branchTerm.FalseBB
}

func execCall(callInfo *bir.Call, frame *Frame, reg *modules.Registry) *bir.BIRBasicBlock {
	values := extractArgs(callInfo.Args, frame)
	result := executeCall(callInfo, values, reg)
	if callInfo.LhsOp != nil {
		frame.SetOperand(callInfo.LhsOp.Index, result)
	}
	return callInfo.ThenBB
}

func executeCall(callInfo *bir.Call, values []any, reg *modules.Registry) any {
	if callInfo.CachedBIRFunc != nil {
		return executeFunction(*callInfo.CachedBIRFunc, values, reg)
	}
	if callInfo.CachedNativeFunc != nil {
		result, err := callInfo.CachedNativeFunc(values)
		if err != nil {
			panic(err)
		}
		return result
	}
	return lookupAndExecute(callInfo, values, reg)
}

func lookupAndExecute(callInfo *bir.Call, values []any, reg *modules.Registry) any {
	lookupKey := callInfo.FunctionLookupKey
	fn := reg.GetBIRFunction(lookupKey)
	if fn != nil {
		callInfo.CachedBIRFunc = fn
		return executeFunction(*fn, values, reg)
	}
	externFn := reg.GetNativeFunction(lookupKey)
	if externFn != nil {
		callInfo.CachedNativeFunc = externFn.Impl
		result, err := externFn.Impl(values)
		if err != nil {
			panic(err)
		}
		return result
	}
	panic("function not found: " + callInfo.Name.Value())
}

func extractArgs(args []bir.BIROperand, frame *Frame) []any {
	values := make([]any, len(args))
	for i, op := range args {
		values[i] = frame.GetOperand(op.Index)
	}
	return values
}
