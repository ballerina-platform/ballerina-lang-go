# `ballerina/http` Client — Support Reference

_Version_: v0.05.0 \
_Created_: 2026/05/14 \
_Updated_: 2026/05/14 

## Interface Definition

```ballerina
// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ballerina/http — supported subset of the HTTP client API.

// ── Types ────────────────────────────────────────────────────────────────────

public type Protocol "SSL"|"TLS"|"DTLS";

public type CertValidationType "OCSP_CRL"|"OCSP_STAPLING";

public type HttpVersion "1.1"|"2.0";

public type HeaderPosition "LEADING"|"TRAILING";

public final HeaderPosition LEADING = "LEADING";
public final HeaderPosition TRAILING = "TRAILING";

public type CertKey record {|
    string certFile;
    string keyFile;
    string keyPassword?;
|};

public type ClientSecureSocket record {|
    boolean enable?;
    string cert?;
    CertKey key?;
    record {| Protocol name; string[] versions; |} protocol?;
    record {| CertValidationType 'type; int cacheSize; int cacheValidityPeriod; |} certValidation?;
    string[] ciphers?;
    boolean verifyHostName?;
    boolean shareSession?;
    decimal handshakeTimeout?;
    decimal sessionTimeout?;
    string serverName?;
|};

public type FollowRedirects record {|
    boolean enabled = false;
    int maxCount = 5;
    boolean allowAuthHeaders = false;
|};

public type ClientConfiguration record {|
    decimal timeout = 30;
    FollowRedirects? followRedirects = ();
    HttpVersion httpVersion = "2.0";
    ClientSecureSocket? secureSocket = ();
|};

// ── Response ─────────────────────────────────────────────────────────────────

public class Response {
    public int statusCode;

    public isolated function getTextPayload() returns string = external;

    public isolated function getJsonPayload() returns json|error = external;

    public isolated function getBinaryPayload() returns byte[]|error = external;

    public isolated function hasHeader(string headerName,
                                       HeaderPosition position = LEADING) returns boolean = external;

    public isolated function getHeader(string headerName,
                                       HeaderPosition position = LEADING) returns string|error = external;

    public isolated function getHeaders(string headerName,
                                        HeaderPosition position = LEADING) returns string[]|error = external;

    public isolated function getHeaderNames(HeaderPosition position = LEADING) returns string[] = external;
}

// ── Client ───────────────────────────────────────────────────────────────────

public client class Client {

    public isolated function init(string url,
                                  ClientConfiguration config = {}) returns error? = external;

    remote isolated function get(string path,
                                 map<string|string[]>? headers = ()) returns Response|error = external;

    remote isolated function post(string path,
                                  json message,
                                  map<string|string[]>? headers = (),
                                  string? mediaType = ()) returns Response|error = external;

    remote isolated function put(string path,
                                 json message,
                                 map<string|string[]>? headers = (),
                                 string? mediaType = ()) returns Response|error = external;

    remote isolated function patch(string path,
                                   json message,
                                   map<string|string[]>? headers = (),
                                   string? mediaType = ()) returns Response|error = external;

    remote isolated function delete(string path,
                                    json? message = (),
                                    map<string|string[]>? headers = (),
                                    string? mediaType = ()) returns Response|error = external;

    remote isolated function head(string path,
                                  map<string|string[]>? headers = ()) returns Response|error = external;

    remote isolated function options(string path,
                                     map<string|string[]>? headers = ()) returns Response|error = external;

    remote isolated function execute(string httpVerb,
                                     string path,
                                     json message,
                                     map<string|string[]>? headers = (),
                                     string? mediaType = ()) returns Response|error = external;
}

// ── Module-level functions ────────────────────────────────────────────────────

public type HeaderValue record {|
    string value;
    map<string> params;
|};

public isolated function parseHeader(string headerValue) returns HeaderValue[]|error = external;
```

---

## Supported

### Client — remote methods

All eight remote methods are supported: `get`, `post`, `put`, `patch`, `delete`, `head`, `options`, and `execute`. Each method accepts optional request headers as a `map<string|string[]>`. Methods that carry a body (`post`, `put`, `patch`, `delete`, `execute`) additionally accept an optional media type override.

### Client — initialisation

The client can be initialised with a URL and an optional `ClientConfiguration` record. The configuration supports:

| Field | Notes |
|---|---|
| `timeout` | Request timeout as a decimal (seconds); default is `30` |
| `httpVersion` | `"1.1"` or `"2.0"` (default). HTTP/2 is enabled over both TLS (via ALPN) and cleartext (h2c) |
| `followRedirects` | Full `FollowRedirects` record: `enabled`, `maxCount` (default 5), `allowAuthHeaders` |
| `secureSocket` | See TLS section below |

### Request message body

The `message` parameter is typed as `json`, which in Ballerina includes `string`, `byte[]`, and all JSON-compatible values. The runtime infers `Content-Type` from the value:

- `string` — sent as `text/plain`
- `byte[]` (a list where every element is an integer in 0–255) — sent as `application/octet-stream`
- All other `json`-compatible values (`nil`, `boolean`, `int`, `float`, `decimal`, nested maps and lists, JSON arrays) — serialised and sent as `application/json`

The `mediaType` parameter overrides the inferred `Content-Type` in all cases.

### Response — fields

| Field | Notes |
|---|---|
| `statusCode` | HTTP status code of the response |

### Response — methods

