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
    map<float?> valuesById = {"a": 4.0, "b": (), "c": 3.5, "d": 2.0};

    float?[] ascFloats = from var value in valuesById
        order by value ascending
        select value;

    float?[] descFloats = from var value in valuesById
        order by value descending
        select value;

    int[][] rows = [[1, 2], [1, 1], [2], [1], [1, 1, 0]];

    int[][] ascRows = from var row in rows
        order by row ascending
        select row;

    int[][] descRows = from var row in rows
        order by row descending
        select row;

    io:println(ascFloats); // @output [2.0,3.5,4.0,null]
    io:println(descFloats); // @output [4.0,3.5,2.0,null]
    io:println(ascRows); // @output [[1],[1,1],[1,1,0],[1,2],[2]]
    io:println(descRows); // @output [[2],[1,2],[1,1,0],[1,1],[1]]
}
