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
	_ model.CollectClauseNode = &BLangCollectClause{}
	_ model.DoClauseNode      = &BLangDoClause{}
	_ model.OnFailClauseNode  = &BLangOnFailClause{}
)

var (
	_ BLangNode = &BLangCollectClause{}
	_ BLangNode = &BLangDoClause{}
	_ BLangNode = &BLangOnFailClause{}
)

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
