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
    io:println(1 < ()); // @output false
    io:println(() < ()); // @output false
    io:println(() <= ()); // @output true

    int[]? a = [1];
    io:println(a < ()); // @output false
    io:println(() < a); // @output false

    int[]?[] b = [[1]];
    io:println(b < ()); // @output false
    io:println(() < b); // @output false

    int[]?[] c = [()];
    io:println(b < c); // @output false

    int?[] d = [1];
    io:println(d < ()); // @output false

    int?[] e = [()];
    io:println(d < e); // @output false
    io:println(e < d); // @output false

    int[] f = [1, 2, 3];
    ()[] g = [(), (), ()];
    io:println(f < g); // @output false
    io:println(g < f); // @output false

    int?[] h = [1, 2, 3];
    ()[] i = [(), (), (), ()];
    io:println(h > i); // @output false
    io:println(i > h); // @output false

    io:println(g < i); // @output true
    io:println(i < g); // @output false

    ()[][] j = [[()], [(), ()]];
    ()[][] k = [[()], [(), (), ()]];
    io:println(j < k); // @output true
    io:println(k < j); // @output false
}
