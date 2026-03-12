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
type F function(int) returns int?;

public function main() {
    F f = foo;
    int? i = f(10); // @output 10
    if (i is ()) {
        io:println("working"); // @output working
    }
    else {
        io:println("unexpected");
    }
    f = bar;
    int? j = f(10);
    if (j is int) {
        io:println(j); // @output 11
    }
    else {
        io:println("unexpected");
    }
}

function foo(int i) {
    io:println(i);
}

function bar(int i) returns int {
    return i + 1;
}
