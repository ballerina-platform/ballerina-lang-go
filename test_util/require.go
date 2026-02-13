/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package test_util

import (
	"reflect"
	"testing"
)

// Require provides fatal assertion methods that stop test execution on failure.
// Usage: require := testutil.NewRequire(t)
type Require struct {
	t *testing.T
}

// NewRequire creates a new Require instance for the given test.
func NewRequire(t *testing.T) *Require {
	t.Helper()
	return &Require{t: t}
}

// True requires that the condition is true, fails the test immediately if false.
func (r *Require) True(condition bool, msgAndArgs ...interface{}) {
	r.t.Helper()
	if !condition {
		r.failNow("expected true but got false", msgAndArgs...)
	}
}

// False requires that the condition is false, fails the test immediately if true.
func (r *Require) False(condition bool, msgAndArgs ...interface{}) {
	r.t.Helper()
	if condition {
		r.failNow("expected false but got true", msgAndArgs...)
	}
}

// Nil requires that the value is nil, fails the test immediately if not.
func (r *Require) Nil(value interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if !isNil(value) {
		r.failNow("expected nil but got non-nil value", msgAndArgs...)
	}
}

// NotNil requires that the value is not nil, fails the test immediately if nil.
func (r *Require) NotNil(value interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if isNil(value) {
		r.failNow("expected non-nil but got nil", msgAndArgs...)
	}
}

// Equal requires that two values are equal, fails the test immediately if not.
func (r *Require) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		r.t.Fatalf("expected %v but got %v", expected, actual)
	}
}

// NotEqual requires that two values are not equal, fails the test immediately if equal.
func (r *Require) NotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		r.failNow("expected values to be different but they are equal", msgAndArgs...)
	}
}

// Same requires that two pointers refer to the same object.
// Note: Panics are prevented by checking comparability first.
func (r *Require) Same(expected, actual interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if !isComparable(expected) || !isComparable(actual) {
		r.failNow("Same() requires comparable types (not slices, maps, or functions)", msgAndArgs...)
		return
	}
	if expected != actual {
		r.failNow("expected same instance but got different instances", msgAndArgs...)
	}
}

// NotSame requires that two pointers refer to different objects.
// Note: Panics are prevented by checking comparability first.
func (r *Require) NotSame(expected, actual interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	if !isComparable(expected) || !isComparable(actual) {
		r.failNow("NotSame() requires comparable types (not slices, maps, or functions)", msgAndArgs...)
		return
	}
	if expected == actual {
		r.failNow("expected different instances but got same instance", msgAndArgs...)
	}
}

// Len requires that the slice/map/string has the expected length.
func (r *Require) Len(object interface{}, expected int, msgAndArgs ...interface{}) {
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

// Empty requires that the slice/map/string is empty.
func (r *Require) Empty(object interface{}, msgAndArgs ...interface{}) {
	r.t.Helper()
	r.Len(object, 0, msgAndArgs...)
}

// NotEmpty requires that the slice/map/string is not empty.
func (r *Require) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
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

// NoError requires that err is nil, fails the test immediately if not.
func (r *Require) NoError(err error, msgAndArgs ...interface{}) {
	r.t.Helper()
	if err != nil {
		r.t.Fatalf("expected no error but got: %v", err)
	}
}

// Error requires that err is not nil, fails the test immediately if nil.
func (r *Require) Error(err error, msgAndArgs ...interface{}) {
	r.t.Helper()
	if err == nil {
		r.failNow("expected an error but got nil", msgAndArgs...)
	}
}

// Fail fails the test immediately.
func (r *Require) Fail(msgAndArgs ...interface{}) {
	r.t.Helper()
	r.failNow("test failed", msgAndArgs...)
}

// failNow is a helper to report fatal test failures.
func (r *Require) failNow(defaultMsg string, msgAndArgs ...interface{}) {
	r.t.Helper()
	if len(msgAndArgs) > 0 {
		r.t.Fatal(formatMessage(msgAndArgs...))
	} else {
		r.t.Fatal(defaultMsg)
	}
}
