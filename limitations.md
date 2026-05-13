# HTTP Client Limitations

Known limitations in the native Ballerina interpreter's `ballerina/http` client support.

---

1. **No mechanism to declare distinct error types from Go stdlib code**
   `semtypes/error.go` has `errorDistinct(distinctId int)` but it is unexported. There is no public API to allocate a new distinct error ID from outside the `semtypes` package. `ClientError` (which upstream declares as `distinct Error`) is currently exposed as a plain `error`, so `c is http:ClientError` type narrowing won't work.

2. **Remote method VTable keys must use the `$remote$` prefix**
   The BIR gen emits `callInfo.Name = "$remote$get"` for `c->get(...)`. `resolveObjectMethod` looks up `string(callInfo.Name)` in the object's `methodKeys` map. The synthetic `BIRClassDef` VTable must use `"$remote$get"` / `"$remote$post"` etc. as keys, not plain method names. This is a non-obvious naming requirement with no existing documentation or helper, and must be remembered when adding future methods.

3. **`Response` objects use the generic `semtypes.OBJECT` runtime type**
   `Response` objects constructed by `Client.get` / `Client.post` on the Go side have `semtypes.OBJECT` as their runtime semtype rather than the specific `responseTy` registered in the symbol space. This means `r is http:Response` type narrowing won't work at runtime, even though `http:Response r = check c->get(...)` compiles correctly.

4. **`http:Client` runtime type narrowing is unreliable**
   Runtime type tests (`c is http:Client`) depend on the object's runtime semtype matching the specific `clientTy` registered during compilation. Objects constructed by `execNewObject` get their type from the LHS variable declaration's determined type, which may differ from `clientTy` depending on how the type resolver propagated the type. This makes `c is http:Client` expressions unreliable.

5. **`secureSocket` TLS configuration is partially supported**
   `secureSocket.enable = false` and `secureSocket.verifyHostName = false` disable TLS verification (`InsecureSkipVerify`). `secureSocket.cert` (string PEM file path) configures a custom CA trust store. `secureSocket.key` (`CertKey` record with `certFile`/`keyFile`) enables mTLS with unencrypted PEM files. Fields `ciphers`, `shareSession`, `handshakeTimeout`, `sessionTimeout`, and `serverName` are accepted at compile time but silently ignored at runtime. Not supported: `crypto:TrustStore` for `cert`, `crypto:KeyStore` for `key`, `protocol` and `certValidation` records. `keyPassword` in `CertKey` is accepted at compile time but ignored at runtime — `tls.X509KeyPair` requires unencrypted PEM files.

6. **`Response` cannot be directly constructed by user code**
   `execNewObject` would require a registered `BIRClassDef` for `Response` if user code ever wrote `new http:Response()`. Since users only receive `Response` from `Client.get` / `Client.post`, the Go side constructs `*values.Object` directly (bypassing `execNewObject`), injecting method keys manually. If a future use case requires user construction of a `Response`, a synthetic `BIRClassDef` would need to be registered.

7. **No corpus test for error propagation through `check c->get(...)`**
   The corpus tests only cover the 200 OK path. Network errors (connection refused, DNS failure) returned as `ClientError` through `check` are tested only in the extern test using `httptest.Server` — they cannot be hermetically reproduced in corpus tests without a PAL stub that returns errors on demand.

8. **`*ClientConfiguration` rest-param inclusion not supported**
   Upstream's `init(string url, *ClientConfiguration config)` uses a rest-param inclusion record (`*ClientConfiguration`), which spreads all record fields as individual named parameters. The interpreter does not support this syntax; callers must pass a record literal or variable directly.

9. **`RequestMessage` for `post` does not support `Request` objects, XML, streaming, or multipart**
   The `message` param is typed as Ballerina `json` at compile time (rejects objects, errors, functions, and xml). At runtime the interpreter serializes: `string` (→ `text/plain`), `byte[]` (→ `application/octet-stream`), `nil` (empty body), and all JSON-compatible values — `map<anydata>`, `boolean`, `int`, `float`, `decimal`, and nested maps/lists — to JSON (`application/json`). Unsupported: `http:Request` objects, `xml` values, `stream<byte[], io:Error?>`, and `mime:Entity[]` multipart. A `json[]` (list of JSON values) is treated as `byte[]` since list type cannot be distinguished at runtime.

