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

function testUndeclaredStructAccess() {
    string name = "";
    dpt1[name] = "HR";
}

function testUndeclaredAttributeAccess() {
    string name;
    Department dpt = {};
    dpt["id"] = "HR";
}
      
type Department record {|
    string dptName = "";
    int count = 0;
|};

function testInvalidTypeOfIndexExpression1() {
    Department dpt = {};
    int index = 1;
    var x = dpt[index]; // incompatible types: expected 'string', found 'int'
}

function testInvalidTypeOfIndexExpression2() {
    Department dpt = {};
    string index = "dptName";
    string x = dpt[index]; // incompatible types: expected 'string', found 'string|int?'
}

const string FIELD_FOUR = "fieldFour";
type FiniteOne "fieldOne"|"fieldTwo"|0;
type FiniteTwo 0|1;
type FiniteThree "fieldOne"|"fieldTwo"|"fieldThree";
type FiniteFour FiniteThree|FIELD_FOUR;
type NoIntersection "F1"|"F2"|"F3"|FIELD_FOUR;

type Foo record {
    string|boolean fieldOne;
    int fieldTwo;
    float fieldThree;
};

type Bar record {|
    string|boolean fieldOne;
    int fieldTwo;
|};

function testFiniteTypeAsIndex() {
    FiniteOne f1 = "fieldOne";
    FiniteTwo f2 = 0;
    FiniteThree f3 = "fieldOne";
    NoIntersection f4 = "fieldFour";
    NoIntersection f5 = "F2";


    Foo foo = { fieldOne: "S", fieldTwo: 12, fieldThree: 98.9 };
    Bar bar = { fieldOne: "S", fieldTwo: 12 };

    string|boolean|int|float? v1 = foo[f1]; //incompatible types: expected 'string', found 'fieldOne|fieldTwo|0'
    string|boolean|int|float? v2 = foo[f2]; // incompatible types: expected 'string', found '0|1'
    string|boolean|int|float? v3 = foo[f3];

    string|boolean|int|float? v4 = bar[f1]; // incompatible types: expected 'string', found 'fieldOne|fieldTwo|0'
    string|boolean|int|float? v5 = bar[f2]; // incompatible types: expected 'string', found '0|1'
    string|boolean|int|float? v6 = bar[f4]; // invalid record index expression: value space 'NoIntersection' out of range
    var v7 = bar[f5]; // invalid record index expression: value space 'NoIntersection' out of range
}
