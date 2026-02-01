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

function lazyInitThrowArrayIndexOutOfBound () {
    int[] arr = [];
    int x;

    // This should panic an exception
    x = arr[0];
}

function lazyInitSizeZero() returns (string[]) {
    string[] names = [];

    return names;
}

function addValueToIntArray() returns (int[]) {
    int[] arr = [];
    int x;
    x = 4;
    arr[0] = -10;
    arr[15] = 2 + 18;
    arr[99] = 2147483647;
    arr[100] = -x;
    arr[115] = -2147483647;
    arr[199] = 6;

    return arr;
}

function addValueToFloatArray() returns (float[]) {
    float[] arr = [];
    arr[0] = -10.0;
    arr[15] = 2.5 ;
    arr[99] = 2147483647.1;
    arr[100] = 4.3;
    arr[115] = -2147483647.7;
    arr[199] = 6.9;

    return arr;
}


function commnetFunction() {
// TODO
// 1) Array add value test
// 2) Array get value test
// 3) Array grow test
// 4) Array maximum size test
// 5) Array grow and size change test
}

function testDefaultValueOfIntArrayElement() returns [int, int, int]{
    int[] a = [];
    a[5] = 45;
    return [a[0], a[1], a[5]];
}

function testDefaultValueOfJsonArrayElement() returns [json, json, json] {
    json[] j = [];
    j[5] = {name:"supun"};
    return [j[0], j[1], j[5]];
}

function testArrayGrowth () returns (int) {
    float value = 100.0;
    int[] ar = [];

    int count = 0;
    while (count < 20) {
        int intValue;
        intValue = <int>value;

        ar[intValue] = 1;
        value = value * 1.7;
        count = count + 1;
    }
    return ar.length();
}