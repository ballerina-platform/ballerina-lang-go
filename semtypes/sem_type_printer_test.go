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

package semtypes

import (
	"math/big"
	"strings"
	"testing"
)

type fakeSemType struct{}

func (fakeSemType) All() int {
	return 0
}

func (fakeSemType) String() string {
	return "fake-semtype"
}

type customProperSubtype struct{}

func (c customProperSubtype) String() string {
	return "customProperSubtype"
}

func TestSemTypeStringAddsSubtypeDetails(t *testing.T) {
	rendered := String(nil, IntConst(42))

	assertTrue(t, strings.Contains(rendered, "((), (INT))"))
	assertTrue(t, strings.Contains(rendered, "subtypes=[INT:ranges=[42..42]]"))
}

func TestSemTypeStringUsesListMembersWhenContextAvailable(t *testing.T) {
	env := CreateTypeEnv()
	ctx := ContextFrom(env)

	ld := NewListDefinition()
	listOfInt := ld.DefineListTypeWrappedWithEnvSemType(env, &INT)
	rendered := String(ctx, listOfInt)

	assertTrue(t, strings.Contains(rendered, "listMembers=[0..*:((INT), ())]"))
	assertFalse(t, strings.Contains(rendered, "subtypes=[LIST:"))
}

func TestSemTypeStringAndCompactStringBranches(t *testing.T) {
	assertEqual(t, String(nil, nil), "<UNKNOWN>")
	assertEqual(t, CompactString(nil), "<UNKNOWN>")

	assertEqual(t, CompactString(&INT), "((INT), ())")
	assertTrue(t, strings.Contains(CompactString(IntConst(5)), "((), (INT))"))

	var fake SemType = fakeSemType{}
	assertEqual(t, CompactString(fake), "fake-semtype")
	assertEqual(t, String(nil, fake), "fake-semtype")
}

func TestRenderSubtypeDataVariants(t *testing.T) {
	booleanSubtype := BooleanSubtypeFrom(true)
	assertTrue(t, strings.Contains(renderSubtypeData(booleanSubtype), "value=true"))
	assertTrue(t, strings.Contains(renderSubtypeData(&booleanSubtype), "value=true"))

	floatValue := EnumerableFloatFrom(3.5)
	floatSubtype := CreateFloatSubtype(true, []EnumerableType[float64]{&floatValue}).(FloatSubtype)
	assertTrue(t, strings.Contains(renderSubtypeData(floatSubtype), "allowed=true"))
	assertTrue(t, strings.Contains(renderSubtypeData(&floatSubtype), "values=[3.5]"))

	decimalRat := new(big.Rat).SetFrac64(1, 3)
	decimalValue := EnumerableDecimalFrom(*decimalRat)
	decimalSubtype := CreateDecimalSubtype(true, []EnumerableType[big.Rat]{&decimalValue}).(DecimalSubtype)
	assertTrue(t, strings.Contains(renderSubtypeData(decimalSubtype), "allowed=true"))
	assertTrue(t, strings.Contains(renderSubtypeData(&decimalSubtype), "values=[1/3]"))

	stringSubtype := StringSubtypeFrom(
		CharStringSubtypeFrom(true, []EnumerableType[string]{EnumerableCharStringFrom("x")}),
		NonCharStringSubtypeFrom(true, []EnumerableType[string]{EnumerableStringFrom("hello")}),
	)
	assertTrue(t, strings.Contains(renderSubtypeData(stringSubtype), "char(allowed=true"))
	assertTrue(t, strings.Contains(renderSubtypeData(&stringSubtype), "nonChar(allowed=true"))
	assertTrue(t, strings.Contains(renderSubtypeData(&stringSubtype), "hello"))

	assertEqual(t, renderSubtypeData(BddAll()), "bdd=true")

	customRendered := renderSubtypeData(customProperSubtype{})
	assertTrue(t, strings.Contains(customRendered, "customProperSubtype"))
}

func TestRenderSubtypeDetailsGuardsAndSkipMask(t *testing.T) {
	assertEqual(t, renderSubtypeDetails(0, nil, 0), "")

	intMask := 1 << BT_INT.Code
	intSubtype := CreateSingleRangeSubtype(1, 2)
	assertEqual(t, renderSubtypeDetails(intMask, []ProperSubtypeData{intSubtype}, intMask), "")

	details := renderSubtypeDetails(intMask, []ProperSubtypeData{intSubtype}, 0)
	assertTrue(t, strings.Contains(details, "INT:ranges=[1..2]"))
}
