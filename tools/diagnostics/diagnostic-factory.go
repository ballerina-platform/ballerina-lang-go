/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package diagnostics

// CreateDiagnostic creates a Diagnostic instance from the given details.
//
// Parameters:
//   - diagnosticInfo: static diagnostic information
//   - location: the location of the diagnostic
//   - args: arguments to diagnostic message format
//
// Returns a Diagnostic instance.
func CreateDiagnostic(diagnosticInfo DiagnosticInfo, location Location, args ...interface{}) Diagnostic {
	return NewDefaultDiagnostic(diagnosticInfo, location, []DiagnosticProperty[any]{}, args...)
}

// CreateDiagnosticWithProperties creates a Diagnostic instance from the given details.
//
// Parameters:
//   - diagnosticInfo: static diagnostic information
//   - location: the location of the diagnostic
//   - properties: properties associated with the diagnostic
//   - args: arguments to diagnostic message format
//
// Returns a Diagnostic instance.
func CreateDiagnosticWithProperties(diagnosticInfo DiagnosticInfo, location Location, properties []DiagnosticProperty[any], args ...interface{}) Diagnostic {
	return NewDefaultDiagnostic(diagnosticInfo, location, properties, args...)
}
