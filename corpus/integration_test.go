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
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

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
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/values"

	_ "ballerina-lang-go/lib/rt"
)

const (
	corpusProjectBaseDir            = "../corpus/project"
	corpusProjectIntegrationBaseDir = "../corpus/integration/project"

	externOrgName    = "ballerina"
	externModuleName = "io"
	externFuncName   = "println"

	panicPrefix = "panic: "
)

var (
	update = flag.Bool("update", false, "update corpus integration test outputs")

	// Skip tests that cause unrecoverable Go runtime errors
	skipIntegrationTests = []string{
		"subset7/07-function/assign12-v.bal",
		"subset7/07-function/assign17-v.bal",
		"subset7/07-function/record3-v.bal",
		"subset7/07-function/record4-v.bal",
		"subset7/07-function/typeCast1-v.bal",
		"subset7/07-function/typeCast2-v.bal",
		"subset7/07-function/typeCast3-p.bal",
		"subset7/07-closure/equality4-v.bal",
		"subset7/07-closure/typeCast1-v.bal",
		"subset7/07-closure/typeCast2-p.bal",
	}
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

type testResult struct {
	success        bool
	expectedStdout string
	actualStdout   string
	expectedStderr string
	actualStderr   string
}

func TestIntegration(t *testing.T) {
	testPairs := test_util.GetTests(t, test_util.Integration, func(path string) bool {
		return true
	})

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testIntegration(t, testPair)
		})
	}
}

func TestProjectIntegration(t *testing.T) {
	if _, err := os.Stat(corpusProjectBaseDir); os.IsNotExist(err) {
		return
	}

	projectDirs := findProjectDirs(corpusProjectBaseDir)

	for _, projDir := range projectDirs {
		dirName := filepath.Base(projDir)
		txtarPath := filepath.Join(corpusProjectIntegrationBaseDir, dirName+".txtar")

		t.Run(dirName, func(t *testing.T) {
			t.Parallel()
			testProjectIntegration(t, dirName, projDir, txtarPath)
		})
	}
}

func testIntegration(t *testing.T, testPair test_util.TestCase) {
	if isTestSkipped(testPair) {
		t.Skipf("Skipping integration test for %s", testPair.InputPath)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", testPair.InputPath, r)
		}
	}()

	if *update {
		stdout, stderr := runIntegrationCase(testPair.InputPath)
		if test_util.UpdateTxtarArchiveIfNeeded(t, testPair.ExpectedPath, test_util.TxtarFilesStdoutStderr(stdout, normalizeIntegrationStderr(stderr))) {
			t.Fatalf("Updated expected file: %s", testPair.ExpectedPath)
		}
		return
	}

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(testPair.ExpectedPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", testPair.ExpectedPath, err)
	}

	result := runTest(testPair.InputPath, expectedStdout, expectedStderr)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != normalizeIntegrationStderr(result.actualStderr)

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(
			normalizeIntegrationStderr(result.expectedStderr),
			normalizeIntegrationStderr(result.actualStderr),
		))
	}
	t.Errorf("%s", msg.String())
}

func splitStderrDiagnostics(stderr string) []string {
	var diagnostics []string
	for part := range strings.SplitSeq(stderr, "\n\n") {
		diagnostic := strings.TrimSpace(part)
		if diagnostic != "" {
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}

func normalizeIntegrationStderr(stderr string) string {
	stderr = strings.TrimSpace(stderr)
	if stderr == "" {
		return ""
	}

	diagnostics := splitStderrDiagnostics(stderr)

	slices.Sort(diagnostics)
	return strings.Join(diagnostics, "\n\n") + "\n"
}

func isTestSkipped(tc test_util.TestCase) bool {
	return slices.Contains(skipIntegrationTests, filepath.ToSlash(tc.Name))
}

func runTest(balFile string, expectedStdout, expectedStderr string) testResult {
	actualStdout, actualStderr := runIntegrationCase(balFile)
	return evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr)
}

func runIntegrationCase(balFile string) (stdout, stderr string) {
	var stdoutBuf, stderrBuf bytes.Buffer

	birPkg, compileErr := runCompilePhase(balFile, &stdoutBuf, &stderrBuf)
	if birPkg == nil || compileErr != nil {
		return stdoutBuf.String(), stderrBuf.String()
	}

	runInterpretPhase(birPkg, &stdoutBuf, &stderrBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr string) testResult {
	stderrMatch := expectedStderr == normalizeIntegrationStderr(actualStderr)
	return testResult{
		success:        actualStdout == expectedStdout && stderrMatch,
		expectedStdout: expectedStdout,
		actualStdout:   actualStdout,
		expectedStderr: expectedStderr,
		actualStderr:   actualStderr,
	}
}

func runCompilePhase(balFile string, stdoutBuf, stderrBuf *bytes.Buffer) (pkg *bir.BIRPackage, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			msg = strings.TrimPrefix(msg, panicPrefix)
			fmt.Fprintf(stdoutBuf, "%s%s\n", panicPrefix, msg)
			err = fmt.Errorf("compile panic")
		}
	}()

	fsys := os.DirFS(filepath.Dir(balFile))

	result, err := directory.LoadProject(fsys, filepath.Base(balFile))
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, err
	}
	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, stderrBuf, compilation.DiagnosticResult())
	if compilation.DiagnosticResult().HasErrors() {
		return nil, nil
	}

	backend := projects.NewBallerinaBackend(compilation)
	return backend.BIR(), nil
}

