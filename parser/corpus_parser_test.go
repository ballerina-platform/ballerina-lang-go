// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package parser

import (
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/tools/text"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update expected JSON files")

// XML parser ignore list (135 files failing)
var xmlParserIgnoreList = []string{
	"action/client_resource_access_return_type_negative_test.bal",
	"bala/test_bala/readonly/test_selectively_immutable_type.bal",
	"bala/test_bala/types/xml_attribute_access_negative.bal",
	"bala/test_bala/types/xml_attribute_access.bal",
	"bala/test_projects/test_project_selectively_immutable/constructs.bal",
	"bala/test_projects/test_project/modules/dependently_typed/interop_funcs.bal",
	"bala/test_projects/test_project/modules/selectively_immutable/constructs.bal",
	"closures/var-mutability-closure.bal",
	"dataflow/analysis/dataflow-analysis-negative.bal",
	"dataflow/analysis/dataflow-analysis-semantics-negative.bal",
	"expressions/access/field_access_negative.bal",
	"expressions/access/field_access.bal",
	"expressions/access/xml_member_access_negative.bal",
	"expressions/access/xml_member_access.bal",
	"expressions/binaryoperations/add-operation-negative.bal",
	"expressions/binaryoperations/add-operation.bal",
	"expressions/binaryoperations/equal_and_not_equal_operation.bal",
	"expressions/binaryoperations/negative-type-test-expr-negative.bal",
	"expressions/binaryoperations/negative-type-test-expr.bal",
	"expressions/binaryoperations/ref_equal_and_not_equal_operation_negative.bal",
	"expressions/binaryoperations/ref_equal_and_not_equal_operation.bal",
	"expressions/binaryoperations/type-test-expr-negative.bal",
	"expressions/binaryoperations/type-test-expr.bal",
	"expressions/builtinoperations/clone-operation.bal",
	"expressions/builtinoperations/freeze-and-isfrozen.bal",
	"expressions/builtinoperations/length-operation.bal",
	"expressions/conversion/native-conversion-negative.bal",
	"expressions/conversion/native-conversion.bal",
	"expressions/elvis/elvis-expr-negative.bal",
	"expressions/elvis/elvis-expr.bal",
	"expressions/lambda/iterable/basic-iterable-with-variable-mutability.bal",
	"expressions/lambda/iterable/basic-iterable.bal",
	"expressions/let/let-expression-negative.bal",
	"expressions/let/let-expression-test.bal",
	"expressions/listconstructor/list_constructor_infer_type.bal",
	"expressions/mappingconstructor/mapping_constructor_infer_record.bal",
	"expressions/stamp/anydata-stamp-expr-test.bal",
	"expressions/stamp/negative/object-stamp-expr-negative-test.bal",
	"expressions/stamp/negative/union-stamp-expr-negative-test.bal",
	"expressions/stamp/negative/xml-stamp-expr-negative-test.bal",
	"expressions/stamp/union-stamp-expr-test.bal",
	"expressions/stamp/xml-stamp-expr-test.bal",
	"expressions/typecast/type_cast_expr_runtime_errors.bal",
	"expressions/typecast/type_cast_expr.bal",
	"expressions/typeof/typeof.bal",
	"functions/different-function-signatures-semantics-negative.bal",
	"functions/expr_bodied_functions.bal",
	"imports/InvalidAutoImportsTestProject/invalid-auto-imports-negative.bal",
	"imports/OverriddenPredeclaredImportsTestProject/overridden-xml.bal",
	"imports/PredeclaredImportsTestProject/predeclared-xml.bal",
	"isolated-objects/isolated_objects_isolation_negative.bal",
	"isolation-analysis/isolation_inference_with_objects_runtime_negative_1.bal",
	"javainterop/ballerina_types_as_interop_types.bal",
	"javainterop/ballerina_types_with_public_api.bal",
	"javainterop/dependently_typed_functions_bir_test.bal",
	"javainterop/dependently_typed_functions_test.bal",
	"javainterop/inferred_dependently_typed_func_signature_negative.bal",
	"javainterop/inferred_dependently_typed_func_signature.bal",
	"jvm/largeMethods3/main.bal",
	"jvm/too-large-method.bal",
	"jvm/too-large-object-field.bal",
	"jvm/too-large-object-method.bal",
	"jvm/too-large-package-variable.bal",
	"jvm/types.bal",
	"jvm/xml-literals-with-namespaces.bal",
	"jvm/xml.bal",
	"module.declarations/client-decl/client_decl_client_prefix_as_xmlns_prefix_negative_test.bal",
	"query/inner-queries.bal",
	"query/order-by-clause.bal",
	"query/query_ambiguous_type_negative.bal",
	"query/query_with_closures.bal",
	"query/query-action.bal",
	"query/query-expr-query-construct-type-negative.bal",
	"query/query-expr-with-query-construct-type.bal",
	"query/query-negative.bal",
	"query/string-query-expression-v2.bal",
	"query/xml-query-expression-negative.bal",
	"query/xml-query-expression.bal",
	"reachability-analysis/reachability_analysis.bal",
	"record/closed_record_type_inclusion.bal",
	"record/map_to_record.bal",
	"record/open_record_type_inclusion.bal",
	"record/readonly_record_fields.bal",
	"statements/arrays/array-fill-test.bal",
	"statements/arrays/array-test.bal",
	"statements/arrays/sealed_array.bal",
	"statements/assign/assign-stmt.bal",
	"statements/comment/comments.bal",
	"statements/compoundassignment/compound_assignment.bal",
	"statements/expression/expression-stmt2-semantics-negative.bal",
	"statements/ifelse/type-guard.bal",
	"typedefs/type-definitions.bal",
	"types/anydata/anydata_conversion_using_ternary.bal",
	"types/anydata/anydata_invalid_conversions.bal",
	"types/anydata/anydata_test.bal",
	"types/future/future_positive.bal",
	"types/never/never-type-negative.bal",
	"types/never/never-type.bal",
	"types/readonly/test_inherently_immutable_type.bal",
	"types/readonly/test_selectively_immutable_type_langlib_negative.bal",
	"types/readonly/test_selectively_immutable_type_negative.bal",
	"types/readonly/test_selectively_immutable_type.bal",
	"types/string/string-value-xml-test.bal",
	"types/table/record-constraint-table-value.bal",
	"types/table/record-type-table-key.bal",
	"types/table/table_key_field_value_test.bal",
	"types/table/table_key_violations.bal",
	"types/table/xml-type-table-key.bal",
	"types/tuples/tuple_basic_test.bal",
	"types/tuples/tuple_negative_test.bal",
	"types/xml/package_level_xml_literals.bal",
	"types/xml/xml_inline_large_literal.bal",
	"types/xml/xml_iteration_negative.bal",
	"types/xml/xml_iteration.bal",
	"types/xml/xml_step_expr_negative.bal",
	"types/xml/xml_text_to_string_conversion-negative.bal",
	"types/xml/xml_type_descriptor_negative.bal",
	"types/xml/xml_type_descriptor.bal",
	"types/xml/xml-attribute-access-lax-behavior.bal",
	"types/xml/xml-attribute-access-syntax-neg.bal",
	"types/xml/xml-attribute-access-syntax.bal",
	"types/xml/xml-attributes.bal",
	"types/xml/xml-element-access.bal",
	"types/xml/xml-indexed-access-negative.bal",
	"types/xml/xml-indexed-access.bal",
	"types/xml/xml-literals-negative.bal",
	"types/xml/xml-literals-with-namespaces.bal",
	"types/xml/xml-literals.bal",
	"types/xml/xml-native-functions.bal",
	"types/xml/xml-nav-access-negative-filter.bal",
	"types/xml/xml-nav-access-negative.bal",
	"types/xml/xml-nav-access-type-check-negative.bal",
	"types/xml/xml-navigation-access.bal",
	"variable/shadowing/shadowing.bal",
	"workers/basic-worker-actions.bal",
}

