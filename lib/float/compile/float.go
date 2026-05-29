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

const PackageName = "lang.float"

var floatPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("float")},
	model.Name("0.0.1"),
)

func GetFloatSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	space := ctx.NewSymbolSpace(*floatPackageID)

	addFloatToFloat := func(name string) {
		sig := model.FunctionSignature{
			ParamTypes: []semtypes.SemType{semtypes.FLOAT},
			ReturnType: semtypes.FLOAT,
			Flags:      model.FuncSymbolFlagIsolated,
		}
		sym := model.NewFunctionSymbol(name, sig, true)
		space.AddSymbol(name, sym)
		ref, _ := space.GetSymbol(name)
		ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &sig))
	}

	addFloatFloatToFloat := func(name string) {
		sig := model.FunctionSignature{
			ParamTypes: []semtypes.SemType{semtypes.FLOAT, semtypes.FLOAT},
			ReturnType: semtypes.FLOAT,
			Flags:      model.FuncSymbolFlagIsolated,
		}
		sym := model.NewFunctionSymbol(name, sig, true)
		space.AddSymbol(name, sym)
		ref, _ := space.GetSymbol(name)
		ctx.SetSymbolType(ref, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &sig))
	}

	addFloatToFloat("abs")
	addFloatToFloat("ceiling")
	addFloatToFloat("floor")
	addFloatToFloat("round")
	addFloatToFloat("sqrt")
	addFloatFloatToFloat("pow")

	return model.NewExportedSymbolSpace(space, nil)
}
