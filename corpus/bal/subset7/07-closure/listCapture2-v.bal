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
type F function(int) returns int;
type F1 function(int) returns F;
public function main() {
    final int[] arr = [5];
    F1 f1 = function(int y) returns F {
        int b = y * 2;
        return function(int x) returns int {
            arr.push(x + b);
            return arr[0] + arr[1];
        };
    };
    F f = f1(10);
    io:println(arr); // @output [5]
    int i = f(20);
    io:println(i); // @output 45
    io:println(arr); // @output [5,40]
}
