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

import ballerina/file;
import ballerina/io;

// sourcePath: a regular file holding "hello\n".
isolated function sourcePath() returns string = external;

// hardLinkPath: a hard link to the source file.
isolated function hardLinkPath() returns string = external;

// symlinkPath: a symlink whose target is the source file.
isolated function symlinkPath() returns string = external;

public function testMain() returns error? {
    string src = sourcePath();

    // copying a file onto itself is a no-op, with or without REPLACE_EXISTING
    check file:copy(src, src);
    check file:copy(src, src, file:REPLACE_EXISTING);
    io:println((check file:getMetaData(src)).size); // @output 6

    // a hard link to the source is the same file
    check file:copy(src, hardLinkPath());
    io:println((check file:getMetaData(src)).size); // @output 6

    // replacing a symlink to the source replaces the link, not the source
    string link = symlinkPath();
    check file:copy(src, link, file:REPLACE_EXISTING);
    io:println((check file:getMetaData(src)).size); // @output 6
    io:println(check file:test(link, file:IS_SYMLINK)); // @output false
}
