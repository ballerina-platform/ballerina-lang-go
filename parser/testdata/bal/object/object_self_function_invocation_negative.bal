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

class Person1 {
    string name;
    int age;

    public function init(string name, int age) {
        self.modify(self);
        self.name = name;
        self.age = age;
    }

    function modify(Person1 person) {
        person.name = "New Name";
    }
}

public function testSelfKeywordInvocationNegative() {
    Person1 _ = new("person1", 32);
}

class Person2 {
    string name;
    string city;

    public function init(string name, string city) {
        int i = 0;
        while (i < 5) {
            if (i/2 == 2) {
                continue;
            }
            self.modify(self);
        }
        self.name = name;
        self.city = city;
    }

    function modify(Person2 person) {
        person.name = "New Name";
    }
}

public function testSelfKeywordInvocationWithLoopAndContinueNegative() {
    Person2 _ = new("person2", "city");
}

class Person3 {
    string name;
    string city;

    public function init(string name, string city) {
        int i = 0;
        while (i < 5) {
            if (i/2 == 2) {
                break;
            }
            self.modify(self);
        }
        self.name = name;
        self.city = city;
    }

    function modify(Person3 person) {
        person.name = "New Name";
    }
}

public function testSelfKeywordInvocationWithLoopAndBreakNegative() {
    Person3 _ = new("person2", "city");
}


class Person4 {
    string name;
    string city;

    public function init(string name) {
        if (name == "") {
            self.name = "tom";
          } else if (name == "tom") {
            self.city = "London";
          } else {
            self.name = "tom";
            self.city = "London";
          }
        self.modify(self);
        self.name = name;
        self.city = "city";
    }

    function modify(Person4 person) {
        person.name = "New Name";
    }
}

public function testSelfKeywordInvocationWithBranchNegative() {
    Person4 _ = new("");
}

class Person5 {
    string name;
    string city;

    public function init(string[] names, string city) {
        foreach string n in names {
            if (n == "person1") {
                self.modify(self);
            }
            self.name = n;
        }
        self.city = city;
    }

    function modify(Person5 person) {
        person.name = "New Name";
    }
}

public function testSelfKeywordInvocationWithForEachNegative() {
    string[] persons = ["person1", "person2", "person3"];
    Person5 _ = new(persons, "city");
}

class Person6 {
    string name;
    string city;

    public function init(string[] names, string city) {
        foreach string n in names {
            if (n == "person1") {
                self.modify("person1");
            }
            self.name = n;
        }
        self.city = city;
    }

    function modify(string name) {
        self.name = name;
    }
}

public function testSelfKeywordInvocationWithInvocationArg() {
    string[] persons = ["person1", "person2", "person3"];
    Person5 _ = new(persons, "city");
}


class Person7 {
    string name;

    public function init(string name) {
        modify(self);
        self.name = name;
    }
}

function modify(Person7 person) {

}

public function testSelfKeywordInvocationWithModuleLevelFunctionInvocation() {
    Person7 _ = new("person");
}

class Person8 {
    string name;

    public function init(string name) {
        change(self);
        self.name = name;
    }
}

function change(Person8... person) {
}

public function testSelfKeywordInvocationWithModuleLevelFunctionInvocationWithRestArgs() {
    Person8 _ = new("person");
}
