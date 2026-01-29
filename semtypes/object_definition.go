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

// Represent object type desc.
//
// migrated from ObjectDefinition.java:50:1
type ObjectDefinition struct {
	// migrated from ObjectDefinition.java:52:5
	mappingDefinition MappingDefinition
}

// migrated from ObjectDefinition.java:50:1
var _ Definition = &ObjectDefinition{}

func NewObjectDefinition() ObjectDefinition {
	this := ObjectDefinition{}
	this.mappingDefinition = NewMappingDefinition()
	return this
}

// migrated from ObjectDefinition.java:54:5
func ObjectDefinitionDistinct(distinctId int) SemType {
	// migrated from ObjectDefinition.java:55:9
	common.Assert(distinctId >= 0)
	bdd := BddAtom(common.ToPointer(CreateDistinctRecAtom(-distinctId - 1)))
	return basicSubtype(BT_OBJECT, bdd)
}

// Each object type is represented as mapping type (with its basic type set to object) as fallows
//
//	{
//	  "$qualifiers": {
//	    boolean isolated,
//	    "client"|"service" network
//	  },
//	   [field_name]: {
//	     "field"|"method" kind,
//	     "public"|"private" visibility,
//	      VAL value;
//	   }
//	   ...{
//	     "field" kind,
//	     "public"|"private" visibility,
//	      VAL value;
//	   } | {
//	      "method" kind,
//	      "public"|"private" visibility,
//	      FUNCTION value;
//	   }
//	}
//
// migrated from ObjectDefinition.java:81:5
func (this *ObjectDefinition) Define(env Env, qualifiers ObjectQualifiers, members []Member) SemType {
	common.Assert(objectDefinitionValidateMembers(members))
	// migrated from ObjectDefinition.java:82:9
	var mut CellMutability
	if qualifiers.readonly {
		mut = CellMutability_CELL_MUT_NONE
	} else {
		mut = CellMutability_CELL_MUT_LIMITED
	}
	// migrated from ObjectDefinition.java:85:9
	var memberStream []CellField
	for _, member := range members {
		memberStream = append(memberStream, memberField(env, &member, mut))
	}
	// migrated from ObjectDefinition.java:87:9
	qualifierStream := []CellField{qualifiers.Field(env)}
	// migrated from ObjectDefinition.java:88:9
	var cellFields []CellField
	cellFields = append(cellFields, memberStream...)
	cellFields = append(cellFields, qualifierStream...)
	mappingType := this.mappingDefinition.Define(env, cellFields, this.restMemberType(env, mut, qualifiers.readonly))
	return this.objectContaining(mappingType)
}

// migrated from ObjectDefinition.java:93:5
func objectDefinitionValidateMembers(members []Member) bool {
	// Check if there are two members with same name
	// migrated from ObjectDefinition.java:95:9
	nameMap := make(map[string]bool)
	for _, member := range members {
		if nameMap[member.Name] {
			return false
		}
		nameMap[member.Name] = true
	}
	return len(nameMap) == len(members)
}

// migrated from ObjectDefinition.java:98:5
func (this *ObjectDefinition) objectContaining(mappingType SemType) SemType {
	// migrated from ObjectDefinition.java:99:9
	bdd := subtypeData(mappingType, BT_MAPPING)
	// migrated from ObjectDefinition.java:100:9
	return CreateBasicSemType(BT_OBJECT, bdd)
}

// migrated from ObjectDefinition.java:104:5
func (this *ObjectDefinition) restMemberType(env Env, mut CellMutability, immutable bool) CellSemType {
	// migrated from ObjectDefinition.java:105:9
	fieldDefn := NewMappingDefinition()
	// migrated from ObjectDefinition.java:106:9
	var fieldValueTy SemType
	if immutable {
		fieldValueTy = VAL_READONLY
	} else {
		fieldValueTy = &VAL
	}
	fieldMemberType := fieldDefn.DefineMappingTypeWrapped(
		env,
		[]Field{
			FieldFrom("value", fieldValueTy, immutable, false),
			common.ToPointer(MemberKindField).field(),
			visibilityAll,
		},
		&NEVER)

	// migrated from ObjectDefinition.java:116:9
	methodDefn := NewMappingDefinition()
	// migrated from ObjectDefinition.java:117:9
	methodMemberType := methodDefn.DefineMappingTypeWrapped(
		env,
		[]Field{
			FieldFrom("value", &FUNCTION, true, false),
			common.ToPointer(MemberKindMethod).field(),
			visibilityAll,
		},
		&NEVER)
	return CellContainingWithEnvSemTypeCellMutability(env, Union(fieldMemberType, methodMemberType), mut)
}

// migrated from ObjectDefinition.java:128:5
func memberField(env Env, member *Member, mut CellMutability) CellField {
	// migrated from ObjectDefinition.java:129:9
	md := NewMappingDefinition()
	// migrated from ObjectDefinition.java:130:9
	var fieldMut CellMutability
	if member.Immutable {
		fieldMut = CellMutability_CELL_MUT_NONE
	} else {
		fieldMut = mut
	}
	// migrated from ObjectDefinition.java:131:9
	semtype := md.DefineMappingTypeWrapped(
		env,
		[]Field{
			FieldFrom("value", member.ValueTy, member.Immutable, false),
			(&member.Kind).field(),
			(&member.Visibility).field(),
		},
		&NEVER)
	return CellFieldFrom(member.Name, CellContainingWithEnvSemTypeCellMutability(env, semtype, fieldMut))
}

// migrated from ObjectDefinition.java:143:5
func (this *ObjectDefinition) GetSemType(env Env) SemType {
	// migrated from ObjectDefinition.java:144:9
	return this.objectContaining(this.mappingDefinition.GetSemType(env))
}
