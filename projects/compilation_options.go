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

// OptionalBool is a tri-state boolean: unknown (unset), true, or false.
// It uses a single byte instead of *bool (pointer + bool = 9+ bytes with heap allocation).
type OptionalBool byte

const (
	// OptionalBoolUnknown is the zero value representing an unset boolean.
	OptionalBoolUnknown OptionalBool = iota
	// OptionalBoolTrue represents an explicitly set true value.
	OptionalBoolTrue
	// OptionalBoolFalse represents an explicitly set false value.
	OptionalBoolFalse
)

// Of converts a bool to an OptionalBool.
func Of(value bool) OptionalBool {
	if value {
		return OptionalBoolTrue
	}
	return OptionalBoolFalse
}

// IsSet returns true if this optional has been explicitly set (true or false).
func (o OptionalBool) IsSet() bool {
	return o != OptionalBoolUnknown
}

// ValueOr returns the bool value, or the given default if unset.
func (o OptionalBool) ValueOr(def bool) bool {
	switch o {
	case OptionalBoolTrue:
		return true
	case OptionalBoolFalse:
		return false
	default:
		return def
	}
}

// PackageLockingMode represents the locking mode for package dependencies.
type PackageLockingMode int

const (
	// PackageLockingModeUnknown represents an unset locking mode.
	PackageLockingModeUnknown PackageLockingMode = iota
	// PackageLockingModeSoft allows minor and patch version updates.
	PackageLockingModeSoft
	// PackageLockingModeMedium allows only patch version updates.
	PackageLockingModeMedium
	// PackageLockingModeHard locks to exact versions.
	PackageLockingModeHard
)

// String returns the string representation of PackageLockingMode.
func (m PackageLockingMode) String() string {
	switch m {
	case PackageLockingModeSoft:
		return "soft"
	case PackageLockingModeMedium:
		return "medium"
	case PackageLockingModeHard:
		return "hard"
	default:
		return "unknown"
	}
}

// CFGFormat represents the output format for CFG (Control Flow Graph) dumps.
type CFGFormat int

const (
	// CFGFormatUnknown represents an unset format (zero value).
	// Callers should treat this as "use default" which is S-expression.
	CFGFormatUnknown CFGFormat = iota
	// CFGFormatSexp uses S-expression format (the default).
	CFGFormatSexp
	// CFGFormatDot uses Graphviz DOT format.
	CFGFormatDot
)

// String returns the string representation of CFGFormat.
func (f CFGFormat) String() string {
	switch f {
	case CFGFormatSexp:
		return "sexp"
	case CFGFormatDot:
		return "dot"
	default:
		return ""
	}
}

// ParseCFGFormat parses a string into a CFGFormat value.
func ParseCFGFormat(s string) CFGFormat {
	switch s {
	case "dot":
		return CFGFormatDot
	case "sexp", "":
		return CFGFormatSexp
	default:
		return CFGFormatUnknown
	}
}

// CompilationOptions holds compilation-specific options.
// BuildOptions contains a CompilationOptions instance and delegates compilation-related methods to it.
type CompilationOptions struct {
	offlineBuild                  OptionalBool
	experimental                  OptionalBool
	observabilityIncluded         OptionalBool
	dumpAST                       OptionalBool
	dumpBIR                       OptionalBool
	dumpBIRFile                   OptionalBool
	dumpCFG                       OptionalBool
	dumpGraph                     OptionalBool
	dumpRawGraphs                 OptionalBool
	dumpTokens                    OptionalBool
	dumpST                        OptionalBool
	listConflictedClasses         OptionalBool
	sticky                        OptionalBool
	withCodeGenerators            OptionalBool
	withCodeModifiers             OptionalBool
	configSchemaGen               OptionalBool
	exportOpenAPI                 OptionalBool
	exportComponentModel          OptionalBool
	disableSyntaxTree             OptionalBool
	remoteManagement              OptionalBool
	optimizeDependencyCompilation OptionalBool
	traceRecovery                 OptionalBool
	cloud                         *string
	dumpCFGFormat                 CFGFormat
	lockingMode                   PackageLockingMode
}

// NewCompilationOptions creates a new CompilationOptions with default values.
func NewCompilationOptions() CompilationOptions {
	return CompilationOptions{}
}

// OfflineBuild returns whether offline build mode is enabled.
func (c CompilationOptions) OfflineBuild() bool {
	return c.offlineBuild.ValueOr(false)
}

// Experimental returns whether experimental features are enabled.
func (c CompilationOptions) Experimental() bool {
	return c.experimental.ValueOr(false)
}

// ObservabilityIncluded returns whether observability is included.
func (c CompilationOptions) ObservabilityIncluded() bool {
	return c.observabilityIncluded.ValueOr(false)
}

