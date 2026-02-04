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

// moduleContext holds internal state for a Module.
// It manages document contexts for source and test documents.
// Java: io.ballerina.projects.ModuleContext
type moduleContext struct {
	project                Project
	moduleID               ModuleID
	moduleDescriptor       ModuleDescriptor
	isDefaultModule        bool
	srcDocContextMap       map[DocumentID]*documentContext
	srcDocIDs              []DocumentID
	testDocContextMap      map[DocumentID]*documentContext
	testSrcDocIDs          []DocumentID
	moduleDescDependencies []ModuleDescriptor
}

// newModuleContext creates a moduleContext from ModuleConfig.
// Java: ModuleContext.from(Project, ModuleConfig, boolean)
func newModuleContext(project Project, moduleConfig ModuleConfig, disableSyntaxTree bool) *moduleContext {
	// Build source document context map
	srcDocContextMap := make(map[DocumentID]*documentContext)
	srcDocIDs := make([]DocumentID, 0, len(moduleConfig.SourceDocs()))
	for _, srcDocConfig := range moduleConfig.SourceDocs() {
		docID := srcDocConfig.DocumentID()
		srcDocContextMap[docID] = newDocumentContext(srcDocConfig, disableSyntaxTree)
		srcDocIDs = append(srcDocIDs, docID)
	}

	// Build test document context map
	testDocContextMap := make(map[DocumentID]*documentContext)
	testSrcDocIDs := make([]DocumentID, 0, len(moduleConfig.TestSourceDocs()))
	for _, testDocConfig := range moduleConfig.TestSourceDocs() {
		docID := testDocConfig.DocumentID()
		testDocContextMap[docID] = newDocumentContext(testDocConfig, disableSyntaxTree)
		testSrcDocIDs = append(testSrcDocIDs, docID)
	}

	// Copy dependencies
	deps := moduleConfig.Dependencies()
	depsCopy := make([]ModuleDescriptor, len(deps))
	copy(depsCopy, deps)

	return &moduleContext{
		project:                project,
		moduleID:               moduleConfig.ModuleID(),
		moduleDescriptor:       moduleConfig.ModuleDescriptor(),
		isDefaultModule:        moduleConfig.IsDefaultModule(),
		srcDocContextMap:       srcDocContextMap,
		srcDocIDs:              srcDocIDs,
		testDocContextMap:      testDocContextMap,
		testSrcDocIDs:          testSrcDocIDs,
		moduleDescDependencies: depsCopy,
	}
}

// newModuleContextFromMaps creates a moduleContext directly from document context maps.
// This is used for creating modified module contexts.
func newModuleContextFromMaps(
	project Project,
	moduleID ModuleID,
	moduleDescriptor ModuleDescriptor,
	isDefaultModule bool,
	srcDocContextMap map[DocumentID]*documentContext,
	testDocContextMap map[DocumentID]*documentContext,
	moduleDescDependencies []ModuleDescriptor,
) *moduleContext {
	// Build srcDocIDs from map keys
	srcDocIDs := make([]DocumentID, 0, len(srcDocContextMap))
	for id := range srcDocContextMap {
		srcDocIDs = append(srcDocIDs, id)
	}

	// Build testSrcDocIDs from map keys
	testSrcDocIDs := make([]DocumentID, 0, len(testDocContextMap))
	for id := range testDocContextMap {
		testSrcDocIDs = append(testSrcDocIDs, id)
	}

	// Copy dependencies
	depsCopy := make([]ModuleDescriptor, len(moduleDescDependencies))
	copy(depsCopy, moduleDescDependencies)

	return &moduleContext{
		project:                project,
		moduleID:               moduleID,
		moduleDescriptor:       moduleDescriptor,
		isDefaultModule:        isDefaultModule,
		srcDocContextMap:       srcDocContextMap,
		srcDocIDs:              srcDocIDs,
		testDocContextMap:      testDocContextMap,
		testSrcDocIDs:          testSrcDocIDs,
		moduleDescDependencies: depsCopy,
	}
}

