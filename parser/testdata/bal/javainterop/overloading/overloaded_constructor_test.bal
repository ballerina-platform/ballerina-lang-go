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

public function testOverloadedConstructorsWithOneParam() returns [handle, handle] {
    handle bufferStrValue = java:fromString("string buffer value");
    handle stringBuffer = newStringBuffer(bufferStrValue);

    handle builderStrValue = java:fromString("string builder value");
    handle stringBuilder = newStringBuilder(builderStrValue);

    handle stringCreatedWithBuffer = newStringWithStringBuffer(stringBuffer);
    handle stringCreatedWithBuilder = newStringWithStringBuilder(stringBuilder);
    return [stringCreatedWithBuffer, stringCreatedWithBuilder];
}

public function newStringBuffer(handle strValue) returns handle = @java:Constructor {
    'class:"java.lang.StringBuffer",
    paramTypes:["java.lang.String"]
} external;

public function newStringBuilder(handle strValue) returns handle = @java:Constructor {
    'class:"java.lang.StringBuilder",
    paramTypes:["java.lang.String"]
} external;

public function newStringWithStringBuffer(handle buffer) returns handle = @java:Constructor {
    'class:"java.lang.String",
    paramTypes:["java.lang.StringBuffer"]
} external;

public function newStringWithStringBuilder(handle builder) returns handle = @java:Constructor {
    'class:"java.lang.String",
    paramTypes:["java.lang.StringBuilder"]
} external;

public function testOverloadedMethodsWithByteArrayParams(string strValue) returns string? {
    handle str = java:fromString(strValue);
    handle bytes = getBytes(str);
    sortByteArray(bytes);
    handle sortedStr = newString(bytes);
    return java:toString(sortedStr);
}

public function testOverloadedMethodsWithDifferentParametersOne(int intValue) {
    string intResult = getIntString(intValue).toString();
    assertEquality("5", intResult);
}

public function testOverloadedMethodsWithDifferentParametersTwo(string strValue) {
    string strResult = getString(strValue).toString();
    assertEquality("BALLERINA", strResult);
}

function getBytes(handle receiver) returns handle = @java:Method {
    'class: "java.lang.String"
} external;

function sortByteArray(handle src) = @java:Method {
    name: "sort",
    'class: "java.util.Arrays",
    paramTypes: [{'class: "byte", dimensions:1}]
} external;

function newString(handle bytes) returns handle = @java:Constructor {
    'class: "java.lang.String",
    paramTypes: [{'class: "byte", dimensions:1}]
} external;

function getString(string str) returns handle = @java:Method {
    name: "moveTo",
    'class: "org.ballerinalang.test.javainterop.overloading.pkg.Vehicle",
    paramTypes: ["io.ballerina.runtime.api.values.BString"]
} external;

function getIntString(int val) returns handle = @java:Method {
    name: "moveTo",
    'class: "org.ballerinalang.test.javainterop.overloading.pkg.Vehicle",
    paramTypes: ["java.lang.Long"]
} external;

const ASSERTION_ERROR_REASON = "AssertionError";

function assertEquality(any|error expected, any|error actual) {
    if expected is anydata && actual is anydata && expected == actual {
        return;
    }
    if expected === actual {
        return;
    }

    string expectedValAsString = expected is error ? expected.toString() : expected.toString();
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expectedValAsString + "', found '" + actualValAsString + "'");
}
