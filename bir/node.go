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

package bir

import (
	"strconv"
	"strings"

	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

type (
	Name                  = model.Name
	PackageID             = model.PackageID
	SymbolOrigin          = model.SymbolOrigin
	AttachPoint           = model.AttachPoint
	MarkdownDocAttachment = model.MarkdownDocAttachment
	Flag                  = model.Flag
	NamedNode             = model.NamedNode

	BType = model.ValueType
)

type BIRNodeData interface {
	SetPos(pos diagnostics.Location)
	GetPos() diagnostics.Location
}

type BIRNodeBase struct {
	Pos diagnostics.Location
}

func (b *BIRNodeBase) SetPos(pos diagnostics.Location) {
	b.Pos = pos
}

func (b *BIRNodeBase) GetPos() diagnostics.Location {
	return b.Pos
}

type BIRNode interface {
	BIRNodeData
	Accept(visitor BIRVisitor)
}

type BIRPackageData interface {
	BIRNodeData
	SetPackageID(packageID PackageID)
	GetPackageID() PackageID
	SetImportModules(importModules *[]BIRImportModule)
	GetImportModules() *[]BIRImportModule
	SetTypeDefs(typeDefs *[]BIRTypeDefinition)
	GetTypeDefs() *[]BIRTypeDefinition
	SetGlobalVars(globalVars *[]BIRGlobalVariableDcl)
	GetGlobalVars() *[]BIRGlobalVariableDcl
	SetImportedGlobalVarsDummyVarDcls(importedGlobalVarsDummyVarDcls *[]BIRGlobalVariableDcl)
	GetImportedGlobalVarsDummyVarDcls() *[]BIRGlobalVariableDcl
	SetFunctions(functions *[]BIRFunction)
	GetFunctions() *[]BIRFunction
	SetAnnotations(annotations *[]BIRAnnotation)
	GetAnnotations() *[]BIRAnnotation
	SetConstants(constants *[]BIRConstant)
	GetConstants() *[]BIRConstant
	SetServiceDecls(serviceDecls *[]BIRServiceDeclaration)
	GetServiceDecls() *[]BIRServiceDeclaration
	SetIsListenerAvailable(isListenerAvailable bool)
	GetIsListenerAvailable() bool
	SetRecordDefaultValueMap(recordDefaultValueMap map[string]map[string]string)
	GetRecordDefaultValueMap() map[string]map[string]string
}

type BIRPackageBase struct {
	BIRNodeBase
	PackageID                      PackageID
	ImportModules                  []BIRImportModule
	TypeDefs                       []BIRTypeDefinition
	GlobalVars                     []BIRGlobalVariableDcl
	ImportedGlobalVarsDummyVarDcls []BIRGlobalVariableDcl
	Functions                      []BIRFunction
	Annotations                    []BIRAnnotation
	Constants                      []BIRConstant
	ServiceDecls                   []BIRServiceDeclaration
	IsListenerAvailable            bool
	RecordDefaultValueMap          map[string]map[string]string
}

func (b *BIRPackageBase) SetPackageID(packageID PackageID) {
	b.PackageID = packageID
}

func (b *BIRPackageBase) GetPackageID() PackageID {
	return b.PackageID
}

func (b *BIRPackageBase) SetImportModules(importModules *[]BIRImportModule) {
	b.ImportModules = *importModules
}

func (b *BIRPackageBase) GetImportModules() *[]BIRImportModule {
	return &b.ImportModules
}

func (b *BIRPackageBase) SetTypeDefs(typeDefs *[]BIRTypeDefinition) {
	b.TypeDefs = *typeDefs
}

func (b *BIRPackageBase) GetTypeDefs() *[]BIRTypeDefinition {
	return &b.TypeDefs
}

func (b *BIRPackageBase) SetGlobalVars(globalVars *[]BIRGlobalVariableDcl) {
	b.GlobalVars = *globalVars
}

func (b *BIRPackageBase) GetGlobalVars() *[]BIRGlobalVariableDcl {
	return &b.GlobalVars
}

func (b *BIRPackageBase) SetImportedGlobalVarsDummyVarDcls(importedGlobalVarsDummyVarDcls *[]BIRGlobalVariableDcl) {
	b.ImportedGlobalVarsDummyVarDcls = *importedGlobalVarsDummyVarDcls
}

func (b *BIRPackageBase) GetImportedGlobalVarsDummyVarDcls() *[]BIRGlobalVariableDcl {
	return &b.ImportedGlobalVarsDummyVarDcls
}

func (b *BIRPackageBase) SetFunctions(functions *[]BIRFunction) {
	b.Functions = *functions
}

func (b *BIRPackageBase) GetFunctions() *[]BIRFunction {
	return &b.Functions
}

func (b *BIRPackageBase) SetAnnotations(annotations *[]BIRAnnotation) {
	b.Annotations = *annotations
}

func (b *BIRPackageBase) GetAnnotations() *[]BIRAnnotation {
	return &b.Annotations
}

func (b *BIRPackageBase) SetConstants(constants *[]BIRConstant) {
	b.Constants = *constants
}

func (b *BIRPackageBase) GetConstants() *[]BIRConstant {
	return &b.Constants
}

func (b *BIRPackageBase) SetServiceDecls(serviceDecls *[]BIRServiceDeclaration) {
	b.ServiceDecls = *serviceDecls
}

func (b *BIRPackageBase) GetServiceDecls() *[]BIRServiceDeclaration {
	return &b.ServiceDecls
}

func (b *BIRPackageBase) SetIsListenerAvailable(isListenerAvailable bool) {
	b.IsListenerAvailable = isListenerAvailable
}

func (b *BIRPackageBase) GetIsListenerAvailable() bool {
	return b.IsListenerAvailable
}

func (b *BIRPackageBase) SetRecordDefaultValueMap(recordDefaultValueMap map[string]map[string]string) {
	b.RecordDefaultValueMap = recordDefaultValueMap
}

func (b *BIRPackageBase) GetRecordDefaultValueMap() map[string]map[string]string {
	return b.RecordDefaultValueMap
}

type BIRPackage interface {
	BIRPackageData
	BIRNode
}

type BIRPackageMethods struct {
	Self BIRPackage
}

func (m *BIRPackageMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRPackage(m.Self)
}

func NewBIRPackage(pos diagnostics.Location, org Name, pkgName Name, name Name, version Name, sourceFileName Name, sourceRoot string, skipTest bool) BIRPackage {
	return NewBIRPackageWithIsTestPkg(pos, org, pkgName, name, version, sourceFileName, sourceRoot, skipTest, false)
}

func NewBIRPackageWithIsTestPkg(pos diagnostics.Location, org Name, pkgName Name, name Name, version Name, sourceFileName Name, sourceRoot string, skipTest bool, isTestPkg bool) BIRPackage {
	pkg := &BIRPackageImpl{
		BIRPackageBase: BIRPackageBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			PackageID:                      model.NewPackageIDFull(org, pkgName, name, version, sourceFileName, sourceRoot, isTestPkg, skipTest),
			ImportModules:                  []BIRImportModule{},
			TypeDefs:                       []BIRTypeDefinition{},
			GlobalVars:                     []BIRGlobalVariableDcl{},
			ImportedGlobalVarsDummyVarDcls: []BIRGlobalVariableDcl{},
			Functions:                      []BIRFunction{},
			Annotations:                    []BIRAnnotation{},
			Constants:                      []BIRConstant{},
			ServiceDecls:                   []BIRServiceDeclaration{},
			RecordDefaultValueMap:          make(map[string]map[string]string),
		},
		BIRPackageMethods: BIRPackageMethods{},
	}
	pkg.BIRPackageMethods.Self = pkg
	return pkg
}

type BIRPackageImpl struct {
	BIRPackageBase
	BIRPackageMethods
}

type BIRImportModuleData interface {
	BIRNodeData
	SetPackageID(packageID PackageID)
	GetPackageID() PackageID
}

type BIRImportModuleBase struct {
	BIRNodeBase
	PackageID PackageID
}

func (b *BIRImportModuleBase) SetPackageID(packageID PackageID) {
	b.PackageID = packageID
}

func (b *BIRImportModuleBase) GetPackageID() PackageID {
	return b.PackageID
}

type BIRImportModule interface {
	BIRImportModuleData
	BIRNode
}

type BIRImportModuleMethods struct {
	Self BIRImportModule
}

func (m *BIRImportModuleMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRImportModule(m.Self)
}

