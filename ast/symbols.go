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

package ast

import (
	"ballerina-lang-go/common"
	"strings"
)

type SymbolKind uint

const (
	SymbolKind_PACKAGE SymbolKind = iota
	SymbolKind_STRUCT
	SymbolKind_OBJECT
	SymbolKind_RECORD
	SymbolKind_CONNECTOR
	SymbolKind_ACTION
	SymbolKind_SERVICE
	SymbolKind_RESOURCE
	SymbolKind_FUNCTION
	SymbolKind_WORKER
	SymbolKind_ANNOTATION
	SymbolKind_ANNOTATION_ATTRIBUTE
	SymbolKind_CONSTANT
	SymbolKind_VARIABLE
	SymbolKind_PACKAGE_VARIABLE
	SymbolKind_TRANSFORMER
	SymbolKind_TYPE_DEF
	SymbolKind_ENUM
	SymbolKind_ERROR

	SymbolKind_PARAMETER
	SymbolKind_PATH_PARAMETER
	SymbolKind_PATH_REST_PARAMETER
	SymbolKind_LOCAL_VARIABLE
	SymbolKind_SERVICE_VARIABLE
	SymbolKind_CONNECTOR_VARIABLE

	SymbolKind_CAST_OPERATOR
	SymbolKind_CONVERSION_OPERATOR
	SymbolKind_TYPEOF_OPERATOR

	SymbolKind_XMLNS
	SymbolKind_SCOPE
	SymbolKind_OTHER

	SymbolKind_INVOKABLE_TYPE

	SymbolKind_RESOURCE_PATH_IDENTIFIER_SEGMENT
	SymbolKind_RESOURCE_PATH_PARAM_SEGMENT
	SymbolKind_RESOURCE_PATH_REST_PARAM_SEGMENT
	SymbolKind_RESOURCE_ROOT_PATH_SEGMENT

	SymbolKind_SEQUENCE
)

type SymbolOrigin uint8

const (
	SymbolOrigin_BUILTIN SymbolOrigin = iota
	SymbolOrigin_SOURCE
	SymbolOrigin_COMPILED_SOURCE
	SymbolOrigin_VIRTUAL
)

type DiagnosticState uint8

const (
	DiagnosticState_VALID DiagnosticState = iota
	DiagnosticState_REDECLARED
	DiagnosticState_UNKNOWN_TYPE
)

type SymTag int64

