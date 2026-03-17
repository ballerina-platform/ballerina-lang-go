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

package projects

import (
	"context"
)

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

// ResolutionStatus indicates whether a module/package was resolved.
// Java source: io.ballerina.projects.environment.ResolutionResponse.ResolutionStatus
type ResolutionStatus int

const (
	// ResolutionStatusResolved indicates the module/package was found.
	ResolutionStatusResolved ResolutionStatus = iota
	// ResolutionStatusUnresolved indicates the module/package was not found.
	ResolutionStatusUnresolved
)

// ImportModuleResponse represents the result of resolving a moduleLoadRequest.
// Java source: io.ballerina.projects.internal.ImportModuleResponse
type ImportModuleResponse struct {
	packageDescriptor *PackageDescriptor // Package containing the module (nil for same package)
	moduleDesc        ModuleDescriptor   // Module descriptor
	resolutionStatus  ResolutionStatus
}

// moduleResolver resolves module dependencies using PackageResolver.
// Java source: io.ballerina.projects.internal.ModuleResolver
type moduleResolver struct {
	rootPkgDesc     PackageDescriptor
	responseMap     map[moduleLoadRequestKey]*ImportModuleResponse
	packageResolver PackageResolver
}

func newModuleResolver(rootPkgDesc PackageDescriptor, env *Environment) *moduleResolver {
	return &moduleResolver{
		rootPkgDesc:     rootPkgDesc,
		responseMap:     make(map[moduleLoadRequestKey]*ImportModuleResponse),
		packageResolver: env.PackageResolver(),
	}
}

func (r *moduleResolver) resolveModuleLoadRequests(ctx context.Context, requests []*moduleLoadRequest) []*ImportModuleResponse {
	responses := make([]*ImportModuleResponse, 0, len(requests))

	for _, request := range requests {
		key := newModuleLoadRequestKey(request)

		// Check cache first
		if cached, ok := r.responseMap[key]; ok {
			responses = append(responses, cached)
			continue
		}

		// Try to resolve the module
		response := r.resolveRequest(ctx, request)
		r.responseMap[key] = response
		responses = append(responses, response)
	}

	return responses
}

func (r *moduleResolver) resolveRequest(ctx context.Context, request *moduleLoadRequest) *ImportModuleResponse {
	// Determine org - use request org or default to root package org
	org := r.rootPkgDesc.Org().Value()
	if request.orgName != nil {
		org = request.orgName.Value()
	}

	// Extract package name from module name (e.g., "http.auth" -> "http")
	pkgName := extractPackageName(request.moduleName)

	// Look up packages via PackageResolver
	packages := r.packageResolver.ResolveByName(ctx, org, pkgName, NewResolutionOptions())
	for _, pkg := range packages {
		// Check if module exists in this package
		for _, mod := range pkg.Modules() {
			if mod.ModuleName().String() == request.moduleName {
				pkgDesc := pkg.Manifest().PackageDescriptor()
				return &ImportModuleResponse{
					packageDescriptor: &pkgDesc,
					moduleDesc:        mod.Descriptor(),
					resolutionStatus:  ResolutionStatusResolved,
				}
			}
		}
	}

	// Module not found
	return &ImportModuleResponse{
		resolutionStatus: ResolutionStatusUnresolved,
	}
}

// extractPackageName extracts the package name from a module name.
// e.g., "http" -> "http", "http.auth" -> "http"
func extractPackageName(moduleName string) string {
	for i, c := range moduleName {
		if c == '.' {
			return moduleName[:i]
		}
	}
	return moduleName
}
