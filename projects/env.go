/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package projects

import (
	"io/fs"

	"ballerina-lang-go/context"
)

type Environment struct {
	fsys        fs.FS
	compilerCtx *context.CompilerContext
}

func newEnvironment(fsys fs.FS, cx *context.CompilerContext) *Environment {
	return &Environment{
		fsys:        fsys,
		compilerCtx: cx,
	}
}

func (e *Environment) compilerContext() *context.CompilerContext {
	return e.compilerCtx
}

func (e *Environment) fs() fs.FS {
	return e.fsys
}
