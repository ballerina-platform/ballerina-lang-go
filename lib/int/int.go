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

package int

import (
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

const PackageName = "lang.int"

var ArrayPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("int")},
	model.Name("0.0.1"),
)

func GetArraySymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	type intType struct {
		name string
		ty   semtypes.SemType
	}
	types := []intType{
		{"Signed8", semtypes.SINT8},
		{"Signed16", semtypes.SINT16},
		{"Signed32", semtypes.SINT32},
		{"Unsigned8", semtypes.UINT8},
		{"Unsigned16", semtypes.UINT16},
		{"Unsigned32", semtypes.UINT32},
	}
	space := ctx.NewSymbolSpace(*ArrayPackageID)
	for _, each := range types {
		tySym := model.NewTypeSymbol(each.name, true)
		tySym.SetType(each.ty)
		space.AddSymbol(each.name, &tySym)
	}

	return model.ExportedSymbolSpace{
		Main: space,
	}
}
