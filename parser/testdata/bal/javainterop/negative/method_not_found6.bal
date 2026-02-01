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

function getFuture(typedesc<anydata> td, future<anydata> f) returns future<anydata> = @java:Method {
    'class: "org.ballerinalang.nativeimpl.jvm.tests.StaticMethods",
    name: "getFuture",
    paramTypes: ["io.ballerina.runtime.api.values.BTypedesc"]
} external;

function getTypeDesc(typedesc<anydata> td, future<anydata> f) returns typedesc<anydata> = @java:Method {
    'class: "org.ballerinalang.nativeimpl.jvm.tests.StaticMethods",
    name: "getTypeDesc",
    paramTypes: ["io.ballerina.runtime.api.values.BFuture"]
} external;

function getFutureOnly(future<anydata> f) returns future<anydata> = @java:Method {
    'class: "org.ballerinalang.nativeimpl.jvm.tests.StaticMethods",
    name: "getFutureOnly",
    paramTypes: ["io.ballerina.runtime.api.values.BFuture", "io.ballerina.runtime.api.values.BTypedesc"]
} external;

function getTypeDescOnly(typedesc<anydata> td) returns typedesc<anydata> = @java:Method {
    'class: "org.ballerinalang.nativeimpl.jvm.tests.StaticMethods",
    name: "getTypeDescOnly",
    paramTypes: ["io.ballerina.runtime.api.values.BTypedesc", "io.ballerina.runtime.api.values.BFuture"]
} external;
