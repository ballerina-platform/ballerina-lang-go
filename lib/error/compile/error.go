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

var ErrorPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("error")},
	model.Name("0.0.1"),
)

const PackageName = "lang.error"

func GetErrorSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	env := ctx.GetTypeEnv()
	space := ctx.NewSymbolSpace(*ErrorPackageID)

	// message: (ERROR) -> STRING
	messageSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.ERROR},
		ReturnType: semtypes.STRING,
		Flags:      model.FuncSymbolFlagIsolated,
	}
	messageSymbol := model.NewFunctionSymbol("message", messageSignature, true)
	space.AddSymbol("message", messageSymbol)
	messageRef, _ := space.GetSymbol("message")
	ctx.SetSymbolType(messageRef, libcommon.FunctionSignatureToSemType(env, &messageSignature))

	return model.NewExportedSymbolSpace(space, nil)
}
