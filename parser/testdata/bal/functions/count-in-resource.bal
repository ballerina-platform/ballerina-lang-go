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

import ballerina/http;

listener http:MockListener testEP = new(9090);

@http:ServiceConfig {
    basePath:"/test"
}
service TestService on testEP {

    @http:ResourceConfig {
        methods:["GET"],
        path:"/resource"
    }
    resource function testResource (http:Caller caller, http:Request req) {
        json[] jsonArray = [];
        string[] strArray = ["foo", "bar"];
        foreach var s in strArray {
            jsonArray[jsonArray.length()] = s;
        }

        http:Response res = new;
        res.setJsonPayload(<@untainted> jsonArray);
        checkpanic caller->respond(res);
    }
}
