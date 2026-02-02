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

import testorg/foo;

public function func1(foo:DummyObject1 obj, foo:Bar b, string str = foo:C1) {
}

public function func2() {
    foo:DummyObject2 obj = new;
    obj.doThatOnObject("");
    int _ = obj.id;

    string _ = foo:deprecated_func();
    future<string> _ = start foo:deprecated_func();

    foo:DummyObject1 _ = new;
    foo:DummyObject1 _ = new ();
}

@foo:deprecatedAnnotation
function testAttachingDeprecatedAnnotation() {

}

function testAccessingDeprecatedAnnotation() {
    _ = (typeof testAttachingDeprecatedAnnotation).@foo:deprecatedAnnotation;
}

function testCallingDeprecatedRemoteMethod() {
    foo:MyClientObject clientObj = client object {
        remote function remoteFunction() {

        }
    };

    clientObj->remoteFunction();
}

function testUsageOfDeprecatedParam(@deprecated foo:C2 a) {
    _ = a;
}
