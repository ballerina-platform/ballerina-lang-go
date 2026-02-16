package projects

import (
	"ballerina-lang-go/context"
)

type Environment struct {
	compilerCtx *context.CompilerContext
}

func NewEnvironment(cx *context.CompilerContext) *Environment {
	return &Environment{
		compilerCtx: cx,
	}
}

func (e *Environment) compilerContext() *context.CompilerContext {
	return e.compilerCtx
}
