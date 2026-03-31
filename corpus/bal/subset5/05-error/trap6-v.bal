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

function foo(int x) returns int {
    if x == 0 {
        panic error("zero!");
    }
    return x * 2;
}

function boom() returns int {
    panic error("deep");
}

function bar(int|error x, int|error y) returns int|error {
    if x is error {
        return x;
    }
    if y is error {
        return y;
    }
    return x + y;
}

public function main() {
    var result = trap bar(
        trap foo(0),
        trap trap boom()
    );

    io:println(result); // @output error("zero!")
}
