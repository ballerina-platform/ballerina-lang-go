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


service class SClass {
    string message  = "";

    remote function foo(int i) returns int {
        return i + 100;
    }

    resource function get barPath() returns string {
        return self.message;
    }

    resource function get foo/path() returns string {
        return self.message + "foo";
    }

    resource function get .() returns string {
        return self.message + "dot";
    }

    resource function get foo/baz(string s) returns string {
        return s;
    }

    resource function get foo/[int i]() returns int {
        return i;
    }

    resource function get foo/[string s]/[string... r]() returns string {
        string result = s + ", ";
        foreach string rdash in r {
            result += rdash;
        }
        return result;
    }

    resource function get foo/[json j]/[anydata a]/[boolean b]/[anydata... r]() returns string {
        return "";
    }

}
