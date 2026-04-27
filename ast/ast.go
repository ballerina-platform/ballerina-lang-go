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
	"iter"
	"strings"

	"ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type Flags uint64

const (
	Flags_PUBLIC   = 1                 //  0
	Flags_NATIVE   = Flags_PUBLIC << 1 //  1
	Flags_FINAL    = Flags_NATIVE << 1 //  2
	Flags_ATTACHED = Flags_FINAL << 1  //  3

	Flags_DEPRECATED     = Flags_ATTACHED << 1       //  4
	Flags_READONLY       = Flags_DEPRECATED << 1     //  5
	Flags_FUNCTION_FINAL = Flags_READONLY << 1       //  6
	Flags_INTERFACE      = Flags_FUNCTION_FINAL << 1 //  7

	// Marks as a field for which the user MUST provide a value
	Flags_REQUIRED = Flags_INTERFACE << 1 //  8

	Flags_RECORD    = Flags_REQUIRED << 1 //  9
	Flags_PRIVATE   = Flags_RECORD << 1   //  10
	Flags_ANONYMOUS = Flags_PRIVATE << 1  //  11

	Flags_OPTIONAL = Flags_ANONYMOUS << 1 //  12
	Flags_TESTABLE = Flags_OPTIONAL << 1  //  13
	Flags_CONSTANT = Flags_TESTABLE << 1  //  14
	Flags_REMOTE   = Flags_CONSTANT << 1  //  15

	Flags_CLIENT   = Flags_REMOTE << 1   //  16
	Flags_RESOURCE = Flags_CLIENT << 1   //  17
	Flags_SERVICE  = Flags_RESOURCE << 1 //  18
	Flags_LISTENER = Flags_SERVICE << 1  //  19

	Flags_LAMBDA     = Flags_LISTENER << 1   //  20
	Flags_TYPE_PARAM = Flags_LAMBDA << 1     //  21
	Flags_LANG_LIB   = Flags_TYPE_PARAM << 1 //  22
	Flags_WORKER     = Flags_LANG_LIB << 1   //  23

	Flags_FORKED        = Flags_WORKER << 1        //  24
	Flags_TRANSACTIONAL = Flags_FORKED << 1        //  25
	Flags_PARAMETERIZED = Flags_TRANSACTIONAL << 1 //  26
	Flags_DISTINCT      = Flags_PARAMETERIZED << 1 //  27

	Flags_CLASS          = Flags_DISTINCT << 1       //  28
	Flags_ISOLATED       = Flags_CLASS << 1          //  29
	Flags_ISOLATED_PARAM = Flags_ISOLATED << 1       //  30
	Flags_CONFIGURABLE   = Flags_ISOLATED_PARAM << 1 //  31
	Flags_OBJECT_CTOR    = Flags_CONFIGURABLE << 1   //  32

	Flags_ENUM               = Flags_OBJECT_CTOR << 1        //  33
	Flags_INCLUDED           = Flags_ENUM << 1               //  34
	Flags_REQUIRED_PARAM     = Flags_INCLUDED << 1           //  35
	Flags_DEFAULTABLE_PARAM  = Flags_REQUIRED_PARAM << 1     //  36
	Flags_REST_PARAM         = Flags_DEFAULTABLE_PARAM << 1  //  37
	Flags_FIELD              = Flags_REST_PARAM << 1         //  38
	Flags_ANY_FUNCTION       = Flags_FIELD << 1              //  39
	Flags_INFER              = Flags_ANY_FUNCTION << 1       //  40
	Flags_ENUM_MEMBER        = Flags_INFER << 1              //  41
	Flags_QUERY_LAMBDA       = Flags_ENUM_MEMBER << 1        //  42
	Flags_EFFECTIVE_TYPE_DEF = Flags_QUERY_LAMBDA << 1       //  43
	Flags_SOURCE_ANNOTATION  = Flags_EFFECTIVE_TYPE_DEF << 1 //  44
)

func AsMask(flagSet common.Set[model.Flag]) Flags {
	mask := Flags(0)
	for flag := range flagSet.Values() {
		mask |= Flags(flag)
	}
	return mask
}

func flagToFlagsBit(flag model.Flag) Flags {
	switch flag {
	case model.Flag_PUBLIC:
		return Flags_PUBLIC
	case model.Flag_PRIVATE:
		return Flags_PRIVATE
	case model.Flag_REMOTE:
		return Flags_REMOTE
	case model.Flag_TRANSACTIONAL:
		return Flags_TRANSACTIONAL
	case model.Flag_NATIVE:
		return Flags_NATIVE
	case model.Flag_FINAL:
		return Flags_FINAL
	case model.Flag_ATTACHED:
		return Flags_ATTACHED
	case model.Flag_LAMBDA:
		return Flags_LAMBDA
	case model.Flag_WORKER:
		return Flags_WORKER
	case model.Flag_LISTENER:
		return Flags_LISTENER
	case model.Flag_READONLY:
		return Flags_READONLY
	case model.Flag_FUNCTION_FINAL:
		return Flags_FUNCTION_FINAL
	case model.Flag_INTERFACE:
		return Flags_INTERFACE
	case model.Flag_REQUIRED:
		return Flags_REQUIRED
	case model.Flag_RECORD:
		return Flags_RECORD
	case model.Flag_ANONYMOUS:
		return Flags_ANONYMOUS
	case model.Flag_OPTIONAL:
		return Flags_OPTIONAL
	case model.Flag_TESTABLE:
		return Flags_TESTABLE
	case model.Flag_CLIENT:
		return Flags_CLIENT
	case model.Flag_RESOURCE:
		return Flags_RESOURCE
	case model.Flag_ISOLATED:
		return Flags_ISOLATED
	case model.Flag_SERVICE:
		return Flags_SERVICE
	case model.Flag_CONSTANT:
		return Flags_CONSTANT
	case model.Flag_TYPE_PARAM:
		return Flags_TYPE_PARAM
	case model.Flag_LANG_LIB:
		return Flags_LANG_LIB
	case model.Flag_FORKED:
		return Flags_FORKED
	case model.Flag_DISTINCT:
		return Flags_DISTINCT
	case model.Flag_CLASS:
		return Flags_CLASS
	case model.Flag_CONFIGURABLE:
		return Flags_CONFIGURABLE
	case model.Flag_OBJECT_CTOR:
		return Flags_OBJECT_CTOR
	case model.Flag_ENUM:
		return Flags_ENUM
	case model.Flag_INCLUDED:
		return Flags_INCLUDED
	case model.Flag_REQUIRED_PARAM:
		return Flags_REQUIRED_PARAM
	case model.Flag_DEFAULTABLE_PARAM:
		return Flags_DEFAULTABLE_PARAM
	case model.Flag_REST_PARAM:
		return Flags_REST_PARAM
	case model.Flag_FIELD:
		return Flags_FIELD
	case model.Flag_ANY_FUNCTION:
		return Flags_ANY_FUNCTION
	case model.Flag_ENUM_MEMBER:
		return Flags_ENUM_MEMBER
	case model.Flag_QUERY_LAMBDA:
		return Flags_QUERY_LAMBDA
	default:
		return 0
	}
}

func UnMask(mask Flags) common.Set[model.Flag] {
	flagSet := common.UnorderedSet[model.Flag]{}
	for flag := model.Flag_PUBLIC; flag <= model.Flag_QUERY_LAMBDA; flag++ {
		flagVal := flagToFlagsBit(flag)
		if flagVal != 0 && (mask&flagVal) == flagVal {
			flagSet.Add(flag)
		}
	}
	return &flagSet
}

type BNodeWithSymbol interface {
	model.NodeWithSymbol
	BLangNode
	SetSymbol(symbolRef model.SymbolRef)
}

// SymbolIsSet returns true if the AST node has its symbol set.
func SymbolIsSet(node model.NodeWithSymbol) bool {
	return node.Symbol() != (model.SymbolRef{})
}

type NodeWithScope interface {
	Scope() model.Scope
	SetScope(scope model.Scope)
}

