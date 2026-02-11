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

// Package internal provides internal implementation details for the projects package.
package internal

import (
	"fmt"
	"path/filepath"
	"slices"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/tomlparser"
	"ballerina-lang-go/tools/diagnostics"
)

// TOML key constants for Ballerina.toml parsing.
const (
	keyPackage     = "package"
	keyOrg         = "org"
	keyName        = "name"
	keyVersion     = "version"
	keyLicense     = "license"
	keyAuthors     = "authors"
	keyKeywords    = "keywords"
	keyRepository  = "repository"
	keyDescription = "description"
	keyVisibility  = "visibility"
	keyTemplate    = "template"

	keyDependency = "dependency"

	keyBuildOptions          = "build-options"
	keyOffline               = "offline"
	keyObservabilityIncluded = "observabilityIncluded"
	keySkipTests             = "skipTests"
	keyTestReport            = "testReport"
	keyCodeCoverage          = "codeCoverage"
	keyCloud                 = "cloud"
	keyGraalVM               = "graalvm"
	keyGraalVMBuildOptions   = "graalvmBuildOptions"
	keySticky                = "sticky"
)

// ManifestBuilder parses Ballerina.toml and produces a PackageManifest and BuildOptions.
// It uses the tomlparser package to read TOML content and constructs the manifest.
// Diagnostics are accumulated for any parsing errors or warnings.
//
// ManifestBuilder combines TOML parsing with chainable With* methods for flexibility.
// It can be used either by parsing TOML or by directly setting values via With* methods.
//
// Java source: io.ballerina.projects.internal.ManifestBuilder
type ManifestBuilder struct {
	toml        *tomlparser.Toml
	projectPath string
	diagnostics []diagnostics.Diagnostic

	// Builder state for With* methods
	packageDesc      projects.PackageDescriptor
	dependencies     []projects.Dependency
	buildOptions     projects.BuildOptions
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
	platforms        map[string]*projects.Platform
	otherEntries     map[string]any
}

// NewManifestBuilder creates a builder from a parsed TOML document.
func NewManifestBuilder(toml *tomlparser.Toml, projectPath string) *ManifestBuilder {
	return &ManifestBuilder{
		toml:         toml,
		projectPath:  projectPath,
		diagnostics:  make([]diagnostics.Diagnostic, 0),
		dependencies: make([]projects.Dependency, 0),
		buildOptions: projects.NewBuildOptions(),
		license:      []string{},
		authors:      []string{},
		keywords:     []string{},
		platforms:    make(map[string]*projects.Platform),
		otherEntries: make(map[string]any),
	}
}

// NewManifestBuilderFromDescriptor creates a builder with a package descriptor (no TOML parsing).
// Use this when you want to build a manifest programmatically via With* methods.
func NewManifestBuilderFromDescriptor(desc projects.PackageDescriptor) *ManifestBuilder {
	return &ManifestBuilder{
		packageDesc:  desc,
		diagnostics:  make([]diagnostics.Diagnostic, 0),
		dependencies: make([]projects.Dependency, 0),
		buildOptions: projects.NewBuildOptions(),
		license:      []string{},
		authors:      []string{},
		keywords:     []string{},
		platforms:    make(map[string]*projects.Platform),
		otherEntries: make(map[string]any),
	}
}

// WithDependencies sets the package dependencies.
func (b *ManifestBuilder) WithDependencies(deps []projects.Dependency) *ManifestBuilder {
	b.dependencies = make([]projects.Dependency, len(deps))
	copy(b.dependencies, deps)
	return b
}

// WithBuildOptions sets the build options.
func (b *ManifestBuilder) WithBuildOptions(opts projects.BuildOptions) *ManifestBuilder {
	b.buildOptions = opts
	return b
}

// WithDiagnostics sets the diagnostics from manifest parsing.
func (b *ManifestBuilder) WithDiagnostics(diags []diagnostics.Diagnostic) *ManifestBuilder {
	b.diagnostics = make([]diagnostics.Diagnostic, len(diags))
	copy(b.diagnostics, diags)
	return b
}

// WithLicense sets the license information.
func (b *ManifestBuilder) WithLicense(license []string) *ManifestBuilder {
	b.license = make([]string, len(license))
	copy(b.license, license)
	return b
}

// WithAuthors sets the package authors.
func (b *ManifestBuilder) WithAuthors(authors []string) *ManifestBuilder {
	b.authors = make([]string, len(authors))
	copy(b.authors, authors)
	return b
}

// WithKeywords sets the package keywords.
func (b *ManifestBuilder) WithKeywords(keywords []string) *ManifestBuilder {
	b.keywords = make([]string, len(keywords))
	copy(b.keywords, keywords)
	return b
}

// WithRepository sets the package repository URL.
func (b *ManifestBuilder) WithRepository(repo string) *ManifestBuilder {
	b.repository = repo
	return b
}

// WithBallerinaVersion sets the required Ballerina version.
func (b *ManifestBuilder) WithBallerinaVersion(version string) *ManifestBuilder {
	b.ballerinaVersion = version
	return b
}

