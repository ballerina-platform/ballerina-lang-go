// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package models

type PackageNameResolutionRequest struct {
	Modules []packageNameResolutionRequestModule `json:"modules"`
}

// PackageNameResolutionRequest.java:36-93
type packageNameResolutionRequestModule struct {
	Organization     string                                        `json:"organization"`
	ModuleName       string                                        `json:"moduleName"`
	PossiblePackages []packageNameResolutionRequestPossiblePackage `json:"possiblePackages,omitempty"`
	Mode             PackageResolutionMode                         `json:"mode,omitempty"`
}

// PackageNameResolutionRequest.java:41-63
type packageNameResolutionRequestPossiblePackage struct {
	Org     string `json:"org"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (r *PackageNameResolutionRequest) AddModule(orgName, moduleName string, possiblePackages []packageNameResolutionRequestPossiblePackage, mode PackageResolutionMode) {
	r.Modules = append(r.Modules, packageNameResolutionRequestModule{
		Organization:     orgName,
		ModuleName:       moduleName,
		PossiblePackages: possiblePackages,
		Mode:             mode,
	})
}

func (r *PackageNameResolutionRequest) AddModuleSimple(orgName, moduleName string) {
	r.Modules = append(r.Modules, packageNameResolutionRequestModule{
		Organization: orgName,
		ModuleName:   moduleName,
	})
}
