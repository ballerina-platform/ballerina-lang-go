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
	"fmt"
	"net/url"
	"strings"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "url"
)

// encodeExtern replicates Java URLEncoder.encode() + post-processing:
//   - space -> %20 (Java uses +, post-processing converts to %20)
//   - * -> %2A (Java keeps *, post-processing converts to %2A; Go encodes directly)
//   - ~ -> ~ (Java encodes as %7E, post-processing reverts; Go never encodes ~)
//
// Only UTF-8 charset is supported; Go strings are always UTF-8.
func encodeExtern() extern.NativeFunc {
	return func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
		value, ok := args[0].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while encoding. invalid string argument"), nil
		}
		charset, ok := args[1].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while encoding. invalid charset argument"), nil
		}
		if !strings.EqualFold(charset, "UTF-8") {
			return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while encoding. %s", charset)), nil
		}
		encoded := url.QueryEscape(value)
		// QueryEscape encodes space as +; convert to %20 to match Java behaviour.
		encoded = strings.ReplaceAll(encoded, "+", "%20")
		return encoded, nil
	}
}

// decodeExtern replicates Java URLDecoder.decode().
// Only UTF-8 charset is supported; Go strings are always UTF-8.
func decodeExtern() extern.NativeFunc {
	return func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
		value, ok := args[0].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while decoding. invalid string argument"), nil
		}
		charset, ok := args[1].(string)
		if !ok {
			return values.NewErrorWithMessage("Error occurred while decoding. invalid charset argument"), nil
		}
		if !strings.EqualFold(charset, "UTF-8") {
			return values.NewErrorWithMessage(fmt.Sprintf("Error occurred while decoding. %s", charset)), nil
		}
		decoded, err := url.QueryUnescape(value)
		if err != nil {
			return values.NewErrorWithMessage("Error occurred while decoding. " + err.Error()), nil
		}
		return decoded, nil
	}
}

func initURLModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "encode", encodeExtern())
	runtime.RegisterExternFunction(rt, orgName, moduleName, "decode", decodeExtern())
}

func init() {
	runtime.RegisterModuleInitializer(initURLModule)
}
