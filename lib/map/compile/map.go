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
	"fmt"
	"sync"

	"ballerina-lang-go/context"
	libcommon "ballerina-lang-go/lib/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

var MapPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("map")},
	model.Name("0.0.1"),
)

const PackageName = "lang.map"

func GetMapSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	env := ctx.GetTypeEnv()
	space := ctx.NewSymbolSpace(*MapPackageID)

	// length: (MAPPING) -> INT
	lengthSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.MAPPING},
		ReturnType: semtypes.INT,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	lengthSymbol := model.NewFunctionSymbol("length", lengthSignature, true)
	space.AddSymbol("length", lengthSymbol)
	lengthRef, _ := space.GetSymbol("length")
	ctx.SetSymbolType(lengthRef, libcommon.FunctionSignatureToSemType(env, &lengthSignature))

	// keys: (MAPPING) -> string[]
	stringArrayLd := semtypes.NewListDefinition()
	stringArrayTy := stringArrayLd.DefineListTypeWrappedWithEnvSemType(env, semtypes.STRING)
	keysSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.MAPPING},
		ReturnType: stringArrayTy,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	keysSymbol := model.NewFunctionSymbol("keys", keysSignature, true)
	space.AddSymbol("keys", keysSymbol)
	keysRef, _ := space.GetSymbol("keys")
	ctx.SetSymbolType(keysRef, libcommon.FunctionSignatureToSemType(env, &keysSignature))

	// remove: generic (mapType, STRING) -> memberType
	removeSymbol := model.NewGenericFunctionSymbol("remove", space, []string{"m", "k"}, createRemoveMonomorphizer(ctx))
	space.AddSymbol("remove", removeSymbol)

	return model.NewExportedSymbolSpace(space, nil)
}

func createRemoveMonomorphizer(ctx *context.CompilerContext) func(s model.GenericFunctionSymbol, args []semtypes.SemType) model.SymbolRef {
	var mut sync.Mutex
	monomorphized := make(map[semtypes.SemType]model.SymbolRef)
	nextIndex := 0

	return func(s model.GenericFunctionSymbol, args []semtypes.SemType) model.SymbolRef {
		if len(args) == 0 {
			ctx.SemanticError("remove() requires at least 1 argument", nil)
			return model.SymbolRef{}
		}
		ty := args[0]
		mut.Lock()
		defer mut.Unlock()
		if ref, ok := monomorphized[ty]; ok {
			return ref
		}
		tyCtx := semtypes.ContextFrom(ctx.GetTypeEnv())
		if !semtypes.IsSubtype(tyCtx, ty, semtypes.MAPPING) {
			ctx.SemanticError("expect first argument to be a subtype of map<any|error>", nil)
			return model.SymbolRef{}
		}
		memberType := semtypes.MappingMemberTypeInnerValProj(tyCtx, ty, semtypes.STRING)
		removeSignature := model.FunctionSignature{
			ParamTypes: []semtypes.SemType{ty, semtypes.STRING},
			ReturnType: memberType,
			Flags:      model.FuncSymbolFlagIsolated,
		}
		removeSymbol := model.NewFunctionSymbol("remove", removeSignature, true)
		symbolName := fmt.Sprintf("remove_%d", nextIndex)
		nextIndex++
		space := s.Space()
		space.AddSymbol(symbolName, removeSymbol)
		ref, _ := space.GetSymbol(symbolName)
		ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &removeSignature))
		monomorphized[ty] = ref
		return ref
	}
}
