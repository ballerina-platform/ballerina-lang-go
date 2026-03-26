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
import ballerina/io;
type F function(int, int) returns int;

public function main() {
    io:println(fold(sum, [1, 2, 3, 4, 5], 0)); // @output 15
    io:println(fold(mul, [1, 2, 3, 4, 5], 1)); // @output 120
}

function fold(F f, int[] a, int init) returns int {
    int result = init;
    foreach int i in 0..<a.length() {
        result = f(result, a[i]);
    }
    return result;
}

function sum(int a, int b) returns int {
    return a + b;
}

function mul(int a, int b) returns int {
    return a * b;
}
