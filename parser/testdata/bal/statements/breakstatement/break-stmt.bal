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
    while (x >= yCopy) {
        yCopy = yCopy + 1;
        if (yCopy == 10) {
            z = 100;
            break;
        } else if (yCopy > 20) {
            z = 1000;
            break;
        }
        z = z + 10;
    }
    return z;
}

function nestedBreakStmt (int x, int y) returns (int) {
    int z = 10;
    int yCopy = y;
    while (x >= yCopy) {
        yCopy = yCopy + 1;
        if (yCopy >= 10) {
            z = z + 100;
            break;
        }
        z = z + 10;
        while (yCopy < x) {
            z = z + 10;
            yCopy = yCopy + 1;
            if (z >= 40) {
                break;
            }
        }
    }
    return z;
}

string output = "";

function tracePath (string path) {
    output = output + "->" + path;
}

function testBreakWithForeach (string command) returns (string) {
    output = "start";
    foreach var i in 0 ... 5 {
        tracePath("foreach" + i.toString());
        if (i == 3 && command == "break") {
            tracePath("break");
            break;
        }
        tracePath("foreachEnd" + i.toString());
    }
    tracePath("end");
    return output;
}
