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

// Supported subset of ballerina/http for the Go runtime.
// See lib/http/client-support.md for the full feature support matrix.

// ── Shared types ─────────────────────────────────────────────────────────────

// Represents the parsed header value details.
// Fields: value - the header value; params - map of header parameters.
public type HeaderValue record {|
    string value;
    map<string> params;
|};

# Parses the header value which contains multiple values or parameters.
# ```ballerina
#  http:HeaderValue[] values = check http:parseHeader("text/plain;level=1;q=0.6, application/xml;level=2");
# ```
#
# + headerValue - The header value
# + return - An array of `http:HeaderValue` typed records containing the value and its parameter map,
#            or an `error` if the header parsing fails
public isolated function parseHeader(string headerValue) returns HeaderValue[]|error = external;

// ── TLS / secure-socket types ─────────────────────────────────────────────────

// Represents a combination of certificate and private key for mutual TLS.
// The key file must be an unencrypted PEM file; password-protected keys are not supported.
// Fields: certFile - path to the certificate chain; keyFile - path to the private key;
//         keyPassword - accepted at compile time, ignored at runtime.
public type CertKey record {|
    string certFile;
    string keyFile;
    string keyPassword?;
|};

// SSL/TLS protocol type. "SSL" and "DTLS" are accepted at compile time but have no effect
// at runtime — the Go TLS stack only supports TLS.
public type Protocol "SSL"|"TLS"|"DTLS";

// Provides configuration for TLS protocol version constraints.
// Fields: name - the SSL protocol ("SSL", "TLS", or "DTLS");
//         versions - list of enabled protocol versions (e.g., ["TLSv1.2", "TLSv1.3"]).
public type ProtocolConfig record {|
    Protocol name;
    string[] versions;
|};

// Certificate validation type. Accepted at compile time; not implemented at runtime
// (OCSP/CRL checks are not supported in Go's standard TLS library).
public type CertValidationType "OCSP_CRL"|"OCSP_STAPLING";

// Provides configuration for certificate revocation validation.
// Accepted at compile time for compatibility; not enforced at runtime.
// Fields: 'type - certificate validation type; cacheSize - maximum cache size;
//         cacheValidityPeriod - cache entry validity period.
public type CertValidation record {|
    CertValidationType 'type;
    int cacheSize;
    int cacheValidityPeriod;
|};

// Provides configurations for facilitating secure communication with a remote HTTP endpoint.
//
// Supported: enable, verifyHostName, cert (PEM path), key (CertKey), serverName,
//            ciphers, handshakeTimeout, shareSession, protocol.versions.
// Not supported: cert as crypto:TrustStore, key as crypto:KeyStore, certValidation,
//               sessionTimeout, protocol.name (Go supports TLS only).
//
// Fields:
//   enable         - Enable SSL validation. Set false to skip certificate verification.
//   cert           - PEM file path for custom CA trust store (crypto:TrustStore not supported).
//   key            - Mutual TLS certificate+key pair (crypto:KeyStore not supported).
//   protocol       - TLS protocol and version constraints.
//   certValidation - OCSP/CRL settings (accepted; not enforced at runtime).
//   ciphers        - IANA cipher suite names for TLS 1.2; unknown names silently skipped.
//   verifyHostName - Enable/disable hostname verification.
//   shareSession   - Enable/disable TLS session ticket reuse.
//   handshakeTimeout - Maximum TLS handshake duration (seconds).
//   sessionTimeout - TLS session timeout (accepted; not configurable in Go).
//   serverName     - SNI override; defaults to hostname from the target URL.
public type ClientSecureSocket record {|
    boolean enable?;
    string cert?;
    CertKey key?;
    ProtocolConfig? protocol?;
    CertValidation? certValidation?;
    string[] ciphers?;
    boolean verifyHostName?;
    boolean shareSession?;
    decimal handshakeTimeout?;
    decimal sessionTimeout?;
    string serverName?;
|};

// ── Client configuration ──────────────────────────────────────────────────────

// Provides configurations for controlling the behaviour in response to HTTP redirect responses
// (status codes 301, 302, 303, 307, 308).
// Fields: enabled - enable/disable redirect following (default: false);
//         maxCount - maximum redirects to follow (default: 5);
//         allowAuthHeaders - forward Authorization headers on redirect (default: false).
public type FollowRedirects record {|
    boolean enabled = false;
    int maxCount = 5;
    boolean allowAuthHeaders = false;
|};

