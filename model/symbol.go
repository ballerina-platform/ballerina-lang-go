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
	"iter"
	"sync"

	"ballerina-lang-go/semtypes"
)

type Scope interface {
	GetSymbol(name string) (SymbolRef, bool)
	GetPrefixedSymbol(prefix, name string) (SymbolRef, bool)
	AddSymbol(name string, symbol Symbol)
}

// SymbolSpaceProvider provides access to symbol spaces for block-level scopes
type SymbolSpaceProvider interface {
	MainSpace() *SymbolSpace
}

// BlockLevelScope combines Scope and SymbolSpaceProvider for block-level scopes
type BlockLevelScope interface {
	Scope
	SymbolSpaceProvider
}

// Symbol methods should never be called directly. Instead call them via the compiler context.
type Symbol interface {
	Name() string
	Type() semtypes.SemType
	Kind() SymbolKind
	SetType(semtypes.SemType)
	IsPublic() bool
	Copy() Symbol
}

// symbolTypeSetter is a private interface for updating symbol types during type resolution.
// All concrete symbol types implement this through symbolBase.
type symbolTypeSetter interface {
	SetType(semtypes.SemType)
}

type FuncSymbolFlags uint8

const (
	FuncSymbolFlagIsolated FuncSymbolFlags = 1 << iota
	FuncSymbolFlagTransactional
)

type FunctionSymbol interface {
	Symbol
	Signature() FunctionSignature
	SetSignature(FunctionSignature)
	DefaultableParams() *DefaultableParamInfo
	SetDefaultableParams(DefaultableParamInfo)
	IncludedRecordParams() *IncludedRecordParamInfo
	SetIncludedRecordParams(IncludedRecordParamInfo)
	ParamNames() []string
}

// GenericFunctionSymbol represents functions with [@typeParam] types
type GenericFunctionSymbol interface {
	FunctionSymbol
	Monomorphize(args []semtypes.SemType) SymbolRef
	Space() *SymbolSpace
}

// DependentlyTypedFunctionSymbol represents a [dependently typed function]. Actual function signature
// is determined at each call site by calling Monomorphize.
// TODO: this is very similar to [GenericFunctionSymbol]; merge both. #389
type DependentlyTypedFunctionSymbol interface {
	FunctionSymbol
	Monomorphize(ctx semtypes.Context, name string, polymorphicRef SymbolRef, argTys []semtypes.SemType) FunctionSymbol
	ParamTypes() []semtypes.SemType
	ReturnType() TypeOp
	NRequiredArgs() int
	FuncFlags() FuncSymbolFlags
	SetParamTypes(types []semtypes.SemType)
	SetReturnType(op TypeOp)
}

// MonomorphicFunctionSymbol represent a polymorphic function after monomrophizisation.
// It carries a reference back to the underlying polymorphic function so BIR can dispatch to it,
// but front end can treat this as non polymorphic.
type MonomorphicFunctionSymbol interface {
	FunctionSymbol
	PolymorphicSymbol() SymbolRef
}

type BinaryTypeOpKind uint8

const (
	TypeOpUnion BinaryTypeOpKind = iota
	TypeOpIntersection
)

// TypeOp represents a partially applied type. This is necessary to represent things like
// return types of dependently typed functions in a serializable manner.
type TypeOp interface {
	Apply(ctx semtypes.Context, args []semtypes.SemType) semtypes.SemType
}

type BinaryTypeOp struct {
	Kind BinaryTypeOpKind
	Lhs  TypeOp
	Rhs  TypeOp
}

func (binary *BinaryTypeOp) Apply(ctx semtypes.Context, args []semtypes.SemType) semtypes.SemType {
	lhs := binary.Lhs.Apply(ctx, args)
	rhs := binary.Rhs.Apply(ctx, args)
	if binary.Kind == TypeOpUnion {
		return semtypes.Union(lhs, rhs)
	}
	return semtypes.Intersect(lhs, rhs)
}

type IdentityTypeOp struct {
	Type semtypes.SemType
}

func (identity *IdentityTypeOp) Apply(_ semtypes.Context, _ []semtypes.SemType) semtypes.SemType {
	return identity.Type
}