// DumpAST returns whether AST dumping is enabled.
func (c CompilationOptions) DumpAST() bool {
	return c.dumpAST.ValueOr(false)
}

// DumpBIR returns whether BIR dumping is enabled.
func (c CompilationOptions) DumpBIR() bool {
	return c.dumpBIR.ValueOr(false)
}

// DumpBIRFile returns whether BIR file dumping is enabled.
func (c CompilationOptions) DumpBIRFile() bool {
	return c.dumpBIRFile.ValueOr(false)
}

// DumpCFG returns whether CFG dumping is enabled.
func (c CompilationOptions) DumpCFG() bool {
	return c.dumpCFG.ValueOr(false)
}

// DumpCFGFormat returns the CFG dump format.
// Returns CFGFormatUnknown if not set (callers should treat as S-expression default).
func (c CompilationOptions) DumpCFGFormat() CFGFormat {
	return c.dumpCFGFormat
}

// DumpGraph returns whether graph dumping is enabled.
func (c CompilationOptions) DumpGraph() bool {
	return c.dumpGraph.ValueOr(false)
}

// DumpRawGraphs returns whether raw graph dumping is enabled.
func (c CompilationOptions) DumpRawGraphs() bool {
	return c.dumpRawGraphs.ValueOr(false)
}

// Cloud returns the cloud target.
// Returns empty string if not set or explicitly cleared.
func (c CompilationOptions) Cloud() string {
	if c.cloud == nil {
		return ""
	}
	return *c.cloud
}

// ListConflictedClasses returns whether conflicted classes should be listed.
func (c CompilationOptions) ListConflictedClasses() bool {
	return c.listConflictedClasses.ValueOr(false)
}

// Sticky returns whether sticky mode is enabled.
// Deprecated: Use LockingMode() instead.
func (c CompilationOptions) Sticky() bool {
	return c.sticky.ValueOr(false)
}

// WithCodeGenerators returns whether code generators should be run.
func (c CompilationOptions) WithCodeGenerators() bool {
	return c.withCodeGenerators.ValueOr(false)
}

// WithCodeModifiers returns whether code modifiers should be run.
func (c CompilationOptions) WithCodeModifiers() bool {
	return c.withCodeModifiers.ValueOr(false)
}

// ConfigSchemaGen returns whether config schema generation is enabled.
func (c CompilationOptions) ConfigSchemaGen() bool {
	return c.configSchemaGen.ValueOr(false)
}

// ExportOpenAPI returns whether OpenAPI export is enabled.
func (c CompilationOptions) ExportOpenAPI() bool {
	return c.exportOpenAPI.ValueOr(false)
}

// ExportComponentModel returns whether component model export is enabled.
func (c CompilationOptions) ExportComponentModel() bool {
	return c.exportComponentModel.ValueOr(false)
}

// DisableSyntaxTree returns whether syntax tree caching is disabled.
func (c CompilationOptions) DisableSyntaxTree() bool {
	return c.disableSyntaxTree.ValueOr(false)
}

// RemoteManagement returns whether remote management is enabled.
func (c CompilationOptions) RemoteManagement() bool {
	return c.remoteManagement.ValueOr(false)
}

// OptimizeDependencyCompilation returns whether dependency compilation optimization is enabled.
func (c CompilationOptions) OptimizeDependencyCompilation() bool {
	return c.optimizeDependencyCompilation.ValueOr(false)
}

// DumpTokens returns whether lexer token dumping is enabled.
func (c CompilationOptions) DumpTokens() bool {
	return c.dumpTokens.ValueOr(false)
}

// DumpST returns whether syntax tree dumping is enabled.
func (c CompilationOptions) DumpST() bool {
	return c.dumpST.ValueOr(false)
}

// TraceRecovery returns whether error recovery tracing is enabled.
func (c CompilationOptions) TraceRecovery() bool {
	return c.traceRecovery.ValueOr(false)
}

// LockingMode returns the package locking mode.
// Returns PackageLockingModeUnknown if not explicitly set.
func (c CompilationOptions) LockingMode() PackageLockingMode {
	return c.lockingMode
}

// acceptOptionalBool returns theirs if set, else ours.
func acceptOptionalBool(ours, theirs OptionalBool) OptionalBool {
	if theirs.IsSet() {
		return theirs
	}
	return ours
}

