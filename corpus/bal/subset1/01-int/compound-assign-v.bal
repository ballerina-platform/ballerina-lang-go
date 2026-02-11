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

public function main() {
    int loopCount = 1000000;
    int i = 0;
    int total = 0;

    while (i < loopCount) {
        int a = sum(i, i + 1);
        int b = addOffset(a);
        total += b;
        i += 1;
    }

    io:println("Loop count: ", loopCount); // @output Loop count:  1000000
    io:println("Final total: ", total); // @output Final total:  1000010000000
}

function sum(int x, int y) returns int {
    return x + y;
}

function addOffset(int value) returns int {
    return value + 10;
}
