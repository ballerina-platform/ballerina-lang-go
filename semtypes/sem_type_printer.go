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

package semtypes

import "strings"

// bitsetToTypeNames converts a bitset to a comma-separated list of type names.
// Returns empty string for empty bitset (0).
// Type names are returned in bitset order without the "BT_" prefix.
func bitsetToTypeNames(bitset int) string {
	if bitset == 0 {
		return ""
	}

	var builder strings.Builder
	first := true

	for i := 0; i < VT_COUNT; i++ {
		if (bitset & (1 << i)) != 0 {
			if !first {
				builder.WriteString(", ")
			}
			typeCode := BasicTypeCodeFrom(i)
			typeName := typeCode.String()
			// Strip "BT_" prefix for clean output
			cleanName := strings.TrimPrefix(typeName, "BT_")
			builder.WriteString(cleanName)
			first = false
		}
	}

	return builder.String()
}
