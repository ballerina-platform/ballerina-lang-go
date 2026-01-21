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

type BIRVisitor interface {
	VisitBIRPackage(birPackage BIRPackage)
	VisitBIRImportModule(birImportModule BIRImportModule)
	VisitBIRTypeDefinition(birTypeDefinition BIRTypeDefinition)
	VisitBIRVariableDcl(birVariableDcl BIRVariableDcl)
	VisitBIRGlobalVariableDcl(birGlobalVariableDcl BIRGlobalVariableDcl)
	VisitBIRFunctionParameter(birFunctionParameter BIRFunctionParameter)
	VisitBIRFunction(birFunction BIRFunction)
	VisitBIRBasicBlock(birBasicBlock BIRBasicBlock)
	VisitBIRParameter(birParameter BIRParameter)
	VisitBIRAnnotation(birAnnotation BIRAnnotation)
	VisitBIRConstant(birConstant BIRConstant)
	VisitBIRAnnotationAttachment(birAnnotAttach BIRAnnotationAttachment)
	VisitBIRConstAnnotationAttachment(birConstAnnotAttach BIRConstAnnotationAttachment)
	VisitBIRErrorEntry(birErrorEntry BIRErrorEntry)
	VisitBIRServiceDeclaration(birServiceDecl BIRServiceDeclaration)
	VisitBIRTerminatorGOTO(birGoto BIRTerminator)
	VisitBIRTerminatorCall(birCall BIRTerminator)
	VisitBIRTerminatorAsyncCall(birCall BIRTerminator)
	VisitBIRTerminatorReturn(birReturn BIRTerminator)
	VisitBIRTerminatorBranch(birBranch BIRTerminator)
	VisitBIRTerminatorFPCall(fpCall BIRTerminator)
	VisitBIRTerminatorLock(lock BIRTerminatorLock)
	VisitBIRTerminatorFieldLock(lock BIRTerminator)
	VisitBIRTerminatorUnlock(unlock BIRTerminator)
	VisitBIRTerminatorPanic(birPanic BIRTerminator)
	VisitBIRTerminatorWait(birWait BIRTerminator)
	VisitBIRTerminatorWaitAll(waitAll BIRTerminator)
	VisitBIRTerminatorFlush(birFlush BIRTerminator)
	VisitBIRTerminatorWorkerReceive(workerReceive BIRTerminator)
	VisitBIRTerminatorWorkerSend(workerSend BIRTerminator)
	VisitBIRTerminatorWorkerAlternateReceive(altReceive BIRTerminator)
	VisitBIRTerminatorWorkerMultipleReceive(multipleReceive BIRTerminator)
	VisitBIRNonTerminatorMove(birMove BIRNonTerminator)
	VisitBIRNonTerminatorBinaryOp(birBinaryOp BIRNonTerminator)
	VisitBIRNonTerminatorUnaryOP(birUnaryOp BIRNonTerminator)
	VisitBIRNonTerminatorConstantLoad(birConstantLoad BIRNonTerminator)
	VisitBIRNonTerminatorNewStructure(birNewStructure BIRNonTerminator)
	VisitBIRNonTerminatorNewArray(birNewArray BIRNonTerminator)
	VisitBIRNonTerminatorFieldAccess(birFieldAccess BIRNonTerminator)
	VisitBIRNonTerminatorNewError(birNewError BIRNonTerminator)
	VisitBIRNonTerminatorFPLoad(fpLoad BIRNonTerminator)
	VisitBIRNonTerminatorTypeCast(birTypeCast BIRNonTerminator)
	VisitBIRNonTerminatorNewInstance(newInstance BIRNonTerminator)
	VisitBIRNonTerminatorIsLike(birIsLike BIRNonTerminator)
	VisitBIRNonTerminatorTypeTest(birTypeTest BIRNonTerminator)
	VisitBIRNonTerminatorNewTable(newTable BIRNonTerminator)
	VisitBIRNonTerminatorNewTypeDesc(newTypeDesc BIRNonTerminator)
	VisitBIRNonTerminatorNewStringXMLQName(newStringXMLQName BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLProcIns(newXMLProcIns BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLComment(newXMLComment BIRNonTerminator)
	VisitBIRNonTerminatorXMLAccess(xmlAccess BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLText(newXMLText BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLSequence(newXMLSequence BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLQName(newXMLQName BIRNonTerminator)
	VisitBIRNonTerminatorNewXMLElement(newXMLElement BIRNonTerminator)
	VisitBIRNonTerminatorNewRegExp(newRegExp BIRNonTerminator)
	VisitBIRNonTerminatorNewReDisjunction(reDisjunction BIRNonTerminator)
	VisitBIRNonTerminatorNewReSequence(reSequence BIRNonTerminator)
	VisitBIRNonTerminatorNewReAssertion(reAssertion BIRNonTerminator)
	VisitBIRNonTerminatorNewReAtomQuantifier(reAtomQuantifier BIRNonTerminator)
	VisitBIRNonTerminatorNewReLiteralCharOrEscape(reLiteralCharOrEscape BIRNonTerminator)
	VisitBIRNonTerminatorNewReCharacterClass(reCharacterClass BIRNonTerminator)
	VisitBIRNonTerminatorNewReCharSet(reCharSet BIRNonTerminator)
	VisitBIRNonTerminatorNewReCharSetRange(reCharSetRange BIRNonTerminator)
	VisitBIRNonTerminatorNewReCapturingGroup(reCapturingGroup BIRNonTerminator)
	VisitBIRNonTerminatorNewReFlagExpression(reFlagExpression BIRNonTerminator)
	VisitBIRNonTerminatorNewReFlagOnOff(reFlagOnOff BIRNonTerminator)
	VisitBIRNonTerminatorNewReQuantifier(reQuantifier BIRNonTerminator)
	VisitBIRNonTerminatorRecordDefaultFPLoad(recordDefaultFPLoad BIRNonTerminator)
	VisitBIROperand(birVarRef BIROperand)
}
