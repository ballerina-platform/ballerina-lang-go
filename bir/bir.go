package bir

import (
	"bytes"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

/**
 * @see <a href="https://github.com/ballerina-platform/ballerina-lang/blob/master/docs/compiler/bir-spec.md">Source</a>
 */

type Bir_TypeTagEnum int

const (
	Bir_TypeTagEnum__TypeTagInt               Bir_TypeTagEnum = 1
	Bir_TypeTagEnum__TypeTagByte              Bir_TypeTagEnum = 2
	Bir_TypeTagEnum__TypeTagFloat             Bir_TypeTagEnum = 3
	Bir_TypeTagEnum__TypeTagDecimal           Bir_TypeTagEnum = 4
	Bir_TypeTagEnum__TypeTagString            Bir_TypeTagEnum = 5
	Bir_TypeTagEnum__TypeTagBoolean           Bir_TypeTagEnum = 6
	Bir_TypeTagEnum__TypeTagJson              Bir_TypeTagEnum = 7
	Bir_TypeTagEnum__TypeTagXml               Bir_TypeTagEnum = 8
	Bir_TypeTagEnum__TypeTagTable             Bir_TypeTagEnum = 9
	Bir_TypeTagEnum__TypeTagNil               Bir_TypeTagEnum = 10
	Bir_TypeTagEnum__TypeTagAnydata           Bir_TypeTagEnum = 11
	Bir_TypeTagEnum__TypeTagRecord            Bir_TypeTagEnum = 12
	Bir_TypeTagEnum__TypeTagTypedesc          Bir_TypeTagEnum = 13
	Bir_TypeTagEnum__TypeTagTyperefdesc       Bir_TypeTagEnum = 14
	Bir_TypeTagEnum__TypeTagStream            Bir_TypeTagEnum = 15
	Bir_TypeTagEnum__TypeTagMap               Bir_TypeTagEnum = 16
	Bir_TypeTagEnum__TypeTagInvokable         Bir_TypeTagEnum = 17
	Bir_TypeTagEnum__TypeTagAny               Bir_TypeTagEnum = 18
	Bir_TypeTagEnum__TypeTagEndpoint          Bir_TypeTagEnum = 19
	Bir_TypeTagEnum__TypeTagArray             Bir_TypeTagEnum = 20
	Bir_TypeTagEnum__TypeTagUnion             Bir_TypeTagEnum = 21
	Bir_TypeTagEnum__TypeTagIntersection      Bir_TypeTagEnum = 22
	Bir_TypeTagEnum__TypeTagPackage           Bir_TypeTagEnum = 23
	Bir_TypeTagEnum__TypeTagNone              Bir_TypeTagEnum = 24
	Bir_TypeTagEnum__TypeTagVoid              Bir_TypeTagEnum = 25
	Bir_TypeTagEnum__TypeTagXmlns             Bir_TypeTagEnum = 26
	Bir_TypeTagEnum__TypeTagAnnotation        Bir_TypeTagEnum = 27
	Bir_TypeTagEnum__TypeTagSemanticError     Bir_TypeTagEnum = 28
	Bir_TypeTagEnum__TypeTagError             Bir_TypeTagEnum = 29
	Bir_TypeTagEnum__TypeTagIterator          Bir_TypeTagEnum = 30
	Bir_TypeTagEnum__TypeTagTuple             Bir_TypeTagEnum = 31
	Bir_TypeTagEnum__TypeTagFuture            Bir_TypeTagEnum = 32
	Bir_TypeTagEnum__TypeTagFinite            Bir_TypeTagEnum = 33
	Bir_TypeTagEnum__TypeTagObjectOrService   Bir_TypeTagEnum = 34
	Bir_TypeTagEnum__TypeTagByteArray         Bir_TypeTagEnum = 35
	Bir_TypeTagEnum__TypeTagFunctionPointer   Bir_TypeTagEnum = 36
	Bir_TypeTagEnum__TypeTagHandle            Bir_TypeTagEnum = 37
	Bir_TypeTagEnum__TypeTagReadonly          Bir_TypeTagEnum = 38
	Bir_TypeTagEnum__TypeTagSigned32Int       Bir_TypeTagEnum = 39
	Bir_TypeTagEnum__TypeTagSigned16Int       Bir_TypeTagEnum = 40
	Bir_TypeTagEnum__TypeTagSigned8Int        Bir_TypeTagEnum = 41
	Bir_TypeTagEnum__TypeTagUnsigned32Int     Bir_TypeTagEnum = 42
	Bir_TypeTagEnum__TypeTagUnsigned16Int     Bir_TypeTagEnum = 43
	Bir_TypeTagEnum__TypeTagUnsigned8Int      Bir_TypeTagEnum = 44
	Bir_TypeTagEnum__TypeTagCharString        Bir_TypeTagEnum = 45
	Bir_TypeTagEnum__TypeTagXmlElement        Bir_TypeTagEnum = 46
	Bir_TypeTagEnum__TypeTagXmlPi             Bir_TypeTagEnum = 47
	Bir_TypeTagEnum__TypeTagXmlComment        Bir_TypeTagEnum = 48
	Bir_TypeTagEnum__TypeTagXmlText           Bir_TypeTagEnum = 49
	Bir_TypeTagEnum__TypeTagNever             Bir_TypeTagEnum = 50
	Bir_TypeTagEnum__TypeTagNullSet           Bir_TypeTagEnum = 51
	Bir_TypeTagEnum__TypeTagParameterizedType Bir_TypeTagEnum = 52
	Bir_TypeTagEnum__TypeTagRegExpType        Bir_TypeTagEnum = 53
)

type Bir struct {
	ConstantPool *Bir_ConstantPoolSet
	Module       *Bir_Module
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir() *Bir {
	return &Bir{}
}

func (this Bir) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1 := NewBir_ConstantPoolSet()
	err = tmp1.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ConstantPool = tmp1
	tmp2 := NewBir_Module()
	err = tmp2.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Module = tmp2
	return err
}

type Bir_Annotation struct {
	PackageIdCpIndex             int32
	NameCpIndex                  int32
	OriginalNameCpIndex          int32
	Flags                        int64
	Origin                       int8
	Position                     *Bir_Position
	AttachPointsCount            int32
	AttachPoints                 []*Bir_AttachPoint
	AnnotationTypeCpIndex        int32
	Doc                          *Bir_Markdown
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Module
}

func NewBir_Annotation() *Bir_Annotation {
	return &Bir_Annotation{}
}

func (this Bir_Annotation) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Annotation) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp3, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PackageIdCpIndex = int32(tmp3)
	tmp4, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp4)
	tmp5, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OriginalNameCpIndex = int32(tmp5)
	tmp6, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp6)
	tmp7, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp7
	tmp8 := NewBir_Position()
	err = tmp8.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp8
	tmp9, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AttachPointsCount = int32(tmp9)
	for i := 0; i < int(this.AttachPointsCount); i++ {
		_ = i
		tmp10 := NewBir_AttachPoint()
		err = tmp10.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.AttachPoints = append(this.AttachPoints, tmp10)
	}
	tmp11, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AnnotationTypeCpIndex = int32(tmp11)
	tmp12 := NewBir_Markdown()
	err = tmp12.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp12
	tmp13 := NewBir_AnnotationAttachmentsContent()
	err = tmp13.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp13
	return err
}

type Bir_AnnotationAttachment struct {
	PackageIdCpIndex    int32
	Position            *Bir_Position
	TagReferenceCpIndex int32
	IsConstAnnot        uint8
	ConstantValue       *Bir_ConstantValue
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_AnnotationAttachmentsContent
}

func NewBir_AnnotationAttachment() *Bir_AnnotationAttachment {
	return &Bir_AnnotationAttachment{}
}

func (this Bir_AnnotationAttachment) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_AnnotationAttachment) Read(io *kaitai.Stream, parent *Bir_AnnotationAttachmentsContent, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp14, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PackageIdCpIndex = int32(tmp14)
	tmp15 := NewBir_Position()
	err = tmp15.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp15
	tmp16, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TagReferenceCpIndex = int32(tmp16)
	tmp17, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsConstAnnot = tmp17
	if this.IsConstAnnot == 1 {
		tmp18 := NewBir_ConstantValue()
		err = tmp18.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValue = tmp18
	}
	return err
}

type Bir_AnnotationAttachmentsContent struct {
	AnnotationAttachmentsContentLength int64
	AttachmentsCount                   int32
	AnnotationAttachments              []*Bir_AnnotationAttachment
	_io                                *kaitai.Stream
	_root                              *Bir
	_parent                            kaitai.Struct
}

func NewBir_AnnotationAttachmentsContent() *Bir_AnnotationAttachmentsContent {
	return &Bir_AnnotationAttachmentsContent{}
}

func (this Bir_AnnotationAttachmentsContent) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_AnnotationAttachmentsContent) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp19, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContentLength = int64(tmp19)
	tmp20, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AttachmentsCount = int32(tmp20)
	for i := 0; i < int(this.AttachmentsCount); i++ {
		_ = i
		tmp21 := NewBir_AnnotationAttachment()
		err = tmp21.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.AnnotationAttachments = append(this.AnnotationAttachments, tmp21)
	}
	return err
}

type Bir_AttachPoint struct {
	PointNameCpIndex int32
	IsSource         uint8
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Annotation
}

func NewBir_AttachPoint() *Bir_AttachPoint {
	return &Bir_AttachPoint{}
}

func (this Bir_AttachPoint) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_AttachPoint) Read(io *kaitai.Stream, parent *Bir_Annotation, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp22, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PointNameCpIndex = int32(tmp22)
	tmp23, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSource = tmp23
	return err
}

type Bir_BasicBlock struct {
	NameCpIndex       int32
	InstructionsCount int32
	Instructions      []*Bir_Instruction
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_BasicBlocksInfo
}

func NewBir_BasicBlock() *Bir_BasicBlock {
	return &Bir_BasicBlock{}
}

func (this Bir_BasicBlock) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_BasicBlock) Read(io *kaitai.Stream, parent *Bir_BasicBlocksInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp24, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp24)
	tmp25, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InstructionsCount = int32(tmp25)
	for i := 0; i < int(this.InstructionsCount); i++ {
		_ = i
		tmp26 := NewBir_Instruction()
		err = tmp26.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Instructions = append(this.Instructions, tmp26)
	}
	return err
}

// TODO: bir model generator

// unsafe cast

// bir
// - internal
// 		- types
// - types
// bir.ksy

// compiler directivies
// go generate

// write as protobuf

// type BasicBlocksInfo struct {
// 	BasicBlocksCount int32
// 	BasicBlocks      []*BasicBlock
// }

type Bir_BasicBlocksInfo struct {
	BasicBlocksCount int32
	BasicBlocks      []*Bir_BasicBlock
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_FunctionBody
}

func NewBir_BasicBlocksInfo() *Bir_BasicBlocksInfo {
	return &Bir_BasicBlocksInfo{}
}

func (this Bir_BasicBlocksInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_BasicBlocksInfo) Read(io *kaitai.Stream, parent *Bir_FunctionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp27, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.BasicBlocksCount = int32(tmp27)
	for i := 0; i < int(this.BasicBlocksCount); i++ {
		_ = i
		tmp28 := NewBir_BasicBlock()
		err = tmp28.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.BasicBlocks = append(this.BasicBlocks, tmp28)
	}
	return err
}

type Bir_BooleanConstantInfo struct {
	ValueBooleanConstant uint8
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              kaitai.Struct
}

func NewBir_BooleanConstantInfo() *Bir_BooleanConstantInfo {
	return &Bir_BooleanConstantInfo{}
}

func (this Bir_BooleanConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_BooleanConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp29, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.ValueBooleanConstant = tmp29
	return err
}

type Bir_BooleanCpInfo struct {
	Value   uint8
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolEntry
}

func NewBir_BooleanCpInfo() *Bir_BooleanCpInfo {
	return &Bir_BooleanCpInfo{}
}

func (this Bir_BooleanCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_BooleanCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp30, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Value = tmp30
	return err
}

type Bir_ByteConstantInfo struct {
	ValueCpIndex int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir_ByteConstantInfo() *Bir_ByteConstantInfo {
	return &Bir_ByteConstantInfo{}
}

func (this Bir_ByteConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ByteConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp31, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueCpIndex = int32(tmp31)
	return err
}

type Bir_ByteCpInfo struct {
	Value   int32
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolEntry
}

func NewBir_ByteCpInfo() *Bir_ByteCpInfo {
	return &Bir_ByteCpInfo{}
}

func (this Bir_ByteCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ByteCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp32, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Value = int32(tmp32)
	return err
}

type Bir_CallInstructionInfo struct {
	IsVirtual       uint8
	PackageIndex    int32
	CallNameCpIndex int32
	ArgumentsCount  int32
	Arguments       []*Bir_Operand
	HasLhsOperand   int8
	LhsOperand      *Bir_Operand
	_io             *kaitai.Stream
	_root           *Bir
	_parent         kaitai.Struct
}

func NewBir_CallInstructionInfo() *Bir_CallInstructionInfo {
	return &Bir_CallInstructionInfo{}
}

func (this Bir_CallInstructionInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_CallInstructionInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp33, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsVirtual = tmp33
	tmp34, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PackageIndex = int32(tmp34)
	tmp35, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.CallNameCpIndex = int32(tmp35)
	tmp36, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ArgumentsCount = int32(tmp36)
	for i := 0; i < int(this.ArgumentsCount); i++ {
		_ = i
		tmp37 := NewBir_Operand()
		err = tmp37.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Arguments = append(this.Arguments, tmp37)
	}
	tmp38, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasLhsOperand = tmp38
	if this.HasLhsOperand != 0 {
		tmp39 := NewBir_Operand()
		err = tmp39.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.LhsOperand = tmp39
	}
	return err
}

type Bir_ClosureSymbolBody struct {
	NameCpIndex  int32
	Flags        int64
	TypeCpIndex  int32
	PkdIdCpIndex int32
	ParamCount   int32
	Params       []*Bir_FunctionParameter
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_DefaultValueBody
}

func NewBir_ClosureSymbolBody() *Bir_ClosureSymbolBody {
	return &Bir_ClosureSymbolBody{}
}

func (this Bir_ClosureSymbolBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ClosureSymbolBody) Read(io *kaitai.Stream, parent *Bir_DefaultValueBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp40, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp40)
	tmp41, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp41)
	tmp42, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp42)
	tmp43, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkdIdCpIndex = int32(tmp43)
	tmp44, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamCount = int32(tmp44)
	for i := 0; i < int(this.ParamCount); i++ {
		_ = i
		tmp45 := NewBir_FunctionParameter()
		err = tmp45.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = append(this.Params, tmp45)
	}
	return err
}

type Bir_Constant struct {
	NameCpIndex                  int32
	Flags                        int64
	Origin                       int8
	Position                     *Bir_Position
	Doc                          *Bir_Markdown
	TypeCpIndex                  int32
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	Length                       int64
	ConstantValue                *Bir_ConstantValue
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Module
}

func NewBir_Constant() *Bir_Constant {
	return &Bir_Constant{}
}

func (this Bir_Constant) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Constant) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp46, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp46)
	tmp47, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp47)
	tmp48, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp48
	tmp49 := NewBir_Position()
	err = tmp49.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp49
	tmp50 := NewBir_Markdown()
	err = tmp50.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp50
	tmp51, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp51)
	tmp52 := NewBir_AnnotationAttachmentsContent()
	err = tmp52.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp52
	tmp53, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Length = int64(tmp53)
	tmp54 := NewBir_ConstantValue()
	err = tmp54.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ConstantValue = tmp54
	return err
}

type Bir_ConstantPoolEntry_TagEnum int

const (
	Bir_ConstantPoolEntry_TagEnum__CpEntryInteger Bir_ConstantPoolEntry_TagEnum = 1
	Bir_ConstantPoolEntry_TagEnum__CpEntryFloat   Bir_ConstantPoolEntry_TagEnum = 2
	Bir_ConstantPoolEntry_TagEnum__CpEntryBoolean Bir_ConstantPoolEntry_TagEnum = 3
	Bir_ConstantPoolEntry_TagEnum__CpEntryString  Bir_ConstantPoolEntry_TagEnum = 4
	Bir_ConstantPoolEntry_TagEnum__CpEntryPackage Bir_ConstantPoolEntry_TagEnum = 5
	Bir_ConstantPoolEntry_TagEnum__CpEntryByte    Bir_ConstantPoolEntry_TagEnum = 6
	Bir_ConstantPoolEntry_TagEnum__CpEntryShape   Bir_ConstantPoolEntry_TagEnum = 7
)

type Bir_ConstantPoolEntry struct {
	Tag     Bir_ConstantPoolEntry_TagEnum
	CpInfo  kaitai.Struct
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolSet
}

func NewBir_ConstantPoolEntry() *Bir_ConstantPoolEntry {
	return &Bir_ConstantPoolEntry{}
}

func (this Bir_ConstantPoolEntry) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ConstantPoolEntry) Read(io *kaitai.Stream, parent *Bir_ConstantPoolSet, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp55, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Tag = Bir_ConstantPoolEntry_TagEnum(tmp55)
	switch this.Tag {
	case Bir_ConstantPoolEntry_TagEnum__CpEntryBoolean:
		tmp56 := NewBir_BooleanCpInfo()
		err = tmp56.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp56
	case Bir_ConstantPoolEntry_TagEnum__CpEntryByte:
		tmp57 := NewBir_ByteCpInfo()
		err = tmp57.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp57
	case Bir_ConstantPoolEntry_TagEnum__CpEntryFloat:
		tmp58 := NewBir_FloatCpInfo()
		err = tmp58.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp58
	case Bir_ConstantPoolEntry_TagEnum__CpEntryInteger:
		tmp59 := NewBir_IntCpInfo()
		err = tmp59.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp59
	case Bir_ConstantPoolEntry_TagEnum__CpEntryPackage:
		tmp60 := NewBir_PackageCpInfo()
		err = tmp60.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp60
	case Bir_ConstantPoolEntry_TagEnum__CpEntryShape:
		tmp61 := NewBir_ShapeCpInfo()
		err = tmp61.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp61
	case Bir_ConstantPoolEntry_TagEnum__CpEntryString:
		tmp62 := NewBir_StringCpInfo()
		err = tmp62.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CpInfo = tmp62
	}
	return err
}

type Bir_ConstantPoolSet struct {
	Magic               []byte
	Version             int32
	ConstantPoolCount   int32
	ConstantPoolEntries []*Bir_ConstantPoolEntry
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir
}

func NewBir_ConstantPoolSet() *Bir_ConstantPoolSet {
	return &Bir_ConstantPoolSet{}
}

func (this Bir_ConstantPoolSet) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ConstantPoolSet) Read(io *kaitai.Stream, parent *Bir, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp63, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	this.Magic = tmp63
	if !(bytes.Equal(this.Magic, []uint8{186, 16, 192, 222})) {
		return kaitai.NewValidationNotEqualError([]uint8{186, 16, 192, 222}, this.Magic, this._io, "/types/constant_pool_set/seq/0")
	}
	tmp64, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Version = int32(tmp64)
	tmp65, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstantPoolCount = int32(tmp65)
	for i := 0; i < int(this.ConstantPoolCount); i++ {
		_ = i
		tmp66 := NewBir_ConstantPoolEntry()
		err = tmp66.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantPoolEntries = append(this.ConstantPoolEntries, tmp66)
	}
	return err
}

type Bir_ConstantValue struct {
	ConstantValueTypeCpIndex int32
	ConstantValueInfo        kaitai.Struct
	_io                      *kaitai.Stream
	_root                    *Bir
	_parent                  kaitai.Struct
	_f_type                  bool
	t                        *Bir_ShapeCpInfo
}

func NewBir_ConstantValue() *Bir_ConstantValue {
	return &Bir_ConstantValue{}
}

func (this Bir_ConstantValue) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ConstantValue) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp67, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstantValueTypeCpIndex = int32(tmp67)
	tmp68, err := this.Type()
	if err != nil {
		return err
	}
	switch tmp68.Shape.TypeTag {
	case Bir_TypeTagEnum__TypeTagBoolean:
		tmp69 := NewBir_BooleanConstantInfo()
		err = tmp69.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp69
	case Bir_TypeTagEnum__TypeTagByte:
		tmp70 := NewBir_ByteConstantInfo()
		err = tmp70.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp70
	case Bir_TypeTagEnum__TypeTagCharString:
		tmp71 := NewBir_StringConstantInfo()
		err = tmp71.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp71
	case Bir_TypeTagEnum__TypeTagDecimal:
		tmp72 := NewBir_DecimalConstantInfo()
		err = tmp72.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp72
	case Bir_TypeTagEnum__TypeTagFloat:
		tmp73 := NewBir_FloatConstantInfo()
		err = tmp73.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp73
	case Bir_TypeTagEnum__TypeTagInt:
		tmp74 := NewBir_IntConstantInfo()
		err = tmp74.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp74
	case Bir_TypeTagEnum__TypeTagIntersection:
		tmp75 := NewBir_IntersectionConstantInfo()
		err = tmp75.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp75
	case Bir_TypeTagEnum__TypeTagNil:
		tmp76 := NewBir_NilConstantInfo()
		err = tmp76.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp76
	case Bir_TypeTagEnum__TypeTagRecord:
		tmp77 := NewBir_MapConstantInfo()
		err = tmp77.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp77
	case Bir_TypeTagEnum__TypeTagSigned16Int:
		tmp78 := NewBir_IntConstantInfo()
		err = tmp78.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp78
	case Bir_TypeTagEnum__TypeTagSigned32Int:
		tmp79 := NewBir_IntConstantInfo()
		err = tmp79.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp79
	case Bir_TypeTagEnum__TypeTagSigned8Int:
		tmp80 := NewBir_IntConstantInfo()
		err = tmp80.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp80
	case Bir_TypeTagEnum__TypeTagString:
		tmp81 := NewBir_StringConstantInfo()
		err = tmp81.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp81
	case Bir_TypeTagEnum__TypeTagUnsigned16Int:
		tmp82 := NewBir_IntConstantInfo()
		err = tmp82.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp82
	case Bir_TypeTagEnum__TypeTagUnsigned32Int:
		tmp83 := NewBir_IntConstantInfo()
		err = tmp83.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp83
	case Bir_TypeTagEnum__TypeTagUnsigned8Int:
		tmp84 := NewBir_IntConstantInfo()
		err = tmp84.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp84
	}
	return err
}

func (this *Bir_ConstantValue) Type() (v *Bir_ShapeCpInfo, err error) {
	if this._f_type {
		return this.t, nil
	}
	this._f_type = true
	this.t = this._root.ConstantPool.ConstantPoolEntries[this.ConstantValueTypeCpIndex].CpInfo.(*Bir_ShapeCpInfo)
	return this.t, nil
}

