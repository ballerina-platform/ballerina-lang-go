// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

// Java compatibility tests — each test corresponds to a file from the Java
// TOML parser test suite at
// misc/toml-parser/src/test/resources/
package tomlparser

import (
	"os"
	"testing"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func mustLoad(t *testing.T, file string) *Toml {
	t.Helper()
	doc, err := Read(os.DirFS("testdata"), file)
	if err != nil {
		t.Fatalf("unexpected parse error for %s: %v", file, err)
	}
	return doc
}

func mustLoadWithError(t *testing.T, file string) *Toml {
	t.Helper()
	doc, err := Read(os.DirFS("testdata"), file)
	if err == nil {
		t.Fatalf("expected a parse error for %s, got nil", file)
	}
	return doc
}

func expectStr(t *testing.T, doc *Toml, key, want string) {
	t.Helper()
	got, ok := doc.GetString(key)
	if !ok {
		t.Errorf("GetString(%q): key not found", key)
		return
	}
	if got != want {
		t.Errorf("GetString(%q) = %q, want %q", key, got, want)
	}
}

func expectInt(t *testing.T, doc *Toml, key string, want int64) {
	t.Helper()
	got, ok := doc.GetInt(key)
	if !ok {
		t.Errorf("GetInt(%q): key not found", key)
		return
	}
	if got != want {
		t.Errorf("GetInt(%q) = %d, want %d", key, got, want)
	}
}

func expectFloat(t *testing.T, doc *Toml, key string, want float64) {
	t.Helper()
	got, ok := doc.GetFloat(key)
	if !ok {
		t.Errorf("GetFloat(%q): key not found", key)
		return
	}
	if got != want {
		t.Errorf("GetFloat(%q) = %v, want %v", key, got, want)
	}
}

func expectBool(t *testing.T, doc *Toml, key string, want bool) {
	t.Helper()
	got, ok := doc.GetBool(key)
	if !ok {
		t.Errorf("GetBool(%q): key not found", key)
		return
	}
	if got != want {
		t.Errorf("GetBool(%q) = %v, want %v", key, got, want)
	}
}

func expectDiagCount(t *testing.T, doc *Toml, min int) {
	t.Helper()
	if n := len(doc.Diagnostics()); n < min {
		t.Errorf("expected at least %d diagnostic(s), got %d", min, n)
	}
}

// ── valid syntax tests ────────────────────────────────────────────────────────

// TestJava_SyntaxKeys — keys.toml: bare, underscore, dash, numeric, and
// quoted key formats.
func TestJava_SyntaxKeys(t *testing.T) {
	doc := mustLoad(t, "syntax-keys.toml")
	expectStr(t, doc, "key", "basic key")
	expectStr(t, doc, "underscore_key", "Underscore Key")
	expectStr(t, doc, "dash-key", "Dash Key")
	expectStr(t, doc, "key2", "value") // 'key2' single-quoted key
}

// TestJava_SyntaxValues — values.toml: integers, floats, booleans, strings,
// hex/octal/binary literals, underscored numeric literals.
func TestJava_SyntaxValues(t *testing.T) {
	doc := mustLoad(t, "syntax-values.toml")

	// Strings
	expectStr(t, doc, "key1", "hello")
	expectStr(t, doc, "str1", "The quick brown fox jumps over the lazy dog.")
	expectStr(t, doc, "lit1", `C:\Users\nodejs\templates`)
	expectStr(t, doc, "lit3", `Tom "Dubs" Preston-Werner`)

	// Integers — basic
	expectInt(t, doc, "key2", 123)
	expectInt(t, doc, "int1", 99) // +99
	expectInt(t, doc, "int2", 42)
	expectInt(t, doc, "int3", 0)
	expectInt(t, doc, "int4", -17)

	// Integers — underscore separators
	expectInt(t, doc, "int5", 1_000)
	expectInt(t, doc, "int6", 5_349_221)

	// Integers — hex / octal / binary
	expectInt(t, doc, "hex1", 0xDEADBEEF)
	expectInt(t, doc, "hex2", 0xdeadbeef)
	expectInt(t, doc, "oct1", 0o01234567)
	expectInt(t, doc, "oct2", 0o755)
	expectInt(t, doc, "bin1", 0b11010110)

	// Floats
	expectFloat(t, doc, "key3", 56.55)
	expectFloat(t, doc, "flt2", 3.1415)
	expectFloat(t, doc, "flt3", -0.01)
	expectFloat(t, doc, "flt7", 6.626e-34)

	// Booleans
	expectBool(t, doc, "key6", false)
	expectBool(t, doc, "key7", true)
}

// TestJava_SyntaxDottedKeys — dotted.toml: dotted key paths, quoted key
// segments, nested tables created from dotted keys.
// The file intentionally contains `"key" = "quoted key"` which re-defines the
// bare key `key` — the Java parser notes this as a semantic error. The Go
// parser emits a parse-time duplicate-key diagnostic, so we use mustLoadWithError.
func TestJava_SyntaxDottedKeys(t *testing.T) {
	doc := mustLoadWithError(t, "syntax-dotted.toml")

	expectStr(t, doc, "key", "no quote")
	expectStr(t, doc, "apple.type", "fruit")
	expectStr(t, doc, "apple.skin", "thin")
	expectStr(t, doc, "apple.color", "red")
	expectStr(t, doc, "orange.type", "fruit")
	expectStr(t, doc, "orange.skin", "thick")
	expectStr(t, doc, "foo.bar", "test")

	// Under [table] header
	expectStr(t, doc, "table.hi", "hii")
	expectStr(t, doc, "table.hello.child", "new")

	// Dotted keys inside [root] table
	expectStr(t, doc, "root.first.second", "value")
	expectStr(t, doc, "root.first.third", "value1")
}

// TestJava_SyntaxArrays — array.toml: homogeneous, mixed-type, nested arrays.
func TestJava_SyntaxArrays(t *testing.T) {
	doc := mustLoad(t, "syntax-array.toml")

	integers, ok := doc.GetArray("integers")
	if !ok {
		t.Fatal("integers array not found")
	}
	if len(integers) != 3 {
		t.Fatalf("integers: len=%d, want 3", len(integers))
	}
	if v, _ := integers[0].(int64); v != 1 {
		t.Errorf("integers[0] = %v, want 1", integers[0])
	}

	colors, ok := doc.GetArray("colors")
	if !ok {
		t.Fatal("colors array not found")
	}
	wantColors := []string{"red", "yellow", "green"}
	if len(colors) != len(wantColors) {
		t.Fatalf("colors: len=%d, want %d", len(colors), len(wantColors))
	}
	for i, want := range wantColors {
		if got, _ := colors[i].(string); got != want {
			t.Errorf("colors[%d] = %q, want %q", i, got, want)
		}
	}

	nested, ok := doc.GetArray("nested_array")
	if !ok {
		t.Fatal("nested_array not found")
	}
	if len(nested) != 4 {
		t.Fatalf("nested_array: len=%d, want 4", len(nested))
	}
	// nested_array = [1, 2, [5,6], 3] — third element is a nested array
	inner, ok := nested[2].([]any)
	if !ok {
		t.Errorf("nested_array[2] should be []any, got %T", nested[2])
	} else if len(inner) != 2 {
		t.Errorf("nested_array[2]: len=%d, want 2", len(inner))
	}

	same, ok := doc.GetArray("same_type")
	if !ok {
		t.Fatal("same_type array not found")
	}
	if len(same) != 3 {
		t.Errorf("same_type: len=%d, want 3", len(same))
	}
}

// TestJava_SyntaxInlineTables — inline-tables.toml: simple inline table,
// nested inline table, empty inline table, array containing inline tables.
func TestJava_SyntaxInlineTables(t *testing.T) {
	doc := mustLoad(t, "syntax-inline-tables.toml")

	expectStr(t, doc, "rootInline.first", "Tom")
	expectStr(t, doc, "rootInline.last", "Preston-Werner")

	expectStr(t, doc, "nestedInlineTableRoot.another", "Another one")

	// emptyInline = {} — key must exist, value is an empty map
	v, ok := doc.Get("emptyInline")
	if !ok {
		t.Error("emptyInline not found")
	} else if m, ok := v.(map[string]any); !ok {
		t.Errorf("emptyInline: want map[string]any, got %T", v)
	} else if len(m) != 0 {
		t.Errorf("emptyInline: want empty map, got %v", m)
	}

	// arrSample has 3 elements
	arr, ok := doc.GetArray("arrSample")
	if !ok {
		t.Fatal("arrSample not found")
	}
	if len(arr) != 3 {
		t.Errorf("arrSample: len=%d, want 3", len(arr))
	}

	// Under [table] header
	expectStr(t, doc, "table.inlineTable.another", "Another one")
}

// TestJava_SyntaxNoNewlineAtEOF — no-newline-end.toml: file without trailing
// newline still parses correctly.
func TestJava_SyntaxNoNewlineAtEOF(t *testing.T) {
	doc := mustLoad(t, "syntax-no-newline-end.toml")
	expectStr(t, doc, "key", "basic key")
}

// TestJava_SyntaxTable — table.toml: root key-value, subtable declared before
// parent, sibling tables.
func TestJava_SyntaxTable(t *testing.T) {
	doc := mustLoad(t, "syntax-table.toml")

	expectInt(t, doc, "rootKey", 22)
	// [first.sub] is declared before [first]; both must be accessible
	expectStr(t, doc, "first.sub.hey", "sister")
	expectStr(t, doc, "first.key", "sdsad")
	expectStr(t, doc, "first.key1", "eeww")
	expectStr(t, doc, "second.key", "yay")
	expectStr(t, doc, "second.key1", "test")
}

// TestJava_SyntaxArrayOfTables — array-of-tables.toml: [[...]] headers,
// nested array-of-tables, subtables inside array entries.
func TestJava_SyntaxArrayOfTables(t *testing.T) {
	doc := mustLoad(t, "syntax-array-of-tables.toml")

	// [products] table with a key
	expectStr(t, doc, "products.hello1", "hi")

	// [[products.hello]] — 3 entries
	productHellos, ok := doc.GetTables("products.hello")
	if !ok {
		t.Fatal("products.hello array of tables not found")
	}
	if len(productHellos) != 3 {
		t.Fatalf("products.hello: len=%d, want 3", len(productHellos))
	}
	if name, _ := productHellos[0].GetString("name"); name != "Hammer" {
		t.Errorf("products.hello[0].name = %q, want \"Hammer\"", name)
	}
	if name, _ := productHellos[2].GetString("name"); name != "Nail" {
		t.Errorf("products.hello[2].name = %q, want \"Nail\"", name)
	}

	// [foo] defined after [[foo.bar]]
	expectStr(t, doc, "foo.name", "Bob")

	// [[fruits]] — 2 top-level entries
	fruits, ok := doc.GetTables("fruits")
	if !ok {
		t.Fatal("fruits array of tables not found")
	}
	if len(fruits) != 2 {
		t.Fatalf("fruits: len=%d, want 2", len(fruits))
	}
	if name, _ := fruits[0].GetString("name"); name != "apple" {
		t.Errorf("fruits[0].name = %q, want \"apple\"", name)
	}
	if name, _ := fruits[1].GetString("name"); name != "banana" {
		t.Errorf("fruits[1].name = %q, want \"banana\"", name)
	}
}

// TestJava_ObjectComplex — object/complex.toml: mixed root KVs, deeply nested
// tables ([table.child.grandchild]), and array of tables with varying entry sizes.
func TestJava_ObjectComplex(t *testing.T) {
	doc := mustLoad(t, "object-complex.toml")

	expectStr(t, doc, "simplekv", "simplekv value")
	expectStr(t, doc, "simplekv1", "simplekv1 value")
	expectInt(t, doc, "simpleint", 11)
	expectFloat(t, doc, "simplefloat", 1.5)
	expectBool(t, doc, "simplebool", false)

	arr, ok := doc.GetArray("simpleArr")
	if !ok {
		t.Fatal("simpleArr not found")
	}
	if len(arr) != 3 {
		t.Errorf("simpleArr: len=%d, want 3", len(arr))
	}

	// Nested tables
	expectStr(t, doc, "table.tableKv", "table kv")
	expectStr(t, doc, "table.child.tableKvChild", "table kv child")
	expectStr(t, doc, "table.child.grandchild.tableKvGrandChild", "table kv grandchild")
	expectStr(t, doc, "table.child.grandchild.tableKv1GrandChild", "table kv1 grandchild")

	// Array of tables — 3 entries
	tableArr, ok := doc.GetTables("tableArr")
	if !ok {
		t.Fatal("tableArr not found")
	}
	if len(tableArr) != 3 {
		t.Fatalf("tableArr: len=%d, want 3", len(tableArr))
	}
	if v, _ := tableArr[0].GetString("tableKv"); v != "tableArr kv first" {
		t.Errorf("tableArr[0].tableKv = %q, want \"tableArr kv first\"", v)
	}
	if v, _ := tableArr[2].GetString("tableKv1"); v != "tableArr kv1 third" {
		t.Errorf("tableArr[2].tableKv1 = %q, want \"tableArr kv1 third\"", v)
	}
}

// TestJava_ModifierDependencies — modifier/Dependencies.toml: [[dependency]]
// array of tables, as used in the Ballerina module modifier.
func TestJava_ModifierDependencies(t *testing.T) {
	doc := mustLoad(t, "modifier-dependencies.toml")

	deps, ok := doc.GetTables("dependency")
	if !ok {
		t.Fatal("dependency array of tables not found")
	}
	if len(deps) != 2 {
		t.Fatalf("dependency: len=%d, want 2", len(deps))
	}

	if org, _ := deps[0].GetString("org"); org != "wso2" {
		t.Errorf("dependency[0].org = %q, want \"wso2\"", org)
	}
	if name, _ := deps[0].GetString("name"); name != "twitter" {
		t.Errorf("dependency[0].name = %q, want \"twitter\"", name)
	}
	if ver, _ := deps[0].GetString("version"); ver != "2.3.4" {
		t.Errorf("dependency[0].version = %q, want \"2.3.4\"", ver)
	}

	if name, _ := deps[1].GetString("name"); name != "github" {
		t.Errorf("dependency[1].name = %q, want \"github\"", name)
	}
	if ver, _ := deps[1].GetString("version"); ver != "1.2.3" {
		t.Errorf("dependency[1].version = %q, want \"1.2.3\"", ver)
	}
}

// ── negative / error tests ────────────────────────────────────────────────────

// TestJava_NegMissingEqual — missing-equal-negative.toml: `key "value"` has
// no '=' and must produce at least one diagnostic.
func TestJava_NegMissingEqual(t *testing.T) {
	doc := mustLoadWithError(t, "neg-missing-equal.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegMissingValue — missing-value-negative.toml: `key =` with no
// value must produce at least one diagnostic.
func TestJava_NegMissingValue(t *testing.T) {
	doc := mustLoadWithError(t, "neg-missing-value.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegMissingKey — missing-key-negative.toml: `="value"` with no key
// must produce at least one diagnostic.
func TestJava_NegMissingKey(t *testing.T) {
	doc := mustLoadWithError(t, "neg-missing-key.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegDuplicateKey — neg-duplicate-key.toml: key1 defined twice must
// produce at least one diagnostic, while key = "hello" survives.
func TestJava_NegDuplicateKey(t *testing.T) {
	doc := mustLoadWithError(t, "neg-duplicate-key.toml")
	expectDiagCount(t, doc, 1)
	// Non-duplicate key still accessible
	expectStr(t, doc, "key", "hello")
}

// TestJava_NegMultiErrors — neg-multi-errors.toml: three error lines
// (key"value", ="value", key3 =) interspersed with valid lines; valid keys
// key1 and key2 must survive error recovery.
func TestJava_NegMultiErrors(t *testing.T) {
	doc := mustLoadWithError(t, "neg-multi-errors.toml")
	expectDiagCount(t, doc, 3)
	expectInt(t, doc, "key1", 12)
	expectInt(t, doc, "key2", 33)
}

// TestJava_NegInlineMissingComma — neg-inline-missing-comma.toml: inline
// tables with missing comma, double comma, and unclosed brace.
func TestJava_NegInlineMissingComma(t *testing.T) {
	doc := mustLoadWithError(t, "neg-inline-missing-comma.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegEmptyTableKey — neg-empty-table-key.toml: `[]` with an empty
// table key must produce at least one diagnostic.
func TestJava_NegEmptyTableKey(t *testing.T) {
	doc := mustLoadWithError(t, "neg-empty-table-key.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegTableKeyConflict — neg-table-key-conflict.toml: `[table.key1]`
// redefines key1 that was already a string value, must produce a diagnostic.
func TestJava_NegTableKeyConflict(t *testing.T) {
	doc := mustLoadWithError(t, "neg-table-key-conflict.toml")
	expectDiagCount(t, doc, 1)
}

// TestJava_NegArrayMissingComma — neg-array-missing-comma.toml: `[1 2]` with
// no comma separator must produce at least one diagnostic.
func TestJava_NegArrayMissingComma(t *testing.T) {
	doc := mustLoadWithError(t, "neg-array-missing-comma.toml")
	expectDiagCount(t, doc, 1)
}
