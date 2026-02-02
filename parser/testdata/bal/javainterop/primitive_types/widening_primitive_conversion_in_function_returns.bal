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

// Interop function that returns a Ballerina int for a Java byte
function testReturningBIntJByte(handle receiver) returns int {
    return getBIntJByte(receiver);
}

public function getBIntJByte(handle receiver) returns int = @java:Method {
    name:"byteValue",
    'class:"java.lang.Long"
} external;


// Interop function that returns a Ballerina int for a Java short
function testReturningBIntJShort(handle receiver) returns int {
    return getBIntJShort(receiver);
}

public function getBIntJShort(handle receiver) returns int = @java:Method {
    name:"shortValue",
    'class:"java.lang.Long"
} external;


// Interop function that returns a Ballerina int for a Java char
function testReturningBIntJChar(handle receiver) returns int {
    return getBIntJChar(receiver) + 3;
}

public function getBIntJChar(handle receiver) returns int = @java:Method {
    name:"charValue",
    'class:"java.lang.Character"
} external;


// Interop function that returns a Ballerina int for a Java int
function testReturningBIntJInt(handle receiver) returns int {
    return getBIntJInt(receiver);
}

public function getBIntJInt(handle receiver) returns int = @java:Method {
    name:"intValue",
    'class:"java.lang.Long"
} external;


// Interop function that returns a Ballerina float for a Java byte
function testReturningBFloatJByte(handle receiver) returns float {
    return getBFloatJByte(receiver);
}

public function getBFloatJByte(handle receiver) returns float = @java:Method {
    name:"byteValue",
    'class:"java.lang.Double"
} external;


// Interop function that returns a Ballerina float for a Java short
function testReturningBFloatJShort(handle receiver) returns float {
    return getBFloatJShort(receiver);
}

public function getBFloatJShort(handle receiver) returns float = @java:Method {
    name:"shortValue",
    'class:"java.lang.Double"
} external;


// Interop function that returns a Ballerina float for a Java char
function testReturningBFloatJChar(handle receiver) returns float {
    return getBFloatJChar(receiver);
}

public function getBFloatJChar(handle receiver) returns float = @java:Method {
    name:"charValue",
    'class:"java.lang.Character"
} external;


// Interop function that returns a Ballerina float for a Java int
function testReturningBFloatJInt(handle receiver) returns float {
    return getBFloatJInt(receiver);
}

public function getBFloatJInt(handle receiver) returns float = @java:Method {
    name:"intValue",
    'class:"java.lang.Double"
} external;


// Interop function that returns a Ballerina float for a Java float
function testReturningBFloatJFloat(handle receiver) returns float {
    return getBFloatJFloat(receiver);
}

public function getBFloatJFloat(handle receiver) returns float = @java:Method {
    name:"floatValue",
    'class:"java.lang.Double"
} external;
