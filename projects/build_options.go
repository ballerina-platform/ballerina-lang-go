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

// BuildOptions represents build options for a Ballerina project.
// It contains a CompilationOptions instance and delegates compilation-related methods to it.
// Java source: io.ballerina.projects.BuildOptions
type BuildOptions struct {
	// Build-specific fields
	testReport                  *bool
	codeCoverage                *bool
	dumpBuildTime               *bool
	skipTests                   *bool
	targetDir                   string
	nativeImage                 *bool
	exportComponentModel        *bool
	graalVMBuildOptions         string
	showDependencyDiagnostics   *bool

	// Composition: BuildOptions contains CompilationOptions
	compilationOptions CompilationOptions
}

// NewBuildOptions creates a new BuildOptions with default values.
func NewBuildOptions() BuildOptions {
	return BuildOptions{
		compilationOptions: NewCompilationOptions(),
	}
}

// TestReport returns whether test report generation is enabled.
func (b BuildOptions) TestReport() bool {
	return toBoolDefaultIfNull(b.testReport)
}

// CodeCoverage returns whether code coverage is enabled.
func (b BuildOptions) CodeCoverage() bool {
	return toBoolDefaultIfNull(b.codeCoverage)
}

// DumpBuildTime returns whether build time dumping is enabled.
func (b BuildOptions) DumpBuildTime() bool {
	return toBoolDefaultIfNull(b.dumpBuildTime)
}

// SkipTests returns whether tests should be skipped.
// By default, tests are skipped (returns true if nil).
func (b BuildOptions) SkipTests() bool {
	if b.skipTests == nil {
		return true // Default: skip tests
	}
	return *b.skipTests
}

// TargetDir returns the target directory path.
func (b BuildOptions) TargetDir() string {
	return b.targetDir
}

// NativeImage returns whether GraalVM native image is enabled.
func (b BuildOptions) NativeImage() bool {
	return toBoolDefaultIfNull(b.nativeImage)
}

// GraalVM returns whether GraalVM native image is enabled.
// Alias for NativeImage().
func (b BuildOptions) GraalVM() bool {
	return b.NativeImage()
}

// GraalVMBuildOptions returns additional GraalVM build options.
func (b BuildOptions) GraalVMBuildOptions() string {
	return b.graalVMBuildOptions
}

// ShowDependencyDiagnostics returns whether dependency diagnostics should be shown.
func (b BuildOptions) ShowDependencyDiagnostics() bool {
	return toBoolDefaultIfNull(b.showDependencyDiagnostics)
}

// CompilationOptions returns the underlying compilation options.
// Java: BuildOptions.compilationOptions()
func (b BuildOptions) CompilationOptions() CompilationOptions {
	return b.compilationOptions
}

// --- Delegated methods to CompilationOptions ---

// OfflineBuild returns whether offline build mode is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) OfflineBuild() bool {
	return b.compilationOptions.OfflineBuild()
}

// Offline is an alias for OfflineBuild.
func (b BuildOptions) Offline() bool {
	return b.OfflineBuild()
}

// Sticky returns whether sticky mode is enabled.
// Deprecated: Use LockingMode() instead.
// Delegated to CompilationOptions.
func (b BuildOptions) Sticky() bool {
	return b.compilationOptions.Sticky()
}

// DisableSyntaxTree returns whether syntax tree caching is disabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DisableSyntaxTree() bool {
	return b.compilationOptions.DisableSyntaxTree()
}

// OptimizeDependencyCompilation returns whether dependency compilation optimization is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) OptimizeDependencyCompilation() bool {
	return b.compilationOptions.OptimizeDependencyCompilation()
}

// LockingMode returns the package locking mode.
// Returns PackageLockingModeMedium if not explicitly set.
// Delegated to CompilationOptions.
func (b BuildOptions) LockingMode() PackageLockingMode {
	mode := b.compilationOptions.LockingMode()
	if mode == PackageLockingModeUnknown {
		return PackageLockingModeMedium
	}
	return mode
}

// RawLockingMode returns the raw package locking mode.
// Returns PackageLockingModeUnknown if not explicitly set.
// Delegated to CompilationOptions.
func (b BuildOptions) RawLockingMode() PackageLockingMode {
	return b.compilationOptions.LockingMode()
}

// Experimental returns whether experimental features are enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) Experimental() bool {
	return b.compilationOptions.Experimental()
}

