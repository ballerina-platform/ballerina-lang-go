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

package ast

import "ballerina-lang-go/model"

// TODO: think of a better way and place to put this
type BTypeSymbolTable struct {
}

var _ model.TypeSymbolTable = &BTypeSymbolTable{}

var (
	booleanType = &BTypeImpl{
		tag:   model.TypeTags_BOOLEAN,
		flags: Flags_READONLY,
	}
	intType = &BTypeImpl{
		tag:   model.TypeTags_INT,
		flags: Flags_READONLY,
	}

	nilType = &BTypeImpl{
		tag:   model.TypeTags_NIL,
		flags: Flags_READONLY,
	}
	stringType = &BTypeImpl{
		tag:   model.TypeTags_STRING,
		flags: Flags_READONLY,
	}
	floatType = &BTypeImpl{
		tag:   model.TypeTags_FLOAT,
		flags: Flags_READONLY,
	}
)

func (this *BTypeSymbolTable) GetTypeFromTag(tag model.TypeTags) model.TypeNode {
	switch tag {
	case model.TypeTags_BOOLEAN:
		return booleanType
	case model.TypeTags_INT:
		return intType
	case model.TypeTags_NIL:
		return nilType
	case model.TypeTags_STRING:
		return stringType
	case model.TypeTags_FLOAT:
		return floatType
	default:
		panic("not implemented")
	}
}
