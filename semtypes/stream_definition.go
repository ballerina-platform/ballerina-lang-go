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

// Represent stream type desc.
//
// @since 2201.12.0
type streamDefinition struct {
	listDefinition ListDefinition
}

var _ Definition = &streamDefinition{}

func newStreamDefinition() streamDefinition {
	this := streamDefinition{}
	this.listDefinition = NewListDefinition()
	return this
}

func streamContaining(tupleType SemType) SemType {
	bdd := subtypeData(tupleType, BTList)
	return createBasicSemType(BTStream, bdd)
}

func (s *streamDefinition) GetSemType(env Env) SemType {
	return streamContaining(s.listDefinition.GetSemType(env))
}

func (s *streamDefinition) Define(env Env, valueTy SemType, completionTy SemType) SemType {
	if common.PointerEqualToValue(VAL, completionTy) && common.PointerEqualToValue(VAL, valueTy) {
		return STREAM
	}
	tuple := s.listDefinition.TupleTypeWrapped(env, valueTy, completionTy)
	return streamContaining(tuple)
}
