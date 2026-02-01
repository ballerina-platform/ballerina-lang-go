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


int glbVarInt = 22343;
string glbVarString = "stringval";
float glbVarFloat = 6342.234234;
any glbVarAny = 572343;

float glbVarFloatChange = 23424.0;

public function getIntValue() returns (int) {
    return 8876;
}

public function getGlbVarInt() returns int {
    return glbVarInt;
}

public function getGlbVarString() returns string {
    return glbVarString;
}

public function getGlbVarFloat() returns float {
    return glbVarFloat;
}

public function getGlbVarFloatChange() returns float {
    return glbVarFloatChange;
}

public function setGlbVarFloatChange(float value) {
    glbVarFloatChange = value;
}
