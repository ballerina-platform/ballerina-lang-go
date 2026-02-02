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

int gInt = 5;
string gStr = "str";
boolean gBool = true;
string gNewStr = callFunc();
json gJson = null;
Person gPerson = new;
byte gByte = 255;

public function main(string... args) {
    int x = 10;
    int z = gInt + x;
    int y = x + z;
    Foo foo = { count: 5, last: "last" };
    Person p = new;
}

type Foo record {
    int count;
    string last;
};

class Person {
    public int age = 0;
    public string name = "";
    public Person? parent = ();
    private string email = "default@abc.com";

    function init() {

    }
}

function callFunc() returns (string) {
    string newStr = "ABCDEFG";
    newStr = newStr + " HIJ";
    return newStr;
}