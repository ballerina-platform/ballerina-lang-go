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

// PackageResolver resolves packages from available repositories.
type PackageResolver interface {
	// ResolvePackages loads the packages specified in the requests.
	ResolvePackages(ctx context.Context, requests []ResolutionRequest, options ResolutionOptions) []ResolutionResponse

	// ResolveByName resolves packages by org and name, without requiring a specific version.
	// Returns all matching packages from cache or repositories.
	ResolveByName(ctx context.Context, org, name string, options ResolutionOptions) []*Package

	// AddRepository adds a repository to the resolver.
	// Repositories are searched in the order they are added.
	AddRepository(repo Repository)

	// Repositories returns the list of repositories in this resolver.
	Repositories() []Repository
}

// defaultPackageResolver is the default implementation of PackageResolver.
type defaultPackageResolver struct {
	cache        *PackageCache
	repositories []Repository // Uses the Repository interface from repository.go
}

// AddRepository adds a repository to the resolver.
// Repositories are searched in the order they are added.
func (r *defaultPackageResolver) AddRepository(repo Repository) {
	r.repositories = append(r.repositories, repo)
}

// Repositories returns the list of repositories in this resolver.
func (r *defaultPackageResolver) Repositories() []Repository {
	return r.repositories
}

// NewPackageResolver creates a new PackageResolver with the given cache and repositories.
func NewPackageResolver(cache *PackageCache, repositories ...Repository) PackageResolver {
	return &defaultPackageResolver{
		cache:        cache,
		repositories: repositories,
	}
}

// ResolvePackages resolves packages by checking cache first, then repositories.
func (r *defaultPackageResolver) ResolvePackages(
	ctx context.Context,
	requests []ResolutionRequest,
	options ResolutionOptions,
) []ResolutionResponse {
	responses := make([]ResolutionResponse, 0, len(requests))

	for _, request := range requests {
		response := r.resolvePackage(ctx, request, options)
		responses = append(responses, response)
	}

	return responses
}

func (r *defaultPackageResolver) resolvePackage(
	ctx context.Context,
	request ResolutionRequest,
	options ResolutionOptions,
) ResolutionResponse {
	desc := request.Descriptor()
	org := desc.Org().Value()
	name := desc.Name().Value()
	version := desc.Version().String()

	// 1. Check cache first (unless disabled)
	if !options.DisableCache() {
		pkg := r.cache.Get(org, name, version)
		if pkg != nil {
			return NewResolvedResponse(pkg, request)
		}
	}

	// 2. Try repositories in order
	for _, repo := range r.repositories {
		select {
		case <-ctx.Done():
			return NewUnresolvedResponse(request)
		default:
		}

		pkg, err := repo.GetPackage(ctx, org, name, version)
		if err != nil {
			continue // Try next repository
		}
		if pkg == nil {
			continue
		}

		// 3. Cache the loaded package
		if !options.DisableCache() {
			r.cache.Cache(pkg)
		}

		return NewResolvedResponse(pkg, request)
	}

	// Not found in any repository
	return NewUnresolvedResponse(request)
}

// ResolveByName resolves packages by org and name without requiring a specific version.
// It first checks the cache, then tries repositories to find all available versions.
func (r *defaultPackageResolver) ResolveByName(
	ctx context.Context,
	org, name string,
	options ResolutionOptions,
) []*Package {
	// 1. Check cache first
	if !options.DisableCache() {
		packages := r.cache.GetPackages(org, name)
		if len(packages) > 0 {
			return packages
		}
	}

	// 2. Try repositories in order to find available versions
	var result []*Package
	for _, repo := range r.repositories {
		select {
		case <-ctx.Done():
			return result
		default:
		}

		// Get all available versions from this repository
		versions, err := repo.GetPackageVersions(ctx, org, name)
		if err != nil {
			continue // Try next repository
		}

		// Load each version
		for _, version := range versions {
			select {
			case <-ctx.Done():
				return result
			default:
			}

			pkg, err := repo.GetPackage(ctx, org, name, version.String())
			if err != nil || pkg == nil {
				continue
			}

			// Cache the loaded package
			if !options.DisableCache() {
				r.cache.Cache(pkg)
			}

			result = append(result, pkg)
		}

		// If we found packages in this repository, stop searching
		if len(result) > 0 {
			break
		}
	}

	return result
}
