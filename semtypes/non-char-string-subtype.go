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

type NonCharStringSubtype struct {
	allowed bool
	values  []EnumerableType[string]
}

var _ EnumerableSubtype[string] = &NonCharStringSubtype{}

func newNonCharStringSubtypeFromBoolEnumerableStrings(allowed bool, values []EnumerableType[string]) NonCharStringSubtype {
	this := NonCharStringSubtype{}
	this.allowed = allowed
	this.values = values
	return this
}

func NonCharStringSubtypeFrom(allowed bool, values []EnumerableType[string]) NonCharStringSubtype {
	// migrated from NonCharStringSubtype.java:39:5
	return newNonCharStringSubtypeFromBoolEnumerableStrings(allowed, values)
}

func (this *NonCharStringSubtype) Allowed() bool {
	// migrated from NonCharStringSubtype.java:43:5
	return this.allowed
}

func (this *NonCharStringSubtype) Values() []EnumerableType[string] {
	// migrated from NonCharStringSubtype.java:48:5
	return this.values
}
