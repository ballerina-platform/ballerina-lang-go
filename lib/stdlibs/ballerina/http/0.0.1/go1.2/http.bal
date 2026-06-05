// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the
// License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

// -- Type aliases ----------------------------------------------------------

public type Protocol "SSL"|"TLS"|"DTLS";
public type CertValidationType "OCSP_CRL"|"OCSP_STAPLING";
public type HttpVersion "1.1"|"2.0";
public type HeaderPosition "LEADING"|"TRAILING";

// -- Constants -------------------------------------------------------------

public const string LEADING = "LEADING";
public const string TRAILING = "TRAILING";

// -- Record types ----------------------------------------------------------

public type HeaderValue record {|
    string value;
    map<string> params;
|};

public type CertKey record {|
    string certFile;
    string keyFile;
    string keyPassword?;
|};

type ProtocolSpec record {|
    Protocol name;
    string[] versions;
|};

type CertValidationSpec record {|
    CertValidationType 'type;
    int cacheSize;
    int cacheValidityPeriod;
|};

public type ClientSecureSocket record {|
    boolean enable?;
    string cert?;
    CertKey key?;
    ProtocolSpec? protocol = ();
    CertValidationSpec? certValidation = ();
    string[] ciphers?;
    boolean verifyHostName?;
    boolean shareSession?;
    decimal handshakeTimeout?;
    decimal sessionTimeout?;
    string serverName?;
|};

public type FollowRedirects record {|
    boolean enabled?;
    int maxCount?;
    boolean allowAuthHeaders?;
|};

public type ClientConfiguration record {|
    decimal timeout?;
    FollowRedirects? followRedirects = ();
    HttpVersion httpVersion?;
    ClientSecureSocket? secureSocket = ();
|};

// -- Response class --------------------------------------------------------

public class Response {
    public int statusCode = 0;

    public isolated function getTextPayload() returns string = external;
    public isolated function getJsonPayload() returns json|error = external;
    public isolated function getBinaryPayload() returns byte[]|error = external;
    public isolated function hasHeader(string headerName, HeaderPosition position = "LEADING") returns boolean = external;
    public isolated function getHeader(string headerName, HeaderPosition position = "LEADING") returns string|error = external;
    public isolated function getHeaders(string headerName, HeaderPosition position = "LEADING") returns string[]|error = external;
    public isolated function getHeaderNames(HeaderPosition position = "LEADING") returns string[] = external;
}

// -- Client class ----------------------------------------------------------

public isolated client class Client {
    public isolated function init(string url, ClientConfiguration config = {}) returns error? {
        return self.initExtern(url, config);
    }

    private isolated function initExtern(string url, ClientConfiguration config) returns error? = external;
    remote isolated function get(string path, map<string|string[]>? headers = ()) returns Response|error = external;
    remote isolated function post(string path, json message, map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
    remote isolated function head(string path, map<string|string[]>? headers = ()) returns Response|error = external;
    remote isolated function options(string path, map<string|string[]>? headers = ()) returns Response|error = external;
    remote isolated function put(string path, json message, map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
    remote isolated function patch(string path, json message, map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
    remote isolated function delete(string path, json? message = (), map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
    remote isolated function execute(string httpVerb, string path, json message, map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
}

// -- Standalone functions --------------------------------------------------

public isolated function parseHeader(string headerValue) returns HeaderValue[]|error = external;