// AcceptTheirs merges the given compilation options by favoring theirs if there are conflicts.
func (c CompilationOptions) AcceptTheirs(theirs CompilationOptions) CompilationOptions {
	merged := CompilationOptions{
		offlineBuild:                  acceptOptionalBool(c.offlineBuild, theirs.offlineBuild),
		experimental:                  acceptOptionalBool(c.experimental, theirs.experimental),
		observabilityIncluded:         acceptOptionalBool(c.observabilityIncluded, theirs.observabilityIncluded),
		dumpAST:                       acceptOptionalBool(c.dumpAST, theirs.dumpAST),
		dumpBIR:                       acceptOptionalBool(c.dumpBIR, theirs.dumpBIR),
		dumpBIRFile:                   acceptOptionalBool(c.dumpBIRFile, theirs.dumpBIRFile),
		dumpCFG:                       acceptOptionalBool(c.dumpCFG, theirs.dumpCFG),
		dumpGraph:                     acceptOptionalBool(c.dumpGraph, theirs.dumpGraph),
		dumpRawGraphs:                 acceptOptionalBool(c.dumpRawGraphs, theirs.dumpRawGraphs),
		dumpTokens:                    acceptOptionalBool(c.dumpTokens, theirs.dumpTokens),
		dumpST:                        acceptOptionalBool(c.dumpST, theirs.dumpST),
		listConflictedClasses:         acceptOptionalBool(c.listConflictedClasses, theirs.listConflictedClasses),
		sticky:                        acceptOptionalBool(c.sticky, theirs.sticky),
		withCodeGenerators:            acceptOptionalBool(c.withCodeGenerators, theirs.withCodeGenerators),
		withCodeModifiers:             acceptOptionalBool(c.withCodeModifiers, theirs.withCodeModifiers),
		configSchemaGen:               acceptOptionalBool(c.configSchemaGen, theirs.configSchemaGen),
		exportOpenAPI:                 acceptOptionalBool(c.exportOpenAPI, theirs.exportOpenAPI),
		exportComponentModel:          acceptOptionalBool(c.exportComponentModel, theirs.exportComponentModel),
		disableSyntaxTree:             acceptOptionalBool(c.disableSyntaxTree, theirs.disableSyntaxTree),
		remoteManagement:              acceptOptionalBool(c.remoteManagement, theirs.remoteManagement),
		optimizeDependencyCompilation: acceptOptionalBool(c.optimizeDependencyCompilation, theirs.optimizeDependencyCompilation),
		traceRecovery:                 acceptOptionalBool(c.traceRecovery, theirs.traceRecovery),
	}

	// Cloud (*string)
	if theirs.cloud != nil {
		merged.cloud = theirs.cloud
	} else {
		merged.cloud = c.cloud
	}

	// CFG format
	if theirs.dumpCFGFormat != CFGFormatUnknown {
		merged.dumpCFGFormat = theirs.dumpCFGFormat
	} else {
		merged.dumpCFGFormat = c.dumpCFGFormat
	}

	// Locking mode with special sticky logic
	if theirs.lockingMode != PackageLockingModeUnknown {
		merged.lockingMode = theirs.lockingMode
	} else if merged.sticky == OptionalBoolTrue {
		merged.lockingMode = PackageLockingModeHard
	} else {
		merged.lockingMode = c.lockingMode
	}

	return merged
}

// CompilationOptionsBuilder provides a builder pattern for CompilationOptions.
type CompilationOptionsBuilder struct {
	options CompilationOptions
}

// NewCompilationOptionsBuilder creates a new CompilationOptionsBuilder.
func NewCompilationOptionsBuilder() *CompilationOptionsBuilder {
	return &CompilationOptionsBuilder{
		options: NewCompilationOptions(),
	}
}

// WithOfflineBuild sets offline build mode.
func (b *CompilationOptionsBuilder) WithOfflineBuild(value bool) *CompilationOptionsBuilder {
	b.options.offlineBuild = Of(value)
	return b
}

// WithExperimental sets experimental features flag.
func (b *CompilationOptionsBuilder) WithExperimental(value bool) *CompilationOptionsBuilder {
	b.options.experimental = Of(value)
	return b
}

// WithObservabilityIncluded sets observability inclusion.
func (b *CompilationOptionsBuilder) WithObservabilityIncluded(value bool) *CompilationOptionsBuilder {
	b.options.observabilityIncluded = Of(value)
	return b
}

// WithDumpAST sets AST dumping flag.
func (b *CompilationOptionsBuilder) WithDumpAST(value bool) *CompilationOptionsBuilder {
	b.options.dumpAST = Of(value)
	return b
}

// WithDumpBIR sets BIR dumping flag.
func (b *CompilationOptionsBuilder) WithDumpBIR(value bool) *CompilationOptionsBuilder {
	b.options.dumpBIR = Of(value)
	return b
}

// WithDumpBIRFile sets BIR file dumping flag.
func (b *CompilationOptionsBuilder) WithDumpBIRFile(value bool) *CompilationOptionsBuilder {
	b.options.dumpBIRFile = Of(value)
	return b
}

