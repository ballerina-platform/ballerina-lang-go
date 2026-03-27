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
    int a = 5;
    F f = function(int x) returns int {
        a += 1;
        return x + a;
    };
    a = 10;
    F f1 = function(int x) returns int {
        a += 2;
        return x + a;
    };
    int b = a + 5;
    int c = f(b);
    io:println(c); // @output 26
    io:println(a); // @output 11
    int d = f1(10);
    io:println(a); // @output 13
    io:println(d); // @output 23
    io:println(f(6)); // @output 20
    io:println(a); // @output 14
}
