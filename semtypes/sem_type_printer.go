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
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

func String(cx Context, ty SemType) string {
	base := CompactString(ty)
	if ty == nil {
		return base
	}
	details := semTypeDetails(cx, ty)
	if len(details) == 0 {
		return base
	}
	return fmt.Sprintf("%s {%s}", base, strings.Join(details, ", "))
}

func CompactString(ty SemType) string {
	if ty == nil {
		return "<UNKNOWN>"
	}
	if bitset, ok := ty.(*BasicTypeBitSet); ok {
		return fmt.Sprintf("((%s), ())", bitsetToTypeNames(bitset.bitset))
	}
	if complexTy, ok := ty.(ComplexSemType); ok {
		return fmt.Sprintf("((%s), (%s))", bitsetToTypeNames(complexTy.All()), bitsetToTypeNames(complexTy.Some()))
	}
	return fmt.Sprintf("%v", ty)
}

func semTypeDetails(cx Context, ty SemType) []string {
	details := make([]string, 0, 2)
	subtypeSkipMask := 0

	if cx != nil && IsSubtypeSimple(ty, LIST) {
		listMemberDetails := renderListMemberDetails(cx, ty)
		if listMemberDetails != "" {
			details = append(details, "listMembers=["+listMemberDetails+"]")
			subtypeSkipMask = 1 << BT_LIST.Code
		}
	}

	if complexTy, ok := ty.(ComplexSemType); ok {
		subtypeDetails := renderSubtypeDetails(complexTy.Some(), complexTy.SubtypeDataList(), subtypeSkipMask)
		if subtypeDetails != "" {
			details = append(details, "subtypes=["+subtypeDetails+"]")
		}
	}

	return details
}

func renderSubtypeDetails(some int, subtypeDataList []ProperSubtypeData, skipMask int) string {
	if some == 0 || len(subtypeDataList) == 0 {
		return ""
	}

	details := make([]string, 0, len(subtypeDataList))
	mask := some
	for i := 0; i < len(subtypeDataList) && mask != 0; i++ {
		bit := bits.TrailingZeros(uint(mask))
		if bit >= VT_COUNT {
			break
		}
		code := BasicTypeCodeFrom(bit)
		if skipMask&(1<<bit) == 0 {
			codeName := strings.TrimPrefix(code.String(), "BT_")
			details = append(details, fmt.Sprintf("%s:%s", codeName, renderSubtypeData(subtypeDataList[i])))
		}
		mask &= ^(1 << bit)
	}
	return strings.Join(details, ", ")
}

func renderSubtypeData(data ProperSubtypeData) string {
	switch d := data.(type) {
	case IntSubtype:
		return renderIntSubtype(d)
	case *IntSubtype:
		return renderIntSubtype(*d)
	case BooleanSubtype:
		return fmt.Sprintf("value=%t", d.value)
	case *BooleanSubtype:
		return fmt.Sprintf("value=%t", d.value)
	case FloatSubtype:
		return renderFloatSubtype(d)
	case *FloatSubtype:
		return renderFloatSubtype(*d)
	case DecimalSubtype:
		return renderDecimalSubtype(d)
	case *DecimalSubtype:
		return renderDecimalSubtype(*d)
	case StringSubtype:
		return renderStringSubtype(d)
	case *StringSubtype:
		return renderStringSubtype(*d)
	case Bdd:
		return "bdd=" + d.canonicalKey()
	default:
		return fmt.Sprintf("%T", data)
	}
}

func renderIntSubtype(subtype IntSubtype) string {
	ranges := make([]string, 0, len(subtype.Ranges))
	for _, r := range subtype.Ranges {
		max := fmt.Sprintf("%d", r.Max)
		if r.Max == MAX_VALUE {
			max = "*"
		}
		ranges = append(ranges, fmt.Sprintf("%d..%s", r.Min, max))
	}
	return "ranges=[" + strings.Join(ranges, ", ") + "]"
}

func renderFloatSubtype(subtype FloatSubtype) string {
	values := make([]string, 0, len(subtype.values))
	for _, v := range subtype.values {
		values = append(values, fmt.Sprintf("%v", v.value))
	}
	sort.Strings(values)
	return fmt.Sprintf("allowed=%t, values=[%s]", subtype.allowed, strings.Join(values, ", "))
}

func renderDecimalSubtype(subtype DecimalSubtype) string {
	values := make([]string, 0, len(subtype.values))
	for _, v := range subtype.values {
		values = append(values, v.value.RatString())
	}
	sort.Strings(values)
	return fmt.Sprintf("allowed=%t, values=[%s]", subtype.allowed, strings.Join(values, ", "))
}

func renderStringSubtype(subtype StringSubtype) string {
	charValues := renderEnumerableStringValues(subtype.charData.values)
	nonCharValues := renderEnumerableStringValues(subtype.nonCharData.values)
	return fmt.Sprintf("char(allowed=%t, values=[%s]), nonChar(allowed=%t, values=[%s])",
		subtype.charData.allowed, charValues, subtype.nonCharData.allowed, nonCharValues)
}

func renderEnumerableStringValues(values []EnumerableType[string]) string {
	rendered := make([]string, 0, len(values))
	for _, each := range values {
		rendered = append(rendered, each.Value())
	}
	sort.Strings(rendered)
	return strings.Join(rendered, ", ")
}

func renderListMemberDetails(cx Context, ty SemType) string {
	memberTypes := ListAllMemberTypesInner(cx, ty)
	if len(memberTypes.SemTypes) == 0 {
		return ""
	}
	details := make([]string, 0, len(memberTypes.SemTypes))
	for i := range memberTypes.SemTypes {
		r := memberTypes.Ranges[i]
		rangeLabel := fmt.Sprintf("%d..%d", r.Min, r.Max)
		if r.Max == MAX_VALUE {
			rangeLabel = fmt.Sprintf("%d..*", r.Min)
		}
		details = append(details, fmt.Sprintf("%s:%v", rangeLabel, memberTypes.SemTypes[i]))
	}
	sort.Strings(details)
	return strings.Join(details, ", ")
}

// bitsetToTypeNames converts a bitset to a comma-separated list of type names.
// Returns empty string for empty bitset (0).
// Type names are returned in bitset order without the "BT_" prefix.
func bitsetToTypeNames(bitset int) string {
	if bitset == 0 {
		return ""
	}

	var builder strings.Builder
	first := true

	for i := 0; i < VT_COUNT; i++ {
		if (bitset & (1 << i)) != 0 {
			if !first {
				builder.WriteString(", ")
			}
			typeCode := BasicTypeCodeFrom(i)
			typeName := typeCode.String()
			// Strip "BT_" prefix for clean output
			cleanName := strings.TrimPrefix(typeName, "BT_")
			builder.WriteString(cleanName)
			first = false
		}
	}

	return builder.String()
}
