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

package native

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"ballerina-lang-go/decimal"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
)

type fileIOTypes struct {
	strArrTy          semtypes.SemType
	byteArrTy         semtypes.SemType
	jsonListTy        semtypes.SemType
	jsonMapTy         semtypes.SemType
	stringMapTy       semtypes.SemType
	stringMapAtomicTy *semtypes.MappingAtomicType
}

func fileIOError(msg string) values.BalValue {
	return values.NewErrorWithMessage(msg)
}

func splitLines(data []byte) []string {
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func toByteSlice(list *values.List) ([]byte, bool) {
	b := make([]byte, list.Len())
	for i := range list.Len() {
		n, ok := list.Get(i).(int64)
		if !ok || n < 0 || n > 255 {
			return nil, false
		}
		b[i] = byte(n)
	}
	return b, true
}

func balValueToGoJSON(v values.BalValue) any {
	switch t := v.(type) {
	case nil:
		return nil
	case bool:
		return t
	case int64:
		return t
	case float64:
		return t
	case *decimal.Decimal:
		return json.RawMessage(t.String())
	case string:
		return t
	case *values.Map:
		m := make(map[string]any, t.Len())
		for _, k := range t.Keys() {
			val, _ := t.Get(k)
			m[k] = balValueToGoJSON(val)
		}
		return m
	case *values.List:
		s := make([]any, t.Len())
		for i := range t.Len() {
			s[i] = balValueToGoJSON(t.Get(i))
		}
		return s
	default:
		return nil
	}
}

func goJSONToBalValue(tc semtypes.Context, v any, jsonListTy, jsonMapTy semtypes.SemType) values.BalValue {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return v
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		f, _ := v.Float64()
		return f
	case string:
		return v
	case []any:
		items := make([]values.BalValue, len(v))
		for i, elem := range v {
			items[i] = goJSONToBalValue(tc, elem, jsonListTy, jsonMapTy)
		}
		return values.NewList(jsonListTy, semtypes.ToListAtomicType(tc, jsonListTy), false, nil, 0, items)
	case map[string]any:
		m := values.NewMap(jsonMapTy, semtypes.ToMappingAtomicType(tc, jsonMapTy), false, nil)
		for k, val := range v {
			m.Put(tc, k, goJSONToBalValue(tc, val, jsonListTy, jsonMapTy))
		}
		return m
	default:
		return nil
	}
}

