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

package compile

import (
	"ballerina-lang-go/context"
	libcommon "ballerina-lang-go/lib/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

var HttpPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("http")},
	model.Name("0.0.1"),
)

func GetHttpSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	space := ctx.NewSymbolSpace(*HttpPackageID)

	addParseHeader(ctx, space)
	configSemType := addClientConfiguration(ctx, space)
	addClient(ctx, space, configSemType)

	return model.NewExportedSymbolSpace(space, nil)
}

// headerValueType returns the closed record type {| string value; map<string> params; |}
// used as the element type of parseHeader's return value.
func headerValueType(env semtypes.Env) semtypes.SemType {
	paramsMd := semtypes.NewMappingDefinition()
	paramsType := paramsMd.DefineMappingTypeWrapped(env, []semtypes.Field{}, semtypes.STRING)
	hvMd := semtypes.NewMappingDefinition()
	return hvMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("value", semtypes.STRING, false, false),
		semtypes.FieldFrom("params", paramsType, false, false),
	}, semtypes.NEVER)
}

func addParseHeader(ctx *context.CompilerContext, space *model.SymbolSpace) {
	env := ctx.GetTypeEnv()

	hvType := headerValueType(env)
	hvSym := model.NewTypeSymbol("HeaderValue", true)
	hvSym.SetType(hvType)
	space.AddSymbol("HeaderValue", &hvSym)

	hvListLd := semtypes.NewListDefinition()
	hvListType := hvListLd.DefineListTypeWrappedWithEnvSemType(env, hvType)

	sig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: semtypes.Union(hvListType, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	sym := model.NewFunctionSymbol("parseHeader", sig, true)
	space.AddSymbol("parseHeader", sym)
	ref, _ := space.GetSymbol("parseHeader")
	ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(env, &sig))
}

