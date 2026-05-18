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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

const PackageName = "lang.xml"

var XMLPackageID = model.NewPackageID(
	model.DefaultPackageIDInterner,
	model.Name("ballerina"),
	[]model.Name{model.Name("lang"), model.Name("xml")},
	model.Name("0.0.1"),
)

func GetXMLSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	type xmlType struct {
		name string
		ty   semtypes.SemType
	}
	types := []xmlType{
		{"Element", semtypes.XML_ELEMENT},
		{"Comment", semtypes.XML_COMMENT},
		{"Text", semtypes.XML_TEXT},
		{"ProcessingInstruction", semtypes.XML_PI},
	}
	space := ctx.NewSymbolSpace(*XMLPackageID)
	for _, each := range types {
		tySym := model.NewTypeSymbol(each.name, true)
		space.AddSymbol(each.name, &tySym)
		ref, _ := space.GetSymbol(each.name)
		ctx.SetSymbolType(ref, each.ty)
	}
	return model.NewExportedSymbolSpace(space, nil)
}
