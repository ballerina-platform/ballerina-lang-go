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

function (string a, int b) returns string glf1 = concat;

function (string a, boolean b) returns string glf2 = function (string a, boolean b) returns string {
                                                               return a + b.toString();
                                                           };

public function anyFunc = <function> glf1;

function (string, int) returns string glf3 = function (string a, int b) returns string {
    return "llll";
};

public function getGlf1() returns function (string a, int b) returns string {
    return glf1;
}

public function getGlf2() returns function (string a, boolean b) returns string {
    return glf2;
}

public function setGlf2(function (string a, boolean b) returns string func) {
    glf2 = func;
}

public function getGlf3() returns function (string, int) returns string {
    return glf3;
}

public function setGlf3(function (string, int) returns string func) {
    glf3 = func;
}

public function anyFunction1() returns function {
    return glf2;
}

public function anyFunction2(function (string, int) returns string func) returns function {
    return <function> func;
}

public function concat (string x, int y) returns string {
    string result = x + y.toString();
    return result;
}

public function (int i = 300, int j = 200) returns int func = (i, j) => i + j;

public function (int i = 300, int j = i + 200) returns int func1 = (i, j) => i + j;
