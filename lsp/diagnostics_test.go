// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lsp

import "testing"

func TestLSPPositionMapperPositionsMatchesSinglePosition(t *testing.T) {
	content := "a😀b\r\nc\n𝄞d\re"
	mapper := newLSPPositionMapper(content)
	offsets := []int{len(content), 0, -1, len(content) + 1}
	for offset := len(content) - 1; offset > 0; offset-- {
		offsets = append(offsets, offset)
	}

	positions := mapper.Positions(offsets)
	for i, offset := range offsets {
		expected := mapper.Position(offset)
		if positions[i] != expected {
			t.Fatalf("position for offset %d = %#v, want %#v", offset, positions[i], expected)
		}
	}
}
