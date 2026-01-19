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

package symbols

type SymbolOrigin uint8

const (
	SYMBOL_ORIGIN_BUILTIN         SymbolOrigin = 1
	SYMBOL_ORIGIN_SOURCE          SymbolOrigin = 2
	SYMBOL_ORIGIN_COMPILED_SOURCE SymbolOrigin = 3
	SYMBOL_ORIGIN_VIRTUAL         SymbolOrigin = 4
)

func (s SymbolOrigin) ToBIROrigin() SymbolOrigin {
	switch s {
	case SYMBOL_ORIGIN_BUILTIN:
		return SYMBOL_ORIGIN_BUILTIN
	case SYMBOL_ORIGIN_SOURCE:
		return SYMBOL_ORIGIN_COMPILED_SOURCE
	case SYMBOL_ORIGIN_COMPILED_SOURCE:
		return SYMBOL_ORIGIN_COMPILED_SOURCE
	case SYMBOL_ORIGIN_VIRTUAL:
		return SYMBOL_ORIGIN_VIRTUAL
	default:
		return s
	}
}

func (s SymbolOrigin) Value() uint8 {
	return uint8(s)
}

func ToOrigin(value uint8) SymbolOrigin {
	switch value {
	case 1:
		return SYMBOL_ORIGIN_BUILTIN
	case 2:
		return SYMBOL_ORIGIN_SOURCE
	case 3:
		return SYMBOL_ORIGIN_COMPILED_SOURCE
	case 4:
		return SYMBOL_ORIGIN_VIRTUAL
	default:
		panic("Invalid symbol origin value")
	}
}
