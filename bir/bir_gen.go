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
	"sort"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
)

// birLoc converts a diagnostics.Location (byte offsets) to a bir.Location (line/column)
// using a DiagnosticEnv to resolve byte offsets.
func birLoc(de *diagnostics.DiagnosticEnv, pos diagnostics.Location) Location {
	return NewLocation(de.FileName(pos), de.StartLine(pos), de.EndLine(pos), de.StartColumn(pos), de.EndColumn(pos))
}

// Since BLangNodeVisitor is anyway deprecated in jBallerina, we'll try to do this more cleanly
// TODO: may be we should have this in a separate package and keep BIR package clean (only definitions)

type Context struct {
	CompilerContext *context.CompilerContext
	importAliasMap  map[string]*model.PackageID // Maps import alias to package ID
	packageID       *model.PackageID            // Current package ID
	birPkg          *BIRPackage
	typeCtx         semtypes.Context
	stringMapTy     semtypes.SemType // Memoized map<string> type
}

func (c *Context) TypeContext() semtypes.Context {
	return c.typeCtx
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
	// activeLockKey is set while emitting BIR for the body of a `lock`
	// statement. Nested locks are rejected by the lock analyzer, so at most
	// one is active at any point during BIR-gen.
	activeLockKey *string
}

// emitLockEndBeforeAbruptExit, when called inside a lock body, closes the
// current BB with a LockEnd terminator and returns a fresh BB whose execution
// resumes the abrupt-exit terminator (Return/Break/Continue). When not inside
// a lock body, returns curBB unchanged.
func emitLockEndBeforeAbruptExit(ctx *stmtContext, curBB *BIRBasicBlock, pos Location) *BIRBasicBlock {
	if ctx.activeLockKey == nil {
		return curBB
	}
	newBB := ctx.addBB()
	curBB.Terminator = NewLockEnd(*ctx.activeLockKey, newBB, pos)
	return newBB
}

func (c *Context) stringMapType() semtypes.SemType {
	if c.stringMapTy == nil {
		md := semtypes.NewMappingDefinition()
		c.stringMapTy = md.DefineMappingTypeWrapped(c.CompilerContext.GetTypeEnv(), nil, semtypes.STRING)
	}
	return c.stringMapTy
}

func (cx *stmtContext) loc(pos diagnostics.Location) Location {
	return birLoc(cx.birCx.CompilerContext.DiagnosticEnv(), pos)
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
	return &BIROperand{VariableDcl: varDcl, Address: RelativeAddress(len(sc.localVars) - 1)}
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

func buildLookupKey(pkg model.PackageIdentifier, qualifiedName string) string {
	return pkg.Organization + "/" + pkg.Package + ":" + qualifiedName
}

func buildFunctionLookupKeyFromSymbol(ctx *Context, symRef model.SymbolRef) string {
	sym := ctx.CompilerContext.GetSymbol(symRef)
	if mono, ok := sym.(model.MonomorphicFunctionSymbol); ok {
		// For monomorphic functions (ex: dependently typed functions), in runtime we dispatch to a single
		// polymorphic function
		origRef := mono.PolymorphicSymbol()
		return buildLookupKey(origRef.Package, ctx.CompilerContext.GetSymbol(origRef).Name())
	}
	return buildLookupKey(symRef.Package, sym.Name())
}

func buildMethodLookupKeyFromSymbol(ctx *Context, className string, symRef model.SymbolRef) string {
	return buildLookupKey(symRef.Package, className+"."+ctx.CompilerContext.GetSymbol(symRef).Name())
}

func buildGlobalVarLookupKey(pkgId *model.PackageID, name model.Name) string {
	return pkgId.OrgName.Value() + "/" + pkgId.PkgName.Value() + ":" + name.Value()
}

func GenBir(ctx *context.CompilerContext, ast *ast.BLangPackage) *BIRPackage {
	birPkg := &BIRPackage{}
	birPkg.PackageID = ast.PackageID
	genCtx := &Context{
		CompilerContext: ctx,
		importAliasMap:  make(map[string]*model.PackageID),
		packageID:       ast.PackageID,
		birPkg:          birPkg,
	}
	genCtx.typeCtx = semtypes.TypeCheckContext(ctx.GetTypeEnv())
	birPkg.GlobalVars = make(map[string]BIRGlobalVariableDcl)
	processImports(ctx, genCtx, ast.Imports, birPkg)
	for _, globalVar := range ast.GlobalVars {
		addGlobalVar(birPkg, TransformGlobalVariableDcl(genCtx, &globalVar))
	}
	for _, constant := range ast.Constants {
		addGlobalVar(birPkg, transformConstantAsGlobal(genCtx, &constant))
	}
	if ast.InitFunction != nil {
		birPkg.InitFunction = TransformFunction(genCtx, ast.InitFunction)
	}
	for i := range ast.ClassDefinitions {
		transformClassDefinition(genCtx, &ast.ClassDefinitions[i], birPkg)
	}
	for _, function := range ast.Functions {
		if function.IsNative() {
			continue
		}
		birFunc := TransformFunction(genCtx, &function)
		birPkg.Functions = append(birPkg.Functions, *birFunc)
		if birFunc.Name.Value() == "main" {
			birPkg.MainFunction = birFunc
		}
	}
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

func addGlobalVar(birPkg *BIRPackage, dcl BIRGlobalVariableDcl) {
	birPkg.GlobalVars[dcl.GlobalVarLookupKey] = dcl
}

func TransformGlobalVariableDcl(ctx *Context, ast *ast.BLangSimpleVariable) BIRGlobalVariableDcl {
	name := model.Name(ast.GetName().GetValue())
	dcl := BIRGlobalVariableDcl{}
	dcl.Pos = birLoc(ctx.CompilerContext.DiagnosticEnv(), ast.GetPosition())
	dcl.Name = name
	dcl.PkgId = ctx.packageID
	dcl.Type = ctx.CompilerContext.SymbolType(ast.Symbol())
	dcl.Flags = ast.Flags()
	dcl.GlobalVarLookupKey = buildGlobalVarLookupKey(ctx.packageID, name)
	return dcl
}

func transformConstantAsGlobal(ctx *Context, c *ast.BLangConstant) BIRGlobalVariableDcl {
	name := model.Name(c.GetName().GetValue())
	dcl := BIRGlobalVariableDcl{}
	dcl.Pos = birLoc(ctx.CompilerContext.DiagnosticEnv(), c.GetPosition())
	dcl.Name = name
	dcl.PkgId = ctx.packageID
	dcl.Type = ctx.CompilerContext.SymbolType(c.Symbol())
	dcl.Flags = c.Flags()
	dcl.GlobalVarLookupKey = buildGlobalVarLookupKey(ctx.packageID, name)
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
	birFunc.Pos = stmtCx.loc(astFunc.GetPosition())
	birFunc.Name = funcName
	birFunc.OriginalName = funcName
	birFunc.Flags = astFunc.Flags()
	ctx := stmtCx.birCx
	birFunc.FunctionLookupKey = buildFunctionLookupKeyFromSymbol(ctx, symRef)
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
			Flags: param.Flags(),
		}
	}
	if astFunc.RestParam != nil {
		restParam := astFunc.RestParam
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
		curBB.Terminator = NewReturn(ctx.loc(ast.GetPosition()))
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
	case *ast.BLangXMLNS:
		// xmlns declarations have no runtime effect.
		return statementEffect{block: curBB}
	case *ast.BLangLock:
		return lockStatement(ctx, curBB, stmt)
	default:
		panic("unexpected statement type")
	}
}

func lockStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangLock) statementEffect {
	pos := ctx.loc(stmt.GetPosition())
	if stmt.LockKey == "" {
		ctx.birCx.CompilerContext.InternalError("lock statement reached BIR-gen without a lock key", stmt.GetPosition())
	}
	key := stmt.LockKey
	bodyEntry := ctx.addBB()
	bb.Terminator = NewLockStart(key, bodyEntry, pos)
	ctx.activeLockKey = &key
	bodyEffect := blockStatement(ctx, bodyEntry, &stmt.Body)
	ctx.activeLockKey = nil
	afterLock := ctx.addBB()
	if bodyEffect.block != nil {
		bodyEffect.block.Terminator = NewLockEnd(key, afterLock, pos)
	}
	return statementEffect{block: afterLock}
}

func compoundAssignment(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangCompoundAssignment) statementEffect {
	pos := ctx.loc(stmt.GetPosition())
	if indexRef, ok := stmt.VarRef.(*ast.BLangIndexBasedAccess); ok {
		return compoundAssignmentToMember(ctx, curBB, stmt, indexRef, pos)
	}
	ref := stmt.VarRef
	valueEffect := binaryExpressionInner(ctx, curBB, stmt.OpKind, ref, stmt.Expr, stmt.Expr.GetDeterminedType(), pos)
	return assignmentStatementInner(ctx, valueEffect.block, ref, valueEffect, pos)
}

// compoundAssignmentToMember handles compound assignment with an index-based access LHS
// (e.g. `x[i] += rhs`). The container reference and index expression must be evaluated
// only once even though the LHS is conceptually both read and written.
func compoundAssignmentToMember(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangCompoundAssignment, ref *ast.BLangIndexBasedAccess, pos Location) statementEffect {
	containerEffect := assignmentContainerReference(ctx, curBB, ref.Expr)
	indexEffect := handleActionOrExpression(ctx, containerEffect.block, ref.IndexExpr)
	curBB = indexEffect.block

	loadKind, storeKind := memberAccessInstructionKinds(ref.Expr.GetDeterminedType())

	lhsValue := ctx.addTempVar(ref.GetDeterminedType())
	load := NewFieldAccess(loadKind, lhsValue, indexEffect.result, containerEffect.result, pos)
	curBB.Instructions = append(curBB.Instructions, load)
	lhsEffect := snapshotIfNeeded(ctx, expressionEffect{result: lhsValue, block: curBB}, pos)
	curBB = lhsEffect.block

	rhsEffect := handleActionOrExpression(ctx, curBB, stmt.Expr)
	rhsEffect = snapshotIfNeeded(ctx, rhsEffect, pos)
	curBB = rhsEffect.block

	resultOperand := ctx.addTempVar(ref.GetDeterminedType())
	binaryOp := NewBinaryOp(operatorKindToBinaryInstructionKind(stmt.OpKind), resultOperand, lhsEffect.result, rhsEffect.result, pos)
	curBB.Instructions = append(curBB.Instructions, binaryOp)

	store := NewFieldAccess(storeKind, containerEffect.result, indexEffect.result, resultOperand, pos)
	curBB.Instructions = append(curBB.Instructions, store)
	return statementEffect{
		block: curBB,
	}
}

func memberAccessInstructionKinds(containerType semtypes.SemType) (loadKind, storeKind InstructionKind) {
	containerType = semtypes.Diff(containerType, semtypes.NIL)
	switch {
	case semtypes.IsSubtypeSimple(containerType, semtypes.LIST):
		return INSTRUCTION_KIND_ARRAY_LOAD, INSTRUCTION_KIND_ARRAY_STORE
	case semtypes.IsSubtypeSimple(containerType, semtypes.OBJECT):
		return INSTRUCTION_KIND_OBJECT_LOAD, INSTRUCTION_KIND_OBJECT_STORE
	default:
		return INSTRUCTION_KIND_MAP_LOAD, INSTRUCTION_KIND_MAP_STORE
	}
}

func continueStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangContinue) statementEffect {
	pos := ctx.loc(stmt.GetPosition())
	curBB = emitLockEndBeforeAbruptExit(ctx, curBB, pos)
	curBB.Instructions = append(curBB.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: pos}}})
	onContinueBB := ctx.loopCtx.onContinueBB
	curBB.Terminator = NewGoto(onContinueBB, pos)
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func breakStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangBreak) statementEffect {
	pos := ctx.loc(stmt.GetPosition())
	curBB = emitLockEndBeforeAbruptExit(ctx, curBB, pos)
	curBB.Instructions = append(curBB.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: pos}}})
	onBreakBB := ctx.loopCtx.onBreakBB
	curBB.Terminator = NewGoto(onBreakBB, pos)
	return statementEffect{
		// We don't know where to add the next statement so we return nil
		block: nil,
	}
}

func whileStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangWhile) statementEffect {
	loopHead := ctx.addBB()
	// jump to loop head
	bb.Terminator = NewGoto(loopHead, ctx.loc(stmt.GetPosition()))
	condEffect := handleActionOrExpression(ctx, loopHead, stmt.Expr)

	loopBody := ctx.addBB()
	loopEnd := ctx.addBB()
	// conditionally jump to loop body
	condEffect.block.Terminator = NewBranch(condEffect.result, loopBody, loopEnd, ctx.loc(stmt.GetPosition()))

	// Push scope frame for loop body — each iteration gets its own frame
	pushScope := &PushScopeFrame{}
	pushScope.Pos = ctx.loc(stmt.GetPosition())
	loopBody.Instructions = append(loopBody.Instructions, pushScope)
	ctx.pushScope()

	ctx.addLoopCtx(loopEnd, loopHead)
	bodyEffect := blockStatement(ctx, loopBody, &stmt.Body)

	// Fill in scope frame local vars now that body has been processed
	pushScope.NumLocals = len(ctx.scopeCtx.localVars)

	// This could happen if the while block always ends return, break or continue
	if bodyEffect.block != nil {
		bodyEffect.block.Instructions = append(bodyEffect.block.Instructions, &PopScopeFrame{BIRInstructionBase: BIRInstructionBase{BIRNodeBase: BIRNodeBase{Pos: ctx.loc(stmt.GetPosition())}}})
		bodyEffect.block.Terminator = NewGoto(loopHead, ctx.loc(stmt.GetPosition()))
	}

	ctx.popLoopCtx()
	ctx.popScope()
	return statementEffect{
		block: loopEnd,
	}
}

func assignmentStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangAssignment) statementEffect {
	valueEffect := handleActionOrExpression(ctx, bb, stmt.Expr)
	return assignmentStatementInner(ctx, valueEffect.block, stmt.VarRef, valueEffect, ctx.loc(stmt.GetPosition()))
}

func assignmentStatementInner(ctx *stmtContext, bb *BIRBasicBlock, ref ast.BLangExpression, valueEffect expressionEffect, pos Location) statementEffect {
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

func assignToWildcardBindingPattern(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangWildCardBindingPattern, valueEffect expressionEffect, pos Location) statementEffect {
	refEffect := wildcardBindingPattern(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := NewMove(valueEffect.result, refEffect.result, pos)
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToSimpleVariable(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangSimpleVarRef, valueEffect expressionEffect, pos Location) statementEffect {
	refEffect := simpleVariableReference(ctx, valueEffect.block, varRef)
	currBB := refEffect.block
	mov := NewMove(valueEffect.result, refEffect.result, pos)
	currBB.Instructions = append(currBB.Instructions, mov)
	return statementEffect{
		block: currBB,
	}
}

func assignToMemberStatement(ctx *stmtContext, bb *BIRBasicBlock, varRef *ast.BLangIndexBasedAccess, valueEffect expressionEffect, pos Location) statementEffect {
	currBB := valueEffect.block
	containerRefEffect := assignmentContainerReference(ctx, currBB, varRef.Expr)
	currBB = containerRefEffect.block
	indexEffect := handleActionOrExpression(ctx, currBB, varRef.IndexExpr)
	currBB = indexEffect.block
	_, storeKind := memberAccessInstructionKinds(varRef.Expr.GetDeterminedType())
	fieldAccess := NewFieldAccess(storeKind, containerRefEffect.result, indexEffect.result, valueEffect.result, pos)
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
	exprResult := handleActionOrExpression(ctx, bb, stmt.Var.Expr)
	curBB := exprResult.block
	lhsOp := ctx.addLocalVar(varName, ty, stmt.Var.Symbol())
	move := NewMove(exprResult.result, lhsOp, ctx.loc(stmt.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, move)
	return statementEffect{
		block: curBB,
	}
}

func returnStatement(ctx *stmtContext, bb *BIRBasicBlock, stmt *ast.BLangReturn) statementEffect {
	curBB := bb
	pos := ctx.loc(stmt.GetPosition())
	if stmt.Expr != nil {
		valueEffect := handleActionOrExpression(ctx, curBB, stmt.Expr)
		curBB = valueEffect.block
		mov := NewMove(valueEffect.result, ctx.scopeCtx.retVar, pos)
		curBB.Instructions = append(curBB.Instructions, mov)
	}
	curBB = emitLockEndBeforeAbruptExit(ctx, curBB, pos)
	ret := NewReturn(pos)
	curBB.Terminator = ret
	return statementEffect{}
}

func panicStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangPanic) statementEffect {
	errorEffect := handleActionOrExpression(ctx, curBB, stmt.Expr)
	curBB = errorEffect.block
	curBB.Terminator = NewPanic(errorEffect.result, ctx.loc(stmt.GetPosition()))
	return statementEffect{}
}

func expressionStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangExpressionStmt) statementEffect {
	result := handleActionOrExpression(ctx, curBB, stmt.Expr)
	// We are ignoring the expression result (We can have one for things like call)
	return statementEffect{
		block: result.block,
	}
}

func ifStatement(ctx *stmtContext, curBB *BIRBasicBlock, stmt *ast.BLangIf) statementEffect {
	cond := handleActionOrExpression(ctx, curBB, stmt.Expr)
	curBB = cond.block
	thenBB := ctx.addBB()
	var finalBB *BIRBasicBlock
	thenEffect := blockStatement(ctx, thenBB, &stmt.Body)
	// TODO: refactor this
	if stmt.ElseStmt != nil {
		elseBB := ctx.addBB()
		// Add branch to current BB
		curBB.Terminator = NewBranch(cond.result, thenBB, elseBB, ctx.loc(stmt.GetPosition()))

		elseEffect := handleStatement(ctx, elseBB, stmt.ElseStmt)
		finalBB = ctx.addBB()
		if elseEffect.block != nil {
			elseEffect.block.Terminator = NewGoto(finalBB, ctx.loc(stmt.GetPosition()))
		}
	} else {
		finalBB = ctx.addBB()
		curBB.Terminator = NewBranch(cond.result, thenBB, finalBB, ctx.loc(stmt.GetPosition()))
	}
	// this could be nil if the control flow moved out of the if (ex: break, continue, return, etc)
	if thenEffect.block != nil {
		thenEffect.block.Terminator = NewGoto(finalBB, ctx.loc(stmt.GetPosition()))
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
	exprEffect := handleActionOrExpression(ctx, curBB, stmt.Expr)
	curBB = exprEffect.block
	matchOperand := exprEffect.result
	finalBB := ctx.addBB()

	for _, clause := range stmt.MatchClauses {
		clauseBodyBB := ctx.addBB()

		if isUnconditionalWildcard(&clause) {
			curBB.Terminator = NewGoto(clauseBodyBB, ctx.loc(stmt.GetPosition()))
			bodyEffect := blockStatement(ctx, clauseBodyBB, &clause.Body)
			if bodyEffect.block != nil {
				bodyEffect.block.Terminator = NewGoto(finalBB, ctx.loc(stmt.GetPosition()))
			}
			continue
		}

		var condOperand *BIROperand
		for _, pattern := range clause.Patterns {
			switch p := pattern.(type) {
			case *ast.BLangConstPattern:
				patternEffect := handleActionOrExpression(ctx, curBB, p.Expr)
				curBB = patternEffect.block
				eqResult := ctx.addTempVar(semtypes.BOOLEAN)
				eqPos := ctx.loc(p.Expr.GetPosition())
				binaryOp := NewBinaryOp(INSTRUCTION_KIND_EQUAL, eqResult, matchOperand, patternEffect.result, eqPos)
				curBB.Instructions = append(curBB.Instructions, binaryOp)
				condOperand = orOperands(ctx, curBB, condOperand, eqResult, eqPos)
			case *ast.BLangWildCardMatchPattern:
				// Wildcard in multi-pattern — always matches; but may have guard
				trueOperand := ctx.addTempVar(semtypes.BOOLEAN)
				constLoad := NewConstantLoad(trueOperand, true, ctx.loc(p.GetPosition()))
				curBB.Instructions = append(curBB.Instructions, constLoad)
				condOperand = orOperands(ctx, curBB, condOperand, trueOperand, ctx.loc(p.GetPosition()))
			default:
				ctx.birCx.CompilerContext.InternalError("unexpected match pattern type", pattern.GetPosition())
			}
		}

		if clause.Guard != nil {
			guardEffect := handleActionOrExpression(ctx, curBB, clause.Guard)
			curBB = guardEffect.block
			condOperand = andOperands(ctx, curBB, condOperand, guardEffect.result, ctx.loc(clause.Guard.GetPosition()))
		}

		nextCheckBB := ctx.addBB()
		curBB.Terminator = NewBranch(condOperand, clauseBodyBB, nextCheckBB, ctx.loc(stmt.GetPosition()))

		bodyEffect := blockStatement(ctx, clauseBodyBB, &clause.Body)
		if bodyEffect.block != nil {
			bodyEffect.block.Terminator = NewGoto(finalBB, ctx.loc(stmt.GetPosition()))
		}

		curBB = nextCheckBB
	}

	if !stmt.IsExhaustive {
		curBB.Terminator = NewGoto(finalBB, ctx.loc(stmt.GetPosition()))
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

func orOperands(ctx *stmtContext, bb *BIRBasicBlock, existing *BIROperand, new *BIROperand, pos Location) *BIROperand {
	if existing == nil {
		return new
	}
	result := ctx.addTempVar(semtypes.BOOLEAN)
	binaryOp := NewBinaryOp(INSTRUCTION_KIND_OR, result, existing, new, pos)
	bb.Instructions = append(bb.Instructions, binaryOp)
	return result
}

func andOperands(ctx *stmtContext, bb *BIRBasicBlock, existing *BIROperand, new *BIROperand, pos Location) *BIROperand {
	result := ctx.addTempVar(semtypes.BOOLEAN)
	binaryOp := NewBinaryOp(INSTRUCTION_KIND_AND, result, existing, new, pos)
	bb.Instructions = append(bb.Instructions, binaryOp)
	return result
}

func handleExprFunctionBody(ctx *stmtContext, body *ast.BLangExprFunctionBody) {
	curBB := ctx.addBB()
	effect := handleActionOrExpression(ctx, curBB, body.Expr)
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
	fpLoad.Pos = ctx.loc(expr.GetPosition())
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
func snapshotIfNeeded(ctx *stmtContext, effect expressionEffect, pos Location) expressionEffect {
	op := effect.result
	if _, isLocal := op.VariableDcl.(*BIRLocalVariableDcl); isLocal && semtypes.HasNoStorageIdentity(op.VariableDcl.GetType()) {
		tempOp := ctx.addTempVar(op.VariableDcl.GetType())
		effect.block.Instructions = append(effect.block.Instructions, NewMove(op, tempOp, pos))
		effect.result = tempOp
	}
	return effect
}

func handleActionOrExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr ast.BLangActionOrExpression) expressionEffect {
	switch expr := expr.(type) {
	case *ast.BLangInvocation:
		return generateCall(ctx, curBB, expr)
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
	case *ast.BLangRemoteMethodCallAction:
		return generateCall(ctx, curBB, expr)
	case *ast.BLangTypedescExpr:
		return typedescExpression(ctx, curBB, expr)
	case *ast.BLangXMLSequenceLiteral:
		return xmlSequenceLiteral(ctx, curBB, expr)
	case *ast.BLangXMLElementLiteral:
		return xmlElementLiteral(ctx, curBB, expr)
	case *ast.BLangXMLPILiteral:
		return xmlPILiteral(ctx, curBB, expr)
	case *ast.BLangXMLCommentLiteral:
		return xmlCommentLiteral(ctx, curBB, expr)
	case *ast.BLangXMLTextLiteral:
		return xmlTextLiteral(ctx, curBB, expr)
	default:
		panic(fmt.Sprintf("unexpected expression type: %T", expr))
	}
}

func typedescExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypedescExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	td := &values.TypeDesc{Type: expr.Constraint}
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(resultOperand, td, ctx.loc(expr.GetPosition())))
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func xmlTextLiteral(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangXMLTextLiteral) expressionEffect {
	pos := ctx.loc(expr.GetPosition())
	bodyOp := ctx.addTempVar(semtypes.STRING)
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(bodyOp, expr.Body, pos))
	resultOp := ctx.addTempVar(expr.GetDeterminedType())
	curBB.Instructions = append(curBB.Instructions, NewXMLTextInstr(resultOp, bodyOp, pos))
	return expressionEffect{result: resultOp, block: curBB}
}

func xmlCommentLiteral(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangXMLCommentLiteral) expressionEffect {
	pos := ctx.loc(expr.GetPosition())
	bodyOp := ctx.addTempVar(semtypes.STRING)
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(bodyOp, expr.Body, pos))
	resultOp := ctx.addTempVar(expr.GetDeterminedType())
	curBB.Instructions = append(curBB.Instructions, NewXMLCommentInstr(resultOp, bodyOp, pos))
	return expressionEffect{result: resultOp, block: curBB}
}

func xmlPILiteral(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangXMLPILiteral) expressionEffect {
	pos := ctx.loc(expr.GetPosition())
	targetOp := ctx.addTempVar(semtypes.STRING)
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(targetOp, expr.Target, pos))
	dataOp := ctx.addTempVar(semtypes.STRING)
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(dataOp, expr.Data, pos))
	resultOp := ctx.addTempVar(expr.GetDeterminedType())
	curBB.Instructions = append(curBB.Instructions, NewXMLPIInstr(resultOp, targetOp, dataOp, pos))
	return expressionEffect{result: resultOp, block: curBB}
}

