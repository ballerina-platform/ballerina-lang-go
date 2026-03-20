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

package maprt

import (
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
	"sync"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.map"
)

func mapLength(args []values.BalValue) (values.BalValue, error) {
	m := args[0].(*values.Map)
	return int64(m.Len()), nil
}

func mapKeys(rt *runtime.Runtime) func(args []values.BalValue) (values.BalValue, error) {
	var stringArrayTy semtypes.SemType
	var stringArrayOnce sync.Once
	return func(args []values.BalValue) (values.BalValue, error) {
		stringArrayOnce.Do(func() {
			env := rt.GetTypeEnv()
			ld := semtypes.NewListDefinition()
			stringArrayTy = ld.DefineListTypeWrappedWithEnvSemType(env, &semtypes.STRING)
		})
		m := args[0].(*values.Map)
		keys := m.Keys()
		items := make([]values.BalValue, len(keys))
		for i, k := range keys {
			items[i] = k
		}
		list := values.NewList(0, stringArrayTy, nil)
		list.Append(items...)
		return list, nil
	}
}

func mapRemove(args []values.BalValue) (values.BalValue, error) {
	m := args[0].(*values.Map)
	key := args[1].(string)
	val, _ := m.Get(key)
	m.Delete(key)
	return val, nil
}

func initMapModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "length", mapLength)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "keys", mapKeys(rt))
	runtime.RegisterExternFunction(rt, orgName, moduleName, "remove", mapRemove)
}

func init() {
	runtime.RegisterModuleInitializer(initMapModule)
}