type SourceKind = model.SourceKind

type CompilerPhase uint8

const (
	CompilerPhase_DEFINE CompilerPhase = iota
	CompilerPhase_TYPE_CHECK
	CompilerPhase_CODE_ANALYZE
	CompilerPhase_DATAFLOW_ANALYZE
	CompilerPhase_ISOLATION_ANALYZE
	CompilerPhase_DOCUMENTATION_ANALYZE
	CompilerPhase_CONSTANT_PROPAGATION
	CompilerPhase_COMPILER_PLUGIN
	CompilerPhase_DESUGAR
	CompilerPhase_BIR_GEN
	CompilerPhase_BIR_EMIT
	CompilerPhase_CODE_GEN
)

type BLangNode interface {
	model.Node
	SetDeterminedType(ty semtypes.SemType)
	SetPosition(pos diagnostics.Location)
}

type (
	bLangNodeBase struct {
		DeterminedType semtypes.SemType
		parent         BLangNode
		pos            diagnostics.Location
	}

	BLangAnnotation struct {
		bLangNodeBase
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		typeDescriptor                  model.TypeDescriptor
		attachPoints                    common.UnorderedSet[model.AttachPoint]
	}

	BLangAnnotationAttachment struct {
		bLangNodeBase
		Expr           BLangExpression
		AnnotationName *BLangIdentifier
		PkgAlias       *BLangIdentifier
		AttachPoints   common.OrderedSet[model.Point]
	}

	bLangFunctionBodyBase struct {
		bLangNodeBase
	}

	BLangBlockFunctionBody struct {
		bLangFunctionBodyBase
		Stmts []BLangStatement
	}

	BLangExprFunctionBody struct {
		bLangFunctionBodyBase
		Expr model.ExpressionNode
	}

	BLangExternFunctionBody struct {
		bLangFunctionBodyBase
		AnnAttachments []BLangAnnotationAttachment
	}

	BLangIdentifier struct {
		bLangNodeBase
		Value         string
		OriginalValue string
		isLiteral     bool
	}

	BLangImportPackage struct {
		bLangNodeBase
		OrgName      *BLangIdentifier
		PkgNameComps []BLangIdentifier
		Alias        *BLangIdentifier
		CompUnit     *BLangIdentifier
		Version      *BLangIdentifier
	}

	BLangClassDefinition struct {
		bLangNodeBase
		Name                            *BLangIdentifier
		symbol                          model.SymbolRef
		scope                           model.Scope
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		InitFunction                    *BLangFunction
		Methods                         map[string]*BLangFunction
		Fields                          []model.SimpleVariableNode
		Inclusions                      []model.SymbolRef       // This needs to be symbol because it could be a class definition as well
		unresolvedInclusions            []*BLangUserDefinedType // we need this because we can't get symbols before the symbol resolution in node_builder. After symbol resolution this field is cleared out
		flags                           nodeFlags
		typeData                        model.TypeData
		Definition                      semtypes.Definition
		CycleDepth                      int
		precedence                      int
	}

	BLangService struct {
		bLangNodeBase
		symbol                          model.SymbolRef
		ServiceVariable                 *BLangSimpleVariable
		AttachedExprs                   []BLangExpression
		ServiceClass                    *BLangClassDefinition
		AbsoluteResourcePath            []model.IdentifierNode
		ServiceNameLiteral              *BLangLiteral
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		ListenerType                    BType
		ResourceFunctions               []BLangFunction
		InferredServiceType             BType
	}

	BLangCompilationUnit struct {
		bLangNodeBase
		TopLevelNodes []model.TopLevelNode
		Name          string
		packageID     *model.PackageID
		sourceKind    SourceKind
	}

	BLangPackage struct {
		bLangNodeBase
		CompUnits               []BLangCompilationUnit
		Imports                 []BLangImportPackage
		XmlnsList               []BLangXMLNS
		Constants               []BLangConstant
		GlobalVars              []BLangSimpleVariable
		Services                []BLangService
		Functions               []BLangFunction
		TypeDefinitions         []BLangTypeDefinition
		Annotations             []BLangAnnotation
		InitFunction            *BLangFunction
		StartFunction           *BLangFunction
		StopFunction            *BLangFunction
		TopLevelNodes           []model.TopLevelNode
		TestablePkgs            []*BLangTestablePackage
		ClassDefinitions        []BLangClassDefinition
		CompletedPhases         common.UnorderedSet[CompilerPhase]
		LambdaFunctions         []BLangLambdaFunction
		PackageID               *model.PackageID
		Scope                   model.Scope
		diagnostics             []diagnostics.Diagnostic
		ModuleContextDataHolder *ModuleContextDataHolder
		errorCount              int
		warnCount               int
	}
	BLangTestablePackage struct {
		BLangPackage
		Parent               *BLangPackage
		mockFunctionNamesMap map[string]string
		isLegacyMockingMap   map[string]bool
	}
	BLangXMLNS struct {
		bLangNodeBase
		namespaceURI BLangExpression
		prefix       *BLangIdentifier
		CompUnit     *BLangIdentifier
	}
	BLangLocalXMLNS struct {
		BLangXMLNS
	}
	BLangPackageXMLNS struct {
		BLangXMLNS
	}
	BLangMarkdownDocumentation struct {
		bLangNodeBase
		DocumentationLines                []BLangMarkdownDocumentationLine
		Parameters                        []BLangMarkdownParameterDocumentation
		References                        []BLangMarkdownReferenceDocumentation
		ReturnParameter                   *BLangMarkdownReturnParameterDocumentation
		DeprecationDocumentation          *BLangMarkDownDeprecationDocumentation
		DeprecatedParametersDocumentation *BLangMarkDownDeprecatedParametersDocumentation
	}
	BLangMarkdownReferenceDocumentation struct {
		bLangNodeBase
		Qualifier         string
		TypeName          string
		Identifier        string
		ReferenceName     string
		Type              model.DocumentationReferenceType
		HasParserWarnings bool
	}

	BLangVariableBase struct {
		bLangNodeBase
		// We are using variable for function paramets and record td fields so we need to have
		// type descriptors here. Not sure this is the best way to do this.
		typeNode                        BType
		AnnAttachments                  []model.AnnotationAttachmentNode
		MarkdownDocumentationAttachment model.MarkdownDocumentationNode
		Expr                            model.ExpressionNode
		flags                           nodeFlags
		IsDeclaredWithVar               bool
		symbol                          model.SymbolRef
	}

	BLangConstant struct {
		BLangVariableBase
		Name *BLangIdentifier
	}

	BLangSimpleVariable struct {
		BLangVariableBase
		Name *BLangIdentifier
	}

	ClosureVarSymbol struct {
		DiagnosticLocation diagnostics.Location
	}

	bLangInvokableNodeBase struct {
		bLangNodeBase
		Name                            BLangIdentifier
		symbol                          model.SymbolRef
		AnnAttachments                  []model.AnnotationAttachmentNode
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		RequiredParams                  []BLangSimpleVariable
		RestParam                       model.SimpleVariableNode
		returnTypeDescriptor            model.TypeDescriptor
		Body                            model.FunctionBodyNode
		flags                           nodeFlags
	}

	BLangFunction struct {
		bLangInvokableNodeBase
		scope model.Scope
	}

	nodeFlags uint64

	BLangTypeDefinition struct {
		bLangNodeBase
		Name                            *BLangIdentifier
		symbol                          model.SymbolRef
		typeData                        model.TypeData
		annAttachments                  []BLangAnnotationAttachment
		markdownDocumentationAttachment *BLangMarkdownDocumentation
		flags                           nodeFlags
		precedence                      int
		CycleDepth                      int
		isBuiltinTypeDef                bool
		hasCyclicReference              bool
		referencedFieldsDefined         bool
	}
)

