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
	"fmt"

	"ballerina-lang-go/ast"
	array "ballerina-lang-go/lib/array/compile"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

func walkExpression(cx *functionContext, node model.ExpressionNode) desugaredNode[model.ExpressionNode] {
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
	case *ast.BLangNewExpression:
		return walkNewExpression(cx, expr)
	case *ast.BLangNamedArgsExpression:
		result := walkExpression(cx, expr.Expr)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
		return desugaredNode[model.ExpressionNode]{
			initStmts:       result.initStmts,
			replacementNode: expr,
		}
	case *ast.BLangWildCardBindingPattern:
		// Wildcard binding pattern can appear in variable references (e.g., _ = expr)
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", node))
	}
}

func walkBinaryExpr(cx *functionContext, expr *ast.BLangBinaryExpr) desugaredNode[model.ExpressionNode] {
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

	if !isNilLiftableBinaryOp(expr.OpKind) {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	lhsTy := expr.LhsExpr.GetDeterminedType()
	rhsTy := expr.RhsExpr.GetDeterminedType()
	if lhsTy == nil || rhsTy == nil {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}
	lhsHasNil := semtypes.ContainsBasicType(lhsTy, semtypes.NIL)
	rhsHasNil := semtypes.ContainsBasicType(rhsTy, semtypes.NIL)

	if !lhsHasNil && !rhsHasNil {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	basePos := expr.GetPosition()
	resultTy := expr.GetDeterminedType()

	// Create temp vars for nullable operands
	var lhsVarName *ast.BLangIdentifier
	var lhsSymbol model.SymbolRef
	if lhsHasNil {
		lhsVarName, lhsSymbol, initStmts = createOperandTempVar(cx, lhsTy, expr.LhsExpr, basePos, initStmts)
	}

	var rhsVarName *ast.BLangIdentifier
	var rhsSymbol model.SymbolRef
	if rhsHasNil {
		rhsVarName, rhsSymbol, initStmts = createOperandTempVar(cx, rhsTy, expr.RhsExpr, basePos, initStmts)
	}

	// Create result temp var initialized to nil
	resultVarName, resultSymbol, initStmts := createNilResultVar(cx, resultTy, basePos, initStmts)

	// Build the nil check condition
	var nilCheckCond ast.BLangExpression
	if lhsHasNil {
		nilCheckCond = createNilTypeTest(lhsVarName, lhsSymbol, lhsTy, basePos)
	}
	if rhsHasNil {
		rhsNilCheck := createNilTypeTest(rhsVarName, rhsSymbol, rhsTy, basePos)
		if nilCheckCond == nil {
			nilCheckCond = rhsNilCheck
		} else {
			orExpr := &ast.BLangBinaryExpr{
				LhsExpr: nilCheckCond,
				RhsExpr: rhsNilCheck,
				OpKind:  model.OperatorKind_OR,
			}
			orExpr.SetDeterminedType(semtypes.BOOLEAN)
			orExpr.SetPosition(basePos)
			nilCheckCond = orExpr
		}
	}

	// Build the operation in the else branch
	var lhsRef ast.BLangExpression
	if lhsHasNil {
		lhsRef = createVarRef(lhsVarName, lhsSymbol, semtypes.Diff(lhsTy, semtypes.NIL))
	} else {
		lhsRef = expr.LhsExpr
	}

	var rhsRef ast.BLangExpression
	if rhsHasNil {
		rhsRef = createVarRef(rhsVarName, rhsSymbol, semtypes.Diff(rhsTy, semtypes.NIL))
	} else {
		rhsRef = expr.RhsExpr
	}

	newBinaryExpr := &ast.BLangBinaryExpr{
		LhsExpr: lhsRef,
		RhsExpr: rhsRef,
		OpKind:  expr.OpKind,
	}
	newBinaryExpr.SetDeterminedType(semtypes.Diff(resultTy, semtypes.NIL))
	newBinaryExpr.SetPosition(basePos)

	resultAssign := createResultAssignment(resultVarName, resultSymbol, resultTy, newBinaryExpr, basePos)

	elseBody := &ast.BLangBlockStmt{
		Stmts: []ast.BLangStatement{resultAssign},
	}
	elseBody.SetDeterminedType(semtypes.NEVER)
	ifStmt := &ast.BLangIf{
		Expr:     nilCheckCond,
		Body:     ast.BLangBlockStmt{},
		ElseStmt: elseBody,
	}
	ifStmt.Body.SetDeterminedType(semtypes.NEVER)
	ifStmt.SetDeterminedType(semtypes.NEVER)
	ifStmt.SetScope(cx.currentScope())
	setPositionIfMissing(ifStmt, basePos)
	initStmts = append(initStmts, ifStmt)

	replacementRef := createVarRef(resultVarName, resultSymbol, resultTy)
	setPositionIfMissing(replacementRef, basePos)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: replacementRef,
	}
}

func walkUnaryExpr(cx *functionContext, expr *ast.BLangUnaryExpr) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	// Unary + is identity — desugar to just the operand (BIR gen doesn't handle unary +)
	if expr.Operator == model.OperatorKind_ADD {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr.Expr,
		}
	}

	if !isNilLiftableUnaryOp(expr.Operator) {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	operandTy := expr.Expr.GetDeterminedType()
	if !semtypes.ContainsBasicType(operandTy, semtypes.NIL) {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	basePos := expr.GetPosition()
	resultTy := expr.GetDeterminedType()

	// Create operand temp var
	operandVarName, operandSymbol, initStmts := createOperandTempVar(cx, operandTy, expr.Expr, basePos, initStmts)

	// Create result temp var initialized to nil
	resultVarName, resultSymbol, initStmts := createNilResultVar(cx, resultTy, basePos, initStmts)

	// Build nil check: if ($operand is ()) { } else { ... }
	nilCheck := createNilTypeTest(operandVarName, operandSymbol, operandTy, basePos)

	// Build the operation for the if-body (operand is not nil)
	nonNilTy := semtypes.Diff(operandTy, semtypes.NIL)
	operandRef := createVarRef(operandVarName, operandSymbol, nonNilTy)

	newUnary := &ast.BLangUnaryExpr{
		Expr:     operandRef,
		Operator: expr.Operator,
	}
	newUnary.SetDeterminedType(semtypes.Diff(resultTy, semtypes.NIL))
	newUnary.SetPosition(basePos)
	var opExpr ast.BLangExpression = newUnary

	resultAssign := createResultAssignment(resultVarName, resultSymbol, resultTy, opExpr, basePos)

	// if ($operand is ()) { } else { $result = op $operand }
	elseBody := &ast.BLangBlockStmt{
		Stmts: []ast.BLangStatement{resultAssign},
	}
	elseBody.SetDeterminedType(semtypes.NEVER)
	ifStmt := &ast.BLangIf{
		Expr:     nilCheck,
		Body:     ast.BLangBlockStmt{},
		ElseStmt: elseBody,
	}
	ifStmt.Body.SetDeterminedType(semtypes.NEVER)
	ifStmt.SetDeterminedType(semtypes.NEVER)
	ifStmt.SetScope(cx.currentScope())
	setPositionIfMissing(ifStmt, basePos)
	initStmts = append(initStmts, ifStmt)

	replacementRef := createVarRef(resultVarName, resultSymbol, resultTy)
	setPositionIfMissing(replacementRef, basePos)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: replacementRef,
	}
}