// ObservabilityIncluded returns whether observability is included.
// Delegated to CompilationOptions.
func (b BuildOptions) ObservabilityIncluded() bool {
	return b.compilationOptions.ObservabilityIncluded()
}

// ListConflictedClasses returns whether conflicted classes should be listed.
// Delegated to CompilationOptions.
func (b BuildOptions) ListConflictedClasses() bool {
	return b.compilationOptions.ListConflictedClasses()
}

// Cloud returns the cloud target.
// Delegated to CompilationOptions.
func (b BuildOptions) Cloud() string {
	return b.compilationOptions.Cloud()
}

// RemoteManagement returns whether remote management is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) RemoteManagement() bool {
	return b.compilationOptions.RemoteManagement()
}

// ExportOpenAPI returns whether OpenAPI export is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) ExportOpenAPI() bool {
	return b.compilationOptions.ExportOpenAPI()
}

// ExportComponentModel returns whether component model export is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) ExportComponentModel() bool {
	return b.compilationOptions.ExportComponentModel()
}

// DumpAST returns whether AST dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpAST() bool {
	return b.compilationOptions.DumpAST()
}

// DumpBIR returns whether BIR dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpBIR() bool {
	return b.compilationOptions.DumpBIR()
}

// DumpBIRFile returns whether BIR file dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpBIRFile() bool {
	return b.compilationOptions.DumpBIRFile()
}

// DumpGraph returns whether graph dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpGraph() bool {
	return b.compilationOptions.DumpGraph()
}

// DumpRawGraphs returns whether raw graph dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpRawGraphs() bool {
	return b.compilationOptions.DumpRawGraphs()
}

// DumpTokens returns whether lexer token dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpTokens() bool {
	return b.compilationOptions.DumpTokens()
}

// DumpST returns whether syntax tree dumping is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) DumpST() bool {
	return b.compilationOptions.DumpST()
}

// TraceRecovery returns whether error recovery tracing is enabled.
// Delegated to CompilationOptions.
func (b BuildOptions) TraceRecovery() bool {
	return b.compilationOptions.TraceRecovery()
}

// AcceptTheirs merges the given build options by favoring theirs if there are conflicts.
// Java: BuildOptions.acceptTheirs(BuildOptions)
func (b BuildOptions) AcceptTheirs(theirs BuildOptions) BuildOptions {
	builder := NewBuildOptionsBuilder()

	// Handle build-specific fields
	if theirs.skipTests != nil {
		builder.WithSkipTests(*theirs.skipTests)
	} else if b.skipTests != nil {
		builder.WithSkipTests(*b.skipTests)
	}

	if theirs.codeCoverage != nil {
		builder.WithCodeCoverage(*theirs.codeCoverage)
	} else if b.codeCoverage != nil {
		builder.WithCodeCoverage(*b.codeCoverage)
	}

	if theirs.testReport != nil {
		builder.WithTestReport(*theirs.testReport)
	} else if b.testReport != nil {
		builder.WithTestReport(*b.testReport)
	}

	if theirs.dumpBuildTime != nil {
		builder.WithDumpBuildTime(*theirs.dumpBuildTime)
	} else if b.dumpBuildTime != nil {
		builder.WithDumpBuildTime(*b.dumpBuildTime)
	}

	if theirs.targetDir != "" {
		builder.WithTargetDir(theirs.targetDir)
	} else {
		builder.WithTargetDir(b.targetDir)
	}

	if theirs.nativeImage != nil {
		builder.WithNativeImage(*theirs.nativeImage)
	} else if b.nativeImage != nil {
		builder.WithNativeImage(*b.nativeImage)
	}

	if theirs.exportComponentModel != nil {
		builder.WithExportComponentModel(*theirs.exportComponentModel)
	} else if b.exportComponentModel != nil {
		builder.WithExportComponentModel(*b.exportComponentModel)
	}

	if theirs.graalVMBuildOptions != "" {
		builder.WithGraalVMBuildOptions(theirs.graalVMBuildOptions)
	} else {
		builder.WithGraalVMBuildOptions(b.graalVMBuildOptions)
	}

	if theirs.showDependencyDiagnostics != nil {
		builder.WithShowDependencyDiagnostics(*theirs.showDependencyDiagnostics)
	} else if b.showDependencyDiagnostics != nil {
		builder.WithShowDependencyDiagnostics(*b.showDependencyDiagnostics)
	}

	// Merge compilation options
	mergedCompilationOptions := b.compilationOptions.AcceptTheirs(theirs.compilationOptions)
	builder.compilationOptionsBuilder = newCompilationOptionsBuilderFrom(mergedCompilationOptions)

	return builder.Build()
}

