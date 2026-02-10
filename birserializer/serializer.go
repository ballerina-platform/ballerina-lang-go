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

package birserializer

import (
	"bytes"
	"fmt"
	"math"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

const (
	BIR_MAGIC   = "\xba\x10\xc0\xde"
	BIR_VERSION = 75
)

type birWriter struct {
	cp *ConstantPool
}

func Marshal(pkg *bir.BIRPackage) ([]byte, error) {
	writer := &birWriter{
		cp: NewConstantPool(),
	}
	return writer.serialize(pkg)
}

func (bw *birWriter) serialize(pkg *bir.BIRPackage) (result []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("BIR serializer failed due to %s", r)
		}
	}()

	birbuf := &bytes.Buffer{}
	bw.writePackageCPEntry(birbuf, pkg.PackageID)
	bw.writeImportModuleDecls(birbuf, pkg)
	bw.writeConstants(birbuf, pkg)
	bw.writeGlobalVars(birbuf, pkg)
	bw.writeFunctions(birbuf, pkg)

	buf := &bytes.Buffer{}
	_, err = buf.Write([]byte(BIR_MAGIC))
	if err != nil {
		panic(fmt.Sprintf("writing BIR magic: %v", err))
	}

	write(buf, int32(BIR_VERSION))

	cpBytes, err := bw.cp.Serialize()
	if err != nil {
		panic(fmt.Sprintf("serializing constant pool: %v", err))
	}

	_, err = buf.Write(cpBytes)
	if err != nil {
		panic(fmt.Sprintf("writing constant pool bytes: %v", err))
	}

	_, err = buf.Write(birbuf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("writing BIR buffer bytes: %v", err))
	}

	return buf.Bytes(), nil
}

func (bw *birWriter) writeImportModuleDecls(buf *bytes.Buffer, pkg *bir.BIRPackage) {
	bw.writeLength(buf, len(pkg.ImportModules))
	for _, imp := range pkg.ImportModules {
		bw.writeStringCPEntry(buf, imp.PackageID.OrgName.Value())
		bw.writeStringCPEntry(buf, imp.PackageID.PkgName.Value())
		bw.writeStringCPEntry(buf, imp.PackageID.Name.Value())
		bw.writeStringCPEntry(buf, imp.PackageID.Version.Value())
	}
}

func (bw *birWriter) writeConstants(buf *bytes.Buffer, pkg *bir.BIRPackage) {
	bw.writeLength(buf, len(pkg.Constants))
	for _, c := range pkg.Constants {
		bw.writeConstant(buf, &c)
	}
}

func (bw *birWriter) writeConstant(buf *bytes.Buffer, constant *bir.BIRConstant) {
	bw.writeStringCPEntry(buf, constant.Name.Value())
	bw.writeFlags(buf, constant.Flags)
	bw.writeOrigin(buf, constant.Origin)
	bw.writePosition(buf, constant.Pos)
	bw.writeType(buf, constant.Type)

	birbuf := &bytes.Buffer{}
	bw.writeType(birbuf, constant.ConstValue.Type)
	bw.writeConstValue(birbuf, &constant.ConstValue)
	bw.writeBufferLength(buf, birbuf)

	_, err := buf.Write(birbuf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("writing constant value bytes: %v", err))
	}
}

func (bw *birWriter) writeGlobalVars(buf *bytes.Buffer, pkg *bir.BIRPackage) {
	bw.writeLength(buf, len(pkg.GlobalVars))
	for _, gv := range pkg.GlobalVars {
		bw.writePosition(buf, gv.Pos)
		bw.writeKind(buf, gv.Kind)
		bw.writeStringCPEntry(buf, gv.Name.Value())
		bw.writeFlags(buf, gv.Flags)
		bw.writeOrigin(buf, gv.Origin)
		bw.writeType(buf, gv.Type)
	}
}

func (bw *birWriter) writeFunctions(buf *bytes.Buffer, pkg *bir.BIRPackage) {
	bw.writeLength(buf, len(pkg.Functions))
	for _, fn := range pkg.Functions {
		bw.writeFunction(buf, &fn)
	}
}

