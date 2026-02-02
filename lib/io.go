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

package lib

import (
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

func GetIoSymbols(ctx *context.CompilerContext, env semtypes.Env) model.ExportedSymbolSpace {
	pkg := model.NewPackageID(
		model.DefaultPackageIDInterner,
		model.Name("ballerina"),
		[]model.Name{model.Name("io")},
		model.Name("0.0.1"),
	)
	space := ctx.NewSymbolSpace(*pkg)
	printLnSignature := model.FunctionSignature{
		RestParamType: &semtypes.ANY,
		ReturnType:    &semtypes.NIL,
	}
	printLnSymbol := model.NewFunctionSymbol("println", printLnSignature, true)
	printLnSymbol.SetType(functionSignatureToSemType(env, &printLnSignature))

	space.AddSymbol("println", &printLnSymbol)
	return model.ExportedSymbolSpace{
		Main: space,
	}
}