const (
	SymTag_NIL                   SymTag = 0
	SymTag_IMPORT                SymTag = 1
	SymTag_ANNOTATION            SymTag = 1 << 1
	SymTag_MAIN                  SymTag = 1 << 2
	SymTag_TYPE                  SymTag = 1<<3 | SymTag_MAIN
	SymTag_VARIABLE_NAME         SymTag = 1<<4 | SymTag_MAIN
	SymTag_VARIABLE              SymTag = 1<<5 | SymTag_VARIABLE_NAME
	SymTag_STRUCT                SymTag = 1<<6 | SymTag_TYPE | SymTag_VARIABLE_NAME
	SymTag_SERVICE               SymTag = 1<<7 | SymTag_MAIN
	SymTag_INVOKABLE             SymTag = 1 << 8
	SymTag_FUNCTION              SymTag = 1<<9 | SymTag_INVOKABLE | SymTag_VARIABLE
	SymTag_WORKER                SymTag = 1<<10 | SymTag_INVOKABLE | SymTag_MAIN
	SymTag_LISTENER              SymTag = 1<<11 | SymTag_MAIN
	SymTag_PACKAGE               SymTag = 1<<12 | SymTag_IMPORT
	SymTag_XMLNS                 SymTag = 1<<13 | SymTag_IMPORT
	SymTag_ENDPOINT              SymTag = 1<<14 | SymTag_VARIABLE
	SymTag_TYPE_DEF              SymTag = 1<<15 | SymTag_TYPE | SymTag_VARIABLE_NAME
	SymTag_OBJECT                SymTag = 1<<16 | SymTag_TYPE_DEF | SymTag_STRUCT
	SymTag_RECORD                SymTag = 1<<17 | SymTag_TYPE_DEF | SymTag_STRUCT
	SymTag_ERROR                 SymTag = 1<<18 | SymTag_TYPE_DEF
	SymTag_FINITE_TYPE           SymTag = 1<<19 | SymTag_TYPE_DEF
	SymTag_UNION_TYPE            SymTag = 1<<20 | SymTag_TYPE_DEF
	SymTag_INTERSECTION_TYPE     SymTag = 1<<21 | SymTag_TYPE_DEF
	SymTag_TUPLE_TYPE            SymTag = 1<<22 | SymTag_TYPE_DEF
	SymTag_ARRAY_TYPE            SymTag = 1<<23 | SymTag_TYPE_DEF
	SymTag_CONSTANT              SymTag = 1<<24 | SymTag_VARIABLE_NAME | SymTag_TYPE
	SymTag_FUNCTION_TYPE         SymTag = 1<<25 | SymTag_TYPE_DEF
	SymTag_CONSTRUCTOR           SymTag = 1<<26 | SymTag_INVOKABLE
	SymTag_LET                   SymTag = 1 << 27
	SymTag_ENUM                  SymTag = 1<<28 | SymTag_TYPE_DEF
	SymTag_TYPE_REF              SymTag = 1 << 29
	SymTag_ANNOTATION_ATTACHMENT SymTag = 1 << 30
	SymTag_RESOURCE_PATH_SEGMENT SymTag = 1 << 31
	SymTag_SEQUENCE              SymTag = 1<<32 | SymTag_MAIN
)

type Symbol interface {
	GetName() Name
	GetOriginalName() Name
	GetKind() SymbolKind
	GetType() Type
	GetFlags() common.Set[Flag]
	GetEnclosingSymbol() Symbol
	GetEnclosedSymbols() []Symbol
	GetPosition() Location
	GetOrigin() SymbolOrigin
}

type TypeSymbol = Symbol

type Annotatable interface {
	AddAnnotation(AnnotationAttachmentSymbol)
	GetAnnotations() []AnnotationAttachmentSymbol
}

type AnnotationSymbol = Annotatable

type ConstantSymbol = Annotatable

type AnnotationAttachmentSymbol interface {
	IsConstAnnotation() bool
}

type SchedulerPolicy uint8

const (
	SchedulerPolicy_PARENT SchedulerPolicy = iota
	SchedulerPolicy_ANY
)

type InvokableSymbol interface {
	Annotatable
	GetParameters() []VariableSymbol
	GetReturnType() Type
}

type VariableSymbol interface {
	Symbol
	GetConstValue() any
}

type BOperatorSymbol = BInvokableSymbol
type (
	BSymbol struct {
		BLangNode
		Tag                   SymTag
		Flags                 Flags
		Name                  *Name
		OriginalName          *Name
		PkgID                 *PackageID
		Kind                  SymbolKind
		Type                  *BType
		Owner                 Symbol
		Tainted               bool
		Closure               bool
		MarkdownDocumentation *MarkdownDocAttachment
		Pos                   Location
		Origin                SymbolOrigin
	}
	BVarSymbol struct {
		BSymbol
		annotationAttachments []BAnnotationAttachmentSymbol
		IsDefaultable         bool
		IsWildcard            bool
		State                 DiagnosticState
	}
	BConstantSymbol struct {
		BVarSymbol
		Value       *BLangConstantValue
		LiteralType *BType
	}
	BTypeSymbol struct {
		BSymbol
		IsTypeParamResolved bool
		TypeParamTSymbol    *BTypeSymbol
		Annotations         *BVarSymbol
	}
	BAnnotationAttachmentSymbol struct {
		BSymbol
		AnnotPkgID *PackageID
		AnnotTag   *Name
	}
	BAnnotationSymbol struct {
		BTypeSymbol
		AttachedType          *BType
		Points                common.Set[AttachPoint]
		MaskedPoints          int
		annotationAttachments []BAnnotationAttachmentSymbol
	}
	BInvokableSymbol struct {
		BVarSymbol
		Params                          []BVarSymbol
		RestParam                       *BVarSymbol
		RetType                         *BType
		ParamDefaultValTypes            map[string]*BType
		ReceiverSymbol                  *BVarSymbol
		BodyExist                       bool
		annotationAttachmentsOnExternal []BAnnotationAttachmentSymbol
		EnclForkName                    string
		Source                          string
		StrandName                      *string
		SchedulerPolicy                 SchedulerPolicy
		DependentGlobalVars             common.Set[*BVarSymbol]
	}
	BPackageSymbol struct {
		BTypeSymbol
	}
)

