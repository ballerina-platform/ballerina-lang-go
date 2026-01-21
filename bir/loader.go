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
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type bbMap map[string]*BIRBasicBlock

func LoadBIRPackageFromReader(r io.Reader) (*BIRPackage, error) {
	// Read all data into a buffer since kaitai.NewStream requires io.ReadSeeker
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading BIR binary data: %w", err)
	}

	stream := kaitai.NewStream(bytes.NewReader(data))

	b := NewBir()
	if err := b.Read(stream, nil, b); err != nil {
		return nil, fmt.Errorf("reading BIR binary with kaitai: %w", err)
	}

	if b.Module == nil {
		return nil, fmt.Errorf("bir module is nil")
	}

	// Build package id from constant pool + module id index.
	pkg, err := buildBIRPackage(b)
	if err != nil {
		return nil, err
	}

	// Imports.
	if err := populateImports(b, pkg); err != nil {
		return nil, err
	}

	// Constants.
	if err := populateConstants(b, pkg); err != nil {
		return nil, err
	}

	// Type defs (container only, no bodies/attached funcs yet).
	if err := populateTypeDefs(b, pkg); err != nil {
		return nil, err
	}

	// Type def bodies (attached functions and referenced types).
	if err := populateTypeDefBodies(b, pkg); err != nil {
		return nil, err
	}

	// Global vars (meta only).
	if err := populateGlobals(b, pkg); err != nil {
		return nil, err
	}

	// Functions (only high‑level meta: names/flags/origin/required params).
	if err := populateFunctions(b, pkg); err != nil {
		return nil, err
	}

	// Annotations.
	if err := populateAnnotations(b, pkg); err != nil {
		return nil, err
	}

	// Service declarations.
	if err := populateServices(b, *pkg); err != nil {
		return nil, err
	}

	// TODO: types, globals, annotations, services etc. can be wired in the
	// same style as needed. For now we keep them empty – the model remains
	// structurally valid and the code builds.

	return pkg, nil
}

// LoadBIRPackageFromFile loads a BIR package from a file in the given filesystem.
// It opens the file at the specified path, reads its contents, and parses it as a BIR package.
func LoadBIRPackageFromFile(fsys fs.FS, path string) (BIRPackage, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening BIR file %q: %w", path, err)
	}
	defer file.Close()

	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("loading BIR package from %q: %w", path, err)
	}

	return pkg, nil
}

// buildBIRPackage constructs the BIRPackage shell: only PackageID and empty
// collections are initialized here; contents are populated by the helpers
// below.
func buildBIRPackage(b *Bir) (*BIRPackage, error) {
	mod := b.Module

	// The module id is a CP index into the constant pool; the corresponding
	// entry is a PackageCP entry.
	pkgCp, err := cpAsPackage(b, mod.IdCpIndex)
	if err != nil {
		return nil, fmt.Errorf("reading module package id: %w", err)
	}

	org := model.Name(cpString(b, pkgCp.OrgIndex))
	pkgName := model.Name(cpString(b, pkgCp.PackageNameIndex))
	name := model.Name(cpString(b, pkgCp.NameIndex))
	version := model.Name(cpString(b, pkgCp.VersionIndex))

	packageID := &model.PackageID{
		OrgName: &org,
		PkgName: &pkgName,
		Name:    &name,
		Version: &version,
	}

	return &BIRPackage{
		PackageID: packageID,
	}, nil
}

// populateImports fills BIRPackage.importModules from Module.Imports.
func populateImports(b *Bir, pkg *BIRPackage) error {
	if b.Module.ImportCount == 0 {
		return nil
	}

	imports := make([]BIRImportModule, 0, len(b.Module.Imports))

	for _, imp := range b.Module.Imports {
		if imp == nil {
			continue
		}
		org := model.Name(cpString(b, imp.OrgIndex))
		// Note: PackageNameIndex is available but NewBIRImportModule doesn't support separate pkgName
		// It uses NewPackageIDWithOrgNameVersion which sets PkgName = name
		name := model.Name(cpString(b, imp.NameIndex))
		version := model.Name(cpString(b, imp.VersionIndex))
		packageID := &model.PackageID{
			OrgName: &org,
			PkgName: &name,
			Name:    &name,
			Version: &version,
		}
		importModule := BIRImportModule{
			PackageID: packageID,
		}
		imports = append(imports, importModule)
	}

	pkg.ImportModules = imports
	return nil
}

