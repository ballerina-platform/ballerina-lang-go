package array

import (
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"
	"fmt"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.array"
)

func initArrayModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "push", func(args []values.BalValue) (values.BalValue, error) {
		if list, ok := args[0].(*values.List); ok {
			list.Append(args[1:]...)
			return nil, nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "length", func(args []values.BalValue) (values.BalValue, error) {
		if list, ok := args[0].(*values.List); ok {
			return int64(list.Len()), nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
}

func init() {
	runtime.RegisterModuleInitializer(initArrayModule)
}