var (
	_ Symbol                     = &BSymbol{}
	_ TypeSymbol                 = &BTypeSymbol{}
	_ AnnotationSymbol           = &BAnnotationSymbol{}
	_ AnnotationAttachmentSymbol = &BAnnotationAttachmentSymbol{}
	_ ConstantSymbol             = &BConstantSymbol{}
	_ InvokableSymbol            = &BInvokableSymbol{}
	_ VariableSymbol             = &BVarSymbol{}
	_ Annotatable                = &BVarSymbol{}
)

func (this *BSymbol) GetName() Name {
	return *this.Name
}

func (this *BSymbol) GetOriginalName() Name {
	if this.OriginalName != nil {
		return *this.OriginalName
	}
	return *this.Name
}

func (this *BSymbol) GetKind() SymbolKind {
	return this.Kind
}

func (this *BSymbol) GetType() Type {
	return this.Type
}

func (this *BSymbol) GetFlags() common.Set[Flag] {
	return UnMask(this.Flags)
}

func (this *BSymbol) GetEnclosingSymbol() Symbol {
	return this.Owner
}

func (this *BSymbol) GetEnclosedSymbols() []Symbol {
	// Returns empty slice as per Java implementation
	return []Symbol{}
}

func (this *BSymbol) GetPosition() Location {
	return this.Pos
}

func (this *BSymbol) GetOrigin() SymbolOrigin {
	return this.Origin
}

func (this *BConstantSymbol) GetKind() SymbolKind {
	return SymbolKind_CONSTANT
}

func (this *BConstantSymbol) GetConstValue() any {
	return this.Value
}

func (this *BVarSymbol) AddAnnotation(symbol AnnotationAttachmentSymbol) {
	if symbol == nil {
		return
	}
	if bSymbol, ok := symbol.(*BAnnotationAttachmentSymbol); ok {
		this.annotationAttachments = append(this.annotationAttachments, *bSymbol)
	} else {
		panic("symbol is not a BAnnotationAttachmentSymbol")
	}
}

func (this *BVarSymbol) GetAnnotations() []AnnotationAttachmentSymbol {
	result := make([]AnnotationAttachmentSymbol, len(this.annotationAttachments))
	for i := range this.annotationAttachments {
		result[i] = &this.annotationAttachments[i]
	}
	return result
}

func (this *BVarSymbol) GetConstValue() any {
	return nil
}

func (this *BAnnotationAttachmentSymbol) IsConstAnnotation() bool {
	return false
}

func (this *BAnnotationSymbol) AddAnnotation(symbol AnnotationAttachmentSymbol) {
	if symbol == nil {
		return
	}
	if bSymbol, ok := symbol.(*BAnnotationAttachmentSymbol); ok {
		this.annotationAttachments = append(this.annotationAttachments, *bSymbol)
	} else {
		panic("symbol is not a BAnnotationAttachmentSymbol")
	}
}

func (this *BAnnotationSymbol) GetAnnotations() []AnnotationAttachmentSymbol {
	result := make([]AnnotationAttachmentSymbol, len(this.annotationAttachments))
	for i := range this.annotationAttachments {
		result[i] = &this.annotationAttachments[i]
	}
	return result
}

