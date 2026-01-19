/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

package model

import (
	"testing"

	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/compiler/model/symbols"
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/tools/text"
)

// MockLocation implements diagnostics.Location for testing
type MockLocation struct{}

func (m MockLocation) LineRange() text.LineRange {
	startPos := text.LinePositionFromLineAndOffset(1, 0)
	endPos := text.LinePositionFromLineAndOffset(1, 0)
	return text.LineRangeFromLinePositions("test.bal", startPos, endPos)
}

func (m MockLocation) TextRange() text.TextRange {
	return text.TextRangeFromStartOffsetAndLength(0, 0)
}

// MockBIRVisitor implements BIRVisitor for testing
type MockBIRVisitor struct {
	visitedNodes []string
}

func (v *MockBIRVisitor) VisitBIRPackage(birPackage BIRPackage) {
	v.visitedNodes = append(v.visitedNodes, "BIRPackage")
}

func (v *MockBIRVisitor) VisitBIRImportModule(birImportModule BIRImportModule) {
	v.visitedNodes = append(v.visitedNodes, "BIRImportModule")
}

func (v *MockBIRVisitor) VisitBIRTypeDefinition(birTypeDefinition BIRTypeDefinition) {
	v.visitedNodes = append(v.visitedNodes, "BIRTypeDefinition")
}

func (v *MockBIRVisitor) VisitBIRVariableDcl(birVariableDcl BIRVariableDcl) {
	v.visitedNodes = append(v.visitedNodes, "BIRVariableDcl")
}

func (v *MockBIRVisitor) VisitBIRGlobalVariableDcl(birGlobalVariableDcl BIRGlobalVariableDcl) {
	v.visitedNodes = append(v.visitedNodes, "BIRGlobalVariableDcl")
}

func (v *MockBIRVisitor) VisitBIRFunctionParameter(birFunctionParameter BIRFunctionParameter) {
	v.visitedNodes = append(v.visitedNodes, "BIRFunctionParameter")
}

func (v *MockBIRVisitor) VisitBIRFunction(birFunction BIRFunction) {
	v.visitedNodes = append(v.visitedNodes, "BIRFunction")
}

func (v *MockBIRVisitor) VisitBIRBasicBlock(birBasicBlock BIRBasicBlock) {
	v.visitedNodes = append(v.visitedNodes, "BIRBasicBlock")
}

func (v *MockBIRVisitor) VisitBIRParameter(birParameter BIRParameter) {
	v.visitedNodes = append(v.visitedNodes, "BIRParameter")
}

func (v *MockBIRVisitor) VisitBIRAnnotation(birAnnotation BIRAnnotation) {
	v.visitedNodes = append(v.visitedNodes, "BIRAnnotation")
}

func (v *MockBIRVisitor) VisitBIRConstant(birConstant BIRConstant) {
	v.visitedNodes = append(v.visitedNodes, "BIRConstant")
}

func (v *MockBIRVisitor) VisitBIRAnnotationAttachment(birAnnotAttach BIRAnnotationAttachment) {
	v.visitedNodes = append(v.visitedNodes, "BIRAnnotationAttachment")
}

func (v *MockBIRVisitor) VisitBIRConstAnnotationAttachment(birConstAnnotAttach BIRConstAnnotationAttachment) {
	v.visitedNodes = append(v.visitedNodes, "BIRConstAnnotationAttachment")
}

func (v *MockBIRVisitor) VisitBIRErrorEntry(birErrorEntry BIRErrorEntry) {
	v.visitedNodes = append(v.visitedNodes, "BIRErrorEntry")
}

func (v *MockBIRVisitor) VisitBIRServiceDeclaration(birServiceDecl BIRServiceDeclaration) {
	v.visitedNodes = append(v.visitedNodes, "BIRServiceDeclaration")
}

func (v *MockBIRVisitor) VisitBIRTerminatorGOTO(birGoto BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorGOTO")
}

