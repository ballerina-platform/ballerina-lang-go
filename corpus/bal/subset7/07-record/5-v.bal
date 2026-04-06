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

public function main() {
    final int c = 10;
    record {|int value = c;|} f = {};
    io:println(f.value); // @output 10

    final string s = "hello";
    record {|int value = c; string name = s;|} f2 = {};
    io:println(f2.value); // @output 10
    io:println(f2.name); // @output hello

    record {|int value = c; string name = s;|} f3 = {value: 20};
    io:println(f3.value); // @output 20
    io:println(f3.name); // @output hello

    record {|int x; int y = c;|} f4 = {x: 1};
    io:println(f4.x); // @output 1
    io:println(f4.y); // @output 10

    record {|int x; int y = c; string name = s;|} f5 = {x: 5};
    io:println(f5.x); // @output 5
    io:println(f5.y); // @output 10
    io:println(f5.name); // @output hello

    record {|int x; int y = c; string name = s;|} f6 = {x: 5, y: 20, name: "world"};
    io:println(f6.x); // @output 5
    io:println(f6.y); // @output 20
    io:println(f6.name); // @output world
}
