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

type Person record {
    readonly string name;
    int age;
};

type Customer record {
    readonly int id;
    readonly string name;
    string lname;
};

type CustomerTable table<Customer> key(id, name);

type PersonTable table<Person> key(name);

PersonTable tab1 = table [
    {name: "AAA", age: 31},
    {name: "AAA", age: 34}
];

CustomerTable tab2 = table [
    {id: 13 , name: "Foo", lname: "QWER"},
    {id: 13 , name: "Foo" , lname: "UYOR"}
];

type Foo record {
    readonly map<string> m;
    int age;
};

type GlobalTable2 table<Foo> key(m);

function testTableConstructExprWithDuplicateKeys() {
    GlobalTable2 _ = table [
        {m: {"AAA": "DDDD"}, age: 31},
        {m: {"AAA": "DDDD"}, age: 34}
    ];
}

const int idNum = 1;
function testVariableNameFieldAsKeyField() {
    table<record {readonly int idNum; string name;}> key (idNum) _ = table [
        {idNum, name: "Jo"},
        {idNum, name: "Chiran"},
        {idNum: 2, name: "Amy"}
    ];
}

function testVariableNameFieldAsKeyField2() {
    table<record {readonly int idNum = 2; string name;}> key (idNum) _ = table [
        {idNum, name: "Jo"},
        {idNum, name: "Chiran"},
        {idNum: 2, name: "Amy"}
    ];
}

function testVariableNameFieldAsKeyField3() {
    table<record {readonly int idNum = 2; readonly string name;}> key (idNum, name) _ = table [
        {idNum, name: "Jo"},
        {idNum, name: "Jo"},
        {idNum: 2, name: "Amy"}
    ];
}

function testVariableNameFieldAsKeyField4() {
    table<record {readonly int idNum = 2; readonly string name = "A";}> key (idNum, name) _ = table [
        {idNum, name: "Jo"},
        {idNum, name: "Jo"},
        {idNum: 2, name: "Amy"}
    ];
}

type Foo2 record {
    readonly map<string> m = {};
    int age;
};

type GlobalTable3 table<Foo2> key(m);

function testTableConstructExprWithDuplicateKeys2() {
    GlobalTable3 _ = table [
        {m: {"AAA": "DDDD"}, age: 31},
        {m: {"AAA": "DDDD"}, age: 34}
    ];
}

type Foo3 record {
    readonly map<string> m = {};
    readonly int age = 18;
};

type GlobalTable4 table<Foo3> key(m, age);

function testTableConstructExprWithDuplicateKeys3() {
    GlobalTable4 _ = table [
        {m: {"AAA": "DDDD"}, age: 11},
        {m: {"AAA": "DDDD"}, age: 11}
    ];
}
