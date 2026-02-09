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

package bir

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"fmt"
)

// Since BLangNodeVisitor is anyway deprecated in jBallerina, we'll try to do this more cleanly
// TODO: may be we should have this in a separate package and keep BIR package clean (only definitions)

type Context struct {
	CompilerContext *context.CompilerContext
	constantMap     map[model.SymbolRef]*BIRConstant
	importAliasMap  map[string]*model.PackageID // Maps import alias to package ID
	packageID       *model.PackageID            // Current package ID
}

type stmtContext struct {
	birCx       *Context
	bbs         []*BIRBasicBlock
	localVars   []*BIRVariableDcl
	retVar      *BIROperand
	scope       *BIRScope
	nextScopeId int
	// TODO: do better
	varMap  map[model.SymbolRef]*BIROperand
	loopCtx *loopContext
}

type loopContext struct {
	onBreakBB    *BIRBasicBlock
	onContinueBB *BIRBasicBlock
	enclosing    *loopContext
}

func (cx *stmtContext) addLoopCtx(onBreakBB *BIRBasicBlock, onContinueBB *BIRBasicBlock) *loopContext {
	newCtx := &loopContext{
		onBreakBB:    onBreakBB,
		onContinueBB: onContinueBB,
		enclosing:    cx.loopCtx,
	}
	cx.loopCtx = newCtx
	return newCtx
}

func (cx *stmtContext) popLoopCtx() {
	if cx.loopCtx == nil {
		panic("no enclosing loop context")
	}
	cx.loopCtx = cx.loopCtx.enclosing
}

func (cx *stmtContext) addLocalVarInner(name model.Name, ty semtypes.SemType, kind VarKind) *BIROperand {
	varDcl := &BIRVariableDcl{}
	varDcl.Name = name
	varDcl.Type = ty
	varDcl.Kind = kind
	varDcl.Scope = VAR_SCOPE_FUNCTION
	varDcl.MetaVarName = name.Value()
	cx.localVars = append(cx.localVars, varDcl)
	return &BIROperand{VariableDcl: varDcl, Index: len(cx.localVars) - 1}
}

func (cx *stmtContext) addTempVar(ty semtypes.SemType) *BIROperand {
	return cx.addLocalVarInner(model.Name(fmt.Sprintf("%%%d", len(cx.localVars))), ty, VAR_KIND_TEMP)
}

func (cx *stmtContext) addLocalVar(name model.Name, ty semtypes.SemType, kind VarKind, symbol model.Symbol) *BIROperand {
	operand := cx.addLocalVarInner(name, ty, kind)
	symRef := cx.birCx.CompilerContext.RefSymbol(symbol)
	cx.varMap[symRef] = operand
	return operand
}

func (cx *stmtContext) addBB() *BIRBasicBlock {
	index := len(cx.bbs)
	bb := BB(index)
	cx.bbs = append(cx.bbs, &bb)
	return &bb
}

func GenBir(ctx *context.CompilerContext, ast *ast.BLangPackage) *BIRPackage {
	birPkg := &BIRPackage{}
	birPkg.PackageID = ast.PackageID
	genCtx := &Context{
		CompilerContext: ctx,
		constantMap:     make(map[model.SymbolRef]*BIRConstant),
		importAliasMap:  make(map[string]*model.PackageID),
		packageID:       ast.PackageID,
	}
	processImports(ctx, genCtx, ast.Imports, birPkg)
	for _, typeDef := range ast.TypeDefinitions {
		birPkg.TypeDefs = appendIfNotNil(birPkg.TypeDefs, TransformTypeDefinition(genCtx, &typeDef))
	}
	for _, globalVar := range ast.GlobalVars {
		birPkg.GlobalVars = appendIfNotNil(birPkg.GlobalVars, TransformGlobalVariableDcl(genCtx, &globalVar))
	}
	for _, constant := range ast.Constants {
		c := TransformConstant(genCtx, &constant)
		symRef := genCtx.CompilerContext.RefSymbol(constant.Symbol())
		genCtx.constantMap[symRef] = c
		birPkg.Constants = appendIfNotNil(birPkg.Constants, c)
	}
	for _, function := range ast.Functions {
		birFunc := TransformFunction(genCtx, &function)
		birPkg.Functions = append(birPkg.Functions, *birFunc)
		if birFunc.Name.Value() == "main" {
			birPkg.MainFunction = birFunc
		}
	}
	return birPkg
}

