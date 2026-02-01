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

function testMultiInteractions(int k) returns int{
    return test(k);
}

function test (int k) returns int {
    @strand{thread:"any"}
    worker mainW returns int {
        int x = 1100;
        x -> w1;
        x = <- w3;
        return x;
    }

    @strand{thread:"any"}
    worker w1 {
        int x = 0;
        x = <- mainW;
        x = x + 1;
        x -> w2;
    }

    @strand{thread:"any"}
    worker w2 {
        int x = 0;
        x = <- w1;
        x = x + 1;
        x -> w3;
    }

    @strand{thread:"any"}
    worker w3 {
        int x = 0;
        x = <- w2;
        x = x + 1;
        x -> mainW;
    }

    return wait mainW;
}
