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

function test () returns int {
    Person p = new(8, new (9));
    return p.emp.age;
}

class Person {
    public int age = 0;
    public Employee emp;

    function init (int age, Employee emp) {
        self.age = age;
        self.emp = emp;
    }
}

class Employee {
    public int age = 0;
    public Foo foo;
    public Bar bar = {};

    function init (int age) {
        self.age = age;
    }
}

class Foo {
    public int calc;

    function init (int calc) {
        self.calc = calc;
    }
}

type Bar record {
    int barVal = 0;
    string name = "";
};