// populateFunctions creates a minimal BIRFunction node for each Bir_Function.
func populateFunctions(b *Bir, pkg *BIRPackage) error {
	if b.Module.FunctionCount == 0 {
		return nil
	}

	funcs := make([]BIRFunction, 0, len(b.Module.Functions))

	for _, f := range b.Module.Functions {
		if f == nil {
			continue
		}

		name := model.Name(cpString(b, f.NameCpIndex))
		origName := model.Name(cpString(b, f.OriginalNameCpIndex))
		// workerName := model.Name(cpString(b, f.WorkerNameCpIndex))
		origin := model.SymbolOrigin(f.Origin)
		pos := positionToLocation(b, f.Position)

		fn := BIRFunction{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:         name,
			OriginalName: origName,
			Flags:        f.Flags,
			Origin:       origin,
		}

		// Required params -> BIRParameter list.
		if f.RequiredParamCount > 0 {
			params := make([]BIRParameter, 0, len(f.RequiredParams))
			for _, rp := range f.RequiredParams {
				if rp == nil {
					continue
				}
				pName := model.Name(cpString(b, rp.ParamNameCpIndex))
				p := BIRParameter{
					Name:  pName,
					Flags: rp.Flags,
				}
				// TODO: Parse parameter annotations from rp if available
				params = append(params, p)
			}
			fn.RequiredParams = params
		}

		// Rest parameter
		if f.HasRestParam != 0 {
			restParamName := model.Name(cpString(b, f.RestParamNameCpIndex))
			restParam := BIRParameter{
				Name:  restParamName,
				Flags: 0,
			}
			// Note: BIRParameter doesn't have SetAnnotAttachments in the current model
			// Rest param annotations are available in f.RestParamAnnotations if needed
			fn.RestParams = &restParam
		}

		// Receiver
		if f.HasReceiver != 0 && f.Reciever != nil {
			panic("receiver not supported")
		}

		// Path parameters (resource function)
		if f.IsResourceFunction != 0 && f.ResourceFunctionContent != nil {
			panic("resource function not supported")
		}

		// Annotation attachments
		if f.AnnotationAttachmentsContent != nil {
			fmt.Println("WARNING: annotation attachments not supported ignoring")
		}

		// Return type annotations
		if f.ReturnTypeAnnotations != nil {
			fmt.Println("WARNING: return type annotations not supported ignoring")
		}

		// Markdown doc attachment
		if f.Doc != nil {
			fmt.Println("WARNING: markdown doc attachment not supported ignoring")
		}

		// Scope entries (instruction vs scope table)
		if f.ScopeEntryCount > 0 && len(f.ScopeEntries) > 0 {
			// TODO: Parse scope entries and map them to instructions
			// Scope entries map instruction offsets to scope IDs
			// This is used for debugging/analysis but not critical for basic functionality
		}

		// Dependent global vars
		if f.DependentGlobalVarLength > 0 && len(f.DependentGlobalVarCpEntry) > 0 {
			dependentVars := make([]BIRGlobalVariableDcl, 0, len(f.DependentGlobalVarCpEntry))
			pkgID := pkg.PackageID
			for _, cpIdx := range f.DependentGlobalVarCpEntry {
				varName := model.Name(cpString(b, cpIdx))
				// Create a minimal global var reference
				gv := BIRGlobalVariableDcl{
					BIRVariableDcl: BIRVariableDcl{
						Name:         varName,
						OriginalName: varName,
						MetaVarName:  varName.Value(),
						Scope:        VAR_SCOPE_GLOBAL,
						Kind:         VAR_KIND_LOCAL,
					},
					Flags:  0,
					PkgId:  pkgID,
					Origin: model.SymbolOrigin_SOURCE,
				}
				dependentVars = append(dependentVars, gv)
			}

			fn.DependentGlobalVars = dependentVars
		}

		// Populate function body (basic blocks, instructions, etc.)
		if f.FunctionBody != nil {
			if err := populateFunctionBody(b, &fn, f.FunctionBody); err != nil {
				return fmt.Errorf("populating function body for %s: %w", name.Value(), err)
			}
		}

		// We also set the function's flags/origin/workerName via the ctor.
		funcs = append(funcs, fn)
	}

	pkg.Functions = funcs
	return nil
}

// populateConstants maps Bir_Constant -> BIRConstant.
func populateConstants(b *Bir, pkg *BIRPackage) error {
	if b.Module.ConstCount == 0 {
		return nil
	}

	consts := make([]BIRConstant, 0, len(b.Module.Constants))

	for _, c := range b.Module.Constants {
		if c == nil {
			continue
		}
		name := model.Name(cpString(b, c.NameCpIndex))
		origin := model.SymbolOrigin(c.Origin)
		pos := positionToLocation(b, c.Position)

		// Type + value + attachments left nil/zeroed; mapping them would
		// require full type + const value decoding.
		var t model.ValueType
		cv := ConstValue{}

		// Parse type
		if c.TypeCpIndex >= 0 {
			t = parseTypeFromCP(b, c.TypeCpIndex)
		}

		// Parse constant value
		if c.ConstantValue != nil {
			cv = parseConstantValue(b, c.ConstantValue)
		}

		bc := BIRConstant{
			BIRDocumentableNodeBase: BIRDocumentableNodeBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			Name:       name,
			Flags:      c.Flags,
			Type:       t,
			ConstValue: cv,
			Origin:     origin,
		}

		// Parse markdown doc attachment
		if c.Doc != nil {
			panic("markdown doc attachment not supported")
		}

		// Parse annotation attachments
		if c.AnnotationAttachmentsContent != nil {
			panic("annotation attachments not supported")
		}

		consts = append(consts, bc)
	}

	pkg.Constants = consts
	return nil
}

// populateTypeDefs maps Bir_TypeDefinition -> BIRTypeDefinition.
func populateTypeDefs(b *Bir, pkg *BIRPackage) error {
	if b.Module.TypeDefinitionCount == 0 {
		return nil
	}

	fmt.Println("WARNING: type definition not supported ignoring")
	return nil
}

