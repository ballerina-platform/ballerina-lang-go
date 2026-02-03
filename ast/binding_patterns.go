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

import "ballerina-lang-go/model"

type (
	BLangBindingPatternBase struct {
		BLangNodeBase
	}

	BLangCaptureBindingPattern struct {
		BLangBindingPatternBase
		Identifier BLangIdentifier
	}

	BLangErrorBindingPattern struct {
		BLangBindingPatternBase
		ErrorTypeReference         *BLangUserDefinedType
		ErrorMessageBindingPattern *BLangErrorMessageBindingPattern
		ErrorCauseBindingPattern   *BLangErrorCauseBindingPattern
		ErrorFieldBindingPatterns  *BLangErrorFieldBindingPatterns
	}

	BLangErrorMessageBindingPattern struct {
		BLangBindingPatternBase
		SimpleBindingPattern *BLangSimpleBindingPattern
	}
	BLangErrorCauseBindingPattern struct {
		BLangBindingPatternBase
		SimpleBindingPattern *BLangSimpleBindingPattern
		ErrorBindingPattern  *BLangErrorBindingPattern
	}

	BLangErrorFieldBindingPatterns struct {
		BLangBindingPatternBase
		NamedArgBindingPatterns []BLangNamedArgBindingPattern
		RestBindingPattern      *BLangRestBindingPattern
	}
	BLangSimpleBindingPattern struct {
		BLangBindingPatternBase
		CaptureBindingPattern  *BLangCaptureBindingPattern
		WildCardBindingPattern *BLangWildCardBindingPattern
	}

	BLangNamedArgBindingPattern struct {
		BLangBindingPatternBase
		ArgName        *BLangIdentifier
		BindingPattern model.BindingPatternNode
	}

	BLangRestBindingPattern struct {
		BLangBindingPatternBase
		VariableName *BLangIdentifier
	}

	BLangWildCardBindingPattern struct {
		BLangBindingPatternBase
	}
)

var (
	_ model.CaptureBindingPatternNode      = &BLangCaptureBindingPattern{}
	_ model.ErrorBindingPatternNode        = &BLangErrorBindingPattern{}
	_ model.ErrorMessageBindingPatternNode = &BLangErrorMessageBindingPattern{}
	_ model.ErrorCauseBindingPatternNode   = &BLangErrorCauseBindingPattern{}
	_ model.ErrorFieldBindingPatternsNode  = &BLangErrorFieldBindingPatterns{}
	_ model.SimpleBindingPatternNode       = &BLangSimpleBindingPattern{}
	_ model.NamedArgBindingPatternNode     = &BLangNamedArgBindingPattern{}
	_ model.RestBindingPatternNode         = &BLangRestBindingPattern{}
	_ model.WildCardBindingPatternNode     = &BLangWildCardBindingPattern{}
)

var (
	_ BLangNode = &BLangCaptureBindingPattern{}
	_ BLangNode = &BLangErrorBindingPattern{}
	_ BLangNode = &BLangErrorMessageBindingPattern{}
	_ BLangNode = &BLangErrorCauseBindingPattern{}
	_ BLangNode = &BLangErrorFieldBindingPatterns{}
	_ BLangNode = &BLangSimpleBindingPattern{}
	_ BLangNode = &BLangNamedArgBindingPattern{}
	_ BLangNode = &BLangRestBindingPattern{}
	_ BLangNode = &BLangWildCardBindingPattern{}
)

func (this *BLangCaptureBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangCaptureBindingPattern.java:55:5
	return model.NodeKind_CAPTURE_BINDING_PATTERN
}

func (this *BLangCaptureBindingPattern) GetIdentifier() model.IdentifierNode {
	// migrated from BLangCaptureBindingPattern.java:60:5
	return &this.Identifier
}

func (this *BLangCaptureBindingPattern) SetIdentifier(identifier model.IdentifierNode) {
	// migrated from BLangCaptureBindingPattern.java:65:5
	if id, ok := identifier.(*BLangIdentifier); ok {
		this.Identifier = *id
		return
	}
	panic("identifier is not a BLangIdentifier")
}

func (this *BLangErrorBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangErrorBindingPattern.java:59:5
	return model.NodeKind_ERROR_BINDING_PATTERN
}

func (this *BLangErrorBindingPattern) GetErrorTypeReference() model.UserDefinedTypeNode {
	// migrated from BLangErrorBindingPattern.java:64:5
	return this.ErrorTypeReference
}

