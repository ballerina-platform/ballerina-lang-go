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

package http

import (
	"fmt"
	"strings"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "http"
)

func parseHeader(input string) (*values.List, error) {
	segments := strings.Split(input, ",")
	list := values.NewList(0, semtypes.LIST, nil)
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			return nil, fmt.Errorf("invalid header value: empty segment")
		}
		parts := strings.Split(seg, ";")
		headerVal := strings.TrimSpace(parts[0])
		if headerVal == "" {
			return nil, fmt.Errorf("invalid header value: missing value before parameters")
		}
		params := values.NewMap(semtypes.MAPPING)
		for _, param := range parts[1:] {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}
			eqIdx := strings.IndexByte(param, '=')
			if eqIdx < 0 {
				params.Put(strings.ToLower(param), "")
				continue
			}
			key := strings.ToLower(strings.TrimSpace(param[:eqIdx]))
			val := strings.TrimSpace(param[eqIdx+1:])
			if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}
			params.Put(key, val)
		}
		entry := values.NewMap(semtypes.MAPPING)
		entry.Put("value", headerVal)
		entry.Put("params", params)
		list.Append(entry)
	}
	return list, nil
}

func initHttpModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "parseHeader",
		func(args []values.BalValue) (values.BalValue, error) {
			input, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("parseHeader: expected string argument")
			}
			result, err := parseHeader(input)
			if err != nil {
				return values.NewErrorWithMessage(err.Error()), nil
			}
			return result, nil
		})
}

func init() {
	runtime.RegisterModuleInitializer(initHttpModule)
}
