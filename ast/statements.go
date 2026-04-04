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

type BLangStatement = model.StatementNode

type (
	bLangStatementBase struct {
		bLangNodeBase
	}
	BLangAssignment struct {
		bLangStatementBase
		VarRef BLangExpression
		Expr   BLangExpression
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
		VarRef       model.ExpressionNode
		Expr         BLangExpression
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
		Expr BLangExpression
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
		Collection        BLangExpression
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
		Expr BLangExpression
	}

	BLangPanic struct {
		bLangStatementBase
		Expr BLangExpression
	}

	BLangMatchStatement struct {
		bLangStatementBase
		Expr         BLangExpression
		MatchClauses []BLangMatchClause
		IsExhaustive bool
	}
)

var (
	_ model.AssignmentNode          = &BLangAssignment{}
	_ model.CompoundAssignmentNode  = &BLangCompoundAssignment{}
	_ model.ContinueNode            = &BLangContinue{}
	_ model.DoNode                  = &BLangDo{}
	_ model.BlockStatementNode      = &BLangBlockStmt{}
	_ model.ExpressionStatementNode = &BLangExpressionStmt{}
	_ model.IfNode                  = &BLangIf{}
	_ model.WhileNode               = &BLangWhile{}
	_ model.ForeachNode             = &BLangForeach{}
	_ model.VariableDefinitionNode  = &BLangSimpleVariableDef{}
	_ model.ReturnNode              = &BLangReturn{}
	_ model.PanicNode               = &BLangPanic{}
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

func (b *BLangAssignment) GetVariable() model.ExpressionNode {
	// migrated from BLangAssignment.java:48:5
	return b.VarRef
}

func (b *BLangAssignment) GetExpression() model.ExpressionNode {
	// migrated from BLangAssignment.java:53:5
	return b.Expr
}

func (b *BLangAssignment) IsDeclaredWithVar() bool {
	// migrated from BLangAssignment.java:58:5
	return false
}

func (b *BLangAssignment) GetKind() model.NodeKind {
	return model.NodeKind_ASSIGNMENT
}

func (b *BLangAssignment) SetExpression(expression model.ExpressionNode) {
	// migrated from BLangAssignment.java:64:5
	if expr, ok := expression.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (b *BLangAssignment) SetDeclaredWithVar(isDeclaredWithVar bool) {
	// migrated from BLangAssignment.java:69:5
}

func (b *BLangAssignment) SetVariable(variableReferenceNode model.VariableReferenceNode) {
	// migrated from BLangAssignment.java:74:5
	if varRef, ok := variableReferenceNode.(BLangExpression); ok {
		b.VarRef = varRef
	} else {
		panic("variableReferenceNode is not a BLangExpression")
	}
}

func (b *BLangBlockStmt) GetKind() model.NodeKind {
	// migrated from BLangBlockStmt.java:83:5
	return model.NodeKind_BLOCK
}

func (b *BLangBlockStmt) GetStatements() []model.StatementNode {
	// migrated from BLangBlockStmt.java:88:5
	return b.Stmts
}

func (b *BLangBlockStmt) AddStatement(statement model.StatementNode) {
	// migrated from BLangBlockStmt.java:93:5
	b.Stmts = append(b.Stmts, statement)
}

func (b *BLangBreak) GetKind() model.NodeKind {
	// migrated from BLangBreak.java:45:5
	return model.NodeKind_BREAK
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

func (b *BLangCompoundAssignment) GetVariable() model.ExpressionNode {
	// migrated from BLangCompoundAssignment.java:64:5
	return b.VarRef
}

func (b *BLangCompoundAssignment) GetExpression() model.ExpressionNode {
	// migrated from BLangCompoundAssignment.java:69:5
	return b.Expr
}

func (b *BLangCompoundAssignment) SetExpression(expression model.ExpressionNode) {
	// migrated from BLangCompoundAssignment.java:74:5
	if exp, ok := expression.(BLangExpression); ok {
		b.Expr = exp
	} else {
		panic("Expected BLangExpression")
	}
}

func (b *BLangCompoundAssignment) SetVariable(variableReferenceNode model.VariableReferenceNode) {
	// migrated from BLangCompoundAssignment.java:79:5
	b.VarRef = variableReferenceNode
}

func (b *BLangCompoundAssignment) GetKind() model.NodeKind {
	// migrated from BLangCompoundAssignment.java:99:5
	return model.NodeKind_COMPOUND_ASSIGNMENT
}

func (b *BLangContinue) GetKind() model.NodeKind {
	// migrated from BLangContinue.java:46:5
	return model.NodeKind_NEXT
}

func (b *BLangDo) GetBody() model.BlockStatementNode {
	// migrated from BLangDo.java:47:5
	return &b.Body
}

func (b *BLangDo) SetBody(body model.BlockStatementNode) {
	// migrated from BLangDo.java:52:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangDo) GetOnFailClause() model.OnFailClauseNode {
	// migrated from BLangDo.java:57:5
	return &b.OnFailClause
}

func (b *BLangDo) SetOnFailClause(onFailClause model.OnFailClauseNode) {
	// migrated from BLangDo.java:62:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		b.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (b *BLangDo) GetKind() model.NodeKind {
	// migrated from BLangDo.java:82:5
	return model.NodeKind_DO_STMT
}

func (b *BLangExpressionStmt) GetExpression() model.ExpressionNode {
	// migrated from BLangExpressionStmt.java:46:5
	return b.Expr
}

func (b *BLangExpressionStmt) GetKind() model.NodeKind {
	return model.NodeKind_EXPRESSION_STATEMENT
}

func (b *BLangIf) Scope() model.Scope {
	return b.scope
}

func (b *BLangIf) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangIf) GetCondition() model.ExpressionNode {
	// migrated from BLangIf.java:47:5
	return b.Expr
}

func (b *BLangIf) GetBody() model.BlockStatementNode {
	// migrated from BLangIf.java:52:5
	return &b.Body
}

func (b *BLangIf) GetElseStatement() model.StatementNode {
	// migrated from BLangIf.java:57:5
	return b.ElseStmt
}

func (b *BLangIf) SetCondition(condition model.ExpressionNode) {
	// migrated from BLangIf.java:62:5
	if expr, ok := condition.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("condition is not a BLangExpression")
	}
}

func (b *BLangIf) SetBody(body model.BlockStatementNode) {
	// migrated from BLangIf.java:67:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangIf) SetElseStatement(elseStatement model.StatementNode) {
	// migrated from BLangIf.java:72:5
	b.ElseStmt = elseStatement
}

func (b *BLangIf) GetKind() model.NodeKind {
	// migrated from BLangIf.java:77:5
	return model.NodeKind_IF
}

func (b *BLangWhile) Scope() model.Scope {
	return b.scope
}

func (b *BLangWhile) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangWhile) GetCondition() model.ExpressionNode {
	// migrated from BLangWhile.java:50:5
	return b.Expr
}

