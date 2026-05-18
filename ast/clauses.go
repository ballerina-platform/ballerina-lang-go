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
	BLangJoinClause struct {
		BLangInputClause
		OnClause        BLangOnClause
		IsOuterJoinFlag bool
	}
	BLangLetClause struct {
		bLangNodeBase
		LetVarDeclarations []model.VariableDefinitionNode
	}
	BLangOnClause struct {
		bLangNodeBase
		OnExpr     BLangExpression
		EqualsExpr BLangExpression
	}
	BLangWhereClause struct {
		bLangNodeBase
		Expression BLangExpression
	}
	BLangGroupByClause struct {
		bLangNodeBase
		GroupingKeyList []BLangGroupingKey
		NonGroupingKeys common.Set[string]
	}
	BLangGroupingKey struct {
		bLangNodeBase
		VariableDef *BLangSimpleVariableDef
		VariableRef *BLangSimpleVarRef
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
		Expression   BLangExpression
		IsDescending bool
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
	_ model.JoinClauseNode    = &BLangJoinClause{}
	_ model.Node              = &BLangLetClause{}
	_ model.OnClauseNode      = &BLangOnClause{}
	_ model.Node              = &BLangWhereClause{}
	_ model.GroupByClauseNode = &BLangGroupByClause{}
	_ model.GroupingKeyNode   = &BLangGroupingKey{}
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
	_ BLangNode = &BLangJoinClause{}
	_ BLangNode = &BLangLetClause{}
	_ BLangNode = &BLangOnClause{}
	_ BLangNode = &BLangWhereClause{}
	_ BLangNode = &BLangGroupByClause{}
	_ BLangNode = &BLangGroupingKey{}
	_ BLangNode = &BLangLimitClause{}
	_ BLangNode = &BLangOrderByClause{}
	_ BLangNode = &BLangOrderKey{}
	_ BLangNode = &BLangSelectClause{}
	_ BLangNode = &BLangOnConflictClause{}
	_ BLangNode = &BLangCollectClause{}
	_ BLangNode = &BLangDoClause{}
	_ BLangNode = &BLangOnFailClause{}
)

func (b *BLangFromClause) GetKind() model.NodeKind {
	return model.NodeKind_FROM
}

func (b *BLangJoinClause) GetKind() model.NodeKind {
	return model.NodeKind_JOIN
}

func (b *BLangJoinClause) GetCollection() model.ExpressionNode {
	return b.Collection
}

func (b *BLangJoinClause) SetCollection(collection model.ExpressionNode) {
	if exp, ok := collection.(BLangExpression); ok {
		b.Collection = exp
		return
	}
	panic("collection is not a BLangExpression")
}

func (b *BLangJoinClause) GetVariableDefinitionNode() model.VariableDefinitionNode {
	return b.VariableDefinitionNode
}

func (b *BLangJoinClause) SetVariableDefinitionNode(variableDefinitionNode model.VariableDefinitionNode) {
	b.VariableDefinitionNode = variableDefinitionNode
}

func (b *BLangJoinClause) IsDeclaredWithVar() bool {
	return b.IsDeclaredWithVarFlag
}

func (b *BLangJoinClause) GetOnClause() model.OnClauseNode {
	if b.OnClause.OnExpr == nil && b.OnClause.EqualsExpr == nil {
		return nil
	}
	return &b.OnClause
}

func (b *BLangJoinClause) IsOuterJoin() bool {
	return b.IsOuterJoinFlag
}

func (b *BLangOnClause) GetKind() model.NodeKind {
	return model.NodeKind_ON
}

func (b *BLangOnClause) GetOnExpression() model.ExpressionNode {
	return b.OnExpr
}

func (b *BLangOnClause) SetOnExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		b.OnExpr = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (b *BLangOnClause) GetEqualsExpression() model.ExpressionNode {
	return b.EqualsExpr
}

func (b *BLangOnClause) SetEqualsExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		b.EqualsExpr = exp
		return
	}
	panic("expression is not a BLangExpression")
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

func (b *BLangGroupByClause) GetKind() model.NodeKind {
	return model.NodeKind_GROUP_BY
}

func (b *BLangGroupByClause) AddGroupingKey(groupingKey model.GroupingKeyNode) {
	if key, ok := groupingKey.(*BLangGroupingKey); ok {
		b.GroupingKeyList = append(b.GroupingKeyList, *key)
		return
	}
	panic("groupingKey is not a BLangGroupingKey")
}

func (b *BLangGroupByClause) GetGroupingKeyList() []model.GroupingKeyNode {
	result := make([]model.GroupingKeyNode, len(b.GroupingKeyList))
	for i := range b.GroupingKeyList {
		result[i] = &b.GroupingKeyList[i]
	}
	return result
}

func (b *BLangGroupingKey) GetKind() model.NodeKind {
	return model.NodeKind_GROUPING_KEY
}

func (b *BLangGroupingKey) SetGroupingKey(groupingKey model.Node) {
	switch key := groupingKey.(type) {
	case *BLangSimpleVariableDef:
		b.VariableDef = key
		b.VariableRef = nil
	case *BLangSimpleVarRef:
		b.VariableRef = key
		b.VariableDef = nil
	default:
		panic("groupingKey is neither a BLangSimpleVariableDef nor a BLangSimpleVarRef")
	}
}

func (b *BLangGroupingKey) GetGroupingKey() model.Node {
	if b.VariableRef != nil {
		return b.VariableRef
	}
	if b.VariableDef != nil {
		return b.VariableDef
	}
	return nil
}

func (b *BLangLimitClause) GetKind() model.NodeKind {
	return model.NodeKind_LIMIT
}

func (b *BLangLimitClause) GetExpression() model.ExpressionNode {
	return b.Expression
}

func (b *BLangLimitClause) SetExpression(expression model.ExpressionNode) {
	if exp, ok := expression.(BLangExpression); ok {
		b.Expression = exp
		return
	}
	panic("expression is not a BLangExpression")
}

func (b *BLangSelectClause) GetKind() model.NodeKind {
	return model.NodeKind_SELECT
}

func (b *BLangOrderByClause) GetKind() model.NodeKind {
	return model.NodeKind_ORDER_BY
}

func (b *BLangOrderKey) GetKind() model.NodeKind {
	return model.NodeKind_ORDER_KEY
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

func (b *BLangOnConflictClause) GetKind() model.NodeKind {
	return model.NodeKind_ON_CONFLICT
}

func (b *BLangOnConflictClause) GetExpression() model.ExpressionNode {
	return b.Expression
}

func (b *BLangOnConflictClause) SetExpression(expression model.ExpressionNode) {
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
