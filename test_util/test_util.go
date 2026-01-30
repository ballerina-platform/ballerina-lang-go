/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package test_util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestKind represents the type of corpus test
type TestKind int

const (
	AST TestKind = iota
	Parser
	BIR
)

// TestCase represents a test case: input file and expected output file
type TestCase struct {
	Name         string
	InputPath    string // Absolute path to .bal file
	ExpectedPath string // Absolute path to expected output (.txt or .json)
}

// GetValidTests returns all valid test pairs for the given test kind
// It only returns test cases where the input file ends with "-v.bal"
func GetValidTests(t *testing.T, kind TestKind) []TestCase {
	return GetTests(t, kind, func(path string) bool {
		return strings.HasSuffix(path, "-v.bal")
	})
}

// GetErrorTests returns all error test pairs for the given test kind
// It only returns test cases where the input file ends with "-e.bal"
func GetErrorTests(t *testing.T, kind TestKind) []TestCase {
	return GetTests(t, kind, func(path string) bool {
		return strings.HasSuffix(path, "-e.bal")
	})
}

// GetTests returns test pairs for the given test kind, filtered by the provided function
func GetTests(t *testing.T, kind TestKind, filterFunc func(string) bool) []TestCase {
	inputBaseDirAlt := "bal"
	var outputBaseDir string
	var outputExt string
	var subsets []string

	switch kind {
	case AST:
		outputBaseDir = "ast"
		outputExt = ".txt"
		subsets = []string{"subset1"}

	case Parser:
		outputBaseDir = "parser"
		outputExt = ".json"
		// Parser processes entire corpus tree, no subset restriction

	case BIR:
		outputBaseDir = "bir"
		outputExt = ".txt"
		subsets = []string{"subset1"}
	}

	// Resolve input and output directories
	resolvedInputDir, resolvedOutputDir := resolveDir(t, inputBaseDirAlt, outputBaseDir)

	// Discover files
	files := discoverFiles(t, resolvedInputDir, subsets, filterFunc)

	// Create test pairs with expected paths
	testPairs := make([]TestCase, 0, len(files))
	for _, inputPath := range files {
		expectedPath := computeExpectedPath(inputPath, resolvedInputDir, resolvedOutputDir, outputExt)
		relPath, _ := filepath.Rel(resolvedInputDir, inputPath)
		testPairs = append(testPairs, TestCase{
			InputPath:    inputPath,
			ExpectedPath: expectedPath,
			Name:         relPath,
		})
	}

	return testPairs
}

// resolveDir resolves the input and output directories
// It tries ../corpus/<inputBaseDir> first, then ./corpus/<inputBaseDir>
func resolveDir(t *testing.T, inputBaseDir, outputBaseDir string) (string, string) {
	// Try ../corpus first (when running from package dir)
	inputDir := filepath.Join("..", "corpus", inputBaseDir)
	if _, err := os.Stat(inputDir); err == nil {
		outputDir := filepath.Join("..", "corpus", outputBaseDir)
		absInput, _ := filepath.Abs(inputDir)
		absOutput, _ := filepath.Abs(outputDir)
		return absInput, absOutput
	}

	// Fall back to ./corpus (when running from repo root)
	inputDir = filepath.Join(".", "corpus", inputBaseDir)
	if _, err := os.Stat(inputDir); err == nil {
		outputDir := filepath.Join(".", "corpus", outputBaseDir)
		absInput, _ := filepath.Abs(inputDir)
		absOutput, _ := filepath.Abs(outputDir)
		return absInput, absOutput
	}

	t.Fatalf("Could not find corpus directory")
	return "", ""
}

// discoverFiles walks the directory tree and collects all .bal files that match the filter
func discoverFiles(t *testing.T, baseDir string, subsets []string, filterFunc func(string) bool) []string {
	var files []string

	// If subsets are specified, only walk those subdirectories
	if len(subsets) > 0 {
		for _, subset := range subsets {
			subsetDir := filepath.Join(baseDir, subset)
			if _, err := os.Stat(subsetDir); err != nil {
				t.Logf("Warning: subset directory %s does not exist, skipping", subsetDir)
				continue
			}
			subsetFiles := walkDir(t, subsetDir, filterFunc)
			files = append(files, subsetFiles...)
		}
	} else {
		// Walk entire base directory
		files = walkDir(t, baseDir, filterFunc)
	}

	return files
}

// walkDir recursively walks a directory and collects all .bal files that match the filter
func walkDir(t *testing.T, dir string, filterFunc func(string) bool) []string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".bal") {
			return nil
		}
		if filterFunc != nil && !filterFunc(path) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk directory %s: %v", dir, err)
	}

	return files
}

// computeExpectedPath converts an input path to the expected output path
// For example: corpus/bal/subset1/foo/bar.bal -> corpus/ast/subset1/foo/bar.txt
func computeExpectedPath(inputPath, inputBaseDir, outputBaseDir, outputExt string) string {
	// Get relative path from input base
	relPath, _ := filepath.Rel(inputBaseDir, inputPath)

	// Replace .bal extension with output extension
	relPath = strings.TrimSuffix(relPath, ".bal") + outputExt

	// Join with output base directory
	return filepath.Join(outputBaseDir, relPath)
}
