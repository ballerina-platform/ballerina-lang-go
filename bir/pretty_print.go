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
	"fmt"
	"sort"
	"strings"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

type PrettyPrinter struct {
	indentLevel int
	sb          strings.Builder
	cx          semtypes.Context
}

// writeLine writes a line with current indentation and newline
func (p *PrettyPrinter) writeLine(s string) {
	for i := 0; i < p.indentLevel; i++ {
		p.sb.WriteString("  ")
	}
	p.sb.WriteString(s)
	p.sb.WriteString("\n")
}

// write writes without indentation or newline
func (p *PrettyPrinter) write(s string) {
	p.sb.WriteString(s)
}

// writeIndent writes current indentation without content or newline
func (p *PrettyPrinter) writeIndent() {
	for i := 0; i < p.indentLevel; i++ {
		p.sb.WriteString("  ")
	}
}

// increaseIndent increases indentation level
func (p *PrettyPrinter) increaseIndent() {
	p.indentLevel++
}

// decreaseIndent decreases indentation level
func (p *PrettyPrinter) decreaseIndent() {
	p.indentLevel--
}

func (p *PrettyPrinter) Print(node BIRPackage) string {
	// Reset the builder
	p.sb.Reset()
	p.cx = semtypes.TypeCheckContext(node.TypeEnv)

	p.write("module ")
	p.write(p.PrintPackageID(node.PackageID))
	p.write(";\n")
	for _, importModule := range node.ImportModules {
		p.write(p.PrintImportModule(importModule))
		p.write(";\n")
	}
	sortedGlobalVars := make([]BIRGlobalVariableDcl, 0, len(node.GlobalVars))
	for _, globalVar := range node.GlobalVars {
		sortedGlobalVars = append(sortedGlobalVars, globalVar)
	}
	sort.Slice(sortedGlobalVars, func(i, j int) bool {
		return string(sortedGlobalVars[i].GetName()) < string(sortedGlobalVars[j].GetName())
	})
	for _, globalVar := range sortedGlobalVars {
		p.write(p.PrintGlobalVar(globalVar))
		p.write(";\n")
	}
	for _, classDef := range node.ClassDefs {
		p.PrintClassDef(classDef)
		p.write("\n")
	}
	for _, function := range node.Functions {
		p.PrintFunction(function)
		p.write("\n")
	}
	return p.sb.String()
}

func (p *PrettyPrinter) PrintFunction(function BIRFunction) {
	p.write(function.Name.Value())
	p.write("(")
	paramStart := 1
	if len(function.LocalVars) > 1 && function.LocalVars[1].GetName() == "self" {
		paramStart = 2
	}
	for i, v := range function.LocalVars[paramStart:] {
		if i < len(function.RequiredParams) {
			if i > 0 {
				p.write(",")
			}
			p.write(p.PrintSemType(v.Type))
		} else {
			break
		}
	}
	if function.RestParams != nil {
		variableIndex := paramStart + len(function.RequiredParams)
		if variableIndex != 1 {
			p.write(",")
		}
		p.write(p.PrintSemType(function.LocalVars[variableIndex].Type))
		p.write("...")
	}
	p.write(")")
	if function.ReturnVariable != nil && function.ReturnVariable.Type != nil {
		p.write(" -> ")
		p.write(p.PrintSemType(function.ReturnVariable.Type))
	}
	p.write("{\n")
	p.increaseIndent()
	for _, basicBlock := range function.BasicBlocks {
		p.PrintBasicBlock(basicBlock)
	}
	if len(function.ErrorTable) > 0 {
		p.writeLine("")
		p.writeLine("error-table {")
		p.increaseIndent()
		for _, entry := range function.ErrorTable {
			p.writeLine(fmt.Sprintf("[bb%d, bb%d] -> bb%d, %s", entry.Start, entry.End, entry.Target, p.PrintOperand(*entry.ErrorOp)))
		}
		p.decreaseIndent()
		p.writeLine("}")
	}
	p.decreaseIndent()
	p.writeIndent()
	p.write("}")
}