func NewBIRImportModule(pos diagnostics.Location, org Name, name Name, version Name) BIRImportModule {
	mod := &BIRImportModuleImpl{
		BIRImportModuleBase: BIRImportModuleBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			PackageID: model.NewPackageIDWithName(org, name, version),
		},
		BIRImportModuleMethods: BIRImportModuleMethods{},
	}
	mod.BIRImportModuleMethods.Self = mod
	return mod
}

type BIRImportModuleImpl struct {
	BIRImportModuleBase
	BIRImportModuleMethods
}

type BIRVariableDclData interface {
	BIRDocumentableNodeData
	SetType(type_ BType)
	GetType() BType
	SetName(name Name)
	GetName() Name
	SetOriginalName(originalName Name)
	GetOriginalName() Name
	SetMetaVarName(metaVarName string)
	GetMetaVarName() string
	SetJvmVarName(jvmVarName string)
	GetJvmVarName() string
	SetKind(kind VarKind)
	GetKind() VarKind
	SetScope(scope VarScope)
	GetScope() VarScope
	SetIgnoreVariable(ignoreVariable bool)
	GetIgnoreVariable() bool
	SetEndBB(endBB BIRBasicBlock)
	GetEndBB() BIRBasicBlock
	SetStartBB(startBB BIRBasicBlock)
	GetStartBB() BIRBasicBlock
	SetInsOffset(insOffset int)
	GetInsOffset() int
	SetOnlyUsedInSingleBB(onlyUsedInSingleBB bool)
	GetOnlyUsedInSingleBB() bool
	SetInitialized(initialized bool)
	GetInitialized() bool
	SetInsScope(insScope BIRScope)
	GetInsScope() BIRScope
}

type BIRVariableDclBase struct {
	BIRDocumentableNodeBase
	Type               BType
	Name               Name
	OriginalName       Name
	MetaVarName        string
	JvmVarName         string
	Kind               VarKind
	Scope              VarScope
	IgnoreVariable     bool
	EndBB              BIRBasicBlock
	StartBB            BIRBasicBlock
	InsOffset          int
	OnlyUsedInSingleBB bool
	Initialized        bool
	InsScope           BIRScope
}

func (b *BIRVariableDclBase) SetType(type_ BType) {
	b.Type = type_
}

func (b *BIRVariableDclBase) GetType() BType {
	return b.Type
}

func (b *BIRVariableDclBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRVariableDclBase) GetName() Name {
	return b.Name
}

func (b *BIRVariableDclBase) SetOriginalName(originalName Name) {
	b.OriginalName = originalName
}

func (b *BIRVariableDclBase) GetOriginalName() Name {
	return b.OriginalName
}

func (b *BIRVariableDclBase) SetMetaVarName(metaVarName string) {
	b.MetaVarName = metaVarName
}

func (b *BIRVariableDclBase) GetMetaVarName() string {
	return b.MetaVarName
}

func (b *BIRVariableDclBase) SetJvmVarName(jvmVarName string) {
	b.JvmVarName = jvmVarName
}

func (b *BIRVariableDclBase) GetJvmVarName() string {
	return b.JvmVarName
}

func (b *BIRVariableDclBase) SetKind(kind VarKind) {
	b.Kind = kind
}

func (b *BIRVariableDclBase) GetKind() VarKind {
	return b.Kind
}

func (b *BIRVariableDclBase) SetScope(scope VarScope) {
	b.Scope = scope
}

func (b *BIRVariableDclBase) GetScope() VarScope {
	return b.Scope
}

func (b *BIRVariableDclBase) SetIgnoreVariable(ignoreVariable bool) {
	b.IgnoreVariable = ignoreVariable
}

func (b *BIRVariableDclBase) GetIgnoreVariable() bool {
	return b.IgnoreVariable
}

func (b *BIRVariableDclBase) SetEndBB(endBB BIRBasicBlock) {
	b.EndBB = endBB
}

func (b *BIRVariableDclBase) GetEndBB() BIRBasicBlock {
	return b.EndBB
}

func (b *BIRVariableDclBase) SetStartBB(startBB BIRBasicBlock) {
	b.StartBB = startBB
}

func (b *BIRVariableDclBase) GetStartBB() BIRBasicBlock {
	return b.StartBB
}

func (b *BIRVariableDclBase) SetInsOffset(insOffset int) {
	b.InsOffset = insOffset
}

func (b *BIRVariableDclBase) GetInsOffset() int {
	return b.InsOffset
}

func (b *BIRVariableDclBase) SetOnlyUsedInSingleBB(onlyUsedInSingleBB bool) {
	b.OnlyUsedInSingleBB = onlyUsedInSingleBB
}

func (b *BIRVariableDclBase) GetOnlyUsedInSingleBB() bool {
	return b.OnlyUsedInSingleBB
}

func (b *BIRVariableDclBase) SetInitialized(initialized bool) {
	b.Initialized = initialized
}

func (b *BIRVariableDclBase) GetInitialized() bool {
	return b.Initialized
}

func (b *BIRVariableDclBase) SetInsScope(insScope BIRScope) {
	b.InsScope = insScope
}

func (b *BIRVariableDclBase) GetInsScope() BIRScope {
	return b.InsScope
}

type BIRVariableDcl interface {
	BIRVariableDclData
	BIRDocumentableNode
}

type BIRVariableDclMethods struct {
	Self BIRVariableDcl
}

func (m *BIRVariableDclMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRVariableDcl(m.Self)
}

func NewBIRVariableDcl(pos diagnostics.Location, type_ BType, name Name, originalName Name, scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	varDecl := &BIRVariableDclImpl{
		BIRVariableDclBase: BIRVariableDclBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Type:         type_,
			Name:         name,
			OriginalName: originalName,
			Scope:        scope,
			Kind:         kind,
			MetaVarName:  metaVarName,
			JvmVarName:   strings.ReplaceAll(name.Value(), "%", "_"),
		},
		BIRVariableDclMethods: BIRVariableDclMethods{},
	}
	varDecl.BIRVariableDclMethods.Self = varDecl
	return varDecl
}

func NewBIRVariableDclWithName(pos diagnostics.Location, type_ BType, name Name, scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	return NewBIRVariableDcl(pos, type_, name, name, scope, kind, metaVarName)
}

func NewBIRVariableDclSimple(type_ BType, name Name, scope VarScope, kind VarKind) BIRVariableDcl {
	return NewBIRVariableDcl(nil, type_, name, name, scope, kind, "")
}

type BIRVariableDclImpl struct {
	BIRVariableDclBase
	BIRVariableDclMethods
}

type BIRParameterData interface {
	BIRNodeData
	SetName(name Name)
	GetName() Name
	SetFlags(flags int64)
	GetFlags() int64
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
}

type BIRParameterBase struct {
	BIRNodeBase
	Name             Name
	Flags            int64
	AnnotAttachments []BIRAnnotationAttachment
}

func (b *BIRParameterBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRParameterBase) GetName() Name {
	return b.Name
}

func (b *BIRParameterBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRParameterBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRParameterBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRParameterBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

type BIRParameter interface {
	BIRParameterData
	BIRNode
}

type BIRParameterMethods struct {
	Self BIRParameter
}

func (m *BIRParameterMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRParameter(m.Self)
}

func NewBIRParameter(pos diagnostics.Location, name Name, flags int64) BIRParameter {
	param := &BIRParameterImpl{
		BIRParameterBase: BIRParameterBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			Name:             name,
			Flags:            flags,
			AnnotAttachments: []BIRAnnotationAttachment{},
		},
		BIRParameterMethods: BIRParameterMethods{},
	}
	param.BIRParameterMethods.Self = param
	return param
}

type BIRParameterImpl struct {
	BIRParameterBase
	BIRParameterMethods
}

type BIRGlobalVariableDclData interface {
	BIRVariableDclData
	SetFlags(flags int64)
	GetFlags() int64
	SetPkgId(pkgId PackageID)
	GetPkgId() PackageID
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
}

type BIRGlobalVariableDclBase struct {
	BIRVariableDclBase
	Flags            int64
	PkgId            PackageID
	Origin           SymbolOrigin
	AnnotAttachments []BIRAnnotationAttachment
}

func (b *BIRGlobalVariableDclBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRGlobalVariableDclBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRGlobalVariableDclBase) SetPkgId(pkgId PackageID) {
	b.PkgId = pkgId
}

func (b *BIRGlobalVariableDclBase) GetPkgId() PackageID {
	return b.PkgId
}

func (b *BIRGlobalVariableDclBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRGlobalVariableDclBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRGlobalVariableDclBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRGlobalVariableDclBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

type BIRGlobalVariableDcl interface {
	BIRGlobalVariableDclData
	BIRVariableDcl
}

type BIRGlobalVariableDclMethods struct {
	Self BIRGlobalVariableDcl
}

func (m *BIRGlobalVariableDclMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRGlobalVariableDcl(m.Self)
}

func NewBIRGlobalVariableDcl(pos diagnostics.Location, flags int64, type_ BType, pkgId PackageID, name Name, originalName Name, scope VarScope, kind VarKind, metaVarName string, origin SymbolOrigin) BIRGlobalVariableDcl {
	globalVar := &BIRGlobalVariableDclImpl{
		BIRGlobalVariableDclBase: BIRGlobalVariableDclBase{
			BIRVariableDclBase: BIRVariableDclBase{
				BIRDocumentableNodeBase: BIRDocumentableNodeBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
				},
				Type:         type_,
				Name:         name,
				OriginalName: originalName,
				Scope:        scope,
				Kind:         kind,
				MetaVarName:  metaVarName,
			},
			Flags:            flags,
			PkgId:            pkgId,
			Origin:           origin,
			AnnotAttachments: []BIRAnnotationAttachment{},
		},
		BIRGlobalVariableDclMethods: BIRGlobalVariableDclMethods{},
	}
	globalVar.BIRGlobalVariableDclMethods.Self = globalVar
	return globalVar
}

type BIRGlobalVariableDclImpl struct {
	BIRGlobalVariableDclBase
	BIRGlobalVariableDclMethods
}

func (v *BIRGlobalVariableDclImpl) String() string {
	return string(v.GetName())
}

type BIRFunctionParameterData interface {
	BIRVariableDclData
	SetHasDefaultExpr(hasDefaultExpr bool)
	GetHasDefaultExpr() bool
	SetIsPathParameter(isPathParameter bool)
	GetIsPathParameter() bool
}

type BIRFunctionParameterBase struct {
	BIRVariableDclBase
	HasDefaultExpr  bool
	IsPathParameter bool
}

func (b *BIRFunctionParameterBase) SetHasDefaultExpr(hasDefaultExpr bool) {
	b.HasDefaultExpr = hasDefaultExpr
}

func (b *BIRFunctionParameterBase) GetHasDefaultExpr() bool {
	return b.HasDefaultExpr
}

func (b *BIRFunctionParameterBase) SetIsPathParameter(isPathParameter bool) {
	b.IsPathParameter = isPathParameter
}

func (b *BIRFunctionParameterBase) GetIsPathParameter() bool {
	return b.IsPathParameter
}

type BIRFunctionParameter interface {
	BIRFunctionParameterData
	BIRVariableDcl
}

type BIRFunctionParameterMethods struct {
	Self BIRFunctionParameter
}

func (m *BIRFunctionParameterMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRFunctionParameter(m.Self)
}

func NewBIRFunctionParameter(pos diagnostics.Location, type_ BType, name Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr bool) BIRFunctionParameter {
	param := &BIRFunctionParameterImpl{
		BIRFunctionParameterBase: BIRFunctionParameterBase{
			BIRVariableDclBase: BIRVariableDclBase{
				BIRDocumentableNodeBase: BIRDocumentableNodeBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
				},
				Type:        type_,
				Name:        name,
				Scope:       scope,
				Kind:        kind,
				MetaVarName: metaVarName,
			},
			HasDefaultExpr: hasDefaultExpr,
		},
		BIRFunctionParameterMethods: BIRFunctionParameterMethods{},
	}
	param.BIRFunctionParameterMethods.Self = param
	return param
}

func NewBIRFunctionParameterWithIsPathParameter(pos diagnostics.Location, type_ BType, name Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr bool, isPathParameter bool) BIRFunctionParameter {
	param := NewBIRFunctionParameter(pos, type_, name, scope, kind, metaVarName, hasDefaultExpr)
	param.SetIsPathParameter(isPathParameter)
	return param
}

type BIRFunctionParameterImpl struct {
	BIRFunctionParameterBase
	BIRFunctionParameterMethods
}

func (v *BIRFunctionParameterImpl) String() string {
	return string(v.GetName())
}

type BIRFunctionData interface {
	BIRDocumentableNodeData
	NamedNode
	SetName(name Name)
	GetName() Name
	SetOriginalName(originalName Name)
	GetOriginalName() Name
	SetFlags(flags int64)
	GetFlags() int64
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetType(type_ BInvokableType)
	GetType() BInvokableType
	SetRequiredParams(requiredParams *[]BIRParameter)
	GetRequiredParams() *[]BIRParameter
	SetReceiver(receiver BIRVariableDcl)
	GetReceiver() BIRVariableDcl
	SetRestParam(restParam BIRParameter)
	GetRestParam() BIRParameter
	SetArgsCount(argsCount int)
	GetArgsCount() int
	SetLocalVars(localVars *[]BIRVariableDcl)
	GetLocalVars() *[]BIRVariableDcl
	SetReturnVariable(returnVariable BIRVariableDcl)
	GetReturnVariable() BIRVariableDcl
	SetParameters(parameters *[]BIRFunctionParameter)
	GetParameters() *[]BIRFunctionParameter
	SetBasicBlocks(basicBlocks *[]BIRBasicBlock)
	GetBasicBlocks() *[]BIRBasicBlock
	SetErrorTable(errorTable *[]BIRErrorEntry)
	GetErrorTable() *[]BIRErrorEntry
	SetWorkerName(workerName Name)
	GetWorkerName() Name
	SetWorkerChannels(workerChannels []ChannelDetails)
	GetWorkerChannels() []ChannelDetails
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
	SetAnnotAttachmentsOnExternal(annotAttachmentsOnExternal *[]BIRAnnotationAttachment)
	GetAnnotAttachmentsOnExternal() *[]BIRAnnotationAttachment
	SetReturnTypeAnnots(returnTypeAnnots *[]BIRAnnotationAttachment)
	GetReturnTypeAnnots() *[]BIRAnnotationAttachment
	SetDependentGlobalVars(dependentGlobalVars *[]BIRGlobalVariableDcl)
	GetDependentGlobalVars() *[]BIRGlobalVariableDcl
	SetPathParams(pathParams *[]BIRVariableDcl)
	GetPathParams() *[]BIRVariableDcl
	SetRestPathParam(restPathParam BIRVariableDcl)
	GetRestPathParam() BIRVariableDcl
	SetResourcePath(resourcePath *[]Name)
	GetResourcePath() *[]Name
	SetResourcePathSegmentPosList(resourcePathSegmentPosList *[]diagnostics.Location)
	GetResourcePathSegmentPosList() *[]diagnostics.Location
	SetAccessor(accessor Name)
	GetAccessor() Name
	SetPathSegmentTypeList(pathSegmentTypeList *[]BType)
	GetPathSegmentTypeList() *[]BType
	SetHasWorkers(hasWorkers bool)
	GetHasWorkers() bool
}

type BIRFunctionBase struct {
	BIRDocumentableNodeBase
	Name                       Name
	OriginalName               Name
	Flags                      int64
	Origin                     SymbolOrigin
	Type                       BInvokableType
	RequiredParams             []BIRParameter
	Receiver                   BIRVariableDcl
	RestParam                  BIRParameter
	ArgsCount                  int
	LocalVars                  []BIRVariableDcl
	ReturnVariable             BIRVariableDcl
	Parameters                 []BIRFunctionParameter
	BasicBlocks                []BIRBasicBlock
	ErrorTable                 []BIRErrorEntry
	WorkerName                 Name
	WorkerChannels             []ChannelDetails
	AnnotAttachments           []BIRAnnotationAttachment
	AnnotAttachmentsOnExternal []BIRAnnotationAttachment
	ReturnTypeAnnots           []BIRAnnotationAttachment
	DependentGlobalVars        []BIRGlobalVariableDcl
	PathParams                 []BIRVariableDcl
	RestPathParam              BIRVariableDcl
	ResourcePath               []Name
	ResourcePathSegmentPosList []diagnostics.Location
	Accessor                   Name
	PathSegmentTypeList        []BType
	HasWorkers                 bool
}

func (b *BIRFunctionBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRFunctionBase) GetName() Name {
	return b.Name
}

func (b *BIRFunctionBase) SetOriginalName(originalName Name) {
	b.OriginalName = originalName
}

func (b *BIRFunctionBase) GetOriginalName() Name {
	return b.OriginalName
}

func (b *BIRFunctionBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRFunctionBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRFunctionBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRFunctionBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRFunctionBase) SetType(type_ BInvokableType) {
	b.Type = type_
}

func (b *BIRFunctionBase) GetType() BInvokableType {
	return b.Type
}

func (b *BIRFunctionBase) SetRequiredParams(requiredParams *[]BIRParameter) {
	b.RequiredParams = *requiredParams
}

