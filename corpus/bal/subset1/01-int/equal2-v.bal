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


// @productions equality boolean if-else-stmt equality-expr return-stmt unary-expr function-call-expr int-literal
import ballerina/io;

public function main() {
    printBoolean(eq(9223372036854775806, 9223372036854775806)); // @output 1
    printBoolean(eq(9223372036854775806, 1)); //@output 0
    printBoolean(eq(9223372036854775806, 0)); //@output 0
    printBoolean(eq(9223372036854775806, -1)); //@output 0
    printBoolean(eq(9223372036854775806, -9223372036854775806)); //@output 0

    printBoolean(eq(1, 9223372036854775806)); // @output 0
    printBoolean(eq(1, 1)); // @output 1
    printBoolean(eq(1, 0)); //@output 0
    printBoolean(eq(1, -1)); //@output 0
    printBoolean(eq(1, -9223372036854775806)); //@output 0

    printBoolean(eq(0, 9223372036854775806)); // @output 0
    printBoolean(eq(0, 1)); // @output 0
    printBoolean(eq(0, 0)); // @output 1
    printBoolean(eq(0, -1)); //@output 0
    printBoolean(eq(0, -9223372036854775806)); //@output 0

    printBoolean(eq(-1, 9223372036854775806)); // @output 0
    printBoolean(eq(-1, 1)); // @output 0
    printBoolean(eq(-1, 0)); // @output 0
    printBoolean(eq(-1, -1)); // @output 1
    printBoolean(eq(-1, -9223372036854775806)); //@output 0

    printBoolean(eq(-9223372036854775806, 9223372036854775806)); // @output 0
    printBoolean(eq(-9223372036854775806, 1)); // @output 0
    printBoolean(eq(-9223372036854775806, 0)); // @output 0
    printBoolean(eq(-9223372036854775806, -1)); // @output 0
    printBoolean(eq(-9223372036854775806, -9223372036854775806)); // @output 1
}

function eq(int a, int b) returns boolean {
    return a == b;
}

function printBoolean(boolean b) {
    if b {
        io:println(1);
    } 
    else {
        io:println(0);
    }
}