type Bir_DecimalConstantInfo struct {
	ValueCpIndex int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir_DecimalConstantInfo() *Bir_DecimalConstantInfo {
	return &Bir_DecimalConstantInfo{}
}

func (this Bir_DecimalConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_DecimalConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp85, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueCpIndex = int32(tmp85)
	return err
}

type Bir_DefaultParameter struct {
	Kind               int8
	TypeCpIndex        int32
	NameCpIndex        int32
	MetaVarNameCpIndex int32
	HasDefaultExpr     uint8
	_io                *kaitai.Stream
	_root              *Bir
	_parent            *Bir_FunctionBody
}

func NewBir_DefaultParameter() *Bir_DefaultParameter {
	return &Bir_DefaultParameter{}
}

func (this Bir_DefaultParameter) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_DefaultParameter) Read(io *kaitai.Stream, parent *Bir_FunctionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp86, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp86
	tmp87, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp87)
	tmp88, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp88)
	if this.Kind == 2 {
		tmp89, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.MetaVarNameCpIndex = int32(tmp89)
	}
	tmp90, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasDefaultExpr = tmp90
	return err
}

type Bir_DefaultValueBody struct {
	ParamNameCpIndex int32
	ClosureSymbol    *Bir_ClosureSymbolBody
	_io              *kaitai.Stream
	_root            *Bir
	_parent          kaitai.Struct
}

func NewBir_DefaultValueBody() *Bir_DefaultValueBody {
	return &Bir_DefaultValueBody{}
}

func (this Bir_DefaultValueBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_DefaultValueBody) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp91, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamNameCpIndex = int32(tmp91)
	tmp92 := NewBir_ClosureSymbolBody()
	err = tmp92.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ClosureSymbol = tmp92
	return err
}

type Bir_EnclosingBasicBlockId struct {
	MetaVarNameCpIndex int32
	EndBbIdCpIndex     int32
	StartBbIdCpIndex   int32
	InstructionOffset  int32
	_io                *kaitai.Stream
	_root              *Bir
	_parent            *Bir_LocalVariable
}

func NewBir_EnclosingBasicBlockId() *Bir_EnclosingBasicBlockId {
	return &Bir_EnclosingBasicBlockId{}
}

func (this Bir_EnclosingBasicBlockId) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_EnclosingBasicBlockId) Read(io *kaitai.Stream, parent *Bir_LocalVariable, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp93, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.MetaVarNameCpIndex = int32(tmp93)
	tmp94, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.EndBbIdCpIndex = int32(tmp94)
	tmp95, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.StartBbIdCpIndex = int32(tmp95)
	tmp96, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InstructionOffset = int32(tmp96)
	return err
}

type Bir_ErrorEntry struct {
	TrapBbIdCpIndex   int32
	EndBbIdCpIndex    int32
	ErrorOperand      *Bir_Operand
	TargetBbIdCpIndex int32
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_ErrorTable
}

func NewBir_ErrorEntry() *Bir_ErrorEntry {
	return &Bir_ErrorEntry{}
}

func (this Bir_ErrorEntry) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ErrorEntry) Read(io *kaitai.Stream, parent *Bir_ErrorTable, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp97, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TrapBbIdCpIndex = int32(tmp97)
	tmp98, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.EndBbIdCpIndex = int32(tmp98)
	tmp99 := NewBir_Operand()
	err = tmp99.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ErrorOperand = tmp99
	tmp100, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TargetBbIdCpIndex = int32(tmp100)
	return err
}

type Bir_ErrorTable struct {
	ErrorEntriesCount int32
	ErrorEntries      []*Bir_ErrorEntry
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_FunctionBody
}

func NewBir_ErrorTable() *Bir_ErrorTable {
	return &Bir_ErrorTable{}
}

func (this Bir_ErrorTable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ErrorTable) Read(io *kaitai.Stream, parent *Bir_FunctionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp101, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ErrorEntriesCount = int32(tmp101)
	for i := 0; i < int(this.ErrorEntriesCount); i++ {
		_ = i
		tmp102 := NewBir_ErrorEntry()
		err = tmp102.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ErrorEntries = append(this.ErrorEntries, tmp102)
	}
	return err
}

type Bir_ExternalTypeDefintionInfo struct {
	ExternalPkgIdCpIndex int32
	ObjectNameCpIndex    int32
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_InstructionNewInstance
}

func NewBir_ExternalTypeDefintionInfo() *Bir_ExternalTypeDefintionInfo {
	return &Bir_ExternalTypeDefintionInfo{}
}

func (this Bir_ExternalTypeDefintionInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ExternalTypeDefintionInfo) Read(io *kaitai.Stream, parent *Bir_InstructionNewInstance, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp103, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ExternalPkgIdCpIndex = int32(tmp103)
	tmp104, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ObjectNameCpIndex = int32(tmp104)
	return err
}

type Bir_FloatConstantInfo struct {
	ValueCpIndex int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir_FloatConstantInfo() *Bir_FloatConstantInfo {
	return &Bir_FloatConstantInfo{}
}

func (this Bir_FloatConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_FloatConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp105, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueCpIndex = int32(tmp105)
	return err
}

type Bir_FloatCpInfo struct {
	Value   float64
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolEntry
}

func NewBir_FloatCpInfo() *Bir_FloatCpInfo {
	return &Bir_FloatCpInfo{}
}

func (this Bir_FloatCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_FloatCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp106, err := this._io.ReadF8be()
	if err != nil {
		return err
	}
	this.Value = float64(tmp106)
	return err
}

type Bir_FpLoadFunctionParam struct {
	Kind        int8
	TypeCpIndex int32
	NameCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_InstructionFpLoad
}

func NewBir_FpLoadFunctionParam() *Bir_FpLoadFunctionParam {
	return &Bir_FpLoadFunctionParam{}
}

func (this Bir_FpLoadFunctionParam) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_FpLoadFunctionParam) Read(io *kaitai.Stream, parent *Bir_InstructionFpLoad, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp107, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp107
	tmp108, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp108)
	tmp109, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp109)
	return err
}

type Bir_Function struct {
	Position                     *Bir_Position
	NameCpIndex                  int32
	OriginalNameCpIndex          int32
	WorkerNameCpIndex            int32
	Flags                        int64
	Origin                       int8
	TypeCpIndex                  int32
	IsResourceFunction           uint8
	ResourceFunctionContent      *Bir_ResourceFunctionContent
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	ReturnTypeAnnotations        *Bir_AnnotationAttachmentsContent
	RequiredParamCount           int32
	RequiredParams               []*Bir_RequiredParam
	HasRestParam                 uint8
	RestParamNameCpIndex         int32
	RestParamAnnotations         *Bir_AnnotationAttachmentsContent
	HasReceiver                  uint8
	Reciever                     *Bir_Reciever
	Doc                          *Bir_Markdown
	DependentGlobalVarLength     int32
	DependentGlobalVarCpEntry    []int32
	ScopeTableLength             int64
	ScopeEntryCount              int32
	ScopeEntries                 []*Bir_ScopeEntry
	FunctionBodyLength           int64
	FunctionBody                 *Bir_FunctionBody
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      kaitai.Struct
	_raw_FunctionBody            []byte
}

func NewBir_Function() *Bir_Function {
	return &Bir_Function{}
}

func (this Bir_Function) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Function) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp110 := NewBir_Position()
	err = tmp110.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp110
	tmp111, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp111)
	tmp112, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OriginalNameCpIndex = int32(tmp112)
	tmp113, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.WorkerNameCpIndex = int32(tmp113)
	tmp114, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp114)
	tmp115, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp115
	tmp116, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp116)
	tmp117, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsResourceFunction = tmp117
	if this.IsResourceFunction == 1 {
		tmp118 := NewBir_ResourceFunctionContent()
		err = tmp118.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ResourceFunctionContent = tmp118
	}
	tmp119 := NewBir_AnnotationAttachmentsContent()
	err = tmp119.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp119
	tmp120 := NewBir_AnnotationAttachmentsContent()
	err = tmp120.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ReturnTypeAnnotations = tmp120
	tmp121, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.RequiredParamCount = int32(tmp121)
	for i := 0; i < int(this.RequiredParamCount); i++ {
		_ = i
		tmp122 := NewBir_RequiredParam()
		err = tmp122.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RequiredParams = append(this.RequiredParams, tmp122)
	}
	tmp123, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasRestParam = tmp123
	if this.HasRestParam != 0 {
		tmp124, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.RestParamNameCpIndex = int32(tmp124)
	}
	if this.HasRestParam != 0 {
		tmp125 := NewBir_AnnotationAttachmentsContent()
		err = tmp125.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RestParamAnnotations = tmp125
	}
	tmp126, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasReceiver = tmp126
	if this.HasReceiver != 0 {
		tmp127 := NewBir_Reciever()
		err = tmp127.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Reciever = tmp127
	}
	tmp128 := NewBir_Markdown()
	err = tmp128.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp128
	tmp129, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DependentGlobalVarLength = int32(tmp129)
	for i := 0; i < int(this.DependentGlobalVarLength); i++ {
		_ = i
		tmp130, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.DependentGlobalVarCpEntry = append(this.DependentGlobalVarCpEntry, tmp130)
	}
	tmp131, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.ScopeTableLength = int64(tmp131)
	tmp132, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ScopeEntryCount = int32(tmp132)
	for i := 0; i < int(this.ScopeEntryCount); i++ {
		_ = i
		tmp133 := NewBir_ScopeEntry()
		err = tmp133.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ScopeEntries = append(this.ScopeEntries, tmp133)
	}
	tmp134, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.FunctionBodyLength = int64(tmp134)
	tmp135, err := this._io.ReadBytes(int(this.FunctionBodyLength))
	if err != nil {
		return err
	}
	this._raw_FunctionBody = tmp135
	_io__raw_FunctionBody := kaitai.NewStream(bytes.NewReader(this._raw_FunctionBody))
	tmp136 := NewBir_FunctionBody()
	err = tmp136.Read(_io__raw_FunctionBody, this, this._root)
	if err != nil {
		return err
	}
	this.FunctionBody = tmp136
	return err
}

type Bir_FunctionBody struct {
	ArgsCount               int32
	HasReturnVar            uint8
	ReturnVar               *Bir_ReturnVar
	DefaultParameterCount   int32
	DefaultParameters       []*Bir_DefaultParameter
	LocalVariablesCount     int32
	LocalVariables          []*Bir_LocalVariable
	FunctionBasicBlocksInfo *Bir_BasicBlocksInfo
	ErrorTable              *Bir_ErrorTable
	WorkerChannelInfo       *Bir_WorkerChannel
	_io                     *kaitai.Stream
	_root                   *Bir
	_parent                 *Bir_Function
}

func NewBir_FunctionBody() *Bir_FunctionBody {
	return &Bir_FunctionBody{}
}

func (this Bir_FunctionBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_FunctionBody) Read(io *kaitai.Stream, parent *Bir_Function, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp137, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ArgsCount = int32(tmp137)
	tmp138, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasReturnVar = tmp138
	if this.HasReturnVar == 1 {
		tmp139 := NewBir_ReturnVar()
		err = tmp139.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ReturnVar = tmp139
	}
	tmp140, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DefaultParameterCount = int32(tmp140)
	for i := 0; i < int(this.DefaultParameterCount); i++ {
		_ = i
		tmp141 := NewBir_DefaultParameter()
		err = tmp141.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DefaultParameters = append(this.DefaultParameters, tmp141)
	}
	tmp142, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.LocalVariablesCount = int32(tmp142)
	for i := 0; i < int(this.LocalVariablesCount); i++ {
		_ = i
		tmp143 := NewBir_LocalVariable()
		err = tmp143.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.LocalVariables = append(this.LocalVariables, tmp143)
	}
	tmp144 := NewBir_BasicBlocksInfo()
	err = tmp144.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FunctionBasicBlocksInfo = tmp144
	tmp145 := NewBir_ErrorTable()
	err = tmp145.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ErrorTable = tmp145
	tmp146 := NewBir_WorkerChannel()
	err = tmp146.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.WorkerChannelInfo = tmp146
	return err
}

type Bir_FunctionParameter struct {
	NameCpIndex int32
	Flags       int64
	Doc         *Bir_Markdown
	TypeCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     kaitai.Struct
}

func NewBir_FunctionParameter() *Bir_FunctionParameter {
	return &Bir_FunctionParameter{}
}

func (this Bir_FunctionParameter) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_FunctionParameter) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp147, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp147)
	tmp148, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp148)
	tmp149 := NewBir_Markdown()
	err = tmp149.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp149
	tmp150, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp150)
	return err
}

type Bir_GlobalVar struct {
	Position                     *Bir_Position
	Kind                         int8
	NameCpIndex                  int32
	Flags                        int64
	Origin                       int8
	Doc                          *Bir_Markdown
	TypeCpIndex                  int32
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Module
}

func NewBir_GlobalVar() *Bir_GlobalVar {
	return &Bir_GlobalVar{}
}

func (this Bir_GlobalVar) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_GlobalVar) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp151 := NewBir_Position()
	err = tmp151.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp151
	tmp152, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp152
	tmp153, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp153)
	tmp154, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp154)
	tmp155, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp155
	tmp156 := NewBir_Markdown()
	err = tmp156.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp156
	tmp157, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp157)
	tmp158 := NewBir_AnnotationAttachmentsContent()
	err = tmp158.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp158
	return err
}

type Bir_GlobalVariable struct {
	PackageIndex int32
	TypeCpIndex  int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Variable
}

func NewBir_GlobalVariable() *Bir_GlobalVariable {
	return &Bir_GlobalVariable{}
}

func (this Bir_GlobalVariable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_GlobalVariable) Read(io *kaitai.Stream, parent *Bir_Variable, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp159, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PackageIndex = int32(tmp159)
	tmp160, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp160)
	return err
}

type Bir_IndexAccess struct {
	LhsOperand *Bir_Operand
	KeyOperand *Bir_Operand
	RhsOperand *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    kaitai.Struct
}

func NewBir_IndexAccess() *Bir_IndexAccess {
	return &Bir_IndexAccess{}
}

func (this Bir_IndexAccess) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_IndexAccess) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp161 := NewBir_Operand()
	err = tmp161.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp161
	tmp162 := NewBir_Operand()
	err = tmp162.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.KeyOperand = tmp162
	tmp163 := NewBir_Operand()
	err = tmp163.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp163
	return err
}

type Bir_Instruction_InstructionKindEnum int

const (
	Bir_Instruction_InstructionKindEnum__InstructionKindGoto                      Bir_Instruction_InstructionKindEnum = 1
	Bir_Instruction_InstructionKindEnum__InstructionKindCall                      Bir_Instruction_InstructionKindEnum = 2
	Bir_Instruction_InstructionKindEnum__InstructionKindBranch                    Bir_Instruction_InstructionKindEnum = 3
	Bir_Instruction_InstructionKindEnum__InstructionKindReturn                    Bir_Instruction_InstructionKindEnum = 4
	Bir_Instruction_InstructionKindEnum__InstructionKindAsyncCall                 Bir_Instruction_InstructionKindEnum = 5
	Bir_Instruction_InstructionKindEnum__InstructionKindWait                      Bir_Instruction_InstructionKindEnum = 6
	Bir_Instruction_InstructionKindEnum__InstructionKindFpCall                    Bir_Instruction_InstructionKindEnum = 7
	Bir_Instruction_InstructionKindEnum__InstructionKindWkReceive                 Bir_Instruction_InstructionKindEnum = 8
	Bir_Instruction_InstructionKindEnum__InstructionKindWkSend                    Bir_Instruction_InstructionKindEnum = 9
	Bir_Instruction_InstructionKindEnum__InstructionKindFlush                     Bir_Instruction_InstructionKindEnum = 10
	Bir_Instruction_InstructionKindEnum__InstructionKindLock                      Bir_Instruction_InstructionKindEnum = 11
	Bir_Instruction_InstructionKindEnum__InstructionKindFieldLock                 Bir_Instruction_InstructionKindEnum = 12
	Bir_Instruction_InstructionKindEnum__InstructionKindUnlock                    Bir_Instruction_InstructionKindEnum = 13
	Bir_Instruction_InstructionKindEnum__InstructionKindWaitAll                   Bir_Instruction_InstructionKindEnum = 14
	Bir_Instruction_InstructionKindEnum__InstructionKindMove                      Bir_Instruction_InstructionKindEnum = 20
	Bir_Instruction_InstructionKindEnum__InstructionKindConstLoad                 Bir_Instruction_InstructionKindEnum = 21
	Bir_Instruction_InstructionKindEnum__InstructionKindNewStructure              Bir_Instruction_InstructionKindEnum = 22
	Bir_Instruction_InstructionKindEnum__InstructionKindMapStore                  Bir_Instruction_InstructionKindEnum = 23
	Bir_Instruction_InstructionKindEnum__InstructionKindMapLoad                   Bir_Instruction_InstructionKindEnum = 24
	Bir_Instruction_InstructionKindEnum__InstructionKindNewArray                  Bir_Instruction_InstructionKindEnum = 25
	Bir_Instruction_InstructionKindEnum__InstructionKindArrayStore                Bir_Instruction_InstructionKindEnum = 26
	Bir_Instruction_InstructionKindEnum__InstructionKindArrayLoad                 Bir_Instruction_InstructionKindEnum = 27
	Bir_Instruction_InstructionKindEnum__InstructionKindNewError                  Bir_Instruction_InstructionKindEnum = 28
	Bir_Instruction_InstructionKindEnum__InstructionKindTypeCast                  Bir_Instruction_InstructionKindEnum = 29
	Bir_Instruction_InstructionKindEnum__InstructionKindIsLike                    Bir_Instruction_InstructionKindEnum = 30
	Bir_Instruction_InstructionKindEnum__InstructionKindTypeTest                  Bir_Instruction_InstructionKindEnum = 31
	Bir_Instruction_InstructionKindEnum__InstructionKindNewInstance               Bir_Instruction_InstructionKindEnum = 32
	Bir_Instruction_InstructionKindEnum__InstructionKindObjectStore               Bir_Instruction_InstructionKindEnum = 33
	Bir_Instruction_InstructionKindEnum__InstructionKindObjectLoad                Bir_Instruction_InstructionKindEnum = 34
	Bir_Instruction_InstructionKindEnum__InstructionKindPanic                     Bir_Instruction_InstructionKindEnum = 35
	Bir_Instruction_InstructionKindEnum__InstructionKindFpLoad                    Bir_Instruction_InstructionKindEnum = 36
	Bir_Instruction_InstructionKindEnum__InstructionKindStringLoad                Bir_Instruction_InstructionKindEnum = 37
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlElement             Bir_Instruction_InstructionKindEnum = 38
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlText                Bir_Instruction_InstructionKindEnum = 39
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlComment             Bir_Instruction_InstructionKindEnum = 40
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlPi                  Bir_Instruction_InstructionKindEnum = 41
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlSequence            Bir_Instruction_InstructionKindEnum = 42
	Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlQname               Bir_Instruction_InstructionKindEnum = 43
	Bir_Instruction_InstructionKindEnum__InstructionKindNewStringXmlQname         Bir_Instruction_InstructionKindEnum = 44
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlSeqStore               Bir_Instruction_InstructionKindEnum = 45
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlSeqLoad                Bir_Instruction_InstructionKindEnum = 46
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlLoad                   Bir_Instruction_InstructionKindEnum = 47
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlLoadAll                Bir_Instruction_InstructionKindEnum = 48
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlAttributeLoad          Bir_Instruction_InstructionKindEnum = 49
	Bir_Instruction_InstructionKindEnum__InstructionKindXmlAttributeStore         Bir_Instruction_InstructionKindEnum = 50
	Bir_Instruction_InstructionKindEnum__InstructionKindNewTable                  Bir_Instruction_InstructionKindEnum = 51
	Bir_Instruction_InstructionKindEnum__InstructionKindNewTypedesc               Bir_Instruction_InstructionKindEnum = 52
	Bir_Instruction_InstructionKindEnum__InstructionKindNewStream                 Bir_Instruction_InstructionKindEnum = 53
	Bir_Instruction_InstructionKindEnum__InstructionKindTableStore                Bir_Instruction_InstructionKindEnum = 54
	Bir_Instruction_InstructionKindEnum__InstructionKindTableLoad                 Bir_Instruction_InstructionKindEnum = 55
	Bir_Instruction_InstructionKindEnum__InstructionKindAdd                       Bir_Instruction_InstructionKindEnum = 61
	Bir_Instruction_InstructionKindEnum__InstructionKindSub                       Bir_Instruction_InstructionKindEnum = 62
	Bir_Instruction_InstructionKindEnum__InstructionKindMul                       Bir_Instruction_InstructionKindEnum = 63
	Bir_Instruction_InstructionKindEnum__InstructionKindDiv                       Bir_Instruction_InstructionKindEnum = 64
	Bir_Instruction_InstructionKindEnum__InstructionKindMod                       Bir_Instruction_InstructionKindEnum = 65
	Bir_Instruction_InstructionKindEnum__InstructionKindEqual                     Bir_Instruction_InstructionKindEnum = 66
	Bir_Instruction_InstructionKindEnum__InstructionKindNotEqual                  Bir_Instruction_InstructionKindEnum = 67
	Bir_Instruction_InstructionKindEnum__InstructionKindGreaterThan               Bir_Instruction_InstructionKindEnum = 68
	Bir_Instruction_InstructionKindEnum__InstructionKindGreaterEqual              Bir_Instruction_InstructionKindEnum = 69
	Bir_Instruction_InstructionKindEnum__InstructionKindLessThan                  Bir_Instruction_InstructionKindEnum = 70
	Bir_Instruction_InstructionKindEnum__InstructionKindLessEqual                 Bir_Instruction_InstructionKindEnum = 71
	Bir_Instruction_InstructionKindEnum__InstructionKindAnd                       Bir_Instruction_InstructionKindEnum = 72
	Bir_Instruction_InstructionKindEnum__InstructionKindOr                        Bir_Instruction_InstructionKindEnum = 73
	Bir_Instruction_InstructionKindEnum__InstructionKindRefEqual                  Bir_Instruction_InstructionKindEnum = 74
	Bir_Instruction_InstructionKindEnum__InstructionKindRefNotEqual               Bir_Instruction_InstructionKindEnum = 75
	Bir_Instruction_InstructionKindEnum__InstructionKindClosedRange               Bir_Instruction_InstructionKindEnum = 76
	Bir_Instruction_InstructionKindEnum__InstructionKindHalfOpenRange             Bir_Instruction_InstructionKindEnum = 77
	Bir_Instruction_InstructionKindEnum__InstructionKindAnnotAccess               Bir_Instruction_InstructionKindEnum = 78
	Bir_Instruction_InstructionKindEnum__InstructionKindTypeof                    Bir_Instruction_InstructionKindEnum = 80
	Bir_Instruction_InstructionKindEnum__InstructionKindNot                       Bir_Instruction_InstructionKindEnum = 81
	Bir_Instruction_InstructionKindEnum__InstructionKindNegate                    Bir_Instruction_InstructionKindEnum = 82
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseAnd                Bir_Instruction_InstructionKindEnum = 83
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseOr                 Bir_Instruction_InstructionKindEnum = 84
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseXor                Bir_Instruction_InstructionKindEnum = 85
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseLeftShift          Bir_Instruction_InstructionKindEnum = 86
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseRightShift         Bir_Instruction_InstructionKindEnum = 87
	Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseUnsignedRightShift Bir_Instruction_InstructionKindEnum = 88
	Bir_Instruction_InstructionKindEnum__InstructionKindNewRegExp                 Bir_Instruction_InstructionKindEnum = 89
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReDisjunction          Bir_Instruction_InstructionKindEnum = 90
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReSequence             Bir_Instruction_InstructionKindEnum = 91
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReAssertion            Bir_Instruction_InstructionKindEnum = 92
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReAtomQuantifier       Bir_Instruction_InstructionKindEnum = 93
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReLiteralCharEscape    Bir_Instruction_InstructionKindEnum = 94
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharClass            Bir_Instruction_InstructionKindEnum = 95
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharSet              Bir_Instruction_InstructionKindEnum = 96
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharSetRange         Bir_Instruction_InstructionKindEnum = 97
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReCapturingGroup       Bir_Instruction_InstructionKindEnum = 98
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReFlagExpr             Bir_Instruction_InstructionKindEnum = 99
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReFlagOnOff            Bir_Instruction_InstructionKindEnum = 100
	Bir_Instruction_InstructionKindEnum__InstructionKindNewReQuantifier           Bir_Instruction_InstructionKindEnum = 101
	Bir_Instruction_InstructionKindEnum__InstructionKindRecordDefaultFpLoad       Bir_Instruction_InstructionKindEnum = 102
	Bir_Instruction_InstructionKindEnum__InstructionKindWkAltReceive              Bir_Instruction_InstructionKindEnum = 103
	Bir_Instruction_InstructionKindEnum__InstructionKindWkMulReceive              Bir_Instruction_InstructionKindEnum = 104
	Bir_Instruction_InstructionKindEnum__InstructionKindPlatform                  Bir_Instruction_InstructionKindEnum = 128
)