// Flag bit constants — bit positions match model.Flag iota values.
const (
	flagPublic        nodeFlags = 1 << model.Flag_PUBLIC
	flagPrivate       nodeFlags = 1 << model.Flag_PRIVATE
	flagRemote        nodeFlags = 1 << model.Flag_REMOTE
	flagTransactional nodeFlags = 1 << model.Flag_TRANSACTIONAL
	flagNative        nodeFlags = 1 << model.Flag_NATIVE
	flagFinal         nodeFlags = 1 << model.Flag_FINAL
	flagAttached      nodeFlags = 1 << model.Flag_ATTACHED
	flagLambda        nodeFlags = 1 << model.Flag_LAMBDA
	flagReadonly      nodeFlags = 1 << model.Flag_READONLY
	flagInterface     nodeFlags = 1 << model.Flag_INTERFACE
	flagAnonymous     nodeFlags = 1 << model.Flag_ANONYMOUS
	flagOptional      nodeFlags = 1 << model.Flag_OPTIONAL
	flagClient        nodeFlags = 1 << model.Flag_CLIENT
	flagResource      nodeFlags = 1 << model.Flag_RESOURCE
	flagIsolated      nodeFlags = 1 << model.Flag_ISOLATED
	flagService       nodeFlags = 1 << model.Flag_SERVICE
	flagConstant      nodeFlags = 1 << model.Flag_CONSTANT
	flagDistinct      nodeFlags = 1 << model.Flag_DISTINCT
	flagClass         nodeFlags = 1 << model.Flag_CLASS
	flagRequiredParam nodeFlags = 1 << model.Flag_REQUIRED_PARAM
	flagDefaultParam  nodeFlags = 1 << model.Flag_DEFAULTABLE_PARAM
	flagRestParam     nodeFlags = 1 << model.Flag_REST_PARAM
	flagAnyFunction   nodeFlags = 1 << model.Flag_ANY_FUNCTION
)

func (f nodeFlags) has(flag nodeFlags) bool { return f&flag != 0 }
func (f nodeFlags) asInt64() int64          { return int64(f) }

// bLangInvokableNodeBase flag methods
func (b *bLangInvokableNodeBase) IsPublic() bool        { return b.flags.has(flagPublic) }
func (b *bLangInvokableNodeBase) IsRemote() bool        { return b.flags.has(flagRemote) }
func (b *bLangInvokableNodeBase) IsTransactional() bool { return b.flags.has(flagTransactional) }
func (b *bLangInvokableNodeBase) IsResource() bool      { return b.flags.has(flagResource) }
func (b *bLangInvokableNodeBase) IsIsolated() bool      { return b.flags.has(flagIsolated) }
func (b *bLangInvokableNodeBase) IsInterface() bool     { return b.flags.has(flagInterface) }
func (b *bLangInvokableNodeBase) IsNative() bool        { return b.flags.has(flagNative) }
func (b *bLangInvokableNodeBase) IsAnonymous() bool     { return b.flags.has(flagLambda) }
func (b *bLangInvokableNodeBase) IsAttached() bool      { return b.flags.has(flagAttached) }

func (b *bLangInvokableNodeBase) SetPublic()          { b.flags |= flagPublic }
func (b *bLangInvokableNodeBase) SetRemote()          { b.flags |= flagRemote }
func (b *bLangInvokableNodeBase) SetTransactional()   { b.flags |= flagTransactional }
func (b *bLangInvokableNodeBase) SetResource()        { b.flags |= flagResource }
func (b *bLangInvokableNodeBase) SetIsolated()        { b.flags |= flagIsolated }
func (b *bLangInvokableNodeBase) SetInterface()       { b.flags |= flagInterface }
func (b *bLangInvokableNodeBase) SetNative()          { b.flags |= flagNative }
func (b *bLangInvokableNodeBase) SetAnonymous()       { b.flags |= flagLambda | flagAnonymous }
func (b *bLangInvokableNodeBase) SetAttached()        { b.flags |= flagAttached }
func (b *bLangInvokableNodeBase) FlagsAsInt64() int64 { return b.flags.asInt64() }

// BLangVariableBase flag methods
func (b *BLangVariableBase) IsPublic() bool           { return b.flags.has(flagPublic) }
func (b *BLangVariableBase) IsFinal() bool            { return b.flags.has(flagFinal) }
func (b *BLangVariableBase) IsDefaultableParam() bool { return b.flags.has(flagDefaultParam) }
func (b *BLangVariableBase) IsRequiredParam() bool    { return b.flags.has(flagRequiredParam) }
func (b *BLangVariableBase) IsRestParam() bool        { return b.flags.has(flagRestParam) }

func (b *BLangVariableBase) SetPublic()           { b.flags |= flagPublic }
func (b *BLangVariableBase) SetPrivate()          { b.flags &^= flagPublic }
func (b *BLangVariableBase) SetFinal()            { b.flags |= flagFinal }
func (b *BLangVariableBase) SetIsolated()         { b.flags |= flagIsolated }
func (b *BLangVariableBase) SetDefaultableParam() { b.flags |= flagDefaultParam }
func (b *BLangVariableBase) SetRequiredParam()    { b.flags |= flagRequiredParam }
func (b *BLangVariableBase) SetRestParam()        { b.flags |= flagRestParam }
func (b *BLangVariableBase) IsReadonly() bool     { return b.flags.has(flagReadonly) }
func (b *BLangVariableBase) FlagsAsInt64() int64  { return b.flags.asInt64() }

func (b *BLangConstant) FlagsAsInt64() int64 { return (b.flags | flagConstant).asInt64() }

// BLangClassDefinition flag methods
func (b *BLangClassDefinition) IsPublic() bool   { return b.flags.has(flagPublic) }
func (b *BLangClassDefinition) IsDistinct() bool { return b.flags.has(flagDistinct) }
func (b *BLangClassDefinition) IsClient() bool   { return b.flags.has(flagClient) }
func (b *BLangClassDefinition) IsReadonly() bool { return b.flags.has(flagReadonly) }
func (b *BLangClassDefinition) IsService() bool  { return b.flags.has(flagService) }
func (b *BLangClassDefinition) IsIsolated() bool { return b.flags.has(flagIsolated) }

func (b *BLangClassDefinition) SetPublic()          { b.flags |= flagPublic }
func (b *BLangClassDefinition) SetDistinct()        { b.flags |= flagDistinct }
func (b *BLangClassDefinition) SetClient()          { b.flags |= flagClient }
func (b *BLangClassDefinition) SetReadonly()        { b.flags |= flagReadonly }
func (b *BLangClassDefinition) SetService()         { b.flags |= flagService }
func (b *BLangClassDefinition) SetIsolated()        { b.flags |= flagIsolated }
func (b *BLangClassDefinition) SetClass()           { b.flags |= flagClass }
func (b *BLangClassDefinition) FlagsAsInt64() int64 { return b.flags.asInt64() }

// BLangTypeDefinition flag methods
func (b *BLangTypeDefinition) IsPublic() bool    { return b.flags.has(flagPublic) }
func (b *BLangTypeDefinition) IsAnonymous() bool { return b.flags.has(flagAnonymous) }
func (b *BLangTypeDefinition) SetPublic()        { b.flags |= flagPublic }
func (b *BLangTypeDefinition) SetAnonymous()     { b.flags |= flagAnonymous }

// Stub IsPublic for types with no flags
func (b *BLangAnnotation) IsPublic() bool     { return false }
func (b *BLangService) IsPublic() bool        { return false }
func (b *BLangInvocation) IsPublic() bool     { return false }
func (b *BLangMemberTypeDesc) IsPublic() bool { return false }

func (b *bLangNodeBase) SetDeterminedType(ty semtypes.SemType) {
	b.DeterminedType = ty
}

func (b *bLangNodeBase) GetDeterminedType() semtypes.SemType {
	return b.DeterminedType
}

func (b *bLangNodeBase) GetPosition() diagnostics.Location {
	return b.pos
}

func (b *bLangNodeBase) SetPosition(pos diagnostics.Location) {
	b.pos = pos
}

func (n *BLangClassDefinition) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *BLangClassDefinition) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

