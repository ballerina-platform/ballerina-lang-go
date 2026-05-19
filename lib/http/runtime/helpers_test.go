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

// Internal test file (package http) so private helpers are directly accessible.

package http

import (
	"encoding/json"
	"testing"

	"ballerina-lang-go/decimal"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

// getParams is a test helper that extracts the "params" field from a header entry map.
func getParams(entry *values.Map) *values.Map {
	v, _ := entry.Get("params")
	return v.(*values.Map)
}

// ---------------------------------------------------------------------------
// splitOutsideQuotes
// ---------------------------------------------------------------------------

func TestSplitOutsideQuotes_Empty(t *testing.T) {
	got := splitOutsideQuotes("", ',')
	if len(got) != 1 || got[0] != "" {
		t.Errorf("expected [\"\"], got %v", got)
	}
}

func TestSplitOutsideQuotes_NoSeparator(t *testing.T) {
	got := splitOutsideQuotes("text/html", ',')
	if len(got) != 1 || got[0] != "text/html" {
		t.Errorf("expected [\"text/html\"], got %v", got)
	}
}

func TestSplitOutsideQuotes_SimpleSplit(t *testing.T) {
	got := splitOutsideQuotes("a,b,c", ',')
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("got[%d] = %q, want %q", i, got[i], w)
		}
	}
}

func TestSplitOutsideQuotes_SeparatorInsideQuotes(t *testing.T) {
	// The comma inside "a,b" should not split.
	got := splitOutsideQuotes(`"a,b",c`, ',')
	want := []string{`"a,b"`, "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("got[%d] = %q, want %q", i, got[i], w)
		}
	}
}

func TestSplitOutsideQuotes_EscapedQuoteInsideQuotes(t *testing.T) {
	// \"  inside a quoted string is an escaped quote, not a quote boundary.
	got := splitOutsideQuotes(`"a\"b",c`, ',')
	want := []string{`"a\"b"`, "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("got[%d] = %q, want %q", i, got[i], w)
		}
	}
}

// ---------------------------------------------------------------------------
// parseHeader
// ---------------------------------------------------------------------------

func TestParseHeader_SimpleValue(t *testing.T) {
	list, err := parseHeader("text/html")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", list.Len())
	}
	entry := list.Get(0).(*values.Map)
	if v, _ := entry.Get("value"); v != "text/html" {
		t.Errorf("expected value \"text/html\", got %v", v)
	}
}

func TestParseHeader_WithParams(t *testing.T) {
	list, err := parseHeader("text/html; charset=utf-8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", list.Len())
	}
	entry := list.Get(0).(*values.Map)
	if v, _ := entry.Get("value"); v != "text/html" {
		t.Errorf("expected value \"text/html\", got %v", v)
	}
	params := getParams(entry)
	if v, _ := params.Get("charset"); v != "utf-8" {
		t.Errorf("expected charset=utf-8, got %v", v)
	}
}