func (this *BAnnotationSymbol) BvmAlias() string {
	pkg := this.getPackageIDStringWithMajorVersion(this.PkgID)
	if pkg != "." {
		if this.Name != nil {
			return pkg + ":" + string(*this.Name)
		}
	}
	if this.Name != nil {
		return string(*this.Name)
	}
	return ""
}

func (this *BAnnotationSymbol) getMaskedPoints(attachPoints common.Set[AttachPoint]) int {
	points := make(map[Point]bool)
	if attachPoints != nil {
		for ap := range attachPoints.Values() {
			if ap.Point != "" {
				points[ap.Point] = true
			}
		}
	}
	return asMask(points)
}

func asMask(points map[Point]bool) int {
	mask := 0
	for point := range points {
		switch point {
		case Point_TYPE:
			mask |= 1
		case Point_OBJECT:
			mask |= 1 << 1
		case Point_FUNCTION:
			mask |= 1 << 2
		case Point_OBJECT_METHOD:
			mask |= 1 << 3
		case Point_SERVICE_REMOTE:
			mask |= 1 << 4
		case Point_PARAMETER:
			mask |= 1 << 5
		case Point_RETURN:
			mask |= 1 << 6
		case Point_SERVICE:
			mask |= 1 << 7
		case Point_FIELD:
			mask |= 1 << 8
		case Point_OBJECT_FIELD:
			mask |= 1 << 9
		case Point_RECORD_FIELD:
			mask |= 1 << 10
		case Point_LISTENER:
			mask |= 1 << 11
		case Point_ANNOTATION:
			mask |= 1 << 12
		case Point_EXTERNAL:
			mask |= 1 << 13
		case Point_VAR:
			mask |= 1 << 14
		case Point_CONST:
			mask |= 1 << 15
		case Point_WORKER:
			mask |= 1 << 16
		case Point_CLASS:
			mask |= 1 << 17
		}
	}
	return mask
}

func NewBAnnotationSymbol(name *Name, originalName *Name, flags Flags, points common.Set[AttachPoint], pkgID *PackageID, bType *BType, owner Symbol, pos Location, origin SymbolOrigin) *BAnnotationSymbol {
	symbol := &BAnnotationSymbol{
		BTypeSymbol: BTypeSymbol{
			BSymbol: BSymbol{
				BLangNode:    BLangNode{pos: pos},
				Tag:          SymTag_ANNOTATION,
				Flags:        flags,
				Name:         name,
				OriginalName: originalName,
				PkgID:        pkgID,
				Type:         bType,
				Owner:        owner,
				Pos:          pos,
				Origin:       origin,
			},
		},
		AttachedType:          bType,
		Points:                points,
		annotationAttachments: []BAnnotationAttachmentSymbol{},
	}
	symbol.MaskedPoints = symbol.getMaskedPoints(points)
	return symbol
}

func NewBConstantSymbol(flags Flags, name *Name, pkgID *PackageID, literalType *BType, bType *BType, owner Symbol, pos Location, origin SymbolOrigin) *BConstantSymbol {
	return NewBConstantSymbolWithOriginalName(flags, name, name, pkgID, literalType, bType, owner, pos, origin)
}

func NewBConstantSymbolWithOriginalName(flags Flags, name *Name, originalName *Name, pkgID *PackageID, literalType *BType, bType *BType, owner Symbol, pos Location, origin SymbolOrigin) *BConstantSymbol {
	symbol := &BConstantSymbol{
		BVarSymbol: BVarSymbol{
			BSymbol: BSymbol{
				BLangNode:    BLangNode{pos: pos},
				Tag:          SymTag_CONSTANT,
				Flags:        flags,
				Name:         name,
				OriginalName: originalName,
				PkgID:        pkgID,
				Type:         bType,
				Owner:        owner,
				Pos:          pos,
				Origin:       origin,
				Kind:         SymbolKind_CONSTANT,
			},
			annotationAttachments: []BAnnotationAttachmentSymbol{},
			State:                 DiagnosticState_VALID,
		},
		LiteralType: literalType,
	}
	symbol.Kind = SymbolKind_CONSTANT
	return symbol
}

