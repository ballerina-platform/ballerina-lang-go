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

func CellContaining(env Env, ty SemType) CellSemType {
	// migrated from CellSubtype.java:42:5
	return CellContainingWithEnvSemTypeCellMutability(env, ty, CellMutability_CELL_MUT_LIMITED)
}

func RoCellContaining(env Env, ty SemType) CellSemType {
	// migrated from CellSubtype.java:46:5
	return CellContainingWithEnvSemTypeCellMutability(env, ty, CellMutability_CELL_MUT_NONE)
}

func CellContainingWithEnvSemTypeCellMutability(env Env, ty SemType, mut CellMutability) CellSemType {
	// migrated from CellSubtype.java:50:5
	common.Assert(IsNever(ty) || !IsSubtypeSimple(ty, CELL))
	atomicCell := CellAtomicTypeFrom(ty, mut)
	atom := env.cellAtom(&atomicCell)
	bdd := BddAtom(&atom)
	complexSemType := basicSubtype(BT_CELL, bdd)
	return CellSemTypeFrom(complexSemType.SubtypeDataList())
}
