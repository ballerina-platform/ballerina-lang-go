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
	"fmt"
)

func walkExpression(cx *FunctionContext, node model.ExpressionNode) desugaredNode[model.ExpressionNode] {
	switch expr := node.(type) {
	case *ast.BLangBinaryExpr:
		return walkBinaryExpr(cx, expr)
	case *ast.BLangUnaryExpr:
		return walkUnaryExpr(cx, expr)
	case *ast.BLangElvisExpr:
		return walkElvisExpr(cx, expr)
	case *ast.BLangGroupExpr:
		return walkGroupExpr(cx, expr)
	case *ast.BLangIndexBasedAccess:
		return walkIndexBasedAccess(cx, expr)
	case *ast.BLangFieldBaseAccess:
		return walkFieldBaseAccess(cx, expr)
	case *ast.BLangInvocation:
		return walkInvocation(cx, expr)
	case *ast.BLangListConstructorExpr:
		return walkListConstructorExpr(cx, expr)
	case *ast.BLangMappingConstructorExpr:
		return walkMappingConstructorExpr(cx, expr)
	case *ast.BLangErrorConstructorExpr:
		return walkErrorConstructorExpr(cx, expr)
	case *ast.BLangCheckedExpr:
		return walkCheckedExpr(cx, expr)
	case *ast.BLangCheckPanickedExpr:
		return walkCheckPanickedExpr(cx, expr)
	case *ast.BLangTrapExpr:
		return walkTrapExpr(cx, expr)
	case *ast.BLangDynamicArgExpr:
		return walkDynamicArgExpr(cx, expr)
	case *ast.BLangLambdaFunction:
		return walkLambdaFunction(cx, expr)
	case *ast.BLangTypeConversionExpr:
		return walkTypeConversionExpr(cx, expr)
	case *ast.BLangTypeTestExpr:
		return walkTypeTestExpr(cx, expr)
	case *ast.BLangAnnotAccessExpr:
		return walkAnnotAccessExpr(cx, expr)
	case *ast.BLangCollectContextInvocation:
		return walkCollectContextInvocation(cx, expr)
	case *ast.BLangArrowFunction:
		return walkArrowFunction(cx, expr)
	case *ast.BLangQueryExpr:
		return walkQueryExpr(cx, expr)
	case *ast.BLangLiteral:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangNumericLiteral:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangSimpleVarRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangLocalVarRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangConstRef:
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	case *ast.BLangWildCardBindingPattern:
		// Wildcard binding pattern can appear in variable references (e.g., _ = expr)
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", node))
	}
}

