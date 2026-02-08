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

listener test:MockListener echoEP = new(9090);
listener test:MockListener echoEP2 = new(9091);

# PizzaService HTTP Service
service /PizzaService on echoEP {

    # Check orderPizza resource.
    # + conn - HTTP connection.
    # + req - In request.
    resource function get orderPizza(string conn, string req) {
    }

    # Check status resource.
    # + conn - HTTP connection.
    # + req - In request.
    resource function get checkStatus(string conn, string req) {

    }
}

# Test type `typeDef`
# Test service `helloWorld`
# Test variable `testVar`
# Test var `testVar`
# Test function `add`
# Test parameter `x`
# Test const `constant`
# Test annotation `annot`
service /PizzaService2 on echoEP2 {

    # Check orderPizza resource.
    # + conn - HTTP connection.
    # + req - In request.
    resource function get orderPizza(string conn, string req) {

    }

    # Check status resource.
    # + conn - HTTP connection.
    # + req - In request.
    resource function get checkStatus(string conn, string req) {

    }
}
