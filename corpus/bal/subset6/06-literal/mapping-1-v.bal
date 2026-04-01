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

type FloatRecord record {|
    float x;
    float y;
|};

type DecimalRecord record {|
    decimal amount;
    decimal tax;
|};

type MixedRecord record {|
    int count;
    float ratio;
    decimal price;
|};

public function main() {
    // int literals widened to float in record context
    FloatRecord fr = {x: 1, y: 2};
    io:println(fr);

    // int literals widened to decimal in record context
    DecimalRecord dr = {amount: 100, tax: 15};
    io:println(dr);

    // int literals widened to matching field types
    MixedRecord mr = {count: 10, ratio: 2, price: 50};
    io:println(mr);

    // float literal widened to decimal in record context
    DecimalRecord dr2 = {amount: 1.5, tax: 0.5};
    io:println(dr2);

    // map with float values - int literals widened
    map<float> mf = {a: 1, b: 2};
    io:println(mf);

    // map with decimal values - int literals widened
    map<decimal> md = {a: 1, b: 2};
    io:println(md);
}
