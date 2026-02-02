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

int[] glbArray = [1, 2];

public function getGlbArray() returns int[] {
    return glbArray;
}

int[3] glbSealedArray = [0, 0, 0];

public function getGlbSealedArray() returns int[3] {
    return glbSealedArray;
}

int[*] glbSealedArray2 = [1, 2, 3, 4];

public function getGlbSealedArray2() returns int[4] {
    return glbSealedArray2;
}

int[2][3] glbSealed2DArray = [[0, 0, 0], [0, 0, 0]];

public function getGlbSealed2DArray() returns int[2][3] {
    return glbSealed2DArray;
}

int[*][2] glbSealed2DArray2 = [[1, 2], [3, 4], [5, 6]];

public function getGlbSealed2DArray2() returns int[3][2] {
    return glbSealed2DArray2;
}
