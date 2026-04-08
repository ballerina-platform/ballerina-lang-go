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
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

type desugaredNode[E model.Node] struct {
	initStmts       []model.StatementNode
	replacementNode E
}

// packageContext holds shared state for desugaring a single package.
type packageContext struct {
	compilerCtx          *context.CompilerContext
	pkg                  *ast.BLangPackage
	importedSymbols      map[string]model.ExportedSymbolSpace
	importMu             sync.Mutex
	addedImplicitImports map[string]bool
	desugarSymbolCounter int
}

var _ desugarContext = &packageContext{}

func newPackageContext(compilerCtx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *packageContext {
	return &packageContext{
		compilerCtx:          compilerCtx,
		pkg:                  pkg,
		importedSymbols:      importedSymbols,
		addedImplicitImports: make(map[string]bool),
	}
}

func (ctx *packageContext) addImplicitImport(pkgName string, imp ast.BLangImportPackage) {
	ctx.importMu.Lock()
	defer ctx.importMu.Unlock()
	if !ctx.addedImplicitImports[pkgName] {
		ctx.addedImplicitImports[pkgName] = true
		ctx.pkg.Imports = append(ctx.pkg.Imports, imp)
	}
}

func (ctx *packageContext) getImportedSymbolSpace(pkgName string) (model.ExportedSymbolSpace, bool) {
	space, ok := ctx.importedSymbols[pkgName]
	return space, ok
}

func (ctx *packageContext) symbolType(symbol model.SymbolRef) semtypes.SemType {
	return ctx.compilerCtx.SymbolType(symbol)
}

func (ctx *packageContext) newFunctionScope(parent model.Scope) *model.FunctionScope {
	return ctx.compilerCtx.NewFunctionScope(parent, *ctx.pkg.PackageID)
}

func (ctx *packageContext) getSymbol(ref model.SymbolRef) model.Symbol {
	return ctx.compilerCtx.GetSymbol(ref)
}

func (ctx *packageContext) getSymbolType(ref model.SymbolRef) semtypes.SemType {
	return ctx.compilerCtx.SymbolType(ref)
}

func (ctx *packageContext) setSymbolType(ref model.SymbolRef, ty semtypes.SemType) {
	ctx.compilerCtx.SetSymbolType(ref, ty)
}

func (ctx *packageContext) typeEnv() semtypes.Env {
	return ctx.compilerCtx.GetTypeEnv()
}

func (ctx *packageContext) nextDesugarSymbolName() string {
	name := fmt.Sprintf("$desugar$%d", ctx.desugarSymbolCounter)
	ctx.desugarSymbolCounter++
	return name
}

func (ctx *packageContext) addSymbolToSameSpace(ref model.SymbolRef, name string, symbol model.Symbol) model.SymbolRef {
	return ctx.compilerCtx.AddSymbolToSameSpace(ref, name, symbol)
}

func (ctx *packageContext) internalError(msg string) {
	ctx.compilerCtx.InternalError(msg, nil)
}

func (ctx *packageContext) unimplemented(msg string) {
	ctx.compilerCtx.Unimplemented(msg, nil)
}

type functionContext struct {
	pkgCtx               *packageContext
	scopeStack           []model.Scope
	desugarSymbolCounter int
	loopVarStack         []ast.BLangExpression // Stack to track loop variables (nil for while, varRef for desugared foreach)
}

var _ desugarContext = &functionContext{}

func (ctx *functionContext) internalError(msg string) {
	ctx.pkgCtx.internalError(msg)
}

func (ctx *functionContext) unimplemented(msg string) {
	ctx.pkgCtx.unimplemented(msg)
}

func (ctx *functionContext) getImportedSymbolSpace(pkgName string) (model.ExportedSymbolSpace, bool) {
	return ctx.pkgCtx.getImportedSymbolSpace(pkgName)
}

func (ctx *functionContext) addImplicitImport(pkgName string, imp ast.BLangImportPackage) {
	ctx.pkgCtx.addImplicitImport(pkgName, imp)
}

func (ctx *functionContext) symbolType(symbol model.SymbolRef) semtypes.SemType {
	return ctx.pkgCtx.symbolType(symbol)
}

func (ctx *functionContext) pushScope(scope model.Scope) {
	ctx.scopeStack = append(ctx.scopeStack, scope)
}

func (ctx *functionContext) popScope() {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("cannot pop from empty scope stack")
	}
	ctx.scopeStack = ctx.scopeStack[:len(ctx.scopeStack)-1]
}