// HTTP protocol version. "1.1" or "2.0" (default).
// HTTP/1.0 is not supported — Go's HTTP client cannot send HTTP/1.0 requests.
public type HttpVersion "1.1"|"2.0";

// Provides a set of configurations for controlling the behaviours when communicating with
// a remote HTTP endpoint.
//
// Supported: timeout, httpVersion, followRedirects, secureSocket.
// Not supported: circuitBreaker, retryConfig, cookieConfig, cache, compression,
//               auth, http1Settings, http2Settings, responseLimits, socketConfig,
//               validation, laxDataBinding.
//
// Fields:
//   timeout        - Max wait time in seconds before request times out (default: 30).
//   followRedirects - Redirect handling configuration; () disables redirect following.
//   httpVersion    - HTTP protocol version: "1.1" or "2.0" (default: "2.0").
//   secureSocket   - TLS settings; () uses default TLS verification.
public type ClientConfiguration record {|
    decimal timeout = 30;
    FollowRedirects? followRedirects = ();
    HttpVersion httpVersion = "2.0";
    ClientSecureSocket? secureSocket = ();
|};

// ── Header position ───────────────────────────────────────────────────────────

// Represents the position of an HTTP header: leading (transport) or trailing (after body).
// TRAILING is accepted at compile time but all runtime operations act on transport headers.
public type HeaderPosition "LEADING"|"TRAILING";

# Represents the leading header position (transport headers). This is the default for all
# header operations.
public const HeaderPosition LEADING = "LEADING";

# Represents the trailing header position. Accepted at compile time for API compatibility;
# at runtime all header operations act on transport (leading) headers.
public const HeaderPosition TRAILING = "TRAILING";

// ── Response ──────────────────────────────────────────────────────────────────

# Represents an HTTP response received from a remote endpoint.
#
# `Response` objects are created by the HTTP client after a successful request —
# they are never constructed directly by user code. All write methods (`addHeader`,
# `setHeader`, `setPayload`, etc.) are not supported in this implementation.
public class Response {
    # The HTTP status code of the response (e.g., 200, 404, 500).
    public int statusCode = 0;

    # Returns the response body as a plain string.
    #
    # + return - The response body as a `string`
    public isolated function getTextPayload() returns string = external;

    # Parses the response body as JSON.
    #
    # + return - The parsed `json` value, or an `error` if the body is not valid JSON
    public isolated function getJsonPayload() returns json|error = external;

    # Returns the response body as a byte array.
    #
    # + return - The response body as `byte[]`, or an `error` if extraction fails
    public isolated function getBinaryPayload() returns byte[]|error = external;

    # Checks whether the specified header is present in the response.
    #
    # + headerName - The header name (case-insensitive)
    # + position - Header position (`LEADING` or `TRAILING`). `TRAILING` is accepted
    #              but all lookups operate on transport headers
    # + return - `true` if the header exists, `false` otherwise
    public isolated function hasHeader(string headerName, HeaderPosition position = LEADING) returns boolean = external;

    # Returns the first value for the specified header.
    #
    # + headerName - The header name (case-insensitive)
    # + position - Header position (`LEADING` or `TRAILING`). `TRAILING` is accepted
    #              but all lookups operate on transport headers
    # + return - The first header value, or an `error` if the header is not found
    public isolated function getHeader(string headerName, HeaderPosition position = LEADING) returns string|error = external;

    # Returns all values for the specified header.
    #
    # + headerName - The header name (case-insensitive)
    # + position - Header position (`LEADING` or `TRAILING`). `TRAILING` is accepted
    #              but all lookups operate on transport headers
    # + return - A `string[]` of all values for the header, or an `error` if the header is not found
    public isolated function getHeaders(string headerName, HeaderPosition position = LEADING) returns string[]|error = external;

    # Returns the names of all response headers.
    #
    # + position - Header position (`LEADING` or `TRAILING`). `TRAILING` is accepted
    #              but all lookups operate on transport headers
    # + return - An array of all header names present in the response
    public isolated function getHeaderNames(HeaderPosition position = LEADING) returns string[] = external;
}

// ── Client ────────────────────────────────────────────────────────────────────

