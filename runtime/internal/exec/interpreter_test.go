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

package exec_test

import (
	"os"
	"strings"
	"testing"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/runtime/internal/exec"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testphases"
)

// compileBIR compiles inline Ballerina source to a BIR package.
// t.Fatal is called on any compilation failure.
func compileBIR(t *testing.T, src string) bir.BIRPackage {
	t.Helper()
	tmp, err := os.CreateTemp("", "exec-test-*.bal")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			t.Fatalf("remove temp file: %v", err)
		}
	}()
	if _, err := tmp.WriteString(src); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}

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
	return *result.BIRPackage
}

// TestInterpret_InitFunctionReturnsError verifies that when a module-level
// init() function explicitly returns an error value, Interpret surfaces it
// as a formatted error (exercises interpreter.go L36-38).
func TestInterpret_InitFunctionReturnsError(t *testing.T) {
	pkg := compileBIR(t, `
function init() returns error? {
    return error("initialization failed");
}

public function main() {
}
`)

	err := exec.Interpret(pkg, modules.NewRegistry())
	if err == nil {
		t.Fatal("expected non-nil error when init function returns an error, got nil")
	}
	if !strings.Contains(err.Error(), "initialization failed") {
		t.Errorf("expected error message to contain 'initialization failed', got: %s", err.Error())
	}
}

// TestInterpret_MainSucceeds verifies a program with only a main function
// (no explicit init) executes without error (exercises the main-function path).
func TestInterpret_MainSucceeds(t *testing.T) {
	pkg := compileBIR(t, `
public function main() {
}
`)
	if err := exec.Interpret(pkg, modules.NewRegistry()); err != nil {
		t.Fatalf("expected nil error for successful main, got: %s", err.Error())
	}
}

// TestInterpret_MainFunctionReturnsError verifies that when main() explicitly
// returns an error, Interpret surfaces it (exercises interpreter.go L49-51).
func TestInterpret_MainFunctionReturnsError(t *testing.T) {
	pkg := compileBIR(t, `
public function main() returns error? {
    return error("main failed");
}
`)
	err := exec.Interpret(pkg, modules.NewRegistry())
	if err == nil {
		t.Fatal("expected non-nil error when main returns an error, got nil")
	}
	if !strings.Contains(err.Error(), "main failed") {
		t.Errorf("expected error message to contain 'main failed', got: %s", err.Error())
	}
}

// TestInterpret_InitFunctionPanics verifies that a panic inside init() is
// caught by the defer/recover block and returned as an error (exercises L30-34).
func TestInterpret_InitFunctionPanics(t *testing.T) {
	pkg := compileBIR(t, `
function init() {
    panic error("init panic");
}

public function main() {
}
`)
	err := exec.Interpret(pkg, modules.NewRegistry())
	if err == nil {
		t.Fatal("expected non-nil error when init panics, got nil")
	}
	if !strings.Contains(err.Error(), "init panic") {
		t.Errorf("expected error message to contain 'init panic', got: %s", err.Error())
	}
}

// TestInterpret_MainFunctionPanics verifies that a panic inside main() is
// caught by the defer/recover block and returned as an error (exercises L43-46).
func TestInterpret_MainFunctionPanics(t *testing.T) {
	pkg := compileBIR(t, `
public function main() {
    panic error("main panic");
}
`)
	err := exec.Interpret(pkg, modules.NewRegistry())
	if err == nil {
		t.Fatal("expected non-nil error when main panics, got nil")
	}
	if !strings.Contains(err.Error(), "main panic") {
		t.Errorf("expected error message to contain 'main panic', got: %s", err.Error())
	}
}
