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

function testInvlaidFloatValue() returns (float) {
    float b;
    b = 010.1;
    return b;
}

function testInvlaidFloatValue2() {
    float b = 999e9999999999;
    float c = 999e-9999999999;
    decimal|float d = 999e9999999999;
    int|decimal|float e = 99.9E99999999;
    int|decimal|float f = 99.9E-99999999;
}

type Bar 0x9999999p999999999999999999999999;

type Bar1 0x9999999p-999999999999999999999999;

type Baz 9999999999e9999999999999999999f;

type Baz1 9999999999e-9999999999999999999f;

0x999.9p999999999999999 x = 0x999.9p999999999999999;

type FloatType 9.99E+6111f;

function testOutOfRangeFloat() {
    float _ = 9.99E+6111f;
    9.99E+6111f _ = 9.99E+6111f;
    float _ = 9.99E+6111;
}
