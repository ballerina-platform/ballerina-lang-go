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

type (
	BLangBindingPatternBase struct {
		bLangNodeBase
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
		BindingPattern BLangBindingPattern
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
	_ CaptureBindingPatternNode      = &BLangCaptureBindingPattern{}
	_ ErrorBindingPatternNode        = &BLangErrorBindingPattern{}
	_ ErrorMessageBindingPatternNode = &BLangErrorMessageBindingPattern{}
	_ ErrorCauseBindingPatternNode   = &BLangErrorCauseBindingPattern{}
	_ ErrorFieldBindingPatternsNode  = &BLangErrorFieldBindingPatterns{}
	_ SimpleBindingPatternNode       = &BLangSimpleBindingPattern{}
	_ NamedArgBindingPatternNode     = &BLangNamedArgBindingPattern{}
	_ RestBindingPatternNode         = &BLangRestBindingPattern{}
	_ WildCardBindingPatternNode     = &BLangWildCardBindingPattern{}
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

func (b *BLangCaptureBindingPattern) GetKind() NodeKind {
	// migrated from BLangCaptureBindingPattern.java:55:5
	return NodeKind_CAPTURE_BINDING_PATTERN
}

func (b *BLangCaptureBindingPattern) GetIdentifier() *BLangIdentifier {
	// migrated from BLangCaptureBindingPattern.java:60:5
	return &b.Identifier
}

func (b *BLangCaptureBindingPattern) SetIdentifier(identifier *BLangIdentifier) {
	// migrated from BLangCaptureBindingPattern.java:65:5
	b.Identifier = *identifier
}

func (b *BLangErrorBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorBindingPattern.java:59:5
	return NodeKind_ERROR_BINDING_PATTERN
}

func (b *BLangErrorBindingPattern) GetErrorTypeReference() UserDefinedTypeNode {
	// migrated from BLangErrorBindingPattern.java:64:5
	return b.ErrorTypeReference
}

func (b *BLangErrorBindingPattern) SetErrorTypeReference(userDefinedTypeNode UserDefinedTypeNode) {
	// migrated from BLangErrorBindingPattern.java:69:5
	if userDefinedTypeNode, ok := userDefinedTypeNode.(*BLangUserDefinedType); ok {
		b.ErrorTypeReference = userDefinedTypeNode
		return
	}
	panic("userDefinedTypeNode is not a BLangUserDefinedType")
}

func (b *BLangErrorBindingPattern) GetErrorMessageBindingPatternNode() ErrorMessageBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:74:5
	return b.ErrorMessageBindingPattern
}

func (b *BLangErrorBindingPattern) SetErrorMessageBindingPatternNode(errorMessageBindingPatternNode ErrorMessageBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:79:5
	if errorMessageBindingPatternNode, ok := errorMessageBindingPatternNode.(*BLangErrorMessageBindingPattern); ok {
		b.ErrorMessageBindingPattern = errorMessageBindingPatternNode
		return
	}
	panic("errorMessageBindingPatternNode is not a BLangErrorMessageBindingPattern")
}

func (b *BLangErrorBindingPattern) GetErrorCauseBindingPatternNode() ErrorCauseBindingPatternNode {
	// migrated from BLangErrorBindingPattern.java:84:5
	return b.ErrorCauseBindingPattern
}

func (b *BLangErrorBindingPattern) SetErrorCauseBindingPatternNode(errorCauseBindingPatternNode ErrorCauseBindingPatternNode) {
	// migrated from BLangErrorBindingPattern.java:89:5
	if errorCauseBindingPatternNode, ok := errorCauseBindingPatternNode.(*BLangErrorCauseBindingPattern); ok {
		b.ErrorCauseBindingPattern = errorCauseBindingPatternNode
		return
	}
	panic("errorCauseBindingPatternNode is not a BLangErrorCauseBindingPattern")
}

func (b *BLangErrorBindingPattern) GetErrorFieldBindingPatternsNode() ErrorFieldBindingPatternsNode {
	// migrated from BLangErrorBindingPattern.java:94:5
	return b.ErrorFieldBindingPatterns
}

func (b *BLangErrorBindingPattern) SetErrorFieldBindingPatternsNode(errorFieldBindingPatternsNode ErrorFieldBindingPatternsNode) {
	// migrated from BLangErrorBindingPattern.java:99:5
	if errorFieldBindingPatternsNode, ok := errorFieldBindingPatternsNode.(*BLangErrorFieldBindingPatterns); ok {
		b.ErrorFieldBindingPatterns = errorFieldBindingPatternsNode
		return
	}
	panic("errorFieldBindingPatternsNode is not a BLangErrorFieldBindingPatterns")
}

func (b *BLangErrorMessageBindingPattern) GetSimpleBindingPattern() SimpleBindingPatternNode {
	// migrated from BLangErrorMessageBindingPattern.java:37:5
	return b.SimpleBindingPattern
}

func (b *BLangErrorMessageBindingPattern) SetSimpleBindingPattern(simpleBindingPattern SimpleBindingPatternNode) {
	// migrated from BLangErrorMessageBindingPattern.java:42:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		b.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (b *BLangErrorMessageBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorMessageBindingPattern.java:62:5
	return NodeKind_ERROR_MESSAGE_BINDING_PATTERN
}

func (b *BLangErrorCauseBindingPattern) GetSimpleBindingPattern() SimpleBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:39:5
	return b.SimpleBindingPattern
}