// Regex parser ignore list (8 files failing)
var regexParserIgnoreList = []string{
	"bala/test_bala/types/regexp_type_test.bal",
	"bala/test_projects/test_project_regexp/regexpTypes.bal",
	"jvm/largeMethods/modules/functions/large-functions.bal",
	"query/query_action_or_expr.bal",
	"query/simple-query-with-defined-type.bal",
	"types/regexp/regexp_type_test.bal",
	"types/regexp/regexp_value_negative_test.bal",
	"types/regexp/regexp_value_test.bal",
}

// Documentation parser ignore list (46 files failing)
var documentationParserIgnoreList = []string{
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
	"documentation/markdown_multiline_documentation.bal",
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

// shouldIgnoreFile checks if a file should be ignored based on the ignore lists
func shouldIgnoreFile(filePath string) bool {
	allIgnoreLists := [][]string{
		xmlParserIgnoreList,
		regexParserIgnoreList,
		documentationParserIgnoreList,
	}

	for _, ignoreList := range allIgnoreLists {
		for _, ignorePath := range ignoreList {
			if strings.HasSuffix(filePath, ignorePath) {
				return true
			}
		}
	}

	return false
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestParseCorpusFiles(t *testing.T) {
	if os.Getenv("GOARCH") == "wasm" {
		t.Skip("skipping parser testing wasm")
	}

	// Parser can parse all .bal files, not just -v.bal
	testPairs := test_util.GetTests(t, test_util.Parser, func(path string) bool {
		return true
	})

	// Create subtests for each file
	// Running in parallel for faster test execution
	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel() // Run in parallel for faster execution (native only)
			parseFile(t, testPair)
		})
	}
}

