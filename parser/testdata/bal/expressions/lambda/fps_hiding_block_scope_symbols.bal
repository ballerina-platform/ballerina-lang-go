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

function name(int y) returns function (int) returns int {

    var func = function (int y) returns int {
        return y + y;
    };

    return func;
}

function nameArrow(int y, int z) returns int {
    function (int, int, int) returns int lambda = (x, y, z) => x + y + z;
    int x = 34; // this is not in a overlapping scope
    return lambda(12, 32, 33);
}

int z = 34;
function nameP() returns function (int) returns int {

    var func = function (int z) returns int {
        return z + z;
    };

    return func;
}

function nameArrowP() returns int {
    function (int, int) returns int lambda = (x, z) => x + z;
    return lambda(12, 32);
}

function nameX(int y) returns function (int) returns int {

    var func = function (int z) returns int {
        int y = 3;
        return y + y;
    };
    return func;
}

public function test2() {
    int a = 10;
    var x = function (function(int a = 10) l) {};
}

public function test1() {
    int a = 10;
    var x = function (function(int a, int b = a) l) {};
}
