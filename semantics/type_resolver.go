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

package semantics

import (
	"fmt"
	"math/big"
	"math/bits"
	"sort"
	"strconv"
	"strings"
	"sync"

	"ballerina-lang-go/ast"
	balCommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"

	array "ballerina-lang-go/lib/array/compile"
	bInt "ballerina-lang-go/lib/int/compile"
)

type typeResolver interface {
	ast.Visitor
	typeContext() semtypes.Context
	expectedReturnType() semtypes.SemType
	parent() typeResolver
	typeEnv() semtypes.Env

	// Error reporting (proxied from CompilerContext)
	semanticError(message string, loc diagnostics.Location)
	internalError(message string, loc diagnostics.Location)
	unimplemented(message string, loc diagnostics.Location)
	syntaxError(message string, loc diagnostics.Location)

	// Symbol management (proxied from CompilerContext)
	symbolType(ref model.SymbolRef) semtypes.SemType
	setSymbolType(ref model.SymbolRef, ty semtypes.SemType)
	getSymbol(ref model.SymbolRef) model.Symbol
	unnarrowedSymbol(ref model.SymbolRef) model.SymbolRef
	symbolName(ref model.SymbolRef) string
	setTypeDefinition(ref model.SymbolRef, defn model.TypeDefinition)
	getTypeDefinition(ref model.SymbolRef) (model.TypeDefinition, bool)
	createNarrowedSymbol(ref model.SymbolRef) model.SymbolRef

	// Import management
	lookupImportedSymbols(pkgName string) (model.ExportedSymbolSpace, bool)
	addImplicitImport(pkgName string, imp ast.BLangImportPackage)
	hasImplicitImport(pkgName string) bool
}

type packageTypeResolver struct {
	ctx             *context.CompilerContext
	tyCtx           semtypes.Context
	importedSymbols map[string]model.ExportedSymbolSpace
	pkg             *ast.BLangPackage
	implicitImports map[string]ast.BLangImportPackage
}

var _ ast.Visitor = &packageTypeResolver{}

func (t *packageTypeResolver) typeContext() semtypes.Context        { return t.tyCtx }
func (t *packageTypeResolver) expectedReturnType() semtypes.SemType { return nil }
func (t *packageTypeResolver) parent() typeResolver                 { return nil }
func (t *packageTypeResolver) typeEnv() semtypes.Env                { return t.ctx.GetTypeEnv() }

func (t *packageTypeResolver) semanticError(msg string, loc diagnostics.Location) {
	t.ctx.SemanticError(msg, loc)
}

func (t *packageTypeResolver) internalError(msg string, loc diagnostics.Location) {
	t.ctx.InternalError(msg, loc)
}

func (t *packageTypeResolver) unimplemented(msg string, loc diagnostics.Location) {
	t.ctx.Unimplemented(msg, loc)
}

func (t *packageTypeResolver) syntaxError(msg string, loc diagnostics.Location) {
	t.ctx.SyntaxError(msg, loc)
}

func (t *packageTypeResolver) symbolType(ref model.SymbolRef) semtypes.SemType {
	return t.ctx.SymbolType(ref)
}

func (t *packageTypeResolver) setSymbolType(ref model.SymbolRef, ty semtypes.SemType) {
	t.ctx.SetSymbolType(ref, ty)
}

func (t *packageTypeResolver) getSymbol(ref model.SymbolRef) model.Symbol {
	return t.ctx.GetSymbol(ref)
}

func (t *packageTypeResolver) unnarrowedSymbol(ref model.SymbolRef) model.SymbolRef {
	return t.ctx.UnnarrowedSymbol(ref)
}

func (t *packageTypeResolver) symbolName(ref model.SymbolRef) string {
	return t.ctx.SymbolName(ref)
}

func (t *packageTypeResolver) setTypeDefinition(ref model.SymbolRef, defn model.TypeDefinition) {
	t.ctx.SetTypeDefinition(ref, defn)
}

func (t *packageTypeResolver) getTypeDefinition(ref model.SymbolRef) (model.TypeDefinition, bool) {
	return t.ctx.GetTypeDefinition(ref)
}

func (t *packageTypeResolver) createNarrowedSymbol(ref model.SymbolRef) model.SymbolRef {
	return t.ctx.CreateNarrowedSymbol(ref)
}

func (t *packageTypeResolver) lookupImportedSymbols(name string) (model.ExportedSymbolSpace, bool) {
	s, ok := t.importedSymbols[name]
	return s, ok
}

func (t *packageTypeResolver) addImplicitImport(name string, imp ast.BLangImportPackage) {
	t.implicitImports[name] = imp
}

func (t *packageTypeResolver) hasImplicitImport(name string) bool {
	_, ok := t.implicitImports[name]
	return ok
}

func (t *packageTypeResolver) Visit(node ast.BLangNode) ast.Visitor { return visit(t, node) }
func (t *packageTypeResolver) VisitTypeData(td *model.TypeData) ast.Visitor {
	return visitTypeData(t, td)
}

type functionTypeResolver struct {
	parentResolver  typeResolver
	tyCtx           semtypes.Context
	retTy           semtypes.SemType
	implicitImports map[string]ast.BLangImportPackage
}

func (f *functionTypeResolver) typeContext() semtypes.Context        { return f.tyCtx }
func (f *functionTypeResolver) expectedReturnType() semtypes.SemType { return f.retTy }
func (f *functionTypeResolver) parent() typeResolver                 { return f.parentResolver }
func (f *functionTypeResolver) typeEnv() semtypes.Env                { return f.parentResolver.typeEnv() }

func (f *functionTypeResolver) semanticError(msg string, loc diagnostics.Location) {
	f.parentResolver.semanticError(msg, loc)
}

func (f *functionTypeResolver) internalError(msg string, loc diagnostics.Location) {
	f.parentResolver.internalError(msg, loc)
}

func (f *functionTypeResolver) unimplemented(msg string, loc diagnostics.Location) {
	f.parentResolver.unimplemented(msg, loc)
}

func (f *functionTypeResolver) syntaxError(msg string, loc diagnostics.Location) {
	f.parentResolver.syntaxError(msg, loc)
}

func (f *functionTypeResolver) symbolType(ref model.SymbolRef) semtypes.SemType {
	return f.parentResolver.symbolType(ref)
}

func (f *functionTypeResolver) setSymbolType(ref model.SymbolRef, ty semtypes.SemType) {
	f.parentResolver.setSymbolType(ref, ty)
}

func (f *functionTypeResolver) getSymbol(ref model.SymbolRef) model.Symbol {
	return f.parentResolver.getSymbol(ref)
}

func (f *functionTypeResolver) unnarrowedSymbol(ref model.SymbolRef) model.SymbolRef {
	return f.parentResolver.unnarrowedSymbol(ref)
}

func (f *functionTypeResolver) symbolName(ref model.SymbolRef) string {
	return f.parentResolver.symbolName(ref)
}

func (f *functionTypeResolver) setTypeDefinition(ref model.SymbolRef, defn model.TypeDefinition) {
	f.parentResolver.setTypeDefinition(ref, defn)
}

func (f *functionTypeResolver) getTypeDefinition(ref model.SymbolRef) (model.TypeDefinition, bool) {
	return f.parentResolver.getTypeDefinition(ref)
}

func (f *functionTypeResolver) createNarrowedSymbol(ref model.SymbolRef) model.SymbolRef {
	return f.parentResolver.createNarrowedSymbol(ref)
}

func (f *functionTypeResolver) lookupImportedSymbols(name string) (model.ExportedSymbolSpace, bool) {
	return f.parentResolver.lookupImportedSymbols(name)
}

func (f *functionTypeResolver) addImplicitImport(name string, imp ast.BLangImportPackage) {
	f.implicitImports[name] = imp
}

func (f *functionTypeResolver) hasImplicitImport(name string) bool {
	_, ok := f.implicitImports[name]
	return ok
}

func (f *functionTypeResolver) Visit(node ast.BLangNode) ast.Visitor { return visit(f, node) }
func (f *functionTypeResolver) VisitTypeData(td *model.TypeData) ast.Visitor {
	return visitTypeData(f, td)
}

func newPackageTypeResolver(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *packageTypeResolver {
	return &packageTypeResolver{
		ctx:             ctx,
		tyCtx:           semtypes.ContextFrom(ctx.GetTypeEnv()),
		importedSymbols: importedSymbols,
		pkg:             pkg,
		implicitImports: make(map[string]ast.BLangImportPackage),
	}
}

// ResolveTopLevelNodes resolves type definitions, function signatures, and constants.
// After this (for the given package) all the semtypes are known. This means after resolving types of all the packages
// it is safe to use the closed world assumption to optimize type checks.
func ResolveTopLevelNodes(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	t := newPackageTypeResolver(ctx, pkg, importedSymbols)
	t.resolveTopLevelTypes(ctx, pkg)
}

// ResolveLocalNodes resolves the types of function bodies and remaining inner nodes.
func ResolveLocalNodes(ctx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) {
	p := newPackageTypeResolver(ctx, pkg, importedSymbols)
	resolvers := make([]*functionTypeResolver, len(pkg.Functions))
	var wg sync.WaitGroup
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resolvers[idx] = resolveFunctionBody(p, fn)
		}(i)
	}
	wg.Wait()

	seen := make(map[string]bool, len(resolvers))
	for _, ft := range resolvers {
		if ft == nil {
			continue
		}
		for name, importNode := range ft.implicitImports {
			if !seen[name] {
				seen[name] = true
				pkg.Imports = append(pkg.Imports, importNode)
			}
		}
	}
}

func resolveFunctionBody(p *packageTypeResolver, fn *ast.BLangFunction) *functionTypeResolver {
	fnSymbol := p.getSymbol(fn.Symbol())
	fnSym, ok := fnSymbol.(model.FunctionSymbol)
	if !ok {
		p.internalError("expected function symbol", fn.GetPosition())
		return nil
	}
	ft := &functionTypeResolver{
		parentResolver:  p,
		tyCtx:           semtypes.ContextFrom(p.typeEnv()),
		retTy:           fnSym.Signature().ReturnType,
		implicitImports: make(map[string]ast.BLangImportPackage),
	}
	switch body := fn.Body.(type) {
	case *ast.BLangExternFunctionBody:
		_ = body
	case *ast.BLangBlockFunctionBody:
		resolveBlockStatements(ft, nil, body.Stmts)
		body.SetDeterminedType(semtypes.NEVER)
	case *ast.BLangExprFunctionBody:
		resolveExpression(ft, nil, body.Expr.(ast.BLangExpression), ft.retTy)
	default:
		p.internalError("unexpected function body kind", fn.Body.GetPosition())
	}
	return ft
}

func (t *packageTypeResolver) resolveTopLevelTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol()
		ctx.SetTypeDefinition(symbol, defn)
	}
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		if _, ok := resolveTypeDefinition(t, defn, 0); !ok {
			return
		}
	}
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		if _, ok := t.resolveFunction(ctx, fn); !ok {
			return
		}
	}
	for i := range pkg.Constants {
		resolveConstant(t, &pkg.Constants[i])
	}
	for i := range pkg.Imports {
		ast.Walk(t, &pkg.Imports[i])
	}
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		fn.SetDeterminedType(semtypes.NEVER)
		fn.Name.SetDeterminedType(semtypes.NEVER)
	}
	pkg.SetDeterminedType(semtypes.NEVER)
	for i := range pkg.CompUnits {
		pkg.CompUnits[i].SetDeterminedType(semtypes.NEVER)
	}
	for i := range pkg.GlobalVars {
		ast.Walk(t, &pkg.GlobalVars[i])
	}

	tctx := t.tyCtx
	for _, defn := range pkg.TypeDefinitions {
		if semtypes.IsEmpty(tctx, defn.DeterminedType) {
			t.semanticError(fmt.Sprintf("type definition %s is empty", defn.Name.GetValue()), defn.GetPosition())
		}
	}
}

