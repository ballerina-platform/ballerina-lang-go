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

import (
	"fmt"

	"ballerina-lang-go/tools/text"
)

// DefaultDiagnostic is an internal implementation of the Diagnostic interface that is used by the DiagnosticFactory
// to create diagnostics.
type DefaultDiagnostic interface {
	Diagnostic
}

type defaultDiagnosticImpl struct {
	diagnosticBase
	diagnosticInfo DiagnosticInfo
	location       Location
	properties     []DiagnosticProperty[any]
	message        string
}

func NewDefaultDiagnostic(diagnosticInfo DiagnosticInfo, location Location, properties []DiagnosticProperty[any], args ...any) DefaultDiagnostic {
	message := formatMessage(diagnosticInfo.MessageFormat(), args...)
	return &defaultDiagnosticImpl{
		diagnosticInfo: diagnosticInfo,
		location:       location,
		properties:     properties,
		message:        message,
	}
}

func (dd *defaultDiagnosticImpl) Location() Location {
	return dd.location
}

func (dd *defaultDiagnosticImpl) DiagnosticInfo() DiagnosticInfo {
	return dd.diagnosticInfo
}

func (dd *defaultDiagnosticImpl) Message() string {
	return dd.message
}

func (dd *defaultDiagnosticImpl) Properties() []DiagnosticProperty[any] {
	return dd.properties
}

func (dd *defaultDiagnosticImpl) String() string {
	lineRange := dd.location.LineRange()
	filePath := lineRange.FileName()

	startLine := lineRange.StartLine()
	endLine := lineRange.EndLine()

	oneBasedStartLine := text.LinePositionFromLineAndOffset(startLine.Line()+1, startLine.Offset()+1)
	oneBasedEndLine := text.LinePositionFromLineAndOffset(endLine.Line()+1, endLine.Offset()+1)
	oneBasedLineRange := text.LineRangeFromLinePositions(filePath, oneBasedStartLine, oneBasedEndLine)

	return fmt.Sprintf("%s [%s:%s] %s",
		dd.diagnosticInfo.Severity().String(),
		filePath,
		oneBasedLineRange.String(),
		dd.Message())
}

func formatMessage(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}