// populateTypeDefBodies populates attached functions and referenced types for type definitions.
func populateTypeDefBodies(b *Bir, pkg *BIRPackage) error {
	if b.Module.TypeDefinitionBodiesCount == 0 {
		return nil
	}
	fmt.Println("WARNING: type definition not supported ignoring")
	return nil
}

// populateGlobals maps Bir_GlobalVar -> BIRGlobalVariableDcl.
func populateGlobals(b *Bir, pkg *BIRPackage) error {
	if b.Module.GlobalVarCount == 0 {
		return nil
	}

	globals := make([]BIRGlobalVariableDcl, 0, len(b.Module.GlobalVars))

	// Use the package's PackageID for all globals.
	pkgID := pkg.PackageID

	for _, gv := range b.Module.GlobalVars {
		if gv == nil {
			continue
		}
		name := model.Name(cpString(b, gv.NameCpIndex))
		origin := model.SymbolOrigin(gv.Origin)
		pos := positionToLocation(b, gv.Position)

		kind := VarKind(gv.Kind)
		scope := VAR_SCOPE_GLOBAL

		var t model.ValueType
		metaVarName := name.Value()

		// Parse type
		if gv.TypeCpIndex >= 0 {
			t = parseTypeFromCP(b, gv.TypeCpIndex)
		}

		g := BIRGlobalVariableDcl{
			BIRVariableDcl: BIRVariableDcl{
				BIRDocumentableNodeBase: BIRDocumentableNodeBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
				},
				Type:         t,
				Name:         name,
				OriginalName: name,
				MetaVarName:  metaVarName,
				Scope:        scope,
				Kind:         kind,
			},
			Flags:  gv.Flags,
			PkgId:  pkgID,
			Origin: origin,
		}
		// Parse markdown doc attachment
		if gv.Doc != nil {
			fmt.Println("WARNING: markdown doc attachment not supported ignoring")
		}

		// Parse annotation attachments
		if gv.AnnotationAttachmentsContent != nil {
			fmt.Println("WARNING: annotation attachments not supported ignoring")
		}

		globals = append(globals, g)
	}

	pkg.GlobalVars = globals
	return nil
}

// populateAnnotations maps Bir_Annotation -> BIRAnnotation.
func populateAnnotations(b *Bir, pkg *BIRPackage) error {
	if b.Module.AnnotationsSize == 0 {
		return nil
	}
	panic("annotations not supported")
}

// populateServices maps Bir_ServiceDeclaration -> BIRServiceDeclaration.
func populateServices(b *Bir, pkg BIRPackage) error {
	if b.Module.ServiceDeclsSize == 0 {
		return nil
	}
	panic("services not supported")
}

// --- Constant‑pool helpers ---------------------------------------------------

// cpEntry returns the constant‑pool entry at idx.
func cpEntry(b *Bir, idx int32) (*Bir_ConstantPoolEntry, error) {
	if idx < 0 || int(idx) >= len(b.ConstantPool.ConstantPoolEntries) {
		return nil, fmt.Errorf("cp index out of bounds: %d", idx)
	}
	return b.ConstantPool.ConstantPoolEntries[idx], nil
}

// cpString resolves a string CP index to its value. Empty string is returned
// if the entry is not a string.
func cpString(b *Bir, idx int32) string {
	if idx < 0 {
		return ""
	}
	e, err := cpEntry(b, idx)
	if err != nil || e == nil {
		return ""
	}
	if s, ok := e.CpInfo.(*Bir_StringCpInfo); ok && s != nil {
		return s.Value
	}
	return ""
}

// cpAsPackage resolves a package CP entry.
func cpAsPackage(b *Bir, idx int32) (*Bir_PackageCpInfo, error) {
	e, err := cpEntry(b, idx)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("nil package cp entry at %d", idx)
	}
	pkg, ok := e.CpInfo.(*Bir_PackageCpInfo)
	if !ok || pkg == nil {
		return nil, fmt.Errorf("cp entry at %d is not a package", idx)
	}
	return pkg, nil
}

