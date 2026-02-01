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

import testorg/foo as foo;

public class NewButSameBin {
    *foo:Bin;
    public int age;
    public string name;
    public int year;
    public string month;

    public function attachFunc1(int add, string value1) returns [int, string] {
        return [2, "dummy"];
    }

    public function attachInterface(int add, string value1) returns [int, string] {
        return [2, "dummy"];
    }

    public function init(int age, string name, int year, string month) {
        self.age = age;
        self.name = name;
        self.year = year;
        self.month = month;
    }
}

public function testSimpleObjectOverridingSimilarObject () {
    NewButSameBin p = new (20, "Bob", 2020, "march");
    assertEquality(20, p.age);
    assertEquality("Bob", p.name);
    assertEquality(2020, p.year);
    assertEquality("march", p.month);
}

public class DustBinOverridingBin {
    public int age = 20;
    public string name = "sample name";
    public int year = 50;
    public string month = "february";

    *foo:Bin;

    public function init (int year, int count, string name = "sample value1", string val1 = "default value") {
        self.year = year;
        self.name = name;
        self.age = self.age + count + 50;
        self.month = val1 + " uuuu";
    }

    public function attachFunc1(int add, string value1) returns [int, string] {
        int count = self.age + add;
        string val2 = value1 + self.month;
        return [count, val2];
    }

    public function attachInterface(int add, string value1) returns [int, string] {
        int count = self.age + add;
        string val2 = value1 + self.month;
        return [count, val2];
    }
}

public function testObjectOverrideInterfaceWithInterface() {
    foo:Bin p = new DustBinOverridingBin(100, 10, val1 = "adding");
    assertEquality(80, p.age);
    assertEquality("sample value1", p.name);
    assertEquality(100, p.year);
    assertEquality("adding uuuu", p.month);
}

public function testObjectWithOverriddenFieldsAndMethods() {
    foo:ManagingDirector p = new foo:ManagingDirector(20, "John");
    assertEquality(2000000.0, p.getBonus(1, 1));
    assertEquality("Hello Director, John", p.getName());
}

readonly service class FooClass {
    *foo:FooObj;

    isolated remote function execute(string aVar, int bVar) {
    }
}

const ASSERTION_ERROR_REASON = "AssertionError";

function assertEquality(any|error expected, any|error actual) {
    if expected is anydata && actual is anydata && expected == actual {
        return;
    }

    if expected === actual {
        return;
    }

    string expectedValAsString = expected is error ? expected.toString() : expected.toString();
    string actualValAsString = actual is error ? actual.toString() : actual.toString();
    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expectedValAsString + "', found '" + actualValAsString + "'");
}