func resolveBlockStatements(t typeResolver, chain *binding, stmts []ast.BLangStatement) (statementEffect, bool) {
	result := chain
	for i, each := range stmts {
		eachResult, ok := resolveStatement(t, result, each)
		if !ok {
			continue
		}
		if !eachResult.nonCompletion {
			result = eachResult.binding
		} else {
			rest := stmts[i+1:]
			if len(rest) > 0 {
				// These are unreachable nodes will be caught later by reachability analysis
				// we are doing type resolution here anyway to give error message to these statements
				resolveBlockStatements(t, chain, rest)
			}
			return statementEffect{result, true}, true
		}
	}
	return statementEffect{result, false}, true
}

func resolveStatement(t typeResolver, chain *binding, stmt ast.BLangStatement) (statementEffect, bool) {
	effect, ok := resolveStatementInner(t, chain, stmt)
	stmt.(ast.BLangNode).SetDeterminedType(semtypes.NEVER)
	return effect, ok
}

func resolveStatementInner(t typeResolver, chain *binding, stmt ast.BLangStatement) (statementEffect, bool) {
	switch s := stmt.(type) {
	case *ast.BLangSimpleVariableDef:
		variable := s.GetVariable().(*ast.BLangSimpleVariable)
		if !resolveSimpleVariable(t, chain, variable) {
			return defaultStmtEffect(chain), false
		}
		return defaultStmtEffect(chain), true
	// PR-TODO: extract assignment out
	case *ast.BLangAssignment:
		lhsTy, _, ok := resolveExpression(t, nil, s.GetVariable().(ast.BLangExpression), nil)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		if _, _, ok := resolveExpression(t, chain, s.GetExpression().(ast.BLangExpression), lhsTy); !ok {
			return defaultStmtEffect(chain), false
		}
		if expr, ok := s.GetVariable().(model.NodeWithSymbol); ok {
			return unnarrowSymbol(t, chain, expr.Symbol()), true
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangCompoundAssignment:
		lhsTy, _, ok := resolveExpression(t, nil, s.GetVariable().(ast.BLangExpression), nil)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		if _, _, ok := resolveExpression(t, chain, s.GetExpression().(ast.BLangExpression), lhsTy); !ok {
			return defaultStmtEffect(chain), false
		}
		if expr, ok := s.GetVariable().(model.NodeWithSymbol); ok {
			return unnarrowSymbol(t, chain, expr.Symbol()), true
		}
		return defaultStmtEffect(chain), true
	case *ast.BLangExpressionStmt:
		if _, _, ok := resolveExpression(t, chain, s.Expr, nil); !ok {
			return defaultStmtEffect(chain), false
		}
		return defaultStmtEffect(chain), true
	// PT-TODO: extract if while out
	case *ast.BLangIf:
		_, exprEffect, ok := resolveExpression(t, chain, s.Expr, semtypes.BOOLEAN)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		ifTrueEffect, ok := resolveBlockStatements(t, exprEffect.ifTrue, s.Body.Stmts)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		s.Body.SetDeterminedType(semtypes.NEVER)
		var ifFalseEffect statementEffect
		if s.ElseStmt != nil {
			ifFalseEffect, ok = resolveStatement(t, exprEffect.ifFalse, s.ElseStmt)
			if !ok {
				return defaultStmtEffect(chain), false
			}
		} else {
			ifFalseEffect = statementEffect{exprEffect.ifFalse, false}
		}
		return mergeStatementEffects(t, ifTrueEffect, ifFalseEffect), true
	case *ast.BLangWhile:
		_, exprEffect, ok := resolveExpression(t, chain, s.Expr, semtypes.BOOLEAN)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		bodyEffect, ok := resolveBlockStatements(t, exprEffect.ifTrue, s.Body.Stmts)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		s.Body.SetDeterminedType(semtypes.NEVER)
		resolveOnFailClause(t, chain, &s.OnFailClause)
		result := exprEffect.ifFalse
		if !bodyEffect.nonCompletion {
			result = mergeChains(t, result, bodyEffect.binding, semtypes.Union)
		}
		return statementEffect{result, false}, true
	case *ast.BLangReturn:
		if s.Expr != nil {
			if _, _, ok := resolveExpression(t, chain, s.Expr, t.expectedReturnType()); !ok {
				return defaultStmtEffect(chain), false
			}
		}
		return statementEffect{nil, true}, true
	case *ast.BLangBlockStmt:
		return resolveBlockStatements(t, chain, s.Stmts)
	case *ast.BLangForeach:
		if s.VariableDef != nil {
			variable := s.VariableDef.GetVariable().(*ast.BLangSimpleVariable)
			if !resolveSimpleVariable(t, chain, variable) {
				return defaultStmtEffect(chain), false
			}
			s.VariableDef.SetDeterminedType(semtypes.NEVER)
		}
		if s.Collection != nil {
			if _, _, ok := resolveExpression(t, chain, s.Collection, nil); !ok {
				return defaultStmtEffect(chain), false
			}
		}
		// Foreach loop can't create a conditional narrowing at the begining so at the end there shouldn't be
		// any narrowing.
		_, ok := resolveBlockStatements(t, chain, s.Body.Stmts)
		s.Body.SetDeterminedType(semtypes.NEVER)
		if s.OnFailClause != nil {
			resolveOnFailClause(t, chain, s.OnFailClause)
		}
		return defaultStmtEffect(chain), ok
	case *ast.BLangPanic:
		if _, _, ok := resolveExpression(t, chain, s.Expr, semtypes.ERROR); !ok {
			return defaultStmtEffect(chain), false
		}
		return statementEffect{nil, true}, true
	case *ast.BLangMatchStatement:
		return resolveMatchStatement(t, chain, s)
	case *ast.BLangBreak, *ast.BLangContinue:
		return defaultStmtEffect(chain), true
	default:
		t.internalError(fmt.Sprintf("unhandled statement type: %T", stmt), stmt.GetPosition())
		return defaultStmtEffect(chain), false
	}
}

func resolveOnFailClause(t typeResolver, chain *binding, clause *ast.BLangOnFailClause) {
	clause.SetDeterminedType(semtypes.NEVER)
	if clause.VariableDefinitionNode != nil {
		varDef := clause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
		variable := varDef.GetVariable().(*ast.BLangSimpleVariable)
		resolveSimpleVariable(t, chain, variable)
		varDef.SetDeterminedType(semtypes.NEVER)
	}
	if clause.Body != nil {
		resolveBlockStatements(t, chain, clause.Body.Stmts)
		clause.Body.SetDeterminedType(semtypes.NEVER)
	}
}

func (t *packageTypeResolver) resolveFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) (semtypes.SemType, bool) {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i := range fn.RequiredParams {
		ast.Walk(t, &fn.RequiredParams[i])
		paramTypes[i] = fn.RequiredParams[i].GetDeterminedType()
	}
	var restTy semtypes.SemType = semtypes.NEVER
	if fn.RestParam != nil {
		t.unimplemented("var args not supported", fn.RestParam.GetPosition())
		return nil, false
	}
	paramListDefn := semtypes.NewListDefinition()
	paramListTy := paramListDefn.DefineListTypeWrapped(t.typeEnv(), paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	if retTd := fn.GetReturnTypeDescriptor(); retTd != nil {
		var ok bool
		returnTy, ok = resolveBType(t, retTd.(ast.BType), 0)
		if !ok {
			return nil, false
		}
		ast.Walk(t, retTd.(ast.BLangNode))
	} else {
		returnTy = semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	fnType := functionDefn.Define(t.typeEnv(), paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.typeEnv(), false, false))

	// Update symbol type for the function
	updateSymbolType(t, fn, fnType)
	fnSymbol := ctx.GetSymbol(fn.Symbol()).(model.FunctionSymbol)
	sig := fnSymbol.Signature()
	sig.ParamTypes = paramTypes
	sig.ReturnType = returnTy
	sig.RestParamType = restTy
	fnSymbol.SetSignature(sig)

	return fnType, true
}

func visitTypeData(t typeResolver, typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return t
	}
	ty, ok := resolveBType(t, typeData.TypeDescriptor.(ast.BType), 0)
	if !ok {
		return nil
	}
	typeData.Type = ty

	// Update symbol type if the type descriptor has a symbol
	if tdNode, ok := typeData.TypeDescriptor.(ast.BLangNode); ok {
		updateSymbolType(t, tdNode, ty)
	}

	return t
}

func visit(t typeResolver, node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.BLangConstant:
		resolveConstant(t, n)
		return nil
	case *ast.BLangSimpleVariable:
		resolveSimpleVariable(t, nil, node.(*ast.BLangSimpleVariable))
	case ast.BType:
		resolveBType(t, node.(ast.BType), 0)
	case *ast.BLangLiteral:
		resolveLiteral(t, n, nil)
		return nil
	case *ast.BLangNumericLiteral:
		resolveNumericLiteral(t, n, nil)
		return nil
	case *ast.BLangTypeDefinition:
		resolveTypeDefinition(t, n, 0)
		return nil
	case *ast.BLangMatchStatement:
		resolveMatchStatement(t, nil, n)
		return nil
	case ast.BLangExpression:
		if _, _, ok := resolveExpression(t, nil, n, nil); !ok {
			return nil
		}
	default:
		// Non-expression nodes with no specific handling: mark as NEVER and continue traversal
	}
	// Set DeterminedType to NEVER as fallback for nodes that didn't get a type assigned.
	if node.GetDeterminedType() == nil {
		node.SetDeterminedType(semtypes.NEVER)
	}
	return t
}

func resolveTypeDefinition(t typeResolver, defn model.TypeDefinition, depth int) (semtypes.SemType, bool) {
	if defn.GetDeterminedType() != nil {
		return defn.GetDeterminedType(), true
	}
	// Walk Name identifier to ensure it gets DeterminedType set
	if defn.GetName() != nil {
		ast.Walk(t, defn.GetName().(ast.BLangNode))
	}
	if depth == defn.GetCycleDepth() {
		t.semanticError(fmt.Sprintf("invalid cycle detected for type definition %s", defn.GetName().GetValue()), defn.GetPosition())
		return nil, false
	}
	defn.SetCycleDepth(depth)
	semType, ok := resolveBType(t, defn.GetTypeData().TypeDescriptor.(ast.BType), depth)
	if !ok {
		return nil, false
	}
	if defn.GetDeterminedType() == nil {
		defn.SetDeterminedType(semType)
		t.setSymbolType(defn.Symbol(), semType)
		defn.SetCycleDepth(-1)
		typeData := defn.GetTypeData()
		typeData.Type = semType
		defn.SetTypeData(typeData)
		return semType, true
	} else {
		// This can happen with recursion
		// We use the first definition we produced
		// and throw away the others
		return defn.GetDeterminedType(), true
	}
}

func resolveLiteral(t typeResolver, n *ast.BLangLiteral, expectedType semtypes.SemType) bool {
	bType := n.GetValueType()
	var ty semtypes.SemType

	switch bType.BTypeGetTag() {
	case model.TypeTags_INT, model.TypeTags_BYTE, model.TypeTags_FLOAT, model.TypeTags_DECIMAL:
		var ok bool
		ty, ok = resolveNumericLiteralValue(t, n, expectedType)
		if !ok {
			return false
		}
	case model.TypeTags_BOOLEAN:
		value := n.GetValue().(bool)
		ty = semtypes.BooleanConst(value)
	case model.TypeTags_STRING:
		value := n.GetValue().(string)
		ty = semtypes.StringConst(value)
	case model.TypeTags_NIL:
		ty = semtypes.NIL
	default:
		t.unimplemented("unsupported literal type", n.GetPosition())
		return false
	}

	setExpectedType(n, ty)

	// Update symbol type if this literal has a symbol
	updateSymbolType(t, n, ty)
	return true
}

func hasFloatTypeSuffix(s string) bool {
	if len(s) == 0 {
		return false
	}
	last := s[len(s)-1]
	return last == 'f' || last == 'F'
}

