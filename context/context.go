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
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
	"strconv"
)

// TODO: consider moving type resolution env in to this
type CompilerContext struct {
	anonTypeCount   map[*model.PackageID]int
	packageInterner *model.PackageIDInterner
	symbolSpaces    []*model.SymbolSpace
}

func (this *CompilerContext) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageId, len(this.symbolSpaces))
	this.symbolSpaces = append(this.symbolSpaces, space)
	return space
}

func (this *CompilerContext) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) GetSymbol(symbol model.Symbol) model.Symbol {
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		symbolSpace := this.symbolSpaces[refSymbol.SpaceIndex]
		return symbolSpace.Symbols[refSymbol.Index]
	}
	return symbol
}

func (this *CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.packageInterner.GetDefaultPackage()
}

func (this *CompilerContext) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(this.packageInterner, orgName, nameComps, version)
}

func (this *CompilerContext) Unimplemented(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Unimplemented: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Unimplemented: %s", message))
}

func (this *CompilerContext) SemanticError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Semantic error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Semantic error: %s", message))
}

// TODO: implement these properly
func (this *CompilerContext) SyntaxError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Syntax error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Syntax error: %s", message))
}

func (this *CompilerContext) InternalError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Internal error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Internal error: %s", message))
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
	}
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := this.anonTypeCount[packageID]
	this.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}
