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
	"fmt"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

// Since BLangNodeVisitor is anyway deprecated in jBallerina, we'll try to do this more cleanly
// TODO: may be we should have this in a separate package and keep BIR package clean (only definitions)

type Context struct {
	CompilerContext *context.CompilerContext
	importAliasMap  map[string]*model.PackageID // Maps import alias to package ID
	packageID       *model.PackageID            // Current package ID
	// Ideally I would like to track this by SymbolRef, but in order to map NewExpr to class definition we will need
	// some changes to sementic analysis. Basically we need a way to map the type to the class declaration
	classDefMap map[*semtypes.MappingAtomicType]*BIRClassDef
	birPkg      *BIRPackage
}

type stmtContext struct {
	birCx        *Context
	bbs          []*BIRBasicBlock
	scope        *BIRScope
	nextScopeId  int
	errorEntries []BIRErrorEntry
	loopCtx      *loopContext
	isClosure    bool          // set to true when a captured variable is resolved across a function boundary
	scopeCtx     *scopeContext // current scope (holds localVars, varMap, retVar)
}

type scopeContext struct {
	localVars          []*BIRLocalVariableDcl
	varMap             map[model.SymbolRef]*BIROperand
	retVar             *BIROperand
	parent             *scopeContext
	isFunctionBoundary bool // true at the root scope of each function
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

func (cx *stmtContext) pushScope() {
	retVar := cx.scopeCtx.retVar
	baseIndex := cx.scopeDepth() + 1
	cx.scopeCtx = &scopeContext{
		varMap: make(map[model.SymbolRef]*BIROperand),
		retVar: &BIROperand{
			VariableDcl: retVar.VariableDcl,
			Address:     absoluteAddress(baseIndex, retVar.Address.FrameIndex),
		},
		parent: cx.scopeCtx,
	}
}

func (cx *stmtContext) popScope() {
	if cx.scopeCtx.parent == nil {
		panic("no enclosing scope")
	}
	cx.scopeCtx = cx.scopeCtx.parent
}

func (cx *stmtContext) addLocalVarInner(name model.Name, ty semtypes.SemType) *BIROperand {
	varDcl := &BIRLocalVariableDcl{}
	varDcl.Name = name
	varDcl.Type = ty
	sc := cx.scopeCtx
	sc.localVars = append(sc.localVars, varDcl)
	return &BIROperand{VariableDcl: varDcl, Address: relativeAddress(len(sc.localVars) - 1)}
}

func (cx *stmtContext) addTempVar(ty semtypes.SemType) *BIROperand {
	return cx.addLocalVarInner(model.Name(fmt.Sprintf("%%%d", len(cx.scopeCtx.localVars))), ty)
}

func (cx *stmtContext) addLocalVar(name model.Name, ty semtypes.SemType, symbol model.SymbolRef) *BIROperand {
	operand := cx.addLocalVarInner(name, ty)
	cx.scopeCtx.varMap[symbol] = operand
	return operand
}

// lookupVariable looks up a variable by symbol, checking the local scope first,
// then walking the parent scope chain (crossing function boundaries for closures).
// Returns the operand, whether it crossed a function boundary, and whether it was found.
func (cx *stmtContext) lookupVariable(symRef model.SymbolRef) (*BIROperand, bool, bool) {
	if operand, ok := cx.scopeCtx.varMap[symRef]; ok {
		return operand, false, true
	}
	levelsUp := 1
	crossedFunction := cx.scopeCtx.isFunctionBoundary
	parent := cx.scopeCtx.parent
	for parent != nil {
		if outerOp, ok := parent.varMap[symRef]; ok {
			baseIndex := levelsUp
			if outerOp.Address.Mode == AddressingModeAbsolute {
				baseIndex = levelsUp + outerOp.Address.BaseIndex
			}
			return &BIROperand{
				VariableDcl: outerOp.VariableDcl,
				Address:     absoluteAddress(baseIndex, outerOp.Address.FrameIndex),
			}, crossedFunction, true
		}
		crossedFunction = crossedFunction || parent.isFunctionBoundary
		levelsUp++
		parent = parent.parent
	}
	return nil, false, false
}

// scopeDepth returns the number of block scopes between the current scope and the function root.
func (cx *stmtContext) scopeDepth() int {
	depth := 0
	scope := cx.scopeCtx
	for !scope.isFunctionBoundary {
		depth++
		scope = scope.parent
	}
	return depth
}

func (cx *stmtContext) addBB() *BIRBasicBlock {
	index := len(cx.bbs)
	bb := BB(index)
	cx.bbs = append(cx.bbs, &bb)
	return &bb
}

func buildFunctionLookupKeyFromSymbol(ctx *Context, symRef model.SymbolRef) string {
	return symRef.Package.Organization + "/" + symRef.Package.Package + ":" + ctx.CompilerContext.GetSymbol(symRef).Name()
}

func GenBir(ctx *context.CompilerContext, ast *ast.BLangPackage) *BIRPackage {
	birPkg := &BIRPackage{}
	birPkg.PackageID = ast.PackageID
	genCtx := &Context{
		CompilerContext: ctx,
		importAliasMap:  make(map[string]*model.PackageID),
		packageID:       ast.PackageID,
		classDefMap:     make(map[*semtypes.MappingAtomicType]*BIRClassDef),
		birPkg:          birPkg,
	}
	birPkg.GlobalVars = make(map[model.SymbolRef]BIRGlobalVariableDcl)
	processImports(ctx, genCtx, ast.Imports, birPkg)
	for _, globalVar := range ast.GlobalVars {
		symRef := globalVar.Symbol()
		addGlobalVar(birPkg, symRef, TransformGlobalVariableDcl(genCtx, &globalVar))
	}
	for _, constant := range ast.Constants {
		symRef := constant.Symbol()
		addGlobalVar(birPkg, symRef, transformConstantAsGlobal(genCtx, &constant))
	}
	if ast.InitFunction != nil {
		birPkg.InitFunction = TransformFunction(genCtx, ast.InitFunction)
	}
	for i := range ast.ClassDefinitions {
		transformClassDefinition(genCtx, &ast.ClassDefinitions[i], birPkg)
	}
	for _, function := range ast.Functions {
		if function.FlagSet.Contains(model.Flag_NATIVE) {
			continue
		}
		birFunc := TransformFunction(genCtx, &function)
		birPkg.Functions = append(birPkg.Functions, *birFunc)
		if birFunc.Name.Value() == "main" {
			birPkg.MainFunction = birFunc
		}
	}
	birPkg.TypeEnv = ctx.GetTypeEnv()
	return birPkg
}

func processImports(compilerCtx *context.CompilerContext, genCtx *Context, imports []ast.BLangImportPackage, birPkg *BIRPackage) {
	for _, importPkg := range imports {
		if importPkg.Alias != nil && importPkg.Alias.Value != "" {
			var orgName model.Name
			if importPkg.OrgName != nil && importPkg.OrgName.Value != "" {
				orgName = model.Name(importPkg.OrgName.Value)
			} else if genCtx.packageID != nil && genCtx.packageID.OrgName != nil {
				orgName = *genCtx.packageID.OrgName
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
			Name:    &pkgName,
			Version: &version,
		},
	}
}

func addGlobalVar(birPkg *BIRPackage, symRef model.SymbolRef, dcl BIRGlobalVariableDcl) {
	birPkg.GlobalVars[symRef] = dcl
}

func flagSetToInt64(flags common.Set[model.Flag]) int64 {
	var result int64
	for f := range flags.Values() {
		result |= 1 << int64(f)
	}
	return result
}

func TransformGlobalVariableDcl(ctx *Context, ast *ast.BLangSimpleVariable) BIRGlobalVariableDcl {
	name := model.Name(ast.GetName().GetValue())
	dcl := BIRGlobalVariableDcl{}
	dcl.Pos = ast.GetPosition()
	dcl.Name = name
	dcl.PkgId = ctx.packageID
	dcl.Type = ctx.CompilerContext.SymbolType(ast.Symbol())
	dcl.Flags = flagSetToInt64(ast.GetFlags())
	dcl.Origin = model.SymbolOrigin_SOURCE
	return dcl
}

func transformConstantAsGlobal(ctx *Context, c *ast.BLangConstant) BIRGlobalVariableDcl {
	name := model.Name(c.GetName().GetValue())
	dcl := BIRGlobalVariableDcl{}
	dcl.Pos = c.GetPosition()
	dcl.Name = name
	dcl.PkgId = ctx.packageID
	dcl.Type = ctx.CompilerContext.SymbolType(c.Symbol())
	dcl.Flags = flagSetToInt64(c.GetFlags())
	dcl.Origin = model.SymbolOrigin_SOURCE
	return dcl
}

func TransformFunction(ctx *Context, astFunc *ast.BLangFunction) *BIRFunction {
	stmtCx := &stmtContext{birCx: ctx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), isFunctionBoundary: true}}
	return transformFunctionInner(stmtCx, astFunc, nil)
}

