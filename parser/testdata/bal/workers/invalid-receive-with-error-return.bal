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

import ballerina/jballerina.java;


type TrxError error<TrxErrorData>;

type TrxErrorData record {|
    string message?;
    error cause?;
    string data = "";
|};

public function main() {
    worker w1 {
        int i = 2;
        TrxError? success = i ->> w2;
        println("w1");
    }

    worker w2 returns boolean|error {
        int j = 25;
        if (0 > 1) {
            return error("trxErr", data = "test");
        }
        j = <- w1;
        println(j);
        return true;
    }
}

public function println(any|error... values) = @java:Method {
    'class: "org.ballerinalang.test.utils.interop.Utils"
} external;
