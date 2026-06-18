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

package desugar

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/decimal"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/values"
)

// materializeConstantRef replaces a reference to a folded constant with a
// literal carrying the folded value — the "apply" of the symbol's stored value
// (the E in the design), realized here because model cannot import ast. Returns
// nil for non-foldable constants (e.g. casts that panic at runtime), which keep
// flowing through BIR as global variables.
func materializeConstantRef(cx *functionContext, ref *ast.BLangSimpleVarRef) ast.BLangExpression {
	constSym, ok := cx.getSymbol(ref.Symbol()).(*model.ConstantValueSymbol)
	if !ok {
		return nil
	}
	value, ok := constSym.ConstantValue()
	if !ok || !values.IsSerializableConstValue(value) {
		return nil
	}
	return constantValueLiteral(value, ref.GetPosition(), ref.GetDeterminedType())
}

func constantValueLiteral(value values.BalValue, pos diagnostics.Location, ty semtypes.SemType) ast.BLangExpression {
	if semtypes.IsZero(ty) {
		ty = values.SemTypeForValue(value)
	}

	var expr ast.BLangExpression
	var lit *ast.BLangLiteral
	if isNumericConstantValue(value) {
		numeric := &ast.BLangNumericLiteral{}
		expr = numeric
		lit = &numeric.BLangLiteral
	} else {
		plain := &ast.BLangLiteral{}
		expr = plain
		lit = plain
	}

	lit.SetValue(value)
	lit.SetOriginalValue(values.String(value, make(map[uintptr]bool)))
	lit.SetIsConstant(true)
	lit.SetDeterminedType(ty)
	lit.SetPosition(pos)
	if tag, ok := constantValueTypeTag(value); ok {
		bt := &ast.BTypeBasic{}
		bt.BTypeSetTag(tag)
		lit.SetValueType(bt)
	}
	return expr
}

func isNumericConstantValue(value values.BalValue) bool {
	switch value.(type) {
	case int, int64, int32, int16, int8, byte, float64, float32, *decimal.Decimal:
		return true
	default:
		return false
	}
}

func constantValueTypeTag(value values.BalValue) (ast.TypeTags, bool) {
	switch value.(type) {
	case nil:
		return ast.TypeTags_NIL, true
	case bool:
		return ast.TypeTags_BOOLEAN, true
	case int, int64, int32, int16, int8:
		return ast.TypeTags_INT, true
	case byte:
		return ast.TypeTags_BYTE, true
	case float64, float32:
		return ast.TypeTags_FLOAT, true
	case *decimal.Decimal:
		return ast.TypeTags_DECIMAL, true
	case string, *string:
		return ast.TypeTags_STRING, true
	default:
		return 0, false
	}
}
