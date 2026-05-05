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
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/values"
)

// MethodHandle is an opaque, non-nil reference to a resolved object method
// (BIR or native). Obtain one from LookupObjectMethod and pass it to
// InvokeObjectMethod.
type MethodHandle struct {
	lookupKey string
	birFunc   *bir.BIRFunction
	nativeFn  extern.NativeFunc
}

// LookupObjectMethod resolves methodName on obj. Returns nil if obj has no
// such method.
func LookupObjectMethod(ctx *extern.Context, obj *values.Object, methodName string) *MethodHandle {
	lookupKey, found := obj.MethodLookupKey(methodName)
	if !found {
		return nil
	}
	reg := ctx.Env.Registry.(*modules.Registry)
	handle := &MethodHandle{lookupKey: lookupKey}
	if fn := reg.GetBIRFunction(lookupKey); fn != nil {
		handle.birFunc = fn
		return handle
	}
	if externFn := reg.GetNativeFunction(lookupKey); externFn != nil {
		handle.nativeFn = externFn.Impl
		return handle
	}
	return nil
}

// InvokeObjectMethod calls a previously looked up method. args is the full
// argument list including the receiver as args[0].
func InvokeObjectMethod(ctx *extern.Context, h *MethodHandle, args []values.BalValue) values.BalValue {
	if h.birFunc != nil {
		return executeFunction(ctx, *h.birFunc, args, nil)
	}
	if h.nativeFn != nil {
		result, err := h.nativeFn(ctx, args)
		if err != nil {
			panic(err)
		}
		return result
	}
	panic(values.NewErrorWithMessage("unexpected function handle: " + h.lookupKey))
}