func determineCandidatesFromLiteral(t typeResolver, n *ast.BLangLiteral) semtypes.SemType {
	switch n.GetValueType().BTypeGetTag() {
	case model.TypeTags_INT, model.TypeTags_BYTE:
		return semtypes.NUMBER
	case model.TypeTags_FLOAT:
		if hasFloatTypeSuffix(n.OriginalValue) {
			return semtypes.FLOAT
		}
		if balCommon.HasHexIndicator(n.OriginalValue) {
			return semtypes.FLOAT
		}
		return semtypes.Union(semtypes.FLOAT, semtypes.DECIMAL)
	case model.TypeTags_DECIMAL:
		return semtypes.DECIMAL
	default:
		t.internalError(fmt.Sprintf("unexpected type tag %v for numeric literal", n.GetValueType().BTypeGetTag()), n.GetPosition())
		return semtypes.NEVER
	}
}

func determineCandidatesFromNumericLiteral(t typeResolver, n *ast.BLangNumericLiteral) semtypes.SemType {
	switch n.Kind {
	case model.NodeKind_INTEGER_LITERAL:
		return semtypes.NUMBER
	case model.NodeKind_DECIMAL_FLOATING_POINT_LITERAL:
		if hasFloatTypeSuffix(n.OriginalValue) {
			return semtypes.FLOAT
		}
		if balCommon.IsDecimalDiscriminated(n.OriginalValue) {
			return semtypes.DECIMAL
		}
		return semtypes.Union(semtypes.FLOAT, semtypes.DECIMAL)
	case model.NodeKind_HEX_FLOATING_POINT_LITERAL:
		return semtypes.FLOAT
	default:
		t.internalError(fmt.Sprintf("unexpected numeric literal kind: %v", n.Kind), n.GetPosition())
		return semtypes.NEVER
	}
}

func narrowCandidates(candidates, expectedType semtypes.SemType) semtypes.SemType {
	if expectedType == nil {
		return candidates
	}
	narrowed := semtypes.Intersect(candidates, expectedType)
	if !semtypes.IsNever(narrowed) {
		return narrowed
	}
	return candidates
}

func pickNumericType(t typeResolver, n *ast.BLangLiteral, candidates semtypes.SemType) (semtypes.SemType, bool) {
	switch {
	case semtypes.ContainsBasicType(candidates, semtypes.INT):
		return resolveAsInt(t, n)
	case semtypes.ContainsBasicType(candidates, semtypes.FLOAT):
		return resolveAsFloat(t, n)
	case semtypes.ContainsBasicType(candidates, semtypes.DECIMAL):
		return resolveAsDecimal(t, n)
	default:
		t.semanticError("no valid candidate to resolve numeric literal", n.GetPosition())
		return nil, false
	}
}

func resolveAsInt(t typeResolver, n *ast.BLangLiteral) (semtypes.SemType, bool) {
	var intVal int64
	switch v := n.GetValue().(type) {
	case int64:
		intVal = v
	case float64:
		intVal = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			t.syntaxError(fmt.Sprintf("invalid int literal: %s", v), n.GetPosition())
			return nil, false
		}
		intVal = parsed
	default:
		t.internalError(fmt.Sprintf("unexpected int literal value type: %T", n.GetValue()), n.GetPosition())
		return nil, false
	}
	n.SetValue(intVal)
	return semtypes.IntConst(intVal), true
}

func resolveAsFloat(t typeResolver, n *ast.BLangLiteral) (semtypes.SemType, bool) {
	var floatVal float64
	switch v := n.GetValue().(type) {
	case string:
		parsed, ok := parseFloatValue(t, v, n.GetPosition())
		if !ok {
			return nil, false
		}
		floatVal = parsed
	case float64:
		floatVal = v
	case int64:
		floatVal = float64(v)
	default:
		t.internalError(fmt.Sprintf("unexpected float literal value type: %T", v), n.GetPosition())
		return nil, false
	}
	n.SetValue(floatVal)
	return semtypes.FloatConst(floatVal), true
}

func resolveAsDecimal(t typeResolver, n *ast.BLangLiteral) (semtypes.SemType, bool) {
	var ratVal *big.Rat
	switch v := n.GetValue().(type) {
	case string:
		parsed, ok := parseDecimalValue(t, stripFloatingPointTypeSuffix(v), n.GetPosition())
		if !ok {
			return nil, false
		}
		ratVal = parsed
	case *big.Rat:
		ratVal = v
	case int64:
		ratVal = new(big.Rat).SetInt64(v)
	case float64:
		ratVal = new(big.Rat).SetFloat64(v)
	default:
		t.internalError(fmt.Sprintf("unexpected decimal literal value type: %T", v), n.GetPosition())
		return nil, false
	}
	n.SetValue(ratVal)
	return semtypes.DecimalConst(*ratVal), true
}

func resolveNumericLiteralValue(t typeResolver, n *ast.BLangLiteral, expectedType semtypes.SemType) (semtypes.SemType, bool) {
	candidates := determineCandidatesFromLiteral(t, n)
	candidates = narrowCandidates(candidates, expectedType)
	return pickNumericType(t, n, candidates)
}

// stripFloatingPointTypeSuffix removes the f/F/d/D type suffix from a floating point literal string
func stripFloatingPointTypeSuffix(s string) string {
	last := s[len(s)-1]
	if last == 'f' || last == 'F' || last == 'd' || last == 'D' {
		return s[:len(s)-1]
	}
	return s
}

func parseFloatValue(t typeResolver, strValue string, pos diagnostics.Location) (float64, bool) {
	strValue = strings.TrimRight(strValue, "fF")
	f, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		t.syntaxError(fmt.Sprintf("invalid float literal: %s", strValue), pos)
		return 0, false
	}
	return f, true
}

func parseDecimalValue(t typeResolver, strValue string, pos diagnostics.Location) (*big.Rat, bool) {
	r := new(big.Rat)
	if _, ok := r.SetString(strValue); !ok {
		t.syntaxError(fmt.Sprintf("invalid decimal literal: %s", strValue), pos)
		return big.NewRat(0, 1), false
	}
	return r, true
}

func resolveNumericLiteral(t typeResolver, n *ast.BLangNumericLiteral, expectedType semtypes.SemType) bool {
	candidates := determineCandidatesFromNumericLiteral(t, n)
	candidates = narrowCandidates(candidates, expectedType)

	ty, ok := pickNumericType(t, &n.BLangLiteral, candidates)
	if !ok {
		return false
	}

	setExpectedType(n, ty)
	updateSymbolType(t, n, ty)
	return true
}

// updateSymbolType updates the symbol's type if the node has an associated symbol.
func updateSymbolType(t typeResolver, node ast.BLangNode, ty semtypes.SemType) {
	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		t.setSymbolType(symbol, ty)
	}
}

func lookupSymbol(chain *binding, ref model.SymbolRef) model.SymbolRef {
	if chain == nil {
		return ref
	}
	narrowedRef, isNarrowed := lookupBinding(chain, ref)
	if isNarrowed {
		return narrowedRef
	}
	return ref
}

func resolveSimpleVariable(t typeResolver, chain *binding, node *ast.BLangSimpleVariable) bool {
	node.Name.SetDeterminedType(semtypes.NEVER)
	typeNode := node.TypeNode()
	if typeNode == nil {
		if node.Expr != nil {
			exprTy, _, ok := resolveExpression(t, chain, node.Expr.(ast.BLangExpression), nil)
			if !ok {
				return false
			}
			setExpectedType(node, exprTy)
			updateSymbolType(t, node, exprTy)
		}
		return true
	}

	semType, ok := resolveBType(t, typeNode, 0)
	if !ok {
		setExpectedType(node, semtypes.NEVER)
		updateSymbolType(t, node, semtypes.NEVER)
		return false
	}

	setExpectedType(node, semType)
	updateSymbolType(t, node, semType)

	if node.Expr != nil {
		if _, _, ok := resolveExpression(t, chain, node.Expr.(ast.BLangExpression), semType); !ok {
			return false
		}
	}

	return true
}

func resolveExpression(t typeResolver, chain *binding, expr ast.BLangExpression, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	// Check if already resolved
	if ty := expr.GetDeterminedType(); ty != nil {
		return ty, defaultExpressionEffect(chain), true
	}

	ty, effect, ok := resolveExpressionInner(t, chain, expr, expectedType)
	if !ok {
		// Mark failed expressions so ast.Walk won't re-process them
		setExpectedType(expr, semtypes.NEVER)
		return nil, expressionEffect{}, false
	}
	if singletonEffect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return ty, singletonEffect, true
	}
	return ty, effect, ok
}

