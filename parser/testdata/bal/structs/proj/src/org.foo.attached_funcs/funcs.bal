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



public type user record {
    int age = 0;
    string name = "";
    string address = "";
    string zipcode = "23468";
};

public function <user u> getName() returns (string) {
    return u.name;
}

function <user u> getAge() returns (int) {
    return u.age;
}



public type person record {
    int age = 0;
    string name = "";
    string address = "";
    string zipcode = "95134";
    private:
    string ssn = "";
        int id = 0;
};


public function <person p> getName() returns (string) {
    return p.name;
}

function <person p> getAge() returns (int) {
    return p.age;
}

public function <person p> getSSN() returns (string) {
    return p.ssn;
}

public function <person p> setSSN(string ssn) {
    p.ssn = ssn;
}

public type employee record {
    int age = 0;
    string name = "";
    string address = "";
    string zipcode = "95134";
    private:
        string ssn = "";
        int id = 0;
        int employeeId = 123456;
};

public function <employee e> getName() returns (string) {
    return e.name;
}

function <employee e> getAge() returns (int) {
    return e.age;
}

public function <employee e> getSSN() returns (string) {
    return e.ssn;
}

public function <employee e> setSSN(string ssn) {
    e.ssn = ssn;
}

public function <employee e> getEmployeeId() returns (int) {
    return e.employeeId;
}