func xmlElementLiteral(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangXMLElementLiteral) expressionEffect {
	pos := ctx.loc(expr.GetPosition())
	nameOp := ctx.addTempVar(semtypes.STRING)
	curBB.Instructions = append(curBB.Instructions, NewConstantLoad(nameOp, expr.Name, pos))
	var contentOp *BIROperand
	if expr.Content != nil {
		eff := handleActionOrExpression(ctx, curBB, expr.Content)
		curBB = eff.block
		contentOp = eff.result
	}
	var attrsOp *BIROperand
	if len(expr.Attrs) > 0 {
		fields := make([]mappingField, 0, len(expr.Attrs))
		for _, attr := range expr.Attrs {
			fields = append(fields, mappingField{key: attr.Name, value: attr.Value})
		}
		attrMapEff := mappingConstructorExpressionInner(ctx, curBB, ctx.birCx.stringMapType(), fields, nil, pos)
		curBB = attrMapEff.block
		attrsOp = attrMapEff.result
	}
	var namespacesOp *BIROperand
	if len(expr.Namespaces) > 0 {
		namespacesOp, curBB = buildXMLNamespacesMap(ctx, curBB, expr.Namespaces, pos)
	}
	resultOp := ctx.addTempVar(expr.GetDeterminedType())
	curBB.Instructions = append(curBB.Instructions, NewXMLElementInstr(resultOp, nameOp, contentOp, attrsOp, namespacesOp, pos))
	return expressionEffect{result: resultOp, block: curBB}
}

