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

// DiagnosticRelatedInformation represents a message and location related to a particular Diagnostic.
// A sample usage would be to record all symbol information related to duplicate symbol error.
type DiagnosticRelatedInformation interface {
	Location() Location
	Message() string
}

type diagnosticRelatedInformationImpl struct {
	location Location
	message  string
}

func NewDiagnosticRelatedInformation(location Location, message string) DiagnosticRelatedInformation {
	return &diagnosticRelatedInformationImpl{
		location: location,
		message:  message,
	}
}

func (dri diagnosticRelatedInformationImpl) Location() Location {
	return dri.location
}

func (dri diagnosticRelatedInformationImpl) Message() string {
	return dri.message
}