func processImports(compilerCtx *context.CompilerContext, genCtx *Context, imports []ast.BLangImportPackage, birPkg *BIRPackage) {
	// Add implicit imports
	for _, importPkg := range imports {
		if importPkg.Alias != nil && importPkg.Alias.Value != "" {
			var orgName model.Name
			if importPkg.OrgName != nil && importPkg.OrgName.Value != "" {
				orgName = model.Name(importPkg.OrgName.Value)
			} else {
				orgName = model.ANON_ORG
			}
			var nameComps []model.Name
			if len(importPkg.PkgNameComps) > 0 {
				for _, comp := range importPkg.PkgNameComps {
					nameComps = append(nameComps, model.Name(comp.Value))
				}
			} else {
				nameComps = []model.Name{model.DEFAULT_PACKAGE}
			}
			var version model.Name
			if importPkg.Version != nil && importPkg.Version.Value != "" {
				version = model.Name(importPkg.Version.Value)
			} else {
				version = model.DEFAULT_VERSION
			}
			pkgID := compilerCtx.NewPackageID(orgName, nameComps, version)
			genCtx.importAliasMap[importPkg.Alias.Value] = pkgID
		}
		birPkg.ImportModules = appendIfNotNil(birPkg.ImportModules, TransformImportModule(genCtx, importPkg))
	}
}

func TransformImportModule(ctx *Context, ast ast.BLangImportPackage) *BIRImportModule {
	// FIXME: fix this when we have symbol resolution, given only import we support is io we are going to hardcode it
	orgName := model.Name("ballerina")
	pkgName := model.Name("io")
	version := model.Name("0.0.0")
	return &BIRImportModule{
		PackageID: &model.PackageID{
			OrgName: &orgName,
			PkgName: &pkgName,
			Version: &version,
		},
	}
}

func TransformTypeDefinition(ctx *Context, ast *ast.BLangTypeDefinition) *BIRTypeDefinition {
	// FIXME: implement this
	return nil
}

func TransformGlobalVariableDcl(ctx *Context, ast *ast.BLangSimpleVariable) *BIRGlobalVariableDcl {
	var name, originalName model.Name
	name = model.Name(ast.GetName().GetValue())
	originalName = name
	birVarDcl := &BIRGlobalVariableDcl{}
	birVarDcl.Pos = ast.GetPosition()
	birVarDcl.Name = name
	birVarDcl.OriginalName = originalName
	birVarDcl.Scope = VAR_SCOPE_GLOBAL
	birVarDcl.Kind = VAR_KIND_GLOBAL
	birVarDcl.MetaVarName = name.Value()
	return birVarDcl
}

func TransformFunction(ctx *Context, astFunc *ast.BLangFunction) *BIRFunction {
	funcName := model.Name(astFunc.GetName().GetValue())
	birFunc := &BIRFunction{}
	birFunc.Pos = astFunc.GetPosition()
	birFunc.Name = funcName
	birFunc.OriginalName = funcName
	orgName := ctx.packageID.OrgName.Value()
	pkgName := ctx.packageID.PkgName.Value()
	moduleKey := orgName + "/" + pkgName
	birFunc.FunctionLookupKey = moduleKey + ":" + funcName.Value()
	common.Assert(astFunc.Receiver == nil)
	stmtCx := &stmtContext{birCx: ctx, varMap: make(map[model.SymbolRef]*BIROperand)}
	stmtCx.retVar = stmtCx.addLocalVarInner(model.Name("%0"), nil, VAR_KIND_RETURN)
	for _, param := range astFunc.RequiredParams {
		stmtCx.addLocalVar(model.Name(param.GetName().GetValue()), nil, VAR_KIND_ARG, param.Symbol())
	}
	switch body := astFunc.Body.(type) {
	case *ast.BLangBlockFunctionBody:
		handleBlockFunctionBody(stmtCx, body)
	case *ast.BLangExprFunctionBody:
		handleExprFunctionBody(stmtCx, body)
	default:
		panic("unexpected function body type")
	}
	for _, bbPtr := range stmtCx.bbs {
		birFunc.BasicBlocks = append(birFunc.BasicBlocks, *bbPtr)
	}
	for _, varPtr := range stmtCx.localVars {
		birFunc.LocalVars = append(birFunc.LocalVars, *varPtr)
	}
	birFunc.ReturnVariable = stmtCx.retVar.VariableDcl
	return birFunc
}

