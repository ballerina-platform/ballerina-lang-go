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
type F function(int...) returns int;
type F1 function(int, int) returns int;

public function main() {
    final int base = 10;
    F1 f1 = function(int... vals) returns int {
       int sum = base;
       foreach int i in 0 ..< vals.length() {
            sum += vals[i];
       }
       return sum;
    };
    F1 f2 = f1;
    io:println(f1 === f2); // @output true
    F1 f3 = function(int... vals) returns int {
       int sum = base;
       foreach int i in 0 ..< vals.length() {
            sum += vals[i];
       }
       return sum;
    };
    io:println(f1 === f3); // @output false
    F f4 = <F>f1;
    io:println(f1 === f4); // @output true
}
