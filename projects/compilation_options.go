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
	offlineBuild                  *bool
	experimental                  *bool
	observabilityIncluded         *bool
	dumpAST                       *bool
	dumpBIR                       *bool
	dumpBIRFile                   *bool
	dumpCFG                       *bool
	dumpGraph                     *bool
	dumpRawGraphs                 *bool
	dumpTokens                    *bool
	dumpST                        *bool
	listConflictedClasses         *bool
	sticky                        *bool
	withCodeGenerators            *bool
	withCodeModifiers             *bool
	configSchemaGen               *bool
	exportOpenAPI                 *bool
	exportComponentModel          *bool
	disableSyntaxTree             *bool
	remoteManagement              *bool
	optimizeDependencyCompilation *bool
	traceRecovery                 *bool
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
	return toBoolDefaultIfNull(c.offlineBuild)
}

// Experimental returns whether experimental features are enabled.
func (c CompilationOptions) Experimental() bool {
	return toBoolDefaultIfNull(c.experimental)
}

// ObservabilityIncluded returns whether observability is included.
func (c CompilationOptions) ObservabilityIncluded() bool {
	return toBoolDefaultIfNull(c.observabilityIncluded)
}

// DumpAST returns whether AST dumping is enabled.
func (c CompilationOptions) DumpAST() bool {
	return toBoolDefaultIfNull(c.dumpAST)
}

// DumpBIR returns whether BIR dumping is enabled.
func (c CompilationOptions) DumpBIR() bool {
	return toBoolDefaultIfNull(c.dumpBIR)
}

// DumpBIRFile returns whether BIR file dumping is enabled.
func (c CompilationOptions) DumpBIRFile() bool {
	return toBoolDefaultIfNull(c.dumpBIRFile)
}

// DumpCFG returns whether CFG dumping is enabled.
func (c CompilationOptions) DumpCFG() bool {
	return toBoolDefaultIfNull(c.dumpCFG)
}

// DumpCFGFormat returns the CFG dump format.
// Returns CFGFormatUnknown if not set (callers should treat as S-expression default).
func (c CompilationOptions) DumpCFGFormat() CFGFormat {
	return c.dumpCFGFormat
}

// DumpGraph returns whether graph dumping is enabled.
func (c CompilationOptions) DumpGraph() bool {
	return toBoolDefaultIfNull(c.dumpGraph)
}

// DumpRawGraphs returns whether raw graph dumping is enabled.
func (c CompilationOptions) DumpRawGraphs() bool {
	return toBoolDefaultIfNull(c.dumpRawGraphs)
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
	return toBoolDefaultIfNull(c.listConflictedClasses)
}

// Sticky returns whether sticky mode is enabled.
// Deprecated: Use LockingMode() instead.
func (c CompilationOptions) Sticky() bool {
	return toBoolDefaultIfNull(c.sticky)
}

// WithCodeGenerators returns whether code generators should be run.
func (c CompilationOptions) WithCodeGenerators() bool {
	return toBoolDefaultIfNull(c.withCodeGenerators)
}

// WithCodeModifiers returns whether code modifiers should be run.
func (c CompilationOptions) WithCodeModifiers() bool {
	return toBoolDefaultIfNull(c.withCodeModifiers)
}

// ConfigSchemaGen returns whether config schema generation is enabled.
func (c CompilationOptions) ConfigSchemaGen() bool {
	return toBoolDefaultIfNull(c.configSchemaGen)
}

// ExportOpenAPI returns whether OpenAPI export is enabled.
func (c CompilationOptions) ExportOpenAPI() bool {
	return toBoolDefaultIfNull(c.exportOpenAPI)
}

// ExportComponentModel returns whether component model export is enabled.
func (c CompilationOptions) ExportComponentModel() bool {
	return toBoolDefaultIfNull(c.exportComponentModel)
}

// DisableSyntaxTree returns whether syntax tree caching is disabled.
func (c CompilationOptions) DisableSyntaxTree() bool {
	return toBoolDefaultIfNull(c.disableSyntaxTree)
}

// RemoteManagement returns whether remote management is enabled.
func (c CompilationOptions) RemoteManagement() bool {
	return toBoolDefaultIfNull(c.remoteManagement)
}

