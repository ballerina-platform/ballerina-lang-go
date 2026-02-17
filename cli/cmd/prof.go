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

package main

import "github.com/spf13/cobra"

// Profiler manages the pprof HTTP server lifecycle
type Profiler interface {
	// RegisterFlags registers profiling-related flags on the given command.
	// In release builds this is a no-op, so the flags are invisible.
	RegisterFlags(cmd *cobra.Command)
	Start() error
	Stop() error
}

// Global profiler instance (implementation set by init() in build-tagged files)
var profiler Profiler
