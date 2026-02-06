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

// @productions equality if-else-stmt equality-expr return-stmt additive-expr function-call-expr int-literal
import ballerina/io;

public function main() {
    // @output 1
    // @output 2
    // @output 3
    // @output 4
    // @output 5
    foo(0, 5);
}

function inc(int n) returns int {
    return n + 1;
}

function foo(int depth, int maxDepth) {
    if depth == maxDepth {
        return;
    }
    io:println(inc(depth));
    foo(depth + 1, maxDepth);
}
