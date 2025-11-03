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

package constants

const (
	Underscore = "_"
	UserHome   = "user.home"
)

// SymbolFlags contains flag constants for symbols.
const (
	PUBLIC        int64 = 1
	NATIVE        int64 = 2
	FINAL         int64 = 4
	ATTACHED      int64 = 8
	READONLY      int64 = 32
	REQUIRED      int64 = 256
	PRIVATE       int64 = 1024
	OPTIONAL      int64 = 4096
	REMOTE        int64 = 32768
	CLIENT        int64 = 65536
	RESOURCE      int64 = 131072
	SERVICE       int64 = 262144
	TRANSACTIONAL int64 = 33554432
	CLASS         int64 = 268435456
	ISOLATED      int64 = 536870912
	ENUM          int64 = 8589934592
	ANY_FUNCTION  int64 = 549755813888
)

// IsFlagOn checks if a specific flag is set in the bitmask.
func IsFlagOn(bitmask, flag int64) bool {
	return (bitmask & flag) == flag
}
