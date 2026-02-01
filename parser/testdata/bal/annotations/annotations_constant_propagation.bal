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

import ballerina/lang.test;

listener test:MockListener pathEP = new(9090);

const string RESOURCE_PATH = "/{id}";
const string SERVICE_BASE_PATH = "/hello";
const string SERVICE_HOST = "b7a.default";
const string RESOURCE_BODY = "person";

type ser service object {
};

@test:ServiceConfig {
    basePath: SERVICE_BASE_PATH,
    host: SERVICE_HOST
}
service ser on pathEP {

    @test:ResourceConfig {
        path: RESOURCE_PATH,
        methods: ["GET", "POST"],
        body: RESOURCE_BODY
    }

     resource function get hello(string req, string person) {
        json responseJson = { "Person": person };
        test:Caller caller = new test:Caller();
        string res = caller->respond(req);
    }
}
