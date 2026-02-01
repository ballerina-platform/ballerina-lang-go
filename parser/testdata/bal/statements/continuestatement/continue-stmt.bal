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

function calculateExp1 (int x, int y) returns (int) {
    int z = 0;
    int yCopy = y;
    while (x > yCopy) {
        yCopy = yCopy + 1;
        if (yCopy == 10) {
            continue;
        }
        z = z + 1;
    }
    return z;
}

function nestedContinueStmt (int x, int y) returns (int) {
    int z = 0;
    int yCopy = y;
    while (x >= yCopy) {
        yCopy = yCopy + 1;
        int i = 0;
        if (yCopy == 10) {
            continue;
        }
        while (i < yCopy) {
            i = i + 1;
            if (i == 10) {
                continue;
            }
            z = z + i;
        }
    }
    return z;
}

string output = "";

function tracePath (string path) {
    output = output + "->" + path;
}