func (v *MockBIRVisitor) VisitBIRTerminatorCall(birCall BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorCall")
}

func (v *MockBIRVisitor) VisitBIRTerminatorAsyncCall(birCall BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorAsyncCall")
}

func (v *MockBIRVisitor) VisitBIRTerminatorReturn(birReturn BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorReturn")
}

func (v *MockBIRVisitor) VisitBIRTerminatorBranch(birBranch BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorBranch")
}

func (v *MockBIRVisitor) VisitBIRTerminatorFPCall(fpCall BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorFPCall")
}

func (v *MockBIRVisitor) VisitBIRTerminatorLock(lock BIRTerminatorLock) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorLock")
}

func (v *MockBIRVisitor) VisitBIRTerminatorFieldLock(lock BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorFieldLock")
}

func (v *MockBIRVisitor) VisitBIRTerminatorUnlock(unlock BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorUnlock")
}

func (v *MockBIRVisitor) VisitBIRTerminatorPanic(birPanic BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorPanic")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWait(birWait BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWait")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWaitAll(waitAll BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWaitAll")
}

func (v *MockBIRVisitor) VisitBIRTerminatorFlush(birFlush BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorFlush")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWorkerReceive(workerReceive BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWorkerReceive")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWorkerSend(workerSend BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWorkerSend")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWorkerAlternateReceive(altReceive BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWorkerAlternateReceive")
}

func (v *MockBIRVisitor) VisitBIRTerminatorWorkerMultipleReceive(multipleReceive BIRTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRTerminatorWorkerMultipleReceive")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorMove(birMove BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorMove")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorBinaryOp(birBinaryOp BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorBinaryOp")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorUnaryOP(birUnaryOp BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorUnaryOP")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorConstantLoad(birConstantLoad BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorConstantLoad")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewStructure(birNewStructure BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewStructure")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewArray(birNewArray BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewArray")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorFieldAccess(birFieldAccess BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorFieldAccess")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewError(birNewError BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewError")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorFPLoad(fpLoad BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorFPLoad")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorTypeCast(birTypeCast BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorTypeCast")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewInstance(newInstance BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewInstance")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorIsLike(birIsLike BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorIsLike")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorTypeTest(birTypeTest BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorTypeTest")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewTable(newTable BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewTable")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewTypeDesc(newTypeDesc BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewTypeDesc")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewStringXMLQName(newStringXMLQName BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewStringXMLQName")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLProcIns(newXMLProcIns BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLProcIns")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLComment(newXMLComment BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLComment")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorXMLAccess(xmlAccess BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorXMLAccess")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLText(newXMLText BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLText")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLSequence(newXMLSequence BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLSequence")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLQName(newXMLQName BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLQName")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewXMLElement(newXMLElement BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewXMLElement")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewRegExp(newRegExp BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewRegExp")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReDisjunction(reDisjunction BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReDisjunction")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReSequence(reSequence BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReSequence")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReAssertion(reAssertion BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReAssertion")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReAtomQuantifier(reAtomQuantifier BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReAtomQuantifier")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReLiteralCharOrEscape(reLiteralCharOrEscape BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReLiteralCharOrEscape")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReCharacterClass(reCharacterClass BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReCharacterClass")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReCharSet(reCharSet BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReCharSet")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReCharSetRange(reCharSetRange BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReCharSetRange")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReCapturingGroup(reCapturingGroup BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReCapturingGroup")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReFlagExpression(reFlagExpression BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReFlagExpression")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReFlagOnOff(reFlagOnOff BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReFlagOnOff")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorNewReQuantifier(reQuantifier BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorNewReQuantifier")
}

func (v *MockBIRVisitor) VisitBIRNonTerminatorRecordDefaultFPLoad(recordDefaultFPLoad BIRNonTerminator) {
	v.visitedNodes = append(v.visitedNodes, "BIRNonTerminatorRecordDefaultFPLoad")
}