func (bw *birWriter) writeFunction(buf *bytes.Buffer, fn *bir.BIRFunction) {
	bw.writePosition(buf, fn.Pos)
	bw.writeStringCPEntry(buf, fn.Name.Value())
	bw.writeStringCPEntry(buf, fn.OriginalName.Value())
	bw.writeFlags(buf, fn.Flags)
	bw.writeOrigin(buf, fn.Origin)

	bw.writeLength(buf, len(fn.RequiredParams))
	for _, requiredParam := range fn.RequiredParams {
		bw.writeStringCPEntry(buf, requiredParam.Name.Value())
		bw.writeFlags(buf, requiredParam.Flags)
	}

	birbuf := &bytes.Buffer{}
	bw.writeLength(birbuf, fn.ArgsCount)

	write(birbuf, fn.ReturnVariable != nil)
	if fn.ReturnVariable != nil {
		bw.writeKind(birbuf, fn.ReturnVariable.Kind)
		bw.writeType(birbuf, fn.ReturnVariable.Type)
		bw.writeStringCPEntry(birbuf, fn.ReturnVariable.Name.Value())
	}

	bw.writeLength(birbuf, len(fn.LocalVars))
	for _, localVar := range fn.LocalVars {
		bw.writeLocalVar(birbuf, &localVar)
	}

	bw.writeLength(birbuf, len(fn.BasicBlocks))
	for _, bb := range fn.BasicBlocks {
		bw.writeBasicBlock(birbuf, &bb)
	}

	bw.writeBufferLength(buf, birbuf)
	_, err := buf.Write(birbuf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("writing function body bytes: %v", err))
	}
}

func (bw *birWriter) writeLocalVar(buf *bytes.Buffer, localVar *bir.BIRVariableDcl) {
	bw.writeKind(buf, localVar.Kind)
	bw.writeType(buf, localVar.Type)
	bw.writeStringCPEntry(buf, localVar.Name.Value())

	if localVar.Kind == bir.VAR_KIND_ARG {
		bw.writeStringCPEntry(buf, localVar.MetaVarName)
	}

	if localVar.Kind == bir.VAR_KIND_LOCAL {
		bw.writeStringCPEntry(buf, localVar.MetaVarName)

		endBBId := ""
		if localVar.EndBB != nil {
			endBBId = localVar.EndBB.Id.Value()
		}
		bw.writeStringCPEntry(buf, endBBId)

		startBBId := ""
		if localVar.StartBB != nil {
			startBBId = localVar.StartBB.Id.Value()
		}
		bw.writeStringCPEntry(buf, startBBId)
		bw.writeLength(buf, localVar.InsOffset)
	}
}

func (bw *birWriter) writeBasicBlock(buf *bytes.Buffer, bb *bir.BIRBasicBlock) {
	bw.writeStringCPEntry(buf, bb.Id.Value())
	bw.writeLength(buf, len(bb.Instructions))

	for _, instr := range bb.Instructions {
		bw.writeInstructionKind(buf, instr.GetKind())
		bw.writeInstruction(buf, instr)
	}

	if bb.Terminator == nil {
		write(buf, uint8(0))
		return
	}
	bw.writeInstructionKind(buf, bb.Terminator.GetKind())
	bw.writeTerminator(buf, bb.Terminator)
}

func (bw *birWriter) writeInstruction(buf *bytes.Buffer, instr bir.BIRInstruction) {
	switch instr := instr.(type) {
	case *bir.Move:
		bw.writeOperand(buf, instr.RhsOp)
		bw.writeOperand(buf, instr.LhsOp)

	case *bir.BinaryOp:
		bw.writeOperand(buf, &instr.RhsOp1)
		bw.writeOperand(buf, &instr.RhsOp2)
		bw.writeOperand(buf, instr.LhsOp)

	case *bir.UnaryOp:
		bw.writeOperand(buf, instr.RhsOp)
		bw.writeOperand(buf, instr.LhsOp)

	case *bir.ConstantLoad:
		bw.writeType(buf, instr.Type)
		bw.writeOperand(buf, instr.LhsOp)

		isWrapped := false
		var tagValue any = instr.Value
		if cv, ok := instr.Value.(bir.ConstValue); ok {
			isWrapped = true
			tagValue = cv.Value
		}

		tag, err := bw.inferTag(tagValue)
		if err != nil {
			panic(fmt.Sprintf("inferring constant load tag: %v", err))
		}

		write(buf, isWrapped)
		write(buf, int8(tag))
		bw.writeConstValueByTag(buf, tag, tagValue)

	case *bir.FieldAccess:
		// TODO: MAP_LOAD and ARRAY_LOAD
		bw.writeOperand(buf, instr.LhsOp)
		bw.writeOperand(buf, instr.KeyOp)
		bw.writeOperand(buf, instr.RhsOp)

	case *bir.NewArray:
		bw.writeType(buf, instr.Type)
		bw.writeOperand(buf, instr.LhsOp)
		bw.writeOperand(buf, instr.SizeOp)

	default:
		panic(fmt.Sprintf("unsupported instruction type: %T", instr))
	}
}