func TestParseHeader_QuotedParamValue(t *testing.T) {
	list, err := parseHeader(`multipart/form-data; boundary="----boundary"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := list.Get(0).(*values.Map)
	params := getParams(entry)
	if v, _ := params.Get("boundary"); v != "----boundary" {
		t.Errorf("expected boundary stripped of quotes, got %v", v)
	}
}

func TestParseHeader_MultipleValues(t *testing.T) {
	list, err := parseHeader("text/html, application/json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", list.Len())
	}
}

func TestParseHeader_EmptySegment(t *testing.T) {
	_, err := parseHeader("text/html,,application/json")
	if err == nil {
		t.Fatal("expected error for empty segment, got nil")
	}
}

func TestParseHeader_MissingValueBeforeParams(t *testing.T) {
	_, err := parseHeader("; charset=utf-8")
	if err == nil {
		t.Fatal("expected error for missing value before parameters, got nil")
	}
}

func TestParseHeader_ParamWithoutValue(t *testing.T) {
	list, err := parseHeader("multipart/form-data; boundary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := list.Get(0).(*values.Map)
	params := getParams(entry)
	if v, _ := params.Get("boundary"); v != "" {
		t.Errorf("expected empty string for valueless param, got %v", v)
	}
}

// ---------------------------------------------------------------------------
// listToBytes
// ---------------------------------------------------------------------------

func TestListToBytes_ValidBytes(t *testing.T) {
	list := values.NewList(3, semtypes.LIST, nil)
	list.FillingSet(0, int64(72))  // 'H'
	list.FillingSet(1, int64(101)) // 'e'
	list.FillingSet(2, int64(108)) // 'l'
	b, ok := listToBytes(list)
	if !ok {
		t.Fatal("expected ok=true for valid byte list")
	}
	if string(b) != "Hel" {
		t.Errorf("expected \"Hel\", got %q", string(b))
	}
}

func TestListToBytes_EmptyList(t *testing.T) {
	list := values.NewList(0, semtypes.LIST, nil)
	b, ok := listToBytes(list)
	if !ok {
		t.Fatal("expected ok=true for empty list")
	}
	if len(b) != 0 {
		t.Errorf("expected empty bytes, got %v", b)
	}
}

func TestListToBytes_OutOfRange(t *testing.T) {
	list := values.NewList(2, semtypes.LIST, nil)
	list.FillingSet(0, int64(100))
	list.FillingSet(1, int64(300)) // > 255
	_, ok := listToBytes(list)
	if ok {
		t.Fatal("expected ok=false for out-of-range value")
	}
}

func TestListToBytes_NegativeValue(t *testing.T) {
	list := values.NewList(1, semtypes.LIST, nil)
	list.FillingSet(0, int64(-1))
	_, ok := listToBytes(list)
	if ok {
		t.Fatal("expected ok=false for negative value")
	}
}

func TestListToBytes_NonIntegerValue(t *testing.T) {
	list := values.NewList(1, semtypes.LIST, nil)
	list.FillingSet(0, "not-an-int")
	_, ok := listToBytes(list)
	if ok {
		t.Fatal("expected ok=false for non-integer value")
	}
}

// ---------------------------------------------------------------------------
// goToBalValue
// ---------------------------------------------------------------------------

func TestGoToBalValue_Nil(t *testing.T) {
	if got := goToBalValue(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestGoToBalValue_Bool(t *testing.T) {
	if got := goToBalValue(true); got != true {
		t.Errorf("expected true, got %v", got)
	}
	if got := goToBalValue(false); got != false {
		t.Errorf("expected false, got %v", got)
	}
}

func TestGoToBalValue_JsonNumberInt(t *testing.T) {
	n := json.Number("42")
	got := goToBalValue(n)
	if got != int64(42) {
		t.Errorf("expected int64(42), got %v (%T)", got, got)
	}
}

func TestGoToBalValue_JsonNumberFloat(t *testing.T) {
	n := json.Number("3.14")
	got := goToBalValue(n)
	f, ok := got.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", got)
	}
	if f < 3.13 || f > 3.15 {
		t.Errorf("expected ~3.14, got %v", f)
	}
}

func TestGoToBalValue_String(t *testing.T) {
	if got := goToBalValue("hello"); got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

func TestGoToBalValue_SliceOfInterface(t *testing.T) {
	input := []interface{}{"a", json.Number("1")}
	got := goToBalValue(input)
	list, ok := got.(*values.List)
	if !ok {
		t.Fatalf("expected *values.List, got %T", got)
	}
	if list.Len() != 2 {
		t.Errorf("expected length 2, got %d", list.Len())
	}
}

func TestGoToBalValue_MapOfStringInterface(t *testing.T) {
	input := map[string]interface{}{"key": "value"}
	got := goToBalValue(input)
	m, ok := got.(*values.Map)
	if !ok {
		t.Fatalf("expected *values.Map, got %T", got)
	}
	if v, _ := m.Get("key"); v != "value" {
		t.Errorf("expected \"value\", got %v", v)
	}
}

func TestGoToBalValue_Unknown(t *testing.T) {
	type custom struct{}
	if got := goToBalValue(custom{}); got != nil {
		t.Errorf("expected nil for unknown type, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// balToGoJSON
// ---------------------------------------------------------------------------

func TestBalToGoJSON_Nil(t *testing.T) {
	if got := balToGoJSON(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestBalToGoJSON_Bool(t *testing.T) {
	if got := balToGoJSON(true); got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestBalToGoJSON_Int64(t *testing.T) {
	if got := balToGoJSON(int64(7)); got != int64(7) {
		t.Errorf("expected int64(7), got %v", got)
	}
}

func TestBalToGoJSON_Float64(t *testing.T) {
	if got := balToGoJSON(float64(1.5)); got != float64(1.5) {
		t.Errorf("expected float64(1.5), got %v", got)
	}
}

func TestBalToGoJSON_Decimal(t *testing.T) {
	d := decimal.FromInt64(42)
	got := balToGoJSON(d)
	raw, ok := got.(json.RawMessage)
	if !ok {
		t.Fatalf("expected json.RawMessage for decimal, got %T", got)
	}
	if string(raw) != "42" {
		t.Errorf("expected \"42\", got %q", string(raw))
	}
}

func TestBalToGoJSON_String(t *testing.T) {
	if got := balToGoJSON("hello"); got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

func TestBalToGoJSON_Map(t *testing.T) {
	m := values.NewMap(semtypes.MAPPING)
	m.Put("a", int64(1))
	got := balToGoJSON(m)
	goMap, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}
	if goMap["a"] != int64(1) {
		t.Errorf("expected a=1, got %v", goMap["a"])
	}
}

func TestBalToGoJSON_List(t *testing.T) {
	list := values.NewList(2, semtypes.LIST, nil)
	list.FillingSet(0, int64(10))
	list.FillingSet(1, "x")
	got := balToGoJSON(list)
	slice, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	if len(slice) != 2 {
		t.Fatalf("expected length 2, got %d", len(slice))
	}
}

func TestBalToGoJSON_Unknown(t *testing.T) {
	type custom struct{}
	if got := balToGoJSON(custom{}); got != nil {
		t.Errorf("expected nil for unknown type, got %v", got)
	}
}
