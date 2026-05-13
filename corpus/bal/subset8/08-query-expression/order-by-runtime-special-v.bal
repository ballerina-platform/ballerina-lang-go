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
    float?[] specialFloats = [4.0, (), 3.5, 2.0];

    float?[] ascSpecialFloats = from var value in specialFloats
        order by value ascending
        select value;

    float?[] descSpecialFloats = from var value in specialFloats
        order by value descending
        select value;

    float?[][] specialRows = [[1.0, ()], [1.0, 2.0], [1.0], [2.0]];

    float?[][] ascSpecialRows = from var row in specialRows
        order by row ascending
        select row;

    float?[][] descSpecialRows = from var row in specialRows
        order by row descending
        select row;

    io:println(ascSpecialFloats); // @output [2.0,3.5,4.0,null]
    io:println(descSpecialFloats); // @output [4.0,3.5,2.0,null]
    io:println(ascSpecialRows); // @output [[1.0],[1.0,2.0],[1.0,null],[2.0]]
    io:println(descSpecialRows); // @output [[2.0],[1.0,2.0],[1.0,null],[1.0]]
}
