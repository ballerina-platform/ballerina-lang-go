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


// @productions multiplicative-expr unary-expr additive-expr local-var-decl-stmt int-literal
import ballerina/io;

public function main() {
  io:println(12 + 6 / 3);
  io:println(30 / 3 + 12);
  io:println(6 * 3 - 2);
  io:println(8 - 4 * 2 );
  io:println(9 + 4 % 3 );
  io:println(4 % 3 + 9);
  io:println(18 % 11 % 3);
  io:println(30 % 18 % 11 % 5);
  io:println(18 % 12 / 3);
  io:println(16 / 8 % 6);
  io:println( 4 + -3);
  io:println(-3 + 4);

  int i = 12;
  int j = 6;
  int k = 3;
  int l = 4;
  io:println(i + j / k);
  io:println(j / k + i);
  io:println(j * k - i);
  io:println(i - j * k );
  io:println(l % k + j);
  io:println(j % l % k);
}
