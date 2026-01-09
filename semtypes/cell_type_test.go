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
// software distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package semtypes

import (
	"testing"
)

// TestTypeCellDisparity tests type and cell type disparity
// Ported from CellTypeTest.java:testTypeCellDisparity()
func TestTypeCellDisparity(t *testing.T) {
	env := GetTypeEnv()
	ctx := ContextFrom(env)

	tests := []struct {
		name     string
		t1       SemType
		t2       SemType
		relation Relation
	}{
		{
			name:     "INT vs cell(INT, NONE)",
			t1:       &INT,
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_NONE),
			relation: RelationNoRelation,
		},
		{
			name:     "INT vs cell(INT, LIMITED)",
			t1:       &INT,
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_LIMITED),
			relation: RelationNoRelation,
		},
		{
			name:     "INT vs cell(INT, UNLIMITED)",
			t1:       &INT,
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationNoRelation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSemTypeRelation(t, ctx, tt.t1, tt.t2, tt.relation)
		})
	}
}

// TestBasicCellSubtyping tests basic cell subtyping
// Ported from CellTypeTest.java:testBasicCellSubtyping()
func TestBasicCellSubtyping(t *testing.T) {
	env := GetTypeEnv()
	ctx := ContextFrom(env)

	tests := []struct {
		name      string
		t1        SemType
		t2        SemType
		relations [3]Relation // [NONE, LIMITED, UNLIMITED]
	}{
		{
			name:      "INT vs INT",
			t1:        &INT,
			t2:        &INT,
			relations: [3]Relation{RelationEqual, RelationEqual, RelationEqual},
		},
		{
			name:      "BOOLEAN vs BOOLEAN",
			t1:        &BOOLEAN,
			t2:        &BOOLEAN,
			relations: [3]Relation{RelationEqual, RelationEqual, RelationEqual},
		},
		{
			name:      "BYTE vs INT",
			t1:        BYTE,
			t2:        &INT,
			relations: [3]Relation{RelationSubtype, RelationSubtype, RelationSubtype},
		},
		{
			name:      "BOOLEAN vs INT",
			t1:        &BOOLEAN,
			t2:        &INT,
			relations: [3]Relation{RelationNoRelation, RelationNoRelation, RelationNoRelation},
		},
		{
			name:      "BOOLEAN vs INT|BOOLEAN",
			t1:        &BOOLEAN,
			t2:        Union(&INT, &BOOLEAN),
			relations: [3]Relation{RelationSubtype, RelationSubtype, RelationSubtype},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutabilities := []CellMutability{
				CellMutability_CELL_MUT_NONE,
				CellMutability_CELL_MUT_LIMITED,
				CellMutability_CELL_MUT_UNLIMITED,
			}

			for i, mut := range mutabilities {
				c1 := testCell(env, tt.t1, mut)
				c2 := testCell(env, tt.t2, mut)
				actual := getSemTypeRelation(ctx, c1, c2)
				if actual != tt.relations[i] {
					t.Errorf("mutability %v: got %v, want %v", mut, actual, tt.relations[i])
				}
			}
		})
	}
}

// TestCellSubtyping1 tests cell subtyping with unions
// Ported from CellTypeTest.java:testCellSubtyping1()
func TestCellSubtyping1(t *testing.T) {
	env := GetTypeEnv()
	ctx := ContextFrom(env)

	tests := []struct {
		name     string
		t1       SemType
		t2       SemType
		relation Relation
	}{
		// Set 1
		{
			name:     "cell(INT,NONE)|cell(BOOLEAN,NONE) vs cell(INT|BOOLEAN,NONE)",
			t1:       Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_NONE)),
			t2:       testCell(env, Union(&INT, &BOOLEAN), CellMutability_CELL_MUT_NONE),
			relation: RelationEqual,
		},
		{
			name:     "cell(INT,LIMITED)|cell(BOOLEAN,LIMITED) vs cell(INT|BOOLEAN,LIMITED)",
			t1:       Union(testCell(env, &INT, CellMutability_CELL_MUT_LIMITED), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, Union(&INT, &BOOLEAN), CellMutability_CELL_MUT_LIMITED),
			relation: RelationSubtype,
		},
		{
			name:     "cell(INT,UNLIMITED)|cell(BOOLEAN,UNLIMITED) vs cell(INT|BOOLEAN,UNLIMITED)",
			t1:       Union(testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, Union(&INT, &BOOLEAN), CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationEqual,
		},
		// Set 2
		{
			name:     "cell(INT,NONE)|cell(BOOLEAN,NONE)|cell(STRING,NONE) vs cell(INT|BOOLEAN|STRING,NONE)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_NONE)), testCell(env, &STRING, CellMutability_CELL_MUT_NONE)),
			t2:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_NONE),
			relation: RelationEqual,
		},
		{
			name:     "cell(INT,LIMITED)|cell(BOOLEAN,LIMITED)|cell(STRING,LIMITED) vs cell(INT|BOOLEAN|STRING,LIMITED)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_LIMITED), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_LIMITED)), testCell(env, &STRING, CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_LIMITED),
			relation: RelationSubtype,
		},
		{
			name:     "cell(INT,UNLIMITED)|cell(BOOLEAN,UNLIMITED)|cell(STRING,UNLIMITED) vs cell(INT|BOOLEAN|STRING,UNLIMITED)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_UNLIMITED)), testCell(env, &STRING, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationEqual,
		},
		// Set 3
		{
			name:     "cell(roTuple(INT),NONE)|cell(roTuple(BOOLEAN),NONE) vs cell(roTuple(INT|BOOLEAN),NONE)",
			t1:       Union(testCell(env, testRoTuple(env, &INT), CellMutability_CELL_MUT_NONE), testCell(env, testRoTuple(env, &BOOLEAN), CellMutability_CELL_MUT_NONE)),
			t2:       testCell(env, testRoTuple(env, Union(&INT, &BOOLEAN)), CellMutability_CELL_MUT_NONE),
			relation: RelationEqual,
		},
		{
			name:     "cell(tuple(INT),LIMITED)|cell(tuple(BOOLEAN),LIMITED) vs cell(tuple(INT|BOOLEAN),LIMITED)",
			t1:       Union(testCell(env, testTuple(env, &INT), CellMutability_CELL_MUT_LIMITED), testCell(env, testTuple(env, &BOOLEAN), CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, testTuple(env, Union(&INT, &BOOLEAN)), CellMutability_CELL_MUT_LIMITED),
			relation: RelationSubtype,
		},
		{
			name:     "cell(tuple(INT),UNLIMITED)|cell(tuple(BOOLEAN),UNLIMITED) vs cell(tuple(INT|BOOLEAN),UNLIMITED)",
			t1:       Union(testCell(env, testTuple(env, &INT), CellMutability_CELL_MUT_UNLIMITED), testCell(env, testTuple(env, &BOOLEAN), CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, testTuple(env, Union(&INT, &BOOLEAN)), CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationSubtype,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSemTypeRelation(t, ctx, tt.t1, tt.t2, tt.relation)
		})
	}
}

