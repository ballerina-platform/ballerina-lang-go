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

function invalidByteLiteral() {
    byte b1 = 345;
    byte b2 = -123;
    byte b3 = 256;
    byte b4 = -1;
    byte b5 = 12.3;
    byte b6 = -34.6;
    byte b7 = 0.0;
    byte[] a1 = [1, 2, -2];
    byte[] a2 = [1, 256, 45];
    byte[] a3 = [1, 2.3, 34, -56, 257];
    byte[] a4 = [1, -2.3, 3.4, 4, 46777];

    int x = 45;
    byte y = x;

    float w = 4.0;
    byte z = w;

    string r = "4";
    byte s = r;

    int x1 = -123;
    byte x2 = trap <byte> x1;

    int x3 = 256;
    byte x4 = trap <byte> x3;

    int x5 = 12345;
    byte x6 = trap <byte> x5;
}

function testUnreachableByteMatchStmt3() {
    var result = foo(2);
    int a = result is int ? 333 : (result is byte ? 777 : (result is string[] ? 666 : result));
}

function testUnreachableByteMatchStmt4() {
    var result = foo(2);
    int a = result is int ? 333 : (result is byte ? 777 : (result is string[] ? 666 : result));
}

function foo (int a) returns (byte|int|string[]) {
    if (a == 1) {
        return check trap <byte>12;
    } else if (a == 3) {
        return 267;
    }
    return ["ba", "le"];
}
