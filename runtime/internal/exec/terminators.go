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
	"ballerina-lang-go/values"
)

func execBranch(branchTerm *bir.Branch, frame *Frame, reg *modules.Registry) *bir.BIRBasicBlock {
	if getOperandValue(branchTerm.Op, frame, reg).(bool) {
		return branchTerm.TrueBB
	}
	return branchTerm.FalseBB
}

func execCall(callInfo *bir.Call, frame *Frame, reg *modules.Registry, callStack *callStack) *bir.BIRBasicBlock {
	args := extractArgs(callInfo.Args, frame, reg)
	result := executeCall(callInfo, args, reg, callStack)
	if callInfo.LhsOp != nil {
		setOperandValue(callInfo.LhsOp, frame, reg, result)
	}
	return callInfo.ThenBB
}

func executeCall(callInfo *bir.Call, args []values.BalValue, reg *modules.Registry, callStack *callStack) values.BalValue {
	if callInfo.IsVirtual {
		fn := resolveVirtualCallTarget(callInfo, args, reg)
		result := executeFunction(*fn, args, reg, callStack)
		if callInfo.Name.Value() == "init" {
			if result != nil {
				return result
			}
			return args[0]
		}
		return result
	}
	if callInfo.CachedBIRFunc != nil {
		return executeFunction(*callInfo.CachedBIRFunc, args, reg, callStack)
	}
	if callInfo.CachedNativeFunc != nil {
		result, err := callInfo.CachedNativeFunc(args)
		if err != nil {
			panic(err)
		}
		return result
	}
	return lookupAndExecute(callInfo, args, reg, callStack)
}

func resolveVirtualCallTarget(callInfo *bir.Call, args []values.BalValue, reg *modules.Registry) *bir.BIRFunction {
	if callInfo.CachedBIRFunc != nil {
		return callInfo.CachedBIRFunc
	}
	receiverObj := args[0].(*values.Object)
	lookupKey, found := receiverObj.MethodLookupKey(string(callInfo.Name))
	if !found {
		panic("function not found: " + callInfo.Name.Value())
	}
	fn := reg.GetBIRFunction(lookupKey)
	callInfo.CachedBIRFunc = fn
	return fn
}

func lookupAndExecute(callInfo *bir.Call, args []values.BalValue, reg *modules.Registry, callStack *callStack) values.BalValue {
	fn := reg.GetBIRFunction(callInfo.FunctionLookupKey)
	if fn != nil {
		callInfo.CachedBIRFunc = fn
		return executeFunction(*fn, args, reg, callStack)
	}
	externFn := reg.GetNativeFunction(callInfo.FunctionLookupKey)
	if externFn != nil {
		callInfo.CachedNativeFunc = externFn.Impl
		result, err := externFn.Impl(args)
		if err != nil {
			panic(err)
		}
		return result
	}
	panic("function not found: " + callInfo.Name.Value())
}

func execFpCall(callInfo *bir.Call, frame *Frame, reg *modules.Registry, callStack *callStack) *bir.BIRBasicBlock {
	args := extractArgs(callInfo.Args, frame, reg)
	fnValue := getOperandValue(callInfo.FpOperand, frame, reg).(*values.Function)
	lookupKey := fnValue.LookupKey
	fn := reg.GetBIRFunction(lookupKey)
	var result values.BalValue
	if fn != nil {
		result = executeFunction(*fn, args, reg, callStack)
	} else {
		externFn := reg.GetNativeFunction(lookupKey)
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
		setOperandValue(callInfo.LhsOp, frame, reg, result)
	}
	return callInfo.ThenBB
}

func extractArgs(args []bir.BIROperand, frame *Frame, reg *modules.Registry) []values.BalValue {
	values := make([]values.BalValue, len(args))
	for i, op := range args {
		values[i] = getOperandValue(&op, frame, reg)
	}
	return values
}