func (ctx *functionContext) currentScope() model.Scope {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("scope stack is empty")
	}
	return ctx.scopeStack[len(ctx.scopeStack)-1]
}

func (ctx *functionContext) pushLoopVar(varRef ast.BLangExpression) {
	ctx.loopVarStack = append(ctx.loopVarStack, varRef)
}

func (ctx *functionContext) popLoopVar() {
	if len(ctx.loopVarStack) == 0 {
		ctx.internalError("cannot pop from empty loopVar stack")
	}
	ctx.loopVarStack = ctx.loopVarStack[:len(ctx.loopVarStack)-1]
}

func (ctx *functionContext) currentLoopVar() ast.BLangExpression {
	if len(ctx.loopVarStack) == 0 {
		return nil
	}
	return ctx.loopVarStack[len(ctx.loopVarStack)-1]
}

func (ctx *functionContext) nextDesugarSymbolName() string {
	name := fmt.Sprintf("$desugar$%d", ctx.desugarSymbolCounter)
	ctx.desugarSymbolCounter++
	return name
}

func (ctx *functionContext) addSymbolToSameSpace(ref model.SymbolRef, name string, symbol model.Symbol) model.SymbolRef {
	return ctx.pkgCtx.addSymbolToSameSpace(ref, name, symbol)
}

func (ctx *functionContext) newFunctionScope(parent model.Scope) *model.FunctionScope {
	return ctx.pkgCtx.newFunctionScope(parent)
}

func (ctx *functionContext) setSymbolType(ref model.SymbolRef, ty semtypes.SemType) {
	ctx.pkgCtx.setSymbolType(ref, ty)
}

func (ctx *functionContext) getSymbol(ref model.SymbolRef) model.Symbol {
	return ctx.pkgCtx.getSymbol(ref)
}

func (ctx *functionContext) typeEnv() semtypes.Env {
	return ctx.pkgCtx.typeEnv()
}

type desugarContext interface {
	nextDesugarSymbolName() string
	addSymbolToSameSpace(ref model.SymbolRef, name string, symbol model.Symbol) model.SymbolRef
	newFunctionScope(parent model.Scope) *model.FunctionScope
	setSymbolType(ref model.SymbolRef, ty semtypes.SemType)
	getSymbol(ref model.SymbolRef) model.Symbol
	typeEnv() semtypes.Env
	internalError(msg string)
}

type desugaredSymbol struct {
	name     string
	ty       semtypes.SemType
	kind     model.SymbolKind
	isPublic bool
}

var _ model.Symbol = &desugaredSymbol{}

func (s *desugaredSymbol) Name() string {
	return s.name
}

func (s *desugaredSymbol) Type() semtypes.SemType {
	return s.ty
}

func (s *desugaredSymbol) Kind() model.SymbolKind {
	return s.kind
}

func (s *desugaredSymbol) SetType(_ semtypes.SemType) {
	panic("SetType is not supported for desugared symbols")
}

func (s *desugaredSymbol) IsPublic() bool {
	return s.isPublic
}

func (s *desugaredSymbol) Copy() model.Symbol {
	cp := *s
	return &cp
}

func (ctx *functionContext) addDesugardSymbol(ty semtypes.SemType, kind model.SymbolKind, isPublic bool) (string, model.SymbolRef) {
	if len(ctx.scopeStack) == 0 {
		ctx.internalError("cannot add desugared symbol when scope stack is empty")
	}
	name := ctx.nextDesugarSymbolName()
	symbol := &desugaredSymbol{
		name:     name,
		ty:       ty,
		kind:     kind,
		isPublic: isPublic,
	}
	ctx.currentScope().AddSymbol(name, symbol)
	ref, _ := ctx.currentScope().GetSymbol(name)
	return name, ref
}

