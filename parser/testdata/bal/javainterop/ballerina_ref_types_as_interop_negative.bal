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

// Finite type interop

type ALL_INT 1|2|3|4|5;
type MIX_TYPE 1 | 2 | "hello" | true | false;

function testAcceptAllInts() returns int {
    ALL_INT i = 4;
    return acceptAllInts(i);
}

function getAllInts() returns ALL_INT = @java:Method {
    name:"getAllFloats",
    'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
} external;

function acceptAllInts(ALL_INT x) returns int = @java:Method {
    name:"acceptAllFloats",
    'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
} external;

function getMixType() returns MIX_TYPE = @java:Method {
    name:"getAllInts",
    'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
} external;

function acceptMixType(MIX_TYPE x) returns any = @java:Method {
    name:"acceptAllInts",
    'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
} external;

function testInvalidIntersectionParamType(map<int> & readonly m) =
    @java:Method {
        name:"acceptImmutableValue",
        'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
    } external;

function testInvalidIntersectionReturnType(int[] & readonly a) returns int[] & readonly =
    @java:Method {
        name:"acceptAndReturnImmutableArray",
        'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
    } external;

function testInvalidReadOnlyParamType(readonly r) =
    @java:Method {
        name:"acceptReadOnlyValue",
        'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
    } external;

function testReturnReadOnlyValue(function () returns int f) returns readonly =
    @java:Method {
        name:"returnReadOnlyValue",
        'class:"org/ballerinalang/test/javainterop/RefTypeNegativeTests"
    } external;
