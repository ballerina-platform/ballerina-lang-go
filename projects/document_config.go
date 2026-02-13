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

import "sync"

// DocumentConfig represents configuration for a source document (.bal file).
// It provides access to document metadata and content, supporting both eager
// and lazy content loading strategies.
// Java source: io.ballerina.projects.DocumentConfig
type DocumentConfig interface {
	// DocumentID returns the unique identifier for this document.
	DocumentID() DocumentID

	// Name returns the document name (filename).
	Name() string

	// Content returns the document content.
	Content() string
}

// eagerDocumentConfig holds document content in memory.
type eagerDocumentConfig struct {
	documentID DocumentID
	name       string
	content    string
}

// NewDocumentConfig creates a DocumentConfig with eager-loaded content.
// Java equivalent: DocumentConfig.from(DocumentId, String content, String name)
func NewDocumentConfig(documentID DocumentID, name string, content string) DocumentConfig {
	return &eagerDocumentConfig{
		documentID: documentID,
		name:       name,
		content:    content,
	}
}

func (d *eagerDocumentConfig) DocumentID() DocumentID {
	return d.documentID
}

func (d *eagerDocumentConfig) Name() string {
	return d.name
}

func (d *eagerDocumentConfig) Content() string {
	return d.content
}

// lazyDocumentConfig loads document content on-demand.
type lazyDocumentConfig struct {
	documentID      DocumentID
	name            string
	contentSupplier func() string
	content         string
	once            sync.Once
}

// NewLazyDocumentConfig creates a DocumentConfig with lazy-loaded content.
// The contentSupplier function is called once when Content() is first accessed.
// Java equivalent: DocumentConfig.from(DocumentId, Supplier<String> content, String name)
func NewLazyDocumentConfig(documentID DocumentID, name string, contentSupplier func() string) DocumentConfig {
	return &lazyDocumentConfig{
		documentID:      documentID,
		name:            name,
		contentSupplier: contentSupplier,
	}
}

func (d *lazyDocumentConfig) DocumentID() DocumentID {
	return d.documentID
}

func (d *lazyDocumentConfig) Name() string {
	return d.name
}

func (d *lazyDocumentConfig) Content() string {
	d.once.Do(func() {
		if d.contentSupplier != nil {
			d.content = d.contentSupplier()
		}
	})
	return d.content
}
