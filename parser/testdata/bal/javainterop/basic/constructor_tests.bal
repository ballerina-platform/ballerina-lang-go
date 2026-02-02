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

public class ClassA {

    public function newClassWithDefaultConstructor() returns handle = @java:Constructor {
        'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithDefaultConstructor"
    } external;

    public function newClassWithOneParamConstructor(handle h) returns handle = @java:Constructor {
        'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithOneParamConstructor"
    } external;

    public function newClassWithTwoParamConstructor(handle h1, handle h2) returns handle = @java:Constructor {
        'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithTwoParamConstructor"
    } external;
}

function testDefaultConstructor() returns handle {
    return newClassWithDefaultConstructor();
}

function testOneParamConstructor(handle h) returns handle {
    return newClassWithOneParamConstructor(h);
}

function testTwoParamConstructor(handle h1, handle h2) returns handle {
    return newClassWithTwoParamConstructor(h1, h2);
}

function testDefaultConstructorForClass() returns handle {
    ClassA classA = new;
    return classA.newClassWithDefaultConstructor();
}

function testOneParamConstructorForClass(handle h) returns handle {
    ClassA classA = new;
    return classA.newClassWithOneParamConstructor(h);
}

function testTwoParamConstructorForClass(handle h1, handle h2) returns handle {
    ClassA classA = new;
    return classA.newClassWithTwoParamConstructor(h1, h2);
}

// Interop functions

public function newClassWithDefaultConstructor() returns handle = @java:Constructor {
    'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithDefaultConstructor"
} external;

public function newClassWithOneParamConstructor(handle h) returns handle = @java:Constructor {
    'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithOneParamConstructor"
} external;

public function newClassWithTwoParamConstructor(handle h1, handle h2) returns handle = @java:Constructor {
    'class:"org/ballerinalang/nativeimpl/jvm/tests/ClassWithTwoParamConstructor"
} external;

