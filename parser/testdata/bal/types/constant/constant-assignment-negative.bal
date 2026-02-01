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

final int PI = 3.14989;

final float floatVariableEx = "not a float";

final int intVariableEx = "not an int";

const int CONST1 = 3 + 5;
type A CONST1;
const A CONST2 = 3;

function testAssignInvalidValue() {
    A a = 3;
    CONST1 c = 4;
}

const map<string> X = {a: "a"};

type Foo record {|
   X x;
   int i;
|};

const Foo F1 = {x: {b : "a"}, i: 1, c: 2};
const record{|X x; int i;|} F2 = {x: {b : "a"}, i: 1, c: 2};

const X F3 = {b : "b"};

const string[] Y = base16 `aabb`;

const [string, int] Z =  base16 `aabb`;

const [170] Z1 =  base16 `aabb`;
