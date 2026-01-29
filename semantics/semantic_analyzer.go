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
	"fmt"
	"math/bits"
	"reflect"
)

type analyzer interface {
	ast.Visitor
	ctx() *context.CompilerContext
	tyCtx() semtypes.Context
	refTy(name UniformRef) semtypes.SemType
	unimplementedErr(message string)
	semanticErr(message string)
	syntaxErr(message string)
	internalErr(message string)
	localRef(name string) UniformRef
	setRefTy(name UniformRef, ty semtypes.SemType)
}

type (
	SemanticAnalyzer struct {
		compilerCtx   *context.CompilerContext
		typeCtx       semtypes.Context
		resolvedTypes TypeResolutionResult
		// TODO: move the constant resolution to type resolver as well so that we can run semantic analyzer in parallel as well
		constants map[UniformRef]*ast.BLangConstant
		pkg       *ast.BLangPackage
	}
	constantAnalyzer struct {
		sa           *SemanticAnalyzer
		constant     *ast.BLangConstant
		expectedType semtypes.SemType
	}

	functionAnalyzer struct {
		sa          *SemanticAnalyzer
		function    *ast.BLangFunction
		localVarsTy map[UniformRef]semtypes.SemType
		retTy       semtypes.SemType
	}
)

var _ analyzer = &SemanticAnalyzer{}
var _ analyzer = &constantAnalyzer{}
var _ analyzer = &functionAnalyzer{}

func (sa *SemanticAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return nil
}

func (fa *functionAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return nil
}

func (fa *functionAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangSimpleVariableDef:
		analyzeSimpleVariableDef(fa, n)
		return fa
	case *ast.BLangAssignment:
		analyzeAssignment(fa, n)
		return fa
	case *ast.BLangReturn:
		if n.Expr == nil {
			if !semtypes.IsSubtypeSimple(fa.retTy, semtypes.NIL) {
				// TODO: should put error at this node
				fa.semanticErr("expect a return value")
			}
			return nil
		}
		analyzeExpression(fa, n.Expr, fa.retTy)
	case *ast.BLangExpressionStmt:
		analyzeExpression(fa, n.Expr, nil)
		return fa
	default:
		return fa
	}
	return fa
}

func (ca *constantAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return ca
}

func (sa *SemanticAnalyzer) localRef(name string) UniformRef {
	return refInPackage(sa.pkg, name)
}

func (ca *constantAnalyzer) localRef(name string) UniformRef {
	return ca.sa.localRef(name)
}

func (sa *SemanticAnalyzer) ctx() *context.CompilerContext {
	return sa.compilerCtx
}

func (sa *SemanticAnalyzer) tyCtx() semtypes.Context {
	return sa.typeCtx
}

func (sa *SemanticAnalyzer) refTy(name UniformRef) semtypes.SemType {
	ty, ok := sa.resolvedTypes.functions[name]
	if ok {
		return ty
	}
	constant, ok := sa.constants[name]
	if ok {
		typeData := constant.GetTypeData()
		return typeData.Type
	}
	sa.semanticErr(fmt.Sprintf("symbol %s not found", name))
	return nil
}

func (fa *functionAnalyzer) refTy(name UniformRef) semtypes.SemType {
	ty, ok := fa.localVarsTy[name]
	if ok {
		return ty
	}
	return fa.sa.refTy(name)
}

func (fa *functionAnalyzer) localRef(name string) UniformRef {
	return fa.sa.localRef(name)
}

func (fa *functionAnalyzer) ctx() *context.CompilerContext {
	return fa.sa.ctx()
}

func (fa *functionAnalyzer) tyCtx() semtypes.Context {
	return fa.sa.tyCtx()
}

func (ca *constantAnalyzer) ctx() *context.CompilerContext {
	return ca.sa.ctx()
}

func (ca *constantAnalyzer) tyCtx() semtypes.Context {
	return ca.sa.tyCtx()
}

func (ca *constantAnalyzer) refTy(name UniformRef) semtypes.SemType {
	return ca.sa.refTy(name)
}

func (sa *SemanticAnalyzer) setRefTy(name UniformRef, ty semtypes.SemType) {
	panic("setRefTy not supported on SemanticAnalyzer")
}

