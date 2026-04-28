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

// Repository provides access to Ballerina packages from a specific source.
// Implementations include FileSystemRepository (local bala cache),
// CentralRepository (Ballerina Central API), MavenRepository, etc.
type Repository interface {
	// GetPackage loads a specific version of a package.
	// Returns (nil, nil) if not found (not an error).
	// Returns (nil, error) on actual errors (IO, parse, etc.)
	GetPackage(ctx context.Context, org, name, version string) (*Package, error)

	// GetPackageVersions returns all available versions for a package.
	// Returns empty slice if package not found.
	GetPackageVersions(ctx context.Context, org, name string) ([]PackageVersion, error)

	// GetLatestVersion returns the latest available version for a package.
	// Returns (zero, false, nil) if not found.
	GetLatestVersion(ctx context.Context, org, name string) (PackageVersion, bool, error)
}
