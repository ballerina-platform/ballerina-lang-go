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
	"ballerina-lang-go/values"
)

func execBranch(ctx *Context, branchTerm *bir.Branch, frame *Frame) *bir.BIRBasicBlock {
	if getOperandValue(ctx, branchTerm.Op, frame).(bool) {
		return branchTerm.TrueBB
	}
	return branchTerm.FalseBB
}

func execCall(ctx *Context, callInfo *bir.Call, frame *Frame) *bir.BIRBasicBlock {
	args := extractArgs(ctx, callInfo.Args, frame)
	result := executeCall(ctx, callInfo, args)
	if callInfo.LhsOp != nil {
		setOperandValue(ctx, callInfo.LhsOp, frame, result)
	}
	return callInfo.ThenBB
}

func executeCall(ctx *Context, callInfo *bir.Call, args []values.BalValue) values.BalValue {
	if callInfo.IsMethodCall {
		fn := resolveObjectMethod(ctx, callInfo, args)
		return executeFunction(ctx, *fn, args, nil)
	}
	if callInfo.CachedBIRFunc != nil {
		return executeFunction(ctx, *callInfo.CachedBIRFunc, args, nil)
	}
	if callInfo.CachedNativeFunc != nil {
		result, err := callInfo.CachedNativeFunc(args)
		if err != nil {
			panic(err)
		}
		return result
	}
	return lookupAndExecute(ctx, callInfo, args)
}

func resolveObjectMethod(ctx *Context, callInfo *bir.Call, args []values.BalValue) *bir.BIRFunction {
	receiverObj := args[0].(*values.Object)
	lookupKey, found := receiverObj.MethodLookupKey(string(callInfo.Name))
	if !found {
		panic("function not found: " + callInfo.Name.Value())
	}

	// The same call site can be polymorphic across executions (e.g., iterating over a list
	// of objects with different concrete types). Cache only when it matches the receiver.
	if callInfo.CachedBIRFunc != nil {
		if callInfo.CachedMethodLookupKey == lookupKey {
			return callInfo.CachedBIRFunc
		}
	}

	fn := ctx.GetBIRFunction(lookupKey)
	callInfo.CachedBIRFunc = fn
	callInfo.CachedMethodLookupKey = lookupKey
	return fn
}

func lookupAndExecute(ctx *Context, callInfo *bir.Call, args []values.BalValue) values.BalValue {
	fn := ctx.GetBIRFunction(callInfo.FunctionLookupKey)
	if fn != nil {
		callInfo.CachedBIRFunc = fn
		return executeFunction(ctx, *fn, args, nil)
	}
	externFn := ctx.GetNativeFunction(callInfo.FunctionLookupKey)
	if externFn != nil {
		callInfo.CachedNativeFunc = externFn.Impl
		result, err := externFn.Impl(args)
		if err != nil {
			panic(err)
		}
		return result
	}
	panic(values.NewErrorWithMessage("function not found: " + callInfo.Name.Value()))
}

func execFpCall(ctx *Context, callInfo *bir.Call, frame *Frame) *bir.BIRBasicBlock {
	args := extractArgs(ctx, callInfo.Args, frame)
	fnValue := getOperandValue(ctx, callInfo.FpOperand, frame).(*values.Function)
	lookupKey := fnValue.LookupKey
	var parentFrame *Frame
	if fnValue.ParentFrame != nil {
		parentFrame = fnValue.ParentFrame.(*Frame)
	}
	fn := ctx.GetBIRFunction(lookupKey)
	var result values.BalValue
	if fn != nil {
		result = executeFunction(ctx, *fn, args, parentFrame)
	} else {
		externFn := ctx.GetNativeFunction(lookupKey)
		if externFn != nil {
			var err error
			result, err = externFn.Impl(args)
			if err != nil {
				panic(err)
			}
		} else {
			panic("function not found: " + callInfo.Name.Value())
		}
	}
	if callInfo.LhsOp != nil {
		setOperandValue(ctx, callInfo.LhsOp, frame, result)
	}
	return callInfo.ThenBB
}

func extractArgs(ctx *Context, args []bir.BIROperand, frame *Frame) []values.BalValue {
	values := make([]values.BalValue, len(args))
	for i, op := range args {
		values[i] = getOperandValue(ctx, &op, frame)
	}
	return values
}

func execPanic(ctx *Context, panicTerm *bir.Panic, frame *Frame) *bir.BIRBasicBlock {
	errVal := getOperandValue(ctx, panicTerm.ErrorOp, frame)
	panic(errVal)
}