func desugarInitFn(pkgCtx *packageContext, compilerCtx *context.CompilerContext, pkg *ast.BLangPackage) {
	var initStmts []ast.BLangStatement

	for i := range pkg.GlobalVars {
		globalVar := &pkg.GlobalVars[i]
		if globalVar.Expr == nil {
			continue
		}
		initExpr := globalVar.Expr.(ast.BLangExpression)
		basePos := initExpr.GetPosition()
		varRef := &ast.BLangSimpleVarRef{
			VariableName: globalVar.Name,
		}
		varRef.SetSymbol(globalVar.Symbol())
		varRef.SetDeterminedType(globalVar.GetDeterminedType())
		assignment := &ast.BLangAssignment{
			VarRef: varRef,
			Expr:   initExpr,
		}
		assignment.SetDeterminedType(semtypes.NEVER)
		setPositionIfMissing(assignment, basePos)

		initStmts = append(initStmts, assignment)
		globalVar.Expr = nil
	}

	for i := range pkg.Constants {
		constant := &pkg.Constants[i]
		if constant.Expr == nil {
			continue
		}
		initExpr := constant.Expr.(ast.BLangExpression)
		basePos := initExpr.GetPosition()
		varRef := &ast.BLangSimpleVarRef{
			VariableName: constant.Name,
		}
		varRef.SetSymbol(constant.Symbol())
		varRef.SetDeterminedType(constant.GetDeterminedType())
		assignment := &ast.BLangAssignment{
			VarRef: varRef,
			Expr:   initExpr,
		}
		assignment.SetDeterminedType(semtypes.NEVER)
		setPositionIfMissing(assignment, basePos)

		initStmts = append(initStmts, assignment)
		constant.Expr = nil
	}

	if len(initStmts) == 0 && pkg.InitFunction == nil {
		return
	}

	if pkg.InitFunction == nil {
		initName := &ast.BLangIdentifier{Value: "init"}
		initName.SetDeterminedType(semtypes.NEVER)
		pkg.InitFunction = &ast.BLangFunction{}
		pkg.InitFunction.Name = initName
		body := &ast.BLangBlockFunctionBody{
			Stmts: initStmts,
		}
		body.SetDeterminedType(semtypes.NEVER)
		pkg.InitFunction.Body = body
		pkg.InitFunction.SetDeterminedType(semtypes.NEVER)
		// Create a proper function symbol and scope for the synthetic init function
		pkgID := pkg.PackageID
		signature := model.FunctionSignature{ReturnType: semtypes.NIL}
		initSymbol := model.NewFunctionSymbol("init", signature, false)
		symbolSpace := compilerCtx.NewSymbolSpace(*pkgID)
		symbolSpace.AddSymbol("init", initSymbol)
		symRef, _ := symbolSpace.GetSymbol("init")
		pkg.InitFunction.SetSymbol(symRef)
		fnScope := compilerCtx.NewFunctionScope(nil, *pkgID)
		pkg.InitFunction.SetScope(fnScope)
	} else {
		body := pkg.InitFunction.Body.(*ast.BLangBlockFunctionBody)
		body.Stmts = append(initStmts, body.Stmts...)
	}

	*pkg.InitFunction = *desugarFunction(pkgCtx, pkg.InitFunction)
}

func createDefaultValueFunction(name string, defaultExpr ast.BLangExpression) *ast.BLangFunction {
	retStmt := &ast.BLangReturn{Expr: defaultExpr}
	retStmt.SetDeterminedType(semtypes.NEVER)
	body := &ast.BLangBlockFunctionBody{Stmts: []ast.BLangStatement{retStmt}}
	body.SetDeterminedType(semtypes.NEVER)

	fn := &ast.BLangFunction{}
	fn.Name = &ast.BLangIdentifier{Value: name}
	fn.Name.SetDeterminedType(semtypes.NEVER)
	fn.Body = body
	fn.SetDeterminedType(semtypes.NEVER)
	setPositionIfMissing(fn, defaultExpr.GetPosition())
	return fn
}

type desugaredRecordFieldResult struct {
	fn     *ast.BLangFunction
	symRef model.SymbolRef
}

type desugaredTypeDescResult struct {
	recordFields []desugaredRecordFieldResult
}

