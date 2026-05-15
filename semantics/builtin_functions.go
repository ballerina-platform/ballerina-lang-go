// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
//
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

package semantics

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/lib/registry"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

// tryResolveBuiltinFunction applies compiler-only call-site typing for embedded platform
// functions that export wide symbols from the registry and need monomorphization here.
func tryResolveBuiltinFunction(t typeResolver, chain *binding, inv invocable, polyRef model.SymbolRef, wideSym model.FunctionSymbol, expectedType semtypes.SemType) ([]semtypes.SemType, model.SymbolRef, *binding, bool, bool) {
	if polyRef.Package.Organization != "ballerina" {
		return nil, polyRef, chain, false, false
	}

	fnName := wideSym.Name()
	if call, ok := inv.(*ast.BLangInvocation); ok && call.Name != nil {
		fnName = call.Name.Value
	}

	pkg := polyRef.Package.Package
	tyCtx := semtypes.ContextFrom(t.typeEnv())

	var sig model.FunctionSignature
	switch {
	case pkg == registry.LangArray && fnName == "push":
		containerTy, chain, ok := builtinFirstArg(t, chain, inv, 1, "missing container argument")
		if !ok {
			return nil, polyRef, chain, true, false
		}
		if !semtypes.IsSubtype(tyCtx, containerTy, semtypes.LIST) {
			t.semanticError("expect first argument to be a subtype of (any|error)[]", inv.GetPosition())
			return nil, polyRef, chain, true, false
		}
		sig = model.FunctionSignature{
			ParamTypes:    []semtypes.SemType{containerTy},
			RestParamType: semtypes.ListProj(tyCtx, containerTy, semtypes.INT),
			ReturnType:    semtypes.NIL,
			Flags:         model.FuncSymbolFlagIsolated,
		}

	case pkg == registry.LangMap && fnName == "remove":
		mapTy, chain, ok := builtinFirstArg(t, chain, inv, 2, "missing container or key argument")
		if !ok {
			return nil, polyRef, chain, true, false
		}
		if !semtypes.IsSubtype(tyCtx, mapTy, semtypes.MAPPING) {
			t.semanticError("expect first argument to be a subtype of map<any|error>", inv.GetPosition())
			return nil, polyRef, chain, true, false
		}
		sig = model.FunctionSignature{
			ParamTypes:    []semtypes.SemType{mapTy, semtypes.STRING},
			RestParamType: semtypes.NEVER,
			ReturnType:    semtypes.MappingMemberTypeInnerValProj(tyCtx, mapTy, semtypes.STRING),
			Flags:         model.FuncSymbolFlagIsolated,
		}

	default:
		return nil, polyRef, chain, false, false
	}

	monoName := t.nextMonoFnName(fnName)
	monoSym := model.NewMonomorphicFromPolymorphic(polyRef, monoName, sig, true)
	monoSym.SetType(typeFromFunctionSignature(t, sig))
	scope := t.currentScope()
	scope.AddSymbol(monoName, monoSym)
	monoRef, ok := scope.GetSymbol(monoName)
	if !ok {
		t.internalError("monomorphized "+fnName+" symbol missing from scope", inv.GetPosition())
		return nil, polyRef, chain, true, false
	}
	argTys, chain, ok := argArray(t, monoSym, sig.ParamTypes, sig.RestParamType, chain, inv, expectedType)
	if !ok {
		return nil, polyRef, chain, true, false
	}
	inv.SetResolvedSymbol(monoRef)
	return argTys, monoRef, chain, true, true
}

func builtinFirstArg(t typeResolver, chain *binding, inv invocable, minArgs int, tooFewMsg string) (semtypes.SemType, *binding, bool) {
	args := inv.CallArgs()
	if len(args) < minArgs {
		t.semanticError(tooFewMsg, inv.GetPosition())
		return semtypes.NEVER, chain, false
	}
	ty, _, ok := resolveActionOrExpression(t, chain, args[0], nil)
	return ty, chain, ok
}