func NewBInvokableSymbol(tag SymTag, flags Flags, name *Name, pkgID *PackageID, bType *BType, owner Symbol, pos Location, origin SymbolOrigin) *BInvokableSymbol {
	return NewBInvokableSymbolWithOriginalName(tag, flags, name, name, pkgID, bType, owner, pos, origin)
}

func NewBInvokableSymbolWithOriginalName(tag SymTag, flags Flags, name *Name, originalName *Name, pkgID *PackageID, bType *BType, owner Symbol, pos Location, origin SymbolOrigin) *BInvokableSymbol {
	symbol := &BInvokableSymbol{
		BVarSymbol: BVarSymbol{
			BSymbol: BSymbol{
				BLangNode:    BLangNode{pos: pos},
				Tag:          tag,
				Flags:        flags,
				Name:         name,
				OriginalName: originalName,
				PkgID:        pkgID,
				Type:         bType,
				Owner:        owner,
				Pos:          pos,
				Origin:       origin,
				Kind:         SymbolKind_FUNCTION,
			},
			annotationAttachments: []BAnnotationAttachmentSymbol{},
			State:                 DiagnosticState_VALID,
		},
		Params:                          []BVarSymbol{},
		ParamDefaultValTypes:            make(map[string]*BType),
		annotationAttachmentsOnExternal: []BAnnotationAttachmentSymbol{},
		SchedulerPolicy:                 SchedulerPolicy_PARENT,
		DependentGlobalVars:             &common.UnorderedSet[*BVarSymbol]{},
	}
	return symbol
}

func (this *BAnnotationSymbol) getPackageIDStringWithMajorVersion(pkgID *PackageID) string {
	if DOT == *pkgID.Name {
		return pkgID.Name.Value()
	}
	org := ""
	if pkgID.OrgName != nil && *pkgID.OrgName != ANON_ORG {
		org = pkgID.OrgName.Value() + string(ORG_NAME_SEPARATOR)
	}
	if *pkgID.Version == EMPTY {
		return org + pkgID.Name.Value()
	}
	return org + pkgID.Name.Value() + string(VERSION_SEPARATOR) + GetMajorVersion(*pkgID.Version)
}

func GetMajorVersion(version Name) string {
	return strings.Split(version.Value(), ".")[0]
}

func (this *BInvokableSymbol) GetParameters() []VariableSymbol {
	result := make([]VariableSymbol, len(this.Params))
	for i := range this.Params {
		result[i] = &this.Params[i]
	}
	return result
}

func (this *BInvokableSymbol) GetReturnType() Type {
	return this.RetType
}

func (this *BInvokableSymbol) SetAnnotationAttachments(annotationAttachments []BAnnotationAttachmentSymbol) {
	this.annotationAttachments = annotationAttachments
}

func (this *BInvokableSymbol) SetAnnotationAttachmentsOnExternal(annotationAttachments []BAnnotationAttachmentSymbol) {
	this.annotationAttachmentsOnExternal = annotationAttachments
}

func (this *BInvokableSymbol) GetAnnotationAttachmentsOnExternal() []AnnotationAttachmentSymbol {
	result := make([]AnnotationAttachmentSymbol, len(this.annotationAttachmentsOnExternal))
	for i := range this.annotationAttachmentsOnExternal {
		result[i] = &this.annotationAttachmentsOnExternal[i]
	}
	return result
}

type MarkdownDocAttachment struct {
	Description             *string
	Parameters              []Parameters
	ReturnValueDescription  *string
	DeprecatedDocumentation *string
	DeprecatedParameters    []Parameters
}

type Parameters struct {
	Name        *string
	Description *string
}