func (n *BLangClassDefinition) Scope() model.Scope {
	return n.scope
}

func (n *BLangClassDefinition) SetScope(scope model.Scope) {
	n.scope = scope
}

func (n *BLangService) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *BLangService) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

func (n *BLangVariableBase) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *BLangVariableBase) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

func (n *BLangVariableBase) TypeNode() BType {
	return n.typeNode
}

func (n *BLangVariableBase) SetTypeNode(bt BType) {
	n.typeNode = bt
}

func (n *bLangInvokableNodeBase) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *bLangInvokableNodeBase) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

func (n *BLangTypeDefinition) Symbol() model.SymbolRef {
	return n.symbol
}

func (n *BLangTypeDefinition) SetSymbol(symbolRef model.SymbolRef) {
	n.symbol = symbolRef
}

var (
	_ model.AnnotationAttachmentNode                    = &BLangAnnotationAttachment{}
	_ model.IdentifierNode                              = &BLangIdentifier{}
	_ model.ImportPackageNode                           = &BLangImportPackage{}
	_ model.ClassDefinition                             = &BLangClassDefinition{}
	_ model.TypeDefinition                              = &BLangClassDefinition{}
	_ NodeWithScope                                     = &BLangClassDefinition{}
	_ model.PackageNode                                 = &BLangPackage{}
	_ model.PackageNode                                 = &BLangTestablePackage{}
	_ model.AnnotationNode                              = &BLangAnnotation{}
	_ model.XMLNSDeclarationNode                        = &BLangXMLNS{}
	_ model.ServiceNode                                 = &BLangService{}
	_ model.CompilationUnitNode                         = &BLangCompilationUnit{}
	_ model.ConstantNode                                = &BLangConstant{}
	_ model.TypeDefinition                              = &BLangTypeDefinition{}
	_ model.SimpleVariableNode                          = &BLangSimpleVariable{}
	_ model.MarkdownDocumentationNode                   = &BLangMarkdownDocumentation{}
	_ model.MarkdownDocumentationReferenceAttributeNode = &BLangMarkdownReferenceDocumentation{}
	_ model.ExprFunctionBodyNode                        = &BLangExprFunctionBody{}
	_ model.FunctionNode                                = &BLangFunction{}
	_ model.FunctionBodyNode                            = &BLangExternFunctionBody{}
)

var (
	_ BLangNode = &BLangAnnotation{}
	_ BLangNode = &BLangAnnotationAttachment{}
	_ BLangNode = &BLangBlockFunctionBody{}
	_ BLangNode = &BLangExprFunctionBody{}
	_ BLangNode = &BLangIdentifier{}
	_ BLangNode = &BLangImportPackage{}
	_ BLangNode = &BLangClassDefinition{}
	_ BLangNode = &BLangService{}
	_ BLangNode = &BLangCompilationUnit{}
	_ BLangNode = &BLangPackage{}
	_ BLangNode = &BLangTestablePackage{}
	_ BLangNode = &BLangXMLNS{}
	_ BLangNode = &BLangLocalXMLNS{}
	_ BLangNode = &BLangPackageXMLNS{}
	_ BLangNode = &BLangMarkdownDocumentation{}
	_ BLangNode = &BLangMarkdownReferenceDocumentation{}
	_ BLangNode = &BLangConstant{}
	_ BLangNode = &BLangSimpleVariable{}
	_ BLangNode = &BLangFunction{}
	_ BLangNode = &BLangTypeDefinition{}
)

var (
	// Assert that concrete types with symbols implement BNodeWithSymbol
	_ BNodeWithSymbol = &BLangClassDefinition{}
	_ BNodeWithSymbol = &BLangService{}
	_ BNodeWithSymbol = &BLangConstant{}
	_ BNodeWithSymbol = &BLangSimpleVariable{}
	_ BNodeWithSymbol = &BLangFunction{}
	_ BNodeWithSymbol = &BLangTypeDefinition{}
)

func (b *BLangAnnotationAttachment) GetKind() model.NodeKind {
	// migrated from BLangAnnotationAttachment.java:89:5
	return model.NodeKind_ANNOTATION_ATTACHMENT
}

func (b *BLangAnnotationAttachment) GetPackageAlias() model.IdentifierNode {
	return b.PkgAlias
}

func (b *BLangAnnotationAttachment) SetPackageAlias(pkgAlias model.IdentifierNode) {
	if id, ok := pkgAlias.(*BLangIdentifier); ok {
		b.PkgAlias = id
	} else {
		panic("pkgAlias is not a BLangIdentifier")
	}
}

func (b *BLangAnnotationAttachment) GetAnnotationName() model.IdentifierNode {
	return b.AnnotationName
}

func (b *BLangAnnotationAttachment) SetAnnotationName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.AnnotationName = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangAnnotationAttachment) GetExpressionNode() model.ExpressionNode {
	return b.Expr
}

func (b *BLangAnnotationAttachment) SetExpressionNode(expr model.ExpressionNode) {
	if expr, ok := expr.(BLangExpression); ok {
		b.Expr = expr
	} else {
		panic("expr is not a BLangExpression")
	}
}

func (b *BLangAnnotation) GetKind() model.NodeKind {
	// migrated from BLangAnnotation.java:135:5
	return model.NodeKind_ANNOTATION
}

func (b *BLangAnnotation) GetName() model.IdentifierNode {
	// migrated from BLangAnnotation.java:80:5
	return b.Name
}

