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

import (
	"ballerina-lang-go/model"
)

// OCEDynamicEnvironmentData represents the dynamic environment data for object constructor expressions
type OCEDynamicEnvironmentData struct {
	// TypeInit                       *BLangTypeInit
	ObjectType *BObjectType
	// AttachedFunctionInvocation     *BLangAttachedFunctionInvocation
	LambdaFunctionsList            []BLangLambdaFunction
	CloneRef                       *OCEDynamicEnvironmentData
	OriginalClass                  *BLangClassDefinition
	Parents                        []*BLangClassDefinition
	DesugaredClosureVars           []*BLangSimpleVarRef
	InitInvocation                 model.ExpressionNode
	CloneAttempt                   int
	ClosureDesugaringInProgress    bool
	IsDirty                        bool
	BlockMapUpdatedInInitMethod    bool
	FunctionMapUpdatedInInitMethod bool
}

// NewOCEDynamicEnvironmentData creates a new OCEDynamicEnvironmentData instance
func NewOCEDynamicEnvironmentData() *OCEDynamicEnvironmentData {
	this := &OCEDynamicEnvironmentData{}
	this.LambdaFunctionsList = make([]BLangLambdaFunction, 0, 1)
	this.Parents = make([]*BLangClassDefinition, 0)
	this.DesugaredClosureVars = make([]*BLangSimpleVarRef, 0)
	this.CloneAttempt = 0
	return this
}

// CleanUp clears all collections and resets pointer fields
func (this *OCEDynamicEnvironmentData) CleanUp() {
	this.Parents = nil
	this.LambdaFunctionsList = nil
	this.DesugaredClosureVars = nil
}
