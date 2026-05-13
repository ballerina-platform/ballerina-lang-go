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

func addParseHeader(ctx *context.CompilerContext, space *model.SymbolSpace) {
	sig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: semtypes.Union(semtypes.LIST, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	sym := model.NewFunctionSymbol("parseHeader", sig, true)
	space.AddSymbol("parseHeader", sym)
	ref, _ := space.GetSymbol("parseHeader")
	ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &sig))
}

func addClientConfiguration(ctx *context.CompilerContext, space *model.SymbolSpace) semtypes.SemType {
	env := ctx.GetTypeEnv()

	// CertKey: simplified mTLS record matching upstream http:CertKey.
	// certFile and keyFile are required; keyPassword is accepted but ignored at runtime
	// (tls.X509KeyPair requires unencrypted PEM files).
	certKeyMd := semtypes.NewMappingDefinition()
	certKeySemType := certKeyMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("certFile",    semtypes.STRING, false, false),
		semtypes.FieldFrom("keyFile",     semtypes.STRING, false, false),
		semtypes.FieldFrom("keyPassword", semtypes.STRING, false, true),
	}, semtypes.NEVER)
	certKeySym := model.NewTypeSymbol("CertKey", true)
	certKeySym.SetType(certKeySemType)
	space.AddSymbol("CertKey", &certKeySym)

	// ClientSecureSocket: matches upstream http:ClientSecureSocket field names.
	// cert accepts string only (not crypto:TrustStore).
	// key accepts CertKey only (not crypto:KeyStore).
	// Fields ciphers/shareSession/handshakeTimeout/sessionTimeout/serverName are
	// accepted at compile time but silently ignored at runtime.
	secureSocketMd := semtypes.NewMappingDefinition()
	secureSocketSemType := secureSocketMd.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("enable",           semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("cert",             semtypes.STRING,  false, true),
		semtypes.FieldFrom("key",              certKeySemType,   false, true),
		semtypes.FieldFrom("verifyHostName",   semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("shareSession",     semtypes.BOOLEAN, false, true),
		semtypes.FieldFrom("handshakeTimeout", semtypes.DECIMAL, false, true),
		semtypes.FieldFrom("sessionTimeout",   semtypes.DECIMAL, false, true),
		semtypes.FieldFrom("serverName",       semtypes.STRING,  false, true),
		semtypes.FieldFrom("ciphers",          semtypes.LIST,    false, true),
	}, semtypes.NEVER)
	secureSocketSym := model.NewTypeSymbol("ClientSecureSocket", true)
	secureSocketSym.SetType(secureSocketSemType)
	space.AddSymbol("ClientSecureSocket", &secureSocketSym)

	// HttpVersion: "1.1"|"2.0". "1.0" is omitted — Go's net/http client cannot send HTTP/1.0.
	httpVersionSemType := semtypes.Union(semtypes.StringConst("1.1"), semtypes.StringConst("2.0"))
	httpVersionSym := model.NewTypeSymbol("HttpVersion", true)
	httpVersionSym.SetType(httpVersionSemType)
	space.AddSymbol("HttpVersion", &httpVersionSym)

	// ClientConfiguration: existing fields + secureSocket?: ClientSecureSocket?
	md := semtypes.NewMappingDefinition()
	configSemType := md.DefineMappingTypeWrapped(env, []semtypes.Field{
		semtypes.FieldFrom("timeout",         semtypes.DECIMAL,                                    false, true),
		semtypes.FieldFrom("followRedirects",  semtypes.BOOLEAN,                                   false, true),
		semtypes.FieldFrom("httpVersion",     httpVersionSemType,                                  false, true),
		semtypes.FieldFrom("secureSocket",    semtypes.Union(secureSocketSemType, semtypes.NIL),   false, true),
	}, semtypes.NEVER)
	configSym := model.NewTypeSymbol("ClientConfiguration", true)
	configSym.SetType(configSemType)
	space.AddSymbol("ClientConfiguration", &configSym)
	return configSemType
}

