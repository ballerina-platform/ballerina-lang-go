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

import ballerina/lang.test;

# Documentation for Tst struct
# + a - `field a` documentation
# + b - `field b` documentation
# + c - `field c` documentation
type Tst record {
    string a;
    string b;
    string c;
};

# Documentation for Test annotation
annotation Test Tst;

# Documentation for Test struct
# + a - struct `field a` documentation
# + b - struct `field b` documentation
# + c - struct `field c` documentation
type Test record {
    int a;
    int b;
    int c;
};

type Person record {
    string firstName;
    string lastName;
    int age;
    string city;
};

type Employee record {
    string name;
    int age;
    string address;
    any ageAny;
};

# PizzaService HTTP Service
service /PizzaService on new test:MockListener(9090){

    # Check orderPizza resource.
    # + conn - connection.
    # + req - In request.
    resource function get orderPizza(string conn, string req) {

    }

    # Check status resource.
    # + conn - connection.
    # + req - In request.
    resource function get checkStatus(string conn, string req) {

    }
}
