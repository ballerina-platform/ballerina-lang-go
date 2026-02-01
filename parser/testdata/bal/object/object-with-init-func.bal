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
    Person p = new Person(99, 7);
    return [p.age, p.name, p.year, p.month];
}

public function testObjectWithSimpleInitWithDiffValues () returns [int, string, int, string] {
    Person p = new Person(675, 27, val1 = "adding value in invocation");
    return [p.age, p.name, p.year, p.month];
}

public function testObjectWithoutRHSType () returns [int, string, int, string] {
    Person p = new (675, 27, val1 = "adding value in invocation");
    return [p.age, p.name, p.year, p.month];
}


class Person {
    public int age = 10;
    public string name = "sample name";

    int year = 50;
    string month = "february";

    function init (int year, int count, string name = "sample value1", string val1 = "default value") {
        self.year = year;
        self.name = name;
        self.age = self.age + count;
        self.month = val1;
    }
}
