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

type FailureBreakMode uint

const (
	FailureBreakMode_NOT_BREAKABLE FailureBreakMode = iota
	FailureBreakMode_BREAK_WITHIN_BLOCK
	FailureBreakMode_BREAK_TO_OUTER_BLOCK
)

type (
	bLangStatementBase struct {
		bLangNodeBase
	}
	BLangAssignment struct {
		bLangStatementBase
		VarRef BLangExpression
		Expr   BLangActionOrExpression
	}
	BLangBlockStmt struct {
		bLangStatementBase
		Stmts            []BLangStatement
		FailureBreakMode FailureBreakMode
		IsLetExpr        bool
	}
	BLangBreak struct {
		bLangStatementBase
	}

	BLangCompoundAssignment struct {
		bLangStatementBase
		VarRef       BLangExpression
		Expr         BLangActionOrExpression
		OpKind       model.OperatorKind
		ModifiedExpr BLangExpression
	}
	BLangContinue struct {
		bLangStatementBase
	}
	BLangDo struct {
		bLangStatementBase
		Body         BLangBlockStmt
		OnFailClause BLangOnFailClause
	}

	BLangExpressionStmt struct {
		bLangStatementBase
		Expr BLangActionOrExpression
	}

	BLangIf struct {
		bLangStatementBase
		scope    model.Scope
		Expr     BLangExpression
		Body     BLangBlockStmt
		ElseStmt BLangStatement
	}

	BLangWhile struct {
		bLangStatementBase
		scope        model.Scope
		Expr         BLangExpression
		Body         BLangBlockStmt
		OnFailClause BLangOnFailClause
	}

	BLangForeach struct {
		bLangStatementBase
		scope             model.Scope
		VariableDef       *BLangSimpleVariableDef
		Collection        BLangActionOrExpression
		Body              BLangBlockStmt
		OnFailClause      *BLangOnFailClause
		IsDeclaredWithVar bool
	}

	BLangSimpleVariableDef struct {
		bLangStatementBase
		Var      *BLangSimpleVariable
		IsInFork bool
		IsWorker bool
	}

	BLangReturn struct {
		bLangStatementBase
		Expr BLangActionOrExpression
	}

	BLangPanic struct {
		bLangStatementBase
		Expr BLangExpression
	}

	BLangMatchStatement struct {
		bLangStatementBase
		Expr         BLangActionOrExpression
		MatchClauses []BLangMatchClause
		IsExhaustive bool
	}
)

var (
	_ AssignmentNode          = &BLangAssignment{}
	_ CompoundAssignmentNode  = &BLangCompoundAssignment{}
	_ ContinueNode            = &BLangContinue{}
	_ DoNode                  = &BLangDo{}
	_ BlockStatementNode      = &BLangBlockStmt{}
	_ ExpressionStatementNode = &BLangExpressionStmt{}
	_ IfNode                  = &BLangIf{}
	_ WhileNode               = &BLangWhile{}
	_ ForeachNode             = &BLangForeach{}
	_ VariableDefinitionNode  = &BLangSimpleVariableDef{}
	_ ReturnNode              = &BLangReturn{}
	_ PanicNode               = &BLangPanic{}
)

var (
	_ NodeWithScope = &BLangIf{}
	_ NodeWithScope = &BLangWhile{}
	_ NodeWithScope = &BLangForeach{}
)

var (
	_ BLangNode = &BLangAssignment{}
	_ BLangNode = &BLangBlockStmt{}
	_ BLangNode = &BLangBreak{}
	_ BLangNode = &BLangCompoundAssignment{}
	_ BLangNode = &BLangContinue{}
	_ BLangNode = &BLangDo{}
	_ BLangNode = &BLangExpressionStmt{}
	_ BLangNode = &BLangIf{}
	_ BLangNode = &BLangWhile{}
	_ BLangNode = &BLangForeach{}
	_ BLangNode = &BLangSimpleVariableDef{}
	_ BLangNode = &BLangReturn{}
	_ BLangNode = &BLangPanic{}
	_ BLangNode = &BLangMatchStatement{}
)

func (b *BLangAssignment) GetVariable() BLangExpression {
	// migrated from BLangAssignment.java:48:5
	return b.VarRef
}

func (b *BLangAssignment) GetExpression() BLangActionOrExpression {
	// migrated from BLangAssignment.java:53:5
	return b.Expr
}

