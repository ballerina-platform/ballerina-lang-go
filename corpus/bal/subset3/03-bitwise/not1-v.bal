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

// @productions bitwise-complement-expr local-var-decl-stmt int-literal unary-expr assign-stmt
import ballerina/io;
public function main() {
    int i = 5;
    io:println(~i);

    i = -5;
    io:println(~i);

    i = 0;
    io:println(~i);

    i = -1;
    io:println(~i);

    -1 j = -1;
    io:println(~j);

    5 k = 5;
    io:println(~k);

    i = 9223372036854775807; // MAX_INT
    io:println(~i);
}
