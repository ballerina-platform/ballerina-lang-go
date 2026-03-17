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
	"strings"
	"testing"
)

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