func (this *BLangErrorBindingPattern) SetErrorTypeReference(userDefinedTypeNode model.UserDefinedTypeNode) {
	// migrated from BLangErrorBindingPattern.java:69:5
	if userDefinedTypeNode, ok := userDefinedTypeNode.(*BLangUserDefinedType); ok {
		this.ErrorTypeReference = userDefinedTypeNode
		return
	}
	panic("userDefinedTypeNode is not a BLangUserDefinedType")
}

func (this *BLangErrorBindingPattern) GetErrorMessageBindingPatternNode() model.ErrorMessageBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:74:5
	return this.ErrorMessageBindingPattern
}

func (this *BLangErrorBindingPattern) SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode model.ErrorMessageBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:79:5
	if errorMessageBindingPatternNode, ok := errorMessageBindingPatternNode.(*BLangErrorMessageBindingPattern); ok {
		this.ErrorMessageBindingPattern = errorMessageBindingPatternNode
		return
	}
	panic("errorMessageBindingPatternNode is not a BLangErrorMessageBindingPattern")
}

func (this *BLangErrorBindingPattern) GetErrorCauseBindingPatternNode() model.ErrorCauseBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:84:5
	return this.ErrorCauseBindingPattern
}

func (this *BLangErrorBindingPattern) SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode model.ErrorCauseBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:89:5
	if errorCauseBindingPatternNode, ok := errorCauseBindingPatternNode.(*BLangErrorCauseBindingPattern); ok {
		this.ErrorCauseBindingPattern = errorCauseBindingPatternNode
		return
	}
	panic("errorCauseBindingPatternNode is not a BLangErrorCauseBindingPattern")
}

func (this *BLangErrorBindingPattern) GetErrorFieldBindingPatternsNode() model.ErrorFieldBindingPatternsNode {
	// migrated from BLangErrorBindingPattern.java:94:5
	return this.ErrorFieldBindingPatterns
}

func (this *BLangErrorBindingPattern) SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode model.ErrorFieldBindingPatternsNode) {
	// migrated from BLangErrorBindingPattern.java:99:5
	if errorFieldBindingPatternsNode, ok := errorFieldBindingPatternsNode.(*BLangErrorFieldBindingPatterns); ok {
		this.ErrorFieldBindingPatterns = errorFieldBindingPatternsNode
		return
	}
	panic("errorFieldBindingPatternsNode is not a BLangErrorFieldBindingPatterns")
}

func (this *BLangErrorMessageBindingPattern) GetSimpleBindingPattern() model.SimpleBindingPatternNode {
	// migrated from BLangErrorMessageBindingPattern.java:37:5
	return this.SimpleBindingPattern
}

func (this *BLangErrorMessageBindingPattern) SetSimpleBindingPattern(simpleBindingPattern model.SimpleBindingPatternNode) {
	// migrated from BLangErrorMessageBindingPattern.java:42:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		this.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (this *BLangErrorMessageBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangErrorMessageBindingPattern.java:62:5
	return model.NodeKind_ERROR_MESSAGE_BINDING_PATTERN
}

func (this *BLangErrorCauseBindingPattern) GetSimpleBindingPattern() model.SimpleBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:39:5
	return this.SimpleBindingPattern
}

func (this *BLangErrorCauseBindingPattern) SetSimpleBindingPattern(simpleBindingPattern model.SimpleBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:44:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		this.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (this *BLangErrorCauseBindingPattern) GetErrorBindingPatternNode() model.ErrorBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:49:5
	return this.ErrorBindingPattern
}

func (this *BLangErrorCauseBindingPattern) SetErrorBindingPatternNode(errorBindingPatternNode model.ErrorBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:54:5
	if errorBindingPatternNode, ok := errorBindingPatternNode.(*BLangErrorBindingPattern); ok {
		this.ErrorBindingPattern = errorBindingPatternNode
		return
	}
	panic("errorBindingPatternNode is not a BLangErrorBindingPattern")
}

func (this *BLangErrorCauseBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangErrorCauseBindingPattern.java:74:5
	return model.NodeKind_ERROR_CAUSE_BINDING_PATTERN
}

func (this *BLangSimpleBindingPattern) GetCaptureBindingPattern() model.CaptureBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:39:5
	return this.CaptureBindingPattern
}

func (this *BLangSimpleBindingPattern) SetCaptureBindingPattern(captureBindingPattern model.CaptureBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:44:5
	if captureBindingPattern, ok := captureBindingPattern.(*BLangCaptureBindingPattern); ok {
		this.CaptureBindingPattern = captureBindingPattern
		return
	}
	panic("captureBindingPattern is not a BLangCaptureBindingPattern")
}

