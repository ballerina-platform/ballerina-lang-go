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
	"reflect"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type analyzer interface {
	ast.Visitor
	ctx() *context.CompilerContext
	tyCtx() semtypes.Context
	getSymbol(ref model.SymbolRef) model.Symbol
	internalError(message string, loc diagnostics.Location)
	importedPackage(alias string) *ast.BLangImportPackage
	unimplementedErr(message string, loc diagnostics.Location)
	semanticErr(message string, loc diagnostics.Location)
	syntaxErr(message string, loc diagnostics.Location)
	internalErr(message string, loc diagnostics.Location)
	parentAnalyzer() analyzer
	loc() diagnostics.Location
}

type (
	analyzerBase struct {
		parent analyzer
	}
	SemanticAnalyzer struct {
		analyzerBase
		compilerCtx *context.CompilerContext
		typeCtx     semtypes.Context
		// TODO: move the constant resolution to type resolver as well so that we can run semantic analyzer in parallel as well
		pkg          *ast.BLangPackage
		importedPkgs map[string]*ast.BLangImportPackage
	}
	constantAnalyzer struct {
		analyzerBase
		constant     *ast.BLangConstant
		expectedType semtypes.SemType
	}

	functionAnalyzer struct {
		analyzerBase
		function *ast.BLangFunction
		retTy    semtypes.SemType
	}

	loopAnalyzer struct {
		analyzerBase
		loop ast.BLangNode
	}
)

var (
	_ analyzer = &constantAnalyzer{}
	_ analyzer = &SemanticAnalyzer{}
	_ analyzer = &functionAnalyzer{}
	_ analyzer = &loopAnalyzer{}
)

// expectedReturnType walks up the analyzer chain and returns the enclosing function's return type.
// Returns nil if not inside a function.
func expectedReturnType(a analyzer) semtypes.SemType {
	current := a
	for current != nil {
		if fa, ok := current.(*functionAnalyzer); ok {
			return fa.retTy
		}
		current = current.parentAnalyzer()
	}
	return nil
}

func returnFound(a analyzer, returnStmt *ast.BLangReturn) bool {
	retTy := expectedReturnType(a)
	if retTy == nil {
		a.ctx().SemanticError("return statement not allowed in this context", a.loc())
		return false
	}
	if returnStmt.Expr == nil {
		if !semtypes.IsSubtypeSimple(retTy, semtypes.NIL) {
			a.ctx().SemanticError("expect a return value", returnStmt.GetPosition())
			return false
		}
	} else if !analyzeExpression(a, returnStmt.Expr, retTy) {
		return false
	}
	return true
}

func (ab *analyzerBase) parentAnalyzer() analyzer {
	return ab.parent
}

func (ab *analyzerBase) importedPackage(alias string) *ast.BLangImportPackage {
	return ab.parentAnalyzer().importedPackage(alias)
}

func (ab *analyzerBase) ctx() *context.CompilerContext {
	return ab.parentAnalyzer().ctx()
}

func (ab *analyzerBase) getSymbol(ref model.SymbolRef) model.Symbol {
	return ab.ctx().GetSymbol(ref)
}

func (ab *analyzerBase) internalError(message string, loc diagnostics.Location) {
	ab.ctx().InternalError(message, loc)
}

func (ab *analyzerBase) tyCtx() semtypes.Context {
	return ab.parentAnalyzer().tyCtx()
}

func (sa *SemanticAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return nil
}

func (fa *functionAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return nil
}

func (la *loopAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return la
}

func (fa *functionAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangReturn:
		if !returnFound(fa, n) {
			return nil
		}
		return fa
	case *ast.BLangIdentifier:
		return nil
	default:
		// Delegate loop creation and common nodes to visitInner
		return visitInner(fa, node)
	}
}

func (la *loopAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	switch node.(type) {
	case *ast.BLangBreak, *ast.BLangContinue:
		return nil
	default:
		// Delegate nested loops and common nodes to visitInner
		return visitInner(la, node)
	}
}

func (fa *functionAnalyzer) loc() diagnostics.Location {
	return fa.function.GetPosition()
}

func (la *loopAnalyzer) loc() diagnostics.Location {
	return la.loop.GetPosition()
}

func (sa *SemanticAnalyzer) loc() diagnostics.Location {
	return sa.pkg.GetPosition()
}

func (ca *constantAnalyzer) loc() diagnostics.Location {
	return ca.constant.GetPosition()
}

func (ca *constantAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return ca
}

func (sa *SemanticAnalyzer) ctx() *context.CompilerContext {
	return sa.compilerCtx
}

func (sa *SemanticAnalyzer) tyCtx() semtypes.Context {
	return sa.typeCtx
}

func (sa *SemanticAnalyzer) importedPackage(alias string) *ast.BLangImportPackage {
	return sa.importedPkgs[alias]
}

func (la *loopAnalyzer) ctx() *context.CompilerContext {
	return la.parent.ctx()
}

func (la *loopAnalyzer) tyCtx() semtypes.Context {
	return la.parent.tyCtx()
}

func (sa *SemanticAnalyzer) unimplementedErr(message string, loc diagnostics.Location) {
	sa.compilerCtx.Unimplemented(message, loc)
}

func (sa *SemanticAnalyzer) semanticErr(message string, loc diagnostics.Location) {
	sa.compilerCtx.SemanticError(message, loc)
}

func (sa *SemanticAnalyzer) syntaxErr(message string, loc diagnostics.Location) {
	sa.compilerCtx.SyntaxError(message, loc)
}

func (sa *SemanticAnalyzer) internalErr(message string, loc diagnostics.Location) {
	sa.compilerCtx.InternalError(message, loc)
}

func (sa *SemanticAnalyzer) internalError(message string, loc diagnostics.Location) {
	sa.compilerCtx.InternalError(message, loc)
}

func (ca *constantAnalyzer) unimplementedErr(message string, loc diagnostics.Location) {
	ca.parentAnalyzer().ctx().Unimplemented(message, loc)
}

