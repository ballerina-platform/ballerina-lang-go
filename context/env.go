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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"strconv"
	"sync"
)

type CompilerEnvironment struct {
	anonTypeCount    map[*model.PackageID]int
	packageInterner  *model.PackageIDInterner
	symbolSpaces     []*model.SymbolSpace
	typeEnv          semtypes.Env
	underlyingSymbol sync.Map
	typeDefns        map[model.SymbolRef]model.TypeDefinition
}

func (c *CompilerEnvironment) NewSymbolSpace(packageID model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageID, len(c.symbolSpaces))
	c.symbolSpaces = append(c.symbolSpaces, space)
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
	symbolSpace := c.symbolSpaces[symbol.SpaceIndex]
	return symbolSpace.SymbolAt(symbol.Index)
}

// CreateNarrowedSymbol create a narrowed symbol for the given baseRef symbol. IMPORTANT: baseRef must be the actual symbol
// not a narrowed symbol.
func (c *CompilerEnvironment) CreateNarrowedSymbol(baseRef model.SymbolRef) model.SymbolRef {
	symbolSpace := c.symbolSpaces[baseRef.SpaceIndex]
	underlyingSymbolCopy := *c.GetSymbol(baseRef).(*model.ValueSymbol)
	symbolIndex := symbolSpace.AppendSymbol(&underlyingSymbolCopy)
	narrowedSymbol := model.SymbolRef{
		Package:    baseRef.Package,
		SpaceIndex: baseRef.SpaceIndex,
		Index:      symbolIndex,
	}
	c.underlyingSymbol.Store(narrowedSymbol, baseRef)
	return narrowedSymbol
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

func (c *CompilerEnvironment) SetTypeDefinition(symbol model.SymbolRef, defn model.TypeDefinition) {
	c.typeDefns[symbol] = defn
}

func (c *CompilerEnvironment) GetTypeDefinition(symbol model.SymbolRef) (model.TypeDefinition, bool) {
	defn, ok := c.typeDefns[symbol]
	return defn, ok
}

func NewCompilerEnvironment(typeEnv semtypes.Env) *CompilerEnvironment {
	return &CompilerEnvironment{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
		typeEnv:         typeEnv,
		typeDefns:       make(map[model.SymbolRef]model.TypeDefinition),
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

func (c *CompilerEnvironment) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := c.anonTypeCount[packageID]
	c.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
