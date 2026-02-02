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
    int x = 1;
    string s = "ten";
    int i = 10;
    int j = 20;
    string a = x == 1 ? s : i;
    any b = x ? s : i;
    any c = true ? s : i;
}


function test1 (int value) {
    string s1 = value > 40 ? true : false ? "morethan40" : "lessthan20";
}

function test2() {
    int[] a = [];
    boolean trueVal = true;
    boolean falseVal = false;

    a.push(trueVal ? 1 : "");

    byte[] b = [];
    b.push(falseVal ? 0 : 256);
}