func addClientConfiguration(ctx *context.CompilerContext, space *model.SymbolSpace) semtypes.SemType {
	env := ctx.GetTypeEnv()

	// CertKey: simplified mTLS record matching upstream http:CertKey.
	// certFile and keyFile are required; keyPassword is accepted but ignored at runtime
	// (tls.X509KeyPair requires unencrypted PEM files).
	certKeyMd := semtypes.NewMappingDefinition()
	certKeySemType := certKeyMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("certFile", semtypes.STRING, false, false),
		semtypes.FieldFrom("keyFile", semtypes.STRING, false, false),
		semtypes.FieldFrom("keyPassword", semtypes.STRING, false, true),
	}, semtypes.NEVER)
	certKeySym := model.NewTypeSymbol("CertKey", true)
	certKeySym.SetType(certKeySemType)
	space.AddSymbol("CertKey", &certKeySym)

	// Protocol: SSL|TLS|DTLS. Go only supports TLS; SSL and DTLS are compile-time-only.
	protocolSemType := semtypes.Union(
		semtypes.Union(semtypes.StringConst("SSL"), semtypes.StringConst("TLS")),
		semtypes.StringConst("DTLS"),
	)
	protocolSym := model.NewTypeSymbol("Protocol", true)
	protocolSym.SetType(protocolSemType)
	space.AddSymbol("Protocol", &protocolSym)

	// protocol record: {| Protocol name; string[] versions; |}
	protocolRecordMd := semtypes.NewMappingDefinition()
	protocolRecordSemType := protocolRecordMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("name", protocolSemType, false, false),
		semtypes.FieldFrom("versions", semtypes.LIST, false, false),
	}, semtypes.NEVER)

	// CertValidationType: OCSP_CRL|OCSP_STAPLING — accepted at compile time, not implemented.
	certValidTypeSemType := semtypes.Union(
		semtypes.StringConst("OCSP_CRL"), semtypes.StringConst("OCSP_STAPLING"),
	)
	certValidTypeSym := model.NewTypeSymbol("CertValidationType", true)
	certValidTypeSym.SetType(certValidTypeSemType)
	space.AddSymbol("CertValidationType", &certValidTypeSym)

	// certValidation record: {| CertValidationType 'type; int cacheSize; int cacheValidityPeriod; |}
	certValidRecordMd := semtypes.NewMappingDefinition()
	certValidRecordSemType := certValidRecordMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("type", certValidTypeSemType, false, false),
		semtypes.FieldFrom("cacheSize", semtypes.INT, false, false),
		semtypes.FieldFrom("cacheValidityPeriod", semtypes.INT, false, false),
	}, semtypes.NEVER)

	// ClientSecureSocket: matches upstream http:ClientSecureSocket field names.
	// cert accepts string only (not crypto:TrustStore).
	// key accepts CertKey only (not crypto:KeyStore).
	// Implemented: enable, verifyHostName, cert, key, serverName, ciphers, handshakeTimeout, shareSession, protocol.versions.
	// Accepted but not implemented: sessionTimeout, keyPassword, certValidation, protocol.name.
	secureSocketMd := semtypes.NewMappingDefinition()
	secureSocketSemType := secureSocketMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("enable", semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("cert", semtypes.STRING, false, true),
		semtypes.FieldFrom("key", certKeySemType, false, true),
		semtypes.FieldFrom("protocol", semtypes.Union(protocolRecordSemType, semtypes.NIL), false, true),
		semtypes.FieldFrom("certValidation", semtypes.Union(certValidRecordSemType, semtypes.NIL), false, true),
		semtypes.FieldFrom("ciphers", semtypes.LIST, false, true),
		semtypes.FieldFrom("verifyHostName", semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("shareSession", semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("handshakeTimeout", semtypes.DECIMAL, false, true),
		semtypes.FieldFrom("sessionTimeout", semtypes.DECIMAL, false, true),
		semtypes.FieldFrom("serverName", semtypes.STRING, false, true),
	}, semtypes.NEVER)
	secureSocketSym := model.NewTypeSymbol("ClientSecureSocket", true)
	secureSocketSym.SetType(secureSocketSemType)
	space.AddSymbol("ClientSecureSocket", &secureSocketSym)

	// FollowRedirects: matches upstream http:FollowRedirects field names.
	// enabled defaults to false (no redirects by default), maxCount defaults to 5,
	// allowAuthHeaders defaults to false (auth headers stripped on redirect).
	followRedirectsMd := semtypes.NewMappingDefinition()
	followRedirectsSemType := followRedirectsMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("enabled", semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("maxCount", semtypes.INT, false, true),
		semtypes.FieldFrom("allowAuthHeaders", semtypes.BOOLEAN, false, true),
	}, semtypes.NEVER)
	followRedirectsSym := model.NewTypeSymbol("FollowRedirects", true)
	followRedirectsSym.SetType(followRedirectsSemType)
	space.AddSymbol("FollowRedirects", &followRedirectsSym)

	// HttpVersion: "1.1"|"2.0". "1.0" is omitted — Go's net/http client cannot send HTTP/1.0.
	httpVersionSemType := semtypes.Union(semtypes.StringConst("1.1"), semtypes.StringConst("2.0"))
	httpVersionSym := model.NewTypeSymbol("HttpVersion", true)
	httpVersionSym.SetType(httpVersionSemType)
	space.AddSymbol("HttpVersion", &httpVersionSym)

	// ClientConfiguration: matching upstream http:ClientConfiguration field names.
	md := semtypes.NewMappingDefinition()
	configSemType := md.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("timeout", semtypes.DECIMAL, false, true),
		semtypes.FieldFrom("followRedirects", semtypes.Union(followRedirectsSemType, semtypes.NIL), false, true),
		semtypes.FieldFrom("httpVersion", httpVersionSemType, false, true),
		semtypes.FieldFrom("secureSocket", semtypes.Union(secureSocketSemType, semtypes.NIL), false, true),
	}, semtypes.NEVER)
	configSym := model.NewTypeSymbol("ClientConfiguration", true)
	configSym.SetType(configSemType)
	space.AddSymbol("ClientConfiguration", &configSym)
	return configSemType
}

// registerDefaultLambda registers a default-parameter lambda function symbol and returns its ref.
// All default lambdas are internal (public=false) and isolated.
func registerDefaultLambda(ctx *context.CompilerContext, space *model.SymbolSpace, name string, sig model.FunctionSignature) model.SymbolRef {
	sym := model.NewFunctionSymbol(name, sig, false)
	space.AddSymbol(name, sym)
	ref, _ := space.GetSymbol(name)
	ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &sig))
	return ref
}