func walkBinaryExpr(cx *FunctionContext, expr *ast.BLangBinaryExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.LhsExpr != nil {
		result := walkExpression(cx, expr.LhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.LhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.RhsExpr != nil {
		result := walkExpression(cx, expr.RhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.RhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkUnaryExpr(cx *FunctionContext, expr *ast.BLangUnaryExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkElvisExpr(cx *FunctionContext, expr *ast.BLangElvisExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.LhsExpr != nil {
		result := walkExpression(cx, expr.LhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.LhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.RhsExpr != nil {
		result := walkExpression(cx, expr.RhsExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.RhsExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkGroupExpr(cx *FunctionContext, expr *ast.BLangGroupExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkIndexBasedAccess(cx *FunctionContext, expr *ast.BLangIndexBasedAccess) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	if expr.IndexExpr != nil {
		result := walkExpression(cx, expr.IndexExpr)
		initStmts = append(initStmts, result.initStmts...)
		expr.IndexExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkFieldBaseAccess(cx *FunctionContext, expr *ast.BLangFieldBaseAccess) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	name := expr.Field.Value
	lit := &ast.BLangLiteral{
		Value:         name,
		OriginalValue: name,
	}
	s := semtypes.STRING
	lit.SetDeterminedType(&s)

	indexAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: lit,
	}
	indexAccess.Expr = expr.Expr
	indexAccess.SetDeterminedType(expr.GetDeterminedType())

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: indexAccess,
	}
}

func walkInvocation(cx *FunctionContext, expr *ast.BLangInvocation) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	for i := range expr.ArgExprs {
		result := walkExpression(cx, expr.ArgExprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.ArgExprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkListConstructorExpr(cx *FunctionContext, expr *ast.BLangListConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	for i := range expr.Exprs {
		result := walkExpression(cx, expr.Exprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.Exprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkErrorConstructorExpr(cx *FunctionContext, expr *ast.BLangErrorConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	//nolint:staticcheck // TODO
	if expr.ErrorTypeRef != nil {
		// ErrorTypeRef is a type descriptor, not an expression, so we don't walk it
	}

	for i := range expr.PositionalArgs {
		result := walkExpression(cx, expr.PositionalArgs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.PositionalArgs[i] = result.replacementNode.(ast.BLangExpression)
	}

	for i := range expr.NamedArgs {
		result := walkExpression(cx, expr.NamedArgs[i].Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.NamedArgs[i].Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCheckedExpr(cx *FunctionContext, expr *ast.BLangCheckedExpr) desugaredNode[model.ExpressionNode] {
	return desugarCheckedExpr(cx, expr, false)
}

func walkCheckPanickedExpr(cx *FunctionContext, expr *ast.BLangCheckPanickedExpr) desugaredNode[model.ExpressionNode] {
	return desugarCheckedExpr(cx, &expr.BLangCheckedExpr, true)
}

func walkTrapExpr(cx *FunctionContext, expr *ast.BLangTrapExpr) desugaredNode[model.ExpressionNode] {
	result := walkExpression(cx, expr.Expr)
	if len(result.initStmts) > 0 {
		// I don't think this can ever happen but if it does we need to think about how to add these statements in to the
		// trap region in BIR gen
		cx.internalError("Init statements will be hoisted outside of trap region")
	}
	expr.Expr = result.replacementNode.(ast.BLangExpression)
	return desugaredNode[model.ExpressionNode]{initStmts: nil, replacementNode: expr}
}

func desugarCheckedExpr(cx *FunctionContext, expr *ast.BLangCheckedExpr, isPanic bool) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	// Walk the inner expression first
	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	innerTy := expr.Expr.GetDeterminedType()
	resultTy := expr.GetDeterminedType()

	// TODO: extract util to add definition and get reference
	// Create temp var: $desugar$N = <inner expr>
	tempName, tempSymbol := cx.addDesugardSymbol(innerTy, model.SymbolKindVariable, false)
	tempVarName := &ast.BLangIdentifier{Value: tempName}
	tempVar := &ast.BLangSimpleVariable{Name: tempVarName}
	tempVar.SetDeterminedType(innerTy)
	tempVar.SetInitialExpression(expr.Expr)
	tempVar.SetSymbol(tempSymbol)
	tempVarDef := &ast.BLangSimpleVariableDef{Var: tempVar}
	initStmts = append(initStmts, tempVarDef)

	// Type test: $desugar$N is error
	tempVarRefForTest := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	tempVarRefForTest.SetSymbol(tempSymbol)
	tempVarRefForTest.SetDeterminedType(innerTy)

	typeTestExpr := &ast.BLangTypeTestExpr{}
	typeTestExpr.Expr = tempVarRefForTest
	typeTestExpr.Type = model.TypeData{Type: &semtypes.ERROR}
	typeTestExpr.SetDeterminedType(&semtypes.BOOLEAN)

	// If body: return or panic
	tempVarRefForBody := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	tempVarRefForBody.SetSymbol(tempSymbol)
	tempVarRefForBody.SetDeterminedType(innerTy)

	var bodyStmt ast.BLangStatement
	if isPanic {
		bodyStmt = &ast.BLangPanic{Expr: tempVarRefForBody}
	} else {
		bodyStmt = &ast.BLangReturn{Expr: tempVarRefForBody}
	}

	ifStmt := &ast.BLangIf{
		Expr: typeTestExpr,
		Body: ast.BLangBlockStmt{
			Stmts: []ast.BLangStatement{bodyStmt},
		},
	}
	initStmts = append(initStmts, ifStmt)

	// Replacement: var ref typed as non-error type
	replacementVarRef := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	replacementVarRef.SetSymbol(tempSymbol)
	replacementVarRef.SetDeterminedType(resultTy)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: replacementVarRef,
	}
}

func walkDynamicArgExpr(cx *FunctionContext, expr *ast.BLangDynamicArgExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Condition != nil {
		result := walkExpression(cx, expr.Condition)
		initStmts = append(initStmts, result.initStmts...)
		expr.Condition = result.replacementNode.(ast.BLangExpression)
	}

	if expr.ConditionalArgument != nil {
		result := walkExpression(cx, expr.ConditionalArgument)
		initStmts = append(initStmts, result.initStmts...)
		expr.ConditionalArgument = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkLambdaFunction(cx *FunctionContext, expr *ast.BLangLambdaFunction) desugaredNode[model.ExpressionNode] {
	// Desugar the function body
	if expr.Function != nil {
		expr.Function = desugarFunction(cx.pkgCtx, expr.Function)
	}

	return desugaredNode[model.ExpressionNode]{
		replacementNode: expr,
	}
}

func walkTypeConversionExpr(cx *FunctionContext, expr *ast.BLangTypeConversionExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkTypeTestExpr(cx *FunctionContext, expr *ast.BLangTypeTestExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkAnnotAccessExpr(cx *FunctionContext, expr *ast.BLangAnnotAccessExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCollectContextInvocation(cx *FunctionContext, expr *ast.BLangCollectContextInvocation) desugaredNode[model.ExpressionNode] {
	// Walk the underlying invocation
	result := walkInvocation(cx, &expr.Invocation)
	expr.Invocation = *result.replacementNode.(*ast.BLangInvocation)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       result.initStmts,
		replacementNode: expr,
	}
}

func walkArrowFunction(cx *FunctionContext, expr *ast.BLangArrowFunction) desugaredNode[model.ExpressionNode] {
	// Arrow functions have a body that may need desugaring
	if expr.Body != nil {
		result := walkExpression(cx, expr.Body.Expr)
		expr.Body.Expr = result.replacementNode.(ast.BLangExpression)
		// Handle initStmts if needed - arrow functions may need special handling
	}

	return desugaredNode[model.ExpressionNode]{
		replacementNode: expr,
	}
}

func walkMappingConstructorExpr(cx *FunctionContext, expr *ast.BLangMappingConstructorExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	for _, field := range expr.Fields {
		kv := field.(*ast.BLangMappingKeyValueField)

		if !kv.Key.ComputedKey {
			if varRef, ok := kv.Key.Expr.(*ast.BLangSimpleVarRef); ok {
				name := varRef.VariableName.Value
				lit := &ast.BLangLiteral{
					Value:         name,
					OriginalValue: name,
				}
				s := semtypes.STRING
				lit.SetDeterminedType(&s)
				kv.Key.Expr = lit
			}
		}

		result := walkExpression(cx, kv.ValueExpr)
		initStmts = append(initStmts, result.initStmts...)
		kv.ValueExpr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkQueryExpr(cx *FunctionContext, expr *ast.BLangQueryExpr) desugaredNode[model.ExpressionNode] {
	fromClause, ok := expr.QueryClauseList[0].(*ast.BLangFromClause)
	if !ok {
		cx.internalError("query expression must start with from clause")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}
	selectClause, ok := expr.QueryClauseList[len(expr.QueryClauseList)-1].(*ast.BLangSelectClause)
	if !ok {
		cx.internalError("query expression must end with select clause")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}
	loopVarDef, ok := fromClause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
	if !ok {
		cx.unimplemented("query from clause currently supports only simple variable definition")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}
	cloneLoopVarDef := cloneSimpleVariableDef(loopVarDef)
	if cloneLoopVarDef == nil || cloneLoopVarDef.Var == nil {
		cx.internalError("failed to clone query from-clause variable definition")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}

	queryTy := expr.GetDeterminedType()

	var initStmts []model.StatementNode

	collResult := walkExpression(cx, fromClause.Collection)
	initStmts = append(initStmts, collResult.initStmts...)
	collExpr := collResult.replacementNode.(ast.BLangExpression)
	collTy := collExpr.GetDeterminedType()

	collName, collSymbol := cx.addDesugardSymbol(collTy, model.SymbolKindVariable, false)
	collVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: collName},
	}
	collVar.SetDeterminedType(collTy)
	collVar.SetInitialExpression(collExpr)
	collVar.SetSymbol(collSymbol)
	collVarDef := &ast.BLangSimpleVariableDef{Var: collVar}
	initStmts = append(initStmts, collVarDef)

	collRef := &ast.BLangSimpleVarRef{
		VariableName: collVar.Name,
	}
	collRef.SetSymbol(collSymbol)
	collRef.SetDeterminedType(collTy)

	resultName, resultSymbol := cx.addDesugardSymbol(queryTy, model.SymbolKindVariable, false)
	resultVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: resultName},
	}
	resultVar.SetDeterminedType(queryTy)
	emptyList := &ast.BLangListConstructorExpr{
		Exprs: []ast.BLangExpression{},
	}
	emptyList.SetDeterminedType(&semtypes.LIST)
	emptyList.AtomicType = semtypes.LIST_ATOMIC_INNER
	resultVar.SetInitialExpression(emptyList)
	resultVar.SetSymbol(resultSymbol)
	resultVarDef := &ast.BLangSimpleVariableDef{Var: resultVar}
	initStmts = append(initStmts, resultVarDef)

	resultRef := &ast.BLangSimpleVarRef{
		VariableName: resultVar.Name,
	}
	resultRef.SetSymbol(resultSymbol)
	resultRef.SetDeterminedType(queryTy)

	zeroLiteral := &ast.BLangNumericLiteral{
		BLangLiteral: ast.BLangLiteral{
			Value:         int64(0),
			OriginalValue: "0",
		},
		Kind: model.NodeKind_NUMERIC_LITERAL,
	}
	zeroLiteral.SetDeterminedType(&semtypes.INT)

	idxName, idxSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	idxVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: idxName},
	}
	idxVar.SetDeterminedType(&semtypes.INT)
	idxVar.SetInitialExpression(zeroLiteral)
	idxVar.SetSymbol(idxSymbol)
	idxVarDef := &ast.BLangSimpleVariableDef{Var: idxVar}
	initStmts = append(initStmts, idxVarDef)

	idxRef := &ast.BLangSimpleVarRef{
		VariableName: idxVar.Name,
	}
	idxRef.SetSymbol(idxSymbol)
	idxRef.SetDeterminedType(&semtypes.INT)

	lengthInvocation := createLengthInvocation(cx, collRef)
	lenName, lenSymbol := cx.addDesugardSymbol(&semtypes.INT, model.SymbolKindVariable, false)
	lenVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: lenName},
	}
	lenVar.SetDeterminedType(&semtypes.INT)
	lenVar.SetInitialExpression(lengthInvocation)
	lenVar.SetSymbol(lenSymbol)
	lenVarDef := &ast.BLangSimpleVariableDef{Var: lenVar}
	initStmts = append(initStmts, lenVarDef)

	lenRef := &ast.BLangSimpleVarRef{
		VariableName: lenVar.Name,
	}
	lenRef.SetSymbol(lenSymbol)
	lenRef.SetDeterminedType(&semtypes.INT)

	condition := &ast.BLangBinaryExpr{
		LhsExpr: idxRef,
		RhsExpr: lenRef,
		OpKind:  model.OperatorKind_LESS_THAN,
	}
	condition.SetDeterminedType(&semtypes.BOOLEAN)

	elementAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: idxRef,
	}
	elementAccess.Expr = collRef
	loopVarSymbol := cloneLoopVarDef.Var.Symbol()
	loopVarTy := cx.symbolType(loopVarSymbol)
	if loopVarTy == nil {
		cx.internalError("query from-clause variable symbol type not found")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}
	elementAccess.SetDeterminedType(loopVarTy)
	cloneLoopVarDef.Var.SetInitialExpression(elementAccess)

	var bodyStmts []ast.BLangStatement
	bodyStmts = append(bodyStmts, cloneLoopVarDef)

	bodyStmts, ok = appendQueryIntermediateClauseStmts(cx, expr, idxRef, bodyStmts)
	if !ok {
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}

	selectResult := walkExpression(cx, selectClause.Expression)
	for _, s := range selectResult.initStmts {
		bodyStmts = append(bodyStmts, s.(ast.BLangStatement))
	}
	pushInvocation := createPushInvocation(cx, resultRef, selectResult.replacementNode.(ast.BLangExpression))
	if pushInvocation == nil {
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}
	bodyStmts = append(bodyStmts, &ast.BLangExpressionStmt{Expr: pushInvocation})
	bodyStmts = append(bodyStmts, createIncrementStmt(idxRef))

	whileStmt := &ast.BLangWhile{
		Expr: condition,
		Body: ast.BLangBlockStmt{
			Stmts: bodyStmts,
		},
	}
	whileStmt.SetScope(cx.currentScope())
	whileStmt.SetDeterminedType(&semtypes.NEVER)
	initStmts = append(initStmts, whileStmt)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: resultRef,
	}
}

func cloneSimpleVariableDef(varDef *ast.BLangSimpleVariableDef) *ast.BLangSimpleVariableDef {
	if varDef == nil {
		return nil
	}
	clone := *varDef
	if varDef.Var == nil {
		return &clone
	}
	cloneVar := *varDef.Var
	if varDef.Var.Name != nil {
		cloneName := *varDef.Var.Name
		cloneVar.Name = &cloneName
	}
	clone.Var = &cloneVar
	return &clone
}

func appendQueryIntermediateClauseStmts(
	cx *FunctionContext,
	queryExpr *ast.BLangQueryExpr,
	idxRef ast.BLangExpression,
	bodyStmts []ast.BLangStatement,
) ([]ast.BLangStatement, bool) {
	for i := 1; i < len(queryExpr.QueryClauseList)-1; i++ {
		switch clause := queryExpr.QueryClauseList[i].(type) {
		case *ast.BLangLetClause:
			for _, variableDef := range clause.LetVarDeclarations {
				varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
				if !ok || varDef.Var == nil || varDef.Var.Expr == nil {
					cx.unimplemented("query let clause currently supports only initialized simple variable declarations")
					return nil, false
				}
				letResult := walkExpression(cx, varDef.Var.Expr.(ast.BLangExpression))
				for _, s := range letResult.initStmts {
					bodyStmts = append(bodyStmts, s.(ast.BLangStatement))
				}
				varDef.Var.SetInitialExpression(letResult.replacementNode.(ast.BLangExpression))
				bodyStmts = append(bodyStmts, varDef)
			}
		case *ast.BLangWhereClause:
			if clause.Expression == nil {
				cx.unimplemented("query where clause requires a condition expression")
				return nil, false
			}
			whereResult := walkExpression(cx, clause.Expression)
			for _, s := range whereResult.initStmts {
				bodyStmts = append(bodyStmts, s.(ast.BLangStatement))
			}
			whereCond := whereResult.replacementNode.(ast.BLangExpression)
			notWhereCond := &ast.BLangUnaryExpr{
				Expr:     whereCond,
				Operator: model.OperatorKind_NOT,
			}
			notWhereCond.SetDeterminedType(&semtypes.BOOLEAN)
			continueStmt := &ast.BLangContinue{}
			continueStmt.SetDeterminedType(&semtypes.NEVER)
			skipBody := ast.BLangBlockStmt{
				Stmts: []ast.BLangStatement{
					createIncrementStmt(idxRef),
					continueStmt,
				},
			}
			filterIf := &ast.BLangIf{
				Expr: notWhereCond,
				Body: skipBody,
			}
			filterIf.SetScope(cx.currentScope())
			filterIf.SetDeterminedType(&semtypes.NEVER)
			bodyStmts = append(bodyStmts, filterIf)
		default:
			cx.unimplemented("query expression currently supports only let + where clauses as intermediate clauses")
			return nil, false
		}
	}
	return bodyStmts, true
}

func createPushInvocation(cx *FunctionContext, listExpr ast.BLangExpression, valueExpr ast.BLangExpression) *ast.BLangInvocation {
	pkgName := array.PackageName
	space, ok := cx.getImportedSymbolSpace(pkgName)
	if !ok {
		cx.internalError(pkgName + " symbol space not found")
		return nil
	}
	symbolRef, ok := space.GetSymbol("push")
	if !ok {
		cx.internalError(pkgName + ":push symbol not found")
		return nil
	}
	cx.addImplicitImport(pkgName, ast.BLangImportPackage{
		OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
		PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
		Alias:        &ast.BLangIdentifier{Value: pkgName},
	})
	inv := &ast.BLangInvocation{
		Name:     &ast.BLangIdentifier{Value: "push"},
		PkgAlias: &ast.BLangIdentifier{Value: pkgName},
		ArgExprs: []ast.BLangExpression{listExpr, valueExpr},
	}
	inv.SetSymbol(symbolRef)
	inv.SetDeterminedType(&semtypes.NIL)
	return inv
}
