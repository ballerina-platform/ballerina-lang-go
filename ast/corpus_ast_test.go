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
	"os"
	"path/filepath"
	"strings"
	"testing"

	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestASTGeneration(t *testing.T) {
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
	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
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
	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())
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

func TestJBalUnitTests(t *testing.T) {
	corpusDir := "./testdata/bal"
	if os.Getenv("GOARCH") == "wasm" {
		t.Skip("skipping AST testing wasm")
	}
	balFiles := getCorpusFiles(t, corpusDir)

	for _, balFile := range balFiles {
		// Skip files in ignore lists
		if shouldIgnoreFile(balFile) {
			t.Run(balFile, func(t *testing.T) {
				t.Skipf("Skipping file in ignore list: %s", balFile)
			})
			continue
		}

		testCase := createASTTestCase(t, balFile, corpusDir)

		t.Run(balFile, func(t *testing.T) {
			t.Parallel()
			testASTFile(t, testCase)
		})
	}
}

func getCorpusFiles(t *testing.T, corpusBalDir string) []string {
	var balFiles []string
	err := filepath.Walk(corpusBalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".bal") {
			balFiles = append(balFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking corpus/bal directory: %v", err)
	}

	if len(balFiles) == 0 {
		t.Fatalf("No .bal files found in %s", corpusBalDir)
	}
	return balFiles
}

// AST markdown documentation ignore list (tests that fail due to unimplemented panics)
var astMarkdownIgnoreList = []string{
	"annotations/deprecation_annotation_crlf.bal",
	"annotations/deprecation_annotation_negative.bal",
	"annotations/deprecation_annotation.bal",
	"bala/test_projects/test_documentation/test_documentation_symbol.bal",
	"bala/test_projects/test_project_errors/errors.bal",
	"bala/test_projects/test_project/deprecation_annotation.bal",
	"bala/test_projects/test_project/modules/errors/errors.bal",
	"documentation/default_value_initialization/main.bal",
	"documentation/deprecated_annotation_project/main.bal",
	"documentation/docerina_project/main.bal",
	"documentation/docerina_project/modules/world/world.bal",
	"documentation/errors_project/errors.bal",
	"documentation/markdown_annotation.bal",
	"documentation/markdown_constant.bal",
	"documentation/markdown_doc_inline_triple.bal",
	"documentation/markdown_doc_inline.bal",
	"documentation/markdown_finite_types.bal",
	"documentation/markdown_function_special.bal",
	"documentation/markdown_function.bal",
	"documentation/markdown_multiple.bal",
	"documentation/markdown_native_function.bal",
	"documentation/markdown_negative.bal",
	"documentation/markdown_object.bal",
	"documentation/markdown_on_disallowed_constructs.bal",
	"documentation/markdown_on_method_object_type_def.bal",
	"documentation/markdown_service.bal",
	"documentation/markdown_type.bal",
	"documentation/markdown_with_lambda.bal",
	"documentation/multi_line_docs_project/main.bal",
	"documentation/record_object_fields_project/main.bal",
	"documentation/type_models_project/type_models.bal",
	"enums/enum_metadata_test.bal",
	"expressions/naturalexpr/natural_expr.bal",
	"jvm/largePackage/modules/records/bigRecord2.bal",
	"jvm/largePackage/modules/records/bigRecord3.bal",
	"object/object_annotation.bal",
	"object/object_doc_annotation.bal",
	"object/object_documentation_negative.bal",
	"record/record_annotation.bal",
	"record/record_doc_annotation.bal",
	"record/record_documentation_negative.bal",
	"runtime/api/types/modules/typeref/typeref.bal",
	"statements/vardeclr/module_error_var_decl_annotation_negetive.bal",
	"statements/vardeclr/module_record_var_decl_annotation_negetive.bal",
	"statements/vardeclr/module_tuple_var_decl_annotation_negetive.bal",
}

func shouldIgnoreFile(filePath string) bool {
	for _, ignorePath := range astMarkdownIgnoreList {
		if strings.HasSuffix(filePath, ignorePath) {
			return true
		}
	}
	return false
}

// createASTTestCase creates a TestCase from a file path and base directory
func createASTTestCase(t *testing.T, filePath string, baseDir string) test_util.TestCase {
	// Get the relative path from baseDir
	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}

	// Replace "bal" directory with "ast" in the base directory path
	astBaseDir := strings.Replace(baseDir, string(filepath.Separator)+"bal", string(filepath.Separator)+"ast", 1)

	// Construct the expected AST path
	expectedASTPath := filepath.Join(astBaseDir, relPath)
	expectedASTPath = strings.TrimSuffix(expectedASTPath, ".bal") + ".txt"

	return test_util.TestCase{
		Name:         filePath,
		InputPath:    filePath,
		ExpectedPath: expectedASTPath,
	}
}

func testASTFile(t *testing.T, testCase test_util.TestCase) {
	// Catch any panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, testCase.InputPath)
	if err != nil {
		t.Fatalf("error getting syntax tree for %s: %v", testCase.InputPath, err)
	}
	compilationUnit := GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Fatalf("compilation unit is nil for %s", testCase.InputPath)
	}
	prettyPrinter := PrettyPrinter{}
	actualAST := prettyPrinter.Print(compilationUnit)

	// If update flag is set, update expected file
	if *update {
		if test_util.UpdateIfNeeded(t, testCase.ExpectedPath, actualAST) {
			t.Errorf("Updated expected AST file: %s", testCase.ExpectedPath)
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