func (p *PrettyPrinter) PrintBasicBlock(basicBlock BIRBasicBlock) {
	p.writeLine(basicBlock.Id.Value() + " {")
	p.increaseIndent()
	for _, instruction := range basicBlock.Instructions {
		p.writeLine(p.PrintInstruction(instruction))
	}
	if basicBlock.Terminator != nil {
		p.writeLine(p.PrintInstruction(basicBlock.Terminator))
	}
	p.decreaseIndent()
	p.writeLine("}")
}

func (p *PrettyPrinter) PrintInstruction(instruction BIRInstruction) string {
	switch instruction := instruction.(type) {
	case *Move:
		return p.PrintMove(instruction)
	case *BinaryOp:
		return p.PrintBinaryOp(instruction)
	case *UnaryOp:
		return p.PrintUnaryOp(instruction)
	case *ConstantLoad:
		return p.PrintConstantLoad(instruction)
	case *Goto:
		return p.PrintGoto(instruction)
	case *Call:
		return p.PrintCall(instruction)
	case *Return:
		return p.PrintReturn(instruction)
	case *Branch:
		return p.PrintBranch(instruction)
	case *FieldAccess:
		return p.PrintFieldAccess(instruction)
	case *NewArray:
		return p.PrintNewArray(instruction)
	case *NewMap:
		return p.PrintNewMap(instruction)
	case *NewError:
		return p.PrintNewError(instruction)
	case *TypeCast:
		return p.PrintTypeCast(instruction)
	case *TypeTest:
		return p.PrintTypeTest(instruction)
	case *Panic:
		return p.PrintPanic(instruction)
	case *NewObject:
		return p.PrintNewObject(instruction)
	case *FPLoad:
		return p.PrintFPLoad(instruction)
	case *PushScopeFrame:
		return p.PrintPushScopeFrame(instruction)
	case *PopScopeFrame:
		return "PopScopeFrame"
	default:
		panic(fmt.Sprintf("unknown instruction type: %T", instruction))
	}
}

func (p *PrettyPrinter) PrintFPLoad(fpLoad *FPLoad) string {
	kind := "fp"
	if fpLoad.IsClosure {
		kind = "closure_fp"
	}
	return fmt.Sprintf("%s = %s %s", p.PrintOperand(*fpLoad.LhsOp), kind, fpLoad.FunctionLookupKey)
}

func (p *PrettyPrinter) PrintPushScopeFrame(push *PushScopeFrame) string {
	return fmt.Sprintf("PushScopeFrame %d", push.NumLocals)
}

func (p *PrettyPrinter) PrintTypeCast(cast *TypeCast) string {
	return fmt.Sprintf("%s = <%s>(%s)", p.PrintOperand(*cast.LhsOp), semtypes.ToString(p.cx, cast.Type), p.PrintOperand(*cast.RhsOp))
}

func (p *PrettyPrinter) PrintTypeTest(test *TypeTest) string {
	op := "is"
	if test.IsNegation {
		op = "!is"
	}
	return fmt.Sprintf("%s = %s %s %s", p.PrintOperand(*test.LhsOp), p.PrintOperand(*test.RhsOp), op, semtypes.ToString(p.cx, test.Type))
}

func (p *PrettyPrinter) PrintNewArray(array *NewArray) string {
	values := strings.Builder{}
	for i, v := range array.Values {
		if i > 0 {
			values.WriteString(", ")
		}
		values.WriteString(p.PrintOperand(*v))
	}
	return fmt.Sprintf("%s = newArray %s[%s]{%s}", p.PrintOperand(*array.LhsOp), p.PrintSemType(array.Type), p.PrintOperand(*array.SizeOp), values.String())
}