// WithVisibility sets the package visibility.
func (b *ManifestBuilder) WithVisibility(visibility string) *ManifestBuilder {
	b.visibility = visibility
	return b
}

// WithIcon sets the package icon path.
func (b *ManifestBuilder) WithIcon(icon string) *ManifestBuilder {
	b.icon = icon
	return b
}

// WithReadme sets the package readme path.
func (b *ManifestBuilder) WithReadme(readme string) *ManifestBuilder {
	b.readme = readme
	return b
}

// WithDescription sets the package description.
func (b *ManifestBuilder) WithDescription(desc string) *ManifestBuilder {
	b.description = desc
	return b
}

// WithTemplate sets whether this is a template package.
func (b *ManifestBuilder) WithTemplate(template bool) *ManifestBuilder {
	b.template = template
	return b
}

// WithPlatforms sets the platform-specific configurations.
func (b *ManifestBuilder) WithPlatforms(platforms map[string]*projects.Platform) *ManifestBuilder {
	b.platforms = make(map[string]*projects.Platform, len(platforms))
	for k, v := range platforms {
		if v != nil {
			b.platforms[k] = projects.NewPlatform(v.Dependencies())
			if v.IsGraalVMCompatibleSet() {
				b.platforms[k] = projects.NewPlatformWithGraalVM(v.Dependencies(), v.GraalVMCompatible())
			}
		}
	}
	return b
}

// Build constructs the PackageManifest.
// If TOML was provided, it parses TOML values first, then applies any With* overrides.
// All errors are captured as diagnostics - use Diagnostics() to check for errors.
func (b *ManifestBuilder) Build() projects.PackageManifest {
	// If we have TOML, parse it first to populate fields
	if b.toml != nil {
		b.parseFromTOML()
	}

	// Build the manifest using PackageManifestParams
	params := projects.PackageManifestParams{
		PackageDesc:      b.packageDesc,
		Dependencies:     b.dependencies,
		BuildOptions:     b.buildOptions,
		Diagnostics:      b.diagnostics,
		License:          b.license,
		Authors:          b.authors,
		Keywords:         b.keywords,
		Repository:       b.repository,
		BallerinaVersion: b.ballerinaVersion,
		Visibility:       b.visibility,
		Icon:             b.icon,
		Readme:           b.readme,
		Description:      b.description,
		Template:         b.template,
		Platforms:        b.platforms,
		OtherEntries:     b.otherEntries,
	}

	return projects.NewPackageManifestFromParams(params)
}

// parseFromTOML parses all values from the TOML document into builder fields.
// All errors are captured as diagnostics.
func (b *ManifestBuilder) parseFromTOML() {
	// Parse package descriptor with defaults
	b.packageDesc = b.parsePackageDescriptor()

	// Parse dependencies
	b.dependencies = b.parseDependencies()

	// Parse build options
	b.buildOptions = b.parseBuildOptions()

	// Parse metadata fields
	b.license = b.parseStringArray(keyPackage + "." + keyLicense)
	b.authors = b.parseStringArray(keyPackage + "." + keyAuthors)
	b.keywords = b.parseStringArray(keyPackage + "." + keyKeywords)
	b.repository = b.parseString(keyPackage + "." + keyRepository)
	b.description = b.parseString(keyPackage + "." + keyDescription)
	b.visibility = b.parseString(keyPackage + "." + keyVisibility)
	if template, ok := b.toml.GetBool(keyPackage + "." + keyTemplate); ok {
		b.template = template
	}
}

// Diagnostics returns accumulated diagnostics.
func (b *ManifestBuilder) Diagnostics() []diagnostics.Diagnostic {
	if b.diagnostics == nil {
		return []diagnostics.Diagnostic{}
	}
	return slices.Clone(b.diagnostics)
}

// parsePackageDescriptor parses the [package] section and returns a PackageDescriptor.
// Applies default values for missing fields:
//   - org: projects.DefaultOrg ("$anon")
//   - name: directory name of projectPath
//   - version: projects.DefaultVersion ("0.0.0")
//
// All errors are captured as diagnostics.
func (b *ManifestBuilder) parsePackageDescriptor() projects.PackageDescriptor {
	// Get org with default
	org := b.parseString(keyPackage + "." + keyOrg)
	if org == "" {
		org = projects.DefaultOrg
	}

	// Get name with default (directory name)
	name := b.parseString(keyPackage + "." + keyName)
	if name == "" {
		name = filepath.Base(b.projectPath)
	}

	// Get version with default
	versionStr := b.parseString(keyPackage + "." + keyVersion)
	if versionStr == "" {
		versionStr = projects.DefaultVersion
	}

	// Parse and validate version
	version, err := projects.NewPackageVersionFromString(versionStr)
	if err != nil {
		b.addDiagnostic(diagnostics.Error, fmt.Sprintf("invalid version '%s': %v", versionStr, err))
		// Use default version on error
		version = projects.DefaultPackageVersion
	}

	return projects.NewPackageDescriptor(
		projects.NewPackageOrg(org),
		projects.NewPackageName(name),
		version)
}