func initFileIOModule(rt *runtime.Runtime) {
	var (
		once  sync.Once
		types fileIOTypes
	)
	ensureTypes := func() {
		once.Do(func() {
			env := rt.GetTypeEnv()
			sld := semtypes.NewListDefinition()
			types.strArrTy = sld.DefineListTypeWrappedWithEnvSemType(env, semtypes.STRING)
			bld := semtypes.NewListDefinition()
			types.byteArrTy = bld.DefineListTypeWrappedWithEnvSemType(env, semtypes.BYTE)
			typCtx := semtypes.ContextFrom(env)
			jsonTy := semtypes.CreateJSON(typCtx)
			jmd := semtypes.NewMappingDefinition()
			types.jsonMapTy = jmd.DefineMappingTypeWrapped(env, nil, jsonTy)
			jld := semtypes.NewListDefinition()
			types.jsonListTy = jld.DefineListTypeWrappedWithEnvSemType(env, jsonTy)
			smd := semtypes.NewMappingDefinition()
			types.stringMapTy = smd.DefineMappingTypeWrapped(env, nil, semtypes.STRING)
			types.stringMapAtomicTy = semtypes.ToMappingAtomicType(typCtx, types.stringMapTy)
		})
	}

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileReadString",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			data, err := rt.Platform().FS.ReadFile(path)
			if err != nil {
				return fileIOError(fmt.Sprintf("error while reading file '%s': %s", path, err.Error())), nil
			}
			return strings.Join(splitLines(data), "\n"), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileReadLines",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			path, _ := args[0].(string)
			data, err := rt.Platform().FS.ReadFile(path)
			if err != nil {
				return fileIOError(fmt.Sprintf("error while reading file '%s': %s", path, err.Error())), nil
			}
			lines := splitLines(data)
			items := make([]values.BalValue, len(lines))
			for i, line := range lines {
				items[i] = line
			}
			return values.NewList(types.strArrTy, semtypes.ToListAtomicType(ctx.TypeCtx, types.strArrTy), false, nil, 0, items), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileReadBytes",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			path, _ := args[0].(string)
			data, err := rt.Platform().FS.ReadFile(path)
			if err != nil {
				return fileIOError(fmt.Sprintf("error while reading file '%s': %s", path, err.Error())), nil
			}
			items := make([]values.BalValue, len(data))
			for i, b := range data {
				items[i] = int64(b)
			}
			return values.NewList(types.byteArrTy, semtypes.ToListAtomicType(ctx.TypeCtx, types.byteArrTy), false, nil, 0, items), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileReadJson",
		func(ctx *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			path, _ := args[0].(string)
			data, err := rt.Platform().FS.ReadFile(path)
			if err != nil {
				return fileIOError(fmt.Sprintf("error while reading file '%s': %s", path, err.Error())), nil
			}
			dec := json.NewDecoder(strings.NewReader(string(data)))
			dec.UseNumber()
			var raw any
			if err := dec.Decode(&raw); err != nil {
				return fileIOError(fmt.Sprintf("error while parsing JSON from file '%s': %s", path, err.Error())), nil
			}
			return goJSONToBalValue(ctx.TypeCtx, raw, types.jsonListTy, types.jsonMapTy), nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileWriteString",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			content, _ := args[1].(string)
			option, _ := args[2].(string)
			data := []byte(content)
			var err error
			if option == "APPEND" {
				err = rt.Platform().FS.AppendFile(path, data)
			} else {
				err = rt.Platform().FS.WriteFile(path, data)
			}
			if err != nil {
				return fileIOError(fmt.Sprintf("error while writing to file '%s': %s", path, err.Error())), nil
			}
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileWriteLines",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			list, _ := args[1].(*values.List)
			option, _ := args[2].(string)
			var sb strings.Builder
			for i := range list.Len() {
				line, _ := list.Get(i).(string)
				sb.WriteString(line)
				sb.WriteByte('\n')
			}
			data := []byte(sb.String())
			var err error
			if option == "APPEND" {
				err = rt.Platform().FS.AppendFile(path, data)
			} else {
				err = rt.Platform().FS.WriteFile(path, data)
			}
			if err != nil {
				return fileIOError(fmt.Sprintf("error while writing to file '%s': %s", path, err.Error())), nil
			}
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileWriteBytes",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			list, _ := args[1].(*values.List)
			option, _ := args[2].(string)
			data, ok := toByteSlice(list)
			if !ok {
				return fileIOError("invalid byte value in content array"), nil
			}
			var err error
			if option == "APPEND" {
				err = rt.Platform().FS.AppendFile(path, data)
			} else {
				err = rt.Platform().FS.WriteFile(path, data)
			}
			if err != nil {
				return fileIOError(fmt.Sprintf("error while writing to file '%s': %s", path, err.Error())), nil
			}
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileWriteJson",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			data, err := json.Marshal(balValueToGoJSON(args[1]))
			if err != nil {
				return fileIOError(fmt.Sprintf("error while serializing JSON for file '%s': %s", path, err.Error())), nil
			}
			if err := rt.Platform().FS.WriteFile(path, data); err != nil {
				return fileIOError(fmt.Sprintf("error while writing to file '%s': %s", path, err.Error())), nil
			}
			return nil, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileReadXml",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			ensureTypes()
			path, _ := args[0].(string)
			data, err := rt.Platform().FS.ReadFile(path)
			if err != nil {
				return fileIOError(fmt.Sprintf("error while reading file '%s': %s", path, err.Error())), nil
			}
			xmlVal, parseErr := parseXMLFromBytes(data, types.stringMapTy, types.stringMapAtomicTy)
			if parseErr != nil {
				return fileIOError(fmt.Sprintf("error while parsing XML from file '%s': %s", path, parseErr.Error())), nil
			}
			return xmlVal, nil
		})

	runtime.RegisterExternFunction(rt, orgName, moduleName, "externFileWriteXml",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			path, _ := args[0].(string)
			content, _ := args[1].(values.XMLValue)
			option, _ := args[2].(string)
			data := []byte(content.XMLString())
			var err error
			if option == "APPEND" {
				err = rt.Platform().FS.AppendFile(path, data)
			} else {
				err = rt.Platform().FS.WriteFile(path, data)
			}
			if err != nil {
				return fileIOError(fmt.Sprintf("error while writing to file '%s': %s", path, err.Error())), nil
			}
			return nil, nil
		})
}

func init() {
	runtime.RegisterModuleInitializer(initFileIOModule)
}
