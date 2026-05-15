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

package floatruntime

import (
	"math"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.float"
)

func initFloatModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "pow", func(args []values.BalValue) (values.BalValue, error) {
		x := args[0].(float64)
		y := args[1].(float64)
		return math.Pow(x, y), nil
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "sqrt", func(args []values.BalValue) (values.BalValue, error) {
		x := args[0].(float64)
		return math.Sqrt(x), nil
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "abs", func(args []values.BalValue) (values.BalValue, error) {
		x := args[0].(float64)
		return math.Abs(x), nil
	})
}

func init() {
	runtime.RegisterModuleInitializer(initFloatModule)
}