func (p *PrettyPrinter) PrintNewMap(m *NewMap) string {
	values := strings.Builder{}
	for i, entry := range m.Values {
		if i > 0 {
			values.WriteString(", ")
		}
		if entry.IsKeyValuePair() {
			kv := entry.(*MappingConstructorKeyValueEntry)
			values.WriteString(p.PrintOperand(*kv.KeyOp()))
			values.WriteString("=")
			values.WriteString(p.PrintOperand(*kv.ValueOp()))
		} else {
			values.WriteString(p.PrintOperand(*entry.ValueOp()))
		}
	}
	defaults := strings.Builder{}
	for i, def := range m.Defaults {
		if i > 0 {
			defaults.WriteString(", ")
		}
		defaults.WriteString(def.FieldName)
		defaults.WriteString("=")
		defaults.WriteString(def.FunctionLookupKey)
	}
	if defaults.Len() > 0 {
		return fmt.Sprintf("%s = newMap %s{%s} defaults{%s}", p.PrintOperand(*m.LhsOp), p.PrintSemType(m.Type), values.String(), defaults.String())
	}
	return fmt.Sprintf("%s = newMap %s{%s}", p.PrintOperand(*m.LhsOp), p.PrintSemType(m.Type), values.String())
}

func (p *PrettyPrinter) PrintNewError(e *NewError) string {
	args := p.PrintOperand(*e.MessageOp)
	if e.CauseOp != nil {
		args += ", " + p.PrintOperand(*e.CauseOp)
	}
	if e.DetailOp != nil {
		args += ", " + p.PrintOperand(*e.DetailOp)
	}
	return fmt.Sprintf("%s = newError %s(%s)", p.PrintOperand(*e.LhsOp), p.PrintSemType(e.Type), args)
}

func (p *PrettyPrinter) PrintFieldAccess(access *FieldAccess) string {
	switch access.Kind {
	case INSTRUCTION_KIND_MAP_STORE, INSTRUCTION_KIND_ARRAY_STORE, INSTRUCTION_KIND_OBJECT_STORE:
		return fmt.Sprintf("%s[%s] = %s;", p.PrintOperand(*access.LhsOp), p.PrintOperand(*access.KeyOp), p.PrintOperand(*access.RhsOp))
	case INSTRUCTION_KIND_MAP_LOAD, INSTRUCTION_KIND_ARRAY_LOAD, INSTRUCTION_KIND_OBJECT_LOAD:
		return fmt.Sprintf("%s = %s[%s];", p.PrintOperand(*access.LhsOp), p.PrintOperand(*access.RhsOp), p.PrintOperand(*access.KeyOp))
	default:
		panic(fmt.Sprintf("unknown field access kind: %d", access.Kind))
	}
}

func (p *PrettyPrinter) PrintNewObject(n *NewObject) string {
	return fmt.Sprintf("%s = newObject %s", p.PrintOperand(*n.LhsOp), n.ClassDef.Name.Value())
}

func (p *PrettyPrinter) PrintClassDef(classDef BIRClassDef) {
	p.write("class ")
	p.write(classDef.Name.Value())
	p.write(" {\n")
	p.increaseIndent()
	for _, field := range classDef.Fields {
		p.writeLine(fmt.Sprintf("%s %s", field.Name, p.PrintSemType(field.Ty)))
	}
	var methodNames []string
	for name := range classDef.VTable {
		methodNames = append(methodNames, name)
	}
	sort.Strings(methodNames)
	for _, name := range methodNames {
		p.write("\n")
		p.writeIndent()
		p.PrintFunction(*classDef.VTable[name])
		p.write("\n")
	}
	p.decreaseIndent()
	p.write("}")
}

func (p *PrettyPrinter) PrintReturn(r *Return) string {
	return "return;"
}

func (p *PrettyPrinter) PrintPanic(pa *Panic) string {
	return fmt.Sprintf("panic %s;", p.PrintOperand(*pa.ErrorOp))
}

func (p *PrettyPrinter) PrintBranch(b *Branch) string {
	return fmt.Sprintf("%s ? %s : %s;", p.PrintOperand(*b.Op), b.TrueBB.Id.Value(), b.FalseBB.Id.Value())
}

func (p *PrettyPrinter) PrintGoto(g *Goto) string {
	return fmt.Sprintf("GOTO %s;", g.ThenBB.Id.Value())
}

