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

package main

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBalRunDumpFlags(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOARCH == "wasm" {
		t.Skip("skipping CLI integration test on WASM (js/wasm)")
	}

	singleBal := filepath.Join("corpus", "cli", "testdata", "run", "single-bal-files", "run-and-print.bal")

	tests := []struct {
		name string
		flag string
		file string
	}{
		{"dump-bir", "--dump-bir", "dump-bir.txtar"},
		{"dump-st", "--dump-st", "dump-st.txtar"},
		{"dump-tokens", "--dump-tokens", "dump-tokens.txtar"},
		{"dump-ast", "--dump-ast", "dump-ast.txtar"},
		{"dump-cfg", "--dump-cfg", "dump-cfg.txtar"},
	}
	flag.Parse()
	balBin, repoRoot, coverDir := integrationTestBal(t, true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBalCommandMatchesTxtarFragmentsForBinary(t, balBin, repoRoot, coverDir,
				[]string{"run", tt.flag, singleBal},
				"run-dump-flags", tt.file)
		})
	}
}