func (ca *constantAnalyzer) semanticErr(message string, loc diagnostics.Location) {
	ca.parentAnalyzer().ctx().SemanticError(message, loc)
}

func (ca *constantAnalyzer) syntaxErr(message string, loc diagnostics.Location) {
	ca.parentAnalyzer().ctx().SyntaxError(message, loc)
}

func (ca *constantAnalyzer) internalErr(message string, loc diagnostics.Location) {
	ca.parentAnalyzer().ctx().InternalError(message, loc)
}

func (fa *functionAnalyzer) unimplementedErr(message string, loc diagnostics.Location) {
	fa.parent.ctx().Unimplemented(message, loc)
}

func (fa *functionAnalyzer) semanticErr(message string, loc diagnostics.Location) {
	fa.parent.ctx().SemanticError(message, loc)
}

func (fa *functionAnalyzer) syntaxErr(message string, loc diagnostics.Location) {
	fa.parent.ctx().SyntaxError(message, loc)
}

func (fa *functionAnalyzer) internalErr(message string, loc diagnostics.Location) {
	fa.parent.ctx().InternalError(message, loc)
}

func (la *loopAnalyzer) unimplementedErr(message string, loc diagnostics.Location) {
	la.parent.ctx().Unimplemented(message, loc)
}

func (la *loopAnalyzer) semanticErr(message string, loc diagnostics.Location) {
	la.parent.ctx().SemanticError(message, loc)
}

func (la *loopAnalyzer) syntaxErr(message string, loc diagnostics.Location) {
	la.parent.ctx().SyntaxError(message, loc)
}

func (la *loopAnalyzer) internalErr(message string, loc diagnostics.Location) {
	la.parent.ctx().InternalError(message, loc)
}

func NewSemanticAnalyzer(ctx *context.CompilerContext) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		compilerCtx:  ctx,
		typeCtx:      semtypes.ContextFrom(ctx.GetTypeEnv()),
		importedPkgs: make(map[string]*ast.BLangImportPackage),
	}
}

func (sa *SemanticAnalyzer) Analyze(pkg *ast.BLangPackage) {
	sa.pkg = pkg
	sa.importedPkgs = make(map[string]*ast.BLangImportPackage)
	ast.Walk(sa, pkg)
	sa.pkg = nil
	sa.importedPkgs = nil
}

func createConstantAnalyzer(parent analyzer, constant *ast.BLangConstant) *constantAnalyzer {
	expectedType := constant.GetAssociatedType()
	return &constantAnalyzer{analyzerBase: analyzerBase{parent: parent}, constant: constant, expectedType: expectedType}
}

func (sa *SemanticAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangImportPackage:
		sa.processImport(n)
		return nil
	case *ast.BLangConstant:
		return createConstantAnalyzer(sa, n)
	case *ast.BLangSimpleVariable:
		// Module-level variables don't need constant-expression validation.
		// Type checking is handled by the type resolver.
		return nil
	case *ast.BLangReturn:
		// Error: return only valid in functions
		sa.semanticErr("return statement outside function", n.GetPosition())
		return nil
	case *ast.BLangWhile:
		// Error: loop only valid in functions
		sa.semanticErr("loop statement outside function", n.GetPosition())
		return nil
	case *ast.BLangIf:
		sa.semanticErr("if statement outside function", n.GetPosition())
		return nil
	default:
		// Now delegates function creation to visitInner
		return visitInner(sa, node)
	}
}

func (sa *SemanticAnalyzer) processImport(importNode *ast.BLangImportPackage) {
	alias := importNode.Alias.GetValue()

	// Check for duplicate imports
	if _, exists := sa.importedPkgs[alias]; exists {
		sa.semanticErr(fmt.Sprintf("import alias '%s' already defined", alias), importNode.GetPosition())
		return
	}

	sa.importedPkgs[alias] = importNode
}

func isIoImport(importNode *ast.BLangImportPackage) bool {
	return len(importNode.PkgNameComps) == 1 && importNode.PkgNameComps[0].GetValue() == "io"
}

func isImplicitImport(importNode *ast.BLangImportPackage) bool {
	return isLangImport(importNode, "array") || isLangImport(importNode, "int") || isLangImport(importNode, "map")
}

func isLangImport(importNode *ast.BLangImportPackage, name string) bool {
	return len(importNode.PkgNameComps) == 2 && importNode.PkgNameComps[0].GetValue() == "lang" && importNode.PkgNameComps[1].GetValue() == name
}

func validateInitFunction(parent analyzer, function *ast.BLangFunction, fnSymbol model.FunctionSymbol, pos diagnostics.Location) {
	if function.IsPublic() {
		parent.semanticErr("'init' function cannot be declared as public", pos)
	}

	expectedReturnType := semtypes.Union(semtypes.ERROR, semtypes.NIL)
	actualReturnType := fnSymbol.Signature().ReturnType
	if actualReturnType != nil && !semtypes.IsSubtype(parent.tyCtx(), actualReturnType, expectedReturnType) {
		parent.semanticErr("'init' function must have return type 'error?'", pos)
	}

	if len(function.RequiredParams) > 0 || function.RestParam != nil {
		parent.semanticErr("'init' function cannot have parameters", pos)
	}
}

func validateMainFunction(parent analyzer, fnSymbol model.FunctionSymbol, pos diagnostics.Location) {
	// Check 1: Must be public
	if !fnSymbol.IsPublic() {
		parent.semanticErr("'main' function must be public", pos)
	}

	// Check 2: Must return error?
	expectedReturnType := semtypes.Union(semtypes.ERROR, semtypes.NIL)
	actualReturnType := fnSymbol.Signature().ReturnType

	if actualReturnType != nil && !semtypes.IsSubtype(parent.tyCtx(), actualReturnType, expectedReturnType) {
		parent.semanticErr("'main' function must have return type 'error?'", pos)
	}
}

