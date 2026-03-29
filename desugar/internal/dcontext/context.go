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
	"fmt"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

// PackageContext holds shared state for desugaring a single package.
// Fields are private to enforce access through methods.
type PackageContext struct {
	compilerCtx          *context.CompilerContext
	pkg                  *ast.BLangPackage
	importedSymbols      map[string]model.ExportedSymbolSpace
	importMu             sync.Mutex
	addedImplicitImports map[string]bool
	desugarSymbolCounter int
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

func (ctx *PackageContext) SymbolType(symbol model.SymbolRef) semtypes.SemType {
	return ctx.compilerCtx.SymbolType(symbol)
}

func (ctx *PackageContext) NewFunctionScope(parent model.Scope) *model.FunctionScope {
	return ctx.compilerCtx.NewFunctionScope(parent, *ctx.pkg.PackageID)
}

func (ctx *PackageContext) GetSymbol(ref model.SymbolRef) model.Symbol {
	return ctx.compilerCtx.GetSymbol(ref)
}

func (ctx *PackageContext) GetSymbolType(ref model.SymbolRef) semtypes.SemType {
	return ctx.compilerCtx.SymbolType(ref)
}

func (ctx *PackageContext) SetSymbolType(ref model.SymbolRef, ty semtypes.SemType) {
	ctx.compilerCtx.SetSymbolType(ref, ty)
}

func (ctx *PackageContext) TypeEnv() semtypes.Env {
	return ctx.compilerCtx.GetTypeEnv()
}

func (ctx *PackageContext) NextDesugarSymbolName() string {
	name := fmt.Sprintf("$desugar$%d", ctx.desugarSymbolCounter)
	ctx.desugarSymbolCounter++
	return name
}

func (ctx *PackageContext) AddSymbolToSameSpace(ref model.SymbolRef, name string, symbol model.Symbol) model.SymbolRef {
	return ctx.compilerCtx.AddSymbolToSameSpace(ref, name, symbol)
}

func (ctx *PackageContext) InternalError(msg string) {
	ctx.compilerCtx.InternalError(msg, nil)
}

func (ctx *PackageContext) Unimplemented(msg string) {
	ctx.compilerCtx.Unimplemented(msg, nil)
}
