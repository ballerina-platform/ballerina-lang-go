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


type Foo distinct error;
type Fee distinct Foo;
type Bar distinct error<record {| int code; |}>;

function testFooError() returns Foo {
    //Foo foo = Foo("error message");
    Bar b = error Bar("message");

    error foo = error Foo("error message"); // Can assign distinct type to default error.
    Foo f = foo; // Cannot assign error to distinct error type
    Fee e = error Fee("message");
    Foo k = e;

    return foo;
}
