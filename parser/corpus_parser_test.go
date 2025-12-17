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
	"ballerina-lang-go/parser/internal"
	"ballerina-lang-go/tools/text"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
func shouldIgnoreFile(filePath string, corpusBalDir string) bool {
	// Get relative path from corpusBalDir
	relPath, err := filepath.Rel(corpusBalDir, filePath)
	if err != nil {
		return false
	}
	// Normalize path separators to forward slashes for comparison
	relPath = filepath.ToSlash(relPath)

	// Check all ignore lists
	allIgnoreLists := [][]string{
		xmlParserIgnoreList,
		regexParserIgnoreList,
		documentationParserIgnoreList,
	}

	for _, ignoreList := range allIgnoreLists {
		for _, ignorePath := range ignoreList {
			if relPath == ignorePath {
				return true
			}
		}
	}

	return false
}

func TestParseCorpusFiles(t *testing.T) {
	// Try both relative paths - from package directory and from project root
	corpusBalDir := "../corpus/bal"
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		// Try alternative path (when running from project root)
		corpusBalDir = "./corpus/bal"
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bal and ./corpus/bal), skipping test")
		}
	}

	// Find all .bal files
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

	// Create subtests for each file
	// Running in parallel for faster test execution
	total := len(balFiles)
	for i, balFile := range balFiles {
		balFile := balFile // capture loop variable
		index := i + 1

		// Skip files in ignore lists
		if shouldIgnoreFile(balFile, corpusBalDir) {
			t.Run(balFile, func(t *testing.T) {
				t.Skipf("Skipping file in ignore list: %s", balFile)
			})
			continue
		}

		t.Run(balFile, func(t *testing.T) {
			t.Parallel() // Run in parallel for faster execution
			parseFile(t, balFile, index, total)
		})
	}
}

func parseFile(t *testing.T, filePath string, index int, total int) {
	// Print file name at the beginning (before parsing) so we can see it even if stack overflow occurs
	fmt.Fprintf(os.Stderr, "[%d/%d] %s ... ", index, total, filePath)

	// Use a helper function to catch panics completely
	err := parseFileWithRecovery(filePath)

	// Print success/failure at the end
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILED: %v\n", err)
		t.Errorf("FAILED: %s - %v", filePath, err)
	} else {
		fmt.Fprintf(os.Stderr, "SUCCESS\n")
	}
}

// parseFileWithRecovery parses a file and returns any error or panic as an error
func parseFileWithRecovery(filePath string) (err error) {
	// Catch any panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Read file content
	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return fmt.Errorf("error reading file: %w", readErr)
	}

	// Create CharReader from file content
	reader := text.CharReaderFromText(string(content))

	// Create Lexer (no debug context needed for tests)
	lexer := NewLexer(reader, nil)

	// Create TokenReader from Lexer
	tokenReader := CreateTokenReader(*lexer, nil)

	// Create Parser from TokenReader
	ballerinaParser := NewBallerinaParserFromTokenReader(tokenReader, nil)

	// Parse the entire file - this may panic
	ast := ballerinaParser.Parse()

	// Verify that Parse() returns a non-nil STNode
	if ast == nil {
		return fmt.Errorf("Parse() returned nil AST")
	}

	// Verify it's a valid STNode by checking its Kind
	if ast.Kind() == 0 {
		return fmt.Errorf("Parse() returned AST with invalid Kind (0)")
	}

	// Generate JSON from the parsed AST
	actualJSON := internal.GenerateJSON(ast)

	// Determine expected JSON file path
	// Replace .bal with .json and change directory from corpus/bal to corpus/parser
	expectedJSONPath := strings.TrimSuffix(filePath, ".bal") + ".json"
	expectedJSONPath = strings.Replace(expectedJSONPath, string(filepath.Separator)+"corpus"+string(filepath.Separator)+"bal"+string(filepath.Separator), string(filepath.Separator)+"corpus"+string(filepath.Separator)+"parser"+string(filepath.Separator), 1)

	// Read expected JSON file
	expectedJSONBytes, readErr := os.ReadFile(expectedJSONPath)
	if readErr != nil {
		// If expected JSON file doesn't exist, skip this file
		if os.IsNotExist(readErr) {
			return fmt.Errorf("expected JSON file not found: %s (skipping)", expectedJSONPath)
		}
		return fmt.Errorf("error reading expected JSON file: %w", readErr)
	}

	expectedJSON := string(expectedJSONBytes)

	// Normalize both JSON strings by parsing and re-marshaling to handle whitespace differences
	var expectedObj, actualObj interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedObj); err == nil {
		if normalized, err := json.MarshalIndent(expectedObj, "", "  "); err == nil {
			expectedJSON = string(normalized)
		}
	}
	if err := json.Unmarshal([]byte(actualJSON), &actualObj); err == nil {
		if normalized, err := json.MarshalIndent(actualObj, "", "  "); err == nil {
			actualJSON = string(normalized)
		}
	}

	// Compare JSON strings exactly (no tolerance for formatting differences)
	if actualJSON != expectedJSON {
		// Split into lines for line-by-line comparison
		expectedLines := strings.Split(expectedJSON, "\n")
		actualLines := strings.Split(actualJSON, "\n")

		// Build detailed diff showing line numbers and differences
		var diffBuilder strings.Builder
		diffBuilder.WriteString("\nJSON mismatch - showing differences:\n\n")

		maxLines := len(expectedLines)
		if len(actualLines) > maxLines {
			maxLines = len(actualLines)
		}

		diffCount := 0
		const maxDiffsToShow = 20

		// Show line-by-line differences
		for i := 0; i < maxLines && diffCount < maxDiffsToShow; i++ {
			lineNum := i + 1
			expectedLine := ""
			actualLine := ""

			if i < len(expectedLines) {
				expectedLine = expectedLines[i]
			}
			if i < len(actualLines) {
				actualLine = actualLines[i]
			}

			if expectedLine != actualLine {
				diffCount++
				diffBuilder.WriteString(fmt.Sprintf("Line %d:\n", lineNum))
				if expectedLine == "" {
					diffBuilder.WriteString("  Expected: (empty)\n")
				} else {
					diffBuilder.WriteString(fmt.Sprintf("  Expected: %s\n", expectedLine))
				}
				if actualLine == "" {
					diffBuilder.WriteString("  Actual:   (empty)\n\n")
				} else {
					diffBuilder.WriteString(fmt.Sprintf("  Actual:   %s\n\n", actualLine))
				}
			}
		}

		if diffCount >= maxDiffsToShow {
			diffBuilder.WriteString(fmt.Sprintf("... (showing first %d differences, more exist)\n", maxDiffsToShow))
		}

		diffBuilder.WriteString(fmt.Sprintf("Total lines different: %d+\n", diffCount))
		diffBuilder.WriteString("Use diff tool for full comparison\n")

		return fmt.Errorf("JSON mismatch for %s\nExpected file: %s\n%s", filePath, expectedJSONPath, diffBuilder.String())
	}

	return nil
}
