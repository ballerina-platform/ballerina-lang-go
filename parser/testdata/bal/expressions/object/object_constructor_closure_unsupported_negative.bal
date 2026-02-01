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

function test1() returns any {
    final int i = 10;           // mapBlock2

    object {    // zz = 2
            object {
                // int x;
                function bar(int b) returns int;
             } yy;

             function bar(int b) returns int;

        } zz = object {
            object { // yy = 1
                // int x;
                function bar(int b) returns int;

            } yy = object {
                int x = 5 + i;
                function bar(int b) returns int { // call me baby 1
                    return b + self.x;
                }
            };

            function bar(int b) returns int {   // call me baby 2
                return b + self.yy.bar(b);
            }
        };

        return zz.bar(4);
}

function test2() returns any {
    final int i = 10;           // mapBlock2

    object {    // zz = 2
            object {
                // int x;
                function bar(int b) returns int;
             } yy;

             function bar(int b) returns int;

        } zz = object {
            object { // yy = 1
                // int x;
                function bar(int b) returns int;

            } yy = object {
                int x = i;
                function bar(int b) returns int { // call me baby 1
                    return b + self.x;
                }
            };

            function bar(int b) returns int {   // call me baby 2
                return b + self.yy.bar(b);
            }
        };

        return zz.bar(4);
}

function testClosuresWithObjectConstrExprUnsupported(int b1) {
    final int a1 = 10;

    var _ = object {
        int a2 = 10;
        object {
            int a3;
            function foo(int b2) returns int;
        } obj2 = object {
            int a3 = a1; // not supported
            function foo(int b2) returns int {
                //return self.a3 + a1; // should also give not supported error
                return self.a3;
            }
        };
    };
}
