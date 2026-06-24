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

package native

import (
	"strings"
	"testing"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testharness"
	"ballerina-lang-go/values"
)

// The log module writes its records to stderr. A corpus -v.bal test cannot
// validate that: the integration harness fatals on any -v test that produces
// non-empty stderr, and there is no corpus category for "valid program that
// writes to stderr". So the level filter, the LOGFMT formatter, and the emit
// path are unit-tested here against a capturing test PAL — they have no corpus
// equivalent.

func TestIsLevelEnabled(t *testing.T) {
	t.Parallel()
	if isLevelEnabled("DEBUG") {
		t.Error("DEBUG should be filtered at the default INFO level")
	}
	for _, level := range []string{"INFO", "WARN", "ERROR"} {
		if !isLevelEnabled(level) {
			t.Errorf("%s should be enabled at the default INFO level", level)
		}
	}
	// Unknown levels are not filtered (default-permit branch).
	if !isLevelEnabled("TRACE") {
		t.Error("unknown level should be enabled")
	}
}

func TestFormatLogfmt(t *testing.T) {
	t.Parallel()
	env := semtypes.CreateTypeEnv()
	ctx := semtypes.ContextFrom(env)

	// Basic record: no error, no key-values.
	got := formatLogfmt("INFO", "2024-01-01T00:00:00.000Z", "hello", nil, nil)
	want := `time=2024-01-01T00:00:00.000Z level=INFO module="" message="hello"`
	if got != want {
		t.Errorf("basic: got %q, want %q", got, want)
	}

	// Special characters in the message are LOGFMT-escaped.
	esc := formatLogfmt("INFO", "T", "a\"b\nc\td", nil, nil)
	if !strings.Contains(esc, `\"`) || !strings.Contains(esc, `\n`) || !strings.Contains(esc, `\t`) {
		t.Errorf("escaping not applied: %q", esc)
	}

	// With an error value.
	errVal := values.NewErrorWithMessage("boom")
	withErr := formatLogfmt("ERROR", "T", "m", errVal, nil)
	if !strings.Contains(withErr, " error=") {
		t.Errorf("error field missing: %q", withErr)
	}

	// With a key-value map holding both a string and a non-string value.
	atomic := semtypes.ToMappingAtomicType(ctx, semtypes.MAPPING)
	kv := values.NewMap(semtypes.MAPPING, atomic, false, []values.MapEntry{
		{Key: "id", Value: int64(42)},
		{Key: "path", Value: "/api"},
	})
	withKv := formatLogfmt("INFO", "T", "m", nil, kv)
	if !strings.Contains(withKv, `path="/api"`) {
		t.Errorf("string kv not quoted: %q", withKv)
	}
	if !strings.Contains(withKv, "id=42") {
		t.Errorf("non-string kv missing: %q", withKv)
	}
}

func TestPrintLog(t *testing.T) {
	t.Parallel()

	// A filtered (DEBUG) record produces no output.
	debugPal := testharness.NewTestPal()
	debugRt := runtime.NewRuntime(debugPal.Platform(), semtypes.CreateTypeEnv())
	if _, err := printLog(debugRt, []values.BalValue{"DEBUG", "skip me", nil, nil}); err != nil {
		t.Fatalf("printLog DEBUG: %v", err)
	}
	if debugPal.Stderr() != "" {
		t.Errorf("DEBUG should not emit, got %q", debugPal.Stderr())
	}

	// An INFO record is formatted and written to stderr.
	infoPal := testharness.NewTestPal()
	infoRt := runtime.NewRuntime(infoPal.Platform(), semtypes.CreateTypeEnv())
	if _, err := printLog(infoRt, []values.BalValue{"INFO", "hello world", nil, nil}); err != nil {
		t.Fatalf("printLog INFO: %v", err)
	}
	out := infoPal.Stderr()
	if !strings.Contains(out, "level=INFO") || !strings.Contains(out, `message="hello world"`) {
		t.Errorf("INFO emit: got %q", out)
	}
}
