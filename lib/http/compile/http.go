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

	parseHeaderSig := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.STRING},
		ReturnType: semtypes.Union(semtypes.LIST, semtypes.ERROR),
		Flags:      model.FuncSymbolFlagIsolated,
	}
	parseHeaderSymbol := model.NewFunctionSymbol("parseHeader", parseHeaderSig, true)
	space.AddSymbol("parseHeader", parseHeaderSymbol)
	parseHeaderRef, _ := space.GetSymbol("parseHeader")
	ctx.SetSymbolType(parseHeaderRef, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &parseHeaderSig))

	return model.NewExportedSymbolSpace(space, nil)
}
