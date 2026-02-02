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

Person p = new;

function testGetDefaultValuesInObjectGlobalVar() returns [int, string, int, string] {
    return [p.age, p.emp.name, p.foo.key, p.bar.address];
}

function testGetDefaultValuesInObject() returns [int, string, int, string] {
    Person p1 = new;
    return [p1.age, p1.emp.name, p1.foo.key, p1.bar.address];
}

class Person {
    public int age = 0;
    public string name = "";
    public Employee emp = new;
    public Foo foo = new;
    public Bar bar = new;
}

class Employee {
    public int age = 0;
    public string name = "";

    isolated function init (int age = 6, string key = "abc") {
        self.age = age;
        self.name = "sample value";
    }
}

class Foo {
    public int key = 0;
    public string value = "";

    isolated function init () {
    }
}

class Bar {
    public string address = "";
}