func initializeFunctionAnalyzer(parent analyzer, function *ast.BLangFunction) *functionAnalyzer {
	fa := initializeFunctionAnalyzerInner(parent, function)
	// Validate main function constraints
	if function.Name.Value == "main" {
		fnSymbol := parent.ctx().GetSymbol(function.Symbol()).(model.FunctionSymbol)
		validateMainFunction(parent, fnSymbol, function.GetPosition())
	}
	if function.Name.Value == "init" {
		// this is to seperate class init from module init
		if _, isTopLevel := parent.(*SemanticAnalyzer); isTopLevel {
			fnSymbol := parent.ctx().GetSymbol(function.Symbol()).(model.FunctionSymbol)
			validateInitFunction(parent, function, fnSymbol, function.GetPosition())
		}
	}

	return fa
}

func initializeFunctionAnalyzerInner(parent analyzer, function *ast.BLangFunction) *functionAnalyzer {
	fa := &functionAnalyzer{analyzerBase: analyzerBase{parent: parent}, function: function}
	fnSymbol := parent.ctx().GetSymbol(function.Symbol()).(model.FunctionSymbol)
	fa.retTy = fnSymbol.Signature().ReturnType
	validateDefaultParamTypes(parent, function)
	return fa
}

func validateDefaultParamTypes(a analyzer, function *ast.BLangFunction) {
	for i := range function.RequiredParams {
		param := &function.RequiredParams[i]
		if !param.IsDefaultableParam() {
			continue
		}
		paramTy := param.GetDeterminedType()
		exprTy := param.Expr.(ast.BLangExpression).GetDeterminedType()
		if exprTy == nil {
			a.internalErr("default expression has no determined type", param.Expr.(ast.BLangNode).GetPosition())
			continue
		}
		if !semtypes.IsSubtype(a.tyCtx(), exprTy, paramTy) {
			a.semanticErr("incompatible default value for parameter '"+param.Name.Value+"'", param.Expr.(ast.BLangNode).GetPosition())
		}
	}
}

func initializeMethodAnalyzer(parent analyzer, function *ast.BLangFunction) *functionAnalyzer {
	return initializeFunctionAnalyzerInner(parent, function)
}

func initializeLoopAnalyzer(parent analyzer, loop ast.BLangNode) *loopAnalyzer {
	return &loopAnalyzer{
		analyzerBase: analyzerBase{parent: parent},
		loop:         loop,
	}
}

func (ca *constantAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangIdentifier:
		return nil
	case *ast.BLangFunction:
		ca.semanticErr("function definition not allowed in constant expression", n.GetPosition())
		return nil
	case *ast.BLangWhile:
		ca.semanticErr("loop not allowed in constant expression", n.GetPosition())
		return nil
	case *ast.BLangIf:
		ca.semanticErr("if statement not allowed in constant expression", n.GetPosition())
		return nil
	case *ast.BLangReturn:
		ca.semanticErr("return statement not allowed in constant expression", n.GetPosition())
		return nil
	case *ast.BLangBreak:
		ca.semanticErr("break statement not allowed in constant expression", n.GetPosition())
		return nil
	case *ast.BLangContinue:
		ca.semanticErr("continue statement not allowed in constant expression", n.GetPosition())
		return nil
	case model.TypeDescriptor:
	case *ast.BLangTypeDefinition:
		// We have set the type at constructor
		return nil
	case model.ExpressionNode:
		bLangExpr := n.(ast.BLangExpression)
		hasErrors := false
		validateConstantExpr(ca.ctx(), bLangExpr, func(e ast.BLangExpression) {
			ca.semanticErr("expression is not a constant expression", e.GetPosition())
			hasErrors = true
		})
		if hasErrors {
			return nil
		}
		analyzeExpression(ca, bLangExpr, ca.expectedType)
		return nil
	}
	return ca
}

func validateConstantExpr(ctx *context.CompilerContext, expr ast.BLangExpression, onNonConst func(ast.BLangExpression)) {
	switch e := expr.(type) {
	case *ast.BLangLiteral, *ast.BLangNumericLiteral:
		// always valid
	case *ast.BLangSimpleVarRef:
		sym := ctx.GetSymbol(e.Symbol())
		if vs, ok := sym.(*model.ValueSymbol); ok && vs.IsConst() {
			return
		}
		onNonConst(expr)
	case *ast.BLangUnaryExpr:
		validateConstantExpr(ctx, e.Expr, onNonConst)
	case *ast.BLangGroupExpr:
		validateConstantExpr(ctx, e.Expression, onNonConst)
	case *ast.BLangBinaryExpr:
		validateConstantExpr(ctx, e.LhsExpr, onNonConst)
		validateConstantExpr(ctx, e.RhsExpr, onNonConst)
	case *ast.BLangListConstructorExpr:
		for _, member := range e.Exprs {
			validateConstantExpr(ctx, member, onNonConst)
		}
	case *ast.BLangMappingConstructorExpr:
		for _, field := range e.Fields {
			if kv, ok := field.(*ast.BLangMappingKeyValueField); ok {
				validateConstantExpr(ctx, kv.ValueExpr, onNonConst)
			}
		}
	default:
		onNonConst(expr)
	}
}

// validateResolvedType validates that a resolved expression type is compatible with the expected type
func validateResolvedType[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) bool {
	resolvedTy := expr.GetDeterminedType()
	if resolvedTy == nil {
		a.internalErr(fmt.Sprintf("expression type not resolved for %T", expr), expr.GetPosition())
		return false
	}

	if expectedType == nil {
		return true
	}

	ctx := a.tyCtx()
	if !semtypes.IsSubtype(ctx, resolvedTy, expectedType) {
		a.semanticErr(formatIncompatibleTypeMessage(ctx, expectedType, resolvedTy), expr.GetPosition())
		return false
	}
	if semtypes.IsNever(resolvedTy) {
		if !semtypes.IsNever(expectedType) {
			a.semanticErr(formatIncompatibleTypeMessage(ctx, expectedType, resolvedTy), expr.GetPosition())
			return false
		}
	}

	return true
}

func formatIncompatibleTypeMessage(ctx semtypes.Context, expectedType semtypes.SemType, actualType semtypes.SemType) string {
	return fmt.Sprintf("incompatible type: expected %s, got %s", semtypes.ToString(ctx, expectedType), semtypes.ToString(ctx, actualType))
}