func (b *BIRFunctionBase) GetRequiredParams() *[]BIRParameter {
	return &b.RequiredParams
}

func (b *BIRFunctionBase) SetReceiver(receiver BIRVariableDcl) {
	b.Receiver = receiver
}

func (b *BIRFunctionBase) GetReceiver() BIRVariableDcl {
	return b.Receiver
}

func (b *BIRFunctionBase) SetRestParam(restParam BIRParameter) {
	b.RestParam = restParam
}

func (b *BIRFunctionBase) GetRestParam() BIRParameter {
	return b.RestParam
}

func (b *BIRFunctionBase) SetArgsCount(argsCount int) {
	b.ArgsCount = argsCount
}

func (b *BIRFunctionBase) GetArgsCount() int {
	return b.ArgsCount
}

func (b *BIRFunctionBase) SetLocalVars(localVars *[]BIRVariableDcl) {
	b.LocalVars = *localVars
}

func (b *BIRFunctionBase) GetLocalVars() *[]BIRVariableDcl {
	return &b.LocalVars
}

func (b *BIRFunctionBase) SetReturnVariable(returnVariable BIRVariableDcl) {
	b.ReturnVariable = returnVariable
}

func (b *BIRFunctionBase) GetReturnVariable() BIRVariableDcl {
	return b.ReturnVariable
}

func (b *BIRFunctionBase) SetParameters(parameters *[]BIRFunctionParameter) {
	b.Parameters = *parameters
}

func (b *BIRFunctionBase) GetParameters() *[]BIRFunctionParameter {
	return &b.Parameters
}

func (b *BIRFunctionBase) SetBasicBlocks(basicBlocks *[]BIRBasicBlock) {
	b.BasicBlocks = *basicBlocks
}

func (b *BIRFunctionBase) GetBasicBlocks() *[]BIRBasicBlock {
	return &b.BasicBlocks
}

func (b *BIRFunctionBase) SetErrorTable(errorTable *[]BIRErrorEntry) {
	b.ErrorTable = *errorTable
}

func (b *BIRFunctionBase) GetErrorTable() *[]BIRErrorEntry {
	return &b.ErrorTable
}

func (b *BIRFunctionBase) SetWorkerName(workerName Name) {
	b.WorkerName = workerName
}

func (b *BIRFunctionBase) GetWorkerName() Name {
	return b.WorkerName
}

func (b *BIRFunctionBase) SetWorkerChannels(workerChannels []ChannelDetails) {
	b.WorkerChannels = workerChannels
}

func (b *BIRFunctionBase) GetWorkerChannels() []ChannelDetails {
	return b.WorkerChannels
}

func (b *BIRFunctionBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRFunctionBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

func (b *BIRFunctionBase) SetAnnotAttachmentsOnExternal(annotAttachmentsOnExternal *[]BIRAnnotationAttachment) {
	b.AnnotAttachmentsOnExternal = *annotAttachmentsOnExternal
}

func (b *BIRFunctionBase) GetAnnotAttachmentsOnExternal() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachmentsOnExternal
}

func (b *BIRFunctionBase) SetReturnTypeAnnots(returnTypeAnnots *[]BIRAnnotationAttachment) {
	b.ReturnTypeAnnots = *returnTypeAnnots
}

func (b *BIRFunctionBase) GetReturnTypeAnnots() *[]BIRAnnotationAttachment {
	return &b.ReturnTypeAnnots
}

func (b *BIRFunctionBase) SetDependentGlobalVars(dependentGlobalVars *[]BIRGlobalVariableDcl) {
	b.DependentGlobalVars = *dependentGlobalVars
}

func (b *BIRFunctionBase) GetDependentGlobalVars() *[]BIRGlobalVariableDcl {
	return &b.DependentGlobalVars
}

func (b *BIRFunctionBase) SetPathParams(pathParams *[]BIRVariableDcl) {
	b.PathParams = *pathParams
}

func (b *BIRFunctionBase) GetPathParams() *[]BIRVariableDcl {
	return &b.PathParams
}

func (b *BIRFunctionBase) SetRestPathParam(restPathParam BIRVariableDcl) {
	b.RestPathParam = restPathParam
}

func (b *BIRFunctionBase) GetRestPathParam() BIRVariableDcl {
	return b.RestPathParam
}

func (b *BIRFunctionBase) SetResourcePath(resourcePath *[]Name) {
	b.ResourcePath = *resourcePath
}

func (b *BIRFunctionBase) GetResourcePath() *[]Name {
	return &b.ResourcePath
}

func (b *BIRFunctionBase) SetResourcePathSegmentPosList(resourcePathSegmentPosList *[]diagnostics.Location) {
	b.ResourcePathSegmentPosList = *resourcePathSegmentPosList
}

func (b *BIRFunctionBase) GetResourcePathSegmentPosList() *[]diagnostics.Location {
	return &b.ResourcePathSegmentPosList
}

func (b *BIRFunctionBase) SetAccessor(accessor Name) {
	b.Accessor = accessor
}

func (b *BIRFunctionBase) GetAccessor() Name {
	return b.Accessor
}

func (b *BIRFunctionBase) SetPathSegmentTypeList(pathSegmentTypeList *[]BType) {
	b.PathSegmentTypeList = *pathSegmentTypeList
}

func (b *BIRFunctionBase) GetPathSegmentTypeList() *[]BType {
	return &b.PathSegmentTypeList
}

func (b *BIRFunctionBase) SetHasWorkers(hasWorkers bool) {
	b.HasWorkers = hasWorkers
}

func (b *BIRFunctionBase) GetHasWorkers() bool {
	return b.HasWorkers
}

type BIRFunction interface {
	BIRFunctionData
	BIRDocumentableNode
}

type BIRFunctionMethods struct {
	Self BIRFunction
}

func (m *BIRFunctionMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRFunction(m.Self)
}

func NewBIRFunction(pos diagnostics.Location, name Name, originalName Name, flags int64, origin SymbolOrigin, type_ BInvokableType, requiredParams []BIRParameter, receiver BIRVariableDcl, restParam BIRParameter, argsCount int, localVars []BIRVariableDcl, returnVariable BIRVariableDcl, parameters []BIRFunctionParameter, basicBlocks []BIRBasicBlock, errorTable []BIRErrorEntry, workerName Name, workerChannels []ChannelDetails, annotAttachments []BIRAnnotationAttachment, returnTypeAnnots []BIRAnnotationAttachment, dependentGlobalVars []BIRGlobalVariableDcl) BIRFunction {
	fn := &BIRFunctionImpl{
		BIRFunctionBase: BIRFunctionBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:                name,
			OriginalName:        originalName,
			Flags:               flags,
			Origin:              origin,
			Type:                type_,
			RequiredParams:      requiredParams,
			Receiver:            receiver,
			RestParam:           restParam,
			ArgsCount:           argsCount,
			LocalVars:           localVars,
			ReturnVariable:      returnVariable,
			Parameters:          parameters,
			BasicBlocks:         basicBlocks,
			ErrorTable:          errorTable,
			WorkerName:          workerName,
			WorkerChannels:      workerChannels,
			AnnotAttachments:    annotAttachments,
			ReturnTypeAnnots:    returnTypeAnnots,
			DependentGlobalVars: dependentGlobalVars,
		},
		BIRFunctionMethods: BIRFunctionMethods{},
	}
	fn.BIRFunctionMethods.Self = fn
	return fn
}

func NewBIRFunctionWithSendInsCount(pos diagnostics.Location, name Name, originalName Name, flags int64, type_ BInvokableType, workerName Name, sendInsCount int, origin SymbolOrigin) BIRFunction {
	fn := &BIRFunctionImpl{
		BIRFunctionBase: BIRFunctionBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:             name,
			OriginalName:     originalName,
			Flags:            flags,
			Type:             type_,
			LocalVars:        []BIRVariableDcl{},
			Parameters:       []BIRFunctionParameter{},
			RequiredParams:   []BIRParameter{},
			BasicBlocks:      []BIRBasicBlock{},
			ErrorTable:       []BIRErrorEntry{},
			WorkerName:       workerName,
			WorkerChannels:   make([]ChannelDetails, sendInsCount),
			AnnotAttachments: []BIRAnnotationAttachment{},
			ReturnTypeAnnots: []BIRAnnotationAttachment{},
			Origin:           origin,
		},
		BIRFunctionMethods: BIRFunctionMethods{},
	}
	fn.BIRFunctionMethods.Self = fn
	return fn
}

