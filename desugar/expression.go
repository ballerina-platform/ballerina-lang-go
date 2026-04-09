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

func walkExpression(cx *functionContext, node ast.BLangActionOrExpression) desugaredNode[ast.BLangActionOrExpression] {
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
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	case *ast.BLangNumericLiteral:
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	case *ast.BLangSimpleVarRef:
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	case *ast.BLangLocalVarRef:
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	case *ast.BLangConstRef:
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	case *ast.BLangNewExpression:
		return walkNewExpression(cx, expr)
	case *ast.BLangNamedArgsExpression:
		result := walkExpression(cx, expr.Expr)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       result.initStmts,
			replacementNode: expr,
		}
	case *ast.BLangWildCardBindingPattern:
		// Wildcard binding pattern can appear in variable references (e.g., _ = expr)
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", node))
	}
}

func walkBinaryExpr(cx *functionContext, expr *ast.BLangBinaryExpr) desugaredNode[ast.BLangActionOrExpression] {
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
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	lhsTy := expr.LhsExpr.GetDeterminedType()
	rhsTy := expr.RhsExpr.GetDeterminedType()
	if lhsTy == nil || rhsTy == nil {
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}
	lhsHasNil := semtypes.ContainsBasicType(lhsTy, semtypes.NIL)
	rhsHasNil := semtypes.ContainsBasicType(rhsTy, semtypes.NIL)

	if !lhsHasNil && !rhsHasNil {
		return desugaredNode[ast.BLangActionOrExpression]{
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: replacementRef,
	}
}

func walkUnaryExpr(cx *functionContext, expr *ast.BLangUnaryExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	// Unary + is identity — desugar to just the operand (BIR gen doesn't handle unary +)
	if expr.Operator == model.OperatorKind_ADD {
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       initStmts,
			replacementNode: expr.Expr,
		}
	}

	if !isNilLiftableUnaryOp(expr.Operator) {
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	operandTy := expr.Expr.GetDeterminedType()
	if !semtypes.ContainsBasicType(operandTy, semtypes.NIL) {
		return desugaredNode[ast.BLangActionOrExpression]{
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: replacementRef,
	}
}

func walkElvisExpr(cx *functionContext, expr *ast.BLangElvisExpr) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkGroupExpr(cx *functionContext, expr *ast.BLangGroupExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkIndexBasedAccess(cx *functionContext, expr *ast.BLangIndexBasedAccess) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkFieldBaseAccess(cx *functionContext, expr *ast.BLangFieldBaseAccess) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: indexAccess,
	}
}

func walkInvocation(cx *functionContext, expr *ast.BLangInvocation) desugaredNode[ast.BLangActionOrExpression] {
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
		return desugaredNode[ast.BLangActionOrExpression]{
			initStmts:       initStmts,
			replacementNode: expr,
		}
	}

	directStmts := walkDirectCallArgs(cx, expr, fnSym)
	initStmts = append(initStmts, directStmts...)

	return desugaredNode[ast.BLangActionOrExpression]{
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

func walkListConstructorExpr(cx *functionContext, expr *ast.BLangListConstructorExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	for i := range expr.Exprs {
		result := walkExpression(cx, expr.Exprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.Exprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkErrorConstructorExpr(cx *functionContext, expr *ast.BLangErrorConstructorExpr) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCheckedExpr(cx *functionContext, expr *ast.BLangCheckedExpr) desugaredNode[ast.BLangActionOrExpression] {
	return desugarCheckedExpr(cx, expr, false)
}

func walkCheckPanickedExpr(cx *functionContext, expr *ast.BLangCheckPanickedExpr) desugaredNode[ast.BLangActionOrExpression] {
	return desugarCheckedExpr(cx, &expr.BLangCheckedExpr, true)
}

func walkTrapExpr(cx *functionContext, expr *ast.BLangTrapExpr) desugaredNode[ast.BLangActionOrExpression] {
	result := walkExpression(cx, expr.Expr)
	if len(result.initStmts) > 0 {
		// I don't think this can ever happen but if it does we need to think about how to add these statements in to the
		// trap region in BIR gen
		cx.internalError("Init statements will be hoisted outside of trap region")
	}
	expr.Expr = result.replacementNode.(ast.BLangExpression)
	return desugaredNode[ast.BLangActionOrExpression]{initStmts: nil, replacementNode: expr}
}

func desugarCheckedExpr(cx *functionContext, expr *ast.BLangCheckedExpr, isPanic bool) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	// Walk the inner expression first
	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: replacementVarRef,
	}
}

func walkDynamicArgExpr(cx *functionContext, expr *ast.BLangDynamicArgExpr) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkLambdaFunction(cx *functionContext, expr *ast.BLangLambdaFunction) desugaredNode[ast.BLangActionOrExpression] {
	// Desugar the function body
	if expr.Function != nil {
		expr.Function = desugarFunction(cx.pkgCtx, expr.Function)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		replacementNode: expr,
	}
}

func walkTypeConversionExpr(cx *functionContext, expr *ast.BLangTypeConversionExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	if expr.Expression != nil {
		result := walkExpression(cx, expr.Expression)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expression = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkTypeTestExpr(cx *functionContext, expr *ast.BLangTypeTestExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkAnnotAccessExpr(cx *functionContext, expr *ast.BLangAnnotAccessExpr) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	if expr.Expr != nil {
		result := walkExpression(cx, expr.Expr)
		initStmts = append(initStmts, result.initStmts...)
		expr.Expr = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkCollectContextInvocation(cx *functionContext, expr *ast.BLangCollectContextInvocation) desugaredNode[ast.BLangActionOrExpression] {
	// Walk the underlying invocation
	result := walkInvocation(cx, &expr.Invocation)
	expr.Invocation = *result.replacementNode.(*ast.BLangInvocation)

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       result.initStmts,
		replacementNode: expr,
	}
}

func walkArrowFunction(cx *functionContext, expr *ast.BLangArrowFunction) desugaredNode[ast.BLangActionOrExpression] {
	// Arrow functions have a body that may need desugaring
	if expr.Body != nil {
		result := walkExpression(cx, expr.Body.Expr.(ast.BLangActionOrExpression))
		expr.Body.Expr = result.replacementNode
		// Handle initStmts if needed - arrow functions may need special handling
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		replacementNode: expr,
	}
}

func walkNewExpression(cx *functionContext, expr *ast.BLangNewExpression) desugaredNode[ast.BLangActionOrExpression] {
	var initStmts []model.StatementNode

	for i := range expr.ArgsExprs {
		result := walkExpression(cx, expr.ArgsExprs[i])
		initStmts = append(initStmts, result.initStmts...)
		expr.ArgsExprs[i] = result.replacementNode.(ast.BLangExpression)
	}

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkMappingConstructorExpr(cx *functionContext, expr *ast.BLangMappingConstructorExpr) desugaredNode[ast.BLangActionOrExpression] {
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

	return desugaredNode[ast.BLangActionOrExpression]{
		initStmts:       initStmts,
		replacementNode: expr,
	}
}

func walkQueryExpr(cx *functionContext, expr *ast.BLangQueryExpr) desugaredNode[ast.BLangActionOrExpression] {
	fromClause, ok := expr.QueryClauseList[0].(*ast.BLangFromClause)
	if !ok {
		cx.internalError("query expression must start with from clause")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}
	selectClause, ok := expr.QueryClauseList[len(expr.QueryClauseList)-1].(*ast.BLangSelectClause)
	if !ok {
		cx.internalError("query expression must end with select clause")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}
	loopVarDef, ok := fromClause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
	if !ok {
		cx.unimplemented("query from clause currently supports only simple variable definition")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}
	cloneLoopVarDef := cloneSimpleVariableDef(loopVarDef)
	if cloneLoopVarDef == nil || cloneLoopVarDef.Var == nil {
		cx.internalError("failed to clone query from clause variable definition")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}

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

	zeroLiteral := &ast.BLangNumericLiteral{
		BLangLiteral: ast.BLangLiteral{
			Value:         int64(0),
			OriginalValue: "0",
		},
		Kind: model.NodeKind_NUMERIC_LITERAL,
	}
	zeroLiteral.SetDeterminedType(semtypes.INT)

	idxName, idxSymbol := cx.addDesugardSymbol(semtypes.INT, model.SymbolKindVariable, false)
	idxVar := &ast.BLangSimpleVariable{
		Name: &ast.BLangIdentifier{Value: idxName},
	}
	idxVar.SetDeterminedType(semtypes.INT)
	idxVar.SetInitialExpression(zeroLiteral)
	idxVar.SetSymbol(idxSymbol)
	idxVarDef := &ast.BLangSimpleVariableDef{Var: idxVar}
	setPositionIfMissing(idxVarDef, basePos)
	initStmts = append(initStmts, idxVarDef)

	idxRef := &ast.BLangSimpleVarRef{
		VariableName: idxVar.Name,
	}
	idxRef.SetSymbol(idxSymbol)
	idxRef.SetDeterminedType(semtypes.INT)

	loopVarSymbol := cloneLoopVarDef.Var.Symbol()
	loopVarTy := cx.symbolType(loopVarSymbol)
	if loopVarTy == nil {
		cx.internalError("query from clause variable symbol type not found")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}

	lengthSource := collRef
	var elementAccess ast.BLangExpression

	switch {
	case semtypes.IsSubtypeSimple(collTy, semtypes.LIST):
		listAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: idxRef,
		}
		listAccess.Expr = collRef
		listAccess.SetDeterminedType(loopVarTy)
		elementAccess = listAccess
	case semtypes.IsSubtypeSimple(collTy, semtypes.MAPPING):
		keysInvocation := createKeysInvocation(cx, collRef)
		if keysInvocation == nil {
			return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
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

		keysRef := &ast.BLangSimpleVarRef{
			VariableName: keysVar.Name,
		}
		keysRef.SetSymbol(keysSymbol)
		keysRef.SetDeterminedType(keysTy)
		lengthSource = keysRef

		keyAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: idxRef,
		}
		keyAccess.Expr = keysRef
		keyAccess.SetDeterminedType(semtypes.STRING)

		mapAccess := &ast.BLangIndexBasedAccess{
			IndexExpr: keyAccess,
		}
		mapAccess.Expr = collRef
		mapAccess.SetDeterminedType(loopVarTy)
		elementAccess = mapAccess
	default:
		cx.unimplemented("query from clause currently supports only list or map collections")
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
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

	condition := &ast.BLangBinaryExpr{
		LhsExpr: idxRef,
		RhsExpr: lenRef,
		OpKind:  model.OperatorKind_LESS_THAN,
	}
	condition.SetDeterminedType(semtypes.BOOLEAN)
	cloneLoopVarDef.Var.SetInitialExpression(elementAccess)

	var bodyStmts []ast.BLangStatement
	bodyStmts = append(bodyStmts, cloneLoopVarDef)

	bodyStmts, ok = appendQueryIntermediateClauseStmts(cx, expr, idxRef, &initStmts, bodyStmts)
	if !ok {
		return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
	}

	selectResult := walkExpression(cx, selectClause.Expression)
	bodyStmts = append(bodyStmts, selectResult.initStmts...)

	selectExpr := selectResult.replacementNode.(ast.BLangExpression)
	switch expr.QueryConstructType {
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

		mapPutStmt := createMapPutAssignment(resultRef, keyAccess, valueAccess)
		bodyStmts = append(bodyStmts, mapPutStmt)
	default:
		pushInvocation := createPushInvocation(cx, resultRef, selectExpr)
		if pushInvocation == nil {
			return desugaredNode[ast.BLangActionOrExpression]{replacementNode: expr}
		}
		bodyStmts = append(bodyStmts, &ast.BLangExpressionStmt{Expr: pushInvocation})
	}
	bodyStmts = append(bodyStmts, createIncrementStmt(idxRef))

	whileStmt := &ast.BLangWhile{
		Expr: condition,
		Body: ast.BLangBlockStmt{
			Stmts: bodyStmts,
		},
	}
	whileStmt.SetScope(cx.currentScope())
	whileStmt.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(whileStmt, basePos)
	initStmts = append(initStmts, whileStmt)

	setPositionIfMissing(resultRef, basePos)

	return desugaredNode[ast.BLangActionOrExpression]{
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
	cx *functionContext,
	queryExpr *ast.BLangQueryExpr,
	idxRef ast.BLangExpression,
	initStmts *[]model.StatementNode,
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
				bodyStmts = append(bodyStmts, letResult.initStmts...)
				varDef.Var.SetInitialExpression(letResult.replacementNode)
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
		default:
			cx.unimplemented("query expression currently supports only let + where + limit clauses as intermediate clauses")
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