// OptimizeDependencyCompilation returns whether dependency compilation optimization is enabled.
func (c CompilationOptions) OptimizeDependencyCompilation() bool {
	return toBoolDefaultIfNull(c.optimizeDependencyCompilation)
}

// DumpTokens returns whether lexer token dumping is enabled.
func (c CompilationOptions) DumpTokens() bool {
	return toBoolDefaultIfNull(c.dumpTokens)
}

// DumpST returns whether syntax tree dumping is enabled.
func (c CompilationOptions) DumpST() bool {
	return toBoolDefaultIfNull(c.dumpST)
}

// TraceRecovery returns whether error recovery tracing is enabled.
func (c CompilationOptions) TraceRecovery() bool {
	return toBoolDefaultIfNull(c.traceRecovery)
}

// LockingMode returns the package locking mode.
// Returns PackageLockingModeUnknown if not explicitly set.
func (c CompilationOptions) LockingMode() PackageLockingMode {
	return c.lockingMode
}

// AcceptTheirs merges the given compilation options by favoring theirs if there are conflicts.
func (c CompilationOptions) AcceptTheirs(theirs CompilationOptions) CompilationOptions {
	builder := NewCompilationOptionsBuilder()

	// For each field, use theirs if set, otherwise use ours
	if theirs.offlineBuild != nil {
		builder.WithOfflineBuild(*theirs.offlineBuild)
	} else if c.offlineBuild != nil {
		builder.WithOfflineBuild(*c.offlineBuild)
	}

	if theirs.experimental != nil {
		builder.WithExperimental(*theirs.experimental)
	} else if c.experimental != nil {
		builder.WithExperimental(*c.experimental)
	}

	if theirs.observabilityIncluded != nil {
		builder.WithObservabilityIncluded(*theirs.observabilityIncluded)
	} else if c.observabilityIncluded != nil {
		builder.WithObservabilityIncluded(*c.observabilityIncluded)
	}

	if theirs.dumpAST != nil {
		builder.WithDumpAST(*theirs.dumpAST)
	} else if c.dumpAST != nil {
		builder.WithDumpAST(*c.dumpAST)
	}

	if theirs.dumpBIR != nil {
		builder.WithDumpBIR(*theirs.dumpBIR)
	} else if c.dumpBIR != nil {
		builder.WithDumpBIR(*c.dumpBIR)
	}

	if theirs.dumpBIRFile != nil {
		builder.WithDumpBIRFile(*theirs.dumpBIRFile)
	} else if c.dumpBIRFile != nil {
		builder.WithDumpBIRFile(*c.dumpBIRFile)
	}

	if theirs.dumpCFG != nil {
		builder.WithDumpCFG(*theirs.dumpCFG)
	} else if c.dumpCFG != nil {
		builder.WithDumpCFG(*c.dumpCFG)
	}

	if theirs.dumpCFGFormat != CFGFormatUnknown {
		builder.WithDumpCFGFormat(theirs.dumpCFGFormat)
	} else if c.dumpCFGFormat != CFGFormatUnknown {
		builder.WithDumpCFGFormat(c.dumpCFGFormat)
	}

	if theirs.dumpGraph != nil {
		builder.WithDumpGraph(*theirs.dumpGraph)
	} else if c.dumpGraph != nil {
		builder.WithDumpGraph(*c.dumpGraph)
	}

	if theirs.dumpRawGraphs != nil {
		builder.WithDumpRawGraphs(*theirs.dumpRawGraphs)
	} else if c.dumpRawGraphs != nil {
		builder.WithDumpRawGraphs(*c.dumpRawGraphs)
	}

	if theirs.cloud != nil {
		builder.WithCloud(*theirs.cloud)
	} else if c.cloud != nil {
		builder.WithCloud(*c.cloud)
	}

	if theirs.listConflictedClasses != nil {
		builder.WithListConflictedClasses(*theirs.listConflictedClasses)
	} else if c.listConflictedClasses != nil {
		builder.WithListConflictedClasses(*c.listConflictedClasses)
	}

	if theirs.sticky != nil {
		builder.WithSticky(*theirs.sticky)
	} else if c.sticky != nil {
		builder.WithSticky(*c.sticky)
	}

	if theirs.withCodeGenerators != nil {
		builder.WithCodeGeneratorsEnabled(*theirs.withCodeGenerators)
	} else if c.withCodeGenerators != nil {
		builder.WithCodeGeneratorsEnabled(*c.withCodeGenerators)
	}

	if theirs.withCodeModifiers != nil {
		builder.WithCodeModifiersEnabled(*theirs.withCodeModifiers)
	} else if c.withCodeModifiers != nil {
		builder.WithCodeModifiersEnabled(*c.withCodeModifiers)
	}

	if theirs.configSchemaGen != nil {
		builder.WithConfigSchemaGen(*theirs.configSchemaGen)
	} else if c.configSchemaGen != nil {
		builder.WithConfigSchemaGen(*c.configSchemaGen)
	}

	if theirs.exportOpenAPI != nil {
		builder.WithExportOpenAPI(*theirs.exportOpenAPI)
	} else if c.exportOpenAPI != nil {
		builder.WithExportOpenAPI(*c.exportOpenAPI)
	}

	if theirs.exportComponentModel != nil {
		builder.WithExportComponentModel(*theirs.exportComponentModel)
	} else if c.exportComponentModel != nil {
		builder.WithExportComponentModel(*c.exportComponentModel)
	}

	if theirs.disableSyntaxTree != nil {
		builder.WithDisableSyntaxTree(*theirs.disableSyntaxTree)
	} else if c.disableSyntaxTree != nil {
		builder.WithDisableSyntaxTree(*c.disableSyntaxTree)
	}

	if theirs.remoteManagement != nil {
		builder.WithRemoteManagement(*theirs.remoteManagement)
	} else if c.remoteManagement != nil {
		builder.WithRemoteManagement(*c.remoteManagement)
	}

	if theirs.optimizeDependencyCompilation != nil {
		builder.WithOptimizeDependencyCompilation(*theirs.optimizeDependencyCompilation)
	} else if c.optimizeDependencyCompilation != nil {
		builder.WithOptimizeDependencyCompilation(*c.optimizeDependencyCompilation)
	}

	if theirs.dumpTokens != nil {
		builder.WithDumpTokens(*theirs.dumpTokens)
	} else if c.dumpTokens != nil {
		builder.WithDumpTokens(*c.dumpTokens)
	}

	if theirs.dumpST != nil {
		builder.WithDumpST(*theirs.dumpST)
	} else if c.dumpST != nil {
		builder.WithDumpST(*c.dumpST)
	}

	if theirs.traceRecovery != nil {
		builder.WithTraceRecovery(*theirs.traceRecovery)
	} else if c.traceRecovery != nil {
		builder.WithTraceRecovery(*c.traceRecovery)
	}

	// Handle locking mode with special sticky logic
	if theirs.lockingMode != PackageLockingModeUnknown {
		builder.WithLockingMode(theirs.lockingMode)
	} else if builder.options.sticky != nil && *builder.options.sticky {
		// If sticky is true, set locking mode to HARD unless theirs has a locking mode
		builder.WithLockingMode(PackageLockingModeHard)
	} else {
		builder.WithLockingMode(c.lockingMode)
	}

	return builder.Build()
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
	b.options.offlineBuild = &value
	return b
}

