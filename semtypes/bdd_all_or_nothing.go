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

type bddAllOrNothing struct {
	isAll bool
}

var (
	all         = newBddAllOrNothingFromBool(true)
	nothing     = newBddAllOrNothingFromBool(false)
	_       Bdd = &bddAllOrNothing{}
)

func newBddAllOrNothingFromBool(isAll bool) bddAllOrNothing {
	this := bddAllOrNothing{}
	this.isAll = isAll
	return this
}

func bddAll() *bddAllOrNothing {
	// migrated from bddAllOrNothing.java:37:5
	return &all
}

func bddNothing() *bddAllOrNothing {
	// migrated from bddAllOrNothing.java:41:5
	return &nothing
}

func (this *bddAllOrNothing) IsAll() bool {
	// migrated from bddAllOrNothing.java:45:5
	return this.isAll
}

func (this *bddAllOrNothing) IsNothing() bool {
	// migrated from bddAllOrNothing.java:49:5
	return (!this.isAll)
}

func (this *bddAllOrNothing) complement() *bddAllOrNothing {
	// migrated from bddAllOrNothing.java:53:5
	if this.isAll {
		return &nothing
	}
	return &all
}

func (this *bddAllOrNothing) canonicalKey() string {
	if this.isAll {
		return "true"
	} else {
		return "false"
	}
}

func (this *bddAllOrNothing) String() string {
	if this.isAll {
		return "all"
	} else {
		return "nothing"
	}
}