func (v *MockBIRVisitor) VisitBIROperand(birVarRef BIROperand) {
	v.visitedNodes = append(v.visitedNodes, "BIROperand")
}

func TestBIRPackage(t *testing.T) {
	pos := MockLocation{}
	org := util.NewName("testorg")
	pkgName := util.NewName("testpkg")
	name := util.NewName("test")
	version := util.NewName("1.0.0")
	sourceFileName := util.NewName("test.bal")
	sourceRoot := "/test"

	pkg := NewBIRPackage(pos, org, pkgName, name, version, sourceFileName, sourceRoot, false)

	if pkg == nil {
		t.Fatal("NewBIRPackage returned nil")
	}

	if pkg.GetPos() == nil {
		t.Error("Position not set")
	}

	pkgID := pkg.GetPackageID()
	if pkgID.OrgName.GetValue() != "testorg" {
		t.Errorf("Expected org 'testorg', got '%s'", pkgID.OrgName.GetValue())
	}

	visitor := &MockBIRVisitor{}
	pkg.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIRPackage" {
		t.Errorf("Expected visitor to visit BIRPackage, got %v", visitor.visitedNodes)
	}
}

func TestBIRImportModule(t *testing.T) {
	pos := MockLocation{}
	org := util.NewName("testorg")
	name := util.NewName("test")
	version := util.NewName("1.0.0")

	mod := NewBIRImportModule(pos, org, name, version)

	if mod == nil {
		t.Fatal("NewBIRImportModule returned nil")
	}

	pkgID := mod.GetPackageID()
	if pkgID.OrgName.GetValue() != "testorg" {
		t.Errorf("Expected org 'testorg', got '%s'", pkgID.OrgName.GetValue())
	}

	// Test equals
	mod2 := NewBIRImportModule(pos, org, name, version)
	if !mod.Equals(mod2) {
		t.Error("Expected equal modules to be equal")
	}

	mod3 := NewBIRImportModule(pos, util.NewName("other"), name, version)
	if mod.Equals(mod3) {
		t.Error("Expected different modules to not be equal")
	}

	visitor := &MockBIRVisitor{}
	mod.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIRImportModule" {
		t.Errorf("Expected visitor to visit BIRImportModule, got %v", visitor.visitedNodes)
	}
}

func TestBIRVariableDcl(t *testing.T) {
	var pos diagnostics.Location = nil
	var type_ BType = nil // BType would need to be implemented
	name := util.NewName("var1")
	originalName := util.NewName("var1")
	scope := VAR_SCOPE_FUNCTION
	kind := VAR_KIND_LOCAL
	metaVarName := "var1"

	varDecl := NewBIRVariableDcl(pos, type_, name, originalName, scope, kind, metaVarName)

	if varDecl == nil {
		t.Fatal("NewBIRVariableDcl returned nil")
	}

	if varDecl.GetName().GetValue() != "var1" {
		t.Errorf("Expected name 'var1', got '%s'", varDecl.GetName().GetValue())
	}

	if varDecl.GetScope() != scope {
		t.Error("Scope not set correctly")
	}

	if varDecl.GetKind() != kind {
		t.Error("Kind not set correctly")
	}

	// Test equals
	varDecl2 := NewBIRVariableDcl(pos, type_, name, originalName, scope, kind, metaVarName)
	if !varDecl.Equals(varDecl2) {
		t.Error("Expected equal variable declarations to be equal")
	}

	visitor := &MockBIRVisitor{}
	varDecl.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIRVariableDcl" {
		t.Errorf("Expected visitor to visit BIRVariableDcl, got %v", visitor.visitedNodes)
	}
}

