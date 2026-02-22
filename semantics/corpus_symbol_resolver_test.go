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

package semantics_test

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/test_util/testphases"
	"flag"
	"testing"
)

func TestSymbolResolver(t *testing.T) {
	flag.Parse()
	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testSymbolResolution(t, testPair)
		})
	}
}

func testSymbolResolution(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Symbol resolution panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
	result, err := testphases.RunPipeline(cx, testphases.PhaseSymbolResolution, testCase.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testCase.InputPath, err)
		return
	}
	validator := &symbolResolutionValidator{t: t, testPath: testCase.InputPath}
	ast.Walk(validator, result.Package)
	// If we reach here, symbol resolution completed without panicking
	t.Logf("Symbol resolution completed successfully for %s", testCase.InputPath)
}

type symbolResolutionValidator struct {
	t        *testing.T
	testPath string
}

func (v *symbolResolutionValidator) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	if invocation, ok := node.(*ast.BLangInvocation); ok {
		if semantics.IsDeferredMethodSymbol(invocation.RawSymbol) {
			return nil
		}
	}
	// Check if this node should have a symbol resolved
	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		if !ast.SymbolIsSet(nodeWithSymbol) {
			v.t.Errorf("Symbol not resolved for %T at %s in %s",
				node, node.GetPosition(), v.testPath)
		}
	}
	if nodeWithScope, ok := node.(ast.NodeWithScope); ok {
		if nodeWithScope.Scope() == nil {
			v.t.Errorf("Scope not set for %T at %s in %s",
				node, node.GetPosition(), v.testPath)
		}
	}
	return v
}

func (v *symbolResolutionValidator) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return nil
	}
	// Check if this type descriptor should have a symbol resolved
	if typeWithSymbol, ok := typeData.TypeDescriptor.(ast.BNodeWithSymbol); ok {
		if !ast.SymbolIsSet(typeWithSymbol) {
			v.t.Errorf("Symbol not resolved for type %T at %s in %s",
				typeData.TypeDescriptor, typeData.TypeDescriptor.GetPosition(), v.testPath)
		}
	}
	return v
}
