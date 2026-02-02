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

function test1() {
    int a = 0xFFFFFFFFFFFFFFFF;
    int b = 9999999999999999999;

    int d = -0xFFFFFFFFFFFFFFFF;
    int e = -9999999999999999999;
}

function test2() {
    // This is to test preceding syntax issues.
    int x = 1

    int g = 0672;
    int h = 0912;

    // This is to verify that the compilation continues.
    int y = 1
}

type testType1 int:Signed32;
type testType2 int:Signed16;
type testType3 int:Signed8;

function testStaticTypeOfUnaryExpr() {
    int:Signed8 a = -128;
    int:Signed8 _ = -a;
    int:Signed16 b = -32768;
    int:Signed16 _ = -b;
    int:Signed32 c = -2147483648;
    int:Signed32 _ = -c;

    testType1 d = -2147483648;
    int:Signed32 _ = -d;
    testType2 e = -32768;
    int:Signed16 _ = -e;
    testType3 f = -128;
    int:Signed8 _ = -f;
}