func resolveExpressionInner(t typeResolver, chain *binding, expr ast.BLangExpression, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	switch e := expr.(type) {
	case *ast.BLangLiteral:
		if ok := resolveLiteral(t, e, expectedType); !ok {
			return nil, expressionEffect{}, false
		}
		return e.GetDeterminedType(), defaultExpressionEffect(chain), true
	case *ast.BLangNumericLiteral:
		if ok := resolveNumericLiteral(t, e, expectedType); !ok {
			return nil, expressionEffect{}, false
		}
		return e.GetDeterminedType(), defaultExpressionEffect(chain), true
	case *ast.BLangSimpleVarRef:
		return resolveSimpleVarRef(t, chain, e)
	case *ast.BLangLocalVarRef:
		return resolveLocalVarRef(t, chain, e)
	case *ast.BLangConstRef:
		return resolveConstRef(t, chain, e)
	case *ast.BLangBinaryExpr:
		return resolveBinaryExpr(t, chain, e)
	case *ast.BLangUnaryExpr:
		return resolveUnaryExpr(t, chain, e)
	case *ast.BLangInvocation:
		return resolveInvocation(t, chain, e)
	case *ast.BLangIndexBasedAccess:
		return resolveIndexBasedAccess(t, chain, e)
	case *ast.BLangFieldBaseAccess:
		return resolveFieldBaseAccess(t, chain, e)
	case *ast.BLangListConstructorExpr:
		return resolveListConstructorExpr(t, chain, e, expectedType)
	case *ast.BLangMappingConstructorExpr:
		return resolveMappingConstructorExpr(t, chain, e, expectedType)
	case *ast.BLangErrorConstructorExpr:
		return resolveErrorConstructorExpr(t, chain, e, expectedType)
	case *ast.BLangGroupExpr:
		return resolveGroupExpr(t, chain, e, expectedType)
	case *ast.BLangQueryExpr:
		return resolveQueryExpr(t, chain, e)
	case *ast.BLangWildCardBindingPattern:
		ty := semtypes.ANY
		setExpectedType(e, ty)
		return ty, defaultExpressionEffect(chain), true
	case *ast.BLangTypeConversionExpr:
		return resolveTypeConversionExpr(t, chain, e)
	case *ast.BLangTypeTestExpr:
		return resolveTypeTestExpr(t, chain, e)
	case *ast.BLangCheckedExpr:
		return resolveCheckedExpr(t, chain, e)
	case *ast.BLangCheckPanickedExpr:
		return resolveCheckedExpr(t, chain, &e.BLangCheckedExpr)
	case *ast.BLangTrapExpr:
		return resolveTrapExpr(t, chain, e)
	case *ast.BLangNamedArgsExpression:
		ty, effect, ok := resolveExpression(t, chain, e.Expr, expectedType)
		if !ok {
			return nil, expressionEffect{}, false
		}
		setExpectedType(e, ty)
		e.Name.SetDeterminedType(semtypes.NEVER)
		return ty, effect, true
	default:
		t.internalError(fmt.Sprintf("unsupported expression type: %T", expr), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
}

func resolveTypeTestExpr(t typeResolver, chain *binding, e *ast.BLangTypeTestExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := resolveExpression(t, chain, e.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	ast.WalkTypeData(t, &e.Type)
	testedTy := e.Type.Type

	var resultTy semtypes.SemType
	if semtypes.IsSubtype(t.typeContext(), exprTy, testedTy) {
		resultTy = semtypes.BooleanConst(!e.IsNegation())
	} else if semtypes.IsEmpty(t.typeContext(), semtypes.Intersect(exprTy, testedTy)) {
		resultTy = semtypes.BooleanConst(e.IsNegation())
	} else {
		resultTy = semtypes.BOOLEAN
	}

	setExpectedType(e, resultTy)

	ref, isVarRef := varRefExp(chain, &e.Expr)
	if !isVarRef {
		return resultTy, defaultExpressionEffect(chain), true
	}
	tx := t.symbolType(ref)
	ref = t.unnarrowedSymbol(ref)
	testTy := e.Type.Type
	trueTy := semtypes.Intersect(tx, testTy)
	trueSym := narrowSymbol(t, ref, trueTy)
	trueChain := &binding{ref: ref, narrowedSymbol: trueSym, prev: chain}
	falseTy := semtypes.Diff(tx, testTy)
	falseSym := narrowSymbol(t, ref, falseTy)
	falseChain := &binding{ref: ref, narrowedSymbol: falseSym, prev: chain}
	if e.IsNegation() {
		return resultTy, expressionEffect{ifTrue: falseChain, ifFalse: trueChain}, true
	}
	return resultTy, expressionEffect{ifTrue: trueChain, ifFalse: falseChain}, true
}

func resolveTrapExpr(t typeResolver, chain *binding, e *ast.BLangTrapExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := resolveExpression(t, chain, e.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	resultTy := semtypes.Union(exprTy, semtypes.ERROR)
	setExpectedType(e, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func resolveCheckedExpr(t typeResolver, chain *binding, e *ast.BLangCheckedExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, _, ok := resolveExpression(t, chain, e.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	errorIntersection := semtypes.Intersect(exprTy, semtypes.ERROR)
	if semtypes.IsEmpty(t.typeContext(), errorIntersection) {
		e.IsRedundantChecking = true
	}
	resultTy := semtypes.Diff(exprTy, semtypes.ERROR)
	setExpectedType(e, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func resolveMappingConstructorExpr(t typeResolver, chain *binding, e *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	if expectedType != nil {
		return resolveMappingConstructorWithExpectedType(t, chain, e, expectedType)
	}
	return resolveMappingConstructorBottomUp(t, chain, e)
}

func resolveMappingConstructorBottomUp(t typeResolver, chain *binding, e *ast.BLangMappingConstructorExpr) (semtypes.SemType, expressionEffect, bool) {
	fields := make([]semtypes.Field, len(e.Fields))
	for i, f := range e.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		valueTy, _, ok := resolveExpression(t, chain, kv.ValueExpr, nil)
		if !ok {
			return nil, expressionEffect{}, false
		}
		var broadTy semtypes.SemType
		if semtypes.SingleShape(valueTy).IsEmpty() {
			broadTy = valueTy
		} else {
			broadTy = semtypes.WidenToBasicTypes(valueTy)
		}
		var keyName string
		switch keyExpr := kv.Key.Expr.(type) {
		case *ast.BLangLiteral:
			keyName = keyExpr.GetOriginalValue()
			resolveLiteral(t, keyExpr, nil)
		case ast.BNodeWithSymbol:
			t.setSymbolType(keyExpr.Symbol(), valueTy)
			keyName = t.symbolName(keyExpr.Symbol())
			if e, ok := keyExpr.(ast.BLangExpression); ok {
				setExpectedType(e, valueTy)
			}
			if ref, ok := keyExpr.(*ast.BLangSimpleVarRef); ok {
				setVarRefIdentifierTypes(ref)
			}
		}
		kv.Key.SetDeterminedType(semtypes.NEVER)
		kv.SetDeterminedType(semtypes.NEVER)
		fields[i] = semtypes.FieldFrom(keyName, broadTy, false, false)
	}
	md := semtypes.NewMappingDefinition()
	mapTy := md.DefineMappingTypeWrapped(t.typeEnv(), fields, semtypes.NEVER)
	setExpectedType(e, mapTy)
	mat := semtypes.ToMappingAtomicType(t.typeContext(), mapTy)
	e.AtomicType = *mat
	return mapTy, defaultExpressionEffect(chain), true
}

func resolveMappingConstructorWithExpectedType(t typeResolver, chain *binding, e *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	for _, f := range e.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		if _, _, ok := resolveExpression(t, chain, kv.ValueExpr, nil); !ok {
			return nil, expressionEffect{}, false
		}
		resolveMappingKey(t, kv)
	}

	resultType, mat, ok := selectMappingInherentType(t, e, expectedType)
	if !ok {
		return resolveMappingConstructorBottomUp(t, chain, e)
	}

	for _, f := range e.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		keyName := recordKeyName(kv.Key)
		requiredType := mat.FieldInnerVal(keyName)
		kv.ValueExpr.SetDeterminedType(nil)
		if _, _, ok := resolveExpression(t, chain, kv.ValueExpr, requiredType); !ok {
			return nil, expressionEffect{}, false
		}
	}

	e.AtomicType = mat
	setExpectedType(e, resultType)
	return resultType, defaultExpressionEffect(chain), true
}

func resolveMappingKey(t typeResolver, kv *ast.BLangMappingKeyValueField) {
	switch keyExpr := kv.Key.Expr.(type) {
	case *ast.BLangLiteral:
		resolveLiteral(t, keyExpr, nil)
	case ast.BNodeWithSymbol:
		valueTy := kv.ValueExpr.GetDeterminedType()
		t.setSymbolType(keyExpr.Symbol(), valueTy)
		if e, ok := keyExpr.(ast.BLangExpression); ok {
			setExpectedType(e, valueTy)
		}
		if ref, ok := keyExpr.(*ast.BLangSimpleVarRef); ok {
			setVarRefIdentifierTypes(ref)
		}
	}
	kv.Key.SetDeterminedType(semtypes.NEVER)
	kv.SetDeterminedType(semtypes.NEVER)
}

func selectMappingInherentType(t typeResolver, expr *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, semtypes.MappingAtomicType, bool) {
	expectedMappingType := semtypes.Intersect(expectedType, semtypes.MAPPING)
	tc := t.typeContext()
	if semtypes.IsEmpty(tc, expectedMappingType) {
		t.semanticError("mapping type not found in expected type", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}, false
	}
	mat := semtypes.ToMappingAtomicType(tc, expectedMappingType)
	if mat != nil {
		return expectedMappingType, *mat, true
	}
	alts := semtypes.MappingAlternatives(tc, expectedType)
	var validAlts []semtypes.MappingAlternative

	fields := make([]semtypes.MappingFieldInfo, len(expr.Fields))
	for i, f := range expr.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		fields[i] = semtypes.MappingFieldInfo{Name: recordKeyName(kv.Key), Ty: kv.ValueExpr.GetDeterminedType()}
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })

	for _, alt := range alts {
		if semtypes.MappingAlternativeAllowsFields(tc, alt, fields) {
			validAlts = append(validAlts, alt)
		}
	}
	if len(validAlts) == 0 {
		t.semanticError("no applicable inherent type for mapping constructor", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}, false
	}
	if len(validAlts) > 1 {
		t.semanticError("ambiguous inherent type for mapping constructor", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}, false
	}

	selectedSemType := validAlts[0].SemType
	mat = semtypes.ToMappingAtomicType(tc, selectedSemType)
	if mat == nil {
		t.semanticError("applicable type for mapping constructor is not atomic", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}, false
	}

	return selectedSemType, *mat, true
}

func resolveTypeConversionExpr(t typeResolver, chain *binding, e *ast.BLangTypeConversionExpr) (semtypes.SemType, expressionEffect, bool) {
	expectedType, ok := resolveBType(t, e.TypeDescriptor.(ast.BType), 0)
	if !ok {
		return nil, expressionEffect{}, false
	}
	_, _, ok = resolveExpression(t, chain, e.Expression, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	setExpectedType(e, expectedType)
	return expectedType, defaultExpressionEffect(chain), true
}

// Helper functions for expression type checking

func setVarRefIdentifierTypes(ref *ast.BLangSimpleVarRef) {
	if ref.PkgAlias != nil {
		ref.PkgAlias.SetDeterminedType(semtypes.NEVER)
	}
	if ref.VariableName != nil {
		ref.VariableName.SetDeterminedType(semtypes.NEVER)
	}
}

type opExpr interface {
	GetOperatorKind() model.OperatorKind
}

func isEqualityExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_EQUAL, model.OperatorKind_EQUALS, model.OperatorKind_NOT_EQUAL, model.OperatorKind_REF_EQUAL, model.OperatorKind_REF_NOT_EQUAL:
		return true
	default:
		return false
	}
}

func isMultiplicativeExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_MUL, model.OperatorKind_DIV, model.OperatorKind_MOD:
		return true
	default:
		return false
	}
}

func isRangeExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_CLOSED_RANGE, model.OperatorKind_HALF_OPEN_RANGE:
		return true
	default:
		return false
	}
}

func isBitWiseExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_BITWISE_AND, model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
		return true
	default:
		return false
	}
}

func isShiftExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_BITWISE_LEFT_SHIFT,
		model.OperatorKind_BITWISE_RIGHT_SHIFT,
		model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
		return true
	default:
		return false
	}
}

func isRelationalExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_LESS_THAN, model.OperatorKind_LESS_EQUAL, model.OperatorKind_GREATER_THAN, model.OperatorKind_GREATER_EQUAL:
		return true
	default:
		return false
	}
}

func isAdditiveExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB:
		return true
	default:
		return false
	}
}

func isLogicalExpression(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_AND, model.OperatorKind_OR:
		return true
	default:
		return false
	}
}

func isNumericType(ty semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(ty, semtypes.NUMBER)
}

func resolveGroupExpr(t typeResolver, chain *binding, expr *ast.BLangGroupExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	innerTy, effect, ok := resolveExpression(t, chain, expr.Expression, expectedType)
	if !ok {
		return nil, expressionEffect{}, false
	}
	setExpectedType(expr, innerTy)
	return innerTy, effect, true
}

