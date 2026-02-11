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

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

type birReader struct {
	r   *bytes.Reader
	cp  []any
	ctx *context.CompilerContext
}

func Unmarshal(ctx *context.CompilerContext, data []byte) (*bir.BIRPackage, error) {
	reader := &birReader{
		r:   bytes.NewReader(data),
		ctx: ctx,
	}
	return reader.readPackage()
}

func (br *birReader) readPackage() (pkg *bir.BIRPackage, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("BIR deserializer failed: %v", r)
		}
	}()

	magic := make([]byte, 4)
	br.read(magic)

	if string(magic) != BIR_MAGIC {
		panic(fmt.Sprintf("invalid BIR magic: %x", magic))
	}

	var version int32
	br.read(&version)

	if version != BIR_VERSION {
		panic(fmt.Sprintf("unsupported BIR version: %d", version))
	}

	br.readConstantPool()

	var pkgIdx int32
	br.read(&pkgIdx)

	pkgID := br.getPackageFromCP(int(pkgIdx))
	imports := br.readImports()
	constants := br.readConstants()
	globalVars := br.readGlobalVars()
	functions := br.readFunctions()

	return &bir.BIRPackage{
		PackageID:     pkgID,
		ImportModules: imports,
		Constants:     constants,
		GlobalVars:    globalVars,
		Functions:     functions,
	}, nil
}

func (br *birReader) readConstantPool() {
	var cpSize int64
	br.read(&cpSize)

	br.cp = make([]any, cpSize)

	for i := 0; i < int(cpSize); i++ {
		var tag int8
		br.read(&tag)
		br.readConstantPoolEntry(tag, i)
	}
}

func (br *birReader) readConstantPoolEntry(tag int8, i int) {
	switch tag {
	case 0:
		br.cp[i] = nil
	case 1:
		var length int64
		br.read(&length)

		if length < 0 {
			br.cp[i] = (*string)(nil)
		} else {
			strBytes := make([]byte, length)
			br.read(strBytes)

			br.cp[i] = string(strBytes)
		}
	case 2:
		var orgIdx int32
		br.read(&orgIdx)

		var pkgNameIdx int32
		br.read(&pkgNameIdx)

		var moduleNameIdx int32
		br.read(&moduleNameIdx)

		var versionIdx int32
		br.read(&versionIdx)

		org := model.Name(br.getStringFromCP(int(orgIdx)))
		pkgName := model.Name(br.getStringFromCP(int(pkgNameIdx)))
		_ = br.getStringFromCP(int(moduleNameIdx))
		version := model.Name(br.getStringFromCP(int(versionIdx)))
		nameComps := model.CreateNameComps(pkgName)
		br.cp[i] = br.ctx.NewPackageID(org, nameComps, version)
	case 3:
		panic("shape not implemented")
	default:
		panic(fmt.Sprintf("unknown CP tag: %d", tag))
	}
}

func (br *birReader) getFromCP(index int) any {
	if index < 0 || index >= len(br.cp) {
		return nil
	}
	return br.cp[index]
}

