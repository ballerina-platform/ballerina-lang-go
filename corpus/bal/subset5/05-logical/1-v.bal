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
    io:println(true && f("a", true));  // @output a
                                       // @output true

    io:println(false && f("b", true));  // @output false

    io:println(false || f("c", false));  // @output c
                                         // @output false

    io:println(true || f("c", false));  // @output true
}

public function f(string s, boolean b) returns boolean {
    io:println(s);
    return b;
}