func TransformConstant(ctx *Context, c *ast.BLangConstant) *BIRConstant {
	valueExpr := c.Expr
	if literal, ok := valueExpr.(*ast.BLangLiteral); ok {
		// FIXME: once we have constant propagation these should be propagated and no longer needed
		return &BIRConstant{
			Name: model.Name(c.GetName().GetValue()),
			ConstValue: ConstValue{
				Value: literal.Value,
			},
		}
	}
	// TODO: need this think how to actually implement constant value initialization. May be we add these to init function?
	panic("unexpected constant value type")
}

func handleBlockFunctionBody(ctx *stmtContext, ast *ast.BLangBlockFunctionBody) {
	curBB := ctx.addBB()
	for _, stmt := range ast.Stmts {
		effect := handleStatement(ctx, curBB, stmt)
		curBB = effect.block
	}
	// Add implicit return
	if curBB != nil {
		curBB.Terminator = &Return{}
	}
}

type statementEffect struct {
	block *BIRBasicBlock
}

func handleStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt ast.BLangStatement) statementEffect {
	switch stmt := stmt.(type) {
	case *ast.BLangExpressionStmt:
		return expressionStatement(ctx, curBB, stmt)
	case *ast.BLangIf:
		return ifStatement(ctx, curBB, stmt)
	case *ast.BLangBlockStmt:
		return blockStatement(ctx, curBB, stmt)
	case *ast.BLangReturn:
		return returnStatement(ctx, curBB, stmt)
	case *ast.BLangSimpleVariableDef:
		return simpleVariableDefinition(ctx, curBB, stmt)
	case *ast.BLangAssignment:
		return assignmentStatement(ctx, curBB, stmt)
	case *ast.BLangCompoundAssignment:
		return compoundAssignment(ctx, curBB, stmt)
	case *ast.BLangWhile:
		return whileStatement(ctx, curBB, stmt)
	case *ast.BLangBreak:
		return breakStatement(ctx, curBB, stmt)
	case *ast.BLangContinue:
		return continueStatement(ctx, curBB, stmt)
	default:
		panic("unexpected statement type")
	}
}

func compoundAssignment(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangCompoundAssignment) statementEffect {
	// First do the operation
	ref := stmt.VarRef.(ast.BLangExpression)
	valueEffect := binaryExpressionInner(ctx, curBB, stmt.OpKind, ref, stmt.Expr, stmt.Expr.GetDeterminedType())
	// Then do the assignment
	return assignmentStatementInner(ctx, valueEffect.block, ref, valueEffect)
}

func continueStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangContinue) statementEffect {
	onContinueBB := ctx.loopCtx.onContinueBB
	curBB.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: onContinueBB}}
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func breakStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangBreak) statementEffect {
	onBreakBB := ctx.loopCtx.onBreakBB
	curBB.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: onBreakBB}}
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func whileStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangWhile) statementEffect {
	loopHead := ctx.addBB()
	// jump to loop head
	bb.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: loopHead}}
	condEffect := handleExpression(ctx, loopHead, stmt.Expr)

	loopBody := ctx.addBB()
	loopEnd := ctx.addBB()
	// conditionally jump to loop body
	branch := &Branch{}
	branch.Op = condEffect.result
	branch.TrueBB = loopBody
	branch.FalseBB = loopEnd
	condEffect.block.Terminator = branch

	ctx.addLoopCtx(loopEnd, loopHead)
	bodyEffect := blockStatement(ctx, loopBody, &stmt.Body)
	// This could happen if the while block always ends return, break or continue
	if bodyEffect.block != nil {
		bodyEffect.block.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: loopHead}}
	}

	ctx.popLoopCtx()
	return statementEffect{
		block: loopEnd,
	}
}

func assignmentStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangAssignment) statementEffect {
	valueEffect := handleExpression(ctx, bb, stmt.Expr)
	return assignmentStatementInner(ctx, valueEffect.block, stmt.VarRef, valueEffect)
}

func assignmentStatementInner(ctx *stmtContext, bb *BIRBasicBlock, ref ast.BLangExpression, valueEffect expressionEffect) statementEffect {
	switch varRef := ref.(type) {
	case *ast.BLangIndexBasedAccess:
		return assignToMemberStatement(ctx, bb, varRef, valueEffect)
	case *ast.BLangWildCardBindingPattern:
		return assignToWildcardBindingPattern(ctx, bb, varRef, valueEffect)
	case *ast.BLangSimpleVarRef:
		return assignToSimpleVariable(ctx, bb, varRef, valueEffect)
	default:
		panic("unexpected variable reference type")
	}
}

func assignToWildcardBindingPattern(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangWildCardBindingPattern, valueEffect expressionEffect) statementEffect {
	refEffect := wildcardBindingPattern(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := &Move{}
	mov.LhsOp = refEffect.result
	mov.RhsOp = valueEffect.result
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToSimpleVariable(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangSimpleVarRef, valueEffect expressionEffect) statementEffect {
	refEffect := simpleVariableReference(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := &Move{}
	mov.LhsOp = refEffect.result
	mov.RhsOp = valueEffect.result
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToMemberStatement(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangIndexBasedAccess, valueEffect expressionEffect) statementEffect {
	currBB := valueEffect.block
	containerRefEffect := handleExpression(ctx, currBB, varRef.Expr)
	currBB = containerRefEffect.block
	indexEffect := handleExpression(ctx, currBB, varRef.IndexExpr)
	currBB = indexEffect.block
	fieldAccess := &FieldAccess{}
	fieldAccess.Kind = INSTRUCTION_KIND_ARRAY_STORE
	fieldAccess.LhsOp = containerRefEffect.result
	fieldAccess.KeyOp = indexEffect.result
	fieldAccess.RhsOp = valueEffect.result
	currBB.Instructions = append(currBB.Instructions, fieldAccess)
	return statementEffect{
		block: currBB,
	}
}

func simpleVariableDefinition(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangSimpleVariableDef) statementEffect {
	varName := model.Name(stmt.Var.GetName().GetValue())
	if stmt.Var.Expr == nil {
		ctx.addLocalVar(varName, nil, VAR_KIND_LOCAL, stmt.Var.Symbol())
		// just declare the variable
		return statementEffect{
			block: bb,
		}
	}
	exprResult := handleExpression(ctx, bb, stmt.Var.Expr.(ast.BLangExpression))
	curBB := exprResult.block
	move := &Move{}
	move.LhsOp = ctx.addLocalVar(varName, nil, VAR_KIND_LOCAL, stmt.Var.Symbol())
	move.RhsOp = exprResult.result
	curBB.Instructions = append(curBB.Instructions, move)
	return statementEffect{
		block: curBB,
	}
}

func returnStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangReturn) statementEffect {
	curBB := bb
	if stmt.Expr != nil {
		valueEffect := handleExpression(ctx, curBB, stmt.Expr)
		curBB = valueEffect.block
		mov := &Move{}
		mov.LhsOp = ctx.retVar
		mov.RhsOp = valueEffect.result
		curBB.Instructions = append(curBB.Instructions, mov)
	}
	curBB.Terminator = &Return{}
	return statementEffect{}
}

func expressionStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangExpressionStmt) statementEffect {
	result := handleExpression(ctx, curBB, stmt.Expr)
	// We are ignoring the expression result (We can have one for things like call)
	return statementEffect{
		block: result.block,
	}
}

func ifStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangIf) statementEffect {
	cond := handleExpression(ctx, curBB, stmt.Expr)
	curBB = cond.block
	thenBB := ctx.addBB()
	var finalBB *BIRBasicBlock
	thenEffect := blockStatement(ctx, thenBB, &stmt.Body)
	// TODO: refactor this
	if stmt.ElseStmt != nil {
		elseBB := ctx.addBB()
		// Add branch to current BB
		branch := &Branch{}
		branch.Op = cond.result
		branch.TrueBB = thenBB
		branch.FalseBB = elseBB
		curBB.Terminator = branch

		elseEffect := handleStatement(ctx, elseBB, stmt.ElseStmt)
		finalBB = ctx.addBB()
		if elseEffect.block != nil {
			elseEffect.block.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: finalBB}}
		}
	} else {
		finalBB = ctx.addBB()
		branch := &Branch{}
		branch.Op = cond.result
		branch.TrueBB = thenBB
		branch.FalseBB = finalBB
		curBB.Terminator = branch
	}
	// this could be nil if the control flow moved out of the if (ex: break, continue, return, etc)
	if thenEffect.block != nil {
		thenEffect.block.Terminator = &Goto{BIRTerminatorBase: BIRTerminatorBase{ThenBB: finalBB}}
	}
	return statementEffect{
		block: finalBB,
	}
}

func blockStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangBlockStmt) statementEffect {
	curBB := bb
	for _, stmt := range stmt.Stmts {
		effect := handleStatement(ctx, curBB, stmt)
		curBB = effect.block
	}
	return statementEffect{
		block: curBB,
	}
}

func handleExprFunctionBody(ctx *stmtContext, ast *ast.BLangExprFunctionBody) {
	panic("unimplemented")
}

type expressionEffect struct {
	result *BIROperand
	block  *BIRBasicBlock
}

func handleExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr ast.BLangExpression) expressionEffect {
	switch expr := expr.(type) {
	case *ast.BLangInvocation:
		return invocation(ctx, curBB, expr)
	case *ast.BLangLiteral:
		return literal(ctx, curBB, expr)
	case *ast.BLangBinaryExpr:
		return binaryExpression(ctx, curBB, expr)
	case *ast.BLangSimpleVarRef:
		return simpleVariableReference(ctx, curBB, expr)
	case *ast.BLangUnaryExpr:
		return unaryExpression(ctx, curBB, expr)
	case *ast.BLangWildCardBindingPattern:
		return wildcardBindingPattern(ctx, curBB, expr)
	case *ast.BLangGroupExpr:
		return groupExpression(ctx, curBB, expr)
	case *ast.BLangIndexBasedAccess:
		return indexBasedAccess(ctx, curBB, expr)
	case *ast.BLangListConstructorExpr:
		return listConstructorExpression(ctx, curBB, expr)
	case *ast.BLangTypeConversionExpr:
		return typeConversionExpression(ctx, curBB, expr)
	default:
		panic("unexpected expression type")
	}
}

func typeConversionExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypeConversionExpr) expressionEffect {
	exprEffect := handleExpression(ctx, curBB, expr.Expression)
	curBB = exprEffect.block
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeCast := &TypeCast{}
	typeCast.RhsOp = resultOperand
	typeCast.LhsOp = exprEffect.result
	typeCast.Type = expr.TypeDescriptor.GetDeterminedType()
	curBB.Instructions = append(curBB.Instructions, typeCast)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func listConstructorExpression(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangListConstructorExpr) expressionEffect {
	initValues := make([]*BIROperand, len(expr.Exprs))
	for i, expr := range expr.Exprs {
		exprEffect := handleExpression(ctx, bb, expr)
		bb = exprEffect.block
		initValues[i] = exprEffect.result
	}
	sizeOperand := ctx.addTempVar(&semtypes.INT)
	constantLoad := &ConstantLoad{}
	constantLoad.Value = int64(len(expr.Exprs))
	constantLoad.LhsOp = sizeOperand
	bb.Instructions = append(bb.Instructions, constantLoad)

	resultOperand := ctx.addTempVar(&semtypes.LIST)
	newArray := &NewArray{}
	newArray.AtomicType = expr.AtomicType
	newArray.LhsOp = resultOperand
	newArray.SizeOp = sizeOperand
	newArray.Values = initValues
	bb.Instructions = append(bb.Instructions, newArray)
	return expressionEffect{
		result: resultOperand,
		block:  bb,
	}
}

func indexBasedAccess(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangIndexBasedAccess) expressionEffect {
	// Assignment is handled in assignmentStatement to this is always a load
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	fieldAccess := &FieldAccess{}
	fieldAccess.Kind = INSTRUCTION_KIND_ARRAY_LOAD
	fieldAccess.LhsOp = resultOperand
	indexEffect := handleExpression(ctx, bb, expr.IndexExpr)
	fieldAccess.KeyOp = indexEffect.result
	containerRefEffect := handleExpression(ctx, indexEffect.block, expr.Expr)
	fieldAccess.RhsOp = containerRefEffect.result
	bb.Instructions = append(bb.Instructions, fieldAccess)
	return expressionEffect{
		result: resultOperand,
		block:  bb,
	}
}

func groupExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangGroupExpr) expressionEffect {
	return handleExpression(ctx, curBB, expr.Expression)
}

func wildcardBindingPattern(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangWildCardBindingPattern) expressionEffect {
	return expressionEffect{
		result: ctx.addTempVar(nil),
		block:  curBB,
	}
}

