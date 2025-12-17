// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

//go:generate ../../tree-gen -config ../nodes.json -type node -template ../../compiler-tools/tree-gen/templates/node.go.tmpl -output node-gen.go

import (
	"ballerina-lang-go/parser/common"
	"ballerina-lang-go/parser/internal"
	"ballerina-lang-go/tools/diagnostics"
	"iter"
)

// This represent red nodes in the syntax tree. Red nodes satisfy fallowing properties:
//   1. Immutable
//   2. Have parent reference
//   3. Know position but not width
//   4. We build them top down
//       -- They are essentially facade nodes for green nodes.
//       -- We rebuild these per keystroke.
// All red nodes satisfy Node interface.
// We need these to be very fast to build but we rebuild the tree per keystroke so no necessarily memory efficient.
//

// TODO: Revisit how we store position information for nodes. This is storing too much information that can be computed
//  no demand and costing us memory.

type Node interface {
	Position() int
	Parent() *NonTerminalNode
	Ancestor(filter func(Node) bool) *Node
	Ancestors() []*Node
	TextRange() TextRange
	TextRangeWithMinutiae() TextRange
	Kind() common.SyntaxKind
	Location() NodeLocation
	Diagnostics() iter.Seq[Diagnostic]
	HasDiagnostics() bool
	IsMissing() bool
	SyntaxTree() *SyntaxTree
	LineRange() LineRange
	LeadingMinutiae() MinutiaeList
	TrailingMinutiae() MinutiaeList
	LeadingInvalidTokens() []Token
	TrailingInvalidTokens() []Token
	// TODO: think about how to do nodetransformer
	InternalNode() internal.STNode
	ToSourceCode() string
}

type MinutiaeList struct{}

func (m *MinutiaeList) Iterator() iter.Seq[Minutiae] {
	panic("not implemented")
}

type Minutiae struct {
	internalMinutiae internal.STMinutiae
	token            Token
	position         int
	textRange        TextRange
	lineRange        LineRange
}

type Location interface {
	LineRange() LineRange
	TextRange() TextRange
}

// TODO: get rid of this unwanted indirection (to get this you had a node get location from that)
type NodeLocation struct {
	node Node
}

func (n *NodeLocation) LineRange() LineRange {
	return n.node.LineRange()
}

func (n *NodeLocation) TextRange() TextRange {
	return n.node.TextRange()
}

type Diagnostic interface {
	Location() Location
	DiagnosticInfo() DiagnosticInfo
	Message() string
	Properties() []DiagnosticProperty[any]
}

type DiagnosticProperty[T any] interface {
	Kind() diagnostics.DiagnosticPropertyKind
	Value() T
}

type SyntaxDiagnostic struct {
	nodeDiagnostic internal.STNodeDiagnostic
	location       NodeLocation
	diagnosticInfo *DiagnosticInfo
}

func (sd *SyntaxDiagnostic) Location() Location {
	return &sd.location
}

func (sd *SyntaxDiagnostic) DiagnosticInfo() DiagnosticInfo {
	if sd.diagnosticInfo != nil {
		return *sd.diagnosticInfo
	}
	diagnosticCode := sd.nodeDiagnostic.DiagnosticCode()
	sd.diagnosticInfo = &DiagnosticInfo{code: diagnosticCode.DiagnosticId(),
		messageFormat: diagnosticCode.MessageKey(), severity: DiagnosticSeverity(diagnosticCode.Severity())}
	return *sd.diagnosticInfo
}

type DiagnosticInfo struct {
	code          string
	messageFormat string
	severity      DiagnosticSeverity
}

type DiagnosticSeverity uint8

const (
	Internal DiagnosticSeverity = iota
	Hint
	Info
	Warning
	Error
)

type NodeBase struct {
	internalNode internal.STNode
	// TODO: does this needs to be int?
	position              int
	parent                *NonTerminalNode
	syntaxTree            *SyntaxTree
	lineRange             LineRange
	textRange             TextRange
	textRangeWithMinutiae TextRange
}

