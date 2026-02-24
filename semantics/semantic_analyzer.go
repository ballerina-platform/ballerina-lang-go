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
	"reflect"
	"sort"

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

// FIXME: this is not correct since const analyzer will propagte to semantic analyzer
func returnFound(analyzer analyzer, returnStmt *ast.BLangReturn) bool {
	if analyzer == nil {
		panic("unexpected")
	}
	if fa, ok := analyzer.(*functionAnalyzer); ok {
		if returnStmt.Expr == nil {
			if !semtypes.IsSubtypeSimple(fa.retTy, semtypes.NIL) {
				fa.ctx().SemanticError("expect a return value", returnStmt.GetPosition())
				return false
			}
		} else if !analyzeExpression(fa, returnStmt.Expr, fa.retTy) {
			return false
		}
	} else if analyzer.parentAnalyzer() != nil {
		return returnFound(analyzer.parentAnalyzer(), returnStmt)
	} else {
		analyzer.ctx().SemanticError("return statement not allowed in this context", analyzer.loc())
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
	expectedType := constant.GetDeterminedType()
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

	// Only support ballerina/io
	if importNode.OrgName == nil || importNode.OrgName.GetValue() != "ballerina" {
		sa.unimplementedErr("unsupported import organization: only 'ballerina' imports are supported", importNode.GetPosition())
		return
	}

	if !isIoImport(importNode) && !isImplicitImport(importNode) {
		sa.unimplementedErr("unsupported import package: only 'ballerina/io' is supported", importNode.GetPosition())
		return
	}

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
	return isLangImport(importNode, "array") || isLangImport(importNode, "int")
}

func isLangImport(importNode *ast.BLangImportPackage, name string) bool {
	return len(importNode.PkgNameComps) == 2 && importNode.PkgNameComps[0].GetValue() == "lang" && importNode.PkgNameComps[1].GetValue() == name
}

func validateMainFunction(parent analyzer, function *ast.BLangFunction, fnSymbol model.FunctionSymbol) {
	// Check 1: Must be public
	if !fnSymbol.IsPublic() {
		parent.semanticErr("'main' function must be public", function.GetPosition())
	}

	// Check 2: Must return error?
	expectedReturnType := semtypes.Union(&semtypes.ERROR, &semtypes.NIL)
	actualReturnType := fnSymbol.Signature().ReturnType

	if actualReturnType != nil && !semtypes.IsSubtype(parent.tyCtx(), actualReturnType, expectedReturnType) {
		parent.semanticErr("'main' function must have return type 'error?'", function.GetPosition())
	}
}

func initializeFunctionAnalyzer(parent analyzer, function *ast.BLangFunction) *functionAnalyzer {
	fa := &functionAnalyzer{analyzerBase: analyzerBase{parent: parent}, function: function}
	fnSymbol := parent.ctx().GetSymbol(function.Symbol()).(model.FunctionSymbol)
	fa.retTy = fnSymbol.Signature().ReturnType

	// Validate main function constraints
	if function.Name.Value == "main" {
		validateMainFunction(parent, function, fnSymbol)
	}

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
	case *ast.BLangTypeDefinition:
		typeData := n.GetTypeData()
		expectedType := typeData.Type
		if expectedType == nil {
			ca.internalErr("type not resolved", n.GetPosition())
			return nil
		}
		ctx := ca.tyCtx()
		if semtypes.IsNever(expectedType) || !semtypes.IsSubtype(ctx, expectedType, semtypes.CreateAnydata(ctx)) {
			ca.syntaxErr("invalid type for constant declaration", n.GetPosition())
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
			if !analyzeExpression(ca, bLangExpr, ca.expectedType) {
				return nil
			}
			exprTy := bLangExpr.GetDeterminedType()
			if ca.expectedType != nil {
				if !semtypes.IsSubtype(ca.tyCtx(), exprTy, ca.expectedType) {
					ca.semanticErr("incompatible type for constant expression", bLangExpr.GetPosition())
					return nil
				}
			} else {
				ca.expectedType = exprTy
			}
		default:
			ca.semanticErr("expression is not a constant expression", n.GetPosition())
			return nil
		}
	}
	return ca
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
		a.semanticErr(fmt.Sprintf("incompatible type: expected %v, got %v", expectedType, resolvedTy), expr.GetPosition())
		return false
	}
	if semtypes.IsNever(resolvedTy) {
		if !semtypes.IsNever(expectedType) {
			a.semanticErr(fmt.Sprintf("incompatible type: expected %v, got %v", expectedType, resolvedTy), expr.GetPosition())
			return false
		}
	}

	return true
}

