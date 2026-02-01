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

import ballerina/lang.'float as floats;
import ballerina/lang.'int as ints;

function intToFloat (int value) returns (float) {
    float result;
    result = <float>value;
    return result;
}

function intToString (int value) returns (string) {
    string result;
    result = value.toString();
    return result;
}

function intToAny (int value) returns (any) {
    any result;
    result = value;
    return result;
}

function floatToInt (float value) returns (int) {
    int result;
    result = <int>value;
    return result;
}

function floatToString(float value) returns (string) {
    string result;
    result = value.toString();
    return result;
}

function floatToAny (float value) returns (any) {
    any result;
    result = value;
    return result;
}

function stringToInt(string value) returns (int|error) {
    int result;
    result = check ints:fromString(value);
    return result;
}

function stringToFloat(string value) returns (float|error) {
    float result;
    result = check floats:fromString(value);
    return result;
}

function stringToAny(string value) returns (any) {
    any result;
    result = value;
    return result;
}

function booleanToString(boolean value) returns (string) {
    string result;
    result = value.toString();
    return result;
}

function booleanToAny(boolean value) returns (any) {
    any result;
    result = value;
    return result;
}

function anyToInt () returns (int|error) {
    int i = 5;
    any a = i;
    int value;
    value = check trap <int>a;
    return value;
}

function anyToFloat () returns (float|error) {
    float f = 5.0;
    any a = f;
    float value;
    value = check trap <float>a;
    return value;
}

function anyToString () returns (string) {
    string s = "test";
    any a = s;
    string value;
    value = <string>a;
    return value;
}

function anyToBoolean () returns (boolean|error) {
    boolean b = false;
    any a = b;
    boolean value;
    value = check trap <boolean>a;
    return value;
}

function booleanappendtostring(boolean value) returns (string) {
    string result;
    result = value.toString() + "-append-" + value.toString();
    return result;
}