func addClient(ctx *context.CompilerContext, space *model.SymbolSpace, configSemType semtypes.SemType) {
	env := ctx.GetTypeEnv()

	// headers: map<string|string[]>? — open mapping (any key, string|list values), optional.
	// Build an explicit open mapping type so the field value type resolves to STRING|LIST
	// rather than NEVER (which happens when the basic MAPPING atom is used directly).
	headersMd := semtypes.NewMappingDefinition()
	headersMapType := headersMd.DefineMappingTypeWrapped(env,
		[]semtypes.Field{},
		semtypes.Union(semtypes.STRING, semtypes.LIST))
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

	// json ≈ NIL|BOOLEAN|INT|FLOAT|DECIMAL|STRING|list|map — approximation of Ballerina json type.
	// Rejects objects, errors, functions, and xml at compile time while accepting all JSON-serializable values.
	// Defined here (before Response) so it can be used in Response payload method signatures.
	jsonType := semtypes.Union(semtypes.SIMPLE_OR_STRING, semtypes.Union(semtypes.LIST, semtypes.MAPPING))

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
		ReturnType: semtypes.Union(semtypes.LIST, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeadersFnSemType := libcommon.FunctionSignatureToSemType(env, &getHeadersSig)

	getHeaderNamesSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{headerPositionSemType},
		ParamNames: []string{"position"},
		ReturnType: semtypes.LIST,
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
		ReturnType: semtypes.Union(semtypes.LIST, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	gbpFnSemType := libcommon.FunctionSignatureToSemType(env, &gbpSig)

	responseOd := semtypes.NewObjectDefinition()
	responseTy := responseOd.Define(env,
		semtypes.ObjectQualifiersDEFAULT,
		[]semtypes.Member{
			{Name: "statusCode",      ValueTy: semtypes.INT,            Kind: semtypes.MemberKindField,  Visibility: semtypes.VisibilityPublic},
			{Name: "getTextPayload",  ValueTy: gtpFnSemType,            Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getJsonPayload",  ValueTy: gjpFnSemType,            Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getBinaryPayload",ValueTy: gbpFnSemType,            Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "hasHeader",       ValueTy: hasHeaderFnSemType,      Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeader",       ValueTy: getHeaderFnSemType,      Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeaders",      ValueTy: getHeadersFnSemType,     Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
			{Name: "getHeaderNames",  ValueTy: getHeaderNamesFnSemType, Kind: semtypes.MemberKindMethod, Visibility: semtypes.VisibilityPublic, Immutable: true},
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

	// hasHeader default lambda: position (index 1) → "LEADING"
	hasHeaderDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headerPositionSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	hasHeaderDefaultSym := model.NewFunctionSymbol("$Response.hasHeader$default$1", hasHeaderDefaultSig, false)
	space.AddSymbol("$Response.hasHeader$default$1", hasHeaderDefaultSym)
	hasHeaderDefaultRef, _ := space.GetSymbol("$Response.hasHeader$default$1")
	ctx.SetSymbolType(hasHeaderDefaultRef, libcommon.FunctionSignatureToSemType(env, &hasHeaderDefaultSig))

	hasHeaderSym := model.NewFunctionSymbol("$Response.hasHeader", hasHeaderSig, false)
	space.AddSymbol("$Response.hasHeader", hasHeaderSym)
	hasHeaderRef, _ := space.GetSymbol("$Response.hasHeader")
	ctx.SetSymbolType(hasHeaderRef, hasHeaderFnSemType)
	hasHeaderDefaultable := model.NewDefaultableParamInfo(len(hasHeaderSig.ParamTypes))
	hasHeaderDefaultable.SetDefaultable(1, hasHeaderDefaultRef)
	hasHeaderSym.SetDefaultableParams(hasHeaderDefaultable)

	// getHeader default lambda: position (index 1) → "LEADING"
	getHeaderDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headerPositionSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeaderDefaultSym := model.NewFunctionSymbol("$Response.getHeader$default$1", getHeaderDefaultSig, false)
	space.AddSymbol("$Response.getHeader$default$1", getHeaderDefaultSym)
	getHeaderDefaultRef, _ := space.GetSymbol("$Response.getHeader$default$1")
	ctx.SetSymbolType(getHeaderDefaultRef, libcommon.FunctionSignatureToSemType(env, &getHeaderDefaultSig))

	getHeaderSym := model.NewFunctionSymbol("$Response.getHeader", getHeaderSig, false)
	space.AddSymbol("$Response.getHeader", getHeaderSym)
	getHeaderRef, _ := space.GetSymbol("$Response.getHeader")
	ctx.SetSymbolType(getHeaderRef, getHeaderFnSemType)
	getHeaderDefaultable := model.NewDefaultableParamInfo(len(getHeaderSig.ParamTypes))
	getHeaderDefaultable.SetDefaultable(1, getHeaderDefaultRef)
	getHeaderSym.SetDefaultableParams(getHeaderDefaultable)

	// getHeaders default lambda: position (index 1) → "LEADING"
	getHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headerPositionSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeadersDefaultSym := model.NewFunctionSymbol("$Response.getHeaders$default$1", getHeadersDefaultSig, false)
	space.AddSymbol("$Response.getHeaders$default$1", getHeadersDefaultSym)
	getHeadersDefaultRef, _ := space.GetSymbol("$Response.getHeaders$default$1")
	ctx.SetSymbolType(getHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &getHeadersDefaultSig))

	getHeadersSym := model.NewFunctionSymbol("$Response.getHeaders", getHeadersSig, false)
	space.AddSymbol("$Response.getHeaders", getHeadersSym)
	getHeadersRef, _ := space.GetSymbol("$Response.getHeaders")
	ctx.SetSymbolType(getHeadersRef, getHeadersFnSemType)
	getHeadersDefaultable := model.NewDefaultableParamInfo(len(getHeadersSig.ParamTypes))
	getHeadersDefaultable.SetDefaultable(1, getHeadersDefaultRef)
	getHeadersSym.SetDefaultableParams(getHeadersDefaultable)

	// getHeaderNames default lambda: position (index 0) → "LEADING"
	getHeaderNamesDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{},
		ReturnType: headerPositionSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getHeaderNamesDefaultSym := model.NewFunctionSymbol("$Response.getHeaderNames$default$0", getHeaderNamesDefaultSig, false)
	space.AddSymbol("$Response.getHeaderNames$default$0", getHeaderNamesDefaultSym)
	getHeaderNamesDefaultRef, _ := space.GetSymbol("$Response.getHeaderNames$default$0")
	ctx.SetSymbolType(getHeaderNamesDefaultRef, libcommon.FunctionSignatureToSemType(env, &getHeaderNamesDefaultSig))

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
	// Signature: preceding params only (url: string), return type = configSemType.
	defaultLambdaSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: configSemType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	defaultLambdaSym := model.NewFunctionSymbol("$Client.init$default$1", defaultLambdaSig, false)
	space.AddSymbol("$Client.init$default$1", defaultLambdaSym)
	defaultLambdaRef, _ := space.GetSymbol("$Client.init$default$1")
	ctx.SetSymbolType(defaultLambdaRef, libcommon.FunctionSignatureToSemType(env, &defaultLambdaSig))

	initSym := model.NewFunctionSymbol("$Client.init", initSig, false)
	space.AddSymbol("$Client.init", initSym)
	initRef, _ := space.GetSymbol("$Client.init")
	ctx.SetSymbolType(initRef, initFnSemType)

	// Mark config (param index 1) as defaultable.
	defaultableInfo := model.NewDefaultableParamInfo(len(initSig.ParamTypes))
	defaultableInfo.SetDefaultable(1, defaultLambdaRef)
	initSym.SetDefaultableParams(defaultableInfo)

	// Default lambda for headers param (index 1): $Client.get$default$1(path) → ()
	getDefaultLambdaSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	getDefaultLambdaSym := model.NewFunctionSymbol("$Client.get$default$1", getDefaultLambdaSig, false)
	space.AddSymbol("$Client.get$default$1", getDefaultLambdaSym)
	getDefaultLambdaRef, _ := space.GetSymbol("$Client.get$default$1")
	ctx.SetSymbolType(getDefaultLambdaRef, libcommon.FunctionSignatureToSemType(env, &getDefaultLambdaSig))

	getSym := model.NewFunctionSymbol("$Client.get", getSig, false)
	space.AddSymbol("$Client.get", getSym)
	getRef, _ := space.GetSymbol("$Client.get")
	ctx.SetSymbolType(getRef, getFnSemType)

	getDefaultableInfo := model.NewDefaultableParamInfo(len(getSig.ParamTypes))
	getDefaultableInfo.SetDefaultable(1, getDefaultLambdaRef)
	getSym.SetDefaultableParams(getDefaultableInfo)

	// post default lambdas: headers at index 2, mediaType at index 3
	postHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	postHeadersDefaultSym := model.NewFunctionSymbol("$Client.post$default$2", postHeadersDefaultSig, false)
	space.AddSymbol("$Client.post$default$2", postHeadersDefaultSym)
	postHeadersDefaultRef, _ := space.GetSymbol("$Client.post$default$2")
	ctx.SetSymbolType(postHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &postHeadersDefaultSig))

	postMediaTypeDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	postMediaTypeDefaultSym := model.NewFunctionSymbol("$Client.post$default$3", postMediaTypeDefaultSig, false)
	space.AddSymbol("$Client.post$default$3", postMediaTypeDefaultSym)
	postMediaTypeDefaultRef, _ := space.GetSymbol("$Client.post$default$3")
	ctx.SetSymbolType(postMediaTypeDefaultRef, libcommon.FunctionSignatureToSemType(env, &postMediaTypeDefaultSig))

	postSym := model.NewFunctionSymbol("$Client.post", postSig, false)
	space.AddSymbol("$Client.post", postSym)
	postRef, _ := space.GetSymbol("$Client.post")
	ctx.SetSymbolType(postRef, postFnSemType)

	postDefaultableInfo := model.NewDefaultableParamInfo(len(postSig.ParamTypes))
	postDefaultableInfo.SetDefaultable(2, postHeadersDefaultRef)
	postDefaultableInfo.SetDefaultable(3, postMediaTypeDefaultRef)
	postSym.SetDefaultableParams(postDefaultableInfo)

	// head: headers at index 1
	headDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	headDefaultSym := model.NewFunctionSymbol("$Client.head$default$1", headDefaultSig, false)
	space.AddSymbol("$Client.head$default$1", headDefaultSym)
	headDefaultRef, _ := space.GetSymbol("$Client.head$default$1")
	ctx.SetSymbolType(headDefaultRef, libcommon.FunctionSignatureToSemType(env, &headDefaultSig))

	headSym := model.NewFunctionSymbol("$Client.head", headSig, false)
	space.AddSymbol("$Client.head", headSym)
	headRef, _ := space.GetSymbol("$Client.head")
	ctx.SetSymbolType(headRef, headFnSemType)
	headDefaultableInfo := model.NewDefaultableParamInfo(len(headSig.ParamTypes))
	headDefaultableInfo.SetDefaultable(1, headDefaultRef)
	headSym.SetDefaultableParams(headDefaultableInfo)

	// options: headers at index 1
	optionsDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	optionsDefaultSym := model.NewFunctionSymbol("$Client.options$default$1", optionsDefaultSig, false)
	space.AddSymbol("$Client.options$default$1", optionsDefaultSym)
	optionsDefaultRef, _ := space.GetSymbol("$Client.options$default$1")
	ctx.SetSymbolType(optionsDefaultRef, libcommon.FunctionSignatureToSemType(env, &optionsDefaultSig))

	optionsSym := model.NewFunctionSymbol("$Client.options", optionsSig, false)
	space.AddSymbol("$Client.options", optionsSym)
	optionsRef, _ := space.GetSymbol("$Client.options")
	ctx.SetSymbolType(optionsRef, optionsFnSemType)
	optionsDefaultableInfo := model.NewDefaultableParamInfo(len(optionsSig.ParamTypes))
	optionsDefaultableInfo.SetDefaultable(1, optionsDefaultRef)
	optionsSym.SetDefaultableParams(optionsDefaultableInfo)

	// put: headers at index 2, mediaType at index 3
	putHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	putHeadersDefaultSym := model.NewFunctionSymbol("$Client.put$default$2", putHeadersDefaultSig, false)
	space.AddSymbol("$Client.put$default$2", putHeadersDefaultSym)
	putHeadersDefaultRef, _ := space.GetSymbol("$Client.put$default$2")
	ctx.SetSymbolType(putHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &putHeadersDefaultSig))

	putMediaTypeDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	putMediaTypeDefaultSym := model.NewFunctionSymbol("$Client.put$default$3", putMediaTypeDefaultSig, false)
	space.AddSymbol("$Client.put$default$3", putMediaTypeDefaultSym)
	putMediaTypeDefaultRef, _ := space.GetSymbol("$Client.put$default$3")
	ctx.SetSymbolType(putMediaTypeDefaultRef, libcommon.FunctionSignatureToSemType(env, &putMediaTypeDefaultSig))

	putSym := model.NewFunctionSymbol("$Client.put", putSig, false)
	space.AddSymbol("$Client.put", putSym)
	putRef, _ := space.GetSymbol("$Client.put")
	ctx.SetSymbolType(putRef, putFnSemType)
	putDefaultableInfo := model.NewDefaultableParamInfo(len(putSig.ParamTypes))
	putDefaultableInfo.SetDefaultable(2, putHeadersDefaultRef)
	putDefaultableInfo.SetDefaultable(3, putMediaTypeDefaultRef)
	putSym.SetDefaultableParams(putDefaultableInfo)

	// patch: headers at index 2, mediaType at index 3
	patchHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	patchHeadersDefaultSym := model.NewFunctionSymbol("$Client.patch$default$2", patchHeadersDefaultSig, false)
	space.AddSymbol("$Client.patch$default$2", patchHeadersDefaultSym)
	patchHeadersDefaultRef, _ := space.GetSymbol("$Client.patch$default$2")
	ctx.SetSymbolType(patchHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &patchHeadersDefaultSig))

	patchMediaTypeDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	patchMediaTypeDefaultSym := model.NewFunctionSymbol("$Client.patch$default$3", patchMediaTypeDefaultSig, false)
	space.AddSymbol("$Client.patch$default$3", patchMediaTypeDefaultSym)
	patchMediaTypeDefaultRef, _ := space.GetSymbol("$Client.patch$default$3")
	ctx.SetSymbolType(patchMediaTypeDefaultRef, libcommon.FunctionSignatureToSemType(env, &patchMediaTypeDefaultSig))

	patchSym := model.NewFunctionSymbol("$Client.patch", patchSig, false)
	space.AddSymbol("$Client.patch", patchSym)
	patchRef, _ := space.GetSymbol("$Client.patch")
	ctx.SetSymbolType(patchRef, patchFnSemType)
	patchDefaultableInfo := model.NewDefaultableParamInfo(len(patchSig.ParamTypes))
	patchDefaultableInfo.SetDefaultable(2, patchHeadersDefaultRef)
	patchDefaultableInfo.SetDefaultable(3, patchMediaTypeDefaultRef)
	patchSym.SetDefaultableParams(patchDefaultableInfo)

	// delete: message at index 1, headers at index 2, mediaType at index 3
	deleteMessageDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: deleteMessageType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	deleteMessageDefaultSym := model.NewFunctionSymbol("$Client.delete$default$1", deleteMessageDefaultSig, false)
	space.AddSymbol("$Client.delete$default$1", deleteMessageDefaultSym)
	deleteMessageDefaultRef, _ := space.GetSymbol("$Client.delete$default$1")
	ctx.SetSymbolType(deleteMessageDefaultRef, libcommon.FunctionSignatureToSemType(env, &deleteMessageDefaultSig))

	deleteHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, deleteMessageType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	deleteHeadersDefaultSym := model.NewFunctionSymbol("$Client.delete$default$2", deleteHeadersDefaultSig, false)
	space.AddSymbol("$Client.delete$default$2", deleteHeadersDefaultSym)
	deleteHeadersDefaultRef, _ := space.GetSymbol("$Client.delete$default$2")
	ctx.SetSymbolType(deleteHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &deleteHeadersDefaultSig))

	deleteMediaTypeDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, deleteMessageType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	deleteMediaTypeDefaultSym := model.NewFunctionSymbol("$Client.delete$default$3", deleteMediaTypeDefaultSig, false)
	space.AddSymbol("$Client.delete$default$3", deleteMediaTypeDefaultSym)
	deleteMediaTypeDefaultRef, _ := space.GetSymbol("$Client.delete$default$3")
	ctx.SetSymbolType(deleteMediaTypeDefaultRef, libcommon.FunctionSignatureToSemType(env, &deleteMediaTypeDefaultSig))

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
	executeHeadersDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, semtypes.STRING, jsonType},
		ReturnType: headersOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	executeHeadersDefaultSym := model.NewFunctionSymbol("$Client.execute$default$3", executeHeadersDefaultSig, false)
	space.AddSymbol("$Client.execute$default$3", executeHeadersDefaultSym)
	executeHeadersDefaultRef, _ := space.GetSymbol("$Client.execute$default$3")
	ctx.SetSymbolType(executeHeadersDefaultRef, libcommon.FunctionSignatureToSemType(env, &executeHeadersDefaultSig))

	executeMediaTypeDefaultSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING, semtypes.STRING, jsonType, headersOptType},
		ReturnType: mediaTypeOptType,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	executeMediaTypeDefaultSym := model.NewFunctionSymbol("$Client.execute$default$4", executeMediaTypeDefaultSig, false)
	space.AddSymbol("$Client.execute$default$4", executeMediaTypeDefaultSym)
	executeMediaTypeDefaultRef, _ := space.GetSymbol("$Client.execute$default$4")
	ctx.SetSymbolType(executeMediaTypeDefaultRef, libcommon.FunctionSignatureToSemType(env, &executeMediaTypeDefaultSig))

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