func NodeFrom(internalNode internal.STNode, position int, parent *NonTerminalNode) Node {
	return &NodeBase{
		internalNode: internalNode,
		position:     position,
		parent:       parent,
	}
}

func (n *NodeBase) Kind() common.SyntaxKind {
	return n.internalNode.Kind()
}

func (n *NodeBase) Position() int {
	return n.position
}

func (n *NodeBase) Parent() *NonTerminalNode {
	return n.parent
}

func (n *NodeBase) Ancestor(filter func(Node) bool) *Node {
	if n.parent == nil {
		return nil
	}
	parent := n.parent
	for parent != nil {
		if filter(parent) {
			var result Node = parent
			return &result
		}
		parentPtr := parent.Parent()
		if parentPtr == nil {
			break
		}
		parent = parentPtr
	}
	return nil
}

func (n *NodeBase) Ancestors() []*Node {
	var ancestors []*Node
	if n.parent == nil {
		return ancestors
	}
	var parent Node = n.parent
	for parent != nil {
		ancestors = append(ancestors, &parent)
		parent = parent.Parent()
	}
	return ancestors
}

func (n *NodeBase) TextRange() TextRange {
	if n.textRange.length != 0 {
		return n.textRange
	}
	leadingMinutiaeDelta := int(n.internalNode.WidthWithLeadingMinutiae()) - int(n.internalNode.Width())
	positionWithoutLeadingMinutiae := n.position + leadingMinutiaeDelta
	n.textRange = TextRange{
		startOffset: positionWithoutLeadingMinutiae,
		endOffset:   positionWithoutLeadingMinutiae + int(n.internalNode.Width()),
		length:      int(n.internalNode.Width()),
	}
	return n.textRange
}

func (n *NodeBase) TextRangeWithMinutiae() TextRange {
	if n.textRangeWithMinutiae.length != 0 {
		return n.textRangeWithMinutiae
	}
	n.textRangeWithMinutiae = TextRange{
		startOffset: n.position,
		endOffset:   n.position + int(n.internalNode.WidthWithMinutiae()),
		length:      int(n.internalNode.WidthWithMinutiae()),
	}
	return n.textRangeWithMinutiae
}

func (n *NodeBase) Location() NodeLocation {
	return NodeLocation{node: n}
}

func (n *NodeBase) Diagnostics() iter.Seq[Diagnostic] {
	panic("Diagnostics() should be implemented by child types")
}

func (n *NonTerminalNode) Diagnostics() iter.Seq[Diagnostic] {
	return func(yield func(Diagnostic) bool) {
		if !n.internalNode.HasDiagnostics() {
			return
		}
		for _, ch := range n.Children() {
			for diagnostic := range ch.Diagnostics() {
				if !yield(diagnostic) {
					return
				}
			}
		}
		for _, diagnostic := range n.internalNode.Diagnostics() {
			if !yield(createSyntaxDiagnostic(diagnostic)) {
				return
			}
		}
	}
}

func createSyntaxDiagnostic(diagnostic internal.STNodeDiagnostic) Diagnostic {
	panic("not implemented")
}

func (n *NonTerminalNode) Children() []Node {
	panic("Children() should be implemented by child types")
}

func (n *NodeBase) HasDiagnostics() bool {
	return n.internalNode.HasDiagnostics()
}

func (n *NodeBase) IsMissing() bool {
	return n.internalNode.IsMissing()
}

func (n *NodeBase) SyntaxTree() *SyntaxTree {
	return n.populateSyntaxTree()
}

func (n *NodeBase) LineRange() LineRange {
	if n.lineRange.startLine.line != 0 || n.lineRange.endLine.line != 0 {
		return n.lineRange
	}

	_ = n.SyntaxTree()
	// TODO: implement line range calculation
	// This requires accessing the text document from the syntax tree
	return n.lineRange
}

func (n *NodeBase) LeadingMinutiae() MinutiaeList {
	panic("LeadingMinutiae() should be implemented by child types")
}

func (n *NodeBase) TrailingMinutiae() MinutiaeList {
	panic("TrailingMinutiae() should be implemented by child types")
}

func (n *NodeBase) LeadingInvalidTokens() []Token {
	panic("LeadingInvalidTokens() should be implemented by child types")
}

