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

package modules

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

type Registry struct {
	birFunctions    map[string]*bir.BIRFunction
	nativeFunctions map[string]*ExternFunction
	typeCtx         semtypes.Context
}

func NewRegistry() *Registry {
	return &Registry{
		birFunctions:    make(map[string]*bir.BIRFunction),
		nativeFunctions: make(map[string]*ExternFunction),
	}
}

func (r *Registry) RegisterModule(id *model.PackageID, m *BIRModule) *BIRModule {
	if m.Pkg != nil && m.Pkg.Functions != nil {
		for i := range m.Pkg.Functions {
			fn := &m.Pkg.Functions[i]
			r.birFunctions[fn.FunctionLookupKey] = fn
		}
	}
	return m
}

func (r *Registry) SetTypeCtx(cx semtypes.Context) {
	r.typeCtx = cx
}

func (r *Registry) GetTypeCtx() semtypes.Context {
	return r.typeCtx
}

func (r *Registry) RegisterExternFunction(orgName string, moduleName string, funcName string, impl func(args []values.BalValue) (values.BalValue, error)) {
	externFn := &ExternFunction{
		Name: funcName,
		Impl: impl,
	}
	moduleKey := orgName + "/" + moduleName
	qualifiedName := moduleKey + ":" + funcName
	r.nativeFunctions[qualifiedName] = externFn
	r.nativeFunctions[funcName] = externFn
}

func (r *Registry) GetBIRFunction(funcName string) *bir.BIRFunction {
	return r.birFunctions[funcName]
}

func (r *Registry) GetNativeFunction(funcName string) *ExternFunction {
	return r.nativeFunctions[funcName]
}
