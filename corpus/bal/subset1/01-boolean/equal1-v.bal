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


// @productions equality boolean if-else-stmt equality-expr boolean-literal function-call-expr int-literal
import ballerina/io;

public function main() {
    printEq(true, true); // @output 1
    printEq(true, false); // @output 0
    printEq(false, true); // @output 0
    printEq(false, false); // @output 1
    printNotEq(true, true); // @output 0
    printNotEq(true, false); // @output 1
    printNotEq(false, true); // @output 1
    printNotEq(false, false); // @output 0
}

function printEq(boolean b1, boolean b2) {
    if b1 == b2 {
        io:println(1);
    }
    else {
        io:println(0);
    }
}

function printNotEq(boolean b1, boolean b2) {
    if b1 != b2 {
        io:println(1);
    }
    else {
        io:println(0);
    }
}