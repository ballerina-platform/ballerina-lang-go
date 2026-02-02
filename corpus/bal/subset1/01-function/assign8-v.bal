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


// @productions list-type-descriptor list-constructor-expr assign-stmt local-var-decl-stmt string-literal function-call-expr
import ballerina/io;

function foo(int[] arr, int i) returns int[] {
    arr[i] = i + 10;
    return arr;
}

public function main() {
    int[] arr = [];
    _ = foo(arr, 0);
    _ = foo(arr, 1);
    _ = foo(arr, 2);
    io:println(arr); // @output [10,11,12]

    string str = "test str";
    _ = str;
    io:println(str); // @output test str
}