// populateFunctionBody populates the function body including basic blocks and instructions.
func populateFunctionBody(b *Bir, fn *BIRFunction, body *Bir_FunctionBody) error {
	if body == nil {
		return nil
	}

	// Args count
	fn.ArgsCount = int(body.ArgsCount)

	// Return variable
	if body.HasReturnVar != 0 && body.ReturnVar != nil {
		returnVar := parseReturnVar(b, body.ReturnVar)
		fn.ReturnVariable = returnVar
	}

	// Function parameters (default parameters)
	if body.DefaultParameterCount > 0 && len(body.DefaultParameters) > 0 {
		params := make([]BIRFunctionParameter, 0, len(body.DefaultParameters))
		for _, dp := range body.DefaultParameters {
			if dp == nil {
				continue
			}
			param := parseFunctionParameter(b, dp)
			if param != nil {
				params = append(params, *param)
			}
		}
		if len(params) > 0 {
			fn.Parameters = params
		}
	}

	// Local variables
	if body.LocalVariablesCount > 0 && len(body.LocalVariables) > 0 {
		localVars := make([]BIRVariableDcl, 0, len(body.LocalVariables))
		for _, lv := range body.LocalVariables {
			if lv == nil {
				continue
			}
			localVar := parseLocalVariable(b, lv)
			if localVar != nil {
				localVars = append(localVars, *localVar)
			}
		}
		if len(localVars) > 0 {
			fn.LocalVars = localVars
		}
	}

	// Populate basic blocks
	if body.FunctionBasicBlocksInfo != nil && body.FunctionBasicBlocksInfo.BasicBlocksCount > 0 {
		basicBlocks, err := populateBasicBlocks(b, body.FunctionBasicBlocksInfo)
		if err != nil {
			return fmt.Errorf("populating basic blocks: %w", err)
		}
		fn.BasicBlocks = basicBlocks
	}

	// Error table
	if body.ErrorTable != nil && body.ErrorTable.ErrorEntriesCount > 0 {
		panic("error table not supported")
	}

	// Worker channels
	if body.WorkerChannelInfo != nil && body.WorkerChannelInfo.ChannelsLength > 0 {
		panic("worker channel not supported")
	}

	return nil
}

// populateBasicBlocks creates BIRBasicBlock instances from the Kaitai model.
func populateBasicBlocks(b *Bir, bbInfo *Bir_BasicBlocksInfo) ([]BIRBasicBlock, error) {
	if bbInfo == nil || bbInfo.BasicBlocksCount == 0 {
		return []BIRBasicBlock{}, nil
	}

	basicBlocks := make([]*BIRBasicBlock, 0, bbInfo.BasicBlocksCount)
	bbMap := make(map[string]*BIRBasicBlock) // Map BB name to BB for terminator references

	// First pass: create all basic blocks
	for i, kaitaiBB := range bbInfo.BasicBlocks {
		if kaitaiBB == nil {
			continue
		}

		bbName := model.Name(cpString(b, kaitaiBB.NameCpIndex))
		bb := BIRBasicBlock{
			Number: i,
			Id:     bbName,
		}
		bbMap[bbName.Value()] = &bb
		basicBlocks = append(basicBlocks, &bb)
	}

	// Second pass: populate instructions for each basic block
	for i, kaitaiBB := range bbInfo.BasicBlocks {
		if kaitaiBB == nil || i >= len(basicBlocks) {
			continue
		}

		bb := basicBlocks[i]

		if kaitaiBB.InstructionsCount == 0 || len(kaitaiBB.Instructions) == 0 {
			// Empty basic block - no instructions or terminator
			continue
		}

		// Separate instructions into non-terminators and terminator
		// The last instruction is always the terminator (per BIRBinaryWriter logic)
		nonTerminators := make([]BIRNonTerminator, 0)
		var terminator BIRTerminator

		instructionCount := len(kaitaiBB.Instructions)
		if instructionCount == 0 {
			continue
		}

		// Process all instructions except the last one as non-terminators
		for j := 0; j < instructionCount-1; j++ {
			kaitaiIns := kaitaiBB.Instructions[j]
			if kaitaiIns == nil {
				continue
			}

			kind := InstructionKind(kaitaiIns.InstructionKind)
			pos := positionToLocation(b, kaitaiIns.Position)

			// Verify this is actually a non-terminator (kind > 16)
			// If somehow a terminator appears before the last instruction, skip it
			if kind > 16 {
				nonTerm := createNonTerminator(b, kind, pos, kaitaiIns)
				if nonTerm != nil {
					nonTerminators = append(nonTerminators, nonTerm)
				}
			}
		}

		// The last instruction is the terminator
		lastIns := kaitaiBB.Instructions[instructionCount-1]
		if lastIns != nil {
			kind := InstructionKind(lastIns.InstructionKind)
			pos := positionToLocation(b, lastIns.Position)
			term := createTerminator(b, kind, pos, lastIns, bbMap)
			if term != nil {
				terminator = term
			}
		}

		// Set instructions and terminator on the basic block
		if len(nonTerminators) > 0 {
			bb.Instructions = nonTerminators
		}
		if terminator != nil {
			bb.Terminator = terminator
		}
	}

	basicBlocksList := make([]BIRBasicBlock, 0, len(basicBlocks))
	for _, bb := range basicBlocks {
		basicBlocksList = append(basicBlocksList, *bb)
	}

	return basicBlocksList, nil
}

// createTerminator creates a BIRTerminator instance from Kaitai instruction data.
func createTerminator(b *Bir, kind InstructionKind, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap bbMap) BIRTerminator {
	if kaitaiIns == nil || kaitaiIns.InstructionStructure == nil {
		// Fallback to minimal implementation
		return nil
	}

	// Parse based on instruction kind
	switch kind {
	case INSTRUCTION_KIND_GOTO:
		return parseGotoTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_RETURN:
		return parseReturnTerminator(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_BRANCH:
		return parseBranchTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_CALL:
		return parseCallTerminator(b, pos, kaitaiIns, bbMap)
	default:
		panic(fmt.Sprintf("unknown terminator kind: %d", kind))
	}
}

// parseGotoTerminator parses a GOTO terminator
func parseGotoTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap bbMap) BIRTerminator {
	gotoIns := kaitaiIns.InstructionStructure.(*Bir_InstructionGoto)
	targetBBName := cpString(b, gotoIns.TargetBbIdNameCpIndex)
	var targetBB *BIRBasicBlock
	if targetBBName != "" {
		targetBB = bbMap[targetBBName]
	}
	return &Goto{
		BIRTerminatorBase: BIRTerminatorBase{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
			},
			ThenBB: targetBB,
		},
	}
}

