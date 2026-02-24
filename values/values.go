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

package values

import (
	"ballerina-lang-go/semtypes"
	"math/big"
)

// Currently this is just an alias on any but I think we will need to add methods to this like type
type BalValue any

func DefaultValueForType(t semtypes.SemType) BalValue {
	if t == nil {
		// TODO: this should panic when our operands properly have types
		return nil
	}
	if semtypes.IsNever(t) {
		return NeverValue
	} else if semtypes.IsSubtypeSimple(t, semtypes.BOOLEAN) {
		return false
	} else if semtypes.IsSubtypeSimple(t, semtypes.INT) {
		return int64(0)
	} else if semtypes.IsSubtypeSimple(t, semtypes.FLOAT) {
		return float64(0)
	} else if semtypes.IsSubtypeSimple(t, semtypes.STRING) {
		return ""
	} else if semtypes.IsSubtypeSimple(t, semtypes.DECIMAL) {
		return big.NewRat(0, 1)
	} else if semtypes.IsSubtypeSimple(t, semtypes.LIST) {
		// TODO: this needs to be properly implemeneted for lists
		return NewList(0, &semtypes.NEVER, NeverValue)
	} else if semtypes.ContainsBasicType(t, semtypes.NIL) {
		return nil
	} else {
		return NeverValue
	}
}

func SemTypeForValue(v BalValue) semtypes.SemType {
	switch v := v.(type) {
	case nil:
		return &semtypes.NIL
	case bool:
		return semtypes.BooleanConst(v)
	case int64:
		return semtypes.IntConst(v)
	case float64:
		return semtypes.FloatConst(v)
	case string:
		return semtypes.StringConst(v)
	case *big.Rat:
		return semtypes.DecimalConst(*v)
	case *List:
		return v.Type
	default:
		return &semtypes.ANY
	}
}