func (this *BLangSimpleBindingPattern) GetWildCardBindingPattern() model.WildCardBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:49:5
	return this.WildCardBindingPattern
}

func (this *BLangSimpleBindingPattern) SetWildCardBindingPattern(wildCardBindingPattern model.WildCardBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:54:5
	if wildCardBindingPatternNode, ok := wildCardBindingPattern.(*BLangWildCardBindingPattern); ok {
		this.WildCardBindingPattern = wildCardBindingPatternNode
		return
	}
	panic("wildCardBindingPatternNode is not a BLangWildCardBindingPattern")
}

func (this *BLangSimpleBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangSimpleBindingPattern.java:74:5
	return model.NodeKind_SIMPLE_BINDING_PATTERN
}

func (this *BLangErrorFieldBindingPatterns) GetNamedArgMatchPatterns() []model.NamedArgBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:42:5
	namedArgBindingPatterns := make([]model.NamedArgBindingPatternNode, len(this.NamedArgBindingPatterns))
	for i, namedArgBindingPattern := range this.NamedArgBindingPatterns {
		namedArgBindingPatterns[i] = &namedArgBindingPattern
	}
	return namedArgBindingPatterns
}

func (this *BLangErrorFieldBindingPatterns) AddNamedArgBindingPattern(namedArgBindingPatternNode model.NamedArgBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:47:5
	if namedArgBindingPatternNode, ok := namedArgBindingPatternNode.(*BLangNamedArgBindingPattern); ok {
		this.NamedArgBindingPatterns = append(this.NamedArgBindingPatterns, *namedArgBindingPatternNode)
		return
	}
	panic("namedArgBindingPatternNode is not a BLangNamedArgBindingPattern")
}

func (this *BLangErrorFieldBindingPatterns) GetRestBindingPattern() model.RestBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:52:5
	return this.RestBindingPattern
}

func (this *BLangErrorFieldBindingPatterns) SetRestBindingPattern(restBindingPattern model.RestBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:57:5
	if restBindingPattern, ok := restBindingPattern.(*BLangRestBindingPattern); ok {
		this.RestBindingPattern = restBindingPattern
		return
	}
	panic("restBindingPattern is not a BLangRestBindingPattern")
}

func (this *BLangErrorFieldBindingPatterns) GetKind() model.NodeKind {
	// migrated from BLangErrorFieldBindingPatterns.java:77:5
	return model.NodeKind_ERROR_FIELD_BINDING_PATTERN
}

func (this *BLangNamedArgBindingPattern) GetIdentifier() model.IdentifierNode {
	// migrated from BLangNamedArgBindingPattern.java:40:5
	return this.ArgName
}

func (this *BLangNamedArgBindingPattern) SetIdentifier(variableName model.IdentifierNode) {
	// migrated from BLangNamedArgBindingPattern.java:45:5
	if variableName, ok := variableName.(*BLangIdentifier); ok {
		this.ArgName = variableName
		return
	}
	panic("variableName is not a BLangIdentifier")
}

func (this *BLangNamedArgBindingPattern) GetBindingPattern() model.BindingPatternNode {
	// migrated from BLangNamedArgBindingPattern.java:50:5
	return this.BindingPattern
}

func (this *BLangNamedArgBindingPattern) SetBindingPattern(bindingPattern model.BindingPatternNode) {
	// migrated from BLangNamedArgBindingPattern.java:55:5
	this.BindingPattern = bindingPattern
}

func (this *BLangNamedArgBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangNamedArgBindingPattern.java:75:5
	return model.NodeKind_NAMED_ARG_BINDING_PATTERN
}

func (this *BLangRestBindingPattern) GetIdentifier() model.IdentifierNode {
	// migrated from BLangRestBindingPattern.java:42:5
	return this.VariableName
}

func (this *BLangRestBindingPattern) SetIdentifier(variableName model.IdentifierNode) {
	// migrated from BLangRestBindingPattern.java:47:5
	if variableName, ok := variableName.(*BLangIdentifier); ok {
		this.VariableName = variableName
		return
	}
	panic("variableName is not a BLangIdentifier")
}

func (this *BLangRestBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangRestBindingPattern.java:67:5
	return model.NodeKind_REST_BINDING_PATTERN
}

func (this *BLangWildCardBindingPattern) GetKind() model.NodeKind {
	// migrated from BLangWildCardBindingPattern.java:48:5
	return model.NodeKind_WILDCARD_BINDING_PATTERN
}

func (this *BLangWildCardBindingPattern) SetTypeCheckedType(ty BType) {
	panic("not implemented")
}
