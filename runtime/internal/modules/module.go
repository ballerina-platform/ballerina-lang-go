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

package modules

import (
	"ballerina-lang-go/bir"
)

type ExternFunction struct {
	Name string
	Impl func(args []any) (any, error)
}

// BIRModule represents a module backed by a BIR package and its functions.
type BIRModule struct {
	Pkg       *bir.BIRPackage
	Functions map[string]*bir.BIRFunction
}

// NativeModule represents a module that only contains native/extern functions.
type NativeModule struct {
	ExternFunctions map[string]*ExternFunction
}

// NewBIRModule creates a BIRModule backed by a BIR package.
func NewBIRModule(pkg *bir.BIRPackage) *BIRModule {
	functions := make(map[string]*bir.BIRFunction)
	if pkg != nil && pkg.PackageID != nil {
		if funcs := pkg.Functions; funcs != nil {
			for i := range funcs {
				fn := funcs[i]
				funcName := fn.Name.Value()
				birFnPtr := &fn
				functions[funcName] = birFnPtr
			}
		}
	}
	return &BIRModule{
		Pkg:       pkg,
		Functions: functions,
	}
}

// NewNativeModule creates a module used for native (Go) extern functions only.
func NewNativeModule() *NativeModule {
	return &NativeModule{
		ExternFunctions: make(map[string]*ExternFunction),
	}
}
