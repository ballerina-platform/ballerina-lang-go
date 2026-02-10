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
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
	"reflect"
)

type analyzer interface {
	ast.Visitor
	ctx() *context.CompilerContext
	tyCtx() semtypes.Context
	importedPackage(alias string) *ast.BLangImportPackage
	unimplementedErr(message string)
	semanticErr(message string)
	syntaxErr(message string)
	internalErr(message string)
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

// FIXME: this is not correct since const analyzer will propagte to semantic analyzer
func returnFound(analyzer analyzer, returnStmt *ast.BLangReturn) {
	if analyzer == nil {
		panic("unexpected")
	}
	if fa, ok := analyzer.(*functionAnalyzer); ok {
		if returnStmt.Expr == nil {
			if !semtypes.IsSubtypeSimple(fa.retTy, semtypes.NIL) {
				fa.ctx().SemanticError("expect a return value", returnStmt.GetPosition())
			}
		}
		analyzeExpression(fa, returnStmt.Expr, fa.retTy)
	} else if analyzer.parentAnalyzer() != nil {
		returnFound(analyzer.parentAnalyzer(), returnStmt)
	} else {
		analyzer.ctx().SemanticError("return statement not allowed in this context", analyzer.loc())
	}
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
		returnFound(fa, n)
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

func (sa *SemanticAnalyzer) unimplementedErr(message string) {
	sa.compilerCtx.Unimplemented(message, nil)
}

func (sa *SemanticAnalyzer) semanticErr(message string) {
	sa.compilerCtx.SemanticError(message, nil)
}

func (sa *SemanticAnalyzer) syntaxErr(message string) {
	sa.compilerCtx.SyntaxError(message, nil)
}

func (sa *SemanticAnalyzer) internalErr(message string) {
	sa.compilerCtx.InternalError(message, nil)
}

func (ca *constantAnalyzer) unimplementedErr(message string) {
	ca.parentAnalyzer().ctx().Unimplemented(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) semanticErr(message string) {
	ca.parentAnalyzer().ctx().SemanticError(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) syntaxErr(message string) {
	ca.parentAnalyzer().ctx().SyntaxError(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) internalErr(message string) {
	ca.parentAnalyzer().ctx().InternalError(message, ca.constant.GetPosition())
}

func (fa *functionAnalyzer) unimplementedErr(message string) {
	fa.parent.ctx().Unimplemented(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) semanticErr(message string) {
	fa.parent.ctx().SemanticError(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) syntaxErr(message string) {
	fa.parent.ctx().SyntaxError(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) internalErr(message string) {
	fa.parent.ctx().InternalError(message, fa.function.GetPosition())
}

func (la *loopAnalyzer) unimplementedErr(message string) {
	la.parent.ctx().Unimplemented(message, la.loop.GetPosition())
}

func (la *loopAnalyzer) semanticErr(message string) {
	la.parent.ctx().SemanticError(message, la.loop.GetPosition())
}

func (la *loopAnalyzer) syntaxErr(message string) {
	la.parent.ctx().SyntaxError(message, la.loop.GetPosition())
}

func (la *loopAnalyzer) internalErr(message string) {
	la.parent.ctx().InternalError(message, la.loop.GetPosition())
}

// When we support multiple packages we need to resolve types of all of them before semantic analysis
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
	expectedType := constant.GetTypeData().Type
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
	case *ast.BLangReturn:
		// Error: return only valid in functions
		sa.semanticErr("return statement outside function")
		return nil
	case *ast.BLangWhile:
		// Error: loop only valid in functions
		sa.semanticErr("loop statement outside function")
		return nil
	case *ast.BLangIf:
		sa.semanticErr("if statement outside function")
		return nil
	default:
		// Now delegates function creation to visitInner
		return visitInner(sa, node)
	}
}

func (sa *SemanticAnalyzer) processImport(importNode *ast.BLangImportPackage) {
	alias := importNode.Alias.GetValue()

	// Only support ballerina/io
	if importNode.OrgName == nil || importNode.OrgName.GetValue() != "ballerina" {
		sa.unimplementedErr("unsupported import organization: only 'ballerina' imports are supported")
		return
	}

	if len(importNode.PkgNameComps) != 1 || importNode.PkgNameComps[0].GetValue() != "io" {
		sa.unimplementedErr("unsupported import package: only 'ballerina/io' is supported")
		return
	}

	// Check for duplicate imports
	if _, exists := sa.importedPkgs[alias]; exists {
		sa.semanticErr(fmt.Sprintf("import alias '%s' already defined", alias))
		return
	}

	sa.importedPkgs[alias] = importNode
}

func initializeFunctionAnalyzer(parent analyzer, function *ast.BLangFunction) *functionAnalyzer {
	fa := &functionAnalyzer{analyzerBase: analyzerBase{parent: parent}, function: function}
	fnSymbol := parent.ctx().GetSymbol(function.Symbol()).(*model.FunctionSymbol)
	fa.retTy = fnSymbol.Signature.ReturnType
	return fa
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
		ca.semanticErr("function definition not allowed in constant expression")
		return nil
	case *ast.BLangWhile:
		ca.semanticErr("loop not allowed in constant expression")
		return nil
	case *ast.BLangIf:
		ca.semanticErr("if statement not allowed in constant expression")
		return nil
	case *ast.BLangReturn:
		ca.semanticErr("return statement not allowed in constant expression")
		return nil
	case *ast.BLangBreak:
		ca.semanticErr("break statement not allowed in constant expression")
		return nil
	case *ast.BLangContinue:
		ca.semanticErr("continue statement not allowed in constant expression")
		return nil
	case *ast.BLangTypeDefinition:
		typeData := n.GetTypeData()
		expectedType := typeData.Type
		if expectedType == nil {
			ca.internalErr("type not resolved")
			return nil
		}
		ctx := ca.tyCtx()
		if semtypes.IsNever(expectedType) || !semtypes.IsSubtype(ctx, expectedType, semtypes.CreateAnydata(ctx)) {
			ca.syntaxErr("invalid type for constant declaration")
			return nil
		}
	case model.ExpressionNode:
		switch n.GetKind() {
		case model.NodeKind_LITERAL,
			model.NodeKind_NUMERIC_LITERAL,
			model.NodeKind_STRING_TEMPLATE_LITERAL,
			model.NodeKind_RECORD_LITERAL_EXPR,
			model.NodeKind_LIST_CONSTRUCTOR_EXPR,
			model.NodeKind_LIST_CONSTRUCTOR_SPREAD_OP,
			model.NodeKind_SIMPLE_VARIABLE_REF,
			model.NodeKind_BINARY_EXPR,
			model.NodeKind_GROUP_EXPR,
			model.NodeKind_UNARY_EXPR:
			bLangExpr := n.(ast.BLangExpression)
			analyzeExpression(ca, bLangExpr, ca.expectedType)
			exprTy := bLangExpr.GetTypeData().Type
			if ca.expectedType != nil {
				if !semtypes.IsSubtype(ca.tyCtx(), exprTy, ca.expectedType) {
					ca.semanticErr("incompatible type for constant expression")
					return nil
				}
			} else {
				ca.expectedType = exprTy
			}
		default:
			ca.semanticErr("expression is not a constant expression")
			return nil
		}
	}
	return ca
}

// validateResolvedType validates that a resolved expression type is compatible with the expected type
func validateResolvedType[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) bool {
	typeData := expr.GetTypeData()
	resolvedTy := typeData.Type
	if resolvedTy == nil {
		a.internalErr(fmt.Sprintf("expression type not resolved for %T at %v", expr, expr.GetPosition()))
		return false
	}

	if expectedType == nil {
		return true
	}

	ctx := a.tyCtx()
	if !semtypes.IsSubtype(ctx, resolvedTy, expectedType) {
		a.semanticErr(fmt.Sprintf("incompatible type: expected %v, got %v", expectedType, resolvedTy))
		return false
	}

	return true
}

func analyzeExpression[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) {
	switch expr := expr.(type) {
	// Literals - just validate
	case *ast.BLangLiteral:
		validateResolvedType(a, expr, expectedType)

	case *ast.BLangNumericLiteral:
		validateResolvedType(a, expr, expectedType)

	// Variable References - just validate
	case *ast.BLangSimpleVarRef:
		validateResolvedType(a, expr, expectedType)

	case *ast.BLangLocalVarRef, *ast.BLangConstRef:
		panic("not implemented")

	// Operators - validate operands and result
	case *ast.BLangBinaryExpr:
		analyzeBinaryExpr(a, expr, expectedType)

	case *ast.BLangUnaryExpr:
		analyzeUnaryExpr(a, expr, expectedType)

	// Function and Method Calls - validate arguments and result
	case *ast.BLangInvocation:
		analyzeInvocation(a, expr, expectedType)

	// Indexing - validate container, index, and result
	case *ast.BLangIndexBasedAccess:
		analyzeIndexBasedAccess(a, expr, expectedType)

	// Collections and Groups - validate members and result
	case *ast.BLangListConstructorExpr:
		analyzeListConstructorExpr(a, expr, expectedType)

	case *ast.BLangErrorConstructorExpr:
		analyzeErrorConstructorExpr(a, expr, expectedType)

	case *ast.BLangGroupExpr:
		analyzeExpression(a, expr.Expression, expectedType)

	case *ast.BLangWildCardBindingPattern:
		// Wildcard patterns have type ANY and are always valid
		validateResolvedType(a, expr, expectedType)

	default:
		a.internalErr("unexpected expression type: " + reflect.TypeOf(expr).String())
	}
}

func analyzeIndexBasedAccess[A analyzer](a A, expr *ast.BLangIndexBasedAccess, expectedType semtypes.SemType) {
	// Validate container expression
	containerExpr := expr.Expr
	analyzeExpression(a, containerExpr, nil)
	containerExprTy := containerExpr.GetTypeData().Type

	// Determine expected key type based on container type
	var keyExprExpectedType semtypes.SemType
	ctx := a.tyCtx()
	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) ||
		semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) ||
		semtypes.IsSubtypeSimple(containerExprTy, semtypes.XML) {
		keyExprExpectedType = &semtypes.INT
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.TABLE) {
		a.unimplementedErr("table not supported")
		return
	} else if semtypes.IsSubtype(ctx, containerExprTy, semtypes.Union(&semtypes.NIL, &semtypes.MAPPING)) {
		keyExprExpectedType = &semtypes.STRING
	} else {
		a.semanticErr("incompatible type for index based access")
		return
	}

	// Validate index expression
	keyExpr := expr.IndexExpr
	analyzeExpression(a, keyExpr, keyExprExpectedType)

	// Validate the resolved result type against expected type
	validateResolvedType(a, expr, expectedType)
}

func analyzeListConstructorExpr[A analyzer](a A, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) {
	// Validate each member expression
	for _, memberExpr := range expr.Exprs {
		analyzeExpression(a, memberExpr, nil)
	}

	// Validate the resolved list type against expected type
	validateResolvedType(a, expr, expectedType)
}

func analyzeErrorConstructorExpr[A analyzer](a A, expr *ast.BLangErrorConstructorExpr, expectedType semtypes.SemType) {
	argCount := len(expr.PositionalArgs)
	if argCount < 1 || argCount > 2 {
		a.semanticErr("error constructor must have at least 1 and at most 2 positional arguments")
		return
	}

	msgArg := expr.PositionalArgs[0]
	analyzeExpression(a, msgArg, &semtypes.STRING)

	// Validate second positional argument if present: must be a subtype of error? (Union of error and nil)
	if argCount == 2 {
		causeArg := expr.PositionalArgs[1]
		analyzeExpression(a, causeArg, semtypes.Union(&semtypes.ERROR, &semtypes.NIL))
	}

	// Validate the resolved error type against expected type
	validateResolvedType(a, expr, expectedType)
}

func analyzeUnaryExpr[A analyzer](a A, unaryExpr *ast.BLangUnaryExpr, expectedType semtypes.SemType) {
	// Validate the operand expression
	analyzeExpression(a, unaryExpr.Expr, nil)

	// Get the operand type
	exprTy := unaryExpr.Expr.GetTypeData().Type

	// Validate operand type based on operator
	switch unaryExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		if !isNumericType(exprTy) {
			a.semanticErr(fmt.Sprintf("expect numeric type for %s", string(unaryExpr.GetOperatorKind())))
			return
		}
	case model.OperatorKind_NOT:
		if !semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			a.semanticErr(fmt.Sprintf("expect boolean type for %s", string(unaryExpr.GetOperatorKind())))
			return
		}
	default:
		a.semanticErr(fmt.Sprintf("unsupported unary operator: %s", string(unaryExpr.GetOperatorKind())))
		return
	}

	// Validate the resolved result type against expected type
	validateResolvedType(a, unaryExpr, expectedType)
}

func analyzeBinaryExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, expectedType semtypes.SemType) {
	// Validate both operand expressions
	analyzeExpression(a, binaryExpr.LhsExpr, nil)
	analyzeExpression(a, binaryExpr.RhsExpr, nil)

	// Get operand types
	lhsTy := binaryExpr.LhsExpr.GetTypeData().Type
	rhsTy := binaryExpr.RhsExpr.GetTypeData().Type

	// Perform semantic validation based on operator type
	if isEqualityExpr(binaryExpr) {
		// For equality operators, ensure types have non-empty intersection
		intersection := semtypes.Intersect(lhsTy, rhsTy)
		if semtypes.IsEmpty(a.tyCtx(), intersection) {
			a.semanticErr(fmt.Sprintf("incompatible types for %s", string(binaryExpr.GetOperatorKind())))
			return
		}
	}

	// for nil lifting expression we do semantic analysis as part of type resolver
	// Validate the resolved result type against expected type
	validateResolvedType(a, binaryExpr, expectedType)
}

