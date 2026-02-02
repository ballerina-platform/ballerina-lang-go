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


public class ParentFoo {

    public int i;
    public ChildFoo c;
    private string s = "";

    function init (int i, ChildFoo c){
        self.i = i;
        self.c = c;
    }
}

class ChildFoo {
    private string name = "";

    function init (string name) {
        self.name = name;
    }
}

class PrivatePerson {

    public int age = 0;
    public string name = "";

    function init (int age, string name){
        self.age = age;
        self.name = name;
    }
    public function getPrivatePersonName() returns string { return self.name; }
}

public function newPrivatePerson() returns (PrivatePerson) {
    return new PrivatePerson(12, "mad");
}

public function privatePersonAsParam(PrivatePerson p) returns (string){
    return p.name;
}

public function privatePersonAsParamAndReturn(PrivatePerson p) returns (PrivatePerson) {
    return p;
}
