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

// Package desugar represents AST-> AST transforms
package desugar

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

func walkStatement(cx *Context, node model.StatementNode) desugaredNode[model.StatementNode] {
	switch stmt := node.(type) {
	case *ast.BLangBlockStmt:
		return walkBlockStmt(cx, stmt)
	case *ast.BLangAssignment:
		return walkAssignment(cx, stmt)
	case *ast.BLangCompoundAssignment:
		return walkCompoundAssignment(cx, stmt)
	case *ast.BLangExpressionStmt:
		return walkExpressionStmt(cx, stmt)
	case *ast.BLangIf:
		return walkIf(cx, stmt)
	case *ast.BLangWhile:
		return walkWhile(cx, stmt)
	case *ast.BLangDo:
		return walkDo(cx, stmt)
	case *ast.BLangForeach:
		return visitForEach(cx, stmt)
	case *ast.BLangSimpleVariableDef:
		return walkSimpleVariableDef(cx, stmt)
	case *ast.BLangReturn:
		return walkReturn(cx, stmt)
	case *ast.BLangBreak:
		return desugaredNode[model.StatementNode]{replacementNode: stmt}
	case *ast.BLangContinue:
		return walkContinue(cx, stmt)
	default:
		panic("unexpected statement type")
	}
}

func walkBlockStmt(cx *Context, stmt *ast.BLangBlockStmt) desugaredNode[model.StatementNode] {
	var allStmts []model.StatementNode

	for _, childStmt := range stmt.Stmts {
		result := walkStatement(cx, childStmt)
		allStmts = append(allStmts, result.initStmts...)
		allStmts = append(allStmts, result.replacementNode)
	}

	stmt.Stmts = allStmts
	return desugaredNode[model.StatementNode]{replacementNode: stmt}
}

func walkBlockFunctionBody(cx *Context, body *ast.BLangBlockFunctionBody) desugaredNode[model.StatementNode] {
	var allStmts []ast.BLangStatement

	for _, stmt := range body.Stmts {
		result := walkStatement(cx, stmt.(model.StatementNode))
		for _, initStmt := range result.initStmts {
			allStmts = append(allStmts, initStmt.(ast.BLangStatement))
		}
		allStmts = append(allStmts, result.replacementNode.(ast.BLangStatement))
	}

	body.Stmts = allStmts
	return desugaredNode[model.StatementNode]{replacementNode: body}
}

