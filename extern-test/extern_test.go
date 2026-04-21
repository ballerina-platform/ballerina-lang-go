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
	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	result, err := directory.LoadProject(fsys, filepath.Base(absPath))
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
	result, err := directory.LoadProject(fsys, filepath.Base(absPath))
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
	result, err := directory.LoadProject(fsys, filepath.Base(absPath))
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

func TestDependentlyTyped(t *testing.T) {
	balFile := filepath.Join(testDataDir, "dependently-typed-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	result, err := directory.LoadProject(fsys, filepath.Base(absPath))
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
	tyCtx := semtypes.ContextFrom(rt.GetTypeEnv())

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

	runtime.RegisterExternFunction(rt, "$anon", "dependently-typed-v", "inferred", func(args []values.BalValue) (values.BalValue, error) {
		td, ok := args[1].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[1])
		}
		switch {
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.INT):
			return int64(1), nil
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.STRING):
			return "foo", nil
		}
		panic(values.NewErrorWithMessage("unsupported inferred typedesc constraint"))
	})

	runtime.RegisterExternFunction(rt, "$anon", "dependently-typed-v", "inferredSubType", func(args []values.BalValue) (values.BalValue, error) {
		td, ok := args[1].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[1])
		}
		if !semtypes.IsSubtype(tyCtx, td.Type, semtypes.INT) {
			panic(values.NewErrorWithMessage("inferredSubType requires typedesc<int>"))
		}
		return int64(1), nil
	})

	runtime.RegisterExternFunction(rt, "$anon", "dependently-typed-v", "inferredPartially", func(args []values.BalValue) (values.BalValue, error) {
		td, ok := args[1].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[1])
		}
		switch {
		case semtypes.IsSubtype(tyCtx, semtypes.INT, td.Type):
			return int64(0), nil
		case semtypes.IsSubtype(tyCtx, semtypes.STRING, td.Type):
			return "bar", nil
		}
		panic(values.NewErrorWithMessage("unsupported inferredPartially typedesc constraint"))
	})

	runtime.RegisterExternFunction(rt, "$anon", "dependently-typed-v", "shiftBy", func(args []values.BalValue) (values.BalValue, error) {
		src, ok := args[0].(*values.Map)
		if !ok {
			return nil, fmt.Errorf("expected record argument, got %T", args[0])
		}
		dx := args[1].(int64)
		dy := args[2].(int64)
		td, ok := args[3].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[3])
		}
		xVal, _ := src.Get("x")
		yVal, _ := src.Get("y")
		out := values.NewMap(td.Type)
		out.Put("x", xVal.(int64)+dx)
		out.Put("y", yVal.(int64)+dy)
		return out, nil
	})

	runtime.RegisterExternFunction(rt, "$anon", "dependently-typed-v", "inferredWithDefault", func(args []values.BalValue) (values.BalValue, error) {
		val, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("expected int argument, got %T", args[0])
		}
		td, ok := args[1].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[1])
		}
		switch {
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.INT):
			return val, nil
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.STRING):
			return fmt.Sprintf("%d", val), nil
		}
		panic(values.NewErrorWithMessage("unsupported inferredWithDefault typedesc constraint"))
	})

	if err := rt.Interpret(*birPkg); err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "1\nfoo\n1\n1\n0\nbar\n11\n22\n1\n6\n102\n42\n100\n7\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