type Bir_Instruction struct {
	Position             *Bir_Position
	InstructionKind      Bir_Instruction_InstructionKindEnum
	InstructionStructure kaitai.Struct
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_BasicBlock
}

func NewBir_Instruction() *Bir_Instruction {
	return &Bir_Instruction{}
}

func (this Bir_Instruction) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Instruction) Read(io *kaitai.Stream, parent *Bir_BasicBlock, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp164 := NewBir_Position()
	err = tmp164.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp164
	tmp165, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.InstructionKind = Bir_Instruction_InstructionKindEnum(tmp165)
	switch this.InstructionKind {
	case Bir_Instruction_InstructionKindEnum__InstructionKindAdd:
		tmp166 := NewBir_InstructionBinaryOperation()
		err = tmp166.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp166
	case Bir_Instruction_InstructionKindEnum__InstructionKindAnd:
		tmp167 := NewBir_InstructionBinaryOperation()
		err = tmp167.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp167
	case Bir_Instruction_InstructionKindEnum__InstructionKindAnnotAccess:
		tmp168 := NewBir_InstructionBinaryOperation()
		err = tmp168.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp168
	case Bir_Instruction_InstructionKindEnum__InstructionKindArrayLoad:
		tmp169 := NewBir_InstructionArrayLoad()
		err = tmp169.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp169
	case Bir_Instruction_InstructionKindEnum__InstructionKindArrayStore:
		tmp170 := NewBir_InstructionArrayStore()
		err = tmp170.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp170
	case Bir_Instruction_InstructionKindEnum__InstructionKindAsyncCall:
		tmp171 := NewBir_InstructionAsyncCall()
		err = tmp171.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp171
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseAnd:
		tmp172 := NewBir_InstructionBinaryOperation()
		err = tmp172.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp172
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseLeftShift:
		tmp173 := NewBir_InstructionBinaryOperation()
		err = tmp173.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp173
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseOr:
		tmp174 := NewBir_InstructionBinaryOperation()
		err = tmp174.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp174
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseRightShift:
		tmp175 := NewBir_InstructionBinaryOperation()
		err = tmp175.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp175
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseUnsignedRightShift:
		tmp176 := NewBir_InstructionBinaryOperation()
		err = tmp176.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp176
	case Bir_Instruction_InstructionKindEnum__InstructionKindBitwiseXor:
		tmp177 := NewBir_InstructionBinaryOperation()
		err = tmp177.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp177
	case Bir_Instruction_InstructionKindEnum__InstructionKindBranch:
		tmp178 := NewBir_InstructionBranch()
		err = tmp178.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp178
	case Bir_Instruction_InstructionKindEnum__InstructionKindCall:
		tmp179 := NewBir_InstructionCall()
		err = tmp179.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp179
	case Bir_Instruction_InstructionKindEnum__InstructionKindClosedRange:
		tmp180 := NewBir_InstructionBinaryOperation()
		err = tmp180.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp180
	case Bir_Instruction_InstructionKindEnum__InstructionKindConstLoad:
		tmp181 := NewBir_InstructionConstLoad()
		err = tmp181.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp181
	case Bir_Instruction_InstructionKindEnum__InstructionKindDiv:
		tmp182 := NewBir_InstructionBinaryOperation()
		err = tmp182.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp182
	case Bir_Instruction_InstructionKindEnum__InstructionKindEqual:
		tmp183 := NewBir_InstructionBinaryOperation()
		err = tmp183.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp183
	case Bir_Instruction_InstructionKindEnum__InstructionKindFieldLock:
		tmp184 := NewBir_InstructionFieldLock()
		err = tmp184.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp184
	case Bir_Instruction_InstructionKindEnum__InstructionKindFlush:
		tmp185 := NewBir_InstructionFlush()
		err = tmp185.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp185
	case Bir_Instruction_InstructionKindEnum__InstructionKindFpCall:
		tmp186 := NewBir_InstructionFpCall()
		err = tmp186.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp186
	case Bir_Instruction_InstructionKindEnum__InstructionKindFpLoad:
		tmp187 := NewBir_InstructionFpLoad()
		err = tmp187.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp187
	case Bir_Instruction_InstructionKindEnum__InstructionKindGoto:
		tmp188 := NewBir_InstructionGoto()
		err = tmp188.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp188
	case Bir_Instruction_InstructionKindEnum__InstructionKindGreaterEqual:
		tmp189 := NewBir_InstructionBinaryOperation()
		err = tmp189.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp189
	case Bir_Instruction_InstructionKindEnum__InstructionKindGreaterThan:
		tmp190 := NewBir_InstructionBinaryOperation()
		err = tmp190.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp190
	case Bir_Instruction_InstructionKindEnum__InstructionKindHalfOpenRange:
		tmp191 := NewBir_InstructionBinaryOperation()
		err = tmp191.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp191
	case Bir_Instruction_InstructionKindEnum__InstructionKindIsLike:
		tmp192 := NewBir_InstructionIsLike()
		err = tmp192.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp192
	case Bir_Instruction_InstructionKindEnum__InstructionKindLessEqual:
		tmp193 := NewBir_InstructionBinaryOperation()
		err = tmp193.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp193
	case Bir_Instruction_InstructionKindEnum__InstructionKindLessThan:
		tmp194 := NewBir_InstructionBinaryOperation()
		err = tmp194.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp194
	case Bir_Instruction_InstructionKindEnum__InstructionKindLock:
		tmp195 := NewBir_InstructionLock()
		err = tmp195.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp195
	case Bir_Instruction_InstructionKindEnum__InstructionKindMapLoad:
		tmp196 := NewBir_InstructionMapLoad()
		err = tmp196.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp196
	case Bir_Instruction_InstructionKindEnum__InstructionKindMapStore:
		tmp197 := NewBir_InstructionMapStore()
		err = tmp197.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp197
	case Bir_Instruction_InstructionKindEnum__InstructionKindMod:
		tmp198 := NewBir_InstructionBinaryOperation()
		err = tmp198.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp198
	case Bir_Instruction_InstructionKindEnum__InstructionKindMove:
		tmp199 := NewBir_InstructionMove()
		err = tmp199.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp199
	case Bir_Instruction_InstructionKindEnum__InstructionKindMul:
		tmp200 := NewBir_InstructionBinaryOperation()
		err = tmp200.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp200
	case Bir_Instruction_InstructionKindEnum__InstructionKindNegate:
		tmp201 := NewBir_InstructionUnaryOperation()
		err = tmp201.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp201
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewArray:
		tmp202 := NewBir_InstructionNewArray()
		err = tmp202.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp202
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewError:
		tmp203 := NewBir_InstructionNewError()
		err = tmp203.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp203
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewInstance:
		tmp204 := NewBir_InstructionNewInstance()
		err = tmp204.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp204
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReAssertion:
		tmp205 := NewBir_InstructionNewReAssertion()
		err = tmp205.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp205
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReAtomQuantifier:
		tmp206 := NewBir_InstructionNewReAtomQuantifier()
		err = tmp206.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp206
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReCapturingGroup:
		tmp207 := NewBir_InstructionNewReCapturingGroup()
		err = tmp207.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp207
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharClass:
		tmp208 := NewBir_InstructionNewReCharClass()
		err = tmp208.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp208
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharSet:
		tmp209 := NewBir_InstructionNewReCharSet()
		err = tmp209.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp209
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReCharSetRange:
		tmp210 := NewBir_InstructionNewReCharSetRange()
		err = tmp210.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp210
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReDisjunction:
		tmp211 := NewBir_InstructionNewReDisjunction()
		err = tmp211.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp211
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReFlagExpr:
		tmp212 := NewBir_InstructionNewReFlagExpr()
		err = tmp212.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp212
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReFlagOnOff:
		tmp213 := NewBir_InstructionNewReFlagOnOff()
		err = tmp213.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp213
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReLiteralCharEscape:
		tmp214 := NewBir_InstructionNewReCharEscape()
		err = tmp214.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp214
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReQuantifier:
		tmp215 := NewBir_InstructionNewReQuantifier()
		err = tmp215.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp215
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewReSequence:
		tmp216 := NewBir_InstructionNewReSequence()
		err = tmp216.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp216
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewRegExp:
		tmp217 := NewBir_InstructionNewRegExp()
		err = tmp217.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp217
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewStringXmlQname:
		tmp218 := NewBir_InstructionNewStringXmlQname()
		err = tmp218.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp218
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewStructure:
		tmp219 := NewBir_InstructionNewStructure()
		err = tmp219.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp219
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewTable:
		tmp220 := NewBir_InstructionNewTable()
		err = tmp220.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp220
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewTypedesc:
		tmp221 := NewBir_InstructionNewTypedesc()
		err = tmp221.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp221
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlComment:
		tmp222 := NewBir_InstructionNewXmlComment()
		err = tmp222.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp222
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlElement:
		tmp223 := NewBir_InstructionNewXmlElement()
		err = tmp223.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp223
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlPi:
		tmp224 := NewBir_InstructionNewXmlProcessIns()
		err = tmp224.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp224
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlQname:
		tmp225 := NewBir_InstructionNewXmlQname()
		err = tmp225.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp225
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlSequence:
		tmp226 := NewBir_InstructionNewXmlSequence()
		err = tmp226.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp226
	case Bir_Instruction_InstructionKindEnum__InstructionKindNewXmlText:
		tmp227 := NewBir_InstructionNewXmlText()
		err = tmp227.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp227
	case Bir_Instruction_InstructionKindEnum__InstructionKindNot:
		tmp228 := NewBir_InstructionUnaryOperation()
		err = tmp228.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp228
	case Bir_Instruction_InstructionKindEnum__InstructionKindNotEqual:
		tmp229 := NewBir_InstructionBinaryOperation()
		err = tmp229.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp229
	case Bir_Instruction_InstructionKindEnum__InstructionKindObjectLoad:
		tmp230 := NewBir_InstructionObjectLoad()
		err = tmp230.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp230
	case Bir_Instruction_InstructionKindEnum__InstructionKindObjectStore:
		tmp231 := NewBir_InstructionObjectStore()
		err = tmp231.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp231
	case Bir_Instruction_InstructionKindEnum__InstructionKindOr:
		tmp232 := NewBir_InstructionBinaryOperation()
		err = tmp232.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp232
	case Bir_Instruction_InstructionKindEnum__InstructionKindPanic:
		tmp233 := NewBir_InstructionPanic()
		err = tmp233.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp233
	case Bir_Instruction_InstructionKindEnum__InstructionKindRecordDefaultFpLoad:
		tmp234 := NewBir_InstructionRecordDefaultFpLoad()
		err = tmp234.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp234
	case Bir_Instruction_InstructionKindEnum__InstructionKindRefEqual:
		tmp235 := NewBir_InstructionBinaryOperation()
		err = tmp235.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp235
	case Bir_Instruction_InstructionKindEnum__InstructionKindRefNotEqual:
		tmp236 := NewBir_InstructionBinaryOperation()
		err = tmp236.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp236
	case Bir_Instruction_InstructionKindEnum__InstructionKindReturn:
		tmp237 := NewBir_InstructionReturn()
		err = tmp237.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp237
	case Bir_Instruction_InstructionKindEnum__InstructionKindStringLoad:
		tmp238 := NewBir_InstructionStringLoad()
		err = tmp238.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp238
	case Bir_Instruction_InstructionKindEnum__InstructionKindSub:
		tmp239 := NewBir_InstructionBinaryOperation()
		err = tmp239.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp239
	case Bir_Instruction_InstructionKindEnum__InstructionKindTableLoad:
		tmp240 := NewBir_InstructionTableLoad()
		err = tmp240.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp240
	case Bir_Instruction_InstructionKindEnum__InstructionKindTableStore:
		tmp241 := NewBir_InstructionTableStore()
		err = tmp241.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp241
	case Bir_Instruction_InstructionKindEnum__InstructionKindTypeCast:
		tmp242 := NewBir_InstructionTypeCast()
		err = tmp242.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp242
	case Bir_Instruction_InstructionKindEnum__InstructionKindTypeTest:
		tmp243 := NewBir_InstructionTypeTest()
		err = tmp243.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp243
	case Bir_Instruction_InstructionKindEnum__InstructionKindTypeof:
		tmp244 := NewBir_InstructionUnaryOperation()
		err = tmp244.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp244
	case Bir_Instruction_InstructionKindEnum__InstructionKindUnlock:
		tmp245 := NewBir_InstructionUnlock()
		err = tmp245.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp245
	case Bir_Instruction_InstructionKindEnum__InstructionKindWait:
		tmp246 := NewBir_InstructionWait()
		err = tmp246.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp246
	case Bir_Instruction_InstructionKindEnum__InstructionKindWaitAll:
		tmp247 := NewBir_InstructionWaitAll()
		err = tmp247.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp247
	case Bir_Instruction_InstructionKindEnum__InstructionKindWkAltReceive:
		tmp248 := NewBir_InstructionWkAltReceive()
		err = tmp248.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp248
	case Bir_Instruction_InstructionKindEnum__InstructionKindWkMulReceive:
		tmp249 := NewBir_InstructionWkMulReceive()
		err = tmp249.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp249
	case Bir_Instruction_InstructionKindEnum__InstructionKindWkReceive:
		tmp250 := NewBir_InstructionWkReceive()
		err = tmp250.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp250
	case Bir_Instruction_InstructionKindEnum__InstructionKindWkSend:
		tmp251 := NewBir_InstructionWkSend()
		err = tmp251.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp251
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlAttributeLoad:
		tmp252 := NewBir_InstructionXmlAttributeLoad()
		err = tmp252.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp252
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlAttributeStore:
		tmp253 := NewBir_InstructionXmlAttributeStore()
		err = tmp253.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp253
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlLoad:
		tmp254 := NewBir_InstructionXmlLoad()
		err = tmp254.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp254
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlLoadAll:
		tmp255 := NewBir_InstructionXmlLoadAll()
		err = tmp255.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp255
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlSeqLoad:
		tmp256 := NewBir_InstructionXmlSeqLoad()
		err = tmp256.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp256
	case Bir_Instruction_InstructionKindEnum__InstructionKindXmlSeqStore:
		tmp257 := NewBir_InstructionXmlSeqStore()
		err = tmp257.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InstructionStructure = tmp257
	}
	return err
}

type Bir_InstructionArrayLoad struct {
	IsOptionalFieldAccess uint8
	IsFillingRead         uint8
	ArrayLoad             *Bir_IndexAccess
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_Instruction
}

func NewBir_InstructionArrayLoad() *Bir_InstructionArrayLoad {
	return &Bir_InstructionArrayLoad{}
}

func (this Bir_InstructionArrayLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionArrayLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp258, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsOptionalFieldAccess = tmp258
	tmp259, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsFillingRead = tmp259
	tmp260 := NewBir_IndexAccess()
	err = tmp260.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ArrayLoad = tmp260
	return err
}

type Bir_InstructionArrayStore struct {
	ArrayStore *Bir_IndexAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionArrayStore() *Bir_InstructionArrayStore {
	return &Bir_InstructionArrayStore{}
}

func (this Bir_InstructionArrayStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionArrayStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp261 := NewBir_IndexAccess()
	err = tmp261.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ArrayStore = tmp261
	return err
}

type Bir_InstructionAsyncCall struct {
	CallInstructionInfo          *Bir_CallInstructionInfo
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	ThenBbIdNameCpIndex          int32
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Instruction
}

func NewBir_InstructionAsyncCall() *Bir_InstructionAsyncCall {
	return &Bir_InstructionAsyncCall{}
}

func (this Bir_InstructionAsyncCall) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionAsyncCall) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp262 := NewBir_CallInstructionInfo()
	err = tmp262.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CallInstructionInfo = tmp262
	tmp263 := NewBir_AnnotationAttachmentsContent()
	err = tmp263.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp263
	tmp264, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp264)
	return err
}

type Bir_InstructionBinaryOperation struct {
	RhsOperandOne *Bir_Operand
	RhsOperandTwo *Bir_Operand
	LhsOperand    *Bir_Operand
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_Instruction
}

func NewBir_InstructionBinaryOperation() *Bir_InstructionBinaryOperation {
	return &Bir_InstructionBinaryOperation{}
}

func (this Bir_InstructionBinaryOperation) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionBinaryOperation) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp265 := NewBir_Operand()
	err = tmp265.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperandOne = tmp265
	tmp266 := NewBir_Operand()
	err = tmp266.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperandTwo = tmp266
	tmp267 := NewBir_Operand()
	err = tmp267.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp267
	return err
}

type Bir_InstructionBranch struct {
	BranchOperand        *Bir_Operand
	TrueBbIdNameCpIndex  int32
	FalseBbIdNameCpIndex int32
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_Instruction
}

func NewBir_InstructionBranch() *Bir_InstructionBranch {
	return &Bir_InstructionBranch{}
}

func (this Bir_InstructionBranch) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionBranch) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp268 := NewBir_Operand()
	err = tmp268.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.BranchOperand = tmp268
	tmp269, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TrueBbIdNameCpIndex = int32(tmp269)
	tmp270, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FalseBbIdNameCpIndex = int32(tmp270)
	return err
}

type Bir_InstructionCall struct {
	CallInstructionInfo *Bir_CallInstructionInfo
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionCall() *Bir_InstructionCall {
	return &Bir_InstructionCall{}
}

func (this Bir_InstructionCall) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionCall) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp271 := NewBir_CallInstructionInfo()
	err = tmp271.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CallInstructionInfo = tmp271
	tmp272, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp272)
	return err
}

type Bir_InstructionConstLoad struct {
	TypeCpIndex       int32
	LhsOperand        *Bir_Operand
	ConstantValueInfo kaitai.Struct
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_Instruction
	_f_type           bool
	t                 *Bir_ShapeCpInfo
}

func NewBir_InstructionConstLoad() *Bir_InstructionConstLoad {
	return &Bir_InstructionConstLoad{}
}

func (this Bir_InstructionConstLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionConstLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp273, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp273)
	tmp274 := NewBir_Operand()
	err = tmp274.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp274
	tmp275, err := this.Type()
	if err != nil {
		return err
	}
	switch tmp275.Shape.TypeTag {
	case Bir_TypeTagEnum__TypeTagBoolean:
		tmp276 := NewBir_BooleanConstantInfo()
		err = tmp276.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp276
	case Bir_TypeTagEnum__TypeTagByte:
		tmp277 := NewBir_ByteConstantInfo()
		err = tmp277.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp277
	case Bir_TypeTagEnum__TypeTagCharString:
		tmp278 := NewBir_StringConstantInfo()
		err = tmp278.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp278
	case Bir_TypeTagEnum__TypeTagDecimal:
		tmp279 := NewBir_DecimalConstantInfo()
		err = tmp279.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp279
	case Bir_TypeTagEnum__TypeTagFloat:
		tmp280 := NewBir_FloatConstantInfo()
		err = tmp280.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp280
	case Bir_TypeTagEnum__TypeTagInt:
		tmp281 := NewBir_IntConstantInfo()
		err = tmp281.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp281
	case Bir_TypeTagEnum__TypeTagNil:
		tmp282 := NewBir_NilConstantInfo()
		err = tmp282.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp282
	case Bir_TypeTagEnum__TypeTagSigned16Int:
		tmp283 := NewBir_IntConstantInfo()
		err = tmp283.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp283
	case Bir_TypeTagEnum__TypeTagSigned32Int:
		tmp284 := NewBir_IntConstantInfo()
		err = tmp284.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp284
	case Bir_TypeTagEnum__TypeTagSigned8Int:
		tmp285 := NewBir_IntConstantInfo()
		err = tmp285.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp285
	case Bir_TypeTagEnum__TypeTagString:
		tmp286 := NewBir_StringConstantInfo()
		err = tmp286.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp286
	case Bir_TypeTagEnum__TypeTagUnsigned16Int:
		tmp287 := NewBir_IntConstantInfo()
		err = tmp287.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp287
	case Bir_TypeTagEnum__TypeTagUnsigned32Int:
		tmp288 := NewBir_IntConstantInfo()
		err = tmp288.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp288
	case Bir_TypeTagEnum__TypeTagUnsigned8Int:
		tmp289 := NewBir_IntConstantInfo()
		err = tmp289.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp289
	}
	return err
}

func (this *Bir_InstructionConstLoad) Type() (v *Bir_ShapeCpInfo, err error) {
	if this._f_type {
		return this.t, nil
	}
	this._f_type = true
	this.t = this._root.ConstantPool.ConstantPoolEntries[this.TypeCpIndex].CpInfo.(*Bir_ShapeCpInfo)
	return this.t, nil
}

type Bir_InstructionFieldLock struct {
	LockVarNameCpIndex  int32
	FieldNameCpIndex    int32
	LockBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionFieldLock() *Bir_InstructionFieldLock {
	return &Bir_InstructionFieldLock{}
}

func (this Bir_InstructionFieldLock) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionFieldLock) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp290, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.LockVarNameCpIndex = int32(tmp290)
	tmp291, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FieldNameCpIndex = int32(tmp291)
	tmp292, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.LockBbIdNameCpIndex = int32(tmp292)
	return err
}

