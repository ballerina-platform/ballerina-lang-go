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

// Package palnative provides the native-CLI implementation of pal.Platform.
// The HTTP factory and its TLS plumbing live in http.go; IO is small enough
// to inline here. Other environments (e.g. WASM/web-editor) supply their own
// pal.Platform without importing this package.
package palnative

import (
	"os"

	"ballerina-lang-go/pal"
)

// NewPlatform returns the native-CLI pal.Platform, wiring os.Stdout/Stderr for
// IO and NewHTTPClient for HTTP.
func NewPlatform() pal.Platform {
	return pal.Platform{
		IO: pal.IO{
			Stdout: func(p []byte) (n int, err error) { return os.Stdout.Write(p) },
			Stderr: func(p []byte) (n int, err error) { return os.Stderr.Write(p) },
		},
		HTTP: pal.HTTP{
			NewClient: NewHTTPClient,
		},
	}
}
