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
type F function(int, int) returns int|error;

public function main() {
    F f = foo;
    int|error res = f(1, 2);
    io:println(res is error); // @output false
    res = f(1, -1);
    io:println(res is error); // @output true
}

function foo(int... vals) returns int|error {
    int sum = 0;
    foreach int i in 0 ..< vals.length() {
        int val = vals[i];
        if val < 0 {
            return error("negative value");
        }
        sum += val;
    }
    return sum;
}
