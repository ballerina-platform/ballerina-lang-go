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

package corpus

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"ballerina-lang-go/bir"
	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/context"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
)

// TestBIRSerializationRoundtrip compiles .bal files to BIR, serializes the BIR, deserializes it
// with a fresh compiler context, executes the deserialized BIR, and validates the output matches
// the expected integration test output.
// birSerializationRoundtripSkipPrefixes is the BIR-roundtrip *additional* skip
// list, on top of the shared test_util.UnsupportedTests baseline (which is
// applied via isTestSkipped below). It skips tests whose BIR cannot be
// round-tripped in isolation because they depend on stdlib modules whose
// Ballerina-level functions (e.g. http:Client.init) live in the http BIR
// rather than being registered as native Go functions. Running only the
// deserialized consumer BIR in a fresh env would panic on those calls.
var birSerializationRoundtripSkipPrefixes = []string{
	"08-network/http-client",
}

func TestBIRSerializationRoundtrip(t *testing.T) {
	testPairs := test_util.GetTests(t, test_util.Integration, func(path string) bool {
		if test_util.IsFutureTest(path) {
			return false
		}
		return strings.HasSuffix(path, "-v.bal") || strings.HasSuffix(path, "-p.bal")
	})

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testBIRSerializationRoundtrip(t, testPair)
		})
	}
}

func testBIRSerializationRoundtrip(t *testing.T, testPair test_util.TestCase) {
	if isTestSkipped(testPair) || test_util.MatchesSkipSubstr(testPair.InputPath, birSerializationRoundtripSkipPrefixes) {
		t.Skipf("Skipping BIR serialization roundtrip test for %s", testPair.InputPath)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic during BIR serialization roundtrip for %s: %v", testPair.InputPath, r)
		}
	}()

	// Step 1: Compile .bal to BIR.
	var stdoutBuf, stderrBuf bytes.Buffer
	birPkgs, tyEnv, _, compileErr := runCompilePhase(testPair.InputPath, &stdoutBuf, &stderrBuf)
	if len(birPkgs) == 0 || compileErr != nil {
		t.Fatalf("compilation failed for %s: %v", testPair.InputPath, compileErr)
	}
	// The consumer's BIR is the last package in topological order.
	consumerPkg := birPkgs[len(birPkgs)-1]

	// Step 2: Serialize BIR.
	serialized, err := bircodec.Marshal(tyEnv, consumerPkg)
	if err != nil {
		t.Fatalf("BIR serialization failed for %s: %v", testPair.InputPath, err)
	}

	// Step 3: Deserialize with a fresh compiler context.
	freshEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	freshCtx := context.NewCompilerContext(freshEnv)
	deserialized, err := bircodec.Unmarshal(freshCtx, serialized)
	if err != nil {
		t.Fatalf("BIR deserialization failed for %s: %v", testPair.InputPath, err)
	}

	// Step 4: Execute the deserialized BIR.
	var rtStdoutBuf, rtStderrBuf bytes.Buffer
	runInterpretPhase([]*bir.BIRPackage{deserialized}, freshEnv.GetTypeEnv(), &rtStdoutBuf, &rtStderrBuf)

	// Step 5: Compare against expected output.
	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(testPair.ExpectedPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", testPair.ExpectedPath, err)
	}

	result := evaluateTestResult(expectedStdout, expectedStderr, rtStdoutBuf.String(), rtStderrBuf.String())
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != normalizeIntegrationStderr(result.actualStderr)

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(
			normalizeIntegrationStderr(result.expectedStderr),
			normalizeIntegrationStderr(result.actualStderr),
		))
	}
	t.Errorf("%s", msg.String())
}
