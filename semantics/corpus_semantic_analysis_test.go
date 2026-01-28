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
	"ballerina-lang-go/parser"
	"ballerina-lang-go/test_util"
	"flag"
	"testing"
)

func TestSemanticAnalysis(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testSemanticAnalysis(t, testPair)
		})
	}
}

func testSemanticAnalysis(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Semantic analysis panicked for %s: %v", testCase.InputPath, r)
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

	// Step 1: Type Resolution
	typeResolver := NewIsolatedTypeResolver(cx)
	resolvedTypes := typeResolver.ResolveTypes(pkg)

	// Step 2: Semantic Analysis
	semanticAnalyzer := NewSemanticAnalyzer(cx, resolvedTypes)
	semanticAnalyzer.Analyze(pkg)

	// If we reach here, semantic analysis completed without panicking
	t.Logf("Semantic analysis completed successfully for %s", testCase.InputPath)
}