func walkElvisExpr(cx *functionContext, expr *ast.BLangElvisExpr) desugaredNode[model.ExpressionNode] {
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

func walkGroupExpr(cx *functionContext, expr *ast.BLangGroupExpr) desugaredNode[model.ExpressionNode] {
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

func walkIndexBasedAccess(cx *functionContext, expr *ast.BLangIndexBasedAccess) desugaredNode[model.ExpressionNode] {
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

func walkFieldBaseAccess(cx *functionContext, expr *ast.BLangFieldBaseAccess) desugaredNode[model.ExpressionNode] {
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
	lit.SetPosition(expr.GetPosition())
	lit.SetDeterminedType(semtypes.STRING)

	indexAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: lit,
	}
	indexAccess.SetPosition(expr.GetPosition())
	indexAccess.Expr = expr.Expr
	indexAccess.SetDeterminedType(expr.GetDeterminedType())

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: indexAccess,
	}
}

func walkInvocation(cx *functionContext, expr *ast.BLangInvocation) desugaredNode[model.ExpressionNode] {
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

	fnSym, isDirectCall := cx.getSymbol(expr.Symbol()).(model.FunctionSymbol)
	if !isDirectCall {
		return desugaredNode[model.ExpressionNode]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	directStmts := walkDirectCallArgs(cx, expr, fnSym)
	initStmts = append(initStmts, directStmts...)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkDirectCallArgs(cx *functionContext, expr *ast.BLangInvocation, fnSym model.FunctionSymbol) []model.StatementNode {
	sig := fnSym.Signature()
	totalParams := len(sig.ParamTypes)
	if totalParams == 0 {
		return nil
	}

	reordered := make([]ast.BLangExpression, totalParams)
	for i, arg := range expr.ArgExprs {
		switch arg := arg.(type) {
		case *ast.BLangNamedArgsExpression:
			for j, name := range sig.ParamNames {
				if name == arg.Name.Value {
					reordered[j] = arg.Expr
					break
				}
			}
		default:
			if i < totalParams {
				reordered[i] = arg
			} else {
				reordered = append(reordered, arg)
			}
		}
	}

	pos := expr.GetPosition()
	defaultableParams := fnSym.DefaultableParams()
	var initStmts []model.StatementNode

	var transformed []ast.BLangExpression
	for i := range totalParams {
		if reordered[i] != nil {
			continue
		}
		for j := len(transformed); j < i; j++ {
			varDef, varRef := assignToLocal(cx, reordered[j], pos)
			initStmts = append(initStmts, varDef)
			reordered[j] = varRef
			transformed = append(transformed, varRef)
		}
		dp, _ := defaultableParams.Get(i)
		defaultInv := &ast.BLangInvocation{
			Name:     &ast.BLangIdentifier{Value: cx.pkgCtx.compilerCtx.GetSymbol(dp.Symbol).Name()},
			ArgExprs: transformed,
		}
		defaultInv.SetSymbol(dp.Symbol)
		defaultInv.SetDeterminedType(sig.ParamTypes[i])
		setPositionIfMissing(defaultInv, pos)

		varDef, varRef := assignToLocal(cx, defaultInv, pos)
		initStmts = append(initStmts, varDef)
		reordered[i] = varRef
		transformed = append(transformed, varRef)
	}

	expr.ArgExprs = reordered
	return initStmts
}

func assignToLocal(cx *functionContext, initExpr ast.BLangExpression, pos diagnostics.Location) (model.StatementNode, *ast.BLangSimpleVarRef) {
	ty := initExpr.GetDeterminedType()
	tempName, tempSymRef := cx.addDesugardSymbol(ty, model.SymbolKindVariable, false)
	tempVar := &ast.BLangSimpleVariable{Name: &ast.BLangIdentifier{Value: tempName}}
	tempVar.SetDeterminedType(ty)
	tempVar.SetInitialExpression(initExpr)
	tempVar.SetSymbol(tempSymRef)
	varDef := &ast.BLangSimpleVariableDef{}
	varDef.SetVariable(tempVar)
	varDef.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(varDef, pos)

	varRef := &ast.BLangSimpleVarRef{VariableName: tempVar.Name}
	varRef.SetSymbol(tempSymRef)
	varRef.SetDeterminedType(ty)
	setPositionIfMissing(varRef, pos)
	return varDef, varRef
}

func walkListConstructorExpr(cx *functionContext, expr *ast.BLangListConstructorExpr) desugaredNode[model.ExpressionNode] {
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

func walkErrorConstructorExpr(cx *functionContext, expr *ast.BLangErrorConstructorExpr) desugaredNode[model.ExpressionNode] {
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

func walkCheckedExpr(cx *functionContext, expr *ast.BLangCheckedExpr) desugaredNode[model.ExpressionNode] {
	return desugarCheckedExpr(cx, expr, false)
}

func walkCheckPanickedExpr(cx *functionContext, expr *ast.BLangCheckPanickedExpr) desugaredNode[model.ExpressionNode] {
	return desugarCheckedExpr(cx, &expr.BLangCheckedExpr, true)
}

func walkTrapExpr(cx *functionContext, expr *ast.BLangTrapExpr) desugaredNode[model.ExpressionNode] {
	result := walkExpression(cx, expr.Expr)
	if len(result.initStmts) > 0 {
		// I don't think this can ever happen but if it does we need to think about how to add these statements in to the
		// trap region in BIR gen
		cx.internalError("Init statements will be hoisted outside of trap region")
	}
	expr.Expr = result.replacementNode.(ast.BLangExpression)
	return desugaredNode[model.ExpressionNode]{initStmts: nil, replacementNode: expr}
}

func desugarCheckedExpr(cx *functionContext, expr *ast.BLangCheckedExpr, isPanic bool) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	// Walk the inner expression first
	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	innerTy := expr.Expr.GetDeterminedType()
	resultTy := expr.GetDeterminedType()

	basePos := expr.Expr.GetPosition()

	// TODO: extract util to add definition and get reference
	// Create temp var: $desugar$N = <inner expr>
	tempName, tempSymbol := cx.addDesugardSymbol(innerTy, model.SymbolKindVariable, false)
	tempVarName := &ast.BLangIdentifier{Value: tempName}
	tempVar := &ast.BLangSimpleVariable{Name: tempVarName}
	tempVar.SetDeterminedType(innerTy)
	tempVar.SetInitialExpression(expr.Expr)
	tempVar.SetSymbol(tempSymbol)
	tempVarDef := &ast.BLangSimpleVariableDef{Var: tempVar}
	setPositionIfMissing(tempVarDef, basePos)
	initStmts = append(initStmts, tempVarDef)

	// Type test: $desugar$N is error
	tempVarRefForTest := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	tempVarRefForTest.SetSymbol(tempSymbol)
	tempVarRefForTest.SetDeterminedType(innerTy)

	typeTestExpr := &ast.BLangTypeTestExpr{}
	typeTestExpr.Expr = tempVarRefForTest
	typeTestExpr.Type = model.TypeData{Type: semtypes.ERROR}
	typeTestExpr.SetDeterminedType(semtypes.BOOLEAN)

	// If body: return or panic
	tempVarRefForBody := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	tempVarRefForBody.SetSymbol(tempSymbol)
	tempVarRefForBody.SetDeterminedType(innerTy)

	var bodyStmt ast.BLangStatement
	if isPanic {
		panicStmt := &ast.BLangPanic{Expr: tempVarRefForBody}
		panicStmt.SetPosition(expr.GetPosition())
		bodyStmt = panicStmt
	} else {
		bodyStmt = &ast.BLangReturn{Expr: tempVarRefForBody}
	}
	setPositionIfMissing(bodyStmt.(ast.BLangNode), basePos)

	ifBody := ast.BLangBlockStmt{
		Stmts: []ast.BLangStatement{bodyStmt},
	}
	ifStmt := &ast.BLangIf{
		Expr: typeTestExpr,
		Body: ifBody,
	}
	initStmts = append(initStmts, ifStmt)
	setPositionIfMissing(ifStmt, basePos)

	// Replacement: var ref typed as non-error type
	replacementVarRef := &ast.BLangSimpleVarRef{VariableName: tempVarName}
	replacementVarRef.SetSymbol(tempSymbol)
	replacementVarRef.SetDeterminedType(resultTy)
	setPositionIfMissing(replacementVarRef, basePos)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: replacementVarRef,
	}
}

func walkDynamicArgExpr(cx *functionContext, expr *ast.BLangDynamicArgExpr) desugaredNode[model.ExpressionNode] {
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

func walkLambdaFunction(cx *functionContext, expr *ast.BLangLambdaFunction) desugaredNode[model.ExpressionNode] {
	// Desugar the function body
	if expr.Function != nil {
		expr.Function = desugarFunction(cx.pkgCtx, expr.Function)
	}

	return desugaredNode[model.ExpressionNode]{
		replacementNode: expr,
	}
}

func walkTypeConversionExpr(cx *functionContext, expr *ast.BLangTypeConversionExpr) desugaredNode[model.ExpressionNode] {
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

func walkTypeTestExpr(cx *functionContext, expr *ast.BLangTypeTestExpr) desugaredNode[model.ExpressionNode] {
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

func walkAnnotAccessExpr(cx *functionContext, expr *ast.BLangAnnotAccessExpr) desugaredNode[model.ExpressionNode] {
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

func walkCollectContextInvocation(cx *functionContext, expr *ast.BLangCollectContextInvocation) desugaredNode[model.ExpressionNode] {
	// Walk the underlying invocation
	result := walkInvocation(cx, &expr.Invocation)
	expr.Invocation = *result.replacementNode.(*ast.BLangInvocation)

	return desugaredNode[model.ExpressionNode]{
		initStmts:       result.initStmts,
		replacementNode: expr,
	}
}

func walkArrowFunction(cx *functionContext, expr *ast.BLangArrowFunction) desugaredNode[model.ExpressionNode] {
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

func walkNewExpression(cx *functionContext, expr *ast.BLangNewExpression) desugaredNode[model.ExpressionNode] {
	var initStmts []model.StatementNode

	for i := range expr.ArgsExprs {
		result := walkExpression(cx, expr.ArgsExprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.ArgsExprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkMappingConstructorExpr(cx *functionContext, expr *ast.BLangMappingConstructorExpr) desugaredNode[model.ExpressionNode] {
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
				lit.SetPosition(varRef.GetPosition())
				lit.SetDeterminedType(semtypes.STRING)
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

func walkQueryExpr(cx *functionContext, expr *ast.BLangQueryExpr) desugaredNode[model.ExpressionNode] {
	fromClause := expr.QueryClauseList[0].(*ast.BLangFromClause)

	selectClauseIndex := len(expr.QueryClauseList) - 1
	var onConflictClause *ast.BLangOnConflictClause
	if clause, isOnConflict := expr.QueryClauseList[selectClauseIndex].(*ast.BLangOnConflictClause); isOnConflict {
		onConflictClause = clause
		selectClauseIndex--
	}

	selectClause := expr.QueryClauseList[selectClauseIndex].(*ast.BLangSelectClause)
	orderByClauseIndices := queryOrderByClauseIndices(expr, 1, selectClauseIndex)

	loopVarDef := fromClause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)

	queryTy := expr.GetDeterminedType()
	basePos := expr.GetPosition()

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
	setPositionIfMissing(collVarDef, basePos)
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
	switch expr.QueryConstructType {
	case model.TypeKind_MAP:
		emptyMap := &ast.BLangMappingConstructorExpr{
			Fields: []model.MappingField{},
		}
		emptyMap.SetDeterminedType(queryTy)
		resultVar.SetInitialExpression(emptyMap)
	default:
		emptyList := &ast.BLangListConstructorExpr{
			Exprs: []ast.BLangExpression{},
		}
		emptyList.SetDeterminedType(semtypes.LIST)
		emptyList.AtomicType = semtypes.LIST_ATOMIC_INNER
		resultVar.SetInitialExpression(emptyList)
	}
	resultVar.SetSymbol(resultSymbol)
	resultVarDef := &ast.BLangSimpleVariableDef{Var: resultVar}
	setPositionIfMissing(resultVarDef, basePos)
	initStmts = append(initStmts, resultVarDef)

	resultRef := &ast.BLangSimpleVarRef{
		VariableName: resultVar.Name,
	}
	resultRef.SetSymbol(resultSymbol)
	resultRef.SetDeterminedType(queryTy)

	var seenKeysRef *ast.BLangSimpleVarRef
	if onConflictClause != nil && expr.QueryConstructType == model.TypeKind_MAP {
		seenKeysRef = createQueryMapStore(cx, &initStmts, basePos)
	}

	loopVarSymbol := loopVarDef.Var.Symbol()
	loopVarTy := cx.symbolType(loopVarSymbol)

	lengthSource := collRef
	var keysRef *ast.BLangSimpleVarRef

	switch {
	case semtypes.IsSubtypeSimple(collTy, semtypes.LIST):
	case semtypes.IsSubtypeSimple(collTy, semtypes.MAPPING):
		keysInvocation := createKeysInvocation(cx, collRef)
		if keysInvocation == nil {
			return desugaredNode[model.ExpressionNode]{replacementNode: expr}
		}
		keysTy := keysInvocation.GetDeterminedType()
		keysName, keysSymbol := cx.addDesugardSymbol(keysTy, model.SymbolKindVariable, false)
		keysVar := &ast.BLangSimpleVariable{
			Name: &ast.BLangIdentifier{Value: keysName},
		}
		keysVar.SetDeterminedType(keysTy)
		keysVar.SetInitialExpression(keysInvocation)
		keysVar.SetSymbol(keysSymbol)
		keysVarDef := &ast.BLangSimpleVariableDef{Var: keysVar}
		setPositionIfMissing(keysVarDef, basePos)
		initStmts = append(initStmts, keysVarDef)

		keysRef = &ast.BLangSimpleVarRef{
			VariableName: keysVar.Name,
		}
		keysRef.SetSymbol(keysSymbol)
		keysRef.SetDeterminedType(keysTy)
		lengthSource = keysRef
	default:
		cx.unimplemented("query from clause currently supports only list or map collections")
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}

	lengthInvocation := createLengthInvocation(cx, lengthSource)
	lenName, lenSymbol := cx.addDesugardSymbol(semtypes.INT, model.SymbolKindVariable, false)
	lenVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: lenName},
	}
	lenVar.SetDeterminedType(semtypes.INT)
	lenVar.SetInitialExpression(lengthInvocation)
	lenVar.SetSymbol(lenSymbol)
	lenVarDef := &ast.BLangSimpleVariableDef{Var: lenVar}
	setPositionIfMissing(lenVarDef, basePos)
	initStmts = append(initStmts, lenVarDef)

	lenRef := &ast.BLangSimpleVarRef{
		VariableName: lenVar.Name,
	}
	lenRef.SetSymbol(lenSymbol)
	lenRef.SetDeterminedType(semtypes.INT)
	stageInput := queryOrderStageInput{
		rowCountRef: lenRef,
	}
	stageStart := 1
	var ok bool
	for _, orderByClauseIndex := range orderByClauseIndices {
		stageInput, ok = appendQueryOrderByStageStmts(
			cx,
			expr,
			collRef,
			keysRef,
			loopVarDef,
			loopVarTy,
			stageStart,
			orderByClauseIndex,
			stageInput,
			&initStmts,
			basePos,
		)
		if !ok {
			return desugaredNode[model.ExpressionNode]{replacementNode: expr}
		}
		stageStart = orderByClauseIndex + 1
	}
	ok = appendQueryFinalStageStmts(
		cx,
		expr,
		collRef,
		keysRef,
		loopVarDef,
		loopVarTy,
		stageStart,
		selectClauseIndex,
		stageInput,
		resultRef,
		selectClause,
		onConflictClause,
		seenKeysRef,
		&initStmts,
		basePos,
	)
	if !ok {
		return desugaredNode[model.ExpressionNode]{replacementNode: expr}
	}

	setPositionIfMissing(resultRef, basePos)
	return desugaredNode[model.ExpressionNode]{
		initStmts:       initStmts,
		replacementNode: resultRef,
	}
}

func cloneSimpleVariableDef(varDef *ast.BLangSimpleVariableDef) *ast.BLangSimpleVariableDef {
	clone := *varDef
	cloneVar := *varDef.Var
	if varDef.Var.Name != nil {
		cloneName := *varDef.Var.Name
		cloneVar.Name = &cloneName
	}
	clone.Var = &cloneVar
	return &clone
}

func queryOrderByClauseIndices(queryExpr *ast.BLangQueryExpr, startClauseIndex int, endClauseIndex int) []int {
	indices := make([]int, 0)
	for i := startClauseIndex; i < endClauseIndex; i++ {
		if _, isOrderBy := queryExpr.QueryClauseList[i].(*ast.BLangOrderByClause); isOrderBy {
			indices = append(indices, i)
		}
	}
	return indices
}

type queryLetStore struct {
	VarDef   *ast.BLangSimpleVariableDef
	StoreRef *ast.BLangSimpleVarRef
	ValueTy  semtypes.SemType
}

type queryOrderStageInput struct {
	indexRowsRef  *ast.BLangSimpleVarRef
	rowCountRef   *ast.BLangSimpleVarRef
	payloadStores []queryLetStore
}

func createQueryCounterRef(
	cx *functionContext,
	initStmts *[]model.StatementNode,
	pos diagnostics.Location,
) *ast.BLangSimpleVarRef {
	counterName, counterSymbol := cx.addDesugardSymbol(semtypes.INT, model.SymbolKindVariable, false)
	counterVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: counterName},
	}
	counterVar.SetDeterminedType(semtypes.INT)
	counterVar.SetInitialExpression(createIntLiteral(0))
	counterVar.SetSymbol(counterSymbol)
	counterVarDef := &ast.BLangSimpleVariableDef{Var: counterVar}
	setPositionIfMissing(counterVarDef, pos)
	*initStmts = append(*initStmts, counterVarDef)

	counterRef := &ast.BLangSimpleVarRef{VariableName: counterVar.Name}
	counterRef.SetSymbol(counterSymbol)
	counterRef.SetDeterminedType(semtypes.INT)
	return counterRef
}

func createQueryLengthRef(
	cx *functionContext,
	initStmts *[]model.StatementNode,
	source ast.BLangExpression,
	pos diagnostics.Location,
) (*ast.BLangSimpleVarRef, bool) {
	lengthInvocation := createLengthInvocation(cx, source)
	if lengthInvocation == nil {
		return nil, false
	}
	lengthName, lengthSymbol := cx.addDesugardSymbol(semtypes.INT, model.SymbolKindVariable, false)
	lengthVar := &ast.BLangSimpleVariable{Name: &ast.BLangIdentifier{Value: lengthName}}
	lengthVar.SetDeterminedType(semtypes.INT)
	lengthVar.SetInitialExpression(lengthInvocation)
	lengthVar.SetSymbol(lengthSymbol)
	lengthVarDef := &ast.BLangSimpleVariableDef{Var: lengthVar}
	setPositionIfMissing(lengthVarDef, pos)
	*initStmts = append(*initStmts, lengthVarDef)
	lengthRef := &ast.BLangSimpleVarRef{VariableName: lengthVar.Name}
	lengthRef.SetSymbol(lengthSymbol)
	lengthRef.SetDeterminedType(semtypes.INT)
	return lengthRef, true
}

func queryStageBaseIndexExpr(loopCounterRef *ast.BLangSimpleVarRef, indexRowsRef *ast.BLangSimpleVarRef) ast.BLangExpression {
	if indexRowsRef == nil {
		return loopCounterRef
	}
	rowIndexAccess := &ast.BLangIndexBasedAccess{IndexExpr: loopCounterRef}
	rowIndexAccess.Expr = indexRowsRef
	rowIndexAccess.SetDeterminedType(semtypes.INT)
	return rowIndexAccess
}

func appendQueryOrderByStageStmts(
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	collRef ast.BLangExpression,
	keysRef *ast.BLangSimpleVarRef,
	loopVarDef *ast.BLangSimpleVariableDef,
	loopVarTy semtypes.SemType,
	startClauseIndex int,
	orderByClauseIndex int,
	stageInput queryOrderStageInput,
	initStmts *[]model.StatementNode,
	basePos diagnostics.Location,
) (queryOrderStageInput, bool) {
	orderByClause := queryExpr.QueryClauseList[orderByClauseIndex].(*ast.BLangOrderByClause)

	orderKeyRowsRef := createQueryListStore(cx, initStmts, basePos)
	sortedIndexRowsRef := createQueryListStore(cx, initStmts, basePos)
	newLetStores, ok := createPreOrderLetStores(cx, queryExpr, startClauseIndex, orderByClauseIndex, initStmts, basePos)
	if !ok {
		return queryOrderStageInput{}, false
	}
	stageStores := make([]queryLetStore, 0, len(stageInput.payloadStores)+len(newLetStores))
	for _, store := range stageInput.payloadStores {
		stageStores = append(stageStores, queryLetStore{
			VarDef:   store.VarDef,
			StoreRef: createQueryListStore(cx, initStmts, basePos),
			ValueTy:  store.ValueTy,
		})
	}
	stageStores = append(stageStores, newLetStores...)
	payloadRowsRef, ok := createQueryPayloadStore(cx, initStmts, basePos, stageStores)
	if !ok {
		return queryOrderStageInput{}, false
	}

	loopCounterRef := createQueryCounterRef(cx, initStmts, basePos)
	baseIndexExpr := queryStageBaseIndexExpr(loopCounterRef, stageInput.indexRowsRef)
	elementAccess := queryElementAccess(collRef, keysRef, baseIndexExpr, loopVarTy)
	loopVarDefClone := cloneSimpleVariableDef(loopVarDef)
	loopVarDefClone.Var.SetInitialExpression(elementAccess)

	var bodyStmts []ast.BLangStatement
	bodyStmts = append(bodyStmts, loopVarDefClone)
	for _, store := range stageInput.payloadStores {
		restoredVarDef := cloneSimpleVariableDef(store.VarDef)
		storeAccess := &ast.BLangIndexBasedAccess{IndexExpr: loopCounterRef}
		storeAccess.Expr = store.StoreRef
		storeAccess.SetDeterminedType(store.ValueTy)
		restoredVarDef.Var.SetInitialExpression(storeAccess)
		bodyStmts = append(bodyStmts, restoredVarDef)
	}
	bodyStmts, ok = appendQueryIntermediateClauseStmts(
		cx,
		queryExpr,
		loopCounterRef,
		initStmts,
		bodyStmts,
		startClauseIndex,
		orderByClauseIndex,
	)
	if !ok {
		return queryOrderStageInput{}, false
	}

	keyTuple, keyInitStmts := buildOrderKeyTupleExpr(cx, orderByClause, basePos)
	bodyStmts = append(bodyStmts, keyInitStmts...)
	if pushKeys := createPushInvocation(cx, orderKeyRowsRef, keyTuple); pushKeys != nil {
		pushStmt := &ast.BLangExpressionStmt{Expr: pushKeys}
		setPositionIfMissing(pushStmt, basePos)
		bodyStmts = append(bodyStmts, pushStmt)
	} else {
		return queryOrderStageInput{}, false
	}
	if pushIndex := createPushInvocation(cx, sortedIndexRowsRef, baseIndexExpr); pushIndex != nil {
		pushStmt := &ast.BLangExpressionStmt{Expr: pushIndex}
		setPositionIfMissing(pushStmt, basePos)
		bodyStmts = append(bodyStmts, pushStmt)
	} else {
		return queryOrderStageInput{}, false
	}
	for _, store := range stageStores {
		valueRef := &ast.BLangSimpleVarRef{
			VariableName: store.VarDef.Var.Name,
		}
		valueRef.SetSymbol(store.VarDef.Var.Symbol())
		valueRef.SetDeterminedType(store.ValueTy)
		pushStore := createPushInvocation(cx, store.StoreRef, valueRef)
		if pushStore == nil {
			return queryOrderStageInput{}, false
		}
		pushStmt := &ast.BLangExpressionStmt{Expr: pushStore}
		setPositionIfMissing(pushStmt, basePos)
		bodyStmts = append(bodyStmts, pushStmt)
	}
	bodyStmts = append(bodyStmts, createIncrementStmt(loopCounterRef))

	stageCondition := &ast.BLangBinaryExpr{
		LhsExpr: loopCounterRef,
		RhsExpr: stageInput.rowCountRef,
		OpKind:  model.OperatorKind_LESS_THAN,
	}
	stageCondition.SetDeterminedType(semtypes.BOOLEAN)
	stageWhile := &ast.BLangWhile{
		Expr: stageCondition,
		Body: ast.BLangBlockStmt{Stmts: bodyStmts},
	}
	stageWhile.SetScope(cx.currentScope())
	stageWhile.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(stageWhile, basePos)
	*initStmts = append(*initStmts, stageWhile)

	directionsExpr := buildOrderDirectionExpr(orderByClause, basePos)
	sortInvocation := createQuerySortInvocation(cx, orderKeyRowsRef, directionsExpr, sortedIndexRowsRef, payloadRowsRef)
	if sortInvocation == nil {
		return queryOrderStageInput{}, false
	}
	sortStmt := &ast.BLangExpressionStmt{Expr: sortInvocation}
	setPositionIfMissing(sortStmt, basePos)
	*initStmts = append(*initStmts, sortStmt)

	sortedLenRef, ok := createQueryLengthRef(cx, initStmts, sortedIndexRowsRef, basePos)
	if !ok {
		return queryOrderStageInput{}, false
	}
	return queryOrderStageInput{
		indexRowsRef:  sortedIndexRowsRef,
		rowCountRef:   sortedLenRef,
		payloadStores: stageStores,
	}, true
}

func appendQueryFinalStageStmts(
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	collRef ast.BLangExpression,
	keysRef *ast.BLangSimpleVarRef,
	loopVarDef *ast.BLangSimpleVariableDef,
	loopVarTy semtypes.SemType,
	startClauseIndex int,
	selectClauseIndex int,
	stageInput queryOrderStageInput,
	resultRef *ast.BLangSimpleVarRef,
	selectClause *ast.BLangSelectClause,
	onConflictClause *ast.BLangOnConflictClause,
	seenKeysRef *ast.BLangSimpleVarRef,
	initStmts *[]model.StatementNode,
	basePos diagnostics.Location,
) bool {
	loopCounterRef := createQueryCounterRef(cx, initStmts, basePos)
	baseIndexExpr := queryStageBaseIndexExpr(loopCounterRef, stageInput.indexRowsRef)
	elementAccess := queryElementAccess(collRef, keysRef, baseIndexExpr, loopVarTy)
	loopVarDefClone := cloneSimpleVariableDef(loopVarDef)
	loopVarDefClone.Var.SetInitialExpression(elementAccess)

	var bodyStmts []ast.BLangStatement
	bodyStmts = append(bodyStmts, loopVarDefClone)
	for _, store := range stageInput.payloadStores {
		restoredVarDef := cloneSimpleVariableDef(store.VarDef)
		storeAccess := &ast.BLangIndexBasedAccess{IndexExpr: loopCounterRef}
		storeAccess.Expr = store.StoreRef
		storeAccess.SetDeterminedType(store.ValueTy)
		restoredVarDef.Var.SetInitialExpression(storeAccess)
		bodyStmts = append(bodyStmts, restoredVarDef)
	}
	var ok bool
	bodyStmts, ok = appendQueryIntermediateClauseStmts(
		cx,
		queryExpr,
		loopCounterRef,
		initStmts,
		bodyStmts,
		startClauseIndex,
		selectClauseIndex,
	)
	if !ok {
		return false
	}
	bodyStmts, ok = appendQuerySelectResultStmts(
		cx,
		queryExpr,
		resultRef,
		selectClause,
		onConflictClause,
		seenKeysRef,
		basePos,
		bodyStmts,
	)
	if !ok {
		return false
	}
	bodyStmts = append(bodyStmts, createIncrementStmt(loopCounterRef))

	finalCondition := &ast.BLangBinaryExpr{
		LhsExpr: loopCounterRef,
		RhsExpr: stageInput.rowCountRef,
		OpKind:  model.OperatorKind_LESS_THAN,
	}
	finalCondition.SetDeterminedType(semtypes.BOOLEAN)
	finalWhile := &ast.BLangWhile{
		Expr: finalCondition,
		Body: ast.BLangBlockStmt{Stmts: bodyStmts},
	}
	finalWhile.SetScope(cx.currentScope())
	finalWhile.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(finalWhile, basePos)
	*initStmts = append(*initStmts, finalWhile)
	return true
}

func createQueryListStore(
	cx *functionContext,
	initStmts *[]model.StatementNode,
	pos diagnostics.Location,
) *ast.BLangSimpleVarRef {
	listName, listSymbol := cx.addDesugardSymbol(semtypes.LIST, model.SymbolKindVariable, false)
	emptyList := &ast.BLangListConstructorExpr{Exprs: []ast.BLangExpression{}}
	emptyList.SetDeterminedType(semtypes.LIST)
	emptyList.AtomicType = semtypes.LIST_ATOMIC_INNER
	setPositionIfMissing(emptyList, pos)
	listVar := &ast.BLangSimpleVariable{Name: &ast.BLangIdentifier{Value: listName}}
	listVar.SetDeterminedType(semtypes.LIST)
	listVar.SetInitialExpression(emptyList)
	listVar.SetSymbol(listSymbol)
	setPositionIfMissing(listVar, pos)
	listVarDef := &ast.BLangSimpleVariableDef{Var: listVar}
	setPositionIfMissing(listVarDef, pos)
	*initStmts = append(*initStmts, listVarDef)
	listRef := &ast.BLangSimpleVarRef{VariableName: listVar.Name}
	listRef.SetSymbol(listSymbol)
	listRef.SetDeterminedType(semtypes.LIST)
	setPositionIfMissing(listRef, pos)
	return listRef
}

func createQueryMapStore(
	cx *functionContext,
	initStmts *[]model.StatementNode,
	pos diagnostics.Location,
) *ast.BLangSimpleVarRef {
	mapName, mapSymbol := cx.addDesugardSymbol(semtypes.MAPPING, model.SymbolKindVariable, false)
	emptyMap := &ast.BLangMappingConstructorExpr{Fields: []model.MappingField{}}
	emptyMap.SetDeterminedType(semtypes.MAPPING)
	setPositionIfMissing(emptyMap, pos)
	mapVar := &ast.BLangSimpleVariable{Name: &ast.BLangIdentifier{Value: mapName}}
	mapVar.SetDeterminedType(semtypes.MAPPING)
	mapVar.SetInitialExpression(emptyMap)
	mapVar.SetSymbol(mapSymbol)
	setPositionIfMissing(mapVar, pos)
	mapVarDef := &ast.BLangSimpleVariableDef{Var: mapVar}
	setPositionIfMissing(mapVarDef, pos)
	*initStmts = append(*initStmts, mapVarDef)
	mapRef := &ast.BLangSimpleVarRef{VariableName: mapVar.Name}
	mapRef.SetSymbol(mapSymbol)
	mapRef.SetDeterminedType(semtypes.MAPPING)
	setPositionIfMissing(mapRef, pos)
	return mapRef
}

func createPreOrderLetStores(
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	startClauseIndex int,
	endClauseIndex int,
	initStmts *[]model.StatementNode,
	pos diagnostics.Location,
) ([]queryLetStore, bool) {
	var stores []queryLetStore
	for i := startClauseIndex; i < endClauseIndex; i++ {
		clause, isLet := queryExpr.QueryClauseList[i].(*ast.BLangLetClause)
		if !isLet {
			continue
		}
		for _, variableDef := range clause.LetVarDeclarations {
			varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
			if !ok || varDef.Var == nil || varDef.Var.Symbol() == (model.SymbolRef{}) {
				cx.unimplemented("query let clause currently supports only initialized simple variable declarations")
				return nil, false
			}
			valueTy := cx.symbolType(varDef.Var.Symbol())
			if valueTy == nil {
				valueTy = varDef.Var.GetDeterminedType()
			}
			if valueTy == nil {
				valueTy = semtypes.ANY
			}
			storeRef := createQueryListStore(cx, initStmts, pos)
			stores = append(stores, queryLetStore{
				VarDef:   varDef,
				StoreRef: storeRef,
				ValueTy:  valueTy,
			})
		}
	}
	return stores, true
}

func createQueryPayloadStore(
	cx *functionContext,
	initStmts *[]model.StatementNode,
	pos diagnostics.Location,
	letStores []queryLetStore,
) (*ast.BLangSimpleVarRef, bool) {
	payloadRef := createQueryListStore(cx, initStmts, pos)
	for _, store := range letStores {
		pushPayload := createPushInvocation(cx, payloadRef, store.StoreRef)
		if pushPayload == nil {
			return nil, false
		}
		pushStmt := &ast.BLangExpressionStmt{Expr: pushPayload}
		setPositionIfMissing(pushStmt, pos)
		*initStmts = append(*initStmts, pushStmt)
	}
	return payloadRef, true
}

func buildOrderKeyTupleExpr(
	cx *functionContext,
	orderByClause *ast.BLangOrderByClause,
	pos diagnostics.Location,
) (*ast.BLangListConstructorExpr, []model.StatementNode) {
	keyExprs := make([]ast.BLangExpression, 0, len(orderByClause.OrderByKeyList))
	var initStmts []model.StatementNode
	for i := range orderByClause.OrderByKeyList {
		keyResult := walkExpression(cx, orderByClause.OrderByKeyList[i].Expression)
		initStmts = append(initStmts, keyResult.initStmts...)
		keyExprs = append(keyExprs, keyResult.replacementNode.(ast.BLangExpression))
	}
	keyTuple := &ast.BLangListConstructorExpr{Exprs: keyExprs}
	keyTuple.SetDeterminedType(semtypes.LIST)
	keyTuple.AtomicType = semtypes.LIST_ATOMIC_INNER
	setPositionIfMissing(keyTuple, pos)
	return keyTuple, initStmts
}

func buildOrderDirectionExpr(orderByClause *ast.BLangOrderByClause, pos diagnostics.Location) *ast.BLangListConstructorExpr {
	directions := make([]ast.BLangExpression, 0, len(orderByClause.OrderByKeyList))
	for i := range orderByClause.OrderByKeyList {
		directions = append(directions, createBoolLiteral(orderByClause.OrderByKeyList[i].IsAscending, pos))
	}
	listExpr := &ast.BLangListConstructorExpr{Exprs: directions}
	listExpr.SetDeterminedType(semtypes.LIST)
	listExpr.AtomicType = semtypes.LIST_ATOMIC_INNER
	setPositionIfMissing(listExpr, pos)
	return listExpr
}

func queryElementAccess(
	collRef ast.BLangExpression,
	keysRef *ast.BLangSimpleVarRef,
	indexExpr ast.BLangExpression,
	elementTy semtypes.SemType,
) ast.BLangExpression {
	if keysRef == nil {
		listAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: indexExpr,
		}
		listAccess.Expr = collRef
		listAccess.SetDeterminedType(elementTy)
		return listAccess
	}
	keyAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: indexExpr,
	}
	keyAccess.Expr = keysRef
	keyAccess.SetDeterminedType(semtypes.STRING)
	mapAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: keyAccess,
	}
	mapAccess.Expr = collRef
	mapAccess.SetDeterminedType(elementTy)
	return mapAccess
}

func appendQuerySelectResultStmts(
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	resultRef *ast.BLangSimpleVarRef,
	selectClause *ast.BLangSelectClause,
	onConflictClause *ast.BLangOnConflictClause,
	seenKeysRef *ast.BLangSimpleVarRef,
	basePos diagnostics.Location,
	bodyStmts []ast.BLangStatement,
) ([]ast.BLangStatement, bool) {
	selectResult := walkExpression(cx, selectClause.Expression)
	bodyStmts = append(bodyStmts, selectResult.initStmts...)
	selectExpr := selectResult.replacementNode.(ast.BLangExpression)

	switch queryExpr.QueryConstructType {
	case model.TypeKind_MAP:
		selectTy := selectExpr.GetDeterminedType()
		pairName, pairSymbol := cx.addDesugardSymbol(selectTy, model.SymbolKindVariable, false)
		pairVar := &ast.BLangSimpleVariable{
			Name: &ast.BLangIdentifier{Value: pairName},
		}
		pairVar.SetDeterminedType(selectTy)
		pairVar.SetInitialExpression(selectExpr)
		pairVar.SetSymbol(pairSymbol)
		pairVarDef := &ast.BLangSimpleVariableDef{Var: pairVar}
		setPositionIfMissing(pairVarDef, basePos)
		bodyStmts = append(bodyStmts, pairVarDef)

		pairRef := &ast.BLangSimpleVarRef{
			VariableName: pairVar.Name,
		}
		pairRef.SetSymbol(pairSymbol)
		pairRef.SetDeterminedType(selectTy)

		keyAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: createIntLiteral(0),
		}
		keyAccess.Expr = pairRef
		keyAccess.SetDeterminedType(semtypes.STRING)

		valueAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: createIntLiteral(1),
		}
		valueAccess.Expr = pairRef
		valueAccess.SetDeterminedType(semtypes.ANY)

		if onConflictClause != nil {
			if seenKeysRef == nil {
				cx.internalError("on conflict query lowering requires seen-key map")
				return nil, false
			}
			seenLookup := &ast.BLangIndexBasedAccess{
				IndexExpr: keyAccess,
			}
			seenLookup.Expr = seenKeysRef
			seenLookup.SetDeterminedType(semtypes.ANY)
			conflictCond := &ast.BLangBinaryExpr{
				LhsExpr: seenLookup,
				RhsExpr: createBoolLiteral(true, basePos),
				OpKind:  model.OperatorKind_EQUAL,
			}
			conflictCond.SetDeterminedType(semtypes.BOOLEAN)

			conflictResult := walkExpression(cx, onConflictClause.Expression)
			conflictBody := make([]ast.BLangStatement, 0, len(conflictResult.initStmts)+2)
			conflictBody = append(conflictBody, conflictResult.initStmts...)

			conflictExpr := conflictResult.replacementNode.(ast.BLangExpression)
			conflictTy := conflictExpr.GetDeterminedType()
			conflictName, conflictSymbol := cx.addDesugardSymbol(conflictTy, model.SymbolKindVariable, false)
			conflictVar := &ast.BLangSimpleVariable{
				Name: &ast.BLangIdentifier{Value: conflictName},
			}
			conflictVar.SetDeterminedType(conflictTy)
			conflictVar.SetInitialExpression(conflictExpr)
			conflictVar.SetSymbol(conflictSymbol)
			conflictVarDef := &ast.BLangSimpleVariableDef{Var: conflictVar}
			setPositionIfMissing(conflictVarDef, basePos)
			conflictBody = append(conflictBody, conflictVarDef)

			conflictRef := &ast.BLangSimpleVarRef{
				VariableName: conflictVar.Name,
			}
			conflictRef.SetSymbol(conflictSymbol)
			conflictRef.SetDeterminedType(conflictTy)

			isErrorExpr := &ast.BLangTypeTestExpr{}
			isErrorExpr.Expr = conflictRef
			isErrorExpr.Type = model.TypeData{Type: semtypes.ERROR}
			isErrorExpr.SetDeterminedType(semtypes.BOOLEAN)

			assignResult := &ast.BLangAssignment{
				VarRef: resultRef,
				Expr:   conflictRef,
			}
			assignResult.SetDeterminedType(semtypes.NEVER)
			breakStmt := &ast.BLangBreak{}
			breakStmt.SetDeterminedType(semtypes.NEVER)
			errorBody := ast.BLangBlockStmt{
				Stmts: []ast.BLangStatement{assignResult, breakStmt},
			}
			errorIf := &ast.BLangIf{
				Expr: isErrorExpr,
				Body: errorBody,
			}
			errorIf.SetScope(cx.currentScope())
			errorIf.SetDeterminedType(semtypes.NEVER)
			conflictBody = append(conflictBody, errorIf)

			onConflictIf := &ast.BLangIf{
				Expr: conflictCond,
				Body: ast.BLangBlockStmt{
					Stmts: conflictBody,
				},
			}
			onConflictIf.SetScope(cx.currentScope())
			onConflictIf.SetDeterminedType(semtypes.NEVER)
			bodyStmts = append(bodyStmts, onConflictIf)

			markSeen := createMapPutAssignment(seenKeysRef, keyAccess, createBoolLiteral(true, basePos))
			setPositionIfMissing(markSeen, basePos)
			bodyStmts = append(bodyStmts, markSeen)
		}

		mapPutStmt := createMapPutAssignment(resultRef, keyAccess, valueAccess)
		setPositionIfMissing(mapPutStmt, basePos)
		bodyStmts = append(bodyStmts, mapPutStmt)
	default:
		pushInvocation := createPushInvocation(cx, resultRef, selectExpr)
		if pushInvocation == nil {
			return nil, false
		}
		bodyStmts = append(bodyStmts, &ast.BLangExpressionStmt{Expr: pushInvocation})
	}
	return bodyStmts, true
}

func appendQueryIntermediateClauseStmts(
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	idxRef ast.BLangExpression,
	initStmts *[]model.StatementNode,
	bodyStmts []ast.BLangStatement,
	startClauseIndex int,
	endClauseIndex int,
) ([]ast.BLangStatement, bool) {
	for i := startClauseIndex; i < endClauseIndex; i++ {
		switch clause := queryExpr.QueryClauseList[i].(type) {
		case *ast.BLangLetClause:
			for _, variableDef := range clause.LetVarDeclarations {
				varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
				if !ok || varDef.Var == nil || varDef.Var.Expr == nil {
					cx.unimplemented("query let clause currently supports only initialized simple variable declarations")
					return nil, false
				}
				letResult := walkExpression(cx, varDef.Var.Expr.(ast.BLangExpression))
				bodyStmts = append(bodyStmts, letResult.initStmts...)
				varDef.Var.SetInitialExpression(letResult.replacementNode.(ast.BLangExpression))
				bodyStmts = append(bodyStmts, varDef)
			}
		case *ast.BLangWhereClause:
			if clause.Expression == nil {
				cx.unimplemented("query where clause requires a condition expression")
				return nil, false
			}
			whereResult := walkExpression(cx, clause.Expression)
			bodyStmts = append(bodyStmts, whereResult.initStmts...)
			whereCond := whereResult.replacementNode.(ast.BLangExpression)
			notWhereCond := &ast.BLangUnaryExpr{
				Expr:     whereCond,
				Operator: model.OperatorKind_NOT,
			}
			notWhereCond.SetDeterminedType(semtypes.BOOLEAN)
			continueStmt := &ast.BLangContinue{}
			continueStmt.SetDeterminedType(semtypes.NEVER)
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
			filterIf.SetDeterminedType(semtypes.NEVER)
			bodyStmts = append(bodyStmts, filterIf)
		case *ast.BLangLimitClause:
			if clause.Expression == nil {
				cx.unimplemented("query limit clause requires a limit expression")
				return nil, false
			}

			limitCounterName, limitCounterSymbol := cx.addDesugardSymbol(semtypes.INT, model.SymbolKindVariable, false)
			limitCounterVar := &ast.BLangSimpleVariable{
				Name: &ast.BLangIdentifier{Value: limitCounterName},
			}
			limitCounterVar.SetDeterminedType(semtypes.INT)
			limitCounterVar.SetInitialExpression(createIntLiteral(0))
			limitCounterVar.SetSymbol(limitCounterSymbol)
			limitCounterVarDef := &ast.BLangSimpleVariableDef{Var: limitCounterVar}
			setPositionIfMissing(limitCounterVarDef, queryExpr.GetPosition())
			*initStmts = append(*initStmts, limitCounterVarDef)

			limitCounterRef := &ast.BLangSimpleVarRef{
				VariableName: limitCounterVar.Name,
			}
			limitCounterRef.SetSymbol(limitCounterSymbol)
			limitCounterRef.SetDeterminedType(semtypes.INT)

			limitResult := walkExpression(cx, clause.Expression)
			bodyStmts = append(bodyStmts, limitResult.initStmts...)
			limitExpr := limitResult.replacementNode.(ast.BLangExpression)

			reachedLimitCond := &ast.BLangBinaryExpr{
				LhsExpr: limitCounterRef,
				RhsExpr: limitExpr,
				OpKind:  model.OperatorKind_GREATER_EQUAL,
			}
			reachedLimitCond.SetDeterminedType(semtypes.BOOLEAN)

			continueStmt := &ast.BLangContinue{}
			continueStmt.SetDeterminedType(semtypes.NEVER)
			skipBody := ast.BLangBlockStmt{
				Stmts: []ast.BLangStatement{
					createIncrementStmt(idxRef),
					continueStmt,
				},
			}
			limitIf := &ast.BLangIf{
				Expr: reachedLimitCond,
				Body: skipBody,
			}
			limitIf.SetScope(cx.currentScope())
			limitIf.SetDeterminedType(semtypes.NEVER)
			bodyStmts = append(bodyStmts, limitIf)

			bodyStmts = append(bodyStmts, createIncrementStmt(limitCounterRef))
		case *ast.BLangOrderByClause:
			cx.unimplemented("query order by clause must be handled by dedicated lowering path")
			return nil, false
		default:
			cx.unimplemented("query expression currently supports only let + where + order by + limit clauses as intermediate clauses")
			return nil, false
		}
	}
	return bodyStmts, true
}

func createIntLiteral(value int64) *ast.BLangNumericLiteral {
	lit := &ast.BLangNumericLiteral{
		BLangLiteral: ast.BLangLiteral{
			Value:         value,
			OriginalValue: fmt.Sprintf("%d", value),
		},
		Kind: model.NodeKind_NUMERIC_LITERAL,
	}
	lit.SetDeterminedType(semtypes.INT)
	return lit
}

func createBoolLiteral(value bool, pos diagnostics.Location) *ast.BLangLiteral {
	originalValue := "false"
	if value {
		originalValue = "true"
	}
	lit := &ast.BLangLiteral{
		Value:         value,
		OriginalValue: originalValue,
	}
	lit.SetDeterminedType(semtypes.BOOLEAN)
	setPositionIfMissing(lit, pos)
	return lit
}

func createMapPutAssignment(mapExpr ast.BLangExpression, keyExpr ast.BLangExpression, valueExpr ast.BLangExpression) *ast.BLangAssignment {
	mapAccess := &ast.BLangIndexBasedAccess{
		IndexExpr: keyExpr,
	}
	mapAccess.Expr = mapExpr
	mapAccess.SetDeterminedType(semtypes.ANY)
	assign := &ast.BLangAssignment{
		VarRef: mapAccess,
		Expr:   valueExpr,
	}
	assign.SetDeterminedType(semtypes.NEVER)
	return assign
}

func createQuerySortInvocation(
	cx *functionContext,
	keysExpr ast.BLangExpression,
	directionsExpr ast.BLangExpression,
	indicesExpr ast.BLangExpression,
	payloadExpr ast.BLangExpression,
) *ast.BLangInvocation {
	pkgName := array.PackageName
	space, ok := cx.getImportedSymbolSpace(pkgName)
	if !ok {
		cx.internalError(pkgName + " symbol space not found")
		return nil
	}
	symbolRef, ok := space.GetSymbol("querySort")
	if !ok {
		cx.internalError(pkgName + ":querySort symbol not found")
		return nil
	}
	cx.addImplicitImport(pkgName, ast.BLangImportPackage{
		OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
		PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
		Alias:        &ast.BLangIdentifier{Value: pkgName},
	})
	inv := &ast.BLangInvocation{
		Name:     &ast.BLangIdentifier{Value: "querySort"},
		PkgAlias: &ast.BLangIdentifier{Value: pkgName},
		ArgExprs: []ast.BLangExpression{keysExpr, directionsExpr, indicesExpr, payloadExpr},
	}
	inv.SetSymbol(symbolRef)
	inv.SetDeterminedType(semtypes.NIL)
	setPositionIfMissing(inv, keysExpr.GetPosition())
	return inv
}

func createPushInvocation(cx *functionContext, listExpr ast.BLangExpression, valueExpr ast.BLangExpression) *ast.BLangInvocation {
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
	inv.SetDeterminedType(semtypes.NIL)
	return inv
}

func isNilLiftableBinaryOp(op model.OperatorKind) bool {
	switch op {
	case model.OperatorKind_ADD, model.OperatorKind_SUB,
		model.OperatorKind_MUL, model.OperatorKind_DIV, model.OperatorKind_MOD,
		model.OperatorKind_BITWISE_LEFT_SHIFT, model.OperatorKind_BITWISE_RIGHT_SHIFT,
		model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT,
		model.OperatorKind_BITWISE_AND, model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
		return true
	default:
		return false
	}
}

func isNilLiftableUnaryOp(op model.OperatorKind) bool {
	switch op {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		return true
	default:
		return false
	}
}

func createOperandTempVar(cx *functionContext, ty semtypes.SemType, initExpr ast.BLangExpression, pos diagnostics.Location, initStmts []model.StatementNode) (*ast.BLangIdentifier, model.SymbolRef, []model.StatementNode) {
	name, symbol := cx.addDesugardSymbol(ty, model.SymbolKindVariable, false)
	varName := &ast.BLangIdentifier{Value: name}
	tempVar := &ast.BLangSimpleVariable{Name: varName}
	tempVar.SetDeterminedType(ty)
	tempVar.SetInitialExpression(initExpr)
	tempVar.SetSymbol(symbol)
	varDef := &ast.BLangSimpleVariableDef{Var: tempVar}
	varDef.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(varDef, pos)
	return varName, symbol, append(initStmts, varDef)
}

func createNilResultVar(cx *functionContext, ty semtypes.SemType, pos diagnostics.Location, initStmts []model.StatementNode) (*ast.BLangIdentifier, model.SymbolRef, []model.StatementNode) {
	nilLit := &ast.BLangLiteral{Value: nil}
	nilLit.SetDeterminedType(semtypes.NIL)
	setPositionIfMissing(nilLit, pos)

	name, symbol := cx.addDesugardSymbol(ty, model.SymbolKindVariable, false)
	varName := &ast.BLangIdentifier{Value: name}
	tempVar := &ast.BLangSimpleVariable{Name: varName}
	tempVar.SetDeterminedType(ty)
	tempVar.SetInitialExpression(nilLit)
	tempVar.SetSymbol(symbol)
	varDef := &ast.BLangSimpleVariableDef{Var: tempVar}
	varDef.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(varDef, pos)
	return varName, symbol, append(initStmts, varDef)
}

func createNilTypeTest(varName *ast.BLangIdentifier, symbol model.SymbolRef, ty semtypes.SemType, pos diagnostics.Location) *ast.BLangTypeTestExpr {
	ref := createVarRef(varName, symbol, ty)
	typeTest := &ast.BLangTypeTestExpr{
		Expr: ref,
		Type: model.TypeData{Type: semtypes.NIL},
	}
	typeTest.SetDeterminedType(semtypes.BOOLEAN)
	setPositionIfMissing(typeTest, pos)
	return typeTest
}

func createVarRef(varName *ast.BLangIdentifier, symbol model.SymbolRef, ty semtypes.SemType) *ast.BLangSimpleVarRef {
	ref := &ast.BLangSimpleVarRef{VariableName: varName}
	ref.SetSymbol(symbol)
	ref.SetDeterminedType(ty)
	return ref
}

func createResultAssignment(resultVarName *ast.BLangIdentifier, resultSymbol model.SymbolRef, resultTy semtypes.SemType, valueExpr ast.BLangExpression, pos diagnostics.Location) *ast.BLangAssignment {
	varRef := createVarRef(resultVarName, resultSymbol, resultTy)
	assign := &ast.BLangAssignment{
		VarRef: varRef,
		Expr:   valueExpr,
	}
	assign.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(assign, pos)
	return assign
}
