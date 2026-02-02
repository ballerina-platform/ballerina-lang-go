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


public function testInitAbstractObject ()  {
    Person1 p1 = new Person1();
    Person2 p2 = new Person2();
}

public function testInitAbstractObjectWithNew () {
    Person1 p1 = new;
    Person2 p2 = new;
}

type Person1 object {
    public int age;
    public string name;

    int year;
    string month;
};

// Abstract object with constructor method
type Person2 object {
    public int age;
    public string name;
    int year;
    string month;

    function init(int age, string name, int year, string month);
};
