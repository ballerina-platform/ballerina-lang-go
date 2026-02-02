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

class Person {
    public string name = "default first name";
    public string lname = "";
    public map<any> adrs = {};
    public int age = 999;

    function init (string name, map<any> adrs, int age) {
        self.name = name;
        self.age = age;
        self.adrs = adrs;
    }
}

class ObjectField {
    public string key = "";

    function init (string key) {
        self.key = key;
    }
}

function testExpressionAsStructIndex () returns (string) {
    ObjectField nameField = new ("name");
    Person emp = new ("Jack", {"country":"USA", "state":"CA"}, 25);
    return emp[nameField.key];
}