func analyzeExpression[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) bool {
	switch expr := expr.(type) {
	case *ast.BLangLiteral:
		return validateResolvedType(a, expr, expectedType)

	case *ast.BLangNumericLiteral:
		return validateResolvedType(a, expr, expectedType)

	case *ast.BLangSimpleVarRef:
		return validateResolvedType(a, expr, expectedType)

	case *ast.BLangLocalVarRef, *ast.BLangConstRef:
		panic("not implemented")

	case *ast.BLangBinaryExpr:
		return analyzeBinaryExpr(a, expr, expectedType)

	case *ast.BLangUnaryExpr:
		return analyzeUnaryExpr(a, expr, expectedType)

	case *ast.BLangInvocation:
		return analyzeInvocation(a, expr, expectedType)

	case *ast.BLangIndexBasedAccess:
		return analyzeIndexBasedAccess(a, expr, expectedType)

	case *ast.BLangFieldBaseAccess:
		return analyzeFieldBasedAccess(a, expr, expectedType)
	// Collections and Groups - validate members and result
	case *ast.BLangListConstructorExpr:
		return analyzeListConstructorExpr(a, expr, expectedType)

	case *ast.BLangMappingConstructorExpr:
		return analyzeMappingConstructorExpr(a, expr, expectedType)

	case *ast.BLangErrorConstructorExpr:
		return analyzeErrorConstructorExpr(a, expr, expectedType)

	case *ast.BLangGroupExpr:
		return analyzeExpression(a, expr.Expression, expectedType)

	case *ast.BLangQueryExpr:
		return analyzeQueryExpr(a, expr, expectedType)

	case *ast.BLangWildCardBindingPattern:
		return validateResolvedType(a, expr, expectedType)

	case *ast.BLangTypeConversionExpr:
		return validateTypeConversionExpr(a, expr, expectedType)

	case *ast.BLangTypeTestExpr:
		return validateResolvedType(a, expr, expectedType)
	case *ast.BLangCheckedExpr:
		return analyzeCheckedExpr(a, expr, expectedType)
	case *ast.BLangCheckPanickedExpr:
		return analyzeCheckPanickedExpr(a, expr, expectedType)
	case *ast.BLangTrapExpr:
		return analyzeTrapExpr(a, expr, expectedType)
	case *ast.BLangNamedArgsExpression:
		return analyzeExpression(a, expr.Expr, expectedType)
	case *ast.BLangNewExpression:
		return analyzeNewExpression(a, expr, expectedType)
	case *ast.BLangLambdaFunction:
		return analyzeLambdaFunction(a, expr)
	default:
		a.internalErr("unexpected expression type: "+reflect.TypeOf(expr).String(), expr.GetPosition())
		return false
	}
}

