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

function testWorkerInVM () returns int {
    int q = 0;
    q = testWorker();
    return q;
}

function testWorker () returns int {

    @strand{thread:"any"}
    worker w1 returns int {
        int result = 0;
        int i = 10;
        i -> sampleWorker;
        result = <- sampleWorker;
        return result;
    }

    @strand{thread:"any"}
    worker sampleWorker {
        int r = 120;
        int i = 0;
        i = <- w1;
        r = changeMessage(i);
        r -> w1;
    }
    return wait w1;
}

function changeMessage (int i) returns int {
    return i + 10;
}