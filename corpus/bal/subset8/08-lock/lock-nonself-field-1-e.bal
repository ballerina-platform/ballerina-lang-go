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

// Field access on a non-`self` base inside a lock: the LHS of an
// assignment to a target defined outside the lock must be a plain
// variable name. `b.val = ...` is rejected even though `b.val` has
// an isolated (readonly) static type.
class Box {
    int val = 0;
}

function bump() {
    Box b = new;
    lock {
        // JBallerina don't treat is as an error but according to spec "an assignment to a variable 
        // defined outside the lock statement is allowed only if left-hand side is just a 
        // variable name and the right hand side is an isolated expression"
        b.val = 5; // @error
    }
}

public function main() {
    bump();
}
