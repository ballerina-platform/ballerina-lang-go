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

package io

import (
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

const (
	orgName    = "ballerina"
	moduleName = "io"
	funcName   = "println"
)

func Println(vals ...values.BalValue) {
	parts := make([]string, len(vals))
	visited := make(map[uintptr]bool)
	for i, v := range vals {
		parts[i] = formatValue(v, visited)
	}
	fmt.Fprintln(os.Stdout, strings.Join(parts, " "))
}

func formatValue(v values.BalValue, visited map[uintptr]bool) string {
	if v == nil {
		return "nil"
	}
	if list, ok := v.(*values.List); ok {
		ptr := uintptr(unsafe.Pointer(list))
		if visited[ptr] {
			return "[...]"
		}
		visited[ptr] = true
		var b strings.Builder
		b.WriteByte('[')
		for i := 0; i < list.Len(); i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(formatValue(list.Get(i), visited))
		}
		b.WriteByte(']')
		return b.String()
	}
	return fmt.Sprintf("%v", v)
}

func printlnExtern(args []values.BalValue) (values.BalValue, error) {
	Println(args...)
	return nil, nil
}

func initIOModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, funcName, printlnExtern)
}

func init() {
	runtime.RegisterModuleInitializer(initIOModule)
}
