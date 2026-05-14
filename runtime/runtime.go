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

package runtime

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/runtime/internal/exec"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

// Runtime represents a Ballerina runtime instance that owns a module registry
// and is used as the execution context for interpreting BIR packages.
type Runtime struct {
	env *extern.Env
}

// ModuleInitializer is a function that can install modules (e.g. stdlibs) into
// a runtime instance during its construction.
type ModuleInitializer func(*Runtime)

var moduleInitializers []ModuleInitializer

// NewRuntime constructs a new runtime with an empty registry and runs all
// registered module initializers.
func NewRuntime(platform pal.Platform, tyEnv semtypes.Env) *Runtime {
	registry := modules.NewRegistry()
	env := extern.InitEnv(platform, tyEnv, registry)
	rt := &Runtime{env}
	for _, init := range moduleInitializers {
		init(rt)
	}
	return rt
}

// Platform returns the platform configuration of this runtime instance.
func (rt *Runtime) Platform() pal.Platform {
	return rt.env.Platform
}

func (rt *Runtime) registry() *modules.Registry {
	return rt.env.Registry.(*modules.Registry)
}

// Interpret interprets a BIR package using this runtime instance.
func (rt *Runtime) Interpret(pkg bir.BIRPackage) (err error) {
	return exec.Interpret(pkg, rt.env)
}

// RegisterModuleInitializer registers a module initializer that will be invoked
// for every newly created runtime.
func RegisterModuleInitializer(init ModuleInitializer) {
	moduleInitializers = append(moduleInitializers, init)
}

// GetTypeEnv returns the semantic type environment.
func (rt *Runtime) GetTypeEnv() semtypes.Env {
	return rt.env.TypeEnv
}

// RegisterExternFunction registers a native (extern) function implementation in
// the given runtime instance so it can be called from interpreted BIR code.
func RegisterExternFunction(rt *Runtime, orgName string, moduleName string, funcName string, impl func(args []values.BalValue) (values.BalValue, error)) {
	rt.registry().RegisterExternFunction(orgName, moduleName, funcName, impl)
}

// RegisterExternClassDef registers a synthetic BIRClassDef for a Go-declared class so
// that execNewObject can resolve it. VTable entries have no BIR body; exec falls through
// to nativeFunctions for method dispatch.
func RegisterExternClassDef(rt *Runtime, def *bir.BIRClassDef) {
	rt.registry().RegisterExternClassDef(def)
}

// RegisterModuleGlobals makes module-level constants accessible at runtime.
// When Ballerina source code accesses an extern package's constant (e.g. http:LEADING),
// the BIR executor looks it up as a global variable in that package's module. Without
// registration, GetModule returns nil and causes a nil dereference panic.
func RegisterModuleGlobals(rt *Runtime, pkgId *model.PackageID, globals map[string]values.BalValue) {
	if existing := rt.registry().GetModule(pkgId); existing != nil {
		if existing.Globals == nil {
			existing.Globals = make(map[string]values.BalValue)
		}
		for k, v := range globals {
			existing.Globals[k] = v
		}
		return
	}
	rt.registry().RegisterModule(pkgId, &modules.BIRModule{Globals: globals})
}
