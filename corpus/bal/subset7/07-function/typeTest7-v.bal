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
type F function(int) returns int;
type F2 function(int, int) returns int;
type F3 function(int, int, int) returns int;

public function main() {
    F|F2|F3 f = f1;
    io:println(call(f)); // @output 1
    f = f2;
    io:println(call(f)); // @output 3
    f = f3;
    io:println(call(f)); // @output 6
}

function call(F|F2|F3 f) returns int {
    if f is F {
        return f(1);
    }
    else if f is F2 {
        return f(1, 2);
    }
    else {
        return f(1, 2, 3);
    }
}

function f1(int i) returns int {
    return i;
}

function f2(int i, int j) returns int {
    return i + j;
}

function f3(int i, int j, int k) returns int {
    return i + j + k;
}
