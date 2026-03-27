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
	"strings"
	"testing"
)

type unsupportedExpr struct {
	ast.BLangLiteral
}

type unsupportedBType struct {
	ast.BLangValueType
}

var queryTestPos = diagnostics.NewBLangDiagnosticLocation("query_test.bal", 0, 0, 0, 0, 0, 0)

func TestResolveQueryExprErrorCases(t *testing.T) {
	testCases := []struct {
		name    string
		query   *ast.BLangQueryExpr
		diagSub string
	}{
		{
			name: "missing select clause list",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), nil, true),
			),
			diagSub: "query expression requires from and select clauses",
		},
		{
			name: "must start with from clause",
			query: newQueryExpr(
				newSelectClause(newIntLiteral(1)),
				newSelectClause(newIntLiteral(2)),
			),
			diagSub: "query expression must start with a from clause",
		},
		{
			name: "requires select clause",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), nil, true),
				newWhereClause(newIntLiteral(1)),
			),
			diagSub: "query expression requires a select clause",
		},
		{
			name: "from collection resolution fails",
			query: newQueryExpr(
				newFromClause(newUnsupportedExprNode(), nil, true),
				newSelectClause(newIntLiteral(1)),
			),
			diagSub: "unsupported expression type",
		},
		{
			name: "from collection is not a list",
			query: newQueryExpr(
				newFromClause(newIntLiteral(42), nil, true),
				newSelectClause(newIntLiteral(1)),
			),
			diagSub: "query from-clause currently supports only list or map collections",
		},
		{
			name: "from binding variable is nil",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), newEmptySimpleVarDef(), true),
				newSelectClause(newIntLiteral(1)),
			),
			diagSub: "only simple variable bindings are supported in from clause",
		},
		{
			name: "from binding type resolution fails",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), newSimpleVarDef("x", newUnsupportedTypeNode(), nil), false),
				newSelectClause(newIntLiteral(1)),
			),
			diagSub: "unsupported type",
		},
		{
			name: "from binding type incompatible",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), newSimpleVarDef("x", newValueType(model.TypeKind_STRING), nil), false),
				newSelectClause(newIntLiteral(1)),
			),
			diagSub: "from-clause variable type is incompatible with collection member type",
		},
		{
			name: "select expression resolution fails",
			query: newQueryExpr(
				newFromClause(newIntListLiteral(1), nil, true),
				newSelectClause(newUnsupportedExprNode()),
			),
			diagSub: "unsupported expression type",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resolver, cx := newTestQueryResolver()
			_, _, ok := resolver.resolveQueryExpr(nil, testCase.query)
			if ok {
				t.Fatalf("expected resolveQueryExpr to fail")
			}
			assertDiagnosticContains(t, cx, testCase.diagSub)
		})
	}
}

func TestResolveQueryIntermediateClauseErrorCases(t *testing.T) {
	testCases := []struct {
		name    string
		clause  ast.BLangNode
		diagSub string
	}{
		{
			name: "let var declaration is nil",
			clause: newLetClause(
				newEmptySimpleVarDef(),
			),
			diagSub: "only simple variable declarations are supported in let clause",
		},
		{
			name: "let var declaration has no initializer",
			clause: newLetClause(
				newSimpleVarDef("y", nil, nil),
			),
			diagSub: "let-clause variable declaration requires an initializer",
		},
		{
			name: "let initializer resolution fails",
			clause: newLetClause(
				newSimpleVarDef("y", nil, newUnsupportedExprNode()),
			),
			diagSub: "unsupported expression type",
		},
		{
			name: "let declared type resolution fails",
			clause: newLetClause(
				newSimpleVarDef("y", newUnsupportedTypeNode(), newIntLiteral(1)),
			),
			diagSub: "unsupported type",
		},
		{
			name: "let declared type incompatible",
			clause: newLetClause(
				newSimpleVarDef("y", newValueType(model.TypeKind_STRING), newIntLiteral(1)),
			),
			diagSub: "let-clause variable type is incompatible with initializer expression",
		},
		{
			name:    "where expression resolution fails",
			clause:  newWhereClause(newUnsupportedExprNode()),
			diagSub: "unsupported expression type",
		},
		{
			name:    "where expression non-boolean",
			clause:  newWhereClause(newIntLiteral(1)),
			diagSub: "where-clause expression must be boolean",
		},
		{
			name:    "unsupported intermediate clause",
			clause:  newCollectClause(),
			diagSub: "only let + where clauses are supported as intermediate query clauses",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			query := newQueryExpr(
				newFromClause(newIntListLiteral(1), nil, true),
				testCase.clause,
				newSelectClause(newIntLiteral(1)),
			)
			resolver, cx := newTestQueryResolver()
			_, ok := resolver.resolveQueryIntermediateClauses(nil, query)
			if ok {
				t.Fatalf("expected resolveQueryIntermediateClauses to fail")
			}
			assertDiagnosticContains(t, cx, testCase.diagSub)
		})
	}
}

