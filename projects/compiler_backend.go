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

package projects

// TargetPlatform represents the unique name of a supported compiler backend target.
// In Java, this is an interface with a code() method. In Go, we use a string type
// since the only operation is retrieving the code string.
// Java source: io.ballerina.projects.CompilerBackend.TargetPlatform
type TargetPlatform string

// Code returns the platform code string.
func (tp TargetPlatform) Code() string {
	return string(tp)
}

// CompilerBackend represents an abstract Ballerina compiler backend.
// The backend is responsible for persisting BIR to the filesystem.
// TODO(P6): Define full interface once backend design is finalized.
// Java source: io.ballerina.projects.CompilerBackend
type CompilerBackend interface {
	// TargetPlatform returns the supported target platform of this compiler backend.
	TargetPlatform() TargetPlatform
}
