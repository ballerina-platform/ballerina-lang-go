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

package native

import (
	"io"
	"strings"
	"testing"
)

func TestNewResponseBodyHolder(t *testing.T) {
	// A nil stream yields an already-materialized empty body.
	empty := newResponseBodyHolder(nil)
	buf, err := empty.materialize()
	if err != nil {
		t.Fatalf("empty holder materialize: %v", err)
	}
	if len(buf) != 0 {
		t.Errorf("empty holder body = %v, want empty", buf)
	}

	// A streaming holder materializes the stream contents once.
	h := newResponseBodyHolder(io.NopCloser(strings.NewReader("payload")))
	got, err := h.materialize()
	if err != nil {
		t.Fatalf("streaming holder materialize: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("streaming holder body = %q, want %q", got, "payload")
	}
	// A second call returns the cached buffer (sync.Once guard).
	again, _ := h.materialize()
	if string(again) != "payload" {
		t.Errorf("cached body = %q, want %q", again, "payload")
	}
}
