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
type R1 record {
    int l1;
    1|2 l2;
};

type R2 record {
    byte l1;
    2|3 l2;
};

public function main() {
    R1&R2 r = { l1: 5, l2: 2, "l3": "l" };
    byte? l1 = r.l1;
    io:println(l1); // @output 5
    io:println(r.l2); // @output 2
    io:println(r["l3"]); // @output l
}
