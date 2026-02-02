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
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
)

const (
	BIR_MAGIC   = "\xba\x10\xc0\xde"
	BIR_VERSION = 75
)

type BIRWriter struct {
	cp *ConstantPool
}

func NewBIRWriter() *BIRWriter {
	return &BIRWriter{
		cp: NewConstantPool(),
	}
}

// Marshal returns the BIR encoding of pkg.
// This is a convenience function that creates a new BIRWriter and calls Serialize.
func Marshal(pkg *bir.BIRPackage) ([]byte, error) {
	return NewBIRWriter().serialize(pkg)
}

func (bw *BIRWriter) serialize(pkg *bir.BIRPackage) ([]byte, error) {
	// Reset constant pool for reuse
	bw.cp = NewConstantPool()

	birbuf := &bytes.Buffer{}

	// Write the package details
	pkgIDIdx := int32(-1)
	if pkg.PackageID != nil {
		pkgIDIdx = bw.cp.AddPackageCPEntry(pkg.PackageID)
	}
	if err := bw.writeInt32(birbuf, pkgIDIdx); err != nil {
		return nil, err
	}

	// Write import module declarations
	if err := bw.writeImportModuleDecls(birbuf, pkg); err != nil {
		return nil, err
	}

	// Write constants
	if err := bw.writeConstants(birbuf, pkg); err != nil {
		return nil, err
	}

	// Write global vars
	if err := bw.writeGlobalVars(birbuf, pkg); err != nil {
		return nil, err
	}

	// Write functions
	if err := bw.writeFunctions(birbuf, pkg); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	_, err := buf.Write([]byte(BIR_MAGIC))
	if err != nil {
		return nil, err
	}

	if err := bw.writeInt32(buf, int32(BIR_VERSION)); err != nil {
		return nil, err
	}

	cpBytes, err := bw.cp.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(cpBytes); err != nil {
		return nil, err
	}

	if _, err := buf.Write(birbuf.Bytes()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (bw *BIRWriter) writeImportModuleDecls(buf *bytes.Buffer, pkg *bir.BIRPackage) error {
	if err := bw.writeInt32(buf, int32(len(pkg.ImportModules))); err != nil {
		return err
	}
	for _, imp := range pkg.ImportModules {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(imp.PackageID.OrgName.Value())); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(imp.PackageID.PkgName.Value())); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(imp.PackageID.Name.Value())); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(imp.PackageID.Version.Value())); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeConstants(buf *bytes.Buffer, pkg *bir.BIRPackage) error {
	if err := bw.writeInt32(buf, int32(len(pkg.Constants))); err != nil {
		return err
	}

	for _, c := range pkg.Constants {
		if err := bw.writeConstant(buf, &c); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeConstant(buf *bytes.Buffer, constant *bir.BIRConstant) error {
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(constant.Name.Value())); err != nil {
		return err
	}

	if err := bw.writeInt64(buf, constant.Flags); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(constant.Origin)); err != nil {
		return err
	}

	err := bw.writeType(buf, constant.Type)
	if err != nil {
		return err
	}

	birbuf := &bytes.Buffer{}
	err = bw.writeType(birbuf, constant.ConstValue.Type)
	if err != nil {
		return err
	}
	err = bw.writeConstValue(birbuf, &constant.ConstValue)
	if err != nil {
		return err
	}

	if err := bw.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err = buf.Write(birbuf.Bytes())
	return err
}

func (bw *BIRWriter) writeType(buf *bytes.Buffer, t any) error {
	return bw.writeInt32(buf, -1)
}

func (bw *BIRWriter) writeConstValue(buf *bytes.Buffer, cv *bir.ConstValue) error {
	tag, err := bw.inferTag(cv.Value)
	if err != nil {
		return err
	}

	valIdx, err := bw.addValueToCP(tag, cv.Value)
	if err != nil {
		return err
	}
	return bw.writeInt32(buf, valIdx)
}

func (bw *BIRWriter) inferTag(value any) (model.TypeTags, error) {
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

func (bw *BIRWriter) addValueToCP(tag model.TypeTags, value any) (int32, error) {
	if cv, ok := value.(bir.ConstValue); ok {
		return bw.addValueToCP(tag, cv.Value)
	}

	switch tag {
	case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
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
			return 0, fmt.Errorf("expected integer for tag %v, got %T", tag, value)
		}
		return bw.cp.AddIntegerCPEntry(val), nil
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
			return 0, fmt.Errorf("expected byte for tag %v, got %T", tag, value)
		}
		return bw.cp.AddByteCPEntry(val), nil
	case model.TypeTags_FLOAT:
		var val float64
		switch v := value.(type) {
		case float64:
			val = v
		case float32:
			val = float64(v)
		default:
			return 0, fmt.Errorf("expected float for tag %v, got %T", tag, value)
		}
		return bw.cp.AddFloatCPEntry(val), nil
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
	case model.TypeTags_BOOLEAN:
		var val bool
		switch v := value.(type) {
		case bool:
			val = v
		default:
			return 0, fmt.Errorf("expected boolean for tag %v, got %T", tag, value)
		}
		return bw.cp.AddBooleanCPEntry(val), nil
	case model.TypeTags_NIL:
		// NIL values don't need to be added to the constant pool
		// Return -1 as the index for nil values
		return -1, nil
	default:
		fmt.Println("Unsupported tag:", tag)
		return 0, fmt.Errorf("unsupported tag for constant value: %v", tag)
	}
}

