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


public function testObjectInsideObject () returns [string, string] {
    Person p = new Person();
    return [p.getNameWrapperInside1(), p.getNameFromDiffObject()];
}


function testGetValueFromPassedSelf() returns string {
    Person p = new Person();
    return p.selfAsValue();
}

class Person {
    public int age = 10;
    public string name = "sample name";

    private int year = 50;
    private string month = "february";

    function getName() returns string {
        return self.name;
    }

    function getNameWrapperInside1() returns string {
        return self.getName();
    }

    function getNameFromDiffObject() returns string {
        Person p = new ();
        p.name = "changed value";
        return p.getName();
    }

    function selfAsValue() returns string {
        return passSelfAsValue(self);
    }

}

function passSelfAsValue(Person p) returns string {
    return p.getName();
}
