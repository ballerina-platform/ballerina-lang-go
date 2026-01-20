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

import "ballerina-lang-go/common"

type Error struct {
}

func ErrorDetail(detail SemType) SemType {
	// migrated from Error.java:39:5
	mappingSd := subtypeData(detail, BT_MAPPING)
	if allOrNothingSubtype, ok := mappingSd.(AllOrNothingSubtype); ok {
		if allOrNothingSubtype.IsAllSubtype() {
			return &ERROR
		} else {
			return &NEVER
		}
	}
	sd := BddIntersect(mappingSd.(Bdd), BDD_SUBTYPE_RO)
	if sd == BDD_SUBTYPE_RO {
		return &ERROR
	}
	return basicSubtype(BT_ERROR, sd.(ProperSubtypeData))
}

func ErrorDistinct(distinctId int) SemType {
	// migrated from Error.java:57:5
	common.Assert(distinctId >= 0)
	bdd := BddAtom(common.ToPointer(CreateDistinctRecAtom(((-distinctId) - 1))))
	return basicSubtype(BT_ERROR, bdd)
}
