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

type AssignmentNode interface {
	GetVariable() ExpressionNode
	GetExpression() ExpressionNode
	IsDeclaredWithVar() bool
	SetExpression(expression Node)
	SetDeclaredWithVar(IsDeclaredWithVar bool)
	SetVariable(variableReferenceNode VariableReferenceNode)
}

type CompoundAssignmentNode interface {
	StatementNode
	GetVariable() ExpressionNode
	GetExpression() ExpressionNode
	SetExpression(expression ExpressionNode)
	SetVariable(variableReferenceNode VariableReferenceNode)
	GetOperatorKind() OperatorKind
}

type DoNode interface {
	StatementNode
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
	GetOnFailClause() OnFailClauseNode
	SetOnFailClause(onFailClause OnFailClauseNode)
}

type BlockStatementNode interface {
	BlockNode
	StatementNode
}

type BlockNode interface {
	Node
	GetStatements() []StatementNode
	AddStatement(statement StatementNode)
}

type VariableDefinitionNode interface {
	StatementNode
	GetVariable() VariableNode
	SetVariable(variable VariableNode)
	GetIsInFork() bool
	GetIsWorker() bool
}

type OnFailClauseNode interface {
	Node
	SetDeclaredWithVar()
	IsDeclaredWithVar() bool
	GetVariableDefinitionNode() VariableDefinitionNode
	SetVariableDefinitionNode(variableDefinitionNode VariableDefinitionNode)
	GetBody() BlockStatementNode
	SetBody(body BlockStatementNode)
}

type BLangStatement = Node
type StatementNode = Node
type ContinueNode = StatementNode
type FailureBreakMode uint

const (
	FailureBreakMode_NOT_BREAKABLE FailureBreakMode = iota
	FailureBreakMode_BREAK_WITHIN_BLOCK
	FailureBreakMode_BREAK_TO_OUTER_BLOCK
)

type (
	BLangStatementBase struct {
		BLangNodeBase
	}
	BLangAssignment struct {
		BLangStatementBase
		VarRef BLangExpression
		Expr   BLangExpression
	}
	BLangBlockStmt struct {
		BLangStatementBase
		Stmts            []BLangStatement
		MapSymbol        BVarSymbol
		FailureBreakMode FailureBreakMode
		IsLetExpr        bool
		Scope            Scope
	}
	BLangBreak struct {
		BLangStatementBase
	}

	BLangCompoundAssignment struct {
		BLangStatementBase
		VarRef       ExpressionNode
		Expr         BLangExpression
		OpKind       OperatorKind
		ModifiedExpr BLangExpression
	}
	BLangContinue struct {
		BLangStatementBase
	}
	BLangDo struct {
		BLangStatementBase
		Body         BLangBlockStmt
		OnFailClause BLangOnFailClause
	}
)

var _ AssignmentNode = &BLangAssignment{}
var _ CompoundAssignmentNode = &BLangCompoundAssignment{}
var _ ContinueNode = &BLangContinue{}
var _ DoNode = &BLangDo{}
var _ BlockStatementNode = &BLangBlockStmt{}

var _ BLangNode = &BLangStatementBase{}
var _ BLangNode = &BLangAssignment{}
var _ BLangNode = &BLangBlockStmt{}
var _ BLangNode = &BLangBreak{}
var _ BLangNode = &BLangCompoundAssignment{}
var _ BLangNode = &BLangContinue{}
var _ BLangNode = &BLangDo{}

func (this *BLangAssignment) GetVariable() ExpressionNode {
	// migrated from BLangAssignment.java:48:5
	return this.VarRef
}

func (this *BLangAssignment) GetExpression() ExpressionNode {
	// migrated from BLangAssignment.java:53:5
	return this.Expr
}

func (this *BLangAssignment) IsDeclaredWithVar() bool {
	// migrated from BLangAssignment.java:58:5
	return false
}

func (this *BLangAssignment) SetExpression(expression ExpressionNode) {
	// migrated from BLangAssignment.java:64:5
	this.Expr = expression
}

func (this *BLangAssignment) SetDeclaredWithVar(isDeclaredWithVar bool) {
	// migrated from BLangAssignment.java:69:5
}

func (this *BLangAssignment) SetVariable(variableReferenceNode VariableReferenceNode) {
	// migrated from BLangAssignment.java:74:5
	this.VarRef = variableReferenceNode
}

func (this *BLangBlockStmt) GetKind() NodeKind {
	// migrated from BLangBlockStmt.java:83:5
	return NodeKind_BLOCK
}

func (this *BLangBlockStmt) GetStatements() []StatementNode {
	// migrated from BLangBlockStmt.java:88:5
	return this.Stmts
}

func (this *BLangBlockStmt) AddStatement(statement StatementNode) {
	// migrated from BLangBlockStmt.java:93:5
	this.Stmts = append(this.Stmts, statement)
}
func (this *BLangBreak) GetKind() NodeKind {
	// migrated from BLangBreak.java:45:5
	return NodeKind_BREAK
}

func (this *BLangCompoundAssignment) GetOperatorKind() OperatorKind {
	// migrated from BLangCompoundAssignment.java:59:5
	return this.OpKind
}

func (this *BLangCompoundAssignment) GetVariable() ExpressionNode {
	// migrated from BLangCompoundAssignment.java:64:5
	return this.VarRef
}

func (this *BLangCompoundAssignment) GetExpression() ExpressionNode {
	// migrated from BLangCompoundAssignment.java:69:5
	return this.Expr
}

func (this *BLangCompoundAssignment) SetExpression(expression ExpressionNode) {
	// migrated from BLangCompoundAssignment.java:74:5
	if exp, ok := expression.(BLangExpression); ok {
		this.Expr = exp
	} else {
		panic("Expected BLangExpression")
	}
}

func (this *BLangCompoundAssignment) SetVariable(variableReferenceNode VariableReferenceNode) {
	// migrated from BLangCompoundAssignment.java:79:5
	this.VarRef = variableReferenceNode
}

func (this *BLangCompoundAssignment) GetKind() NodeKind {
	// migrated from BLangCompoundAssignment.java:99:5
	return NodeKind_COMPOUND_ASSIGNMENT
}

func (this *BLangContinue) GetKind() NodeKind {
	// migrated from BLangContinue.java:46:5
	return NodeKind_NEXT
}

func (this *BLangDo) GetBody() BlockStatementNode {
	// migrated from BLangDo.java:47:5
	return &this.Body
}

func (this *BLangDo) SetBody(body BlockStatementNode) {
	// migrated from BLangDo.java:52:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		this.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (this *BLangDo) GetOnFailClause() OnFailClauseNode {
	// migrated from BLangDo.java:57:5
	return &this.OnFailClause
}

func (this *BLangDo) SetOnFailClause(onFailClause OnFailClauseNode) {
	// migrated from BLangDo.java:62:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		this.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (this *BLangDo) GetKind() NodeKind {
	// migrated from BLangDo.java:82:5
	return NodeKind_DO_STMT
}