func widenNumericLiteral[A analyzer](a A, expr *ast.BLangLiteral, expectedType semtypes.SemType) {
	if expectedType == nil {
		return
	}
	resolvedTy := expr.GetDeterminedType()
	// If already compatible, no widening needed
	if semtypes.IsSubtype(a.tyCtx(), resolvedTy, expectedType) {
		return
	}
	// Determine the single target numeric type
	singleNumType := semtypes.SingleNumericType(expectedType)
	if !singleNumType.IsPresent() {
		return
	}
	targetNumType := singleNumType.Get()

	// int → float
	if targetNumType == semtypes.FLOAT {
		if intVal, ok := expr.Value.(int64); ok {
			floatVal := float64(intVal)
			expr.Value = floatVal
			expr.SetDeterminedType(semtypes.FloatConst(floatVal))
		}
		return
	}
	// int → decimal OR float → decimal
	if targetNumType == semtypes.DECIMAL {
		if intVal, ok := expr.Value.(int64); ok {
			ratVal := new(big.Rat).SetInt64(intVal)
			expr.Value = ratVal
			expr.SetDeterminedType(semtypes.DecimalConst(*ratVal))
		} else if floatVal, ok := expr.Value.(float64); ok {
			ratVal := new(big.Rat).SetFloat64(floatVal)
			expr.Value = ratVal
			expr.SetDeterminedType(semtypes.DecimalConst(*ratVal))
		}
		return
	}
}

func analyzeExpression[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) bool {
	switch expr := expr.(type) {
	case *ast.BLangLiteral:
		widenNumericLiteral(a, expr, expectedType)
		return analyzeLiteral(a, expr, expectedType)

	case *ast.BLangNumericLiteral:
		widenNumericLiteral(a, &expr.BLangLiteral, expectedType)
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

	case *ast.BLangWildCardBindingPattern:
		return validateResolvedType(a, expr, expectedType)

	case *ast.BLangTypeConversionExpr:
		return validateTypeConversionExpr(a, expr, expectedType)

	case *ast.BLangTypeTestExpr:
		return validateResolvedType(a, expr, expectedType)

	default:
		a.internalErr("unexpected expression type: "+reflect.TypeOf(expr).String(), expr.GetPosition())
		return false
	}
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
		a.semanticErr(fmt.Sprintf("incompatible type: expected %v, got %v", expectedType, exprTy), expr.GetPosition())
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
		keyExprExpectedType = &semtypes.INT
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.TABLE) {
		a.unimplementedErr("table not supported", expr.GetPosition())
		return false
	} else if semtypes.IsSubtype(ctx, containerExprTy, semtypes.Union(&semtypes.NIL, &semtypes.MAPPING)) {
		keyExprExpectedType = &semtypes.STRING
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
	for _, memberExpr := range expr.Exprs {
		if !analyzeExpression(a, memberExpr, nil) {
			return false
		}
	}

	if expectedType != nil {
		resultType, listAtomicType := selectListInherentType(a, expr, expectedType)
		if resultType == nil {
			return false
		}
		for i, memberExpr := range expr.Exprs {
			requiredType := listAtomicType.MemberAtInnerVal(i)
			if semtypes.IsNever(requiredType) {
				a.semanticErr("too many members in list constructor", expr.GetPosition())
				return false
			}
			if !analyzeExpression(a, memberExpr, requiredType) {
				return false
			}
		}
		expr.AtomicType = listAtomicType
		setExpectedType(expr, resultType)
	} else {
		for _, memberExpr := range expr.Exprs {
			if !analyzeExpression(a, memberExpr, nil) {
				return false
			}
		}
	}
	return validateResolvedType(a, expr, expectedType)
}