func runInterpretPhase(birPkg *bir.BIRPackage, stdoutBuf, stderrBuf *bytes.Buffer) {
	if birPkg == nil {
		return
	}
	rt := runtime.NewRuntime()
	runtime.RegisterExternFunction(rt, externOrgName, externModuleName, externFuncName, capturePrintlnOutput(stdoutBuf))
	if err := rt.Interpret(*birPkg); err != nil {
		// For now just write the error string to stderr to match corpus expectations
		fmt.Fprintln(stderrBuf, err.Error())
	}
}

func capturePrintlnOutput(stdoutBuf *bytes.Buffer) func(args []values.BalValue) (values.BalValue, error) {
	return func(args []values.BalValue) (values.BalValue, error) {
		var b strings.Builder
		visited := make(map[uintptr]bool)
		for _, arg := range args {
			b.WriteString(values.String(arg, visited))
		}
		b.WriteByte('\n')
		stdoutBuf.WriteString(b.String())

		return nil, nil
	}
}

func findProjectDirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, "-v") || strings.HasSuffix(name, "-e") || strings.HasSuffix(name, "-p") {
			dirs = append(dirs, filepath.Join(dir, name))
		}
	}
	return dirs
}

func testProjectIntegration(t *testing.T, dirName, projDir, txtarPath string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", dirName, r)
		}
	}()

	if *update {
		stdout, stderr := runProjectIntegrationCase(projDir)
		if test_util.UpdateTxtarArchiveIfNeeded(t, txtarPath, test_util.TxtarFilesStdoutStderr(stdout, normalizeIntegrationStderr(stderr))) {
			t.Fatalf("Updated expected file: %s", txtarPath)
		}
		return
	}

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(txtarPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", txtarPath, err)
	}

	stdout, stderr := runProjectIntegrationCase(projDir)
	result := evaluateTestResult(expectedStdout, expectedStderr, stdout, stderr)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != result.actualStderr

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(
			result.expectedStdout,
			result.actualStdout,
		))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(
			result.expectedStderr,
			normalizeIntegrationStderr(result.actualStderr),
		))
	}
	t.Errorf("%s", msg.String())
}

func runProjectIntegrationCase(projectDir string) (stdout, stderr string) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	birPkgs, compileErr := runProjectCompilePhase(projectDir, &stdoutBuf, &stderrBuf)
	if birPkgs == nil || compileErr != nil {
		return stdoutBuf.String(), stderrBuf.String()
	}

	runProjectInterpretPhase(birPkgs, &stdoutBuf, &stderrBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func runProjectCompilePhase(projectDir string, stdoutBuf, stderrBuf *bytes.Buffer) (pkgs []*bir.BIRPackage, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			msg = strings.TrimPrefix(msg, panicPrefix)
			fmt.Fprintf(stdoutBuf, "%s%s\n", panicPrefix, msg)
			err = fmt.Errorf("compile panic")
		}
	}()

	fsys := os.DirFS(projectDir)

	result, err := directory.LoadProject(fsys, ".")
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, err
	}
	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, stderrBuf, compilation.DiagnosticResult())
	if compilation.DiagnosticResult().HasErrors() {
		return nil, nil
	}

	backend := projects.NewBallerinaBackend(compilation)
	return backend.BIRPackages(), nil
}