// RefTypeOp references a typedesc parameter by position. Apply resolves to the constraint T of args[Index] (typedesc<T>).
type RefTypeOp struct {
	Index int
}

func (ref *RefTypeOp) Apply(ctx semtypes.Context, args []semtypes.SemType) semtypes.SemType {
	return semtypes.TypedescConstraint(ctx, args[ref.Index])
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
		mu          sync.RWMutex
		Pkg         PackageIdentifier
		lookupTable map[string]int
		symbols     []Symbol
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

	// memberHolderBase carries direct + type-inclusion-inherited members
	// (fields and optional rest-type for records; fields + methods for classes
	// and object type aliases).
	memberHolderBase struct {
		members []InclusionMember
	}

	ClassSymbol struct {
		TypeSymbol
		memberHolderBase
		methods map[string]SymbolRef
	}

	RecordSymbol struct {
		TypeSymbol
		memberHolderBase
	}

	ObjectTypeSymbol struct {
		TypeSymbol
		memberHolderBase
	}

	FieldDefault struct {
		FieldName string
		FnRef     SymbolRef
	}

	ValueSymbol struct {
		symbolBase
		isConst     bool
		isParameter bool
	}

	functionSymbol struct {
		symbolBase
		signature            FunctionSignature
		defaultableParams    DefaultableParamInfo
		includedRecordParams IncludedRecordParamInfo
	}

	monomorphicFunctionSymbol struct {
		functionSymbol
		polymorhpicFn SymbolRef
	}

	genericFunctionSymbol struct {
		name          string
		space         *SymbolSpace
		monomorphizer func(s GenericFunctionSymbol, args []semtypes.SemType) SymbolRef
		paramNames    []string
	}

	dependentlyTypedFunctionSymbol struct {
		symbolBase
		paramNames           []string
		nRequiredArgs        int
		Flags                FuncSymbolFlags
		defaultable          DefaultableParamInfo
		includedRecordParams IncludedRecordParamInfo

		// Populated by type resolver at stage 4.
		paramTypes []semtypes.SemType
		retType    TypeOp
	}

	FunctionSignature struct {
		ParamTypes    []semtypes.SemType
		ParamNames    []string
		ReturnType    semtypes.SemType
		RestParamType semtypes.SemType
		Flags         FuncSymbolFlags
	}

	DefaultableParamKind uint8

	DefaultableParam struct {
		Symbol SymbolRef
		Kind   DefaultableParamKind
	}

	DefaultableParamInfo struct {
		params      []DefaultableParam
		defaultable []bool
	}

	IncludedRecordParamInfo struct {
		params     []bool
		fieldNames [][]string
	}
)

const (
	DefaultableParamKindExpr DefaultableParamKind = iota
	DefaultableParamKindInferredTypedesc
)

type InclusionMemberKind uint8

const (
	InclusionMemberKindField InclusionMemberKind = iota
	InclusionMemberKindMethod
	InclusionMemberKindRemoteMethod
	InclusionMemberKindResourceMethod
	InclusionMemberKindRestType
)

type InclusionMember interface {
	MemberName() string
	MemberKind() InclusionMemberKind
	MemberType() semtypes.SemType
	SetMemberType(semtypes.SemType)
}

type FieldDescriptorFlag uint8

const (
	FieldDescriptorReadonly FieldDescriptorFlag = 1 << iota
	FieldDescriptorOptional
	FieldDescriptorHasDefault
)

type FieldDescriptor struct {
	name         string
	ty           semtypes.SemType
	flags        FieldDescriptorFlag
	DefaultFnRef SymbolRef
	visibility   Visibility
}

func NewFieldDescriptor(name string, flags FieldDescriptorFlag, visibility Visibility) FieldDescriptor {
	return FieldDescriptor{name: name, flags: flags, visibility: visibility}
}

