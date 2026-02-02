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

import testorg/foo;
import testorg/foo.records;

const ASSERTION_ERROR_REASON = "AssertionError";

type FunctionTypeDesc typedesc<function () returns (never)>;

function testTypeOfNeverReturnTypedFunction() {
    any|error expectedFunctionType = FunctionTypeDesc;

    typedesc <any|error> actualFunctionType = typeof foo:sigma;

    if (actualFunctionType is typedesc<function () returns (never)>) {
        return;
    }

    string expectedValAsString =
                expectedFunctionType is error ? expectedFunctionType.toString() : expectedFunctionType.toString();
    panic error(ASSERTION_ERROR_REASON,
                message = "expected '" + expectedValAsString + "', found '" + actualFunctionType.toString() + "'");
}

function testNeverReturnTypedFunctionCall() {
    error e = trap foo:sigma();
}

function testInclusiveRecord() {
    records:VehicleWithNever vehicle = {j:0, "q":1};
}

function testExclusiveRecord() {
    records:ClosedVehicleWithNever closedVehicle = {j:0};
}

type SomePersonalTable table<records:SomePerson> key<never>;

function testNeverWithKeyLessTable() {
    SomePersonalTable personalTable = table [
        { name: "DD", age: 33},
        { name: "XX", age: 34}
    ];
}
