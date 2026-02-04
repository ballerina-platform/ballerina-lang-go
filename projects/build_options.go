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

// BuildOptions represents the [build-options] section in Ballerina.toml.
// Java source: io.ballerina.projects.BuildOptions
type BuildOptions struct {
	observabilityIncluded *bool
	offline               *bool
	skipTests             *bool
	testReport            *bool
	codeCoverage          *bool
	cloud                 string
	listConflictedClasses *bool
	dumpBIR               *bool
	dumpBIRFile           *bool
	dumpGraph             *bool
	dumpRawGraphs         *bool
	sticky                *bool
	graalVM               *bool
	graalVMBuildOptions   string
	remoteManagement      *bool
	exportOpenAPI         *bool
	exportComponentModel  *bool
}

// NewBuildOptions creates a new BuildOptions with default values.
func NewBuildOptions() BuildOptions {
	return BuildOptions{}
}

// ObservabilityIncluded returns whether observability is included.
func (b BuildOptions) ObservabilityIncluded() bool {
	if b.observabilityIncluded == nil {
		return false
	}
	return *b.observabilityIncluded
}

// Offline returns whether offline mode is enabled.
func (b BuildOptions) Offline() bool {
	if b.offline == nil {
		return false
	}
	return *b.offline
}

// SkipTests returns whether tests should be skipped.
func (b BuildOptions) SkipTests() bool {
	if b.skipTests == nil {
		return false
	}
	return *b.skipTests
}

// TestReport returns whether test report generation is enabled.
func (b BuildOptions) TestReport() bool {
	if b.testReport == nil {
		return false
	}
	return *b.testReport
}

// CodeCoverage returns whether code coverage is enabled.
func (b BuildOptions) CodeCoverage() bool {
	if b.codeCoverage == nil {
		return false
	}
	return *b.codeCoverage
}

// Cloud returns the cloud target (e.g., "k8s", "docker").
func (b BuildOptions) Cloud() string {
	return b.cloud
}

// GraalVM returns whether GraalVM native image is enabled.
func (b BuildOptions) GraalVM() bool {
	if b.graalVM == nil {
		return false
	}
	return *b.graalVM
}

// GraalVMBuildOptions returns additional GraalVM build options.
func (b BuildOptions) GraalVMBuildOptions() string {
	return b.graalVMBuildOptions
}

// Sticky returns whether sticky mode is enabled.
func (b BuildOptions) Sticky() bool {
	if b.sticky == nil {
		return false
	}
	return *b.sticky
}

// BuildOptionsBuilder provides a builder pattern for BuildOptions.
type BuildOptionsBuilder struct {
	options BuildOptions
}

// NewBuildOptionsBuilder creates a new BuildOptionsBuilder.
func NewBuildOptionsBuilder() *BuildOptionsBuilder {
	return &BuildOptionsBuilder{
		options: NewBuildOptions(),
	}
}

// WithObservabilityIncluded sets whether observability is included.
func (b *BuildOptionsBuilder) WithObservabilityIncluded(value bool) *BuildOptionsBuilder {
	b.options.observabilityIncluded = &value
	return b
}

// WithOffline sets whether offline mode is enabled.
func (b *BuildOptionsBuilder) WithOffline(value bool) *BuildOptionsBuilder {
	b.options.offline = &value
	return b
}

// WithSkipTests sets whether tests should be skipped.
func (b *BuildOptionsBuilder) WithSkipTests(value bool) *BuildOptionsBuilder {
	b.options.skipTests = &value
	return b
}

// WithTestReport sets whether test report generation is enabled.
func (b *BuildOptionsBuilder) WithTestReport(value bool) *BuildOptionsBuilder {
	b.options.testReport = &value
	return b
}

// WithCodeCoverage sets whether code coverage is enabled.
func (b *BuildOptionsBuilder) WithCodeCoverage(value bool) *BuildOptionsBuilder {
	b.options.codeCoverage = &value
	return b
}

// WithCloud sets the cloud target.
func (b *BuildOptionsBuilder) WithCloud(value string) *BuildOptionsBuilder {
	b.options.cloud = value
	return b
}

// WithGraalVM sets whether GraalVM native image is enabled.
func (b *BuildOptionsBuilder) WithGraalVM(value bool) *BuildOptionsBuilder {
	b.options.graalVM = &value
	return b
}

// WithGraalVMBuildOptions sets additional GraalVM build options.
func (b *BuildOptionsBuilder) WithGraalVMBuildOptions(value string) *BuildOptionsBuilder {
	b.options.graalVMBuildOptions = value
	return b
}

// WithSticky sets whether sticky mode is enabled.
func (b *BuildOptionsBuilder) WithSticky(value bool) *BuildOptionsBuilder {
	b.options.sticky = &value
	return b
}

// Build creates the BuildOptions instance.
func (b *BuildOptionsBuilder) Build() BuildOptions {
	return b.options
}
