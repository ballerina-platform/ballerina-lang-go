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
	"os"
	"path/filepath"
	"testing"

	"ballerina-lang-go/bir"
	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
)

type serializationFixture struct {
	birPkg          *bir.BIRPackage
	exportedSymbols model.ExportedSymbolSpace
	tyEnv           semtypes.Env
	cenv            *context.CompilerEnvironment
}

func compileForSerializationBench(b *testing.B, tc test_util.TestCase) *serializationFixture {
	b.Helper()

	fsys := os.DirFS(filepath.Dir(tc.InputPath))
	entry := filepath.Base(tc.InputPath)
	if tc.IsProject {
		fsys = os.DirFS(tc.InputPath)
		entry = "."
	}

	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		b.Fatalf("getBallerinaEnvPath: %v", err)
	}
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, entry, projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		b.Fatalf("projects.Load(%s): %v", tc.InputPath, err)
	}
	tyEnv := result.Project().Environment().TypeEnv()
	cenv := result.Project().Environment().CompilerEnvironment()
	compilation := result.Project().CurrentPackage().Compilation()

	var stderrBuf bytes.Buffer
	printDiagnostics(fsys, &stderrBuf, compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	if compilation.DiagnosticResult().HasErrors() {
		b.Fatalf("compile errors for %s:\n%s", tc.InputPath, stderrBuf.String())
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()
	if birPkg == nil {
		b.Fatalf("nil BIR for %s", tc.InputPath)
	}

	pkgIdent := semantics.PackageIdentifier{
		OrgName:    birPkg.PackageID.OrgName.Value(),
		ModuleName: birPkg.PackageID.PkgName.Value(),
	}
	exported, ok := backend.ExportedSymbols()[pkgIdent]
	if !ok {
		b.Fatalf("exported symbols not found for %s/%s", pkgIdent.OrgName, pkgIdent.ModuleName)
	}

	return &serializationFixture{birPkg: birPkg, exportedSymbols: exported, tyEnv: tyEnv, cenv: cenv}
}

func benchTestPairs(b *testing.B) []test_util.TestCase {
	return test_util.GetTests(b, test_util.Bench, func(path string) bool { return true })
}

func BenchmarkBIRMarshal(b *testing.B) {
	for _, tc := range benchTestPairs(b) {
		fixture := compileForSerializationBench(b, tc)
		b.Run(tc.Name, func(b *testing.B) {
			for b.Loop() {
				if _, err := bircodec.Marshal(fixture.tyEnv, fixture.birPkg); err != nil {
					b.Fatalf("BIR Marshal: %v", err)
				}
			}
		})
	}
}

func BenchmarkBIRUnmarshal(b *testing.B) {
	for _, tc := range benchTestPairs(b) {
		fixture := compileForSerializationBench(b, tc)
		data, err := bircodec.Marshal(fixture.tyEnv, fixture.birPkg)
		if err != nil {
			b.Fatalf("BIR Marshal setup: %v", err)
		}
		b.Run(tc.Name, func(b *testing.B) {
			for b.Loop() {
				freshEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
				freshCtx := context.NewCompilerContext(freshEnv)
				if _, err := bircodec.Unmarshal(freshCtx, data); err != nil {
					b.Fatalf("BIR Unmarshal: %v", err)
				}
			}
		})
	}
}

func BenchmarkSymbolMarshal(b *testing.B) {
	for _, tc := range benchTestPairs(b) {
		fixture := compileForSerializationBench(b, tc)
		b.Run(tc.Name, func(b *testing.B) {
			for b.Loop() {
				if _, err := symbolpool.Marshal(fixture.exportedSymbols, fixture.cenv); err != nil {
					b.Fatalf("Symbol Marshal: %v", err)
				}
			}
		})
	}
}

func BenchmarkSymbolUnmarshal(b *testing.B) {
	for _, tc := range benchTestPairs(b) {
		fixture := compileForSerializationBench(b, tc)
		data, err := symbolpool.Marshal(fixture.exportedSymbols, fixture.cenv)
		if err != nil {
			b.Fatalf("Symbol Marshal setup: %v", err)
		}
		b.Run(tc.Name, func(b *testing.B) {
			for b.Loop() {
				freshEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
				if _, err := symbolpool.Unmarshal(freshEnv, data); err != nil {
					b.Fatalf("Symbol Unmarshal: %v", err)
				}
			}
		})
	}
}
