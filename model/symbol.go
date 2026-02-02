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

package model

import (
	"ballerina-lang-go/semtypes"
)

type Scope interface {
	GetSymbol(name string) (Symbol, bool)
	GetPrefixedSymbol(prefix, name string) (Symbol, bool)
	AddSymbol(name string, symbol Symbol)
}

type Symbol interface {
	Name() string
	Type() semtypes.SemType
	Kind() SymbolKind
	IsPublic() bool
}

type SymbolKind uint

const (
	SymbolKindType SymbolKind = iota
	SymbolKindConstant
	SymbolKindVariable
	SymbolKindParemeter
	SymbolKindFunction
)

type (
	PackageIdentifier struct {
		Organization string
		Package      string
		Version      string
	}

	SymbolRef struct {
		Package PackageIdentifier
		Index   int
		symbol  Symbol
	}

	ModuleScope struct {
		Main       SymbolSpace
		Prefix     map[string]ExportedSymbolSpace
		Annotation SymbolSpace
	}

	// ExportedSymbolSpace is a readonly representation of symbols exported by a Module
	ExportedSymbolSpace struct {
		Main       *SymbolSpace
		Annotation *SymbolSpace
	}

	blockScopeBase struct {
		parent Scope
		Main   SymbolSpace
	}

	// This is a delimiter to help detect if we need to capture a symbol as a closure
	// TODO: need to think how to implement closures correctly
	FunctionScope struct {
		blockScopeBase
	}

	BlockScope struct {
		blockScopeBase
	}

	SymbolSpace struct {
		pkg         PackageIdentifier
		lookupTable map[string]int
		symbols     []Symbol
	}

	symbolBase struct {
		name     string
		ty       semtypes.SemType
		isPublic bool
	}

	TypeSymbol struct {
		symbolBase
	}

	ValueSymbol struct {
		symbolBase
		isConst     bool
		isParameter bool
	}

	FunctionSymbol struct {
		symbolBase
		Signature FunctionSignature
	}

	FunctionSignature struct {
		ParamTypes []semtypes.SemType
		ReturnType semtypes.SemType
		// RestParamType is nil if there is no rest param
		RestParamType semtypes.SemType
	}
)

var _ Scope = &ModuleScope{}
var _ Scope = &FunctionScope{}
var _ Scope = &BlockScope{}
var _ Symbol = &TypeSymbol{}
var _ Symbol = &ValueSymbol{}
var _ Symbol = &FunctionSymbol{}
var _ Symbol = &SymbolRef{}

func (space *SymbolSpace) AddSymbol(name string, symbol Symbol) {
	if _, ok := symbol.(*SymbolRef); ok {
		panic("SymbolRef cannot be added to a SymbolSpace")
	}
	space.lookupTable[name] = len(space.symbols)
	space.symbols = append(space.symbols, symbol)
}

func (space *SymbolSpace) GetSymbol(name string) (SymbolRef, bool) {
	index, ok := space.lookupTable[name]
	if !ok {
		return SymbolRef{}, false
	}
	return SymbolRef{Package: space.pkg, Index: index, symbol: space.symbols[index]}, true
}

func NewSymbolSpace(packageId PackageID) *SymbolSpace {
  pkg := PackageIdentifier{
    Organization: packageId.OrgName.Value(),
    Package: packageId.PkgName.Value(),
    Version: packageId.Version.Value(),
  }
	return &SymbolSpace{pkg: pkg, lookupTable: make(map[string]int), symbols: make([]Symbol, 0)}
}

func NewFunctionScope(parent Scope, pkg PackageID) *FunctionScope {
	return &FunctionScope{
		blockScopeBase: blockScopeBase{
			parent: parent,
			Main:   *NewSymbolSpace(pkg),
		},
	}
}

func NewBlockScope(parent Scope, pkg PackageID) *BlockScope {
	return &BlockScope{
		blockScopeBase: blockScopeBase{
			parent: parent,
			Main:   *NewSymbolSpace(pkg),
		},
	}
}

func (ms *ModuleScope) Exports() ExportedSymbolSpace {
	// FIXME: this needs to only export public symbols but this means we need to correct indexes in symbol refs how to do that?
	// -- Or do we need to do that correction
	return ExportedSymbolSpace{
		Main:       &ms.Main,
		Annotation: &ms.Annotation,
	}
}

func (ms *ModuleScope) GetSymbol(name string) (Symbol, bool) {
  ref, ok := ms.Main.GetSymbol(name)
  if !ok {
    return nil, false
  }
  return &ref, true
}

func (ms *ModuleScope) GetPrefixedSymbol(prefix, name string) (Symbol, bool) {
  if prefix == "" {
    return ms.GetSymbol(name)
  }
  exported, ok := ms.Prefix[prefix]
  if !ok {
    return nil, false
  }
  ref, ok := exported.Main.GetSymbol(name)
  if !ok {
    return nil, false
  }
  return &ref, true
}

func (ms *ModuleScope) AddSymbol(name string, symbol Symbol) {
  ms.Main.AddSymbol(name, symbol)
}

func (ms *ModuleScope) AddAnnotationSymbol(name string, symbol Symbol) {
  ms.Annotation.AddSymbol(name, symbol)
}

func (space *ExportedSymbolSpace) GetSymbol(name string) (SymbolRef, bool) {
  return space.Main.GetSymbol(name)
}

func (bs *blockScopeBase) GetSymbol(name string) (Symbol, bool) {
	ref, ok := bs.Main.GetSymbol(name)
	if ok {
		return &ref, true
	}
	return bs.parent.GetSymbol(name)
}

func (bs *blockScopeBase) GetPrefixedSymbol(prefix, name string) (Symbol, bool) {
	return bs.parent.GetPrefixedSymbol(prefix, name)
}

func (bs *blockScopeBase) AddSymbol(name string, symbol Symbol) {
	bs.Main.AddSymbol(name, symbol)
}


func (ba *symbolBase) Name() string {
	return ba.name
}

func (ba *symbolBase) Type() semtypes.SemType {
	return ba.ty
}

func (ba *symbolBase) IsPublic() bool {
	return ba.isPublic
}

func (ref *SymbolRef) Name() string {
	return ref.symbol.Name()
}

func (ref *SymbolRef) Type() semtypes.SemType {
	return ref.symbol.Type()
}

func (ref *SymbolRef) Kind() SymbolKind {
	return ref.symbol.Kind()
}

func (ref *SymbolRef) IsPublic() bool {
	return ref.symbol.IsPublic()
}

func (ts *TypeSymbol) Kind() SymbolKind {
	return SymbolKindType
}

func (vs *ValueSymbol) Kind() SymbolKind {
	if vs.isConst {
		return SymbolKindConstant
	}
	if vs.isParameter {
		return SymbolKindParemeter
	}
	return SymbolKindVariable
}

func (fs *FunctionSymbol) Kind() SymbolKind {
	return SymbolKindFunction
}

func NewFunctionSymbol(name string, signature FunctionSignature, isPublic bool) FunctionSymbol {
	return FunctionSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		Signature: signature,
	}
}

func NewValueSymbol(name string, isPublic bool, isConst bool, isParameter bool) ValueSymbol {
	return ValueSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		isConst: isConst,
		isParameter: isParameter,
	}
}

func NewTypeSymbol(name string, isPublic bool) TypeSymbol {
	return TypeSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
	}
}
