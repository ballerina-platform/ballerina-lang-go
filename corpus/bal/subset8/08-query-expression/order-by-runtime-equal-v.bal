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
    float?[] withNils = [(), (), 1.0];
    float?[] sortedNils = from var value in withNils
        order by value ascending
        select value;

    int[][] dupRows = [[1], [1], [1, 0]];
    int[][] sortedRows = from var row in dupRows
        order by row ascending
        select row;

    io:println(sortedNils); // @output [1.0,null,null]
    io:println(sortedRows); // @output [[1],[1],[1,0]]
}
