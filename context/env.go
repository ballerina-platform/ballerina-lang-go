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
	"fmt"
	"strconv"
	"sync"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// TODO: consider moving type resolution env in to this
type CompilerEnvironment struct {
	anonTypeCount    map[*model.PackageID]int
	packageInterner  *model.PackageIDInterner
	symbolSpaces     []*model.SymbolSpace
	typeEnv          semtypes.Env
	underlyingSymbol sync.Map
}

func (this *CompilerEnvironment) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageId, len(this.symbolSpaces))
	this.symbolSpaces = append(this.symbolSpaces, space)
	return space
}

func (this *CompilerEnvironment) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerEnvironment) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerEnvironment) GetSymbol(symbol model.SymbolRef) model.Symbol {
	symbolSpace := this.symbolSpaces[symbol.SpaceIndex]
	return symbolSpace.Symbols[symbol.Index]
}

// CreateNarrowedSymbol create a narrowed symbol for the given baseRef symbol. IMPORTANT: baseRef must be the actual symbol
// not a narrowed symbol.
func (this *CompilerEnvironment) CreateNarrowedSymbol(baseRef model.SymbolRef) model.SymbolRef {
	symbolSpace := this.symbolSpaces[baseRef.SpaceIndex]
	underlyingSymbolCopy := *this.GetSymbol(baseRef).(*model.ValueSymbol)
	symbolIndex := symbolSpace.AppendSymbol(&underlyingSymbolCopy)
	narrowedSymbol := model.SymbolRef{
		Package:    baseRef.Package,
		SpaceIndex: baseRef.SpaceIndex,
		Index:      symbolIndex,
	}
	this.underlyingSymbol.Store(narrowedSymbol, baseRef)
	return narrowedSymbol
}

func (this *CompilerEnvironment) UnnarrowedSymbol(symbol model.SymbolRef) model.SymbolRef {
	if underlying, ok := this.underlyingSymbol.Load(symbol); ok {
		return underlying.(model.SymbolRef)
	}
	return symbol
}

func (this *CompilerEnvironment) SymbolName(symbol model.SymbolRef) string {
	return this.GetSymbol(symbol).Name()
}

func (this *CompilerEnvironment) SymbolType(symbol model.SymbolRef) semtypes.SemType {
	return this.GetSymbol(symbol).Type()
}

func (this *CompilerEnvironment) SymbolKind(symbol model.SymbolRef) model.SymbolKind {
	return this.GetSymbol(symbol).Kind()
}

func (this *CompilerEnvironment) SymbolIsPublic(symbol model.SymbolRef) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this *CompilerEnvironment) SetSymbolType(symbol model.SymbolRef, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this *CompilerEnvironment) GetDefaultPackage() *model.PackageID {
	return this.packageInterner.GetDefaultPackage()
}

func (this *CompilerEnvironment) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(this.packageInterner, orgName, nameComps, version)
}

func (this *CompilerEnvironment) Unimplemented(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Unimplemented: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Unimplemented: %s", message))
}

func (this *CompilerEnvironment) SemanticError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Semantic error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Semantic error: %s", message))
}

// TODO: implement these properly
func (this *CompilerEnvironment) SyntaxError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Syntax error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Syntax error: %s", message))
}

func (this *CompilerEnvironment) InternalError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Internal error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Internal error: %s", message))
}

func NewCompilerEnvironment(typeEnv semtypes.Env) *CompilerEnvironment {
	return &CompilerEnvironment{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
		typeEnv:         typeEnv,
	}
}

// GetTypeEnv returns the type environment for this context
func (this *CompilerEnvironment) GetTypeEnv() semtypes.Env {
	return this.typeEnv
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerEnvironment) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := this.anonTypeCount[packageID]
	this.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