func TestJBalUnitTests(t *testing.T) {
	corpusDir := "./testdata/bal"
	if os.Getenv("GOARCH") == "wasm" {
		t.Skip("skipping parser testing wasm")
	}
	balFiles := getCorpusFiles(t, corpusDir)

	// Create subtests for each file
	// Running in parallel for faster test execution
	for _, balFile := range balFiles {
		// Skip files in ignore lists
		if shouldIgnoreFile(balFile) {
			t.Run(balFile, func(t *testing.T) {
				t.Skipf("Skipping file in ignore list: %s", balFile)
			})
			continue
		}

		// Create TestCase from file path
		testCase := createTestCase(t, balFile, corpusDir)

		t.Run(balFile, func(t *testing.T) {
			t.Parallel() // Run in parallel for faster execution (native only)
			parseFile(t, testCase)
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

func parseFile(t *testing.T, testCase test_util.TestCase) {
	// Catch any panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()

	// Read file content
	content, readErr := os.ReadFile(testCase.InputPath)
	if readErr != nil {
		t.Fatalf("error reading file: %v", readErr)
	}

	reader := text.CharReaderFromText(string(content))

	lexer := NewLexer(reader, nil)

	tokenReader := CreateTokenReader(*lexer, nil)

	ballerinaParser := NewBallerinaParserFromTokenReader(tokenReader, nil)

	ast := ballerinaParser.Parse()

	actualJSON := tree.GenerateJSON(ast)

	normalizedJSON := normalizeJSON(actualJSON)

	// If update flag is set, check if update is needed and update if necessary
	if *update {
		if test_util.UpdateIfNeeded(t, testCase.ExpectedPath, normalizedJSON, normalizeJSON) {
			t.Fatalf("Updated expected JSON file: %s", testCase.ExpectedPath)
		}
		return
	}

	expectedJSON := expectedJSON(t, testCase.ExpectedPath)

	// Compare JSON strings exactly (no tolerance for formatting differences)
	if normalizedJSON != expectedJSON {
		diff := getDiff(expectedJSON, normalizedJSON)
		t.Errorf("JSON mismatch for %s\nExpected file: %s\n%s", testCase.InputPath, testCase.ExpectedPath, diff)
		return

	}
}

// createTestCase creates a TestCase from a file path and base directory
func createTestCase(t *testing.T, filePath string, baseDir string) test_util.TestCase {
	// Get the relative path from baseDir
	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}

	// Replace "bal" directory with "parser" in the base directory path
	parserBaseDir := strings.Replace(baseDir, string(filepath.Separator)+"bal", string(filepath.Separator)+"parser", 1)

	// Construct the expected JSON path
	expectedJSONPath := filepath.Join(parserBaseDir, relPath)
	expectedJSONPath = strings.TrimSuffix(expectedJSONPath, ".bal") + ".json"

	return test_util.TestCase{
		Name:         filePath,
		InputPath:    filePath,
		ExpectedPath: expectedJSONPath,
	}
}

func normalizeJSON(jsonStr string) string {
	var obj interface{}
	normalizedJSON := jsonStr
	if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
		if normalized, err := json.MarshalIndent(obj, "", "  "); err == nil {
			normalizedJSON = string(normalized)
		}
	}
	return normalizedJSON
}

func expectedJSON(t *testing.T, expectedJSONPath string) string {
	// Check if file exists
	expectedJSONBytes, readErr := os.ReadFile(expectedJSONPath)
	if readErr != nil {
		t.Fatalf("error reading expected JSON file: %v", readErr)
		return ""
	}

	// File exists - normalize and compare
	expectedJSON := string(expectedJSONBytes)
	var expectedObj interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedObj); err == nil {
		if normalized, err := json.MarshalIndent(expectedObj, "", "  "); err == nil {
			expectedJSON = string(normalized)
		}
	}
	return expectedJSON
}

// getDiff generates a detailed diff string showing differences between expected and actual AST strings.
func getDiff(expectedAST, actualAST string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedAST, actualAST, false)
	return dmp.DiffPrettyText(diffs)
}