// BuildOptionsBuilder provides a builder pattern for BuildOptions.
// It contains a CompilationOptionsBuilder to build the embedded CompilationOptions.
// Java: io.ballerina.projects.BuildOptions.BuildOptionsBuilder
type BuildOptionsBuilder struct {
	// Build-specific fields
	testReport                *bool
	codeCoverage              *bool
	dumpBuildTime             *bool
	skipTests                 *bool
	targetDir                 string
	nativeImage               *bool
	exportComponentModel      *bool
	graalVMBuildOptions       string
	showDependencyDiagnostics *bool

	// Embedded builder for compilation options
	compilationOptionsBuilder *CompilationOptionsBuilder
}

// NewBuildOptionsBuilder creates a new BuildOptionsBuilder.
func NewBuildOptionsBuilder() *BuildOptionsBuilder {
	return &BuildOptionsBuilder{
		compilationOptionsBuilder: NewCompilationOptionsBuilder(),
	}
}

// WithTestReport sets whether test report generation is enabled.
func (b *BuildOptionsBuilder) WithTestReport(value bool) *BuildOptionsBuilder {
	b.testReport = &value
	return b
}

// WithCodeCoverage sets whether code coverage is enabled.
func (b *BuildOptionsBuilder) WithCodeCoverage(value bool) *BuildOptionsBuilder {
	b.codeCoverage = &value
	return b
}

// WithDumpBuildTime sets whether build time dumping is enabled.
func (b *BuildOptionsBuilder) WithDumpBuildTime(value bool) *BuildOptionsBuilder {
	b.dumpBuildTime = &value
	return b
}

// WithSkipTests sets whether tests should be skipped.
func (b *BuildOptionsBuilder) WithSkipTests(value bool) *BuildOptionsBuilder {
	b.skipTests = &value
	return b
}

// WithTargetDir sets the target directory path.
func (b *BuildOptionsBuilder) WithTargetDir(path string) *BuildOptionsBuilder {
	b.targetDir = path
	return b
}

// WithNativeImage sets whether GraalVM native image is enabled.
func (b *BuildOptionsBuilder) WithNativeImage(value bool) *BuildOptionsBuilder {
	b.nativeImage = &value
	return b
}

// WithGraalVM sets whether GraalVM native image is enabled.
// Alias for WithNativeImage.
func (b *BuildOptionsBuilder) WithGraalVM(value bool) *BuildOptionsBuilder {
	return b.WithNativeImage(value)
}

// WithGraalVMBuildOptions sets additional GraalVM build options.
func (b *BuildOptionsBuilder) WithGraalVMBuildOptions(value string) *BuildOptionsBuilder {
	b.graalVMBuildOptions = value
	return b
}

// WithShowDependencyDiagnostics sets whether dependency diagnostics should be shown.
func (b *BuildOptionsBuilder) WithShowDependencyDiagnostics(value bool) *BuildOptionsBuilder {
	b.showDependencyDiagnostics = &value
	return b
}

// WithExportComponentModel sets whether component model export is enabled.
// This sets both the build-level flag and delegates to compilation options.
func (b *BuildOptionsBuilder) WithExportComponentModel(value bool) *BuildOptionsBuilder {
	b.exportComponentModel = &value
	b.compilationOptionsBuilder.WithExportComponentModel(value)
	return b
}

// --- Methods that delegate to CompilationOptionsBuilder ---

// WithDisableSyntaxTreeCaching sets whether syntax tree caching is disabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDisableSyntaxTreeCaching(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDisableSyntaxTree(value)
	return b
}

// WithOffline sets whether offline mode is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithOffline(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithOfflineBuild(value)
	return b
}

// WithSticky sets whether sticky mode is enabled.
// Deprecated: Use WithLockingMode() instead.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithSticky(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithSticky(value)
	return b
}

// WithListConflictedClasses sets whether conflicted classes should be listed.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithListConflictedClasses(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithListConflictedClasses(value)
	return b
}

// WithExperimental sets whether experimental features are enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithExperimental(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithExperimental(value)
	return b
}

// WithObservabilityIncluded sets whether observability is included.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithObservabilityIncluded(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithObservabilityIncluded(value)
	return b
}

// WithCloud sets the cloud target.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithCloud(value string) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithCloud(value)
	return b
}

