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

package registry

import (
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

// Embedded .sym files store langlib type aliases (e.g. `public type Char string`) as their
// underlying type. LoadSymbols calls applyExportedTypeSemtypes to set the intended semtypes
// on exported TypeSymbols so Char and the int width aliases type-check as distinct subtypes.
var exportedTypeSemtypes = map[string]map[string]semtypes.SemType{
	"lang.string": {
		"Char": semtypes.CHAR,
	},
	"lang.int": {
		"Signed8":    semtypes.SINT8,
		"Signed16":   semtypes.SINT16,
		"Signed32":   semtypes.SINT32,
		"Unsigned8":  semtypes.UINT8,
		"Unsigned16": semtypes.UINT16,
		"Unsigned32": semtypes.UINT32,
	},
}

func applyExportedTypeSemtypes(moduleName string, exp model.ExportedSymbolSpace) {
	types, ok := exportedTypeSemtypes[moduleName]
	if !ok || exp.Main == nil {
		return
	}
	main := exp.Main
	for name, ty := range types {
		ref, ok := main.GetSymbol(name)
		if !ok {
			continue
		}
		sym := main.SymbolAt(ref.Index)
		if ts, ok := sym.(*model.TypeSymbol); ok {
			ts.SetType(ty)
		}
	}
}
