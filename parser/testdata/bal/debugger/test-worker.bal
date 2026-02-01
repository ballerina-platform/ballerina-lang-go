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

import ballerina/io;
public function main(string... args) {
    simpleWorkers();
    io:println("worker run finished");
}

function simpleWorkers() {
    worker w1 {
        int p = 15;
        int q = 5;
        // Invoke Some random Logic.
        int a = calculateExp1(p , q);
        io:println("worker 1 - " + a);
        a -> w2;
        a = <- w2;
    }
    worker w2 {
        int a = 0;
        int p = 15;
        int q = 5;
        // Invoke Some random Logic.
        int b = calculateExp3(p , q);
        io:println("worker 2 - " + b);
        a = <- w1;
        b -> w1;
    }
}

function calculateExp1(int x, int y) returns (int) {
    int z = 0;
    int a = y;
    while(x >= a) {
        a = a + 1;
        if(a == 10) {
            z = 100;
            break;
        }
        z = z + 10;
    }
    return z;
}

function calculateExp3(int x, int y) returns (int) {
    int z = 0;
    int a = y;
    while(x >= a) {
        a = a + 1;
        if(a == 10) {
            z = 100;
            break;
        }
        z = z + 10;
    }
    return z;
}