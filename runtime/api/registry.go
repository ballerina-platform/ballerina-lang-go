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

import "ballerina-lang-go/runtime/internal/modules"

type NativeModule struct {
	Name string
}

func RegisterNativeModule(orgName, moduleName string) *NativeModule {
	// Native modules are backed by Go implementations only.
	moduleKey := modules.GetRegistry().RegisterNativeModule(orgName, moduleName, modules.NewNativeModule())
	return &NativeModule{Name: moduleKey}
}

func RegisterExternFunction(module *NativeModule, funcName string, impl func(args []any) (any, error)) {
	modules.GetRegistry().RegisterExternFunction(module.Name, funcName, impl)
}
