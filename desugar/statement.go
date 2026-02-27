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
	array "ballerina-lang-go/lib/array/compile"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

func walkStatement(cx *FunctionContext, node model.StatementNode) desugaredNode[model.StatementNode] {
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

func walkBlockStmt(cx *FunctionContext, stmt *ast.BLangBlockStmt) desugaredNode[model.StatementNode] {
	var allStmts []model.StatementNode

	for _, childStmt := range stmt.Stmts {
		result := walkStatement(cx, childStmt)
		allStmts = append(allStmts, result.initStmts...)
		allStmts = append(allStmts, result.replacementNode)
	}

	stmt.Stmts = allStmts
	return desugaredNode[model.StatementNode]{replacementNode: stmt}
}

func walkBlockFunctionBody(cx *FunctionContext, body *ast.BLangBlockFunctionBody) {
	var allStmts []ast.BLangStatement

	for _, stmt := range body.Stmts {
		result := walkStatement(cx, stmt.(model.StatementNode))
		for _, initStmt := range result.initStmts {
			allStmts = append(allStmts, initStmt.(ast.BLangStatement))
		}
		allStmts = append(allStmts, result.replacementNode.(ast.BLangStatement))
	}

	body.Stmts = allStmts
}

func walkAssignment(cx *FunctionContext, stmt *ast.BLangAssignment) desugaredNode[model.StatementNode] {
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

func walkCompoundAssignment(cx *FunctionContext, stmt *ast.BLangCompoundAssignment) desugaredNode[model.StatementNode] {
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

func walkExpressionStmt(cx *FunctionContext, stmt *ast.BLangExpressionStmt) desugaredNode[model.StatementNode] {
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

func walkIf(cx *FunctionContext, stmt *ast.BLangIf) desugaredNode[model.StatementNode] {
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
		if len(elseResult.initStmts) > 0 {
			stmt.ElseStmt = &ast.BLangBlockStmt{
				Stmts: append(elseResult.initStmts, elseResult.replacementNode),
			}
		} else {
			stmt.ElseStmt = elseResult.replacementNode.(ast.BLangStatement)
		}
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkWhile(cx *FunctionContext, stmt *ast.BLangWhile) desugaredNode[model.StatementNode] {
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
		walkOnFailClause(cx, &stmt.OnFailClause)
	}

	return desugaredNode[model.StatementNode]{
		initStmts:       initStmts,
		replacementNode: stmt,
	}
}

func walkDo(cx *FunctionContext, stmt *ast.BLangDo) desugaredNode[model.StatementNode] {
	bodyResult := walkBlockStmt(cx, &stmt.Body)
	stmt.Body = *bodyResult.replacementNode.(*ast.BLangBlockStmt)

	// Only walk onFail clause if it has a body
	if stmt.OnFailClause.Body != nil {
		walkOnFailClause(cx, &stmt.OnFailClause)
	}

	return desugaredNode[model.StatementNode]{
		replacementNode: stmt,
	}
}

func walkOnFailClause(cx *FunctionContext, clause *ast.BLangOnFailClause) {
	bodyResult := walkBlockStmt(cx, clause.Body)
	clause.Body = bodyResult.replacementNode.(*ast.BLangBlockStmt)
}

func walkSimpleVariableDef(cx *FunctionContext, stmt *ast.BLangSimpleVariableDef) desugaredNode[model.StatementNode] {
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

func walkReturn(cx *FunctionContext, stmt *ast.BLangReturn) desugaredNode[model.StatementNode] {
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

func walkContinue(cx *FunctionContext, stmt *ast.BLangContinue) desugaredNode[model.StatementNode] {
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

func visitForEach(cx *FunctionContext, stmt *ast.BLangForeach) desugaredNode[model.StatementNode] {
	cx.pushScope(stmt.Scope())
	defer cx.popScope()
	if isRangeExpr(stmt.Collection) {
		rangeExpr := stmt.Collection.(*ast.BLangBinaryExpr)
		return desugarForEachOnRange(cx, rangeExpr, stmt.VariableDef, &stmt.Body, stmt.Scope())
	}
	if semtypes.IsSubtypeSimple(stmt.Collection.GetDeterminedType(), semtypes.LIST) {
		return desugarForEachOnList(cx, stmt.Collection, stmt.VariableDef, &stmt.Body, stmt.Scope())
	}
	cx.unimplemented("unsupported collection type in foreach")
	return desugaredNode[model.StatementNode]{}
}

func desugarForEachOnList(cx *FunctionContext, collection ast.BLangExpression, loopVarDef *ast.BLangSimpleVariableDef, body *ast.BLangBlockStmt, foreachScope model.Scope) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	// Step 1: evaluate collection once into a temp variable
	collResult := walkExpression(cx, collection)
	initStmts = append(initStmts, collResult.initStmts...)
	collExpr := collResult.replacementNode

	collType := collExpr.GetDeterminedType()
	collName, collVarSymbol := cx.addDesugardSymbol(collType, model.SymbolKindVariable, false)
	collVarName := &ast.BLangIdentifier{Value: collName}
	collVar := &ast.BLangSimpleVariable{Name: collVarName}
	collVar.SetDeterminedType(collType)
	collVar.SetInitialExpression(collExpr)
	collVar.SetSymbol(collVarSymbol)
	collVarDef := &ast.BLangSimpleVariableDef{Var: collVar}
	initStmts = append(initStmts, collVarDef)

	collVarRef := &ast.BLangSimpleVarRef{VariableName: collVarName}
	collVarRef.SetSymbol(collVarSymbol)
	collVarRef.SetDeterminedType(collType)

	// Step 2: index variable ($desugar$N = 0)
	zeroLiteral := &ast.BLangNumericLiteral{
		BLangLiteral: ast.BLangLiteral{
			Value:         int64(0),
			OriginalValue: "0",
		},
		Kind: model.NodeKind_NUMERIC_LITERAL,
	}
	zeroLiteral.SetDeterminedType(&semtypes.INT)

	idxName, idxVarSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	idxVarName := &ast.BLangIdentifier{Value: idxName}
	idxVar := &ast.BLangSimpleVariable{Name: idxVarName}
	idxVar.SetDeterminedType(&semtypes.INT)
	idxVar.SetInitialExpression(zeroLiteral)
	idxVar.SetSymbol(idxVarSymbol)
	idxVarDef := &ast.BLangSimpleVariableDef{Var: idxVar}
	initStmts = append(initStmts, idxVarDef)

	idxVarRef := &ast.BLangSimpleVarRef{VariableName: idxVarName}
	idxVarRef.SetSymbol(idxVarSymbol)
	idxVarRef.SetDeterminedType(&semtypes.INT)

	// Step 3: length variable ($desugar$M = length(collVar))
	lengthInvocation := createLengthInvocation(cx, collVarRef)

	lenName, lenVarSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	lenVarName := &ast.BLangIdentifier{Value: lenName}
	lenVar := &ast.BLangSimpleVariable{Name: lenVarName}
	lenVar.SetDeterminedType(&semtypes.INT)
	lenVar.SetInitialExpression(lengthInvocation)
	lenVar.SetSymbol(lenVarSymbol)
	lenVarDef := &ast.BLangSimpleVariableDef{Var: lenVar}
	initStmts = append(initStmts, lenVarDef)

	lenVarRef := &ast.BLangSimpleVarRef{VariableName: lenVarName}
	lenVarRef.SetSymbol(lenVarSymbol)
	lenVarRef.SetDeterminedType(&semtypes.INT)

	// Step 4: while condition ($idx < $len)
	whileCondition := &ast.BLangBinaryExpr{
		LhsExpr: idxVarRef,
		RhsExpr: lenVarRef,
		OpKind:  model.OperatorKind_LESS_THAN,
	}
	whileCondition.SetDeterminedType(&semtypes.BOOLEAN)

	// Step 5: element access (collVar[$idx])
	elementAccess := &ast.BLangIndexBasedAccess{
		BLangAccessExpressionBase: ast.BLangAccessExpressionBase{
			Expr: collVarRef,
		},
		IndexExpr: idxVarRef,
	}
	elementAccess.SetDeterminedType(loopVarDef.Var.GetDeterminedType())

	// Step 6: patch loop var def initial expression
	loopVarDef.Var.SetInitialExpression(elementAccess)

	// Step 7: build body
	incrementStmt := createIncrementStmt(idxVarRef)
	cx.pushLoopVar(idxVarRef)

	newBodyStmts := make([]model.StatementNode, 0, len(body.Stmts)+2)
	newBodyStmts = append(newBodyStmts, loopVarDef)
	newBodyStmts = append(newBodyStmts, body.Stmts...)
	if len(newBodyStmts) > 0 {
		if isAppentReachable(newBodyStmts[len(newBodyStmts)-1]) {
			newBodyStmts = append(newBodyStmts, incrementStmt)
		}
	}
	body.Stmts = newBodyStmts

	bodyResult := walkBlockStmt(cx, body)
	newBody := bodyResult.replacementNode.(*ast.BLangBlockStmt)

	cx.popLoopVar()

	// Step 8: create while loop
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

func createLengthInvocation(cx *FunctionContext, collection ast.BLangExpression) *ast.BLangInvocation {
	pkgName := array.PackageName
	space, ok := cx.getImportedSymbolSpace(pkgName)
	if !ok {
		cx.internalError(pkgName + " symbol space not found")
		return nil
	}
	symbolRef, ok := space.GetSymbol("length")
	if !ok {
		cx.internalError(pkgName + ":length symbol not found")
		return nil
	}
	cx.addImplicitImport(pkgName, ast.BLangImportPackage{
		OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
		PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
		Alias:        &ast.BLangIdentifier{Value: pkgName},
	})
	inv := &ast.BLangInvocation{
		Name:     &ast.BLangIdentifier{Value: "length"},
		PkgAlias: &ast.BLangIdentifier{Value: pkgName},
		ArgExprs: []ast.BLangExpression{collection},
	}
	inv.SetSymbol(symbolRef)
	inv.SetDeterminedType(&semtypes.INT)
	return inv
}

func desugarForEachOnRange(cx *FunctionContext, rangeExpr *ast.BLangBinaryExpr, loopVarDef *ast.BLangSimpleVariableDef, body *ast.BLangBlockStmt, foreachScope model.Scope) desugaredNode[model.StatementNode] {
	var initStmts []model.StatementNode

	startResult := walkExpression(cx, rangeExpr.LhsExpr)
	initStmts = append(initStmts, startResult.initStmts...)
	startExpr := startResult.replacementNode

	endResult := walkExpression(cx, rangeExpr.RhsExpr)
	initStmts = append(initStmts, endResult.initStmts...)
	endExpr := endResult.replacementNode

	loopVarDef.Var.SetInitialExpression(startExpr)
	initStmts = append(initStmts, loopVarDef)

	loopVarRef := &ast.BLangSimpleVarRef{
		VariableName: loopVarDef.Var.Name,
	}
	loopVarRef.SetSymbol(loopVarDef.Var.Symbol())

	endName, endVarSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	endVarName := &ast.BLangIdentifier{Value: endName}
	endVar := &ast.BLangSimpleVariable{Name: endVarName}
	endVar.SetDeterminedType(&semtypes.INT)
	endVar.SetInitialExpression(endExpr)
	endVar.SetSymbol(endVarSymbol)

	endVarDef := &ast.BLangSimpleVariableDef{
		Var: endVar,
	}
	initStmts = append(initStmts, endVarDef)

	endVarRef := &ast.BLangSimpleVarRef{
		VariableName: endVarName,
	}
	endVarRef.SetSymbol(endVarSymbol)

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
	if len(newBodyStmts) > 0 {
		if isAppentReachable(newBodyStmts[len(newBodyStmts)-1]) {
			newBodyStmts = append(newBodyStmts, incrementStmt)
		}
	} else {
		// just replace it with a no-op
		return desugaredNode[model.StatementNode]{
			replacementNode: &ast.BLangBlockStmt{},
		}
	}
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

// TODO: do we need to think about if-else here as well?
// If the last statement in a block is something like panic, return, continue or break, then we shouldn't append
// nodes after that. I would make that node unreacheable. We need to make sure desugared AST is still valid.
func isAppentReachable(stmt ast.BLangStatement) bool {
	switch stmt := stmt.(type) {
	case *ast.BLangReturn, *ast.BLangContinue, *ast.BLangBreak:
		return false
	case *ast.BLangBlockStmt:
		if len(stmt.Stmts) == 0 {
			return true
		}
		lastChild := stmt.Stmts[len(stmt.Stmts)-1]
		return isAppentReachable(lastChild)
	default:
		return true
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
