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

import testorg/modulevar.vars;

public function main() {
    // Simple variable without init
    io:println(vars:count); // @output 10
    io:println(vars:label); // @output initialized

    // Simple variable with init
    io:println(vars:maxRetries); // @output 3
    io:println(vars:greeting); // @output hello
    io:println(vars:verbose); // @output true

    // List/mapping constructor
    io:println(vars:primes); // @output [2,3,5,7,11]
    io:println(vars:limits); // @output {"min":0,"max":100}

    // Query expression
    io:println(vars:oddNumbers); // @output [1,3,5]
    io:println(vars:tripled); // @output [3,6,9,12,15,18]
}
