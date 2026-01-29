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
)

var SINT8 = IntWidthSigned(8)
var SINT16 = IntWidthSigned(16)
var SINT32 = IntWidthSigned(32)
var UINT8 = BYTE
var UINT16 = IntWidthUnsigned(16)
var UINT32 = IntWidthUnsigned(32)
var CHAR = STRING_CHAR

func DecimalConstFromStringValue(value string) SemType {
	// migrated from SemTypes.java:68:5
	if strings.Contains(value, "d") || strings.Contains(value, "D") {
		value = value[:len(value)-1]
	}
	d := new(big.Rat)
	d, ok := d.SetString(value)
	if !ok {
		panic("failed to set string to big.Rat")
	}
	return DecimalConst(*d)
}

func UnionWithSemTypeSemTypesSemType(first SemType, second SemType, rest ...SemType) SemType {
	// migrated from SemTypes.java:80:5
	u := Union(first, second)
	for _, s := range rest {
		u = Union(u, s)
	}
	return u
}

func IntersectWithSemTypeSemTypesSemType(first SemType, second SemType, rest ...SemType) SemType {
	// migrated from SemTypes.java:92:5
	i := Intersect(first, second)
	for _, s := range rest {
		i = Intersect(i, s)
	}
	return i
}

func IsSubtypeSimpleNotNever(t1 SemType, t2 BasicTypeBitSet) bool {
	// migrated from SemTypes.java:108:5
	return ((!IsNever(t1)) && IsSubtypeSimple(t1, t2))
}

func ContainsBasicType(t1 SemType, t2 BasicTypeBitSet) bool {
	// migrated from SemTypes.java:112:5
	return ((widenToBasicTypes(t1).bitset & t2.bitset) != 0)
}

func ContainsType(context Context, ty SemType, typeToBeContained SemType) bool {
	// migrated from SemTypes.java:116:5
	return IsSameType(context, Intersect(ty, typeToBeContained), typeToBeContained)
}

func ListProj(context Context, t SemType, key SemType) SemType {
	// migrated from SemTypes.java:160:5
	return ListProjInnerVal(context, t, key)
}

func ListMemberType(context Context, t SemType, key SemType) SemType {
	// migrated from SemTypes.java:164:5
	return listMemberTypeInnerVal(context, t, key)
}