// parseReturnTerminator parses a Return terminator
func parseReturnTerminator(_ *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRTerminator {
	if _, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionReturn); ok {
		return &Return{
			BIRTerminatorBase: BIRTerminatorBase{
				BIRInstructionBase: BIRInstructionBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
				},
			},
		}
	}
	panic("unexpected")
}

// parseBranchTerminator parses a Branch terminator
func parseBranchTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap bbMap) BIRTerminator {
	if branchIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionBranch); ok && branchIns != nil {
		op := parseOperand(b, branchIns.BranchOperand)
		trueBBName := cpString(b, branchIns.TrueBbIdNameCpIndex)
		falseBBName := cpString(b, branchIns.FalseBbIdNameCpIndex)
		var trueBB, falseBB *BIRBasicBlock
		if trueBBName != "" {
			trueBB = bbMap[trueBBName]
		}
		if falseBBName != "" {
			falseBB = bbMap[falseBBName]
		}
		return &Branch{
			BIRTerminatorBase: BIRTerminatorBase{
				BIRInstructionBase: BIRInstructionBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
				},
			},
			Op:      op,
			TrueBB:  trueBB,
			FalseBB: falseBB,
		}
	}
	panic("unexpected")
}

// parseCallTerminator parses a Call terminator
func parseCallTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap bbMap) BIRTerminator {
	if callIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionCall); ok && callIns != nil {
		callInfo := callIns.CallInstructionInfo
		if callInfo == nil {
			return nil
		}

		isVirtual := callInfo.IsVirtual != 0
		// Parse package ID
		var calleePkg model.PackageID
		if callInfo.PackageIndex >= 0 {
			pkgCp, err := cpAsPackage(b, callInfo.PackageIndex)
			if err == nil && pkgCp != nil {
				org := model.Name(cpString(b, pkgCp.OrgIndex))
				pkgName := model.Name(cpString(b, pkgCp.PackageNameIndex))
				namePart := model.Name(cpString(b, pkgCp.NameIndex))
				version := model.Name(cpString(b, pkgCp.VersionIndex))
				calleePkg = model.NewPackageID(org, []model.Name{pkgName, namePart}, version)
			}
		}

		name := model.Name(cpString(b, callInfo.CallNameCpIndex))
		args := make([]BIROperand, 0, len(callInfo.Arguments))
		for _, arg := range callInfo.Arguments {
			if arg != nil {
				args = append(args, *parseOperand(b, arg))
			}
		}

		var lhsOp *BIROperand
		if callInfo.HasLhsOperand != 0 && callInfo.LhsOperand != nil {
			lhsOp = parseOperand(b, callInfo.LhsOperand)
		}

		thenBBName := cpString(b, callIns.ThenBbIdNameCpIndex)
		var thenBB *BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}

		return &Call{
			BIRTerminatorBase: BIRTerminatorBase{
				BIRInstructionBase: BIRInstructionBase{
					BIRNodeBase: BIRNodeBase{
						Pos: pos,
					},
					LhsOp: lhsOp,
				},
				ThenBB: thenBB,
			},
			Kind:      INSTRUCTION_KIND_CALL,
			IsVirtual: isVirtual,
			CalleePkg: calleePkg,
			Name:      name,
			Args:      args,
		}
	}
	panic("unexpected")
}

// createNonTerminator creates a BIRNonTerminator instance from Kaitai instruction data.
func createNonTerminator(b *Bir, kind InstructionKind, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if kaitaiIns == nil || kaitaiIns.InstructionStructure == nil {
		// Fallback to minimal implementation
		return nil
	}

	// Parse based on instruction kind
	switch kind {
	case INSTRUCTION_KIND_MOVE:
		return parseMoveInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_CONST_LOAD:
		return parseConstantLoadInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_ADD, INSTRUCTION_KIND_SUB, INSTRUCTION_KIND_MUL, INSTRUCTION_KIND_DIV, INSTRUCTION_KIND_MOD,
		INSTRUCTION_KIND_EQUAL, INSTRUCTION_KIND_NOT_EQUAL, INSTRUCTION_KIND_GREATER_THAN, INSTRUCTION_KIND_GREATER_EQUAL,
		INSTRUCTION_KIND_LESS_THAN, INSTRUCTION_KIND_LESS_EQUAL, INSTRUCTION_KIND_AND, INSTRUCTION_KIND_OR,
		INSTRUCTION_KIND_REF_EQUAL, INSTRUCTION_KIND_REF_NOT_EQUAL, INSTRUCTION_KIND_CLOSED_RANGE, INSTRUCTION_KIND_HALF_OPEN_RANGE,
		INSTRUCTION_KIND_ANNOT_ACCESS:
		return parseBinaryOpInstruction(b, pos, kind, kaitaiIns)
	case INSTRUCTION_KIND_TYPEOF, INSTRUCTION_KIND_NOT, INSTRUCTION_KIND_NEGATE:
		return parseUnaryOpInstruction(b, pos, kind, kaitaiIns)
	default:
		return nil
	}
}

