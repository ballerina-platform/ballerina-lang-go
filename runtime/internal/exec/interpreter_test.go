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

// These two test cases cannot be expressed as corpus/bal -p tests: when
// init()/main() returns an error value (rather than panicking), ctx.PopFrame()
// has already run before getFormattedError is called, so no stack frame exists
// in the output. The corpus @panic annotation validator requires a frame, so
// they live here instead.
package exec_test

import (
	"os"
	"strings"
	"testing"

	"ballerina-lang-go/context"
	"ballerina-lang-go/runtime/internal/exec"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testphases"
)

// compileAndInterpret compiles inline Ballerina source and interprets it.
// Returns the error from Interpret, or nil on success.
func compileAndInterpret(t *testing.T, src string) error {
	t.Helper()
	tmp, err := os.CreateTemp("", "exec-test-*.bal")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(src); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	tmp.Close()

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	cx := context.NewCompilerContext(env)
	result, err := testphases.RunPipeline(cx, testphases.PhaseBIR, tmp.Name())
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}
	if cx.HasDiagnostics() {
		t.Fatalf("compilation produced unexpected diagnostics")
	}
	if result.BIRPackage == nil {
		t.Fatal("BIR package is nil after compilation")
	}
	return exec.Interpret(*result.BIRPackage, modules.NewRegistry())
}

// TestInterpret_InitFunctionReturnsError verifies that when init() explicitly
// returns an error value, Interpret surfaces it (exercises interpreter.go L36-38).
// This cannot be a corpus -p test: no stack frame is emitted for return-error
// (as opposed to panic), so the @panic annotation validator would fail.
func TestInterpret_InitFunctionReturnsError(t *testing.T) {
	err := compileAndInterpret(t, `
function init() returns error? {
    return error("initialization failed");
}
public function main() {
}
`)
	if err == nil {
		t.Fatal("expected non-nil error when init function returns an error, got nil")
	}
	if !strings.Contains(err.Error(), "initialization failed") {
		t.Errorf("expected error message to contain 'initialization failed', got: %s", err.Error())
	}
}

// TestInterpret_MainFunctionReturnsError verifies that when main() explicitly
// returns an error value, Interpret surfaces it (exercises interpreter.go L49-51).
// Same reason as TestInterpret_InitFunctionReturnsError for not being a corpus test.
func TestInterpret_MainFunctionReturnsError(t *testing.T) {
	err := compileAndInterpret(t, `
public function main() returns error? {
    return error("main failed");
}
`)
	if err == nil {
		t.Fatal("expected non-nil error when main returns an error, got nil")
	}
	if !strings.Contains(err.Error(), "main failed") {
		t.Errorf("expected error message to contain 'main failed', got: %s", err.Error())
	}
}