func (bw *BIRWriter) writeGlobalVars(buf *bytes.Buffer, pkg *bir.BIRPackage) error {
	if err := bw.writeInt32(buf, int32(len(pkg.GlobalVars))); err != nil {
		return err
	}

	for _, gv := range pkg.GlobalVars {
		if err := bw.writeUInt8(buf, uint8(gv.Kind)); err != nil {
			return err
		}

		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(gv.Name.Value())); err != nil {
			return err
		}

		if err := bw.writeInt64(buf, gv.Flags); err != nil {
			return err
		}

		if err := bw.writeUInt8(buf, uint8(gv.Origin)); err != nil {
			return err
		}

		bw.writeType(buf, gv.Type)
	}

	return nil
}

func (bw *BIRWriter) writeFunctions(buf *bytes.Buffer, pkg *bir.BIRPackage) error {
	if err := bw.writeInt32(buf, int32(len(pkg.Functions))); err != nil {
		return err
	}

	for _, fn := range pkg.Functions {
		if err := bw.writeFunction(buf, &fn); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeFunction(buf *bytes.Buffer, fn *bir.BIRFunction) error {
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(fn.Name.Value())); err != nil {
		return err
	}

	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(fn.OriginalName.Value())); err != nil {
		return err
	}

	if err := bw.writeInt64(buf, fn.Flags); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(fn.Origin)); err != nil {
		return err
	}

	if err := bw.writeInt32(buf, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, requiredParam := range fn.RequiredParams {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(requiredParam.Name.Value())); err != nil {
			return err
		}
		if err := bw.writeInt64(buf, requiredParam.Flags); err != nil {
			return err
		}
	}

	birbuf := &bytes.Buffer{}

	if err := bw.writeInt32(birbuf, int32(fn.ArgsCount)); err != nil {
		return err
	}
	if err := bw.writeBool(birbuf, fn.ReturnVariable != nil); err != nil {
		return err
	}
	if fn.ReturnVariable != nil {
		if err := bw.writeUInt8(birbuf, uint8(fn.ReturnVariable.Kind)); err != nil {
			return err
		}

		err := bw.writeType(birbuf, fn.ReturnVariable.Type)
		if err != nil {
			return err
		}
		if err := bw.writeInt32(birbuf, bw.cp.AddStringCPEntry(fn.ReturnVariable.Name.Value())); err != nil {
			return err
		}
	}

	if err := bw.writeInt32(birbuf, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	for _, localVar := range fn.LocalVars {
		if err := bw.writeLocalVar(birbuf, &localVar); err != nil {
			return err
		}
	}

	if err := bw.writeInt32(birbuf, int32(len(fn.BasicBlocks))); err != nil {
		return err
	}

	for _, bb := range fn.BasicBlocks {
		if err := bw.writeBasicBlock(birbuf, &bb); err != nil {
			return err
		}
	}

	if err := bw.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err := buf.Write(birbuf.Bytes())

	return err
}

func (bw *BIRWriter) writeLocalVar(buf *bytes.Buffer, localVar *bir.BIRVariableDcl) error {
	if err := bw.writeUInt8(buf, uint8(localVar.Kind)); err != nil {
		return err
	}

	if err := bw.writeType(buf, localVar.Type); err != nil {
		return err
	}
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(localVar.Name.Value())); err != nil {
		return err
	}

	if localVar.Kind == bir.VAR_KIND_ARG {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(localVar.MetaVarName)); err != nil {
			return err
		}
	}

	if localVar.Kind == bir.VAR_KIND_LOCAL {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(localVar.MetaVarName)); err != nil {
			return err
		}

		endBBId := ""
		if localVar.EndBB != nil {
			endBBId = localVar.EndBB.Id.Value()
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(endBBId)); err != nil {
			return err
		}

		startBBId := ""
		if localVar.StartBB != nil {
			startBBId = localVar.StartBB.Id.Value()
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(startBBId)); err != nil {
			return err
		}

		if err := bw.writeInt32(buf, int32(localVar.InsOffset)); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BIRWriter) writeBasicBlock(buf *bytes.Buffer, bb *bir.BIRBasicBlock) error {
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(bb.Id.Value())); err != nil {
		return err
	}
	// TODO: Adding the terminator instruction as well! Why?
	if err := bw.writeInt32(buf, int32(len(bb.Instructions))); err != nil {
		return err
	}

	for _, instr := range bb.Instructions {
		if err := bw.writeUInt8(buf, uint8(instr.GetKind())); err != nil {
			return err
		}
		if err := bw.writeInstruction(buf, instr); err != nil {
			return err
		}
	}

	if bb.Terminator == nil {
		return bw.writeUInt8(buf, 0)
	}
	if err := bw.writeUInt8(buf, uint8(bb.Terminator.GetKind())); err != nil {
		return err
	}

	return bw.writeTerminator(buf, bb.Terminator)
}

