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

type Person record {
    string name;
    int age;
    Person? parent;
    json info;
    map<any> address;
    int[] marks;
};


type Student record {
    string name;
    int age;
    map<any> address;
    int[] marks;
};

function testStructToStruct() returns (Person) {
    Student s = { name:"Supun", 
                  age:25, 
                  address:{"city":"Kandy", "country":"SriLanka"}, 
                  marks:[24, 81]
                };
    Person p = s;
    return p;
}

function intToFloatImpCast() {
    float[] numbers;
    int a = 999;
    int b = 95;
    int c = 889;
    numbers = [a, b, c];
    float val1 = 160.0;
    float val2 = 160;
    int d = 0;
    float val3 = d;
}
