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
	int[] a = [1,2,3];
	int?[] b = [1,(),3];
	int?[] c = [1,(),4];

	io:println(b < b); // @output false
	io:println(b <= b); // @output true
	io:println(b > b); // @output false
	io:println(b >= b); // @output true

	io:println(a < b); // @output false
	io:println(a <= b); // @output false
	io:println(a > b); // @output false
	io:println(a >= b); // @output false

	io:println(b < a); // @output false
	io:println(b <= a); // @output false
	io:println(b > a); // @output false
	io:println(b >= a); // @output false

	io:println(b < c); // @output true
	io:println(b <= c); // @output true
	io:println(b > c); // @output false
	io:println(b >= c); // @output false

	io:println(c < b); // @output false
	io:println(c <= b); // @output false
	io:println(c > b); // @output true
	io:println(c >= b); // @output true
}
