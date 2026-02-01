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



function testVariableIfScope () returns (int) {
    int a = 90;
    int b = 50;
    if (a > 20) {
        int c = 20;
    } else {
        int k = 60;
        if (b < 100) {
            k = b + 30;
        }
    }
    return k;
}

function testVariableElseScope() returns (int) {
    int a = 10;
    if(a > 20) {
        a = 50;
    } else {
        int b = 30;
    }

    return b;
}

function testVariableWhileScope() {
    int a = 0;
    while( a < 5) {
        a = a + 1;
        int b = a + 20;
    }
    int sum = b;
}


service object {} myService = service object {

    resource function get res(string s) {
        string res = "abc";
        int b = a + 50;
    }

    resource function post res(string s) {
        string res = "abc";
        int c = b + 50;
    }
};