func runProjectInterpretPhase(birPkgs []*bir.BIRPackage, stdoutBuf, stderrBuf *bytes.Buffer) {
	if len(birPkgs) == 0 {
		return
	}
	rt := runtime.NewRuntime()
	runtime.RegisterExternFunction(rt, externOrgName, externModuleName, externFuncName, capturePrintlnOutput(stdoutBuf))
	for _, birPkg := range birPkgs {
		if err := rt.Interpret(*birPkg); err != nil {
			fmt.Fprintln(stderrBuf, err.Error())
			return
		}
	}
}

func TestProjectSerializationRoundtrip(t *testing.T) {
	flag.Parse()

	if _, err := os.Stat(corpusProjectBaseDir); os.IsNotExist(err) {
		return
	}

	projectDirs := findProjectDirs(corpusProjectBaseDir)

	for _, projDir := range projectDirs {
		dirName := filepath.Base(projDir)
		if !strings.HasSuffix(dirName, "-v") {
			continue
		}
		txtarPath := filepath.Join(corpusProjectIntegrationBaseDir, dirName+".txtar")

		t.Run(dirName, func(t *testing.T) {
			t.Parallel()
			testProjectSerializationRoundtrip(t, dirName, projDir, txtarPath)
		})
	}
}

func testProjectSerializationRoundtrip(t *testing.T, dirName, projDir, txtarPath string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", dirName, r)
		}
	}()

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(txtarPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", txtarPath, err)
	}

	stdout, stderr := runProjectSerializationRoundtrip(projDir)
	result := evaluateTestResult(expectedStdout, expectedStderr, stdout, stderr)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != result.actualStderr

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(result.expectedStderr, result.actualStderr))
	}
	t.Errorf("%s", msg.String())
}

