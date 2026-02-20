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
    int:Unsigned8 u8 = 200;
    int:Unsigned8 a = u8 >> 2;
    io:println(a); // @output 50

    int:Unsigned8 b = u8 >>> 2;
    io:println(b); // @output 50

    int:Unsigned16 u16 = 2000;
    int:Unsigned16 c = u16 >> 2;
    io:println(c); // @output 500

    int:Unsigned32 u32 = 200000;
    int:Unsigned32 d = u32 >> 2;
    io:println(d); // @output 50000

    int e = u8 << 2;
    io:println(e); // @output 800
}
