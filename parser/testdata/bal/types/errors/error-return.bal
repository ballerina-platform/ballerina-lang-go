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

public type InvalidNameError error<record { string companyName; }>;

function getQuote(string name) returns (float|InvalidNameError) {

    if (name == "FOO") {
        return 10.5;
    } else if (name == "BAR") {
        return 11.5;
    }

    InvalidNameError err = error InvalidNameError("invalid name", companyName = name);
    return err;
}

function testReturnError() returns [string, string, string, string]|error {

    string a;
    string b;
    string c;
    string d;

    float quoteValue;
    // Special identifier "=?" will be used to ignore values.

    quoteValue = check getQuote("FOO");
    a = "FOO:" + quoteValue.toString();

    // Ignore error.
    var r = getQuote("QUX");

    if (r is float) {
        b = "QUX:" + r.toString();
    } else {
        b = "QUX:ERROR";
    }

    // testing for errors.
    // error occurred. Recover from the error by assigning 0.
    var q = getQuote("BAZ");

    if (q is float) {
        c = "BAZ:" + quoteValue.toString();
    } else {
        quoteValue = 0.0;
        c = "BAZ:" + quoteValue.toString();
    }

    var p = getQuote("BAR");

    if (p is float) {
        d = "BAR:" + p.toString();
    } else {
        quoteValue = 0.0;
        d = "BAR:ERROR";
    }

    return [a,b,c,d];
}

function testReturnAndThrowError() returns (string){
    error? e = trap checkAndThrow();

    if (e is error) {
        return e.message();
    }

    return "OK";
}

function checkAndThrow(){
    float qVal;
    var p = getQuote("BAZ");

    if (p is float) {
        qVal = p;
    } else {
        panic p;
    }
}
