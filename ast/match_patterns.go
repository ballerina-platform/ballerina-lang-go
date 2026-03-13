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

import (
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

type BLangMatchPattern interface {
	model.MatchPatternNode
	SetAcceptedType(semtypes.SemType)
}

type BLangMatchGuard BLangExpression

type (
	BLangMatchClause struct {
		bLangNodeBase
		Guard        BLangExpression
		Body         BLangBlockStmt
		Patterns     []BLangMatchPattern
		AcceptedType semtypes.SemType
	}

	bLangMatchPatternBase struct {
		bLangNodeBase
		AcceptedType semtypes.SemType
	}
	BLangConstPattern struct {
		bLangMatchPatternBase
		Expr BLangExpression
	}

	BLangWildCardMatchPattern struct {
		bLangMatchPatternBase
	}
)

var _ model.ConstPatternNode = &BLangConstPattern{}
var _ model.MatchClause = &BLangMatchClause{}
var _ BLangMatchPattern = &BLangConstPattern{}
var _ BLangMatchPattern = &BLangWildCardMatchPattern{}

var _ BLangNode = &BLangConstPattern{}
var _ BLangNode = &BLangMatchClause{}
var _ BLangNode = &BLangWildCardMatchPattern{}

func (b *BLangConstPattern) GetKind() model.NodeKind {
	// migrated from BLangConstPattern.java:53:5
	return model.NodeKind_CONST_MATCH_PATTERN
}

func (b *BLangConstPattern) GetExpression() model.ExpressionNode {
	// migrated from BLangConstPattern.java:58:5
	return b.Expr
}

func (b *BLangConstPattern) SetExpression(expression model.ExpressionNode) {
	// migrated from BLangConstPattern.java:63:5
	if expr, ok := expression.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("Expected BLangExpression")
	}
}

func (b *BLangWildCardMatchPattern) GetKind() model.NodeKind {
	return model.NodeKind_WILDCARD_MATCH_PATTERN
}

func (b *BLangMatchClause) GetKind() model.NodeKind {
	return model.NodeKind_MATCH_CLAUSE
}

func (b *BLangMatchClause) GetMatchGuard() model.MatchGuard {
	return b.Guard
}

func (b *BLangMatchClause) GetBlockStatementNode() model.BlockStatementNode {
	return &b.Body
}

func (b *BLangMatchClause) GetMatchPatterns() []model.MatchPatternNode {
	result := make([]model.MatchPatternNode, len(b.Patterns))
	for i, p := range b.Patterns {
		result[i] = p
	}
	return result
}

func (b *BLangMatchClause) SetMatchClause(node BLangMatchGuard) {
	b.Guard = node
}

func (b *BLangMatchClause) SetBlockStatementNode(node BLangBlockStmt) {
	b.Body = node
}

func (b *BLangMatchClause) SetMatchPatterns(nodes []BLangMatchPattern) {
	b.Patterns = nodes
}

func (b *BLangMatchClause) GetAcceptedType() semtypes.SemType {
	return b.AcceptedType
}

func (b *bLangMatchPatternBase) GetAcceptedType() semtypes.SemType {
	return b.AcceptedType
}

func (b *bLangMatchPatternBase) SetAcceptedType(t semtypes.SemType) {
	b.AcceptedType = t
}
