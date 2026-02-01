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

function stringTemplate1() returns (string) {
    string s = string `Hello ${name}`;
    return s;
}

function stringTemplate2() returns (string) {
    json name = {};
    string s = string `Hello ${name}`;
    return s;
}

public type Foo int|float|decimal|string|boolean|();

function stringTemplate3() returns string {
    Foo foo = 4;
    return string`${foo}`;
}

function stringTemplate4() returns string {
    () foo = ();
    return string`${foo}`;
}

function stringTemplate5() returns string {
    int[]|string[] x = [1, 2, 3];
    return string`${x}`;
}
