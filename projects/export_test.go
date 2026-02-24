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

// This file exports internal types for testing purposes only.
// It is only compiled during test runs.

// TopologicallySortedModuleNames returns module names in dependency order for testing.
func (r *PackageResolution) TopologicallySortedModuleNames() []string {
	names := make([]string, len(r.topologicallySortedModuleList))
	for i, modCtx := range r.topologicallySortedModuleList {
		names[i] = modCtx.getModuleName().String()
	}
	return names
}
