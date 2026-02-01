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

function testVarDeclarationWithAllDeclaredSymbols () returns [int, string] {
    int a;
    string s;
    var [a, s] = unionReturnTest();
    return [a, s];
}

function unionReturnTest() returns [int, string] {
    return [5, "hello"];
}

function testVarDeclarationWithAtLeaseOneNonDeclaredSymbol () returns [int, error] {
    int a;
    var [a, err] = returnTupleForVarAssignment();
    return [a, err];
}

function returnTupleForVarAssignment() returns [int, error] {
    int a = 10;
    error er = error("");
    return [a, er];
}