func (b *BLangAssignment) IsDeclaredWithVar() bool {
	// migrated from BLangAssignment.java:58:5
	return false
}

func (b *BLangAssignment) GetKind() NodeKind {
	return NodeKind_ASSIGNMENT
}

func (b *BLangAssignment) SetActionOrExpression(actionOrExpression BLangActionOrExpression) {
	b.Expr = actionOrExpression
}

func (b *BLangAssignment) SetDeclaredWithVar(isDeclaredWithVar bool) {
	// migrated from BLangAssignment.java:69:5
}

func (b *BLangAssignment) SetVariable(variableReferenceNode VariableReferenceNode) {
	// migrated from BLangAssignment.java:74:5
	b.VarRef = variableReferenceNode
}

func (b *BLangBlockStmt) GetKind() NodeKind {
	// migrated from BLangBlockStmt.java:83:5
	return NodeKind_BLOCK
}

func (b *BLangBlockStmt) GetStatements() []StatementNode {
	// migrated from BLangBlockStmt.java:88:5
	return b.Stmts
}

func (b *BLangBlockStmt) AddStatement(statement StatementNode) {
	// migrated from BLangBlockStmt.java:93:5
	b.Stmts = append(b.Stmts, statement)
}

func (b *BLangBreak) GetKind() NodeKind {
	// migrated from BLangBreak.java:45:5
	return NodeKind_BREAK
}

func (b *BLangCompoundAssignment) IsDeclaredWithVar() bool {
	return false
}

func (b *BLangCompoundAssignment) SetDeclaredWithVar(_ bool) {
	panic("compound assignemnt can't be declared with var")
}

func (b *BLangCompoundAssignment) GetOperatorKind() model.OperatorKind {
	// migrated from BLangCompoundAssignment.java:59:5
	return b.OpKind
}

func (b *BLangCompoundAssignment) GetVariable() BLangExpression {
	// migrated from BLangCompoundAssignment.java:64:5
	return b.VarRef
}

func (b *BLangCompoundAssignment) GetExpression() BLangActionOrExpression {
	// migrated from BLangCompoundAssignment.java:69:5
	return b.Expr
}

func (b *BLangCompoundAssignment) SetActionOrExpression(actionOrExpression BLangActionOrExpression) {
	b.Expr = actionOrExpression
}

func (b *BLangCompoundAssignment) SetVariable(variableReferenceNode VariableReferenceNode) {
	b.VarRef = variableReferenceNode
}

func (b *BLangCompoundAssignment) GetKind() NodeKind {
	// migrated from BLangCompoundAssignment.java:99:5
	return NodeKind_COMPOUND_ASSIGNMENT
}

func (b *BLangContinue) GetKind() NodeKind {
	// migrated from BLangContinue.java:46:5
	return NodeKind_NEXT
}

func (b *BLangDo) GetBody() BlockStatementNode {
	// migrated from BLangDo.java:47:5
	return &b.Body
}

func (b *BLangDo) SetBody(body BlockStatementNode) {
	// migrated from BLangDo.java:52:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangDo) GetOnFailClause() OnFailClauseNode {
	// migrated from BLangDo.java:57:5
	return &b.OnFailClause
}

func (b *BLangDo) SetOnFailClause(onFailClause OnFailClauseNode) {
	// migrated from BLangDo.java:62:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		b.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (b *BLangDo) GetKind() NodeKind {
	// migrated from BLangDo.java:82:5
	return NodeKind_DO_STMT
}

func (b *BLangExpressionStmt) GetExpression() BLangActionOrExpression {
	// migrated from BLangExpressionStmt.java:46:5
	return b.Expr
}

func (b *BLangExpressionStmt) GetKind() NodeKind {
	return NodeKind_EXPRESSION_STATEMENT
}

func (b *BLangIf) Scope() model.Scope {
	return b.scope
}

func (b *BLangIf) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangIf) GetCondition() BLangExpression {
	// migrated from BLangIf.java:47:5
	return b.Expr
}

func (b *BLangIf) GetBody() BlockStatementNode {
	// migrated from BLangIf.java:52:5
	return &b.Body
}

func (b *BLangIf) GetElseStatement() StatementNode {
	// migrated from BLangIf.java:57:5
	return b.ElseStmt
}

func (b *BLangIf) SetCondition(condition BLangExpression) {
	// migrated from BLangIf.java:62:5
	b.Expr = condition
}

