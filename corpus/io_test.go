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

package corpus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIOFileWriteStringAndReadString(t *testing.T) {
	outputFile := filepath.Join(externTestDataDir, "io-file-string-output.txt")
	_ = os.Remove(outputFile)
	t.Cleanup(func() { _ = os.Remove(outputFile) })

	run := runIntegrationCase(filepath.Join(externTestDataDir, "io-file-string-v.bal"))
	if run.stderr != "" {
		t.Fatalf("unexpected stderr: %s", run.stderr)
	}
	if run.stdout != "true\n" {
		t.Fatalf("unexpected stdout: %q", run.stdout)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if string(content) != "Hello\nWorld" {
		t.Fatalf("unexpected file content: %q", content)
	}
}
