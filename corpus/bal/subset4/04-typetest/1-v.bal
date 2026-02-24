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
    any x = 1;
    if x is int {
        io:println("int"); // @output int
    }
    if x is string {
        io:println("string");
    }

    any s = "hello";
    if s is string {
        io:println("string"); // @output string
    }
    if s is decimal {
        io:println("decimal");
    }

    any f = 1.0;
    if f is float {
        io:println("float"); // @output float
    }
    if f is boolean {
        io:println("int");
    }

    any y = true;
    if y is boolean {
        io:println("boolean"); // @output boolean
    }
    if y is float {
        io:println("int");
    }

    any n = <any>();
    if n is () {
        io:println("nil"); // @output nil
    }
    if n is any {
        io:println("any"); // @output any
    }

    any z = <decimal>1;
    if z is decimal {
        io:println("decimal"); // @output decimal
    }
    if z is int {
        io:println("int");
    }
}