// WithDumpAST sets whether AST dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpAST(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpAST(value)
	return b
}

// WithDumpBIR sets whether BIR dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpBIR(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpBIR(value)
	return b
}

// WithDumpBIRFile sets whether BIR file dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpBIRFile(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpBIRFile(value)
	return b
}

// WithDumpGraph sets whether graph dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpGraph(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpGraph(value)
	return b
}

// WithDumpRawGraphs sets whether raw graph dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpRawGraphs(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpRawGraphs(value)
	return b
}

// WithDumpTokens sets whether lexer token dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpTokens(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpTokens(value)
	return b
}

// WithDumpST sets whether syntax tree dumping is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithDumpST(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithDumpST(value)
	return b
}

// WithTraceRecovery sets whether error recovery tracing is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithTraceRecovery(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithTraceRecovery(value)
	return b
}

// WithConfigSchemaGen sets whether config schema generation is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithConfigSchemaGen(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithConfigSchemaGen(value)
	return b
}

// WithExportOpenAPI sets whether OpenAPI export is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithExportOpenAPI(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithExportOpenAPI(value)
	return b
}

// WithRemoteManagement sets whether remote management is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithRemoteManagement(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithRemoteManagement(value)
	return b
}

// WithOptimizeDependencyCompilation sets whether dependency compilation optimization is enabled.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithOptimizeDependencyCompilation(value bool) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithOptimizeDependencyCompilation(value)
	return b
}

// WithLockingMode sets the package locking mode.
// Delegates to CompilationOptionsBuilder.
func (b *BuildOptionsBuilder) WithLockingMode(mode PackageLockingMode) *BuildOptionsBuilder {
	b.compilationOptionsBuilder.WithLockingMode(mode)
	return b
}

// Build creates the BuildOptions instance.
// First builds CompilationOptions, then includes it in the BuildOptions.
func (b *BuildOptionsBuilder) Build() BuildOptions {
	compilationOptions := b.compilationOptionsBuilder.Build()
	return BuildOptions{
		testReport:                b.testReport,
		codeCoverage:              b.codeCoverage,
		dumpBuildTime:             b.dumpBuildTime,
		skipTests:                 b.skipTests,
		targetDir:                 b.targetDir,
		nativeImage:               b.nativeImage,
		exportComponentModel:      b.exportComponentModel,
		graalVMBuildOptions:       b.graalVMBuildOptions,
		showDependencyDiagnostics: b.showDependencyDiagnostics,
		compilationOptions:        compilationOptions,
	}
}

// newCompilationOptionsBuilderFrom creates a CompilationOptionsBuilder initialized
// with values from an existing CompilationOptions.
func newCompilationOptionsBuilderFrom(opts CompilationOptions) *CompilationOptionsBuilder {
	builder := NewCompilationOptionsBuilder()
	// Copy all fields from opts to builder
	builder.options = opts
	return builder
}

// OptionName represents build option names.
// Java: io.ballerina.projects.BuildOptions.OptionName
type OptionName string

const (
	OptionNameOffline                       OptionName = "offline"
	OptionNameSticky                        OptionName = "sticky"
	OptionNameLockingMode                   OptionName = "lockingMode"
	OptionNameObservabilityIncluded         OptionName = "observabilityIncluded"
	OptionNameExperimental                  OptionName = "experimental"
	OptionNameSkipTests                     OptionName = "skipTests"
	OptionNameTestReport                    OptionName = "testReport"
	OptionNameCodeCoverage                  OptionName = "codeCoverage"
	OptionNameListConflictedClasses         OptionName = "listConflictedClasses"
	OptionNameDumpBuildTime                 OptionName = "dumpBuildTime"
	OptionNameTargetDir                     OptionName = "targetDir"
	OptionNameNativeImage                   OptionName = "graalvm"
	OptionNameExportComponentModel          OptionName = "exportComponentModel"
	OptionNameGraalVMBuildOptions           OptionName = "graalvmBuildOptions"
	OptionNameShowDependencyDiagnostics     OptionName = "showDependencyDiagnostics"
	OptionNameOptimizeDependencyCompilation OptionName = "optimizeDependencyCompilation"
	OptionNameRemoteManagement              OptionName = "remoteManagement"
	OptionNameCloud                         OptionName = "cloud"
)

// String returns the string representation of the option name.
func (n OptionName) String() string {
	return string(n)
}
