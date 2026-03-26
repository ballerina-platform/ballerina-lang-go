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

public function main() {
    int val = 10;
    F a = function(int p) returns int {
        val *= 2;
        return val + p;
    };
    F b = function(int p) returns int {
        val = val + a(p);
        return val;
    };
    F c = function(int p) returns int {
        val = a(p) + val;
        return val;
    };
    io:println(b(3)); // @output 33
    io:println(val); // @output 33
    val = 10;
    io:println(c(3)); // @output 43
    io:println(val); // @output 43
}
