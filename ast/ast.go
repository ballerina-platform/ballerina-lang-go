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
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"strings"
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
	SetSymbol(symbol model.Symbol)
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

type Location = diagnostics.Location

type BLangNode interface {
	model.Node
	SetTypeData(ty model.TypeData)
	SetDeterminedType(ty semtypes.SemType)
	SetPosition(pos Location)
}

type (
	BLangNodeBase struct {
		ty             model.TypeData
		DeterminedType semtypes.SemType

		parent BLangNode

		pos                Location
		desugared          bool
		constantPropagated bool
		internal           bool
	}

	// TODO: look into what is the difference between the public TypeData fields and ty
	BLangAnnotation struct {
		BLangNodeBase
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		TypeData                        model.TypeData
		FlagSet                         common.UnorderedSet[model.Flag]
		attachPoints                    common.UnorderedSet[model.AttachPoint]
	}

	BLangAnnotationAttachment struct {
		BLangNodeBase
		Expr           BLangExpression
		AnnotationName *BLangIdentifier
		PkgAlias       *BLangIdentifier
		AttachPoints   common.OrderedSet[model.Point]
	}

	BLangFunctionBodyBase struct {
		BLangNodeBase
	}

	BLangBlockFunctionBody struct {
		BLangFunctionBodyBase
		Stmts []BLangStatement
	}

	BLangExprFunctionBody struct {
		BLangFunctionBodyBase
		Expr model.ExpressionNode
	}

	BLangIdentifier struct {
		BLangNodeBase
		Value         string
		OriginalValue string
		isLiteral     bool
	}

	BLangImportPackage struct {
		BLangNodeBase
		OrgName      *BLangIdentifier
		PkgNameComps []BLangIdentifier
		Alias        *BLangIdentifier
		CompUnit     *BLangIdentifier
		Version      *BLangIdentifier
	}

	BLangClassDefinition struct {
		BLangNodeBase
		Name                            *BLangIdentifier
		symbol                          model.Symbol
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		InitFunction                    *BLangFunction
		Functions                       []BLangFunction
		Fields                          []model.SimpleVariableNode
		TypeRefs                        []model.TypeDescriptor
		FlagSet                         common.Set[model.Flag]
		GeneratedInitFunction           *BLangFunction
		Receiver                        *BLangSimpleVariable
		ReferencedFields                []BLangSimpleVariable
		LocalVarRefs                    []BLangLocalVarRef
		OceEnvData                      *OCEDynamicEnvironmentData
		ObjectType                      *BObjectType
		CycleDepth                      int
		Precedence                      int
		IsServiceDecl                   bool
		HasClosureVars                  bool
		IsObjectContructorDecl          bool
		DefinitionCompleted             bool
	}

	BLangService struct {
		BLangNodeBase
		symbol                          model.Symbol
		ServiceVariable                 *BLangSimpleVariable
		AttachedExprs                   []BLangExpression
		ServiceClass                    *BLangClassDefinition
		AbsoluteResourcePath            []model.IdentifierNode
		ServiceNameLiteral              *BLangLiteral
		Name                            *BLangIdentifier
		AnnAttachments                  []BLangAnnotationAttachment
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		FlagSet                         common.UnorderedSet[model.Flag]
		ListenerType                    BType
		ResourceFunctions               []BLangFunction
		InferredServiceType             BType
	}

	BLangCompilationUnit struct {
		BLangNodeBase
		TopLevelNodes []model.TopLevelNode
		Name          string
		packageID     *model.PackageID
		sourceKind    SourceKind
	}

	BLangPackage struct {
		BLangNodeBase
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
		FlagSet                 common.UnorderedSet[model.Flag]
		CompletedPhases         common.UnorderedSet[CompilerPhase]
		LambdaFunctions         []BLangLambdaFunction
		PackageID               *model.PackageID
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
		BLangNodeBase
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
		BLangNodeBase
		DocumentationLines                []BLangMarkdownDocumentationLine
		Parameters                        []BLangMarkdownParameterDocumentation
		References                        []BLangMarkdownReferenceDocumentation
		ReturnParameter                   *BLangMarkdownReturnParameterDocumentation
		DeprecationDocumentation          *BLangMarkDownDeprecationDocumentation
		DeprecatedParametersDocumentation *BLangMarkDownDeprecatedParametersDocumentation
	}
	BLangMarkdownReferenceDocumentation struct {
		BLangNodeBase
		Qualifier         string
		TypeName          string
		Identifier        string
		ReferenceName     string
		Type              model.DocumentationReferenceType
		HasParserWarnings bool
	}

	BLangVariableBase struct {
		BLangNodeBase
		AnnAttachments                  []model.AnnotationAttachmentNode
		MarkdownDocumentationAttachment model.MarkdownDocumentationNode
		Expr                            model.ExpressionNode
		FlagSet                         common.Set[model.Flag]
		IsDeclaredWithVar               bool
		symbol                          model.Symbol
	}

	BLangConstant struct {
		BLangVariableBase
		Name                     *BLangIdentifier
		AssociatedTypeDefinition *BLangTypeDefinition
	}

	BLangSimpleVariable struct {
		BLangVariableBase
		Name *BLangIdentifier
	}

	ClosureVarSymbol struct {
		DiagnosticLocation Location
	}

	BLangInvokableNodeBase struct {
		BLangNodeBase
		Name                            *BLangIdentifier
		symbol                          model.Symbol
		AnnAttachments                  []model.AnnotationAttachmentNode
		MarkdownDocumentationAttachment *BLangMarkdownDocumentation
		RequiredParams                  []BLangSimpleVariable
		RestParam                       model.SimpleVariableNode
		ReturnTypeData                  model.TypeData
		ReturnTypeAnnAttachments        []model.AnnotationAttachmentNode
		Body                            model.FunctionBodyNode
		DefaultWorkerName               model.IdentifierNode
		FlagSet                         common.UnorderedSet[model.Flag]
		DesugaredReturnType             bool
	}

	BLangFunction struct {
		BLangInvokableNodeBase
		scope             model.Scope
		Receiver          *BLangSimpleVariable
		ClosureVarSymbols common.OrderedSet[ClosureVarSymbol]
		SendsToThis       common.OrderedSet[Channel]
		AnonForkName      string
		MapSymbolUpdated  bool
		AttachedFunction  bool
		ObjInitFunction   bool
		InterfaceFunction bool
	}

	BLangTypeDefinition struct {
		BLangNodeBase
		Name                            *BLangIdentifier
		symbol                          model.Symbol
		typeData                        model.TypeData
		annAttachments                  []BLangAnnotationAttachment
		markdownDocumentationAttachment *BLangMarkdownDocumentation
		FlagSet                         common.UnorderedSet[model.Flag]
		precedence                      int
		CycleDepth                      int
		isBuiltinTypeDef                bool
		hasCyclicReference              bool
		referencedFieldsDefined         bool
	}
)

