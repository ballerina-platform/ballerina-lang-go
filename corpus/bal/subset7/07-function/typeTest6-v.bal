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

public function main() {
    F|F2 f = f1;
    if f is F {
        io:println(f(1)); // @output 1
    }
    else {
        io:println(f(1, 2));
    }
    f = f2;
    if f is F {
        io:println(f(1));
    }
    else {
        io:println(f(1, 2)); // @output 3
    }
}

function f1(int i) returns int {
    return i;
}

function f2(int i, int j) returns int {
    return i + j;
}
