# Ballerina HTTP Library

## Overview

This module provides the HTTP client and listener APIs for building and consuming HTTP services. The full jBallerina `http` module covers two sides:

**Client** — response data binding to custom types and status code records, authentication (Basic, Bearer, JWT, OAuth2), resiliency patterns (circuit breaker, retry, failover, load balancer), cookie management, HTTP response caching, compression negotiation, connection pooling, async requests, HTTP/2 server push, Server-Sent Events, multipart payloads, and streaming I/O.

**Service / Listener** — an HTTP listener with configurable host, TLS, HTTP version, and request limits; service definition with path-based routing and resource function dispatch; automatic binding of path parameters, query parameters, headers, and payloads in resource signatures; caller-based response dispatch; request/response interceptor pipeline; service-level and resource-level annotations (`@http:ServiceConfig`, `@http:ResourceConfig`, `@http:Payload`, `@http:Header`, `@http:Query`, `@http:Cache`); CORS configuration; listener authentication and authorization (File user store, LDAP, JWT, OAuth2); status code response types from resources; and SSE streaming responses.

The Go Native Interpreter currently supports the **HTTP client subset only**: the eight core remote methods, TLS/mTLS (PEM-based), redirect following, and manual payload extraction from responses. The service/listener side is not yet implemented.

## Key Functionalities

- Send HTTP requests using GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, and custom verbs.
- Configure request timeout, HTTP version (1.1 / 2.0), and redirect behaviour.
- Secure connections with TLS and mutual TLS using PEM certificate and key files.
- Set custom request headers and override the inferred Content-Type.
- Read the response status code, text, JSON, or binary payload.
- Inspect response headers by name or enumerate all header names.
- Parse structured header values (value + parameter map) with the header parsing utility.

## Examples

```ballerina
import ballerina/http;
import ballerina/io;

public function main() returns error? {
    // Plain HTTP client with a 10-second timeout
    http:Client client = check new ("http://httpbin.org", {timeout: 10});

    // GET request
    http:Response getResp = check client->get("/get");
    io:println("Status: ", getResp.statusCode);
    json body = check getResp.getJsonPayload();
    io:println("Body: ", body);

    // POST request with a JSON payload
    json payload = {name: "Alice", age: 30};
    http:Response postResp = check client->post("/post", payload);
    io:println("POST status: ", postResp.statusCode);

    // TLS client with a custom CA certificate
    http:Client secureClient = check new ("https://example.com", {
        secureSocket: {
            cert: "/path/to/ca.pem"
        }
    });

    http:Response secureResp = check secureClient->get("/");
    io:println("Secure status: ", secureResp.statusCode);
}
```

## Go Native Interpreter Support Status

This library is currently being migrated to Go to support the Ballerina Native Interpreter. The tables below outline the current support level for various features of this library in the Go implementation.

Support Levels:

- **Supported**: Fully implemented and tested in the Go version.
- **Partially Supported**: Implemented but lacking some edge cases, options, or sub-features. (See comments).
- **Not Yet Supported**: Planned for migration, but not yet implemented.
- **Cannot Support**: Cannot be implemented in the Go version due to technical limitations or architectural differences. (See comments).

