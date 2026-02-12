package array

import (
	"ballerina-lang-go/runtime"
	"fmt"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.array"
)

func initArrayModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "push", func(args []any) (any, error) {
		if arr, ok := args[0].(*[]any); ok {
			*arr = append(*arr, args[1:]...)
			return nil, nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "length", func(args []any) (any, error) {
		if arr, ok := args[0].(*[]any); ok {
			return len(*arr), nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
}

func init() {
	runtime.RegisterModuleInitializer(initArrayModule)
}