type Bir_InstructionFlush struct {
	WorkerChannelDetail *Bir_WorkerChannel
	LhsOperand          *Bir_Operand
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionFlush() *Bir_InstructionFlush {
	return &Bir_InstructionFlush{}
}

func (this Bir_InstructionFlush) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionFlush) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp293 := NewBir_WorkerChannel()
	err = tmp293.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.WorkerChannelDetail = tmp293
	tmp294 := NewBir_Operand()
	err = tmp294.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp294
	tmp295, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp295)
	return err
}

type Bir_InstructionFpCall struct {
	FpOperand                    *Bir_Operand
	FpArgumentsCount             int32
	FpArguments                  []*Bir_Operand
	HasLhsOperand                int8
	LhsOperand                   *Bir_Operand
	IsAsynch                     uint8
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	ThenBbIdNameCpIndex          int32
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Instruction
}

func NewBir_InstructionFpCall() *Bir_InstructionFpCall {
	return &Bir_InstructionFpCall{}
}

func (this Bir_InstructionFpCall) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionFpCall) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp296 := NewBir_Operand()
	err = tmp296.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FpOperand = tmp296
	tmp297, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FpArgumentsCount = int32(tmp297)
	for i := 0; i < int(this.FpArgumentsCount); i++ {
		_ = i
		tmp298 := NewBir_Operand()
		err = tmp298.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.FpArguments = append(this.FpArguments, tmp298)
	}
	tmp299, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasLhsOperand = tmp299
	if this.HasLhsOperand == 1 {
		tmp300 := NewBir_Operand()
		err = tmp300.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.LhsOperand = tmp300
	}
	tmp301, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsAsynch = tmp301
	tmp302 := NewBir_AnnotationAttachmentsContent()
	err = tmp302.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp302
	tmp303, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp303)
	return err
}

type Bir_InstructionFpLoad struct {
	LhsOperand                *Bir_Operand
	PkgIndexCpIndex           int32
	FunctionNameCpIndex       int32
	ReturnTypeCpIndex         int32
	ClosureMapsSize           int32
	ClosureMapOperand         []*Bir_Operand
	FpLoadFunctionParamsCount int32
	FpLoadFunctionParams      []*Bir_FpLoadFunctionParam
	_io                       *kaitai.Stream
	_root                     *Bir
	_parent                   *Bir_Instruction
}

func NewBir_InstructionFpLoad() *Bir_InstructionFpLoad {
	return &Bir_InstructionFpLoad{}
}

func (this Bir_InstructionFpLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionFpLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp304 := NewBir_Operand()
	err = tmp304.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp304
	tmp305, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkgIndexCpIndex = int32(tmp305)
	tmp306, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FunctionNameCpIndex = int32(tmp306)
	tmp307, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ReturnTypeCpIndex = int32(tmp307)
	tmp308, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ClosureMapsSize = int32(tmp308)
	for i := 0; i < int(this.ClosureMapsSize); i++ {
		_ = i
		tmp309 := NewBir_Operand()
		err = tmp309.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ClosureMapOperand = append(this.ClosureMapOperand, tmp309)
	}
	tmp310, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FpLoadFunctionParamsCount = int32(tmp310)
	for i := 0; i < int(this.FpLoadFunctionParamsCount); i++ {
		_ = i
		tmp311 := NewBir_FpLoadFunctionParam()
		err = tmp311.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.FpLoadFunctionParams = append(this.FpLoadFunctionParams, tmp311)
	}
	return err
}

type Bir_InstructionGoto struct {
	TargetBbIdNameCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_Instruction
}

func NewBir_InstructionGoto() *Bir_InstructionGoto {
	return &Bir_InstructionGoto{}
}

func (this Bir_InstructionGoto) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionGoto) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp312, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TargetBbIdNameCpIndex = int32(tmp312)
	return err
}

type Bir_InstructionIsLike struct {
	TypeCpIndex int32
	LhsOperand  *Bir_Operand
	RhsOperand  *Bir_Operand
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionIsLike() *Bir_InstructionIsLike {
	return &Bir_InstructionIsLike{}
}

func (this Bir_InstructionIsLike) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionIsLike) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp313, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp313)
	tmp314 := NewBir_Operand()
	err = tmp314.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp314
	tmp315 := NewBir_Operand()
	err = tmp315.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp315
	return err
}

type Bir_InstructionLock struct {
	LockBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionLock() *Bir_InstructionLock {
	return &Bir_InstructionLock{}
}

func (this Bir_InstructionLock) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionLock) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp316, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.LockBbIdNameCpIndex = int32(tmp316)
	return err
}

type Bir_InstructionMapLoad struct {
	IsOptionalFieldAccess uint8
	IsFillingRead         uint8
	MapLoad               *Bir_IndexAccess
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_Instruction
}

func NewBir_InstructionMapLoad() *Bir_InstructionMapLoad {
	return &Bir_InstructionMapLoad{}
}

func (this Bir_InstructionMapLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionMapLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp317, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsOptionalFieldAccess = tmp317
	tmp318, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsFillingRead = tmp318
	tmp319 := NewBir_IndexAccess()
	err = tmp319.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.MapLoad = tmp319
	return err
}

type Bir_InstructionMapStore struct {
	MapStore *Bir_IndexAccess
	_io      *kaitai.Stream
	_root    *Bir
	_parent  *Bir_Instruction
}

func NewBir_InstructionMapStore() *Bir_InstructionMapStore {
	return &Bir_InstructionMapStore{}
}

func (this Bir_InstructionMapStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionMapStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp320 := NewBir_IndexAccess()
	err = tmp320.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.MapStore = tmp320
	return err
}

type Bir_InstructionMove struct {
	RhsOperand *Bir_Operand
	LhsOperand *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionMove() *Bir_InstructionMove {
	return &Bir_InstructionMove{}
}

func (this Bir_InstructionMove) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionMove) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp321 := NewBir_Operand()
	err = tmp321.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp321
	tmp322 := NewBir_Operand()
	err = tmp322.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp322
	return err
}

type Bir_InstructionNewArray struct {
	TypeCpIndex               int32
	LhsOperand                *Bir_Operand
	HasTypedescOperand        int8
	TypedescOperand           *Bir_Operand
	HasElementTypedescOperand int8
	ElementTypedescOperand    *Bir_Operand
	SizeOperand               *Bir_Operand
	InitValuesCount           int32
	InitValues                []*Bir_Operand
	_io                       *kaitai.Stream
	_root                     *Bir
	_parent                   *Bir_Instruction
}

func NewBir_InstructionNewArray() *Bir_InstructionNewArray {
	return &Bir_InstructionNewArray{}
}

func (this Bir_InstructionNewArray) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewArray) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp323, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp323)
	tmp324 := NewBir_Operand()
	err = tmp324.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp324
	tmp325, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasTypedescOperand = tmp325
	if this.HasTypedescOperand == 1 {
		tmp326 := NewBir_Operand()
		err = tmp326.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypedescOperand = tmp326
	}
	tmp327, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasElementTypedescOperand = tmp327
	if this.HasElementTypedescOperand == 1 {
		tmp328 := NewBir_Operand()
		err = tmp328.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ElementTypedescOperand = tmp328
	}
	tmp329 := NewBir_Operand()
	err = tmp329.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.SizeOperand = tmp329
	tmp330, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InitValuesCount = int32(tmp330)
	for i := 0; i < int(this.InitValuesCount); i++ {
		_ = i
		tmp331 := NewBir_Operand()
		err = tmp331.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InitValues = append(this.InitValues, tmp331)
	}
	return err
}

type Bir_InstructionNewError struct {
	ErrorTypeCpIndex int32
	LhsOperand       *Bir_Operand
	MessageOperand   *Bir_Operand
	CauseOperand     *Bir_Operand
	DetailOperand    *Bir_Operand
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Instruction
}

func NewBir_InstructionNewError() *Bir_InstructionNewError {
	return &Bir_InstructionNewError{}
}

func (this Bir_InstructionNewError) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewError) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp332, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ErrorTypeCpIndex = int32(tmp332)
	tmp333 := NewBir_Operand()
	err = tmp333.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp333
	tmp334 := NewBir_Operand()
	err = tmp334.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.MessageOperand = tmp334
	tmp335 := NewBir_Operand()
	err = tmp335.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CauseOperand = tmp335
	tmp336 := NewBir_Operand()
	err = tmp336.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.DetailOperand = tmp336
	return err
}

type Bir_InstructionNewInstance struct {
	IsExternalDefinition      uint8
	ExternalTypeDefintionInfo *Bir_ExternalTypeDefintionInfo
	DefinitionIndex           int32
	LhsOperand                *Bir_Operand
	_io                       *kaitai.Stream
	_root                     *Bir
	_parent                   *Bir_Instruction
}

func NewBir_InstructionNewInstance() *Bir_InstructionNewInstance {
	return &Bir_InstructionNewInstance{}
}

func (this Bir_InstructionNewInstance) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewInstance) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp337, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsExternalDefinition = tmp337
	if this.IsExternalDefinition == 1 {
		tmp338 := NewBir_ExternalTypeDefintionInfo()
		err = tmp338.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ExternalTypeDefintionInfo = tmp338
	}
	if this.IsExternalDefinition == 0 {
		tmp339, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.DefinitionIndex = int32(tmp339)
	}
	tmp340 := NewBir_Operand()
	err = tmp340.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp340
	return err
}

type Bir_InstructionNewReAssertion struct {
	Assertion *Bir_Operand
	LhsOp     *Bir_Operand
	_io       *kaitai.Stream
	_root     *Bir
	_parent   *Bir_Instruction
}

func NewBir_InstructionNewReAssertion() *Bir_InstructionNewReAssertion {
	return &Bir_InstructionNewReAssertion{}
}

func (this Bir_InstructionNewReAssertion) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReAssertion) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp341 := NewBir_Operand()
	err = tmp341.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Assertion = tmp341
	tmp342 := NewBir_Operand()
	err = tmp342.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp342
	return err
}

type Bir_InstructionNewReAtomQuantifier struct {
	LhsOp      *Bir_Operand
	Atom       *Bir_Operand
	Quantifier *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionNewReAtomQuantifier() *Bir_InstructionNewReAtomQuantifier {
	return &Bir_InstructionNewReAtomQuantifier{}
}

func (this Bir_InstructionNewReAtomQuantifier) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReAtomQuantifier) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp343 := NewBir_Operand()
	err = tmp343.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp343
	tmp344 := NewBir_Operand()
	err = tmp344.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Atom = tmp344
	tmp345 := NewBir_Operand()
	err = tmp345.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Quantifier = tmp345
	return err
}

type Bir_InstructionNewReCapturingGroup struct {
	LhsOp         *Bir_Operand
	OpenParen     *Bir_Operand
	FlagExpr      *Bir_Operand
	ReDisjunction *Bir_Operand
	CloseParen    *Bir_Operand
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_Instruction
}

func NewBir_InstructionNewReCapturingGroup() *Bir_InstructionNewReCapturingGroup {
	return &Bir_InstructionNewReCapturingGroup{}
}

func (this Bir_InstructionNewReCapturingGroup) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReCapturingGroup) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp346 := NewBir_Operand()
	err = tmp346.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp346
	tmp347 := NewBir_Operand()
	err = tmp347.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.OpenParen = tmp347
	tmp348 := NewBir_Operand()
	err = tmp348.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FlagExpr = tmp348
	tmp349 := NewBir_Operand()
	err = tmp349.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ReDisjunction = tmp349
	tmp350 := NewBir_Operand()
	err = tmp350.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CloseParen = tmp350
	return err
}

type Bir_InstructionNewReCharClass struct {
	LhsOp      *Bir_Operand
	ClassStart *Bir_Operand
	Negation   *Bir_Operand
	CharSet    *Bir_Operand
	ClassEnd   *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionNewReCharClass() *Bir_InstructionNewReCharClass {
	return &Bir_InstructionNewReCharClass{}
}

func (this Bir_InstructionNewReCharClass) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReCharClass) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp351 := NewBir_Operand()
	err = tmp351.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp351
	tmp352 := NewBir_Operand()
	err = tmp352.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ClassStart = tmp352
	tmp353 := NewBir_Operand()
	err = tmp353.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Negation = tmp353
	tmp354 := NewBir_Operand()
	err = tmp354.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CharSet = tmp354
	tmp355 := NewBir_Operand()
	err = tmp355.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ClassEnd = tmp355
	return err
}

type Bir_InstructionNewReCharEscape struct {
	LhsOp        *Bir_Operand
	CharOrEscape *Bir_Operand
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Instruction
}

func NewBir_InstructionNewReCharEscape() *Bir_InstructionNewReCharEscape {
	return &Bir_InstructionNewReCharEscape{}
}

func (this Bir_InstructionNewReCharEscape) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReCharEscape) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp356 := NewBir_Operand()
	err = tmp356.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp356
	tmp357 := NewBir_Operand()
	err = tmp357.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CharOrEscape = tmp357
	return err
}

type Bir_InstructionNewReCharSet struct {
	LhsOp        *Bir_Operand
	CharSetAtoms *Bir_Operand
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Instruction
}

func NewBir_InstructionNewReCharSet() *Bir_InstructionNewReCharSet {
	return &Bir_InstructionNewReCharSet{}
}

func (this Bir_InstructionNewReCharSet) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReCharSet) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp358 := NewBir_Operand()
	err = tmp358.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp358
	tmp359 := NewBir_Operand()
	err = tmp359.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.CharSetAtoms = tmp359
	return err
}

type Bir_InstructionNewReCharSetRange struct {
	LhsOp          *Bir_Operand
	LhsCharSetAtom *Bir_Operand
	Dash           *Bir_Operand
	RhsCharSetAtom *Bir_Operand
	_io            *kaitai.Stream
	_root          *Bir
	_parent        *Bir_Instruction
}

func NewBir_InstructionNewReCharSetRange() *Bir_InstructionNewReCharSetRange {
	return &Bir_InstructionNewReCharSetRange{}
}

func (this Bir_InstructionNewReCharSetRange) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReCharSetRange) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp360 := NewBir_Operand()
	err = tmp360.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp360
	tmp361 := NewBir_Operand()
	err = tmp361.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsCharSetAtom = tmp361
	tmp362 := NewBir_Operand()
	err = tmp362.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Dash = tmp362
	tmp363 := NewBir_Operand()
	err = tmp363.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsCharSetAtom = tmp363
	return err
}

type Bir_InstructionNewReDisjunction struct {
	Sequences *Bir_Operand
	LhsOp     *Bir_Operand
	_io       *kaitai.Stream
	_root     *Bir
	_parent   *Bir_Instruction
}

func NewBir_InstructionNewReDisjunction() *Bir_InstructionNewReDisjunction {
	return &Bir_InstructionNewReDisjunction{}
}

func (this Bir_InstructionNewReDisjunction) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReDisjunction) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp364 := NewBir_Operand()
	err = tmp364.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Sequences = tmp364
	tmp365 := NewBir_Operand()
	err = tmp365.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp365
	return err
}

type Bir_InstructionNewReFlagExpr struct {
	LhsOp        *Bir_Operand
	QuestionMark *Bir_Operand
	FlagsOnOff   *Bir_Operand
	Colon        *Bir_Operand
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Instruction
}

func NewBir_InstructionNewReFlagExpr() *Bir_InstructionNewReFlagExpr {
	return &Bir_InstructionNewReFlagExpr{}
}

func (this Bir_InstructionNewReFlagExpr) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReFlagExpr) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp366 := NewBir_Operand()
	err = tmp366.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp366
	tmp367 := NewBir_Operand()
	err = tmp367.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.QuestionMark = tmp367
	tmp368 := NewBir_Operand()
	err = tmp368.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.FlagsOnOff = tmp368
	tmp369 := NewBir_Operand()
	err = tmp369.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Colon = tmp369
	return err
}

type Bir_InstructionNewReFlagOnOff struct {
	LhsOp   *Bir_Operand
	Flags   *Bir_Operand
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_Instruction
}

func NewBir_InstructionNewReFlagOnOff() *Bir_InstructionNewReFlagOnOff {
	return &Bir_InstructionNewReFlagOnOff{}
}

func (this Bir_InstructionNewReFlagOnOff) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReFlagOnOff) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp370 := NewBir_Operand()
	err = tmp370.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp370
	tmp371 := NewBir_Operand()
	err = tmp371.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Flags = tmp371
	return err
}

type Bir_InstructionNewReQuantifier struct {
	LhsOp         *Bir_Operand
	Quantifier    *Bir_Operand
	NonGreedyChar *Bir_Operand
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_Instruction
}

func NewBir_InstructionNewReQuantifier() *Bir_InstructionNewReQuantifier {
	return &Bir_InstructionNewReQuantifier{}
}

func (this Bir_InstructionNewReQuantifier) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReQuantifier) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp372 := NewBir_Operand()
	err = tmp372.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp372
	tmp373 := NewBir_Operand()
	err = tmp373.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Quantifier = tmp373
	tmp374 := NewBir_Operand()
	err = tmp374.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.NonGreedyChar = tmp374
	return err
}

type Bir_InstructionNewReSequence struct {
	Terms   *Bir_Operand
	LhsOp   *Bir_Operand
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_Instruction
}

func NewBir_InstructionNewReSequence() *Bir_InstructionNewReSequence {
	return &Bir_InstructionNewReSequence{}
}

func (this Bir_InstructionNewReSequence) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewReSequence) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp375 := NewBir_Operand()
	err = tmp375.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Terms = tmp375
	tmp376 := NewBir_Operand()
	err = tmp376.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp376
	return err
}

type Bir_InstructionNewRegExp struct {
	LhsOp         *Bir_Operand
	ReDisjunction *Bir_Operand
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_Instruction
}

func NewBir_InstructionNewRegExp() *Bir_InstructionNewRegExp {
	return &Bir_InstructionNewRegExp{}
}

func (this Bir_InstructionNewRegExp) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewRegExp) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp377 := NewBir_Operand()
	err = tmp377.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp377
	tmp378 := NewBir_Operand()
	err = tmp378.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ReDisjunction = tmp378
	return err
}

type Bir_InstructionNewStringXmlQname struct {
	LhsOperand         *Bir_Operand
	StringQnameOperand *Bir_Operand
	_io                *kaitai.Stream
	_root              *Bir
	_parent            *Bir_Instruction
}

func NewBir_InstructionNewStringXmlQname() *Bir_InstructionNewStringXmlQname {
	return &Bir_InstructionNewStringXmlQname{}
}

func (this Bir_InstructionNewStringXmlQname) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewStringXmlQname) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp379 := NewBir_Operand()
	err = tmp379.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp379
	tmp380 := NewBir_Operand()
	err = tmp380.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.StringQnameOperand = tmp380
	return err
}

type Bir_InstructionNewStructure struct {
	RhsOperand      *Bir_Operand
	LhsOperand      *Bir_Operand
	InitValuesCount int32
	InitValues      []*Bir_MappingConstructor
	_io             *kaitai.Stream
	_root           *Bir
	_parent         *Bir_Instruction
}

func NewBir_InstructionNewStructure() *Bir_InstructionNewStructure {
	return &Bir_InstructionNewStructure{}
}

func (this Bir_InstructionNewStructure) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewStructure) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp381 := NewBir_Operand()
	err = tmp381.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp381
	tmp382 := NewBir_Operand()
	err = tmp382.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp382
	tmp383, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InitValuesCount = int32(tmp383)
	for i := 0; i < int(this.InitValuesCount); i++ {
		_ = i
		tmp384 := NewBir_MappingConstructor()
		err = tmp384.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InitValues = append(this.InitValues, tmp384)
	}
	return err
}

type Bir_InstructionNewTable struct {
	TypeCpIndex      int32
	LhsOperand       *Bir_Operand
	KeyColumnOperand *Bir_Operand
	DataOperand      *Bir_Operand
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Instruction
}

func NewBir_InstructionNewTable() *Bir_InstructionNewTable {
	return &Bir_InstructionNewTable{}
}

func (this Bir_InstructionNewTable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewTable) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp385, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp385)
	tmp386 := NewBir_Operand()
	err = tmp386.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp386
	tmp387 := NewBir_Operand()
	err = tmp387.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.KeyColumnOperand = tmp387
	tmp388 := NewBir_Operand()
	err = tmp388.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.DataOperand = tmp388
	return err
}

type Bir_InstructionNewTypedesc struct {
	LhsOperand  *Bir_Operand
	TypeCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionNewTypedesc() *Bir_InstructionNewTypedesc {
	return &Bir_InstructionNewTypedesc{}
}

func (this Bir_InstructionNewTypedesc) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewTypedesc) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp389 := NewBir_Operand()
	err = tmp389.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp389
	tmp390, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp390)
	return err
}

type Bir_InstructionNewXmlComment struct {
	LhsOperand  *Bir_Operand
	TextOperand *Bir_Operand
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionNewXmlComment() *Bir_InstructionNewXmlComment {
	return &Bir_InstructionNewXmlComment{}
}

func (this Bir_InstructionNewXmlComment) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlComment) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp391 := NewBir_Operand()
	err = tmp391.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp391
	tmp392 := NewBir_Operand()
	err = tmp392.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TextOperand = tmp392
	return err
}

type Bir_InstructionNewXmlElement struct {
	LhsOperand          *Bir_Operand
	StartTagOperand     *Bir_Operand
	DefaultNsUriOperand *Bir_Operand
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionNewXmlElement() *Bir_InstructionNewXmlElement {
	return &Bir_InstructionNewXmlElement{}
}

func (this Bir_InstructionNewXmlElement) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlElement) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp393 := NewBir_Operand()
	err = tmp393.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp393
	tmp394 := NewBir_Operand()
	err = tmp394.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.StartTagOperand = tmp394
	tmp395 := NewBir_Operand()
	err = tmp395.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.DefaultNsUriOperand = tmp395
	return err
}

type Bir_InstructionNewXmlProcessIns struct {
	LhsOperand    *Bir_Operand
	DataOperand   *Bir_Operand
	TargetOperand *Bir_Operand
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_Instruction
}

func NewBir_InstructionNewXmlProcessIns() *Bir_InstructionNewXmlProcessIns {
	return &Bir_InstructionNewXmlProcessIns{}
}

func (this Bir_InstructionNewXmlProcessIns) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlProcessIns) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp396 := NewBir_Operand()
	err = tmp396.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp396
	tmp397 := NewBir_Operand()
	err = tmp397.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.DataOperand = tmp397
	tmp398 := NewBir_Operand()
	err = tmp398.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TargetOperand = tmp398
	return err
}

type Bir_InstructionNewXmlQname struct {
	LhsOperand       *Bir_Operand
	LocalNameOperand *Bir_Operand
	NsUriOperand     *Bir_Operand
	PrefixOperand    *Bir_Operand
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Instruction
}