// TestCellSubtyping2 tests cell subtyping with different mutability
// Ported from CellTypeTest.java:testCellSubtyping2()
func TestCellSubtyping2(t *testing.T) {
	env := GetTypeEnv()
	ctx := ContextFrom(env)

	tests := []struct {
		name     string
		t1       SemType
		t2       SemType
		relation Relation
	}{
		// test 1
		{
			name:     "cell(INT,NONE)|cell(BOOLEAN,UNLIMITED)|cell(STRING,LIMITED) vs cell(INT|BOOLEAN|STRING,UNLIMITED)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_UNLIMITED)), testCell(env, &STRING, CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationSubtype,
		},
		// test 2
		{
			name:     "cell(INT|BOOLEAN|STRING,NONE) vs cell(INT,NONE)|cell(BOOLEAN,UNLIMITED)|cell(STRING,LIMITED)",
			t1:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_NONE),
			t2:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_UNLIMITED)), testCell(env, &STRING, CellMutability_CELL_MUT_LIMITED)),
			relation: RelationSubtype,
		},
		// test 3
		{
			name:     "cell(INT,NONE)|cell(BOOLEAN,UNLIMITED)|cell(STRING,LIMITED) vs cell(INT|BOOLEAN|STRING,LIMITED)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_UNLIMITED)), testCell(env, &STRING, CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, Union(Union(&INT, &BOOLEAN), &STRING), CellMutability_CELL_MUT_LIMITED),
			relation: RelationNoRelation,
		},
		// test 4
		{
			name:     "cell(INT,NONE)|cell(INT,LIMITED)|cell(INT,UNLIMITED) vs cell(INT,UNLIMITED)",
			t1:       Union(Union(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &INT, CellMutability_CELL_MUT_LIMITED)), testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationEqual,
		},
		// test 5
		{
			name:     "cell(INT,NONE)∩cell(INT,LIMITED)∩cell(INT,UNLIMITED) vs cell(INT,UNLIMITED)",
			t1:       Intersect(Intersect(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &INT, CellMutability_CELL_MUT_LIMITED)), testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED),
			relation: RelationSubtype,
		},
		// test 6
		{
			name:     "cell(INT,NONE)∩cell(INT,LIMITED)∩cell(INT,UNLIMITED) vs cell(INT,LIMITED)",
			t1:       Intersect(Intersect(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &INT, CellMutability_CELL_MUT_LIMITED)), testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_LIMITED),
			relation: RelationSubtype,
		},
		// test 7
		{
			name:     "cell(INT,NONE)∩cell(INT,LIMITED)∩cell(INT,UNLIMITED) vs cell(INT,NONE)",
			t1:       Intersect(Intersect(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &INT, CellMutability_CELL_MUT_LIMITED)), testCell(env, &INT, CellMutability_CELL_MUT_UNLIMITED)),
			t2:       testCell(env, &INT, CellMutability_CELL_MUT_NONE),
			relation: RelationEqual,
		},
		// test 8
		{
			name:     "cell(INT,NONE)∩cell(INT,LIMITED)∩cell(BYTE,LIMITED) vs cell(BYTE,LIMITED)",
			t1:       Intersect(Intersect(testCell(env, &INT, CellMutability_CELL_MUT_NONE), testCell(env, &INT, CellMutability_CELL_MUT_LIMITED)), testCell(env, BYTE, CellMutability_CELL_MUT_LIMITED)),
			t2:       testCell(env, BYTE, CellMutability_CELL_MUT_LIMITED),
			relation: RelationSubtype,
		},
		// test 9
		{
			name:     "cell(INT,NONE)∩(cell(BYTE,LIMITED)|cell(BOOLEAN,LIMITED)) vs cell(BYTE,NONE)",
			t1:       Intersect(testCell(env, &INT, CellMutability_CELL_MUT_NONE), Union(testCell(env, BYTE, CellMutability_CELL_MUT_LIMITED), testCell(env, &BOOLEAN, CellMutability_CELL_MUT_LIMITED))),
			t2:       testCell(env, BYTE, CellMutability_CELL_MUT_NONE),
			relation: RelationEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSemTypeRelation(t, ctx, tt.t1, tt.t2, tt.relation)
		})
	}
}
