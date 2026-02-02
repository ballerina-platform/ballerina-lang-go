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


function testForkReturnAnyType() returns [int, string]|error {
    int p = 0;
    string q = "";
    string r;
    float t;

    fork {
        @strand{thread:"any"}
        worker W1 returns [int, string] {
            int x = 23;
            string a = "aaaaa";
            return [x, a];
        }
        @strand{thread:"any"}
        worker W2 returns [string, float] {
            string s = "test";
            float u = 10.23;
            return [s, u];
        }
    }
    map<any> results = wait {W1, W2};
    any t1 = results["W1"];
    if t1 is [int,string] {
        [p, q] = t1;
    }
    any t2 = results["W2"];
    if t2 is [string,float] {
        [r, t] = t2;
    }
    return [p, q];
}