func (bw *birWriter) writeTerminator(buf *bytes.Buffer, term bir.BIRTerminator) {
	switch term := term.(type) {
	case *bir.Goto:
		bw.writeStringCPEntry(buf, term.ThenBB.Id.Value())

	case *bir.Branch:
		bw.writeOperand(buf, term.Op)
		bw.writeStringCPEntry(buf, term.TrueBB.Id.Value())
		bw.writeStringCPEntry(buf, term.FalseBB.Id.Value())

	case *bir.Call:
		write(buf, term.IsVirtual)
		bw.writePackageCPEntry(buf, term.CalleePkg)
		bw.writeStringCPEntry(buf, term.Name.Value())

		bw.writeLength(buf, len(term.Args))
		for _, arg := range term.Args {
			bw.writeOperand(buf, &arg)
		}

		if term.LhsOp != nil {
			write(buf, uint8(1))
			bw.writeOperand(buf, term.LhsOp)
		} else {
			write(buf, uint8(0))
		}

		bw.writeStringCPEntry(buf, term.ThenBB.Id.Value())

	case *bir.Return:

	default:
		panic(fmt.Sprintf("unsupported terminator type: %T", term))
	}
}

func (bw *birWriter) writeOperand(buf *bytes.Buffer, op *bir.BIROperand) {
	if op == nil || op.VariableDcl == nil {
		write(buf, false)
		write(buf, uint8(bir.VAR_KIND_TEMP))
		bw.writeScope(buf, bir.VAR_SCOPE_FUNCTION)
		bw.writeStringCPEntry(buf, "")
		return
	}

	if op.VariableDcl.IgnoreVariable {
		write(buf, true)
		bw.writeType(buf, op.VariableDcl.Type)
		return
	}

	write(buf, false)
	bw.writeKind(buf, op.VariableDcl.Kind)
	bw.writeScope(buf, op.VariableDcl.Scope)
	bw.writeStringCPEntry(buf, op.VariableDcl.Name.Value())
}

func (bw *birWriter) writeConstValue(buf *bytes.Buffer, cv *bir.ConstValue) {
	tag, err := bw.inferTag(cv.Value)
	if err != nil {
		panic(fmt.Sprintf("inferring const value tag: %v", err))
	}

	write(buf, int8(tag))
	bw.writeConstValueByTag(buf, tag, cv.Value)
}

func (bw *birWriter) writeConstValueByTag(buf *bytes.Buffer, tag model.TypeTags, value any) {
	if cv, ok := value.(bir.ConstValue); ok {
		bw.writeConstValueByTag(buf, tag, cv.Value)
		return
	}

	switch tag {
	case model.TypeTags_INT,
		model.TypeTags_SIGNED32_INT,
		model.TypeTags_SIGNED16_INT,
		model.TypeTags_SIGNED8_INT,
		model.TypeTags_UNSIGNED32_INT,
		model.TypeTags_UNSIGNED16_INT,
		model.TypeTags_UNSIGNED8_INT:
		var val int64
		switch v := value.(type) {
		case int64:
			val = v
		case int:
			val = int64(v)
		case int32:
			val = int64(v)
		case int16:
			val = int64(v)
		case int8:
			val = int64(v)
		default:
			panic(fmt.Sprintf("expected integer for tag %v, got %T", tag, value))
		}
		bw.writeIntegerInline(buf, val)

	case model.TypeTags_BYTE:
		var val byte
		switch v := value.(type) {
		case byte:
			val = v
		case int:
			val = byte(v)
		case int32:
			val = byte(v)
		default:
			panic(fmt.Sprintf("expected byte for tag %v, got %T", tag, value))
		}
		bw.writeByteInline(buf, val)

	case model.TypeTags_FLOAT:
		var val float64
		switch v := value.(type) {
		case float64:
			val = v
		case float32:
			val = float64(v)
		default:
			panic(fmt.Sprintf("expected float for tag %v, got %T", tag, value))
		}
		bw.writeFloatInline(buf, val)

	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		var val string
		switch v := value.(type) {
		case string:
			val = v
		case *string:
			if v != nil {
				val = *v
			} else {
				val = ""
			}
		default:
			panic(fmt.Sprintf("expected string for tag %v, got %T", tag, value))
		}
		cpIdx := bw.cp.AddStringCPEntry(val)
		write(buf, cpIdx)

	case model.TypeTags_BOOLEAN:
		var val bool
		switch v := value.(type) {
		case bool:
			val = v
		default:
			panic(fmt.Sprintf("expected boolean for tag %v, got %T", tag, value))
		}
		bw.writeBooleanInline(buf, val)

	case model.TypeTags_NIL:
		write(buf, int32(-1))

	default:
		panic(fmt.Sprintf("unsupported tag for constant value: %v", tag))
	}
}

