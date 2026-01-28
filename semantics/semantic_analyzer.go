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
	"ballerina-lang-go/semtypes"
	"fmt"
	"reflect"
)

// TODO: move the constant resolution to type resolver as well so that we can run semantic analyzer in parallel as well
type analyzer interface {
	ast.Visitor
	ctx() *context.CompilerContext
	tyCtx() semtypes.Context
	symbolTy(name string) semtypes.SemType
	unimplementedErr(message string)
	semanticErr(message string)
	syntaxErr(message string)
	internalErr(message string)
}

type (
	SemanticAnalyzer struct {
		compilerCtx   *context.CompilerContext
		typeCtx       semtypes.Context
		resolvedTypes TypeResolutionResult
		constants     map[string]*ast.BLangConstant
	}
	constantAnalyzer struct {
		sa           *SemanticAnalyzer
		constant     *ast.BLangConstant
		expectedType semtypes.SemType
	}
)

var _ analyzer = &SemanticAnalyzer{}
var _ analyzer = &constantAnalyzer{}

func (sa *SemanticAnalyzer) ctx() *context.CompilerContext {
	return sa.compilerCtx
}

func (sa *SemanticAnalyzer) tyCtx() semtypes.Context {
	return sa.typeCtx
}

func (sa *SemanticAnalyzer) symbolTy(name string) semtypes.SemType {
	ty, ok := sa.resolvedTypes.functions[name]
	if ok {
		return ty
	}
	constant, ok := sa.constants[name]
	if ok {
		return constant.GetBType().(ast.BType).SemType()
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

func (ca *constantAnalyzer) symbolTy(name string) semtypes.SemType {
	return ca.sa.symbolTy(name)
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
		constants:     make(map[string]*ast.BLangConstant),
	}
}

func (sa *SemanticAnalyzer) Analyze(pkg *ast.BLangPackage) {
	ast.Walk(sa, pkg)
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
		expectedType := n.GetBType().(ast.BType).SemType()
		ctx := ca.sa.tyCtx()
		if semtypes.IsNever(expectedType) || !semtypes.IsSubtype(ctx, expectedType, semtypes.CreateAnydata(ctx)) {
			ca.syntaxErr("invalid type for constant declaration")
			return nil
		}
		ca.expectedType = expectedType
	case *ast.BLangIdentifier:
		name := n.GetValue()
		if _, ok := ca.sa.constants[name]; ok {
			ca.syntaxErr("constant already declared")
			return nil
		}
		ca.sa.constants[name] = ca.constant
	}
	return ca
}

func analyzeExpression[A analyzer](a A, expr ast.BLangExpression, expectedType semtypes.SemType) {
	switch expr := expr.(type) {
	// Literals
	case *ast.BLangLiteral:
		// TODO: implement semantic analysis
	case *ast.BLangNumericLiteral:
		// TODO: implement semantic analysis

	// Variable References
	case *ast.BLangSimpleVarRef:
		// TODO: implement semantic analysis
	case *ast.BLangLocalVarRef:
		// TODO: implement semantic analysis
	case *ast.BLangConstRef:
		// TODO: implement semantic analysis

	// Operators
	case *ast.BLangBinaryExpr:
		// TODO: implement semantic analysis
	case *ast.BLangUnaryExpr:
		// TODO: implement semantic analysis

	// Function and Method Calls
	case *ast.BLangInvocation:
		// TODO: implement semantic analysis

	// Indexing
	case *ast.BLangIndexBasedAccess:
		// TODO: implement semantic analysis

	// Collections and Groups
	case *ast.BLangListConstructorExpr:
		// TODO: implement semantic analysis
	case *ast.BLangGroupExpr:
		// TODO: implement semantic analysis

	default:
		a.internalErr("unexpected expression type: " + reflect.TypeOf(expr).String())
	}
}
