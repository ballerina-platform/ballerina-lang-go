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

func ErrorDetailAtomicType(ctx Context, errorType SemType) (MappingAtomicType, bool) {
	errorType = Intersect(errorType, ERROR)
	if IsNever(errorType) || !IsSubtype(ctx, errorType, ERROR) {
		return MappingAtomicType{}, false
	}

	if IsSameType(ctx, errorType, ERROR) {
		return mappingAtomicTypeFrom(nil, nil, cellContaining(ctx.Env(), CreateCloneable(ctx))), true
	}
	mappingSd := subtypeData(errorType, BTError)
	if bddNode, ok := mappingSd.(BddNode); ok {
		if bddNode.Atom().index() != 0 {
			// Not readonly. Not sure if this can happen (due to ErroWithDetail) but just in case
			return MappingAtomicType{}, false
		}
		if !isNothing(bddNode.Middle()) || !isNothing(bddNode.Right()) {
			// Not atomic
			return MappingAtomicType{}, false
		}
		if leftNode, ok := bddNode.Left().(BddNode); ok {
			if !isSimpleNode(leftNode.Left(), leftNode.Middle(), leftNode.Right()) {
				// Also not atomic
				return MappingAtomicType{}, false
			}
			return *ctx.MappingAtomType(leftNode.Atom()), true
		} else {
			return MappingAtomicType{}, false
		}
	}
	return MappingAtomicType{}, false
}

func ErrorWithDetail(detail SemType) SemType {
	mappingSd := subtypeData(detail, BTMapping)
	if allOrNothingSubtype, ok := mappingSd.(allOrNothingSubtype); ok {
		if allOrNothingSubtype.IsAllSubtype() {
			return ERROR
		} else {
			return NEVER
		}
	}
	sd := bddIntersect(mappingSd.(Bdd), BDD_SUBTYPE_RO)
	if sd == BDD_SUBTYPE_RO {
		return ERROR
	}
	return getBasicSubtype(BTError, sd.(ProperSubtypeData))
}

func errorDistinct(distinctId int) SemType {
	common.Assert(distinctId >= 0)
	bdd := bddAtom(new(createDistinctRecAtom(((-distinctId) - 1))))
	return getBasicSubtype(BTError, bdd)
}