func NewBIRFunctionSimple(pos diagnostics.Location, name Name, flags int64, type_ BInvokableType, workerName Name, sendInsCount int, origin SymbolOrigin) BIRFunction {
	return NewBIRFunctionWithSendInsCount(pos, name, name, flags, type_, workerName, sendInsCount, origin)
}

type BIRFunctionImpl struct {
	BIRFunctionBase
	BIRFunctionMethods
}

func (f *BIRFunctionImpl) Duplicate() BIRFunction {
	newFn := NewBIRFunctionSimple(f.GetPos(), f.GetName(), f.GetFlags(), f.GetType(), f.GetWorkerName(), 0, f.GetOrigin())
	newFn.SetLocalVars(f.GetLocalVars())
	newFn.SetParameters(f.GetParameters())
	newFn.SetRequiredParams(f.GetRequiredParams())
	newFn.SetBasicBlocks(f.GetBasicBlocks())
	newFn.SetErrorTable(f.GetErrorTable())
	newFn.SetWorkerChannels(f.GetWorkerChannels())
	newFn.SetAnnotAttachments(f.GetAnnotAttachments())
	newFn.SetAnnotAttachmentsOnExternal(f.GetAnnotAttachmentsOnExternal())
	newFn.SetReturnTypeAnnots(f.GetReturnTypeAnnots())
	return newFn
}

type BIRBasicBlockData interface {
	BIRNodeData
	SetNumber(number int)
	GetNumber() int
	SetId(id Name)
	GetId() Name
	SetInstructions(instructions *[]BIRNonTerminator)
	GetInstructions() *[]BIRNonTerminator
	SetTerminator(terminator BIRTerminator)
	GetTerminator() BIRTerminator
}

type BIRBasicBlockBase struct {
	BIRNodeBase
	Number       int
	Id           Name
	Instructions []BIRNonTerminator
	Terminator   BIRTerminator
}

const BIR_BASIC_BLOCK_PREFIX = "bb"

func (b *BIRBasicBlockBase) SetNumber(number int) {
	b.Number = number
}

func (b *BIRBasicBlockBase) GetNumber() int {
	return b.Number
}

func (b *BIRBasicBlockBase) SetId(id Name) {
	b.Id = id
}

func (b *BIRBasicBlockBase) GetId() Name {
	return b.Id
}

func (b *BIRBasicBlockBase) SetInstructions(instructions *[]BIRNonTerminator) {
	b.Instructions = *instructions
}

func (b *BIRBasicBlockBase) GetInstructions() *[]BIRNonTerminator {
	return &b.Instructions
}

func (b *BIRBasicBlockBase) SetTerminator(terminator BIRTerminator) {
	b.Terminator = terminator
}

func (b *BIRBasicBlockBase) GetTerminator() BIRTerminator {
	return b.Terminator
}

type BIRBasicBlock interface {
	BIRBasicBlockData
	BIRNode
}

type BIRBasicBlockMethods struct {
	Self BIRBasicBlock
}

func (m *BIRBasicBlockMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRBasicBlock(m.Self)
}

func NewBIRBasicBlock(id Name, number int) BIRBasicBlock {
	bb := &BIRBasicBlockImpl{
		BIRBasicBlockBase: BIRBasicBlockBase{
			BIRNodeBase:  BIRNodeBase{},
			Number:       number,
			Id:           id,
			Instructions: []BIRNonTerminator{},
		},
		BIRBasicBlockMethods: BIRBasicBlockMethods{},
	}
	bb.BIRBasicBlockMethods.Self = bb
	return bb
}

func NewBIRBasicBlockWithNumber(number int) BIRBasicBlock {
	return NewBIRBasicBlock(Name(BIR_BASIC_BLOCK_PREFIX+strconv.Itoa(number)), number)
}

func NewBIRBasicBlockWithIdPrefix(idPrefix string, number int) BIRBasicBlock {
	return NewBIRBasicBlock(Name(idPrefix+strconv.Itoa(number)), number)
}

type BIRBasicBlockImpl struct {
	BIRBasicBlockBase
	BIRBasicBlockMethods
}

type BIRTypeDefinitionData interface {
	BIRDocumentableNodeData
	NamedNode
	SetName(name Name)
	GetName() Name
	SetOriginalName(originalName Name)
	GetOriginalName() Name
	SetInternalName(internalName Name)
	GetInternalName() Name
	SetAttachedFuncs(attachedFuncs *[]BIRFunction)
	GetAttachedFuncs() *[]BIRFunction
	SetFlags(flags int64)
	GetFlags() int64
	SetType(type_ BType)
	GetType() BType
	SetIsBuiltin(isBuiltin bool)
	GetIsBuiltin() bool
	SetReferencedTypes(referencedTypes *[]BType)
	GetReferencedTypes() *[]BType
	SetReferenceType(referenceType BType)
	GetReferenceType() BType
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
	SetIndex(index int)
	GetIndex() int
}

type BIRTypeDefinitionBase struct {
	BIRDocumentableNodeBase
	Name             Name
	OriginalName     Name
	InternalName     Name
	AttachedFuncs    []BIRFunction
	Flags            int64
	Type             BType
	IsBuiltin        bool
	ReferencedTypes  []BType
	ReferenceType    BType
	Origin           SymbolOrigin
	AnnotAttachments []BIRAnnotationAttachment
	Index            int
}

func (b *BIRTypeDefinitionBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRTypeDefinitionBase) GetName() Name {
	return b.Name
}

func (b *BIRTypeDefinitionBase) SetOriginalName(originalName Name) {
	b.OriginalName = originalName
}

func (b *BIRTypeDefinitionBase) GetOriginalName() Name {
	return b.OriginalName
}

func (b *BIRTypeDefinitionBase) SetInternalName(internalName Name) {
	b.InternalName = internalName
}

func (b *BIRTypeDefinitionBase) GetInternalName() Name {
	return b.InternalName
}

func (b *BIRTypeDefinitionBase) SetAttachedFuncs(attachedFuncs *[]BIRFunction) {
	b.AttachedFuncs = *attachedFuncs
}

func (b *BIRTypeDefinitionBase) GetAttachedFuncs() *[]BIRFunction {
	return &b.AttachedFuncs
}

func (b *BIRTypeDefinitionBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRTypeDefinitionBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRTypeDefinitionBase) SetType(type_ BType) {
	b.Type = type_
}

func (b *BIRTypeDefinitionBase) GetType() BType {
	return b.Type
}

func (b *BIRTypeDefinitionBase) SetIsBuiltin(isBuiltin bool) {
	b.IsBuiltin = isBuiltin
}

func (b *BIRTypeDefinitionBase) GetIsBuiltin() bool {
	return b.IsBuiltin
}

func (b *BIRTypeDefinitionBase) SetReferencedTypes(referencedTypes *[]BType) {
	b.ReferencedTypes = *referencedTypes
}

func (b *BIRTypeDefinitionBase) GetReferencedTypes() *[]BType {
	return &b.ReferencedTypes
}

func (b *BIRTypeDefinitionBase) SetReferenceType(referenceType BType) {
	b.ReferenceType = referenceType
}

func (b *BIRTypeDefinitionBase) GetReferenceType() BType {
	return b.ReferenceType
}

func (b *BIRTypeDefinitionBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRTypeDefinitionBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRTypeDefinitionBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRTypeDefinitionBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

func (b *BIRTypeDefinitionBase) SetIndex(index int) {
	b.Index = index
}

func (b *BIRTypeDefinitionBase) GetIndex() int {
	return b.Index
}

type BIRTypeDefinition interface {
	BIRTypeDefinitionData
	BIRDocumentableNode
}

type BIRTypeDefinitionMethods struct {
	Self BIRTypeDefinition
}

func (m *BIRTypeDefinitionMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTypeDefinition(m.Self)
}

func NewBIRTypeDefinition(pos diagnostics.Location, internalName Name, flags int64, isBuiltin bool, type_ BType, attachedFuncs []BIRFunction, origin SymbolOrigin, name Name, originalName Name) BIRTypeDefinition {
	td := &BIRTypeDefinitionImpl{
		BIRTypeDefinitionBase: BIRTypeDefinitionBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			InternalName:     internalName,
			Flags:            flags,
			IsBuiltin:        isBuiltin,
			Type:             type_,
			AttachedFuncs:    attachedFuncs,
			ReferencedTypes:  []BType{},
			Origin:           origin,
			Name:             name,
			OriginalName:     originalName,
			AnnotAttachments: []BIRAnnotationAttachment{},
		},
		BIRTypeDefinitionMethods: BIRTypeDefinitionMethods{},
	}
	td.BIRTypeDefinitionMethods.Self = td
	return td
}

