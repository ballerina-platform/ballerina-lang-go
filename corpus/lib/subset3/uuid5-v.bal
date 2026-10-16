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
import ballerina/uuid;

public function main() returns error? {
    uuid:Uuid r3 = check uuid:createType3AsRecord(uuid:NAME_SPACE_URL, "ballerina.io");
    io:println(check uuid:toString(r3)); // @output 8ffc8eb7-a6f3-33ba-a8f6-8d7e9bced7db
    uuid:Uuid r5 = check uuid:createType5AsRecord(uuid:NAME_SPACE_OID, "ballerina.io");
    io:println(check uuid:toString(r5)); // @output 39ae198f-502e-5ad8-9630-0839c55bd63c
    io:println(check uuid:createType5AsString(uuid:NAME_SPACE_X500, "ballerina.io")); // @output 655b2306-9bed-5f3d-9142-b4a76ad5514f
    io:println(check uuid:createType3AsString(uuid:NAME_SPACE_NIL, "ballerina.io")); // @output e1b734f8-51d8-3a5d-bc92-db87f8085e15
    uuid:Uuid|uuid:Error e5 = uuid:createType5AsRecord(uuid:NAME_SPACE_DNS, " ");
    io:println(e5 is uuid:Error ? e5.message() : ""); // @output Failed to create UUID of type 5
}
