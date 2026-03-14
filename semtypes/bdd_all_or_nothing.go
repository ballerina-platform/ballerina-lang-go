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

type BddAllOrNothing struct {
	isAll bool
}

var (
	all         = newBddAllOrNothingFromBool(true)
	nothing     = newBddAllOrNothingFromBool(false)
	_       Bdd = &BddAllOrNothing{}
)

func newBddAllOrNothingFromBool(isAll bool) BddAllOrNothing {
	this := BddAllOrNothing{}
	this.isAll = isAll
	return this
}

func BddAll() *BddAllOrNothing {
	// migrated from BddAllOrNothing.java:37:5
	return &all
}

func BddNothing() *BddAllOrNothing {
	// migrated from BddAllOrNothing.java:41:5
	return &nothing
}

func (b *BddAllOrNothing) IsAll() bool {
	// migrated from BddAllOrNothing.java:45:5
	return b.isAll
}

func (b *BddAllOrNothing) IsNothing() bool {
	// migrated from BddAllOrNothing.java:49:5
	return (!b.isAll)
}

func (b *BddAllOrNothing) Complement() *BddAllOrNothing {
	// migrated from BddAllOrNothing.java:53:5
	if b.isAll {
		return &nothing
	}
	return &all
}

func (b *BddAllOrNothing) canonicalKey() string {
	if b.isAll {
		return "true"
	} else {
		return "false"
	}
}
