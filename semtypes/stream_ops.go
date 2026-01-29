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

type StreamOps struct {
	CommonOps
}

var _ BasicTypeOps = &StreamOps{}

func streamSubtypeComplement(t SubtypeData) SubtypeData {
	// migrated from StreamOps.java:38:5
	return bddSubtypeDiff(LIST_SUBTYPE_TWO_ELEMENT, t)
}

func streamSubtypeIsEmpty(cx Context, t SubtypeData) bool {
	// migrated from StreamOps.java:42:5
	b := Bdd(t)
	if bddPosMaybeEmpty(b) {
		b = BddIntersect(b, LIST_SUBTYPE_TWO_ELEMENT)
	}
	return listSubtypeIsEmpty(cx, b)
}

func NewStreamOps() StreamOps {
	this := StreamOps{}
	return this
}

func (this *StreamOps) Complement(t SubtypeData) SubtypeData {
	// migrated from StreamOps.java:51:5
	return streamSubtypeComplement(t)
}

func (this *StreamOps) IsEmpty(cx Context, t SubtypeData) bool {
	// migrated from StreamOps.java:56:5
	return streamSubtypeIsEmpty(cx, t)
}
