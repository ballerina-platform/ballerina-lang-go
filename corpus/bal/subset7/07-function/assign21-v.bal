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
type F function(int, int, int...) returns int;
type G function(int, int, int...) returns int[];

public function main() {
    F f = foo;
    int result = f(1, 2, 3, 4, 5);
    io:println(result); // @output 15
    F g = bar;
    io:println(g(1, 2, 3, 4, 5)); // @output 120
    io:println(f(g(1,2), g(3,4), 5)); // @output 19
    G h = baz;
    io:println(h(1, 2, 3)); // @output [1,3,6]
}

function foo(int init, int... rest) returns int {
    int result = init;
    foreach int i in 0..< rest.length() {
        result += rest[i];
    }
    return result;
}

function bar(int... rest) returns int {
    int result = 1;
    foreach int i in 0..< rest.length() {
        result *= rest[i];
    }
    return result;
}

function baz(int init, int... rest) returns int[] {
    int[] result = [init];
    int current = init;
    foreach int i in 0..< rest.length() {
        current += rest[i];
        result.push(current);
    }
    return result;
}
