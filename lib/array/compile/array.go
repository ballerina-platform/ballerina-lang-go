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
	"fmt"
	"sync"
)

var ArrayPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("array")},
	model.Name("0.0.1"),
)

const PackageName = "lang.array"

func GetArraySymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	space := ctx.NewSymbolSpace(*ArrayPackageID)
	pushSymbol := model.NewGenericFunctionSymbol("push", space, createPushMonomorphizer(ctx))
	space.AddSymbol("push", pushSymbol)
	lenghtSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{&semtypes.LIST},
		ReturnType: &semtypes.INT,
	}

	lengthSymbol := model.NewFunctionSymbol("length", lenghtSignature, true)
	space.AddSymbol("length", lengthSymbol)
	lengthRef, _ := space.GetSymbol("length")
	ctx.SetSymbolType(lengthRef, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &lenghtSignature))
	return model.ExportedSymbolSpace{
		Main: space,
	}
}

func createPushMonomorphizer(ctx *context.CompilerContext) func(s model.GenericFunctionSymbol, args []semtypes.SemType) model.SymbolRef {
	var mut sync.Mutex
	monomorphized := make(map[semtypes.SemType]model.SymbolRef)
	nextIndex := 0

	return func(s model.GenericFunctionSymbol, args []semtypes.SemType) model.SymbolRef {

		if len(args) == 0 {
			ctx.SemanticError("push() requires at least 1 argument", nil)
		}
		ty := args[0]
		mut.Lock()
		defer mut.Unlock()
		if _, ok := monomorphized[ty]; ok {
			return monomorphized[ty]
		}
		topType := &semtypes.LIST
		tyCtx := semtypes.ContextFrom(ctx.GetTypeEnv())
		if !semtypes.IsSubtype(tyCtx, ty, topType) {
			ctx.SemanticError("expect first argument to be a subtype of (any|error)[]", nil)
		}
		// Is this is correct or do we need to take the list atomic type for this?
		valType := semtypes.ListProj(tyCtx, ty, semtypes.IntConst(0))
		pushSignature := model.FunctionSignature{
			ParamTypes:    []semtypes.SemType{ty},
			RestParamType: valType,
			ReturnType:    &semtypes.NIL,
		}
		pushSymbol := model.NewFunctionSymbol("push", pushSignature, true)
		symbolName := fmt.Sprintf("push_%d", nextIndex)
		nextIndex++
		space := s.Space()
		space.AddSymbol(symbolName, pushSymbol)
		ref, _ := space.GetSymbol(symbolName)
		ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &pushSignature))
		monomorphized[ty] = ref
		return ref
	}
}