func addClient(ctx *context.CompilerContext, space *model.SymbolSpace, configSemType semtypes.SemType) {
	env := ctx.GetTypeEnv()

	// headers: map<string|string[]>? — open mapping (any key, string|string[] values), optional.
	// Build an explicit open mapping type so the field value type resolves to STRING|string[]
	// rather than NEVER (which happens when the basic MAPPING atom is used directly), and so
	// the list arm rejects non-string lists like int[] at compile time.
	stringArrayLd := semtypes.NewListDefinition()
	stringArrayType := stringArrayLd.DefineListTypeWrappedWithEnvSemType(env, semtypes.STRING)
	byteArrayLd := semtypes.NewListDefinition()
	byteArrayType := byteArrayLd.DefineListTypeWrappedWithEnvSemType(env, semtypes.BYTE)
	headersMd := semtypes.NewMappingDefinition()
	headersMapType := headersMd.DefineMappingTypeWrapped(env,
		[]semtypes.Field{},
		semtypes.Union(semtypes.STRING, stringArrayType))
	headersOptType := semtypes.Union(headersMapType, semtypes.NIL)

	// HeaderPosition: "LEADING"|"TRAILING". TRAILING is accepted at compile time but ignored at runtime.
	headerPositionSemType := semtypes.Union(semtypes.StringConst("LEADING"), semtypes.StringConst("TRAILING"))
	headerPositionSym := model.NewTypeSymbol("HeaderPosition", true)
	headerPositionSym.SetType(headerPositionSemType)
	space.AddSymbol("HeaderPosition", &headerPositionSym)

	leadingSym := model.NewValueSymbol("LEADING", true, true, false)
	leadingSym.SetType(semtypes.StringConst("LEADING"))
	space.AddSymbol("LEADING", &leadingSym)

	trailingSym := model.NewValueSymbol("TRAILING", true, true, false)
	trailingSym.SetType(semtypes.StringConst("TRAILING"))
	space.AddSymbol("TRAILING", &trailingSym)

	// json — the proper recursive Ballerina json type: nil|boolean|int|float|decimal|string|list(json)|map(json).
	jsonType := semtypes.CreateJSON(semtypes.ContextFrom(ctx.GetTypeEnv()))

	// Response: declared as a class so the type checker knows about statusCode and
	// the header API. Users never write `new http:Response()` — Response objects are
	// only constructed on the Go side by Client.get. The raw headers map is intentionally
	// not exposed; use hasHeader/getHeader/getHeaders/getHeaderNames instead.
	gtpSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{},
		ReturnType: semtypes.STRING,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	gtpFnSemType := libcommon.FunctionSignatureToSemType(env, &gtpSig)

	hasHeaderSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headerPositionSemType},
		ParamNames: []string{"headerName", "position"},
		ReturnType: semtypes.BOOLEAN,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	hasHeaderFnSemType := libcommon.FunctionSignatureToSemType(env, &hasHeaderSig)

	getHeaderSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headerPositionSemType},
		ParamNames: []string{"headerName", "position"},
		ReturnType: semtypes.Union(semtypes.STRING, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeaderFnSemType := libcommon.FunctionSignatureToSemType(env, &getHeaderSig)

	getHeadersSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headerPositionSemType},
		ParamNames: []string{"headerName", "position"},
		ReturnType: semtypes.Union(stringArrayType, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeadersFnSemType := libcommon.FunctionSignatureToSemType(env, &getHeadersSig)

	getHeaderNamesSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{headerPositionSemType},
		ParamNames: []string{"position"},
		ReturnType: stringArrayType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeaderNamesFnSemType := libcommon.FunctionSignatureToSemType(env, &getHeaderNamesSig)

	gjpSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{},
		ReturnType: semtypes.Union(jsonType, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	gjpFnSemType := libcommon.FunctionSignatureToSemType(env, &gjpSig)

	gbpSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{},
		ReturnType: semtypes.Union(byteArrayType, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	gbpFnSemType := libcommon.FunctionSignatureToSemType(env, &gbpSig)

	responseOd := semtypes.NewObjectDefinition()
	responseTy := responseOd.Define(env,
		semtypes.ObjectQualifiersDEFAULT,
		[]semtypes.Member{
			{Name: "statusCode", ValueTy: semtypes.INT, Kind: semtypes.MemberKindField, Visibility: semtypes.VisibilityPublic},
			{Name: "getTextPayload", ValueTy: gtpFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getJsonPayload", ValueTy: gjpFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getBinaryPayload", ValueTy: gbpFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "hasHeader", ValueTy: hasHeaderFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeader", ValueTy: getHeaderFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeaders", ValueTy: getHeadersFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeaderNames", ValueTy: getHeaderNamesFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
		})

	gtpSym := model.NewFunctionSymbol("$Response.getTextPayload", gtpSig, false)
	space.AddSymbol("$Response.getTextPayload", gtpSym)
	gtpRef, _ := space.GetSymbol("$Response.getTextPayload")
	ctx.SetSymbolType(gtpRef, gtpFnSemType)

	gjpSym := model.NewFunctionSymbol("$Response.getJsonPayload", gjpSig, false)
	space.AddSymbol("$Response.getJsonPayload", gjpSym)
	gjpRef, _ := space.GetSymbol("$Response.getJsonPayload")
	ctx.SetSymbolType(gjpRef, gjpFnSemType)

	gbpSym := model.NewFunctionSymbol("$Response.getBinaryPayload", gbpSig, false)
	space.AddSymbol("$Response.getBinaryPayload", gbpSym)
	gbpRef, _ := space.GetSymbol("$Response.getBinaryPayload")
	ctx.SetSymbolType(gbpRef, gbpFnSemType)

	// Response header method default lambdas: position param → "LEADING"
	posDefault1 := model.FunctionSignature{ParamTypes: []semtypes.SemType{semtypes.STRING}, ReturnType: headerPositionSemType, Flags: model.FuncSymbolFlagIsolated}
	posDefault0 := model.FunctionSignature{ParamTypes: []semtypes.SemType{}, ReturnType: headerPositionSemType, Flags: model.FuncSymbolFlagIsolated}

	hasHeaderDefaultRef := registerDefaultLambda(ctx, space, "$Response.hasHeader$default$1", posDefault1)
	hasHeaderSym := model.NewFunctionSymbol("$Response.hasHeader", hasHeaderSig, false)
	space.AddSymbol("$Response.hasHeader", hasHeaderSym)
	hasHeaderRef, _ := space.GetSymbol("$Response.hasHeader")
	ctx.SetSymbolType(hasHeaderRef, hasHeaderFnSemType)
	hasHeaderDefaultable := model.NewDefaultableParamInfo(len(hasHeaderSig.ParamTypes))
	hasHeaderDefaultable.SetDefaultable(1, hasHeaderDefaultRef)
	hasHeaderSym.SetDefaultableParams(hasHeaderDefaultable)

	getHeaderDefaultRef := registerDefaultLambda(ctx, space, "$Response.getHeader$default$1", posDefault1)
	getHeaderSym := model.NewFunctionSymbol("$Response.getHeader", getHeaderSig, false)
	space.AddSymbol("$Response.getHeader", getHeaderSym)
	getHeaderRef, _ := space.GetSymbol("$Response.getHeader")
	ctx.SetSymbolType(getHeaderRef, getHeaderFnSemType)
	getHeaderDefaultable := model.NewDefaultableParamInfo(len(getHeaderSig.ParamTypes))
	getHeaderDefaultable.SetDefaultable(1, getHeaderDefaultRef)
	getHeaderSym.SetDefaultableParams(getHeaderDefaultable)

	getHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Response.getHeaders$default$1", posDefault1)
	getHeadersSym := model.NewFunctionSymbol("$Response.getHeaders", getHeadersSig, false)
	space.AddSymbol("$Response.getHeaders", getHeadersSym)
	getHeadersRef, _ := space.GetSymbol("$Response.getHeaders")
	ctx.SetSymbolType(getHeadersRef, getHeadersFnSemType)
	getHeadersDefaultable := model.NewDefaultableParamInfo(len(getHeadersSig.ParamTypes))
	getHeadersDefaultable.SetDefaultable(1, getHeadersDefaultRef)
	getHeadersSym.SetDefaultableParams(getHeadersDefaultable)

	getHeaderNamesDefaultRef := registerDefaultLambda(ctx, space, "$Response.getHeaderNames$default$0", posDefault0)
	getHeaderNamesSym := model.NewFunctionSymbol("$Response.getHeaderNames", getHeaderNamesSig, false)
	space.AddSymbol("$Response.getHeaderNames", getHeaderNamesSym)
	getHeaderNamesRef, _ := space.GetSymbol("$Response.getHeaderNames")
	ctx.SetSymbolType(getHeaderNamesRef, getHeaderNamesFnSemType)
	getHeaderNamesDefaultable := model.NewDefaultableParamInfo(len(getHeaderNamesSig.ParamTypes))
	getHeaderNamesDefaultable.SetDefaultable(0, getHeaderNamesDefaultRef)
	getHeaderNamesSym.SetDefaultableParams(getHeaderNamesDefaultable)

	responseSym := model.NewClassSymbol("Response", true)
	responseSym.SetType(responseTy)
	responseSym.SetMethods(map[string]model.SymbolRef{
		"getTextPayload":   gtpRef,
		"getJsonPayload":   gjpRef,
		"getBinaryPayload": gbpRef,
		"hasHeader":        hasHeaderRef,
		"getHeader":        getHeaderRef,
		"getHeaders":       getHeadersRef,
		"getHeaderNames":   getHeaderNamesRef,
	})
	space.AddSymbol("Response", &responseSym)

	// Member-level signatures: self is NOT included here because the BIR gen prepends
	// the receiver object automatically. The type checker only sees user-provided args.
	initSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, configSemType},
		ReturnType: semtypes.Union(semtypes.NIL, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	initFnSemType := libcommon.FunctionSignatureToSemType(env, &initSig)

	getSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headersOptType},
		ParamNames: []string{"path", "headers"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getFnSemType := libcommon.FunctionSignatureToSemType(env, &getSig)

	// post: path(string), message(json), headers?(map<string|string[]>?), mediaType?(string?)
	mediaTypeOptType := semtypes.Union(semtypes.STRING, semtypes.NIL)
	postSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType, mediaTypeOptType},
		ParamNames: []string{"path", "message", "headers", "mediaType"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	postFnSemType := libcommon.FunctionSignatureToSemType(env, &postSig)

	// head / options — body-less, like get
	headSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headersOptType},
		ParamNames: []string{"path", "headers"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	headFnSemType := libcommon.FunctionSignatureToSemType(env, &headSig)

	optionsSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, headersOptType},
		ParamNames: []string{"path", "headers"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	optionsFnSemType := libcommon.FunctionSignatureToSemType(env, &optionsSig)

	// put / patch — body required, like post
	putSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType, mediaTypeOptType},
		ParamNames: []string{"path", "message", "headers", "mediaType"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	putFnSemType := libcommon.FunctionSignatureToSemType(env, &putSig)

	patchSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType, mediaTypeOptType},
		ParamNames: []string{"path", "message", "headers", "mediaType"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	patchFnSemType := libcommon.FunctionSignatureToSemType(env, &patchSig)

	// delete — message is optional (defaults to ())
	deleteMessageType := semtypes.Union(jsonType, semtypes.NIL)
	deleteSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, deleteMessageType, headersOptType, mediaTypeOptType},
		ParamNames: []string{"path", "message", "headers", "mediaType"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	deleteFnSemType := libcommon.FunctionSignatureToSemType(env, &deleteSig)

	// execute — explicit httpVerb as first param, message required
	executeSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, semtypes.STRING, jsonType, headersOptType, mediaTypeOptType},
		ParamNames: []string{"httpVerb", "path", "message", "headers", "mediaType"},
		ReturnType: semtypes.Union(responseTy, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	executeFnSemType := libcommon.FunctionSignatureToSemType(env, &executeSig)

	// Build a proper client-qualified object semtype so the type checker
	// accepts c->get(...), c->post(...), and new http:Client(...) correctly.
	od := semtypes.NewObjectDefinition()
	clientTy := od.Define(env,
		semtypes.ObjectQualifiersFrom(true, false, semtypes.NetworkQualifierClient),
		[]semtypes.Member{
			{Name: "init", ValueTy: initFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("get"), ValueTy: getFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("post"), ValueTy: postFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("head"), ValueTy: headFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("options"), ValueTy: optionsFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("put"), ValueTy: putFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("patch"), ValueTy: patchFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("delete"), ValueTy: deleteFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: model.RemoteMethodName("execute"), ValueTy: executeFnSemType, Kind: semtypes.MemberKindRemoteMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
		})

	// Default lambda for the config param (index 1): $Client.init$default$1(url) → {}
	initDefaultRef := registerDefaultLambda(ctx, space, "$Client.init$default$1", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: configSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	initSym := model.NewFunctionSymbol("$Client.init", initSig, false)
	space.AddSymbol("$Client.init", initSym)
	initRef, _ := space.GetSymbol("$Client.init")
	ctx.SetSymbolType(initRef, initFnSemType)
	initDefaultableInfo := model.NewDefaultableParamInfo(len(initSig.ParamTypes))
	initDefaultableInfo.SetDefaultable(1, initDefaultRef)
	initSym.SetDefaultableParams(initDefaultableInfo)

	// get: headers at index 1
	getDefaultRef := registerDefaultLambda(ctx, space, "$Client.get$default$1", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	getSym := model.NewFunctionSymbol("$Client.get", getSig, false)
	space.AddSymbol("$Client.get", getSym)
	getRef, _ := space.GetSymbol("$Client.get")
	ctx.SetSymbolType(getRef, getFnSemType)
	getDefaultableInfo := model.NewDefaultableParamInfo(len(getSig.ParamTypes))
	getDefaultableInfo.SetDefaultable(1, getDefaultRef)
	getSym.SetDefaultableParams(getDefaultableInfo)

	// post: headers at index 2, mediaType at index 3
	postHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Client.post$default$2", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	postMediaTypeDefaultRef := registerDefaultLambda(ctx, space, "$Client.post$default$3", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	postSym := model.NewFunctionSymbol("$Client.post", postSig, false)
	space.AddSymbol("$Client.post", postSym)
	postRef, _ := space.GetSymbol("$Client.post")
	ctx.SetSymbolType(postRef, postFnSemType)
	postDefaultableInfo := model.NewDefaultableParamInfo(len(postSig.ParamTypes))
	postDefaultableInfo.SetDefaultable(2, postHeadersDefaultRef)
	postDefaultableInfo.SetDefaultable(3, postMediaTypeDefaultRef)
	postSym.SetDefaultableParams(postDefaultableInfo)

	// head: headers at index 1
	headDefaultRef := registerDefaultLambda(ctx, space, "$Client.head$default$1", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	headSym := model.NewFunctionSymbol("$Client.head", headSig, false)
	space.AddSymbol("$Client.head", headSym)
	headRef, _ := space.GetSymbol("$Client.head")
	ctx.SetSymbolType(headRef, headFnSemType)
	headDefaultableInfo := model.NewDefaultableParamInfo(len(headSig.ParamTypes))
	headDefaultableInfo.SetDefaultable(1, headDefaultRef)
	headSym.SetDefaultableParams(headDefaultableInfo)

	// options: headers at index 1
	optionsDefaultRef := registerDefaultLambda(ctx, space, "$Client.options$default$1", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	optionsSym := model.NewFunctionSymbol("$Client.options", optionsSig, false)
	space.AddSymbol("$Client.options", optionsSym)
	optionsRef, _ := space.GetSymbol("$Client.options")
	ctx.SetSymbolType(optionsRef, optionsFnSemType)
	optionsDefaultableInfo := model.NewDefaultableParamInfo(len(optionsSig.ParamTypes))
	optionsDefaultableInfo.SetDefaultable(1, optionsDefaultRef)
	optionsSym.SetDefaultableParams(optionsDefaultableInfo)

	// put: headers at index 2, mediaType at index 3
	putHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Client.put$default$2", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	putMediaTypeDefaultRef := registerDefaultLambda(ctx, space, "$Client.put$default$3", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	putSym := model.NewFunctionSymbol("$Client.put", putSig, false)
	space.AddSymbol("$Client.put", putSym)
	putRef, _ := space.GetSymbol("$Client.put")
	ctx.SetSymbolType(putRef, putFnSemType)
	putDefaultableInfo := model.NewDefaultableParamInfo(len(putSig.ParamTypes))
	putDefaultableInfo.SetDefaultable(2, putHeadersDefaultRef)
	putDefaultableInfo.SetDefaultable(3, putMediaTypeDefaultRef)
	putSym.SetDefaultableParams(putDefaultableInfo)

	// patch: headers at index 2, mediaType at index 3
	patchHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Client.patch$default$2", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	patchMediaTypeDefaultRef := registerDefaultLambda(ctx, space, "$Client.patch$default$3", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	patchSym := model.NewFunctionSymbol("$Client.patch", patchSig, false)
	space.AddSymbol("$Client.patch", patchSym)
	patchRef, _ := space.GetSymbol("$Client.patch")
	ctx.SetSymbolType(patchRef, patchFnSemType)
	patchDefaultableInfo := model.NewDefaultableParamInfo(len(patchSig.ParamTypes))
	patchDefaultableInfo.SetDefaultable(2, patchHeadersDefaultRef)
	patchDefaultableInfo.SetDefaultable(3, patchMediaTypeDefaultRef)
	patchSym.SetDefaultableParams(patchDefaultableInfo)

	// delete: message at index 1, headers at index 2, mediaType at index 3
	deleteMessageDefaultRef := registerDefaultLambda(ctx, space, "$Client.delete$default$1", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: deleteMessageType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	deleteHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Client.delete$default$2", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, deleteMessageType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	deleteMediaTypeDefaultRef := registerDefaultLambda(ctx, space, "$Client.delete$default$3", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, deleteMessageType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	deleteSym := model.NewFunctionSymbol("$Client.delete", deleteSig, false)
	space.AddSymbol("$Client.delete", deleteSym)
	deleteRef, _ := space.GetSymbol("$Client.delete")
	ctx.SetSymbolType(deleteRef, deleteFnSemType)
	deleteDefaultableInfo := model.NewDefaultableParamInfo(len(deleteSig.ParamTypes))
	deleteDefaultableInfo.SetDefaultable(1, deleteMessageDefaultRef)
	deleteDefaultableInfo.SetDefaultable(2, deleteHeadersDefaultRef)
	deleteDefaultableInfo.SetDefaultable(3, deleteMediaTypeDefaultRef)
	deleteSym.SetDefaultableParams(deleteDefaultableInfo)

	// execute: headers at index 3, mediaType at index 4
	executeHeadersDefaultRef := registerDefaultLambda(ctx, space, "$Client.execute$default$3", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	executeMediaTypeDefaultRef := registerDefaultLambda(ctx, space, "$Client.execute$default$4", model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	})
	executeSym := model.NewFunctionSymbol("$Client.execute", executeSig, false)
	space.AddSymbol("$Client.execute", executeSym)
	executeRef, _ := space.GetSymbol("$Client.execute")
	ctx.SetSymbolType(executeRef, executeFnSemType)
	executeDefaultableInfo := model.NewDefaultableParamInfo(len(executeSig.ParamTypes))
	executeDefaultableInfo.SetDefaultable(3, executeHeadersDefaultRef)
	executeDefaultableInfo.SetDefaultable(4, executeMediaTypeDefaultRef)
	executeSym.SetDefaultableParams(executeDefaultableInfo)

	clientSym := model.NewClassSymbol("Client", true)
	clientSym.SetType(clientTy)
	clientSym.SetMethods(map[string]model.SymbolRef{
		"init":                            initRef,
		model.RemoteMethodName("get"):     getRef,
		model.RemoteMethodName("post"):    postRef,
		model.RemoteMethodName("head"):    headRef,
		model.RemoteMethodName("options"): optionsRef,
		model.RemoteMethodName("put"):     putRef,
		model.RemoteMethodName("patch"):   patchRef,
		model.RemoteMethodName("delete"):  deleteRef,
		model.RemoteMethodName("execute"): executeRef,
	})
	space.AddSymbol("Client", &clientSym)
}
