/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/internal/modules"
	"fmt"
)

func execBranch(branchTerm *bir.Branch, frame *Frame) *bir.BIRBasicBlock {
	opIndex := branchTerm.Op.Index
	value := frame.locals[opIndex]
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
	funcName := callInfo.Name.Value()
	values := make([]any, len(callInfo.Args))
	for i, op := range callInfo.Args {
		values[i] = frame.locals[op.Index]
	}
	orgName := callInfo.CalleePkg.OrgName.Value()
	pkgName := callInfo.CalleePkg.PkgName.Value()
	moduleKey := orgName + "/" + pkgName
	qualifiedName := moduleKey + ":" + funcName

	fn := reg.GetBIRFunction(qualifiedName)
	if fn != nil {
		result := executeFunction(*fn, values, reg)
		if callInfo.LhsOp != nil {
			frame.locals[callInfo.LhsOp.Index] = result
		}
		return callInfo.ThenBB
	}

	externFn := reg.GetNativeFunction(qualifiedName)
	if externFn != nil {
		result, err := externFn.Impl(values)
		if err != nil {
			panic(err)
		}
		if callInfo.LhsOp != nil {
			frame.locals[callInfo.LhsOp.Index] = result
		}
		return callInfo.ThenBB
	}

	panic("function not found: " + funcName)
}