func transformFunctionInner(stmtCx *stmtContext, astFunc *ast.BLangFunction, selfSymbolRef *model.SymbolRef) *BIRFunction {
	symRef := astFunc.Symbol()
	funcName := model.Name(astFunc.GetName().GetValue())
	birFunc := &BIRFunction{}
	birFunc.Pos = astFunc.GetPosition()
	birFunc.Name = funcName
	birFunc.OriginalName = funcName
	birFunc.Flags = flagSetToInt64(astFunc.GetFlags())
	ctx := stmtCx.birCx
	birFunc.FunctionLookupKey = buildFunctionLookupKeyFromSymbol(ctx, symRef)
	common.Assert(astFunc.Receiver == nil)
	funcSym := ctx.CompilerContext.GetSymbol(astFunc.Symbol()).(model.FunctionSymbol)
	stmtCx.scopeCtx.retVar = stmtCx.addLocalVarInner(model.Name("%0"), funcSym.Signature().ReturnType)
	if selfSymbolRef != nil {
		stmtCx.addLocalVar(model.Name("self"), ctx.CompilerContext.SymbolType(*selfSymbolRef), *selfSymbolRef)
	}
	requiredParams := make([]BIRParameter, len(astFunc.RequiredParams))
	for i, param := range astFunc.RequiredParams {
		stmtCx.addLocalVar(model.Name(param.GetName().GetValue()), ctx.CompilerContext.SymbolType(param.Symbol()), param.Symbol())
		requiredParams[i] = BIRParameter{
			Name:  model.Name(param.GetName().GetValue()),
			Flags: flagSetToInt64(param.GetFlags()),
		}
	}
	if astFunc.RestParam != nil {
		restParam := astFunc.RestParam.(*ast.BLangSimpleVariable)
		ty := ctx.CompilerContext.SymbolType(restParam.Symbol())
		stmtCx.addLocalVar(model.Name(restParam.GetName().GetValue()), ty, restParam.Symbol())
		birFunc.RestParams = &BIRParameter{Name: model.Name(restParam.GetName().GetValue())}
	}
	birFunc.RequiredParams = requiredParams
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
	for _, varPtr := range stmtCx.scopeCtx.localVars {
		birFunc.LocalVars = append(birFunc.LocalVars, *varPtr)
	}
	birFunc.ErrorTable = stmtCx.errorEntries
	birFunc.ReturnVariable = stmtCx.scopeCtx.retVar.VariableDcl.(*BIRLocalVariableDcl)
	return birFunc
}

