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

function testConstrainedMapAssignNegative() returns (map<int>) {
    map<any> testMap = {};
    return testMap;
}

function testConstrainedMapRecordLiteralNegative() returns (map<int>) {
    map<int> testMap = {index_1:1, index_2:"two"};
    return testMap;
}

function testConstrainedMapIndexBasedAssignNegative() returns (map<string>) {
    map<string> testMap = {};
    testMap["name"] = 24;
    return testMap;
}

function testConstrainedMapAssignDifferentConstraintsNegative() returns (map<int>) {
    map<string> testMap = {};
    return testMap;
}

type Person record {
    string name;
    int age;
    string address;
};

type Employee record {
    string name;
    int age;
};

function testInvalidMapPassAsArgument() returns (map<Person>) {
    map<Employee> testMap = {};
    map<Person> m = returnMap(testMap);
    return m;
}

function returnMap(map<Person> m) returns (map<Person>) {
    return m;
}

function testInvalidAnyMapPassAsArgument() returns (map<Person>) {
    map<any> testMap = {};
    map<Person> m = returnMap(testMap);
    return m;
}

function testInvalidStructEquivalentCast() returns (map<Person>) {
    map<Employee> testEMap = {};
    map<Person> testPMap = <map<Person>>testEMap;
    return testPMap;
}

function testInvalidCastAnyToConstrainedMap() returns (map<Employee>) {
    map<Employee> testMap = {};
    any m = testMap;
    map<Employee> castMap = <map<Employee>>m;
    return castMap;
}

type Student record {
    int index;
    int age;
};

type IntMap map<int>;

function testInvalidStructToConstrainedMapSafeConversion() returns (map<int>|error) {
    Student s = {index:100, age:25};
    map<int> imap = check s.cloneWithType(IntMap);
    return imap;
}

function testInvalidStructEquivalentCastCaseTwo() returns (map<Student>) {
    map<Person> testPMap = {};
    map<Student> testSMap = <map<Student>>testPMap;
    return testSMap;
}

type StudentTypedesc typedesc<Student>;

function testMapToStructConversionNegative () returns (Student|error) {
    map<string> testMap = {};
    testMap["index"] = "100";
    testMap["age"] = "63";
    return check testMap.cloneWithType(StudentTypedesc);
}
