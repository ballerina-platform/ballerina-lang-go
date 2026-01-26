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
	"ballerina-lang-go/model"
	"fmt"
)

// Registry holds all loaded modules
type Registry struct {
	// birModules holds BIR-backed modules keyed by a stable string derived from PackageID.
	birModules map[*model.PackageID]*BIRModule
	// nativeModules holds Go-native modules keyed by their module name.
	nativeModules map[string]*NativeModule
}

var globalRegistry *Registry

func init() {
	// Initialize registry on package load to support native module registration in init() functions
	InitRegistry()
}

func GetRegistry() *Registry {
	return globalRegistry
}

// InitRegistry initializes a fresh registry
func InitRegistry() {
	globalRegistry = &Registry{
		birModules:    make(map[*model.PackageID]*BIRModule),
		nativeModules: make(map[string]*NativeModule),
	}
}

// RegisterModule registers a BIR-backed module using its PackageID.
func (r *Registry) RegisterModule(id *model.PackageID, m *BIRModule) *BIRModule {
	if existing, exists := r.birModules[id]; exists {
		return existing
	}
	r.birModules[id] = m
	return m
}

// RegisterNativeModule registers a Go-native module using orgName and moduleName.
// The moduleKey is calculated as "orgName/moduleName" internally.
// Returns the moduleKey that was used to register the module.
func (r *Registry) RegisterNativeModule(orgName, moduleName string, m *NativeModule) string {
	moduleKey := orgName + "/" + moduleName
	if _, exists := r.nativeModules[moduleKey]; exists {
		return moduleKey
	}
	r.nativeModules[moduleKey] = m
	return moduleKey
}

// GetBIRModule returns a BIR-backed module for the given PackageID.
func (r *Registry) GetBIRModule(id *model.PackageID) (*BIRModule, error) {
	if m, ok := r.birModules[id]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("module not found for given PackageID")
}

func (r *Registry) RegisterExternFunction(moduleName, funcName string, impl func(args []any) (any, error)) error {
	// Extern functions are registered only on native (Go) modules.
	m, ok := r.nativeModules[moduleName]
	if !ok {
		return fmt.Errorf("module not found: %s", moduleName)
	}
	if m.ExternFunctions == nil {
		m.ExternFunctions = make(map[string]*ExternFunction)
	}
	m.ExternFunctions[funcName] = &ExternFunction{
		Name: funcName,
		Impl: impl,
	}
	return nil
}

func (r *Registry) GetNativeModuleByPkgID(id *model.PackageID) (*NativeModule, error) {
	nameKey := id.OrgName.Value() + "/" + id.PkgName.Value()
	if m, ok := r.nativeModules[nameKey]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("module not found: %s", nameKey)
}

// GetAllBIRModules returns all registered BIR modules
func (r *Registry) GetAllBIRModules() []*BIRModule {
	modules := make([]*BIRModule, 0, len(r.birModules))
	for _, m := range r.birModules {
		modules = append(modules, m)
	}
	return modules
}

// GetAllNativeModules returns all registered native modules
func (r *Registry) GetAllNativeModules() []*NativeModule {
	modules := make([]*NativeModule, 0, len(r.nativeModules))
	for _, m := range r.nativeModules {
		modules = append(modules, m)
	}
	return modules
}