func (bw *BIRWriter) writeTerminator(buf *bytes.Buffer, term bir.BIRTerminator) error {
	switch term := term.(type) {
	case *bir.Goto:
		id := term.ThenBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(id)); err != nil {
			return err
		}
	case *bir.Branch:
		if err := bw.writeOperand(buf, term.Op); err != nil {
			return err
		}

		trueId := term.TrueBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(trueId)); err != nil {
			return err
		}
		falseId := term.FalseBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(falseId)); err != nil {
			return err
		}
	case *bir.Call:
		if err := bw.writeBool(buf, term.IsVirtual); err != nil {
			return err
		}
		pkgIdx := int32(-1)
		if term.CalleePkg != nil {
			pkgIdx = bw.cp.AddPackageCPEntry(term.CalleePkg)
		}
		if err := bw.writeInt32(buf, pkgIdx); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(term.Name.Value())); err != nil {
			return err
		}

		if err := bw.writeInt32(buf, int32(len(term.Args))); err != nil {
			return err
		}

		for _, arg := range term.Args {
			if err := bw.writeOperand(buf, &arg); err != nil {
				return err
			}
		}

		if term.LhsOp != nil {
			if err := bw.writeUInt8(buf, uint8(1)); err != nil {
				return err
			}
			if err := bw.writeOperand(buf, term.LhsOp); err != nil {
				return err
			}
		} else {
			if err := bw.writeUInt8(buf, uint8(0)); err != nil {
				return err
			}
		}

		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(term.ThenBB.Id.Value())); err != nil {
			return err
		}

	case *bir.Return:
		// Nothing to writ for return terminator
	default:
		return fmt.Errorf("unsupported terminator type: %T", term)
	}
	return nil
}

func (bw *BIRWriter) writeInstruction(buf *bytes.Buffer, instr bir.BIRInstruction) error {
	switch instr := instr.(type) {
	case *bir.Move:
		if err := bw.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.BinaryOp:
		if err := bw.writeOperand(buf, &instr.RhsOp1); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, &instr.RhsOp2); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.UnaryOp:
		if err := bw.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.ConstantLoad:
		if err := bw.writeType(buf, instr.Type); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, instr.LhsOp); err != nil {
			return err
		}

		isWrapped := false
		var tagValue any = instr.Value
		if cv, ok := instr.Value.(bir.ConstValue); ok {
			isWrapped = true
			tagValue = cv.Value
		}

		tag, err := bw.inferTag(tagValue)
		if err != nil {
			return err
		}

		valIdx, err := bw.addValueToCP(tag, tagValue)
		if err != nil {
			return err
		}
		if err := bw.writeBool(buf, isWrapped); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, valIdx); err != nil {
			return err
		}

	case *bir.FieldAccess:
		// TODO: MAP_LOAD and ARRAY_LOAD
		if err := bw.writeOperand(buf, instr.LhsOp); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, instr.KeyOp); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
	case *bir.NewArray:
		if err := bw.writeType(buf, instr.Type); err != nil {
			return err
		}

		if err := bw.writeOperand(buf, instr.LhsOp); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, instr.SizeOp); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BIRWriter) writeOperand(buf *bytes.Buffer, op *bir.BIROperand) error {
	if op == nil || op.VariableDcl == nil {
		// Should not happen based on current bir-gen, but let's be safe
		if err := bw.writeBool(buf, false); err != nil {
			return err
		}
		if err := bw.writeUInt8(buf, uint8(bir.VAR_KIND_TEMP)); err != nil {
			return err
		}
		if err := bw.writeUInt8(buf, uint8(bir.VAR_SCOPE_FUNCTION)); err != nil {
			return err
		}
		return bw.writeInt32(buf, bw.cp.AddStringCPEntry(""))
	}

	if op.VariableDcl.IgnoreVariable {
		if err := bw.writeBool(buf, true); err != nil {
			return err
		}

		return bw.writeType(buf, op.VariableDcl.Type)
	}

	if err := bw.writeBool(buf, false); err != nil {
		return err
	}
	if err := bw.writeUInt8(buf, uint8(op.VariableDcl.Kind)); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(op.VariableDcl.Scope)); err != nil {
		return err
	}

	nameIdx := bw.cp.AddStringCPEntry(op.VariableDcl.Name.Value())
	return bw.writeInt32(buf, nameIdx)
}

func (bw *BIRWriter) writeInt8(buf *bytes.Buffer, val int8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeUInt8(buf *bytes.Buffer, val uint8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeInt32(buf *bytes.Buffer, val int32) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeInt64(buf *bytes.Buffer, val int64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeFloat64(buf *bytes.Buffer, val float64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeBool(buf *bytes.Buffer, val bool) error {
	return binary.Write(buf, binary.BigEndian, val)
}
