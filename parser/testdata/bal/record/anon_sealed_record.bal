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


function testAnonStructAsFuncParam() returns int {
    return testAnonStructFunc(10, {k:14, s:"sameera"});
}

function testAnonStructFunc(int i, record {int k = 10; string s;} anonSt) returns (int) {
    return anonSt.k + i;
}


function testAnonStructAsLocalVar() returns int {
    record {| int k = 11; string s; |} anonSt = {};

    return anonSt.k;
}


record {| string fname; string lname; int age; |} person;

function testAnonStructAsPkgVar() returns string {

    person = {fname:"sameera", lname:"jaya"};
    person.lname = person.lname + "soma";
    person.age = 100;
    return person.fname + ":" + person.lname + ":" + person.age;
}

type employee record {
    string fname;
    string lname;
    int age;
    record {|
        string line01;
        string line02;
        string city;
        string state;
        string zipcode;
    |} address;

    record {|
        string month = "JAN";
        string day = "01";
        string year = "1970";
    |} dateOfBirth;
};

function testAnonStructAsStructField() returns string {

    employee e = {fname:"sam", lname:"json", age:100,
                     address:{line01:"12 Gemba St APT 134", city:"Los Altos", state:"CA", zipcode:"95123"},
                    dateOfBirth:{}};
    return e.dateOfBirth.month + ":" + e.address.line01 + ":" + e.address["state"] + ":" + e.fname;
}