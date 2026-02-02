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


public function main(string... args) {
    var a = testObjectWithInterface();
}

public function testObjectWithInterface () returns (int, string) {
    Person p = new Person();
    p.init(6);
    Person p1 = new Person();
    p1.init("");
    Person p2 = new Person();
    p2.init(true);
    Foo foo = {count:5, last:" last"};
    Person p3 = new Person();
    p3.init(foo);
    return (p.attachInterface(7), p.month);
}


class Person {

    public int age;
    public string name;


    string month = "february";


    function init () {
        self.name = "llll";
        self.age  = 6;
    }

    function _init_(int | string | boolean | Foo inVal) {
        if inVal is int {

            self.age = self.age + inVal;
        }
        else if inVal is string {
            self.name = self.name + inVal;
        }
        else if inVal is boolean {
            self.age = self.age + 10;
            self.name = self.name + " hello";
        }
        else {
            self.age = self.age + inVal.count;
            self.name = self.name + inVal.last;
        }

    }

    function incrementCount(int increment) returns int {
        int retVal = self.age;
        if (self.age > increment) {
            retVal = retVal + increment;
        }
        return retVal;
    }

    function attachInterface(int add) returns int;
}


function Person.attachInterface(int add) returns int {
    int count = self.age + add;
    count  = count + self.incrementCount(100);
    return count;
}

type Foo record {
    int count;
    string last;
};