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

// @productions equality boolean equality-expr boolean-literal return-stmt any function-call-expr assign-stmt local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    io:println(mkInt(2) == 2); // @output true
    io:println(17 != mkInt(17)); // @output false
    io:println(mkBoolean(true) == true); // @output true
    io:println(false != mkBoolean(false)); // @output false
    int n = 2;
    io:println(mkInt(2) == n); // @output true
    n = 17;
    io:println(n != mkInt(17)); // @output false
    boolean b = true;
    io:println(mkBoolean(true) == b); // @output true
    b = false;
    io:println(b != mkBoolean(false)); // @output false
}

function mkInt(int n) returns any {
    return n;
}

function mkBoolean(boolean b) returns any {
    return b;
}