func (b *BLangErrorCauseBindingPattern) SetSimpleBindingPattern(simpleBindingPattern SimpleBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:44:5
	if simpleBindingPattern, ok := simpleBindingPattern.(*BLangSimpleBindingPattern); ok {
		b.SimpleBindingPattern = simpleBindingPattern
		return
	}
	panic("simpleBindingPattern is not a BLangSimpleBindingPattern")
}

func (b *BLangErrorCauseBindingPattern) GetErrorBindingPatternNode() ErrorBindingPatternNode {
	// migrated from BLangErrorCauseBindingPattern.java:49:5
	return b.ErrorBindingPattern
}

func (b *BLangErrorCauseBindingPattern) SetErrorBindingPatternNode(errorBindingPatternNode ErrorBindingPatternNode) {
	// migrated from BLangErrorCauseBindingPattern.java:54:5
	if errorBindingPatternNode, ok := errorBindingPatternNode.(*BLangErrorBindingPattern); ok {
		b.ErrorBindingPattern = errorBindingPatternNode
		return
	}
	panic("errorBindingPatternNode is not a BLangErrorBindingPattern")
}

func (b *BLangErrorCauseBindingPattern) GetKind() NodeKind {
	// migrated from BLangErrorCauseBindingPattern.java:74:5
	return NodeKind_ERROR_CAUSE_BINDING_PATTERN
}

func (b *BLangSimpleBindingPattern) GetCaptureBindingPattern() CaptureBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:39:5
	return b.CaptureBindingPattern
}

func (b *BLangSimpleBindingPattern) SetCaptureBindingPattern(captureBindingPattern CaptureBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:44:5
	if captureBindingPattern, ok := captureBindingPattern.(*BLangCaptureBindingPattern); ok {
		b.CaptureBindingPattern = captureBindingPattern
		return
	}
	panic("captureBindingPattern is not a BLangCaptureBindingPattern")
}

func (b *BLangSimpleBindingPattern) GetWildCardBindingPattern() WildCardBindingPatternNode {
	// migrated from BLangSimpleBindingPattern.java:49:5
	return b.WildCardBindingPattern
}

func (b *BLangSimpleBindingPattern) SetWildCardBindingPattern(wildCardBindingPattern WildCardBindingPatternNode) {
	// migrated from BLangSimpleBindingPattern.java:54:5
	if wildCardBindingPatternNode, ok := wildCardBindingPattern.(*BLangWildCardBindingPattern); ok {
		b.WildCardBindingPattern = wildCardBindingPatternNode
		return
	}
	panic("wildCardBindingPatternNode is not a BLangWildCardBindingPattern")
}