// WithExperimental sets experimental features flag.
func (b *CompilationOptionsBuilder) WithExperimental(value bool) *CompilationOptionsBuilder {
	b.options.experimental = &value
	return b
}

// WithObservabilityIncluded sets observability inclusion.
func (b *CompilationOptionsBuilder) WithObservabilityIncluded(value bool) *CompilationOptionsBuilder {
	b.options.observabilityIncluded = &value
	return b
}

// WithDumpAST sets AST dumping flag.
func (b *CompilationOptionsBuilder) WithDumpAST(value bool) *CompilationOptionsBuilder {
	b.options.dumpAST = &value
	return b
}

// WithDumpBIR sets BIR dumping flag.
func (b *CompilationOptionsBuilder) WithDumpBIR(value bool) *CompilationOptionsBuilder {
	b.options.dumpBIR = &value
	return b
}

// WithDumpBIRFile sets BIR file dumping flag.
func (b *CompilationOptionsBuilder) WithDumpBIRFile(value bool) *CompilationOptionsBuilder {
	b.options.dumpBIRFile = &value
	return b
}

// WithDumpCFG sets CFG dumping flag.
func (b *CompilationOptionsBuilder) WithDumpCFG(value bool) *CompilationOptionsBuilder {
	b.options.dumpCFG = &value
	return b
}

