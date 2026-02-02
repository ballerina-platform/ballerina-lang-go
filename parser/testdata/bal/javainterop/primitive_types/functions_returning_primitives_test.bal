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

import ballerina/jballerina.java;

// Interop function that returns a Ballerina boolean for a Java boolean
function testReturningBBooleanJBoolean(handle receiver, handle strValue) returns boolean {
    return contentEquals(receiver, strValue);
}

public function contentEquals(handle receiver, handle strValue) returns boolean = @java:Method {
    'class:"java.lang.String",
    paramTypes: ["java.lang.String"]
} external;


// Interop function that returns a Ballerina byte for a Java byte
function testReturningBByteJByte(handle receiver) returns byte {
    return byteValue(receiver);
}

public function byteValue(handle receiver) returns byte = @java:Method {
    'class:"java.lang.Long"
} external;


// Interop function that returns a Ballerina int for a Java long
function testReturningBIntJLong(handle receiver) returns int {
    return longValue(receiver);
}

public function longValue(handle receiver) returns int = @java:Method {
    name:"longValue",
    'class:"java.lang.Long"
} external;


// Interop function that returns a Ballerina float for a Java double
function testReturningBFloatJDouble(handle receiver) returns float {
    return doubleValue(receiver);
}

public function doubleValue(handle receiver) returns float = @java:Method {
    name:"doubleValue",
    'class:"java.lang.Double"
} external;

