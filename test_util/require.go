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

package test_util

import (
	"reflect"
	"testing"
)

// Require provides fatal assertion methods that stop test execution on failure.
type Require struct {
	t *testing.T
}

// NewRequire creates a new Require instance for the given test.
func NewRequire(t *testing.T) *Require {
	t.Helper()
	return &Require{t: t}
}

// NoError requires that err is nil, fails the test immediately if not.
func (r *Require) NoError(err error, msgAndArgs ...any) {
	r.t.Helper()
	if err != nil {
		if len(msgAndArgs) > 0 {
			r.t.Fatalf("%s: %v", formatMessage(msgAndArgs...), err)
		}
		r.t.Fatalf("expected no error but got: %v", err)
	}
}

// NotNil requires that the value is not nil, fails the test immediately if nil.
func (r *Require) NotNil(value any, msgAndArgs ...any) {
	r.t.Helper()
	if isNil(value) {
		r.failNow("expected non-nil but got nil", msgAndArgs...)
	}
}

// NotEqual requires that two values are not equal, fails the test immediately if equal.
func (r *Require) NotEqual(expected, actual any, msgAndArgs ...any) {
	r.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		r.failNow("expected values to be different but they are equal", msgAndArgs...)
	}
}

// Len requires that the slice/map/string has the expected length.
func (r *Require) Len(object any, expected int, msgAndArgs ...any) {
	r.t.Helper()
	v := reflect.ValueOf(object)
	var actual int
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		actual = v.Len()
	default:
		r.failNow("cannot get length of non-collection type", msgAndArgs...)
		return
	}
	if actual != expected {
		r.t.Fatalf("expected length %d but got %d", expected, actual)
	}
}

// NotEmpty requires that the slice/map/string is not empty.
func (r *Require) NotEmpty(object any, msgAndArgs ...any) {
	r.t.Helper()
	v := reflect.ValueOf(object)
	var length int
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		length = v.Len()
	default:
		r.failNow("cannot get length of non-collection type", msgAndArgs...)
		return
	}
	if length == 0 {
		r.failNow("expected non-empty but got empty", msgAndArgs...)
	}
}

func (r *Require) failNow(defaultMsg string, msgAndArgs ...any) {
	r.t.Helper()
	if len(msgAndArgs) > 0 {
		r.t.Fatal(formatMessage(msgAndArgs...))
	} else {
		r.t.Fatal(defaultMsg)
	}
}
