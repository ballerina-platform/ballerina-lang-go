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

int[*] moduleInts = [4, 5, 6];

public function main() {
    int[*] x = [1, 2];
    io:println(x); // @output [1,2]
    io:println(x.length()); // @output 2
    io:println(moduleInts.length()); // @output 3
    boolean[*] flags = [false, true];
    io:println(flags); // @output [false,true]
    int[*][2] grid = [[1, 2], [3, 4], [5, 6]];
    io:println(grid.length()); // @output 3
    io:println(grid[2][1]); // @output 6
    int[*] empty = [];
    io:println(empty.length()); // @output 0
}
