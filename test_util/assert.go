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
	"fmt"
	"reflect"
	"testing"
)

// Assert provides assertion methods for testing.
type Assert struct {
	t *testing.T
}

// New creates a new Assert instance for the given test.
func New(t *testing.T) *Assert {
	t.Helper()
	return &Assert{t: t}
}

// True asserts that the condition is true.
func (a *Assert) True(condition bool, msgAndArgs ...any) {
	a.t.Helper()
	if !condition {
		a.fail("expected true but got false", msgAndArgs...)
	}
}

// False asserts that the condition is false.
func (a *Assert) False(condition bool, msgAndArgs ...any) {
	a.t.Helper()
	if condition {
		a.fail("expected false but got true", msgAndArgs...)
	}
}

// NotNil asserts that the value is not nil.
func (a *Assert) NotNil(value any, msgAndArgs ...any) {
	a.t.Helper()
	if isNil(value) {
		a.fail("expected non-nil but got nil", msgAndArgs...)
	}
}

// Equal asserts that two values are equal.
func (a *Assert) Equal(expected, actual any, msgAndArgs ...any) {
	a.t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("expected %v but got %v", expected, actual)
		if len(msgAndArgs) > 0 {
			a.t.Log(formatMessage(msgAndArgs...))
		}
	}
}

// NotEqual asserts that two values are not equal.
func (a *Assert) NotEqual(expected, actual any, msgAndArgs ...any) {
	a.t.Helper()
	if reflect.DeepEqual(expected, actual) {
		a.fail("expected values to be different but they are equal", msgAndArgs...)
	}
}

// Same asserts that two pointers refer to the same object.
func (a *Assert) Same(expected, actual any, msgAndArgs ...any) {
	a.t.Helper()
	if !isComparable(expected) || !isComparable(actual) {
		a.fail("Same() requires comparable types (not slices, maps, or functions)", msgAndArgs...)
		return
	}
	if expected != actual {
		a.fail("expected same instance but got different instances", msgAndArgs...)
	}
}

// NotSame asserts that two pointers refer to different objects.
func (a *Assert) NotSame(expected, actual any, msgAndArgs ...any) {
	a.t.Helper()
	if !isComparable(expected) || !isComparable(actual) {
		a.fail("NotSame() requires comparable types (not slices, maps, or functions)", msgAndArgs...)
		return
	}
	if expected == actual {
		a.fail("expected different instances but got same instance", msgAndArgs...)
	}
}

// Len asserts that the slice/map/string has the expected length.
func (a *Assert) Len(object any, expected int, msgAndArgs ...any) {
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

// NotEmpty asserts that the slice/map/string is not empty.
func (a *Assert) NotEmpty(object any, msgAndArgs ...any) {
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

func (a *Assert) fail(defaultMsg string, msgAndArgs ...any) {
	a.t.Helper()
	if len(msgAndArgs) > 0 {
		a.t.Error(formatMessage(msgAndArgs...))
	} else {
		a.t.Error(defaultMsg)
	}
}

func formatMessage(msgAndArgs ...any) string {
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
		return fmt.Sprintf(format, msgAndArgs[1:]...)
	}
	return ""
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	}
	return false
}

func isComparable(value any) bool {
	if value == nil {
		return true
	}
	return reflect.TypeOf(value).Comparable()
}