func (ca *constantAnalyzer) setRefTy(name UniformRef, ty semtypes.SemType) {
	panic("setRefTy not supported on constantAnalyzer")
}

func (fa *functionAnalyzer) setRefTy(name UniformRef, ty semtypes.SemType) {
	fa.localVarsTy[name] = ty
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
	ca.sa.compilerCtx.Unimplemented(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) semanticErr(message string) {
	ca.sa.compilerCtx.SemanticError(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) syntaxErr(message string) {
	ca.sa.compilerCtx.SyntaxError(message, ca.constant.GetPosition())
}

func (ca *constantAnalyzer) internalErr(message string) {
	ca.sa.compilerCtx.InternalError(message, ca.constant.GetPosition())
}

func (fa *functionAnalyzer) unimplementedErr(message string) {
	fa.sa.compilerCtx.Unimplemented(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) semanticErr(message string) {
	fa.sa.compilerCtx.SemanticError(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) syntaxErr(message string) {
	fa.sa.compilerCtx.SyntaxError(message, fa.function.GetPosition())
}

func (fa *functionAnalyzer) internalErr(message string) {
	fa.sa.compilerCtx.InternalError(message, fa.function.GetPosition())
}

// When we support multiple packages we need to resolve types of all of them before semantic analysis
func NewSemanticAnalyzer(ctx *context.CompilerContext, resolvedTypes TypeResolutionResult) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		compilerCtx:   ctx,
		typeCtx:       semtypes.ContextFrom(semtypes.GetTypeEnv()),
		resolvedTypes: resolvedTypes,
		constants:     make(map[UniformRef]*ast.BLangConstant),
	}
}

func (sa *SemanticAnalyzer) Analyze(pkg *ast.BLangPackage) {
	sa.pkg = pkg
	ast.Walk(sa, pkg)
	sa.pkg = nil
}

func (sa *SemanticAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangConstant:
		return &constantAnalyzer{sa: sa, constant: n}
	case ast.BLangExpression:
		analyzeExpression(sa, n, nil)
		return nil
	case *ast.BLangFunction:
		return initializeFunctionAnalyzer(sa, n)
	default:
		return sa
	}
}

func initializeFunctionAnalyzer(sa *SemanticAnalyzer, function *ast.BLangFunction) *functionAnalyzer {
	fa := &functionAnalyzer{sa: sa, function: function, localVarsTy: make(map[UniformRef]semtypes.SemType)}
	for _, param := range function.RequiredParams {
		name := param.GetName().GetValue()
		fa.setRefTy(fa.localRef(name), param.GetTypeData().Type)
	}
	fa.retTy = function.ReturnTypeData.Type
	return fa
}

func (ca *constantAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		setExpectedType(ca.constant, ca.expectedType)
		ca.constant.TypeData.Type = ca.expectedType
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangTypeDefinition:
		typeData := n.GetTypeData()
		expectedType := typeData.Type
		if expectedType == nil {
			ca.syntaxErr("type not resolved")
			return nil
		}
		ctx := ca.sa.tyCtx()
		if semtypes.IsNever(expectedType) || !semtypes.IsSubtype(ctx, expectedType, semtypes.CreateAnydata(ctx)) {
			ca.syntaxErr("invalid type for constant declaration")
			return nil
		}
		ca.expectedType = expectedType
	case *ast.BLangIdentifier:
		name := n.GetValue()
		ref := ca.localRef(name)
		if _, ok := ca.sa.constants[ref]; ok {
			ca.syntaxErr("constant already declared")
			return nil
		}
		ca.sa.constants[ref] = ca.constant
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
			exprTy := bLangExpr.GetBType().Type
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

func analyzeExpression[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) {
	switch expr := expr.(type) {
	// Literals
	case *ast.BLangLiteral:
		if expectedType == nil {
			return
		}
		typeData := expr.GetBType()
		ty := typeData.Type
		ctx := a.tyCtx()
		if !semtypes.IsSubtype(ctx, ty, expectedType) {
			a.semanticErr("incompatible type for literal")
			return
		}
	case *ast.BLangNumericLiteral:
		if expectedType == nil {
			return
		}
		typeData := expr.GetBType()
		ty := typeData.Type
		ctx := a.tyCtx()
		if !semtypes.IsSubtype(ctx, ty, expectedType) {
			a.semanticErr("incompatible type for literal")
			return
		}

	// Variable References
	case *ast.BLangSimpleVarRef:
		ref := a.localRef(expr.VariableName.GetValue())
		ty := a.refTy(ref)
		if ty == nil {
			a.semanticErr("variable not found: " + expr.VariableName.GetValue())
			return
		}
		if expectedType != nil {
			if !semtypes.IsSubtype(a.tyCtx(), ty, expectedType) {
				a.semanticErr("incompatible type for variable reference")
				return
			}
		}
		setExpectedType(expr, ty)
	case *ast.BLangLocalVarRef, *ast.BLangConstRef:
		panic("not implemented")
	// Operators
	case *ast.BLangBinaryExpr:
		analyzeBinaryExpr(a, expr, expectedType)
	case *ast.BLangUnaryExpr:
		analyzeUnaryExpr(a, expr, expectedType)

	// Function and Method Calls
	case *ast.BLangInvocation:
		analyzeInvocation(a, expr, expectedType)

	// Indexing
	case *ast.BLangIndexBasedAccess:
		analyzeIndexBasedAccess(a, expr, expectedType)

	// Collections and Groups
	case *ast.BLangListConstructorExpr:
		analyzeListConstructorExpr(a, expr, expectedType)
	case *ast.BLangGroupExpr:
		panic("not implemented")
	case *ast.BLangWildCardBindingPattern:
		setExpectedType(expr, &semtypes.ANY)
	default:
		a.internalErr("unexpected expression type: " + reflect.TypeOf(expr).String())
	}
}

func analyzeIndexBasedAccess[A analyzer](a A, expr *ast.BLangIndexBasedAccess, expectedType semtypes.SemType) {
	containerExpr := expr.Expr
	analyzeExpression(a, containerExpr, nil)
	containerExprTy := containerExpr.GetBType().Type
	var keyExprExpectedType semtypes.SemType
	ctx := a.tyCtx()
	if !semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) || !semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) || !semtypes.IsSubtypeSimple(containerExprTy, semtypes.XML) {
		keyExprExpectedType = &semtypes.INT
	} else if !semtypes.IsSubtypeSimple(containerExprTy, semtypes.TABLE) {
		a.unimplementedErr("table not supported")
	} else if !semtypes.IsSubtype(ctx, containerExprTy, semtypes.Union(&semtypes.NIL, &semtypes.MAPPING)) {
		keyExprExpectedType = &semtypes.STRING
	} else {
		a.semanticErr("incompatible type for index based access")
		return
	}
	keyExpr := expr.IndexExpr
	analyzeExpression(a, keyExpr, keyExprExpectedType)
	keyExprTy := keyExpr.GetBType().Type
	var resultTy semtypes.SemType
	if semtypes.IsSubtypeSimple(containerExprTy, semtypes.LIST) {
		resultTy = semtypes.ListProjInnerVal(ctx, containerExprTy, keyExprTy)
	} else if semtypes.IsSubtypeSimple(containerExprTy, semtypes.STRING) {
		resultTy = &semtypes.STRING
	} else {
		a.unimplementedErr("unsupported container type for index based access")
		return
	}
	if expectedType != nil {
		if !semtypes.IsSubtype(ctx, resultTy, expectedType) {
			a.semanticErr("incompatible type for index based access")
			return
		}
	}
	setExpectedType(expr, resultTy)
}