| Method | Notes |
|---|---|
| `getTextPayload()` | Returns the response body as a string |
| `getJsonPayload()` | Parses the body as JSON; returns `json\|error` |
| `getBinaryPayload()` | Returns the body as a byte array; returns `byte[]\|error` |
| `hasHeader(name, position?)` | Returns `true` if the header is present |
| `getHeader(name, position?)` | Returns the first value for the header, or an error if absent |
| `getHeaders(name, position?)` | Returns all values for the header, or an error if absent |
| `getHeaderNames(position?)` | Returns the names of all response headers |

The `position` parameter accepts `http:LEADING` (default) or `http:TRAILING`. Trailing headers are accepted at compile time but not modelled at runtime — all operations act on transport headers.

### TLS (`secureSocket`)

| Setting | Notes |
|---|---|
| `enable` / `verifyHostName` | Disabling either turns off certificate/hostname verification |
| `cert` (string path) | Custom CA trust store from a PEM file; CN-based hostname fallback supported for legacy self-signed certificates |
| `key` (`CertKey`) | Mutual TLS using `certFile` and `keyFile` (unencrypted PEM) |
| `serverName` | Overrides the SNI hostname sent during the TLS handshake |
| `ciphers` | IANA cipher suite names applied to TLS 1.2 connections; unknown names are silently skipped |
| `handshakeTimeout` | Maximum duration allowed for the TLS handshake |
| `shareSession = false` | Disables TLS session ticket reuse |
| `protocol.versions` | Accepts `"TLSv1.0"`, `"TLSv1.1"`, `"TLSv1.2"`, `"TLSv1.3"` to set minimum and maximum TLS versions |

---

## Not Supported

### Client — remote methods

| Method | Reason |
|---|---|
| `forward` | Requires an `http:Request` object, which is not implemented |
| `submit` / `getResponse` | Asynchronous request model (`HttpFuture`) is not implemented |
| `hasPromise` / `getNextPromise` / `getPromisedResponse` / `rejectPromise` | HTTP/2 server push is not implemented |

Resource function syntax (`client->/path.get(...)`) is not supported; use the remote method form instead.

### Client — configuration

| Field | Reason |
|---|---|
| `circuitBreaker` | Circuit breaker pattern not implemented |
| `retryConfig` | Automatic retry not implemented |
| `cookieConfig` | Cookie store not implemented |
| `cache` | HTTP response caching not implemented |
| `compression` | Compression negotiation not implemented |
| `auth` | Auth handlers not implemented |
| `http1Settings` | HTTP/1.x-specific settings (keep-alive, chunking, proxy) not implemented |
| `http2Settings` | HTTP/2-specific settings (prior knowledge, window size) not implemented |
| `responseLimits` | Response size limits not implemented |
| `socketConfig` | TCP socket configuration not implemented |
| `validation` / `laxDataBinding` | Payload validation not implemented |

`httpVersion: "1.0"` is a compile error — Go's HTTP client does not support sending HTTP/1.0 requests.

### Response — fields

`reasonPhrase`, `resolvedRequestedURI`, `server`, and `cacheControl` are not exposed. The raw `headers` map is not exposed; use the header methods instead.

### Response — methods

All write methods (`addHeader`, `setHeader`, `removeHeader`, `removeAllHeaders`, `setJsonPayload`, `setPayload`, etc.) are not supported — `Response` objects are only received from the server, never constructed by user code.

| Method | Reason |
|---|---|
| `getXmlPayload` | XML values are not representable at runtime |
| `getByteStream` | Streaming response body not implemented |
| `getSseEventStream` | Server-Sent Events not implemented |
| `getBodyParts` | Multipart (`mime:Entity[]`) not implemented |
| `getContentType` / `setContentType` | Not exposed |
| `getEntity` / `setEntity` | MIME entity access not exposed |
| `getStatusCodeRecord` | Status code response type binding not implemented |
| Cookie methods | Cookie handling not implemented |

### Response data binding (`targetType`)

The `targetType` parameter present on upstream methods (which enables automatic binding of the response body to `string`, `byte[]`, `json`, custom record types, or `stream<SseEvent, error?>`) is not supported. All methods return `Response|error` and the caller must extract the payload explicitly using `getTextPayload()`, `getJsonPayload()`, or `getBinaryPayload()`.

### Request message body types

| Type | Reason |
|---|---|
| `http:Request` objects | `Request` class is not implemented |
| `xml` | XML values are not representable at runtime |
| `stream<byte[], io:Error?>` | Streaming request body not implemented |
| `mime:Entity[]` | Multipart not implemented; requires `ballerina/mime` |

### TLS (`secureSocket`)

| Setting | Reason |
|---|---|
| `cert` as `crypto:TrustStore` | Requires `ballerina/crypto`, which is not implemented |
| `key` as `crypto:KeyStore` | Requires `ballerina/crypto`, which is not implemented |
| `keyPassword` in `CertKey` | Password-protected private keys are not supported; the key file must be unencrypted PEM |
| `certValidation` | OCSP/CRL certificate revocation checks are not supported in Go's standard TLS library |
| `sessionTimeout` | Not configurable in Go's TLS stack |
| `protocol.name` | Go only supports TLS; `"SSL"` and `"DTLS"` are accepted at compile time but have no effect at runtime |

### Error types

Errors returned by client methods are plain `error` values. The upstream distinct error types (`http:ClientError`, `http:HeaderNotFoundError`, etc.) are not declared — type narrowing with `is http:ClientError` or `is http:HeaderNotFoundError` will not work.
