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

import (
	"maps"
	"slices"

	"ballerina-lang-go/tools/diagnostics"
)

// PackageManifest represents parsed Ballerina.toml content.
// It is an immutable data structure containing package metadata, dependencies,
// build options, and any diagnostics from parsing.
// Java source: io.ballerina.projects.PackageManifest
type PackageManifest struct {
	packageDesc      PackageDescriptor
	dependencies     []Dependency
	buildOptions     BuildOptions
	diagnostics      []diagnostics.Diagnostic
	license          []string
	authors          []string
	keywords         []string
	repository       string
	ballerinaVersion string
	visibility       string
	icon             string
	readme           string
	description      string
	template         bool
	platforms        map[string]*Platform
}

// Dependency represents a package dependency declared in Ballerina.toml.
// Java source: io.ballerina.projects.PackageManifest.Dependency
type Dependency struct {
	name       PackageName
	org        PackageOrg
	version    PackageVersion
	repository string
}

// Platform represents platform-specific configuration in Ballerina.toml.
// Java source: io.ballerina.projects.PackageManifest.Platform
type Platform struct {
	dependencies      []map[string]any
	graalvmCompatible *bool
}

// NewDependency creates a new Dependency with the given components.
func NewDependency(org PackageOrg, name PackageName, version PackageVersion) Dependency {
	return Dependency{
		org:     org,
		name:    name,
		version: version,
	}
}

// NewDependencyWithRepository creates a new Dependency with a repository URL.
func NewDependencyWithRepository(org PackageOrg, name PackageName, version PackageVersion, repository string) Dependency {
	return Dependency{
		org:        org,
		name:       name,
		version:    version,
		repository: repository,
	}
}

// Org returns the dependency's organization.
func (d Dependency) Org() PackageOrg {
	return d.org
}

// Name returns the dependency's package name.
func (d Dependency) Name() PackageName {
	return d.name
}

// Version returns the dependency's version.
func (d Dependency) Version() PackageVersion {
	return d.version
}

// Repository returns the dependency's repository URL.
func (d Dependency) Repository() string {
	return d.repository
}

// NewPlatform creates a new Platform with the given dependencies.
func NewPlatform(dependencies []map[string]any) *Platform {
	depsCopy := make([]map[string]any, len(dependencies))
	for i, dep := range dependencies {
		depsCopy[i] = maps.Clone(dep)
	}
	return &Platform{
		dependencies: depsCopy,
	}
}

// NewPlatformWithGraalVM creates a new Platform with GraalVM compatibility info.
func NewPlatformWithGraalVM(dependencies []map[string]any, graalvmCompatible bool) *Platform {
	p := NewPlatform(dependencies)
	p.graalvmCompatible = &graalvmCompatible
	return p
}

// Dependencies returns a copy of the platform dependencies.
func (p *Platform) Dependencies() []map[string]any {
	if p == nil || p.dependencies == nil {
		return nil
	}
	result := make([]map[string]any, len(p.dependencies))
	for i, dep := range p.dependencies {
		result[i] = maps.Clone(dep)
	}
	return result
}

// GraalVMCompatible returns whether this platform is GraalVM compatible.
func (p *Platform) GraalVMCompatible() bool {
	if p == nil || p.graalvmCompatible == nil {
		return false
	}
	return *p.graalvmCompatible
}

// IsGraalVMCompatibleSet returns true if GraalVM compatibility was explicitly set.
func (p *Platform) IsGraalVMCompatibleSet() bool {
	return p != nil && p.graalvmCompatible != nil
}

// NewPackageManifest creates a new PackageManifest with default values.
func NewPackageManifest(desc PackageDescriptor) PackageManifest {
	return PackageManifest{
		packageDesc:  desc,
		dependencies: []Dependency{},
		buildOptions: NewBuildOptions(),
		diagnostics:  []diagnostics.Diagnostic{},
		license:      []string{},
		authors:      []string{},
		keywords:     []string{},
		platforms:    make(map[string]*Platform),
	}
}

// PackageDescriptor returns the package descriptor (org/name/version).
func (m PackageManifest) PackageDescriptor() PackageDescriptor {
	return m.packageDesc
}

// Org returns the package organization.
func (m PackageManifest) Org() PackageOrg {
	return m.packageDesc.Org()
}

// Name returns the package name.
func (m PackageManifest) Name() PackageName {
	return m.packageDesc.Name()
}

// Version returns the package version.
func (m PackageManifest) Version() PackageVersion {
	return m.packageDesc.Version()
}

// Dependencies returns a copy of the package dependencies.
func (m PackageManifest) Dependencies() []Dependency {
	if m.dependencies == nil {
		return []Dependency{}
	}
	return slices.Clone(m.dependencies)
}

// BuildOptions returns the build options.
func (m PackageManifest) BuildOptions() BuildOptions {
	return m.buildOptions
}

