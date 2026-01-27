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

package api

import (
	"ballerina-lang-go/runtime/internal/modules"
)

type ModuleInitializer func(*Runtime)

var moduleInitializers []ModuleInitializer

type Runtime struct {
	Registry *modules.Registry
}

func RegisterModuleInitializer(init ModuleInitializer) {
	moduleInitializers = append(moduleInitializers, init)
}

func NewRuntime() *Runtime {
	rt := &Runtime{
		Registry: modules.NewRegistry(),
	}
	for _, init := range moduleInitializers {
		init(rt)
	}
	return rt
}

func RegisterExternFunction(reg *modules.Registry, orgName string, moduleName string, funcName string, impl func(args []any) (any, error)) {
	reg.RegisterExternFunction(orgName, moduleName, funcName, impl)
}