func (b *BLangSimpleBindingPattern) GetKind() NodeKind {
	// migrated from BLangSimpleBindingPattern.java:74:5
	return NodeKind_SIMPLE_BINDING_PATTERN
}

func (b *BLangErrorFieldBindingPatterns) GetNamedArgMatchPatterns() []NamedArgBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:42:5
	namedArgBindingPatterns := make([]NamedArgBindingPatternNode, len(b.NamedArgBindingPatterns))
	for i, namedArgBindingPattern := range b.NamedArgBindingPatterns {
		namedArgBindingPatterns[i] = &namedArgBindingPattern
	}
	return namedArgBindingPatterns
}

func (b *BLangErrorFieldBindingPatterns) AddNamedArgBindingPattern(namedArgBindingPatternNode NamedArgBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:47:5
	if namedArgBindingPatternNode, ok := namedArgBindingPatternNode.(*BLangNamedArgBindingPattern); ok {
		b.NamedArgBindingPatterns = append(b.NamedArgBindingPatterns, *namedArgBindingPatternNode)
		return
	}
	panic("namedArgBindingPatternNode is not a BLangNamedArgBindingPattern")
}

func (b *BLangErrorFieldBindingPatterns) GetRestBindingPattern() RestBindingPatternNode {
	// migrated from BLangErrorFieldBindingPatterns.java:52:5
	return b.RestBindingPattern
}

func (b *BLangErrorFieldBindingPatterns) SetRestBindingPattern(restBindingPattern RestBindingPatternNode) {
	// migrated from BLangErrorFieldBindingPatterns.java:57:5
	if restBindingPattern, ok := restBindingPattern.(*BLangRestBindingPattern); ok {
		b.RestBindingPattern = restBindingPattern
		return
	}
	panic("restBindingPattern is not a BLangRestBindingPattern")
}

func (b *BLangErrorFieldBindingPatterns) GetKind() NodeKind {
	// migrated from BLangErrorFieldBindingPatterns.java:77:5
	return NodeKind_ERROR_FIELD_BINDING_PATTERN
}

func (b *BLangNamedArgBindingPattern) GetIdentifier() *BLangIdentifier {
	// migrated from BLangNamedArgBindingPattern.java:40:5
	return b.ArgName
}

func (b *BLangNamedArgBindingPattern) SetIdentifier(variableName *BLangIdentifier) {
	// migrated from BLangNamedArgBindingPattern.java:45:5
	b.ArgName = variableName
}

func (b *BLangNamedArgBindingPattern) GetBindingPattern() BindingPatternNode {
	// migrated from BLangNamedArgBindingPattern.java:50:5
	return b.BindingPattern
}

func (b *BLangNamedArgBindingPattern) SetBindingPattern(bindingPattern BindingPatternNode) {
	// migrated from BLangNamedArgBindingPattern.java:55:5
	b.BindingPattern = bindingPattern
}

func (b *BLangNamedArgBindingPattern) GetKind() NodeKind {
	// migrated from BLangNamedArgBindingPattern.java:75:5
	return NodeKind_NAMED_ARG_BINDING_PATTERN
}

func (b *BLangRestBindingPattern) GetIdentifier() *BLangIdentifier {
	// migrated from BLangRestBindingPattern.java:42:5
	return b.VariableName
}

func (b *BLangRestBindingPattern) SetIdentifier(variableName *BLangIdentifier) {
	// migrated from BLangRestBindingPattern.java:47:5
	b.VariableName = variableName
}

func (b *BLangRestBindingPattern) GetKind() NodeKind {
	// migrated from BLangRestBindingPattern.java:67:5
	return NodeKind_REST_BINDING_PATTERN
}

func (b *BLangWildCardBindingPattern) GetKind() NodeKind {
	// migrated from BLangWildCardBindingPattern.java:48:5
	return NodeKind_WILDCARD_BINDING_PATTERN
}

func (*BLangWildCardBindingPattern) actionOrExpression() {}
func (*BLangWildCardBindingPattern) expressionNode()     {}
