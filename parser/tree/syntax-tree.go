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

package tree

import (
	"ballerina-lang-go/parser/common"
	"iter"
)

type SyntaxTree struct {
	rootNode     Node
	filePath     string
	textDocument TextDocument
	// lineRange    LineRange
}

func NewSyntaxTreeFromNodeTextDocumentStringBool(rootNode Node, textDocument TextDocument, filePath string, clone bool) SyntaxTree {
	this := SyntaxTree{}
	this.rootNode = modifyWithSyntaxTree(rootNode, clone, &this)
	this.textDocument = textDocument
	this.filePath = filePath
	return this
}

func (this *SyntaxTree) ContainsModulePart() bool {
	// migrated from SyntaxTree.java:76:5
	return this.rootNode.Kind() == common.MODULE_PART
}

func (this *SyntaxTree) FilePath() string {
	// migrated from SyntaxTree.java:85:5
	return this.filePath
}

func (this *SyntaxTree) ModifyWith(rootNode Node) SyntaxTree {
	// migrated from SyntaxTree.java:91:5
	panic("not implemented")
}

func (this *SyntaxTree) ReplaceNode(target Node, replacement Node) SyntaxTree {
	// migrated from SyntaxTree.java:95:5
	panic("not implemented")
}

func (this *SyntaxTree) Diagnostics() iter.Seq[Diagnostic] {
	// migrated from SyntaxTree.java:105:5
	return this.rootNode.Diagnostics()
}

func (this *SyntaxTree) HasDiagnostics() bool {
	// migrated from SyntaxTree.java:109:5
	return this.rootNode.HasDiagnostics()
}

func (this *SyntaxTree) ToSourceCode() string {
	// migrated from SyntaxTree.java:123:5
	return this.rootNode.ToSourceCode()
}

func modifyWithSyntaxTree[T Node](node T, clone bool, syntaxTree *SyntaxTree) T {
	var clonedNode T
	if clone {
		clonedNode = node.InternalNode().CreateFacade(0, nil).(T)
	} else {
		clonedNode = node
	}
	clonedNode.SetSyntaxTree(syntaxTree)
	return clonedNode
}