func TestResolveQueryExprMapCollection(t *testing.T) {
	resolver, cx := newTestQueryResolver()

	space := cx.NewSymbolSpace(*cx.GetDefaultPackage())
	mapSymbol := model.NewValueSymbol("m", false, false, false)
	space.AddSymbol("m", &mapSymbol)
	mapSymbolRef, _ := space.GetSymbol("m")
	cx.SetSymbolType(mapSymbolRef, semtypes.MAPPING)

	mapRef := &ast.BLangSimpleVarRef{
		VariableName: &ast.BLangIdentifier{
			Value: "m",
		},
	}
	mapRef.SetPosition(queryTestPos)
	mapRef.SetSymbol(mapSymbolRef)

	query := newQueryExpr(
		newFromClause(mapRef, nil, true),
		newSelectClause(newIntLiteral(1)),
	)
	queryTy, _, ok := resolver.resolveQueryExpr(nil, query)
	if !ok {
		t.Fatalf("expected resolveQueryExpr to succeed for map collection")
	}
	if !semtypes.IsSubtypeSimple(queryTy, semtypes.LIST) {
		t.Fatalf("expected query result type to be a list, got %v", queryTy)
	}
}

func TestResolveQueryExprMapConstructType(t *testing.T) {
	resolver, cx := newTestQueryResolver()

	query := newQueryExpr(
		newFromClause(newIntListLiteral(1), nil, true),
		newSelectClause(newListLiteral(newStringLiteral("k"), newIntLiteral(1))),
	)
	query.QueryConstructType = model.TypeKind_MAP

	queryTy, _, ok := resolver.resolveQueryExpr(nil, query)
	if !ok {
		t.Fatalf("expected resolveQueryExpr to succeed for map construct type")
	}
	if !semtypes.IsSubtypeSimple(queryTy, semtypes.MAPPING) {
		t.Fatalf("expected query result type to be mapping, got %v", queryTy)
	}
	if len(cx.Diagnostics()) > 0 {
		t.Fatalf("expected no diagnostics, got %v", cx.Diagnostics())
	}
}

func TestResolveQueryExprMapConstructTypeInvalidSelect(t *testing.T) {
	resolver, cx := newTestQueryResolver()

	query := newQueryExpr(
		newFromClause(newIntListLiteral(1), nil, true),
		newSelectClause(newIntLiteral(1)),
	)
	query.QueryConstructType = model.TypeKind_MAP

	_, _, ok := resolver.resolveQueryExpr(nil, query)
	if ok {
		t.Fatalf("expected resolveQueryExpr to fail for invalid map select expression")
	}
	assertDiagnosticContains(t, cx, "incompatible type")
}

func newTestQueryResolver() (*TypeResolver, *context.CompilerContext) {
	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv())
	cx := context.NewCompilerContext(env)
	return newTypeResolver(cx, &ast.BLangPackage{}, nil), cx
}

func assertDiagnosticContains(t *testing.T, cx *context.CompilerContext, substr string) {
	t.Helper()
	diagnosticsList := cx.Diagnostics()
	if len(diagnosticsList) == 0 {
		t.Fatalf("expected at least one diagnostic containing %q, but diagnostics are empty", substr)
	}
	for _, diag := range diagnosticsList {
		if strings.Contains(diag.Message(), substr) {
			return
		}
	}

	messages := make([]string, len(diagnosticsList))
	for i, diag := range diagnosticsList {
		messages[i] = diag.Message()
	}
	t.Fatalf("expected diagnostic containing %q, got: %v", substr, messages)
}

