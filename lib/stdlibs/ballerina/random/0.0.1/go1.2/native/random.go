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
	cryptorand "crypto/rand"
	"encoding/binary"
	mathrand "math/rand/v2"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "random"
)

func randomError(msg string) values.BalValue {
	return values.NewErrorWithMessage(msg)
}

// cryptoFloat64 returns a uniformly distributed float64 in [0.0, 1.0) sourced from
// the OS cryptographic entropy pool (equivalent to java.security.SecureRandom.nextFloat()).
func cryptoFloat64() float64 {
	var b [8]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		return mathrand.Float64()
	}
	n := binary.BigEndian.Uint64(b[:]) >> 11
	return float64(n) / (1 << 53)
}

func initRandomModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "externCreateDecimal",
		func(_ *extern.Context, _ []values.BalValue) (values.BalValue, error) {
			return cryptoFloat64(), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externCreateIntInRange",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			startRange, _ := args[0].(int64)
			endRange, _ := args[1].(int64)
			if startRange >= endRange {
				return randomError("End range value must be greater than the start range value"), nil
			}
			return startRange + mathrand.Int64N(endRange-startRange), nil
		})
}

func init() {
	runtime.RegisterModuleInitializer(initRandomModule)
}
