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
// migrated from StreamDefinition.java:37:1
type StreamDefinition struct {
	// migrated from StreamDefinition.java:39:5
	listDefinition ListDefinition
}

// migrated from StreamDefinition.java:37:1
var _ Definition = &StreamDefinition{}

func NewStreamDefinition() StreamDefinition {
	this := StreamDefinition{}
	this.listDefinition = NewListDefinition()
	return this
}

// migrated from StreamDefinition.java:54:5
func streamContaining(tupleType SemType) SemType {
	// migrated from StreamDefinition.java:55:9
	bdd := subtypeData(tupleType, BT_LIST)
	// migrated from StreamDefinition.java:56:9
	return CreateBasicSemType(BT_STREAM, bdd)
}

// migrated from StreamDefinition.java:42:5
func (this *StreamDefinition) GetSemType(env Env) SemType {
	// migrated from StreamDefinition.java:43:9
	return streamContaining(this.listDefinition.GetSemType(env))
}

// migrated from StreamDefinition.java:46:5
func (this *StreamDefinition) Define(env Env, valueTy SemType, completionTy SemType) SemType {
	// migrated from StreamDefinition.java:47:9
	if common.PointerEqualToValue(VAL, completionTy) && common.PointerEqualToValue(VAL, valueTy) {
		return &STREAM
	}
	// migrated from StreamDefinition.java:50:9
	tuple := this.listDefinition.TupleTypeWrapped(env, valueTy, completionTy)
	return streamContaining(tuple)
}
