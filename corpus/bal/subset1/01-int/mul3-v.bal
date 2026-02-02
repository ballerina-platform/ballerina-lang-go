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


// @productions multiplicative-expr return-stmt unary-expr function-call-expr int-literal
import ballerina/io;

public function main() {
    io:println(mul(9223372036854775806, 1)); // @output 9223372036854775806
    io:println(mul(9223372036854775806, 0)); // @output 0
    io:println(mul(9223372036854775806, -1)); // @output -9223372036854775806

    io:println(mul(1, 1)); // @output 1
    io:println(mul(1, 0)); // @output 0
    io:println(mul(1, -1)); // @output -1

    io:println(mul(0, 1)); // @output 0
    io:println(mul(0, 0)); // @output 0
    io:println(mul(0, -1)); // @output 0

    io:println(mul(-1, 1)); // @output -1
    io:println(mul(-1, 0)); // @output 0
    io:println(mul(-1, -1)); // @output 1

    io:println(mul(-9223372036854775806, 1)); // @output -9223372036854775806
    io:println(mul(-9223372036854775806, 0)); // @output 0
    io:println(mul(-9223372036854775806, -1)); // @output 9223372036854775806
}

function mul(int a, int b) returns int {
    return a * b;
}
