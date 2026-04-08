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

type Base record {|
    int x = 10;
    string y = "hello";
|};

type Inherited record {|
    *Base;
    string z;
|};

type OverriddenWithNewDefault record {|
    *Base;
    int x = 42;
    string z;
|};

public function main() {
    Inherited r1 = {z: "world"};
    io:println(r1.x); // @output 10
    io:println(r1.y); // @output hello
    io:println(r1.z); // @output world

    Inherited r2 = {x: 5, z: "world"};
    io:println(r2.x); // @output 5
    io:println(r2.y); // @output hello
    io:println(r2.z); // @output world

    OverriddenWithNewDefault r3 = {z: "test"};
    io:println(r3.x); // @output 42
    io:println(r3.y); // @output hello
    io:println(r3.z); // @output test

    OverriddenWithNewDefault r4 = {x: 99, z: "test"};
    io:println(r4.x); // @output 99
    io:println(r4.y); // @output hello
    io:println(r4.z); // @output test
}
