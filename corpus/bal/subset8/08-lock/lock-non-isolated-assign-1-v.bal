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

// Inside the lock, `arr` is a captured mutable variable from the
// enclosing function: it isn't the lock's restricted variable and
// it isn't an isolated module variable. The transfer-in walker
// rejects the LHS read because `int[]` is not a subtype of
// `Isolated`.
import ballerina/io;

function bump() {
    int[] arr = [];
    lock {
        arr = [1, 2, 3];
    }
    io:println(arr); // @output [1,2,3]
}

public function main() {
    bump();
}