// buildXMLNamespacesMap constructs a string map of XML namespace declarations
// from an element's resolved Namespaces map. Keys are stored in already
// printable form ("xmlns" or "xmlns:<prefix>"). Iteration is sorted by key
// for deterministic output.
func buildXMLNamespacesMap(ctx *stmtContext, curBB *BIRBasicBlock, ns map[string]string, pos Location) (*BIROperand, *BIRBasicBlock) {
	keys := make([]string, 0, len(ns))
	for k := range ns {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	entries := make([]MappingConstructorEntry, 0, len(keys))
	for _, k := range keys {
		keyOp := ctx.addTempVar(semtypes.STRING)
		curBB.Instructions = append(curBB.Instructions, NewConstantLoad(keyOp, k, pos))
		valOp := ctx.addTempVar(semtypes.STRING)
		curBB.Instructions = append(curBB.Instructions, NewConstantLoad(valOp, ns[k], pos))
		entries = append(entries, &MappingConstructorKeyValueEntry{keyOp: keyOp, valueOp: valOp})
	}
	resultOp := ctx.addTempVar(ctx.birCx.stringMapType())
	curBB.Instructions = append(curBB.Instructions, NewMapConstructor(ctx.birCx.stringMapType(), resultOp, entries, nil, false, pos))
	return resultOp, curBB
}

func xmlSequenceLiteral(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangXMLSequenceLiteral) expressionEffect {
	if len(expr.Children) == 1 {
		return handleActionOrExpression(ctx, curBB, expr.Children[0])
	}
	pos := ctx.loc(expr.GetPosition())
	var childOps []*BIROperand
	for _, child := range expr.Children {
		eff := handleActionOrExpression(ctx, curBB, child)
		curBB = eff.block
		childOps = append(childOps, eff.result)
	}
	resultOp := ctx.addTempVar(expr.GetDeterminedType())
	curBB.Instructions = append(curBB.Instructions, NewXMLSequenceInstr(resultOp, childOps, pos))
	return expressionEffect{result: resultOp, block: curBB}
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
	var defaults []MappingConstructorDefaultEntry
	for _, fd := range expr.FieldDefaults {
		defaults = append(defaults, MappingConstructorDefaultEntry{
			FieldName:         fd.FieldName,
			FunctionLookupKey: buildFunctionLookupKeyFromSymbol(ctx.birCx, fd.FnRef),
		})
	}
	return mappingConstructorExpressionInner(ctx, curBB, expr.GetDeterminedType(), fields, defaults, ctx.loc(expr.GetPosition()))
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

func mappingConstructorExpressionInner(ctx *stmtContext, curBB *BIRBasicBlock, mapType semtypes.SemType, fields []mappingField, defaults []MappingConstructorDefaultEntry, pos Location) expressionEffect {
	var entries []MappingConstructorEntry
	for _, field := range fields {
		keyOperand := ctx.addTempVar(semtypes.STRING)
		keyLoad := NewConstantLoad(keyOperand, field.key, pos)
		curBB.Instructions = append(curBB.Instructions, keyLoad)

		valueEffect := handleActionOrExpression(ctx, curBB, field.value)
		curBB = valueEffect.block
		entries = append(entries, &MappingConstructorKeyValueEntry{
			keyOp:   keyOperand,
			valueOp: valueEffect.result,
		})
	}
	resultOperand := ctx.addTempVar(mapType)
	isReadonly := semtypes.IsSubtype(ctx.birCx.typeCtx, mapType, semtypes.VAL_READONLY)
	newMap := NewMapConstructor(mapType, resultOperand, entries, defaults, isReadonly, pos)
	curBB.Instructions = append(curBB.Instructions, newMap)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func errorConstructorExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangErrorConstructorExpr) expressionEffect {
	// Message is the first positional arg
	msgEffect := handleActionOrExpression(ctx, curBB, expr.PositionalArgs[0])
	curBB = msgEffect.block

	// Cause is the optional second positional arg
	var causeOp *BIROperand
	if len(expr.PositionalArgs) > 1 {
		causeEffect := handleActionOrExpression(ctx, curBB, expr.PositionalArgs[1])
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
		detailEffect := mappingConstructorExpressionInner(ctx, curBB, semtypes.MAPPING, fields, nil, ctx.loc(expr.GetPosition()))
		curBB = detailEffect.block
		detailOp = detailEffect.result
	}

	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeName := ""
	if expr.ErrorTypeRef != nil {
		typeName = expr.ErrorTypeRef.TypeName.Value
	}
	newError := NewErrorConstructor(expr.GetDeterminedType(), typeName, resultOperand, msgEffect.result, causeOp, detailOp, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, newError)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func typeConversionExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypeConversionExpr) expressionEffect {
	exprEffect := handleActionOrExpression(ctx, curBB, expr.Expression)
	curBB = exprEffect.block
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeCast := NewTypeCast(expr.TypeDescriptor.GetDeterminedType(), resultOperand, exprEffect.result, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, typeCast)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func typeTestExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTypeTestExpr) expressionEffect {
	exprEffect := handleActionOrExpression(ctx, curBB, expr.Expr)
	curBB = exprEffect.block
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	typeTest := &TypeTest{}
	typeTest.Pos = ctx.loc(expr.GetPosition())
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

// materializeFiller emits BIR instructions that construct a fresh filler
// value for ty at runtime.
func materializeFiller(ctx *stmtContext, bb *BIRBasicBlock, ty semtypes.SemType, f semtypes.Filler, pos Location) (*BIROperand, *BIRBasicBlock) {
	tyCx := ctx.birCx.typeCtx
	switch f := f.(type) {
	case semtypes.SingleValueFiller:
		operand := ctx.addTempVar(ty)
		bb.Instructions = append(bb.Instructions, NewConstantLoad(operand, f.Value, pos))
		return operand, bb
	case semtypes.MappingFiller:
		operand := ctx.addTempVar(f.Type)
		mapReadonly := semtypes.IsSubtype(tyCx, f.Type, semtypes.VAL_READONLY)
		bb.Instructions = append(bb.Instructions, NewMapConstructor(f.Type, operand, nil, nil, mapReadonly, pos))
		return operand, bb
	case semtypes.ListFiller:
		memberOperands := make([]*BIROperand, len(f.Members))
		for i, memberFiller := range f.Members {
			memberOperands[i], bb = materializeFiller(ctx, bb, f.Atomic.MemberAtInnerVal(i), memberFiller, pos)
		}
		sizeOperand := ctx.addTempVar(semtypes.INT)
		bb.Instructions = append(bb.Instructions, NewConstantLoad(sizeOperand, int64(len(memberOperands)), pos))
		restFiller, _ := values.FillerFactoryFor(tyCx, f.Atomic.Rest())
		operand := ctx.addTempVar(f.Type)
		listReadonly := semtypes.IsSubtype(tyCx, f.Type, semtypes.VAL_READONLY)
		bb.Instructions = append(bb.Instructions, NewArrayConstructor(f.Type, operand, sizeOperand, memberOperands, restFiller, listReadonly, pos))
		return operand, bb
	default:
		panic(fmt.Sprintf("unsupported filler kind %T in BIR generation", f))
	}
}

func listConstructorExpression(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangListConstructorExpr) expressionEffect {
	initValues := make([]*BIROperand, len(expr.Exprs))
	for i, expr := range expr.Exprs {
		exprEffect := handleActionOrExpression(ctx, bb, expr)
		bb = exprEffect.block
		initValues[i] = exprEffect.result
	}

	lat := expr.AtomicType
	exprPos := ctx.loc(expr.GetPosition())
	tyCx := ctx.birCx.typeCtx
	for i := len(expr.Exprs); i < lat.Members.FixedLength; i++ {
		ty := lat.MemberAtInnerVal(i)
		filler, ok := semtypes.FillerValue(tyCx, ty)
		if !ok {
			ctx.birCx.CompilerContext.InternalError("no filler value for list member type; semantic analysis should have rejected this", expr.GetPosition())
		}
		var fillerOperand *BIROperand
		fillerOperand, bb = materializeFiller(ctx, bb, ty, filler, exprPos)
		initValues = append(initValues, fillerOperand)
	}
	restFiller, _ := values.FillerFactoryFor(tyCx, lat.Rest())

	sizeOperand := ctx.addTempVar(semtypes.INT)
	constantLoad := NewConstantLoad(sizeOperand, int64(len(initValues)), exprPos)
	bb.Instructions = append(bb.Instructions, constantLoad)

	resultOperand := ctx.addTempVar(semtypes.LIST)
	listTy := expr.GetDeterminedType()
	isReadonly := semtypes.IsSubtype(tyCx, listTy, semtypes.VAL_READONLY)
	newArray := NewArrayConstructor(listTy, resultOperand, sizeOperand, initValues, restFiller, isReadonly, exprPos)
	bb.Instructions = append(bb.Instructions, newArray)
	return expressionEffect{
		result: resultOperand,
		block:  bb,
	}
}

// assignmentContainerReference produces the container reference for an indexed assignment LHS.
// When the container is itself an index-based access on a list or mapping, the inner read
// must be a filling load so that intermediate arrays grow (and fill) and absent map keys
// are populated with a filler value before storing.
func assignmentContainerReference(ctx *stmtContext, bb *BIRBasicBlock, expr ast.BLangExpression) expressionEffect {
	inner, ok := expr.(*ast.BLangIndexBasedAccess)
	if !ok {
		return handleActionOrExpression(ctx, bb, expr)
	}
	// The container of an indexed lvalue access can show up as a nilable type
	// when it itself comes from another map index (e.g. `m["a"]["b"]` where the
	// inner lookup nominally yields `T?`). After filling, the container is
	// guaranteed non-nil, so we strip `()` before classifying.
	containerType := semtypes.Diff(inner.Expr.GetDeterminedType(), semtypes.NIL)
	var fillingKind InstructionKind
	var filler values.FillerFactory
	switch {
	case semtypes.IsSubtypeSimple(containerType, semtypes.LIST):
		fillingKind = INSTRUCTION_KIND_ARRAY_FILLING_LOAD
	case semtypes.IsSubtypeSimple(containerType, semtypes.MAPPING):
		fillingKind = INSTRUCTION_KIND_MAP_FILLING_LOAD
		tyCx := semtypes.TypeCheckContext(ctx.birCx.CompilerContext.GetTypeEnv())
		valueType := semtypes.MappingMemberTypeInnerVal(tyCx, containerType, semtypes.STRING)
		filler, _ = values.FillerFactoryFor(tyCx, valueType)
	default:
		return handleActionOrExpression(ctx, bb, expr)
	}
	resultOperand := ctx.addTempVar(inner.GetDeterminedType())
	indexEffect := handleActionOrExpression(ctx, bb, inner.IndexExpr)
	containerRefEffect := assignmentContainerReference(ctx, indexEffect.block, inner.Expr)
	fieldAccess := NewFieldAccess(fillingKind, resultOperand, indexEffect.result, containerRefEffect.result, ctx.loc(inner.GetPosition()))
	fieldAccess.Filler = filler
	containerRefEffect.block.Instructions = append(containerRefEffect.block.Instructions, fieldAccess)
	return expressionEffect{
		result: resultOperand,
		block:  containerRefEffect.block,
	}
}

func indexBasedAccess(ctx *stmtContext, bb *BIRBasicBlock, expr *ast.BLangIndexBasedAccess) expressionEffect {
	// Assignment is handled in assignmentStatement to this is always a load
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	loadKind, _ := memberAccessInstructionKinds(expr.Expr.GetDeterminedType())
	indexEffect := handleActionOrExpression(ctx, bb, expr.IndexExpr)
	containerRefEffect := handleActionOrExpression(ctx, indexEffect.block, expr.Expr)
	currBB := containerRefEffect.block
	fieldAccess := NewFieldAccess(loadKind, resultOperand, indexEffect.result, containerRefEffect.result, ctx.loc(expr.GetPosition()))
	currBB.Instructions = append(currBB.Instructions, fieldAccess)
	return expressionEffect{
		result: resultOperand,
		block:  currBB,
	}
}

func groupExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangGroupExpr) expressionEffect {
	return handleActionOrExpression(ctx, curBB, expr.Expression)
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
	opEffect := handleActionOrExpression(ctx, bb, expr.Expr)

	resultOperand := ctx.addTempVar(expr.GetDeterminedType())
	curBB := opEffect.block
	unaryOp := NewUnaryOp(kind, resultOperand, opEffect.result, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, unaryOp)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

type callable interface {
	ast.BLangActionOrExpression
	ResolvedSymbol() model.SymbolRef
	Receiver() ast.BLangExpression
	CallArgs() []ast.BLangExpression
	GetName() ast.IdentifierNode
}

func generateCall(ctx *stmtContext, bb *BIRBasicBlock, callable callable) expressionEffect {
	curBB := bb
	if ast.IsStreamOperation(callable) {
		return streamMethodCall(ctx, curBB, callable)
	}
	var args []BIROperand
	isMethodCall := false

	if callable.Receiver() != nil {
		effect := handleActionOrExpression(ctx, curBB, callable.Receiver())
		curBB = effect.block
		args = append(args, *effect.result)
		isMethodCall = true
	}

	for _, arg := range callable.CallArgs() {
		effect := handleActionOrExpression(ctx, curBB, arg)
		effect = snapshotIfNeeded(ctx, effect, ctx.loc(callable.GetPosition()))
		curBB = effect.block
		args = append(args, *effect.result)
	}

	thenBB := ctx.addBB()
	resultOperand := ctx.addTempVar(callable.GetDeterminedType())
	callName := callable.GetName().GetValue()
	if _, isRemote := callable.(*ast.BLangRemoteMethodCallAction); isRemote {
		callName = model.RemoteMethodName(callName)
	}
	call := NewCall(INSTRUCTION_KIND_CALL, args, model.Name(callName), thenBB, resultOperand, ctx.loc(callable.GetPosition()))
	call.IsMethodCall = isMethodCall

	symRef := callable.ResolvedSymbol()
	sym := ctx.birCx.CompilerContext.GetSymbol(symRef)
	if sym.Kind() == model.SymbolKindFunction {
		call.FunctionLookupKey = buildFunctionLookupKeyFromSymbol(ctx.birCx, symRef)
		if inv, ok := callable.(*ast.BLangInvocation); ok && inv.PkgAlias != nil && inv.PkgAlias.Value != "" {
			call.CalleePkg = ctx.birCx.importAliasMap[inv.PkgAlias.Value]
		} else if ctx.birCx.packageID != nil {
			call.CalleePkg = ctx.birCx.packageID
		}
	} else {
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
	constantLoad := NewConstantLoad(resultOperand, expr.Value, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, constantLoad)
	return expressionEffect{
		result: resultOperand,
		block:  curBB,
	}
}

func operatorKindToBinaryInstructionKind(opKind model.OperatorKind) InstructionKind {
	switch opKind {
	case model.OperatorKind_ADD:
		return INSTRUCTION_KIND_ADD
	case model.OperatorKind_SUB:
		return INSTRUCTION_KIND_SUB
	case model.OperatorKind_MUL:
		return INSTRUCTION_KIND_MUL
	case model.OperatorKind_DIV:
		return INSTRUCTION_KIND_DIV
	case model.OperatorKind_MOD:
		return INSTRUCTION_KIND_MOD
	case model.OperatorKind_EQUAL:
		return INSTRUCTION_KIND_EQUAL
	case model.OperatorKind_NOT_EQUAL:
		return INSTRUCTION_KIND_NOT_EQUAL
	case model.OperatorKind_GREATER_THAN:
		return INSTRUCTION_KIND_GREATER_THAN
	case model.OperatorKind_GREATER_EQUAL:
		return INSTRUCTION_KIND_GREATER_EQUAL
	case model.OperatorKind_LESS_THAN:
		return INSTRUCTION_KIND_LESS_THAN
	case model.OperatorKind_LESS_EQUAL:
		return INSTRUCTION_KIND_LESS_EQUAL
	case model.OperatorKind_REF_EQUAL:
		return INSTRUCTION_KIND_REF_EQUAL
	case model.OperatorKind_REF_NOT_EQUAL:
		return INSTRUCTION_KIND_REF_NOT_EQUAL
	case model.OperatorKind_BITWISE_AND:
		return INSTRUCTION_KIND_BITWISE_AND
	case model.OperatorKind_BITWISE_OR:
		return INSTRUCTION_KIND_BITWISE_OR
	case model.OperatorKind_BITWISE_XOR:
		return INSTRUCTION_KIND_BITWISE_XOR
	case model.OperatorKind_BITWISE_LEFT_SHIFT:
		return INSTRUCTION_KIND_BITWISE_LEFT_SHIFT
	case model.OperatorKind_BITWISE_RIGHT_SHIFT:
		return INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT
	case model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
		return INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT
	default:
		panic("unexpected binary operator kind")
	}
}

func binaryExpressionInner(ctx *stmtContext, curBB *BIRBasicBlock, opKind model.OperatorKind, lhsExpr ast.BLangExpression, rhsExpr ast.BLangActionOrExpression, resultType semtypes.SemType, pos Location) expressionEffect {
	kind := operatorKindToBinaryInstructionKind(opKind)
	resultOperand := ctx.addTempVar(resultType)
	op1Effect := handleActionOrExpression(ctx, curBB, lhsExpr)
	op1Effect = snapshotIfNeeded(ctx, op1Effect, pos)
	curBB = op1Effect.block
	op2Effect := handleActionOrExpression(ctx, curBB, rhsExpr)
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
		return binaryExpressionInner(ctx, curBB, expr.OpKind, expr.LhsExpr, expr.RhsExpr, expr.GetDeterminedType(), ctx.loc(expr.GetPosition()))
	}
}

func logicalAndExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	lhsEffect := handleActionOrExpression(ctx, curBB, expr.LhsExpr)
	curBB = lhsEffect.block

	mov := NewMove(lhsEffect.result, resultOperand, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, mov)

	evalRhsBB := ctx.addBB()
	doneBB := ctx.addBB()

	curBB.Terminator = NewBranch(lhsEffect.result, evalRhsBB, doneBB, ctx.loc(expr.GetPosition()))

	rhsEffect := handleActionOrExpression(ctx, evalRhsBB, expr.RhsExpr)
	rhsBB := rhsEffect.block

	rhsMov := NewMove(rhsEffect.result, resultOperand, ctx.loc(expr.GetPosition()))

	rhsBB.Instructions = append(rhsBB.Instructions, rhsMov)
	rhsBB.Terminator = NewGoto(doneBB, ctx.loc(expr.GetPosition()))

	return expressionEffect{
		result: resultOperand,
		block:  doneBB,
	}
}

func logicalOrExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangBinaryExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	lhsEffect := handleActionOrExpression(ctx, curBB, expr.LhsExpr)
	curBB = lhsEffect.block

	mov := NewMove(lhsEffect.result, resultOperand, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, mov)

	evalRhsBB := ctx.addBB()
	doneBB := ctx.addBB()

	curBB.Terminator = NewBranch(lhsEffect.result, doneBB, evalRhsBB, ctx.loc(expr.GetPosition()))

	rhsEffect := handleActionOrExpression(ctx, evalRhsBB, expr.RhsExpr)
	rhsBB := rhsEffect.block

	rhsMov := NewMove(rhsEffect.result, resultOperand, ctx.loc(expr.GetPosition()))
	rhsBB.Instructions = append(rhsBB.Instructions, rhsMov)
	rhsBB.Terminator = NewGoto(doneBB, ctx.loc(expr.GetPosition()))

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
		fpLoad := NewFPLoad(lookupKey, funcType, resultOperand, ctx.loc(expr.GetPosition()))
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
	gv.GlobalVarLookupKey = buildGlobalVarLookupKey(pkgId, gv.Name)
	return expressionEffect{
		result: &BIROperand{VariableDcl: gv},
		block:  curBB,
	}
}

func trapExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangTrapExpr) expressionEffect {
	resultOperand := ctx.addTempVar(expr.GetDeterminedType())

	trapStartBB := ctx.addBB()
	curBB.Terminator = NewGoto(trapStartBB, ctx.loc(expr.GetPosition()))

	innerEffect := handleActionOrExpression(ctx, trapStartBB, expr.Expr)
	trapEndBB := innerEffect.block

	mov := NewMove(innerEffect.result, resultOperand, ctx.loc(expr.GetPosition()))
	trapEndBB.Instructions = append(trapEndBB.Instructions, mov)

	afterTrapBB := ctx.addBB()
	trapEndBB.Terminator = NewGoto(afterTrapBB, ctx.loc(expr.GetPosition()))

	ctx.errorEntries = append(ctx.errorEntries, BIRErrorEntry{
		Start:   trapStartBB.Number,
		End:     trapEndBB.Number,
		Target:  afterTrapBB.Number,
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

	classLookupKey := buildLookupKey(class.Symbol().Package, ctx.CompilerContext.SymbolName(class.Symbol()))
	birClassDef := &BIRClassDef{
		Name:      className,
		LookupKey: classLookupKey,
		VTable:    make(map[string]*BIRFunction),
	}

	for _, field := range class.Fields {
		birClassDef.Fields = append(birClassDef.Fields, ObjectField{
			Name: field.GetName().GetValue(),
			Ty:   ctx.CompilerContext.SymbolType(field.Symbol()),
		})
	}

	initFunc := transformFunctionInner(&stmtContext{birCx: ctx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), isFunctionBoundary: true}}, class.InitFunction, &selfRef)
	initFunc.FunctionLookupKey = buildMethodLookupKeyFromSymbol(ctx, className.Value(), class.InitFunction.Symbol())
	birClassDef.VTable["init"] = initFunc

	for methodName, method := range class.Methods {
		lookupKey := buildMethodLookupKeyFromSymbol(ctx, className.Value(), method.Symbol())
		var fn *BIRFunction
		if method.IsNative() {
			fn = &BIRFunction{
				Name:              model.Name(method.GetName().GetValue()),
				OriginalName:      model.Name(method.GetName().GetValue()),
				Flags:             method.Flags(),
				FunctionLookupKey: lookupKey,
			}
			fn.Pos = birLoc(ctx.CompilerContext.DiagnosticEnv(), method.GetPosition())
		} else {
			fn = transformFunctionInner(&stmtContext{birCx: ctx, scopeCtx: &scopeContext{varMap: make(map[model.SymbolRef]*BIROperand), isFunctionBoundary: true}}, method, &selfRef)
			fn.FunctionLookupKey = lookupKey
		}
		birClassDef.VTable[methodName] = fn
	}

	birPkg.ClassDefs = append(birPkg.ClassDefs, *birClassDef)
}

func newExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangNewExpression) expressionEffect {
	if semtypes.IsSubtypeSimple(expr.GetDeterminedType(), semtypes.STREAM) {
		return newStreamExpression(ctx, curBB, expr)
	}
	classSymbol := expr.ClassSymbol
	className := ctx.birCx.CompilerContext.SymbolName(classSymbol)
	classLookupKey := buildLookupKey(classSymbol.Package, className)

	object := ctx.addTempVar(expr.GetDeterminedType())
	newObj := NewObjectConstructor(classLookupKey, object, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, newObj)

	var args []BIROperand
	args = append(args, *object)
	for _, arg := range expr.ArgsExprs {
		argEffect := handleActionOrExpression(ctx, curBB, arg)
		curBB = argEffect.block
		args = append(args, *argEffect.result)
	}

	initMethodLookupKey := classLookupKey + ".init"
	initResult := ctx.addTempVar(semtypes.Union(semtypes.NIL, semtypes.ERROR))
	initDoneBB := ctx.addBB()
	call := NewCall(INSTRUCTION_KIND_CALL, args, model.Name("init"), initDoneBB, initResult, ctx.loc(expr.GetPosition()))
	call.IsMethodCall = true
	call.CachedMethodLookupKey = initMethodLookupKey
	curBB.Terminator = call

	result := ctx.addTempVar(expr.DeterminedType)
	isInitResultNil := ctx.addTempVar(semtypes.BOOLEAN)
	nilCheck := NewTypeTest(semtypes.NIL, isInitResultNil, initResult, ctx.loc(expr.GetPosition()))
	initDoneBB.Instructions = append(initDoneBB.Instructions, nilCheck)

	assignObjectBB := ctx.addBB()
	assignErrorBB := ctx.addBB()
	thenBB := ctx.addBB()
	initDoneBB.Terminator = NewBranch(isInitResultNil, assignObjectBB, assignErrorBB, ctx.loc(expr.GetPosition()))

	assignObjectBB.Instructions = append(assignObjectBB.Instructions, NewMove(object, result, ctx.loc(expr.GetPosition())))
	assignObjectBB.Terminator = NewGoto(thenBB, ctx.loc(expr.GetPosition()))

	assignErrorBB.Instructions = append(assignErrorBB.Instructions, NewMove(initResult, result, ctx.loc(expr.GetPosition())))
	assignErrorBB.Terminator = NewGoto(thenBB, ctx.loc(expr.GetPosition()))

	return expressionEffect{
		result: result,
		block:  thenBB,
	}
}

