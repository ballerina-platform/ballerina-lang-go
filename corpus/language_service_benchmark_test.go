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
	"os"
	"testing"

	"ballerina-lang-go/context"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testphases"
)

func BenchmarkLanguageService(b *testing.B) {
	for _, tc := range benchTestPairs(b) {
		contentBytes, err := os.ReadFile(tc.InputPath)
		if err != nil {
			b.Fatalf("read %s: %v", tc.InputPath, err)
		}
		content := string(contentBytes)

		b.Run(tc.Name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
				cx := context.NewCompilerContext(env)
				_, err := testphases.RunPipelineWithContent(env, cx, nil, testphases.PhaseDesugar, tc.InputPath, content)
				if err != nil {
					b.Fatalf("language service benchmark failed for %s: %v", tc.InputPath, err)
				}
				if cx.HasDiagnostics() {
					b.Fatalf("language service benchmark produced diagnostics for %s", tc.InputPath)
				}
			}
		})
	}
}