func (f *FieldDescriptor) MemberName() string                { return f.name }
func (f *FieldDescriptor) MemberKind() InclusionMemberKind   { return InclusionMemberKindField }
func (f *FieldDescriptor) MemberType() semtypes.SemType      { return f.ty }
func (f *FieldDescriptor) SetMemberType(ty semtypes.SemType) { f.ty = ty }
func (f *FieldDescriptor) Visibility() Visibility            { return f.visibility }
func (f *FieldDescriptor) IsReadonly() bool                  { return f.flags&FieldDescriptorReadonly != 0 }
func (f *FieldDescriptor) IsOptional() bool                  { return f.flags&FieldDescriptorOptional != 0 }
func (f *FieldDescriptor) HasDefault() bool                  { return f.flags&FieldDescriptorHasDefault != 0 }

type MethodDescriptor struct {
	name       string
	kind       InclusionMemberKind
	ty         semtypes.SemType
	MethodRef  SymbolRef
	visibility Visibility
}

func NewMethodDescriptor(name string, kind InclusionMemberKind, visibility Visibility, methodRef SymbolRef) MethodDescriptor {
	return MethodDescriptor{name: name, kind: kind, visibility: visibility, MethodRef: methodRef}
}

func (m *MethodDescriptor) MemberName() string                { return m.name }
func (m *MethodDescriptor) MemberKind() InclusionMemberKind   { return m.kind }
func (m *MethodDescriptor) MemberType() semtypes.SemType      { return m.ty }
func (m *MethodDescriptor) SetMemberType(ty semtypes.SemType) { m.ty = ty }
func (m *MethodDescriptor) Visibility() Visibility            { return m.visibility }

type RestTypeDescriptor struct {
	ty semtypes.SemType
}

func NewRestTypeDescriptor() RestTypeDescriptor {
	return RestTypeDescriptor{}
}

func (r *RestTypeDescriptor) MemberName() string                { panic("RestTypeDescriptor has no name") }
func (r *RestTypeDescriptor) MemberKind() InclusionMemberKind   { return InclusionMemberKindRestType }
func (r *RestTypeDescriptor) MemberType() semtypes.SemType      { return r.ty }
func (r *RestTypeDescriptor) SetMemberType(ty semtypes.SemType) { r.ty = ty }

var (
	_ InclusionMember = &FieldDescriptor{}
	_ InclusionMember = &MethodDescriptor{}
	_ InclusionMember = &RestTypeDescriptor{}
)

var (
	_ Scope                          = &ModuleScope{}
	_ Scope                          = &FunctionScope{}
	_ Scope                          = &BlockScope{}
	_ Symbol                         = &TypeSymbol{}
	_ Symbol                         = &ClassSymbol{}
	_ Symbol                         = &RecordSymbol{}
	_ Symbol                         = &ObjectTypeSymbol{}
	_ MemberCarrier                  = &ClassSymbol{}
	_ MemberCarrier                  = &RecordSymbol{}
	_ MemberCarrier                  = &ObjectTypeSymbol{}
	_ Symbol                         = &ValueSymbol{}
	_ Symbol                         = &functionSymbol{}
	_ FunctionSymbol                 = &functionSymbol{}
	_ GenericFunctionSymbol          = &genericFunctionSymbol{}
	_ DependentlyTypedFunctionSymbol = &dependentlyTypedFunctionSymbol{}
	_ MonomorphicFunctionSymbol      = &monomorphicFunctionSymbol{}
	_ Symbol                         = &SymbolRef{}
	_ SymbolSpaceProvider            = &ModuleScope{}
)

func (space *SymbolSpace) AddSymbol(name string, symbol Symbol) {
	if _, ok := symbol.(*SymbolRef); ok {
		panic("SymbolRef cannot be added to a SymbolSpace")
	}
	space.mu.Lock()
	space.lookupTable[name] = len(space.symbols)
	space.symbols = append(space.symbols, symbol)
	space.mu.Unlock()
}

func (space *SymbolSpace) GetSymbol(name string) (SymbolRef, bool) {
	space.mu.RLock()
	index, ok := space.lookupTable[name]
	space.mu.RUnlock()
	if !ok {
		return SymbolRef{}, false
	}
	return SymbolRef{Package: space.Pkg, Index: index, SpaceIndex: space.index}, true
}