func (b *BLangAnnotation) SetName(name model.IdentifierNode) {
	// migrated from BLangAnnotation.java:85:5
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (b *BLangAnnotation) GetTypeDescriptor() model.TypeDescriptor {
	return b.typeDescriptor
}

func (b *BLangAnnotation) SetTypeDescriptor(typeDescriptor model.TypeDescriptor) {
	b.typeDescriptor = typeDescriptor
}

func (b *BLangAnnotation) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangAnnotation.java:100:5
	attachments := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	for i, attachment := range b.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (b *BLangAnnotation) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangAnnotation.java:105:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.AnnAttachments = append(b.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (b *BLangAnnotation) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangAnnotation.java:110:5
	return b.MarkdownDocumentationAttachment
}

func (b *BLangAnnotation) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangAnnotation.java:115:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (b *BLangBlockFunctionBody) GetKind() model.NodeKind {
	// migrated from BLangBlockFunctionBody.java:73:5
	return model.NodeKind_BLOCK_FUNCTION_BODY
}

func (b *BLangExternFunctionBody) GetKind() model.NodeKind {
	return model.NodeKind_EXTERN_FUNCTION_BODY
}

func (b *BLangExprFunctionBody) GetKind() model.NodeKind {
	// migrated from BLangExprFunctionBody.java:50:5
	return model.NodeKind_EXPR_FUNCTION_BODY
}

func (b *BLangExprFunctionBody) GetExpr() model.ExpressionNode {
	// migrated from BLangExprFunctionBody.java:55:5
	return b.Expr
}

func (b *BLangIdentifier) GetValue() string {
	// migrated from BLangIdentifier.java:32:5
	return b.Value
}

func (b *BLangIdentifier) GetKind() model.NodeKind {
	// migrated from BLangIdentifier.java:32:5
	return model.NodeKind_IDENTIFIER
}

func (b *BLangIdentifier) SetValue(value string) {
	// migrated from BLangIdentifier.java:37:5
	b.Value = value
}

func (b *BLangIdentifier) SetOriginalValue(value string) {
	// migrated from BLangIdentifier.java:42:5
	b.OriginalValue = value
}

func (b *BLangIdentifier) IsLiteral() bool {
	// migrated from BLangIdentifier.java:47:5
	return b.isLiteral
}

func (b *BLangIdentifier) SetLiteral(isLiteral bool) {
	// migrated from BLangIdentifier.java:52:5
	b.isLiteral = isLiteral
}

func (b *BLangImportPackage) GetKind() model.NodeKind {
	return model.NodeKind_IMPORT
}

func (b *BLangImportPackage) GetOrgName() model.IdentifierNode {
	return b.OrgName
}

func (b *BLangImportPackage) GetPackageName() []model.IdentifierNode {
	result := make([]model.IdentifierNode, len(b.PkgNameComps))
	for i := range b.PkgNameComps {
		result[i] = &b.PkgNameComps[i]
	}
	return result
}

func (b *BLangImportPackage) SetPackageName(nameParts []model.IdentifierNode) {
	b.PkgNameComps = make([]BLangIdentifier, 0, len(nameParts))
	for _, namePart := range nameParts {
		if id, ok := namePart.(*BLangIdentifier); ok {
			b.PkgNameComps = append(b.PkgNameComps, *id)
		} else {
			panic("namePart is not a BLangIdentifier")
		}
	}
}

func (b *BLangImportPackage) GetPackageVersion() model.IdentifierNode {
	return b.Version
}

func (b *BLangImportPackage) SetPackageVersion(version model.IdentifierNode) {
	if id, ok := version.(*BLangIdentifier); ok {
		b.Version = id
	} else {
		panic("version is not a BLangIdentifier")
	}
}

func (b *BLangImportPackage) GetAlias() model.IdentifierNode {
	return b.Alias
}

func (b *BLangImportPackage) SetAlias(alias model.IdentifierNode) {
	if id, ok := alias.(*BLangIdentifier); ok {
		b.Alias = id
	} else {
		panic("alias is not a BLangIdentifier")
	}
}

func NewBLangClassDefinition() BLangClassDefinition {
	this := BLangClassDefinition{}
	this.CycleDepth = -1
	this.Methods = map[string]*BLangFunction{}
	this.SetClass()
	return this
}

func (b *BLangClassDefinition) PopUnresolvedInclusions() []*BLangUserDefinedType {
	inclusions := b.unresolvedInclusions
	b.unresolvedInclusions = nil
	return inclusions
}

func (b *BLangClassDefinition) GetName() model.IdentifierNode {
	// migrated from BLangClassDefinition.java:88:5
	return b.Name
}

func (b *BLangClassDefinition) SetName(name model.IdentifierNode) {
	// migrated from BLangClassDefinition.java:93:5
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (b *BLangClassDefinition) GetMethods() iter.Seq2[string, model.FunctionNode] {
	return func(yield func(string, model.FunctionNode) bool) {
		for name, method := range b.Methods {
			if !yield(name, method) {
				return
			}
		}
	}
}

func (b *BLangClassDefinition) GetMethod(name string) model.FunctionNode {
	if method, ok := b.Methods[name]; ok {
		return method
	}
	return nil
}

func (b *BLangClassDefinition) AddMethod(name string, function *BLangFunction) {
	if b.Methods == nil {
		b.Methods = map[string]*BLangFunction{}
	}
	b.Methods[name] = function
}

func (b *BLangClassDefinition) GetInitFunction() model.FunctionNode {
	// migrated from BLangClassDefinition.java:108:5
	return b.InitFunction
}

func (b *BLangClassDefinition) AddField(field model.VariableNode) {
	// migrated from BLangClassDefinition.java:113:5
	if field, ok := field.(*BLangSimpleVariable); ok {
		b.Fields = append(b.Fields, field)
		return
	}
	panic("field is not a BLangSimpleVariable")
}

func (b *BLangClassDefinition) AddInclusion(symbolRef model.SymbolRef) {
	b.Inclusions = append(b.Inclusions, symbolRef)
}

func (b *BLangClassDefinition) GetKind() model.NodeKind {
	// migrated from BLangClassDefinition.java:138:5
	return model.NodeKind_CLASS_DEFN
}

func (b *BLangClassDefinition) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangClassDefinition.java:168:5
	attachments := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	for i, attachment := range b.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (b *BLangClassDefinition) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangClassDefinition.java:173:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.AnnAttachments = append(b.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (b *BLangClassDefinition) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangClassDefinition.java:178:5
	return b.MarkdownDocumentationAttachment
}

func (b *BLangClassDefinition) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangClassDefinition.java:183:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (b *BLangClassDefinition) GetPrecedence() int {
	return b.precedence
}

func (b *BLangClassDefinition) SetPrecedence(precedence int) {
	b.precedence = precedence
}

func (b *BLangClassDefinition) GetTypeData() model.TypeData {
	return b.typeData
}

func (b *BLangClassDefinition) SetTypeData(typeData model.TypeData) {
	b.typeData = typeData
}

func (b *BLangClassDefinition) GetCycleDepth() int {
	return b.CycleDepth
}

func (b *BLangClassDefinition) SetCycleDepth(depth int) {
	b.CycleDepth = depth
}

func (b *BLangCompilationUnit) AddTopLevelNode(node model.TopLevelNode) {
	// migrated from BLangCompilationUnit.java:48:5
	b.TopLevelNodes = append(b.TopLevelNodes, node)
}

func (b *BLangCompilationUnit) GetTopLevelNodes() []model.TopLevelNode {
	// migrated from BLangCompilationUnit.java:53:5
	return b.TopLevelNodes
}

func (b *BLangCompilationUnit) GetName() string {
	// migrated from BLangCompilationUnit.java:58:5
	return b.Name
}

func (b *BLangCompilationUnit) SetName(name string) {
	// migrated from BLangCompilationUnit.java:63:5
	b.Name = name
}

func (b *BLangCompilationUnit) GetPackageID() *model.PackageID {
	// migrated from BLangCompilationUnit.java:68:5
	return b.packageID
}

func (b *BLangCompilationUnit) SetPackageID(packageID *model.PackageID) {
	// migrated from BLangCompilationUnit.java:72:5
	b.packageID = packageID
}

func (b *BLangCompilationUnit) GetKind() model.NodeKind {
	// migrated from BLangCompilationUnit.java:76:5
	return model.NodeKind_COMPILATION_UNIT
}

func (b *BLangCompilationUnit) SetSourceKind(kind SourceKind) {
	// migrated from BLangCompilationUnit.java:81:5
	b.sourceKind = kind
}

func (b *BLangCompilationUnit) GetSourceKind() SourceKind {
	// migrated from BLangCompilationUnit.java:86:5
	return b.sourceKind
}

func (b *BLangConstant) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangConstant) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (b *BLangConstant) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangConstant.java:88:5
	return b.AnnAttachments
}

func (b *BLangConstant) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangConstant.java:93:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.AnnAttachments = append(b.AnnAttachments, annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (b *BLangConstant) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangConstant.java:98:5
	return b.MarkdownDocumentationAttachment
}

func (b *BLangConstant) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangConstant.java:103:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (b *BLangConstant) GetKind() model.NodeKind {
	// migrated from BLangConstant.java:108:5
	return model.NodeKind_CONSTANT
}

func (b *BLangConstant) GetAssociatedType() semtypes.SemType {
	if b.TypeNode() != nil {
		return b.TypeNode().GetTypeData().Type
	}
	return nil
}

func (b *BLangConstant) GetPrecedence() int {
	// migrated from BLangConstant.java:144:5
	return 0
}

func (b *BLangConstant) SetPrecedence(precedence int) {
	// migrated from BLangConstant.java:149:5
}

func (b *BLangSimpleVariable) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangSimpleVariable) GetKind() model.NodeKind {
	return model.NodeKind_VARIABLE
}

func (b *BLangSimpleVariable) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (b *BLangMarkdownDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_MARKDOWN_DOCUMENTATION
}

func (b *BLangMarkdownDocumentation) GetDocumentationLines() []model.MarkdownDocumentationTextAttributeNode {
	result := make([]model.MarkdownDocumentationTextAttributeNode, len(b.DocumentationLines))
	for i := range b.DocumentationLines {
		result[i] = &b.DocumentationLines[i]
	}
	return result
}

func (b *BLangMarkdownDocumentation) AddDocumentationLine(documentationText model.MarkdownDocumentationTextAttributeNode) {
	if line, ok := documentationText.(*BLangMarkdownDocumentationLine); ok {
		b.DocumentationLines = append(b.DocumentationLines, *line)
	} else {
		panic("documentationText is not a BLangMarkdownDocumentationLine")
	}
}

