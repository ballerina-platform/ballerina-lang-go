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


public function testObjectWithSimpleInit () returns [int, string, int, string] {
    Person p = new Person(99);
    return [p.age, p.name, p.year, p.month];
}

class Person {
    public int age = 10;
    public string name = "";

    int year = 0;
    string month = "february";

    function init (int count, int year = 50, string name = "sample value1", string val1 = "default value") {
        self.year = year;
        self.name = name;
        self.age = self.age + count;
        self.month = val1;
    }
}

final int classI = 111222;

class ModuleVariableReferencingClass {
    int i = classI;
}

function value(int k = classI) returns int {
    return k;
}

ModuleVariableReferencingClass c1 = new;

function testClassWithModuleLevelVarAsDefaultValue() {
    ModuleVariableReferencingClass c = new;
    assertEquality(111222, c.i);
    assertEquality(111222, c1.i);
}

function testObjectConstructorWithModuleLevelVarAsDefaultValue() {
    var value = object {
        int i = classI;
    };
    assertEquality(111222, value.i);
}

const ASSERTION_ERROR_REASON = "AssertionError";

function assertEquality(anydata expected, anydata actual) {
    if expected == actual {
        return;
    }

    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expected.toString() + "', found '" + actual.toString () + "'");
}
