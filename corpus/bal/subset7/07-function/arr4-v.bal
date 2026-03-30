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
    F[] funcs = [foo, bar, baz];
    int[] results = exec(funcs, 2, 3);
    io:println(results); // @output [5,5,5]
}

function exec(F[] funcs, int a, int b) returns int[] {
    int[] results = [];
    foreach int i in 0..< funcs.length() {
        F func = funcs[i];
        results.push(func(a, b));
    }
    return results;
}

function foo(int|string a, int|string b) returns int {
    if (a is int && b is int) {
        return a + b;
    }
    return -1;
}

function bar(int... vals) returns int {
    int sum = 0;
    foreach int i in 0 ..< vals.length() {
        sum += vals[i];
    }
    return sum;
}

function baz(int|string... vals) returns int {
    int|string init = "";
    if vals[0] is int {
        init = 0;
    }
    int result = 0;
    foreach int i in 0 ..< vals.length() {
        result += foo(init, vals[i]);
    }
    return result;
}
