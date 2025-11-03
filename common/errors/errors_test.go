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

package errors

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestDecorateWithStackTrace_WithBacktraceDisabled(t *testing.T) {
	// Ensure BAL_BACKTRACE is not set
	os.Unsetenv("BAL_BACKTRACE")

	err := errors.New("test error")
	decorated := DecorateWithStackTrace(err)

	// When backtrace is disabled, should return the original error
	if decorated.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", decorated.Error())
	}

	// Should not be an ErrorWithStackTrace
	_, ok := decorated.(*ErrorWithStackTrace)
	if ok {
		t.Error("Expected original error, got ErrorWithStackTrace when backtrace is disabled")
	}
}

func TestDecorateWithStackTrace_WithBacktraceEnabled(t *testing.T) {
	// Enable backtrace
	os.Setenv("BAL_BACKTRACE", "true")
	defer os.Unsetenv("BAL_BACKTRACE")

	err := errors.New("test error with stack")
	decorated := DecorateWithStackTrace(err)

	// Should be an ErrorWithStackTrace
	stackErr, ok := decorated.(*ErrorWithStackTrace)
	if !ok {
		t.Fatal("Expected ErrorWithStackTrace when backtrace is enabled")
	}

	// Error message should contain original error
	if !strings.Contains(stackErr.Error(), "test error with stack") {
		t.Errorf("Error message should contain original error text")
	}

	// Stack trace should be present
	if len(stackErr.stack) == 0 {
		t.Error("Expected non-empty stack trace")
	}

	// Error message should contain stack trace
	errorMsg := stackErr.Error()
	if !strings.Contains(errorMsg, "\n\tat ") {
		t.Error("Error message should contain formatted stack trace")
	}
}

func TestDecorateWithStackTrace_WithNilError(t *testing.T) {
	os.Setenv("BAL_BACKTRACE", "true")
	defer os.Unsetenv("BAL_BACKTRACE")

	decorated := DecorateWithStackTrace(nil)

	if decorated != nil {
		t.Error("Expected nil when decorating nil error")
	}
}

func TestErrorWithStackTrace_StackTrace(t *testing.T) {
	os.Setenv("BAL_BACKTRACE", "true")
	defer os.Unsetenv("BAL_BACKTRACE")

	err := errors.New("test error")
	decorated := DecorateWithStackTrace(err).(*ErrorWithStackTrace)

	stackTrace := decorated.StackTrace()

	if stackTrace == nil {
		t.Error("Expected non-nil stack trace")
	}

	stackTraceStr := string(stackTrace)

	// Should contain function name
	if !strings.Contains(stackTraceStr, "\n\tat ") {
		t.Error("Stack trace should contain function frames")
	}

	// Should contain file path
	if !strings.Contains(stackTraceStr, ".go:") {
		t.Error("Stack trace should contain file paths and line numbers")
	}
}