func NewBIRTypeDefinitionSimple(pos diagnostics.Location, name Name, originalName Name, flags int64, isBuiltin bool, type_ BType, attachedFuncs []BIRFunction, origin SymbolOrigin) BIRTypeDefinition {
	return NewBIRTypeDefinition(pos, name, flags, isBuiltin, type_, attachedFuncs, origin, name, originalName)
}

type BIRTypeDefinitionImpl struct {
	BIRTypeDefinitionBase
	BIRTypeDefinitionMethods
}

type BIRErrorEntryData interface {
	BIRNodeData
	SetTrapBB(trapBB BIRBasicBlock)
	GetTrapBB() BIRBasicBlock
	SetEndBB(endBB BIRBasicBlock)
	GetEndBB() BIRBasicBlock
	SetErrorOp(errorOp BIROperand)
	GetErrorOp() BIROperand
	SetTargetBB(targetBB BIRBasicBlock)
	GetTargetBB() BIRBasicBlock
}

type BIRErrorEntryBase struct {
	BIRNodeBase
	TrapBB   BIRBasicBlock
	EndBB    BIRBasicBlock
	ErrorOp  BIROperand
	TargetBB BIRBasicBlock
}

func (b *BIRErrorEntryBase) SetTrapBB(trapBB BIRBasicBlock) {
	b.TrapBB = trapBB
}

func (b *BIRErrorEntryBase) GetTrapBB() BIRBasicBlock {
	return b.TrapBB
}

func (b *BIRErrorEntryBase) SetEndBB(endBB BIRBasicBlock) {
	b.EndBB = endBB
}

func (b *BIRErrorEntryBase) GetEndBB() BIRBasicBlock {
	return b.EndBB
}

func (b *BIRErrorEntryBase) SetErrorOp(errorOp BIROperand) {
	b.ErrorOp = errorOp
}

func (b *BIRErrorEntryBase) GetErrorOp() BIROperand {
	return b.ErrorOp
}

func (b *BIRErrorEntryBase) SetTargetBB(targetBB BIRBasicBlock) {
	b.TargetBB = targetBB
}

func (b *BIRErrorEntryBase) GetTargetBB() BIRBasicBlock {
	return b.TargetBB
}

type BIRErrorEntry interface {
	BIRErrorEntryData
	BIRNode
}

type BIRErrorEntryMethods struct {
	Self BIRErrorEntry
}

func (m *BIRErrorEntryMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRErrorEntry(m.Self)
}

func NewBIRErrorEntry(trapBB BIRBasicBlock, endBB BIRBasicBlock, errorOp BIROperand, targetBB BIRBasicBlock) BIRErrorEntry {
	entry := &BIRErrorEntryImpl{
		BIRErrorEntryBase: BIRErrorEntryBase{
			BIRNodeBase: BIRNodeBase{},
			TrapBB:      trapBB,
			EndBB:       endBB,
			ErrorOp:     errorOp,
			TargetBB:    targetBB,
		},
		BIRErrorEntryMethods: BIRErrorEntryMethods{},
	}
	entry.BIRErrorEntryMethods.Self = entry
	return entry
}

type BIRErrorEntryImpl struct {
	BIRErrorEntryBase
	BIRErrorEntryMethods
}

type ChannelDetails struct {
	Name                string
	ChannelInSameStrand bool
	Send                bool
}

func NewChannelDetails(name string, channelInSameStrand bool, send bool) *ChannelDetails {
	return &ChannelDetails{
		Name:                name,
		ChannelInSameStrand: channelInSameStrand,
		Send:                send,
	}
}

func (c *ChannelDetails) String() string {
	return c.Name
}

type BIRAnnotationData interface {
	BIRDocumentableNodeData
	SetName(name Name)
	GetName() Name
	SetOriginalName(originalName Name)
	GetOriginalName() Name
	SetFlags(flags int64)
	GetFlags() int64
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetAttachPoints(attachPoints *[]AttachPoint)
	GetAttachPoints() *[]AttachPoint
	SetAnnotationType(annotationType BType)
	GetAnnotationType() BType
	SetPackageID(packageID PackageID)
	GetPackageID() PackageID
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
}

type BIRAnnotationBase struct {
	BIRDocumentableNodeBase
	Name             Name
	OriginalName     Name
	Flags            int64
	Origin           SymbolOrigin
	AttachPoints     []AttachPoint
	AnnotationType   BType
	PackageID        PackageID
	AnnotAttachments []BIRAnnotationAttachment
}

func (b *BIRAnnotationBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRAnnotationBase) GetName() Name {
	return b.Name
}

func (b *BIRAnnotationBase) SetOriginalName(originalName Name) {
	b.OriginalName = originalName
}

func (b *BIRAnnotationBase) GetOriginalName() Name {
	return b.OriginalName
}

func (b *BIRAnnotationBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRAnnotationBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRAnnotationBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRAnnotationBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRAnnotationBase) SetAttachPoints(attachPoints *[]AttachPoint) {
	b.AttachPoints = *attachPoints
}

func (b *BIRAnnotationBase) GetAttachPoints() *[]AttachPoint {
	return &b.AttachPoints
}

func (b *BIRAnnotationBase) SetAnnotationType(annotationType BType) {
	b.AnnotationType = annotationType
}

func (b *BIRAnnotationBase) GetAnnotationType() BType {
	return b.AnnotationType
}

func (b *BIRAnnotationBase) SetPackageID(packageID PackageID) {
	b.PackageID = packageID
}

func (b *BIRAnnotationBase) GetPackageID() PackageID {
	return b.PackageID
}

func (b *BIRAnnotationBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRAnnotationBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

type BIRAnnotation interface {
	BIRAnnotationData
	BIRDocumentableNode
}

type BIRAnnotationMethods struct {
	Self BIRAnnotation
}

func (m *BIRAnnotationMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRAnnotation(m.Self)
}

func NewBIRAnnotation(pos diagnostics.Location, name Name, originalName Name, flags int64, points []AttachPoint, annotationType BType, origin SymbolOrigin) BIRAnnotation {
	ann := &BIRAnnotationImpl{
		BIRAnnotationBase: BIRAnnotationBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:             name,
			OriginalName:     originalName,
			Flags:            flags,
			AttachPoints:     points,
			AnnotationType:   annotationType,
			Origin:           origin,
			AnnotAttachments: []BIRAnnotationAttachment{},
		},
		BIRAnnotationMethods: BIRAnnotationMethods{},
	}
	ann.BIRAnnotationMethods.Self = ann
	return ann
}

type BIRAnnotationImpl struct {
	BIRAnnotationBase
	BIRAnnotationMethods
}

type BIRConstantData interface {
	BIRDocumentableNodeData
	SetName(name Name)
	GetName() Name
	SetFlags(flags int64)
	GetFlags() int64
	SetType(type_ BType)
	GetType() BType
	SetConstValue(constValue ConstValue)
	GetConstValue() ConstValue
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment)
	GetAnnotAttachments() *[]BIRAnnotationAttachment
}

type BIRConstantBase struct {
	BIRDocumentableNodeBase
	Name             Name
	Flags            int64
	Type             BType
	ConstValue       ConstValue
	Origin           SymbolOrigin
	AnnotAttachments []BIRAnnotationAttachment
}

func (b *BIRConstantBase) SetName(name Name) {
	b.Name = name
}

func (b *BIRConstantBase) GetName() Name {
	return b.Name
}

func (b *BIRConstantBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRConstantBase) GetFlags() int64 {
	return b.Flags
}

func (b *BIRConstantBase) SetType(type_ BType) {
	b.Type = type_
}

func (b *BIRConstantBase) GetType() BType {
	return b.Type
}

func (b *BIRConstantBase) SetConstValue(constValue ConstValue) {
	b.ConstValue = constValue
}

func (b *BIRConstantBase) GetConstValue() ConstValue {
	return b.ConstValue
}

func (b *BIRConstantBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRConstantBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRConstantBase) SetAnnotAttachments(annotAttachments *[]BIRAnnotationAttachment) {
	b.AnnotAttachments = *annotAttachments
}

func (b *BIRConstantBase) GetAnnotAttachments() *[]BIRAnnotationAttachment {
	return &b.AnnotAttachments
}

type BIRConstant interface {
	BIRConstantData
	BIRDocumentableNode
}

type BIRConstantMethods struct {
	Self BIRConstant
}

func (m *BIRConstantMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRConstant(m.Self)
}

