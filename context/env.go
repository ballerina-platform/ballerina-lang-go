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

package context

import (
	"strconv"
	"sync"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type CompilerEnvironment struct {
	anonTypeCount     map[*model.PackageID]int
	anonFuncCount     map[*model.PackageID]int
	packageInterner   *model.PackageIDInterner
	symbolSpaces      []*model.SymbolSpace
	symbolSpacesMu    sync.RWMutex // we need this because desugaring add new init functions concurrently we shouldn't need this if the spaces are scoped to the module, may be we should do that?
	typeEnv           semtypes.Env
	underlyingSymbol  sync.Map
	statsEnabled      bool
	diagnosticContext *diagnostics.DiagnosticEnv
}

func (c *CompilerEnvironment) DiagnosticEnv() *diagnostics.DiagnosticEnv {
	return c.diagnosticContext
}

func (c *CompilerEnvironment) NewSymbolSpace(packageID model.PackageID) *model.SymbolSpace {
	c.symbolSpacesMu.Lock()
	space := model.NewSymbolSpaceInner(packageID, len(c.symbolSpaces))
	c.symbolSpaces = append(c.symbolSpaces, space)
	c.symbolSpacesMu.Unlock()
	return space
}

func (c *CompilerEnvironment) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   c.NewSymbolSpace(pkg),
		},
	}
}

func (c *CompilerEnvironment) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   c.NewSymbolSpace(pkg),
		},
	}
}

func (c *CompilerEnvironment) GetSymbol(symbol model.SymbolRef) model.Symbol {
	c.symbolSpacesMu.RLock()
	symbolSpace := c.symbolSpaces[symbol.SpaceIndex]
	c.symbolSpacesMu.RUnlock()
	return symbolSpace.SymbolAt(symbol.Index)
}

func (c *CompilerEnvironment) AddSymbolToSameSpace(ref model.SymbolRef, name string, symbol model.Symbol) model.SymbolRef {
	c.symbolSpacesMu.RLock()
	space := c.symbolSpaces[ref.SpaceIndex]
	c.symbolSpacesMu.RUnlock()
	space.AddSymbol(name, symbol)
	newRef, _ := space.GetSymbol(name)
	return newRef
}

// CreateNarrowedSymbol create a narrowed symbol for the given baseRef symbol. IMPORTANT: baseRef must be the actual symbol
// not a narrowed symbol.
func (c *CompilerEnvironment) CreateNarrowedSymbol(baseRef model.SymbolRef) model.SymbolRef {
	c.symbolSpacesMu.RLock()
	symbolSpace := c.symbolSpaces[baseRef.SpaceIndex]
	c.symbolSpacesMu.RUnlock()
	underlyingSymbolCopy := c.GetSymbol(baseRef).Copy()
	symbolIndex := symbolSpace.AppendSymbol(underlyingSymbolCopy)
	narrowedSymbol := model.SymbolRef{
		Package:    baseRef.Package,
		SpaceIndex: baseRef.SpaceIndex,
		Index:      symbolIndex,
	}
	c.underlyingSymbol.Store(narrowedSymbol, baseRef)
	return narrowedSymbol
}

func (c *CompilerEnvironment) CreateFunctionSymbol(space *model.SymbolSpace, name string, signature model.FunctionSignature, fnTy semtypes.SemType) model.SymbolRef {
	sym := model.NewFunctionSymbol(name, signature, false)
	sym.SetType(fnTy)
	symbolIndex := space.AppendSymbol(sym)
	return space.RefAt(symbolIndex)
}

func (c *CompilerEnvironment) UnnarrowedSymbol(symbol model.SymbolRef) model.SymbolRef {
	if underlying, ok := c.underlyingSymbol.Load(symbol); ok {
		return underlying.(model.SymbolRef)
	}
	return symbol
}

func (c *CompilerEnvironment) SymbolName(symbol model.SymbolRef) string {
	return c.GetSymbol(symbol).Name()
}

func (c *CompilerEnvironment) SymbolType(symbol model.SymbolRef) semtypes.SemType {
	return c.GetSymbol(symbol).Type()
}

func (c *CompilerEnvironment) SymbolKind(symbol model.SymbolRef) model.SymbolKind {
	return c.GetSymbol(symbol).Kind()
}

func (c *CompilerEnvironment) SymbolIsPublic(symbol model.SymbolRef) bool {
	return c.GetSymbol(symbol).IsPublic()
}

func (c *CompilerEnvironment) SetSymbolType(symbol model.SymbolRef, ty semtypes.SemType) {
	c.GetSymbol(symbol).SetType(ty)
}

func (c *CompilerEnvironment) GetDefaultPackage() *model.PackageID {
	return c.packageInterner.GetDefaultPackage()
}

func (c *CompilerEnvironment) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(c.packageInterner, orgName, nameComps, version)
}

func NewCompilerEnvironment(typeEnv semtypes.Env, statsEnabled bool) *CompilerEnvironment {
	return &CompilerEnvironment{
		anonTypeCount:     make(map[*model.PackageID]int),
		anonFuncCount:     make(map[*model.PackageID]int),
		packageInterner:   model.DefaultPackageIDInterner,
		typeEnv:           typeEnv,
		statsEnabled:      statsEnabled,
		diagnosticContext: diagnostics.NewDiagnosticEnv(),
	}
}

// GetTypeEnv returns the type environment for this context
func (c *CompilerEnvironment) GetTypeEnv() semtypes.Env {
	return c.typeEnv
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (c *CompilerEnvironment) GetNextAnonymousFunctionKey(packageID *model.PackageID) string {
	nextValue := c.anonFuncCount[packageID]
	c.anonFuncCount[packageID] = nextValue + 1
	return ANON_PREFIX + "Func$_" + strconv.Itoa(nextValue)
}

func (c *CompilerEnvironment) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := c.anonTypeCount[packageID]
	c.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
