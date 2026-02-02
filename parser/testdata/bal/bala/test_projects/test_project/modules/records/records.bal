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



public type Corge record {|
   string a = "hello";
   string b = "world";
   string[] c = <readonly> ["x", "y"];
|};

public type Config record {|
    readonly readonly & Interceptor[] interceptors = [];
    Foo foo;
|};

public type Interceptor distinct service object {
    isolated remote function execute() returns anydata|error;
};

public type Foo readonly & object {
   public isolated function fooFunc() returns string;
};

public type Foo1 record {|
    any[] x = [1, 2];
|};

public type Foo2 record {|
    *Foo1;
    any[] y = ["abc", "xyz"];
|};