func NewBir_InstructionNewXmlQname() *Bir_InstructionNewXmlQname {
	return &Bir_InstructionNewXmlQname{}
}

func (this Bir_InstructionNewXmlQname) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlQname) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp399 := NewBir_Operand()
	err = tmp399.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp399
	tmp400 := NewBir_Operand()
	err = tmp400.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LocalNameOperand = tmp400
	tmp401 := NewBir_Operand()
	err = tmp401.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.NsUriOperand = tmp401
	tmp402 := NewBir_Operand()
	err = tmp402.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.PrefixOperand = tmp402
	return err
}

type Bir_InstructionNewXmlSequence struct {
	LhsOperand *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionNewXmlSequence() *Bir_InstructionNewXmlSequence {
	return &Bir_InstructionNewXmlSequence{}
}

func (this Bir_InstructionNewXmlSequence) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlSequence) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp403 := NewBir_Operand()
	err = tmp403.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp403
	return err
}

type Bir_InstructionNewXmlText struct {
	LhsOperand  *Bir_Operand
	TextOperand *Bir_Operand
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionNewXmlText() *Bir_InstructionNewXmlText {
	return &Bir_InstructionNewXmlText{}
}

func (this Bir_InstructionNewXmlText) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionNewXmlText) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp404 := NewBir_Operand()
	err = tmp404.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp404
	tmp405 := NewBir_Operand()
	err = tmp405.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TextOperand = tmp405
	return err
}

type Bir_InstructionObjectLoad struct {
	ObjectLoad *Bir_IndexAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionObjectLoad() *Bir_InstructionObjectLoad {
	return &Bir_InstructionObjectLoad{}
}

func (this Bir_InstructionObjectLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionObjectLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp406 := NewBir_IndexAccess()
	err = tmp406.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ObjectLoad = tmp406
	return err
}

type Bir_InstructionObjectStore struct {
	ObjectStore *Bir_IndexAccess
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionObjectStore() *Bir_InstructionObjectStore {
	return &Bir_InstructionObjectStore{}
}

func (this Bir_InstructionObjectStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionObjectStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp407 := NewBir_IndexAccess()
	err = tmp407.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ObjectStore = tmp407
	return err
}

type Bir_InstructionPanic struct {
	ErrorOperand *Bir_Operand
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Instruction
}

func NewBir_InstructionPanic() *Bir_InstructionPanic {
	return &Bir_InstructionPanic{}
}

func (this Bir_InstructionPanic) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionPanic) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp408 := NewBir_Operand()
	err = tmp408.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ErrorOperand = tmp408
	return err
}

type Bir_InstructionRecordDefaultFpLoad struct {
	LhsOp             *Bir_Operand
	EnclosedTypeIndex int32
	FieldName         int32
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_Instruction
}

func NewBir_InstructionRecordDefaultFpLoad() *Bir_InstructionRecordDefaultFpLoad {
	return &Bir_InstructionRecordDefaultFpLoad{}
}

func (this Bir_InstructionRecordDefaultFpLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionRecordDefaultFpLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp409 := NewBir_Operand()
	err = tmp409.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOp = tmp409
	tmp410, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.EnclosedTypeIndex = int32(tmp410)
	tmp411, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FieldName = int32(tmp411)
	return err
}

type Bir_InstructionReturn struct {
	NoValue []byte
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_Instruction
}

func NewBir_InstructionReturn() *Bir_InstructionReturn {
	return &Bir_InstructionReturn{}
}

func (this Bir_InstructionReturn) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionReturn) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp412, err := this._io.ReadBytes(int(0))
	if err != nil {
		return err
	}
	this.NoValue = tmp412
	return err
}

type Bir_InstructionStringLoad struct {
	StringLoad *Bir_IndexAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionStringLoad() *Bir_InstructionStringLoad {
	return &Bir_InstructionStringLoad{}
}

func (this Bir_InstructionStringLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionStringLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp413 := NewBir_IndexAccess()
	err = tmp413.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.StringLoad = tmp413
	return err
}

type Bir_InstructionTableLoad struct {
	TableLoad *Bir_IndexAccess
	_io       *kaitai.Stream
	_root     *Bir
	_parent   *Bir_Instruction
}

func NewBir_InstructionTableLoad() *Bir_InstructionTableLoad {
	return &Bir_InstructionTableLoad{}
}

func (this Bir_InstructionTableLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionTableLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp414 := NewBir_IndexAccess()
	err = tmp414.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TableLoad = tmp414
	return err
}

type Bir_InstructionTableStore struct {
	TableStore *Bir_IndexAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionTableStore() *Bir_InstructionTableStore {
	return &Bir_InstructionTableStore{}
}

func (this Bir_InstructionTableStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionTableStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp415 := NewBir_IndexAccess()
	err = tmp415.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TableStore = tmp415
	return err
}

type Bir_InstructionTypeCast struct {
	LhsOperand   *Bir_Operand
	RhsOperand   *Bir_Operand
	TypeCpIndex  int32
	IsCheckTypes uint8
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_Instruction
}

func NewBir_InstructionTypeCast() *Bir_InstructionTypeCast {
	return &Bir_InstructionTypeCast{}
}

func (this Bir_InstructionTypeCast) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionTypeCast) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp416 := NewBir_Operand()
	err = tmp416.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp416
	tmp417 := NewBir_Operand()
	err = tmp417.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp417
	tmp418, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp418)
	tmp419, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsCheckTypes = tmp419
	return err
}

type Bir_InstructionTypeTest struct {
	TypeCpIndex int32
	LhsOperand  *Bir_Operand
	RhsOperand  *Bir_Operand
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionTypeTest() *Bir_InstructionTypeTest {
	return &Bir_InstructionTypeTest{}
}

func (this Bir_InstructionTypeTest) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionTypeTest) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp420, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp420)
	tmp421 := NewBir_Operand()
	err = tmp421.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp421
	tmp422 := NewBir_Operand()
	err = tmp422.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp422
	return err
}

type Bir_InstructionUnaryOperation struct {
	RhsOperand *Bir_Operand
	LhsOperand *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionUnaryOperation() *Bir_InstructionUnaryOperation {
	return &Bir_InstructionUnaryOperation{}
}

func (this Bir_InstructionUnaryOperation) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionUnaryOperation) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp423 := NewBir_Operand()
	err = tmp423.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp423
	tmp424 := NewBir_Operand()
	err = tmp424.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp424
	return err
}

type Bir_InstructionUnlock struct {
	UnlockBbIdNameCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_Instruction
}

func NewBir_InstructionUnlock() *Bir_InstructionUnlock {
	return &Bir_InstructionUnlock{}
}

func (this Bir_InstructionUnlock) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionUnlock) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp425, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.UnlockBbIdNameCpIndex = int32(tmp425)
	return err
}

type Bir_InstructionWait struct {
	WaitExpressionsCount int32
	WaitExpressions      []*Bir_Operand
	LhsOperand           *Bir_Operand
	ThenBbIdNameCpIndex  int32
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_Instruction
}

func NewBir_InstructionWait() *Bir_InstructionWait {
	return &Bir_InstructionWait{}
}

func (this Bir_InstructionWait) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWait) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp426, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.WaitExpressionsCount = int32(tmp426)
	for i := 0; i < int(this.WaitExpressionsCount); i++ {
		_ = i
		tmp427 := NewBir_Operand()
		err = tmp427.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.WaitExpressions = append(this.WaitExpressions, tmp427)
	}
	tmp428 := NewBir_Operand()
	err = tmp428.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp428
	tmp429, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp429)
	return err
}

type Bir_InstructionWaitAll struct {
	LhsOperand            *Bir_Operand
	KeyNamesCount         int32
	KeyNameCpIndex        []int32
	ValueExpressionsCount int32
	ValueExpression       []*Bir_Operand
	ThenBbIdNameCpIndex   int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_Instruction
}

func NewBir_InstructionWaitAll() *Bir_InstructionWaitAll {
	return &Bir_InstructionWaitAll{}
}

func (this Bir_InstructionWaitAll) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWaitAll) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp430 := NewBir_Operand()
	err = tmp430.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp430
	tmp431, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.KeyNamesCount = int32(tmp431)
	for i := 0; i < int(this.KeyNamesCount); i++ {
		_ = i
		tmp432, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.KeyNameCpIndex = append(this.KeyNameCpIndex, tmp432)
	}
	tmp433, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueExpressionsCount = int32(tmp433)
	for i := 0; i < int(this.ValueExpressionsCount); i++ {
		_ = i
		tmp434 := NewBir_Operand()
		err = tmp434.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ValueExpression = append(this.ValueExpression, tmp434)
	}
	tmp435, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp435)
	return err
}

type Bir_InstructionWkAltReceive struct {
	ChannelNameCount    int32
	ChannelNameCpIndex  []int32
	LhsOperand          *Bir_Operand
	IsSameStrand        uint8
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionWkAltReceive() *Bir_InstructionWkAltReceive {
	return &Bir_InstructionWkAltReceive{}
}

func (this Bir_InstructionWkAltReceive) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWkAltReceive) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp436, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ChannelNameCount = int32(tmp436)
	for i := 0; i < int(this.ChannelNameCount); i++ {
		_ = i
		tmp437, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.ChannelNameCpIndex = append(this.ChannelNameCpIndex, tmp437)
	}
	tmp438 := NewBir_Operand()
	err = tmp438.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp438
	tmp439, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSameStrand = tmp439
	tmp440, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp440)
	return err
}

type Bir_InstructionWkMulReceive struct {
	ChannelFieldCount   int32
	ChannelFieldCpIndex []*Bir_ReceiveField
	TypeCpIndex         int32
	LhsOperand          *Bir_Operand
	IsSameStrand        uint8
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionWkMulReceive() *Bir_InstructionWkMulReceive {
	return &Bir_InstructionWkMulReceive{}
}

func (this Bir_InstructionWkMulReceive) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWkMulReceive) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp441, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ChannelFieldCount = int32(tmp441)
	for i := 0; i < int(this.ChannelFieldCount); i++ {
		_ = i
		tmp442 := NewBir_ReceiveField()
		err = tmp442.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ChannelFieldCpIndex = append(this.ChannelFieldCpIndex, tmp442)
	}
	tmp443, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp443)
	tmp444 := NewBir_Operand()
	err = tmp444.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp444
	tmp445, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSameStrand = tmp445
	tmp446, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp446)
	return err
}

type Bir_InstructionWkReceive struct {
	WorkerNameCpIndex   int32
	LhsOperand          *Bir_Operand
	IsSameStrand        uint8
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionWkReceive() *Bir_InstructionWkReceive {
	return &Bir_InstructionWkReceive{}
}

func (this Bir_InstructionWkReceive) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWkReceive) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp447, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.WorkerNameCpIndex = int32(tmp447)
	tmp448 := NewBir_Operand()
	err = tmp448.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp448
	tmp449, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSameStrand = tmp449
	tmp450, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp450)
	return err
}

type Bir_InstructionWkSend struct {
	ChannelNameCpIndex  int32
	WorkerDataOperand   *Bir_Operand
	IsSameStrand        uint8
	IsSynch             uint8
	LhsOperand          *Bir_Operand
	ThenBbIdNameCpIndex int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_Instruction
}

func NewBir_InstructionWkSend() *Bir_InstructionWkSend {
	return &Bir_InstructionWkSend{}
}

func (this Bir_InstructionWkSend) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionWkSend) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp451, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ChannelNameCpIndex = int32(tmp451)
	tmp452 := NewBir_Operand()
	err = tmp452.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.WorkerDataOperand = tmp452
	tmp453, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSameStrand = tmp453
	tmp454, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSynch = tmp454
	if this.IsSynch == 1 {
		tmp455 := NewBir_Operand()
		err = tmp455.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.LhsOperand = tmp455
	}
	tmp456, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ThenBbIdNameCpIndex = int32(tmp456)
	return err
}

type Bir_InstructionXmlAttributeLoad struct {
	XmlAttributeLoad *Bir_IndexAccess
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Instruction
}

func NewBir_InstructionXmlAttributeLoad() *Bir_InstructionXmlAttributeLoad {
	return &Bir_InstructionXmlAttributeLoad{}
}

func (this Bir_InstructionXmlAttributeLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlAttributeLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp457 := NewBir_IndexAccess()
	err = tmp457.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlAttributeLoad = tmp457
	return err
}

type Bir_InstructionXmlAttributeStore struct {
	XmlAttributeStore *Bir_IndexAccess
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_Instruction
}

func NewBir_InstructionXmlAttributeStore() *Bir_InstructionXmlAttributeStore {
	return &Bir_InstructionXmlAttributeStore{}
}

func (this Bir_InstructionXmlAttributeStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlAttributeStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp458 := NewBir_IndexAccess()
	err = tmp458.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlAttributeStore = tmp458
	return err
}

type Bir_InstructionXmlLoad struct {
	XmlLoad *Bir_IndexAccess
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_Instruction
}

func NewBir_InstructionXmlLoad() *Bir_InstructionXmlLoad {
	return &Bir_InstructionXmlLoad{}
}

func (this Bir_InstructionXmlLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp459 := NewBir_IndexAccess()
	err = tmp459.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlLoad = tmp459
	return err
}

type Bir_InstructionXmlLoadAll struct {
	XmlLoadAll *Bir_XmlAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionXmlLoadAll() *Bir_InstructionXmlLoadAll {
	return &Bir_InstructionXmlLoadAll{}
}

func (this Bir_InstructionXmlLoadAll) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlLoadAll) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp460 := NewBir_XmlAccess()
	err = tmp460.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlLoadAll = tmp460
	return err
}

type Bir_InstructionXmlSeqLoad struct {
	XmlSeqLoad *Bir_IndexAccess
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_Instruction
}

func NewBir_InstructionXmlSeqLoad() *Bir_InstructionXmlSeqLoad {
	return &Bir_InstructionXmlSeqLoad{}
}

func (this Bir_InstructionXmlSeqLoad) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlSeqLoad) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp461 := NewBir_IndexAccess()
	err = tmp461.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlSeqLoad = tmp461
	return err
}

type Bir_InstructionXmlSeqStore struct {
	XmlSeqStore *Bir_XmlAccess
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Instruction
}

func NewBir_InstructionXmlSeqStore() *Bir_InstructionXmlSeqStore {
	return &Bir_InstructionXmlSeqStore{}
}

func (this Bir_InstructionXmlSeqStore) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InstructionXmlSeqStore) Read(io *kaitai.Stream, parent *Bir_Instruction, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp462 := NewBir_XmlAccess()
	err = tmp462.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.XmlSeqStore = tmp462
	return err
}

type Bir_IntConstantInfo struct {
	ValueCpIndex int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir_IntConstantInfo() *Bir_IntConstantInfo {
	return &Bir_IntConstantInfo{}
}

func (this Bir_IntConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_IntConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp463, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueCpIndex = int32(tmp463)
	return err
}

type Bir_IntCpInfo struct {
	Value   int64
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolEntry
}

func NewBir_IntCpInfo() *Bir_IntCpInfo {
	return &Bir_IntCpInfo{}
}

func (this Bir_IntCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_IntCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp464, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Value = int64(tmp464)
	return err
}

type Bir_IntersectionConstantInfo struct {
	ConstantValueInfo   kaitai.Struct
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_ConstantValue
	_f_effectiveType    bool
	effectiveType       *Bir_ShapeCpInfo
	_f_intersectionType bool
	intersectionType    *Bir_TypeIntersection
	_f_type             bool
	t                   *Bir_ShapeCpInfo
}

func NewBir_IntersectionConstantInfo() *Bir_IntersectionConstantInfo {
	return &Bir_IntersectionConstantInfo{}
}

func (this Bir_IntersectionConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_IntersectionConstantInfo) Read(io *kaitai.Stream, parent *Bir_ConstantValue, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp465, err := this.EffectiveType()
	if err != nil {
		return err
	}
	switch tmp465.Shape.TypeTag {
	case Bir_TypeTagEnum__TypeTagRecord:
		tmp466 := NewBir_MapConstantInfo()
		err = tmp466.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp466
	case Bir_TypeTagEnum__TypeTagTuple:
		tmp467 := NewBir_ListConstantInfo()
		err = tmp467.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ConstantValueInfo = tmp467
	}
	return err
}

func (this *Bir_IntersectionConstantInfo) EffectiveType() (v *Bir_ShapeCpInfo, err error) {
	if this._f_effectiveType {
		return this.effectiveType, nil
	}
	this._f_effectiveType = true
	tmp468, err := this.IntersectionType()
	if err != nil {
		return nil, err
	}
	this.effectiveType = this._root.ConstantPool.ConstantPoolEntries[tmp468.EffectiveTypeCpIndex].CpInfo.(*Bir_ShapeCpInfo)
	return this.effectiveType, nil
}

func (this *Bir_IntersectionConstantInfo) IntersectionType() (v *Bir_TypeIntersection, err error) {
	if this._f_intersectionType {
		return this.intersectionType, nil
	}
	this._f_intersectionType = true
	tmp469, err := this.Type()
	if err != nil {
		return nil, err
	}
	this.intersectionType = tmp469.Shape.TypeStructure.(*Bir_TypeIntersection)
	return this.intersectionType, nil
}

func (this *Bir_IntersectionConstantInfo) Type() (v *Bir_ShapeCpInfo, err error) {
	if this._f_type {
		return this.t, nil
	}
	this._f_type = true
	this.t = this._root.ConstantPool.ConstantPoolEntries[this._parent.ConstantValueTypeCpIndex].CpInfo.(*Bir_ShapeCpInfo)
	return this.t, nil
}

type Bir_InvokableTypeSymbolBody struct {
	ParamCount    int32
	Params        []*Bir_FunctionParameter
	HasRestType   uint8
	RestParam     *Bir_FunctionParameter
	DefaultValues int32
	DefaultValue  []*Bir_DefaultValueBody
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_TypeInvokableBody
}

func NewBir_InvokableTypeSymbolBody() *Bir_InvokableTypeSymbolBody {
	return &Bir_InvokableTypeSymbolBody{}
}

func (this Bir_InvokableTypeSymbolBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_InvokableTypeSymbolBody) Read(io *kaitai.Stream, parent *Bir_TypeInvokableBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp470, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamCount = int32(tmp470)
	for i := 0; i < int(this.ParamCount); i++ {
		_ = i
		tmp471 := NewBir_FunctionParameter()
		err = tmp471.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = append(this.Params, tmp471)
	}
	tmp472, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasRestType = tmp472
	if this.HasRestType == 1 {
		tmp473 := NewBir_FunctionParameter()
		err = tmp473.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RestParam = tmp473
	}
	tmp474, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DefaultValues = int32(tmp474)
	for i := 0; i < int(this.DefaultValues); i++ {
		_ = i
		tmp475 := NewBir_DefaultValueBody()
		err = tmp475.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DefaultValue = append(this.DefaultValue, tmp475)
	}
	return err
}

type Bir_ListConstantInfo struct {
	ListConstantSize    int32
	ListMemberValueInfo []*Bir_ConstantValue
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_IntersectionConstantInfo
}

func NewBir_ListConstantInfo() *Bir_ListConstantInfo {
	return &Bir_ListConstantInfo{}
}

func (this Bir_ListConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ListConstantInfo) Read(io *kaitai.Stream, parent *Bir_IntersectionConstantInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp476, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ListConstantSize = int32(tmp476)
	for i := 0; i < int(this.ListConstantSize); i++ {
		_ = i
		tmp477 := NewBir_ConstantValue()
		err = tmp477.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ListMemberValueInfo = append(this.ListMemberValueInfo, tmp477)
	}
	return err
}

type Bir_ListenerType struct {
	TypeCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_ServiceDeclaration
}

func NewBir_ListenerType() *Bir_ListenerType {
	return &Bir_ListenerType{}
}

func (this Bir_ListenerType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ListenerType) Read(io *kaitai.Stream, parent *Bir_ServiceDeclaration, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp478, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp478)
	return err
}

type Bir_LocalVariable struct {
	Kind                  int8
	TypeCpIndex           int32
	NameCpIndex           int32
	MetaVarNameCpIndex    int32
	EnclosingBasicBlockId *Bir_EnclosingBasicBlockId
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_FunctionBody
}

func NewBir_LocalVariable() *Bir_LocalVariable {
	return &Bir_LocalVariable{}
}

func (this Bir_LocalVariable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_LocalVariable) Read(io *kaitai.Stream, parent *Bir_FunctionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp479, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp479
	tmp480, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp480)
	tmp481, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp481)
	if this.Kind == 2 {
		tmp482, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.MetaVarNameCpIndex = int32(tmp482)
	}
	if this.Kind == 1 {
		tmp483 := NewBir_EnclosingBasicBlockId()
		err = tmp483.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.EnclosingBasicBlockId = tmp483
	}
	return err
}

type Bir_MapConstantInfo struct {
	MapConstantSize int32
	MapKeyValues    []*Bir_MapKeyValue
	_io             *kaitai.Stream
	_root           *Bir
	_parent         kaitai.Struct
}

func NewBir_MapConstantInfo() *Bir_MapConstantInfo {
	return &Bir_MapConstantInfo{}
}

func (this Bir_MapConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MapConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp484, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.MapConstantSize = int32(tmp484)
	for i := 0; i < int(this.MapConstantSize); i++ {
		_ = i
		tmp485 := NewBir_MapKeyValue()
		err = tmp485.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.MapKeyValues = append(this.MapKeyValues, tmp485)
	}
	return err
}

type Bir_MapKeyValue struct {
	KeyNameCpIndex int32
	KeyValueInfo   *Bir_ConstantValue
	_io            *kaitai.Stream
	_root          *Bir
	_parent        *Bir_MapConstantInfo
}

func NewBir_MapKeyValue() *Bir_MapKeyValue {
	return &Bir_MapKeyValue{}
}

func (this Bir_MapKeyValue) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MapKeyValue) Read(io *kaitai.Stream, parent *Bir_MapConstantInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp486, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.KeyNameCpIndex = int32(tmp486)
	tmp487 := NewBir_ConstantValue()
	err = tmp487.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.KeyValueInfo = tmp487
	return err
}

type Bir_MappingConstructor_MappingConstructorBodyKind int

const (
	Bir_MappingConstructor_MappingConstructorBodyKind__MappingConstructorSpreadFieldKind Bir_MappingConstructor_MappingConstructorBodyKind = 0
	Bir_MappingConstructor_MappingConstructorBodyKind__MappingConstructorKeyValueKind    Bir_MappingConstructor_MappingConstructorBodyKind = 1
)

type Bir_MappingConstructor struct {
	MappingConstructorKind Bir_MappingConstructor_MappingConstructorBodyKind
	MappingConstructorBody kaitai.Struct
	_io                    *kaitai.Stream
	_root                  *Bir
	_parent                *Bir_InstructionNewStructure
}

func NewBir_MappingConstructor() *Bir_MappingConstructor {
	return &Bir_MappingConstructor{}
}