func handleBlockFunctionBody(ctx *stmtContext, ast *ast.BLangBlockFunctionBody) {
	curBB := ctx.addBB()
	for _, stmt := range ast.Stmts {
		effect := handleStatement(ctx, curBB, stmt)
		curBB = effect.block
	}
	// Add implicit return
	if curBB != nil {
		curBB.Terminator = NewReturn(ast.GetPosition())
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
	case *ast.BLangPanic:
		return panicStatement(ctx, curBB, stmt)
	case *ast.BLangMatchStatement:
		return matchStatement(ctx, curBB, stmt)
	default:
		panic("unexpected statement type")
	}
}

func compoundAssignment(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangCompoundAssignment) statementEffect {
	// First do the operation
	ref := stmt.VarRef.(ast.BLangExpression)
	valueEffect := binaryExpressionInner(ctx, curBB, stmt.OpKind, ref, stmt.Expr, stmt.Expr.GetDeterminedType(), stmt.GetPosition())
	// Then do the assignment
	return assignmentStatementInner(ctx, valueEffect.block, ref, valueEffect, stmt.GetPosition())
}

func continueStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangContinue) statementEffect {
	curBB.Instructions = append(curBB.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: stmt.GetPosition()}}})
	onContinueBB := ctx.loopCtx.onContinueBB
	curBB.Terminator = NewGoto(onContinueBB, stmt.GetPosition())
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func breakStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangBreak) statementEffect {
	curBB.Instructions = append(curBB.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: stmt.GetPosition()}}})
	onBreakBB := ctx.loopCtx.onBreakBB
	curBB.Terminator = NewGoto(onBreakBB, stmt.GetPosition())
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func whileStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangWhile) statementEffect {
	loopHead := ctx.addBB()
	// jump to loop head
	bb.Terminator = NewGoto(loopHead, stmt.GetPosition())
	condEffect := handleExpression(ctx, loopHead, stmt.Expr)

	loopBody := ctx.addBB()
	loopEnd := ctx.addBB()
	// conditionally jump to loop body
	condEffect.block.Terminator = NewBranch(condEffect.result, loopBody, loopEnd, stmt.GetPosition())

	// Push scope frame for loop body — each iteration gets its own frame
	pushScope := &PushScopeFrame{}
	pushScope.Pos = stmt.GetPosition()
	loopBody.Instructions = append(loopBody.Instructions, pushScope)
	ctx.pushScope()

	ctx.addLoopCtx(loopEnd, loopHead)
	bodyEffect := blockStatement(ctx, loopBody, &stmt.Body)

	// Fill in scope frame local vars now that body has been processed
	pushScope.NumLocals = len(ctx.scopeCtx.localVars)

	// This could happen if the while block always ends return, break or continue
	if bodyEffect.block != nil {
		bodyEffect.block.Instructions = append(bodyEffect.block.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: stmt.GetPosition()}}})
		bodyEffect.block.Terminator = NewGoto(loopHead, stmt.GetPosition())
	}

	ctx.popLoopCtx()
	ctx.popScope()
	return statementEffect{
		block: loopEnd,
	}
}

func assignmentStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangAssignment) statementEffect {
	valueEffect := handleExpression(ctx, bb, stmt.Expr)
	return assignmentStatementInner(ctx, valueEffect.block, stmt.VarRef, valueEffect, stmt.GetPosition())
}

func assignmentStatementInner(ctx *stmtContext, bb *BIRBasicBlock, ref ast.BLangExpression, valueEffect expressionEffect, pos ast.Location) statementEffect {
	switch varRef := ref.(type) {
	case *ast.BLangIndexBasedAccess:
		return assignToMemberStatement(ctx, bb, varRef, valueEffect, pos)
	case *ast.BLangWildCardBindingPattern:
		return assignToWildcardBindingPattern(ctx, bb, varRef, valueEffect, pos)
	case *ast.BLangSimpleVarRef:
		return assignToSimpleVariable(ctx, bb, varRef, valueEffect, pos)
	default:
		panic("unexpected variable reference type")
	}
}

