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


// @productions while-stmt break-stmt equality if-else-stmt equality-expr relational-expr return-stmt unary-expr additive-expr function-call-expr assign-stmt local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    // @output 8
    // @output 6
    // @output 4
    // @output 2
    // @output 0
    // @output -1
    io:println(foo(10));
}

function foo(int x) returns int {
    int i = x;
    while (i >= 0) {
        i = i - 1;
        if (x - i == 2) {
            io:println(i);
            break;
        }
    }
    if (i < 0) {
        return -1;
    }
    return foo(i);
}