func (this *BLangNodeBase) SetTypeData(ty model.TypeData) {
	this.ty = ty
}

func (this *BLangNodeBase) GetTypeData() model.TypeData {
	return this.ty
}

func (this *BLangNodeBase) SetDeterminedType(ty semtypes.SemType) {
	this.DeterminedType = ty
}

func (this *BLangNodeBase) GetDeterminedType() semtypes.SemType {
	return this.DeterminedType
}

func (this *BLangNodeBase) GetPosition() Location {
	return this.pos
}

func (this *BLangNodeBase) SetPosition(pos Location) {
	this.pos = pos
}

func (n *BLangClassDefinition) Symbol() model.Symbol {
	return n.symbol
}

func (n *BLangClassDefinition) SetSymbol(symbol model.Symbol) {
	n.symbol = symbol
}

func (n *BLangService) Symbol() model.Symbol {
	return n.symbol
}

func (n *BLangService) SetSymbol(symbol model.Symbol) {
	n.symbol = symbol
}

func (n *BLangVariableBase) Symbol() model.Symbol {
	return n.symbol
}

func (n *BLangVariableBase) SetSymbol(symbol model.Symbol) {
	n.symbol = symbol
}

func (n *BLangInvokableNodeBase) Symbol() model.Symbol {
	return n.symbol
}

func (n *BLangInvokableNodeBase) SetSymbol(symbol model.Symbol) {
	n.symbol = symbol
}

func (n *BLangTypeDefinition) Symbol() model.Symbol {
	return n.symbol
}

func (n *BLangTypeDefinition) SetSymbol(symbol model.Symbol) {
	n.symbol = symbol
}

var (
	_ model.AnnotationAttachmentNode                    = &BLangAnnotationAttachment{}
	_ model.IdentifierNode                              = &BLangIdentifier{}
	_ model.ImportPackageNode                           = &BLangImportPackage{}
	_ model.ClassDefinition                             = &BLangClassDefinition{}
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

func (this *BLangAnnotationAttachment) GetKind() model.NodeKind {
	// migrated from BLangAnnotationAttachment.java:89:5
	return model.NodeKind_ANNOTATION_ATTACHMENT
}

func (this *BLangAnnotationAttachment) GetPackgeAlias() model.IdentifierNode {
	return this.PkgAlias
}

func (this *BLangAnnotationAttachment) SetPackageAlias(pkgAlias model.IdentifierNode) {
	if id, ok := pkgAlias.(*BLangIdentifier); ok {
		this.PkgAlias = id
	} else {
		panic("pkgAlias is not a BLangIdentifier")
	}
}

func (this *BLangAnnotationAttachment) GetAnnotationName() model.IdentifierNode {
	return this.AnnotationName
}

func (this *BLangAnnotationAttachment) SetAnnotationName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.AnnotationName = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangAnnotationAttachment) GetExpressionNode() model.ExpressionNode {
	return this.Expr
}

func (this *BLangAnnotationAttachment) SetExpressionNode(expr model.ExpressionNode) {
	if expr, ok := expr.(BLangExpression); ok {
		this.Expr = expr
	} else {
		panic("expr is not a BLangExpression")
	}
}

func (this *BLangAnnotation) GetKind() model.NodeKind {
	// migrated from BLangAnnotation.java:135:5
	return model.NodeKind_ANNOTATION
}

func (this *BLangAnnotation) GetName() model.IdentifierNode {
	// migrated from BLangAnnotation.java:80:5
	return this.Name
}

func (this *BLangAnnotation) SetName(name model.IdentifierNode) {
	// migrated from BLangAnnotation.java:85:5
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangAnnotation) GetTypeData() model.TypeData {
	// migrated from BLangAnnotation.java:70:5
	return this.TypeData
}

func (this *BLangAnnotation) SetTypeData(typeData model.TypeData) {
	// migrated from BLangAnnotation.java:75:5
	this.TypeData = typeData
}

func (this *BLangAnnotation) GetFlags() common.Set[model.Flag] {
	// migrated from BLangAnnotation.java:90:5
	return &this.FlagSet
}

func (this *BLangAnnotation) AddFlag(flag model.Flag) {
	// migrated from BLangAnnotation.java:95:5
	(&this.FlagSet).Add(flag)
}

func (this *BLangAnnotation) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangAnnotation.java:100:5
	attachments := make([]model.AnnotationAttachmentNode, len(this.AnnAttachments))
	for i, attachment := range this.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (this *BLangAnnotation) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangAnnotation.java:105:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangAnnotation) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangAnnotation.java:110:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangAnnotation) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangAnnotation.java:115:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangBlockFunctionBody) GetKind() model.NodeKind {
	// migrated from BLangBlockFunctionBody.java:73:5
	return model.NodeKind_BLOCK_FUNCTION_BODY
}

func (this *BLangExprFunctionBody) GetKind() model.NodeKind {
	// migrated from BLangExprFunctionBody.java:50:5
	return model.NodeKind_EXPR_FUNCTION_BODY
}

func (this *BLangExprFunctionBody) GetExpr() model.ExpressionNode {
	// migrated from BLangExprFunctionBody.java:55:5
	return this.Expr
}