func resolveQueryExpr(t typeResolver, chain *binding, expr *ast.BLangQueryExpr) (semtypes.SemType, expressionEffect, bool) {
	if len(expr.QueryClauseList) < 2 {
		t.semanticError("query expression requires from and select clauses", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	fromClause, ok := expr.QueryClauseList[0].(*ast.BLangFromClause)
	if !ok {
		t.semanticError("query expression must start with a from clause", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	fromClause.SetDeterminedType(semtypes.NEVER)

	selectClause, ok := expr.QueryClauseList[len(expr.QueryClauseList)-1].(*ast.BLangSelectClause)
	if !ok {
		t.semanticError("query expression requires a select clause", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	selectClause.SetDeterminedType(semtypes.NEVER)

	collectionTy, _, ok := resolveExpression(t, chain, fromClause.Collection, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	var elementTy semtypes.SemType
	switch {
	case semtypes.IsSubtypeSimple(collectionTy, semtypes.LIST):
		memberTypes := semtypes.ListAllMemberTypesInner(t.typeContext(), collectionTy)
		var result semtypes.SemType = semtypes.NEVER
		for _, each := range memberTypes.SemTypes {
			result = semtypes.Union(result, each)
		}
		elementTy = result
	default:
		t.unimplemented("query from-clause currently supports only list collections", fromClause.GetPosition())
		return nil, expressionEffect{}, false
	}

	if fromClause.VariableDefinitionNode != nil {
		varDef, ok := fromClause.VariableDefinitionNode.(*ast.BLangSimpleVariableDef)
		if !ok || varDef.Var == nil {
			t.unimplemented("only simple variable bindings are supported in from clause", fromClause.GetPosition())
			return nil, expressionEffect{}, false
		}
		varDef.SetDeterminedType(semtypes.NEVER)

		var variableTy semtypes.SemType = elementTy
		if !fromClause.IsDeclaredWithVarFlag && varDef.Var.TypeNode() != nil {
			variableTy, ok = resolveBType(t, varDef.Var.TypeNode(), 0)
			if !ok {
				return nil, expressionEffect{}, false
			}
			if !semtypes.IsSubtype(t.typeContext(), elementTy, variableTy) {
				t.semanticError("from-clause variable type is incompatible with collection member type",
					varDef.GetPosition())
				return nil, expressionEffect{}, false
			}
		}

		if varDef.Var.Name != nil {
			varDef.Var.Name.SetDeterminedType(semtypes.NEVER)
		}
		varDef.Var.SetDeterminedType(semtypes.NEVER)
		updateSymbolType(t, varDef.Var, variableTy)
	}

	queryChain, ok := resolveQueryIntermediateClauses(t, chain, expr)
	if !ok {
		return nil, expressionEffect{}, false
	}

	selectTy, _, ok := resolveExpression(t, queryChain, selectClause.Expression, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	ld := semtypes.NewListDefinition()
	queryTy := ld.DefineListTypeWrappedWithEnvSemType(t.typeEnv(), selectTy)
	setExpectedType(expr, queryTy)
	return queryTy, defaultExpressionEffect(chain), true
}

func resolveQueryIntermediateClauses(t typeResolver, chain *binding, queryExpr *ast.BLangQueryExpr) (*binding, bool) {
	currentChain := chain
	for i := 1; i < len(queryExpr.QueryClauseList)-1; i++ {
		switch clause := queryExpr.QueryClauseList[i].(type) {
		case *ast.BLangLetClause:
			clause.SetDeterminedType(semtypes.NEVER)
			for _, variableDef := range clause.LetVarDeclarations {
				varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
				if !ok || varDef.Var == nil {
					t.unimplemented("only simple variable declarations are supported in let clause",
						clause.GetPosition())
					return nil, false
				}
				varDef.SetDeterminedType(semtypes.NEVER)
				if varDef.Var.Expr == nil {
					t.semanticError("let-clause variable declaration requires an initializer",
						varDef.GetPosition())
					return nil, false
				}
				initTy, _, ok := resolveExpression(t, currentChain, varDef.Var.Expr.(ast.BLangExpression), nil)
				if !ok {
					return nil, false
				}
				var variableTy semtypes.SemType = initTy
				if !varDef.Var.GetIsDeclaredWithVar() && varDef.Var.TypeNode() != nil {
					variableTy, ok = resolveBType(t, varDef.Var.TypeNode(), 0)
					if !ok {
						return nil, false
					}
					if !semtypes.IsSubtype(t.typeContext(), initTy, variableTy) {
						t.semanticError("let-clause variable type is incompatible with initializer expression",
							varDef.GetPosition())
						return nil, false
					}
				}
				if varDef.Var.Name != nil {
					varDef.Var.Name.SetDeterminedType(semtypes.NEVER)
				}
				varDef.Var.SetDeterminedType(semtypes.NEVER)
				updateSymbolType(t, varDef.Var, variableTy)
			}
		case *ast.BLangWhereClause:
			clause.SetDeterminedType(semtypes.NEVER)
			whereTy, effect, ok := resolveExpression(t, currentChain, clause.Expression, semtypes.BOOLEAN)
			if !ok {
				return nil, false
			}
			if !semtypes.IsSubtypeSimple(whereTy, semtypes.BOOLEAN) {
				t.semanticError("where-clause expression must be boolean", clause.GetPosition())
				return nil, false
			}
			currentChain = effect.ifTrue
		default:
			t.unimplemented("only let + where clauses are supported as intermediate query clauses", clause.GetPosition())
			return nil, false
		}
	}
	return currentChain, true
}

func resolveSimpleVarRef(t typeResolver, chain *binding, expr *ast.BLangSimpleVarRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.symbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(expr)
	return ty, defaultExpressionEffect(chain), true
}

func resolveLocalVarRef(t typeResolver, chain *binding, expr *ast.BLangLocalVarRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.symbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(&expr.BLangSimpleVarRef)
	return ty, defaultExpressionEffect(chain), true
}

func resolveConstRef(t typeResolver, chain *binding, expr *ast.BLangConstRef) (semtypes.SemType, expressionEffect, bool) {
	sym, isNarrowed := lookupBinding(chain, expr.Symbol())
	if isNarrowed {
		expr.SetSymbol(sym)
	}
	ty := t.symbolType(sym)
	setExpectedType(expr, ty)
	setVarRefIdentifierTypes(&expr.BLangSimpleVarRef)
	return ty, defaultExpressionEffect(chain), true
}

func resolveListConstructorExpr(t typeResolver, chain *binding, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	if expectedType != nil {
		return resolveListConstructorWithExpectedType(t, chain, expr, expectedType)
	}
	return resolveListConstructorBottomUp(t, chain, expr)
}

func resolveListConstructorBottomUp(t typeResolver, chain *binding, expr *ast.BLangListConstructorExpr) (semtypes.SemType, expressionEffect, bool) {
	memberTypes := make([]semtypes.SemType, len(expr.Exprs))
	for i, memberExpr := range expr.Exprs {
		memberTy, _, ok := resolveExpression(t, chain, memberExpr, nil)
		if !ok {
			return nil, expressionEffect{}, false
		}
		var broadTy semtypes.SemType
		if semtypes.SingleShape(memberTy).IsEmpty() {
			broadTy = memberTy
		} else {
			broadTy = semtypes.WidenToBasicTypes(memberTy)
		}
		memberTypes[i] = broadTy
	}

	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(t.typeEnv(), memberTypes, len(memberTypes), semtypes.NEVER, semtypes.CellMutability_CELL_MUT_LIMITED)

	setExpectedType(expr, listTy)
	lat := semtypes.ToListAtomicType(t.typeContext(), listTy)
	expr.AtomicType = *lat

	return listTy, defaultExpressionEffect(chain), true
}

func resolveListConstructorWithExpectedType(t typeResolver, chain *binding, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	for _, memberExpr := range expr.Exprs {
		if _, _, ok := resolveExpression(t, chain, memberExpr, nil); !ok {
			return nil, expressionEffect{}, false
		}
	}

	resultType, lat, ok := selectListInherentType(t, expr, expectedType)
	if !ok {
		return resolveListConstructorBottomUp(t, chain, expr)
	}

	for i, memberExpr := range expr.Exprs {
		requiredType := lat.MemberAtInnerVal(i)
		if semtypes.IsNever(requiredType) {
			t.semanticError("too many members in list constructor", expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		memberExpr.SetDeterminedType(nil)
		if _, _, ok := resolveExpression(t, chain, memberExpr, requiredType); !ok {
			return nil, expressionEffect{}, false
		}
	}

	expr.AtomicType = lat
	setExpectedType(expr, resultType)
	return resultType, defaultExpressionEffect(chain), true
}

func selectListInherentType(t typeResolver, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, semtypes.ListAtomicType, bool) {
	expectedListType := semtypes.Intersect(expectedType, semtypes.LIST)
	tc := t.typeContext()
	if semtypes.IsEmpty(tc, expectedListType) {
		t.semanticError("list type not found in expected type", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}, false
	}
	lat := semtypes.ToListAtomicType(tc, expectedListType)
	if lat != nil {
		return expectedListType, *lat, true
	}

	alts := semtypes.ListAlternatives(tc, expectedListType)

	members := make([]semtypes.ListMemberInfo, len(expr.Exprs))
	for i, expr := range expr.Exprs {
		members[i] = semtypes.ListMemberInfo{Index: i, ValType: expr.GetDeterminedType()}
	}

	var validAlts []semtypes.ListAlternative
	for _, alt := range alts {
		if semtypes.ListAlternativeAllowsMembers(tc, alt, members) {
			validAlts = append(validAlts, alt)
		}
	}

	if len(validAlts) == 0 {
		t.semanticError("no applicable inherent type for list constructor", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}, false
	}
	if len(validAlts) > 1 {
		t.semanticError("ambiguous inherent type for list constructor", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}, false
	}

	selectedSemType := validAlts[0].SemType
	lat = semtypes.ToListAtomicType(tc, selectedSemType)
	if lat == nil {
		t.semanticError("applicable type for list constructor is not atomic", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}, false
	}

	return selectedSemType, *lat, true
}

func resolveErrorConstructorExpr(t typeResolver, chain *binding, expr *ast.BLangErrorConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, expressionEffect, bool) {
	var errorTy semtypes.SemType

	if expr.ErrorTypeRef != nil {
		refTy, ok := resolveBType(t, expr.ErrorTypeRef, 0)
		if !ok {
			return nil, expressionEffect{}, false
		}
		if !semtypes.IsSubtypeSimple(refTy, semtypes.ERROR) {
			t.semanticError("error type parameter must be a subtype of error", expr.ErrorTypeRef.GetPosition())
			return nil, expressionEffect{}, false
		} else {
			errorTy = refTy
		}
	} else {
		errorTy = semtypes.ERROR
	}

	if expectedType != nil && semtypes.IsSameType(t.typeContext(), errorTy, semtypes.ERROR) {
		errorPart := semtypes.Intersect(expectedType, semtypes.ERROR)
		if !semtypes.IsEmpty(t.typeContext(), errorPart) {
			errorTy = errorPart
		}
	}

	setExpectedType(expr, errorTy)

	for _, arg := range expr.PositionalArgs {
		if _, _, ok := resolveExpression(t, chain, arg, nil); !ok {
			return nil, expressionEffect{}, false
		}
	}
	for _, arg := range expr.NamedArgs {
		if _, _, ok := resolveExpression(t, chain, arg, nil); !ok {
			return nil, expressionEffect{}, false
		}
	}
	return errorTy, defaultExpressionEffect(chain), true
}

func resolveUnaryExpr(t typeResolver, chain *binding, expr *ast.BLangUnaryExpr) (semtypes.SemType, expressionEffect, bool) {
	exprTy, innerEffect, ok := resolveExpression(t, chain, expr.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType
	switch expr.GetOperatorKind() {
	case model.OperatorKind_SUB:
		if numLit, ok := expr.Expr.(*ast.BLangNumericLiteral); ok {
			resultValue := numLit.Value.(int64) * -1
			resultTy = semtypes.IntConst(resultValue)
		} else if lit, ok := expr.Expr.(*ast.BLangLiteral); semtypes.IsSubtypeSimple(exprTy, semtypes.INT) && ok {
			resultValue := lit.Value.(int64) * -1
			resultTy = semtypes.IntConst(resultValue)
		} else {
			resultTy = exprTy
		}
	case model.OperatorKind_ADD:
		resultTy = exprTy

	case model.OperatorKind_BITWISE_COMPLEMENT:
		if !semtypes.IsSubtypeSimple(exprTy, semtypes.INT) {
			t.semanticError(fmt.Sprintf("expect int type for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		if semtypes.IsSameType(t.typeContext(), exprTy, semtypes.INT) {
			resultTy = exprTy
			break
		}
		shape := semtypes.SingleShape(exprTy)
		if !shape.IsEmpty() {
			value, ok := shape.Get().Value.(int64)
			if !ok {
				t.internalError(fmt.Sprintf("unexpected singleton type for %s: %T", string(expr.GetOperatorKind()), shape.Get().Value), expr.GetPosition())
				return nil, expressionEffect{}, false
			}
			resultTy = semtypes.IntConst(^value)
		} else {
			resultTy = exprTy
		}

	case model.OperatorKind_NOT:
		if semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			if semtypes.IsSameType(t.typeContext(), exprTy, semtypes.BOOLEAN) {
				resultTy = semtypes.BOOLEAN
			} else {
				resultTy = semtypes.Diff(semtypes.BOOLEAN, exprTy)
			}
		} else {
			t.semanticError(fmt.Sprintf("expect boolean type for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		setExpectedType(expr, resultTy)
		return resultTy, expressionEffect{ifTrue: innerEffect.ifFalse, ifFalse: innerEffect.ifTrue}, true
	default:
		t.internalError(fmt.Sprintf("unsupported unary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func resolveBinaryExpr(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	if isLogicalExpression(expr) {
		return resolveLogicalExpr(t, chain, expr)
	}
	lhsTy, _, ok := resolveExpression(t, chain, expr.LhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, _, ok := resolveExpression(t, chain, expr.RhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType

	if isEqualityExpr(expr) {
		return resolveEqualityExpr(t, chain, expr)
	} else if isRangeExpr(expr) {
		resultTy = createIteratorType(t.typeEnv(), semtypes.INT, semtypes.NIL)
	} else {
		var nilLifted bool
		resultTy, nilLifted = nilLiftingExprResultTy(t, lhsTy, rhsTy, expr)
		if resultTy == nil {
			return nil, expressionEffect{}, false
		}
		if nilLifted {
			resultTy = semtypes.Union(semtypes.NIL, resultTy)
		}
	}

	setExpectedType(expr, resultTy)
	effect := defaultExpressionEffect(chain)
	return resultTy, effect, true
}

func isSingletonBool(ty semtypes.SemType, value bool) bool {
	singleShape := semtypes.SingleShape(ty)
	if singleShape.IsPresent() {
		return singleShape.Get().Value == value
	} else {
		return false
	}
}

func resolveEqualityExpr(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	var effect expressionEffect
	if expr.OpKind == model.OperatorKind_EQUAL || expr.OpKind == model.OperatorKind_NOT_EQUAL {
		effect = equalityNarrowingEffect(t, chain, expr)
	} else {
		effect = defaultExpressionEffect(chain)
	}
	resultTy := semtypes.BOOLEAN
	expr.SetDeterminedType(resultTy)
	return resultTy, effect, true
}

func resolveLogicalExpr(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	switch expr.OpKind {
	case model.OperatorKind_AND:
		return resolveAndExpr(t, chain, expr)
	case model.OperatorKind_OR:
		return resolveOrExpr(t, chain, expr)
	default:
		t.internalError(fmt.Sprintf("Unexpected logical expression op %s", string(expr.OpKind)), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
}

func resolveAndExpr(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	lhsTy, lhsEffect, ok := resolveExpression(t, chain, expr.LhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, rhsEffect, ok := resolveExpression(t, lhsEffect.ifTrue, expr.RhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType = semtypes.BOOLEAN
	if isSingletonBool(lhsTy, false) || isSingletonBool(rhsTy, false) {
		resultTy = semtypes.BooleanConst(false)
	} else if isSingletonBool(lhsTy, true) && isSingletonBool(rhsTy, true) {
		resultTy = semtypes.BooleanConst(true)
	} else if isSingletonBool(lhsTy, true) {
		resultTy = rhsTy
	}
	setExpectedType(expr, resultTy)

	if effect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return resultTy, effect, true
	}

	rhsDiffTrue := diff(rhsEffect.ifTrue, lhsEffect.ifTrue)
	rhsDiffFalse := diff(rhsEffect.ifFalse, lhsEffect.ifTrue)
	ifTrue := mergeChains(t, lhsEffect.ifTrue, rhsDiffTrue, semtypes.Intersect)
	ifFalse := mergeChains(t, lhsEffect.ifFalse, mergeChains(t, lhsEffect.ifTrue, rhsDiffFalse, semtypes.Intersect), semtypes.Union)
	return resultTy, expressionEffect{ifTrue: ifTrue, ifFalse: ifFalse}, true
}

func resolveOrExpr(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) (semtypes.SemType, expressionEffect, bool) {
	lhsTy, lhsEffect, ok := resolveExpression(t, chain, expr.LhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	rhsTy, rhsEffect, ok := resolveExpression(t, lhsEffect.ifFalse, expr.RhsExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType = semtypes.BOOLEAN
	if isSingletonBool(lhsTy, true) || isSingletonBool(rhsTy, true) {
		resultTy = semtypes.BooleanConst(true)
	} else if isSingletonBool(lhsTy, false) && isSingletonBool(rhsTy, false) {
		resultTy = semtypes.BooleanConst(false)
	} else if isSingletonBool(lhsTy, false) {
		resultTy = rhsTy
	}
	setExpectedType(expr, resultTy)

	if effect, isSingleton := singletonExprEffect(chain, expr); isSingleton {
		return resultTy, effect, true
	}

	rhsDiffTrue := diff(rhsEffect.ifTrue, lhsEffect.ifFalse)
	rhsDiffFalse := diff(rhsEffect.ifFalse, lhsEffect.ifFalse)
	ifTrue := mergeChains(t, lhsEffect.ifTrue, mergeChains(t, lhsEffect.ifFalse, rhsDiffTrue, semtypes.Intersect), semtypes.Union)
	ifFalse := mergeChains(t, lhsEffect.ifFalse, rhsDiffFalse, semtypes.Intersect)
	return resultTy, expressionEffect{ifTrue: ifTrue, ifFalse: ifFalse}, true
}

func equalityNarrowingEffect(t typeResolver, chain *binding, expr *ast.BLangBinaryExpr) expressionEffect {
	lhsRef, lhsIsVarRef := varRefExp(chain, &expr.LhsExpr)
	rhsTy := expr.RhsExpr.GetDeterminedType()
	rhsIsSingleton := semtypes.SingleShape(rhsTy).IsPresent()
	if lhsIsVarRef && rhsIsSingleton {
		effect := buildEqualityNarrowing(t, chain, lhsRef, rhsTy)
		if expr.OpKind == model.OperatorKind_NOT_EQUAL {
			return expressionEffect{ifTrue: effect.ifFalse, ifFalse: effect.ifTrue}
		}
		return effect
	}
	rhsRef, rhsIsVarRef := varRefExp(chain, &expr.RhsExpr)
	lhsTy := expr.LhsExpr.GetDeterminedType()
	lhsIsSingleton := semtypes.SingleShape(lhsTy).IsPresent()
	if rhsIsVarRef && lhsIsSingleton {
		effect := buildEqualityNarrowing(t, chain, rhsRef, lhsTy)
		if expr.OpKind == model.OperatorKind_NOT_EQUAL {
			return expressionEffect{ifTrue: effect.ifFalse, ifFalse: effect.ifTrue}
		}
		return effect
	}
	return defaultExpressionEffect(chain)
}

func buildEqualityNarrowing(t typeResolver, chain *binding, ref model.SymbolRef, singletonTy semtypes.SemType) expressionEffect {
	symbolTy := t.symbolType(ref)
	trueTy := semtypes.Intersect(symbolTy, singletonTy)
	trueSym := narrowSymbol(t, ref, trueTy)
	trueChain := &binding{ref: ref, narrowedSymbol: trueSym, prev: chain}
	falseTy := semtypes.Diff(symbolTy, singletonTy)
	falseSym := narrowSymbol(t, ref, falseTy)
	falseChain := &binding{ref: ref, narrowedSymbol: falseSym, prev: chain}
	return expressionEffect{ifTrue: trueChain, ifFalse: falseChain}
}

var additiveSupportedTypes = semtypes.Union(semtypes.NUMBER, semtypes.STRING)

var bitWiseOpLookOrder = []semtypes.SemType{semtypes.UINT8, semtypes.UINT16, semtypes.UINT32}

func nilLiftingExprResultTy(t typeResolver, lhsTy, rhsTy semtypes.SemType, expr *ast.BLangBinaryExpr) (semtypes.SemType, bool) {
	nilLifted := false

	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		nilLifted = true
		lhsTy = semtypes.Diff(lhsTy, semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, semtypes.NIL)
	}

	lhsBasicTy := semtypes.WidenToBasicTypes(lhsTy)
	rhsBasicTy := semtypes.WidenToBasicTypes(rhsTy)

	numLhsBits := bits.OnesCount(uint(lhsBasicTy.All()))
	numRhsBits := bits.OnesCount(uint(rhsBasicTy.All()))

	if numLhsBits > 1 || numRhsBits > 1 {
		t.semanticError(fmt.Sprintf("union types not supported for %s", string(expr.GetOperatorKind())), expr.GetPosition())
		return nil, false
	}

	if isRelationalExpr(expr) {
		if semtypes.Comparable(t.typeContext(), lhsBasicTy, rhsBasicTy) {
			return semtypes.BOOLEAN, false
		}
		t.semanticError("values are not comparable", expr.GetPosition())
		return nil, false
	}

	if isMultiplicativeExpr(expr) {
		if !isNumericType(lhsBasicTy) || !isNumericType(rhsBasicTy) {
			t.semanticError(fmt.Sprintf("expect numeric types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return lhsBasicTy, nilLifted
		}
		ctx := t.typeContext()
		if semtypes.IsSubtype(ctx, rhsBasicTy, semtypes.INT) ||
			(expr.GetOperatorKind() == model.OperatorKind_MUL && semtypes.IsSubtype(ctx, lhsBasicTy, semtypes.INT)) {
			t.unimplemented("type coercion not supported", expr.GetPosition())
			return nil, false
		}
		t.semanticError("both operands must belong to same basic type", expr.GetPosition())
		return nil, false
	}

	if isAdditiveExpr(expr) {
		ctx := t.typeContext()
		if !semtypes.IsSubtype(ctx, lhsBasicTy, additiveSupportedTypes) || !semtypes.IsSubtype(ctx, rhsBasicTy, additiveSupportedTypes) {
			t.semanticError(fmt.Sprintf("expect numeric or string types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		if lhsBasicTy == rhsBasicTy {
			return lhsBasicTy, nilLifted
		}
		t.semanticError("both operands must belong to same basic type", expr.GetPosition())
		return nil, false
	}

	if isShiftExpr(expr) {
		ctx := t.typeContext()
		if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
			t.semanticError(fmt.Sprintf("expect integer types for %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}
		var resultTy semtypes.SemType = semtypes.INT
		switch expr.GetOperatorKind() {
		case model.OperatorKind_BITWISE_RIGHT_SHIFT, model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) {
					resultTy = ty
					break
				}
			}
		}
		return resultTy, nilLifted
	}

	if isBitWiseExpr(expr) {
		ctx := t.typeContext()
		if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
			t.semanticError("expect integer types for bitwise operators", expr.GetPosition())
			return nil, false
		}

		var resultTy semtypes.SemType = semtypes.INT
		switch expr.GetOperatorKind() {
		case model.OperatorKind_BITWISE_AND:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) || semtypes.IsSubtype(ctx, rhsTy, ty) {
					resultTy = ty
					break
				}
			}
		case model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
			for _, ty := range bitWiseOpLookOrder {
				if semtypes.IsSubtype(ctx, lhsTy, ty) && semtypes.IsSubtype(ctx, rhsTy, ty) {
					resultTy = ty
					break
				}
			}
		default:
			t.internalError(fmt.Sprintf("unsupported bitwise operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
			return nil, false
		}

		return resultTy, nilLifted
	}

	t.internalError(fmt.Sprintf("unsupported binary operator: %s", string(expr.GetOperatorKind())), expr.GetPosition())
	return nil, false
}

func createIteratorType(env semtypes.Env, t, c semtypes.SemType) semtypes.SemType {
	od := semtypes.NewObjectDefinition()

	fields := []semtypes.Field{
		semtypes.FieldFrom("value", t, false, false),
	}
	var rest semtypes.SemType = semtypes.NEVER
	recordTy := createClosedRecordType(env, fields, rest)

	resultTy := semtypes.Union(recordTy, c)

	ld := semtypes.NewListDefinition()
	listTy := ld.DefineListTypeWrapped(env, []semtypes.SemType{}, 0, semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
	fd := semtypes.NewFunctionDefinition()
	fnTy := fd.Define(env, listTy, resultTy, semtypes.FunctionQualifiersFrom(env, false, false))

	members := []semtypes.Member{
		{
			Name:       "next",
			ValueTy:    fnTy,
			Kind:       semtypes.MemberKindMethod,
			Visibility: semtypes.VisibilityPublic,
			Immutable:  true,
		},
	}
	return od.Define(env, semtypes.ObjectQualifiersDEFAULT, members)
}

func createClosedRecordType(env semtypes.Env, fields []semtypes.Field, rest semtypes.SemType) semtypes.SemType {
	md := semtypes.NewMappingDefinition()
	return md.DefineMappingTypeWrapped(env, fields, rest)
}

func resolveIndexBasedAccess(t typeResolver, chain *binding, expr *ast.BLangIndexBasedAccess) (semtypes.SemType, expressionEffect, bool) {
	containerExpr := expr.Expr
	containerExprTy, _, ok := resolveExpression(t, chain, containerExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	keyExpr := expr.IndexExpr
	keyExprTy, _, ok := resolveExpression(t, chain, keyExpr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}

	var resultTy semtypes.SemType

	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) {
		resultTy = semtypes.ListMemberTypeInnerVal(t.typeContext(), containerExprTy, keyExprTy)
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.MAPPING) {
		memberTy := semtypes.MappingMemberTypeInner(t.typeContext(), containerExprTy, keyExprTy)
		maybeMissing := semtypes.ContainsUndef(memberTy)
		if maybeMissing {
			memberTy = semtypes.Union(semtypes.Diff(memberTy, semtypes.UNDEF), semtypes.NIL)
		}
		resultTy = memberTy
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) {
		resultTy = semtypes.STRING
	} else {
		t.semanticError("unsupported container type for index based access", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, resultTy)
	return resultTy, defaultExpressionEffect(chain), true
}

func resolveFieldBaseAccess(t typeResolver, chain *binding, expr *ast.BLangFieldBaseAccess) (semtypes.SemType, expressionEffect, bool) {
	containerExprTy, _, ok := resolveExpression(t, chain, expr.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	keyTy := semtypes.StringConst(expr.Field.Value)

	if !semtypes.IsSubtypeSimple(containerExprTy, semtypes.MAPPING) {
		t.semanticError("unsupported container type for field access", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	memberTy := semtypes.MappingMemberTypeInner(t.typeContext(), containerExprTy, keyTy)
	maybeMissing := semtypes.ContainsUndef(memberTy)
	if maybeMissing {
		t.semanticError("field base access is only possible for required fields", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	setExpectedType(expr, memberTy)
	expr.Field.SetDeterminedType(semtypes.NEVER)
	return memberTy, defaultExpressionEffect(chain), true
}

func resolveInvocation(t typeResolver, chain *binding, expr *ast.BLangInvocation) (semtypes.SemType, expressionEffect, bool) {
	symbol := expr.RawSymbol
	if symbol == nil {
		t.internalError("invocation has no symbol", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	var (
		ty       semtypes.SemType
		effect   expressionEffect
		resolved bool
	)
	switch s := symbol.(type) {
	case *deferredMethodSymbol:
		ty, effect, resolved = resolveMethodCall(t, chain, expr, s)
	case *model.SymbolRef:
		ty, effect, resolved = resolveFunctionCall(t, chain, expr, *s)
	default:
		t.internalError(fmt.Sprintf("expected *model.SymbolRef, got %T", symbol), expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	if !resolved {
		return nil, expressionEffect{}, false
	}
	if expr.PkgAlias != nil {
		expr.PkgAlias.SetDeterminedType(semtypes.NEVER)
	}
	if expr.Name != nil {
		expr.Name.SetDeterminedType(semtypes.NEVER)
	}
	return ty, effect, true
}

func resolveMethodCall(t typeResolver, chain *binding, expr *ast.BLangInvocation, methodSymbol *deferredMethodSymbol) (semtypes.SemType, expressionEffect, bool) {
	recieverTy, _, ok := resolveExpression(t, chain, expr.Expr, nil)
	if !ok {
		return nil, expressionEffect{}, false
	}
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.OBJECT) {
		t.unimplemented("method calls not implemented", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	var symbolSpace model.ExportedSymbolSpace
	var pkgAlias ast.BLangIdentifier
	if semtypes.IsSubtypeSimple(recieverTy, semtypes.LIST) {
		pkgName := array.PackageName
		space, ok := t.lookupImportedSymbols(pkgName)
		if !ok {
			t.internalError(fmt.Sprintf("%s symbol space not found", pkgName), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		symbolSpace = space
		pkgAlias = ast.BLangIdentifier{Value: pkgName}
		if !t.hasImplicitImport(pkgName) {
			importNode := ast.BLangImportPackage{
				OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "array"}},
				Alias:        &pkgAlias,
			}
			ast.Walk(t, &importNode)
			t.addImplicitImport(pkgName, importNode)
		}
	} else if semtypes.IsSubtypeSimple(recieverTy, semtypes.INT) {
		pkgName := bInt.PackageName
		space, ok := t.lookupImportedSymbols(pkgName)
		if !ok {
			t.internalError(fmt.Sprintf("%s symbol space not found", pkgName), expr.GetPosition())
			return nil, expressionEffect{}, false
		}
		symbolSpace = space
		pkgAlias = ast.BLangIdentifier{Value: pkgName}
		if !t.hasImplicitImport(pkgName) {
			importNode := ast.BLangImportPackage{
				OrgName:      &ast.BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []ast.BLangIdentifier{{Value: "lang"}, {Value: "int"}},
				Alias:        &pkgAlias,
			}
			ast.Walk(t, &importNode)
			t.addImplicitImport(pkgName, importNode)
		}
	} else {
		t.unimplemented("lang.value not implemented", expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	symbolRef, ok := symbolSpace.GetSymbol(methodSymbol.name)
	if !ok {
		t.semanticError("method not found: "+methodSymbol.name, expr.GetPosition())
		return nil, expressionEffect{}, false
	}
	argExprs := make([]ast.BLangExpression, len(expr.ArgExprs)+1)
	argExprs[0] = expr.Expr
	for i, arg := range expr.ArgExprs {
		argExprs[i+1] = arg
	}
	expr.SetSymbol(symbolRef)
	expr.ArgExprs = argExprs
	expr.Expr = nil
	expr.PkgAlias = &pkgAlias
	return resolveFunctionCall(t, chain, expr, symbolRef)
}

func resolveFunctionCall(t typeResolver, chain *binding, expr *ast.BLangInvocation, symbolRef model.SymbolRef) (semtypes.SemType, expressionEffect, bool) {
	argTys := make([]semtypes.SemType, len(expr.ArgExprs))
	for i, arg := range expr.ArgExprs {
		argTy, _, ok := resolveExpression(t, chain, arg, nil)
		if !ok {
			return nil, expressionEffect{}, false
		}
		argTys[i] = argTy
	}

	baseSymbol := t.getSymbol(symbolRef)
	if genericFn, ok := baseSymbol.(model.GenericFunctionSymbol); ok {
		symbolRef = genericFn.Monomorphize(argTys)
		expr.SetSymbol(symbolRef)
	}

	symbolRef = lookupSymbol(chain, symbolRef)
	fnTy := t.symbolType(symbolRef)
	if fnTy == nil {
		t.internalError("function symbol has no type", expr.GetPosition())
		return nil, expressionEffect{}, false
	}

	paramListTy := semtypes.FunctionParamListType(t.typeContext(), fnTy)
	if paramListTy != nil {
		for i, arg := range expr.ArgExprs {
			key := semtypes.IntConst(int64(i))
			paramTy := semtypes.ListMemberTypeInnerVal(t.typeContext(), paramListTy, key)
			arg.SetDeterminedType(nil)
			argTy, _, ok := resolveExpression(t, chain, arg, paramTy)
			if !ok {
				return nil, expressionEffect{}, false
			}
			argTys[i] = argTy
		}
	}

	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(t.typeEnv(), argTys, len(argTys), semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)

	retTy := semtypes.FunctionReturnType(t.typeContext(), fnTy, argListTy)

	setExpectedType(expr, retTy)
	return retTy, defaultExpressionEffect(chain), true
}

func resolveBType(t typeResolver, btype ast.BType, depth int) (semtypes.SemType, bool) {
	bLangNode := btype.(ast.BLangNode)
	if bLangNode.GetDeterminedType() != nil {
		return bLangNode.GetDeterminedType(), true
	}
	res, ok := resolveBTypeInner(t, btype, depth)
	if !ok {
		return nil, false
	}
	bLangNode.SetDeterminedType(res)
	typeData := btype.GetTypeData()
	typeData.Type = res
	btype.SetTypeData(typeData)
	return res, true
}

func resolveTypeDataPair(t typeResolver, typeData *model.TypeData, depth int) (semtypes.SemType, bool) {
	ty, ok := resolveBType(t, typeData.TypeDescriptor.(ast.BType), depth)
	if !ok {
		return nil, false
	}
	typeData.Type = ty
	return ty, true
}

func resolveBTypeInner(t typeResolver, btype ast.BType, depth int) (semtypes.SemType, bool) {
	switch ty := btype.(type) {
	case *ast.BLangValueType:
		switch ty.TypeKind {
		case model.TypeKind_BOOLEAN:
			return semtypes.BOOLEAN, true
		case model.TypeKind_INT:
			return semtypes.INT, true
		case model.TypeKind_FLOAT:
			return semtypes.FLOAT, true
		case model.TypeKind_STRING:
			return semtypes.STRING, true
		case model.TypeKind_NIL:
			return semtypes.NIL, true
		case model.TypeKind_ANY:
			return semtypes.ANY, true
		case model.TypeKind_DECIMAL:
			return semtypes.DECIMAL, true
		case model.TypeKind_BYTE:
			return semtypes.BYTE, true
		case model.TypeKind_ANYDATA:
			return semtypes.CreateAnydata(t.typeContext()), true
		default:
			t.internalError("unexpected type kind", nil)
			return nil, false
		}
	case *ast.BLangArrayType:
		defn := ty.Definition
		var semTy semtypes.SemType
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			elemTy, ok := resolveTypeDataPair(t, &ty.Elemtype, depth+1)
			if !ok {
				return nil, false
			}
			for i := len(ty.Sizes); i > 0; i-- {
				lenExp := ty.Sizes[i-1]
				if lenExp == nil {
					elemTy = d.DefineListTypeWrappedWithEnvSemType(t.typeEnv(), elemTy)
				} else {
					length := int(lenExp.(*ast.BLangLiteral).Value.(int64))
					elemTy = d.DefineListTypeWrappedWithEnvSemTypesInt(t.typeEnv(), []semtypes.SemType{elemTy}, length)
				}
			}
			semTy = elemTy
		} else {
			semTy = defn.GetSemType(t.typeEnv())
		}
		return semTy, true
	case *ast.BLangUnionTypeNode:
		lhs, ok := resolveTypeDataPair(t, ty.Lhs(), depth+1)
		if !ok {
			return nil, false
		}
		rhs, ok := resolveTypeDataPair(t, ty.Rhs(), depth+1)
		if !ok {
			return nil, false
		}
		return semtypes.Union(lhs, rhs), true
	case *ast.BLangIntersectionTypeNode:
		lhs, ok := resolveTypeDataPair(t, ty.Lhs(), depth+1)
		if !ok {
			return nil, false
		}
		rhs, ok := resolveTypeDataPair(t, ty.Rhs(), depth+1)
		if !ok {
			return nil, false
		}
		result := semtypes.Intersect(lhs, rhs)
		if semtypes.IsEmpty(t.typeContext(), result) {
			t.semanticError("intersection type is empty (equivalent to never)", ty.GetPosition())
			return nil, false
		}
		return result, true
	case *ast.BLangErrorTypeNode:
		if ty.IsDistinct() {
			panic("distinct error types not supported")
		}
		if ty.IsTop() {
			return semtypes.ERROR, true
		} else {
			detailTy, ok := resolveBType(t, ty.DetailType.TypeDescriptor.(ast.BType), depth+1)
			if !ok {
				return nil, false
			}
			ty.DetailType.Type = detailTy
			return semtypes.ErrorWithDetail(detailTy), true
		}
	case *ast.BLangUserDefinedType:
		ast.Walk(t, &ty.TypeName)
		ast.Walk(t, &ty.PkgAlias)
		symbol := ty.Symbol()
		if ty.PkgAlias.Value != "" {
			return t.symbolType(symbol), true
		}
		defn, ok := t.getTypeDefinition(symbol)
		if !ok {
			t.internalError("type definition not found", nil)
			return nil, false
		}
		return resolveTypeDefinition(t, defn, depth)
	case *ast.BLangFiniteTypeNode:
		var result semtypes.SemType = semtypes.NEVER
		for _, value := range ty.ValueSpace {
			valueTy, _, ok := resolveExpression(t, nil, value, nil)
			if !ok {
				return nil, false
			}
			result = semtypes.Union(result, valueTy)
		}
		return result, true
	case *ast.BLangConstrainedType:
		if _, ok := resolveTypeDataPair(t, &ty.Type, depth+1); !ok {
			return nil, false
		}
		defn := ty.Definition
		if defn == nil {
			switch ty.GetTypeKind() {
			case model.TypeKind_MAP:
				d := semtypes.NewMappingDefinition()
				ty.Definition = &d
				rest, ok := resolveTypeDataPair(t, &ty.Constraint, depth+1)
				if !ok {
					return nil, false
				}
				return d.DefineMappingTypeWrapped(t.typeEnv(), nil, rest), true
			default:
				t.unimplemented("unsupported base type kind", nil)
				return nil, false
			}
		} else {
			return defn.GetSemType(t.typeEnv()), true
		}
	case *ast.BLangBuiltInRefTypeNode:
		switch ty.TypeKind {
		case model.TypeKind_MAP:
			return semtypes.MAPPING, true
		default:
			t.internalError("Unexpected builtin type kind", ty.GetPosition())
		}
		return nil, false
	case *ast.BLangTupleTypeNode:
		defn := ty.Definition
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			members := make([]semtypes.SemType, len(ty.Members))
			for i, member := range ty.Members {
				memberTy, ok := resolveBType(t, member.TypeDesc.(ast.BType), depth+1)
				if !ok {
					return nil, false
				}
				members[i] = memberTy
			}
			rest, ok := semtypes.SemType(semtypes.NEVER), true
			if ty.Rest != nil {
				rest, ok = resolveBType(t, ty.Rest.(ast.BType), depth+1)
				if !ok {
					return nil, false
				}
			}
			return d.DefineListTypeWrappedWithEnvSemTypesSemType(t.typeEnv(), members, rest), true
		}
		return defn.GetSemType(t.typeEnv()), true
	case *ast.BLangRecordType:
		defn := ty.Definition
		if defn != nil {
			return defn.GetSemType(t.typeEnv()), true
		}
		d := semtypes.NewMappingDefinition()
		ty.Definition = &d

		includedFields := make(map[string][]ast.BField)
		needsRestOverride, includedRest, ok := accumIncludedFields(t, ty, includedFields, false, nil)
		if !ok {
			return nil, false
		}
		seen := make(map[string]bool)
		var fields []semtypes.Field
		for name, field := range ty.Fields() {
			if seen[name] {
				t.semanticError(fmt.Sprintf("duplicate field name '%s'", name), field.GetPosition())
				return nil, false
			}
			seen[name] = true
			fieldTy, ok := resolveBType(t, field.Type, depth+1)
			if !ok {
				return nil, false
			}
			if overridden, exists := includedFields[name]; exists {
				for _, incField := range overridden {
					incFieldTy, ok := resolveBType(t, incField.Type, depth+1)
					if !ok {
						return nil, false
					}
					if !semtypes.IsSubtype(t.typeContext(), fieldTy, incFieldTy) {
						t.semanticError(
							fmt.Sprintf("field '%s' of type that overrides included field is not a subtype of the included field type", name),
							field.GetPosition(),
						)
					}
				}
				delete(includedFields, name)
			}
			ro := field.FlagSet.Contains(model.Flag_READONLY)
			opt := field.FlagSet.Contains(model.Flag_OPTIONAL)
			fields = append(fields, semtypes.FieldFrom(name, fieldTy, ro, opt))
		}

		for name, incFields := range includedFields {
			if len(incFields) > 1 {
				t.semanticError(fmt.Sprintf("included field '%s' declared in multiple type inclusions must be overridden", name), ty.GetPosition())
			}
		}

		for name, incFields := range includedFields {
			if len(incFields) > 1 {
				continue
			}
			field := incFields[0]
			fieldTy, ok := resolveBType(t, field.Type, depth+1)
			if !ok {
				return nil, false
			}
			ro := field.FlagSet.Contains(model.Flag_READONLY)
			opt := field.FlagSet.Contains(model.Flag_OPTIONAL)
			fields = append(fields, semtypes.FieldFrom(name, fieldTy, ro, opt))
		}

		var rest semtypes.SemType
		if ty.RestType != nil {
			var ok bool
			rest, ok = resolveBType(t, ty.RestType, depth+1)
			if !ok {
				return nil, false
			}
		} else if ty.IsOpen {
			rest = semtypes.CreateAnydata(t.typeContext())
		} else if needsRestOverride {
			t.semanticError("included rest type declared in multiple type inclusions must be overridden", ty.GetPosition())
			rest = semtypes.NEVER
		} else if includedRest != nil {
			var ok bool
			rest, ok = resolveBType(t, includedRest, depth+1)
			if !ok {
				return nil, false
			}
		} else {
			rest = semtypes.NEVER
		}
		return d.DefineMappingTypeWrapped(t.typeEnv(), fields, rest), true
	default:
		t.unimplemented("unsupported type", nil)
		return nil, false
	}
}

func accumIncludedFields(t typeResolver, recordTy *ast.BLangRecordType, includedFields map[string][]ast.BField, needsRestOverride bool, includedRest ast.BType) (bool, ast.BType, bool) {
	for _, inc := range recordTy.TypeInclusions {
		udt, ok := inc.(*ast.BLangUserDefinedType)
		if !ok {
			t.semanticError("type inclusion must be a user-defined type", inc.(ast.BLangNode).GetPosition())
			continue
		}

		_, ok = resolveBType(t, inc, 0)
		if !ok {
			return false, nil, false
		}

		symbol := udt.Symbol()
		tDefn, ok := t.getTypeDefinition(symbol)
		if !ok {
			t.internalError("type definition not found for inclusion", udt.GetPosition())
			continue
		}
		recTy, ok := tDefn.GetTypeData().TypeDescriptor.(*ast.BLangRecordType)
		if !ok {
			t.semanticError("included type is not a record type", udt.GetPosition())
			continue
		}

		needsRestOverride, includedRest, ok = accumIncludedFields(t, recTy, includedFields, needsRestOverride, includedRest)
		if !ok {
			return false, nil, false
		}

		for name, field := range recTy.Fields() {
			includedFields[name] = append(includedFields[name], field)
		}

		if recTy.RestType != nil {
			if includedRest != nil {
				needsRestOverride = true
			}
			includedRest = recTy.RestType
		}
	}
	return needsRestOverride, includedRest, true
}

func resolveConstant(t typeResolver, constant *ast.BLangConstant) bool {
	if constant.Expr == nil {
		t.internalError("constant expression is nil", constant.GetPosition())
		return false
	}
	if constant.Name != nil {
		ast.Walk(t, constant.Name)
	}

	var annotationType semtypes.SemType
	if typeNode := constant.TypeNode(); typeNode != nil {
		var ok bool
		annotationType, ok = resolveBType(t, typeNode, 0)
		if !ok {
			return false
		}
	}

	exprTy, _, ok := resolveExpression(t, nil, constant.Expr.(ast.BLangExpression), annotationType)
	if !ok {
		return false
	}

	var expectedType semtypes.SemType
	if annotationType != nil {
		expectedType = annotationType
	} else {
		expectedType = exprTy
	}
	setExpectedType(constant, expectedType)
	symbol := constant.Symbol()
	t.setSymbolType(symbol, expectedType)

	return true
}

func resolveMatchStatement(t typeResolver, chain *binding, stmt *ast.BLangMatchStatement) (statementEffect, bool) {
	_, exprEffect, ok := resolveExpression(t, chain, stmt.Expr, nil)
	if !ok {
		return defaultStmtEffect(chain), false
	}
	chain = exprEffect.ifTrue

	exprRef, isVarRef := varRefExp(chain, &stmt.Expr)
	var remainingType semtypes.SemType
	if isVarRef {
		remainingType = t.symbolType(exprRef)
	} else {
		remainingType = stmt.Expr.GetDeterminedType()
	}
	allNonCompletion := true
	var bodyEffects []statementEffect

	tyCtx := semtypes.ContextFrom(t.typeEnv())

	for i := range stmt.MatchClauses {
		clause := &stmt.MatchClauses[i]

		if semtypes.IsEmpty(tyCtx, remainingType) {
			t.semanticError("unreachable match clause", clause.GetPosition())
		}

		var bodyChain *binding
		var ok bool
		clause.AcceptedType, bodyChain, ok = matchClauseAcceptedType(t, chain, clause)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		clauseAcceptedType := semtypes.Intersect(remainingType, clause.AcceptedType)

		clauseIsEmpty := semtypes.IsEmpty(tyCtx, clauseAcceptedType)
		if clauseIsEmpty {
			t.semanticError("unmatchable match clause", clause.GetPosition())
		}

		clause.AcceptedType = clauseAcceptedType

		if clauseIsEmpty {
			_, ok := resolveMatchClause(t, bodyChain, clause)
			if !ok {
				return defaultStmtEffect(chain), false
			}
			continue
		}

		if isVarRef {
			baseRef := t.unnarrowedSymbol(exprRef)
			narrowedSym := narrowSymbol(t, baseRef, clauseAcceptedType)
			bodyChain = &binding{
				ref:            baseRef,
				narrowedSymbol: narrowedSym,
				prev:           bodyChain,
			}
		}

		bodyEffect, ok := resolveMatchClause(t, bodyChain, clause)
		if !ok {
			return defaultStmtEffect(chain), false
		}
		bodyEffects = append(bodyEffects, bodyEffect)
		if !bodyEffect.nonCompletion {
			allNonCompletion = false
		}

		remainingType = semtypes.Diff(remainingType, clause.AcceptedType)
	}

	stmt.IsExhaustive = semtypes.IsEmpty(tyCtx, remainingType)

	if stmt.IsExhaustive && allNonCompletion {
		return statementEffect{chain, true}, true
	}

	var result *binding
	first := true
	for _, effect := range bodyEffects {
		if effect.nonCompletion {
			continue
		}
		if first {
			result = effect.binding
			first = false
		} else {
			result = mergeChains(t, result, effect.binding, semtypes.Union)
		}
	}
	return statementEffect{result, false}, true
}

func matchClauseAcceptedType(t typeResolver, chain *binding, clause *ast.BLangMatchClause) (semtypes.SemType, *binding, bool) {
	var acceptedTy semtypes.SemType = semtypes.NEVER
	for _, pattern := range clause.Patterns {
		patternTy, ok := resolveMatchPattern(t, chain, pattern)
		if !ok {
			return nil, nil, false
		}
		acceptedTy = semtypes.Union(acceptedTy, patternTy)
	}
	if clause.Guard != nil {
		_, guardEffect, _ := resolveExpression(t, chain, clause.Guard, nil)
		return acceptedTy, guardEffect.ifTrue, true
	}
	return acceptedTy, chain, true
}

func resolveMatchClause(t typeResolver, chain *binding, clause *ast.BLangMatchClause) (statementEffect, bool) {
	bodyEffect, ok := resolveBlockStatements(t, chain, clause.Body.Stmts)
	if !ok {
		return defaultStmtEffect(chain), false
	}
	clause.Body.SetDeterminedType(semtypes.NEVER)
	clause.SetDeterminedType(semtypes.NEVER)
	return bodyEffect, true
}

func resolveMatchPattern(t typeResolver, chain *binding, pattern ast.BLangMatchPattern) (semtypes.SemType, bool) {
	switch p := pattern.(type) {
	case *ast.BLangConstPattern:
		ty, _, ok := resolveExpression(t, chain, p.Expr, nil)
		if !ok {
			return nil, false
		}
		p.SetAcceptedType(ty)
		p.SetDeterminedType(semtypes.NEVER)
		return ty, true
	case *ast.BLangWildCardMatchPattern:
		ty := semtypes.ANY
		p.SetAcceptedType(ty)
		p.SetDeterminedType(semtypes.NEVER)
		return ty, true
	default:
		t.internalError(fmt.Sprintf("unexpected match pattern type: %T", pattern), pattern.GetPosition())
		return semtypes.NEVER, false
	}
}
