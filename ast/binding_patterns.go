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

package ast

type CaptureBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
}

type WildCardBindingPatternNode = Node
type BindingPatternNode = Node

type SimpleBindingPatternNode interface {
	Node
	GetCaptureBindingPattern() CaptureBindingPatternNode
	SetCaptureBindingPattern(captureBindingPatternNode CaptureBindingPatternNode)
	GetWildCardBindingPattern() WildCardBindingPatternNode
	SetWildCardBindingPattern(wildCardBindingPatternNode WildCardBindingPatternNode)
}

type ErrorMessageBindingPatternNode interface {
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
}

type ErrorBindingPatternNode interface {
	Node
	GetErrorTypeReference() UserDefinedTypeNode
	SetErrorTypeReference(userDefinedTypeNode UserDefinedTypeNode)
	GetErrorMessageBindingPatternNode() ErrorMessageBindingPatternNode
	SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode ErrorMessageBindingPatternNode)
	GetErrorCauseBindingPatternNode() ErrorCauseBindingPatternNode
	SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode ErrorCauseBindingPatternNode)
	GetErrorFieldBindingPatternsNode() ErrorFieldBindingPatternsNode
	SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode ErrorFieldBindingPatternsNode)
}

type ErrorCauseBindingPatternNode interface {
	Node
	GetSimpleBindingPattern() SimpleBindingPatternNode
	SetSimpleBindingPattern(simpleBindingPatternNode SimpleBindingPatternNode)
	GetErrorBindingPatternNode() ErrorBindingPatternNode
	SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode)
}

type ErrorFieldBindingPatternsNode interface {
	Node
	GetNamedArgMatchPatterns() []NamedArgBindingPatternNode
	AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode)
	GetRestBindingPattern() RestBindingPatternNode
	SetRestBindingPattern(restBindingPattern RestBindingPatternNode)
}

type NamedArgBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
	GetBindingPattern() BindingPatternNode
	SetBindingPattern(bindingPattern BindingPatternNode)
}

type RestBindingPatternNode interface {
	Node
	GetIdentifier() IdentifierNode
	SetIdentifier(identifier IdentifierNode)
}

type (
	BLangBindingPattern struct {
		BLangNodeBase
		DeclaredVars map[string]BVarSymbol
	}

	BLangCaptureBindingPattern struct {
		BLangBindingPattern
		Identifier BLangIdentifier
		Symbol     BVarSymbol
	}

	BLangErrorBindingPattern struct {
		BLangBindingPattern
		ErrorTypeReference         *BLangUserDefinedType
		ErrorMessageBindingPattern *BLangErrorMessageBindingPattern
		ErrorCauseBindingPattern   *BLangErrorCauseBindingPattern
		ErrorFieldBindingPatterns  *BLangErrorFieldBindingPatterns
	}

	BLangErrorMessageBindingPattern struct {
		BLangBindingPattern
		SimpleBindingPattern *BLangSimpleBindingPattern
	}
	BLangErrorCauseBindingPattern struct {
		BLangBindingPattern
		SimpleBindingPattern *BLangSimpleBindingPattern
		ErrorBindingPattern  *BLangErrorBindingPattern
	}

	BLangErrorFieldBindingPatterns struct {
		BLangBindingPattern
		NamedArgBindingPatterns []BLangNamedArgBindingPattern
		RestBindingPattern      *BLangRestBindingPattern
	}
	BLangSimpleBindingPattern struct {
		BLangBindingPattern
		CaptureBindingPattern  *BLangCaptureBindingPattern
		WildCardBindingPattern *BLangWildCardBindingPattern
	}

	BLangNamedArgBindingPattern struct {
		BLangBindingPattern
		ArgName        *BLangIdentifier
		BindingPattern BindingPatternNode
	}

	BLangRestBindingPattern struct {
		BLangBindingPattern
		VariableName *BLangIdentifier
		Symbol       *BVarSymbol
	}

	BLangWildCardBindingPattern struct {
		BLangBindingPattern
	}
)

var _ CaptureBindingPatternNode = &BLangCaptureBindingPattern{}
var _ ErrorBindingPatternNode = &BLangErrorBindingPattern{}
var _ ErrorMessageBindingPatternNode = &BLangErrorMessageBindingPattern{}
var _ ErrorCauseBindingPatternNode = &BLangErrorCauseBindingPattern{}
var _ ErrorFieldBindingPatternsNode = &BLangErrorFieldBindingPatterns{}
var _ SimpleBindingPatternNode = &BLangSimpleBindingPattern{}
var _ NamedArgBindingPatternNode = &BLangNamedArgBindingPattern{}
var _ RestBindingPatternNode = &BLangRestBindingPattern{}
var _ WildCardBindingPatternNode = &BLangWildCardBindingPattern{}

var _ BLangNode = &BLangBindingPattern{}
var _ BLangNode = &BLangCaptureBindingPattern{}
var _ BLangNode = &BLangErrorBindingPattern{}
var _ BLangNode = &BLangErrorMessageBindingPattern{}
var _ BLangNode = &BLangErrorCauseBindingPattern{}
var _ BLangNode = &BLangErrorFieldBindingPatterns{}
var _ BLangNode = &BLangSimpleBindingPattern{}
var _ BLangNode = &BLangNamedArgBindingPattern{}
var _ BLangNode = &BLangRestBindingPattern{}
var _ BLangNode = &BLangWildCardBindingPattern{}

