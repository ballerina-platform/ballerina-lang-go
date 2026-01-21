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
	"testing"
)

const testBIR = "testdata/input/bal_test.bir"

func TestLoadBIRPackageFromReader(t *testing.T) {
	file, err := os.Open(testBIR)
	if err != nil {
		t.Fatalf("failed to open test BIR file: %v", err)
	}
	defer file.Close()
	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil || pkg == nil {
		t.Fatalf("failed to load BIR package: %v", err)
	}
	prettyPrinter := PrettyPrinter{}
	actualOutput := prettyPrinter.Print(*pkg)
	expectedOutput := expectedOutput(t)
	if actualOutput != expectedOutput {
		t.Fatalf("actual output does not match expected output: %v", err)
	}
	t.Logf("loaded BIR package: %v", pkg)
}

func expectedOutput(t *testing.T) string {
	outputFile, err := os.ReadFile("testdata/output/bal_test.bir")
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}
	return string(outputFile)
}
