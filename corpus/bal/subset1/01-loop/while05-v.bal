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


// @productions while-stmt relational-expr return-stmt additive-expr function-call-expr assign-stmt local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
    printInts(5);
    // @output 5
    // @output 4
    // @output 3
    // @output 2
    // @output 1
}

function printInts(int maxExclusive) {
    int i = maxExclusive;
    while 0 <= decrease(i) {
        io:println(i);
        i = decrease(i);
    }
}

function decrease(int x) returns int {
    return x - 1;
}
