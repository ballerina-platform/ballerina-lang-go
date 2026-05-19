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

package runtime_test

import (
	"io"
	"testing"

	"ballerina-lang-go/model"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/values"
)

func newTestRuntime() *runtime.Runtime {
	return runtime.NewRuntime(test_util.TestPal(io.Discard, io.Discard))
}

// testPkgID creates a PackageID using the default interner. Test org/pkg names
// use the "testorg" prefix to avoid clashing with real package IDs.
func testPkgID(org, pkg string) *model.PackageID {
	return model.NewPackageID(model.DefaultPackageIDInterner, model.Name(org), []model.Name{model.Name(pkg)}, model.Name("0.0.1"))
}

// TestRegisterModuleGlobals_NewModule verifies that calling RegisterModuleGlobals
// for a previously unknown package ID creates a new module with the supplied globals.
func TestRegisterModuleGlobals_NewModule(t *testing.T) {
	rt := newTestRuntime()
	pkgId := testPkgID("test", "newmod")

	runtime.RegisterModuleGlobals(rt, pkgId, map[string]values.BalValue{
		"test/newmod:KEY": "value1",
	})

	globals := runtime.GetModuleGlobalsForTest(rt, pkgId)
	if globals == nil {
		t.Fatal("expected globals map to be non-nil")
	}
	if got := globals["test/newmod:KEY"]; got != "value1" {
		t.Errorf("expected globals[KEY] = %q, got %v", "value1", got)
	}
}

// TestRegisterModuleGlobals_MergesIntoExisting verifies that a second call to
// RegisterModuleGlobals for the same package merges new keys and overwrites
// existing ones (exercises runtime.go L99–100).
func TestRegisterModuleGlobals_MergesIntoExisting(t *testing.T) {
	rt := newTestRuntime()
	pkgId := testPkgID("test", "merge")

	runtime.RegisterModuleGlobals(rt, pkgId, map[string]values.BalValue{
		"test/merge:A": "original",
		"test/merge:B": "keep",
	})
	runtime.RegisterModuleGlobals(rt, pkgId, map[string]values.BalValue{
		"test/merge:A": "overwritten",
		"test/merge:C": "new",
	})

	globals := runtime.GetModuleGlobalsForTest(rt, pkgId)
	if globals == nil {
		t.Fatal("expected globals map to be non-nil after merge")
	}
	if got := globals["test/merge:A"]; got != "overwritten" {
		t.Errorf("expected A to be overwritten, got %v", got)
	}
	if got := globals["test/merge:B"]; got != "keep" {
		t.Errorf("expected B to be kept, got %v", got)
	}
	if got := globals["test/merge:C"]; got != "new" {
		t.Errorf("expected C to be present, got %v", got)
	}
}

// TestRegisterModuleGlobals_NilGlobalsInit verifies that RegisterModuleGlobals
// initialises a nil Globals map when the module already exists without one
// (exercises runtime.go L96–97).
func TestRegisterModuleGlobals_NilGlobalsInit(t *testing.T) {
	rt := newTestRuntime()
	pkgId := testPkgID("test", "nilglobals")

	// Pre-register the module with nil Globals (simulates a module registered
	// before globals are available, e.g. via BIR deserialization).
	runtime.RegisterNilGlobalsModuleForTest(rt, pkgId)

	// Confirm that the module exists but has nil Globals before the call.
	if existing := runtime.GetModuleGlobalsForTest(rt, pkgId); existing != nil {
		t.Fatalf("pre-condition failed: expected nil Globals, got %v", existing)
	}

	runtime.RegisterModuleGlobals(rt, pkgId, map[string]values.BalValue{
		"test/nilglobals:X": int64(42),
	})

	globals := runtime.GetModuleGlobalsForTest(rt, pkgId)
	if globals == nil {
		t.Fatal("expected Globals to be initialised after RegisterModuleGlobals")
	}
	if got := globals["test/nilglobals:X"]; got != int64(42) {
		t.Errorf("expected X = 42, got %v", got)
	}
}
