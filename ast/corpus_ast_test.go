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

package ast

import (
	"flag"
	"fmt"
	"testing"

	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestASTGeneration(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testASTGeneration(t, testPair)
		})
	}
}

func testASTGeneration(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while testing AST generation for %s: %v", testCase.InputPath, r)
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
	}
	compilationUnit := GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", testCase.InputPath)
	}
	prettyPrinter := PrettyPrinter{}
	actualAST := prettyPrinter.Print(compilationUnit)

	// If update flag is set, update expected file
	if *update {
		if test_util.UpdateIfNeeded(t, testCase.ExpectedPath, actualAST) {
			t.Errorf("updated expected AST file: %s", testCase.ExpectedPath)
		}
		return
	}

	// Read expected AST file
	expectedAST := test_util.ReadExpectedFile(t, testCase.ExpectedPath)

	// Compare AST strings exactly
	if actualAST != expectedAST {
		diff := getDiff(expectedAST, actualAST)
		t.Errorf("AST mismatch for %s\nExpected file: %s\n%s", testCase.InputPath, testCase.ExpectedPath, diff)
		return
	}
}

var update = flag.Bool("update", false, "update expected AST files")

// getDiff generates a detailed diff string showing differences between expected and actual AST strings.
func getDiff(expectedAST, actualAST string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedAST, actualAST, false)
	return dmp.DiffPrettyText(diffs)
}

// walkTestVisitor tracks node types visited during Walk traversal
type walkTestVisitor struct {
	visitedTypes map[string]int
	nodeCount    int
}

func (v *walkTestVisitor) Visit(node BLangNode) Visitor {
	if node == nil {
		return nil
	}
	v.nodeCount++
	typeName := fmt.Sprintf("%T", node)
	v.visitedTypes[typeName]++
	return v
}

func (v *walkTestVisitor) VisitTypeData(typeData *model.TypeData) Visitor {
	return v
}

func TestWalkTraversal(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidTests(t, test_util.AST)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testWalkTraversal(t, testPair)
		})
	}
}

func testWalkTraversal(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Walk panicked for %s: %v", testCase.InputPath, r)
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
	compilationUnit := GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", testCase.InputPath)
		return
	}

	visitor := &walkTestVisitor{visitedTypes: make(map[string]int)}
	Walk(visitor, compilationUnit)

	if visitor.nodeCount == 0 {
		t.Errorf("Walk visited 0 nodes for %s", testCase.InputPath)
	}

	if testing.Verbose() {
		t.Logf("File: %s, Total nodes: %d", testCase.InputPath, visitor.nodeCount)
		for typeName, count := range visitor.visitedTypes {
			t.Logf("  %s: %d nodes", typeName, count)
		}
	}
}