func (this *BLangIdentifier) GetValue() string {
	// migrated from BLangIdentifier.java:32:5
	return this.Value
}

func (this *BLangIdentifier) GetKind() model.NodeKind {
	// migrated from BLangIdentifier.java:32:5
	return model.NodeKind_IDENTIFIER
}

func (this *BLangIdentifier) SetValue(value string) {
	// migrated from BLangIdentifier.java:37:5
	this.Value = value
}

func (this *BLangIdentifier) SetOriginalValue(value string) {
	// migrated from BLangIdentifier.java:42:5
	this.OriginalValue = value
}

func (this *BLangIdentifier) IsLiteral() bool {
	// migrated from BLangIdentifier.java:47:5
	return this.isLiteral
}

func (this *BLangIdentifier) SetLiteral(isLiteral bool) {
	// migrated from BLangIdentifier.java:52:5
	this.isLiteral = isLiteral
}

func (this *BLangImportPackage) GetKind() model.NodeKind {
	return model.NodeKind_IMPORT
}

func (this *BLangImportPackage) GetOrgName() model.IdentifierNode {
	return this.OrgName
}

func (this *BLangImportPackage) GetPackageName() []model.IdentifierNode {
	result := make([]model.IdentifierNode, len(this.PkgNameComps))
	for i := range this.PkgNameComps {
		result[i] = &this.PkgNameComps[i]
	}
	return result
}

func (this *BLangImportPackage) SetPackageName(nameParts []model.IdentifierNode) {
	this.PkgNameComps = make([]BLangIdentifier, 0, len(nameParts))
	for _, namePart := range nameParts {
		if id, ok := namePart.(*BLangIdentifier); ok {
			this.PkgNameComps = append(this.PkgNameComps, *id)
		} else {
			panic("namePart is not a BLangIdentifier")
		}
	}
}

func (this *BLangImportPackage) GetPackageVersion() model.IdentifierNode {
	return this.Version
}

func (this *BLangImportPackage) SetPackageVersion(version model.IdentifierNode) {
	if id, ok := version.(*BLangIdentifier); ok {
		this.Version = id
	} else {
		panic("version is not a BLangIdentifier")
	}
}

func (this *BLangImportPackage) GetAlias() model.IdentifierNode {
	return this.Alias
}

func (this *BLangImportPackage) SetAlias(alias model.IdentifierNode) {
	if id, ok := alias.(*BLangIdentifier); ok {
		this.Alias = id
	} else {
		panic("alias is not a BLangIdentifier")
	}
}

func NewBLangClassDefinition() BLangClassDefinition {
	this := BLangClassDefinition{}
	this.CycleDepth = (-1)
	this.IsObjectContructorDecl = false
	// Default field initializations
	this.FlagSet = &common.UnorderedSet[model.Flag]{}
	this.FlagSet.Add(model.Flag_CLASS)

	return this
}

func (this *BLangClassDefinition) GetName() model.IdentifierNode {
	// migrated from BLangClassDefinition.java:88:5
	return this.Name
}

func (this *BLangClassDefinition) SetName(name model.IdentifierNode) {
	// migrated from BLangClassDefinition.java:93:5
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangClassDefinition) GetFunctions() []model.FunctionNode {
	// migrated from BLangClassDefinition.java:98:5
	result := make([]model.FunctionNode, len(this.Functions))
	for i := range this.Functions {
		result[i] = &this.Functions[i]
	}
	return result
}

func (this *BLangClassDefinition) AddFunction(function model.FunctionNode) {
	// migrated from BLangClassDefinition.java:103:5
	if function, ok := function.(*BLangFunction); ok {
		this.Functions = append(this.Functions, *function)
		return
	}
	panic("function is not a BLangFunction")
}

func (this *BLangClassDefinition) GetInitFunction() model.FunctionNode {
	// migrated from BLangClassDefinition.java:108:5
	return this.InitFunction
}

func (this *BLangClassDefinition) AddField(field model.VariableNode) {
	// migrated from BLangClassDefinition.java:113:5
	if field, ok := field.(*BLangSimpleVariable); ok {
		this.Fields = append(this.Fields, field)
		return
	}
	panic("field is not a BLangSimpleVariable")
}

func (this *BLangClassDefinition) AddTypeReference(typeRef *model.TypeData) {
	// migrated from BLangClassDefinition.java:118:5
	this.TypeRefs = append(this.TypeRefs, typeRef.TypeDescriptor)
}

func (this *BLangClassDefinition) GetKind() model.NodeKind {
	// migrated from BLangClassDefinition.java:138:5
	return model.NodeKind_CLASS_DEFN
}

func (this *BLangClassDefinition) GetFlags() common.Set[model.Flag] {
	// migrated from BLangClassDefinition.java:158:5
	return this.FlagSet
}

func (this *BLangClassDefinition) AddFlag(flag model.Flag) {
	// migrated from BLangClassDefinition.java:163:5
	this.FlagSet.Add(flag)
}

func (this *BLangClassDefinition) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangClassDefinition.java:168:5
	attachments := make([]model.AnnotationAttachmentNode, len(this.AnnAttachments))
	for i, attachment := range this.AnnAttachments {
		attachments[i] = &attachment
	}
	return attachments
}

func (this *BLangClassDefinition) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangClassDefinition.java:173:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangClassDefinition) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangClassDefinition.java:178:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangClassDefinition) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangClassDefinition.java:183:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangClassDefinition) GetPrecedence() int {
	// migrated from BLangClassDefinition.java:188:5
	return this.Precedence
}

func (this *BLangClassDefinition) SetPrecedence(precedence int) {
	// migrated from BLangClassDefinition.java:193:5
	this.Precedence = precedence
}

func (this *BLangCompilationUnit) AddTopLevelNode(node model.TopLevelNode) {
	// migrated from BLangCompilationUnit.java:48:5
	this.TopLevelNodes = append(this.TopLevelNodes, node)
}