func TestExternHandle(t *testing.T) {
	balFile := filepath.Join(testDataDir, "4-v.bal")
	absPath, err := filepath.Abs(balFile)
	if err != nil {
		t.Fatal(err)
	}

	fsys := os.DirFS(filepath.Dir(absPath))
	result, err := directory.LoadProject(fsys, filepath.Base(absPath))
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

// TestDependentlyTypedCrossModuleRoundtrip manually compiles a helper module
// containing a dependently-typed extern function, serializes its exported
// symbols and BIR, then in a fresh compiler environment deserializes those
// bytes and compiles the main module that imports the helper. Finally it
// interprets both BIR packages with the helper's native implementation
// registered and asserts the output is correct. The roundtrip validates that
// dependent-return and inferred-typedesc-default metadata survive
// serialization and drive correct type-resolution + BIR generation when the
// dependent-fn is declared in a dependency module.
func TestDependentlyTypedCrossModuleRoundtrip(t *testing.T) {
	helperBalPath := filepath.Join(testDataDir, "cross-module-dependent-fn-v", "modules", "helper", "helper.bal")
	mainBalPath := filepath.Join(testDataDir, "cross-module-dependent-fn-v", "main.bal")

	const (
		org         = "testorg"
		packageRoot = "crossmoduledependentfn"
		helperMod   = "crossmoduledependentfn.helper"
	)

	// Stage 1: compile the helper module standalone in env1.
	env1 := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	helperExported, helperBIR := compileSingleFileModule(t, env1, helperBalPath,
		model.Name(org),
		[]model.Name{model.Name(packageRoot), model.Name("helper")},
		nil,
		org,
	)

	// Stage 2: serialize helper's exported symbols and BIR.
	symBytes, err := symbolpool.Marshal(helperExported, helperBIR.TypeEnv)
	if err != nil {
		t.Fatalf("helper symbol Marshal: %v", err)
	}
	birBytes, err := bircodec.Marshal(helperBIR)
	if err != nil {
		t.Fatalf("helper BIR Marshal: %v", err)
	}

	// Stage 3: fresh env. Deserialize helper symbols and BIR.
	env2 := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	deserializedHelperExported, err := symbolpool.Unmarshal(env2, symBytes)
	if err != nil {
		t.Fatalf("helper symbol Unmarshal: %v", err)
	}
	deserializedHelperCtx := context.NewCompilerContext(env2)
	deserializedHelperBIR, err := bircodec.Unmarshal(deserializedHelperCtx, birBytes)
	if err != nil {
		t.Fatalf("helper BIR Unmarshal: %v", err)
	}

	// Stage 4: compile the main module in env2 against the deserialized helper
	// symbols. If dependent-fn metadata failed to survive serialization, type
	// resolution of `int a = helper:inferred(0)` would widen to VAL (or error)
	// and main compilation would fail.
	publicSymbols := map[semantics.PackageIdentifier]model.ExportedSymbolSpace{
		{OrgName: org, ModuleName: helperMod}: deserializedHelperExported,
	}
	_, mainBIR := compileSingleFileModule(t, env2, mainBalPath,
		model.Name(org),
		[]model.Name{model.Name(packageRoot)},
		publicSymbols,
		org,
	)

	// Stage 5: interpret [deserialized helper BIR, freshly compiled main BIR]
	// with helper's native registered. Validate stdout.
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
	tyCtx := semtypes.ContextFrom(rt.GetTypeEnv())
	runtime.RegisterExternFunction(rt, org, helperMod, "inferred", func(args []values.BalValue) (values.BalValue, error) {
		td, ok := args[1].(*values.TypeDesc)
		if !ok {
			return nil, fmt.Errorf("expected typedesc argument, got %T", args[1])
		}
		switch {
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.INT):
			return int64(1), nil
		case semtypes.IsSubtype(tyCtx, td.Type, semtypes.STRING):
			return "foo", nil
		}
		panic(values.NewErrorWithMessage("unsupported inferred typedesc constraint"))
	})

	for _, pkg := range []*bir.BIRPackage{deserializedHelperBIR, mainBIR} {
		if err := rt.Interpret(*pkg); err != nil {
			t.Fatalf("runtime error: %v", err)
		}
	}

	const expected = "1\nfoo\n"
	if stdoutBuf.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdoutBuf.String())
	}
}

// compileSingleFileModule parses a .bal file and runs the full compilation
// pipeline for it as its own module with the given package identity. The
// publicSymbols map is consulted for any imports. Returns the module's
// exported symbol space and BIR package.
func compileSingleFileModule(
	t *testing.T,
	env *context.CompilerEnvironment,
	balPath string,
	orgName model.Name,
	nameComps []model.Name,
	publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace,
	defaultOrg string,
) (model.ExportedSymbolSpace, *bir.BIRPackage) {
	t.Helper()
	absPath, err := filepath.Abs(balPath)
	if err != nil {
		t.Fatal(err)
	}
	cx := context.NewCompilerContext(env)
	st, err := parser.GetSyntaxTree(cx, absPath)
	if err != nil {
		t.Fatalf("parsing %s: %v", balPath, err)
	}
	cu := ast.GetCompilationUnit(cx, st)
	pkg := ast.ToPackage(cu)
	pkg.PackageID = cx.NewPackageID(orgName, nameComps, model.DEFAULT_VERSION)

	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx), publicSymbols, defaultOrg)
	exported := semantics.ResolveSymbols(cx, pkg, importedSymbols)
	assertNoDiagnostics(t, cx, "ResolveSymbols")
	semantics.ResolveTopLevelNodes(cx, pkg, importedSymbols)
	assertNoDiagnostics(t, cx, "ResolveTopLevelNodes")
	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	assertNoDiagnostics(t, cx, "ResolveLocalNodes")
	analyzer := semantics.NewSemanticAnalyzer(cx)
	analyzer.Analyze(pkg)
	assertNoDiagnostics(t, cx, "SemanticAnalyzer")
	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	assertNoDiagnostics(t, cx, "CreateControlFlowGraph")
	semantics.AnalyzeCFG(cx, pkg, cfg)
	assertNoDiagnostics(t, cx, "AnalyzeCFG")
	pkg = desugar.DesugarPackage(cx, pkg, importedSymbols)
	return exported, bir.GenBir(cx, pkg)
}

func assertNoDiagnostics(t *testing.T, cx *context.CompilerContext, stage string) {
	t.Helper()
	if !cx.HasDiagnostics() {
		return
	}
	for _, d := range cx.Diagnostics() {
		t.Logf("%s diagnostic: %s", stage, d)
	}
	t.Fatalf("%s produced diagnostics", stage)
}