func selectListInherentType[A analyzer](a A, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, semtypes.ListAtomicType) {
	expectedListType := semtypes.Intersect(expectedType, &semtypes.LIST)
	tc := a.tyCtx()
	if semtypes.IsEmpty(tc, expectedListType) {
		a.semanticErr("list type not found in expected type", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}
	}
	lat := semtypes.ToListAtomicType(tc, expectedListType)
	if lat != nil {
		return expectedListType, *lat
	}

	alts := semtypes.ListAlternatives(tc, expectedListType)

	// Filter alternatives by length compatibility
	var validAlts []semtypes.ListAlternative

	// FIXME: is this needed?
	for _, expr := range expr.Exprs {
		analyzeExpression(a, expr, nil)
	}
	for _, alt := range alts {
		if semtypes.ListAlternativeAllowsLength(alt, len(expr.Exprs)) {
			if alt.Pos != nil {
				isValid := true
				lat := alt.Pos
				for i, expr := range expr.Exprs {
					exprTy := expr.GetDeterminedType()
					ty := lat.MemberAtInnerVal(i)
					if !semtypes.IsSubtype(tc, exprTy, ty) {
						isValid = false
						break
					}
				}
				if isValid {
					validAlts = append(validAlts, alt)
				}
			} else {
				validAlts = append(validAlts, alt)
			}
		}
	}

	// Validate uniqueness
	if len(validAlts) == 0 {
		a.semanticErr("no applicable inherent type for list constructor", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}
	}
	if len(validAlts) > 1 {
		a.semanticErr("ambiguous inherent type for list constructor", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}
	}

	// Extract atomic type from selected alternative
	selectedSemType := validAlts[0].SemType
	lat = semtypes.ToListAtomicType(tc, selectedSemType)
	if lat == nil {
		a.semanticErr("applicable type for list constructor is not atomic", expr.GetPosition())
		return nil, semtypes.ListAtomicType{}
	}

	return selectedSemType, *lat
}

func analyzeMappingConstructorExpr[A analyzer](a A, expr *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) bool {
	if expectedType != nil {
		resultType, mappingAtomicType := selectMappingInherentType(a, expr, expectedType)
		if resultType == nil {
			return false
		}
		for _, f := range expr.Fields {
			kv := f.(*ast.BLangMappingKeyValueField)
			var keyName string
			switch keyExpr := kv.Key.Expr.(type) {
			case *ast.BLangLiteral:
				keyName = keyExpr.GetOriginalValue()
			case ast.BNodeWithSymbol:
				keyName = a.ctx().SymbolName(keyExpr.Symbol())
			}
			requiredType := mappingAtomicType.FieldInnerVal(keyName)
			if !analyzeExpression(a, kv.ValueExpr, requiredType) {
				return false
			}
		}
		expr.AtomicType = mappingAtomicType
		setExpectedType(expr, resultType)
	} else {
		for _, f := range expr.Fields {
			kv := f.(*ast.BLangMappingKeyValueField)
			if !analyzeExpression(a, kv.ValueExpr, nil) {
				return false
			}
		}
	}
	return validateResolvedType(a, expr, expectedType)
}