// parseMoveInstruction parses a Move instruction
func parseMoveInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if moveIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionMove); ok && moveIns != nil {
		rhsOp := parseOperand(b, moveIns.RhsOperand)
		lhsOp := parseOperand(b, moveIns.LhsOperand)
		return &Move{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				LhsOp: lhsOp,
			},
			RhsOp: rhsOp,
		}
	}
	panic("unexpected")
}

// parseConstantLoadInstruction parses a ConstantLoad instruction
func parseConstantLoadInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if constIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionConstLoad); ok && constIns != nil {
		lhsOp := parseOperand(b, constIns.LhsOperand)
		var constType model.ValueType
		if constIns.TypeCpIndex >= 0 {
			constType = parseTypeFromCP(b, constIns.TypeCpIndex)
		}
		// Parse constant value
		var value any
		if constIns.ConstantValueInfo != nil {
			value = parseConstantValueInfo(b, constIns.ConstantValueInfo, constType)
		}
		return &ConstantLoad{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				LhsOp: lhsOp,
			},
			Value: value,
			Type:  constType,
		}
	}
	panic("unexpected")
}

// parseBinaryOpInstruction parses a BinaryOp instruction
func parseBinaryOpInstruction(b *Bir, pos diagnostics.Location, kind InstructionKind, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if binOpIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionBinaryOperation); ok && binOpIns != nil {
		rhsOp1 := parseOperand(b, binOpIns.RhsOperandOne)
		rhsOp2 := parseOperand(b, binOpIns.RhsOperandTwo)
		lhsOp := parseOperand(b, binOpIns.LhsOperand)
		return &BinaryOp{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				LhsOp: lhsOp,
			},
			Kind:   kind,
			RhsOp1: *rhsOp1,
			RhsOp2: *rhsOp2,
		}
	}
	panic("unexpected")
}

// parseUnaryOpInstruction parses a UnaryOp instruction
func parseUnaryOpInstruction(b *Bir, pos diagnostics.Location, kind InstructionKind, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if unaryOpIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionUnaryOperation); ok && unaryOpIns != nil {
		rhsOp := parseOperand(b, unaryOpIns.RhsOperand)
		lhsOp := parseOperand(b, unaryOpIns.LhsOperand)
		return &UnaryOp{
			BIRInstructionBase: BIRInstructionBase{
				BIRNodeBase: BIRNodeBase{
					Pos: pos,
				},
				LhsOp: lhsOp,
			},
			Kind:  kind,
			RhsOp: rhsOp,
		}
	}
	panic("unexpected")
}

// positionToLocation converts a Bir_Position to diagnostics.Location.
// For now, returns nil as we don't have a full Location implementation.
func positionToLocation(_ *Bir, pos *Bir_Position) diagnostics.Location {
	if pos == nil {
		return nil
	}
	// TODO: Implement full Location conversion from Bir_Position
	return nil
}

// Minimal implementations to satisfy interfaces

type terminatorImpl struct {
	BIRTerminatorBase
	kind InstructionKind
}

func (t *terminatorImpl) GetNextBasicBlocks() []BIRBasicBlock {
	if t.ThenBB != nil {
		return []BIRBasicBlock{*t.ThenBB}
	}
	return nil
}

// populatePathParameters populates path parameters for resource functions.
func populatePathParameters(b *Bir, fn BIRFunction, rfc *Bir_ResourceFunctionContent) error {
	if rfc == nil {
		return nil
	}

	panic("path parameters not supported")
}

// parseAnnotationAttachments parses annotation attachments.
func parseAnnotationAttachments(b *Bir, aac *Bir_AnnotationAttachmentsContent) []any {
	if aac == nil || aac.AttachmentsCount == 0 {
		return nil
	}

	panic("annotation attachments not supported")
}

// parseReturnVar parses a return variable.
func parseReturnVar(b *Bir, rv *Bir_ReturnVar) *BIRVariableDcl {
	if rv == nil {
		return nil
	}

	name := model.Name(cpString(b, rv.NameCpIndex))
	kind := VarKind(rv.Kind)
	var t model.ValueType
	if rv.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, rv.TypeCpIndex)
	}
	return &BIRVariableDcl{
		Type:  t,
		Name:  name,
		Scope: VAR_SCOPE_FUNCTION,
		Kind:  kind,
	}
}

// parseFunctionParameter parses a function parameter (default parameter).
func parseFunctionParameter(b *Bir, dp *Bir_DefaultParameter) *BIRFunctionParameter {
	if dp == nil {
		return nil
	}

	name := model.Name(cpString(b, dp.NameCpIndex))
	kind := VarKind(dp.Kind)
	var t model.ValueType
	if dp.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, dp.TypeCpIndex)
	}
	metaVarName := cpString(b, dp.MetaVarNameCpIndex)
	hasDefaultExpr := dp.HasDefaultExpr != 0
	return &BIRFunctionParameter{
		BIRVariableDcl: BIRVariableDcl{
			Type:         t,
			Name:         name,
			OriginalName: name,
			MetaVarName:  metaVarName,
			Scope:        VAR_SCOPE_FUNCTION,
			Kind:         kind,
		},
		HasDefaultExpr: hasDefaultExpr,
	}
}

