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

import (
	"ballerina-lang-go/common"
	"testing"
)

// TestBddDiff tests BDD operations
// Ported from SemTypeBddTest.java:bddTest()
func TestBddDiff(t *testing.T) {
	// Create two BDD atoms from different rec atoms
	b1 := BddAtom(common.ToPointer(CreateRecAtom(1)))
	b2 := BddAtom(common.ToPointer(CreateRecAtom(2)))

	// Intersect them
	b1and2 := BddIntersect(b1, b2)

	// Calculate difference: (b1 ∩ b2) - b1
	r := BddDiff(b1and2, b1)

	// Type assert to BddAllOrNothing
	allOrNothing, ok := r.(BddAllOrNothing)
	if !ok {
		// Try pointer type
		allOrNothingPtr, ok := r.(*BddAllOrNothing)
		if !ok {
			t.Fatalf("expected BddAllOrNothing or *BddAllOrNothing, got %T", r)
		}
		allOrNothing = *allOrNothingPtr
	}

	// Assert that the result is not "all" (should be "nothing")
	if allOrNothing.IsAll() {
		t.Error("expected IsAll() to be false")
	}
}
