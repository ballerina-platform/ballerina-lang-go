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
	"ballerina-lang-go/common"
	"ballerina-lang-go/model"
)

type TypeParamEntry struct {
	TypeParam BType
	BoundType BType
}

type (
	BLangInputClause struct {
		bLangNodeBase
		Collection             BLangExpression
		VariableDefinitionNode model.VariableDefinitionNode
		IsDeclaredWithVarFlag  bool
	}
	BLangFromClause struct {
		BLangInputClause
	}
	BLangLetClause struct {
		bLangNodeBase
		LetVarDeclarations []model.VariableDefinitionNode
	}
	BLangWhereClause struct {
		bLangNodeBase
		Expression BLangExpression
	}
	BLangSelectClause struct {
		bLangNodeBase
		Expression BLangExpression
	}
	BLangCollectClause struct {
		bLangNodeBase
		Expression      model.ExpressionNode
		NonGroupingKeys common.Set[string]
	}
	BLangDoClause struct {
		bLangNodeBase
		Body *BLangBlockStmt
	}
	BLangOnFailClause struct {
		bLangNodeBase
		Body                   *BLangBlockStmt
		VariableDefinitionNode model.VariableDefinitionNode
		VarType                BType
		BodyContainsFail       bool
		IsInternal             bool
		isDeclaredWithVar      bool
	}
)

var (
	_ model.FromClauseNode    = &BLangFromClause{}
	_ model.Node              = &BLangLetClause{}
	_ model.Node              = &BLangWhereClause{}
	_ model.SelectClauseNode  = &BLangSelectClause{}
	_ model.CollectClauseNode = &BLangCollectClause{}
	_ model.DoClauseNode      = &BLangDoClause{}
	_ model.OnFailClauseNode  = &BLangOnFailClause{}
)

var (
	_ BLangNode = &BLangFromClause{}
	_ BLangNode = &BLangLetClause{}
	_ BLangNode = &BLangWhereClause{}
	_ BLangNode = &BLangSelectClause{}
	_ BLangNode = &BLangCollectClause{}
	_ BLangNode = &BLangDoClause{}
	_ BLangNode = &BLangOnFailClause{}
)

func (b *BLangFromClause) GetKind() model.NodeKind {
	return model.NodeKind_FROM
}

func (b *BLangFromClause) GetCollection() model.ExpressionNode {
	return b.Collection
}

func (b *BLangFromClause) SetCollection(collection model.ExpressionNode) {
	if exp, ok := collection.(BLangExpression); ok {
		b.Collection = exp
		return
	}
	panic("collection is not a BLangExpression")
}

func (b *BLangFromClause) GetVariableDefinitionNode() model.VariableDefinitionNode {
	return b.VariableDefinitionNode
}

func (b *BLangFromClause) SetVariableDefinitionNode(variableDefinitionNode model.VariableDefinitionNode) {
	b.VariableDefinitionNode = variableDefinitionNode
}

func (b *BLangFromClause) IsDeclaredWithVar() bool {
	return b.IsDeclaredWithVarFlag
}

func (b *BLangLetClause) GetKind() model.NodeKind {
	return model.NodeKind_LET_CLAUSE
}

func (b *BLangWhereClause) GetKind() model.NodeKind {
	return model.NodeKind_WHERE
}

func (b *BLangSelectClause) GetKind() model.NodeKind {
	return model.NodeKind_SELECT
}

func (b *BLangSelectClause) GetExpression() model.ExpressionNode {
	return b.Expression
}

func (b *BLangSelectClause) SetExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		b.Expression = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (b *BLangCollectClause) GetKind() model.NodeKind {
	// migrated from BLangCollectClause.java:48:5
	return model.NodeKind_COLLECT
}

func (b *BLangCollectClause) GetExpression() model.ExpressionNode {
	// migrated from BLangCollectClause.java:68:5
	return b.Expression
}

func (b *BLangCollectClause) SetExpression(expression model.ExpressionNode) {
	// migrated from BLangCollectClause.java:73:5
	if exp, ok := expression.(BLangExpression); ok {
		b.Expression = exp
	} else {
		panic("Expected BLangExpression")
	}
}

func (b *BLangDoClause) GetBody() model.BlockStatementNode {
	// migrated from BLangDoClause.java:46:5
	return b.Body
}

func (b *BLangDoClause) SetBody(body model.BlockStatementNode) {
	// migrated from BLangDoClause.java:51:5
	if body, ok := body.(*BLangBlockStmt); ok {
		b.Body = body
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangDoClause) GetKind() model.NodeKind {
	// migrated from BLangDoClause.java:66:5
	return model.NodeKind_DO
}

func (b *BLangOnFailClause) SetDeclaredWithVar() {
	// migrated from BLangOnFailClause.java:53:5
	b.isDeclaredWithVar = true
}

func (b *BLangOnFailClause) IsDeclaredWithVar() bool {
	// migrated from BLangOnFailClause.java:58:5
	return b.isDeclaredWithVar
}

func (b *BLangOnFailClause) GetVariableDefinitionNode() model.VariableDefinitionNode {
	// migrated from BLangOnFailClause.java:63:5
	return b.VariableDefinitionNode
}

func (b *BLangOnFailClause) SetVariableDefinitionNode(variableDefinitionNode model.VariableDefinitionNode) {
	// migrated from BLangOnFailClause.java:68:5
	b.VariableDefinitionNode = variableDefinitionNode
}

func (b *BLangOnFailClause) GetBody() model.BlockStatementNode {
	// migrated from BLangOnFailClause.java:73:5
	return b.Body
}

func (b *BLangOnFailClause) SetBody(body model.BlockStatementNode) {
	// migrated from BLangOnFailClause.java:98:5
	if body, ok := body.(*BLangBlockStmt); ok {
		b.Body = body
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangOnFailClause) GetKind() model.NodeKind {
	// migrated from BLangOnFailClause.java:93:5
	return model.NodeKind_ON_FAIL
}
