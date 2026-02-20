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
	"flag"
	"testing"

	"ballerina-lang-go/ast"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
)

func TestTypeResolver(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testTypeResolution(t, testPair)
		})
	}
}

func testTypeResolution(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Type resolution panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv())
	cx := context.NewCompilerContext(env)
	syntaxTree, err := parser.GetSyntaxTree(cx, &debugCtx, testCase.InputPath)
	if err != nil {
		t.Errorf("error getting syntax tree for %s: %v", testCase.InputPath, err)
		return
	}
	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", testCase.InputPath)
		return
	}
	pkg := ast.ToPackage(compilationUnit)
	importedSymbols := ResolveImports(cx, pkg, GetImplicitImports(cx))
	ResolveSymbols(cx, pkg, importedSymbols)
	typeResolver := NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, pkg)
	tyCtx := semtypes.ContextFrom(cx.GetTypeEnv())
	validator := &typeResolutionValidator{t: t, ctx: cx, tyCtx: tyCtx}
	ast.Walk(validator, pkg)

	// If we reach here, type resolution completed without panicking
	t.Logf("Type resolution completed successfully for %s", testCase.InputPath)
}

type typeResolutionValidator struct {
	t     *testing.T
	ctx   *context.CompilerContext
	tyCtx semtypes.Context
}

func (v *typeResolutionValidator) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}

	// Validate that all BLangExpression nodes have their determined type set
	if expr, ok := node.(ast.BLangExpression); ok {
		determinedType := expr.GetDeterminedType()
		if determinedType == nil {
			v.t.Errorf("expression %T at %v does not have determined type set", expr, expr.GetPosition())
		}
		if semtypes.IsNever(determinedType) {
			v.t.Errorf("expression %T at %v has determined type NEVER", expr, expr.GetPosition())
		}

	}

	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		// Skip constant symbols (kind: 1) since they're resolved during semantic analysis
		if v.ctx.SymbolKind(symbol) == model.SymbolKindConstant {
			return v
		}
		if v.ctx.SymbolType(symbol) == nil {
			// FIXME: get rid of this
			if _, ok := node.(*ast.BLangConstant); ok {
				// constants will get their type set during semantic analysis
				return v
			}
			v.t.Errorf("symbol %s (kind: %v) does not have type set for node %T",
				v.ctx.SymbolName(symbol), v.ctx.SymbolKind(symbol), node)
		}
	}

	return v
}

func (v *typeResolutionValidator) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return nil
	}
	if typeData.Type == nil {
		v.t.Errorf("type not resolved for %+v", typeData)
	}
	return v
}
