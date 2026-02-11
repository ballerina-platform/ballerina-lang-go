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

// Package test_util provides assertion utilities for testing.
// Java: org.testng.Assert equivalent
package test_util

import (
	"reflect"
	"testing"
)

// Assert provides assertion methods for testing.
// Usage: assert := testutil.New(t)
type Assert struct {
	t *testing.T
}

// New creates a new Assert instance for the given test.
func New(t *testing.T) *Assert {
	t.Helper()
	return &Assert{t: t}
}

// True asserts that the condition is true.
// Java: Assert.assertTrue(condition)
func (a *Assert) True(condition bool, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !condition {
		a.fail("expected true but got false", msgAndArgs...)
	}
}

// False asserts that the condition is false.
// Java: Assert.assertFalse(condition)
func (a *Assert) False(condition bool, msgAndArgs ...interface{}) {
	a.t.Helper()
	if condition {
		a.fail("expected false but got true", msgAndArgs...)
	}
}

// Nil asserts that the value is nil.
// Java: Assert.assertNull(value)
func (a *Assert) Nil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !isNil(value) {
		a.fail("expected nil but got non-nil value", msgAndArgs...)
	}
}

// NotNil asserts that the value is not nil.
// Java: Assert.assertNotNull(value)
func (a *Assert) NotNil(value interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if isNil(value) {
		a.fail("expected non-nil but got nil", msgAndArgs...)
	}
}

// Equal asserts that two values are equal.
// Java: Assert.assertEquals(actual, expected)
func (a *Assert) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("expected %v but got %v", expected, actual)
		if len(msgAndArgs) > 0 {
			a.t.Log(formatMessage(msgAndArgs...))
		}
	}
}

// NotEqual asserts that two values are not equal.
// Java: Assert.assertNotEquals(actual, expected)
func (a *Assert) NotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		a.fail("expected values to be different but they are equal", msgAndArgs...)
	}
}

// Same asserts that two pointers refer to the same object.
// Java: Assert.assertSame(actual, expected)
func (a *Assert) Same(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if expected != actual {
		a.fail("expected same instance but got different instances", msgAndArgs...)
	}
}

// NotSame asserts that two pointers refer to different objects.
// Java: Assert.assertNotSame(actual, expected)
func (a *Assert) NotSame(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if expected == actual {
		a.fail("expected different instances but got same instance", msgAndArgs...)
	}
}

// Contains asserts that the string contains the substring.
func (a *Assert) Contains(s, substr string, msgAndArgs ...interface{}) {
	a.t.Helper()
	if len(s) == 0 || len(substr) == 0 {
		a.fail("expected string to contain substring", msgAndArgs...)
		return
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return
		}
	}
	a.t.Errorf("expected %q to contain %q", s, substr)
	if len(msgAndArgs) > 0 {
		a.t.Log(formatMessage(msgAndArgs...))
	}
}

// Len asserts that the slice/map/string has the expected length.
func (a *Assert) Len(object interface{}, expected int, msgAndArgs ...interface{}) {
	a.t.Helper()
	v := reflect.ValueOf(object)
	var actual int
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		actual = v.Len()
	default:
		a.fail("cannot get length of non-collection type", msgAndArgs...)
		return
	}
	if actual != expected {
		a.t.Errorf("expected length %d but got %d", expected, actual)
		if len(msgAndArgs) > 0 {
			a.t.Log(formatMessage(msgAndArgs...))
		}
	}
}

// Empty asserts that the slice/map/string is empty.
func (a *Assert) Empty(object interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	a.Len(object, 0, msgAndArgs...)
}

// NotEmpty asserts that the slice/map/string is not empty.
func (a *Assert) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	v := reflect.ValueOf(object)
	var length int
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		length = v.Len()
	default:
		a.fail("cannot get length of non-collection type", msgAndArgs...)
		return
	}
	if length == 0 {
		a.fail("expected non-empty but got empty", msgAndArgs...)
	}
}

// NoError asserts that err is nil.
func (a *Assert) NoError(err error, msgAndArgs ...interface{}) {
	a.t.Helper()
	if err != nil {
		a.t.Errorf("expected no error but got: %v", err)
		if len(msgAndArgs) > 0 {
			a.t.Log(formatMessage(msgAndArgs...))
		}
	}
}

// Error asserts that err is not nil.
func (a *Assert) Error(err error, msgAndArgs ...interface{}) {
	a.t.Helper()
	if err == nil {
		a.fail("expected an error but got nil", msgAndArgs...)
	}
}

// Fail fails the test immediately.
func (a *Assert) Fail(msgAndArgs ...interface{}) {
	a.t.Helper()
	a.fail("test failed", msgAndArgs...)
}

// FailNow fails the test immediately and stops execution.
func (a *Assert) FailNow(msgAndArgs ...interface{}) {
	a.t.Helper()
	if len(msgAndArgs) > 0 {
		a.t.Fatal(formatMessage(msgAndArgs...))
	} else {
		a.t.FailNow()
	}
}

// fail is a helper to report test failures.
func (a *Assert) fail(defaultMsg string, msgAndArgs ...interface{}) {
	a.t.Helper()
	if len(msgAndArgs) > 0 {
		a.t.Error(formatMessage(msgAndArgs...))
	} else {
		a.t.Error(defaultMsg)
	}
}

// formatMessage formats the message and arguments.
func formatMessage(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if msg, ok := msgAndArgs[0].(string); ok {
			return msg
		}
		return ""
	}
	if format, ok := msgAndArgs[0].(string); ok {
		return format // Just return format, let caller handle args
	}
	return ""
}

// isNil checks if a value is nil, handling interface nil correctly.
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	}
	return false
}
