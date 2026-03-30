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
type F function(byte);
type F1 function(int);
type F2 function(int|float);
type F3 function(int|string);

public function main() {
    F f = foo;
    f(3); // @output foo called with 3
    f(10); // @output foo called with 10
    f = bar;
    f(3); // @output bar called with 3
    f(10); // @output bar called with 10
    F1 a = fooBar;
    a(1); // @output int
    F2 b = fooBar;
    b(1); // @output int
    b(1.0); // @output float
    F3 c = fooBar;
    c(1); // @output int
    c("c"); // @output string
}

function foo(int a) {
    io:println("foo called with ", a);
}

function bar(int b) {
    io:println("bar called with ", b);
}

function fooBar(int|float|string a) {
    if a is int {
        io:println("int");
    }
    else if a is float {
        io:println("float");
    }
    else {
        io:println("string");
    }
}
