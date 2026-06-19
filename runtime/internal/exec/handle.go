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

// methodHandleImpl is the concrete payload behind extern.MethodHandle. The
// flavour of the resolved method (BIR, native, resource) is encoded entirely
// in the captured closure; no flavour discriminator is stored on the struct.
type methodHandleImpl struct {
	invoke func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error)
}

func newBIRHandle(fn *bir.BIRFunction) *methodHandleImpl {
	return &methodHandleImpl{
		invoke: func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return executeFunction(ctx, *fn, args, nil), nil
		},
	}
}

func newNativeHandle(fn extern.NativeFunc) *methodHandleImpl {
	return &methodHandleImpl{
		invoke: func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			return fn(ctx, args)
		},
	}
}

func newResourceHandle(receiver *values.Object, match *values.ResourceEntry, path []values.BalValue) *methodHandleImpl {
	return &methodHandleImpl{
		invoke: func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			full := buildResourceCallArgs(ctx, receiver, match, path, args)
			return lookupAndExecute(ctx, nil, full, match.FunctionLookupKey)
		},
	}
}