# The HTTP client provides functionality to connect to remote HTTP services and perform
# requests using the standard HTTP methods GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS,
# and EXECUTE.
#
# **Supported methods:** `get`, `post`, `put`, `patch`, `delete`, `head`, `options`, `execute`.
#
# **Not supported:** `forward`, `submit`/`getResponse`, HTTP/2 server push methods
# (`hasPromise`, `getNextPromise`, etc.), and resource function syntax (`client->/path.get(...)`).
#
# **Message body:** The `message` parameter accepts `json`, which in Ballerina includes
# `string`, `byte[]`, and all JSON-compatible values. The `Content-Type` is inferred:
# - `string` → `text/plain`
# - `byte[]` → `application/octet-stream`
# - everything else → serialised as `application/json`
#
# The `mediaType` parameter overrides the inferred `Content-Type` in all cases.
#
# **Return type:** All methods return `Response|error`. Automatic data binding via
# `targetType` is not supported — use `getTextPayload()`, `getJsonPayload()`, or
# `getBinaryPayload()` to extract the response body.
#
# **Error types:** Errors are plain `error` values. The upstream distinct error types
# (`http:ClientError`, `http:HeaderNotFoundError`, etc.) are not declared.
public isolated client class Client {

    # Gets invoked to initialize the `client`. During initialization, the configurations
    # provided through the `config` record control timeout, TLS, HTTP version, and redirect
    # behaviour.
    #
    # + url - The base URL of the target service
    # + config - The configurations to be used when initializing the `client`.
    #            Unsupported fields (`circuitBreaker`, `retryConfig`, `cookieConfig`, etc.)
    #            are not available in this implementation
    # + return - `()` on success, or an `error` if initialisation fails
    public isolated function init(string url, ClientConfiguration config = {}) returns error? {
        return self.initNative(url, config);
    }

    private isolated function initNative(string url, ClientConfiguration config) returns error? = external;

    # Retrieves a representation of the specified resource from the remote HTTP endpoint.
    #
    # + path - The request path (appended to the base URL)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function get(string path, map<string|string[]>? headers = ()) returns Response|error = external;

    # Creates a new resource or submits data to a resource for processing.
    #
    # + path - The request path (appended to the base URL)
    # + message - The request body (`string`, `byte[]`, or any JSON-compatible value)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + mediaType - Optional `Content-Type` override; inferred from `message` if omitted
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function post(string path, json message, map<string|string[]>? headers = (),
            string? mediaType = ()) returns Response|error = external;

    # Creates a new resource or replaces a representation of the specified resource.
    #
    # + path - The request path (appended to the base URL)
    # + message - The request body (`string`, `byte[]`, or any JSON-compatible value)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + mediaType - Optional `Content-Type` override; inferred from `message` if omitted
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function put(string path, json message, map<string|string[]>? headers = (),
            string? mediaType = ()) returns Response|error = external;

    # Applies a partial modification to the specified resource.
    #
    # + path - The request path (appended to the base URL)
    # + message - The request body (`string`, `byte[]`, or any JSON-compatible value)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + mediaType - Optional `Content-Type` override; inferred from `message` if omitted
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function patch(string path, json message, map<string|string[]>? headers = (),
            string? mediaType = ()) returns Response|error = external;

    # Deletes the specified resource.
    #
    # + path - The request path (appended to the base URL)
    # + message - Optional request body (`string`, `byte[]`, or any JSON-compatible value)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + mediaType - Optional `Content-Type` override; inferred from `message` if omitted
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function delete(string path, json? message = (), map<string|string[]>? headers = (),
            string? mediaType = ()) returns Response|error = external;

    # Requests headers from the specified resource without fetching the response body.
    # Identical to `get` but the server must not return a message body.
    #
    # + path - The request path (appended to the base URL)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function head(string path, map<string|string[]>? headers = ()) returns Response|error = external;

    # Requests the communication options available for the specified resource.
    #
    # + path - The request path (appended to the base URL)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function options(string path, map<string|string[]>? headers = ()) returns Response|error = external;

    # Sends an HTTP request with an explicit verb to the specified path.
    # Use this for HTTP methods not covered by the dedicated remote functions.
    #
    # + httpVerb - The HTTP method to use (e.g., `"GET"`, `"POST"`, `"PUT"`)
    # + path - The request path (appended to the base URL)
    # + message - The request body (`string`, `byte[]`, or any JSON-compatible value)
    # + headers - Optional request headers as a `map<string|string[]>`
    # + mediaType - Optional `Content-Type` override; inferred from `message` if omitted
    # + return - The `http:Response` or an `error` if the request fails
    remote isolated function execute(string httpVerb, string path, json message,
            map<string|string[]>? headers = (), string? mediaType = ()) returns Response|error = external;
}
