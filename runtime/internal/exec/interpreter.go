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

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/runtime/internal/modules"
)

func Interpret(pkg bir.BIRPackage, env *extern.Env) (err error) {
	ctx := extern.CreateContext(env)
	cs := &callStack{elements: make([]callStackEntry, 0, 32)}
	ctx.CallStack = cs
	birModule, err := modules.NewBIRModule(ctx, &pkg)
	if err != nil {
		return err
	}
	env.Registry.(*modules.Registry).RegisterModule(pkg.PackageID, birModule)
	defer func() {
		if r := recover(); r != nil {
			ctx.ReleaseAllHeldLocks()
			err = getFormattedError(cs, r)
		}
	}()
	if pkg.InitFunction != nil {
		defer func() {
			if r := recover(); r != nil {
				ctx.ReleaseAllHeldLocks()
				err = getFormattedError(cs, r)
			}
		}()
		if result := executeFunction(ctx, *pkg.InitFunction, nil, nil); result != nil {
			panic(result)
		}
	}
	if pkg.MainFunction != nil {
		defer func() {
			if r := recover(); r != nil {
				ctx.ReleaseAllHeldLocks()
				err = getFormattedError(cs, r)
			}
		}()
		if result := executeFunction(ctx, *pkg.MainFunction, nil, nil); result != nil {
			panic(result)
		}
	}
	return err
}
