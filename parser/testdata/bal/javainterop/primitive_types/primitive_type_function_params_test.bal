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

// Test a function that accepts a Ballerina boolean for a Java boolean
public function testCreateBoxedBooleanFromBBoolean(boolean value) returns handle {
    return createBoxedBooleanFromBBoolean(value);
}

public function createBoxedBooleanFromBBoolean(boolean value) returns handle = @java:Constructor {
    'class:"java.lang.Boolean",
    paramTypes:["boolean"]
} external;


// Test a function that accepts a Ballerina byte for a Java byte
public function testCreateBoxedByteFromBByte(byte value) returns handle {
    return createBoxedByteFromBByte(value);
}

public function createBoxedByteFromBByte(byte value) returns handle = @java:Constructor {
    'class:"java.lang.Byte",
    paramTypes:["byte"]
} external;


// Test a function that accepts a Ballerina int for a Java long
public function testCreateBoxedLongFromBInt(int value) returns handle {
    return createBoxedLongFromBInt(value);
}

public function createBoxedLongFromBInt(int value) returns handle = @java:Constructor {
    'class:"java.lang.Long",
    paramTypes:["long"]
} external;


// Test a function that accepts a Ballerina float for a Java double
public function testCreateBoxedDoubleFromBFloat(float value) returns handle {
    return createBoxedDoubleFromBFloat(value);
}

public function createBoxedDoubleFromBFloat(float value) returns handle = @java:Constructor {
    'class:"java.lang.Double",
    paramTypes:["double"]
} external;