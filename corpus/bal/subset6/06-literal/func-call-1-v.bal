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

function takeInt(int x) returns int {
    return x;
}

function takeFloat(float x) returns float {
    return x;
}

function takeDecimal(decimal x) returns decimal {
    return x;
}

public function main() {
    // int literal narrowed by function parameter type
    io:println(takeInt(42));
    io:println(takeFloat(42));
    io:println(takeDecimal(42));

    // float literal narrowed by function parameter type
    io:println(takeFloat(1.5));
    io:println(takeDecimal(1.5));

    // suffixed literals
    io:println(takeFloat(1.5f));
    io:println(takeDecimal(100d));
}