func (b *BLangMarkdownDocumentation) GetParameters() []model.MarkdownDocumentationParameterAttributeNode {
	result := make([]model.MarkdownDocumentationParameterAttributeNode, len(b.Parameters))
	for i := range b.Parameters {
		result[i] = &b.Parameters[i]
	}
	return result
}

func (b *BLangMarkdownDocumentation) AddParameter(parameter model.MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		b.Parameters = append(b.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (b *BLangMarkdownDocumentation) GetReturnParameter() model.MarkdownDocumentationReturnParameterAttributeNode {
	return b.ReturnParameter
}

func (b *BLangMarkdownDocumentation) GetDeprecationDocumentation() model.MarkDownDocumentationDeprecationAttributeNode {
	return b.DeprecationDocumentation
}

func (b *BLangMarkdownDocumentation) SetReturnParameter(returnParameter model.MarkdownDocumentationReturnParameterAttributeNode) {
	if param, ok := returnParameter.(*BLangMarkdownReturnParameterDocumentation); ok {
		b.ReturnParameter = param
	} else {
		panic("returnParameter is not a BLangMarkdownReturnParameterDocumentation")
	}
}

func (b *BLangMarkdownDocumentation) SetDeprecationDocumentation(deprecationDocumentation model.MarkDownDocumentationDeprecationAttributeNode) {
	if doc, ok := deprecationDocumentation.(*BLangMarkDownDeprecationDocumentation); ok {
		b.DeprecationDocumentation = doc
	} else {
		panic("deprecationDocumentation is not a BLangMarkDownDeprecationDocumentation")
	}
}

func (b *BLangMarkdownDocumentation) SetDeprecatedParametersDocumentation(deprecatedParametersDocumentation model.MarkDownDocumentationDeprecatedParametersAttributeNode) {
	if doc, ok := deprecatedParametersDocumentation.(*BLangMarkDownDeprecatedParametersDocumentation); ok {
		b.DeprecatedParametersDocumentation = doc
	} else {
		panic("deprecatedParametersDocumentation is not a BLangMarkDownDeprecatedParametersDocumentation")
	}
}

func (b *BLangMarkdownDocumentation) GetDeprecatedParametersDocumentation() model.MarkDownDocumentationDeprecatedParametersAttributeNode {
	return b.DeprecatedParametersDocumentation
}

func (b *BLangMarkdownDocumentation) GetDocumentation() string {
	var lines []string
	for i := range b.DocumentationLines {
		lines = append(lines, b.DocumentationLines[i].GetText())
	}
	result := strings.Join(lines, "\n")
	return strings.ReplaceAll(result, "\r", "")
}

func (b *BLangMarkdownDocumentation) GetParameterDocumentations() map[string]model.MarkdownDocumentationParameterAttributeNode {
	result := make(map[string]model.MarkdownDocumentationParameterAttributeNode)
	for _, parameter := range b.Parameters {
		paramName := parameter.GetParameterName()
		result[paramName.GetValue()] = &parameter
	}
	return result
}

func (b *BLangMarkdownDocumentation) GetReturnParameterDocumentation() *string {
	if b.ReturnParameter == nil {
		return nil
	}
	return new(b.ReturnParameter.GetReturnParameterDocumentation())
}

func (b *BLangMarkdownDocumentation) GetReferences() []model.MarkdownDocumentationReferenceAttributeNode {
	result := make([]model.MarkdownDocumentationReferenceAttributeNode, len(b.References))
	for i := range b.References {
		result[i] = &b.References[i]
	}
	return result
}

func (b *BLangMarkdownDocumentation) AddReference(reference model.MarkdownDocumentationReferenceAttributeNode) {
	if ref, ok := reference.(*BLangMarkdownReferenceDocumentation); ok {
		b.References = append(b.References, *ref)
	} else {
		panic("reference is not a BLangMarkdownReferenceDocumentation")
	}
}

func (b *BLangMarkdownReferenceDocumentation) GetType() model.DocumentationReferenceType {
	return b.Type
}

func (b *BLangMarkdownReferenceDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_REFERENCE
}

// BLangService methods

func (b *BLangService) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangService) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangService) GetResources() []model.FunctionNode {
	return []model.FunctionNode{}
}

func (b *BLangService) IsAnonymousService() bool {
	return false
}

func (b *BLangService) GetAttachedExprs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(b.AttachedExprs))
	for i := range b.AttachedExprs {
		result[i] = b.AttachedExprs[i]
	}
	return result
}

func (b *BLangService) GetServiceClass() model.ClassDefinition {
	return b.ServiceClass
}

func (b *BLangService) GetAbsolutePath() []model.IdentifierNode {
	return b.AbsoluteResourcePath
}

func (b *BLangService) GetServiceNameLiteral() model.LiteralNode {
	return b.ServiceNameLiteral
}

func (b *BLangService) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	for i := range b.AnnAttachments {
		result[i] = &b.AnnAttachments[i]
	}
	return result
}

func (b *BLangService) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.AnnAttachments = append(b.AnnAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (b *BLangService) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangService) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (b *BLangService) GetKind() model.NodeKind {
	return model.NodeKind_SERVICE
}

func (b *BLangFunction) GetKind() model.NodeKind {
	return model.NodeKind_FUNCTION
}

func (b *BLangFunction) Scope() model.Scope {
	return b.scope
}

func (b *BLangFunction) SetScope(scope model.Scope) {
	b.scope = scope
}

var _ NodeWithScope = &BLangFunction{}

func (b *bLangInvokableNodeBase) GetName() model.IdentifierNode {
	return &b.Name
}

func (b *bLangInvokableNodeBase) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = *id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *bLangInvokableNodeBase) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *bLangInvokableNodeBase) GetAnnAttachments() []model.AnnotationAttachmentNode {
	attachments := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	copy(attachments, b.AnnAttachments)
	return attachments
}

func (b *bLangInvokableNodeBase) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	b.AnnAttachments = append(b.AnnAttachments, annAttachment)
}

func (b *bLangInvokableNodeBase) SetAnnAttachments(annAttachments []model.AnnotationAttachmentNode) {
	b.AnnAttachments = annAttachments
}

func (b *bLangInvokableNodeBase) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *bLangInvokableNodeBase) SetMarkdownDocumentationAttachment(markdownDocumentationAttachment model.MarkdownDocumentationNode) {
	if doc, ok := markdownDocumentationAttachment.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = doc
	} else {
		panic("markdownDocumentationAttachment is not a BLangMarkdownDocumentation")
	}
}

func (b *bLangInvokableNodeBase) GetParameters() []model.SimpleVariableNode {
	result := make([]model.SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *bLangInvokableNodeBase) AddParameter(param model.SimpleVariableNode) {
	if blangParam, ok := param.(*BLangSimpleVariable); ok {
		b.RequiredParams = append(b.RequiredParams, *blangParam)
	} else {
		panic("param is not a BLangSimpleVariable")
	}
}

func (b *bLangInvokableNodeBase) GetRequiredParams() []model.SimpleVariableNode {
	result := make([]model.SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *bLangInvokableNodeBase) SetRequiredParams(requiredParams []model.SimpleVariableNode) {
	b.RequiredParams = make([]BLangSimpleVariable, len(requiredParams))
	for i, param := range requiredParams {
		if blangParam, ok := param.(*BLangSimpleVariable); ok {
			b.RequiredParams[i] = *blangParam
		} else {
			panic("requiredParams contains element that is not a BLangSimpleVariable")
		}
	}
}

func (b *bLangInvokableNodeBase) GetRestParameters() model.SimpleVariableNode {
	return b.RestParam
}

func (b *bLangInvokableNodeBase) GetRestParam() model.SimpleVariableNode {
	return b.RestParam
}

func (b *bLangInvokableNodeBase) SetRestParameter(restParam model.SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *bLangInvokableNodeBase) SetRestParam(restParam model.SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *bLangInvokableNodeBase) HasBody() bool {
	return b.Body != nil
}

func (b *bLangInvokableNodeBase) GetReturnTypeDescriptor() model.TypeDescriptor {
	return b.returnTypeDescriptor
}

func (b *bLangInvokableNodeBase) SetReturnTypeDescriptor(typeDescriptor model.TypeDescriptor) {
	b.returnTypeDescriptor = typeDescriptor
}

func (b *bLangInvokableNodeBase) GetBody() model.FunctionBodyNode {
	return b.Body
}

func (b *bLangInvokableNodeBase) SetBody(body model.FunctionBodyNode) {
	b.Body = body
}

func (b *BLangVariableBase) GetAnnAttachments() []model.AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BLangVariableBase) SetAnnAttachments(annAttachments []model.AnnotationAttachmentNode) {
	b.AnnAttachments = annAttachments
}

func (b *BLangVariableBase) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangVariableBase) SetMarkdownDocumentationAttachment(markdownDocumentationAttachment model.MarkdownDocumentationNode) {
	b.MarkdownDocumentationAttachment = markdownDocumentationAttachment
}

