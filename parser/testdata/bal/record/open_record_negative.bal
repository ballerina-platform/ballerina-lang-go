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

type Person record {|
    string name = "";
    int age = 0;
    string...;
|};

function invalidRestField() {
    Person p = { name: "John", age: 20, "height": 6, "employed": false, "city": "Colombo" };
}

type PersonA record {|
    string name = "";
    int age = 0;
    anydata|json...;
|};

function ambiguousEmptyRecordForRestField() {
    PersonA p = { name:"John", "misc": {} };
}

type Pet record {
    Animal lion;
};

class Bar {
    int a = 0;
}

function testInvalidRestFieldAddition() {
    PersonA p = {};
    p["invField"] = new Bar();
}

type Baz record {|
    int a;
    anydata...;
|};

type Qux record {
    string s;
};

type MyError error;

function testErrorAdditionForInvalidRestField() {
    error e1 = error("test reason");
    MyError e2 = error("test reason 2", err = e1);
    Baz b = { a: 1 };
    b["err1"] = e1;
    b["err2"] = e2;

    Qux q = { s: "hello" };
    q["e1"] = e1;
    q["e2"] = e2;
}

function testAnydataOrErrorRestFieldRHSAccess() {
    Person p = {};
    anydata|error name = p?.firstName;
}

type Teacher record {
    int toJson = 44;
};

function testLangFuncOnRecord() returns json{
    Teacher p = {};
    json toJsonResult = p.toValue();
    return toJsonResult;
}