func desugarTypeDesc(ctx desugarContext, typeDesc ast.BType, ownerSymRef model.SymbolRef, parentScope model.Scope) desugaredTypeDescResult {
	switch td := typeDesc.(type) {
	case *ast.BLangRecordType:
		return desugarRecordTypeDesc(ctx, td, ownerSymRef, parentScope)
	}
	return desugaredTypeDescResult{}
}

func desugarRecordTypeDesc(ctx desugarContext, recType *ast.BLangRecordType, ownerSymRef model.SymbolRef, parentScope model.Scope) desugaredTypeDescResult {
	var fields []desugaredRecordFieldResult
	for _, field := range recType.FieldPtrs() {
		if field.DefaultExpr == nil {
			continue
		}
		fieldSemType := field.Type.(ast.BLangNode).GetDeterminedType()
		fnName := ctx.nextDesugarSymbolName()
		signature := model.FunctionSignature{ReturnType: fieldSemType}
		fnSymbol := model.NewFunctionSymbol(fnName, signature, false)
		env := ctx.typeEnv()
		paramListDefn := semtypes.NewListDefinition()
		paramListTy := paramListDefn.TupleTypeWrapped(env)
		fnDefn := semtypes.NewFunctionDefinition()
		fnType := fnDefn.Define(env, paramListTy, fieldSemType, semtypes.FunctionQualifiersFrom(env, false, false))
		fnSymbol.SetType(fnType)
		symRef := ctx.addSymbolToSameSpace(ownerSymRef, fnName, fnSymbol)
		fnScope := ctx.newFunctionScope(parentScope)

		fn := createDefaultValueFunction(fnName, field.DefaultExpr)
		fn.SetSymbol(symRef)
		fn.SetScope(fnScope)

		field.DefaultFnRef = symRef
		fields = append(fields, desugaredRecordFieldResult{fn: fn, symRef: symRef})
	}
	return desugaredTypeDescResult{recordFields: fields}
}

func desugarTopLevelTypeDescs(cx *packageContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		typeDesc, ok := defn.GetTypeData().TypeDescriptor.(ast.BType)
		if !ok {
			cx.internalError("type definition has no BType type descriptor")
			return
		}
		result := desugarTypeDesc(cx, typeDesc, defn.Symbol(), nil)
		for _, rf := range result.recordFields {
			pkg.Functions = append(pkg.Functions, *rf.fn)
		}
	}
}

// DesugarPackage returns a desugared package (may be new or same instance)
func DesugarPackage(compilerCtx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *ast.BLangPackage {
	if importedSymbols == nil {
		importedSymbols = make(map[string]model.ExportedSymbolSpace)
	}
	pkgCtx := newPackageContext(compilerCtx, pkg, importedSymbols)

	var wg sync.WaitGroup
	var panicErr any

	desugarFn := func(fn *ast.BLangFunction) {
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			*fn = *desugarFunction(pkgCtx, fn)
		})
	}

	// Desugar type definition default expressions into standalone functions
	desugarTopLevelTypeDescs(pkgCtx, pkg)

	// Desugar all functions
	for i := range pkg.Functions {
		desugarFn(&pkg.Functions[i])
	}

	// Desugar class definitions (each class concurrently, members sequentially)
	for i := range pkg.ClassDefinitions {
		class := &pkg.ClassDefinitions[i]
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			desugarClassDefinition(pkgCtx, class)
			for name, method := range class.Methods {
				class.Methods[name] = desugarFunction(pkgCtx, method)
			}
			*class.InitFunction = *desugarFunction(pkgCtx, class.InitFunction)
		})
	}

	// Desugar init, start, stop functions
	desugarInitFn(pkgCtx, compilerCtx, pkg)
	if pkg.StartFunction != nil {
		desugarFn(pkg.StartFunction)
	}
	if pkg.StopFunction != nil {
		desugarFn(pkg.StopFunction)
	}

	wg.Wait()
	if panicErr != nil {
		panic(panicErr)
	}

	return pkg
}