func (b *BLangWhile) SetCondition(condition model.ExpressionNode) {
	// migrated from BLangWhile.java:60:5
	if expr, ok := condition.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("condition is not a BLangExpression")
	}
}

func (b *BLangWhile) GetBody() model.BlockStatementNode {
	// migrated from BLangWhile.java:55:5
	return &b.Body
}

func (b *BLangWhile) SetBody(body model.BlockStatementNode) {
	// migrated from BLangWhile.java:65:5
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangWhile) GetOnFailClause() model.OnFailClauseNode {
	// migrated from BLangWhile.java:70:5
	return &b.OnFailClause
}

func (b *BLangWhile) SetOnFailClause(onFailClause model.OnFailClauseNode) {
	// migrated from BLangWhile.java:75:5
	if onFailClause, ok := onFailClause.(*BLangOnFailClause); ok {
		b.OnFailClause = *onFailClause
		return
	}
	panic("onFailClause is not a BLangOnFailClause")
}

func (b *BLangWhile) GetKind() model.NodeKind {
	// migrated from BLangWhile.java:95:5
	return model.NodeKind_WHILE
}

func (b *BLangForeach) Scope() model.Scope {
	return b.scope
}

func (b *BLangForeach) SetScope(scope model.Scope) {
	b.scope = scope
}

func (b *BLangForeach) GetKind() model.NodeKind {
	return model.NodeKind_FOREACH
}

func (b *BLangForeach) GetVariableDefinitionNode() model.VariableDefinitionNode {
	return b.VariableDef
}

func (b *BLangForeach) SetVariableDefinitionNode(node model.VariableDefinitionNode) {
	if varDef, ok := node.(*BLangSimpleVariableDef); ok {
		b.VariableDef = varDef
		return
	}
	panic("node is not a *BLangSimpleVariableDef")
}

func (b *BLangForeach) GetCollection() model.ExpressionNode {
	return b.Collection
}

func (b *BLangForeach) SetCollection(collection model.ExpressionNode) {
	if expr, ok := collection.(BLangExpression); ok {
		b.Collection = expr
	} else {
		panic("collection is not a BLangExpression")
	}
}

func (b *BLangForeach) GetBody() model.BlockStatementNode {
	return &b.Body
}

func (b *BLangForeach) SetBody(body model.BlockStatementNode) {
	if blockStmt, ok := body.(*BLangBlockStmt); ok {
		b.Body = *blockStmt
		return
	}
	panic("body is not a BLangBlockStmt")
}

func (b *BLangForeach) GetIsDeclaredWithVar() bool {
	return b.IsDeclaredWithVar
}

func (b *BLangForeach) GetOnFailClause() model.OnFailClauseNode {
	if b.OnFailClause == nil {
		return nil
	}
	return b.OnFailClause
}

func (b *BLangForeach) SetOnFailClause(onFailClause model.OnFailClauseNode) {
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

func (b *BLangSimpleVariableDef) GetKind() model.NodeKind {
	return model.NodeKind_VARIABLE_DEF
}

func (b *BLangSimpleVariableDef) GetVariable() model.VariableNode {
	return b.Var
}

func (b *BLangSimpleVariableDef) SetVariable(variable model.VariableNode) {
	if v, ok := variable.(*BLangSimpleVariable); ok {
		b.Var = v
	} else {
		panic("variable is not a BLangSimpleVariable")
	}
}

func (b *BLangReturn) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangReturn) SetExpression(expression model.ExpressionNode) {
	if expr, ok := expression.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("expression is not a BLangExpression")
	}
}

func (b *BLangReturn) GetKind() model.NodeKind {
	return model.NodeKind_RETURN
}

func (b *BLangPanic) GetExpression() model.ExpressionNode {
	return b.Expr
}

func (b *BLangPanic) GetKind() model.NodeKind {
	return model.NodeKind_PANIC
}
func (b *BLangMatchStatement) GetKind() model.NodeKind {
	return model.NodeKind_MATCH_STATEMENT
}