func (br *birReader) getStringFromCP(index int) string {
	v := br.getFromCP(index)
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func (br *birReader) getPackageFromCP(index int) *model.PackageID {
	v := br.getFromCP(index)
	if pkg, ok := v.(*model.PackageID); ok {
		return pkg
	}
	return nil
}

// FIXME: Read actual type
func (br *birReader) getTypeFromCP(_ int) ast.BType {
	return nil
}

func (br *birReader) readImports() []bir.BIRImportModule {
	count := br.readLength()
	imports := make([]bir.BIRImportModule, count)
	for i := 0; i < int(count); i++ {
		org := br.readStringCPEntry()
		pkgName := br.readStringCPEntry()
		_ = br.readStringCPEntry()
		version := br.readStringCPEntry()

		nameComps := model.CreateNameComps(pkgName)
		imports[i] = bir.BIRImportModule{
			PackageID: br.ctx.NewPackageID(org, nameComps, version),
		}
	}
	return imports
}

func (br *birReader) readConstants() []bir.BIRConstant {
	count := br.readLength()
	constants := make([]bir.BIRConstant, count)
	for i := 0; i < int(count); i++ {
		name := br.readStringCPEntry()
		flags := br.readFlags()
		origin := br.readOrigin()
		pos := br.readPosition()

		var typeIdx int32
		br.read(&typeIdx)

		// TODO: Implement Types
		// t := br.getTypeFromCP(int(typeIdx))

		constant := bir.BIRConstant{
			BIRDocumentableNodeBase: bir.BIRDocumentableNodeBase{
				BIRNodeBase: bir.BIRNodeBase{
					Pos: pos,
				},
			},
			Name:   name,
			Flags:  flags,
			Origin: origin,
			Type:   nil,
		}

		br.readLength()

		var cTypeIdx int32
		br.read(&cTypeIdx)

		cv := br.getTypeFromCP(int(cTypeIdx))
		value := br.readConstValue(cv)

		constant.ConstValue = bir.ConstValue{
			Type:  cv,
			Value: value,
		}

		constants[i] = constant
	}
	return constants
}

func (br *birReader) readGlobalVars() []bir.BIRGlobalVariableDcl {
	count := br.readLength()
	variables := make([]bir.BIRGlobalVariableDcl, count)
	for i := 0; i < int(count); i++ {
		pos := br.readPosition()
		kind := br.readKind()
		name := br.readStringCPEntry()
		flags := br.readFlags()
		origin := br.readOrigin()

		var typeIdx int32
		br.read(&typeIdx)

		// TODO: Implement Types
		// t := br.getTypeFromCP(int(typeIdx))

		variables[i] = bir.BIRGlobalVariableDcl{
			BIRVariableDcl: bir.BIRVariableDcl{
				BIRDocumentableNodeBase: bir.BIRDocumentableNodeBase{
					BIRNodeBase: bir.BIRNodeBase{
						Pos: pos,
					},
				},
				Kind: kind,
				Name: name,
				Type: nil,
			},
			Flags:  flags,
			Origin: origin,
		}
	}
	return variables
}

func (br *birReader) readFunctions() []bir.BIRFunction {
	count := br.readLength()
	functions := make([]bir.BIRFunction, count)
	for i := 0; i < int(count); i++ {
		fn := br.readFunction()
		functions[i] = *fn
	}
	return functions
}

func (br *birReader) readFunction() *bir.BIRFunction {
	pos := br.readPosition()
	name := br.readStringCPEntry()
	originalName := br.readStringCPEntry()
	flag := br.readFlags()
	origin := br.readOrigin()
	requiredParamsCount := br.readLength()

	requiredParams := make([]bir.BIRParameter, requiredParamsCount)
	for j := 0; j < int(requiredParamsCount); j++ {
		paramName := br.readStringCPEntry()
		paramFlags := br.readFlags()

		requiredParams[j] = bir.BIRParameter{
			Name:  paramName,
			Flags: paramFlags,
		}
	}

	_ = br.readLength() // Unused?

	argsCount := br.readLength()

	varMap := make(map[string]*bir.BIRVariableDcl)
	bbMap := make(map[string]*bir.BIRBasicBlock)

	var hasReturnVar bool
	br.read(&hasReturnVar)

	var returnVar *bir.BIRVariableDcl
	if hasReturnVar {
		returnVarKind := br.readKind()
		var returnVarTypeIdx int32
		br.read(&returnVarTypeIdx)

		// TODO: Implement Types
		// returnVarType := br.getTypeFromCP(int(returnVarTypeIdx))

		returnVarName := br.readStringCPEntry()

		returnVar = &bir.BIRVariableDcl{
			Kind: returnVarKind,
			Name: returnVarName,
			Type: nil,
		}
		varMap[returnVarName.Value()] = returnVar
	}

	localVarCount := br.readLength()
	localVars := make([]bir.BIRVariableDcl, localVarCount)
	for j := 0; j < int(localVarCount); j++ {
		localVar := br.readLocalVar(varMap)
		localVars[j] = *localVar
	}

	basicBlockCount := br.readLength()
	basicBlocks := make([]bir.BIRBasicBlock, basicBlockCount)
	for j := 0; j < int(basicBlockCount); j++ {
		block := br.readBasicBlock(varMap)
		basicBlocks[j] = *block
		bbMap[block.Id.Value()] = &basicBlocks[j]
	}

	for j := range basicBlocks {
		bb := &basicBlocks[j]
		if bb.Terminator != nil {
			switch t := bb.Terminator.(type) {
			case *bir.Goto:
				if target, ok := bbMap[t.ThenBB.Id.Value()]; ok {
					t.ThenBB = target
				}
			case *bir.Branch:
				if target, ok := bbMap[t.TrueBB.Id.Value()]; ok {
					t.TrueBB = target
				}
				if target, ok := bbMap[t.FalseBB.Id.Value()]; ok {
					t.FalseBB = target
				}
			case *bir.Call:
				if target, ok := bbMap[t.ThenBB.Id.Value()]; ok {
					t.ThenBB = target
				}
			}
		}
	}

	for j := range localVars {
		lv := &localVars[j]
		if lv.StartBB != nil {
			if target, ok := bbMap[lv.StartBB.Id.Value()]; ok {
				lv.StartBB = target
			}
		}
		if lv.EndBB != nil {
			if target, ok := bbMap[lv.EndBB.Id.Value()]; ok {
				lv.EndBB = target
			}
		}
	}

	return &bir.BIRFunction{
		BIRDocumentableNodeBase: bir.BIRDocumentableNodeBase{
			BIRNodeBase: bir.BIRNodeBase{
				Pos: pos,
			},
		},
		Name:           name,
		OriginalName:   originalName,
		Flags:          flag,
		Origin:         origin,
		RequiredParams: requiredParams,
		ArgsCount:      int(argsCount),
		ReturnVariable: returnVar,
		LocalVars:      localVars,
		BasicBlocks:    basicBlocks,
	}
}

func (br *birReader) readLocalVar(varMap map[string]*bir.BIRVariableDcl) *bir.BIRVariableDcl {
	kind := br.readKind()
	var typeIdx int32
	br.read(&typeIdx)

	// TODO: Implement Types
	// t := br.getTypeFromCP(int(typeIdx))

	name := br.readStringCPEntry()

	localVar := &bir.BIRVariableDcl{
		Kind: kind,
		Name: name,
		Type: nil,
	}

	if kind == bir.VAR_KIND_ARG {
		metaVarName := br.readStringCPEntry()
		localVar.MetaVarName = metaVarName.Value()
	} else if kind == bir.VAR_KIND_LOCAL {
		metaVarName := br.readStringCPEntry()
		localVar.MetaVarName = metaVarName.Value()

		endBBId := br.readStringCPEntry()
		localVar.EndBB = &bir.BIRBasicBlock{Id: endBBId}

		startBBId := br.readStringCPEntry()
		localVar.StartBB = &bir.BIRBasicBlock{Id: startBBId}

		insOffset := br.readLength()
		localVar.InsOffset = int(insOffset)
	}
	varMap[name.Value()] = localVar
	return localVar
}

func (br *birReader) readBasicBlock(varMap map[string]*bir.BIRVariableDcl) *bir.BIRBasicBlock {
	id := br.readStringCPEntry()
	instructionCount := br.readLength()

	instructions := make([]bir.BIRInstruction, instructionCount)
	for k := 0; k < int(instructionCount); k++ {
		ins := br.readInstruction(varMap)
		instructions[k] = ins
	}

	term := br.readTerminator(varMap)

	return &bir.BIRBasicBlock{
		Id:           id,
		Instructions: instructions,
		Terminator:   term,
	}
}

func (br *birReader) readInstruction(varMap map[string]*bir.BIRVariableDcl) bir.BIRInstruction {
	instructionKind := br.readInstructionKind()

	switch instructionKind {
	case bir.INSTRUCTION_KIND_MOVE:
		rhsOp := br.readOperand(varMap)
		lhsOp := br.readOperand(varMap)
		return &bir.Move{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			RhsOp: rhsOp,
		}

	case bir.INSTRUCTION_KIND_ADD, bir.INSTRUCTION_KIND_SUB, bir.INSTRUCTION_KIND_MUL,
		bir.INSTRUCTION_KIND_DIV, bir.INSTRUCTION_KIND_MOD, bir.INSTRUCTION_KIND_EQUAL,
		bir.INSTRUCTION_KIND_NOT_EQUAL, bir.INSTRUCTION_KIND_GREATER_THAN,
		bir.INSTRUCTION_KIND_GREATER_EQUAL, bir.INSTRUCTION_KIND_LESS_THAN,
		bir.INSTRUCTION_KIND_LESS_EQUAL, bir.INSTRUCTION_KIND_AND, bir.INSTRUCTION_KIND_OR,
		bir.INSTRUCTION_KIND_REF_EQUAL, bir.INSTRUCTION_KIND_REF_NOT_EQUAL,
		bir.INSTRUCTION_KIND_CLOSED_RANGE, bir.INSTRUCTION_KIND_HALF_OPEN_RANGE,
		bir.INSTRUCTION_KIND_ANNOT_ACCESS, bir.INSTRUCTION_KIND_BITWISE_AND,
		bir.INSTRUCTION_KIND_BITWISE_OR, bir.INSTRUCTION_KIND_BITWISE_XOR,
		bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT, bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT,
		bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
		rhsOp1 := br.readOperand(varMap)
		rhsOp2 := br.readOperand(varMap)
		lhsOp := br.readOperand(varMap)
		return &bir.BinaryOp{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Kind:   instructionKind,
			RhsOp1: *rhsOp1,
			RhsOp2: *rhsOp2,
		}

	case bir.INSTRUCTION_KIND_TYPEOF, bir.INSTRUCTION_KIND_NOT, bir.INSTRUCTION_KIND_NEGATE:
		rhsOp := br.readOperand(varMap)
		lhsOp := br.readOperand(varMap)
		return &bir.UnaryOp{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Kind:  instructionKind,
			RhsOp: rhsOp,
		}

	case bir.INSTRUCTION_KIND_CONST_LOAD:
		var constLoadTypeIdx int32
		br.read(&constLoadTypeIdx)

		constLoadType := br.getTypeFromCP(int(constLoadTypeIdx))

		lhsOp := br.readOperand(varMap)

		var isWrapped bool
		br.read(&isWrapped)

		var tagByte int8
		br.read(&tagByte)

		tag := model.TypeTags(tagByte)
		value := br.readConstValueByTag(tag)

		if isWrapped {
			value = bir.ConstValue{
				Type:  nil,
				Value: value,
			}
		}

		return &bir.ConstantLoad{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Type:  constLoadType,
			Value: value,
		}

	case bir.INSTRUCTION_KIND_MAP_STORE, bir.INSTRUCTION_KIND_MAP_LOAD,
		bir.INSTRUCTION_KIND_ARRAY_STORE, bir.INSTRUCTION_KIND_ARRAY_LOAD:
		lhsOp := br.readOperand(varMap)
		keyOp := br.readOperand(varMap)
		rhsOp := br.readOperand(varMap)
		return &bir.FieldAccess{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Kind:  instructionKind,
			KeyOp: keyOp,
			RhsOp: rhsOp,
		}

	case bir.INSTRUCTION_KIND_NEW_ARRAY:
		var typeIdx int32
		br.read(&typeIdx)

		t := br.getTypeFromCP(int(typeIdx))

		lhsOp := br.readOperand(varMap)
		sizeOp := br.readOperand(varMap)
		return &bir.NewArray{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Type:   t,
			SizeOp: sizeOp,
		}

	case bir.INSTRUCTION_KIND_TYPE_CAST:
		lhsOp := br.readOperand(varMap)
		rhsOp := br.readOperand(varMap)

		var typeIdx int32
		br.read(&typeIdx)

		return &bir.TypeCast{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			RhsOp: rhsOp,
			Type:  nil,
		}

	default:
		panic(fmt.Sprintf("unsupported instruction kind: %d", instructionKind))
	}
}

func (br *birReader) readTerminator(varMap map[string]*bir.BIRVariableDcl) bir.BIRTerminator {
	var terminatorKind uint8
	br.read(&terminatorKind)

	if terminatorKind == 0 {
		return nil
	}

	termInstructionKind := bir.InstructionKind(terminatorKind)

	switch termInstructionKind {
	case bir.INSTRUCTION_KIND_RETURN:
		return &bir.Return{}

	case bir.INSTRUCTION_KIND_GOTO:
		id := br.readStringCPEntry()
		return &bir.Goto{
			BIRTerminatorBase: bir.BIRTerminatorBase{
				ThenBB: &bir.BIRBasicBlock{
					Id: id,
				},
			},
		}

	case bir.INSTRUCTION_KIND_BRANCH:
		op := br.readOperand(varMap)
		trueBBId := br.readStringCPEntry()
		falseBBId := br.readStringCPEntry()

		return &bir.Branch{
			Op: op,
			TrueBB: &bir.BIRBasicBlock{
				Id: trueBBId,
			},
			FalseBB: &bir.BIRBasicBlock{
				Id: falseBBId,
			},
		}

	case bir.INSTRUCTION_KIND_CALL:
		var isVirtual bool
		br.read(&isVirtual)

		pkg := br.readPackageCPEntry()
		name := br.readStringCPEntry()
		argsCount := br.readLength()

		args := make([]bir.BIROperand, argsCount)
		for k := 0; k < int(argsCount); k++ {
			arg := br.readOperand(varMap)
			args[k] = *arg
		}

		var lshOpExists bool
		br.read(&lshOpExists)

		var lhsOp *bir.BIROperand
		if lshOpExists {
			lhsOp = br.readOperand(varMap)
		}

		thenBBId := br.readStringCPEntry()

		return &bir.Call{
			Kind:      termInstructionKind,
			IsVirtual: isVirtual,
			CalleePkg: pkg,
			Name:      name,
			Args:      args,
			BIRTerminatorBase: bir.BIRTerminatorBase{
				ThenBB: &bir.BIRBasicBlock{
					Id: thenBBId,
				},
				BIRInstructionBase: bir.BIRInstructionBase{
					LhsOp: lhsOp,
				},
			},
		}

	default:
		panic(fmt.Sprintf("unsupported terminator kind: %d", termInstructionKind))
	}
}

func (br *birReader) readOperand(varMap map[string]*bir.BIRVariableDcl) *bir.BIROperand {
	var ignoreVariable bool
	br.read(&ignoreVariable)

	if ignoreVariable {
		var varTypeIdx int32
		br.read(&varTypeIdx)

		// TODO: Implement Types
		// varType := br.getTypeFromCP(int(varTypeIdx))
		return &bir.BIROperand{
			VariableDcl: &bir.BIRVariableDcl{
				Type: nil,
			},
		}
	}

	varKind := br.readKind()
	scope := br.readScope()
	name := br.readStringCPEntry()

	varDcl, ok := varMap[name.Value()]
	if !ok {
		varDcl = &bir.BIRVariableDcl{
			Kind:  varKind,
			Scope: scope,
			Name:  name,
		}
		varMap[name.Value()] = varDcl
	}

	return &bir.BIROperand{
		VariableDcl: varDcl,
	}
}

func (br *birReader) readConstValue(valueType ast.BType) any {
	var tagByte int8
	br.read(&tagByte)

	tag := model.TypeTags(tagByte)
	return br.readConstValueByTag(tag)
}

func (br *birReader) readConstValueByTag(tag model.TypeTags) any {
	switch tag {
	case model.TypeTags_INT,
		model.TypeTags_SIGNED32_INT,
		model.TypeTags_SIGNED16_INT,
		model.TypeTags_SIGNED8_INT,
		model.TypeTags_UNSIGNED32_INT,
		model.TypeTags_UNSIGNED16_INT,
		model.TypeTags_UNSIGNED8_INT:
		var val int64
		br.read(&val)
		return val

	case model.TypeTags_BYTE:
		var val byte
		br.read(&val)
		return val

	case model.TypeTags_FLOAT:
		var val float64
		br.read(&val)
		return val

	case model.TypeTags_BOOLEAN:
		var val bool
		br.read(&val)
		return val

	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		var idx int32
		br.read(&idx)
		return br.getStringFromCP(int(idx))

	case model.TypeTags_NIL:
		var idx int32
		br.read(&idx)

		if idx != -1 {
			panic(fmt.Sprintf("expected -1 for nil value, got %d", idx))
		}
		return nil

	default:
		var idx int32
		br.read(&idx)

		if idx == -1 {
			return nil
		}
		return br.getFromCP(int(idx))
	}
}

func (br *birReader) read(v any) {
	err := binary.Read(br.r, binary.BigEndian, v)
	if err != nil {
		panic(fmt.Sprintf("binary read failed: %v", err))
	}
}

func (br *birReader) readKind() bir.VarKind {
	var val uint8
	br.read(&val)
	return bir.VarKind(val)
}

func (br *birReader) readFlags() int64 {
	var val int64
	br.read(&val)
	return val
}

func (br *birReader) readOrigin() model.SymbolOrigin {
	var val uint8
	br.read(&val)
	return model.SymbolOrigin(val)
}

func (br *birReader) readStringCPEntry() model.Name {
	var idx int32
	br.read(&idx)
	return model.Name(br.getStringFromCP(int(idx)))
}

func (br *birReader) readLength() int64 {
	var val int64
	br.read(&val)
	return val
}

func (br *birReader) readInstructionKind() bir.InstructionKind {
	var val uint8
	br.read(&val)
	return bir.InstructionKind(val)
}

func (br *birReader) readScope() bir.VarScope {
	var val uint8
	br.read(&val)
	return bir.VarScope(val)
}

func (br *birReader) readPackageCPEntry() *model.PackageID {
	var idx int32
	br.read(&idx)
	return br.getPackageFromCP(int(idx))
}

func (br *birReader) readPosition() diagnostics.Location {
	var sourceFileIdx int32
	br.read(&sourceFileIdx)

	sourceFileName := br.getStringFromCP(int(sourceFileIdx))

	var sLine int32
	br.read(&sLine)
	var sCol int32
	br.read(&sCol)
	var eLine int32
	br.read(&eLine)
	var eCol int32
	br.read(&eCol)

	return diagnostics.NewBLangDiagnosticLocation(sourceFileName, int(sLine), int(eLine), int(sCol), int(eCol), 0, 0)
}