func desugarClassDefinition(pkgCtx *packageContext, class *ast.BLangClassDefinition) {
	if class.InitFunction == nil {
		fn := ast.BLangFunction{
			ObjInitFunction: true,
		}
		fn.FlagSet.Add(model.Flag_ATTACHED)
		fn.Name = &ast.BLangIdentifier{Value: "init"}
		fn.Body = &ast.BLangBlockFunctionBody{}
		fn.SetDeterminedType(semtypes.NEVER)
		fn.SetScope(pkgCtx.newFunctionScope(class.Scope()))
		initSymbol := model.NewFunctionSymbol("init", model.FunctionSignature{ReturnType: semtypes.NIL}, false)
		classScope := class.Scope()
		classScope.AddSymbol("init", initSymbol)
		symRef, _ := classScope.GetSymbol("init")
		fn.SetSymbol(symRef)
		class.InitFunction = &fn
	}

	var initStmts []ast.BLangStatement
	classScope := class.Scope()
	selfRef, ok := classScope.GetSymbol("self")
	if !ok {
		pkgCtx.internalError("self symbol not found in class scope")
	}
	classType := pkgCtx.getSymbol(selfRef).Type()

	for _, field := range class.Fields {
		initExpr := field.GetInitialExpression()
		if initExpr == nil {
			continue
		}
		initExprBal := initExpr.(ast.BLangExpression)
		basePos := initExprBal.GetPosition()

		selfVarRef := &ast.BLangSimpleVarRef{
			VariableName: &ast.BLangIdentifier{Value: "self"},
		}
		selfVarRef.SetSymbol(selfRef)
		selfVarRef.SetDeterminedType(classType)

		fieldAccess := &ast.BLangFieldBaseAccess{
			Field: ast.BLangIdentifier{Value: field.GetName().GetValue()},
		}
		fieldAccess.Field.SetDeterminedType(semtypes.NEVER)
		fieldAccess.Expr = selfVarRef
		fieldAccess.SetDeterminedType(pkgCtx.getSymbolType(field.Symbol()))

		assignment := &ast.BLangAssignment{
			VarRef: fieldAccess,
			Expr:   initExprBal,
		}
		assignment.SetDeterminedType(semtypes.NEVER)
		setPositionIfMissing(assignment, basePos)

		initStmts = append(initStmts, assignment)
		field.SetInitialExpression(nil)
	}

	if len(initStmts) > 0 {
		body := class.InitFunction.Body.(*ast.BLangBlockFunctionBody)
		body.Stmts = append(initStmts, body.Stmts...)
	}
}

// desugarFunction returns a desugared function (may be same or new instance)
func desugarFunction(pkgCtx *packageContext, fn *ast.BLangFunction) *ast.BLangFunction {
	if fn.Body == nil {
		return fn
	}

	cx := &functionContext{
		pkgCtx: pkgCtx,
	}

	// Push function scope
	cx.pushScope(fn.Scope())
	defer cx.popScope()

	switch body := fn.Body.(type) {
	case *ast.BLangBlockFunctionBody:
		result := walkBlockFunctionBody(cx, body)
		if newBody, ok := result.replacementNode.(*ast.BLangBlockFunctionBody); ok {
			fn.Body = newBody
		}
	case *ast.BLangExprFunctionBody:
		if body.Expr != nil {
			result := walkExpression(cx, body.Expr)
			// For expression bodies, init statements need special handling
			// They should be converted to a block body with statements
			if len(result.initStmts) > 0 {
				fn.Body = convertExprBodyToBlockBody(body, result)
			} else {
				body.Expr = result.replacementNode.(ast.BLangExpression)
			}
		}
	case *ast.BLangExternFunctionBody:
		// Nothing to desugar
	}

	return fn
}

// convertExprBodyToBlockBody converts expression function body to block body
// when there are init statements from desugaring
func convertExprBodyToBlockBody(
	exprBody *ast.BLangExprFunctionBody,
	result desugaredNode[model.ExpressionNode],
) *ast.BLangBlockFunctionBody {
	// Create return statement with the desugared expression
	returnStmt := &ast.BLangReturn{
		Expr: result.replacementNode.(ast.BLangExpression),
	}

	// Build block with init statements + return
	stmts := make([]ast.BLangStatement, 0, len(result.initStmts)+1)
	stmts = append(stmts, result.initStmts...)
	stmts = append(stmts, returnStmt)

	return &ast.BLangBlockFunctionBody{
		Stmts: stmts,
	}
}
