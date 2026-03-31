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

package extern_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"

	_ "ballerina-lang-go/lib/rt"
)

const testDataDir = "testdata"

func TestExternValid(t *testing.T) {
	balFile := filepath.Join(testDataDir, "1-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))

	ballerinaHomePath := getBallerinaHomePath(t)
	ballerinaHomeFs := os.DirFS(ballerinaHomePath)

	result, err := projects.Load(fsys, ballerinaHomeFs, filepath.Base(absPath))
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	var stdoutBuf bytes.Buffer
	rt := runtime.NewRuntime()

	// Register println to capture output
	runtime.RegisterExternFunction(rt, "ballerina", "io", "println", func(args []values.BalValue) (values.BalValue, error) {
		var b strings.Builder
		visited := make(map[uintptr]bool)
		for _, arg := range args {
			b.WriteString(values.String(arg, visited))
		}
		b.WriteByte('\n')
		stdoutBuf.WriteString(b.String())
		return nil, nil
	})

	// Register foo() returns "$foo"
	runtime.RegisterExternFunction(rt, "$anon", "1-v", "foo", func(args []values.BalValue) (values.BalValue, error) {
		return "$foo", nil
	})

	// Register bar(a, b) returns a + ", " + b
	runtime.RegisterExternFunction(rt, "$anon", "1-v", "bar", func(args []values.BalValue) (values.BalValue, error) {
		a := values.String(args[0], nil)
		b := values.String(args[1], nil)
		return a + ", " + b, nil
	})

	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "$foo, $foo\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func TestExternTypeMismatchArg(t *testing.T) {
	balFile := filepath.Join(testDataDir, "2-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))

	ballerinaHomePath := getBallerinaHomePath(t)
	ballerinaHomeFs := os.DirFS(ballerinaHomePath)

	result, err := projects.Load(fsys, ballerinaHomeFs, filepath.Base(absPath))
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if !compilation.DiagnosticResult().HasErrors() {
		t.Fatal("expected compilation errors for type mismatch in arguments")
	}

	foundError := false
	for _, d := range compilation.DiagnosticResult().Diagnostics() {
		msg := fmt.Sprintf("%v", d)
		if strings.Contains(msg, "incompatible") || strings.Contains(msg, "type") {
			foundError = true
		}
	}
	if !foundError {
		t.Error("expected a type-related error diagnostic")
	}
}

func TestExternTypeMismatchReturn(t *testing.T) {
	balFile := filepath.Join(testDataDir, "3-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))

	ballerinaHomePath := getBallerinaHomePath(t)
	ballerinaHomeFs := os.DirFS(ballerinaHomePath)

	result, err := projects.Load(fsys, ballerinaHomeFs, filepath.Base(absPath))
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if !compilation.DiagnosticResult().HasErrors() {
		t.Fatal("expected compilation errors for type mismatch in return type")
	}

	foundError := false
	for _, d := range compilation.DiagnosticResult().Diagnostics() {
		msg := fmt.Sprintf("%v", d)
		if strings.Contains(msg, "incompatible") || strings.Contains(msg, "type") {
			foundError = true
		}
	}
	if !foundError {
		t.Error("expected a type-related error diagnostic")
	}
}

func TestExternHandle(t *testing.T) {
	balFile := filepath.Join(testDataDir, "4-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	result, err := projects.Load(fsys, nil, filepath.Base(absPath))
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()
	if compilation.DiagnosticResult().HasErrors() {
		for _, d := range compilation.DiagnosticResult().Diagnostics() {
			t.Logf("diagnostic: %v", d)
		}
		t.Fatal("compilation had errors")
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	var stdoutBuf bytes.Buffer
	rt := runtime.NewRuntime()

	runtime.RegisterExternFunction(rt, "ballerina", "io", "println", func(args []values.BalValue) (values.BalValue, error) {
		var b strings.Builder
		visited := make(map[uintptr]bool)
		for _, arg := range args {
			b.WriteString(values.String(arg, visited))
		}
		b.WriteByte('\n')
		stdoutBuf.WriteString(b.String())
		return nil, nil
	})

	type myHandle struct {
		data string
	}

	runtime.RegisterExternFunction(rt, "$anon", "4-v", "createHandle", func(args []values.BalValue) (values.BalValue, error) {
		return &myHandle{data: "handle_value"}, nil
	})

	runtime.RegisterExternFunction(rt, "$anon", "4-v", "useHandle", func(args []values.BalValue) (values.BalValue, error) {
		h := args[0].(*myHandle)
		return h.data, nil
	})

	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "handle_value\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func getBallerinaHomePath(t *testing.T) string {
	if balHome := os.Getenv(projects.BallerinaHomeEnvVar); balHome != "" {
		return balHome
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home: %v", err)
	}

	return filepath.Join(userHome, projects.UserHomeDirName)
}