func (b *BLangVariableBase) GetExpr() model.ExpressionNode {
	return b.Expr
}

func (b *BLangVariableBase) SetExpr(expr model.ExpressionNode) {
	b.Expr = expr
}

func (b *BLangVariableBase) GetIsDeclaredWithVar() bool {
	return b.IsDeclaredWithVar
}

func (b *BLangVariableBase) SetIsDeclaredWithVar(isDeclaredWithVar bool) {
	b.IsDeclaredWithVar = isDeclaredWithVar
}

func (m *BLangVariableBase) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangVariable.java:83:5
	m.AnnAttachments = append(m.AnnAttachments, annAttachment)
}

func (m *BLangVariableBase) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return m.AnnAttachments
}

func (m *BLangVariableBase) GetInitialExpression() model.ExpressionNode {
	return m.Expr
}

func (m *BLangVariableBase) SetInitialExpression(expr model.ExpressionNode) {
	m.Expr = expr
}

// BLangTypeDefinition methods

func NewBLangTypeDefinition() *BLangTypeDefinition {
	this := &BLangTypeDefinition{}
	this.annAttachments = []BLangAnnotationAttachment{}
	this.CycleDepth = -1
	this.hasCyclicReference = false
	return this
}

func (b *BLangTypeDefinition) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangTypeDefinition) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangTypeDefinition) GetTypeData() model.TypeData {
	return b.typeData
}

func (b *BLangTypeDefinition) SetTypeData(typeData model.TypeData) {
	b.typeData = typeData
}

func (b *BLangTypeDefinition) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(b.annAttachments))
	for i := range b.annAttachments {
		result[i] = &b.annAttachments[i]
	}
	return result
}

func (b *BLangTypeDefinition) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		b.annAttachments = append(b.annAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (b *BLangTypeDefinition) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.markdownDocumentationAttachment
}

func (b *BLangTypeDefinition) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		b.markdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (b *BLangTypeDefinition) GetPrecedence() int {
	return b.precedence
}

func (b *BLangTypeDefinition) SetPrecedence(precedence int) {
	b.precedence = precedence
}

func (b *BLangTypeDefinition) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_DEFINITION
}

func (b *BLangTypeDefinition) GetCycleDepth() int {
	return b.CycleDepth
}

func (b *BLangTypeDefinition) SetCycleDepth(depth int) {
	b.CycleDepth = depth
}

func (b *BLangXMLNS) GetNamespaceURI() model.ExpressionNode {
	return b.namespaceURI
}

func (b *BLangXMLNS) GetPrefix() model.IdentifierNode {
	return b.prefix
}

func (b *BLangXMLNS) SetNamespaceURI(namespaceURI model.ExpressionNode) {
	if expr, ok := namespaceURI.(BLangExpression); ok {
		b.namespaceURI = expr
	} else {
		panic("namespaceURI is not a BLangExpression")
	}
}

func (b *BLangXMLNS) SetPrefix(prefix model.IdentifierNode) {
	if ident, ok := prefix.(*BLangIdentifier); ok {
		b.prefix = ident
	} else {
		panic("prefix is not a BLangIdentifier")
	}
}

func (b *BLangXMLNS) GetKind() model.NodeKind {
	return model.NodeKind_XMLNS
}

func (b *BLangPackage) GetCompilationUnits() []model.CompilationUnitNode {
	result := make([]model.CompilationUnitNode, len(b.CompUnits))
	for i := range b.CompUnits {
		result[i] = &b.CompUnits[i]
	}
	return result
}

func (b *BLangPackage) AddCompilationUnit(compUnit model.CompilationUnitNode) {
	if cu, ok := compUnit.(*BLangCompilationUnit); ok {
		b.CompUnits = append(b.CompUnits, *cu)
	} else {
		panic("compUnit is not a BLangCompilationUnit")
	}
}

func (b *BLangPackage) GetImports() []model.ImportPackageNode {
	result := make([]model.ImportPackageNode, len(b.Imports))
	for i := range b.Imports {
		result[i] = &b.Imports[i]
	}
	return result
}

func (b *BLangPackage) AddImport(importPkg model.ImportPackageNode) {
	if imp, ok := importPkg.(*BLangImportPackage); ok {
		b.Imports = append(b.Imports, *imp)
	} else {
		panic("importPkg is not a BLangImportPackage")
	}
}

func (b *BLangPackage) GetNamespaceDeclarations() []model.XMLNSDeclarationNode {
	result := make([]model.XMLNSDeclarationNode, len(b.XmlnsList))
	for i := range b.XmlnsList {
		result[i] = &b.XmlnsList[i]
	}
	return result
}

func (b *BLangPackage) AddNamespaceDeclaration(xmlnsDecl model.XMLNSDeclarationNode) {
	if xmlns, ok := xmlnsDecl.(*BLangXMLNS); ok {
		b.XmlnsList = append(b.XmlnsList, *xmlns)
		b.TopLevelNodes = append(b.TopLevelNodes, xmlnsDecl)
	} else {
		panic("xmlnsDecl is not a BLangXMLNS")
	}
}

func (b *BLangPackage) GetConstants() []model.ConstantNode {
	result := make([]model.ConstantNode, len(b.Constants))
	for i := range b.Constants {
		result[i] = &b.Constants[i]
	}
	return result
}

func (b *BLangPackage) GetGlobalVariables() []model.VariableNode {
	result := make([]model.VariableNode, len(b.GlobalVars))
	for i := range b.GlobalVars {
		result[i] = &b.GlobalVars[i]
	}
	return result
}

func (b *BLangPackage) AddGlobalVariable(globalVar model.SimpleVariableNode) {
	if sv, ok := globalVar.(*BLangSimpleVariable); ok {
		b.GlobalVars = append(b.GlobalVars, *sv)
		b.TopLevelNodes = append(b.TopLevelNodes, globalVar)
	} else {
		panic("globalVar is not a BLangSimpleVariable")
	}
}

func (b *BLangPackage) GetServices() []model.ServiceNode {
	result := make([]model.ServiceNode, len(b.Services))
	for i := range b.Services {
		result[i] = &b.Services[i]
	}
	return result
}

func (b *BLangPackage) AddService(service model.ServiceNode) {
	if svc, ok := service.(*BLangService); ok {
		b.Services = append(b.Services, *svc)
		b.TopLevelNodes = append(b.TopLevelNodes, service)
	} else {
		panic("service is not a BLangService")
	}
}

func (b *BLangPackage) GetFunctions() []model.FunctionNode {
	result := make([]model.FunctionNode, len(b.Functions))
	for i := range b.Functions {
		result[i] = &b.Functions[i]
	}
	return result
}

func (b *BLangPackage) AddFunction(function model.FunctionNode) {
	if fn, ok := function.(*BLangFunction); ok {
		b.Functions = append(b.Functions, *fn)
		b.TopLevelNodes = append(b.TopLevelNodes, function)
	} else {
		panic("function is not a BLangFunction")
	}
}

