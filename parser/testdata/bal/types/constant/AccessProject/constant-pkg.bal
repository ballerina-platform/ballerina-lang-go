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

import constant_types.variable;

final int constNegativeInt = -342;

final int constNegativeIntWithSpace = -     88;

final float constNegativeFloat = -88.2;

final float constNegativeFloatWithSpace = -      3343.88;

float glbVarFloat = variable:getConstFloat();

function accessConstantFromOtherPkg() returns (float) {
    return variable:getConstFloat();
}

function assignConstFromOtherPkgToGlobalVar() returns (float) {
    return glbVarFloat;
}

function getNegativeConstants() returns [int, int, float, float] {
    return [constNegativeInt, constNegativeIntWithSpace, constNegativeFloat, constNegativeFloatWithSpace];
}


final float a = 4.0;

function floatIntConversion() returns [float, float, float]{
    float[] f = [1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0];
    return [a, f[5], 10.0];
}

function accessPublicConstantFromOtherPackage() returns string {
    return variable:name;
}

function accessPublicConstantTypeFromOtherPackage() returns variable:AB {
    return variable:A;
}

type AB "A";

function testTypeAssignment() returns AB {
    AB ab = variable:A; // This is valid because this is equivalant to `AB ab = "A";`.
    return ab;
}
