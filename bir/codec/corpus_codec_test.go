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

package codec

import (
	"testing"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/test_util/testphases"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// getBIRDiff generates a detailed diff string showing differences between expected and actual BIR text.
func getBIRDiff(expectedText, actualText string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, false)
	return dmp.DiffPrettyText(diffs)
}

// TestBIRSerialization tests BIR serialization and deserialization roundtrip from .bal source files in the corpus.
func TestBIRSerialization(t *testing.T) {
	testPairs := test_util.GetValidTests(t, test_util.BIR)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testBIRSerialization(t, testPair)
		})
	}
}

func testBIRSerialization(t *testing.T, testPair test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic during BIR serialization roundtrip for %s: %v", testPair.InputPath, r)
		}
	}()

	initialEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv())
	initialContext := context.NewCompilerContext(initialEnv)
	result, err := testphases.RunPipeline(initialContext, testphases.PhaseBIR, testPair.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testPair.InputPath, err)
		return
	}

	if result.BIRPackage == nil {
		t.Errorf("BIR package is nil for %s", testPair.InputPath)
		return
	}

	prettyPrinter := bir.PrettyPrinter{}
	expectedBIR := prettyPrinter.Print(*result.BIRPackage)

	serializedBIR, err := Marshal(result.BIRPackage)
	if err != nil {
		t.Errorf("error serializing BIR package for %s: %v", testPair.InputPath, err)
		return
	}
	// Make sure we are using a different type env.
	initialEnv = nil     //nolint:ineffassign
	initialContext = nil //nolint:ineffassign

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv())
	tyCtx := context.NewCompilerContext(env)
	deserializedBIRPkg, err := Unmarshal(tyCtx, serializedBIR)
	if err != nil {
		t.Errorf("error deserializing BIR package for %s: %v", testPair.InputPath, err)
		return
	}

	actualBIR := prettyPrinter.Print(*deserializedBIRPkg)

	if expectedBIR != actualBIR {
		diff := getBIRDiff(expectedBIR, actualBIR)
		t.Errorf("BIR roundtrip mismatch for %s\n%s", testPair.InputPath, diff)
	}
}
