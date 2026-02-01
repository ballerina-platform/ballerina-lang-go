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

public type CompensationStatus record {
    boolean scopeACompensated;
    boolean scopeBCompensated;
};

public type R record {
    boolean compensated;
};

public function main(string... argv) {
    var r = testNestedScopes();
    io:println(r);
}

function testNestedScopes() returns (CompensationStatus) {
    int a = 2;
    CompensationStatus st = {scopeACompensated:false, scopeBCompensated:false};
    CompensationStatus st2 = {scopeACompensated:false, scopeBCompensated:false};


    scope scopeA {
        scope scopeB {
            string s = "abc";
            a = 5;
        } compensation {
            st.scopeACompensated = true;
        }
    } compensation {
        st.scopeBCompensated = true;
    }
    a = 3;
    compensate scopeA;
    return st;
}

public function OtherFunc(boolean comp) returns (R) {
    int b = 0;
    string s = "s";
    R status = { compensated: false };
    scope scp {

    } compensation {
        s = s + "2";
        status.compensated = true;
    }
    
    if (comp) {
        compensate scp;
    }
    return status;
}

function testLoopScopes() returns (int) {
    int a = 2;

    int k = 1;
    while (k < 5) {
        scope ScopeA {
            scope ScopeB {
                a = 5;
            } compensation {
                string abc = "abc";
            }
        } compensation {

        }
        k = k +2;

        compensate ScopeA;
    }
    a = 3;

    return a;
}