func assignToWildcardBindingPattern(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangWildCardBindingPattern, valueEffect expressionEffect, pos ast.Location) statementEffect {
	refEffect := wildcardBindingPattern(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := NewMove(valueEffect.result, refEffect.result, pos)
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToSimpleVariable(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangSimpleVarRef, valueEffect expressionEffect, pos ast.Location) statementEffect {
	refEffect := simpleVariableReference(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := NewMove(valueEffect.result, refEffect.result, pos)
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToMemberStatement(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangIndexBasedAccess, valueEffect expressionEffect, pos ast.Location) statementEffect {
	currBB := valueEffect.block
	containerRefEffect := handleExpression(ctx, currBB, varRef.Expr)
	currBB = containerRefEffect.block
	indexEffect := handleExpression(ctx, currBB, varRef.IndexExpr)
	currBB = indexEffect.block
	containerType := varRef.Expr.GetDeterminedType()
	var fieldAccessKind InstructionKind
	if semtypes.IsSubtypeSimple(containerType, semtypes.LIST) {
		fieldAccessKind = INSTRUCTION_KIND_ARRAY_STORE
	} else if semtypes.IsSubtypeSimple(containerType, semtypes.OBJECT) {
		fieldAccessKind = INSTRUCTION_KIND_OBJECT_STORE
	} else {
		fieldAccessKind = INSTRUCTION_KIND_MAP_STORE
	}
	fieldAccess := NewFieldAccess(fieldAccessKind, containerRefEffect.result, indexEffect.result, valueEffect.result, pos)
	currBB.Instructions = append(currBB.Instructions, fieldAccess)
	return statementEffect{
		block: currBB,
	}
}

func simpleVariableDefinition(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangSimpleVariableDef) statementEffect {
	ty := ctx.birCx.CompilerContext.SymbolType(stmt.Var.Symbol())
	varName := model.Name(stmt.Var.GetName().GetValue())
	if stmt.Var.Expr == nil {
		ctx.addLocalVar(varName, ty, stmt.Var.Symbol())
		// just declare the variable
		return statementEffect{
			block: bb,
		}
	}
	exprResult := handleExpression(ctx, bb, stmt.Var.Expr.(ast.BLangExpression))
	curBB := exprResult.block
	lhsOp := ctx.addLocalVar(varName, ty, stmt.Var.Symbol())
	move := NewMove(exprResult.result, lhsOp, stmt.GetPosition())
	curBB.Instructions = append(curBB.Instructions, move)
	return statementEffect{
		block: curBB,
	}
}

func returnStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangReturn) statementEffect {
	curBB := bb
	pos := stmt.GetPosition()
	if stmt.Expr != nil {
		valueEffect := handleExpression(ctx, curBB, stmt.Expr)
		curBB = valueEffect.block
		mov := NewMove(valueEffect.result, ctx.scopeCtx.retVar, pos)
		curBB.Instructions = append(curBB.Instructions, mov)
	}
	ret := NewReturn(pos)
	curBB.Terminator = ret
	return statementEffect{}
}

func panicStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangPanic) statementEffect {
	errorEffect := handleExpression(ctx, curBB, stmt.Expr)
	curBB = errorEffect.block
	curBB.Terminator = NewPanic(errorEffect.result, stmt.GetPosition())
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
		curBB.Terminator = NewBranch(cond.result, thenBB, elseBB, stmt.GetPosition())

		elseEffect := handleStatement(ctx, elseBB, stmt.ElseStmt)
		finalBB = ctx.addBB()
		if elseEffect.block != nil {
			elseEffect.block.Terminator = NewGoto(finalBB, stmt.GetPosition())
		}
	} else {
		finalBB = ctx.addBB()
		curBB.Terminator = NewBranch(cond.result, thenBB, finalBB, stmt.GetPosition())
	}
	// this could be nil if the control flow moved out of the if (ex: break, continue, return, etc)
	if thenEffect.block != nil {
		thenEffect.block.Terminator = NewGoto(finalBB, stmt.GetPosition())
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

func matchStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangMatchStatement) statementEffect {
	exprEffect := handleExpression(ctx, curBB, stmt.Expr)
	curBB = exprEffect.block
	matchOperand := exprEffect.result
	finalBB := ctx.addBB()

	for _, clause := range stmt.MatchClauses {
		clauseBodyBB := ctx.addBB()

		if isUnconditionalWildcard(&clause) {
			curBB.Terminator = NewGoto(clauseBodyBB, stmt.GetPosition())
			bodyEffect := blockStatement(ctx, clauseBodyBB, &clause.Body)
			if bodyEffect.block != nil {
				bodyEffect.block.Terminator = NewGoto(finalBB, stmt.GetPosition())
			}
			continue
		}

		var condOperand *BIROperand
		for _, pattern := range clause.Patterns {
			switch p := pattern.(type) {
			case *ast.BLangConstPattern:
				patternEffect := handleExpression(ctx, curBB, p.Expr)
				curBB = patternEffect.block
				eqResult := ctx.addTempVar(semtypes.BOOLEAN)
				eqPos := p.Expr.GetPosition()
				binaryOp := NewBinaryOp(INSTRUCTION_KIND_EQUAL, eqResult, matchOperand, patternEffect.result, eqPos)
				curBB.Instructions = append(curBB.Instructions, binaryOp)
				condOperand = orOperands(ctx, curBB, condOperand, eqResult, eqPos)
			case *ast.BLangWildCardMatchPattern:
				// Wildcard in multi-pattern — always matches; but may have guard
				trueOperand := ctx.addTempVar(semtypes.BOOLEAN)
				constLoad := NewConstantLoad(trueOperand, true, p.GetPosition())
				curBB.Instructions = append(curBB.Instructions, constLoad)
				condOperand = orOperands(ctx, curBB, condOperand, trueOperand, p.GetPosition())
			default:
				ctx.birCx.CompilerContext.InternalError("unexpected match pattern type", pattern.GetPosition())
			}
		}

		if clause.Guard != nil {
			guardEffect := handleExpression(ctx, curBB, clause.Guard)
			curBB = guardEffect.block
			condOperand = andOperands(ctx, curBB, condOperand, guardEffect.result, clause.Guard.GetPosition())
		}

		nextCheckBB := ctx.addBB()
		curBB.Terminator = NewBranch(condOperand, clauseBodyBB, nextCheckBB, stmt.GetPosition())

		bodyEffect := blockStatement(ctx, clauseBodyBB, &clause.Body)
		if bodyEffect.block != nil {
			bodyEffect.block.Terminator = NewGoto(finalBB, stmt.GetPosition())
		}

		curBB = nextCheckBB
	}

	if !stmt.IsExhaustive {
		curBB.Terminator = NewGoto(finalBB, stmt.GetPosition())
	}

	return statementEffect{block: finalBB}
}

func isUnconditionalWildcard(clause *ast.BLangMatchClause) bool {
	if clause.Guard != nil {
		return false
	}
	if len(clause.Patterns) != 1 {
		return false
	}
	_, ok := clause.Patterns[0].(*ast.BLangWildCardMatchPattern)
	return ok
}

func orOperands(ctx *stmtContext, bb *BIRBasicBlock, existing *BIROperand, new *BIROperand, pos ast.Location) *BIROperand {
	if existing == nil {
		return new
	}
	result := ctx.addTempVar(semtypes.BOOLEAN)
	binaryOp := NewBinaryOp(INSTRUCTION_KIND_OR, result, existing, new, pos)
	bb.Instructions = append(bb.Instructions, binaryOp)
	return result
}

func andOperands(ctx *stmtContext, bb *BIRBasicBlock, existing *BIROperand, new *BIROperand, pos ast.Location) *BIROperand {
	result := ctx.addTempVar(semtypes.BOOLEAN)
	binaryOp := NewBinaryOp(INSTRUCTION_KIND_AND, result, existing, new, pos)
	bb.Instructions = append(bb.Instructions, binaryOp)
	return result
}

func handleExprFunctionBody(ctx *stmtContext, body *ast.BLangExprFunctionBody) {
	curBB := ctx.addBB()
	effect := handleExpression(ctx, curBB, body.Expr.(ast.BLangExpression))
	curBB = effect.block
	if curBB != nil {
		retAssign := &Move{}
		retAssign.LhsOp = ctx.scopeCtx.retVar
		retAssign.RhsOp = effect.result
		curBB.Instructions = append(curBB.Instructions, retAssign)
		curBB.Terminator = &Return{}
	}
}

func lambdaFunction(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangLambdaFunction) expressionEffect {
	innerCtx := &stmtContext{birCx: ctx.birCx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), parent: ctx.scopeCtx, isFunctionBoundary: true}}
	birFunc := transformFunctionInner(innerCtx, expr.Function, nil)
	ctx.birCx.birPkg.Functions = append(ctx.birCx.birPkg.Functions, *birFunc)
	funcType := expr.GetDeterminedType()
	resultOperand := ctx.addTempVar(funcType)
	fpLoad := &FPLoad{}
	fpLoad.Pos = expr.GetPosition()
	fpLoad.FunctionLookupKey = birFunc.FunctionLookupKey
	fpLoad.Type = funcType
	fpLoad.IsClosure = innerCtx.isClosure
	fpLoad.LhsOp = resultOperand
	curBB.Instructions = append(curBB.Instructions, fpLoad)
	// If the inner function is a closure, this function also needs parent frame
	// access to maintain the frame chain for nested closures
	ctx.isClosure = ctx.isClosure || innerCtx.isClosure
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

type expressionEffect struct {
	result *BIROperand
	block  *BIRBasicBlock
}

// snapshotIfNeeded stores values without storage identity in a temp var before referencing so that modification in one part
// of an expression dont' affect the other.
func snapshotIfNeeded(ctx *stmtContext, effect expressionEffect, pos ast.Location) expressionEffect {
	op := effect.result
	if _, isLocal := op.VariableDcl.(*BIRLocalVariableDcl); isLocal && hasNoStorageIdentity(op.VariableDcl.GetType()) {
		tempOp := ctx.addTempVar(op.VariableDcl.GetType())
		effect.block.Instructions = append(effect.block.Instructions, NewMove(op, tempOp, pos))
		effect.result = tempOp
	}
	return effect
}

func handleExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr ast.BLangExpression) expressionEffect {
	switch expr := expr.(type) {
	case *ast.BLangInvocation:
		return invocation(ctx, curBB, expr)
	case *ast.BLangLiteral:
		return literal(ctx, curBB, expr)
	case *ast.BLangNumericLiteral:
		return literal(ctx, curBB, &expr.BLangLiteral)
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
	case *ast.BLangTypeTestExpr:
		return typeTestExpression(ctx, curBB, expr)
	case *ast.BLangMappingConstructorExpr:
		return mappingConstructorExpression(ctx, curBB, expr)
	case *ast.BLangErrorConstructorExpr:
		return errorConstructorExpression(ctx, curBB, expr)
	case *ast.BLangTrapExpr:
		return trapExpression(ctx, curBB, expr)
	case *ast.BLangNewExpression:
		return newExpression(ctx, curBB, expr)
	case *ast.BLangLambdaFunction:
		return lambdaFunction(ctx, curBB, expr)
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", expr))
	}
}

