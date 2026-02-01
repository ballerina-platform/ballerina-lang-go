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

function forkWithWorkers() returns int {
    fork {
        @strand{thread:"any"}
        worker forkWorker1 returns int {
            int a = 1;
            return a;
        }
        @strand{thread:"any"}
        worker forkWorker2 returns int {
            int a = 2;
            return a;
        }
    }
    @strand{thread:"any"}
    worker newWorker1 returns int {
        int a = 3;
        return a;
    }
    @strand{thread:"any"}
    worker newWorker2 returns int {
        int a = 4;
        return a;
    }

    map<int> result = wait {forkWorker1, forkWorker2, newWorker1, newWorker2};

    return result.get("forkWorker1") + result.get("forkWorker2") + result.get("newWorker1") + result.get("newWorker2");
}