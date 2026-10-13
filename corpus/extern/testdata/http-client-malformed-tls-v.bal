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
import ballerina/io;

// Malformed TLS material must fail Client.init rather than fall back to the
// system trust store or connect without the requested client certificate.
public function main() {
    http:Client|error caCert = new ("http://testserver", {
        secureSocket: {cert: "testdata/certs/malformed.pem"}
    });
    io:println(caCert is error); // @output true

    http:Client|error clientKey = new ("http://testserver", {
        secureSocket: {key: {certFile: "testdata/certs/malformed.pem", keyFile: "testdata/certs/malformed.pem"}}
    });
    io:println(clientKey is error); // @output true
}
