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

function scopeIfValue(int a, int b, int c) returns int {
    int k = 0;
    if(a > b) {
        k = k + c;
    } else {
        k = 1;
        if (c > b) {
            k = k + a;
        } else {
            k = k + 99999;
        }
    }

    return k;
}

function scopeWhileScope (int a, int b, int c) returns int {
    int k = 5;
    int b1 = b;
    while (a > b1 ) {
        b1 = b1 + 1;
        int i  = 10;

        if (c < a) {
            i = i + 10;
        } else {
            i = i + 20;
        }
        k = k + i;
    }
    return k;
}