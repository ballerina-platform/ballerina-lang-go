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


// @productions boolean if-else-stmt relational-expr boolean-literal return-stmt function-call-expr int-literal
import ballerina/io;

public function main() {
    printBoolean(greaterThan(true, false));
    printBoolean(greaterThan(true, true));
    printBoolean(greaterThan(false, false));
    printBoolean(lessThan(true, false));
    printBoolean(lessThan(false, true));
    printBoolean(lessThan(true, true));
    printBoolean(lessThan(false, false));
}

function printBoolean(boolean b) {
    if b {
        io:println(1);
    }
    else {
        io:println(0);
    }
}

function lessThan(boolean x, boolean y) returns boolean {
    return x < y;
}

function greaterThan(boolean x, boolean y) returns boolean {
    return x > y;
}