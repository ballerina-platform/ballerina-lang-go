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

int glbVarInt = 800;
string glbVarString = "value";
float glbVarFloat = 99.34323;
any glbVarAny = 88343;

float glbVarFloatChange = 99.0;

float glbVarFloat1 = glbVarFloat;

listener http:MockListener echoEP = new(9090);

@http:ServiceConfig {basePath:"/globalvar"}
service GlobalVar on echoEP {

    @http:ResourceConfig {
        methods:["GET"],
        path:"/defined"
    }
    resource function defineGlobalVar (http:Caller conn, http:Request request) {
        http:Response res = new;
        json responseJson = {"glbVarInt":glbVarInt, "glbVarString":glbVarString, "glbVarFloat":glbVarFloat};
        res.setJsonPayload(responseJson);
        checkpanic conn->respond(res);
    }

    @http:ResourceConfig {
        methods:["GET"],
        path:"/access-service-level"
    }
    resource function accessGlobalVarAtServiceLevel (http:Caller conn, http:Request request) {
        http:Response res = new;
        json responseJson = {"glbVarFloat":glbVarFloat.toString()};
        res.setJsonPayload(responseJson);
        checkpanic conn->respond(res);
    }

    @http:ResourceConfig {
        methods:["GET"],
        path:"/change-resource-level"
    }
    resource function changeGlobalVarAtResourceLevel (http:Caller conn, http:Request request) {
        http:Response res = new;
        glbVarFloatChange = 77.87;
        json responseJson = {"glbVarFloatChange":glbVarFloatChange};
        res.setJsonPayload(responseJson);
        checkpanic conn->respond(res);
    }

    @http:ResourceConfig {
        methods:["GET"],
        path:"/get-changed-resource-level"
    }
    resource function getChangedGlobalVarAtResourceLevel (http:Caller conn, http:Request request) {
        http:Response res = new;
        json responseJson = {"glbVarFloatChange":glbVarFloatChange};
        res.setJsonPayload(responseJson);
        checkpanic conn->respond(res);
    }

}

@http:ServiceConfig {basePath:"/globalvar-second"}
service GlobalVarSecond on echoEP {

    @http:ResourceConfig {
        methods:["GET"],
        path:"/get-changed-resource-level"
    }
    resource function getChangedGlobalVarAtResourceLevel (http:Caller conn, http:Request request) {
        http:Response res = new;
        json responseJson = {"glbVarFloatChange":glbVarFloatChange};
        res.setJsonPayload(responseJson);
        checkpanic conn->respond(res);
    }

}