// WithDumpCFG sets CFG dumping flag.
func (b *CompilationOptionsBuilder) WithDumpCFG(value bool) *CompilationOptionsBuilder {
	b.options.dumpCFG = Of(value)
	return b
}

// WithDumpCFGFormat sets CFG dump format.
func (b *CompilationOptionsBuilder) WithDumpCFGFormat(value CFGFormat) *CompilationOptionsBuilder {
	b.options.dumpCFGFormat = value
	return b
}

// WithDumpGraph sets graph dumping flag.
func (b *CompilationOptionsBuilder) WithDumpGraph(value bool) *CompilationOptionsBuilder {
	b.options.dumpGraph = Of(value)
	return b
}

// WithDumpRawGraphs sets raw graph dumping flag.
func (b *CompilationOptionsBuilder) WithDumpRawGraphs(value bool) *CompilationOptionsBuilder {
	b.options.dumpRawGraphs = Of(value)
	return b
}

// WithCloud sets the cloud target.
func (b *CompilationOptionsBuilder) WithCloud(value string) *CompilationOptionsBuilder {
	b.options.cloud = &value
	return b
}

// WithListConflictedClasses sets conflicted classes listing flag.
func (b *CompilationOptionsBuilder) WithListConflictedClasses(value bool) *CompilationOptionsBuilder {
	b.options.listConflictedClasses = Of(value)
	return b
}

// WithSticky sets sticky mode.
// Deprecated: Use WithLockingMode() instead.
func (b *CompilationOptionsBuilder) WithSticky(value bool) *CompilationOptionsBuilder {
	b.options.sticky = Of(value)
	return b
}

// WithCodeGeneratorsEnabled sets whether code generators should run.
func (b *CompilationOptionsBuilder) WithCodeGeneratorsEnabled(value bool) *CompilationOptionsBuilder {
	b.options.withCodeGenerators = Of(value)
	return b
}

// WithCodeModifiersEnabled sets whether code modifiers should run.
func (b *CompilationOptionsBuilder) WithCodeModifiersEnabled(value bool) *CompilationOptionsBuilder {
	b.options.withCodeModifiers = Of(value)
	return b
}

// WithConfigSchemaGen sets config schema generation flag.
func (b *CompilationOptionsBuilder) WithConfigSchemaGen(value bool) *CompilationOptionsBuilder {
	b.options.configSchemaGen = Of(value)
	return b
}

// WithExportOpenAPI sets OpenAPI export flag.
func (b *CompilationOptionsBuilder) WithExportOpenAPI(value bool) *CompilationOptionsBuilder {
	b.options.exportOpenAPI = Of(value)
	return b
}

// WithExportComponentModel sets component model export flag.
func (b *CompilationOptionsBuilder) WithExportComponentModel(value bool) *CompilationOptionsBuilder {
	b.options.exportComponentModel = Of(value)
	return b
}

// WithDisableSyntaxTree sets syntax tree caching disabled flag.
func (b *CompilationOptionsBuilder) WithDisableSyntaxTree(value bool) *CompilationOptionsBuilder {
	b.options.disableSyntaxTree = Of(value)
	return b
}

// WithRemoteManagement sets remote management flag.
func (b *CompilationOptionsBuilder) WithRemoteManagement(value bool) *CompilationOptionsBuilder {
	b.options.remoteManagement = Of(value)
	return b
}

// WithOptimizeDependencyCompilation sets dependency compilation optimization flag.
func (b *CompilationOptionsBuilder) WithOptimizeDependencyCompilation(value bool) *CompilationOptionsBuilder {
	b.options.optimizeDependencyCompilation = Of(value)
	return b
}

// WithDumpTokens sets lexer token dumping flag.
func (b *CompilationOptionsBuilder) WithDumpTokens(value bool) *CompilationOptionsBuilder {
	b.options.dumpTokens = Of(value)
	return b
}

// WithDumpST sets syntax tree dumping flag.
func (b *CompilationOptionsBuilder) WithDumpST(value bool) *CompilationOptionsBuilder {
	b.options.dumpST = Of(value)
	return b
}

// WithTraceRecovery sets error recovery tracing flag.
func (b *CompilationOptionsBuilder) WithTraceRecovery(value bool) *CompilationOptionsBuilder {
	b.options.traceRecovery = Of(value)
	return b
}

// WithLockingMode sets the package locking mode.
func (b *CompilationOptionsBuilder) WithLockingMode(mode PackageLockingMode) *CompilationOptionsBuilder {
	b.options.lockingMode = mode
	return b
}

// Build creates the CompilationOptions instance.
func (b *CompilationOptionsBuilder) Build() CompilationOptions {
	return b.options
}
