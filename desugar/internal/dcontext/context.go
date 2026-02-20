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

package dcontext

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"sync"
)

// PackageContext holds shared state for desugaring a single package.
// Fields are private to enforce access through methods.
type PackageContext struct {
	compilerCtx          *context.CompilerContext
	pkg                  *ast.BLangPackage
	importedSymbols      map[string]model.ExportedSymbolSpace
	importMu             sync.Mutex
	addedImplicitImports map[string]bool
}

func NewPackageContext(compilerCtx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) *PackageContext {
	return &PackageContext{
		compilerCtx:          compilerCtx,
		pkg:                  pkg,
		importedSymbols:      importedSymbols,
		addedImplicitImports: make(map[string]bool),
	}
}

func (ctx *PackageContext) AddImplicitImport(pkgName string, imp ast.BLangImportPackage) {
	ctx.importMu.Lock()
	defer ctx.importMu.Unlock()
	if !ctx.addedImplicitImports[pkgName] {
		ctx.addedImplicitImports[pkgName] = true
		ctx.pkg.Imports = append(ctx.pkg.Imports, imp)
	}
}

func (ctx *PackageContext) GetImportedSymbolSpace(pkgName string) (model.ExportedSymbolSpace, bool) {
	space, ok := ctx.importedSymbols[pkgName]
	return space, ok
}

func (ctx *PackageContext) InternalError(msg string) {
	ctx.compilerCtx.InternalError(msg, nil)
}

func (ctx *PackageContext) Unimplemented(msg string) {
	ctx.compilerCtx.Unimplemented(msg, nil)
}
