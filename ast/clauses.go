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
	BLangLimitClause struct {
		bLangNodeBase
		Expression BLangExpression
	}
	BLangOrderByClause struct {
		bLangNodeBase
		OrderByKeyList []BLangOrderKey
	}
	BLangOrderKey struct {
		bLangNodeBase
		Expression  BLangExpression
		IsAscending bool
	}
	BLangSelectClause struct {
		bLangNodeBase
		Expression BLangExpression
	}
	BLangOnConflictClause struct {
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
	_ model.Node              = &BLangLimitClause{}
	_ model.Node              = &BLangOrderByClause{}
	_ model.Node              = &BLangOrderKey{}
	_ model.SelectClauseNode  = &BLangSelectClause{}
	_ model.Node              = &BLangOnConflictClause{}
	_ model.CollectClauseNode = &BLangCollectClause{}
	_ model.DoClauseNode      = &BLangDoClause{}
	_ model.OnFailClauseNode  = &BLangOnFailClause{}
)

var (
	_ BLangNode = &BLangFromClause{}
	_ BLangNode = &BLangLetClause{}
	_ BLangNode = &BLangWhereClause{}
	_ BLangNode = &BLangLimitClause{}
	_ BLangNode = &BLangOrderByClause{}
	_ BLangNode = &BLangOrderKey{}
	_ BLangNode = &BLangSelectClause{}
	_ BLangNode = &BLangOnConflictClause{}
	_ BLangNode = &BLangCollectClause{}
	_ BLangNode = &BLangDoClause{}
	_ BLangNode = &BLangOnFailClause{}
)

func (this *BLangFromClause) GetKind() model.NodeKind {
	return model.NodeKind_FROM
}

func (this *BLangFromClause) GetCollection() model.ExpressionNode {
	return this.Collection
}

func (this *BLangFromClause) SetCollection(collection model.ExpressionNode) {
	if exp, ok := collection.(BLangExpression); ok {
		this.Collection = exp
		return
	}
	panic("collection is not a BLangExpression")
}

func (this *BLangFromClause) GetVariableDefinitionNode() model.VariableDefinitionNode {
	return this.VariableDefinitionNode
}

func (this *BLangFromClause) SetVariableDefinitionNode(variableDefinitionNode model.VariableDefinitionNode) {
	this.VariableDefinitionNode = variableDefinitionNode
}

func (this *BLangFromClause) IsDeclaredWithVar() bool {
	return this.IsDeclaredWithVarFlag
}

func (this *BLangLetClause) GetKind() model.NodeKind {
	return model.NodeKind_LET_CLAUSE
}

func (this *BLangWhereClause) GetKind() model.NodeKind {
	return model.NodeKind_WHERE
}

func (this *BLangLimitClause) GetKind() model.NodeKind {
	return model.NodeKind_LIMIT
}

func (this *BLangLimitClause) GetExpression() model.ExpressionNode {
	return this.Expression
}

func (this *BLangLimitClause) SetExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		this.Expression = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (this *BLangSelectClause) GetKind() model.NodeKind {
	return model.NodeKind_SELECT
}

func (this *BLangOrderByClause) GetKind() model.NodeKind {
	return model.NodeKind_ORDER_BY
}

func (this *BLangOrderKey) GetKind() model.NodeKind {
	return model.NodeKind_ORDER_KEY
}

func (this *BLangSelectClause) GetExpression() model.ExpressionNode {
	return this.Expression
}

func (this *BLangSelectClause) SetExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		this.Expression = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (this *BLangOnConflictClause) GetKind() model.NodeKind {
	return model.NodeKind_ON_CONFLICT
}

func (this *BLangOnConflictClause) GetExpression() model.ExpressionNode {
	return this.Expression
}

func (this *BLangOnConflictClause) SetExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		this.Expression = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (this *BLangCollectClause) GetKind() model.NodeKind {
	// migrated from BLangCollectClause.java:48:5
	return model.NodeKind_COLLECT
}

func (this *BLangCollectClause) GetExpression() model.ExpressionNode {
	// migrated from BLangCollectClause.java:68:5
	return this.Expression
}

func (this *BLangCollectClause) SetExpression(expression model.ExpressionNode) {
	// migrated from BLangCollectClause.java:73:5
	if exp, ok := expression.(BLangExpression); ok {
		this.Expression = exp
	} else {
		panic("Expected BLangExpression")
	}
}

func (this *BLangDoClause) GetBody() model.BlockStatementNode {
	// migrated from BLangDoClause.java:46:5
	return this.Body
}

func (this *BLangDoClause) SetBody(body model.BlockStatementNode) {
	// migrated from BLangDoClause.java:51:5
	if body, ok := body.(*BLangBlockStmt); ok {
		this.Body = body
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (this *BLangDoClause) GetKind() model.NodeKind {
	// migrated from BLangDoClause.java:66:5
	return model.NodeKind_DO
}

func (this *BLangOnFailClause) SetDeclaredWithVar() {
	// migrated from BLangOnFailClause.java:53:5
	this.isDeclaredWithVar = true
}

func (this *BLangOnFailClause) IsDeclaredWithVar() bool {
	// migrated from BLangOnFailClause.java:58:5
	return this.isDeclaredWithVar
}

func (this *BLangOnFailClause) GetVariableDefinitionNode() model.VariableDefinitionNode {
	// migrated from BLangOnFailClause.java:63:5
	return this.VariableDefinitionNode
}

func (this *BLangOnFailClause) SetVariableDefinitionNode(variableDefinitionNode model.VariableDefinitionNode) {
	// migrated from BLangOnFailClause.java:68:5
	this.VariableDefinitionNode = variableDefinitionNode
}

func (this *BLangOnFailClause) GetBody() model.BlockStatementNode {
	// migrated from BLangOnFailClause.java:73:5
	return this.Body
}

func (this *BLangOnFailClause) SetBody(body model.BlockStatementNode) {
	// migrated from BLangOnFailClause.java:98:5
	if body, ok := body.(*BLangBlockStmt); ok {
		this.Body = body
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (this *BLangOnFailClause) GetKind() model.NodeKind {
	// migrated from BLangOnFailClause.java:93:5
	return model.NodeKind_ON_FAIL
}
