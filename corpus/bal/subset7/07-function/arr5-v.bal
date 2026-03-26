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
type Point int[];
type Point2D [int, int];

type IncrementFn function(Point2D, int) returns Point;

public function main() {
    Point2D val = [1, 2];
    IncrementFn fn = point2DIncOverride;
    Point result = fn(val, 1);
    io:println(result); // @output [2,3]
}

function point2DIncOverride(Point point, int increment) returns Point {
    return [point[0] + increment, point[1] + increment];
}
