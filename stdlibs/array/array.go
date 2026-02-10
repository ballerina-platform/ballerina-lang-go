package array

import (
	"ballerina-lang-go/runtime/api"
	"ballerina-lang-go/runtime/values"
	"fmt"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.array"
)

func initArrayModule(rt *api.Runtime) {
	api.RegisterExternFunction(rt.Registry, orgName, moduleName, "push", func(args []any) (any, error) {
		if list, ok := args[0].(*values.List); ok {
			list.Push(args[1:]...)
			return nil, nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	api.RegisterExternFunction(rt.Registry, orgName, moduleName, "length", func(args []any) (any, error) {
		if list, ok := args[0].(*values.List); ok {
			return int64(list.Len()), nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
}

func init() {
	api.RegisterModuleInitializer(initArrayModule)
}