// parseDependencies parses the [[dependency]] array from the TOML document.
func (b *ManifestBuilder) parseDependencies() []projects.Dependency {
	deps := make([]projects.Dependency, 0)

	tables, ok := b.toml.GetTables(keyDependency)
	if !ok {
		return deps
	}

	for _, table := range tables {
		dep, err := b.parseDependency(table)
		if err != nil {
			b.addDiagnostic(diagnostics.Error, fmt.Sprintf("invalid dependency: %v", err))
			continue
		}
		deps = append(deps, dep)
	}

	return deps
}

// parseDependency parses a single dependency table.
func (b *ManifestBuilder) parseDependency(table *tomlparser.Toml) (projects.Dependency, error) {
	org, ok := table.GetString(keyOrg)
	if !ok || org == "" {
		return projects.Dependency{}, fmt.Errorf("missing required field 'org'")
	}

	name, ok := table.GetString(keyName)
	if !ok || name == "" {
		return projects.Dependency{}, fmt.Errorf("missing required field 'name'")
	}

	versionStr, ok := table.GetString(keyVersion)
	if !ok || versionStr == "" {
		return projects.Dependency{}, fmt.Errorf("missing required field 'version'")
	}

	version, err := projects.NewPackageVersionFromString(versionStr)
	if err != nil {
		return projects.Dependency{}, fmt.Errorf("invalid version '%s': %w", versionStr, err)
	}

	// Repository is optional
	repository, _ := table.GetString(keyRepository)

	if repository != "" {
		return projects.NewDependencyWithRepository(
			projects.NewPackageOrg(org),
			projects.NewPackageName(name),
			version,
			repository,
		), nil
	}

	return projects.NewDependency(
		projects.NewPackageOrg(org),
		projects.NewPackageName(name),
		version,
	), nil
}

// parseBuildOptions parses the [build-options] section from the TOML document.
// Uses projects.BuildOptionsBuilder to construct the result.
func (b *ManifestBuilder) parseBuildOptions() projects.BuildOptions {
	builder := projects.NewBuildOptionsBuilder()

	// Check if build-options section exists
	_, ok := b.toml.GetTable(keyBuildOptions)
	if !ok {
		return builder.Build()
	}

	// Parse each build option if present
	if offline, ok := b.toml.GetBool(keyBuildOptions + "." + keyOffline); ok {
		builder.WithOffline(offline)
	}

	if observability, ok := b.toml.GetBool(keyBuildOptions + "." + keyObservabilityIncluded); ok {
		builder.WithObservabilityIncluded(observability)
	}

	if skipTests, ok := b.toml.GetBool(keyBuildOptions + "." + keySkipTests); ok {
		builder.WithSkipTests(skipTests)
	}

	if testReport, ok := b.toml.GetBool(keyBuildOptions + "." + keyTestReport); ok {
		builder.WithTestReport(testReport)
	}

	if codeCoverage, ok := b.toml.GetBool(keyBuildOptions + "." + keyCodeCoverage); ok {
		builder.WithCodeCoverage(codeCoverage)
	}

	if cloud, ok := b.toml.GetString(keyBuildOptions + "." + keyCloud); ok {
		builder.WithCloud(cloud)
	}

	if graalVM, ok := b.toml.GetBool(keyBuildOptions + "." + keyGraalVM); ok {
		builder.WithGraalVM(graalVM)
	}

	if graalVMOpts, ok := b.toml.GetString(keyBuildOptions + "." + keyGraalVMBuildOptions); ok {
		builder.WithGraalVMBuildOptions(graalVMOpts)
	}

	if sticky, ok := b.toml.GetBool(keyBuildOptions + "." + keySticky); ok {
		builder.WithSticky(sticky)
	}

	return builder.Build()
}

// parseString retrieves a string value from the TOML document.
// Returns empty string if the key does not exist.
func (b *ManifestBuilder) parseString(key string) string {
	value, ok := b.toml.GetString(key)
	if !ok {
		return ""
	}
	return value
}

// parseStringArray retrieves a string array from the TOML document.
// Returns an empty slice if the key does not exist or is not an array.
func (b *ManifestBuilder) parseStringArray(key string) []string {
	arr, ok := b.toml.GetArray(key)
	if !ok {
		return []string{}
	}

	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// addDiagnostic adds a diagnostic message to the builder.
func (b *ManifestBuilder) addDiagnostic(severity diagnostics.DiagnosticSeverity, message string) {
	info := diagnostics.NewDiagnosticInfo(nil, message, severity)
	loc := diagnostics.NewBLangDiagnosticLocation(
		filepath.Join(b.projectPath, projects.BallerinaTomlFile),
		0, 0, 0, 0, 0, 0,
	)
	diag := diagnostics.NewDefaultDiagnostic(info, loc, nil)
	b.diagnostics = append(b.diagnostics, diag)
}