// WithDumpCFGFormat sets CFG dump format.
func (b *CompilationOptionsBuilder) WithDumpCFGFormat(value CFGFormat) *CompilationOptionsBuilder {
	b.options.dumpCFGFormat = value
	return b
}

// WithDumpGraph sets graph dumping flag.
func (b *CompilationOptionsBuilder) WithDumpGraph(value bool) *CompilationOptionsBuilder {
	b.options.dumpGraph = &value
	return b
}

// WithDumpRawGraphs sets raw graph dumping flag.
func (b *CompilationOptionsBuilder) WithDumpRawGraphs(value bool) *CompilationOptionsBuilder {
	b.options.dumpRawGraphs = &value
	return b
}

// WithCloud sets the cloud target.
func (b *CompilationOptionsBuilder) WithCloud(value string) *CompilationOptionsBuilder {
	b.options.cloud = &value
	return b
}

// WithListConflictedClasses sets conflicted classes listing flag.
func (b *CompilationOptionsBuilder) WithListConflictedClasses(value bool) *CompilationOptionsBuilder {
	b.options.listConflictedClasses = &value
	return b
}

// WithSticky sets sticky mode.
// Deprecated: Use WithLockingMode() instead.
func (b *CompilationOptionsBuilder) WithSticky(value bool) *CompilationOptionsBuilder {
	b.options.sticky = &value
	return b
}

// WithCodeGeneratorsEnabled sets whether code generators should run.
func (b *CompilationOptionsBuilder) WithCodeGeneratorsEnabled(value bool) *CompilationOptionsBuilder {
	b.options.withCodeGenerators = &value
	return b
}

// WithCodeModifiersEnabled sets whether code modifiers should run.
func (b *CompilationOptionsBuilder) WithCodeModifiersEnabled(value bool) *CompilationOptionsBuilder {
	b.options.withCodeModifiers = &value
	return b
}

// WithConfigSchemaGen sets config schema generation flag.
func (b *CompilationOptionsBuilder) WithConfigSchemaGen(value bool) *CompilationOptionsBuilder {
	b.options.configSchemaGen = &value
	return b
}

// WithExportOpenAPI sets OpenAPI export flag.
func (b *CompilationOptionsBuilder) WithExportOpenAPI(value bool) *CompilationOptionsBuilder {
	b.options.exportOpenAPI = &value
	return b
}

// WithExportComponentModel sets component model export flag.
func (b *CompilationOptionsBuilder) WithExportComponentModel(value bool) *CompilationOptionsBuilder {
	b.options.exportComponentModel = &value
	return b
}

// WithDisableSyntaxTree sets syntax tree caching disabled flag.
func (b *CompilationOptionsBuilder) WithDisableSyntaxTree(value bool) *CompilationOptionsBuilder {
	b.options.disableSyntaxTree = &value
	return b
}

// WithRemoteManagement sets remote management flag.
func (b *CompilationOptionsBuilder) WithRemoteManagement(value bool) *CompilationOptionsBuilder {
	b.options.remoteManagement = &value
	return b
}

// WithOptimizeDependencyCompilation sets dependency compilation optimization flag.
func (b *CompilationOptionsBuilder) WithOptimizeDependencyCompilation(value bool) *CompilationOptionsBuilder {
	b.options.optimizeDependencyCompilation = &value
	return b
}

// WithDumpTokens sets lexer token dumping flag.
func (b *CompilationOptionsBuilder) WithDumpTokens(value bool) *CompilationOptionsBuilder {
	b.options.dumpTokens = &value
	return b
}

// WithDumpST sets syntax tree dumping flag.
func (b *CompilationOptionsBuilder) WithDumpST(value bool) *CompilationOptionsBuilder {
	b.options.dumpST = &value
	return b
}

// WithTraceRecovery sets error recovery tracing flag.
func (b *CompilationOptionsBuilder) WithTraceRecovery(value bool) *CompilationOptionsBuilder {
	b.options.traceRecovery = &value
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

// toBoolDefaultIfNull returns false if the pointer is nil, otherwise returns the value.
func toBoolDefaultIfNull(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