func analyzeListConstructorExpr[A analyzer](a A, expr *ast.BLangListConstructorExpr, expectedType semtypes.SemType) {
	memberTypes := make([]semtypes.SemType, len(expr.Exprs))
	for i, expr := range expr.Exprs {
		analyzeExpression(a, expr, nil)
		memberTypes[i] = expr.GetBType().Type
	}
	ld := semtypes.NewListDefinition()
	valueTy := ld.DefineListTypeWrapped(a.tyCtx().Env(), memberTypes, len(memberTypes), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
	if expectedType != nil {
		if !semtypes.IsSubtype(a.tyCtx(), valueTy, expectedType) {
			a.semanticErr("incompatible type for list constructor expression")
			return
		}
	}
	setExpectedType(expr, valueTy)
}

func analyzeUnaryExpr[A analyzer](a A, unaryExpr *ast.BLangUnaryExpr, expectedType semtypes.SemType) {
	analyzeExpression(a, unaryExpr.Expr, expectedType)
	exprTy := unaryExpr.Expr.GetBType().Type
	var resultTy semtypes.SemType
	switch unaryExpr.GetOperatorKind() {
	case model.OperatorKind_ADD, model.OperatorKind_SUB, model.OperatorKind_BITWISE_COMPLEMENT:
		if !isNumericType(exprTy) {
			a.semanticErr(fmt.Sprintf("expect numeric type for %s", string(unaryExpr.GetOperatorKind())))
			return
		}
		resultTy = exprTy
	case model.OperatorKind_NOT:
		if !semtypes.IsSubtypeSimple(exprTy, semtypes.BOOLEAN) {
			a.semanticErr(fmt.Sprintf("expect boolean type for %s", string(unaryExpr.GetOperatorKind())))
			return
		}
		resultTy = exprTy
	default:
		a.semanticErr(fmt.Sprintf("unsupported unary operator: %s", string(unaryExpr.GetOperatorKind())))
		return
	}
	if expectedType != nil {
		if !semtypes.IsSubtype(a.tyCtx(), resultTy, expectedType) {
			a.semanticErr("incompatible result type for unary expression")
			return
		}
	}
	setExpectedType(unaryExpr, resultTy)
}

func analyzeBinaryExpr[A analyzer](a A, binaryExpr *ast.BLangBinaryExpr, expectedType semtypes.SemType) {
	analyzeExpression(a, binaryExpr.LhsExpr, nil)
	analyzeExpression(a, binaryExpr.RhsExpr, nil)
	lhsTy := binaryExpr.LhsExpr.GetBType().Type
	rhsTy := binaryExpr.RhsExpr.GetBType().Type
	var resultTy semtypes.BasicTypeBitSet
	if isEqualityExpr(binaryExpr) {
		intersection := semtypes.Intersect(lhsTy, rhsTy)
		if semtypes.IsEmpty(a.tyCtx(), intersection) {
			a.semanticErr(fmt.Sprintf("expect same type for %s", string(binaryExpr.GetOperatorKind())))
			return
		}
		resultTy = semtypes.BOOLEAN
	} else {
		lhsBasicTy := semtypes.WidenToBasicTypes(lhsTy)
		rhsBasicTy := semtypes.WidenToBasicTypes(rhsTy)
		numLhsBits := bits.OnesCount(uint(lhsBasicTy.All()))
		numRhsBits := bits.OnesCount(uint(rhsBasicTy.All()))
		if numLhsBits != 1 || numRhsBits != 1 {
			a.semanticErr(fmt.Sprintf("union types not supported for %s", string(binaryExpr.GetOperatorKind())))
			return
		}
		if semtypes.IsSubtypeSimple(&lhsBasicTy, semtypes.NIL) || semtypes.IsSubtypeSimple(&rhsBasicTy, semtypes.NIL) {
			a.unimplementedErr("nil lifting not supported")
			return
		}
		if isMultipcativeExpr(binaryExpr) {
			if !isNumericType(&lhsBasicTy) || !isNumericType(&rhsBasicTy) {
				a.semanticErr(fmt.Sprintf("expect numeric types for %s", string(binaryExpr.GetOperatorKind())))
				return
			}
			if lhsBasicTy == rhsBasicTy {
				resultTy = lhsBasicTy
			} else {
				a.unimplementedErr("type coercion not supported")
			}
		} else if isAdditiveExpr(binaryExpr) {
			supportedTypes := semtypes.Union(&semtypes.NUMBER, &semtypes.STRING)
			ctx := a.tyCtx()
			if !semtypes.IsSubtype(ctx, &lhsBasicTy, supportedTypes) || !semtypes.IsSubtype(ctx, &rhsBasicTy, supportedTypes) {
				a.semanticErr(fmt.Sprintf("expect numeric or string types for %s", string(binaryExpr.GetOperatorKind())))
				return
			}
			if lhsBasicTy == rhsBasicTy {
				resultTy = lhsBasicTy
			} else {
				a.unimplementedErr("type coercion not supported")
			}
		} else if isRelationalExpr(binaryExpr) {
			if !semtypes.Comparable(a.tyCtx(), &lhsBasicTy, &rhsBasicTy) {
				a.semanticErr(fmt.Sprintf("expect comparable types for %s", string(binaryExpr.GetOperatorKind())))
				return
			}
			resultTy = semtypes.BOOLEAN
		} else {
			a.unimplementedErr(fmt.Sprintf("unsupported operator: %s", string(binaryExpr.GetOperatorKind())))
			return
		}
	}
	if expectedType != nil {
		if !semtypes.IsSubtype(a.tyCtx(), &resultTy, expectedType) {
			a.semanticErr("incompatible result type for binary expression")
			return
		}
	}
	setExpectedType(binaryExpr, &resultTy)
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

func isMultipcativeExpr(opExpr opExpr) bool {
	switch opExpr.GetOperatorKind() {
	case model.OperatorKind_MUL, model.OperatorKind_DIV, model.OperatorKind_MOD:
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

func isNumericType(ty semtypes.SemType) bool {
	return semtypes.IsSubtypeSimple(ty, semtypes.NUMBER)
}

func analyzeInvocation[A analyzer](a A, invocation *ast.BLangInvocation, expectedType semtypes.SemType) {
	var retTy semtypes.SemType
	// TODO: fix this when we properly support libraries
	if invocation.PkgAlias != nil && invocation.PkgAlias.GetValue() != "" {
		if invocation.PkgAlias.GetValue() != "io" {
			a.unimplementedErr("unsupported package alias: " + invocation.PkgAlias.GetValue())
		} else if invocation.Name.GetValue() == "println" {
			retTy = &semtypes.NIL
		} else {
			a.unimplementedErr("unsupported io function: " + invocation.Name.GetValue())
		}
	} else {
		fnTy := a.refTy(a.localRef(invocation.Name.GetValue()))
		if fnTy == nil || !semtypes.IsSubtypeSimple(fnTy, semtypes.FUNCTION) {
			a.semanticErr("function not found: " + invocation.Name.GetValue())
			return
		}
		argTys := make([]semtypes.SemType, len(invocation.ArgExprs))
		for i, arg := range invocation.ArgExprs {
			analyzeExpression(a, arg, nil)
			typeData := arg.GetBType()
			argTys[i] = typeData.Type
		}
		paramListTy := semtypes.FunctionParamListType(a.tyCtx(), fnTy)
		argLd := semtypes.NewListDefinition()
		argListTy := argLd.DefineListTypeWrapped(a.tyCtx().Env(), argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE)
		if !semtypes.IsSubtype(a.tyCtx(), argListTy, paramListTy) {
			a.semanticErr("incompatible arguments for function call")
			return
		}
		retTy = semtypes.FunctionReturnType(a.tyCtx(), fnTy, argListTy)
	}
	if expectedType != nil {
		if !semtypes.IsSubtype(a.tyCtx(), retTy, expectedType) {
			a.semanticErr("incompatible return type for function call")
			return
		}
	}
	setExpectedType(invocation, retTy)
}

func analyzeSimpleVariableDef[A analyzer](a A, simpleVariableDef *ast.BLangSimpleVariableDef) {
	variable := simpleVariableDef.GetVariable().(*ast.BLangSimpleVariable)
	expectedType := variable.GetTypeData().Type
	if variable.Expr != nil {
		analyzeExpression(a, variable.Expr.(ast.BLangExpression), expectedType)
	}
	name := variable.Name.GetValue()
	ref := a.localRef(name)
	a.setRefTy(ref, expectedType)
	setExpectedType(simpleVariableDef, expectedType)
}

func analyzeAssignment[A analyzer](a A, assignment *ast.BLangAssignment) {
	analyzeExpression(a, assignment.VarRef, nil)
	expectedType := assignment.VarRef.GetBType().Type
	analyzeExpression(a, assignment.Expr, expectedType)
}


func setExpectedType[E ast.BLangNode](e E, expectedType semtypes.SemType) {
	typeData := e.GetBType()
	typeData.Type = expectedType
	e.SetBType(typeData)
	e.SetDeterminedType(expectedType)
}
