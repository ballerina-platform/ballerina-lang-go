package array

import (
	"ballerina-lang-go/runtime/api"
	"fmt"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.array"
)

func initArrayModule(rt *api.Runtime) {
	api.RegisterExternFunction(rt.Registry, orgName, moduleName, "push", func(args []any) (any, error) {
		if arr, ok := args[0].(*[]any); ok {
			*arr = append(*arr, args[1:]...)
			return nil, nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	api.RegisterExternFunction(rt.Registry, orgName, moduleName, "length", func(args []any) (any, error) {
		if arr, ok := args[0].(*[]any); ok {
			return int64(len(*arr)), nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
}

func init() {
	api.RegisterModuleInitializer(initArrayModule)
}
