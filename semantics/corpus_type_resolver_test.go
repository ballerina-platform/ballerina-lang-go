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
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"flag"
	"testing"
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
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, testCase.InputPath)
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
	env := semtypes.GetIsolatedTypeEnv()
	importedSymbols := ResolveImports(env, pkg)
	ResolveSymbols(cx, pkg, importedSymbols)
	typeResolver := NewIsolatedTypeResolver(cx)
	typeResolver.ResolveTypes(pkg)
	validator := &typeResolutionValidator{t: t}
	ast.Walk(validator, pkg)

	// If we reach here, type resolution completed without panicking
	t.Logf("Type resolution completed successfully for %s", testCase.InputPath)
}

type typeResolutionValidator struct {
	t *testing.T
}

func (v *typeResolutionValidator) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}

	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		// Skip constant symbols (kind: 1) since they're resolved during semantic analysis
		if symbol.Kind() == model.SymbolKindConstant {
			return v
		}
		if symbol.Type() == nil {
			if isExpr(node) {
				// expressions will get their type set during semantic analysis
				return v
			} else if _, ok := node.(*ast.BLangConstant); ok {
				// constants will get their type set during semantic analysis
				return v
			}
			v.t.Errorf("symbol %s (kind: %v) does not have type set for node %T",
				symbol.Name(), symbol.Kind(), node)
		}
	}

	return v
}

func isExpr(node ast.BLangNode) bool {
	if _, ok := node.(model.ExpressionNode); ok {
		return true
	}
	return false
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