// AppendSymbol appends a symbol to the space and returns its index. Thread-safe.
func (space *SymbolSpace) AppendSymbol(symbol Symbol) int {
	// We really need this lock only for module level symbols but we don't distinguish between module level space and other spaces
	space.mu.Lock()
	defer space.mu.Unlock()
	index := len(space.symbols)
	space.symbols = append(space.symbols, symbol)
	return index
}

// RefAt returns a SymbolRef for the symbol at the given index.
func (space *SymbolSpace) RefAt(index int) SymbolRef {
	return SymbolRef{Package: space.Pkg, Index: index, SpaceIndex: space.index}
}

// SymbolAt returns the symbol at the given index. Thread-safe.
func (space *SymbolSpace) SpaceIndex() int {
	return space.index
}

func (space *SymbolSpace) SymbolAt(index int) Symbol {
	space.mu.RLock()
	defer space.mu.RUnlock()
	return space.symbols[index]
}

func (space *SymbolSpace) Len() int {
	space.mu.RLock()
	defer space.mu.RUnlock()
	return len(space.symbols)
}

// Symbols returns an iterator over the symbols in the space. This is for
// reading only — callers must not modify the yielded symbols or add new symbols
// to the space during iteration.
func (space *SymbolSpace) Symbols() iter.Seq2[int, Symbol] {
	return func(yield func(int, Symbol) bool) {
		space.mu.RLock()
		defer space.mu.RUnlock()
		for i, sym := range space.symbols {
			if !yield(i, sym) {
				return
			}
		}
	}
}

func NewSymbolSpaceInner(packageID PackageID, index int) *SymbolSpace {
	pkg := PackageIdentifier{
		Organization: packageID.OrgName.Value(),
		Package:      packageID.PkgName.Value(),
		Version:      packageID.Version.Value(),
	}
	return &SymbolSpace{index: index, Pkg: pkg, lookupTable: make(map[string]int), symbols: make([]Symbol, 0)}
}

func (ms *ModuleScope) Exports() ExportedSymbolSpace {
	return NewExportedSymbolSpace(ms.Main, ms.Annotation)
}

func (ms *ModuleScope) GetSymbol(name string) (SymbolRef, bool) {
	return ms.Main.GetSymbol(name)
}

func (ms *ModuleScope) MainSpace() *SymbolSpace {
	return ms.Main
}

func mapToLangPrefixIfNeeded(prefix string) string {
	switch prefix {
	case "int":
		return "lang.int"
	case "array":
		return "lang.array"
	default:
		return prefix
	}
}

func (ms *ModuleScope) GetPrefixedSymbol(prefix, name string) (SymbolRef, bool) {
	if prefix == "" {
		return ms.Main.GetSymbol(name)
	}
	exported, ok := ms.Prefix[prefix]
	if !ok {
		exported, ok = ms.Prefix[mapToLangPrefixIfNeeded(prefix)]
		if !ok {
			return SymbolRef{}, false
		}
	}
	return exported.GetSymbol(name)
}

func (ms *ModuleScope) AddSymbol(name string, symbol Symbol) {
	ms.Main.AddSymbol(name, symbol)
}

func (ms *ModuleScope) AddAnnotationSymbol(name string, symbol Symbol) {
	ms.Annotation.AddSymbol(name, symbol)
}

func NewExportedSymbolSpace(main, annotation *SymbolSpace) ExportedSymbolSpace {
	return ExportedSymbolSpace{Main: main, Annotation: annotation}
}

func (space *ExportedSymbolSpace) PublicMainSymbols() iter.Seq2[SymbolRef, Symbol] {
	return func(yield func(SymbolRef, Symbol) bool) {
		for i, sym := range space.Main.Symbols() {
			if !sym.IsPublic() {
				continue
			}
			if !yield(space.Main.RefAt(i), sym) {
				return
			}
		}
	}
}

func (space *ExportedSymbolSpace) GetSymbol(name string) (SymbolRef, bool) {
	ref, ok := space.Main.GetSymbol(name)
	if !ok {
		return SymbolRef{}, false
	}
	sym := space.Main.SymbolAt(ref.Index)
	if !sym.IsPublic() {
		return SymbolRef{}, false
	}
	return ref, true
}

