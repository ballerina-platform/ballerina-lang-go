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

// Type aliases for model interfaces
type (
	AssignmentNode          = model.AssignmentNode
	CompoundAssignmentNode  = model.CompoundAssignmentNode
	DoNode                  = model.DoNode
	BlockStatementNode      = model.BlockStatementNode
	ExpressionStatementNode = model.ExpressionStatementNode
	IfNode                  = model.IfNode
	BlockNode               = model.BlockNode
	VariableDefinitionNode  = model.VariableDefinitionNode
	ReturnNode              = model.ReturnNode
	WhileNode               = model.WhileNode
	BLangStatement          = model.StatementNode
	StatementNode           = model.StatementNode
	ContinueNode            = model.ContinueNode
)
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

	BLangExpressionStmt struct {
		BLangStatementBase
		Expr BLangExpression
	}

	BLangIf struct {
		BLangStatementBase
		Expr     BLangExpression
		Body     BLangBlockStmt
		ElseStmt BLangStatement
	}

	BLangWhile struct {
		BLangStatementBase
		Expr         BLangExpression
		Body         BLangBlockStmt
		OnFailClause BLangOnFailClause
	}

	BLangSimpleVariableDef struct {
		BLangStatementBase
		Var      BLangSimpleVariable
		IsInFork bool
		IsWorker bool
	}

	BLangReturn struct {
		BLangStatementBase
		Expr BLangExpression
	}
)

var _ AssignmentNode = &BLangAssignment{}
var _ CompoundAssignmentNode = &BLangCompoundAssignment{}
var _ ContinueNode = &BLangContinue{}
var _ DoNode = &BLangDo{}
var _ BlockStatementNode = &BLangBlockStmt{}
var _ ExpressionStatementNode = &BLangExpressionStmt{}
var _ IfNode = &BLangIf{}
var _ WhileNode = &BLangWhile{}
var _ VariableDefinitionNode = &BLangSimpleVariableDef{}
var _ ReturnNode = &BLangReturn{}

var _ BLangNode = &BLangAssignment{}
var _ BLangNode = &BLangBlockStmt{}
var _ BLangNode = &BLangBreak{}
var _ BLangNode = &BLangCompoundAssignment{}
var _ BLangNode = &BLangContinue{}
var _ BLangNode = &BLangDo{}
var _ BLangNode = &BLangExpressionStmt{}
var _ BLangNode = &BLangIf{}
var _ BLangNode = &BLangWhile{}
var _ BLangNode = &BLangSimpleVariableDef{}

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

func (this *BLangAssignment) GetKind() NodeKind {
	return NodeKind_ASSIGNMENT
}

func (this *BLangAssignment) SetExpression(expression ExpressionNode) {
	// migrated from BLangAssignment.java:64:5
	if expr, ok := expression.(BLangExpression); ok {
		this.Expr = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (this *BLangAssignment) SetDeclaredWithVar(isDeclaredWithVar bool) {
	// migrated from BLangAssignment.java:69:5
}

func (this *BLangAssignment) SetVariable(variableReferenceNode VariableReferenceNode) {
	// migrated from BLangAssignment.java:74:5
	if varRef, ok := variableReferenceNode.(BLangExpression); ok {
		this.VarRef = varRef
	} else {
		panic("variableReferenceNode is not a BLangExpression")
	}
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

func (this *BLangExpressionStmt) GetExpression() ExpressionNode {
	// migrated from BLangExpressionStmt.java:46:5
	return this.Expr
}

func (this *BLangExpressionStmt) GetKind() NodeKind {
	return NodeKind_EXPRESSION_STATEMENT
}

func (this *BLangIf) GetCondition() ExpressionNode {
	// migrated from BLangIf.java:47:5
	return this.Expr
}

func (this *BLangIf) GetBody() BlockStatementNode {
	// migrated from BLangIf.java:52:5
	return &this.Body
}

func (this *BLangIf) GetElseStatement() StatementNode {
	// migrated from BLangIf.java:57:5
	return this.ElseStmt
}

func (this *BLangIf) SetCondition(condition ExpressionNode) {
	// migrated from BLangIf.java:62:5
	if expr, ok := condition.(BLangExpression); ok {
		this.Expr = expr
	} else {
		panic("condition is not a BLangExpression")
	}
}

func (this *BLangIf) SetBody(body BlockStatementNode) {
	// migrated from BLangIf.java:67:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		this.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (this *BLangIf) SetElseStatement(elseStatement StatementNode) {
	// migrated from BLangIf.java:72:5
	if elseStmt, ok := elseStatement.(BLangStatement); ok {
		this.ElseStmt = elseStmt
		return
	}
	panic("elseStatement is not a BLangStatement")
}

func (this *BLangIf) GetKind() NodeKind {
	// migrated from BLangIf.java:77:5
	return NodeKind_IF
}

func (this *BLangWhile) GetCondition() ExpressionNode {
	// migrated from BLangWhile.java:50:5
	return this.Expr
}

func (this *BLangWhile) SetCondition(condition ExpressionNode) {
	// migrated from BLangWhile.java:60:5
	if expr, ok := condition.(BLangExpression); ok {
		this.Expr = expr
	} else {
		panic("condition is not a BLangExpression")
	}
}

func (this *BLangWhile) GetBody() BlockStatementNode {
	// migrated from BLangWhile.java:55:5
	return &this.Body
}

func (this *BLangWhile) SetBody(body BlockStatementNode) {
	// migrated from BLangWhile.java:65:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		this.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (this *BLangWhile) GetOnFailClause() OnFailClauseNode {
	// migrated from BLangWhile.java:70:5
	return &this.OnFailClause
}

func (this *BLangWhile) SetOnFailClause(onFailClause OnFailClauseNode) {
	// migrated from BLangWhile.java:75:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		this.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (this *BLangWhile) GetKind() NodeKind {
	// migrated from BLangWhile.java:95:5
	return NodeKind_WHILE
}

func (this *BLangSimpleVariableDef) GetIsInFork() bool {
	return this.IsInFork
}

func (this *BLangSimpleVariableDef) GetIsWorker() bool {
	return this.IsWorker
}

func (this *BLangSimpleVariableDef) GetKind() NodeKind {
	return NodeKind_VARIABLE_DEF
}

func (this *BLangSimpleVariableDef) GetVariable() VariableNode {
	return &this.Var
}

func (this *BLangSimpleVariableDef) SetVariable(variable VariableNode) {
	if v, ok := variable.(*BLangSimpleVariable); ok {
		this.Var = *v
	} else {
		panic("variable is not a BLangSimpleVariable")
	}
}

func (this *BLangReturn) GetExpression() ExpressionNode {
	return this.Expr
}

func (this *BLangReturn) SetExpression(expression ExpressionNode) {
	if expr, ok := expression.(BLangExpression); ok {
		this.Expr = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (this *BLangReturn) GetKind() NodeKind {
	return NodeKind_RETURN
}
