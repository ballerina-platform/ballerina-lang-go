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

function ifWithWhileComplex(int arg1, int arg2) returns int {
    int a = 10;
    int b = 30;
    int x = 100;
    int y = 10;

    if a == b  {
        a = b;
        if (a != b) {
            a = a * b;
        } else if (a < b) {
            while y < 100 {
                y = y + 1;
            }
            return a;
        }
        while x < 1000 {
            x = x + 1;
        }
    } else if (a > b) {
        a = a - b;
    } else {
        a = a / b;
        return a + b;
    }

    return a + b;
}

function whileWithIf() {
    int x = 100;
    int y = 10;

    while x < 1000 {
        x = x + y;

        if (y < 100) {
            y = y + 1;
        }

        y = y * 1;

        //while y < 100 {
        //    y = y + 1;
        //}
    }

    y = x * y;
}

function whileSample() {
    int x = 100;
    int y = 10;

    while x < 1000 {
        x = x + y;

        while y < 100 {
            y = y + 1;
        }
    }

    y = x * y;
}

function whileSimple() {
    int x = 100;
    int y = 10;

    while x < 1000 {
        x = x + y;
        y = y - 1;
    }

    y = x * y;
}

function ifSimple() returns int {
    int x = 100;
    int y = 10;

    if (x > y) {
        y = x + 1;
        return y;
    } else {
        x = y + 1;
    }

    return x * y;
}


function dump(int arg1, int arg2) returns int {
    int a = 10;
    int b = 30;

    if a == b  {
        a = b;
        if (a != b) {
            a = a * b;
        } else if (a < b) {
            return a;
        }
    } else if (a > b) {
        a = a - b;
    } else {
        a = a / b;
        return a + b;
    }

    return a + b;
}


public function main(int arg) returns error? {
    int a = 10;
    boolean b = a > 100;
}

function genComplex(int arg1, int arg2) returns int {
    int a = 10;
    int b = a + arg1;
    int c = a - b + arg2;
    b = b + c;
    return a + b;
}

type Employee record {
    int id;
    string name;
};

function newTable() returns table<Employee> {
    table<Employee> employeeTable = table [
            {id: 1, name: "Employee1"},
            {id: 2, name: "Employee2"}
        ];

    return employeeTable;
}

function getDecimal() returns decimal {
    decimal d = 10;
    return d;
}
