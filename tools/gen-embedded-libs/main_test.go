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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/context"
	"ballerina-lang-go/lib/registry"
	"ballerina-lang-go/semtypes"
)

func TestCompileAndWrite_ok(t *testing.T) {
	outDir := t.TempDir()
	if err := compileAndWrite("testdata", "langlib/ok/bal", outDir); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"ballerina.lang.embedtest.platform.sym", "ballerina.lang.embedtest.platform.bir"} {
		info, err := os.Stat(filepath.Join(outDir, name))
		if err != nil || info.Size() == 0 {
			t.Fatalf("%s: %v", name, err)
		}
	}

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	exp, ok, err := registry.LoadSymbols(env, registry.ID{OrgName: "ballerina", ModuleName: "lang.embedtest"})
	if err != nil || !ok {
		t.Fatalf("LoadSymbols: ok=%v err=%v", ok, err)
	}
	if _, ok := exp.GetSymbol("EmbedTest"); !ok {
		t.Fatal("missing EmbedTest export")
	}
}

func TestCompileAndWrite_compileErrors(t *testing.T) {
	err := compileAndWrite("testdata", "langlib/bad/bal", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "langlib/bad/bal: compile errors:") {
		t.Fatalf("got err=%v", err)
	}
}
