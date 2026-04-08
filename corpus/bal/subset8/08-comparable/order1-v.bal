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
    int[][] a = [[1], [2]];
    int[][] b = [[1], [3]];
    io:println(a < b); // @output true
    io:println(a > b); // @output false

    int[]?[] c = [[0], [-1]];
    int[][] d = [[0], [-2]];
    io:println(c < d); // @output false
    io:println(c > d); // @output true

    int[]?[] e = [(), [1]];
    int[][] f = [[0], [1]];
    io:println(e == f); // @output false
    io:println(e < f); // @output false
    io:println(e > f); // @output false
}
