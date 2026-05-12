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

var PackageID = model.INTERNAL_PKG

const PackageName = "lang.__internal"

func GetInternalSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	space := ctx.NewSymbolSpace(*PackageID)
	querySortSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{semtypes.LIST, semtypes.LIST, semtypes.LIST, semtypes.LIST},
		ReturnType: semtypes.NIL,
	}
	querySortSymbol := model.NewFunctionSymbol("querySort", querySortSignature, true)
	space.AddSymbol("querySort", querySortSymbol)
	querySortRef, _ := space.GetSymbol("querySort")
	ctx.SetSymbolType(querySortRef, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &querySortSignature))
	return model.NewExportedSymbolSpace(space, nil)
}