func (this Bir_MappingConstructor) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MappingConstructor) Read(io *kaitai.Stream, parent *Bir_InstructionNewStructure, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp488, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.MappingConstructorKind = Bir_MappingConstructor_MappingConstructorBodyKind(tmp488)
	switch this.MappingConstructorKind {
	case Bir_MappingConstructor_MappingConstructorBodyKind__MappingConstructorKeyValueKind:
		tmp489 := NewBir_MappingConstructorKeyValueBody()
		err = tmp489.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.MappingConstructorBody = tmp489
	case Bir_MappingConstructor_MappingConstructorBodyKind__MappingConstructorSpreadFieldKind:
		tmp490 := NewBir_MappingConstructorSpreadFieldBody()
		err = tmp490.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.MappingConstructorBody = tmp490
	}
	return err
}

type Bir_MappingConstructorKeyValueBody struct {
	KeyOperand   *Bir_Operand
	ValueOperand *Bir_Operand
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_MappingConstructor
}

func NewBir_MappingConstructorKeyValueBody() *Bir_MappingConstructorKeyValueBody {
	return &Bir_MappingConstructorKeyValueBody{}
}

func (this Bir_MappingConstructorKeyValueBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MappingConstructorKeyValueBody) Read(io *kaitai.Stream, parent *Bir_MappingConstructor, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp491 := NewBir_Operand()
	err = tmp491.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.KeyOperand = tmp491
	tmp492 := NewBir_Operand()
	err = tmp492.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ValueOperand = tmp492
	return err
}

type Bir_MappingConstructorSpreadFieldBody struct {
	ExprOperand *Bir_Operand
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_MappingConstructor
}

func NewBir_MappingConstructorSpreadFieldBody() *Bir_MappingConstructorSpreadFieldBody {
	return &Bir_MappingConstructorSpreadFieldBody{}
}

func (this Bir_MappingConstructorSpreadFieldBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MappingConstructorSpreadFieldBody) Read(io *kaitai.Stream, parent *Bir_MappingConstructor, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp493 := NewBir_Operand()
	err = tmp493.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ExprOperand = tmp493
	return err
}

type Bir_Markdown struct {
	Length          int32
	HasDoc          uint8
	MarkdownContent *Bir_MarkdownContent
	_io             *kaitai.Stream
	_root           *Bir
	_parent         kaitai.Struct
}

func NewBir_Markdown() *Bir_Markdown {
	return &Bir_Markdown{}
}

func (this Bir_Markdown) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Markdown) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp494, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Length = int32(tmp494)
	tmp495, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasDoc = tmp495
	if this.HasDoc == 1 {
		tmp496 := NewBir_MarkdownContent()
		err = tmp496.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.MarkdownContent = tmp496
	}
	return err
}

type Bir_MarkdownContent struct {
	DescriptionCpIndex            int32
	ReturnValueDescriptionCpIndex int32
	ParametersCount               int32
	Parameters                    []*Bir_MarkdownParameter
	DeprecatedDocsCpIndex         int32
	DeprecatedParamsCount         int32
	DeprecatedParams              []*Bir_MarkdownParameter
	_io                           *kaitai.Stream
	_root                         *Bir
	_parent                       *Bir_Markdown
}

func NewBir_MarkdownContent() *Bir_MarkdownContent {
	return &Bir_MarkdownContent{}
}

func (this Bir_MarkdownContent) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MarkdownContent) Read(io *kaitai.Stream, parent *Bir_Markdown, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp497, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DescriptionCpIndex = int32(tmp497)
	tmp498, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ReturnValueDescriptionCpIndex = int32(tmp498)
	tmp499, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParametersCount = int32(tmp499)
	for i := 0; i < int(this.ParametersCount); i++ {
		_ = i
		tmp500 := NewBir_MarkdownParameter()
		err = tmp500.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Parameters = append(this.Parameters, tmp500)
	}
	tmp501, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DeprecatedDocsCpIndex = int32(tmp501)
	tmp502, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DeprecatedParamsCount = int32(tmp502)
	for i := 0; i < int(this.DeprecatedParamsCount); i++ {
		_ = i
		tmp503 := NewBir_MarkdownParameter()
		err = tmp503.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DeprecatedParams = append(this.DeprecatedParams, tmp503)
	}
	return err
}

type Bir_MarkdownParameter struct {
	NameCpIndex        int32
	DescriptionCpIndex int32
	_io                *kaitai.Stream
	_root              *Bir
	_parent            *Bir_MarkdownContent
}

func NewBir_MarkdownParameter() *Bir_MarkdownParameter {
	return &Bir_MarkdownParameter{}
}

func (this Bir_MarkdownParameter) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_MarkdownParameter) Read(io *kaitai.Stream, parent *Bir_MarkdownContent, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp504, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp504)
	tmp505, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DescriptionCpIndex = int32(tmp505)
	return err
}

type Bir_Module struct {
	IdCpIndex                 int32
	ImportCount               int32
	Imports                   []*Bir_PackageCpInfo
	ConstCount                int32
	Constants                 []*Bir_Constant
	TypeDefinitionCount       int32
	TypeDefinitions           []*Bir_TypeDefinition
	GlobalVarCount            int32
	GlobalVars                []*Bir_GlobalVar
	TypeDefinitionBodiesCount int32
	TypeDefinitionBodies      []*Bir_TypeDefinitionBody
	FunctionCount             int32
	Functions                 []*Bir_Function
	AnnotationsSize           int32
	Annotations               []*Bir_Annotation
	ServiceDeclsSize          int32
	ServiceDeclarations       []*Bir_ServiceDeclaration
	_io                       *kaitai.Stream
	_root                     *Bir
	_parent                   *Bir
}

func NewBir_Module() *Bir_Module {
	return &Bir_Module{}
}

func (this Bir_Module) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Module) Read(io *kaitai.Stream, parent *Bir, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp506, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.IdCpIndex = int32(tmp506)
	tmp507, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ImportCount = int32(tmp507)
	for i := 0; i < int(this.ImportCount); i++ {
		_ = i
		tmp508 := NewBir_PackageCpInfo()
		err = tmp508.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Imports = append(this.Imports, tmp508)
	}
	tmp509, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstCount = int32(tmp509)
	for i := 0; i < int(this.ConstCount); i++ {
		_ = i
		tmp510 := NewBir_Constant()
		err = tmp510.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Constants = append(this.Constants, tmp510)
	}
	tmp511, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeDefinitionCount = int32(tmp511)
	for i := 0; i < int(this.TypeDefinitionCount); i++ {
		_ = i
		tmp512 := NewBir_TypeDefinition()
		err = tmp512.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeDefinitions = append(this.TypeDefinitions, tmp512)
	}
	tmp513, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.GlobalVarCount = int32(tmp513)
	for i := 0; i < int(this.GlobalVarCount); i++ {
		_ = i
		tmp514 := NewBir_GlobalVar()
		err = tmp514.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.GlobalVars = append(this.GlobalVars, tmp514)
	}
	tmp515, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeDefinitionBodiesCount = int32(tmp515)
	for i := 0; i < int(this.TypeDefinitionBodiesCount); i++ {
		_ = i
		tmp516 := NewBir_TypeDefinitionBody()
		err = tmp516.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeDefinitionBodies = append(this.TypeDefinitionBodies, tmp516)
	}
	tmp517, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FunctionCount = int32(tmp517)
	for i := 0; i < int(this.FunctionCount); i++ {
		_ = i
		tmp518 := NewBir_Function()
		err = tmp518.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Functions = append(this.Functions, tmp518)
	}
	tmp519, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AnnotationsSize = int32(tmp519)
	for i := 0; i < int(this.AnnotationsSize); i++ {
		_ = i
		tmp520 := NewBir_Annotation()
		err = tmp520.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Annotations = append(this.Annotations, tmp520)
	}
	tmp521, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ServiceDeclsSize = int32(tmp521)
	for i := 0; i < int(this.ServiceDeclsSize); i++ {
		_ = i
		tmp522 := NewBir_ServiceDeclaration()
		err = tmp522.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ServiceDeclarations = append(this.ServiceDeclarations, tmp522)
	}
	return err
}

type Bir_NilConstantInfo struct {
	ValueNilConstant []byte
	_io              *kaitai.Stream
	_root            *Bir
	_parent          kaitai.Struct
}

func NewBir_NilConstantInfo() *Bir_NilConstantInfo {
	return &Bir_NilConstantInfo{}
}

func (this Bir_NilConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_NilConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp523, err := this._io.ReadBytes(int(0))
	if err != nil {
		return err
	}
	this.ValueNilConstant = tmp523
	return err
}

type Bir_NullableStrInfo struct {
	HasNonNullString uint8
	StrCpIndex       int32
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_SemNamedType
}

func NewBir_NullableStrInfo() *Bir_NullableStrInfo {
	return &Bir_NullableStrInfo{}
}

func (this Bir_NullableStrInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_NullableStrInfo) Read(io *kaitai.Stream, parent *Bir_SemNamedType, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp524, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasNonNullString = tmp524
	if this.HasNonNullString == 1 {
		tmp525, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.StrCpIndex = int32(tmp525)
	}
	return err
}

type Bir_ObjectAttachedFunction struct {
	NameCpIndex         int32
	OriginalNameCpIndex int32
	Flags               int64
	TypeCpIndex         int32
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_TypeObjectOrService
}

func NewBir_ObjectAttachedFunction() *Bir_ObjectAttachedFunction {
	return &Bir_ObjectAttachedFunction{}
}

func (this Bir_ObjectAttachedFunction) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ObjectAttachedFunction) Read(io *kaitai.Stream, parent *Bir_TypeObjectOrService, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp526, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp526)
	tmp527, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OriginalNameCpIndex = int32(tmp527)
	tmp528, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp528)
	tmp529, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp529)
	return err
}

type Bir_ObjectField struct {
	NameCpIndex   int32
	Flags         int64
	IsDefaultable uint8
	Doc           *Bir_Markdown
	TypeCpIndex   int32
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_TypeObjectOrService
}

func NewBir_ObjectField() *Bir_ObjectField {
	return &Bir_ObjectField{}
}

func (this Bir_ObjectField) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ObjectField) Read(io *kaitai.Stream, parent *Bir_TypeObjectOrService, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp530, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp530)
	tmp531, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp531)
	tmp532, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsDefaultable = tmp532
	tmp533 := NewBir_Markdown()
	err = tmp533.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp533
	tmp534, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp534)
	return err
}

type Bir_Operand struct {
	IgnoredVariable    uint8
	IgnoredTypeCpIndex int32
	Variable           *Bir_Variable
	_io                *kaitai.Stream
	_root              *Bir
	_parent            kaitai.Struct
}

func NewBir_Operand() *Bir_Operand {
	return &Bir_Operand{}
}

func (this Bir_Operand) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Operand) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp535, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IgnoredVariable = tmp535
	if this.IgnoredVariable == 1 {
		tmp536, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.IgnoredTypeCpIndex = int32(tmp536)
	}
	if this.IgnoredVariable == 0 {
		tmp537 := NewBir_Variable()
		err = tmp537.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Variable = tmp537
	}
	return err
}

type Bir_PackageCpInfo struct {
	OrgIndex         int32
	PackageNameIndex int32
	NameIndex        int32
	VersionIndex     int32
	_io              *kaitai.Stream
	_root            *Bir
	_parent          kaitai.Struct
}

func NewBir_PackageCpInfo() *Bir_PackageCpInfo {
	return &Bir_PackageCpInfo{}
}

func (this Bir_PackageCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_PackageCpInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp538, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OrgIndex = int32(tmp538)
	tmp539, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PackageNameIndex = int32(tmp539)
	tmp540, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameIndex = int32(tmp540)
	tmp541, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.VersionIndex = int32(tmp541)
	return err
}

type Bir_PathParam struct {
	PathParamNameCpIndex int32
	PathParamTypeCpIndex int32
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_ResourceFunctionContent
}

func NewBir_PathParam() *Bir_PathParam {
	return &Bir_PathParam{}
}

func (this Bir_PathParam) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_PathParam) Read(io *kaitai.Stream, parent *Bir_ResourceFunctionContent, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp542, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PathParamNameCpIndex = int32(tmp542)
	tmp543, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PathParamTypeCpIndex = int32(tmp543)
	return err
}

type Bir_Position struct {
	SourceFileCpIndex int32
	SLine             int32
	SCol              int32
	ELine             int32
	ECol              int32
	_io               *kaitai.Stream
	_root             *Bir
	_parent           kaitai.Struct
}

func NewBir_Position() *Bir_Position {
	return &Bir_Position{}
}

func (this Bir_Position) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Position) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp544, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.SourceFileCpIndex = int32(tmp544)
	tmp545, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.SLine = int32(tmp545)
	tmp546, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.SCol = int32(tmp546)
	tmp547, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ELine = int32(tmp547)
	tmp548, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ECol = int32(tmp548)
	return err
}

type Bir_ReceiveField struct {
	FieldName   int32
	ChannelName int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_InstructionWkMulReceive
}

func NewBir_ReceiveField() *Bir_ReceiveField {
	return &Bir_ReceiveField{}
}

func (this Bir_ReceiveField) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ReceiveField) Read(io *kaitai.Stream, parent *Bir_InstructionWkMulReceive, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp549, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FieldName = int32(tmp549)
	tmp550, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ChannelName = int32(tmp550)
	return err
}

type Bir_Reciever struct {
	Kind        int8
	TypeCpIndex int32
	NameCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_Function
}

func NewBir_Reciever() *Bir_Reciever {
	return &Bir_Reciever{}
}

func (this Bir_Reciever) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Reciever) Read(io *kaitai.Stream, parent *Bir_Function, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp551, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp551
	tmp552, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp552)
	tmp553, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp553)
	return err
}

type Bir_RecordField struct {
	NameCpIndex                  int32
	Flags                        int64
	Doc                          *Bir_Markdown
	TypeCpIndex                  int32
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_TypeRecord
}

func NewBir_RecordField() *Bir_RecordField {
	return &Bir_RecordField{}
}

func (this Bir_RecordField) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_RecordField) Read(io *kaitai.Stream, parent *Bir_TypeRecord, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp554, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp554)
	tmp555, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp555)
	tmp556 := NewBir_Markdown()
	err = tmp556.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp556
	tmp557, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp557)
	tmp558 := NewBir_AnnotationAttachmentsContent()
	err = tmp558.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp558
	return err
}

type Bir_RecordInitFunction struct {
	NameCpIndex int32
	Flags       int64
	TypeCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     kaitai.Struct
}

func NewBir_RecordInitFunction() *Bir_RecordInitFunction {
	return &Bir_RecordInitFunction{}
}

func (this Bir_RecordInitFunction) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_RecordInitFunction) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp559, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp559)
	tmp560, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp560)
	tmp561, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp561)
	return err
}

type Bir_ReferencedType struct {
	TypeCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_TypeDefinitionBody
}

func NewBir_ReferencedType() *Bir_ReferencedType {
	return &Bir_ReferencedType{}
}

func (this Bir_ReferencedType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ReferencedType) Read(io *kaitai.Stream, parent *Bir_TypeDefinitionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp562, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp562)
	return err
}

type Bir_RequiredParam struct {
	ParamNameCpIndex int32
	Flags            int64
	ParamAnnotations *Bir_AnnotationAttachmentsContent
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_Function
}

func NewBir_RequiredParam() *Bir_RequiredParam {
	return &Bir_RequiredParam{}
}

func (this Bir_RequiredParam) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_RequiredParam) Read(io *kaitai.Stream, parent *Bir_Function, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp563, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamNameCpIndex = int32(tmp563)
	tmp564, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp564)
	tmp565 := NewBir_AnnotationAttachmentsContent()
	err = tmp565.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ParamAnnotations = tmp565
	return err
}

type Bir_ResourceFunctionContent struct {
	PathParamsCount          int32
	PathParams               []*Bir_PathParam
	HasRestPathParam         uint8
	RestPathParam            *Bir_PathParam
	ResourcePathSegmentCount int32
	ResourcePathSegments     []*Bir_ResourcePathSegment
	ResourceAccessor         int32
	_io                      *kaitai.Stream
	_root                    *Bir
	_parent                  *Bir_Function
}

func NewBir_ResourceFunctionContent() *Bir_ResourceFunctionContent {
	return &Bir_ResourceFunctionContent{}
}

func (this Bir_ResourceFunctionContent) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ResourceFunctionContent) Read(io *kaitai.Stream, parent *Bir_Function, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp566, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PathParamsCount = int32(tmp566)
	for i := 0; i < int(this.PathParamsCount); i++ {
		_ = i
		tmp567 := NewBir_PathParam()
		err = tmp567.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.PathParams = append(this.PathParams, tmp567)
	}
	tmp568, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasRestPathParam = tmp568
	if this.HasRestPathParam == 1 {
		tmp569 := NewBir_PathParam()
		err = tmp569.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RestPathParam = tmp569
	}
	tmp570, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ResourcePathSegmentCount = int32(tmp570)
	for i := 0; i < int(this.ResourcePathSegmentCount); i++ {
		_ = i
		tmp571 := NewBir_ResourcePathSegment()
		err = tmp571.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ResourcePathSegments = append(this.ResourcePathSegments, tmp571)
	}
	tmp572, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ResourceAccessor = int32(tmp572)
	return err
}

type Bir_ResourcePathSegment struct {
	ResourcePathSegmentCpIndex int32
	ResourcePathSegmentPos     *Bir_Position
	ResourcePathSegmentType    int32
	_io                        *kaitai.Stream
	_root                      *Bir
	_parent                    *Bir_ResourceFunctionContent
}

func NewBir_ResourcePathSegment() *Bir_ResourcePathSegment {
	return &Bir_ResourcePathSegment{}
}

func (this Bir_ResourcePathSegment) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ResourcePathSegment) Read(io *kaitai.Stream, parent *Bir_ResourceFunctionContent, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp573, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ResourcePathSegmentCpIndex = int32(tmp573)
	tmp574 := NewBir_Position()
	err = tmp574.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ResourcePathSegmentPos = tmp574
	tmp575, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ResourcePathSegmentType = int32(tmp575)
	return err
}

type Bir_ReturnVar struct {
	Kind        int8
	TypeCpIndex int32
	NameCpIndex int32
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_FunctionBody
}

func NewBir_ReturnVar() *Bir_ReturnVar {
	return &Bir_ReturnVar{}
}

func (this Bir_ReturnVar) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ReturnVar) Read(io *kaitai.Stream, parent *Bir_FunctionBody, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp576, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp576
	tmp577, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp577)
	tmp578, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp578)
	return err
}

type Bir_ScopeEntry struct {
	CurrentScopeIndex int32
	InstructionOffset int32
	HasParent         uint8
	ParentScopeIndex  int32
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_Function
}

func NewBir_ScopeEntry() *Bir_ScopeEntry {
	return &Bir_ScopeEntry{}
}

func (this Bir_ScopeEntry) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ScopeEntry) Read(io *kaitai.Stream, parent *Bir_Function, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp579, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.CurrentScopeIndex = int32(tmp579)
	tmp580, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InstructionOffset = int32(tmp580)
	tmp581, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasParent = tmp581
	if this.HasParent == 1 {
		tmp582, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.ParentScopeIndex = int32(tmp582)
	}
	return err
}

type Bir_SemNamedType struct {
	Semtype      *Bir_SemtypeInfo
	OptionalName *Bir_NullableStrInfo
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_TypeFinite
}

func NewBir_SemNamedType() *Bir_SemNamedType {
	return &Bir_SemNamedType{}
}

func (this Bir_SemNamedType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemNamedType) Read(io *kaitai.Stream, parent *Bir_TypeFinite, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp583 := NewBir_SemtypeInfo()
	err = tmp583.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Semtype = tmp583
	tmp584 := NewBir_NullableStrInfo()
	err = tmp584.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.OptionalName = tmp584
	return err
}

type Bir_SemtypeBdd struct {
	IsBddNode       uint8
	BddNode         *Bir_SemtypeBddNode
	BddAllOrNothing uint8
	_io             *kaitai.Stream
	_root           *Bir
	_parent         kaitai.Struct
}

func NewBir_SemtypeBdd() *Bir_SemtypeBdd {
	return &Bir_SemtypeBdd{}
}

func (this Bir_SemtypeBdd) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeBdd) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp585, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsBddNode = tmp585
	if this.IsBddNode == 1 {
		tmp586 := NewBir_SemtypeBddNode()
		err = tmp586.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.BddNode = tmp586
	}
	if this.IsBddNode == 0 {
		tmp587, err := this._io.ReadU1()
		if err != nil {
			return err
		}
		this.BddAllOrNothing = tmp587
	}
	return err
}

type Bir_SemtypeBddNode struct {
	IsRecAtom     uint8
	RecAtomIndex  int32
	TargetKind    int32
	TypeAtom      *Bir_SemtypeTypeAtom
	BddNodeLeft   *Bir_SemtypeBdd
	BddNodeMiddle *Bir_SemtypeBdd
	BddNodeRight  *Bir_SemtypeBdd
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_SemtypeBdd
}

func NewBir_SemtypeBddNode() *Bir_SemtypeBddNode {
	return &Bir_SemtypeBddNode{}
}

func (this Bir_SemtypeBddNode) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeBddNode) Read(io *kaitai.Stream, parent *Bir_SemtypeBdd, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp588, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsRecAtom = tmp588
	if this.IsRecAtom == 1 {
		tmp589, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.RecAtomIndex = int32(tmp589)
	}
	if (this.IsRecAtom == 1) && (this.RecAtomIndex > 1) {
		tmp590, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.TargetKind = int32(tmp590)
	}
	if this.IsRecAtom == 0 {
		tmp591 := NewBir_SemtypeTypeAtom()
		err = tmp591.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeAtom = tmp591
	}
	tmp592 := NewBir_SemtypeBdd()
	err = tmp592.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.BddNodeLeft = tmp592
	tmp593 := NewBir_SemtypeBdd()
	err = tmp593.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.BddNodeMiddle = tmp593
	tmp594 := NewBir_SemtypeBdd()
	err = tmp594.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.BddNodeRight = tmp594
	return err
}

