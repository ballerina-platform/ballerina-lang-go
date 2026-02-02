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

int glbVarInt = 800;
string glbVarString = "value";
float glbVarFloat = 99.34323;
any glbVarAny = 88343;

float glbVarFloatChange = 99.0;

float glbVarFloat1 = glbVarFloat;

json glbVarJson = {};

float glbVarFloatLater = 0.0;

function getGlobalVars() returns [int, string, float, any] {
    return [glbVarInt, glbVarString, glbVarFloat, glbVarAny];
}

function accessGlobalVar() returns int | error {
    int value;
    value = check trap <int>glbVarAny;
    return (glbVarInt + value);
}

function changeGlobalVar(int addVal) returns float {
    glbVarFloatChange = 77.0 + <float> addVal;
    float value = glbVarFloatChange;
    return value;
}

function getGlobalFloatVar() returns float {
    _ = changeGlobalVar(3);
    return glbVarFloatChange;
}

function getGlobalVarFloat1() returns float {
    return glbVarFloat1;
}

function initializeGlobalVarSeparately() returns [json, float] {
    glbVarJson = {"name" : "James", "age": 30};
    glbVarFloatLater = 3432.3423;
    return [glbVarJson, glbVarFloatLater];
}