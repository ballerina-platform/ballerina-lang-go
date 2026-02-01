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


// @productions return-stmt unary-expr additive-expr function-call-expr int-literal
import ballerina/io;

public function main() {
    io:println(8 + 5 - 11); // @output 2
    io:println(8 - 11 + 5); // @output 2
    io:println(3 + 5 - 5 + 9); // @output 12
    io:println(add(3, 5, -5, 9)); // @output 12
    io:println(3 - 5 + 9 - 4); // @output 3
    io:println(add(3, -5, 9, -4)); // @output 3
}

function add(int a, int b, int c, int d) returns int {
    return a + b + c + d;
}
