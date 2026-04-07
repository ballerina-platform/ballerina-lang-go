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
    decimal d1 = 1d;
    decimal d2 = 2d;
    decimal d3 = 3d;
    decimal d15 = 1.5d;
    decimal[] x = [d15];
    decimal[] y = [d3];
    decimal[] z = [3d];
    decimal[] a1 = [d1, d2];
    decimal[] a2 = [d2, d3];
    io:println(x < y); // @output true
    io:println(x > y); // @output false
    io:println(y <= x); // @output false
    io:println(y >= x); // @output true
    io:println(y >= z); // @output true
    io:println(y <= z); // @output true
    io:println(a1 < a2); // @output true
    io:println(a1 >= a2); // @output false
}
