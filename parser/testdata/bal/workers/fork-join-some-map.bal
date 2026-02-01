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

function testForkJoinReturnAnyType()(int, string) {
    return testForkJoinReturnAnyTypeVM();
}

function testForkJoinReturnAnyTypeVM()(int, string) {
    int p;
    string q;
    string r;
    float t;
    fork {
        @strand{thread:"any"}
        worker W1 {
            println("Worker W1 started");
            int x = 23;
            string a = "aaaaa";
            x, a -> fork;
        }
        @strand{thread:"any"}
        worker W2 {
            println("Worker W2 started");
            string s = "test";
            float u = 10.23;
            s, u -> fork;
        }
    } join (some 1) (map<any> results) {
        println("Within join block");
        //any[] t1;
        //t1,_ = (any[]) results["W1"];
        println("After any array cast");
        p = (int) results[0];
        //println("After int cast");
        //q = (string) results[1][0];
        //r = (string) results[0][1];
        //t = (float) results[1][1];
        println(p);
        //println(r);
        //println(q);
        //println(t);
    }

    println("After the fork-join statement");
    p = 111;
    q = "eeee";
    return p, q;
}