type mappingField struct {
	key   string
	value ast.BLangExpression
}

func mappingConstructorExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangMappingConstructorExpr) expressionEffect {
	var fields []mappingField
	for _, field := range expr.Fields {
		switch f := field.(type) {
		case *ast.BLangMappingKeyValueField:
			keyName := mappingKeyName(f.Key)
			fields = append(fields, mappingField{key: keyName, value: f.ValueExpr})
		default:
			ctx.birCx.CompilerContext.Unimplemented("non-key-value record field not implemented", expr.GetPosition())
		}
	}
	return mappingConstructorExpressionInner(ctx, curBB, expr.GetDeterminedType(), fields, expr.GetPosition())
}

func mappingKeyName(key *ast.BLangMappingKey) string {
	switch expr := key.Expr.(type) {
	case *ast.BLangLiteral:
		return expr.Value.(string)
	case *ast.BLangSimpleVarRef:
		return expr.VariableName.Value
	default:
		panic(fmt.Sprintf("unexpected mapping key expression type: %T", key.Expr))
	}
}

func mappingConstructorExpressionInner(ctx *stmtContext, curBB *BIRBasicBlock, mapType semtypes.SemType, fields []mappingField, pos ast.Location) expressionEffect {
	var entries []MappingConstructorEntry
	for _, field := range fields {
		keyOperand := ctx.addTempVar(semtypes.STRING)
		keyLoad := NewConstantLoad(keyOperand, field.key, pos)
		curBB.Instructions = append(curBB.Instructions, keyLoad)

		valueEffect := handleExpression(ctx, curBB, field.value)
		curBB = valueEffect.block
		entries = append(entries, &MappingConstructorKeyValueEntry{
			keyOp:   keyOperand,
			valueOp: valueEffect.result,
		})
	}
	resultOperand := ctx.addTempVar(mapType)
	newMap := NewMapConstructor(mapType, resultOperand, entries, pos)
	curBB.Instructions = append(curBB.Instructions, newMap)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func errorConstructorExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangErrorConstructorExpr) expressionEffect {
	// Message is the first positional arg
	msgEffect := handleExpression(ctx, curBB, expr.PositionalArgs[0])
	curBB = msgEffect.block

	// Cause is the optional second positional arg
	var causeOp *BIROperand
	if len(expr.PositionalArgs) > 1 {
		causeEffect := handleExpression(ctx, curBB, expr.PositionalArgs[1])
		curBB = causeEffect.block
		causeOp = causeEffect.result
	}

	// Detail from named args
	var detailOp *BIROperand
	if len(expr.NamedArgs) > 0 {
		var fields []mappingField
		for _, namedArg := range expr.NamedArgs {
			fields = append(fields, mappingField{key: namedArg.Name.Value, value: namedArg.Expr})
		}
		detailEffect := mappingConstructorExpressionInner(ctx, curBB, semtypes.MAPPING, fields, expr.GetPosition())
		curBB = detailEffect.block
		detailOp = detailEffect.result
	}

	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeName := ""
	if expr.ErrorTypeRef != nil {
		typeName = expr.ErrorTypeRef.TypeName.Value
	}
	newError := NewErrorConstructor(expr.GetDeterminedType(), typeName, resultOperand, msgEffect.result, causeOp, detailOp, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, newError)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func typeConversionExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypeConversionExpr) expressionEffect {
	exprEffect := handleExpression(ctx, curBB, expr.Expression)
	curBB = exprEffect.block
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeCast := &TypeCast{}
	typeCast.Pos = expr.GetPosition()
	typeCast.RhsOp = exprEffect.result
	typeCast.LhsOp = resultOperand
	typeCast.Type = expr.TypeDescriptor.GetDeterminedType()
	curBB.Instructions = append(curBB.Instructions, typeCast)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func typeTestExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypeTestExpr) expressionEffect {
	exprEffect := handleExpression(ctx, curBB, expr.Expr)
	curBB = exprEffect.block
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeTest := &TypeTest{}
	typeTest.Pos = expr.GetPosition()
	typeTest.LhsOp = resultOperand
	typeTest.RhsOp = exprEffect.result
	typeTest.Type = expr.Type.Type
	typeTest.IsNegation = expr.IsNegation()
	curBB.Instructions = append(curBB.Instructions, typeTest)
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

	lat := expr.AtomicType
	exprPos := expr.GetPosition()
	for i := len(expr.Exprs); i < lat.Members.FixedLength; i++ {
		ty := lat.MemberAt(i)
		fillerVal := values.DefaultValueForType(ty)
		fillerOperand := ctx.addTempVar(ty)
		fillerLoad := NewConstantLoad(fillerOperand, fillerVal, exprPos)
		bb.Instructions = append(bb.Instructions, fillerLoad)
		initValues = append(initValues, fillerOperand)
	}
	fillerVal := values.DefaultValueForType(semtypes.CellInnerVal(lat.Rest))

	sizeOperand := ctx.addTempVar(semtypes.INT)
	constantLoad := NewConstantLoad(sizeOperand, int64(len(initValues)), exprPos)
	bb.Instructions = append(bb.Instructions, constantLoad)

	resultOperand := ctx.addTempVar(semtypes.LIST)
	newArray := NewArrayConstructor(expr.GetDeterminedType(), resultOperand, sizeOperand, initValues, fillerVal, exprPos)
	bb.Instructions = append(bb.Instructions, newArray)
	return expressionEffect{
		result: resultOperand,
		block:  bb,
	}
}

