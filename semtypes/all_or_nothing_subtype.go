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

type allOrNothingSubtype struct {
	isAll bool
}

var (
	allOrNothingSubtypeAll                 = newAllOrNothingSubtypeFromBool(true)
	allOrNothingSubtypeNothing             = newAllOrNothingSubtypeFromBool(false)
	_                          SubtypeData = &allOrNothingSubtype{}
	_                          Bdd         = &allOrNothingSubtype{}
)

func newAllOrNothingSubtypeFromBool(isAll bool) allOrNothingSubtype {
	this := allOrNothingSubtype{}
	this.isAll = isAll
	return this
}

func createAll() allOrNothingSubtype {
	return allOrNothingSubtypeAll
}

func createNothing() allOrNothingSubtype {
	return allOrNothingSubtypeNothing
}

func (this *allOrNothingSubtype) IsAllSubtype() bool {
	return this.isAll
}

func (this *allOrNothingSubtype) IsNothingSubtype() bool {
	return (!this.isAll)
}

func (this *allOrNothingSubtype) canonicalKey() string {
	if this.isAll {
		return "true"
	} else {
		return "false"
	}
}

func (this *allOrNothingSubtype) String() string {
	if this.isAll {
		return "all"
	} else {
		return "nothing"
	}
}
