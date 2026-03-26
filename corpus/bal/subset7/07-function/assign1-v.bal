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
type F function() returns boolean;
type B function(int, int) returns boolean;

public function main() {
    F f = foo;
    boolean x = f();
    io:println(x); // @output true

    B b = bar;
    boolean y = b(1, 2);
    io:println(y); // @output false
    io:println(b(2, 1)); // @output true
    io:println(FB(foo, bar, 1, 2)); // @output true
    io:println(FB(notFoo, bar, 1, 2)); // @output false
    io:println(FB(notFoo, bar, 2, 1)); // @output true

    B b2 = bar;
    io:println(b2(2, 1)); // @output true
}

function foo() returns boolean {
    return true;
}

function notFoo() returns boolean {
    return false;
}

function bar(int a, int b) returns boolean {
    return a > b;
}

function FB(F a, B b, int x, int y) returns boolean {
    return a() || b(x, y);
}
