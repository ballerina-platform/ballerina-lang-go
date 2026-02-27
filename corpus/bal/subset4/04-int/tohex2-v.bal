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

// @productions string-literal equality equality-expr unary-expr additive-expr local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    int i = 42;
    string s = i.toHexString();
    io:println(s); // @output 2a
    io:println(s == "2a"); // @output true

    io:println(0.toHexString());    // @output 0
    io:println((-1).toHexString()); // @output -1

    io:println(0xff.toHexString());    // @output ff
    io:println((-0xaa).toHexString()); // @output -aa

    io:println(9223372036854775807.toHexString());        // @output  7fffffffffffffff
    io:println((-9223372036854775807 - 1).toHexString()); // @output -8000000000000000
}
