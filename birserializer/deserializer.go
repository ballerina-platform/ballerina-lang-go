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

type BIRReader struct {
	r   *bytes.Reader
	cp  []any
	ctx *context.CompilerContext
}

func Unmarshal(ctx *context.CompilerContext, data []byte) (*bir.BIRPackage, error) {
	reader := &BIRReader{
		r:   bytes.NewReader(data),
		ctx: ctx,
	}

	return reader.readPackage()
}

func (br *BIRReader) readPackage() (*bir.BIRPackage, error) {
	var errMsg string
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				errMsg = msg
			} else {
				errMsg = fmt.Sprintf("%v", r)
			}
		}
	}()

	magic := make([]byte, 4)
	_, err := br.r.Read(magic)
	if err != nil {
		panic(fmt.Sprintf("reading BIR magic: %v", err))
	}

	if string(magic) != BIR_MAGIC {
		panic(fmt.Sprintf("invalid BIR magic: %x", magic))
	}

	var version int32
	if err := br.read(&version); err != nil {
		panic(fmt.Sprintf("reading BIR version: %v", err))
	}

	if version != BIR_VERSION {
		panic(fmt.Sprintf("unsupported BIR version: %d", version))
	}

	br.readConstantPool()

	var pkgCPIndex int32
	if err := br.read(&pkgCPIndex); err != nil {
		panic(fmt.Sprintf("reading package CP index: %v", err))
	}

	pkgID := br.getPackageFromCP(int(pkgCPIndex))

	imports := br.readImports()
	constants := br.readConstants()
	globalVars := br.readGlobalVars()
	functions := br.readFunctions()

	if errMsg != "" {
		return nil, fmt.Errorf("BIR deserializer failed due to %s", errMsg)
	}

	return &bir.BIRPackage{
		PackageID:     pkgID,
		ImportModules: imports,
		Constants:     constants,
		GlobalVars:    globalVars,
		Functions:     functions,
	}, nil
}

func (br *BIRReader) readConstantPool() {
	var cpSize int64
	if err := br.read(&cpSize); err != nil {
		panic(fmt.Sprintf("reading constant pool size: %v", err))
	}

	br.cp = make([]any, cpSize)

	for i := 0; i < int(cpSize); i++ {
		var tag int8
		if err := br.read(&tag); err != nil {
			panic(fmt.Sprintf("reading CP entry %d tag: %v", i, err))
		}

		br.readConstantPoolEntry(tag, i)
	}
}

func (br *BIRReader) readConstantPoolEntry(tag int8, i int) {
	switch tag {
	case 0: // NULL
		br.cp[i] = nil
	case 1: // INTEGER
		var value int64
		if err := br.read(&value); err != nil {
			panic(fmt.Sprintf("reading CP entry %d integer: %v", i, err))
		}
		br.cp[i] = value
	case 2: // FLOAT
		var value float64
		if err := br.read(&value); err != nil {
			panic(fmt.Sprintf("reading CP entry %d float: %v", i, err))
		}
		br.cp[i] = value
	case 3: // BOOLEAN
		var b uint8
		if err := br.read(&b); err != nil {
			panic(fmt.Sprintf("reading CP entry %d boolean: %v", i, err))
		}
		br.cp[i] = b != 0
	case 4: // STRING
		var length int64
		if err := br.read(&length); err != nil {
			panic(fmt.Sprintf("reading CP entry %d string length: %v", i, err))
		}
		if length < 0 {
			br.cp[i] = (*string)(nil)
		} else {
			strBytes := make([]byte, length)
			if _, err := br.r.Read(strBytes); err != nil {
				panic(fmt.Sprintf("reading CP entry %d string bytes: %v", i, err))
			}
			br.cp[i] = string(strBytes)
		}
	case 5: // PACKAGE
		var orgIdx int32
		if err := br.read(&orgIdx); err != nil {
			panic(fmt.Sprintf("reading CP entry %d package org index: %v", i, err))
		}
		var pkgNameIdx int32
		if err := br.read(&pkgNameIdx); err != nil {
			panic(fmt.Sprintf("reading CP entry %d package name index: %v", i, err))
		}
		var moduleNameIdx int32
		if err := br.read(&moduleNameIdx); err != nil {
			panic(fmt.Sprintf("reading CP entry %d module name index: %v", i, err))
		}
		var versionIdx int32
		if err := br.read(&versionIdx); err != nil {
			panic(fmt.Sprintf("reading CP entry %d version index: %v", i, err))
		}
		org := model.Name(br.getStringFromCP(int(orgIdx)))
		pkgName := model.Name(br.getStringFromCP(int(pkgNameIdx)))
		_ = br.getStringFromCP(int(moduleNameIdx)) // moduleName - not used, pkgName contains full path
		version := model.Name(br.getStringFromCP(int(versionIdx)))
		nameComps := model.CreateNameComps(pkgName)
		br.cp[i] = br.ctx.NewPackageID(org, nameComps, version)
	case 6: // BYTE
		var value uint8
		if err := br.read(&value); err != nil {
			panic(fmt.Sprintf("reading CP entry %d byte: %v", i, err))
		}
		br.cp[i] = value
	case 7: // SHAPE
		panic("shape not implemented")
	default:
		panic(fmt.Sprintf("unknown CP tag: %d", tag))
	}
}