func (this *BLangCompilationUnit) GetTopLevelNodes() []model.TopLevelNode {
	// migrated from BLangCompilationUnit.java:53:5
	return this.TopLevelNodes
}

func (this *BLangCompilationUnit) GetName() string {
	// migrated from BLangCompilationUnit.java:58:5
	return this.Name
}

func (this *BLangCompilationUnit) SetName(name string) {
	// migrated from BLangCompilationUnit.java:63:5
	this.Name = name
}

func (this *BLangCompilationUnit) GetPackageID() *model.PackageID {
	// migrated from BLangCompilationUnit.java:68:5
	return this.packageID
}

func (this *BLangCompilationUnit) SetPackageID(packageID *model.PackageID) {
	// migrated from BLangCompilationUnit.java:72:5
	this.packageID = packageID
}

func (this *BLangCompilationUnit) GetKind() model.NodeKind {
	// migrated from BLangCompilationUnit.java:76:5
	return model.NodeKind_COMPILATION_UNIT
}

func (this *BLangCompilationUnit) SetSourceKind(kind SourceKind) {
	// migrated from BLangCompilationUnit.java:81:5
	this.sourceKind = kind
}

func (this *BLangCompilationUnit) GetSourceKind() SourceKind {
	// migrated from BLangCompilationUnit.java:86:5
	return this.sourceKind
}

func (this *BLangConstant) GetName() model.IdentifierNode {
	return this.Name
}

func (this *BLangConstant) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangConstant) GetFlags() common.Set[model.Flag] {
	// migrated from BLangConstant.java:78:5
	return this.FlagSet
}

func (this *BLangConstant) AddFlag(flag model.Flag) {
	// migrated from BLangConstant.java:83:5
	this.FlagSet.Add(flag)
}

func (this *BLangConstant) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	// migrated from BLangConstant.java:88:5
	return this.AnnAttachments
}

func (this *BLangConstant) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	// migrated from BLangConstant.java:93:5
	if annAttachment, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, annAttachment)
		return
	}
	panic("annAttachment is not a BLangAnnotationAttachment")
}

func (this *BLangConstant) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	// migrated from BLangConstant.java:98:5
	return this.MarkdownDocumentationAttachment
}

func (this *BLangConstant) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	// migrated from BLangConstant.java:103:5
	if documentationNode, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = documentationNode
		return
	}
	panic("documentationNode is not a BLangMarkdownDocumentation")
}

func (this *BLangConstant) GetKind() model.NodeKind {
	// migrated from BLangConstant.java:108:5
	return model.NodeKind_CONSTANT
}

func (this *BLangConstant) GetAssociatedTypeDefinition() model.TypeDefinition {
	// migrated from BLangConstant.java:139:5
	return this.AssociatedTypeDefinition
}

func (this *BLangConstant) GetPrecedence() int {
	// migrated from BLangConstant.java:144:5
	return 0
}

func (this *BLangConstant) SetPrecedence(precedence int) {
	// migrated from BLangConstant.java:149:5
}

func (this *BLangSimpleVariable) GetName() model.IdentifierNode {
	return this.Name
}

func (this *BLangSimpleVariable) GetKind() model.NodeKind {
	return model.NodeKind_VARIABLE
}

func (this *BLangSimpleVariable) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
		return
	}
	panic("name is not a BLangIdentifier")
}

func (this *BLangMarkdownDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_MARKDOWN_DOCUMENTATION
}