// Diagnostics returns a copy of the parsing diagnostics.
func (m PackageManifest) Diagnostics() []diagnostics.Diagnostic {
	if m.diagnostics == nil {
		return []diagnostics.Diagnostic{}
	}
	return slices.Clone(m.diagnostics)
}

// HasDiagnostics returns true if there are any diagnostics.
func (m PackageManifest) HasDiagnostics() bool {
	return len(m.diagnostics) > 0
}

// License returns a copy of the license information.
func (m PackageManifest) License() []string {
	if m.license == nil {
		return []string{}
	}
	return slices.Clone(m.license)
}

// Authors returns a copy of the package authors.
func (m PackageManifest) Authors() []string {
	if m.authors == nil {
		return []string{}
	}
	return slices.Clone(m.authors)
}

// Keywords returns a copy of the package keywords.
func (m PackageManifest) Keywords() []string {
	if m.keywords == nil {
		return []string{}
	}
	return slices.Clone(m.keywords)
}

// Repository returns the package repository URL.
func (m PackageManifest) Repository() string {
	return m.repository
}

// BallerinaVersion returns the required Ballerina version.
func (m PackageManifest) BallerinaVersion() string {
	return m.ballerinaVersion
}

// Visibility returns the package visibility.
func (m PackageManifest) Visibility() string {
	return m.visibility
}

// Icon returns the package icon path.
func (m PackageManifest) Icon() string {
	return m.icon
}

// Readme returns the package readme path.
func (m PackageManifest) Readme() string {
	return m.readme
}

// Description returns the package description.
func (m PackageManifest) Description() string {
	return m.description
}

// Template returns whether this is a template package.
func (m PackageManifest) Template() bool {
	return m.template
}

// Platform returns the platform configuration for the given name.
// Returns nil if no platform configuration exists for that name.
func (m PackageManifest) Platform(name string) *Platform {
	p, ok := m.platforms[name]
	if !ok {
		return nil
	}
	result := NewPlatform(p.dependencies)
	if p.graalvmCompatible != nil {
		graalvm := *p.graalvmCompatible
		result.graalvmCompatible = &graalvm
	}
	return result
}

// Platforms returns a copy of all platform configurations.
func (m PackageManifest) Platforms() map[string]*Platform {
	result := make(map[string]*Platform, len(m.platforms))
	for k, v := range m.platforms {
		if v != nil {
			result[k] = NewPlatform(v.dependencies)
			if v.graalvmCompatible != nil {
				graalvm := *v.graalvmCompatible
				result[k].graalvmCompatible = &graalvm
			}
		}
	}
	return result
}

// PackageManifestParams contains all parameters needed to construct a PackageManifest.
// This struct is used by internal packages that need to build PackageManifest instances.
// All fields are exported to allow cross-package construction.
type PackageManifestParams struct {
	PackageDesc      PackageDescriptor
	Dependencies     []Dependency
	BuildOptions     BuildOptions
	Diagnostics      []diagnostics.Diagnostic
	License          []string
	Authors          []string
	Keywords         []string
	Repository       string
	BallerinaVersion string
	Visibility       string
	Icon             string
	Readme           string
	Description      string
	Template         bool
	Platforms        map[string]*Platform
	OtherEntries     map[string]any
}

// NewPackageManifestFromParams creates a PackageManifest from the given parameters.
// This function is intended for use by internal packages that need to construct
// PackageManifest instances with full control over all fields.
func NewPackageManifestFromParams(params PackageManifestParams) PackageManifest {
	// Copy dependencies
	deps := make([]Dependency, len(params.Dependencies))
	copy(deps, params.Dependencies)

	// Copy diagnostics
	diags := make([]diagnostics.Diagnostic, len(params.Diagnostics))
	copy(diags, params.Diagnostics)

	// Copy license
	license := make([]string, len(params.License))
	copy(license, params.License)

	// Copy authors
	authors := make([]string, len(params.Authors))
	copy(authors, params.Authors)

	// Copy keywords
	keywords := make([]string, len(params.Keywords))
	copy(keywords, params.Keywords)

	// Copy platforms
	platforms := make(map[string]*Platform, len(params.Platforms))
	for k, v := range params.Platforms {
		if v != nil {
			platforms[k] = NewPlatform(v.dependencies)
			if v.graalvmCompatible != nil {
				graalvm := *v.graalvmCompatible
				platforms[k].graalvmCompatible = &graalvm
			}
		}
	}

	return PackageManifest{
		packageDesc:      params.PackageDesc,
		dependencies:     deps,
		buildOptions:     params.BuildOptions,
		diagnostics:      diags,
		license:          license,
		authors:          authors,
		keywords:         keywords,
		repository:       params.Repository,
		ballerinaVersion: params.BallerinaVersion,
		visibility:       params.Visibility,
		icon:             params.Icon,
		readme:           params.Readme,
		description:      params.Description,
		template:         params.Template,
		platforms:        platforms,
	}
}
