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

// Skipped on Windows (see test_util.WindowsUnsupportedTests): `echo` is a shell
// builtin there, not an executable os:exec can spawn.
import ballerina/os;
import ballerina/io;

public function main() returns error? {
    // Run `echo`, wait for it to exit, and read its captured stdout.
    os:Process p = check os:exec({value: "echo", arguments: ["hello"]});
    int code = check p.waitForExit();
    io:println(code);             // @output 0
    byte[] out = check p.output();
    io:println(out.length() > 0); // @output true

    // The process's stderr stream is empty for a successful echo.
    os:Process p2 = check os:exec({value: "echo", arguments: ["world"]});
    int _ = check p2.waitForExit();
    byte[] err = check p2.output(io:stderr);
    io:println(err.length());     // @output 0

    // exec of a non-existent command returns an error.
    os:Process|os:Error bad = os:exec({value: "no_such_command_xyz_123"});
    io:println(bad is os:Error);  // @output true

    // exit() terminates a process.
    os:Process p3 = check os:exec({value: "echo", arguments: ["bye"]});
    p3.exit();
    io:println("exited");         // @output exited
}