func runProjectSerializationRoundtrip(projectDir string) (stdout, stderr string) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "%s\n", err.Error())
		return stdoutBuf.String(), stderrBuf.String()
	}

	fsys := os.DirFS(projectDir)
	result, err := directory.LoadProject(fsys, ".")
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "%s\n", err.Error())
		return stdoutBuf.String(), stderrBuf.String()
	}
	project := result.Project()
	currentPkg := project.CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, &stderrBuf, compilation.DiagnosticResult())
	if compilation.DiagnosticResult().HasErrors() {
		return stdoutBuf.String(), stderrBuf.String()
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkgs := backend.BIRPackages()
	exportedSymbols := backend.ExportedSymbols()

	if len(birPkgs) == 0 {
		return stdoutBuf.String(), stderrBuf.String()
	}

	deps := birPkgs[:len(birPkgs)-1]

	// Step 1: Serialize dep symbols and BIR to byte arrays
	type serializedModule struct {
		symBytes []byte
		birBytes []byte
	}
	serializedDeps := make([]serializedModule, 0, len(deps))

	for _, dep := range deps {
		pkgIdent := semantics.PackageIdentifier{
			OrgName:    dep.PackageID.OrgName.Value(),
			ModuleName: dep.PackageID.PkgName.Value(),
		}
		exported, ok := exportedSymbols[pkgIdent]
		if !ok {
			fmt.Fprintf(&stdoutBuf, "exported symbols not found for %s/%s\n", pkgIdent.OrgName, pkgIdent.ModuleName)
			return stdoutBuf.String(), stderrBuf.String()
		}

		symBytes, err := symbolpool.Marshal(exported, dep.TypeEnv)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "symbol serialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		birBytes, err := bircodec.Marshal(dep)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "BIR serialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		serializedDeps = append(serializedDeps, serializedModule{symBytes: symBytes, birBytes: birBytes})
	}

	// Step 2: Create fresh compiler and deserialize dep symbols + BIR
	freshEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	publicSymbols := make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace)
	deserialized := make([]*bir.BIRPackage, 0, len(birPkgs))

	for i, sd := range serializedDeps {
		exported, err := symbolpool.Unmarshal(freshEnv, sd.symBytes)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "symbol deserialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		dep := deps[i]
		pkgIdent := semantics.PackageIdentifier{
			OrgName:    dep.PackageID.OrgName.Value(),
			ModuleName: dep.PackageID.PkgName.Value(),
		}
		publicSymbols[pkgIdent] = exported

		freshCtx := context.NewCompilerContext(freshEnv)
		deserializedPkg, err := bircodec.Unmarshal(freshCtx, sd.birBytes)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "BIR deserialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		deserialized = append(deserialized, deserializedPkg)
	}

	// Step 3: Recompile the main (default) module from source using deserialized dep symbols
	defaultModule := currentPkg.DefaultModule()
	defaultDesc := defaultModule.Descriptor()
	defaultOrg := defaultDesc.Org().Value()

	mainBirPkg, err := compileModuleFromSource(freshEnv, project, defaultModule, absProjectDir, publicSymbols, defaultOrg)
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "main module recompilation failed: %v\n", err)
		return stdoutBuf.String(), stderrBuf.String()
	}

	deserialized = append(deserialized, mainBirPkg)

	runProjectInterpretPhase(deserialized, &stdoutBuf, &stderrBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func compileModuleFromSource(env *context.CompilerEnvironment, project projects.Project, module *projects.Module,
	absProjectDir string, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace, defaultOrg string,
) (*bir.BIRPackage, error) {
	cx := context.NewCompilerContext(env)

	// Parse all source files in the module
	docIDs := module.DocumentIDs()
	var syntaxTrees []*ast.BLangCompilationUnit
	for _, docID := range docIDs {
		relPath := project.DocumentPath(docID)
		absPath := filepath.Join(absProjectDir, relPath)
		st, err := parser.GetSyntaxTree(cx, absPath)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", relPath, err)
		}
		cu := ast.GetCompilationUnit(cx, st)
		syntaxTrees = append(syntaxTrees, cu)
	}

	// Build package from compilation units
	var pkg *ast.BLangPackage
	if len(syntaxTrees) == 1 {
		pkg = ast.ToPackage(syntaxTrees[0])
	} else {
		pkg = &ast.BLangPackage{}
		for _, cu := range syntaxTrees {
			if pkg.PackageID == nil {
				pkg.PackageID = cu.GetPackageID()
			}
			for _, node := range cu.GetTopLevelNodes() {
				switch n := node.(type) {
				case *ast.BLangImportPackage:
					pkg.Imports = append(pkg.Imports, *n)
				case *ast.BLangConstant:
					pkg.Constants = append(pkg.Constants, *n)
				case *ast.BLangService:
					pkg.Services = append(pkg.Services, *n)
				case *ast.BLangFunction:
					pkg.Functions = append(pkg.Functions, *n)
				case *ast.BLangTypeDefinition:
					pkg.TypeDefinitions = append(pkg.TypeDefinitions, *n)
				case *ast.BLangAnnotation:
					pkg.Annotations = append(pkg.Annotations, *n)
				case *ast.BLangClassDefinition:
					pkg.ClassDefinitions = append(pkg.ClassDefinitions, *n)
				default:
					pkg.TopLevelNodes = append(pkg.TopLevelNodes, node)
				}
			}
		}
	}

	// Set the package ID to match the module descriptor
	desc := module.Descriptor()
	orgName := model.Name(desc.Org().Value())
	moduleName := desc.Name().String()
	nameComps := make([]model.Name, 0)
	for _, part := range strings.Split(moduleName, ".") {
		nameComps = append(nameComps, model.Name(part))
	}
	version := model.Name(desc.Version().String())
	if version == "" {
		version = model.DEFAULT_VERSION
	}
	pkg.PackageID = cx.NewPackageID(orgName, nameComps, version)

	// Run compilation pipeline
	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx), publicSymbols, defaultOrg)
	semantics.ResolveSymbols(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("symbol resolution failed")
	}

	semantics.ResolveTopLevelNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("top-level type resolution failed")
	}

	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("local type resolution failed")
	}

	analyzer := semantics.NewSemanticAnalyzer(cx)
	analyzer.Analyze(pkg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("semantic analysis failed")
	}

	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("CFG creation failed")
	}

	semantics.AnalyzeCFG(cx, pkg, cfg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("CFG analysis failed")
	}

	pkg = desugar.DesugarPackage(cx, pkg, importedSymbols)

	return bir.GenBir(cx, pkg), nil
}
