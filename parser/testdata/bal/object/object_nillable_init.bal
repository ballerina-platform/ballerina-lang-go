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

Person? p1 = new;
Person? p2 = new ();

class Person {
    public int age = 0;
}

class Employee {

    public Person? p3 = new;
    public Person? p4 = new ();
    public Person? p5 = ();
    public Person? p6 = ();

    function init () {
        self.p5 = new;
        self.p6 = new();
    }
}

function getEmployeeInstance() returns Employee {
    return new Employee();
}

function getPersonInstances() returns [Person?, Person?, Person?, Person?, Person?, Person?] {
    Person? p3 = new;
    Person? p4 = new();
    Person? p5 = ();
    p5 = new;
    Person? p6 = ();
    p6 = new ();
    return [p1, p2, p3, p4, p5, p6];
}