func NewBIRConstant(pos diagnostics.Location, name Name, flags int64, type_ BType, constValue ConstValue, origin SymbolOrigin) BIRConstant {
	constant := &BIRConstantImpl{
		BIRConstantBase: BIRConstantBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:             name,
			Flags:            flags,
			Type:             type_,
			ConstValue:       constValue,
			Origin:           origin,
			AnnotAttachments: []BIRAnnotationAttachment{},
		},
		BIRConstantMethods: BIRConstantMethods{},
	}
	constant.BIRConstantMethods.Self = constant
	return constant
}

type BIRConstantImpl struct {
	BIRConstantBase
	BIRConstantMethods
}

type BIRAnnotationAttachmentData interface {
	BIRNodeData
	SetAnnotPkgId(annotPkgId PackageID)
	GetAnnotPkgId() PackageID
	SetAnnotTagRef(annotTagRef Name)
	GetAnnotTagRef() Name
}

type BIRAnnotationAttachmentBase struct {
	BIRNodeBase
	AnnotPkgId  PackageID
	AnnotTagRef Name
}

func (b *BIRAnnotationAttachmentBase) SetAnnotPkgId(annotPkgId PackageID) {
	b.AnnotPkgId = annotPkgId
}

func (b *BIRAnnotationAttachmentBase) GetAnnotPkgId() PackageID {
	return b.AnnotPkgId
}

func (b *BIRAnnotationAttachmentBase) SetAnnotTagRef(annotTagRef Name) {
	b.AnnotTagRef = annotTagRef
}

func (b *BIRAnnotationAttachmentBase) GetAnnotTagRef() Name {
	return b.AnnotTagRef
}

type BIRAnnotationAttachment interface {
	BIRAnnotationAttachmentData
	BIRNode
}

type BIRAnnotationAttachmentMethods struct {
	Self BIRAnnotationAttachment
}

func (m *BIRAnnotationAttachmentMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRAnnotationAttachment(m.Self)
}

func NewBIRAnnotationAttachment(pos diagnostics.Location, annotPkgId PackageID, annotTagRef Name) BIRAnnotationAttachment {
	att := &BIRAnnotationAttachmentImpl{
		BIRAnnotationAttachmentBase: BIRAnnotationAttachmentBase{
			BIRNodeBase: BIRNodeBase{
				Pos: pos,
			},
			AnnotPkgId:  annotPkgId,
			AnnotTagRef: annotTagRef,
		},
		BIRAnnotationAttachmentMethods: BIRAnnotationAttachmentMethods{},
	}
	att.BIRAnnotationAttachmentMethods.Self = att
	return att
}

type BIRAnnotationAttachmentImpl struct {
	BIRAnnotationAttachmentBase
	BIRAnnotationAttachmentMethods
}

type BIRConstAnnotationAttachmentData interface {
	BIRAnnotationAttachmentData
	SetAnnotValue(annotValue ConstValue)
	GetAnnotValue() ConstValue
}

type BIRConstAnnotationAttachmentBase struct {
	BIRAnnotationAttachmentBase
	AnnotValue ConstValue
}

func (b *BIRConstAnnotationAttachmentBase) SetAnnotValue(annotValue ConstValue) {
	b.AnnotValue = annotValue
}

func (b *BIRConstAnnotationAttachmentBase) GetAnnotValue() ConstValue {
	return b.AnnotValue
}

type BIRConstAnnotationAttachment interface {
	BIRConstAnnotationAttachmentData
	BIRAnnotationAttachment
}

type BIRConstAnnotationAttachmentMethods struct {
	Self BIRConstAnnotationAttachment
}

func (m *BIRConstAnnotationAttachmentMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRConstAnnotationAttachment(m.Self)
}

func NewBIRConstAnnotationAttachment(pos diagnostics.Location, annotPkgId PackageID, annotTagRef Name, annotValue ConstValue) BIRConstAnnotationAttachment {
	att := &BIRConstAnnotationAttachmentImpl{
		BIRConstAnnotationAttachmentBase: BIRConstAnnotationAttachmentBase{
			BIRAnnotationAttachmentBase: BIRAnnotationAttachmentBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				AnnotPkgId:  annotPkgId,
				AnnotTagRef: annotTagRef,
			},
			AnnotValue: annotValue,
		},
		BIRConstAnnotationAttachmentMethods: BIRConstAnnotationAttachmentMethods{},
	}
	att.BIRConstAnnotationAttachmentMethods.Self = att
	return att
}

type BIRConstAnnotationAttachmentImpl struct {
	BIRConstAnnotationAttachmentBase
	BIRConstAnnotationAttachmentMethods
}

type ConstValue struct {
	Type  BType
	Value interface{}
}

func NewConstValue(value interface{}, type_ BType) *ConstValue {
	return &ConstValue{
		Value: value,
		Type:  type_,
	}
}

type BIRDocumentableNodeData interface {
	BIRNodeData
	SetMarkdownDocAttachment(markdownDocAttachment MarkdownDocAttachment)
	GetMarkdownDocAttachment() MarkdownDocAttachment
}

type BIRDocumentableNodeBase struct {
	BIRNodeBase
	MarkdownDocAttachment MarkdownDocAttachment
}

func (b *BIRDocumentableNodeBase) SetMarkdownDocAttachment(markdownDocAttachment MarkdownDocAttachment) {
	b.MarkdownDocAttachment = markdownDocAttachment
}

func (b *BIRDocumentableNodeBase) GetMarkdownDocAttachment() MarkdownDocAttachment {
	return b.MarkdownDocAttachment
}

type BIRDocumentableNode interface {
	BIRDocumentableNodeData
	BIRNode
}

type BIRLockDetailsHolder struct {
	locks []BIRTerminatorLock
}

func NewBIRLockDetailsHolder() *BIRLockDetailsHolder {
	return &BIRLockDetailsHolder{
		locks: []BIRTerminatorLock{},
	}
}

func (h *BIRLockDetailsHolder) IsEmpty() bool {
	return len(h.locks) == 0
}

func (h *BIRLockDetailsHolder) RemoveLastLock() {
	h.locks = h.locks[:len(h.locks)-1]
}

func (h *BIRLockDetailsHolder) GetLock(index int) BIRTerminatorLock {
	return h.locks[index]
}

func (h *BIRLockDetailsHolder) AddLock(lock BIRTerminatorLock) {
	h.locks = append(h.locks, lock)
}

func (h *BIRLockDetailsHolder) Size() int {
	return len(h.locks)
}

type BIRMappingConstructorEntryData interface {
	IsKeyValuePair() bool
}

type BIRMappingConstructorEntry interface {
	BIRMappingConstructorEntryData
}

type BIRMappingConstructorEntryMethods struct {
	Self BIRMappingConstructorEntry
}

func (m *BIRMappingConstructorEntryMethods) IsKeyValuePair() bool {
	return true
}

type BIRMappingConstructorKeyValueEntryData interface {
	BIRMappingConstructorEntryData
	SetKeyOp(keyOp BIROperand)
	GetKeyOp() BIROperand
	SetValueOp(valueOp BIROperand)
	GetValueOp() BIROperand
}

type BIRMappingConstructorKeyValueEntryBase struct {
	KeyOp   BIROperand
	ValueOp BIROperand
}

func (b *BIRMappingConstructorKeyValueEntryBase) SetKeyOp(keyOp BIROperand) {
	b.KeyOp = keyOp
}

func (b *BIRMappingConstructorKeyValueEntryBase) GetKeyOp() BIROperand {
	return b.KeyOp
}

func (b *BIRMappingConstructorKeyValueEntryBase) SetValueOp(valueOp BIROperand) {
	b.ValueOp = valueOp
}

func (b *BIRMappingConstructorKeyValueEntryBase) GetValueOp() BIROperand {
	return b.ValueOp
}

type BIRMappingConstructorKeyValueEntry interface {
	BIRMappingConstructorKeyValueEntryData
	BIRMappingConstructorEntry
}

type BIRMappingConstructorKeyValueEntryMethods struct {
	BIRMappingConstructorEntryMethods
	Self BIRMappingConstructorKeyValueEntry
}

func NewBIRMappingConstructorKeyValueEntry(keyOp BIROperand, valueOp BIROperand) BIRMappingConstructorKeyValueEntry {
	entry := &BIRMappingConstructorKeyValueEntryImpl{
		BIRMappingConstructorKeyValueEntryBase: BIRMappingConstructorKeyValueEntryBase{
			KeyOp:   keyOp,
			ValueOp: valueOp,
		},
		BIRMappingConstructorKeyValueEntryMethods: BIRMappingConstructorKeyValueEntryMethods{
			BIRMappingConstructorEntryMethods: BIRMappingConstructorEntryMethods{},
		},
	}
	entry.BIRMappingConstructorKeyValueEntryMethods.Self = entry
	entry.BIRMappingConstructorKeyValueEntryMethods.BIRMappingConstructorEntryMethods.Self = entry
	return entry
}

