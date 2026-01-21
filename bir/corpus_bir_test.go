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

package bir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var supportedSubsets = []string{"subset1"}

func getCorpusDir(t *testing.T) string {
	corpusBirDir := "../corpus/bir"
	if _, err := os.Stat(corpusBirDir); os.IsNotExist(err) {
		// Try alternative path (when running from project root)
		corpusBirDir = "./corpus/bir"
		if _, err := os.Stat(corpusBirDir); os.IsNotExist(err) {
			t.Skipf("Corpus directory not found (tried ../corpus/bir and ./corpus/bir), skipping test")
		}
	}
	return corpusBirDir
}

func getCorpusFiles(t *testing.T) []string {
	corpusBirDir := getCorpusDir(t)
	// Find all .bir files
	var birFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(corpusBirDir, subset)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".bir") {
				birFiles = append(birFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error walking corpus/bir/%s directory: %v", subset, err)
		}
	}

	if len(birFiles) == 0 {
		t.Fatalf("No .bir files found in %s", corpusBirDir)
	}
	return birFiles
}

func TestBIRPackageLoading(t *testing.T) {
	birFiles := getCorpusFiles(t)
	for _, birFile := range birFiles {
		t.Run(birFile, func(t *testing.T) {
			t.Parallel()
			testBIRPackageLoading(t, birFile)
		})
	}
}

func testBIRPackageLoading(t *testing.T, birFile string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while loading BIR package from %s: %v", birFile, r)
		}
	}()

	file, err := os.Open(birFile)
	if err != nil {
		t.Fatalf("failed to open test BIR file: %v", err)
	}
	defer file.Close()
	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil {
		t.Errorf("error loading BIR package from %s: %v", birFile, err)
		return
	}

	if pkg == nil {
		t.Errorf("BIR package is nil for %s", birFile)
		return
	}
}