func TestBIRParameter(t *testing.T) {
	pos := MockLocation{}
	name := util.NewName("param1")
	flags := int64(0)

	param := NewBIRParameter(pos, name, flags)

	if param == nil {
		t.Fatal("NewBIRParameter returned nil")
	}

	if param.GetName().GetValue() != "param1" {
		t.Errorf("Expected name 'param1', got '%s'", param.GetName().GetValue())
	}

	if param.GetFlags() != flags {
		t.Error("Flags not set correctly")
	}

	visitor := &MockBIRVisitor{}
	param.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIRParameter" {
		t.Errorf("Expected visitor to visit BIRParameter, got %v", visitor.visitedNodes)
	}
}

func TestBIRBasicBlock(t *testing.T) {
	id := util.NewName("bb1")
	number := 1

	bb := NewBIRBasicBlock(id, number)

	if bb == nil {
		t.Fatal("NewBIRBasicBlock returned nil")
	}

	if bb.GetId().GetValue() != "bb1" {
		t.Errorf("Expected id 'bb1', got '%s'", bb.GetId().GetValue())
	}

	if bb.GetNumber() != number {
		t.Error("Number not set correctly")
	}

	// Test string representation
	if bb.String() != "bb1" {
		t.Errorf("Expected string 'bb1', got '%s'", bb.String())
	}

	visitor := &MockBIRVisitor{}
	bb.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIRBasicBlock" {
		t.Errorf("Expected visitor to visit BIRBasicBlock, got %v", visitor.visitedNodes)
	}
}

func TestBIROperand(t *testing.T) {
	var pos diagnostics.Location = nil
	var type_ BType = nil // BType would need to be implemented
	name := util.NewName("var1")
	scope := VAR_SCOPE_FUNCTION
	kind := VAR_KIND_LOCAL

	varDecl := NewBIRVariableDcl(pos, type_, name, name, scope, kind, "var1")
	operand := NewBIROperand(varDecl)

	if operand == nil {
		t.Fatal("NewBIROperand returned nil")
	}

	if operand.GetVariableDcl() == nil {
		t.Error("Variable declaration not set correctly")
	}

	// Test equals
	operand2 := NewBIROperand(varDecl)
	if !operand.Equals(operand2) {
		t.Error("Expected equal operands to be equal")
	}

	visitor := &MockBIRVisitor{}
	operand.Accept(visitor)

	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "BIROperand" {
		t.Errorf("Expected visitor to visit BIROperand, got %v", visitor.visitedNodes)
	}
}

func TestChannelDetails(t *testing.T) {
	cd := NewChannelDetails("channel1", true, false)

	if cd == nil {
		t.Fatal("NewChannelDetails returned nil")
	}

	if cd.Name != "channel1" {
		t.Errorf("Expected name 'channel1', got '%s'", cd.Name)
	}

	if !cd.ChannelInSameStrand {
		t.Error("ChannelInSameStrand not set correctly")
	}

	if cd.Send {
		t.Error("Send not set correctly")
	}

	// Test string representation
	if cd.String() != "channel1" {
		t.Errorf("Expected string 'channel1', got '%s'", cd.String())
	}
}

func TestBIRLockDetailsHolder(t *testing.T) {
	holder := NewBIRLockDetailsHolder()

	if holder == nil {
		t.Fatal("NewBIRLockDetailsHolder returned nil")
	}

	if !holder.IsEmpty() {
		t.Error("Expected holder to be empty initially")
	}

	if holder.Size() != 0 {
		t.Error("Expected size to be 0 initially")
	}

	// Note: We can't test AddLock/GetLock without a BIRTerminatorLock implementation
}

func TestVarKind(t *testing.T) {
	if VAR_KIND_LOCAL.GetValue() != 1 {
		t.Error("VAR_KIND_LOCAL value incorrect")
	}
	if VAR_KIND_ARG.GetValue() != 2 {
		t.Error("VAR_KIND_ARG value incorrect")
	}
	if VAR_KIND_TEMP.GetValue() != 3 {
		t.Error("VAR_KIND_TEMP value incorrect")
	}
}

