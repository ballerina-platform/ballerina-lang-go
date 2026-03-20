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

import ballerina/io;

public function main() {
    // int literal widened to float and decimal in tuple context
    [int, float, decimal] t1 = [1, 2, 3];
    io:println(t1);

    // float literal widened to decimal in tuple context
    [float, decimal] t2 = [1.5, 1.5];
    io:println(t2);

    // int array - literals stay as int
    int[] a1 = [1, 2, 3];
    io:println(a1);

    // float array - int literals widened to float
    float[] a2 = [1, 2, 3];
    io:println(a2);

    // decimal array - int literals widened to decimal
    decimal[] a3 = [1, 2, 3];
    io:println(a3);
}
