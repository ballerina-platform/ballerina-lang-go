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

int recurs1 = 0;
int recurs2 = 0;

public function chain(int chainBreak, boolean startFunc) {
    if (chainBreak == 1 && !startFunc) {
        return;
    }

    recurs1 += 10;
    chain2(chainBreak, false);
}

public function chain2(int chainBreak, boolean startFunc) {
    if (chainBreak == 2 && !startFunc) {
        return;
    }

    recurs2 += 10;
    chain3(chainBreak, false);
}

public function chain3(int chainBreak, boolean startFunc) {
    if (chainBreak == 3 && !startFunc) {
        return;
    }

    chain(chainBreak, false);
}

public function getValues() returns [int, int] {
    return [recurs1, recurs2];
}
