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
import ballerina/lang.runtime;

isolated boolean woke = false;

function sleeper(decimal seconds) {
    runtime:sleep(seconds);
    lock {
        woke = true;
    }
}

// A duration beyond int64 nanoseconds must sleep (effectively) forever, not
// overflow into an immediate return. main exits without waiting on the sleepers.
public function main() {
    _ = start sleeper(1e300);
    _ = start sleeper(99999999999);
    runtime:sleep(0.05);
    lock {
        io:println(woke); // @output false
    }
}
