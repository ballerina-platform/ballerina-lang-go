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

function testUninitializedClosureVars() {
    string a;

    var bazz = function () {
        string _ = a + "aa";
    };

    bazz();

    string b;
    int count;

    var bar = function () {
        count += 10;

        b = b + "bb";
    };

    bar();
}

public function lambdaInitTest() {
    int localVar;
    int localVar2;

    worker w1 {
        localVar = 4;
        int _ = localVar;
        int _ = localVar2;
    }

    var _ = function () {
        localVar = 2;
        int _ = localVar;
        int _ = localVar2;
    };

    int _ = localVar;
    int _ = localVar2;
}
