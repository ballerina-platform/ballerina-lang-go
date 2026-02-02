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

public function testSimpleObjectAsStruct () {
    Person p = new Person();
    assertEquality(p.age, 10);
    assertEquality(p.name, "sample name");
    assertEquality(p.year, 50);
    assertEquality(p.month, "february");
}

public function testSimpleObjectAsStructWithNew() {
    Person p = new;
    assertEquality(p.age, 10);
    assertEquality(p.name, "sample name");
    assertEquality(p.year, 50);
    assertEquality(p.month, "february");
}

public function testTypeRefInClass() {
    Student s = new Student(10, "sample name", 50, "feb");
    assertEquality(s.age, 10);
    assertEquality(s.name, "sample name");
    assertEquality(s.year, 50);
    assertEquality(s.month, "feb");
    assertEquality(s.grade, 1);

    var x = new Stu(0);
    assertEquality(x.i, 0);
}

public function testUsingClassValueAsRecordField() {
    Rec r = {p: new(), s: new(1, "Guido", 1970, "Feb")};
    assertEquality(r.p.name, "sample name");
    assertEquality(r.s.name, "Guido");
    assertEquality(r.p.age, 10);
    assertEquality(r.s.age, 1);
}

type Rec record {
    Person p;
    Student s;
};

class Person {
    public int age = 10;
    public string name = "sample name";
    int year = 50;
    string month = "february";
}

class Student {
    *Person;
    int grade = 1;

    function init(int age, string name, int year, string month = "Jan")  {
        self.age = age;
        self.name = name;
        self.year = year;
        self.month = month;
    }
}

class Stu {
    *Per;

    function init(int i) {
        self.i = i;
    }
}

class Per {
    int i = 0;
}

function assertEquality(any|error expected, any|error actual) {
    if expected is anydata && actual is anydata && expected == actual {
        return;
    }

    if expected === actual {
        return;
    }

    string expectedValAsString = expected is error ? expected.toString() : expected.toString();
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error("AssertionError",
            message = "expected '" + expectedValAsString + "', found '" + actualValAsString + "'");
}
