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

// Package repository provides abstractions for accessing Ballerina package repositories.
//
// A Repository represents a source of Ballerina packages, such as a local cache
// or a remote registry like Ballerina Central. The Repository interface defines
// the common operations that all repository implementations must support.
//
// Currently implemented:
//   - LocalCacheRepository: Scans the local package cache (~/.ballerina/repositories/)
//
// Planned implementations:
//   - CentralRepository: Queries Ballerina Central API
//   - DistributionRepository: Provides langlib and built-in packages
//   - CustomRepository: Supports per-dependency custom registry URLs
package repository

import "context"

// Repository provides read access to Ballerina package versions.
//
// All methods accept a context.Context for cancellation support.
// Implementations should check for context cancellation before
// performing expensive operations.
//
// For "not found" conditions (package or version doesn't exist),
// implementations should return empty values (empty slice, empty string,
// false) rather than errors. Errors are reserved for actual failures
// like I/O errors or context cancellation.
type Repository interface {
	// Name returns the repository identifier for logging and debugging.
	// Examples: "local-cache", "central", "distribution"
	Name() string

	// GetVersions returns all available versions for a package.
	//
	// Versions are sorted in semver order (oldest first, latest last).
	// Returns an empty slice if the package is not found in this repository.
	// Returns an error only for actual failures (I/O error, context cancelled).
	GetVersions(ctx context.Context, org, name string) ([]string, error)

	// GetLatestVersion returns the latest (highest semver) version for a package.
	//
	// Returns an empty string if the package is not found in this repository.
	// Returns an error only for actual failures (I/O error, context cancelled).
	GetLatestVersion(ctx context.Context, org, name string) (string, error)

	// Exists checks if a specific package version exists in this repository.
	//
	// Returns false if the package or version is not found.
	// Returns an error only for actual failures (I/O error, context cancelled).
	Exists(ctx context.Context, org, name, version string) (bool, error)
}