// parseLocalVariable parses a local variable.
func parseLocalVariable(b *Bir, lv *Bir_LocalVariable) *BIRVariableDcl {
	if lv == nil {
		return nil
	}

	name := model.Name(cpString(b, lv.NameCpIndex))
	kind := VarKind(lv.Kind)
	var t model.ValueType
	if lv.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, lv.TypeCpIndex)
	}
	localVar := &BIRVariableDcl{
		Type:  t,
		Name:  name,
		Scope: VAR_SCOPE_FUNCTION,
		Kind:  kind,
	}

	// Parse enclosing basic block info if this is a LOCAL variable
	if kind == VAR_KIND_LOCAL && lv.EnclosingBasicBlockId != nil {
		// TODO: Set startBB, endBB, and insOffset if needed
		// This requires access to the basic blocks map
	}

	return localVar
}

// Helper functions for parsing complex structures

// parseTypeFromCP parses a BType from a constant pool index (shape_cp_info).
func parseTypeFromCP(b *Bir, cpIndex int32) model.ValueType {
	if cpIndex < 0 {
		return nil
	}

	e, err := cpEntry(b, cpIndex)
	if err != nil || e == nil {
		return nil
	}

	// Check if it's a shape CP entry
	shapeCp, ok := e.CpInfo.(*Bir_ShapeCpInfo)
	if !ok || shapeCp == nil || shapeCp.Shape == nil {
		return nil
	}

	// Create a minimal BType implementation from the type info
	return createBTypeFromTypeInfo(b, shapeCp.Shape)
}

// createBTypeFromTypeInfo creates a minimal BType from Bir_TypeInfo.
// This is a placeholder implementation until full type system is available.
func createBTypeFromTypeInfo(b *Bir, ti *Bir_TypeInfo) model.ValueType {
	if ti == nil {
		return nil
	}

	// Create a minimal type implementation
	// The actual type structure is complex and would require full type system
	// For now, we create a stub that satisfies the BType interface
	return &minimalBType{
		tag:   int(ti.TypeTag),
		name:  model.Name(cpString(b, ti.NameIndex)),
		flags: ti.TypeFlag,
	}
}

// minimalBType is a minimal implementation of BType for parsing purposes.
type minimalBType struct {
	tag     int
	name    model.Name
	flags   int64
	semType any
	tsymbol any
}

func (t *minimalBType) GetTag() int                    { return t.tag }
func (t *minimalBType) SetTag(tag int)                 { t.tag = tag }
func (t *minimalBType) GetName() model.Name            { return t.name }
func (t *minimalBType) SetName(name model.Name)        { t.name = name }
func (t *minimalBType) GetFlags() int64                { return t.flags }
func (t *minimalBType) SetFlags(flags int64)           { t.flags = flags }
func (t *minimalBType) GetSemType() any                { return t.semType }
func (t *minimalBType) SetSemType(semType any)         { t.semType = semType }
func (t *minimalBType) GetTsymbol() any                { return t.tsymbol }
func (t *minimalBType) SetTsymbol(tsymbol any)         { t.tsymbol = tsymbol }
func (t *minimalBType) GetReturnType() model.ValueType { return nil }
func (t *minimalBType) GetTypeKind() model.TypeKind {
	// Return a default type kind - this is a minimal implementation
	return model.TypeKind_OTHER
}

func (t *minimalBType) String() string {
	if t.name.Value() != "" {
		return t.name.Value()
	}
	return fmt.Sprintf("type_%d", t.tag)
}

// parseMarkdown parses markdown documentation from Bir_Markdown.
func parseMarkdown(b *Bir, md *Bir_Markdown) model.MarkdownDocAttachment {
	panic("markdown not supported")
}

// parseConstantValue parses a constant value from Bir_ConstantValue.
func parseConstantValue(b *Bir, cv *Bir_ConstantValue) ConstValue {
	if cv == nil {
		return ConstValue{}
	}

	var constType model.ValueType
	if cv.ConstantValueTypeCpIndex >= 0 {
		constType = parseTypeFromCP(b, cv.ConstantValueTypeCpIndex)
	}

	// Parse constant value info based on type
	var value any
	if cv.ConstantValueInfo != nil {
		// The constant value info is a switch-on type structure
		// For now, we extract basic values based on the type
		// Full parsing would require understanding all constant value types
		value = parseConstantValueInfo(b, cv.ConstantValueInfo, constType)
	}

	return ConstValue{
		Type:  constType,
		Value: value,
	}
}