func newQueryExpr(clauses ...ast.BLangNode) *ast.BLangQueryExpr {
	query := &ast.BLangQueryExpr{
		QueryClauseList: clauses,
	}
	query.SetPosition(queryTestPos)
	return query
}

func newFromClause(collection ast.BLangExpression, varDef model.VariableDefinitionNode, declaredWithVar bool) *ast.BLangFromClause {
	fromClause := &ast.BLangFromClause{
		BLangInputClause: ast.BLangInputClause{
			VariableDefinitionNode: varDef,
			IsDeclaredWithVarFlag:  declaredWithVar,
		},
	}
	fromClause.SetPosition(queryTestPos)
	fromClause.SetCollection(collection)
	return fromClause
}

func newSelectClause(expr ast.BLangExpression) *ast.BLangSelectClause {
	selectClause := &ast.BLangSelectClause{}
	selectClause.SetPosition(queryTestPos)
	selectClause.SetExpression(expr)
	return selectClause
}

func newLetClause(defs ...model.VariableDefinitionNode) *ast.BLangLetClause {
	letClause := &ast.BLangLetClause{
		LetVarDeclarations: defs,
	}
	letClause.SetPosition(queryTestPos)
	return letClause
}

func newWhereClause(expr ast.BLangExpression) *ast.BLangWhereClause {
	whereClause := &ast.BLangWhereClause{
		Expression: expr,
	}
	whereClause.SetPosition(queryTestPos)
	return whereClause
}

func newCollectClause() *ast.BLangCollectClause {
	collectClause := &ast.BLangCollectClause{}
	collectClause.SetPosition(queryTestPos)
	collectClause.SetExpression(newIntLiteral(1))
	return collectClause
}

func newEmptySimpleVarDef() *ast.BLangSimpleVariableDef {
	varDef := &ast.BLangSimpleVariableDef{}
	varDef.SetPosition(queryTestPos)
	return varDef
}

func newSimpleVarDef(name string, typeNode ast.BType, expr ast.BLangExpression) *ast.BLangSimpleVariableDef {
	ident := &ast.BLangIdentifier{
		Value:         name,
		OriginalValue: name,
	}
	ident.SetPosition(queryTestPos)

	variable := &ast.BLangSimpleVariable{
		Name: ident,
	}
	variable.SetPosition(queryTestPos)
	variable.SetTypeNode(typeNode)
	if expr != nil {
		variable.SetExpr(expr)
	}

	varDef := &ast.BLangSimpleVariableDef{
		Var: variable,
	}
	varDef.SetPosition(queryTestPos)
	return varDef
}

func newValueType(typeKind model.TypeKind) ast.BType {
	ty := &ast.BLangValueType{
		TypeKind: typeKind,
	}
	ty.SetPosition(queryTestPos)
	return ty
}

func newUnsupportedTypeNode() ast.BType {
	ty := &unsupportedBType{}
	ty.SetPosition(queryTestPos)
	return ty
}

func newIntLiteral(value int64) *ast.BLangLiteral {
	literal := &ast.BLangLiteral{}
	literal.SetPosition(queryTestPos)
	literal.SetValue(value)
	literal.SetValueType(ast.NewBType(model.TypeTags_INT, "", 0))
	return literal
}

func newIntListLiteral(values ...int64) *ast.BLangListConstructorExpr {
	exprs := make([]ast.BLangExpression, 0, len(values))
	for _, value := range values {
		exprs = append(exprs, newIntLiteral(value))
	}
	listExpr := &ast.BLangListConstructorExpr{
		Exprs: exprs,
	}
	listExpr.SetPosition(queryTestPos)
	return listExpr
}

func newStringLiteral(value string) *ast.BLangLiteral {
	literal := &ast.BLangLiteral{}
	literal.SetPosition(queryTestPos)
	literal.SetValue(value)
	literal.SetValueType(ast.NewBType(model.TypeTags_STRING, "", 0))
	return literal
}

func newListLiteral(values ...ast.BLangExpression) *ast.BLangListConstructorExpr {
	listExpr := &ast.BLangListConstructorExpr{
		Exprs: values,
	}
	listExpr.SetPosition(queryTestPos)
	return listExpr
}

func newUnsupportedExprNode() ast.BLangExpression {
	expr := &unsupportedExpr{}
	expr.SetPosition(queryTestPos)
	return expr
}