func selectMappingInherentType[A analyzer](a A, expr *ast.BLangMappingConstructorExpr, expectedType semtypes.SemType) (semtypes.SemType, semtypes.MappingAtomicType) {
	expectedMappingType := semtypes.Intersect(expectedType, &semtypes.MAPPING)
	tc := a.tyCtx()
	if semtypes.IsEmpty(tc, expectedMappingType) {
		a.semanticErr("mapping type not found in expected type", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}
	}
	mat := semtypes.ToMappingAtomicType(tc, expectedMappingType)
	if mat != nil {
		return expectedMappingType, *mat
	}
	alts := semtypes.MappingAlternatives(tc, expectedType)
	var validAlts []semtypes.MappingAlternative

	fieldNames := make([]string, len(expr.Fields))
	for i, f := range expr.Fields {
		kv := f.(*ast.BLangMappingKeyValueField)
		fieldNames[i] = recordKeyName(kv.Key)
	}
	sort.Strings(fieldNames)

	for _, alt := range alts {
		if semtypes.MappingAlternativeAllowsFields(alt, fieldNames) {
			if alt.Pos != nil {
				isValid := true
				mat := alt.Pos
				for _, f := range expr.Fields {
					kv := f.(*ast.BLangMappingKeyValueField)
					keyName := recordKeyName(kv.Key)
					exprTy := kv.ValueExpr.GetDeterminedType()
					ty := mat.FieldInnerVal(keyName)
					if !semtypes.IsSubtype(tc, exprTy, ty) {
						isValid = false
						break
					}
				}
				if isValid {
					validAlts = append(validAlts, alt)
				}
			} else {
				validAlts = append(validAlts, alt)
			}
		}
	}
	if len(validAlts) == 0 {
		a.semanticErr("no applicable inherent type for mapping constructor", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}
	}
	if len(validAlts) > 1 {
		a.semanticErr("ambiguous inherent type for mapping constructor", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}
	}

	// Extract atomic type from selected alternative
	selectedSemType := validAlts[0].SemType
	mat = semtypes.ToMappingAtomicType(tc, selectedSemType)
	if mat == nil {
		a.semanticErr("applicable type for mapping constructor is not atomic", expr.GetPosition())
		return nil, semtypes.MappingAtomicType{}
	}

	return selectedSemType, *mat
}