func (this *BLangCaptureBindingPattern) GetKind() NodeKind {
	// migrated from BLangCaptureBindingPattern.java:55:5
	return NodeKind_CAPTURE_BINDING_PATTERN
}

func (this *BLangCaptureBindingPattern) GetIdentifier() IdentifierNode {
	// migrated from BLangCaptureBindingPattern.java:60:5
	return &this.Identifier
}

func (this *BLangCaptureBindingPattern) SetIdentifier(identifier IdentifierNode) {
	// migrated from BLangCaptureBindingPattern.java:65:5
	if id, ok := identifier.(*BLangIdentifier); ok {
		this.Identifier = *id
		return
	}
	panic("identifier is not a BLangIdentifier")
}

func (this *BLangErrorBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorBindingPattern.java:59:5
	return NodeKind_ERROR_BINDING_PATTERN
}

func (this *BLangErrorBindingPattern) GetErrorTypeReference() UserDefinedTypeNode {
	// migrated from BLangErrorBindingPattern.java:64:5
	return this.ErrorTypeReference
}

func (this *BLangErrorBindingPattern) SetErrorTypeReference(userDefinedTypeNode UserDefinedTypeNode) {
	// migrated from BLangErrorBindingPattern.java:69:5
	if userDefinedTypeNode, ok := userDefinedTypeNode.(*BLangUserDefinedType); ok {
		this.ErrorTypeReference = userDefinedTypeNode
		return
	}
	panic("userDefinedTypeNode is not a BLangUserDefinedType")
}

func (this *BLangErrorBindingPattern) GetErrorMessageBindingPatternNode() ErrorMessageBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:74:5
	return this.ErrorMessageBindingPattern
}

func (this *BLangErrorBindingPattern) SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode ErrorMessageBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:79:5
	if errorMessageBindingPatternNode, ok := errorMessageBindingPatternNode.(*BLangErrorMessageBindingPattern); ok {
		this.ErrorMessageBindingPattern = errorMessageBindingPatternNode
		return
	}
	panic("errorMessageBindingPatternNode is not a BLangErrorMessageBindingPattern")
}

func (this *BLangErrorBindingPattern) GetErrorCauseBindingPatternNode() ErrorCauseBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:84:5
	return this.ErrorCauseBindingPattern
}

func (this *BLangErrorBindingPattern) SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode ErrorCauseBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:89:5
	if errorCauseBindingPatternNode, ok := errorCauseBindingPatternNode.(*BLangErrorCauseBindingPattern); ok {
		this.ErrorCauseBindingPattern = errorCauseBindingPatternNode
		return
	}
	panic("errorCauseBindingPatternNode is not a BLangErrorCauseBindingPattern")
}

func (this *BLangErrorBindingPattern) GetErrorFieldBindingPatternsNode() ErrorFieldBindingPatternsNode {
	// migrated from BLangErrorBindingPattern.java:94:5
	return this.ErrorFieldBindingPatterns
}

func (this *BLangErrorBindingPattern) SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode ErrorFieldBindingPatternsNode) {
	// migrated from BLangErrorBindingPattern.java:99:5
	if errorFieldBindingPatternsNode, ok := errorFieldBindingPatternsNode.(*BLangErrorFieldBindingPatterns); ok {
		this.ErrorFieldBindingPatterns = errorFieldBindingPatternsNode
		return
	}
	panic("errorFieldBindingPatternsNode is not a BLangErrorFieldBindingPatterns")
}

func (this *BLangErrorMessageBindingPattern) GetSimpleBindingPattern() SimpleBindingPatternNode {
	// migrated from BLangErrorMessageBindingPattern.java:37:5
	return this.SimpleBindingPattern
}

func (this *BLangErrorMessageBindingPattern) SetSimpleBindingPattern(simpleBindingPattern SimpleBindingPatternNode) {
	// migrated from BLangErrorMessageBindingPattern.java:42:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		this.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (this *BLangErrorMessageBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorMessageBindingPattern.java:62:5
	return NodeKind_ERROR_MESSAGE_BINDING_PATTERN
}

func (this *BLangErrorCauseBindingPattern) GetSimpleBindingPattern() SimpleBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:39:5
	return this.SimpleBindingPattern
}

func (this *BLangErrorCauseBindingPattern) SetSimpleBindingPattern(simpleBindingPattern SimpleBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:44:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		this.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (this *BLangErrorCauseBindingPattern) GetErrorBindingPatternNode() ErrorBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:49:5
	return this.ErrorBindingPattern
}

func (this *BLangErrorCauseBindingPattern) SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:54:5
	if errorBindingPatternNode, ok := errorBindingPatternNode.(*BLangErrorBindingPattern); ok {
		this.ErrorBindingPattern = errorBindingPatternNode
		return
	}
	panic("errorBindingPatternNode is not a BLangErrorBindingPattern")
}