type BIRMappingConstructorKeyValueEntryImpl struct {
	BIRMappingConstructorKeyValueEntryBase
	BIRMappingConstructorKeyValueEntryMethods
}

type BIRMappingConstructorSpreadFieldEntryData interface {
	BIRMappingConstructorEntryData
	SetExprOp(exprOp BIROperand)
	GetExprOp() BIROperand
}

type BIRMappingConstructorSpreadFieldEntryBase struct {
	ExprOp BIROperand
}

func (b *BIRMappingConstructorSpreadFieldEntryBase) SetExprOp(exprOp BIROperand) {
	b.ExprOp = exprOp
}

func (b *BIRMappingConstructorSpreadFieldEntryBase) GetExprOp() BIROperand {
	return b.ExprOp
}

type BIRMappingConstructorSpreadFieldEntry interface {
	BIRMappingConstructorSpreadFieldEntryData
	BIRMappingConstructorEntry
}

type BIRMappingConstructorSpreadFieldEntryMethods struct {
	BIRMappingConstructorEntryMethods
	Self BIRMappingConstructorSpreadFieldEntry
}

func (m *BIRMappingConstructorSpreadFieldEntryMethods) IsKeyValuePair() bool {
	return false
}

func NewBIRMappingConstructorSpreadFieldEntry(exprOp BIROperand) BIRMappingConstructorSpreadFieldEntry {
	entry := &BIRMappingConstructorSpreadFieldEntryImpl{
		BIRMappingConstructorSpreadFieldEntryBase: BIRMappingConstructorSpreadFieldEntryBase{
			ExprOp: exprOp,
		},
		BIRMappingConstructorSpreadFieldEntryMethods: BIRMappingConstructorSpreadFieldEntryMethods{
			BIRMappingConstructorEntryMethods: BIRMappingConstructorEntryMethods{},
		},
	}
	entry.BIRMappingConstructorSpreadFieldEntryMethods.Self = entry
	entry.BIRMappingConstructorSpreadFieldEntryMethods.BIRMappingConstructorEntryMethods.Self = entry
	return entry
}

type BIRMappingConstructorSpreadFieldEntryImpl struct {
	BIRMappingConstructorSpreadFieldEntryBase
	BIRMappingConstructorSpreadFieldEntryMethods
}

type BIRListConstructorEntryData interface {
	SetExprOp(exprOp BIROperand)
	GetExprOp() BIROperand
}

type BIRListConstructorEntryBase struct {
	ExprOp BIROperand
}

func (b *BIRListConstructorEntryBase) SetExprOp(exprOp BIROperand) {
	b.ExprOp = exprOp
}

func (b *BIRListConstructorEntryBase) GetExprOp() BIROperand {
	return b.ExprOp
}

type BIRListConstructorEntry interface {
	BIRListConstructorEntryData
}

type BIRListConstructorSpreadMemberEntryData interface {
	BIRListConstructorEntryData
}

type BIRListConstructorSpreadMemberEntry interface {
	BIRListConstructorSpreadMemberEntryData
	BIRListConstructorEntry
}

type BIRListConstructorSpreadMemberEntryImpl struct {
	BIRListConstructorEntryBase
}

func NewBIRListConstructorSpreadMemberEntry(exprOp BIROperand) BIRListConstructorSpreadMemberEntry {
	return &BIRListConstructorSpreadMemberEntryImpl{
		BIRListConstructorEntryBase: BIRListConstructorEntryBase{
			ExprOp: exprOp,
		},
	}
}

type BIRListConstructorExprEntryData interface {
	BIRListConstructorEntryData
}

type BIRListConstructorExprEntry interface {
	BIRListConstructorExprEntryData
	BIRListConstructorEntry
}

type BIRListConstructorExprEntryImpl struct {
	BIRListConstructorEntryBase
}

func NewBIRListConstructorExprEntry(exprOp BIROperand) BIRListConstructorExprEntry {
	return &BIRListConstructorExprEntryImpl{
		BIRListConstructorEntryBase: BIRListConstructorEntryBase{
			ExprOp: exprOp,
		},
	}
}

type BIRServiceDeclarationData interface {
	BIRDocumentableNodeData
	SetAttachPoint(attachPoint *[]string)
	GetAttachPoint() *[]string
	SetAttachPointLiteral(attachPointLiteral string)
	GetAttachPointLiteral() string
	SetListenerTypes(listenerTypes *[]BType)
	GetListenerTypes() *[]BType
	SetGeneratedName(generatedName Name)
	GetGeneratedName() Name
	SetAssociatedClassName(associatedClassName Name)
	GetAssociatedClassName() Name
	SetType(type_ BType)
	GetType() BType
	SetOrigin(origin SymbolOrigin)
	GetOrigin() SymbolOrigin
	SetFlags(flags int64)
	GetFlags() int64
}

type BIRServiceDeclarationBase struct {
	BIRDocumentableNodeBase
	AttachPoint         []string
	AttachPointLiteral  string
	ListenerTypes       []BType
	GeneratedName       Name
	AssociatedClassName Name
	Type                BType
	Origin              SymbolOrigin
	Flags               int64
}

func (b *BIRServiceDeclarationBase) SetAttachPoint(attachPoint *[]string) {
	b.AttachPoint = *attachPoint
}

func (b *BIRServiceDeclarationBase) GetAttachPoint() *[]string {
	return &b.AttachPoint
}

func (b *BIRServiceDeclarationBase) SetAttachPointLiteral(attachPointLiteral string) {
	b.AttachPointLiteral = attachPointLiteral
}

func (b *BIRServiceDeclarationBase) GetAttachPointLiteral() string {
	return b.AttachPointLiteral
}

func (b *BIRServiceDeclarationBase) SetListenerTypes(listenerTypes *[]BType) {
	b.ListenerTypes = *listenerTypes
}

func (b *BIRServiceDeclarationBase) GetListenerTypes() *[]BType {
	return &b.ListenerTypes
}

func (b *BIRServiceDeclarationBase) SetGeneratedName(generatedName Name) {
	b.GeneratedName = generatedName
}

func (b *BIRServiceDeclarationBase) GetGeneratedName() Name {
	return b.GeneratedName
}

func (b *BIRServiceDeclarationBase) SetAssociatedClassName(associatedClassName Name) {
	b.AssociatedClassName = associatedClassName
}

func (b *BIRServiceDeclarationBase) GetAssociatedClassName() Name {
	return b.AssociatedClassName
}

func (b *BIRServiceDeclarationBase) SetType(type_ BType) {
	b.Type = type_
}

func (b *BIRServiceDeclarationBase) GetType() BType {
	return b.Type
}

func (b *BIRServiceDeclarationBase) SetOrigin(origin SymbolOrigin) {
	b.Origin = origin
}

func (b *BIRServiceDeclarationBase) GetOrigin() SymbolOrigin {
	return b.Origin
}

func (b *BIRServiceDeclarationBase) SetFlags(flags int64) {
	b.Flags = flags
}

func (b *BIRServiceDeclarationBase) GetFlags() int64 {
	return b.Flags
}

type BIRServiceDeclaration interface {
	BIRServiceDeclarationData
	BIRDocumentableNode
}

type BIRServiceDeclarationMethods struct {
	Self BIRServiceDeclaration
}

func (m *BIRServiceDeclarationMethods) Accept(visitor BIRVisitor) {
	visitor.VisitBIRServiceDeclaration(m.Self)
}

func NewBIRServiceDeclaration(attachPoint []string, attachPointLiteral string, listenerTypes []BType, generatedName Name, associatedClassName Name, type_ BType, origin SymbolOrigin, flags int64, location diagnostics.Location) BIRServiceDeclaration {
	svc := &BIRServiceDeclarationImpl{
		BIRServiceDeclarationBase: BIRServiceDeclarationBase{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: location,
				},
			},
			AttachPoint:         attachPoint,
			AttachPointLiteral:  attachPointLiteral,
			ListenerTypes:       listenerTypes,
			GeneratedName:       generatedName,
			AssociatedClassName: associatedClassName,
			Type:                type_,
			Origin:              origin,
			Flags:               flags,
		},
		BIRServiceDeclarationMethods: BIRServiceDeclarationMethods{},
	}
	svc.BIRServiceDeclarationMethods.Self = svc
	return svc
}

type BIRServiceDeclarationImpl struct {
	BIRServiceDeclarationBase
	BIRServiceDeclarationMethods
}