func (b *BLangPackage) GetTypeDefinitions() []model.TypeDefinition {
	result := make([]model.TypeDefinition, len(b.TypeDefinitions))
	for i := range b.TypeDefinitions {
		result[i] = &b.TypeDefinitions[i]
	}
	return result
}

func (b *BLangPackage) AddTypeDefinition(typeDefinition model.TypeDefinition) {
	if td, ok := typeDefinition.(*BLangTypeDefinition); ok {
		b.TypeDefinitions = append(b.TypeDefinitions, *td)
		b.TopLevelNodes = append(b.TopLevelNodes, typeDefinition)
	} else {
		panic("typeDefinition is not a BLangTypeDefinition")
	}
}

func (b *BLangPackage) GetAnnotations() []model.AnnotationNode {
	result := make([]model.AnnotationNode, len(b.Annotations))
	for i := range b.Annotations {
		result[i] = &b.Annotations[i]
	}
	return result
}

func (b *BLangPackage) AddAnnotation(annotation model.AnnotationNode) {
	if ann, ok := annotation.(*BLangAnnotation); ok {
		b.Annotations = append(b.Annotations, *ann)
		b.TopLevelNodes = append(b.TopLevelNodes, annotation)
	} else {
		panic("annotation is not a BLangAnnotation")
	}
}

func (b *BLangPackage) GetClassDefinitions() []model.ClassDefinition {
	result := make([]model.ClassDefinition, len(b.ClassDefinitions))
	for i := range b.ClassDefinitions {
		result[i] = &b.ClassDefinitions[i]
	}
	return result
}

func (b *BLangPackage) GetKind() model.NodeKind {
	return model.NodeKind_PACKAGE
}

func (b *BLangPackage) AddTestablePkg(testablePkg *BLangTestablePackage) {
	b.TestablePkgs = append(b.TestablePkgs, testablePkg)
}

func (b *BLangPackage) GetTestablePkgs() []*BLangTestablePackage {
	return b.TestablePkgs
}

func (b *BLangPackage) GetTestablePkg() *BLangTestablePackage {
	if len(b.TestablePkgs) > 0 {
		return b.TestablePkgs[0]
	}
	return nil
}

func (b *BLangPackage) ContainsTestablePkg() bool {
	return len(b.TestablePkgs) > 0
}

func (b *BLangPackage) HasTestablePackage() bool {
	return len(b.TestablePkgs) > 0
}

func (b *BLangPackage) AddClassDefinition(classDefNode *BLangClassDefinition) {
	b.TopLevelNodes = append(b.TopLevelNodes, classDefNode)
	b.ClassDefinitions = append(b.ClassDefinitions, *classDefNode)
}

func (b *BLangPackage) AddDiagnostic(diagnostic diagnostics.Diagnostic) {
	// Check if diagnostic already exists
	for _, existing := range b.diagnostics {
		if diagnosticEqual(existing, diagnostic) {
			return
		}
	}
	b.diagnostics = append(b.diagnostics, diagnostic)
	severity := diagnostic.DiagnosticInfo().Severity()
	switch severity {
	case diagnostics.Error:
		b.errorCount++
	case diagnostics.Warning:
		b.warnCount++
	}
}

func diagnosticEqual(d1, d2 diagnostics.Diagnostic) bool {
	info1 := d1.DiagnosticInfo()
	info2 := d2.DiagnosticInfo()
	return info1.Code() == info2.Code() &&
		info1.MessageFormat() == info2.MessageFormat() &&
		info1.Severity() == info2.Severity()
}

func (b *BLangPackage) GetDiagnostics() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(b.diagnostics))
	copy(result, b.diagnostics)
	return result
}

func (b *BLangPackage) GetErrorCount() int {
	return b.errorCount
}

func (b *BLangPackage) GetWarnCount() int {
	return b.warnCount
}

func (b *BLangPackage) HasErrors() bool {
	return b.errorCount > 0
}

func NewBLangPackage(env semtypes.Env) *BLangPackage {
	this := &BLangPackage{}
	this.CompUnits = []BLangCompilationUnit{}
	this.Imports = []BLangImportPackage{}
	this.XmlnsList = []BLangXMLNS{}
	this.Constants = []BLangConstant{}
	this.GlobalVars = []BLangSimpleVariable{}
	this.Services = []BLangService{}
	this.Functions = []BLangFunction{}
	this.TypeDefinitions = []BLangTypeDefinition{}
	this.Annotations = []BLangAnnotation{}
	this.TopLevelNodes = []model.TopLevelNode{}
	this.TestablePkgs = []*BLangTestablePackage{}
	this.ClassDefinitions = []BLangClassDefinition{}
	this.CompletedPhases = common.UnorderedSet[CompilerPhase]{}
	this.LambdaFunctions = []BLangLambdaFunction{}
	this.errorCount = 0
	this.warnCount = 0
	this.diagnostics = []diagnostics.Diagnostic{}
	return this
}

func (b *BLangTestablePackage) GetMockFunctionNamesMap() map[string]string {
	return b.mockFunctionNamesMap
}

func (b *BLangTestablePackage) AddMockFunction(id string, function string) {
	if b.mockFunctionNamesMap == nil {
		b.mockFunctionNamesMap = make(map[string]string)
	}
	b.mockFunctionNamesMap[id] = function
}

func (b *BLangTestablePackage) GetIsLegacyMockingMap() map[string]bool {
	return b.isLegacyMockingMap
}

func (b *BLangTestablePackage) AddIsLegacyMockingMap(id string, isLegacy bool) {
	if b.isLegacyMockingMap == nil {
		b.isLegacyMockingMap = make(map[string]bool)
	}
	b.isLegacyMockingMap[id] = isLegacy
}

func createSimpleVariableNodeWithLocationTokenLocation(location diagnostics.Location, identifier tree.Token, identifierPos diagnostics.Location) *BLangSimpleVariable {
	memberVar := createSimpleVariableNode()
	memberVar.pos = location
	name := createIdentifierFromToken(identifierPos, identifier)
	BLangNode(&name).SetPosition(identifierPos)
	memberVar.SetName(&name)
	return memberVar
}

func createSimpleVariableNode() *BLangSimpleVariable {
	return &BLangSimpleVariable{}
}

func createConstantNode() *BLangConstant {
	return &BLangConstant{}
}

func GetCompilationUnit(cx *context.CompilerContext, syntaxTree *tree.SyntaxTree) *BLangCompilationUnit {
	nodeBuilder := NewNodeBuilder(cx)
	compilationUnit := nodeBuilder.TransformModulePart(syntaxTree.RootNode.(*tree.ModulePart))
	return compilationUnit.(*BLangCompilationUnit)
}

// TODO: get rid of this once we have a proper project api. This just remaps compilation unit to a BLangPackage.
func ToPackage(compilationUnit *BLangCompilationUnit) *BLangPackage {
	p := BLangPackage{}
	p.PackageID = compilationUnit.packageID
	for _, node := range compilationUnit.TopLevelNodes {
		switch node := node.(type) {
		case *BLangImportPackage:
			p.Imports = append(p.Imports, *node)
		case *BLangConstant:
			p.Constants = append(p.Constants, *node)
		case *BLangService:
			p.Services = append(p.Services, *node)
		case *BLangSimpleVariable:
			p.GlobalVars = append(p.GlobalVars, *node)
		case *BLangFunction:
			if node.Name.Value == "init" {
				p.InitFunction = node
			} else {
				p.Functions = append(p.Functions, *node)
			}
		case *BLangTypeDefinition:
			p.TypeDefinitions = append(p.TypeDefinitions, *node)
		case *BLangAnnotation:
			p.Annotations = append(p.Annotations, *node)
		case *BLangClassDefinition:
			p.ClassDefinitions = append(p.ClassDefinitions, *node)
		default:
			p.TopLevelNodes = append(p.TopLevelNodes, node)
		}
	}
	return &p
}