// parseConstantValueInfo parses the actual constant value from the info structure.
func parseConstantValueInfo(b *Bir, cvi kaitai.Struct, constType model.ValueType) any {
	if cvi == nil {
		return nil
	}

	// Handle different constant value types based on type assertion
	switch v := cvi.(type) {
	case *Bir_IntConstantInfo:
		if v.ValueCpIndex >= 0 {
			// Get integer value from CP
			e, err := cpEntry(b, v.ValueCpIndex)
			if err == nil && e != nil {
				if intCp, ok := e.CpInfo.(*Bir_IntCpInfo); ok && intCp != nil {
					return intCp.Value
				}
			}
		}
	case *Bir_StringConstantInfo:
		if v.ValueCpIndex >= 0 {
			return cpString(b, v.ValueCpIndex)
		}
	case *Bir_BooleanConstantInfo:
		return v.ValueBooleanConstant != 0
	case *Bir_FloatConstantInfo:
		if v.ValueCpIndex >= 0 {
			// Get float value from CP
			e, err := cpEntry(b, v.ValueCpIndex)
			if err == nil && e != nil {
				if floatCp, ok := e.CpInfo.(*Bir_FloatCpInfo); ok && floatCp != nil {
					return floatCp.Value
				}
			}
		}
	case *Bir_ByteConstantInfo:
		if v.ValueCpIndex >= 0 {
			// Get byte value from CP
			e, err := cpEntry(b, v.ValueCpIndex)
			if err == nil && e != nil {
				if byteCp, ok := e.CpInfo.(*Bir_ByteCpInfo); ok && byteCp != nil {
					return byteCp.Value
				}
			}
		}
	case *Bir_DecimalConstantInfo:
		if v.ValueCpIndex >= 0 {
			// Decimal values are stored as CP indices - would need decimal parsing
			// For now, return the CP index as a placeholder
			return v.ValueCpIndex
		}
	case *Bir_NilConstantInfo:
		return nil
	case *Bir_MapConstantInfo:
		// Map constants - parse key-value pairs
		if v.MapConstantSize > 0 && len(v.MapKeyValues) > 0 {
			result := make(map[string]any)
			for _, kv := range v.MapKeyValues {
				if kv == nil {
					continue
				}
				keyName := cpString(b, kv.KeyNameCpIndex)
				keyValue := parseConstantValueInfo(b, kv.KeyValueInfo, nil)
				result[keyName] = keyValue
			}
			return result
		}
	case *Bir_ListConstantInfo:
		// List constants - parse member values
		if v.ListConstantSize > 0 && len(v.ListMemberValueInfo) > 0 {
			result := make([]any, 0, len(v.ListMemberValueInfo))
			for _, member := range v.ListMemberValueInfo {
				if member == nil {
					continue
				}
				memberValue := parseConstantValue(b, member)
				result = append(result, memberValue.Value)
			}
			return result
		}
	case *Bir_IntersectionConstantInfo:
		// Intersection constants - recursively parse the underlying constant value
		if v.ConstantValueInfo != nil {
			return parseConstantValueInfo(b, v.ConstantValueInfo, constType)
		}
	}

	return nil
}

// parseReceiver parses a receiver from Bir_Reciever.
func parseReceiver(b *Bir, rec *Bir_Reciever) *BIRVariableDcl {
	if rec == nil {
		return nil
	}

	name := model.Name(cpString(b, rec.NameCpIndex))
	kind := VarKind(rec.Kind)
	var t model.ValueType
	if rec.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, rec.TypeCpIndex)
	}

	return &BIRVariableDcl{
		Type:  t,
		Name:  name,
		Scope: VAR_SCOPE_FUNCTION,
		Kind:  kind,
	}
}

// parseOperand parses a BIROperand from Bir_Operand.
func parseOperand(b *Bir, op *Bir_Operand) *BIROperand {
	if op == nil {
		return nil
	}

	// Check if variable is ignored
	if op.IgnoredVariable != 0 {
		// Ignored variable - create a minimal variable declaration
		var ignoredType model.ValueType
		if op.IgnoredTypeCpIndex >= 0 {
			ignoredType = parseTypeFromCP(b, op.IgnoredTypeCpIndex)
		}
		// Create a minimal variable for ignored operands
		ignoredVar := &BIRVariableDcl{
			Type:  ignoredType,
			Name:  model.Name("_"),
			Scope: VAR_SCOPE_FUNCTION,
			Kind:  VAR_KIND_LOCAL,
		}
		return &BIROperand{
			VariableDcl: ignoredVar,
		}
	}

	// Parse variable
	if op.Variable == nil {
		return nil
	}

	var varType model.ValueType
	varName := model.Name(cpString(b, op.Variable.VariableDclNameCpIndex))
	kind := VarKind(op.Variable.Kind)
	scope := VarScope(op.Variable.Scope)

	// For global/constant variables, get type from GlobalOrConstantVariable
	if op.Variable.Kind == 5 || op.Variable.Kind == 7 {
		if op.Variable.GlobalOrConstantVariable != nil {
			if op.Variable.GlobalOrConstantVariable.TypeCpIndex >= 0 {
				varType = parseTypeFromCP(b, op.Variable.GlobalOrConstantVariable.TypeCpIndex)
			}
		}
	}

	// Create variable declaration
	varDcl := &BIRVariableDcl{
		Type:  varType,
		Name:  varName,
		Scope: scope,
		Kind:  kind,
	}
	return &BIROperand{
		VariableDcl: varDcl,
	}
}