type Bir_SemtypeBooleanSubtype struct {
	Value   uint8
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeBooleanSubtype() *Bir_SemtypeBooleanSubtype {
	return &Bir_SemtypeBooleanSubtype{}
}

func (this Bir_SemtypeBooleanSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeBooleanSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp595, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Value = tmp595
	return err
}

type Bir_SemtypeCellAtomicType struct {
	Ty      *Bir_SemtypeInfo
	Mut     int8
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_SemtypeTypeAtom
}

func NewBir_SemtypeCellAtomicType() *Bir_SemtypeCellAtomicType {
	return &Bir_SemtypeCellAtomicType{}
}

func (this Bir_SemtypeCellAtomicType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeCellAtomicType) Read(io *kaitai.Stream, parent *Bir_SemtypeTypeAtom, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp596 := NewBir_SemtypeInfo()
	err = tmp596.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Ty = tmp596
	tmp597, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Mut = tmp597
	return err
}

type Bir_SemtypeComplex struct {
	AllBitSet             int32
	SomeBitSet            int32
	SubtypeDataListLength int8
	ProperSubtypeData     []*Bir_SemtypeProperSubtypeData
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_SemtypeInternal
}

func NewBir_SemtypeComplex() *Bir_SemtypeComplex {
	return &Bir_SemtypeComplex{}
}

func (this Bir_SemtypeComplex) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeComplex) Read(io *kaitai.Stream, parent *Bir_SemtypeInternal, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp598, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AllBitSet = int32(tmp598)
	tmp599, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.SomeBitSet = int32(tmp599)
	tmp600, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.SubtypeDataListLength = tmp600
	for i := 0; i < int(this.SubtypeDataListLength); i++ {
		_ = i
		tmp601 := NewBir_SemtypeProperSubtypeData()
		err = tmp601.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ProperSubtypeData = append(this.ProperSubtypeData, tmp601)
	}
	return err
}

type Bir_SemtypeDecimalSubtype struct {
	Allowed      uint8
	ValuesLength int32
	Values       []*Bir_SemtypeEnumerableDecimal
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeDecimalSubtype() *Bir_SemtypeDecimalSubtype {
	return &Bir_SemtypeDecimalSubtype{}
}

func (this Bir_SemtypeDecimalSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeDecimalSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp602, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Allowed = tmp602
	tmp603, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValuesLength = int32(tmp603)
	for i := 0; i < int(this.ValuesLength); i++ {
		_ = i
		tmp604 := NewBir_SemtypeEnumerableDecimal()
		err = tmp604.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Values = append(this.Values, tmp604)
	}
	return err
}

type Bir_SemtypeEnumerableDecimal struct {
	Scale                    int32
	UnscaledValueBytesLength int32
	UnscaledValueBytes       []byte
	_io                      *kaitai.Stream
	_root                    *Bir
	_parent                  *Bir_SemtypeDecimalSubtype
}

func NewBir_SemtypeEnumerableDecimal() *Bir_SemtypeEnumerableDecimal {
	return &Bir_SemtypeEnumerableDecimal{}
}

func (this Bir_SemtypeEnumerableDecimal) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeEnumerableDecimal) Read(io *kaitai.Stream, parent *Bir_SemtypeDecimalSubtype, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp605, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Scale = int32(tmp605)
	tmp606, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.UnscaledValueBytesLength = int32(tmp606)
	tmp607, err := this._io.ReadBytes(int(this.UnscaledValueBytesLength))
	if err != nil {
		return err
	}
	this.UnscaledValueBytes = tmp607
	return err
}

type Bir_SemtypeEnumerableString struct {
	StringCpIndex int32
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_SemtypeStringSubtype
}

func NewBir_SemtypeEnumerableString() *Bir_SemtypeEnumerableString {
	return &Bir_SemtypeEnumerableString{}
}

func (this Bir_SemtypeEnumerableString) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeEnumerableString) Read(io *kaitai.Stream, parent *Bir_SemtypeStringSubtype, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp608, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.StringCpIndex = int32(tmp608)
	return err
}

type Bir_SemtypeFloatSubtype struct {
	Allowed      uint8
	ValuesLength int32
	Values       []float64
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeFloatSubtype() *Bir_SemtypeFloatSubtype {
	return &Bir_SemtypeFloatSubtype{}
}

func (this Bir_SemtypeFloatSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeFloatSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp609, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Allowed = tmp609
	tmp610, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValuesLength = int32(tmp610)
	for i := 0; i < int(this.ValuesLength); i++ {
		_ = i
		tmp611, err := this._io.ReadF8be()
		if err != nil {
			return err
		}
		this.Values = append(this.Values, tmp611)
	}
	return err
}

type Bir_SemtypeFunctionAtomicType struct {
	ParamType     *Bir_SemtypeInfo
	RetType       *Bir_SemtypeInfo
	QualifierType *Bir_SemtypeInfo
	IsGeneric     uint8
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_SemtypeTypeAtom
}

func NewBir_SemtypeFunctionAtomicType() *Bir_SemtypeFunctionAtomicType {
	return &Bir_SemtypeFunctionAtomicType{}
}

func (this Bir_SemtypeFunctionAtomicType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeFunctionAtomicType) Read(io *kaitai.Stream, parent *Bir_SemtypeTypeAtom, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp612 := NewBir_SemtypeInfo()
	err = tmp612.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.ParamType = tmp612
	tmp613 := NewBir_SemtypeInfo()
	err = tmp613.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RetType = tmp613
	tmp614 := NewBir_SemtypeInfo()
	err = tmp614.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.QualifierType = tmp614
	tmp615, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsGeneric = tmp615
	return err
}

type Bir_SemtypeInfo struct {
	HasSemtype uint8
	Semtype    *Bir_SemtypeInternal
	_io        *kaitai.Stream
	_root      *Bir
	_parent    kaitai.Struct
}

func NewBir_SemtypeInfo() *Bir_SemtypeInfo {
	return &Bir_SemtypeInfo{}
}

func (this Bir_SemtypeInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp616, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasSemtype = tmp616
	if this.HasSemtype == 1 {
		tmp617 := NewBir_SemtypeInternal()
		err = tmp617.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Semtype = tmp617
	}
	return err
}

type Bir_SemtypeIntSubtype struct {
	RangesLength int32
	X            []*Bir_SemtypeRange
	_io          *kaitai.Stream
	_root        *Bir
	_parent      *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeIntSubtype() *Bir_SemtypeIntSubtype {
	return &Bir_SemtypeIntSubtype{}
}

func (this Bir_SemtypeIntSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeIntSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp618, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.RangesLength = int32(tmp618)
	for i := 0; i < int(this.RangesLength); i++ {
		_ = i
		tmp619 := NewBir_SemtypeRange()
		err = tmp619.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.X = append(this.X, tmp619)
	}
	return err
}

type Bir_SemtypeInternal struct {
	IsUniformTypeBitSet uint8
	UniformTypeBitSet   int32
	ComplexSemtype      *Bir_SemtypeComplex
	_io                 *kaitai.Stream
	_root               *Bir
	_parent             *Bir_SemtypeInfo
}

func NewBir_SemtypeInternal() *Bir_SemtypeInternal {
	return &Bir_SemtypeInternal{}
}

func (this Bir_SemtypeInternal) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeInternal) Read(io *kaitai.Stream, parent *Bir_SemtypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp620, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsUniformTypeBitSet = tmp620
	if this.IsUniformTypeBitSet == 1 {
		tmp621, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.UniformTypeBitSet = int32(tmp621)
	}
	if this.IsUniformTypeBitSet == 0 {
		tmp622 := NewBir_SemtypeComplex()
		err = tmp622.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ComplexSemtype = tmp622
	}
	return err
}

type Bir_SemtypeListAtomicType struct {
	InitialListSize int32
	Initial         []*Bir_SemtypeInfo
	FixedLength     int32
	Rest            *Bir_SemtypeInfo
	_io             *kaitai.Stream
	_root           *Bir
	_parent         *Bir_SemtypeTypeAtom
}

func NewBir_SemtypeListAtomicType() *Bir_SemtypeListAtomicType {
	return &Bir_SemtypeListAtomicType{}
}

func (this Bir_SemtypeListAtomicType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeListAtomicType) Read(io *kaitai.Stream, parent *Bir_SemtypeTypeAtom, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp623, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.InitialListSize = int32(tmp623)
	for i := 0; i < int(this.InitialListSize); i++ {
		_ = i
		tmp624 := NewBir_SemtypeInfo()
		err = tmp624.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Initial = append(this.Initial, tmp624)
	}
	tmp625, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.FixedLength = int32(tmp625)
	tmp626 := NewBir_SemtypeInfo()
	err = tmp626.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Rest = tmp626
	return err
}

type Bir_SemtypeMappingAtomicType struct {
	NamesLength int32
	Names       []int32
	TypesLength int32
	Types       []*Bir_SemtypeInfo
	Rest        *Bir_SemtypeInfo
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_SemtypeTypeAtom
}

func NewBir_SemtypeMappingAtomicType() *Bir_SemtypeMappingAtomicType {
	return &Bir_SemtypeMappingAtomicType{}
}

func (this Bir_SemtypeMappingAtomicType) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeMappingAtomicType) Read(io *kaitai.Stream, parent *Bir_SemtypeTypeAtom, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp627, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NamesLength = int32(tmp627)
	for i := 0; i < int(this.NamesLength); i++ {
		_ = i
		tmp628, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.Names = append(this.Names, tmp628)
	}
	tmp629, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypesLength = int32(tmp629)
	for i := 0; i < int(this.TypesLength); i++ {
		_ = i
		tmp630 := NewBir_SemtypeInfo()
		err = tmp630.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Types = append(this.Types, tmp630)
	}
	tmp631 := NewBir_SemtypeInfo()
	err = tmp631.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Rest = tmp631
	return err
}

type Bir_SemtypeProperSubtypeData struct {
	ProperSubtypeDataKind int8
	Bdd                   *Bir_SemtypeBdd
	IntSubtype            *Bir_SemtypeIntSubtype
	BooleanSubtype        *Bir_SemtypeBooleanSubtype
	FloatSubtype          *Bir_SemtypeFloatSubtype
	DecimalSubtype        *Bir_SemtypeDecimalSubtype
	StringSubtype         *Bir_SemtypeStringSubtype
	XmlSubtype            *Bir_SemtypeXmlSubtype
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_SemtypeComplex
}

func NewBir_SemtypeProperSubtypeData() *Bir_SemtypeProperSubtypeData {
	return &Bir_SemtypeProperSubtypeData{}
}

func (this Bir_SemtypeProperSubtypeData) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeProperSubtypeData) Read(io *kaitai.Stream, parent *Bir_SemtypeComplex, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp632, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.ProperSubtypeDataKind = tmp632
	if this.ProperSubtypeDataKind == 1 {
		tmp633 := NewBir_SemtypeBdd()
		err = tmp633.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Bdd = tmp633
	}
	if this.ProperSubtypeDataKind == 2 {
		tmp634 := NewBir_SemtypeIntSubtype()
		err = tmp634.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.IntSubtype = tmp634
	}
	if this.ProperSubtypeDataKind == 3 {
		tmp635 := NewBir_SemtypeBooleanSubtype()
		err = tmp635.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.BooleanSubtype = tmp635
	}
	if this.ProperSubtypeDataKind == 4 {
		tmp636 := NewBir_SemtypeFloatSubtype()
		err = tmp636.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.FloatSubtype = tmp636
	}
	if this.ProperSubtypeDataKind == 5 {
		tmp637 := NewBir_SemtypeDecimalSubtype()
		err = tmp637.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DecimalSubtype = tmp637
	}
	if this.ProperSubtypeDataKind == 6 {
		tmp638 := NewBir_SemtypeStringSubtype()
		err = tmp638.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.StringSubtype = tmp638
	}
	if this.ProperSubtypeDataKind == 7 {
		tmp639 := NewBir_SemtypeXmlSubtype()
		err = tmp639.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.XmlSubtype = tmp639
	}
	return err
}

type Bir_SemtypeRange struct {
	Min     int64
	Max     int64
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_SemtypeIntSubtype
}

func NewBir_SemtypeRange() *Bir_SemtypeRange {
	return &Bir_SemtypeRange{}
}

func (this Bir_SemtypeRange) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeRange) Read(io *kaitai.Stream, parent *Bir_SemtypeIntSubtype, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp640, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Min = int64(tmp640)
	tmp641, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Max = int64(tmp641)
	return err
}

type Bir_SemtypeStringSubtype struct {
	Allowed       uint8
	ValuesLength  int32
	Values        []*Bir_SemtypeEnumerableString
	Allowed1      uint8
	ValuesLength1 int32
	Values1       []*Bir_SemtypeEnumerableString
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeStringSubtype() *Bir_SemtypeStringSubtype {
	return &Bir_SemtypeStringSubtype{}
}

func (this Bir_SemtypeStringSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeStringSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp642, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Allowed = tmp642
	tmp643, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValuesLength = int32(tmp643)
	for i := 0; i < int(this.ValuesLength); i++ {
		_ = i
		tmp644 := NewBir_SemtypeEnumerableString()
		err = tmp644.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Values = append(this.Values, tmp644)
	}
	tmp645, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.Allowed1 = tmp645
	tmp646, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValuesLength1 = int32(tmp646)
	for i := 0; i < int(this.ValuesLength1); i++ {
		_ = i
		tmp647 := NewBir_SemtypeEnumerableString()
		err = tmp647.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Values1 = append(this.Values1, tmp647)
	}
	return err
}

type Bir_SemtypeTypeAtom struct {
	TypeAtomIndex      int32
	TypeAtomKind       int8
	MappingAtomicType  *Bir_SemtypeMappingAtomicType
	ListAtomicType     *Bir_SemtypeListAtomicType
	FunctionAtomicType *Bir_SemtypeFunctionAtomicType
	CellAtomicType     *Bir_SemtypeCellAtomicType
	_io                *kaitai.Stream
	_root              *Bir
	_parent            *Bir_SemtypeBddNode
}

func NewBir_SemtypeTypeAtom() *Bir_SemtypeTypeAtom {
	return &Bir_SemtypeTypeAtom{}
}

func (this Bir_SemtypeTypeAtom) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeTypeAtom) Read(io *kaitai.Stream, parent *Bir_SemtypeBddNode, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp648, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeAtomIndex = int32(tmp648)
	tmp649, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.TypeAtomKind = tmp649
	if this.TypeAtomKind == 1 {
		tmp650 := NewBir_SemtypeMappingAtomicType()
		err = tmp650.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.MappingAtomicType = tmp650
	}
	if this.TypeAtomKind == 2 {
		tmp651 := NewBir_SemtypeListAtomicType()
		err = tmp651.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ListAtomicType = tmp651
	}
	if this.TypeAtomKind == 3 {
		tmp652 := NewBir_SemtypeFunctionAtomicType()
		err = tmp652.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.FunctionAtomicType = tmp652
	}
	if this.TypeAtomKind == 4 {
		tmp653 := NewBir_SemtypeCellAtomicType()
		err = tmp653.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.CellAtomicType = tmp653
	}
	return err
}

type Bir_SemtypeXmlSubtype struct {
	Primitives int32
	Sequence   *Bir_SemtypeBdd
	_io        *kaitai.Stream
	_root      *Bir
	_parent    *Bir_SemtypeProperSubtypeData
}

func NewBir_SemtypeXmlSubtype() *Bir_SemtypeXmlSubtype {
	return &Bir_SemtypeXmlSubtype{}
}

func (this Bir_SemtypeXmlSubtype) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_SemtypeXmlSubtype) Read(io *kaitai.Stream, parent *Bir_SemtypeProperSubtypeData, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp654, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Primitives = int32(tmp654)
	tmp655 := NewBir_SemtypeBdd()
	err = tmp655.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Sequence = tmp655
	return err
}

type Bir_ServiceDeclaration struct {
	NameCpIndex                int32
	AssociatedClassNameCpIndex int32
	Flags                      int64
	Origin                     int8
	Position                   *Bir_Position
	HasType                    uint8
	TypeCpIndex                int32
	HasAttachPoint             uint8
	AttachPointCount           int32
	AttachPoints               []int32
	HasAttachPointLiteral      uint8
	AttachPointLiteral         int32
	ListenerTypesCount         int32
	ListenerTypes              []*Bir_ListenerType
	_io                        *kaitai.Stream
	_root                      *Bir
	_parent                    *Bir_Module
}

func NewBir_ServiceDeclaration() *Bir_ServiceDeclaration {
	return &Bir_ServiceDeclaration{}
}

func (this Bir_ServiceDeclaration) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ServiceDeclaration) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp656, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp656)
	tmp657, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AssociatedClassNameCpIndex = int32(tmp657)
	tmp658, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp658)
	tmp659, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp659
	tmp660 := NewBir_Position()
	err = tmp660.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp660
	tmp661, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasType = tmp661
	if this.HasType != 0 {
		tmp662, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.TypeCpIndex = int32(tmp662)
	}
	tmp663, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasAttachPoint = tmp663
	if this.HasAttachPoint != 0 {
		tmp664, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.AttachPointCount = int32(tmp664)
	}
	if this.HasAttachPoint != 0 {
		for i := 0; i < int(this.AttachPointCount); i++ {
			_ = i
			tmp665, err := this._io.ReadS4be()
			if err != nil {
				return err
			}
			this.AttachPoints = append(this.AttachPoints, tmp665)
		}
	}
	tmp666, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasAttachPointLiteral = tmp666
	if this.HasAttachPointLiteral != 0 {
		tmp667, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.AttachPointLiteral = int32(tmp667)
	}
	tmp668, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ListenerTypesCount = int32(tmp668)
	for i := 0; i < int(this.ListenerTypesCount); i++ {
		_ = i
		tmp669 := NewBir_ListenerType()
		err = tmp669.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ListenerTypes = append(this.ListenerTypes, tmp669)
	}
	return err
}

type Bir_ShapeCpInfo struct {
	ShapeLength int32
	Shape       *Bir_TypeInfo
	_io         *kaitai.Stream
	_root       *Bir
	_parent     *Bir_ConstantPoolEntry
}

func NewBir_ShapeCpInfo() *Bir_ShapeCpInfo {
	return &Bir_ShapeCpInfo{}
}

func (this Bir_ShapeCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_ShapeCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp670, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ShapeLength = int32(tmp670)
	tmp671 := NewBir_TypeInfo()
	err = tmp671.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Shape = tmp671
	return err
}

type Bir_StringConstantInfo struct {
	ValueCpIndex int32
	_io          *kaitai.Stream
	_root        *Bir
	_parent      kaitai.Struct
}

func NewBir_StringConstantInfo() *Bir_StringConstantInfo {
	return &Bir_StringConstantInfo{}
}

func (this Bir_StringConstantInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_StringConstantInfo) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp672, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueCpIndex = int32(tmp672)
	return err
}

type Bir_StringCpInfo struct {
	StrLen  int32
	Value   string
	_io     *kaitai.Stream
	_root   *Bir
	_parent *Bir_ConstantPoolEntry
}

func NewBir_StringCpInfo() *Bir_StringCpInfo {
	return &Bir_StringCpInfo{}
}

func (this Bir_StringCpInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_StringCpInfo) Read(io *kaitai.Stream, parent *Bir_ConstantPoolEntry, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp673, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.StrLen = int32(tmp673)
	tmp674, err := this._io.ReadBytes(int(this.StrLen))
	if err != nil {
		return err
	}
	this.Value = string(tmp674)
	return err
}

type Bir_TableFieldNameList struct {
	Size             int32
	FieldNameCpIndex []int32
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_TypeTable
}

func NewBir_TableFieldNameList() *Bir_TableFieldNameList {
	return &Bir_TableFieldNameList{}
}

func (this Bir_TableFieldNameList) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TableFieldNameList) Read(io *kaitai.Stream, parent *Bir_TypeTable, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp675, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Size = int32(tmp675)
	for i := 0; i < int(this.Size); i++ {
		_ = i
		tmp676, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.FieldNameCpIndex = append(this.FieldNameCpIndex, tmp676)
	}
	return err
}

type Bir_TupleMember struct {
	NameCpIndex                  int32
	Flags                        int64
	TypeCpIndex                  int32
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_TypeTuple
}

func NewBir_TupleMember() *Bir_TupleMember {
	return &Bir_TupleMember{}
}

func (this Bir_TupleMember) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TupleMember) Read(io *kaitai.Stream, parent *Bir_TypeTuple, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp677, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp677)
	tmp678, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp678)
	tmp679, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp679)
	tmp680 := NewBir_AnnotationAttachmentsContent()
	err = tmp680.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp680
	return err
}

type Bir_TypeArray struct {
	State            int8
	Size             int32
	ElementTypeIndex int32
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_TypeInfo
}

func NewBir_TypeArray() *Bir_TypeArray {
	return &Bir_TypeArray{}
}

func (this Bir_TypeArray) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeArray) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp681, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.State = tmp681
	tmp682, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.Size = int32(tmp682)
	tmp683, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ElementTypeIndex = int32(tmp683)
	return err
}

type Bir_TypeDefinition struct {
	Position                     *Bir_Position
	NameCpIndex                  int32
	OriginalNameCpIndex          int32
	Flags                        int64
	Origin                       int8
	Doc                          *Bir_Markdown
	TypeCpIndex                  int32
	HasReferenceType             uint8
	AnnotationAttachmentsContent *Bir_AnnotationAttachmentsContent
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_Module
}

func NewBir_TypeDefinition() *Bir_TypeDefinition {
	return &Bir_TypeDefinition{}
}

func (this Bir_TypeDefinition) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeDefinition) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp684 := NewBir_Position()
	err = tmp684.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Position = tmp684
	tmp685, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp685)
	tmp686, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OriginalNameCpIndex = int32(tmp686)
	tmp687, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp687)
	tmp688, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Origin = tmp688
	tmp689 := NewBir_Markdown()
	err = tmp689.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Doc = tmp689
	tmp690, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeCpIndex = int32(tmp690)
	tmp691, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasReferenceType = tmp691
	tmp692 := NewBir_AnnotationAttachmentsContent()
	err = tmp692.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.AnnotationAttachmentsContent = tmp692
	return err
}

type Bir_TypeDefinitionBody struct {
	AttachedFunctionsCount int32
	AttachedFunctions      []*Bir_Function
	ReferencedTypesCount   int32
	ReferencedTypes        []*Bir_ReferencedType
	_io                    *kaitai.Stream
	_root                  *Bir
	_parent                *Bir_Module
}

func NewBir_TypeDefinitionBody() *Bir_TypeDefinitionBody {
	return &Bir_TypeDefinitionBody{}
}

func (this Bir_TypeDefinitionBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeDefinitionBody) Read(io *kaitai.Stream, parent *Bir_Module, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp693, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.AttachedFunctionsCount = int32(tmp693)
	for i := 0; i < int(this.AttachedFunctionsCount); i++ {
		_ = i
		tmp694 := NewBir_Function()
		err = tmp694.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.AttachedFunctions = append(this.AttachedFunctions, tmp694)
	}
	tmp695, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ReferencedTypesCount = int32(tmp695)
	for i := 0; i < int(this.ReferencedTypesCount); i++ {
		_ = i
		tmp696 := NewBir_ReferencedType()
		err = tmp696.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ReferencedTypes = append(this.ReferencedTypes, tmp696)
	}
	return err
}

type Bir_TypeError struct {
	PkgIdCpIndex         int32
	ErrorTypeNameCpIndex int32
	DetailTypeCpIndex    int32
	TypeIds              *Bir_TypeId
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              *Bir_TypeInfo
}

func NewBir_TypeError() *Bir_TypeError {
	return &Bir_TypeError{}
}