func streamMethodCall(ctx *stmtContext, curBB *BIRBasicBlock, callable callable) expressionEffect {
	recvEffect := handleActionOrExpression(ctx, curBB, callable.Receiver())
	curBB = recvEffect.block
	result := ctx.addTempVar(callable.GetDeterminedType())
	pos := ctx.loc(callable.GetPosition())
	switch callable.GetName().GetValue() {
	case "next":
		curBB.Instructions = append(curBB.Instructions, NewStreamNext(result, recvEffect.result, pos))
	case "close":
		curBB.Instructions = append(curBB.Instructions, NewStreamClose(result, recvEffect.result, pos))
	default:
		ctx.birCx.CompilerContext.InternalError("unexpected stream method: "+callable.GetName().GetValue(), callable.GetPosition())
	}
	return expressionEffect{result: result, block: curBB}
}

func newStreamExpression(ctx *stmtContext, curBB *BIRBasicBlock, expr *ast.BLangNewExpression) expressionEffect {
	argEffect := handleActionOrExpression(ctx, curBB, expr.ArgsExprs[0])
	curBB = argEffect.block
	result := ctx.addTempVar(expr.GetDeterminedType())
	instr := NewStreamConstructor(expr.GetDeterminedType(), result, argEffect.result, ctx.loc(expr.GetPosition()))
	curBB.Instructions = append(curBB.Instructions, instr)
	return expressionEffect{result: result, block: curBB}
}

func appendIfNotNil[T any](slice []T, item *T) []T {
	if item != nil {
		slice = append(slice, *item)
	}
	return slice
}
