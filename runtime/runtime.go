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

package runtime

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/api"
	"ballerina-lang-go/runtime/internal/exec"
	_ "ballerina-lang-go/stdlibs/io" // Import to trigger io module's init() registration
	"fmt"
	"os"
)

// Interpret interprets a BIR package and returns the runtime instance and any error.
// When the runtime is not needed, the first return value can be ignored.
func Interpret(pkg bir.BIRPackage) (*api.Runtime, error) {
	rt := api.NewRuntime()
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	exec.Interpret(pkg, rt)
	return rt, err
}
