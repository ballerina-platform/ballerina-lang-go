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

public function main(string... args){
    int p = 15;
    int q = 5;
    // Invoke Some random Logic.
    int r = calculateExp1(p , q);
    p = p - q;
    // Invoke another random Logic.
    string s = testCalculateExp2(p);
    s = "done";
}
function calculateExp1(int x, int y) returns (int) {
    int z = 0;
    int a = y;
    while(x >= a) {
        a = a + 1;
        if(a == 10){
            z = 100;
            break;
        }
        z = z + 10;
    }
    return z;
}
function calculateExp2(int a, int b, int c) returns (int, int) {
    int x;
    x = 10;
    int e = a;
    if (e == b) {
        e = 100;
    } else if (e == b + 1){
        e = 200;
    } else  if (e == b + 2){
        e = 300;
    }  else {
        e = 400;
    }
    int d = c;
    return (e + x, d + 1);
}
function testCalculateExp2 (int x) returns (string) {
    var (v1, v2) = calculateExp2(x, 9, 15);
    if (v1 > 200) {
        return "large";
    } else {
        return "small";
    }
}