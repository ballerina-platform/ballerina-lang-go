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

package main

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	baseRef    = "HEAD~1"
	headRef    = "HEAD"
	targetPath = "../../corpus/bal/subset6/06-bench/4-v.bal"
)

func TestBenchmarkRunExportsHTML(t *testing.T) {
	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("benchmark target is unavailable: %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "output.html")
	b := &benchmark{
		config: config{
			baseRef:    baseRef,
			headRef:    headRef,
			target:     targetPath,
			warmup:     2,
			runs:       10,
			exportPath: outputPath,
		},
	}
	if err := b.run(); err != nil {
		t.Fatalf("benchmark run failed: %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("expected html report at %q: %v", outputPath, err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected non-empty html report at %q", outputPath)
	}
}