func walkAssignment(cx *Context, stmt *ast.BLangAssignment) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.VarRef != nil {
		result := walkExpression(cx, stmt.VarRef)
		initStmts = append(initStmts, result.initStmts...)
		stmt.VarRef = result.replacementNode.(ast.BLangExpression)
	}

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkCompoundAssignment(cx *Context, stmt *ast.BLangCompoundAssignment) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.VarRef != nil {
		result := walkExpression(cx, stmt.VarRef.(ast.BLangExpression))
		initStmts = append(initStmts, result.initStmts...)
		stmt.VarRef = result.replacementNode.(ast.BLangExpression)
	}

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkExpressionStmt(cx *Context, stmt *ast.BLangExpressionStmt) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkIf(cx *Context, stmt *ast.BLangIf) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	// Push if scope before visiting body
	cx.pushScope(stmt.Scope())
	bodyResult := walkBlockStmt(cx, &stmt.Body)
	stmt.Body = *bodyResult.replacementNode.(*ast.BLangBlockStmt)
	cx.popScope()

	if stmt.ElseStmt != nil {
		elseResult := walkStatement(cx, stmt.ElseStmt)
		stmt.ElseStmt = elseResult.replacementNode.(ast.BLangStatement)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkWhile(cx *Context, stmt *ast.BLangWhile) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	// Push nil to loopVarStack to indicate this is a native while (not desugared foreach)
	cx.pushLoopVar(nil)
	// Push while scope before visiting body
	cx.pushScope(stmt.Scope())
	bodyResult := walkBlockStmt(cx, &stmt.Body)
	stmt.Body = *bodyResult.replacementNode.(*ast.BLangBlockStmt)
	cx.popScope()
	cx.popLoopVar()

	// Only walk onFail clause if it has a body
	if stmt.OnFailClause.Body != nil {
		onFailResult := walkOnFailClause(cx, &stmt.OnFailClause)
		stmt.OnFailClause = *onFailResult.replacementNode.(*ast.BLangOnFailClause)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkDo(cx *Context, stmt *ast.BLangDo) desugaredNode[model.StatementNode] {
	bodyResult := walkBlockStmt(cx, &stmt.Body)
	stmt.Body = *bodyResult.replacementNode.(*ast.BLangBlockStmt)

	// Only walk onFail clause if it has a body
	if stmt.OnFailClause.Body != nil {
		onFailResult := walkOnFailClause(cx, &stmt.OnFailClause)
		stmt.OnFailClause = *onFailResult.replacementNode.(*ast.BLangOnFailClause)
	}

	return desugaredNode[model.StatementNode]{
		replacementNode: stmt,
	}
}

func walkOnFailClause(cx *Context, clause *ast.BLangOnFailClause) desugaredNode[model.StatementNode] {
	bodyResult := walkBlockStmt(cx, clause.Body)
	clause.Body = bodyResult.replacementNode.(*ast.BLangBlockStmt)

	return desugaredNode[model.StatementNode]{
		replacementNode: clause,
	}
}

func walkSimpleVariableDef(cx *Context, stmt *ast.BLangSimpleVariableDef) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.Var != nil && stmt.Var.Expr != nil {
		result := walkExpression(cx, stmt.Var.Expr.(ast.BLangExpression))
		initStmts = append(initStmts, result.initStmts...)
		stmt.Var.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkReturn(cx *Context, stmt *ast.BLangReturn) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	if stmt.Expr != nil {
		result := walkExpression(cx, stmt.Expr)
		initStmts = append(initStmts, result.initStmts...)
		stmt.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func createIncrementStmt(loopVar ast.BLangExpression) *ast.BLangAssignment {
	oneLiteral := &ast.BLangNumericLiteral{
		BLangLiteral: ast.BLangLiteral{
			Value:         int64(1),
			OriginalValue: "1",
		},
		Kind: model.NodeKind_NUMERIC_LITERAL,
	}
	oneLiteral.SetDeterminedType(&semtypes.INT)
	addExpr := &ast.BLangBinaryExpr{
		LhsExpr: loopVar,
		RhsExpr: oneLiteral,
		OpKind:  model.OperatorKind_ADD,
	}
	addExpr.SetDeterminedType(&semtypes.INT)
	incrementStmt := &ast.BLangAssignment{
		VarRef: loopVar,
		Expr:   addExpr,
	}
	incrementStmt.SetDeterminedType(&semtypes.NEVER)
	return incrementStmt
}

func walkContinue(cx *Context, stmt *ast.BLangContinue) desugaredNode[model.StatementNode] {
	// Check if we're in a desugared foreach (has a loop variable)
	loopVar := cx.currentLoopVar()
	if loopVar != nil {
		// For desugared foreach, we need to add increment before continue
		incrementStmt := createIncrementStmt(loopVar)

		// Return increment as initStmts and continue as replacement
		return desugaredNode[model.StatementNode]{
			initStmts:       []model.StatementNode{incrementStmt},
			replacementNode: stmt,
		}
	}

	// For native while loops, continue as-is
	return desugaredNode[model.StatementNode]{
		initStmts:       []model.StatementNode{},
		replacementNode: stmt,
	}
}

func visitForEach(cx *Context, stmt *ast.BLangForeach) desugaredNode[model.StatementNode] {
	cx.pushScope(stmt.Scope())
	defer cx.popScope()
	if isRangeExpr(stmt.Collection) {
		rangeExpr := stmt.Collection.(*ast.BLangBinaryExpr)
		return desugarForEachOnRange(cx, rangeExpr, stmt.VariableDef, &stmt.Body, stmt.Scope())
	}
	cx.compilerCtx.Unimplemented("unsupported collection type in foreach", nil)
	return desugaredNode[model.StatementNode]{}
}

func desugarForEachOnRange(cx *Context, rangeExpr *ast.BLangBinaryExpr, loopVarDef *ast.BLangSimpleVariableDef, body *ast.BLangBlockStmt, foreachScope model.Scope) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	startExpr := rangeExpr.LhsExpr
	endExpr := rangeExpr.RhsExpr

	loopVarDef.Var.SetInitialExpression(startExpr)
	initStmts = append(initStmts, loopVarDef)

	loopVarRef := &ast.BLangSimpleVarRef{
		VariableName: loopVarDef.Var.Name,
	}
	loopVarRef.SetSymbol(loopVarDef.Var.Symbol())

	endVarName := &ast.BLangIdentifier{
		Value: cx.nextDesugarSymbolName(),
	}
	endVar := &ast.BLangSimpleVariable{
		Name: endVarName,
	}
	endVar.SetDeterminedType(&semtypes.INT)
	endVar.SetInitialExpression(endExpr)
	endVarSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	endVar.SetSymbol(&endVarSymbol)

	endVarDef := &ast.BLangSimpleVariableDef{
		Var: endVar,
	}
	initStmts = append(initStmts, endVarDef)

	endVarRef := &ast.BLangSimpleVarRef{
		VariableName: endVarName,
	}
	endVarRef.SetSymbol(&endVarSymbol)

	var compOp model.OperatorKind
	if rangeExpr.GetOperatorKind() == model.OperatorKind_CLOSED_RANGE {
		compOp = model.OperatorKind_LESS_EQUAL // <= for closed range
	} else {
		compOp = model.OperatorKind_LESS_THAN // < for half-open range
	}

	whileCondition := &ast.BLangBinaryExpr{
		LhsExpr: loopVarRef,
		RhsExpr: endVarRef,
		OpKind:  compOp,
	}
	whileCondition.SetDeterminedType(&semtypes.BOOLEAN)

	incrementStmt := createIncrementStmt(loopVarRef)

	// Note: foreach scope is already pushed by visitForEach at the top level
	cx.pushLoopVar(loopVarRef)

	newBodyStmts := make([]model.StatementNode, len(body.Stmts))
	copy(newBodyStmts, body.Stmts)
	newBodyStmts = append(newBodyStmts, incrementStmt)
	body.Stmts = newBodyStmts

	bodyResult := walkBlockStmt(cx, body)
	newBody := bodyResult.replacementNode.(*ast.BLangBlockStmt)

	cx.popLoopVar()

	// 10. Create the while loop using the foreach scope
	whileStmt := &ast.BLangWhile{
		Expr: whileCondition,
		Body: *newBody,
	}
	whileStmt.SetScope(foreachScope)
	whileStmt.SetDeterminedType(&semtypes.NEVER)

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: whileStmt,
	}
}

func isRangeExpr(expr ast.BLangExpression) bool {
	if binaryExpr, ok := expr.(*ast.BLangBinaryExpr); ok {
		switch binaryExpr.GetOperatorKind() {
		case model.OperatorKind_CLOSED_RANGE, model.OperatorKind_HALF_OPEN_RANGE:
			return true
		default:
			return false
		}
	}
	return false
}