func (bs *BlockScopeBase) GetSymbol(name string) (SymbolRef, bool) {
	ref, ok := bs.Main.GetSymbol(name)
	if ok {
		return ref, true
	}
	return bs.Parent.GetSymbol(name)
}

func (bs *BlockScopeBase) GetPrefixedSymbol(prefix, name string) (SymbolRef, bool) {
	return bs.Parent.GetPrefixedSymbol(prefix, name)
}

func (bs *BlockScopeBase) AddSymbol(name string, symbol Symbol) {
	bs.Main.AddSymbol(name, symbol)
}

func (bs *BlockScopeBase) MainSpace() *SymbolSpace {
	return bs.Main
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

func (ref *SymbolRef) Copy() Symbol {
	panic("SymbolRef can't be copied")
}

func (ts *TypeSymbol) Kind() SymbolKind {
	return SymbolKindType
}

func (ts *TypeSymbol) Copy() Symbol {
	panic("TypeSymbol cannot be copied")
}

// MemberCarrier is implemented by symbols that carry direct + inclusion-inherited members.
// TypeSymbol does not implement this; only RecordSymbol, ClassSymbol, and ObjectTypeSymbol do.
type MemberCarrier interface {
	Members() []InclusionMember
	AddMember(InclusionMember)
	FieldDefaults() []FieldDefault
}

func (m *memberHolderBase) Members() []InclusionMember { return m.members }
func (m *memberHolderBase) AddMember(im InclusionMember) {
	m.members = append(m.members, im)
}

func (m *memberHolderBase) FieldDefaults() []FieldDefault {
	var defaults []FieldDefault
	for _, im := range m.members {
		if fd, ok := im.(*FieldDescriptor); ok && fd.DefaultFnRef != (SymbolRef{}) {
			defaults = append(defaults, FieldDefault{FieldName: fd.name, FnRef: fd.DefaultFnRef})
		}
	}
	return defaults
}

func (r *RecordSymbol) Fields() iter.Seq2[string, *FieldDescriptor] {
	return func(yield func(string, *FieldDescriptor) bool) {
		for _, m := range r.members {
			fd, ok := m.(*FieldDescriptor)
			if !ok {
				continue
			}
			if !yield(fd.name, fd) {
				return
			}
		}
	}
}

func (r *RecordSymbol) Field(name string) (*FieldDescriptor, bool) {
	for _, m := range r.members {
		if fd, ok := m.(*FieldDescriptor); ok && fd.name == name {
			return fd, true
		}
	}
	return nil, false
}

func (r *RecordSymbol) RestField() (*RestTypeDescriptor, bool) {
	for _, m := range r.members {
		if rd, ok := m.(*RestTypeDescriptor); ok {
			return rd, true
		}
	}
	return nil, false
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

func (vs *ValueSymbol) IsConst() bool {
	return vs.isConst || vs.isParameter
}

func (vs *ValueSymbol) Copy() Symbol {
	cp := *vs
	return &cp
}

func (fs *functionSymbol) Kind() SymbolKind {
	return SymbolKindFunction
}

func (fs *functionSymbol) Copy() Symbol {
	cp := *fs
	return &cp
}

func (fs *functionSymbol) Signature() FunctionSignature {
	return fs.signature
}

func (fs *functionSymbol) SetSignature(sig FunctionSignature) {
	fs.signature = sig
}

func (fs *functionSymbol) DefaultableParams() *DefaultableParamInfo {
	return &fs.defaultableParams
}

func (fs *functionSymbol) SetDefaultableParams(info DefaultableParamInfo) {
	fs.defaultableParams = info
}

func (fs *functionSymbol) IncludedRecordParams() *IncludedRecordParamInfo {
	return &fs.includedRecordParams
}

func (fs *functionSymbol) SetIncludedRecordParams(info IncludedRecordParamInfo) {
	fs.includedRecordParams = info
}

func (fs *functionSymbol) ParamNames() []string {
	return fs.Signature().ParamNames
}

func NewFunctionSymbol(name string, signature FunctionSignature, isPublic bool) FunctionSymbol {
	return &functionSymbol{
		symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		signature:  signature,
	}
}

func NewDefaultableParamInfo(paramCount int) DefaultableParamInfo {
	return DefaultableParamInfo{
		params:      make([]DefaultableParam, paramCount),
		defaultable: make([]bool, paramCount),
	}
}

func (d *DefaultableParamInfo) Get(index int) (DefaultableParam, bool) {
	if index >= len(d.defaultable) || !d.defaultable[index] {
		return DefaultableParam{}, false
	}
	return d.params[index], true
}

func (d *DefaultableParamInfo) SetDefaultable(index int, symbol SymbolRef) {
	d.defaultable[index] = true
	d.params[index] = DefaultableParam{Symbol: symbol, Kind: DefaultableParamKindExpr}
}

func (d *DefaultableParamInfo) SetInferredTypedesc(index int) {
	d.defaultable[index] = true
	d.params[index] = DefaultableParam{Kind: DefaultableParamKindInferredTypedesc}
}

func NewIncludedRecordParamInfo(paramCount int) IncludedRecordParamInfo {
	return IncludedRecordParamInfo{
		params:     make([]bool, paramCount),
		fieldNames: make([][]string, paramCount),
	}
}

func (i *IncludedRecordParamInfo) Set(index int) {
	i.params[index] = true
}

func (i *IncludedRecordParamInfo) SetFields(index int, names []string) {
	i.fieldNames[index] = names
}

func (i *IncludedRecordParamInfo) Fields(index int) []string {
	if index >= len(i.fieldNames) {
		return nil
	}
	return i.fieldNames[index]
}

func (i *IncludedRecordParamInfo) LookupField(name string) (int, bool) {
	for idx, names := range i.fieldNames {
		for _, n := range names {
			if n == name {
				return idx, true
			}
		}
	}
	return -1, false
}

func (i *IncludedRecordParamInfo) IsIncluded(index int) bool {
	if index >= len(i.params) {
		return false
	}
	return i.params[index]
}

func (i *IncludedRecordParamInfo) Len() int {
	return len(i.params)
}

func (fs *FunctionSignature) IsIsolated() bool {
	return fs.Flags&FuncSymbolFlagIsolated != 0
}

func (fs *FunctionSignature) IsTransactional() bool {
	return fs.Flags&FuncSymbolFlagTransactional != 0
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

func NewClassSymbol(name string, isPublic bool) ClassSymbol {
	return ClassSymbol{
		TypeSymbol: TypeSymbol{
			symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		},
	}
}

func NewRecordSymbol(name string, isPublic bool) RecordSymbol {
	return RecordSymbol{
		TypeSymbol: TypeSymbol{
			symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		},
	}
}

func NewObjectTypeSymbol(name string, isPublic bool) ObjectTypeSymbol {
	return ObjectTypeSymbol{
		TypeSymbol: TypeSymbol{
			symbolBase: symbolBase{name: name, ty: nil, isPublic: isPublic},
		},
	}
}

func (c *ClassSymbol) SetMethods(methods map[string]SymbolRef) {
	c.methods = methods
}

func (c *ClassSymbol) MethodSymbol(name string) (SymbolRef, bool) {
	ref, ok := c.methods[name]
	return ref, ok
}

func NewGenericFunctionSymbol(name string, space *SymbolSpace, paramNames []string, monomorphizer func(s GenericFunctionSymbol, args []semtypes.SemType) SymbolRef) GenericFunctionSymbol {
	return &genericFunctionSymbol{name: name, space: space, paramNames: paramNames, monomorphizer: monomorphizer}
}

func (s *genericFunctionSymbol) Name() string {
	return s.name
}

func (s *genericFunctionSymbol) ParamNames() []string {
	return s.paramNames
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

func (s *genericFunctionSymbol) DefaultableParams() *DefaultableParamInfo {
	return &DefaultableParamInfo{}
}

func (s *genericFunctionSymbol) SetDefaultableParams(_ DefaultableParamInfo) {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) IncludedRecordParams() *IncludedRecordParamInfo {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) SetIncludedRecordParams(_ IncludedRecordParamInfo) {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) Copy() Symbol {
	panic("GenericSymbol must be Monomorphized")
}

func (s *genericFunctionSymbol) Monomorphize(args []semtypes.SemType) SymbolRef {
	return s.monomorphizer(s, args)
}

func (s *genericFunctionSymbol) Space() *SymbolSpace {
	return s.space
}

func NewDependentlyTypedFunctionSymbol(name string, paramNames []string, nRequiredArgs int, flags FuncSymbolFlags, isPublic bool) DependentlyTypedFunctionSymbol {
	return &dependentlyTypedFunctionSymbol{
		symbolBase:    symbolBase{name: name, ty: nil, isPublic: isPublic},
		paramNames:    paramNames,
		nRequiredArgs: nRequiredArgs,
		Flags:         flags,
	}
}

func (s *dependentlyTypedFunctionSymbol) Kind() SymbolKind { return SymbolKindFunction }

func (s *dependentlyTypedFunctionSymbol) Type() semtypes.SemType {
	panic("DependentlyTypedFunctionSymbol must be Monomorphized")
}

func (s *dependentlyTypedFunctionSymbol) SetType(_ semtypes.SemType) {
	panic("DependentlyTypedFunctionSymbol must be Monomorphized")
}

func (s *dependentlyTypedFunctionSymbol) Signature() FunctionSignature {
	panic("DependentlyTypedFunctionSymbol must be Monomorphized")
}

func (s *dependentlyTypedFunctionSymbol) SetSignature(_ FunctionSignature) {
	panic("DependentlyTypedFunctionSymbol must be Monomorphized")
}

func (s *dependentlyTypedFunctionSymbol) DefaultableParams() *DefaultableParamInfo {
	return &s.defaultable
}

func (s *dependentlyTypedFunctionSymbol) SetDefaultableParams(info DefaultableParamInfo) {
	s.defaultable = info
}

func (s *dependentlyTypedFunctionSymbol) IncludedRecordParams() *IncludedRecordParamInfo {
	return &s.includedRecordParams
}

func (s *dependentlyTypedFunctionSymbol) SetIncludedRecordParams(info IncludedRecordParamInfo) {
	s.includedRecordParams = info
}

func (s *dependentlyTypedFunctionSymbol) ParamNames() []string { return s.paramNames }

func (s *dependentlyTypedFunctionSymbol) Copy() Symbol {
	cp := *s
	return &cp
}

func (s *dependentlyTypedFunctionSymbol) ParamTypes() []semtypes.SemType { return s.paramTypes }
func (s *dependentlyTypedFunctionSymbol) ReturnType() TypeOp             { return s.retType }
func (s *dependentlyTypedFunctionSymbol) NRequiredArgs() int             { return s.nRequiredArgs }
func (s *dependentlyTypedFunctionSymbol) FuncFlags() FuncSymbolFlags     { return s.Flags }

func (s *dependentlyTypedFunctionSymbol) SetParamTypes(types []semtypes.SemType) {
	s.paramTypes = types
}

func (s *dependentlyTypedFunctionSymbol) SetReturnType(op TypeOp) {
	s.retType = op
}

func (s *dependentlyTypedFunctionSymbol) Monomorphize(ctx semtypes.Context, name string, origRef SymbolRef, argTys []semtypes.SemType) FunctionSymbol {
	var rest semtypes.SemType = semtypes.NEVER
	if len(argTys) > s.nRequiredArgs {
		for _, each := range argTys[s.nRequiredArgs:] {
			rest = semtypes.Union(rest, each)
		}
	}
	returnType := s.retType.Apply(ctx, argTys)
	sig := FunctionSignature{
		ParamTypes:    argTys,
		ParamNames:    s.paramNames,
		RestParamType: rest,
		ReturnType:    returnType,
		Flags:         s.Flags,
	}
	return &monomorphicFunctionSymbol{
		functionSymbol: functionSymbol{
			symbolBase:           symbolBase{name: name, ty: nil, isPublic: s.isPublic},
			signature:            sig,
			defaultableParams:    s.defaultable,
			includedRecordParams: s.includedRecordParams,
		},
		polymorhpicFn: origRef,
	}
}

func (m *monomorphicFunctionSymbol) PolymorphicSymbol() SymbolRef { return m.polymorhpicFn }

func (m *monomorphicFunctionSymbol) Copy() Symbol {
	cp := *m
	return &cp
}
