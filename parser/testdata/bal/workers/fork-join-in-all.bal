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

function testForkAndWaitForAll() returns int[]|error {

    int[] results = [];

    fork {
        @strand{thread:"any"}
        worker ABC_Airline returns int {
            int x = 234;
            return x;
        }

        @strand{thread:"any"}
        worker XYZ_Airline returns int {
            int x = 500;
            return x;
        }
    }
    map<int> resultRecode = wait {ABC_Airline, XYZ_Airline};
    results[0] = resultRecode["ABC_Airline"] ?: 0;
    results[1] = resultRecode["XYZ_Airline"] ?: 0;
    return results;
}

