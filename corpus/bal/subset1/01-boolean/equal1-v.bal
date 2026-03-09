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
    printEq(true, true);
    printEq(true, false);
    printEq(false, true);
    printEq(false, false);
    printNotEq(true, true);
    printNotEq(true, false);
    printNotEq(false, true);
    printNotEq(false, false);
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