func (b *BLangIf) SetBody(body BlockStatementNode) {
	// migrated from BLangIf.java:67:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangIf) SetElseStatement(elseStatement StatementNode) {
	// migrated from BLangIf.java:72:5
	b.ElseStmt = elseStatement
}

func (b *BLangIf) GetKind() NodeKind {
	// migrated from BLangIf.java:77:5
	return NodeKind_IF
}

func (b *BLangWhile) Scope() model.Scope {
	return b.scope
}

func (b *BLangWhile) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangWhile) GetCondition() BLangExpression {
	// migrated from BLangWhile.java:50:5
	return b.Expr
}

func (b *BLangWhile) SetCondition(condition BLangExpression) {
	// migrated from BLangWhile.java:60:5
	b.Expr = condition
}

func (b *BLangWhile) GetBody() BlockStatementNode {
	// migrated from BLangWhile.java:55:5
	return &b.Body
}

func (b *BLangWhile) SetBody(body BlockStatementNode) {
	// migrated from BLangWhile.java:65:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangWhile) GetOnFailClause() OnFailClauseNode {
	// migrated from BLangWhile.java:70:5
	return &b.OnFailClause
}

func (b *BLangWhile) SetOnFailClause(onFailClause OnFailClauseNode) {
	// migrated from BLangWhile.java:75:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		b.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (b *BLangWhile) GetKind() NodeKind {
	// migrated from BLangWhile.java:95:5
	return NodeKind_WHILE
}

func (b *BLangForeach) Scope() model.Scope {
	return b.scope
}

func (b *BLangForeach) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangForeach) GetKind() NodeKind {
	return NodeKind_FOREACH
}

func (b *BLangForeach) GetVariableDefinitionNode() VariableDefinitionNode {
	return b.VariableDef
}

func (b *BLangForeach) SetVariableDefinitionNode(node VariableDefinitionNode) {
	if varDef, ok := node.(*BLangSimpleVariableDef); ok {
		b.VariableDef = varDef
		return
	}
	panic("node is not a *BLangSimpleVariableDef")
}

func (b *BLangForeach) GetCollection() BLangActionOrExpression {
	return b.Collection
}

func (b *BLangForeach) SetCollection(collection BLangActionOrExpression) {
	b.Collection = collection
}

func (b *BLangForeach) GetBody() BlockStatementNode {
	return &b.Body
}

func (b *BLangForeach) SetBody(body BlockStatementNode) {
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangForeach) GetIsDeclaredWithVar() bool {
	return b.IsDeclaredWithVar
}

func (b *BLangForeach) GetOnFailClause() OnFailClauseNode {
	if b.OnFailClause == nil {
		return nil
	}
	return b.OnFailClause
}

func (b *BLangForeach) SetOnFailClause(onFailClause OnFailClauseNode) {
	if clause, ok := onFailClause.(*BLangOnFailClause); ok {
		b.OnFailClause = clause
		return
	}
	panic("onFailClause is not a *BLangOnFailClause")
}

func (b *BLangSimpleVariableDef) GetIsInFork() bool {
	return b.IsInFork
}

func (b *BLangSimpleVariableDef) GetIsWorker() bool {
	return b.IsWorker
}

func (b *BLangSimpleVariableDef) GetKind() NodeKind {
	return NodeKind_VARIABLE_DEF
}

func (b *BLangSimpleVariableDef) GetVariable() VariableNode {
	return b.Var
}

func (b *BLangSimpleVariableDef) SetVariable(variable VariableNode) {
	if v, ok := variable.(*BLangSimpleVariable); ok {
		b.Var = v
	} else {
		panic("variable is not a BLangSimpleVariable")
	}
}

func (b *BLangReturn) GetExpression() BLangActionOrExpression {
	return b.Expr
}

func (b *BLangReturn) SetActionOrExpression(actionOrExpression BLangActionOrExpression) {
	b.Expr = actionOrExpression
}

func (b *BLangReturn) GetKind() NodeKind {
	return NodeKind_RETURN
}

func (b *BLangPanic) GetExpression() BLangExpression {
	return b.Expr
}

func (b *BLangPanic) GetKind() NodeKind {
	return NodeKind_PANIC
}

func (b *BLangMatchStatement) GetKind() NodeKind {
	return NodeKind_MATCH_STATEMENT
}
