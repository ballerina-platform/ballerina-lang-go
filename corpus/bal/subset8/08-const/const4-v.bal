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
// @productions float module-const-decl exact-equality equality-expr multiplicative-expr floating-point-literal unary-expr
import ballerina/io;

const A = 0.0 / 0.0;
const B = -0.0 / 0.0;
const C = 1.0 / 0.0;
const D = -1.0 / 0.0;
const E = -0.0;
const F = 0.0;

public function main() {
    io:println(A === B); // @output true
    io:println(C === D); // @output false
    io:println(E === F); // @output false
    io:println(A === C); // @output false
    io:println(B === D); // @output false
    io:println(E !== F); // @output true
    io:println(A !== B); // @output false
}
