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

type Customer record {
    int id;
    string firstName;
    readonly Address? address;
};

type Address record {
    int no;
    string street;
    string town;
    Customer? customer;
};

function testCyclesInRecords() returns string {
    Address a1 = {no: 1, street: "E", town: "A", customer: ()};
    Customer c1 = {id: 1, firstName: "G", address: a1};

    a1.customer = c1;

    table<Customer> key(address) t = table [];
    t[a1] = c1;
    return t.toString();
}
