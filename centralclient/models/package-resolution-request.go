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

import "net/url"

type PackageResolutionMode string

const (
	PackageResolutionModeSoft   PackageResolutionMode = "soft"
	PackageResolutionModeMedium PackageResolutionMode = "medium"
	PackageResolutionModeHard   PackageResolutionMode = "hard"
	PackageResolutionModeLocked PackageResolutionMode = "locked"
)

type PackageResolutionRequest struct {
	Packages []packageResolutionRequestPackage `json:"packages"`
}

// PackageResolutionRequest.java:49-95
type packageResolutionRequestPackage struct {
	Org     string                `json:"org"`
	Name    string                `json:"name"`
	Version string                `json:"version"`
	Mode    PackageResolutionMode `json:"mode,omitempty"`
}

func (r *PackageResolutionRequest) AddPackage(orgName, name, version string, mode PackageResolutionMode) {
	encodedVersion := url.QueryEscape(version)
	r.Packages = append(r.Packages, packageResolutionRequestPackage{
		Org:     orgName,
		Name:    name,
		Version: encodedVersion,
		Mode:    mode,
	})
}