func analyzeErrorConstructorExpr[A analyzer](a A, expr *ast.BLangErrorConstructorExpr, expectedType semtypes.SemType) bool {
	argCount := len(expr.PositionalArgs)
	if argCount < 1 || argCount > 2 {
		a.semanticErr("error constructor must have at least 1 and at most 2 positional arguments", expr.GetPosition())
		return false
	}

	msgArg := expr.PositionalArgs[0]
	if !analyzeExpression(a, msgArg, &semtypes.STRING) {
		return false
	}

	if argCount == 2 {
		causeArg := expr.PositionalArgs[1]
		if !analyzeExpression(a, causeArg, semtypes.Union(&semtypes.ERROR, &semtypes.NIL)) {
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

	switch unaryExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		if !isNumericType(exprTy) {
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
		case model.OperatorKind_EQUALS, model.OperatorKind_NOT_EQUAL:
			anyData := semtypes.CreateAnydata(ctx)
			if !semtypes.IsSubtype(ctx, lhsTy, anyData) && !semtypes.IsSubtype(ctx, rhsTy, anyData) {
				a.semanticErr(fmt.Sprintf("expect anydata types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
				return false
			}
		}
	} else if isBitWiseExpr(binaryExpr) {
		if !analyzeBitWiseExpr(a, binaryExpr, lhsTy, rhsTy, expectedType) {
			return false
		}
	} else if isRangeExpr(binaryExpr) {
		if !semtypes.IsSubtypeSimple(lhsTy, semtypes.INT) || !semtypes.IsSubtypeSimple(rhsTy, semtypes.INT) {
			a.semanticErr(fmt.Sprintf("expect int types for %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
			return false
		}
	} else if isShiftExpr(binaryExpr) {
		if !analyzeShiftExpr(a, binaryExpr, lhsTy, rhsTy, expectedType) {
			return false
		}
	}

	// for nil lifting expression we do semantic analysis as part of type resolver
	// Validate the resolved result type against expected type
	return validateResolvedType(a, binaryExpr, expectedType)
}

var bitWiseOpLookOrder = []semtypes.SemType{semtypes.UINT8, semtypes.UINT16, semtypes.UINT32}

func analyzeBitWiseExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, lhsTy, rhsTy semtypes.SemType, expectedType semtypes.SemType) bool {
	ctx := a.tyCtx()
	nilLifted := false
	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		nilLifted = true
		lhsTy = semtypes.Diff(lhsTy, &semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, &semtypes.NIL)
	}
	if !semtypes.IsSubtype(ctx, lhsTy, &semtypes.INT) || !semtypes.IsSubtype(ctx, rhsTy, &semtypes.INT) {
		a.semanticErr("expect integer types for bitwise operators", binaryExpr.GetPosition())
		return false
	}
	var resultTy semtypes.SemType
	switch binaryExpr.GetOperatorKind() {
	case model.OperatorKind_BITWISE_AND:
		resultTy = &semtypes.INT
		for _, ty := range bitWiseOpLookOrder {
			if semtypes.IsSubtype(ctx, lhsTy, ty) || semtypes.IsSubtype(ctx, rhsTy, ty) {
				resultTy = ty
				break
			}
		}
	case model.OperatorKind_BITWISE_OR, model.OperatorKind_BITWISE_XOR:
		resultTy = &semtypes.INT
		for _, ty := range bitWiseOpLookOrder {
			if semtypes.IsSubtype(ctx, lhsTy, ty) && semtypes.IsSubtype(ctx, rhsTy, ty) {
				resultTy = ty
				break
			}
		}
	default:
		a.internalErr(fmt.Sprintf("unsupported bitwise operator: %s", string(binaryExpr.GetOperatorKind())), binaryExpr.GetPosition())
		return false
	}
	if nilLifted {
		resultTy = semtypes.Union(&semtypes.NIL, resultTy)
	}
	setExpectedType(binaryExpr, resultTy)
	return true
}

func analyzeShiftExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, lhsTy, rhsTy semtypes.SemType, expectedType semtypes.SemType) bool {
	ctx := a.tyCtx()
	nilLifted := false
	if semtypes.ContainsBasicType(lhsTy, semtypes.NIL) || semtypes.ContainsBasicType(rhsTy, semtypes.NIL) {
		nilLifted = true
		lhsTy = semtypes.Diff(lhsTy, &semtypes.NIL)
		rhsTy = semtypes.Diff(rhsTy, &semtypes.NIL)
	}
	op := binaryExpr.GetOperatorKind()
	var resultTy semtypes.SemType = &semtypes.INT

	switch op {
	case model.OperatorKind_BITWISE_RIGHT_SHIFT,
		model.OperatorKind_BITWISE_UNSIGNED_RIGHT_SHIFT:

		for _, ty := range bitWiseOpLookOrder {
			if semtypes.IsSubtype(ctx, lhsTy, ty) {
				resultTy = ty
				break
			}
		}
	default:
		resultTy = &semtypes.INT
	}

	if nilLifted {
		resultTy = semtypes.Union(&semtypes.NIL, resultTy)
	}

	setExpectedType(binaryExpr, resultTy)
	return true
}

func analyzeInvocation[A analyzer](a A, invocation *ast.BLangInvocation, expectedType semtypes.SemType) bool {
	// Get the function type from the symbol
	symbol := invocation.Symbol()
	fnTy := a.ctx().SymbolType(symbol)
	if fnTy == nil || !semtypes.IsSubtypeSimple(fnTy, semtypes.FUNCTION) {
		a.semanticErr("function not found: "+invocation.Name.GetValue(), invocation.GetPosition())
		return false
	}

	// Validate each argument expression
	argTys := make([]semtypes.SemType, len(invocation.ArgExprs))
	for i, arg := range invocation.ArgExprs {
		if !analyzeExpression(a, arg, nil) {
			return false
		}
		argTys[i] = arg.GetDeterminedType()
	}

	// Validate argument types against function parameter types
	paramListTy := semtypes.FunctionParamListType(a.tyCtx(), fnTy)
	argLd := semtypes.NewListDefinition()
	argListTy := argLd.DefineListTypeWrapped(a.tyCtx().Env(), argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
	if !semtypes.IsSubtype(a.tyCtx(), argListTy, paramListTy) {
		a.semanticErr("incompatible arguments for function call", invocation.GetPosition())
		return false
	}

	// Validate the resolved return type against expected type
	return validateResolvedType(a, invocation, expectedType)
}

func analyzeSimpleVariableDef[A analyzer](a A, simpleVariableDef *ast.BLangSimpleVariableDef) bool {
	variable := simpleVariableDef.GetVariable().(*ast.BLangSimpleVariable)
	expectedType := variable.GetDeterminedType()
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
		validateForeach(a, n)
		return initializeLoopAnalyzer(a, n)
	case *ast.BLangIf:
		if !analyzeIf(a, n) {
			return nil
		}
		return a
	case *ast.BLangBreak, *ast.BLangContinue:
		return nil
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
		res := analyzeExpression(a, n.Expr, &semtypes.NIL)
		if !res {
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
	if !analyzeExpression(a, expression, expectedType) {
		return false
	}
	return true
}

func analyzeIf[A analyzer](a A, ifStmt *ast.BLangIf) bool {
	return analyzeExpression(a, ifStmt.Expr, &semtypes.BOOLEAN)
}

func analyzeWhile[A analyzer](a A, whileStmt *ast.BLangWhile) bool {
	return analyzeExpression(a, whileStmt.Expr, &semtypes.BOOLEAN)
}

func validateForeach[A analyzer](a A, foreachStmt *ast.BLangForeach) {
	collection := foreachStmt.Collection
	if !analyzeExpression(a, collection, nil) {
		return
	}
	variable := foreachStmt.VariableDef.GetVariable().(*ast.BLangSimpleVariable)
	variableType := a.ctx().SymbolType(variable.Symbol())
	if binExpr, ok := collection.(*ast.BLangBinaryExpr); ok && isRangeExpr(binExpr) {
		if !semtypes.IsSubtypeSimple(variableType, semtypes.INT) {
			a.semanticErr("foreach variable must be a subtype of int for range expression", collection.GetPosition())
		}
	} else {
		collectionType := collection.GetDeterminedType()
		var expectedValueType semtypes.SemType
		switch {
		case semtypes.IsSubtypeSimple(collectionType, semtypes.LIST):
			memberTypes := semtypes.ListAllMemberTypesInner(a.tyCtx(), collectionType)
			var result semtypes.SemType = &semtypes.NEVER
			for _, each := range memberTypes.SemTypes {
				result = semtypes.Union(result, each)
			}
			expectedValueType = result
		default:
			a.unimplementedErr("unsupported foreach collection", collection.GetPosition())
		}
		if !semtypes.IsSubtype(a.tyCtx(), expectedValueType, variableType) {
			a.ctx().SemanticError("invalid type for variable", variable.GetPosition())
		}
	}
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

func analyzeLiteral[A analyzer](a A, expr *ast.BLangLiteral, expectedType semtypes.SemType) bool {
	if expectedType != nil && semtypes.IsSubtypeSimple(expectedType, semtypes.FLOAT) {
		if intVal, ok := expr.GetValue().(int64); ok {
			floatVal := float64(intVal)
			setExpectedType(expr, semtypes.FloatConst(floatVal))
			expr.SetValue(floatVal)
			expr.GetValueType().BTypeSetTag(model.TypeTags_FLOAT)
		}
	}
	return validateResolvedType(a, expr, expectedType)
}

func setExpectedType[E ast.BLangNode](e E, expectedType semtypes.SemType) {
	e.SetDeterminedType(expectedType)
}
