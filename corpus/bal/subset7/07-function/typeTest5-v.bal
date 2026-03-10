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

public function main() {
    F|int f = foo;
    if f is function {
        io:println(f(1)); // @output 2
    }
    else {
        io:println("unexpected");
    }
    F? g = foo;
    if g == () {
        io:println("unexpected");
    }
    else {
        io:println(g(2)); // @output 3
    }
    F? h = ();
    if h == () {
        io:println("working"); // @output working
    }
    else {
        io:println(h(3));
    }
}

function foo(int a) returns int {
    return a + 1;
}
