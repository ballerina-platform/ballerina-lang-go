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

int glbVarInt = 800;

public function getGlbVarInt() returns int {
    return glbVarInt;
}

string glbVarString = "value";

public function getGlbVarString() returns string {
    return glbVarString;
}

float glbVarFloat = 99.34323;

public function getGlbVarFloat() returns float {
    return glbVarFloat;
}

any glbVarAny = 88343;

public function getGlbVarAny() returns any {
    return glbVarAny;
}

float glbVarFloatChange = 99.0;

public function getGlbVarFloatChange() returns float {
    return glbVarFloatChange;
}

public function setGlbVarFloatChange(float value) {
    glbVarFloatChange = value;
}

float glbVarFloat1 = glbVarFloat;

public function getGlbVarFloat1() returns float {
    return glbVarFloat1;
}

json glbVarJson = {};

public function getGlbVarJson() returns json {
    return glbVarJson;
}

public function setGlbVarJson(json value) {
    glbVarJson = value;
}

float glbVarFloatLater = 0.0;

public function getGlbVarFloatLater() returns float {
    return glbVarFloatLater;
}

public function setGlbVarFloatLater(float value) {
    glbVarFloatLater = value;
}

byte glbByte = 234;

public function getGlbByte() returns byte {
    return glbByte;
}

byte[] glbByteArray1 = [2,3,4,67,89];

public function getGlbByteArray1() returns byte[] {
    return glbByteArray1;
}

byte[] glbByteArray2 = base64 `afcd34abcdef+dfginermkmf123w/bc234cd/1a4bdfaaFGTdaKMN8923as=`;

public function getGlbByteArray2() returns byte[] {
    return glbByteArray2;
}

byte[] glbByteArray3 = base16 `afcd34abcdef123abc234bcd1a4bdfaaabadabcd892312df`;

public function getGlbByteArray3() returns byte[] {
    return glbByteArray3;
}

