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

// ModuleConfig represents configuration for a Ballerina module.
// It contains the module's source documents, test documents, and dependencies.
// Java source: io.ballerina.projects.ModuleConfig
type ModuleConfig struct {
	moduleID         ModuleID
	moduleDescriptor ModuleDescriptor
	sourceDocs       []DocumentConfig
	testSourceDocs   []DocumentConfig
	dependencies     []ModuleDescriptor
	readmeMd         DocumentConfig // can be nil
}

// NewModuleConfig creates a new ModuleConfig.
// Java equivalent: ModuleConfig.from(ModuleId, ModuleDescriptor, List<DocumentConfig>, List<DocumentConfig>, DocumentConfig, List<ModuleDescriptor>)
func NewModuleConfig(
	moduleID ModuleID,
	moduleDescriptor ModuleDescriptor,
	sourceDocs []DocumentConfig,
	testSourceDocs []DocumentConfig,
	readmeMd DocumentConfig,
	dependencies []ModuleDescriptor,
) ModuleConfig {
	// Make defensive copies
	srcCopy := make([]DocumentConfig, len(sourceDocs))
	copy(srcCopy, sourceDocs)

	testCopy := make([]DocumentConfig, len(testSourceDocs))
	copy(testCopy, testSourceDocs)

	depsCopy := make([]ModuleDescriptor, len(dependencies))
	copy(depsCopy, dependencies)

	return ModuleConfig{
		moduleID:         moduleID,
		moduleDescriptor: moduleDescriptor,
		sourceDocs:       srcCopy,
		testSourceDocs:   testCopy,
		dependencies:     depsCopy,
		readmeMd:         readmeMd,
	}
}

// ModuleID returns the unique identifier for this module.
func (m ModuleConfig) ModuleID() ModuleID {
	return m.moduleID
}

// ModuleDescriptor returns the module descriptor (metadata).
func (m ModuleConfig) ModuleDescriptor() ModuleDescriptor {
	return m.moduleDescriptor
}

// IsDefaultModule returns true if this is the default module of the package.
func (m ModuleConfig) IsDefaultModule() bool {
	return m.moduleDescriptor.Name().IsDefaultModuleName()
}

// SourceDocs returns a copy of the source document configurations.
func (m ModuleConfig) SourceDocs() []DocumentConfig {
	result := make([]DocumentConfig, len(m.sourceDocs))
	copy(result, m.sourceDocs)
	return result
}

// TestSourceDocs returns a copy of the test source document configurations.
func (m ModuleConfig) TestSourceDocs() []DocumentConfig {
	result := make([]DocumentConfig, len(m.testSourceDocs))
	copy(result, m.testSourceDocs)
	return result
}

// Dependencies returns a copy of the module dependencies.
func (m ModuleConfig) Dependencies() []ModuleDescriptor {
	result := make([]ModuleDescriptor, len(m.dependencies))
	copy(result, m.dependencies)
	return result
}

// ReadmeMd returns the README.md document config, or nil if not present.
func (m ModuleConfig) ReadmeMd() DocumentConfig {
	return m.readmeMd
}

// HasReadmeMd returns true if this module has a README.md file.
func (m ModuleConfig) HasReadmeMd() bool {
	return m.readmeMd != nil
}
