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
	AllOrNothingSubtypeAll                 = newAllOrNothingSubtypeFromBool(true)
	AllOrNothingSubtypeNothing             = newAllOrNothingSubtypeFromBool(false)
	_                          SubtypeData = &allOrNothingSubtype{}
	_                          Bdd         = &allOrNothingSubtype{}
)

func newAllOrNothingSubtypeFromBool(isAll bool) allOrNothingSubtype {
	this := allOrNothingSubtype{}
	this.isAll = isAll
	return this
}

func createAll() allOrNothingSubtype {
	// migrated from allOrNothingSubtype.java:38:5
	return AllOrNothingSubtypeAll
}

func createNothing() allOrNothingSubtype {
	// migrated from allOrNothingSubtype.java:42:5
	return AllOrNothingSubtypeNothing
}

func (this *allOrNothingSubtype) IsAllSubtype() bool {
	// migrated from allOrNothingSubtype.java:46:5
	return this.isAll
}

func (this *allOrNothingSubtype) IsNothingSubtype() bool {
	// migrated from allOrNothingSubtype.java:50:5
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
