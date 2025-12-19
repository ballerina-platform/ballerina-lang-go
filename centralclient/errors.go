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

package centralclient

type CentralClientError struct {
	message string
}

func (e *CentralClientError) Error() string {
	return e.message
}

func NewCentralClientError(message string) *CentralClientError {
	return &CentralClientError{message: message}
}

type ConnectionError struct {
	CentralClientError
}

func NewConnectionError(message string) *ConnectionError {
	return &ConnectionError{
		CentralClientError: CentralClientError{message: message},
	}
}

type NoPackageError struct {
	CentralClientError
}

func NewNoPackageError(message string) *NoPackageError {
	return &NoPackageError{
		CentralClientError: CentralClientError{message: message},
	}
}

type PackageAlreadyExistsError struct {
	CentralClientError
	version string
}

func (e *PackageAlreadyExistsError) Version() string {
	return e.version
}

func NewPackageAlreadyExistsError(message string, version string) *PackageAlreadyExistsError {
	return &PackageAlreadyExistsError{
		CentralClientError: CentralClientError{message: message},
		version:            version,
	}
}
