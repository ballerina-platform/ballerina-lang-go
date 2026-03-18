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
//
// @productions int-literal boolean relational-expr function-call-expr local-var-decl-stmt
import ballerina/io;

public function main() {
    int? a = ();
    int? b = 1;
    int? c = ();

    // comparisons between nil and non-nil operands are always false
    io:println(a < b);
    io:println(a <= b);
    io:println(a > b);
    io:println(a >= b);

    io:println(b < a);
    io:println(b <= a);
    io:println(b > a);
    io:println(b >= a);

    // comparisons between two nil operands behave like equality/inequality
    io:println(a < c);
    io:println(a <= c);
    io:println(a > c);
    io:println(a >= c);
}

