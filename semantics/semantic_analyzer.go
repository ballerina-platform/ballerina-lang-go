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
}

type (
	SemanticAnalyzer struct {
		compilerCtx   *context.CompilerContext
		typeCtx       semtypes.Context
		resolvedTypes TypeResolutionResult
		// TODO: move the constant resolution to type resolver as well so that we can run semantic analyzer in parallel as well
		constants     map[UniformRef]*ast.BLangConstant
		pkg           *ast.BLangPackage
	}
	constantAnalyzer struct {
		sa           *SemanticAnalyzer
		constant     *ast.BLangConstant
		expectedType semtypes.SemType
	}
)

var _ analyzer = &SemanticAnalyzer{}
var _ analyzer = &constantAnalyzer{}

func (sa *SemanticAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	sa.unimplementedErr("type data not supported")
	return nil
}

func (ca *constantAnalyzer) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	ca.unimplementedErr("type data not supported")
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

func (ca *constantAnalyzer) ctx() *context.CompilerContext {
	return ca.sa.ctx()
}

func (ca *constantAnalyzer) tyCtx() semtypes.Context {
	return ca.sa.tyCtx()
}

func (ca *constantAnalyzer) refTy(name UniformRef) semtypes.SemType {
	return ca.sa.refTy(name)
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
	default:
		return sa
	}
}

func (ca *constantAnalyzer) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
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
            case model.NodeKind_LITERAL:
            case model.NodeKind_NUMERIC_LITERAL:
            case model.NodeKind_STRING_TEMPLATE_LITERAL:
            case model.NodeKind_RECORD_LITERAL_EXPR:
            case model.NodeKind_LIST_CONSTRUCTOR_EXPR:
            case model.NodeKind_LIST_CONSTRUCTOR_SPREAD_OP:
            case model.NodeKind_SIMPLE_VARIABLE_REF:
            case model.NodeKind_BINARY_EXPR:
            case model.NodeKind_GROUP_EXPR:
            case model.NodeKind_UNARY_EXPR:
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
		typeData := expr.GetBType()
		ty := typeData.Type
		ctx := a.tyCtx()
		if !semtypes.IsSubtype(ctx, ty, expectedType) {
			a.semanticErr("incompatible type for literal")
			return
		}
	case *ast.BLangNumericLiteral:
		typeData := expr.GetBType()
		ty := typeData.Type
		ctx := a.tyCtx()
		if !semtypes.IsSubtype(ctx, ty, expectedType) {
			a.semanticErr("incompatible type for literal")
			return
		}

	// Variable References
	case *ast.BLangSimpleVarRef:
	case *ast.BLangLocalVarRef:
	case *ast.BLangConstRef:

	// Operators
	case *ast.BLangBinaryExpr:
	case *ast.BLangUnaryExpr:

	// Function and Method Calls
	case *ast.BLangInvocation:

	// Indexing
	case *ast.BLangIndexBasedAccess:

	// Collections and Groups
	case *ast.BLangListConstructorExpr:
	case *ast.BLangGroupExpr:

	default:
		a.internalErr("unexpected expression type: " + reflect.TypeOf(expr).String())
	}
}

func analyzeInvocation[A analyzer](a A, invocation *ast.BLangInvocation, expectedType semtypes.SemType) {
	var retTy semtypes.SemType
	// TODO: fix this when we properly support libraries
	if invocation.PkgAlias != nil {
		if invocation.PkgAlias.GetValue() != "io" {
			a.unimplementedErr("unsupported package alias: " + invocation.PkgAlias.GetValue())
		} else if invocation.Name.GetValue() != "println" {
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
		argListTy := argLd.DefineListTypeWrapped(a.tyCtx().Env(), argTys, len(argTys), &semtypes.NEVER, semtypes.CellMutability_CELL_MUT_NONE);
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
	typeData := invocation.GetBType()
	typeData.Type = retTy
	invocation.SetBType(typeData)
}
