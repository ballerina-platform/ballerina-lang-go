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

function incorrectArrayAccessTest() returns (string) {
    string[] animals;
    animals = ["Dog", "Cat"];
    return animals["cat"];
}

type IntOrString 1|"two";

function testIncorrectArrayAccessByFiniteType() {
    string[] animals = ["Dog", "Cat"];
    IntOrString x = 1;
    string s =  animals[x]; // should fail since x could be a string
}

const INDEX = -2;

public function testNegativeIndexArrayAcess() {
    byte[2] d = [1, 33];
    d[-1] = 3;
    d[INDEX] = 12;
}

function testUnaryConstExpressionInIndexAccess() {
    int[2] a = [1,2];
    int _ = a[-1];
}
