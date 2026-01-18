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
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	supportedSubsets = []string{"subset1"}
)

func getCorpusDir(t *testing.T) string {
	corpusBalDir := "../corpus/bal"
	if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
		// Try alternative path (when running from project root)
		corpusBalDir = "./corpus/bal"
		if _, err := os.Stat(corpusBalDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bal and ./corpus/bal), skipping test")
		}
	}
	return corpusBalDir
}

func getCorpusFiles(t *testing.T) []string {

	corpusBalDir := getCorpusDir(t)
	// Find all .bal files
	var balFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(corpusBalDir, subset)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".bal") {
				balFiles = append(balFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error walking corpus/bal/%s directory: %v", subset, err)
		}
	}

	if len(balFiles) == 0 {
		t.Fatalf("No .bal files found in %s", corpusBalDir)
	}
	return balFiles
}

func TestASTGeneration(t *testing.T) {
	balFiles := getCorpusFiles(t)
	for _, balFile := range balFiles {
		t.Run(balFile, func(t *testing.T) {
			t.Parallel()
			testASTGeneration(t, balFile)
		})
	}
}

func testASTGeneration(t *testing.T, balFile string) {
	if !strings.HasSuffix(balFile, "-v.bal") {
		t.Skipf("Skipping %s", balFile)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while testing AST generation for %s: %v", balFile, r)
		}
	}()

	debugCtx := debugcommon.DebugContext{
		Channel: make(chan string),
	}
	cx := context.NewCompilerContext()
	syntaxTree, err := parser.GetSyntaxTree(&debugCtx, balFile)
	if err != nil {
		t.Errorf("error getting syntax tree for %s: %v", balFile, err)
	}
	compilationUnit := GetCompilationUnit(cx, syntaxTree)
	if compilationUnit == nil {
		t.Errorf("compilation unit is nil for %s", balFile)
	}
}
