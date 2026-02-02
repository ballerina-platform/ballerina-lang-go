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

function incorrectMapAccessTest() returns (string?) {
    map<any> animals = {};
    animals["dog"] = "Jimmy";
    return animals[0];
}

function accessAllFields() {
    map<any> fruits = {"name":"John", "address":"unkown"};
    any a = fruits.*;
}

function accessUninitMap() {
    map<int> ints;
    ints["0"] = 0;

    map<int> m1 = getUninitializedMap11(ints);

    map<int> m2 = getUninitializedMap21();
}


function getUninitializedMap11(map<int> m) returns map<int> {
    map<int> m2 = getUninitializedMap12(m);
    return m2;
}

function getUninitializedMap12(map<int> m) returns map<int> {
    map<int> m2 = m;
    return m2;
}

function getUninitializedMap21() returns map<int> {
    map<int> m3 = getUninitializedMap22();
    return m3;
}

function getUninitializedMap22() returns map<int> {
    map<int> m4;
    return m4;
}
