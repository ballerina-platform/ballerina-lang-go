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

class Person {
    public int age = 0;

    function init (int age) {
        self.age = age;
    }
}

class Employee {
    public int age = 0;

    function init (int age, int addVal) {
        self.age = age + addVal;
    }
}

function returnDifferentObectInit() returns Person {
    return new (5, 7);
}

function returnDifferentObectInit1() returns Person | () {
    return new (5);
}

function returnDifferentObectInit2() {
    Person | () person = new (5);
    var person1 = new (5);
}

class Student {
    string name = "";
}

function returnDifferentObectInit3() {
    Student? student = new("student-1");
}

function returnDifferentObectInit4() {
    Student student = new("student-2");
}
