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

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var lspLogMu sync.Mutex

func logLS(root string, format string, args ...any) {
	if !lspLoggingEnabled() || root == "" {
		return
	}
	if filepath.Ext(root) != "" {
		root = filepath.Dir(root)
	}
	logDir := filepath.Join(root, ".bal")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return
	}
	file, err := os.OpenFile(filepath.Join(logDir, "lsp.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	line := fmt.Sprintf(format, args...)
	lspLogMu.Lock()
	defer lspLogMu.Unlock()
	_, _ = fmt.Fprintf(file, "%s %s\n", time.Now().Format(time.RFC3339Nano), line)
}

func lspLoggingEnabled() bool {
	value := os.Getenv("BAL_LSP_LOG")
	return value == "1" || value == "true" || value == "TRUE"
}

func projectKindString(kind ProjectKind) string {
	switch kind {
	case ProjectKindBuild:
		return "build"
	case ProjectKindSingleFile:
		return "single-file"
	default:
		return "unknown"
	}
}