func (bw *birWriter) writeIntegerInline(buf *bytes.Buffer, value int64) {
	write(buf, value)
}

func (bw *birWriter) writeFloatInline(buf *bytes.Buffer, value float64) {
	write(buf, value)
}

func (bw *birWriter) writeBooleanInline(buf *bytes.Buffer, value bool) {
	write(buf, value)
}

func (bw *birWriter) writeByteInline(buf *bytes.Buffer, value byte) {
	write(buf, value)
}

func (bw *birWriter) addValueToCP(tag model.TypeTags, value any) (int32, error) {
	if cv, ok := value.(bir.ConstValue); ok {
		return bw.addValueToCP(tag, cv.Value)
	}

	switch tag {
	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		var val string
		switch v := value.(type) {
		case string:
			val = v
		case *string:
			if v != nil {
				val = *v
			} else {
				val = ""
			}
		default:
			return 0, fmt.Errorf("expected string for tag %v, got %T", tag, value)
		}
		return bw.cp.AddStringCPEntry(val), nil
	default:
		return 0, fmt.Errorf("primitive values should not be added to CP, use writeConstValueByTag instead")
	}
}

// FIXME: Remove this after implementing types
func (bw *birWriter) inferTag(value any) (model.TypeTags, error) {
	switch v := value.(type) {
	case bir.ConstValue:
		return bw.inferTag(v.Value)
	case int, int64, int32, int16, int8:
		return model.TypeTags_INT, nil
	case float64, float32:
		return model.TypeTags_FLOAT, nil
	case string, *string:
		return model.TypeTags_STRING, nil
	case bool:
		return model.TypeTags_BOOLEAN, nil
	case byte:
		return model.TypeTags_BYTE, nil
	case nil:
		return model.TypeTags_NIL, nil
	default:
		return 0, fmt.Errorf("cannot infer tag for value %v (%T)", value, value)
	}
}

func (bw *birWriter) writeKind(buf *bytes.Buffer, kind bir.VarKind) {
	write(buf, uint8(kind))
}

func (bw *birWriter) writeFlags(buf *bytes.Buffer, flags int64) {
	write(buf, flags)
}

func (bw *birWriter) writeOrigin(buf *bytes.Buffer, origin model.SymbolOrigin) {
	write(buf, uint8(origin))
}

func (bw *birWriter) writeStringCPEntry(buf *bytes.Buffer, str string) {
	write(buf, bw.cp.AddStringCPEntry(str))
}

func (bw *birWriter) writeLength(buf *bytes.Buffer, length int) {
	write(buf, int64(length))
}

func (bw *birWriter) writeInstructionKind(buf *bytes.Buffer, kind bir.InstructionKind) {
	write(buf, uint8(kind))
}

func (bw *birWriter) writeScope(buf *bytes.Buffer, scope bir.VarScope) {
	write(buf, uint8(scope))
}

func (bw *birWriter) writePackageCPEntry(buf *bytes.Buffer, pkgID *model.PackageID) {
	pkgIdx := int32(-1)
	if pkgID != nil {
		pkgIdx = bw.cp.AddPackageCPEntry(pkgID)
	}
	write(buf, pkgIdx)
}

func (bw *birWriter) writeBufferLength(buf *bytes.Buffer, birbuf *bytes.Buffer) {
	write(buf, int64(birbuf.Len()))
}

// FIXME: Write actual type
func (bw *birWriter) writeType(buf *bytes.Buffer, _ any) {
	write(buf, int32(-1))
}

func (bw *birWriter) writePosition(buf *bytes.Buffer, pos diagnostics.Location) {
	var sLine int32 = math.MaxInt32
	var eLine int32 = math.MaxInt32
	var sCol int32 = math.MaxInt32
	var eCol int32 = math.MaxInt32
	var sourceFileName string = ""

	if pos != nil {
		sLine = int32(pos.LineRange().StartLine().Line())
		eLine = int32(pos.LineRange().EndLine().Line())
		sCol = int32(pos.LineRange().StartLine().Offset())
		eCol = int32(pos.LineRange().EndLine().Offset())
		if (pos.LineRange().FileName()) != "" {
			sourceFileName = pos.LineRange().FileName()
		}
	}

	bw.writeStringCPEntry(buf, sourceFileName)
	write(buf, sLine)
	write(buf, sCol)
	write(buf, eLine)
	write(buf, eCol)
}
