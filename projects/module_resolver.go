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

package projects

// moduleLoadRequestKey is used as a map key for caching module load responses.
type moduleLoadRequestKey struct {
	orgName    string
	moduleName string
}

func newModuleLoadRequestKey(request *moduleLoadRequest) moduleLoadRequestKey {
	orgName := ""
	if request.orgName != nil {
		orgName = request.orgName.Value()
	}
	return moduleLoadRequestKey{
		orgName:    orgName,
		moduleName: request.moduleName,
	}
}

// importModuleResponse represents the result of resolving a moduleLoadRequest.
type importModuleResponse struct {
	moduleDesc ModuleDescriptor
	resolved   bool
}

// moduleResolver resolves module dependencies within a single package.
type moduleResolver struct {
	rootPkgDesc PackageDescriptor
	moduleNames map[string]ModuleDescriptor
	responseMap map[moduleLoadRequestKey]*importModuleResponse
}

func newModuleResolver(
	rootPkgDesc PackageDescriptor,
	moduleDescriptors []ModuleDescriptor,
) *moduleResolver {
	moduleNames := make(map[string]ModuleDescriptor, len(moduleDescriptors))
	for _, modDesc := range moduleDescriptors {
		moduleNames[modDesc.Name().String()] = modDesc
	}

	return &moduleResolver{
		rootPkgDesc: rootPkgDesc,
		moduleNames: moduleNames,
		responseMap: make(map[moduleLoadRequestKey]*importModuleResponse),
	}
}

func (r *moduleResolver) resolveModuleLoadRequests(requests []*moduleLoadRequest) []*importModuleResponse {
	responses := make([]*importModuleResponse, 0, len(requests))

	for _, request := range requests {
		key := newModuleLoadRequestKey(request)

		// Check cache first
		if cached, ok := r.responseMap[key]; ok {
			responses = append(responses, cached)
			continue
		}

		// Try to resolve the module
		response := r.resolveRequest(request)
		r.responseMap[key] = response
		responses = append(responses, response)
	}

	return responses
}

func (r *moduleResolver) resolveRequest(request *moduleLoadRequest) *importModuleResponse {
	// Check if this is a root package module
	if r.isRootPackageModule(request.orgName, request.moduleName) {
		modDesc, ok := r.moduleNames[request.moduleName]
		if ok {
			return &importModuleResponse{
				moduleDesc: modDesc,
				resolved:   true,
			}
		}
	}

	// Module not found in root package
	return &importModuleResponse{
		resolved: false,
	}
}

// isRootPackageModule checks if the given org and module name belong to the root package.
// A module belongs to the root package if:
// 1. The org matches the root package org (or org is nil)
// 2. The module name is found in the moduleNames map
func (r *moduleResolver) isRootPackageModule(orgName *PackageOrg, moduleName string) bool {
	if orgName != nil && !orgName.Equals(r.rootPkgDesc.Org()) {
		return false
	}
	_, exists := r.moduleNames[moduleName]
	return exists
}