// getFromCP safely retrieves a value from the constant pool at the given index.
// Returns nil if the index is out of bounds.
func (r *BIRReader) getFromCP(index int) any {
	if index < 0 || index >= len(r.cp) {
		return nil
	}
	return r.cp[index]
}

func (r *BIRReader) getStringFromCP(index int) string {
	v := r.getFromCP(index)
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func (r *BIRReader) getPackageFromCP(index int) *model.PackageID {
	v := r.getFromCP(index)
	if pkg, ok := v.(*model.PackageID); ok {
		return pkg
	}
	return nil
}

func (r *BIRReader) getTypeFromCP(index int) ast.BType {
	if index == -1 {
		return nil
	}
	v := r.getFromCP(index)
	if t, ok := v.(ast.BType); ok {
		return t
	}
	return nil
}

func (r *BIRReader) getIntegerFromCP(index int) any {
	v := r.getFromCP(index)
	if v == nil {
		return int64(0)
	}
	if val, ok := v.(int64); ok {
		return val
	}
	return v
}

func (r *BIRReader) getByteFromCP(index int) uint8 {
	v := r.getFromCP(index)
	if val, ok := v.(uint8); ok {
		return val
	}
	return 0
}

func (r *BIRReader) getFloatFromCP(index int) float64 {
	v := r.getFromCP(index)
	if val, ok := v.(float64); ok {
		return val
	}
	return 0
}

func (r *BIRReader) getBooleanFromCP(index int) bool {
	v := r.getFromCP(index)
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}

func (br *BIRReader) readImports() []bir.BIRImportModule {
	count := br.readLength()
	imports := make([]bir.BIRImportModule, count)
	for i := 0; i < int(count); i++ {
		org := br.readStringCPEntry()
		pkgName := br.readStringCPEntry()
		_ = br.readStringCPEntry() // moduleName - not used, pkgName contains full path
		version := br.readStringCPEntry()

		nameComps := model.CreateNameComps(pkgName)
		imports[i] = bir.BIRImportModule{
			PackageID: br.ctx.NewPackageID(org, nameComps, version),
		}
	}

	return imports
}

func (br *BIRReader) readConstants() []bir.BIRConstant {
	count := br.readLength()
	constants := make([]bir.BIRConstant, count)
	for i := 0; i < int(count); i++ {
		name := br.readStringCPEntry()
		flags := br.readFlags()
		origin := br.readOrigin()
		pos := br.readPosition()

		constant := bir.BIRConstant{
			BIRDocumentableNodeBase: bir.BIRDocumentableNodeBase{
				BIRNodeBase: bir.BIRNodeBase{
					Pos: pos,
				},
			},
			Name:   name,
			Flags:  flags,
			Origin: origin,
		}

		var typeIdx int32
		if err := br.read(&typeIdx); err != nil {
			panic(fmt.Sprintf("reading constant %d type index: %v", i, err))
		}

		t := br.getTypeFromCP(int(typeIdx))
		constant.Type = t

		br.readLength()

		var cTypeIdx int32
		if err := br.read(&cTypeIdx); err != nil {
			panic(fmt.Sprintf("reading constant %d value type index: %v", i, err))
		}

		cv := br.getTypeFromCP(int(cTypeIdx))
		var value any
		// Type is nil, read value from CP directly
		var valueIdx int32
		if err := br.read(&valueIdx); err != nil {
			panic(fmt.Sprintf("reading constant %d value index: %v", i, err))
		}
		value = br.cp[int(valueIdx)]

		constant.ConstValue = bir.ConstValue{
			Type:  cv,
			Value: value,
		}

		constants[i] = constant
	}

	return constants
}

func (br *BIRReader) readGlobalVars() []bir.BIRGlobalVariableDcl {
	count := br.readLength()
	variables := make([]bir.BIRGlobalVariableDcl, count)
	for i := 0; i < int(count); i++ {
		pos := br.readPosition()
		kind := br.readKind()
		name := br.readStringCPEntry()
		flags := br.readFlags()
		origin := br.readOrigin()

		var typeIdx int32
		if err := br.read(&typeIdx); err != nil {
			panic(fmt.Sprintf("reading global var %d type index: %v", i, err))
		}

		t := br.getTypeFromCP(int(typeIdx))

		variables[i] = bir.BIRGlobalVariableDcl{
			BIRVariableDcl: bir.BIRVariableDcl{
				BIRDocumentableNodeBase: bir.BIRDocumentableNodeBase{
					BIRNodeBase: bir.BIRNodeBase{
						Pos: pos,
					},
				},
				Kind: kind,
				Name: name,
				Type: t,
			},
			Flags:  flags,
			Origin: origin,
		}
	}

	return variables
}

func (br *BIRReader) readFunctions() []bir.BIRFunction {
	count := br.readLength()
	functions := make([]bir.BIRFunction, count)
	for i := 0; i < int(count); i++ {
		fn := br.readFunction()
		functions[i] = *fn
	}

	return functions
}

func (br *BIRReader) readFunction() *bir.BIRFunction {
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

	var length int64
	if err := br.read(&length); err != nil { // length, unused?
		panic(fmt.Sprintf("reading function length: %v", err))
	}

	argsCount := br.readLength()

	// Create local maps for variable and basic block lookups
	varMap := make(map[string]*bir.BIRVariableDcl)
	bbMap := make(map[string]*bir.BIRBasicBlock)

	var hasReturnVar bool
	if err := br.read(&hasReturnVar); err != nil {
		panic(fmt.Sprintf("reading function has return var: %v", err))
	}

	var returnVar *bir.BIRVariableDcl
	if hasReturnVar {
		returnVarKind := br.readKind()
		var returnVarTypeIdx int32
		if err := br.read(&returnVarTypeIdx); err != nil {
			panic(fmt.Sprintf("reading return var type index: %v", err))
		}
		returnVarType := br.getTypeFromCP(int(returnVarTypeIdx))

		returnVarName := br.readStringCPEntry()

		returnVar = &bir.BIRVariableDcl{
			Kind: returnVarKind,
			Name: returnVarName,
			Type: returnVarType,
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

func (br *BIRReader) readLocalVar(varMap map[string]*bir.BIRVariableDcl) *bir.BIRVariableDcl {
	kind := br.readKind()
	var typeIdx int32
	if err := br.read(&typeIdx); err != nil {
		panic(fmt.Sprintf("reading local var type index: %v", err))
	}
	t := br.getTypeFromCP(int(typeIdx))

	name := br.readStringCPEntry()

	localVar := &bir.BIRVariableDcl{
		Kind: kind,
		Name: name,
		Type: t,
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

func (br *BIRReader) readBasicBlock(varMap map[string]*bir.BIRVariableDcl) *bir.BIRBasicBlock {
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

func (br *BIRReader) readInstruction(varMap map[string]*bir.BIRVariableDcl) bir.BIRInstruction {
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
	case bir.INSTRUCTION_KIND_ADD, bir.INSTRUCTION_KIND_SUB, bir.INSTRUCTION_KIND_MUL, bir.INSTRUCTION_KIND_DIV, bir.INSTRUCTION_KIND_MOD, bir.INSTRUCTION_KIND_EQUAL, bir.INSTRUCTION_KIND_NOT_EQUAL, bir.INSTRUCTION_KIND_GREATER_THAN, bir.INSTRUCTION_KIND_GREATER_EQUAL, bir.INSTRUCTION_KIND_LESS_THAN, bir.INSTRUCTION_KIND_LESS_EQUAL, bir.INSTRUCTION_KIND_AND, bir.INSTRUCTION_KIND_OR, bir.INSTRUCTION_KIND_REF_EQUAL, bir.INSTRUCTION_KIND_REF_NOT_EQUAL, bir.INSTRUCTION_KIND_CLOSED_RANGE, bir.INSTRUCTION_KIND_HALF_OPEN_RANGE, bir.INSTRUCTION_KIND_ANNOT_ACCESS, bir.INSTRUCTION_KIND_BITWISE_AND, bir.INSTRUCTION_KIND_BITWISE_OR, bir.INSTRUCTION_KIND_BITWISE_XOR, bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT, bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT, bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
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
		if err := br.read(&constLoadTypeIdx); err != nil {
			panic(fmt.Sprintf("reading const load type index: %v", err))
		}
		constLoadType := br.getTypeFromCP(int(constLoadTypeIdx))

		lhsOp := br.readOperand(varMap)

		var isWrapped bool
		if err := br.read(&isWrapped); err != nil {
			panic(fmt.Sprintf("reading const load is wrapped: %v", err))
		}

		var valueIdx int32
		if err := br.read(&valueIdx); err != nil {
			panic(fmt.Sprintf("reading const load value index: %v", err))
		}

		var value any
		// Check if this is a NIL value (index -1)
		if valueIdx == -1 {
			value = nil
		} else {
			// Type info missing, read from CP and infer
			value = br.getIntegerFromCP(int(valueIdx))
		}

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
	case bir.INSTRUCTION_KIND_MAP_STORE, bir.INSTRUCTION_KIND_MAP_LOAD, bir.INSTRUCTION_KIND_ARRAY_STORE, bir.INSTRUCTION_KIND_ARRAY_LOAD:
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
		if err := br.read(&typeIdx); err != nil {
			panic(fmt.Sprintf("reading new array type index: %v", err))
		}
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
	default:
		panic(fmt.Sprintf("unsupported instruction kind: %d", instructionKind))
	}
}

func (br *BIRReader) readTerminator(varMap map[string]*bir.BIRVariableDcl) bir.BIRTerminator {
	var terminatorKind uint8
	if err := br.read(&terminatorKind); err != nil {
		panic(fmt.Sprintf("reading terminator kind: %v", err))
	}

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
		if err := br.read(&isVirtual); err != nil {
			panic(fmt.Sprintf("reading call is virtual: %v", err))
		}

		pkg := br.readPackageCPEntry()
		name := br.readStringCPEntry()
		argsCount := br.readLength()

		args := make([]bir.BIROperand, argsCount)
		for k := 0; k < int(argsCount); k++ {
			arg := br.readOperand(varMap)
			args[k] = *arg
		}

		var lshOpExists bool
		if err := br.read(&lshOpExists); err != nil {
			panic(fmt.Sprintf("reading call lhs op exists: %v", err))
		}

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

func (br *BIRReader) readOperand(varMap map[string]*bir.BIRVariableDcl) *bir.BIROperand {
	var ignoreVariable bool
	if err := br.read(&ignoreVariable); err != nil {
		panic(fmt.Sprintf("reading operand ignore variable: %v", err))
	}

	if ignoreVariable {
		var varTypeIdx int32
		if err := br.read(&varTypeIdx); err != nil {
			panic(fmt.Sprintf("reading operand type index: %v", err))
		}
		varType := br.getTypeFromCP(int(varTypeIdx))
		return &bir.BIROperand{
			VariableDcl: &bir.BIRVariableDcl{
				Type: varType,
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
	}

	return &bir.BIROperand{
		VariableDcl: varDcl,
	}
}

func (br *BIRReader) read(v any) error {
	return binary.Read(br.r, binary.BigEndian, v)
}

func (br *BIRReader) readKind() bir.VarKind {
	var val uint8
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading kind: %v", err))
	}
	return bir.VarKind(val)
}

func (br *BIRReader) readFlags() int64 {
	var val int64
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading flags: %v", err))
	}
	return val
}

func (br *BIRReader) readOrigin() model.SymbolOrigin {
	var val uint8
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading origin: %v", err))
	}
	return model.SymbolOrigin(val)
}

func (br *BIRReader) readStringCPEntry() model.Name {
	var idx int32
	if err := br.read(&idx); err != nil {
		panic(fmt.Sprintf("reading string CP entry index: %v", err))
	}
	return model.Name(br.getStringFromCP(int(idx)))
}

func (br *BIRReader) readLength() int64 {
	var val int64
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading length: %v", err))
	}
	return val
}

func (br *BIRReader) readInstructionKind() bir.InstructionKind {
	var val uint8
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading instruction kind: %v", err))
	}
	return bir.InstructionKind(val)
}

func (br *BIRReader) readScope() bir.VarScope {
	var val uint8
	if err := br.read(&val); err != nil {
		panic(fmt.Sprintf("reading scope: %v", err))
	}
	return bir.VarScope(val)
}

func (br *BIRReader) readPackageCPEntry() *model.PackageID {
	var idx int32
	if err := br.read(&idx); err != nil {
		panic(fmt.Sprintf("reading package CP entry index: %v", err))
	}
	return br.getPackageFromCP(int(idx))
}

func (br *BIRReader) readPosition() diagnostics.Location {
	var sourceFileIdx int32
	if err := br.read(&sourceFileIdx); err != nil {
		panic(fmt.Sprintf("reading position source file index: %v", err))
	}
	sourceFileName := br.getStringFromCP(int(sourceFileIdx))

	var sLine int32
	if err := br.read(&sLine); err != nil {
		panic(fmt.Sprintf("reading position start line: %v", err))
	}
	var sCol int32
	if err := br.read(&sCol); err != nil {
		panic(fmt.Sprintf("reading position start column: %v", err))
	}
	var eLine int32
	if err := br.read(&eLine); err != nil {
		panic(fmt.Sprintf("reading position end line: %v", err))
	}
	var eCol int32
	if err := br.read(&eCol); err != nil {
		panic(fmt.Sprintf("reading position end column: %v", err))
	}

	return diagnostics.NewBLangDiagnosticLocation(sourceFileName, int(sLine), int(eLine), int(sCol), int(eCol), 0, 0)
}
