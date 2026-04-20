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

int counter = 0;

function nextVal() returns int {
    counter = counter + 1;
    return counter;
}

function foo(int x, int y = x + 1) returns int {
    return x + y;
}

public function main() {
    // nextVal() should be called exactly once, returning 1
    // y defaults to x + 1 = 2, so result = 1 + 2 = 3
    io:println(foo(nextVal())); // @output 3
    // counter should be 1 after the call above
    io:println(counter); // @output 1
}
