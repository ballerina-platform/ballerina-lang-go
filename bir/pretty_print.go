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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"fmt"
	"strings"
)

type PrettyPrinter struct {
	indentLevel int
	sb          strings.Builder
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

	p.write("module ")
	p.write(p.PrintPackageID(node.PackageID))
	p.write(";\n")
	for _, importModule := range node.ImportModules {
		p.write(p.PrintImportModule(importModule))
		p.write(";\n")
	}
	for _, globalVar := range node.GlobalVars {
		p.write(p.PrintGlobalVar(globalVar))
		p.write(";\n")
	}
	for _, function := range node.Functions {
		p.PrintFunction(function)
		p.write("\n")
	}
	return p.sb.String()
}

func (p *PrettyPrinter) PrintFunction(function BIRFunction) {
	p.write(function.Name.Value())
	ty := function.Type
	if ty != nil {
		p.write("(")
		for i, param := range ty.GetParameterTypes() {
			if i > 0 {
				p.write(",")
			}
			p.write(p.PrintType(param))
		}
		p.write(")")
		if ty.GetReturnType() != nil {
			p.write(" -> ")
			p.write(p.PrintType(ty.GetReturnType()))
		}
	} else {
		p.write("<NIL>")
	}
	p.write("{\n")
	p.increaseIndent()
	for _, basicBlock := range function.BasicBlocks {
		p.PrintBasicBlock(basicBlock)
	}
	p.decreaseIndent()
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
	switch instruction.(type) {
	case *Move:
		return p.PrintMove(instruction.(*Move))
	case *BinaryOp:
		return p.PrintBinaryOp(instruction.(*BinaryOp))
	case *UnaryOp:
		return p.PrintUnaryOp(instruction.(*UnaryOp))
	case *ConstantLoad:
		return p.PrintConstantLoad(instruction.(*ConstantLoad))
	case *Goto:
		return p.PrintGoto(instruction.(*Goto))
	case *Call:
		return p.PrintCall(instruction.(*Call))
	case *Return:
		return p.PrintReturn(instruction.(*Return))
	case *Branch:
		return p.PrintBranch(instruction.(*Branch))
	case *FieldAccess:
		return p.PrintFieldAccess(instruction.(*FieldAccess))
	case *NewArray:
		return p.PrintNewArray(instruction.(*NewArray))
	case *NewMap:
		return p.PrintNewMap(instruction.(*NewMap))
	case *TypeCast:
		return p.PrintTypeCast(instruction.(*TypeCast))
	case *TypeTest:
		return p.PrintTypeTest(instruction.(*TypeTest))
	default:
		panic(fmt.Sprintf("unknown instruction type: %T", instruction))
	}
}

func (p *PrettyPrinter) PrintTypeCast(cast *TypeCast) string {
	return fmt.Sprintf("%s = <%s>(%s)", p.PrintOperand(*cast.LhsOp), cast.Type.String(), p.PrintOperand(*cast.RhsOp))
}

func (p *PrettyPrinter) PrintTypeTest(test *TypeTest) string {
	op := "is"
	if test.IsNegation {
		op = "!is"
	}
	return fmt.Sprintf("%s = %s %s %s", p.PrintOperand(*test.LhsOp), p.PrintOperand(*test.RhsOp), op, test.Type.String())
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
	return fmt.Sprintf("%s = newMap %s{%s}", p.PrintOperand(*m.LhsOp), p.PrintSemType(m.Type), values.String())
}

func (p *PrettyPrinter) PrintFieldAccess(access *FieldAccess) string {
	switch access.Kind {
	case INSTRUCTION_KIND_ARRAY_STORE:
		return fmt.Sprintf("%s[%s] = %s;", p.PrintOperand(*access.LhsOp), p.PrintOperand(*access.KeyOp), p.PrintOperand(*access.RhsOp))
	case INSTRUCTION_KIND_ARRAY_LOAD:
		return fmt.Sprintf("%s = %s[%s];", p.PrintOperand(*access.LhsOp), p.PrintOperand(*access.RhsOp), p.PrintOperand(*access.KeyOp))
	default:
		panic(fmt.Sprintf("unknown field access kind: %d", access.Kind))
	}
}

func (p *PrettyPrinter) PrintReturn(r *Return) string {
	return "return;"
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
	return operand.VariableDcl.Name.Value()
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

func (p *PrettyPrinter) PrintType(typeNode model.ValueType) string {
	if typeNode == nil {
		return "<UNKNOWN>"
	}
	sb := strings.Builder{}
	sb.WriteString(string(typeNode.GetTypeKind()))
	return sb.String()
}

func (p *PrettyPrinter) PrintSemType(typeNode semtypes.SemType) string {
	if typeNode == nil {
		return "<UNKNOWN>"
	}
	return typeNode.String()
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
