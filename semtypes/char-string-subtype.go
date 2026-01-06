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

type CharStringSubtype struct {
	allowed bool
	values  []EnumerableType[string]
}

// FIXME: when fixing this we need to think about Values
var _ EnumerableSubtype[string] = &CharStringSubtype{}

func CharStringSubtypeFrom(allowed bool, values []EnumerableType[string]) CharStringSubtype {
	// migrated from CharStringSubtype.java:39:5
	return CharStringSubtype{
		allowed: allowed,
		values:  values,
	}
}

func (this *CharStringSubtype) Allowed() bool {
	// migrated from CharStringSubtype.java:43:5
	return this.allowed
}

func (this *CharStringSubtype) Values() []EnumerableType[string] {
	// migrated from CharStringSubtype.java:48:5
	return this.values
}