// getModuleID returns the module identifier.
func (m *moduleContext) getModuleID() ModuleID {
	return m.moduleID
}

// getDescriptor returns the module descriptor.
func (m *moduleContext) getDescriptor() ModuleDescriptor {
	return m.moduleDescriptor
}

// getModuleName returns the module name.
func (m *moduleContext) getModuleName() ModuleName {
	return m.moduleDescriptor.Name()
}

// getSrcDocumentIDs returns a defensive copy of source document IDs.
func (m *moduleContext) getSrcDocumentIDs() []DocumentID {
	result := make([]DocumentID, len(m.srcDocIDs))
	copy(result, m.srcDocIDs)
	return result
}

// getTestSrcDocumentIDs returns a defensive copy of test document IDs.
func (m *moduleContext) getTestSrcDocumentIDs() []DocumentID {
	result := make([]DocumentID, len(m.testSrcDocIDs))
	copy(result, m.testSrcDocIDs)
	return result
}

// getDocumentContext returns the context for a document.
// It searches both source and test document maps.
func (m *moduleContext) getDocumentContext(documentID DocumentID) *documentContext {
	if ctx, ok := m.srcDocContextMap[documentID]; ok {
		return ctx
	}
	return m.testDocContextMap[documentID]
}

// getSrcDocContextMap returns the source document context map.
// This returns a shallow copy of the map.
func (m *moduleContext) getSrcDocContextMap() map[DocumentID]*documentContext {
	result := make(map[DocumentID]*documentContext, len(m.srcDocContextMap))
	for k, v := range m.srcDocContextMap {
		result[k] = v
	}
	return result
}

// getTestDocContextMap returns the test document context map.
// This returns a shallow copy of the map.
func (m *moduleContext) getTestDocContextMap() map[DocumentID]*documentContext {
	result := make(map[DocumentID]*documentContext, len(m.testDocContextMap))
	for k, v := range m.testDocContextMap {
		result[k] = v
	}
	return result
}

// isDefault returns true if this is the default module.
func (m *moduleContext) isDefault() bool {
	return m.isDefaultModule
}

// getProject returns the project reference.
func (m *moduleContext) getProject() Project {
	return m.project
}

// getModuleDescDependencies returns a defensive copy of module descriptor dependencies.
func (m *moduleContext) getModuleDescDependencies() []ModuleDescriptor {
	result := make([]ModuleDescriptor, len(m.moduleDescDependencies))
	copy(result, m.moduleDescDependencies)
	return result
}

// duplicate creates a copy of the context.
// The duplicated context has all document contexts duplicated as well.
// Java: ModuleContext.duplicate(Project)
func (m *moduleContext) duplicate(project Project) *moduleContext {
	// Duplicate source document contexts
	srcDocContextMap := make(map[DocumentID]*documentContext, len(m.srcDocIDs))
	for _, docID := range m.srcDocIDs {
		if docCtx := m.srcDocContextMap[docID]; docCtx != nil {
			srcDocContextMap[docID] = docCtx.duplicate()
		}
	}

	// Duplicate test document contexts
	testDocContextMap := make(map[DocumentID]*documentContext, len(m.testSrcDocIDs))
	for _, docID := range m.testSrcDocIDs {
		if docCtx := m.testDocContextMap[docID]; docCtx != nil {
			testDocContextMap[docID] = docCtx.duplicate()
		}
	}

	return newModuleContextFromMaps(
		project,
		m.moduleID,
		m.moduleDescriptor,
		m.isDefaultModule,
		srcDocContextMap,
		testDocContextMap,
		m.moduleDescDependencies,
	)
}

// containsDocument checks if the module contains the given document ID.
func (m *moduleContext) containsDocument(documentID DocumentID) bool {
	_, inSrc := m.srcDocContextMap[documentID]
	_, inTest := m.testDocContextMap[documentID]
	return inSrc || inTest
}

// isTestDocument returns true if the given document ID is a test document.
func (m *moduleContext) isTestDocument(documentID DocumentID) bool {
	_, ok := m.testDocContextMap[documentID]
	return ok
}
