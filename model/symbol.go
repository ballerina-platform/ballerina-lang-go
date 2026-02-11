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

// These methods should never be called directly. Instead call them via the compiler context.
type Symbol interface {
	Name() string
	Type() semtypes.SemType
	Kind() SymbolKind
	SetType(semtypes.SemType)
	IsPublic() bool
}

// symbolTypeSetter is a private interface for updating symbol types during type resolution.
// All concrete symbol types implement this through symbolBase.
type symbolTypeSetter interface {
	SetType(semtypes.SemType)
}

type FunctionSymbol interface {
	Symbol
	Signature() FunctionSignature
	SetSignature(FunctionSignature)
}

// GenericFunctionSymbol represents functions with [@typeParam] types
type GenericFunctionSymbol interface {
	FunctionSymbol
	// FIXME: this shouldn't care about the return type
	Monomorphize(args []semtypes.SemType, returnTy semtypes.SemType) SymbolRef
	Space() *SymbolSpace
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

	// We are using indeces here with the same rational as RefAtoms, instead of pointers
	SymbolRef struct {
		Package    PackageIdentifier
		Index      int
		SpaceIndex int
	}

	ModuleScope struct {
		Main       *SymbolSpace
		Prefix     map[string]ExportedSymbolSpace
		Annotation *SymbolSpace
	}

	// ExportedSymbolSpace is a readonly representation of symbols exported by a Module
	ExportedSymbolSpace struct {
		Main       *SymbolSpace
		Annotation *SymbolSpace
	}

	BlockScopeBase struct {
		Parent Scope
		Main   *SymbolSpace
	}

	// This is a delimiter to help detect if we need to capture a symbol as a closure
	// TODO: need to think how to implement closures correctly
	FunctionScope struct {
		BlockScopeBase
	}

	BlockScope struct {
		BlockScopeBase
	}

	SymbolSpace struct {
		pkg         PackageIdentifier
		lookupTable map[string]int
		Symbols     []Symbol
		index       int
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

	functionSymbol struct {
		symbolBase
		signature FunctionSignature
	}

	genericFunctionSymbol struct {
		space         *SymbolSpace
		monomorphizer func(s GenericFunctionSymbol, args []semtypes.SemType, ret semtypes.SemType) SymbolRef
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
var _ Symbol = &functionSymbol{}
var _ FunctionSymbol = &functionSymbol{}
var _ GenericFunctionSymbol = &genericFunctionSymbol{}
var _ Symbol = &SymbolRef{}

func (space *SymbolSpace) AddSymbol(name string, symbol Symbol) {
	if _, ok := symbol.(*SymbolRef); ok {
		panic("SymbolRef cannot be added to a SymbolSpace")
	}
	space.lookupTable[name] = len(space.Symbols)
	space.Symbols = append(space.Symbols, symbol)
}

func (space *SymbolSpace) GetSymbol(name string) (SymbolRef, bool) {
	index, ok := space.lookupTable[name]
	if !ok {
		return SymbolRef{}, false
	}
	return SymbolRef{Package: space.pkg, Index: index, SpaceIndex: space.index}, true
}

func NewSymbolSpaceInner(packageId PackageID, index int) *SymbolSpace {
	pkg := PackageIdentifier{
		Organization: packageId.OrgName.Value(),
		Package:      packageId.PkgName.Value(),
		Version:      packageId.Version.Value(),
	}
	return &SymbolSpace{index: index, pkg: pkg, lookupTable: make(map[string]int), Symbols: make([]Symbol, 0)}
}

func (ms *ModuleScope) Exports() ExportedSymbolSpace {
	// FIXME: this needs to only export public symbols but this means we need to correct indexes in symbol refs how to do that?
	// -- Or do we need to do that correction
	// I think the correct way to do this is for references to fail on lookup if the symbol is not exported
	return ExportedSymbolSpace{
		Main:       ms.Main,
		Annotation: ms.Annotation,
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

func (bs *BlockScopeBase) GetSymbol(name string) (Symbol, bool) {
	ref, ok := bs.Main.GetSymbol(name)
	if ok {
		return &ref, true
	}
	return bs.Parent.GetSymbol(name)
}

func (bs *BlockScopeBase) GetPrefixedSymbol(prefix, name string) (Symbol, bool) {
	return bs.Parent.GetPrefixedSymbol(prefix, name)
}

func (bs *BlockScopeBase) AddSymbol(name string, symbol Symbol) {
	bs.Main.AddSymbol(name, symbol)
}

func (ba *symbolBase) Name() string {
	return ba.name
}

func (ba *symbolBase) Type() semtypes.SemType {
	return ba.ty
}

func (ba *symbolBase) SetType(ty semtypes.SemType) {
	ba.ty = ty
}

func (ba *symbolBase) IsPublic() bool {
	return ba.isPublic
}

func (ref *SymbolRef) Name() string {
	panic("unexpected")
}

func (ref *SymbolRef) Type() semtypes.SemType {
	panic("unexpected")
}

func (ref *SymbolRef) SetType(ty semtypes.SemType) {
	panic("unexpected")
}

func (ref *SymbolRef) Kind() SymbolKind {
	panic("unexpected")
}

func (ref *SymbolRef) IsPublic() bool {
	panic("unexpected")
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

func (fs *functionSymbol) Kind() SymbolKind {
	return SymbolKindFunction
}

func (fs *functionSymbol) Signature() FunctionSignature {
	return fs.signature
}

func (fs *functionSymbol) SetSignature(sig FunctionSignature) {
	fs.signature = sig
}

func NewFunctionSymbol(name string, signature FunctionSignature, isPublic bool) FunctionSymbol {
	return &functionSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		signature:  signature,
	}
}

func NewValueSymbol(name string, isPublic bool, isConst bool, isParameter bool) ValueSymbol {
	return ValueSymbol{
		symbolBase:  symbolBase{name: name, ty: nil, isPublic: isPublic},
		isConst:     isConst,
		isParameter: isParameter,
	}
}

func NewTypeSymbol(name string, isPublic bool) TypeSymbol {
	return TypeSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
	}
}

func NewGenericFunctionSymbol(space *SymbolSpace, monomorphizer func(s GenericFunctionSymbol, args []semtypes.SemType, ret semtypes.SemType) SymbolRef) GenericFunctionSymbol {
	return &genericFunctionSymbol{space: space, monomorphizer: monomorphizer}
}

func (s *genericFunctionSymbol) Name() string {
	return "push"
}

func (s *genericFunctionSymbol) Type() semtypes.SemType {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) Kind() SymbolKind {
	return SymbolKindFunction
}

func (s *genericFunctionSymbol) SetType(_ semtypes.SemType) {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) IsPublic() bool {
	return true
}

func (s *genericFunctionSymbol) Signature() FunctionSignature {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) SetSignature(_ FunctionSignature) {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) Monomorphize(args []semtypes.SemType, ret semtypes.SemType) SymbolRef {
	return s.monomorphizer(s, args, ret)
}

func (s *genericFunctionSymbol) Space() *SymbolSpace {
	return s.space
}