func (n *NodeBase) TrailingInvalidTokens() []Token {
	panic("TrailingInvalidTokens() should be implemented by child types")
}

func (n *NodeBase) InternalNode() internal.STNode {
	return n.internalNode
}

func (n *NodeBase) ToSourceCode() string {
	panic("ToSourceCode() should be implemented by child types")
}

func (n *NodeBase) populateSyntaxTree() *SyntaxTree {
	if n.syntaxTree != nil {
		return n.syntaxTree
	}

	if n.parent == nil {
		// This is a detached node. Create a new SyntaxTree with this node being the root.
		n.syntaxTree = &SyntaxTree{
			rootNode: n,
		}
	} else {
		parent := *n.parent
		n.syntaxTree = parent.SyntaxTree()
	}
	return n.syntaxTree
}

type NonTerminalNode struct {
	NodeBase
	childBuckets []Node
}

func (n *NonTerminalNode) bucketCount() int {
	return n.internalNode.BucketCount()
}

func (n *NonTerminalNode) ChildNodes() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for i := range n.childBuckets {
			if !yield(n.loadNode(i)) {
				return
			}
		}
	}
}

// FIXME: this don't fully implement ChildNodeList.loadNode but do we need to?
func (n *NonTerminalNode) loadNode(childIndex int) Node {
	index := 0
	for i := range n.internalNode.BucketCount() {
		child := n.internalNode.ChildInBucket(i)
		if !internal.IsSTNodePresent(child) {
			continue
		}
		if child.Kind() == common.LIST {
			if childIndex < index+child.BucketCount() {
				listChildIndex := childIndex - index
				return n.ChildInBucket(listChildIndex)
			}
			index += child.BucketCount()
		} else {
			if childIndex == index {
				return n.ChildInBucket(i)
			}
			index++
		}
	}
	panic("failed to load node")
}

func into[T Node](node Node) T {
	typed, ok := node.(T)
	if !ok {
		panic("failed to cast node to type")
	}
	return typed
}

func (n *NonTerminalNode) ChildInBucket(bucket int) Node {
	child := n.childBuckets[bucket]
	if child != nil {
		return child
	}
	internalChild := n.internalNode.ChildInBucket(bucket)
	if !internal.IsSTNodePresent(internalChild) {
		return nil
	}
	child = createFacade[Node](internalChild, n.position, *n)
	n.childBuckets[bucket] = child
	return child

}

type Token struct {
	NodeBase
	leadingMinutiaeList  MinutiaeList
	trailingMinutiaeList MinutiaeList
}

func (t *Token) Text() string {
	stToken, ok := t.internalNode.(internal.STToken)
	if !ok {
		panic("expected STToken")
	}
	return stToken.Text()
}

type SyntaxTree struct {
	rootNode     Node
	filePath     string
	textDocument TextDocument
	lineRange    LineRange
}

type TextDocument interface {
}

type LineRange struct {
	// In java version there is fileNmae as well I think we can get this from textDocument
	startLine LinePosition
	endLine   LinePosition
}

// TODO: int to match with java, i think a pair of u16 is enough
type LinePosition struct {
	line   int
	column int
}
type TextRange struct {
	startOffset int
	endOffset   int
	length      int
}

func createFacade[T Node](node internal.STNode, position int, parent NonTerminalNode) T {
	panic("not implemented")
}

type NodeList[T Node] struct {
	internalListNode internal.STNodeList
	nonTerminalNode  NonTerminalNode
	size             int
}

func nodeListFrom[T Node](nonTerminalNode *NonTerminalNode) NodeList[T] {
	size := nonTerminalNode.bucketCount()
	internalListNode, ok := nonTerminalNode.internalNode.(*internal.STNodeList)
	if !ok {
		panic("expected STNodeList")
	}
	return NodeList[T]{
		internalListNode: *internalListNode,
		nonTerminalNode:  *nonTerminalNode,
		size:             size,
	}
}

type DocumentMemberDeclarationNode struct {
	NonTerminalNode
}

type IdentifierToken struct {
	Token
}
