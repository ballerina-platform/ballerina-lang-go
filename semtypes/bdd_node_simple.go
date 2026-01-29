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

type BddNodeSimple struct {
	atom Atom
}

var _ BddNode = &BddNodeSimple{}

func (this *BddNodeSimple) Left() Bdd {
	// migrated from BddNodeSimple.java:32:5
	return BddAll()
}

func (this *BddNodeSimple) Middle() Bdd {
	// migrated from BddNodeSimple.java:37:5
	return BddNothing()
}

func (this *BddNodeSimple) Right() Bdd {
	// migrated from BddNodeSimple.java:42:5
	return BddNothing()
}

func (this *BddNodeSimple) Atom() Atom {
	return this.atom
}