func analyzeInvocation[A analyzer](a A, invocation *ast.BLangInvocation, expectedType semtypes.SemType) {
	// Get the function type from the symbol
	symbol := invocation.Symbol()
	fnTy := a.ctx().SymbolType(symbol)
	if fnTy == nil || !semtypes.IsSubtypeSimple(fnTy, semtypes.FUNCTION) {
		a.semanticErr("function not found: " + invocation.Name.GetValue())
		return
	}

	// Validate each argument expression
	argTys := make([]semtypes.SemType, len(invocation.ArgExprs))
	for i, arg := range invocation.ArgExprs {
		analyzeExpression(a, arg, nil)
		typeData := arg.GetTypeData()
		argTys[i] = typeData.Type
	}

	// Validate argument types against function parameter types
	paramListTy := semtypes.FunctionParamListType(a.tyCtx(), fnTy)
	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(a.tyCtx().Env(), argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
	if !semtypes.IsSubtype(a.tyCtx(), argListTy, paramListTy) {
		a.semanticErr("incompatible arguments for function call")
		return
	}

	// Validate the resolved return type against expected type
	validateResolvedType(a, invocation, expectedType)
}

func analyzeSimpleVariableDef[A analyzer](a A, simpleVariableDef *ast.BLangSimpleVariableDef) {
	variable := simpleVariableDef.GetVariable().(*ast.BLangSimpleVariable)
	expectedType := variable.GetTypeData().Type
	if variable.Expr != nil {
		analyzeExpression(a, variable.Expr.(ast.BLangExpression), expectedType)
	}
	setExpectedType(simpleVariableDef, expectedType)
}

func visitInner[A analyzer](a A, node ast.BLangNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.BLangFunction:
		return initializeFunctionAnalyzer(a, n)
	case *ast.BLangWhile:
		analyzeWhile(a, n)
		return initializeLoopAnalyzer(a, n)
	case *ast.BLangIf:
		analyzeIf(a, n)
		return a
	case *ast.BLangBreak, *ast.BLangContinue:
		return nil
	case *ast.BLangSimpleVariableDef:
		analyzeSimpleVariableDef(a, n)
		return a
	case *ast.BLangAssignment:
		analyzeAssignment(a, n)
		return a
	case *ast.BLangCompoundAssignment:
		analyzeAssignment(a, n)
		return a
	case *ast.BLangExpressionStmt:
		analyzeExpression(a, n.Expr, &semtypes.NIL)
		return a
	case ast.BLangExpression:
		analyzeExpression(a, n, nil)
		return a
	case *ast.BLangReturn:
		returnFound(a, n)
		return nil
	default:
		return a
	}
}

type assignmentNode interface {
	GetVariable() model.ExpressionNode
	GetExpression() model.ExpressionNode
}

func analyzeAssignment[A analyzer](a A, assignment assignmentNode) {
	variable := assignment.GetVariable().(ast.BLangExpression)
	analyzeExpression(a, variable, nil)
	expectedType := variable.GetTypeData().Type
	expression := assignment.GetExpression().(ast.BLangExpression)
	analyzeExpression(a, expression, expectedType)
}

func analyzeIf[A analyzer](a A, ifStmt *ast.BLangIf) {
	analyzeExpression(a, ifStmt.Expr, &semtypes.BOOLEAN)
}

func analyzeWhile[A analyzer](a A, whileStmt *ast.BLangWhile) {
	analyzeExpression(a, whileStmt.Expr, &semtypes.BOOLEAN)
}

func setExpectedType[E ast.BLangNode](e E, expectedType semtypes.SemType) {
	typeData := e.GetTypeData()
	typeData.Type = expectedType
	e.SetTypeData(typeData)
	e.SetDeterminedType(expectedType)
}