### Client

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| Core HTTP request methods | Supported | `get`, `post`, `put`, `patch`, `delete`, `head`, `options` are all implemented. |
| Custom HTTP verb execution | Supported | `execute` accepts any HTTP verb string. |
| Request timeout | Supported | Configured via `timeout` in `ClientConfiguration` (decimal seconds, default 30). |
| HTTP version selection | Supported | `"1.1"` and `"2.0"` (default) are supported; includes cleartext HTTP/2 (h2c). `"1.0"` is a compile error. |
| Redirect following | Supported | Full `FollowRedirects` record supported: `enabled`, `maxCount` (default 5), `allowAuthHeaders`. |
| Custom request headers | Supported | Accepted as `map<string\|string[]>` on every method. |
| Content-Type inference from payload type | Supported | `string` → `text/plain`, `byte[]` → `application/octet-stream`, all other `json`-compatible values → `application/json`. |
| Media type override | Supported | `mediaType` parameter on body-carrying methods overrides the inferred Content-Type. |
| TLS and mutual TLS (mTLS) | Partially Supported | PEM-file-based CA trust (`cert` as a string path) and client certificate/key pairs (`key` as `CertKey`) are supported. `crypto:TrustStore`, `crypto:KeyStore`, password-protected private keys (`keyPassword`), OCSP/CRL certificate revocation (`certValidation`), and TLS session timeout (`sessionTimeout`) are not supported. |
| Client-side response data binding | Not Yet Supported | The `targetType` parameter is absent; callers must extract the payload explicitly via `getTextPayload`, `getJsonPayload`, or `getBinaryPayload`. Binding to custom record types, `xml`, or `stream<SseEvent, error?>` is not available. |
| Status code response binding | Not Yet Supported | `StatusCodeClient` and `getStatusCodeRecord()` are not implemented. |
| Client authentication | Not Yet Supported | The `auth` field in `ClientConfiguration` is absent. BasicAuth (`CredentialsConfig`), BearerToken, self-signed JWT (`JwtIssuerConfig`), and all OAuth2 grant types are not supported. |
| Circuit breaker | Not Yet Supported | `circuitBreaker` configuration and `CircuitBreakerClient` are not implemented. |
| Automatic retry | Not Yet Supported | `retryConfig` configuration and `RetryClient` are not implemented. |
| Failover client | Not Yet Supported | `FailoverClient` is not implemented. |
| Load balancer client | Not Yet Supported | `LoadBalanceClient` is not implemented. |
| Cookie management | Not Yet Supported | `cookieConfig`, `CookieStore`, and `getCookieStore()` are not implemented. |
| HTTP response caching | Not Yet Supported | The `cache` (`CacheConfig`) configuration is not implemented. |
| Compression negotiation | Not Yet Supported | The `compression` configuration (`accept-encoding` header negotiation) is not implemented. |
| HTTP/1.x protocol settings | Not Yet Supported | `http1Settings` (keep-alive, chunking, proxy) is not implemented. |
| HTTP/2 protocol settings | Not Yet Supported | `http2Settings` (prior knowledge, initial window size) is not implemented. |
| Connection pooling | Not Yet Supported | `poolConfig` (`PoolConfiguration`) is not implemented. |
| Response size limits | Not Yet Supported | `responseLimits` (`ResponseLimitConfigs`) is not implemented. |
| TCP socket configuration | Not Yet Supported | `socketConfig` (`ClientSocketConfig`) is not implemented. |
| Client-side payload validation | Not Yet Supported | The `validation` and `laxDataBinding` flags in `ClientConfiguration` are not implemented. |
| Proxy support | Not Yet Supported | `ProxyConfig` (available via `http1Settings.proxy` in jBallerina) is not implemented. |
| Request forwarding via incoming request | Not Yet Supported | The `forward` remote method requires an `http:Request` object, which is not yet implemented. |
| Async request submission | Not Yet Supported | `submit`, `getResponse`, and `HttpFuture` are not implemented. |
| HTTP/2 server push | Not Yet Supported | `hasPromise`, `getNextPromise`, `getPromisedResponse`, and `rejectPromise` are not implemented. |
| Resource function call syntax | Not Yet Supported | The `client->/path.get(...)` path-template syntax is not supported; use the remote method form instead. |

### Request

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| Request object | Not Yet Supported | The `Request` class (raw path, method, HTTP version, headers, payload read methods, query params, mutual SSL handshake info) is not implemented. It cannot be used as a client message body or accessed in resource function signatures. |
| Path parameter binding | Not Yet Supported | Automatic extraction of URL path segments into resource function parameters is not implemented. |
| Query parameter binding | Not Yet Supported | Automatic binding of URL query parameters to resource function parameters is not implemented. |
| Inbound header binding | Not Yet Supported | Automatic binding of request headers to resource function parameters via `@http:Header` is not implemented. |
| Inbound payload binding | Not Yet Supported | Automatic deserialization of the request body into typed resource function parameters via `@http:Payload` is not implemented. |
| Multipart and form-data payload | Not Yet Supported | `mime:Entity[]` as a request body type and the associated `getBodyParts()` response method are not implemented. |
| Streaming request body | Not Yet Supported | `stream<byte[], io:Error?>` as a request payload type is not implemented. |

### Response

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| Response status code access | Supported | Exposed as the `statusCode` field on `Response`. |
| Response payload as text | Supported | `getTextPayload()` returns the body as a `string`. |
| Response payload as JSON | Supported | `getJsonPayload()` parses the body and returns `json\|error`. |
| Response payload as raw bytes | Supported | `getBinaryPayload()` returns `byte[]\|error`. |
| Response header inspection | Supported | `hasHeader`, `getHeader`, `getHeaders`, and `getHeaderNames` operate on transport (leading) headers. Trailing header position is accepted at compile time but has no runtime effect. |
| Response write methods | Not Yet Supported | `setPayload`, `setJsonPayload`, `setHeader`, `addHeader`, `removeHeader`, `setStatusCode`, etc. are not implemented; `Response` is read-only in the current client-only implementation. |
| Streaming response body | Not Yet Supported | `getByteStream()` is not implemented. |
| Server-Sent Events | Not Yet Supported | `getSseEventStream()` and consuming a `stream<SseEvent, error?>` response are not implemented. |
| Response XML payload | Not Yet Supported | The `xml` type and related payload handling methods (`getXmlPayload()`, `setXmlPayload()`) are not implemented due to the lack of XML support in the Go runtime. |

