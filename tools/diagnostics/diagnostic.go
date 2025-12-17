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

import "fmt"

// Diagnostic represents a diagnostic message (error, warning, etc.) with location information.
// A diagnostic represents a compiler error, a warning or a message at a specific location in the source file.
type Diagnostic interface {
	Location() Location
	DiagnosticInfo() DiagnosticInfo
	Message() string
	Properties() []DiagnosticProperty[any]
	String() string
}

type diagnosticBase struct{}

// String returns a string representation of the diagnostic.
// This is the default implementation from the abstract Diagnostic class.
func (db diagnosticBase) String(d Diagnostic) string {
	var location string
	if d.Location().LineRange().FileName() == "" {
		location = ""
	} else {
		location = fmt.Sprintf(" [%s:%s]",
			d.Location().LineRange().FileName(),
			d.Location().LineRange().String())
	}
	return fmt.Sprintf("%s%s %s",
		d.DiagnosticInfo().Severity().String(),
		location,
		d.Message())
}