func indexBasedAccess(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangIndexBasedAccess) expressionEffect {
	// Assignment is handled in assignmentStatement to this is always a load
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	containerType := expr.Expr.GetDeterminedType()
	var fieldAccessKind InstructionKind
	if semtypes.IsSubtypeSimple(containerType, semtypes.LIST) {
		fieldAccessKind = INSTRUCTION_KIND_ARRAY_LOAD
	} else if semtypes.IsSubtypeSimple(containerType, semtypes.OBJECT) {
		fieldAccessKind = INSTRUCTION_KIND_OBJECT_LOAD
	} else {
		fieldAccessKind = INSTRUCTION_KIND_MAP_LOAD
	}
	indexEffect := handleExpression(ctx, bb, expr.IndexExpr)
	containerRefEffect := handleExpression(ctx, indexEffect.block, expr.Expr)
	fieldAccess := NewFieldAccess(fieldAccessKind, resultOperand, indexEffect.result, containerRefEffect.result, expr.GetPosition())
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
	case model.OperatorKind_BITWISE_COMPLEMENT:
		kind = INSTRUCTION_KIND_BITWISE_COMPLEMENT
	default:
		panic("unexpected unary operator kind")
	}
	opEffect := handleExpression(ctx, bb, expr.Expr)

	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	curBB := opEffect.block
	unaryOp := NewUnaryOp(kind, resultOperand, opEffect.result, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, unaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func invocation(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangInvocation) expressionEffect {
	curBB := bb
	var args []BIROperand

	if expr.Expr != nil {
		receiverEffect := handleExpression(ctx, curBB, expr.Expr)
		curBB = receiverEffect.block
		args = append(args, *receiverEffect.result)
	}

	for _, arg := range expr.ArgExprs {
		argEffect := handleExpression(ctx, curBB, arg)
		argEffect = snapshotIfNeeded(ctx, argEffect, expr.GetPosition())
		curBB = argEffect.block
		args = append(args, *argEffect.result)
	}
	thenBB := ctx.addBB()
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	call := NewCall(INSTRUCTION_KIND_CALL, args, model.Name(expr.GetName().GetValue()), thenBB, resultOperand, expr.GetPosition())

	if expr.Expr != nil {
		call.IsVirtual = true
	}

	symRef := expr.Symbol()
	sym := ctx.birCx.CompilerContext.GetSymbol(symRef)
	if sym.Kind() == model.SymbolKindFunction {
		// Regular function call
		call.Kind = INSTRUCTION_KIND_CALL
		call.FunctionLookupKey = buildFunctionLookupKeyFromSymbol(ctx.birCx, symRef)
		if expr.PkgAlias != nil && expr.PkgAlias.Value != "" {
			call.CalleePkg = ctx.birCx.importAliasMap[expr.PkgAlias.Value]
		} else if ctx.birCx.packageID != nil {
			call.CalleePkg = ctx.birCx.packageID
		}
	} else {
		// Function pointer call through a variable
		call.Kind = INSTRUCTION_KIND_FP_CALL
		unnarrowedRef := ctx.birCx.CompilerContext.UnnarrowedSymbol(symRef)
		if op, crossedFunction, ok := ctx.lookupVariable(unnarrowedRef); ok {
			call.FpOperand = op
			ctx.isClosure = ctx.isClosure || crossedFunction
		}
	}
	curBB.Terminator = call
	return expressionEffect{
		result: resultOperand,
		block:  thenBB,
	}
}

func literal(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangLiteral) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	constantLoad := NewConstantLoad(resultOperand, expr.Value, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, constantLoad)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func binaryExpressionInner(ctx *stmtContext, curBB *BIRBasicBlock, opKind model.OperatorKind, lhsExpr, rhsExpr ast.BLangExpression, resultType semtypes.SemType, pos ast.Location) expressionEffect {
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
	case model.OperatorKind_BITWISE_LEFT_SHIFT:
		kind = INSTRUCTION_KIND_BITWISE_LEFT_SHIFT
	case model.OperatorKind_BITWISE_RIGHT_SHIFT:
		kind = INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT
	case model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
		kind = INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT
	default:
		panic("unexpected binary operator kind")
	}
	resultOperand := ctx.addTempVar(resultType)
	op1Effect := handleExpression(ctx, curBB, lhsExpr)
	op1Effect = snapshotIfNeeded(ctx, op1Effect, pos)
	curBB = op1Effect.block
	op2Effect := handleExpression(ctx, curBB, rhsExpr)
	op2Effect = snapshotIfNeeded(ctx, op2Effect, pos)
	curBB = op2Effect.block
	binaryOp := NewBinaryOp(kind, resultOperand, op1Effect.result, op2Effect.result, pos)
	curBB.Instructions = append(curBB.Instructions, binaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func binaryExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	switch expr.OpKind {
	case model.OperatorKind_AND:
		return logicalAndExpression(ctx, curBB, expr)
	case model.OperatorKind_OR:
		return logicalOrExpression(ctx, curBB, expr)
	default:
		return binaryExpressionInner(ctx, curBB, expr.OpKind, expr.LhsExpr, expr.RhsExpr, expr.GetDeterminedType(), expr.GetPosition())
	}
}

func logicalAndExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	lhsEffect := handleExpression(ctx, curBB, expr.LhsExpr)
	curBB = lhsEffect.block

	mov := NewMove(lhsEffect.result, resultOperand, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, mov)

	evalRhsBB := ctx.addBB()
	doneBB := ctx.addBB()

	curBB.Terminator = NewBranch(lhsEffect.result, evalRhsBB, doneBB, expr.GetPosition())

	rhsEffect := handleExpression(ctx, evalRhsBB, expr.RhsExpr)
	rhsBB := rhsEffect.block

	rhsMov := NewMove(rhsEffect.result, resultOperand, expr.GetPosition())

	rhsBB.Instructions = append(rhsBB.Instructions, rhsMov)
	rhsBB.Terminator = NewGoto(doneBB, expr.GetPosition())

	return expressionEffect{
		result: resultOperand,
		block:  doneBB,
	}
}

func logicalOrExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	lhsEffect := handleExpression(ctx, curBB, expr.LhsExpr)
	curBB = lhsEffect.block

	mov := NewMove(lhsEffect.result, resultOperand, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, mov)

	evalRhsBB := ctx.addBB()
	doneBB := ctx.addBB()

	curBB.Terminator = NewBranch(lhsEffect.result, doneBB, evalRhsBB, expr.GetPosition())

	rhsEffect := handleExpression(ctx, evalRhsBB, expr.RhsExpr)
	rhsBB := rhsEffect.block

	rhsMov := NewMove(rhsEffect.result, resultOperand, expr.GetPosition())
	rhsBB.Instructions = append(rhsBB.Instructions, rhsMov)
	rhsBB.Terminator = NewGoto(doneBB, expr.GetPosition())

	return expressionEffect{
		result: resultOperand,
		block:  doneBB,
	}
}

func simpleVariableReference(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangSimpleVarRef) expressionEffect {
	varName := expr.VariableName.GetValue()
	symRef := ctx.birCx.CompilerContext.UnnarrowedSymbol(expr.Symbol())

	if operand, crossedFunction, ok := ctx.lookupVariable(symRef); ok {
		ctx.isClosure = ctx.isClosure || crossedFunction
		return expressionEffect{
			result: operand,
			block:  curBB,
		}
	}

	// Try function lookup
	sym := ctx.birCx.CompilerContext.GetSymbol(symRef)
	if sym.Kind() == model.SymbolKindFunction {
		funcType := ctx.birCx.CompilerContext.SymbolType(symRef)
		lookupKey := buildFunctionLookupKeyFromSymbol(ctx.birCx, symRef)
		resultOperand := ctx.addTempVar(funcType)
		fpLoad := NewFPLoad(lookupKey, funcType, resultOperand, expr.GetPosition())
		curBB.Instructions = append(curBB.Instructions, fpLoad)
		return expressionEffect{
			result: resultOperand,
			block:  curBB,
		}
	}

	// Global variable reference
	var pkgId *model.PackageID
	if expr.PkgAlias != nil && expr.PkgAlias.Value != "" {
		pkgId = ctx.birCx.importAliasMap[expr.PkgAlias.Value]
	} else {
		pkgId = ctx.birCx.packageID
	}
	gv := &BIRGlobalVariableDcl{}
	gv.Name = model.Name(varName)
	gv.PkgId = pkgId
	return expressionEffect{
		result: &BIROperand{VariableDcl: gv, SymRef: &symRef},
		block:  curBB,
	}
}

func trapExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTrapExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	trapStartBB := ctx.addBB()
	curBB.Terminator = NewGoto(trapStartBB, expr.GetPosition())

	innerEffect := handleExpression(ctx, trapStartBB, expr.Expr)
	trapEndBB := innerEffect.block

	mov := NewMove(innerEffect.result, resultOperand, expr.GetPosition())
	trapEndBB.Instructions = append(trapEndBB.Instructions, mov)

	afterTrapBB := ctx.addBB()
	trapEndBB.Terminator = NewGoto(afterTrapBB, expr.GetPosition())

	ctx.errorEntries = append(ctx.errorEntries, BIRErrorEntry{
		Start:   trapStartBB,
		End:     trapEndBB,
		Target:  afterTrapBB,
		ErrorOp: resultOperand,
	})

	return expressionEffect{
		result: resultOperand,
		block:  afterTrapBB,
	}
}

func transformClassDefinition(ctx *Context, class *ast.BLangClassDefinition, birPkg *BIRPackage) {
	className := model.Name(class.GetName().GetValue())
	classScope := class.Scope()
	selfRef, ok := classScope.GetSymbol("self")
	if !ok {
		ctx.CompilerContext.InternalError("self symbol not found in class scope", class.GetPosition())
	}

	birClassDef := &BIRClassDef{
		Name:   className,
		VTable: make(map[string]*BIRFunction),
	}

	for _, field := range class.Fields {
		birClassDef.Fields = append(birClassDef.Fields, ObjectField{
			Name: field.GetName().GetValue(),
			Ty:   ctx.CompilerContext.SymbolType(field.Symbol()),
		})
	}

	initFunc := transformFunctionInner(&stmtContext{birCx: ctx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), isFunctionBoundary: true}}, class.InitFunction, &selfRef)
	birClassDef.VTable["init"] = initFunc

	for methodName, method := range class.Methods {
		birClassDef.VTable[methodName] = transformFunctionInner(&stmtContext{birCx: ctx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), isFunctionBoundary: true}}, method, &selfRef)
	}

	semCtx := semtypes.ContextFrom(ctx.CompilerContext.GetTypeEnv())
	classType := class.Definition.GetSemType(ctx.CompilerContext.GetTypeEnv())
	atomicType := semtypes.ToObjectAtomicType(semCtx, classType)
	if atomicType != nil {
		ctx.classDefMap[atomicType] = birClassDef
	} else {
		ctx.CompilerContext.InternalError("failed to determine object atomic type", class.GetPosition())
	}

	birPkg.ClassDefs = append(birPkg.ClassDefs, *birClassDef)
}

func newExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangNewExpression) expressionEffect {
	classDef := ctx.birCx.classDefMap[expr.AtomicType]
	if classDef == nil {
		ctx.birCx.CompilerContext.InternalError("failed to find the class definition", expr.GetPosition())
		return expressionEffect{}
	}

	object := ctx.addTempVar(expr.GetDeterminedType())
	newObj := NewObjectConstructor(classDef, object, expr.GetPosition())
	curBB.Instructions = append(curBB.Instructions, newObj)

	var args []BIROperand
	args = append(args, *object)
	for _, arg := range expr.ArgsExprs {
		argEffect := handleExpression(ctx, curBB, arg)
		curBB = argEffect.block
		args = append(args, *argEffect.result)
	}

	initFunc := classDef.VTable["init"]
	initResult := ctx.addTempVar(initFunc.ReturnVariable.Type)
	initDoneBB := ctx.addBB()
	call := NewCall(INSTRUCTION_KIND_CALL, args, initFunc.Name, initDoneBB, initResult, expr.GetPosition())
	call.IsVirtual = true
	call.CachedBIRFunc = initFunc
	curBB.Terminator = call

	result := ctx.addTempVar(expr.DeterminedType)
	isInitResultNil := ctx.addTempVar(semtypes.BOOLEAN)
	nilCheck := &TypeTest{}
	nilCheck.Pos = expr.GetPosition()
	nilCheck.LhsOp = isInitResultNil
	nilCheck.RhsOp = initResult
	nilCheck.Type = semtypes.NIL
	initDoneBB.Instructions = append(initDoneBB.Instructions, nilCheck)

	assignObjectBB := ctx.addBB()
	assignErrorBB := ctx.addBB()
	thenBB := ctx.addBB()
	initDoneBB.Terminator = NewBranch(isInitResultNil, assignObjectBB, assignErrorBB, expr.GetPosition())

	assignObjectBB.Instructions = append(assignObjectBB.Instructions, NewMove(object, result, expr.GetPosition()))
	assignObjectBB.Terminator = NewGoto(thenBB, expr.GetPosition())

	assignErrorBB.Instructions = append(assignErrorBB.Instructions, NewMove(initResult, result, expr.GetPosition()))
	assignErrorBB.Terminator = NewGoto(thenBB, expr.GetPosition())

	return expressionEffect{
		result: result,
		block:  thenBB,
	}
}

func appendIfNotNil[T any](slice []T, item *T) []T {
	if item != nil {
		slice = append(slice, *item)
	}
	return slice
}

func hasNoStorageIdentity(ty semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(ty, semtypes.SIMPLE_BASIC)
}
