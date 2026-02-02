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


function test () returns int {
    Person p = new();
    return 6;
}

class Person {

    public int age = 0;

    //return param mismatch
    function test1(int a, string name) returns string {
        return 5;
    }

    // return missing: moved to object_attached_func_def_negative
    // function test2(int a, string name) returns string {
    //
    // }

    // return mismatch
    function test3(int a, string name) returns [int, string] {
        return ["a", "b"];
    }
}

class Foo {
    public int age = 0;
}

class Bar {
    public int age = 0;
}
