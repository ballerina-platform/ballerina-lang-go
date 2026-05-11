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
    int[] nums = [3, 1, 2, 1, 4];
    int[] listResult = from var n in nums
        let int doubled = n * 2
        order by doubled ascending
        select doubled + 1;

    map<int> values = {"a": 3, "b": 1, "c": 2};
    int[] mapResult = from var v in values
        let int adjusted = v + 10
        order by adjusted descending
        select adjusted;

    io:println(listResult); // @output [3,3,5,7,9]
    io:println(mapResult); // @output [13,12,11]
}