### Listener

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| HTTP Listener | Not Yet Supported | The `Listener` class (start, graceful stop, attach, detach) is not implemented; no server-side listener can be created. |
| Listener configuration | Not Yet Supported | `ListenerConfiguration` (host, timeout, HTTP version, HTTP/1.x settings, HTTP/2 window size, graceful stop timeout, request limits, server name, socket config) is not implemented. |
| Listener TLS / mTLS | Not Yet Supported | `ListenerSecureSocket` (server certificate, mutual TLS, protocol, ciphers, etc.) is not implemented. |
| Default listener | Not Yet Supported | The module-level default listener (`http:defaultListener`) is not implemented. |
| Listener authentication and authorization | Not Yet Supported | `ListenerAuthConfig` and listener-side auth handlers (file user store, LDAP, JWT, OAuth2) are not implemented. |

### Service

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| HTTP service definition and routing | Not Yet Supported | Declaring `service on listener` with path-based routing is not implemented. |
| Resource function dispatch | Not Yet Supported | Resource functions with path parameters, accessor methods, and typed response returns are not implemented. |
| Caller-based response dispatch | Not Yet Supported | The `Caller` class and its `respond()` method for sending responses back to the client are not implemented. |
| Status code response types from resources | Not Yet Supported | Returning `http:Ok`, `http:Created`, `http:NotFound`, and other `StatusCodeResponse` subtypes from resource functions is not implemented. |
| Service-level annotation | Not Yet Supported | `@http:ServiceConfig` (host, compression, chunking, CORS, auth, validation, lax data binding) is not implemented. |
| Resource-level annotation | Not Yet Supported | `@http:ResourceConfig` (name, consumes, produces, CORS, auth, linked resources) is not implemented. |
| Response cache annotation | Not Yet Supported | `@http:Cache` on resource return types is not implemented. |
| CORS configuration | Not Yet Supported | Cross-origin resource sharing configuration at service and resource level is not implemented. |
| Request and response interceptors | Not Yet Supported | `RequestInterceptor`, `ResponseInterceptor`, `RequestErrorInterceptor`, `ResponseErrorInterceptor`, and the `InterceptableService` type are not implemented. |
| Request context | Not Yet Supported | `RequestContext` for passing data through the interceptor pipeline is not implemented. |
| Service contract type | Not Yet Supported | `ServiceContract` type for contract-first service definitions is not implemented. |
| Service-level compression and chunking | Not Yet Supported | Response compression and chunking configuration on the service side are not implemented. |
| Inbound payload validation | Not Yet Supported | Automatic constraint validation of inbound request payloads via `ballerina/constraint` is not implemented. |

### Common

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| Header value parsing utility | Supported | `parseHeader()` parses comma-separated header values with parameters into `HeaderValue[]`. |
| Distinct HTTP error types | Not Yet Supported | All errors surface as the generic `error` type; `http:ClientError`, `http:HeaderNotFoundError`, and similar subtypes are not declared — `is http:ClientError` type checks will not work. |
| Observability and metrics | Not Yet Supported | Metrics and tracing integration via `ballerina/observe` is not implemented. |
| XML payloads | Not Yet Supported | The `xml` type and related payload handling methods (`getXmlPayload()`, `setXmlPayload()`) are not implemented due to the lack of XML support in the Go runtime. |

### Notable Behavioural Changes

- **HTTP/1.0 is a compile error.** Specifying `httpVersion: "1.0"` in `ClientConfiguration` is rejected at compile time. Go's HTTP client cannot send HTTP/1.0 requests, so this is a permanent restriction rather than a missing runtime feature.
- **Trailing headers are not modelled.** The `TRAILING` header position constant is accepted at compile time for API compatibility, but all header operations (`getHeader`, `getHeaders`, `hasHeader`, `getHeaderNames`) act on transport (leading) headers at runtime. HTTP trailers sent by the server are silently discarded.
- **TLS protocol name has no effect.** The `protocol.name` field accepts `"SSL"`, `"TLS"`, and `"DTLS"` at compile time, but only TLS is supported at runtime. `"SSL"` and `"DTLS"` values are ignored because Go's standard TLS stack does not expose separate SSL or DTLS stacks.