func (this *BLangErrorCauseBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorCauseBindingPattern.java:74:5
	return NodeKind_ERROR_CAUSE_BINDING_PATTERN
}

func (this *BLangSimpleBindingPattern) GetCaptureBindingPattern() CaptureBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:39:5
	return this.CaptureBindingPattern
}

func (this *BLangSimpleBindingPattern) SetCaptureBindingPattern(captureBindingPattern CaptureBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:44:5
	if captureBindingPattern, ok := captureBindingPattern.(*BLangCaptureBindingPattern); ok {
		this.CaptureBindingPattern = captureBindingPattern
		return
	}
	panic("captureBindingPattern is not a BLangCaptureBindingPattern")
}

func (this *BLangSimpleBindingPattern) GetWildCardBindingPattern() WildCardBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:49:5
	return this.WildCardBindingPattern
}

func (this *BLangSimpleBindingPattern) SetWildCardBindingPattern(wildCardBindingPattern WildCardBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:54:5
	if wildCardBindingPatternNode, ok := wildCardBindingPattern.(*BLangWildCardBindingPattern); ok {
		this.WildCardBindingPattern = wildCardBindingPatternNode
		return
	}
	panic("wildCardBindingPatternNode is not a BLangWildCardBindingPattern")
}

func (this *BLangSimpleBindingPattern) GetKind() NodeKind {
	// migrated from BLangSimpleBindingPattern.java:74:5
	return NodeKind_SIMPLE_BINDING_PATTERN
}

func (this *BLangErrorFieldBindingPatterns) GetNamedArgMatchPatterns() []NamedArgBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:42:5
	namedArgBindingPatterns := make([]NamedArgBindingPatternNode, len(this.NamedArgBindingPatterns))
	for i, namedArgBindingPattern := range this.NamedArgBindingPatterns {
		namedArgBindingPatterns[i] = &namedArgBindingPattern
	}
	return namedArgBindingPatterns
}

func (this *BLangErrorFieldBindingPatterns) AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:47:5
	if namedArgBindingPatternNode, ok := namedArgBindingPatternNode.(*BLangNamedArgBindingPattern); ok {
		this.NamedArgBindingPatterns = append(this.NamedArgBindingPatterns, *namedArgBindingPatternNode)
		return
	}
	panic("namedArgBindingPatternNode is not a BLangNamedArgBindingPattern")
}

func (this *BLangErrorFieldBindingPatterns) GetRestBindingPattern() RestBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:52:5
	return this.RestBindingPattern
}

func (this *BLangErrorFieldBindingPatterns) SetRestBindingPattern(restBindingPattern RestBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:57:5
	if restBindingPattern, ok := restBindingPattern.(*BLangRestBindingPattern); ok {
		this.RestBindingPattern = restBindingPattern
		return
	}
	panic("restBindingPattern is not a BLangRestBindingPattern")
}

func (this *BLangErrorFieldBindingPatterns) GetKind() NodeKind {
	// migrated from BLangErrorFieldBindingPatterns.java:77:5
	return NodeKind_ERROR_FIELD_BINDING_PATTERN
}

func (this *BLangNamedArgBindingPattern) GetIdentifier() IdentifierNode {
	// migrated from BLangNamedArgBindingPattern.java:40:5
	return this.ArgName
}

func (this *BLangNamedArgBindingPattern) SetIdentifier(variableName IdentifierNode) {
	// migrated from BLangNamedArgBindingPattern.java:45:5
	if variableName, ok := variableName.(*BLangIdentifier); ok {
		this.ArgName = variableName
		return
	}
	panic("variableName is not a BLangIdentifier")
}

func (this *BLangNamedArgBindingPattern) GetBindingPattern() BindingPatternNode {
	// migrated from BLangNamedArgBindingPattern.java:50:5
	return this.BindingPattern
}

func (this *BLangNamedArgBindingPattern) SetBindingPattern(bindingPattern BindingPatternNode) {
	// migrated from BLangNamedArgBindingPattern.java:55:5
	this.BindingPattern = bindingPattern
}

func (this *BLangNamedArgBindingPattern) GetKind() NodeKind {
	// migrated from BLangNamedArgBindingPattern.java:75:5
	return NodeKind_NAMED_ARG_BINDING_PATTERN
}

func (this *BLangRestBindingPattern) GetIdentifier() IdentifierNode {
	// migrated from BLangRestBindingPattern.java:42:5
	return this.VariableName
}

func (this *BLangRestBindingPattern) SetIdentifier(variableName IdentifierNode) {
	// migrated from BLangRestBindingPattern.java:47:5
	if variableName, ok := variableName.(*BLangIdentifier); ok {
		this.VariableName = variableName
		return
	}
	panic("variableName is not a BLangIdentifier")
}

func (this *BLangRestBindingPattern) GetKind() NodeKind {
	// migrated from BLangRestBindingPattern.java:67:5
	return NodeKind_REST_BINDING_PATTERN
}

func (this *BLangWildCardBindingPattern) GetKind() NodeKind {
	// migrated from BLangWildCardBindingPattern.java:48:5
	return NodeKind_WILDCARD_BINDING_PATTERN
}
