/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package io

import (
	"ballerina-lang-go/runtime/api"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	// wasmOutputBuffer is used in WASM mode to capture output
	wasmOutputBuffer *bytes.Buffer
	wasmOutputMu     sync.RWMutex
	wasmOutputOnce   sync.Once
)

// Println prints the given values to the output writer, separated by spaces and followed by a newline.
func Println(values ...any) {
	writer := getOutputWriter()
	var output bytes.Buffer
	for i, v := range values {
		if i > 0 {
			output.WriteString(" ")
		}
		fmt.Fprintf(&output, "%v", v)
	}
	output.WriteString("\n")
	writer.Write(output.Bytes())
}

// GetWASMOutputBuffer returns the output buffer used in WASM mode.
// This allows tests to read captured output when running in WASM.
func GetWASMOutputBuffer() *bytes.Buffer {
	wasmOutputMu.RLock()
	defer wasmOutputMu.RUnlock()
	return wasmOutputBuffer
}

// isWASM returns true if running in WASM mode
func isWASM() bool {
	return os.Getenv("GOARCH") == "wasm"
}

// ensureWASMBuffer initializes the WASM output buffer if it doesn't exist
func ensureWASMBuffer() {
	wasmOutputOnce.Do(func() {
		wasmOutputBuffer = &bytes.Buffer{}
	})
}

// getOutputWriter returns the output writer to use (WASM buffer or os.Stdout)
func getOutputWriter() io.Writer {
	if isWASM() {
		ensureWASMBuffer()
		wasmOutputMu.RLock()
		buf := wasmOutputBuffer
		wasmOutputMu.RUnlock()
		return buf
	}
	return os.Stdout
}

func initIOModule(rt *api.Runtime) {
	api.RegisterExternFunction(rt.Registry, "ballerina", "io", "println", func(args []any) (any, error) {
		Println(args...)
		return nil, nil
	})
}

func init() {
	// In WASM mode, automatically set up output buffer for testing
	// This allows tests to capture output without explicit setup
	if isWASM() {
		ensureWASMBuffer()
	}
	api.RegisterModuleInitializer(initIOModule)
}