func analyzeCheckedExpr[A analyzer](a A, expr *ast.BLangCheckedExpr, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, expr.Expr, nil) {
		return false
	}
	retTy := expectedReturnType(a)
	if retTy == nil {
		a.ctx().SemanticError("check expression not allowed outside a function", expr.GetPosition())
		return false
	}
	exprTy := expr.Expr.GetDeterminedType()
	errorPart := semtypes.Intersect(exprTy, semtypes.ERROR)
	if !semtypes.IsEmpty(a.tyCtx(), errorPart) {
		if !semtypes.IsSubtype(a.tyCtx(), errorPart, retTy) {
			a.ctx().SemanticError("error type of check expression is not a subtype of the enclosing function's return type", expr.GetPosition())
		}
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeTrapExpr[A analyzer](a A, expr *ast.BLangTrapExpr, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, expr.Expr, nil) {
		return false
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeCheckPanickedExpr[A analyzer](a A, expr *ast.BLangCheckPanickedExpr, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, expr.Expr, nil) {
		return false
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeQueryExpr[A analyzer](a A, queryExpr *ast.BLangQueryExpr, expectedType semtypes.SemType) bool {
	fromClause, ok := queryExpr.QueryClauseList[0].(*ast.BLangFromClause)
	if !ok {
		a.semanticErr("query expression must start with a from clause", queryExpr.GetPosition())
		return false
	}
	if !analyzeExpression(a, fromClause.Collection, nil) {
		return false
	}

	selectClause, ok := queryExpr.QueryClauseList[len(queryExpr.QueryClauseList)-1].(*ast.BLangSelectClause)
	if !ok {
		a.semanticErr("query expression requires a select clause", queryExpr.GetPosition())
		return false
	}

	for i := 1; i < len(queryExpr.QueryClauseList)-1; i++ {
		switch clause := queryExpr.QueryClauseList[i].(type) {
		case *ast.BLangLetClause:
			for _, variableDef := range clause.LetVarDeclarations {
				varDef, ok := variableDef.(*ast.BLangSimpleVariableDef)
				if !ok || varDef.Var == nil || varDef.Var.Expr == nil {
					a.semanticErr("let clause supports only initialized simple variable declarations", clause.GetPosition())
					return false
				}
				var expectedType semtypes.SemType
				if ast.SymbolIsSet(varDef.Var) {
					expectedType = a.ctx().SymbolType(varDef.Var.Symbol())
				}
				if !analyzeExpression(a, varDef.Var.Expr.(ast.BLangExpression), expectedType) {
					return false
				}
			}
		case *ast.BLangWhereClause, *ast.BLangLimitClause:
			// Query clause type and shape validation already happen in type resolution.
		}
	}

	var selectExpectedTy semtypes.SemType
	if queryExpr.QueryConstructType == model.TypeKind_MAP {
		selectExpectedTy = mapQuerySelectExpectedType(a.tyCtx().Env())
	}

	if !analyzeExpression(a, selectClause.Expression, selectExpectedTy) {
		return false
	}
	return validateResolvedType(a, queryExpr, expectedType)
}

func analyzeNewExpression[A analyzer](a A, expr *ast.BLangNewExpression, expectedType semtypes.SemType) bool {
	return validateResolvedType(a, expr, expectedType)
}

func analyzeLambdaFunction[A analyzer](a A, expr *ast.BLangLambdaFunction) bool {
	fa := initializeFunctionAnalyzer(a, expr.Function)
	ast.Walk(fa, expr.Function)
	return true
}

func validateTypeConversionExpr[A analyzer](a A, expr *ast.BLangTypeConversionExpr, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, expr.Expression, nil) {
		return false
	}
	exprTy := expr.Expression.GetDeterminedType()
	targetType := expr.TypeDescriptor.GetDeterminedType()
	intersection := semtypes.Intersect(exprTy, targetType)
	if semtypes.IsEmpty(a.tyCtx(), intersection) && !hasPotentialNumericConversions(exprTy, targetType) {
		a.semanticErr("impossible type conversion, intersection is empty", expr.GetPosition())
		return false
	}
	if expectedType != nil && !semtypes.IsSubtype(a.tyCtx(), targetType, expectedType) {
		a.semanticErr(formatIncompatibleTypeMessage(a.tyCtx(), expectedType, targetType), expr.GetPosition())
		return false
	}
	return validateResolvedType(a, expr, expectedType)
}

func hasPotentialNumericConversions(exprTy, targetType semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(exprTy, semtypes.NUMBER) && semtypes.SingleNumericType(targetType).IsPresent()
}

func analyzeFieldBasedAccess[A analyzer](a A, expr *ast.BLangFieldBaseAccess, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, expr.Expr, nil) {
		return false
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeIndexBasedAccess[A analyzer](a A, expr *ast.BLangIndexBasedAccess, expectedType semtypes.SemType) bool {
	// Validate container expression
	containerExpr := expr.Expr
	if !analyzeExpression(a, containerExpr, nil) {
		return false
	}
	containerExprTy := containerExpr.GetDeterminedType()

	var keyExprExpectedType semtypes.SemType
	ctx := a.tyCtx()
	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) ||
		semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) ||
		semtypes.IsSubtypeSimple(containerExprTy, semtypes.XML) {
		keyExprExpectedType = semtypes.INT
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.TABLE) {
		a.unimplementedErr("table not supported", expr.GetPosition())
		return false
	} else if semtypes.IsSubtype(ctx, containerExprTy, semtypes.Union(semtypes.NIL, semtypes.MAPPING)) {
		keyExprExpectedType = semtypes.STRING
	} else {
		a.semanticErr("incompatible type for index based access", expr.GetPosition())
		return false
	}

	keyExpr := expr.IndexExpr
	if !analyzeExpression(a, keyExpr, keyExprExpectedType) {
		return false
	}

	return validateResolvedType(a, expr, expectedType)
}

func analyzeListConstructorExpr[A analyzer](a A, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) bool {
	// The type resolver has already selected the inherent type and re-resolved members
	// with per-member expected types. We only need to validate members here.
	lat := expr.AtomicType
	for i, memberExpr := range expr.Exprs {
		memberExpectedType := lat.MemberAtInnerVal(i)
		if !analyzeExpression(a, memberExpr, memberExpectedType) {
			return false
		}
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeMappingConstructorExpr[A analyzer](a A, expr *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) bool {
	// The type resolver has already selected the inherent type and re-resolved field values
	// with per-field expected types. We only need to validate fields here.
	mat := expr.AtomicType
	hasValue := make(map[string]bool, len(expr.Fields)+len(expr.FieldDefaults))
	for _, fd := range expr.FieldDefaults {
		hasValue[fd.FieldName] = true
	}
	for _, f := range expr.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		keyName := recordKeyName(kv.Key)
		hasValue[keyName] = true
		fieldExpectedType := mat.FieldInnerVal(keyName)
		if !analyzeExpression(a, kv.ValueExpr, fieldExpectedType) {
			return false
		}
	}
	for _, name := range mat.Names {
		if hasValue[name] {
			continue
		}
		if mat.IsOptional(a.tyCtx(), name) {
			continue
		}
		a.semanticErr(fmt.Sprintf("missing non-defaultable required record field '%s'", name), expr.GetPosition())
		return false
	}
	return validateResolvedType(a, expr, expectedType)
}

func analyzeErrorConstructorExpr[A analyzer](a A, expr *ast.BLangErrorConstructorExpr, expectedType semtypes.SemType) bool {
	argCount := len(expr.PositionalArgs)
	if argCount < 1 || argCount > 2 {
		a.semanticErr("error constructor must have at least 1 and at most 2 positional arguments", expr.GetPosition())
		return false
	}
	tyCtx := a.tyCtx()

	msgArg := expr.PositionalArgs[0]
	if !analyzeExpression(a, msgArg, semtypes.STRING) {
		return false
	}
	mat, ok := semtypes.ErrorDetailAtomicType(tyCtx, expr.DeterminedType)
	if !ok {
		a.unimplementedErr("non-atomic detail types not supported", expr.GetPosition())
		return false
	}
	seen := make(map[string]bool, len(expr.NamedArgs))
	clonableTy := semtypes.CreateCloneable(tyCtx)
	for _, namedArg := range expr.NamedArgs {
		name := namedArg.Name.GetValue()
		if seen[name] {
			a.semanticErr(fmt.Sprintf("duplicate named argument '%s' in error constructor", name), namedArg.GetPosition())
			return false
		}
		seen[name] = true
		fieldType := mat.FieldInnerVal(name)
		if !analyzeExpression(a, namedArg.Expr, fieldType) {
			return false
		}
		if !semtypes.IsSubtype(tyCtx, namedArg.Expr.GetDeterminedType(), clonableTy) {
			a.semanticErr("named arguments must be subtypes of cloneable", namedArg.GetPosition())
			return false
		}
	}

	// Every field in the atom must be provided
	for _, name := range mat.Names {
		if !seen[name] {
			a.semanticErr(fmt.Sprintf("missing required field '%s' in error constructor", name), expr.GetPosition())
			return false
		}
	}

	if argCount == 2 {
		causeArg := expr.PositionalArgs[1]
		if !analyzeExpression(a, causeArg, semtypes.Union(semtypes.ERROR, semtypes.NIL)) {
			return false
		}
	}

	return validateResolvedType(a, expr, expectedType)
}

func analyzeUnaryExpr[A analyzer](a A, unaryExpr *ast.BLangUnaryExpr, expectedType semtypes.SemType) bool {
	if !analyzeExpression(a, unaryExpr.Expr, nil) {
		return false
	}

	exprTy := unaryExpr.Expr.GetDeterminedType()
	// Strip nil for nil-lifted numeric/bitwise unary operations
	underlyingTy := exprTy
	if semtypes.ContainsBasicType(exprTy, semtypes.NIL) {
		underlyingTy = semtypes.Diff(exprTy, semtypes.NIL)
	}

	switch unaryExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		if !isNumericType(underlyingTy) {
			a.semanticErr(fmt.Sprintf("expect numeric type for %s", string(unaryExpr.GetOperatorKind())), unaryExpr.GetPosition())
			return false
		}
	case model.OperatorKind_NOT:
		if !semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			a.semanticErr(fmt.Sprintf("expect boolean type for %s", string(unaryExpr.GetOperatorKind())), unaryExpr.GetPosition())
			return false
		}
	default:
		a.semanticErr(fmt.Sprintf("unsupported unary operator: %s", string(unaryExpr.GetOperatorKind())), unaryExpr.GetPosition())
		return false
	}

	return validateResolvedType(a, unaryExpr, expectedType)
}

func analyzeBinaryExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, expectedType semtypes.SemType) bool {
	// Validate both operand expressions
	if !analyzeExpression(a, binaryExpr.LhsExpr, nil) {
		return false
	}
	if !analyzeExpression(a, binaryExpr.RhsExpr, nil) {
		return false
	}

	// Get operand types
	lhsTy := binaryExpr.LhsExpr.GetDeterminedType()
	rhsTy := binaryExpr.RhsExpr.GetDeterminedType()

	ctx := a.tyCtx()
	// Perform semantic validation based on operator type
	if isEqualityExpr(binaryExpr) {
		// For equality operators, ensure types have non-empty intersection
		intersection := semtypes.Intersect(lhsTy, rhsTy)
		if semtypes.IsEmpty(ctx, intersection) {
			a.semanticErr(fmt.Sprintf("incompatible types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
			return false
		}
		switch binaryExpr.GetOperatorKind() {
		case model.OperatorKind_EQUAL, model.OperatorKind_NOT_EQUAL:
			anyData := semtypes.CreateAnydata(ctx)
			if !semtypes.IsSubtype(ctx, lhsTy, anyData) && !semtypes.IsSubtype(ctx, rhsTy, anyData) {
				a.semanticErr(fmt.Sprintf("expect anydata types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
				return false
			}
		}
	} else if isBitWiseExpr(binaryExpr) {
		if !analyzeBitWiseExpr(a, binaryExpr, lhsTy, rhsTy) {
			return false
		}
	} else if isRangeExpr(binaryExpr) {
		if !semtypes.IsSubtypeSimple(lhsTy, semtypes.INT) || !semtypes.IsSubtypeSimple(rhsTy, semtypes.INT) {
			a.semanticErr(fmt.Sprintf("expect int types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
			return false
		}
	} else if isShiftExpr(binaryExpr) {
		if !analyzeShiftExpr(a, lhsTy, rhsTy) {
			return false
		}
	} else if isLogicalExpression(binaryExpr) {
		if !semtypes.IsSubtypeSimple(lhsTy, semtypes.BOOLEAN) || !semtypes.IsSubtypeSimple(rhsTy, semtypes.BOOLEAN) {
			a.semanticErr(fmt.Sprintf("expect boolean types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
			return false
		}
	}
	// for nil lifting expression we do semantic analysis as part of type resolver
	// Validate the resolved result type against expected type
	return validateResolvedType(a, binaryExpr, expectedType)
}

func analyzeBitWiseExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, lhsTy, rhsTy semtypes.SemType) bool {
	ctx := a.tyCtx()
	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		lhsTy = semtypes.Diff(lhsTy, semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, semtypes.NIL)
	}
	if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
		a.semanticErr("expect integer types for bitwise operators", binaryExpr.GetPosition())
		return false
	}
	return true
}

func analyzeShiftExpr[A analyzer](a A, lhsTy, rhsTy semtypes.SemType) bool {
	ctx := a.tyCtx()
	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		lhsTy = semtypes.Diff(lhsTy, semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, semtypes.NIL) //nolint:staticcheck,ineffassign // rhsTy will be used when nil-lifted binary ops are fully implemented
	}
	if !semtypes.IsSubtype(ctx, lhsTy, semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, semtypes.INT) {
		return false
	}
	return true
}

func analyzeInvocation[A analyzer](a A, invocation *ast.BLangInvocation, expectedType semtypes.SemType) bool {
	// Get the function type from the symbol
	symbol := invocation.Symbol()
	fnTy := a.ctx().SymbolType(symbol)
	paramListTy := semtypes.FunctionParamListType(a.tyCtx(), fnTy)

	fnSymbol, isDirectCall := a.ctx().GetSymbol(symbol).(model.FunctionSymbol)
	// TODO: ideally we need to unify these when we no longer has restrictions on lambdas
	if !isDirectCall {
		return analyzeLambdaInvocation(a, invocation, paramListTy, expectedType)
	}
	return analyzeDirectInvocation(a, invocation, fnSymbol, paramListTy, expectedType)
}

func analyzeDirectInvocation[A analyzer](a A, invocation *ast.BLangInvocation, fnSymbol model.FunctionSymbol, paramListTy, expectedType semtypes.SemType) bool {
	signature := fnSymbol.Signature()
	tyCtx := a.tyCtx()
	for i, arg := range invocation.ArgExprs {
		switch arg := arg.(type) {
		case *ast.BLangNamedArgsExpression:
			name := arg.Name.Value
			targetIndex := -1
			for j, each := range signature.ParamNames {
				if each == name {
					targetIndex = j
					break
				}
			}
			key := semtypes.IntConst(int64(targetIndex))
			if !analyzeExpression(a, arg, semtypes.ListMemberTypeInnerVal(tyCtx, paramListTy, key)) {
				return false
			}
		default:
			key := semtypes.IntConst(int64(i))
			if !analyzeExpression(a, arg, semtypes.ListMemberTypeInnerVal(tyCtx, paramListTy, key)) {
				return false
			}
		}
	}

	// Validate the resolved return type against expected type
	return validateResolvedType(a, invocation, expectedType)
}

func analyzeLambdaInvocation[A analyzer](a A, invocation *ast.BLangInvocation, paramListTy, expectedType semtypes.SemType) bool {
	tyCtx := a.tyCtx()

	// Validate each argument expression
	for i, arg := range invocation.ArgExprs {
		key := semtypes.IntConst(int64(i))
		if !analyzeExpression(a, arg, semtypes.ListMemberTypeInnerVal(tyCtx, paramListTy, key)) {
			return false
		}
	}

	// Validate the resolved return type against expected type
	return validateResolvedType(a, invocation, expectedType)
}

func analyzeSimpleVariableDef[A analyzer](a A, simpleVariableDef *ast.BLangSimpleVariableDef) bool {
	variable := simpleVariableDef.GetVariable().(*ast.BLangSimpleVariable)
	expectedType := variable.GetDeterminedType()
	if variable.GetName().GetValue() == string(model.IGNORE) {
		if !semtypes.IsSubtypeSimple(expectedType, semtypes.ANY) {
			a.semanticErr("wildcard binding pattern type must be a subtype of 'any'", variable.GetPosition())
			return false
		}
	}
	if ast.SymbolIsSet(variable) {
		symbolType := a.ctx().SymbolType(variable.Symbol())
		if symbolType != nil {
			expectedType = symbolType
		}
	}
	if variable.Expr != nil && !analyzeExpression(a, variable.Expr.(ast.BLangExpression), expectedType) {
		return false
	}
	setExpectedType(simpleVariableDef, expectedType)
	return true
}

func visitInner[A analyzer](a A, node ast.BLangNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.BLangFunction:
		return initializeFunctionAnalyzer(a, n)
	case *ast.BLangWhile:
		if !analyzeWhile(a, n) {
			return nil
		}
		return initializeLoopAnalyzer(a, n)
	case *ast.BLangForeach:
		if !validateForeach(a, n) {
			return nil
		}
		return initializeLoopAnalyzer(a, n)
	case *ast.BLangIf:
		if !analyzeIf(a, n) {
			return nil
		}
		return a
	case *ast.BLangBreak, *ast.BLangContinue:
		return nil
	case *ast.BLangMatchStatement:
		return a
	case *ast.BLangSimpleVariableDef:
		if !analyzeSimpleVariableDef(a, n) {
			return nil
		}
		return a
	case *ast.BLangAssignment:
		if !analyzeAssignment(a, n) {
			return nil
		}
		return a
	case *ast.BLangCompoundAssignment:
		if !analyzeAssignment(a, n) {
			return nil
		}
		return a
	case *ast.BLangExpressionStmt:
		if !analyzeExpression(a, n.Expr, nil) {
			return nil
		}
		exprType := n.Expr.GetDeterminedType()
		if !semtypes.IsSubtype(a.tyCtx(), exprType, semtypes.NIL) {
			a.semanticErr("expression value must be assigned", n.Expr.GetPosition())
			return nil
		}
		return a
	case ast.BLangExpression:
		if !analyzeExpression(a, n, nil) {
			return nil
		}
		return a
	case *ast.BLangReturn:
		if !returnFound(a, n) {
			return nil
		}
		return nil
	case *ast.BLangPanic:
		analyzeExpression(a, n.Expr, semtypes.ERROR)
		return nil
	case *ast.BLangRecordType:
		validateRecordFieldDefaults(a, n)
		return nil
	case *ast.BLangClassDefinition:
		for _, f := range n.Fields {
			field := f.(*ast.BLangSimpleVariable)
			if field.Expr != nil {
				expectedType := a.ctx().SymbolType(field.Symbol())
				analyzeExpression(a, field.Expr.(ast.BLangExpression), expectedType)
			}
		}
		if n.InitFunction != nil {
			fa := initializeMethodAnalyzer(a, n.InitFunction)
			ast.Walk(fa, n.InitFunction)
		}
		for name := range n.Methods {
			method := n.Methods[name]
			fa := initializeMethodAnalyzer(a, method)
			ast.Walk(fa, method)
		}
		return nil
	default:
		return a
	}
}

type assignmentNode interface {
	GetVariable() model.ExpressionNode
	GetExpression() model.ExpressionNode
}

func analyzeAssignment[A analyzer](a A, assignment assignmentNode) bool {
	variable := assignment.GetVariable().(ast.BLangExpression)
	if symbolNode, ok := variable.(ast.BNodeWithSymbol); ok {
		symbol := symbolNode.Symbol()
		if !ast.SymbolIsSet(symbolNode) {
			a.internalErr("unexpected nil symbol", variable.GetPosition())
			return false
		}
		ctx := a.ctx()
		switch ctx.SymbolKind(symbol) {
		case model.SymbolKindConstant:
			a.semanticErr("cannot assign to constant", variable.GetPosition())
			return false
		case model.SymbolKindParemeter:
			a.semanticErr("cannot assign to parameter", variable.GetPosition())
			return false
		case model.SymbolKindFunction:
			a.semanticErr("cannot assign to function", variable.GetPosition())
			return false
		case model.SymbolKindType:
			a.semanticErr("cannot assign to type", variable.GetPosition())
			return false
		}
	}
	if !analyzeExpression(a, variable, nil) {
		return false
	}
	expectedType := variable.GetDeterminedType()
	expression := assignment.GetExpression().(ast.BLangExpression)
	return analyzeExpression(a, expression, expectedType)
}

func analyzeIf[A analyzer](a A, ifStmt *ast.BLangIf) bool {
	return analyzeExpression(a, ifStmt.Expr, semtypes.BOOLEAN)
}

func analyzeWhile[A analyzer](a A, whileStmt *ast.BLangWhile) bool {
	return analyzeExpression(a, whileStmt.Expr, semtypes.BOOLEAN)
}

func validateForeach[A analyzer](a A, foreachStmt *ast.BLangForeach) bool {
	collection := foreachStmt.Collection
	if !analyzeExpression(a, collection, nil) {
		return false
	}
	variable := foreachStmt.VariableDef.GetVariable().(*ast.BLangSimpleVariable)
	variableType := a.ctx().SymbolType(variable.Symbol())
	if binExpr, ok := collection.(*ast.BLangBinaryExpr); ok && isRangeExpr(binExpr) {
		if !semtypes.IsSubtypeSimple(variableType, semtypes.INT) {
			a.semanticErr("foreach variable must be a subtype of int for range expression", collection.GetPosition())
			return false
		}
	} else {
		collectionType := collection.GetDeterminedType()
		var expectedValueType semtypes.SemType
		switch {
		case semtypes.IsSubtypeSimple(collectionType, semtypes.LIST):
			memberTypes := semtypes.ListAllMemberTypesInner(a.tyCtx(), collectionType)
			var result semtypes.SemType = semtypes.NEVER
			for _, each := range memberTypes.SemTypes {
				result = semtypes.Union(result, each)
			}
			expectedValueType = result
		case semtypes.IsSubtypeSimple(collectionType, semtypes.MAPPING):
			expectedValueType = semtypes.MappingMemberTypeInnerVal(a.tyCtx(), collectionType, semtypes.STRING)
		default:
			a.unimplementedErr("unsupported foreach collection", collection.GetPosition())
			return false
		}
		if !semtypes.IsSubtype(a.tyCtx(), expectedValueType, variableType) {
			a.ctx().SemanticError("invalid type for variable", variable.GetPosition())
			return false
		}
	}
	return true
}

func recordKeyName(key *ast.BLangMappingKey) string {
	switch expr := key.Expr.(type) {
	case *ast.BLangLiteral:
		return expr.Value.(string)
	case *ast.BLangSimpleVarRef:
		return expr.VariableName.Value
	default:
		panic(fmt.Sprintf("unexpected record key expression type: %T", key.Expr))
	}
}

func setExpectedType[E ast.BLangNode](e E, expectedType semtypes.SemType) {
	e.SetDeterminedType(expectedType)
}

// validateRecordFieldDefaults checks that all record field default expressions satisfy
// isolation rules: all variable references must be const, regardless of scope.
// Record field defaults are evaluated at record construction time and must not capture
// mutable state from any scope.
func validateRecordFieldDefaults[A analyzer](a A, node *ast.BLangRecordType) {
	for _, field := range node.Fields() {
		if field.DefaultExpr != nil {
			if !isIsolatedFuncInner(a, nil, field.DefaultExpr.(ast.BLangNode)) {
				a.semanticErr("not an isolated expression", field.DefaultExpr.(ast.BLangNode).GetPosition())
			}
		}
	}
}

// ancestorSpaceIndices collects the SpaceIndex values of all scopes above the given
// function scope (its parent chain up to the module scope). Variables from these scopes
// are non-local (module-level or captured) and must be const in an isolated function.
// If funcScope is nil (module-level context), returns nil to signal that all refs must be const.
// NOTE: this works with narrowing because we create narrowed symbol in the same space as the symbol being narrowed.
func ancestorSpaceIndices(funcScope *model.FunctionScope) map[int]struct{} {
	if funcScope == nil {
		return nil
	}
	ancestors := make(map[int]struct{})
	scope := funcScope.Parent
	for scope != nil {
		switch s := scope.(type) {
		case *model.ModuleScope:
			ancestors[s.MainSpace().SpaceIndex()] = struct{}{}
			return ancestors
		case *model.FunctionScope:
			ancestors[s.MainSpace().SpaceIndex()] = struct{}{}
			scope = s.Parent
		case *model.BlockScope:
			ancestors[s.MainSpace().SpaceIndex()] = struct{}{}
			scope = s.Parent
		default:
			return ancestors
		}
	}
	return ancestors
}

// TODO: Make this generic over Expressions and statements
func isIsolatedFuncInner[A analyzer](a A, funcScope *model.FunctionScope, node ast.BLangNode) bool {
	ancestors := ancestorSpaceIndices(funcScope)
	return everyNode(a, node, func(analyzer A, inner ast.BLangNode) bool {
		switch inner := inner.(type) {
		case *ast.BLangInvocation:
			analyzer.unimplementedErr("isolated functions not implemented", inner.GetPosition())
			return false
		case *ast.BLangSimpleVarRef:
			sym := a.ctx().GetSymbol(inner.Symbol())
			if varSym, ok := sym.(*model.ValueSymbol); ok {
				if ancestors == nil {
					return varSym.IsConst()
				}
				if _, isAncestor := ancestors[inner.Symbol().SpaceIndex]; isAncestor {
					return varSym.IsConst()
				}
				return true
			} else {
				analyzer.unimplementedErr("unsupported reference in isolated function body", inner.GetPosition())
				return false
			}
		default:
			return true
		}
	})
}

type everyNodeVisitor[A analyzer] struct {
	analyzer  A
	predicate func(A, ast.BLangNode) bool
	result    bool
}

func (v *everyNodeVisitor[A]) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return v
	}
	if !v.predicate(v.analyzer, node) {
		v.result = false
		return nil
	}
	return v
}

func (v *everyNodeVisitor[A]) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return v
}

func everyNode[A analyzer](a A, node ast.BLangNode, predicate func(A, ast.BLangNode) bool) bool {
	visitor := &everyNodeVisitor[A]{analyzer: a, predicate: predicate, result: true}
	ast.Walk(visitor, node)
	return visitor.result
}
