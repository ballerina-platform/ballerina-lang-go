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


// @productions exact-equality boolean equality-expr boolean-literal return-stmt any function-call-expr int-literal
import ballerina/io;

public function main() {
    io:println(mkNil() === mkNil()); // @output true
    io:println(mkInt(1) !== mkInt(1)); // @output false
    io:println(mkBoolean(true) === mkBoolean(true)); // @output true

    // following are the boundaries of immediate vs heap int
    io:println(mkInt(-36028797018963969) === mkInt(-36028797018963969)); // @output true
    io:println(mkInt(-36028797018963968) === mkInt(-36028797018963968)); // @output true
    io:println(mkInt(36028797018963967) === mkInt(36028797018963967)); // @output true
    io:println(mkInt(36028797018963968) === mkInt(36028797018963968)); // @output true
}

function mkNil() returns any {
    return ();
}

function mkInt(int n) returns any {
    return n;
}

function mkBoolean(boolean b) returns any {
    return b;
}
