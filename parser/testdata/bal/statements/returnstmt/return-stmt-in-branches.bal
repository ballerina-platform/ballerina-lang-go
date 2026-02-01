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

function returnStmtBranch1(int value, int b) returns (int) {

    if( value > 10) {
        return 100;

    } else if ( value == 10) {
        return 200;

    } else {

        if (b > 10) {
            return 300;

        } else if ( b == 10){
            return 400;
        }

        return 500;
    }
}

function returnStmtBranch2(int value, int b) returns (int) {
    int a = value + b;
    int c = returnStmtBranch1(9, 10);
    c = c + c + a;

    return c;
}