func TestVarScope(t *testing.T) {
	if VAR_SCOPE_FUNCTION.GetValue() != 1 {
		t.Error("VAR_SCOPE_FUNCTION value incorrect")
	}
	if VAR_SCOPE_GLOBAL.GetValue() != 2 {
		t.Error("VAR_SCOPE_GLOBAL value incorrect")
	}
}

func TestBirScope(t *testing.T) {
	scope1 := NewBirScope(1, nil)
	scope2 := NewBirScope(2, &scope1)

	if scope1.Id != 1 {
		t.Error("Scope ID not set correctly")
	}

	if scope2.Id != 2 {
		t.Error("Scope ID not set correctly")
	}

	if scope2.Parent == nil || scope2.Parent.Id != 1 {
		t.Error("Parent scope not set correctly")
	}
}

func TestName(t *testing.T) {
	name := util.NewName("test")
	if name.GetValue() != "test" {
		t.Error("Name value not set correctly")
	}

	name2 := util.NewName("test")
	if !name.Equals(name2) {
		t.Error("Expected equal names to be equal")
	}

	name3 := util.NewName("other")
	if name.Equals(name3) {
		t.Error("Expected different names to not be equal")
	}

	if name.String() != "test" {
		t.Errorf("Expected string 'test', got '%s'", name.String())
	}
}

func TestPackageID(t *testing.T) {
	org := util.NewName("testorg")
	name := util.NewName("test")
	version := util.NewName("1.0.0")

	pkgID := elements.NewPackageIDWithOrgNameVersion(org, name, version)

	if pkgID.OrgName.GetValue() != "testorg" {
		t.Errorf("Expected org 'testorg', got '%s'", pkgID.OrgName.GetValue())
	}

	if pkgID.Name.GetValue() != "test" {
		t.Errorf("Expected name 'test', got '%s'", pkgID.Name.GetValue())
	}

	if pkgID.Version.GetValue() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", pkgID.Version.GetValue())
	}

	// Test equals
	pkgID2 := elements.NewPackageIDWithOrgNameVersion(org, name, version)
	if !pkgID.Equals(pkgID2) {
		t.Error("Expected equal package IDs to be equal")
	}

	pkgID3 := elements.NewPackageIDWithOrgNameVersion(util.NewName("other"), name, version)
	if pkgID.Equals(pkgID3) {
		t.Error("Expected different package IDs to not be equal")
	}
}

func TestSymbolOrigin(t *testing.T) {
	if symbols.SYMBOL_ORIGIN_BUILTIN.Value() != 1 {
		t.Error("SYMBOL_ORIGIN_BUILTIN value incorrect")
	}
	if symbols.SYMBOL_ORIGIN_SOURCE.Value() != 2 {
		t.Error("SYMBOL_ORIGIN_SOURCE value incorrect")
	}

	origin := symbols.SYMBOL_ORIGIN_SOURCE.ToBIROrigin()
	if origin != symbols.SYMBOL_ORIGIN_COMPILED_SOURCE {
		t.Error("ToBIROrigin conversion incorrect")
	}

	converted := symbols.ToOrigin(1)
	if converted != symbols.SYMBOL_ORIGIN_BUILTIN {
		t.Error("ToOrigin conversion incorrect")
	}
}

func TestInstructionKind(t *testing.T) {
	if INSTRUCTION_KIND_GOTO.GetValue() != 1 {
		t.Error("INSTRUCTION_KIND_GOTO value incorrect")
	}
	if INSTRUCTION_KIND_CALL.GetValue() != 2 {
		t.Error("INSTRUCTION_KIND_CALL value incorrect")
	}
	if INSTRUCTION_KIND_MOVE.GetValue() != 20 {
		t.Error("INSTRUCTION_KIND_MOVE value incorrect")
	}
}

func TestFlag(t *testing.T) {
	// Test that flags are defined
	_ = elements.FLAG_PUBLIC
	_ = elements.FLAG_PRIVATE
	_ = elements.FLAG_REMOTE
	_ = elements.FLAG_ISOLATED
	_ = elements.FLAG_CONSTANT
	// Just verify they compile
}
