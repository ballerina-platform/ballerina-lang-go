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




function readBoolean(handle dataIS) returns boolean | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readByte(handle dataIS) returns byte | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readShort(handle dataIS) returns int | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readChar(handle dataIS) returns int | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readInt(handle dataIS) returns int | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readLong(handle dataIS) returns int | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readFloat(handle dataIS) returns float | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;

function readDouble(handle dataIS) returns float | error = @java:Method {
    'class: "java.io.DataInputStream"
} external;
