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

import "ballerina-lang-go/common"

type TypedescSubtype struct {
}

func newTypedescSubtype() TypedescSubtype {
	this := TypedescSubtype{}
	return this
}

func TypedescContaining(env Env, constraint SemType) SemType {
	// migrated from TypedescSubtype.java:44:5
	if common.PointerEqualToValue(VAL, constraint) {
		return &TYPEDESC
	}

	mappingDef := NewMappingDefinition()
	mappingType := mappingDef.DefineMappingTypeWrappedWithEnvFieldsSemTypeCellMutability(env, nil, constraint, CellMutability_CELL_MUT_NONE)
	bdd := subtypeData(mappingType, BT_MAPPING).(Bdd)
	return CreateBasicSemType(BT_TYPEDESC, bdd)
}
