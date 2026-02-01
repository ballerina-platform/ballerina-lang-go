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


function testCyclicReferenceWithDefaultable () returns int {
    Person p = new();
    Employee em = new;
    em.age = 89;
    p.emp = em;
    Employee e = <Employee>p.emp;

    return e.age;
}

class Person {
    public int age = 0;
    public Employee? emp = ();
}

class Employee {
    public int age = 0;
    public Foo? foo = ();
    public Bar? bar = ();
}

class Foo {
    public int calc = 0;
    public Bar? bar1 = ();
}

type Bar record {
    int barVal = 0;
    string name = "";
    Person? person = ();
};