func (this *BLangMarkdownDocumentation) GetDocumentationLines() []model.MarkdownDocumentationTextAttributeNode {
	result := make([]model.MarkdownDocumentationTextAttributeNode, len(this.DocumentationLines))
	for i := range this.DocumentationLines {
		result[i] = &this.DocumentationLines[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddDocumentationLine(documentationText model.MarkdownDocumentationTextAttributeNode) {
	if line, ok := documentationText.(*BLangMarkdownDocumentationLine); ok {
		this.DocumentationLines = append(this.DocumentationLines, *line)
	} else {
		panic("documentationText is not a BLangMarkdownDocumentationLine")
	}
}

func (this *BLangMarkdownDocumentation) GetParameters() []model.MarkdownDocumentationParameterAttributeNode {
	result := make([]model.MarkdownDocumentationParameterAttributeNode, len(this.Parameters))
	for i := range this.Parameters {
		result[i] = &this.Parameters[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddParameter(parameter model.MarkdownDocumentationParameterAttributeNode) {
	if param, ok := parameter.(*BLangMarkdownParameterDocumentation); ok {
		this.Parameters = append(this.Parameters, *param)
	} else {
		panic("parameter is not a BLangMarkdownParameterDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) GetReturnParameter() model.MarkdownDocumentationReturnParameterAttributeNode {
	return this.ReturnParameter
}

func (this *BLangMarkdownDocumentation) GetDeprecationDocumentation() model.MarkDownDocumentationDeprecationAttributeNode {
	return this.DeprecationDocumentation
}

func (this *BLangMarkdownDocumentation) SetReturnParameter(returnParameter model.MarkdownDocumentationReturnParameterAttributeNode) {
	if param, ok := returnParameter.(*BLangMarkdownReturnParameterDocumentation); ok {
		this.ReturnParameter = param
	} else {
		panic("returnParameter is not a BLangMarkdownReturnParameterDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) SetDeprecationDocumentation(deprecationDocumentation model.MarkDownDocumentationDeprecationAttributeNode) {
	if doc, ok := deprecationDocumentation.(*BLangMarkDownDeprecationDocumentation); ok {
		this.DeprecationDocumentation = doc
	} else {
		panic("deprecationDocumentation is not a BLangMarkDownDeprecationDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) SetDeprecatedParametersDocumentation(deprecatedParametersDocumentation model.MarkDownDocumentationDeprecatedParametersAttributeNode) {
	if doc, ok := deprecatedParametersDocumentation.(*BLangMarkDownDeprecatedParametersDocumentation); ok {
		this.DeprecatedParametersDocumentation = doc
	} else {
		panic("deprecatedParametersDocumentation is not a BLangMarkDownDeprecatedParametersDocumentation")
	}
}

func (this *BLangMarkdownDocumentation) GetDeprecatedParametersDocumentation() model.MarkDownDocumentationDeprecatedParametersAttributeNode {
	return this.DeprecatedParametersDocumentation
}

func (this *BLangMarkdownDocumentation) GetDocumentation() string {
	var lines []string
	for i := range this.DocumentationLines {
		lines = append(lines, this.DocumentationLines[i].GetText())
	}
	result := strings.Join(lines, "\n")
	return strings.ReplaceAll(result, "\r", "")
}

func (this *BLangMarkdownDocumentation) GetParameterDocumentations() map[string]model.MarkdownDocumentationParameterAttributeNode {
	result := make(map[string]model.MarkdownDocumentationParameterAttributeNode)
	for _, parameter := range this.Parameters {
		paramName := parameter.GetParameterName()
		result[paramName.GetValue()] = &parameter
	}
	return result
}

func (this *BLangMarkdownDocumentation) GetReturnParameterDocumentation() *string {
	if this.ReturnParameter == nil {
		return nil
	}
	return common.ToPointer(this.ReturnParameter.GetReturnParameterDocumentation())
}

func (this *BLangMarkdownDocumentation) GetReferences() []model.MarkdownDocumentationReferenceAttributeNode {
	result := make([]model.MarkdownDocumentationReferenceAttributeNode, len(this.References))
	for i := range this.References {
		result[i] = &this.References[i]
	}
	return result
}

func (this *BLangMarkdownDocumentation) AddReference(reference model.MarkdownDocumentationReferenceAttributeNode) {
	if ref, ok := reference.(*BLangMarkdownReferenceDocumentation); ok {
		this.References = append(this.References, *ref)
	} else {
		panic("reference is not a BLangMarkdownReferenceDocumentation")
	}
}

func (this *BLangMarkdownReferenceDocumentation) GetType() model.DocumentationReferenceType {
	return this.Type
}

func (this *BLangMarkdownReferenceDocumentation) GetKind() model.NodeKind {
	return model.NodeKind_DOCUMENTATION_REFERENCE
}

// BLangService methods

func (this *BLangService) GetName() model.IdentifierNode {
	return this.Name
}

func (this *BLangService) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangService) GetResources() []model.FunctionNode {
	return []model.FunctionNode{}
}

func (this *BLangService) IsAnonymousService() bool {
	return false
}

func (this *BLangService) GetAttachedExprs() []model.ExpressionNode {
	result := make([]model.ExpressionNode, len(this.AttachedExprs))
	for i := range this.AttachedExprs {
		result[i] = this.AttachedExprs[i]
	}
	return result
}

func (this *BLangService) GetServiceClass() model.ClassDefinition {
	return this.ServiceClass
}

func (this *BLangService) GetAbsolutePath() []model.IdentifierNode {
	return this.AbsoluteResourcePath
}

func (this *BLangService) GetServiceNameLiteral() model.LiteralNode {
	return this.ServiceNameLiteral
}

func (this *BLangService) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BLangService) AddFlag(flag model.Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangService) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(this.AnnAttachments))
	for i := range this.AnnAttachments {
		result[i] = &this.AnnAttachments[i]
	}
	return result
}

func (this *BLangService) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.AnnAttachments = append(this.AnnAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangService) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return this.MarkdownDocumentationAttachment
}

func (this *BLangService) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.MarkdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (this *BLangService) GetKind() model.NodeKind {
	return model.NodeKind_SERVICE
}

func (this *BLangFunction) GetReceiver() model.SimpleVariableNode {
	return this.Receiver
}

func (this *BLangFunction) SetReceiver(receiver model.SimpleVariableNode) {
	if rec, ok := receiver.(*BLangSimpleVariable); ok {
		this.Receiver = rec
	} else {
		panic("receiver is not a BLangSimpleVariable")
	}
}

func (this *BLangFunction) GetKind() model.NodeKind {
	return model.NodeKind_FUNCTION
}

func (this *BLangFunction) Scope() model.Scope {
	return this.scope
}

func (this *BLangFunction) SetScope(scope model.Scope) {
	this.scope = scope
}

var _ NodeWithScope = &BLangFunction{}

func (b *BLangInvokableNodeBase) GetName() model.IdentifierNode {
	return b.Name
}

func (b *BLangInvokableNodeBase) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		b.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (b *BLangInvokableNodeBase) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return b.AnnAttachments
}

func (b *BLangInvokableNodeBase) GetAnnAttachments() []model.AnnotationAttachmentNode {
	attachments := make([]model.AnnotationAttachmentNode, len(b.AnnAttachments))
	for i, attachment := range b.AnnAttachments {
		attachments[i] = attachment
	}
	return attachments
}

func (b *BLangInvokableNodeBase) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	b.AnnAttachments = append(b.AnnAttachments, annAttachment)
}

func (b *BLangInvokableNodeBase) SetAnnAttachments(annAttachments []model.AnnotationAttachmentNode) {
	b.AnnAttachments = annAttachments
}

func (b *BLangInvokableNodeBase) AddFlag(flag model.Flag) {
	b.FlagSet.Add(flag)
}

func (b *BLangInvokableNodeBase) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return b.MarkdownDocumentationAttachment
}

func (b *BLangInvokableNodeBase) SetMarkdownDocumentationAttachment(markdownDocumentationAttachment model.MarkdownDocumentationNode) {
	if doc, ok := markdownDocumentationAttachment.(*BLangMarkdownDocumentation); ok {
		b.MarkdownDocumentationAttachment = doc
	} else {
		panic("markdownDocumentationAttachment is not a BLangMarkdownDocumentation")
	}
}

func (b *BLangInvokableNodeBase) GetParameters() []model.SimpleVariableNode {
	result := make([]model.SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *BLangInvokableNodeBase) AddParameter(param model.SimpleVariableNode) {
	if blangParam, ok := param.(*BLangSimpleVariable); ok {
		b.RequiredParams = append(b.RequiredParams, *blangParam)
	} else {
		panic("param is not a BLangSimpleVariable")
	}
}

func (b *BLangInvokableNodeBase) GetRequiredParams() []model.SimpleVariableNode {
	result := make([]model.SimpleVariableNode, len(b.RequiredParams))
	for i, param := range b.RequiredParams {
		result[i] = &param
	}
	return result
}

func (b *BLangInvokableNodeBase) SetRequiredParams(requiredParams []model.SimpleVariableNode) {
	b.RequiredParams = make([]BLangSimpleVariable, len(requiredParams))
	for i, param := range requiredParams {
		if blangParam, ok := param.(*BLangSimpleVariable); ok {
			b.RequiredParams[i] = *blangParam
		} else {
			panic("requiredParams contains element that is not a BLangSimpleVariable")
		}
	}
}

func (b *BLangInvokableNodeBase) GetRestParameters() model.SimpleVariableNode {
	return b.RestParam
}

func (b *BLangInvokableNodeBase) GetRestParam() model.SimpleVariableNode {
	return b.RestParam
}

func (b *BLangInvokableNodeBase) SetRestParameter(restParam model.SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *BLangInvokableNodeBase) SetRestParam(restParam model.SimpleVariableNode) {
	b.RestParam = restParam
}

func (b *BLangInvokableNodeBase) HasBody() bool {
	return b.Body != nil
}

func (b *BLangInvokableNodeBase) GetReturnTypeData() model.TypeData {
	return b.ReturnTypeData
}

func (b *BLangInvokableNodeBase) SetReturnTypeData(returnTypeData model.TypeData) {
	b.ReturnTypeData = returnTypeData
}

func (b *BLangInvokableNodeBase) GetReturnTypeAnnotationAttachments() []model.AnnotationAttachmentNode {
	return b.ReturnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) GetReturnTypeAnnAttachments() []model.AnnotationAttachmentNode {
	return b.ReturnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) AddReturnTypeAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	b.ReturnTypeAnnAttachments = append(b.ReturnTypeAnnAttachments, annAttachment)
}

func (b *BLangInvokableNodeBase) SetReturnTypeAnnAttachments(returnTypeAnnAttachments []model.AnnotationAttachmentNode) {
	b.ReturnTypeAnnAttachments = returnTypeAnnAttachments
}

func (b *BLangInvokableNodeBase) GetBody() model.FunctionBodyNode {
	return b.Body
}

func (b *BLangInvokableNodeBase) SetBody(body model.FunctionBodyNode) {
	b.Body = body
}

func (b *BLangInvokableNodeBase) GetDefaultWorkerName() model.IdentifierNode {
	return b.DefaultWorkerName
}

func (b *BLangInvokableNodeBase) SetDefaultWorkerName(defaultWorkerName model.IdentifierNode) {
	b.DefaultWorkerName = defaultWorkerName
}

func (b *BLangInvokableNodeBase) GetFlags() common.Set[model.Flag] {
	return &b.FlagSet
}

func (b *BLangInvokableNodeBase) GetFlagSet() common.Set[model.Flag] {
	return &b.FlagSet
}

func (b *BLangInvokableNodeBase) SetFlagSet(flagSet common.Set[model.Flag]) {
	if set, ok := flagSet.(*common.UnorderedSet[model.Flag]); ok {
		b.FlagSet = *set
	} else {
		panic("flagSet is not a common.UnorderedSet[Flag]")
	}
}

func (b *BLangInvokableNodeBase) GetDesugaredReturnType() bool {
	return b.DesugaredReturnType
}

func (b *BLangInvokableNodeBase) SetDesugaredReturnType(desugaredReturnType bool) {
	b.DesugaredReturnType = desugaredReturnType
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

func (b *BLangVariableBase) GetFlagSet() common.Set[model.Flag] {
	return b.FlagSet
}

func (b *BLangVariableBase) SetFlagSet(flagSet common.Set[model.Flag]) {
	b.FlagSet = flagSet
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

func (m *BLangVariableBase) AddFlag(flag model.Flag) {
	m.FlagSet.Add(flag)
}

func (m *BLangVariableBase) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	return m.AnnAttachments
}

func (m *BLangVariableBase) GetFlags() common.Set[model.Flag] {
	return m.FlagSet
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
	this.FlagSet = common.UnorderedSet[model.Flag]{}
	this.CycleDepth = -1
	this.hasCyclicReference = false
	return this
}

func (this *BLangTypeDefinition) GetName() model.IdentifierNode {
	return this.Name
}

func (this *BLangTypeDefinition) SetName(name model.IdentifierNode) {
	if id, ok := name.(*BLangIdentifier); ok {
		this.Name = id
	} else {
		panic("name is not a BLangIdentifier")
	}
}

func (this *BLangTypeDefinition) GetTypeData() model.TypeData {
	return this.typeData
}

func (this *BLangTypeDefinition) SetTypeData(typeData model.TypeData) {
	this.typeData = typeData
}

func (this *BLangTypeDefinition) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BLangTypeDefinition) AddFlag(flag model.Flag) {
	this.FlagSet.Add(flag)
}

func (this *BLangTypeDefinition) GetAnnotationAttachments() []model.AnnotationAttachmentNode {
	result := make([]model.AnnotationAttachmentNode, len(this.annAttachments))
	for i := range this.annAttachments {
		result[i] = &this.annAttachments[i]
	}
	return result
}

func (this *BLangTypeDefinition) AddAnnotationAttachment(annAttachment model.AnnotationAttachmentNode) {
	if ann, ok := annAttachment.(*BLangAnnotationAttachment); ok {
		this.annAttachments = append(this.annAttachments, *ann)
	} else {
		panic("annAttachment is not a BLangAnnotationAttachment")
	}
}

func (this *BLangTypeDefinition) GetMarkdownDocumentationAttachment() model.MarkdownDocumentationNode {
	return this.markdownDocumentationAttachment
}

func (this *BLangTypeDefinition) SetMarkdownDocumentationAttachment(documentationNode model.MarkdownDocumentationNode) {
	if doc, ok := documentationNode.(*BLangMarkdownDocumentation); ok {
		this.markdownDocumentationAttachment = doc
	} else {
		panic("documentationNode is not a BLangMarkdownDocumentation")
	}
}

func (this *BLangTypeDefinition) GetPrecedence() int {
	return this.precedence
}

func (this *BLangTypeDefinition) SetPrecedence(precedence int) {
	this.precedence = precedence
}

func (this *BLangTypeDefinition) GetKind() model.NodeKind {
	return model.NodeKind_TYPE_DEFINITION
}

func (this *BLangXMLNS) GetNamespaceURI() model.ExpressionNode {
	return this.namespaceURI
}

func (this *BLangXMLNS) GetPrefix() model.IdentifierNode {
	return this.prefix
}

func (this *BLangXMLNS) SetNamespaceURI(namespaceURI model.ExpressionNode) {
	if expr, ok := namespaceURI.(BLangExpression); ok {
		this.namespaceURI = expr
	} else {
		panic("namespaceURI is not a BLangExpression")
	}
}

func (this *BLangXMLNS) SetPrefix(prefix model.IdentifierNode) {
	if ident, ok := prefix.(*BLangIdentifier); ok {
		this.prefix = ident
	} else {
		panic("prefix is not a BLangIdentifier")
	}
}

func (this *BLangXMLNS) GetKind() model.NodeKind {
	return model.NodeKind_XMLNS
}

func (this *BLangPackage) GetCompilationUnits() []model.CompilationUnitNode {
	result := make([]model.CompilationUnitNode, len(this.CompUnits))
	for i := range this.CompUnits {
		result[i] = &this.CompUnits[i]
	}
	return result
}

func (this *BLangPackage) AddCompilationUnit(compUnit model.CompilationUnitNode) {
	if cu, ok := compUnit.(*BLangCompilationUnit); ok {
		this.CompUnits = append(this.CompUnits, *cu)
	} else {
		panic("compUnit is not a BLangCompilationUnit")
	}
}

func (this *BLangPackage) GetImports() []model.ImportPackageNode {
	result := make([]model.ImportPackageNode, len(this.Imports))
	for i := range this.Imports {
		result[i] = &this.Imports[i]
	}
	return result
}

func (this *BLangPackage) AddImport(importPkg model.ImportPackageNode) {
	if imp, ok := importPkg.(*BLangImportPackage); ok {
		this.Imports = append(this.Imports, *imp)
	} else {
		panic("importPkg is not a BLangImportPackage")
	}
}

func (this *BLangPackage) GetNamespaceDeclarations() []model.XMLNSDeclarationNode {
	result := make([]model.XMLNSDeclarationNode, len(this.XmlnsList))
	for i := range this.XmlnsList {
		result[i] = &this.XmlnsList[i]
	}
	return result
}

func (this *BLangPackage) AddNamespaceDeclaration(xmlnsDecl model.XMLNSDeclarationNode) {
	if xmlns, ok := xmlnsDecl.(*BLangXMLNS); ok {
		this.XmlnsList = append(this.XmlnsList, *xmlns)
		this.TopLevelNodes = append(this.TopLevelNodes, xmlnsDecl)
	} else {
		panic("xmlnsDecl is not a BLangXMLNS")
	}
}

func (this *BLangPackage) GetConstants() []model.ConstantNode {
	result := make([]model.ConstantNode, len(this.Constants))
	for i := range this.Constants {
		result[i] = &this.Constants[i]
	}
	return result
}

func (this *BLangPackage) GetGlobalVariables() []model.VariableNode {
	result := make([]model.VariableNode, len(this.GlobalVars))
	for i := range this.GlobalVars {
		result[i] = &this.GlobalVars[i]
	}
	return result
}

func (this *BLangPackage) AddGlobalVariable(globalVar model.SimpleVariableNode) {
	if sv, ok := globalVar.(*BLangSimpleVariable); ok {
		this.GlobalVars = append(this.GlobalVars, *sv)
		this.TopLevelNodes = append(this.TopLevelNodes, globalVar)
	} else {
		panic("globalVar is not a BLangSimpleVariable")
	}
}

func (this *BLangPackage) GetServices() []model.ServiceNode {
	result := make([]model.ServiceNode, len(this.Services))
	for i := range this.Services {
		result[i] = &this.Services[i]
	}
	return result
}

func (this *BLangPackage) AddService(service model.ServiceNode) {
	if svc, ok := service.(*BLangService); ok {
		this.Services = append(this.Services, *svc)
		this.TopLevelNodes = append(this.TopLevelNodes, service)
	} else {
		panic("service is not a BLangService")
	}
}

func (this *BLangPackage) GetFunctions() []model.FunctionNode {
	result := make([]model.FunctionNode, len(this.Functions))
	for i := range this.Functions {
		result[i] = &this.Functions[i]
	}
	return result
}

func (this *BLangPackage) AddFunction(function model.FunctionNode) {
	if fn, ok := function.(*BLangFunction); ok {
		this.Functions = append(this.Functions, *fn)
		this.TopLevelNodes = append(this.TopLevelNodes, function)
	} else {
		panic("function is not a BLangFunction")
	}
}

func (this *BLangPackage) GetTypeDefinitions() []model.TypeDefinition {
	result := make([]model.TypeDefinition, len(this.TypeDefinitions))
	for i := range this.TypeDefinitions {
		result[i] = &this.TypeDefinitions[i]
	}
	return result
}

func (this *BLangPackage) AddTypeDefinition(typeDefinition model.TypeDefinition) {
	if td, ok := typeDefinition.(*BLangTypeDefinition); ok {
		this.TypeDefinitions = append(this.TypeDefinitions, *td)
		this.TopLevelNodes = append(this.TopLevelNodes, typeDefinition)
	} else {
		panic("typeDefinition is not a BLangTypeDefinition")
	}
}

func (this *BLangPackage) GetAnnotations() []model.AnnotationNode {
	result := make([]model.AnnotationNode, len(this.Annotations))
	for i := range this.Annotations {
		result[i] = &this.Annotations[i]
	}
	return result
}

func (this *BLangPackage) AddAnnotation(annotation model.AnnotationNode) {
	if ann, ok := annotation.(*BLangAnnotation); ok {
		this.Annotations = append(this.Annotations, *ann)
		this.TopLevelNodes = append(this.TopLevelNodes, annotation)
	} else {
		panic("annotation is not a BLangAnnotation")
	}
}

func (this *BLangPackage) GetClassDefinitions() []model.ClassDefinition {
	result := make([]model.ClassDefinition, len(this.ClassDefinitions))
	for i := range this.ClassDefinitions {
		result[i] = &this.ClassDefinitions[i]
	}
	return result
}

func (this *BLangPackage) GetKind() model.NodeKind {
	return model.NodeKind_PACKAGE
}

func (this *BLangPackage) AddTestablePkg(testablePkg *BLangTestablePackage) {
	this.TestablePkgs = append(this.TestablePkgs, testablePkg)
}

func (this *BLangPackage) GetTestablePkgs() []*BLangTestablePackage {
	return this.TestablePkgs
}

func (this *BLangPackage) GetTestablePkg() *BLangTestablePackage {
	if len(this.TestablePkgs) > 0 {
		return this.TestablePkgs[0]
	}
	return nil
}

func (this *BLangPackage) ContainsTestablePkg() bool {
	return len(this.TestablePkgs) > 0
}

func (this *BLangPackage) GetFlags() common.Set[model.Flag] {
	return &this.FlagSet
}

func (this *BLangPackage) HasTestablePackage() bool {
	return len(this.TestablePkgs) > 0
}

func (this *BLangPackage) AddClassDefinition(classDefNode *BLangClassDefinition) {
	this.TopLevelNodes = append(this.TopLevelNodes, classDefNode)
	this.ClassDefinitions = append(this.ClassDefinitions, *classDefNode)
}

func (this *BLangPackage) AddDiagnostic(diagnostic diagnostics.Diagnostic) {
	// Check if diagnostic already exists
	for _, existing := range this.diagnostics {
		if diagnosticEqual(existing, diagnostic) {
			return
		}
	}
	this.diagnostics = append(this.diagnostics, diagnostic)
	severity := diagnostic.DiagnosticInfo().Severity()
	if severity == diagnostics.Error {
		this.errorCount++
	} else if severity == diagnostics.Warning {
		this.warnCount++
	}
}

func diagnosticEqual(d1, d2 diagnostics.Diagnostic) bool {
	info1 := d1.DiagnosticInfo()
	info2 := d2.DiagnosticInfo()
	return info1.Code() == info2.Code() &&
		info1.MessageFormat() == info2.MessageFormat() &&
		info1.Severity() == info2.Severity()
}

func (this *BLangPackage) GetDiagnostics() []diagnostics.Diagnostic {
	result := make([]diagnostics.Diagnostic, len(this.diagnostics))
	copy(result, this.diagnostics)
	return result
}

func (this *BLangPackage) GetErrorCount() int {
	return this.errorCount
}

func (this *BLangPackage) GetWarnCount() int {
	return this.warnCount
}

func (this *BLangPackage) HasErrors() bool {
	return this.errorCount > 0
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
	this.FlagSet = common.UnorderedSet[model.Flag]{}
	this.CompletedPhases = common.UnorderedSet[CompilerPhase]{}
	this.LambdaFunctions = []BLangLambdaFunction{}
	this.errorCount = 0
	this.warnCount = 0
	this.diagnostics = []diagnostics.Diagnostic{}
	return this
}

func (this *BLangTestablePackage) GetMockFunctionNamesMap() map[string]string {
	return this.mockFunctionNamesMap
}

func (this *BLangTestablePackage) AddMockFunction(id string, function string) {
	if this.mockFunctionNamesMap == nil {
		this.mockFunctionNamesMap = make(map[string]string)
	}
	this.mockFunctionNamesMap[id] = function
}

func (this *BLangTestablePackage) GetIsLegacyMockingMap() map[string]bool {
	return this.isLegacyMockingMap
}

func (this *BLangTestablePackage) AddIsLegacyMockingMap(id string, isLegacy bool) {
	if this.isLegacyMockingMap == nil {
		this.isLegacyMockingMap = make(map[string]bool)
	}
	this.isLegacyMockingMap[id] = isLegacy
}

func createSimpleVariableNodeWithLocationTokenLocation(location Location, identifier tree.Token, identifierPos Location) *BLangSimpleVariable {
	memberVar := createSimpleVariableNode()
	memberVar.pos = location
	name := createIdentifierFromToken(identifierPos, identifier)
	BLangNode(&name).SetPosition(identifierPos)
	memberVar.SetName(&name)
	return memberVar
}

func createSimpleVariableNode() *BLangSimpleVariable {
	return &BLangSimpleVariable{
		BLangVariableBase: BLangVariableBase{
			FlagSet: &common.UnorderedSet[model.Flag]{},
		},
	}
}

func createConstantNode() *BLangConstant {
	return &BLangConstant{
		BLangVariableBase: BLangVariableBase{
			FlagSet: &common.UnorderedSet[model.Flag]{},
		},
	}
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
		switch node.(type) {
		case *BLangImportPackage:
			p.Imports = append(p.Imports, *node.(*BLangImportPackage))
		case *BLangConstant:
			p.Constants = append(p.Constants, *node.(*BLangConstant))
		case *BLangService:
			p.Services = append(p.Services, *node.(*BLangService))
		case *BLangFunction:
			p.Functions = append(p.Functions, *node.(*BLangFunction))
		case *BLangTypeDefinition:
			p.TypeDefinitions = append(p.TypeDefinitions, *node.(*BLangTypeDefinition))
		case *BLangAnnotation:
			p.Annotations = append(p.Annotations, *node.(*BLangAnnotation))
		default:
			p.TopLevelNodes = append(p.TopLevelNodes, node)
		}
	}
	return &p
}
