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

import ballerina/jballerina.java;

function testPrintingObject() returns Person {
	Person p = new(20, "John", 2018, "March");
	println(p);
	return p;
}

class Person {
    public int age = 10;
    public string name = "sample name";

    int year = 50;
    string month = "February";

    function init (int age, string name, int year, string month) {
        self.age = age;
        self.name = name;
        self.year = year;
        self.month = month;
    }

    function toString() returns string {
        json result = {age: self.age, name: self.name};
        return result.toJsonString();
    }
}

public function println(any|error... values) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
