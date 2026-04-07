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
	"io/fs"

	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semantics"
)

// Environment represents an environment shared by a set of projects.
// It maintains a global package cache for all loaded packages (internal and external).
type Environment struct {
	fsys              fs.FS
	compilerEnv       *context.CompilerEnvironment
	packageCache      *PackageCache
	packageResolver   PackageResolver
	resolutionOptions ResolutionOptions
	// TODO: find better place to put this
	publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace
}

// NewEnvironment creates a new Environment.
func NewEnvironment(fsys fs.FS, env *context.CompilerEnvironment) *Environment {
	cache := newPackageCache()
	return &Environment{
		fsys:            fsys,
		compilerEnv:     env,
		packageCache:    cache,
		packageResolver: NewPackageResolver(cache),
		publicSymbols:   make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace),
	}
}

func newEnvironment(fsys fs.FS, env *context.CompilerEnvironment) *Environment {
	return NewEnvironment(fsys, env)
}

func (e *Environment) compilerEnvironment() *context.CompilerEnvironment {
	return e.compilerEnv
}

func (e *Environment) fs() fs.FS {
	return e.fsys
}

// PackageCache returns the environment's package cache.
func (e *Environment) PackageCache() *PackageCache {
	return e.packageCache
}

// PackageResolver returns the environment's package resolver.
func (e *Environment) PackageResolver() PackageResolver {
	return e.packageResolver
}

// ResolutionOptions returns the environment's resolution options.
func (e *Environment) ResolutionOptions() ResolutionOptions {
	return e.resolutionOptions
}

// addPublicSymbolsFrom copies public symbols from another environment.
// This is used internally when integrating external packages.
func (e *Environment) addPublicSymbolsFrom(other *Environment) {
	for k, v := range other.publicSymbols {
		e.publicSymbols[k] = v
	}
}

// addRepository adds a repository to the environment's package resolver.
// Repositories are searched in the order they are added.
// This is private to enforce Environment immutability - repositories should be
// configured during Environment creation via RepositoryFactories.
func (e *Environment) addRepository(repo Repository) {
	e.packageResolver.AddRepository(repo)
}

// Duplicate creates a new Environment with fresh caches but the same repository
// configuration and resolution options.
func (e *Environment) Duplicate() *Environment {
	newEnv := NewEnvironment(e.fsys, e.compilerEnv)
	newEnv.resolutionOptions = e.resolutionOptions

	// Copy repositories from original resolver
	for _, repo := range e.packageResolver.Repositories() {
		newEnv.addRepository(repo)
	}

	return newEnv
}
