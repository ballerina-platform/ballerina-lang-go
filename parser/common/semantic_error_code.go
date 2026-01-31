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

package common

import "ballerina-lang-go/tools/diagnostics"

// SemanticErrorCode represents semantic error codes for type checking and semantic analysis.
type SemanticErrorCode struct {
	diagnosticId  string
	messageFormat string
}

func newSemanticErrorCode(diagnosticId, messageFormat string) *SemanticErrorCode {
	return &SemanticErrorCode{
		diagnosticId:  diagnosticId,
		messageFormat: messageFormat,
	}
}

func (s *SemanticErrorCode) DiagnosticId() string {
	return s.diagnosticId
}

func (s *SemanticErrorCode) MessageKey() string {
	return s.messageFormat
}

func (s *SemanticErrorCode) Severity() diagnostics.DiagnosticSeverity {
	return diagnostics.Error
}

// Semantic error codes - matching Ballerina Java compiler (DiagnosticErrorCode.java)
// Reference: https://github.com/ballerina-platform/ballerina-lang/blob/master/compiler/ballerina-lang/src/main/java/org/ballerinalang/util/diagnostic/DiagnosticErrorCode.java

// Type mismatch errors
var INCOMPATIBLE_TYPES = newSemanticErrorCode("BCE2066", "incompatible types: expected '%s', found '%s'")

// Operator errors - matching Java: BINARY_OP_INCOMPATIBLE_TYPES, UNARY_OP_INCOMPATIBLE_TYPES
var BINARY_OP_INCOMPATIBLE_TYPES = newSemanticErrorCode("BCE2070", "operator '%s' not defined for '%s' and '%s'")
var UNARY_OP_INCOMPATIBLE_TYPES = newSemanticErrorCode("BCE2071", "operator '%s' not defined for '%s'")

// Aliases for backward compatibility
var INVALID_BINARY_OP = BINARY_OP_INCOMPATIBLE_TYPES
var INVALID_UNARY_OP = UNARY_OP_INCOMPATIBLE_TYPES

// Control flow errors - matching Java: BREAK_CANNOT_BE_OUTSIDE_LOOP, CONTINUE_CANNOT_BE_OUTSIDE_LOOP
var BREAK_CANNOT_BE_OUTSIDE_LOOP = newSemanticErrorCode("BCE2108", "break cannot be used outside of a loop")
var CONTINUE_CANNOT_BE_OUTSIDE_LOOP = newSemanticErrorCode("BCE2107", "continue cannot be used outside of a loop")

// Aliases for backward compatibility
var BREAK_OUTSIDE_LOOP = BREAK_CANNOT_BE_OUTSIDE_LOOP
var CONTINUE_OUTSIDE_LOOP = CONTINUE_CANNOT_BE_OUTSIDE_LOOP

// Function call errors - matching Java: TOO_MANY_ARGS_FUNC_CALL, MISSING_REQUIRED_PARAMETER
var TOO_MANY_ARGS_FUNC_CALL = newSemanticErrorCode("BCE2524", "too many arguments in call to '%s'")
var MISSING_REQUIRED_PARAMETER = newSemanticErrorCode("BCE2063", "missing required parameter '%s' in call to '%s'")

// Aliases for backward compatibility
var TOO_MANY_ARGS = TOO_MANY_ARGS_FUNC_CALL
var NOT_ENOUGH_ARGS = newSemanticErrorCode("BCE2063", "not enough arguments in call to '%s'")

// Return type errors - matching Java: INVOKABLE_MUST_RETURN
var INVOKABLE_MUST_RETURN = newSemanticErrorCode("BCE2095", "this function must return a result")
var RETURN_VALUE_REQUIRED = newSemanticErrorCode("BCE2068", "incompatible types: expected '%s', found 'nil'")
var RETURN_IN_NEVER_RETURNING_FUNCTION = newSemanticErrorCode("BCE2069", "return statement not allowed in a function that returns 'never'")

// Undefined symbol errors
var UNDEFINED_SYMBOL = newSemanticErrorCode("BCE2010", "undefined symbol '%s'")
var UNDEFINED_FUNCTION = newSemanticErrorCode("BCE2011", "undefined function '%s'")

// Variable errors
var VARIABLE_NOT_INITIALIZED = newSemanticErrorCode("BCE2520", "variable '%s' is not initialized")

// Unused return value - custom (not directly in Java but similar behavior)
var UNUSED_RETURN_VALUE = newSemanticErrorCode("BCE2067", "return value of function is not used")

// Undefined module - when using a module prefix without importing it
// Reference: BCE2000 in DiagnosticErrorCode.java
var UNDEFINED_MODULE = newSemanticErrorCode("BCE2000", "undefined module '%s'")

// Unreachable code - code after return, break, continue, or infinite loop
// Reference: BCE2106 in DiagnosticErrorCode.java
var UNREACHABLE_CODE = newSemanticErrorCode("BCE2106", "unreachable code")
