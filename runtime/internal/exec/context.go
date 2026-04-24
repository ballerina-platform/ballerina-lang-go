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

package exec

import (
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/runtime/internal/modules"
	"ballerina-lang-go/semtypes"
)

type Context struct {
	registry  *modules.Registry
	callStack callStack
	typeCtx   semtypes.Context
}

func NewContext(reg *modules.Registry) *Context {
	return &Context{
		registry:  reg,
		callStack: callStack{elements: make([]*Frame, 0, 32)},
		typeCtx:   semtypes.TypeCheckContext(reg.GetTypeEnv()),
	}
}

func (ctx *Context) RegisterModule(id *model.PackageID, m *modules.BIRModule) *modules.BIRModule {
	return ctx.registry.RegisterModule(id, m)
}

func (ctx *Context) GetModule(pkgId *model.PackageID) *modules.BIRModule {
	return ctx.registry.GetModule(pkgId)
}

func (ctx *Context) GetBIRFunction(funcName string) *bir.BIRFunction {
	return ctx.registry.GetBIRFunction(funcName)
}

func (ctx *Context) GetNativeFunction(funcName string) *modules.ExternFunction {
	return ctx.registry.GetNativeFunction(funcName)
}

func (ctx *Context) GetClassDef(lookupKey string) *bir.BIRClassDef {
	return ctx.registry.GetClassDef(lookupKey)
}

func (ctx *Context) TypeCheckContext() semtypes.Context {
	return ctx.typeCtx
}

func (ctx *Context) PushFrame(frame *Frame) {
	ctx.callStack.Push(frame)
}

func (ctx *Context) PopFrame() {
	ctx.callStack.Pop()
}

func (ctx *Context) Frames() []*Frame {
	return ctx.callStack.Frames()
}

func (ctx *Context) CallStackDepth() int {
	return len(ctx.callStack.elements)
}