func (p *PrettyPrinter) PrintCall(call *Call) string {
	args := strings.Builder{}
	for i, arg := range call.Args {
		if i > 0 {
			args.WriteString(",")
		}
		args.WriteString(p.PrintOperand(arg))
	}
	return fmt.Sprintf("%s = %s(%s) -> %s;", p.PrintOperand(*call.LhsOp), call.Name.Value(), args.String(), call.ThenBB.Id.Value())
}

func (p *PrettyPrinter) PrintOperand(operand BIROperand) string {
	name := operand.VariableDcl.GetName()
	if operand.Address.Mode == AddressingModeAbsolute {
		return fmt.Sprintf("(%d, %s)", operand.Address.BaseIndex, name)
	}
	return name.Value()
}

func (p *PrettyPrinter) PrintConstantLoad(load *ConstantLoad) string {
	return fmt.Sprintf("%s = ConstantLoad %v", p.PrintOperand(*load.LhsOp), load.Value)
}

func (p *PrettyPrinter) PrintUnaryOp(op *UnaryOp) string {
	return fmt.Sprintf("%s = %s %s;", p.PrintOperand(*op.LhsOp), p.PrintInstructionKind(op.Kind), p.PrintOperand(*op.RhsOp))
}

func (p *PrettyPrinter) PrintBinaryOp(op *BinaryOp) string {
	return fmt.Sprintf("%s = %s %s %s;", p.PrintOperand(*op.LhsOp), p.PrintInstructionKind(op.Kind), p.PrintOperand(op.RhsOp1), p.PrintOperand(op.RhsOp2))
}

func (p *PrettyPrinter) PrintInstructionKind(kind InstructionKind) string {
	switch kind {
	case INSTRUCTION_KIND_ADD:
		return "+"
	case INSTRUCTION_KIND_SUB:
		return "-"
	case INSTRUCTION_KIND_MUL:
		return "*"
	case INSTRUCTION_KIND_DIV:
		return "/"
	case INSTRUCTION_KIND_MOD:
		return "%"
	case INSTRUCTION_KIND_AND:
		return "&&"
	case INSTRUCTION_KIND_OR:
		return "||"
	case INSTRUCTION_KIND_LESS_THAN:
		return "<"
	case INSTRUCTION_KIND_LESS_EQUAL:
		return "<="
	case INSTRUCTION_KIND_GREATER_THAN:
		return ">"
	case INSTRUCTION_KIND_GREATER_EQUAL:
		return ">="
	case INSTRUCTION_KIND_EQUAL:
		return "=="
	case INSTRUCTION_KIND_NOT_EQUAL:
		return "!="
	case INSTRUCTION_KIND_NOT:
		return "!"
	case INSTRUCTION_KIND_BITWISE_COMPLEMENT:
		return "~"
	}
	return "unknown"
}

func (p *PrettyPrinter) PrintMove(move *Move) string {
	return fmt.Sprintf("%s = %s;", p.PrintOperand(*move.LhsOp), p.PrintOperand(*move.RhsOp))
}

func (p *PrettyPrinter) PrintGlobalVar(globalVar BIRGlobalVariableDcl) string {
	sb := strings.Builder{}
	sb.WriteString(globalVar.Name.Value())
	sb.WriteString("  ")
	sb.WriteString(p.PrintSemType(globalVar.Type))
	return sb.String()
}

func (p *PrettyPrinter) PrintSemType(typeNode semtypes.SemType) string {
	if typeNode == nil {
		return "<UNKNOWN>"
	}
	return semtypes.ToString(p.cx, typeNode)
}

func (p *PrettyPrinter) PrintImportModule(importModules BIRImportModule) string {
	sb := strings.Builder{}
	sb.WriteString("import ")
	sb.WriteString(p.PrintPackageID(importModules.PackageID))
	return sb.String()
}

func (p *PrettyPrinter) PrintPackageID(packageID *model.PackageID) string {
	if packageID.IsUnnamed() {
		return "$anon-package"
	}
	orgName := string(*packageID.OrgName)
	pkgName := string(*packageID.PkgName)
	version := string(*packageID.Version)
	return fmt.Sprintf("%s.%s v %s", orgName, pkgName, version)
}