func unaryExpression(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangUnaryExpr) expressionEffect {
	var kind InstructionKind
	switch expr.Operator {
	case model.OperatorKind_NOT:
		kind = INSTRUCTION_KIND_NOT
	case model.OperatorKind_SUB:
		kind = INSTRUCTION_KIND_NEGATE
	default:
		panic("unexpected unary operator kind")
	}
	opEffect := handleExpression(ctx, bb, expr.Expr)

	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	unaryOp := &UnaryOp{}
	unaryOp.Kind = kind
	unaryOp.LhsOp = resultOperand
	curBB := opEffect.block
	unaryOp.RhsOp = opEffect.result
	curBB.Instructions = append(curBB.Instructions, unaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func invocation(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangInvocation) expressionEffect {
	curBB := bb
	var args []BIROperand
	for _, arg := range expr.ArgExprs {
		argEffect := handleExpression(ctx, curBB, arg)
		curBB = argEffect.block
		args = append(args, *argEffect.result)
	}
	thenBB := ctx.addBB()
	// TODO: deal with type
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	call := &Call{}
	call.Kind = INSTRUCTION_KIND_CALL
	call.Args = args
	call.Name = model.Name(expr.GetName().GetValue())
	call.ThenBB = thenBB
	call.LhsOp = resultOperand

	// Package qualified call - look up package ID from import alias
	if expr.PkgAlias != nil && expr.PkgAlias.Value != "" {
		// Qualified call - look up package ID from import alias
		if pkgID, found := ctx.birCx.importAliasMap[expr.PkgAlias.Value]; found {
			call.CalleePkg = pkgID
		} else {
			panic("unexpected")
		}
	} else {
		// Unqualified call (no PkgAlias) - assume same-module call and use current package
		if ctx.birCx.packageID != nil {
			call.CalleePkg = ctx.birCx.packageID
		}
	}
	orgName := call.CalleePkg.OrgName.Value()
	pkgName := call.CalleePkg.PkgName.Value()
	moduleKey := orgName + "/" + pkgName
	call.FunctionLookupKey = moduleKey + ":" + call.Name.Value()
	curBB.Terminator = call
	return expressionEffect{
		result: resultOperand,
		block:  thenBB,
	}
}

func literal(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangLiteral) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	constantLoad := &ConstantLoad{}
	constantLoad.Value = expr.Value
	constantLoad.LhsOp = resultOperand
	curBB.Instructions = append(curBB.Instructions, constantLoad)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func binaryExpressionInner(ctx *stmtContext, curBB *BIRBasicBlock, opKind model.OperatorKind, lhsExpr, rhsExpr ast.BLangExpression, resultType semtypes.SemType) expressionEffect {
	var kind InstructionKind
	switch opKind {
	case model.OperatorKind_ADD:
		kind = INSTRUCTION_KIND_ADD
	case model.OperatorKind_SUB:
		kind = INSTRUCTION_KIND_SUB
	case model.OperatorKind_MUL:
		kind = INSTRUCTION_KIND_MUL
	case model.OperatorKind_DIV:
		kind = INSTRUCTION_KIND_DIV
	case model.OperatorKind_MOD:
		kind = INSTRUCTION_KIND_MOD
	case model.OperatorKind_AND:
		kind = INSTRUCTION_KIND_AND
	case model.OperatorKind_OR:
		kind = INSTRUCTION_KIND_OR
	case model.OperatorKind_EQUAL:
		kind = INSTRUCTION_KIND_EQUAL
	case model.OperatorKind_NOT_EQUAL:
		kind = INSTRUCTION_KIND_NOT_EQUAL
	case model.OperatorKind_GREATER_THAN:
		kind = INSTRUCTION_KIND_GREATER_THAN
	case model.OperatorKind_GREATER_EQUAL:
		kind = INSTRUCTION_KIND_GREATER_EQUAL
	case model.OperatorKind_LESS_THAN:
		kind = INSTRUCTION_KIND_LESS_THAN
	case model.OperatorKind_LESS_EQUAL:
		kind = INSTRUCTION_KIND_LESS_EQUAL
	case model.OperatorKind_REF_EQUAL:
		kind = INSTRUCTION_KIND_REF_EQUAL
	case model.OperatorKind_REF_NOT_EQUAL:
		kind = INSTRUCTION_KIND_REF_NOT_EQUAL
	case model.OperatorKind_BITWISE_AND:
		kind = INSTRUCTION_KIND_BITWISE_AND
	case model.OperatorKind_BITWISE_OR:
		kind = INSTRUCTION_KIND_BITWISE_OR
	case model.OperatorKind_BITWISE_XOR:
		kind = INSTRUCTION_KIND_BITWISE_XOR
	default:
		panic("unexpected binary operator kind")
	}
	resultOperand := ctx.addTempVar(resultType)
	binaryOp := &BinaryOp{}
	binaryOp.Kind = kind
	binaryOp.LhsOp = resultOperand
	op1Effect := handleExpression(ctx, curBB, lhsExpr)
	curBB = op1Effect.block
	op2Effect := handleExpression(ctx, curBB, rhsExpr)
	curBB = op2Effect.block
	binaryOp.RhsOp1 = *op1Effect.result
	binaryOp.RhsOp2 = *op2Effect.result
	curBB.Instructions = append(curBB.Instructions, binaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func binaryExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	return binaryExpressionInner(ctx, curBB, expr.OpKind, expr.LhsExpr, expr.RhsExpr, expr.GetDeterminedType())
}

func simpleVariableReference(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangSimpleVarRef) expressionEffect {
	varName := expr.VariableName.GetValue()
	symRef := ctx.birCx.CompilerContext.RefSymbol(expr.Symbol())

	// Try local variable lookup first
	if operand, ok := ctx.varMap[symRef]; ok {
		return expressionEffect{
			result: operand,
			block:  curBB,
		}
	}

	// Try constant lookup
	// FIXME: this is a hack until we have package level variable initialization
	if constant, ok := ctx.birCx.constantMap[symRef]; ok {
		resultOperand := ctx.addTempVar(constant.Type)
		constantLoad := &ConstantLoad{}
		constantLoad.Value = constant.ConstValue.Value
		constantLoad.LhsOp = resultOperand
		curBB.Instructions = append(curBB.Instructions, constantLoad)
		return expressionEffect{
			result: resultOperand,
			block:  curBB,
		}
	}

	panic(fmt.Sprintf("variable %s not found (SymbolRef: Pkg=%v Index=%d SpaceIndex=%d)",
		varName, symRef.Package, symRef.Index, symRef.SpaceIndex))
}

func appendIfNotNil[T any](slice []T, item *T) []T {
	if item != nil {
		slice = append(slice, *item)
	}
	return slice
}