10. **`xml` values are not representable at runtime**
    `semtypes.XML` exists as a type-system entry but there is no `values.Xml` runtime type. The interpreter cannot construct or pass XML values. This blocks the `xml` variant of `RequestMessage` and general XML processing.

11. **Default Content-Type inference for `post` is simplified**
    The interpreter infers `text/plain` for string bodies, `application/octet-stream` for `byte[]`, and `application/json` for JSON-compatible values. The real stdlib infers additional types (e.g. `application/xml` for `xml` values). The `mediaType` parameter can override this in all cases.

12. **`mediaType` with charset or other directives is passed verbatim**
    The upstream merges charset and other parameters into `Content-Type` (e.g. `application/json; charset=utf-8`). The interpreter passes the `mediaType` string directly to `req.Header.Set("Content-Type", ...)` without merging or normalizing directives.

13. **`getJsonPayload()` and `getBinaryPayload()` are supported; full `TargetType` inference is not**
    `r.getJsonPayload()` returns `json|error` (parses response body as JSON using Go's `encoding/json`). `r.getBinaryPayload()` returns `byte[]|error` (returns body as a byte array). The upstream pattern `json r = check c->get(...)` where `targetType = <>` is inferred from the LHS is not yet supported — it requires `typedesc<T>` runtime values and parameterized return type inference in the type resolver, which are unimplemented (`typeof` panics, no `TypeDesc` value type).

14. **`http:Request` object not supported as `post` message**
    Upstream accepts a fully constructed `http:Request` object (with headers, MIME parts, and body already attached) as the `message` argument. The interpreter has no `Request` class, so this usage is unsupported.

15. **No streaming body support for `post`**
    `stream<byte[], io:Error?>` as a `post` message body is not handled. This is important for large file uploads where buffering the entire body in memory is undesirable.

16. **`mime:Entity[]` multipart not supported for `post`**
    Multi-part form data uploads using `mime:Entity[]` require the `ballerina/mime` module, which is not implemented.

17. **Defaultable param lambda indices are hardcoded by convention**
    Default lambda names (`$Client.get$default$1`, `$Client.post$default$2`, etc.) encode the parameter index as a suffix. There is no validation that the registered extern name matches the BIR-generated `FunctionLookupKey`. If params are reordered, the mismatch would silently produce wrong defaults.

18. **`httpVersion: "1.0"` is not supported**
    Go's `net/http` client cannot send HTTP/1.0 request lines — `ProtoMajor`/`ProtoMinor` are ignored on client requests by design. The `HttpVersion` type is therefore `"1.1"|"2.0"` only; using `"1.0"` is a compile error. `"2.0"` enables HTTP/2 over TLS (via ALPN) and HTTP/2 cleartext (h2c) using Go 1.24+ built-in support.

19. **`http:Response.getHeader` / `getHeaders` return plain `error`, not `http:HeaderNotFoundError`**
    `HeaderNotFoundError` is declared as a distinct error type in upstream. Limitation #1 means distinct error types cannot be declared from Go; both methods return a plain `error`. `check r.getHeader(...)` works correctly; `r.getHeader(...) is http:HeaderNotFoundError` type narrowing does not.

20. **`HeaderPosition.TRAILING` is not supported on `http:Response`**
    The `position` parameter (`http:LEADING` or `http:TRAILING`) is accepted at compile time but silently ignored at runtime — all operations act on LEADING (transport) headers. Trailing headers (entity-body headers) are not modelled.

21. **`http:Response.headers` raw field is not exposed**
    The raw header map is an internal implementation detail. Use `hasHeader`, `getHeader`, `getHeaders`, and `getHeaderNames` to access response headers.