func (this Bir_TypeError) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeError) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp697, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkgIdCpIndex = int32(tmp697)
	tmp698, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ErrorTypeNameCpIndex = int32(tmp698)
	tmp699, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DetailTypeCpIndex = int32(tmp699)
	tmp700 := NewBir_TypeId()
	err = tmp700.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TypeIds = tmp700
	return err
}

type Bir_TypeFinite struct {
	NameCpIndex    int32
	Flags          int64
	ValueSpaceSize int32
	ValueSpace     []*Bir_SemNamedType
	_io            *kaitai.Stream
	_root          *Bir
	_parent        *Bir_TypeInfo
}

func NewBir_TypeFinite() *Bir_TypeFinite {
	return &Bir_TypeFinite{}
}

func (this Bir_TypeFinite) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeFinite) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp701, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp701)
	tmp702, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.Flags = int64(tmp702)
	tmp703, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ValueSpaceSize = int32(tmp703)
	for i := 0; i < int(this.ValueSpaceSize); i++ {
		_ = i
		tmp704 := NewBir_SemNamedType()
		err = tmp704.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ValueSpace = append(this.ValueSpace, tmp704)
	}
	return err
}

type Bir_TypeFuture struct {
	ConstraintTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeFuture() *Bir_TypeFuture {
	return &Bir_TypeFuture{}
}

func (this Bir_TypeFuture) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeFuture) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp705, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp705)
	return err
}

type Bir_TypeId struct {
	PrimaryTypeIdCount   int32
	PrimaryTypeId        []*Bir_TypeIdSet
	SecondaryTypeIdCount int32
	SecondaryTypeId      []*Bir_TypeIdSet
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              kaitai.Struct
}

func NewBir_TypeId() *Bir_TypeId {
	return &Bir_TypeId{}
}

func (this Bir_TypeId) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeId) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp706, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PrimaryTypeIdCount = int32(tmp706)
	for i := 0; i < int(this.PrimaryTypeIdCount); i++ {
		_ = i
		tmp707 := NewBir_TypeIdSet()
		err = tmp707.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.PrimaryTypeId = append(this.PrimaryTypeId, tmp707)
	}
	tmp708, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.SecondaryTypeIdCount = int32(tmp708)
	for i := 0; i < int(this.SecondaryTypeIdCount); i++ {
		_ = i
		tmp709 := NewBir_TypeIdSet()
		err = tmp709.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.SecondaryTypeId = append(this.SecondaryTypeId, tmp709)
	}
	return err
}

type Bir_TypeIdSet struct {
	PkgIdCpIndex      int32
	TypeIdNameCpIndex int32
	IsPublicId        uint8
	_io               *kaitai.Stream
	_root             *Bir
	_parent           *Bir_TypeId
}

func NewBir_TypeIdSet() *Bir_TypeIdSet {
	return &Bir_TypeIdSet{}
}

func (this Bir_TypeIdSet) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeIdSet) Read(io *kaitai.Stream, parent *Bir_TypeId, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp710, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkgIdCpIndex = int32(tmp710)
	tmp711, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeIdNameCpIndex = int32(tmp711)
	tmp712, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsPublicId = tmp712
	return err
}

type Bir_TypeInfo struct {
	TypeTag       Bir_TypeTagEnum
	NameIndex     int32
	TypeFlag      int64
	TypeStructure kaitai.Struct
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_ShapeCpInfo
	_f_nameAsStr  bool
	nameAsStr     string
}

func NewBir_TypeInfo() *Bir_TypeInfo {
	return &Bir_TypeInfo{}
}

func (this Bir_TypeInfo) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeInfo) Read(io *kaitai.Stream, parent *Bir_ShapeCpInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp713, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.TypeTag = Bir_TypeTagEnum(tmp713)
	tmp714, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameIndex = int32(tmp714)
	tmp715, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.TypeFlag = int64(tmp715)
	switch this.TypeTag {
	case Bir_TypeTagEnum__TypeTagArray:
		tmp716 := NewBir_TypeArray()
		err = tmp716.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp716
	case Bir_TypeTagEnum__TypeTagError:
		tmp717 := NewBir_TypeError()
		err = tmp717.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp717
	case Bir_TypeTagEnum__TypeTagFinite:
		tmp718 := NewBir_TypeFinite()
		err = tmp718.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp718
	case Bir_TypeTagEnum__TypeTagFuture:
		tmp719 := NewBir_TypeFuture()
		err = tmp719.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp719
	case Bir_TypeTagEnum__TypeTagIntersection:
		tmp720 := NewBir_TypeIntersection()
		err = tmp720.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp720
	case Bir_TypeTagEnum__TypeTagInvokable:
		tmp721 := NewBir_TypeInvokable()
		err = tmp721.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp721
	case Bir_TypeTagEnum__TypeTagMap:
		tmp722 := NewBir_TypeMap()
		err = tmp722.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp722
	case Bir_TypeTagEnum__TypeTagObjectOrService:
		tmp723 := NewBir_TypeObjectOrService()
		err = tmp723.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp723
	case Bir_TypeTagEnum__TypeTagParameterizedType:
		tmp724 := NewBir_TypeParameterized()
		err = tmp724.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp724
	case Bir_TypeTagEnum__TypeTagRecord:
		tmp725 := NewBir_TypeRecord()
		err = tmp725.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp725
	case Bir_TypeTagEnum__TypeTagStream:
		tmp726 := NewBir_TypeStream()
		err = tmp726.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp726
	case Bir_TypeTagEnum__TypeTagTable:
		tmp727 := NewBir_TypeTable()
		err = tmp727.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp727
	case Bir_TypeTagEnum__TypeTagTuple:
		tmp728 := NewBir_TypeTuple()
		err = tmp728.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp728
	case Bir_TypeTagEnum__TypeTagTypedesc:
		tmp729 := NewBir_TypeTypedesc()
		err = tmp729.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp729
	case Bir_TypeTagEnum__TypeTagTyperefdesc:
		tmp730 := NewBir_TypeTyperefdesc()
		err = tmp730.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp730
	case Bir_TypeTagEnum__TypeTagUnion:
		tmp731 := NewBir_TypeUnion()
		err = tmp731.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp731
	case Bir_TypeTagEnum__TypeTagXml:
		tmp732 := NewBir_TypeXml()
		err = tmp732.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TypeStructure = tmp732
	}
	return err
}

func (this *Bir_TypeInfo) NameAsStr() (v string, err error) {
	if this._f_nameAsStr {
		return this.nameAsStr, nil
	}
	this._f_nameAsStr = true
	this.nameAsStr = string(this._root.ConstantPool.ConstantPoolEntries[this.NameIndex].CpInfo.(*Bir_StringCpInfo).Value)
	return this.nameAsStr, nil
}

type Bir_TypeIntersection struct {
	ConstituentTypesCount  int32
	ConstituentTypeCpIndex []int32
	EffectiveTypeCpIndex   int32
	_io                    *kaitai.Stream
	_root                  *Bir
	_parent                *Bir_TypeInfo
}

func NewBir_TypeIntersection() *Bir_TypeIntersection {
	return &Bir_TypeIntersection{}
}

func (this Bir_TypeIntersection) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeIntersection) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp733, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstituentTypesCount = int32(tmp733)
	for i := 0; i < int(this.ConstituentTypesCount); i++ {
		_ = i
		tmp734, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.ConstituentTypeCpIndex = append(this.ConstituentTypeCpIndex, tmp734)
	}
	tmp735, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.EffectiveTypeCpIndex = int32(tmp735)
	return err
}

type Bir_TypeInvokable struct {
	IsAnyFunction uint8
	InvokableKind *Bir_TypeInvokableBody
	_io           *kaitai.Stream
	_root         *Bir
	_parent       *Bir_TypeInfo
}

func NewBir_TypeInvokable() *Bir_TypeInvokable {
	return &Bir_TypeInvokable{}
}

func (this Bir_TypeInvokable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeInvokable) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp736, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsAnyFunction = tmp736
	if this.IsAnyFunction == 0 {
		tmp737 := NewBir_TypeInvokableBody()
		err = tmp737.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InvokableKind = tmp737
	}
	return err
}

type Bir_TypeInvokableBody struct {
	ParamTypesCount        int32
	ParamTypeCpIndex       []int32
	HasRestType            uint8
	RestTypeCpIndex        int32
	ReturnTypeCpIndex      int32
	HasInvokableTypeSymbol uint8
	InvokableTypeSymbol    *Bir_InvokableTypeSymbolBody
	_io                    *kaitai.Stream
	_root                  *Bir
	_parent                *Bir_TypeInvokable
}

func NewBir_TypeInvokableBody() *Bir_TypeInvokableBody {
	return &Bir_TypeInvokableBody{}
}

func (this Bir_TypeInvokableBody) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeInvokableBody) Read(io *kaitai.Stream, parent *Bir_TypeInvokable, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp738, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamTypesCount = int32(tmp738)
	for i := 0; i < int(this.ParamTypesCount); i++ {
		_ = i
		tmp739, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.ParamTypeCpIndex = append(this.ParamTypeCpIndex, tmp739)
	}
	tmp740, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasRestType = tmp740
	if this.HasRestType == 1 {
		tmp741, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.RestTypeCpIndex = int32(tmp741)
	}
	tmp742, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ReturnTypeCpIndex = int32(tmp742)
	tmp743, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasInvokableTypeSymbol = tmp743
	if this.HasInvokableTypeSymbol == 1 {
		tmp744 := NewBir_InvokableTypeSymbolBody()
		err = tmp744.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InvokableTypeSymbol = tmp744
	}
	return err
}

type Bir_TypeMap struct {
	ConstraintTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeMap() *Bir_TypeMap {
	return &Bir_TypeMap{}
}

func (this Bir_TypeMap) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeMap) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp745, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp745)
	return err
}

type Bir_TypeObjectOrService struct {
	PkdIdCpIndex                 int32
	NameCpIndex                  int32
	ObjectSymbolFlags            int64
	ObjectFieldsCount            int32
	ObjectFields                 []*Bir_ObjectField
	HasGeneratedInitFunction     int8
	GeneratedInitFunction        *Bir_ObjectAttachedFunction
	HasInitFunction              int8
	InitFunction                 *Bir_ObjectAttachedFunction
	ObjectAttachedFunctionsCount int32
	ObjectAttachedFunctions      []*Bir_ObjectAttachedFunction
	TypeInclusionsCount          int32
	TypeInclusionsCpIndex        []int32
	TypeIds                      *Bir_TypeId
	_io                          *kaitai.Stream
	_root                        *Bir
	_parent                      *Bir_TypeInfo
}

func NewBir_TypeObjectOrService() *Bir_TypeObjectOrService {
	return &Bir_TypeObjectOrService{}
}

func (this Bir_TypeObjectOrService) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeObjectOrService) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp746, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkdIdCpIndex = int32(tmp746)
	tmp747, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp747)
	tmp748, err := this._io.ReadS8be()
	if err != nil {
		return err
	}
	this.ObjectSymbolFlags = int64(tmp748)
	tmp749, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ObjectFieldsCount = int32(tmp749)
	for i := 0; i < int(this.ObjectFieldsCount); i++ {
		_ = i
		tmp750 := NewBir_ObjectField()
		err = tmp750.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ObjectFields = append(this.ObjectFields, tmp750)
	}
	tmp751, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasGeneratedInitFunction = tmp751
	if this.HasGeneratedInitFunction == 1 {
		tmp752 := NewBir_ObjectAttachedFunction()
		err = tmp752.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.GeneratedInitFunction = tmp752
	}
	tmp753, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasInitFunction = tmp753
	if this.HasInitFunction == 1 {
		tmp754 := NewBir_ObjectAttachedFunction()
		err = tmp754.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.InitFunction = tmp754
	}
	tmp755, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ObjectAttachedFunctionsCount = int32(tmp755)
	for i := 0; i < int(this.ObjectAttachedFunctionsCount); i++ {
		_ = i
		tmp756 := NewBir_ObjectAttachedFunction()
		err = tmp756.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.ObjectAttachedFunctions = append(this.ObjectAttachedFunctions, tmp756)
	}
	tmp757, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeInclusionsCount = int32(tmp757)
	for i := 0; i < int(this.TypeInclusionsCount); i++ {
		_ = i
		tmp758, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.TypeInclusionsCpIndex = append(this.TypeInclusionsCpIndex, tmp758)
	}
	tmp759 := NewBir_TypeId()
	err = tmp759.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.TypeIds = tmp759
	return err
}

type Bir_TypeParameterized struct {
	ParamValueTypeCpIndex int32
	ParamIndex            int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeParameterized() *Bir_TypeParameterized {
	return &Bir_TypeParameterized{}
}

func (this Bir_TypeParameterized) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeParameterized) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp760, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamValueTypeCpIndex = int32(tmp760)
	tmp761, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ParamIndex = int32(tmp761)
	return err
}

type Bir_TypeRecord struct {
	PkdIdCpIndex          int32
	NameCpIndex           int32
	IsSealed              uint8
	RestFieldTypeCpIndex  int32
	RecordFieldsCount     int32
	RecordFields          []*Bir_RecordField
	TypeInclusionsCount   int32
	TypeInclusionsCpIndex []int32
	DefaultValues         int32
	DefaultValue          []*Bir_DefaultValueBody
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeRecord() *Bir_TypeRecord {
	return &Bir_TypeRecord{}
}

func (this Bir_TypeRecord) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeRecord) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp762, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkdIdCpIndex = int32(tmp762)
	tmp763, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp763)
	tmp764, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSealed = tmp764
	tmp765, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.RestFieldTypeCpIndex = int32(tmp765)
	tmp766, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.RecordFieldsCount = int32(tmp766)
	for i := 0; i < int(this.RecordFieldsCount); i++ {
		_ = i
		tmp767 := NewBir_RecordField()
		err = tmp767.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RecordFields = append(this.RecordFields, tmp767)
	}
	tmp768, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TypeInclusionsCount = int32(tmp768)
	for i := 0; i < int(this.TypeInclusionsCount); i++ {
		_ = i
		tmp769, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.TypeInclusionsCpIndex = append(this.TypeInclusionsCpIndex, tmp769)
	}
	tmp770, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.DefaultValues = int32(tmp770)
	for i := 0; i < int(this.DefaultValues); i++ {
		_ = i
		tmp771 := NewBir_DefaultValueBody()
		err = tmp771.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DefaultValue = append(this.DefaultValue, tmp771)
	}
	return err
}

type Bir_TypeStream struct {
	ConstraintTypeCpIndex int32
	CompletionTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeStream() *Bir_TypeStream {
	return &Bir_TypeStream{}
}

func (this Bir_TypeStream) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeStream) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp772, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp772)
	tmp773, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.CompletionTypeCpIndex = int32(tmp773)
	return err
}

type Bir_TypeTable struct {
	ConstraintTypeCpIndex    int32
	HasFieldNameList         uint8
	FieldNameList            *Bir_TableFieldNameList
	HasKeyConstraintType     uint8
	KeyConstraintTypeCpIndex int32
	_io                      *kaitai.Stream
	_root                    *Bir
	_parent                  *Bir_TypeInfo
}

func NewBir_TypeTable() *Bir_TypeTable {
	return &Bir_TypeTable{}
}

func (this Bir_TypeTable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeTable) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp774, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp774)
	tmp775, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasFieldNameList = tmp775
	if this.HasFieldNameList == 1 {
		tmp776 := NewBir_TableFieldNameList()
		err = tmp776.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.FieldNameList = tmp776
	}
	tmp777, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasKeyConstraintType = tmp777
	if this.HasKeyConstraintType == 1 {
		tmp778, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.KeyConstraintTypeCpIndex = int32(tmp778)
	}
	return err
}

type Bir_TypeTuple struct {
	TupleTypesCount  int32
	TupleTypeCpIndex []*Bir_TupleMember
	HasRestType      uint8
	RestTypeCpIndex  int32
	_io              *kaitai.Stream
	_root            *Bir
	_parent          *Bir_TypeInfo
}

func NewBir_TypeTuple() *Bir_TypeTuple {
	return &Bir_TypeTuple{}
}

func (this Bir_TypeTuple) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeTuple) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp779, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.TupleTypesCount = int32(tmp779)
	for i := 0; i < int(this.TupleTypesCount); i++ {
		_ = i
		tmp780 := NewBir_TupleMember()
		err = tmp780.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.TupleTypeCpIndex = append(this.TupleTypeCpIndex, tmp780)
	}
	tmp781, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.HasRestType = tmp781
	if this.HasRestType == 1 {
		tmp782, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.RestTypeCpIndex = int32(tmp782)
	}
	return err
}

type Bir_TypeTypedesc struct {
	ConstraintTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeTypedesc() *Bir_TypeTypedesc {
	return &Bir_TypeTypedesc{}
}

func (this Bir_TypeTypedesc) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeTypedesc) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp783, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp783)
	return err
}

type Bir_TypeTyperefdesc struct {
	PkdIdCpIndex          int32
	NameCpIndex           int32
	ConstraintTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeTyperefdesc() *Bir_TypeTyperefdesc {
	return &Bir_TypeTyperefdesc{}
}

func (this Bir_TypeTyperefdesc) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeTyperefdesc) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp784, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.PkdIdCpIndex = int32(tmp784)
	tmp785, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp785)
	tmp786, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp786)
	return err
}

type Bir_TypeUnion struct {
	IsCyclic                  uint8
	HasName                   int8
	PkdIdCpIndex              int32
	NameCpIndex               int32
	MemberTypesCount          int32
	MemberTypeCpIndex         []int32
	OriginalMemberTypesCount  int32
	OriginalMemberTypeCpIndex []int32
	IsEnumType                uint8
	PkgCpIndex                int32
	EnumName                  int32
	EnumMembersSize           int32
	EnumMembers               []int32
	_io                       *kaitai.Stream
	_root                     *Bir
	_parent                   *Bir_TypeInfo
}

func NewBir_TypeUnion() *Bir_TypeUnion {
	return &Bir_TypeUnion{}
}

func (this Bir_TypeUnion) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeUnion) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp787, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsCyclic = tmp787
	tmp788, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.HasName = tmp788
	if this.HasName == 1 {
		tmp789, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.PkdIdCpIndex = int32(tmp789)
	}
	if this.HasName == 1 {
		tmp790, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.NameCpIndex = int32(tmp790)
	}
	tmp791, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.MemberTypesCount = int32(tmp791)
	for i := 0; i < int(this.MemberTypesCount); i++ {
		_ = i
		tmp792, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.MemberTypeCpIndex = append(this.MemberTypeCpIndex, tmp792)
	}
	tmp793, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.OriginalMemberTypesCount = int32(tmp793)
	for i := 0; i < int(this.OriginalMemberTypesCount); i++ {
		_ = i
		tmp794, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.OriginalMemberTypeCpIndex = append(this.OriginalMemberTypeCpIndex, tmp794)
	}
	tmp795, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsEnumType = tmp795
	if this.IsEnumType == 1 {
		tmp796, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.PkgCpIndex = int32(tmp796)
	}
	if this.IsEnumType == 1 {
		tmp797, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.EnumName = int32(tmp797)
	}
	if this.IsEnumType == 1 {
		tmp798, err := this._io.ReadS4be()
		if err != nil {
			return err
		}
		this.EnumMembersSize = int32(tmp798)
	}
	if this.IsEnumType == 1 {
		for i := 0; i < int(this.EnumMembersSize); i++ {
			_ = i
			tmp799, err := this._io.ReadS4be()
			if err != nil {
				return err
			}
			this.EnumMembers = append(this.EnumMembers, tmp799)
		}
	}
	return err
}

type Bir_TypeXml struct {
	ConstraintTypeCpIndex int32
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_TypeInfo
}

func NewBir_TypeXml() *Bir_TypeXml {
	return &Bir_TypeXml{}
}

func (this Bir_TypeXml) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_TypeXml) Read(io *kaitai.Stream, parent *Bir_TypeInfo, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp800, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ConstraintTypeCpIndex = int32(tmp800)
	return err
}

type Bir_Variable struct {
	Kind                     int8
	Scope                    int8
	VariableDclNameCpIndex   int32
	GlobalOrConstantVariable *Bir_GlobalVariable
	_io                      *kaitai.Stream
	_root                    *Bir
	_parent                  *Bir_Operand
}

func NewBir_Variable() *Bir_Variable {
	return &Bir_Variable{}
}

func (this Bir_Variable) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_Variable) Read(io *kaitai.Stream, parent *Bir_Operand, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp801, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Kind = tmp801
	tmp802, err := this._io.ReadS1()
	if err != nil {
		return err
	}
	this.Scope = tmp802
	tmp803, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.VariableDclNameCpIndex = int32(tmp803)
	if (this.Kind == 5) || (this.Kind == 7) {
		tmp804 := NewBir_GlobalVariable()
		err = tmp804.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.GlobalOrConstantVariable = tmp804
	}
	return err
}

type Bir_WorkerChannel struct {
	ChannelsLength       int32
	WorkerChannelDetails []*Bir_WorkerChannelDetail
	_io                  *kaitai.Stream
	_root                *Bir
	_parent              kaitai.Struct
}

func NewBir_WorkerChannel() *Bir_WorkerChannel {
	return &Bir_WorkerChannel{}
}

func (this Bir_WorkerChannel) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_WorkerChannel) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp805, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.ChannelsLength = int32(tmp805)
	for i := 0; i < int(this.ChannelsLength); i++ {
		_ = i
		tmp806 := NewBir_WorkerChannelDetail()
		err = tmp806.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.WorkerChannelDetails = append(this.WorkerChannelDetails, tmp806)
	}
	return err
}

type Bir_WorkerChannelDetail struct {
	NameCpIndex           int32
	IsChannelInSameStrand uint8
	IsSend                uint8
	_io                   *kaitai.Stream
	_root                 *Bir
	_parent               *Bir_WorkerChannel
}

func NewBir_WorkerChannelDetail() *Bir_WorkerChannelDetail {
	return &Bir_WorkerChannelDetail{}
}

func (this Bir_WorkerChannelDetail) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_WorkerChannelDetail) Read(io *kaitai.Stream, parent *Bir_WorkerChannel, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp807, err := this._io.ReadS4be()
	if err != nil {
		return err
	}
	this.NameCpIndex = int32(tmp807)
	tmp808, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsChannelInSameStrand = tmp808
	tmp809, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.IsSend = tmp809
	return err
}

type Bir_XmlAccess struct {
	LhsOperand *Bir_Operand
	RhsOperand *Bir_Operand
	_io        *kaitai.Stream
	_root      *Bir
	_parent    kaitai.Struct
}

func NewBir_XmlAccess() *Bir_XmlAccess {
	return &Bir_XmlAccess{}
}

func (this Bir_XmlAccess) IO_() *kaitai.Stream {
	return this._io
}

func (this *Bir_XmlAccess) Read(io *kaitai.Stream, parent kaitai.Struct, root *Bir) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp810 := NewBir_Operand()
	err = tmp810.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LhsOperand = tmp810
	tmp811 := NewBir_Operand()
	err = tmp811.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.RhsOperand = tmp811
	return err
}
