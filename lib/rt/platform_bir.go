// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
//
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

package rt

import (
	"sync"

	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lib/registry"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
)

var platformBIRUnmarshalCtx = sync.OnceValue(func() *context.CompilerContext {
	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	return context.NewCompilerContext(env)
})

func init() {
	runtime.RegisterModuleInitializer(loadEmbeddedPlatformModules)
}

func loadEmbeddedPlatformModules(rt *runtime.Runtime) {
	cc := platformBIRUnmarshalCtx()
	registry.ForEachEmbeddedPlatformBIR(func(birBytes []byte) {
		pkg, err := bircodec.Unmarshal(cc, birBytes)
		if err != nil {
			panic("rt: embedded platform bir: " + err.Error())
		}
		runtime.LoadPlatformModule(rt, pkg)
	})
}
