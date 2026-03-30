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

type P function(int) returns int;
type Q function(int, int) returns int;
type R function(P, Q, int, int...) returns int;

public function main() {
    R x = fooBar;
    io:println(x(foo, bar, 1)); // @output 2
    io:println(x(foo, bar, 1, 2)); // @output 3
    io:println(x(foo, bar, 1, 2, 3)); // @output 6
    io:println(x(foo, bar, 1, 2, 3, 4)); // @output 10
}

function bar(int x, int y) returns int {
    return x + y;
}

function foo(int x) returns int {
    return x + 1;
}

function fooBar(P single, Q pair, int x, int... y) returns int {
    if y.length() == 0 {
        return single(x);
    }
    int curr = pair(x, y[0]);
    foreach int i in 1 ..< y.length() {
        curr = pair(curr, y[i]);
    }
    return curr;
}
