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

	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func LoadBIRPackageFromReader(r io.Reader) (BIRPackage, error) {
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
	if err := populateServices(b, pkg); err != nil {
		return nil, err
	}

	// TODO: types, globals, annotations, services etc. can be wired in the
	// same style as needed. For now we keep them empty – the model remains
	// structurally valid and the code builds.

	return pkg, nil
}

// buildBIRPackage constructs the BIRPackage shell: only PackageID and empty
// collections are initialized here; contents are populated by the helpers
// below.
func buildBIRPackage(b *Bir) (BIRPackage, error) {
	mod := b.Module

	// The module id is a CP index into the constant pool; the corresponding
	// entry is a PackageCP entry.
	pkgCp, err := cpAsPackage(b, mod.IdCpIndex)
	if err != nil {
		return nil, fmt.Errorf("reading module package id: %w", err)
	}

	org := Name(cpString(b, pkgCp.OrgIndex))
	pkgName := Name(cpString(b, pkgCp.PackageNameIndex))
	name := Name(cpString(b, pkgCp.NameIndex))
	version := Name(cpString(b, pkgCp.VersionIndex))

	// We don't have source file name / root at this level; keep them empty.
	sourceFileName := Name("")
	sourceRoot := ""
	skipTest := false
	isTest := false

	p := NewBIRPackageWithIsTestPkg(
		nil, // pos
		org,
		pkgName,
		name,
		version,
		sourceFileName,
		sourceRoot,
		skipTest,
		isTest,
	)

	return p, nil
}

// populateImports fills BIRPackage.importModules from Module.Imports.
func populateImports(b *Bir, pkg BIRPackage) error {
	if b.Module.ImportCount == 0 {
		return nil
	}

	imports := make([]BIRImportModule, 0, len(b.Module.Imports))

	for _, imp := range b.Module.Imports {
		if imp == nil {
			continue
		}
		org := Name(cpString(b, imp.OrgIndex))
		// Note: PackageNameIndex is available but NewBIRImportModule doesn't support separate pkgName
		// It uses NewPackageIDWithOrgNameVersion which sets PkgName = name
		name := Name(cpString(b, imp.NameIndex))
		version := Name(cpString(b, imp.VersionIndex))

		mod := NewBIRImportModule(nil, org, name, version)
		imports = append(imports, mod)
		// TODO: Use imp.PackageNameIndex if BIRImportModule API is extended to support separate pkgName
	}

	pkg.SetImportModules(&imports)
	return nil
}

// populateFunctions creates a minimal BIRFunction node for each Bir_Function.
func populateFunctions(b *Bir, pkg BIRPackage) error {
	if b.Module.FunctionCount == 0 {
		return nil
	}

	funcs := make([]BIRFunction, 0, len(b.Module.Functions))

	for _, f := range b.Module.Functions {
		if f == nil {
			continue
		}

		name := Name(cpString(b, f.NameCpIndex))
		origName := Name(cpString(b, f.OriginalNameCpIndex))
		workerName := Name(cpString(b, f.WorkerNameCpIndex))
		origin := model.SymbolOrigin(f.Origin)
		pos := positionToLocation(b, f.Position)

		// Type is left nil for now; wiring BType from the shape CP entries is
		// a larger task and not required for the structural mapping.
		var fnType BInvokableType

		// Worker channel count is not directly known at this level; create
		// with zero sendInsCount – this matches the secondary Java ctor used
		// by BIRGen when only structure matters.
		fn := NewBIRFunctionWithSendInsCount(
			pos,
			name,
			origName,
			f.Flags,
			fnType,
			workerName,
			0, // sendInsCount
			origin,
		)

		// Required params -> BIRParameter list.
		if f.RequiredParamCount > 0 {
			params := make([]BIRParameter, 0, len(f.RequiredParams))
			for _, rp := range f.RequiredParams {
				if rp == nil {
					continue
				}
				pName := Name(cpString(b, rp.ParamNameCpIndex))
				p := NewBIRParameter(nil, pName, rp.Flags)
				// TODO: Parse parameter annotations from rp if available
				params = append(params, p)
			}
			fn.SetRequiredParams(&params)
		}

		// Rest parameter
		if f.HasRestParam != 0 {
			restParamName := Name(cpString(b, f.RestParamNameCpIndex))
			restParam := NewBIRParameter(nil, restParamName, 0)
			// Note: BIRParameter doesn't have SetAnnotAttachments in the current model
			// Rest param annotations are available in f.RestParamAnnotations if needed
			fn.SetRestParam(restParam)
		}

		// Receiver
		if f.HasReceiver != 0 && f.Reciever != nil {
			receiver := parseReceiver(b, f.Reciever)
			if receiver != nil {
				fn.SetReceiver(receiver)
			}
		}

		// Path parameters (resource function)
		if f.IsResourceFunction != 0 && f.ResourceFunctionContent != nil {
			if err := populatePathParameters(b, fn, f.ResourceFunctionContent); err != nil {
				return fmt.Errorf("populating path parameters: %w", err)
			}
		}

		// Annotation attachments
		if f.AnnotationAttachmentsContent != nil {
			annots := parseAnnotationAttachments(b, f.AnnotationAttachmentsContent)
			fn.SetAnnotAttachments(&annots)
		}

		// Return type annotations
		if f.ReturnTypeAnnotations != nil {
			annots := parseAnnotationAttachments(b, f.ReturnTypeAnnotations)
			fn.SetReturnTypeAnnots(&annots)
		}

		// Markdown doc attachment
		if f.Doc != nil {
			mdDoc := parseMarkdown(b, f.Doc)
			fn.SetMarkdownDocAttachment(mdDoc)
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
			pkgID := pkg.GetPackageID()
			for _, cpIdx := range f.DependentGlobalVarCpEntry {
				varName := Name(cpString(b, cpIdx))
				// Create a minimal global var reference
				gv := NewBIRGlobalVariableDcl(
					nil,
					0,   // flags
					nil, // type
					pkgID,
					varName,
					varName,
					VAR_SCOPE_GLOBAL,
					VAR_KIND_LOCAL, // kind doesn't matter for dependency reference
					varName.Value(),
					model.SymbolOrigin_SOURCE,
				)
				dependentVars = append(dependentVars, gv)
			}
			fn.SetDependentGlobalVars(&dependentVars)
		}

		// Populate function body (basic blocks, instructions, etc.)
		if f.FunctionBody != nil {
			if err := populateFunctionBody(b, fn, f.FunctionBody); err != nil {
				return fmt.Errorf("populating function body for %s: %w", name.Value(), err)
			}
		}

		// We also set the function's flags/origin/workerName via the ctor.
		funcs = append(funcs, fn)
	}

	pkg.SetFunctions(&funcs)
	return nil
}

// populateConstants maps Bir_Constant -> BIRConstant.
func populateConstants(b *Bir, pkg BIRPackage) error {
	if b.Module.ConstCount == 0 {
		empty := []BIRConstant{}
		pkg.SetConstants(&empty)
		return nil
	}

	consts := make([]BIRConstant, 0, len(b.Module.Constants))

	for _, c := range b.Module.Constants {
		if c == nil {
			continue
		}
		name := Name(cpString(b, c.NameCpIndex))
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

		bc := NewBIRConstant(pos, name, c.Flags, t, cv, origin)

		// Parse markdown doc attachment
		if c.Doc != nil {
			mdDoc := parseMarkdown(b, c.Doc)
			bc.SetMarkdownDocAttachment(mdDoc)
		}

		// Parse annotation attachments
		if c.AnnotationAttachmentsContent != nil {
			annots := parseAnnotationAttachments(b, c.AnnotationAttachmentsContent)
			if len(annots) > 0 {
				bc.SetAnnotAttachments(&annots)
			}
		}

		consts = append(consts, bc)
	}

	pkg.SetConstants(&consts)
	return nil
}

// populateTypeDefs maps Bir_TypeDefinition -> BIRTypeDefinition.
func populateTypeDefs(b *Bir, pkg BIRPackage) error {
	if b.Module.TypeDefinitionCount == 0 {
		empty := []BIRTypeDefinition{}
		pkg.SetTypeDefs(&empty)
		return nil
	}

	defs := make([]BIRTypeDefinition, 0, len(b.Module.TypeDefinitions))

	for _, td := range b.Module.TypeDefinitions {
		if td == nil {
			continue
		}

		name := Name(cpString(b, td.NameCpIndex))
		origName := Name(cpString(b, td.OriginalNameCpIndex))
		origin := model.SymbolOrigin(td.Origin)
		pos := positionToLocation(b, td.Position)

		// Internal name is the same as name for now.
		internalName := name

		isBuiltin := false
		var t model.ValueType
		attachedFuncs := []BIRFunction{}

		// Parse type
		if td.TypeCpIndex >= 0 {
			t = parseTypeFromCP(b, td.TypeCpIndex)
		}

		bt := NewBIRTypeDefinition(
			pos,
			internalName,
			td.Flags,
			isBuiltin,
			t,
			attachedFuncs,
			origin,
			name,
			origName,
		)

		// Parse markdown doc attachment
		if td.Doc != nil {
			mdDoc := parseMarkdown(b, td.Doc)
			bt.SetMarkdownDocAttachment(mdDoc)
		}

		// Parse annotation attachments
		if td.AnnotationAttachmentsContent != nil {
			annots := parseAnnotationAttachments(b, td.AnnotationAttachmentsContent)
			if len(annots) > 0 {
				bt.SetAnnotAttachments(&annots)
			}
		}

		// Parse reference type (if hasReferenceType is true, it's written separately)
		// For now, we don't parse it as it requires understanding the type structure

		defs = append(defs, bt)
	}

	pkg.SetTypeDefs(&defs)
	return nil
}

// populateTypeDefBodies populates attached functions and referenced types for type definitions.
func populateTypeDefBodies(b *Bir, pkg BIRPackage) error {
	if b.Module.TypeDefinitionBodiesCount == 0 {
		return nil
	}

	typeDefs := pkg.GetTypeDefs()
	if typeDefs == nil || len(*typeDefs) == 0 {
		return nil
	}

	// Filter type defs to only OBJECT and RECORD types (as per Java code)
	// For now, we'll match type def bodies to type defs by index
	// This assumes the order matches between TypeDefinitions and TypeDefinitionBodies
	objectRecordDefs := make([]BIRTypeDefinition, 0)
	for _, td := range *typeDefs {
		// TODO: Check if type is OBJECT or RECORD - for now, include all
		objectRecordDefs = append(objectRecordDefs, td)
	}

	if len(b.Module.TypeDefinitionBodies) != len(objectRecordDefs) {
		// Mismatch - skip for now
		return nil
	}

	for i, tdb := range b.Module.TypeDefinitionBodies {
		if tdb == nil || i >= len(objectRecordDefs) {
			continue
		}

		td := objectRecordDefs[i]

		// Parse attached functions - these are full Bir_Function objects
		if tdb.AttachedFunctionsCount > 0 && len(tdb.AttachedFunctions) > 0 {
			attachedFuncs := make([]BIRFunction, 0, len(tdb.AttachedFunctions))
			for _, af := range tdb.AttachedFunctions {
				if af == nil {
					continue
				}
				// Parse as a full function (same as regular functions)
				funcName := Name(cpString(b, af.NameCpIndex))
				origName := Name(cpString(b, af.OriginalNameCpIndex))
				workerName := Name(cpString(b, af.WorkerNameCpIndex))
				origin := model.SymbolOrigin(af.Origin)
				pos := positionToLocation(b, af.Position)
				var fnType BInvokableType

				attachedFn := NewBIRFunctionWithSendInsCount(
					pos,
					funcName,
					origName,
					af.Flags,
					fnType,
					workerName,
					0, // sendInsCount
					origin,
				)

				// Parse required params, rest param, etc. (same as regular functions)
				if af.RequiredParamCount > 0 {
					params := make([]BIRParameter, 0, len(af.RequiredParams))
					for _, rp := range af.RequiredParams {
						if rp == nil {
							continue
						}
						pName := Name(cpString(b, rp.ParamNameCpIndex))
						p := NewBIRParameter(nil, pName, rp.Flags)
						// Parse parameter annotations
						// Note: BIRParameter doesn't have SetAnnotAttachments in the current model
						// This would need to be added if parameter annotations are needed
						_ = rp.ParamAnnotations
						params = append(params, p)
					}
					attachedFn.SetRequiredParams(&params)
				}

				if af.HasRestParam != 0 {
					restParamName := Name(cpString(b, af.RestParamNameCpIndex))
					restParam := NewBIRParameter(nil, restParamName, 0)
					attachedFn.SetRestParam(restParam)
				}

				// Parse function type
				if af.TypeCpIndex >= 0 {
					fnType := parseTypeFromCP(b, af.TypeCpIndex)
					if invokableType, ok := fnType.(BInvokableType); ok {
						attachedFn.SetType(invokableType)
					}
				}

				// Parse annotation attachments
				if af.AnnotationAttachmentsContent != nil {
					annots := parseAnnotationAttachments(b, af.AnnotationAttachmentsContent)
					if len(annots) > 0 {
						attachedFn.SetAnnotAttachments(&annots)
					}
				}

				// Parse return type annotations
				if af.ReturnTypeAnnotations != nil {
					annots := parseAnnotationAttachments(b, af.ReturnTypeAnnotations)
					if len(annots) > 0 {
						attachedFn.SetReturnTypeAnnots(&annots)
					}
				}

				// Parse markdown doc
				if af.Doc != nil {
					mdDoc := parseMarkdown(b, af.Doc)
					attachedFn.SetMarkdownDocAttachment(mdDoc)
				}

				// Parse function body if available
				if af.FunctionBody != nil {
					if err := populateFunctionBody(b, attachedFn, af.FunctionBody); err != nil {
						// Log error but continue
						continue
					}
				}

				attachedFuncs = append(attachedFuncs, attachedFn)
			}
			if len(attachedFuncs) > 0 {
				td.SetAttachedFuncs(&attachedFuncs)
			}
		}

		// Parse referenced types
		if tdb.ReferencedTypesCount > 0 && len(tdb.ReferencedTypes) > 0 {
			referencedTypes := make([]model.ValueType, 0, len(tdb.ReferencedTypes))
			for _, rt := range tdb.ReferencedTypes {
				if rt == nil {
					continue
				}
				if rt.TypeCpIndex >= 0 {
					refType := parseTypeFromCP(b, rt.TypeCpIndex)
					if refType != nil {
						referencedTypes = append(referencedTypes, refType)
					}
				}
			}
			if len(referencedTypes) > 0 {
				td.SetReferencedTypes(&referencedTypes)
			}
		}
	}

	return nil
}

// populateGlobals maps Bir_GlobalVar -> BIRGlobalVariableDcl.
func populateGlobals(b *Bir, pkg BIRPackage) error {
	if b.Module.GlobalVarCount == 0 {
		empty := []BIRGlobalVariableDcl{}
		pkg.SetGlobalVars(&empty)
		return nil
	}

	globals := make([]BIRGlobalVariableDcl, 0, len(b.Module.GlobalVars))

	// Use the package's PackageID for all globals.
	pkgID := pkg.GetPackageID()

	for _, gv := range b.Module.GlobalVars {
		if gv == nil {
			continue
		}
		name := Name(cpString(b, gv.NameCpIndex))
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

		g := NewBIRGlobalVariableDcl(
			pos,
			gv.Flags,
			t,
			pkgID,
			name,
			name,
			scope,
			kind,
			metaVarName,
			origin,
		)

		// Parse markdown doc attachment
		if gv.Doc != nil {
			mdDoc := parseMarkdown(b, gv.Doc)
			g.SetMarkdownDocAttachment(mdDoc)
		}

		// Parse annotation attachments
		if gv.AnnotationAttachmentsContent != nil {
			annots := parseAnnotationAttachments(b, gv.AnnotationAttachmentsContent)
			if len(annots) > 0 {
				g.SetAnnotAttachments(&annots)
			}
		}

		globals = append(globals, g)
	}

	pkg.SetGlobalVars(&globals)
	return nil
}

// populateAnnotations maps Bir_Annotation -> BIRAnnotation.
func populateAnnotations(b *Bir, pkg BIRPackage) error {
	if b.Module.AnnotationsSize == 0 {
		empty := []BIRAnnotation{}
		pkg.SetAnnotations(&empty)
		return nil
	}

	anns := make([]BIRAnnotation, 0, len(b.Module.Annotations))

	for _, a := range b.Module.Annotations {
		if a == nil {
			continue
		}

		name := Name(cpString(b, a.NameCpIndex))
		origName := Name(cpString(b, a.OriginalNameCpIndex))
		origin := model.SymbolOrigin(a.Origin)
		pos := positionToLocation(b, a.Position)

		// Parse PackageID from PackageIdCpIndex
		var annotPkgID model.PackageID
		if a.PackageIdCpIndex >= 0 {
			pkgCp, err := cpAsPackage(b, a.PackageIdCpIndex)
			if err == nil && pkgCp != nil {
				org := Name(cpString(b, pkgCp.OrgIndex))
				pkgName := Name(cpString(b, pkgCp.PackageNameIndex))
				namePart := Name(cpString(b, pkgCp.NameIndex))
				version := Name(cpString(b, pkgCp.VersionIndex))
				annotPkgID = model.NewPackageID(org, []model.Name{pkgName, namePart}, version)
			}
		}

		// Attach points -> []AttachPoint.
		points := make([]AttachPoint, 0, len(a.AttachPoints))
		for _, ap := range a.AttachPoints {
			if ap == nil {
				continue
			}
			pointName := cpString(b, ap.PointNameCpIndex)
			apStruct := getAttachmentPoint(pointName, ap.IsSource != 0)
			if apStruct != nil {
				points = append(points, *apStruct)
			}
		}

		var t model.ValueType

		ba := NewBIRAnnotation(
			pos,
			name,
			origName,
			a.Flags,
			pointsToSet(points),
			t,
			origin,
		)

		// Set PackageID
		ba.SetPackageID(annotPkgID)

		// Parse annotation type
		if a.AnnotationTypeCpIndex >= 0 {
			annotType := parseTypeFromCP(b, a.AnnotationTypeCpIndex)
			ba.SetAnnotationType(annotType)
		}

		// Parse markdown doc attachment
		if a.Doc != nil {
			mdDoc := parseMarkdown(b, a.Doc)
			ba.SetMarkdownDocAttachment(mdDoc)
		}

		// Parse annotation attachments
		if a.AnnotationAttachmentsContent != nil {
			annots := parseAnnotationAttachments(b, a.AnnotationAttachmentsContent)
			if len(annots) > 0 {
				ba.SetAnnotAttachments(&annots)
			}
		}

		anns = append(anns, ba)
	}

	pkg.SetAnnotations(&anns)
	return nil
}

// populateServices maps Bir_ServiceDeclaration -> BIRServiceDeclaration.
func populateServices(b *Bir, pkg BIRPackage) error {
	if b.Module.ServiceDeclsSize == 0 {
		empty := []BIRServiceDeclaration{}
		pkg.SetServiceDecls(&empty)
		return nil
	}

	services := make([]BIRServiceDeclaration, 0, len(b.Module.ServiceDeclarations))

	for _, s := range b.Module.ServiceDeclarations {
		if s == nil {
			continue
		}

		genName := Name(cpString(b, s.NameCpIndex))
		assocName := Name(cpString(b, s.AssociatedClassNameCpIndex))
		origin := model.SymbolOrigin(s.Origin)
		pos := positionToLocation(b, s.Position)

		// Attach points as strings.
		var attachPoint []string
		if s.HasAttachPoint != 0 && s.AttachPointCount > 0 {
			attachPoint = make([]string, 0, len(s.AttachPoints))
			for _, cpIdx := range s.AttachPoints {
				attachPoint = append(attachPoint, cpString(b, cpIdx))
			}
		}

		var attachPointLiteral string
		if s.HasAttachPointLiteral != 0 {
			attachPointLiteral = cpString(b, s.AttachPointLiteral)
		}

		var t model.ValueType
		if s.HasType != 0 && s.TypeCpIndex >= 0 {
			t = parseTypeFromCP(b, s.TypeCpIndex)
		}

		listenerTypes := make([]model.ValueType, 0, s.ListenerTypesCount)
		for _, lt := range s.ListenerTypes {
			if lt == nil {
				continue
			}
			if lt.TypeCpIndex >= 0 {
				listenerType := parseTypeFromCP(b, lt.TypeCpIndex)
				if listenerType != nil {
					listenerTypes = append(listenerTypes, listenerType)
				}
			}
		}

		svc := NewBIRServiceDeclaration(
			attachPoint,
			attachPointLiteral,
			listenerTypes,
			genName,
			assocName,
			t,
			origin,
			s.Flags,
			pos,
		)

		services = append(services, svc)
	}

	pkg.SetServiceDecls(&services)
	return nil
}

// pointsToSet converts a slice of AttachPoint to a set. Our Go model uses a
// set‑like collection conceptually; we store it as a slice here.
func pointsToSet(points []AttachPoint) []AttachPoint {
	if len(points) == 0 {
		return []AttachPoint{}
	}
	return points
}

// getAttachmentPoint creates an AttachPoint from a point name and source flag.
func getAttachmentPoint(pointName string, isSource bool) *AttachPoint {
	point := model.Point(pointName)
	return &AttachPoint{
		Point:  point,
		Source: isSource,
	}
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
func populateFunctionBody(b *Bir, fn BIRFunction, body *Bir_FunctionBody) error {
	if body == nil {
		return nil
	}

	// Args count
	fn.SetArgsCount(int(body.ArgsCount))

	// Return variable
	if body.HasReturnVar != 0 && body.ReturnVar != nil {
		returnVar := parseReturnVar(b, body.ReturnVar)
		if returnVar != nil {
			fn.SetReturnVariable(returnVar)
		}
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
				params = append(params, param)
			}
		}
		if len(params) > 0 {
			fn.SetParameters(&params)
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
				localVars = append(localVars, localVar)
			}
		}
		if len(localVars) > 0 {
			fn.SetLocalVars(&localVars)
		}
	}

	// Populate basic blocks
	if body.FunctionBasicBlocksInfo != nil && body.FunctionBasicBlocksInfo.BasicBlocksCount > 0 {
		basicBlocks, err := populateBasicBlocks(b, body.FunctionBasicBlocksInfo)
		if err != nil {
			return fmt.Errorf("populating basic blocks: %w", err)
		}
		fn.SetBasicBlocks(&basicBlocks)
	}

	// Error table
	if body.ErrorTable != nil && body.ErrorTable.ErrorEntriesCount > 0 {
		errorEntries := make([]BIRErrorEntry, 0, body.ErrorTable.ErrorEntriesCount)
		for _, ee := range body.ErrorTable.ErrorEntries {
			if ee == nil {
				continue
			}
			errorEntry := parseErrorEntry(b, ee)
			if errorEntry != nil {
				errorEntries = append(errorEntries, errorEntry)
			}
		}
		if len(errorEntries) > 0 {
			fn.SetErrorTable(&errorEntries)
		}
	}

	// Worker channels
	if body.WorkerChannelInfo != nil && body.WorkerChannelInfo.ChannelsLength > 0 {
		channels := make([]ChannelDetails, 0, body.WorkerChannelInfo.ChannelsLength)
		for _, ch := range body.WorkerChannelInfo.WorkerChannelDetails {
			if ch == nil {
				continue
			}
			channelName := cpString(b, ch.NameCpIndex)
			channel := ChannelDetails{
				Name:                channelName,
				ChannelInSameStrand: ch.IsChannelInSameStrand != 0,
				Send:                ch.IsSend != 0,
			}
			channels = append(channels, channel)
		}
		if len(channels) > 0 {
			fn.SetWorkerChannels(channels)
		}
	}

	return nil
}

// populateBasicBlocks creates BIRBasicBlock instances from the Kaitai model.
func populateBasicBlocks(b *Bir, bbInfo *Bir_BasicBlocksInfo) ([]BIRBasicBlock, error) {
	if bbInfo == nil || bbInfo.BasicBlocksCount == 0 {
		return []BIRBasicBlock{}, nil
	}

	basicBlocks := make([]BIRBasicBlock, 0, bbInfo.BasicBlocksCount)
	bbMap := make(map[string]BIRBasicBlock) // Map BB name to BB for terminator references

	// First pass: create all basic blocks
	for i, kaitaiBB := range bbInfo.BasicBlocks {
		if kaitaiBB == nil {
			continue
		}

		bbName := Name(cpString(b, kaitaiBB.NameCpIndex))
		bb := NewBIRBasicBlock(bbName, i)
		bbMap[bbName.Value()] = bb
		basicBlocks = append(basicBlocks, bb)
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
			bb.SetInstructions(&nonTerminators)
		}
		if terminator != nil {
			bb.SetTerminator(terminator)
		}
	}

	return basicBlocks, nil
}

// createTerminator creates a BIRTerminator instance from Kaitai instruction data.
func createTerminator(b *Bir, kind InstructionKind, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if kaitaiIns == nil || kaitaiIns.InstructionStructure == nil {
		// Fallback to minimal implementation
		return &terminatorImpl{
			BIRTerminatorBase: NewBIRTerminatorBase(pos),
			kind:              kind,
		}
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
	case INSTRUCTION_KIND_ASYNC_CALL:
		return parseAsyncCallTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_FP_CALL:
		return parseFPCallTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_LOCK:
		return parseLockTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_FIELD_LOCK:
		return parseFieldLockTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_UNLOCK:
		return parseUnlockTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_PANIC:
		return parsePanicTerminator(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_WAIT:
		return parseWaitTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_FLUSH:
		return parseFlushTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_WK_RECEIVE:
		return parseWorkerReceiveTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_WK_SEND:
		return parseWorkerSendTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_WK_ALT_RECEIVE:
		return parseWorkerAlternateReceiveTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE:
		return parseWorkerMultipleReceiveTerminator(b, pos, kaitaiIns, bbMap)
	case INSTRUCTION_KIND_WAIT_ALL:
		return parseWaitAllTerminator(b, pos, kaitaiIns, bbMap)
	default:
		// For unknown terminators, return minimal implementation
		return &terminatorImpl{
			BIRTerminatorBase: NewBIRTerminatorBase(pos),
			kind:              kind,
		}
	}
}

// parseGotoTerminator parses a GOTO terminator
func parseGotoTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if gotoIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionGoto); ok && gotoIns != nil {
		targetBBName := cpString(b, gotoIns.TargetBbIdNameCpIndex)
		var targetBB BIRBasicBlock
		if targetBBName != "" {
			targetBB = bbMap[targetBBName]
		}
		return NewBIRTerminatorGOTO(pos, targetBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_GOTO,
	}
}

// parseReturnTerminator parses a Return terminator
func parseReturnTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRTerminator {
	if _, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionReturn); ok {
		return NewBIRTerminatorReturn(pos)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_RETURN,
	}
}

// parseBranchTerminator parses a Branch terminator
func parseBranchTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if branchIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionBranch); ok && branchIns != nil {
		op := parseOperand(b, branchIns.BranchOperand)
		trueBBName := cpString(b, branchIns.TrueBbIdNameCpIndex)
		falseBBName := cpString(b, branchIns.FalseBbIdNameCpIndex)
		var trueBB, falseBB BIRBasicBlock
		if trueBBName != "" {
			trueBB = bbMap[trueBBName]
		}
		if falseBBName != "" {
			falseBB = bbMap[falseBBName]
		}
		return NewBIRTerminatorBranch(pos, op, trueBB, falseBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_BRANCH,
	}
}

// parseCallTerminator parses a Call terminator
func parseCallTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if callIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionCall); ok && callIns != nil {
		callInfo := callIns.CallInstructionInfo
		if callInfo == nil {
			return &terminatorImpl{
				BIRTerminatorBase: NewBIRTerminatorBase(pos),
				kind:              INSTRUCTION_KIND_CALL,
			}
		}

		isVirtual := callInfo.IsVirtual != 0
		// Parse package ID
		var calleePkg model.PackageID
		if callInfo.PackageIndex >= 0 {
			pkgCp, err := cpAsPackage(b, callInfo.PackageIndex)
			if err == nil && pkgCp != nil {
				org := Name(cpString(b, pkgCp.OrgIndex))
				pkgName := Name(cpString(b, pkgCp.PackageNameIndex))
				namePart := Name(cpString(b, pkgCp.NameIndex))
				version := Name(cpString(b, pkgCp.VersionIndex))
				calleePkg = model.NewPackageID(org, []model.Name{pkgName, namePart}, version)
			}
		}

		name := Name(cpString(b, callInfo.CallNameCpIndex))
		args := make([]BIROperand, 0, len(callInfo.Arguments))
		for _, arg := range callInfo.Arguments {
			if arg != nil {
				args = append(args, parseOperand(b, arg))
			}
		}

		var lhsOp BIROperand
		if callInfo.HasLhsOperand != 0 && callInfo.LhsOperand != nil {
			lhsOp = parseOperand(b, callInfo.LhsOperand)
		}

		thenBBName := cpString(b, callIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}

		return NewBIRTerminatorCall(pos, INSTRUCTION_KIND_CALL, isVirtual, calleePkg, name, args, lhsOp, thenBB, []BIRAnnotationAttachment{}, []model.Flag{})
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_CALL,
	}
}

// parseAsyncCallTerminator parses an AsyncCall terminator
func parseAsyncCallTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if asyncCallIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionAsyncCall); ok && asyncCallIns != nil {
		callInfo := asyncCallIns.CallInstructionInfo
		if callInfo == nil {
			return &terminatorImpl{
				BIRTerminatorBase: NewBIRTerminatorBase(pos),
				kind:              INSTRUCTION_KIND_ASYNC_CALL,
			}
		}

		isVirtual := callInfo.IsVirtual != 0
		// Parse package ID
		var calleePkg model.PackageID
		if callInfo.PackageIndex >= 0 {
			pkgCp, err := cpAsPackage(b, callInfo.PackageIndex)
			if err == nil && pkgCp != nil {
				org := Name(cpString(b, pkgCp.OrgIndex))
				pkgName := Name(cpString(b, pkgCp.PackageNameIndex))
				namePart := Name(cpString(b, pkgCp.NameIndex))
				version := Name(cpString(b, pkgCp.VersionIndex))
				calleePkg = model.NewPackageID(org, []model.Name{pkgName, namePart}, version)
			}
		}

		name := Name(cpString(b, callInfo.CallNameCpIndex))
		args := make([]BIROperand, 0, len(callInfo.Arguments))
		for _, arg := range callInfo.Arguments {
			if arg != nil {
				args = append(args, parseOperand(b, arg))
			}
		}

		var lhsOp BIROperand
		if callInfo.HasLhsOperand != 0 && callInfo.LhsOperand != nil {
			lhsOp = parseOperand(b, callInfo.LhsOperand)
		}

		// Parse annotation attachments
		annots := []BIRAnnotationAttachment{}
		if asyncCallIns.AnnotationAttachmentsContent != nil {
			annots = parseAnnotationAttachments(b, asyncCallIns.AnnotationAttachmentsContent)
		}

		thenBBName := cpString(b, asyncCallIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}

		return NewBIRTerminatorAsyncCall(pos, INSTRUCTION_KIND_ASYNC_CALL, isVirtual, calleePkg, name, args, lhsOp, thenBB, annots, []BIRAnnotationAttachment{}, []model.Flag{})
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_ASYNC_CALL,
	}
}

// parseFPCallTerminator parses an FPCall terminator
func parseFPCallTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if fpCallIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionFpCall); ok && fpCallIns != nil {
		fp := parseOperand(b, fpCallIns.FpOperand)
		args := make([]BIROperand, 0, len(fpCallIns.FpArguments))
		for _, arg := range fpCallIns.FpArguments {
			if arg != nil {
				args = append(args, parseOperand(b, arg))
			}
		}

		var lhsOp BIROperand
		if fpCallIns.HasLhsOperand == 1 && fpCallIns.LhsOperand != nil {
			lhsOp = parseOperand(b, fpCallIns.LhsOperand)
		}

		isAsync := fpCallIns.IsAsynch != 0
		annots := []BIRAnnotationAttachment{}
		if fpCallIns.AnnotationAttachmentsContent != nil {
			annots = parseAnnotationAttachments(b, fpCallIns.AnnotationAttachmentsContent)
		}

		thenBBName := cpString(b, fpCallIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}

		return NewBIRTerminatorFPCall(pos, INSTRUCTION_KIND_FP_CALL, fp, args, lhsOp, isAsync, thenBB, annots)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_FP_CALL,
	}
}

// parseLockTerminator parses a Lock terminator
func parseLockTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if lockIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionLock); ok && lockIns != nil {
		lockedBBName := cpString(b, lockIns.LockBbIdNameCpIndex)
		var lockedBB BIRBasicBlock
		if lockedBBName != "" {
			lockedBB = bbMap[lockedBBName]
		}
		lock := NewBIRTerminatorLock(pos, lockedBB)
		// Cast to BIRTerminator since NewBIRTerminatorLock returns BIRTerminatorLock interface
		if term, ok := lock.(BIRTerminator); ok {
			return term
		}
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_LOCK,
	}
}

// parseFieldLockTerminator parses a FieldLock terminator
func parseFieldLockTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if fieldLockIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionFieldLock); ok && fieldLockIns != nil {
		// Parse local var name (for now, create a minimal operand)
		lockVarName := cpString(b, fieldLockIns.LockVarNameCpIndex)
		fieldName := cpString(b, fieldLockIns.FieldNameCpIndex)
		lockedBBName := cpString(b, fieldLockIns.LockBbIdNameCpIndex)
		var lockedBB BIRBasicBlock
		if lockedBBName != "" {
			lockedBB = bbMap[lockedBBName]
		}
		// Create a minimal operand for the local var
		localVar := NewBIROperand(NewBIRVariableDclSimple(nil, Name(lockVarName), VAR_SCOPE_FUNCTION, VAR_KIND_LOCAL))
		return NewBIRTerminatorFieldLock(pos, localVar, fieldName, lockedBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_FIELD_LOCK,
	}
}

// parseUnlockTerminator parses an Unlock terminator
func parseUnlockTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if unlockIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionUnlock); ok && unlockIns != nil {
		unlockBBName := cpString(b, unlockIns.UnlockBbIdNameCpIndex)
		var unlockBB BIRBasicBlock
		if unlockBBName != "" {
			unlockBB = bbMap[unlockBBName]
		}
		return NewBIRTerminatorUnlock(pos, unlockBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_UNLOCK,
	}
}

// parsePanicTerminator parses a Panic terminator
func parsePanicTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRTerminator {
	if panicIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionPanic); ok && panicIns != nil {
		errorOp := parseOperand(b, panicIns.ErrorOperand)
		return NewBIRTerminatorPanic(pos, errorOp)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_PANIC,
	}
}

// parseWaitTerminator parses a Wait terminator
func parseWaitTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if waitIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWait); ok && waitIns != nil {
		exprList := make([]BIROperand, 0, len(waitIns.WaitExpressions))
		for _, expr := range waitIns.WaitExpressions {
			if expr != nil {
				exprList = append(exprList, parseOperand(b, expr))
			}
		}
		lhsOp := parseOperand(b, waitIns.LhsOperand)
		thenBBName := cpString(b, waitIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorWait(pos, exprList, lhsOp, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WAIT,
	}
}

// parseFlushTerminator parses a Flush terminator
func parseFlushTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if flushIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionFlush); ok && flushIns != nil {
		channels := []ChannelDetails{}
		if flushIns.WorkerChannelDetail != nil && len(flushIns.WorkerChannelDetail.WorkerChannelDetails) > 0 {
			for _, wcd := range flushIns.WorkerChannelDetail.WorkerChannelDetails {
				if wcd != nil {
					channelName := cpString(b, wcd.NameCpIndex)
					channel := ChannelDetails{
						Name:                channelName,
						ChannelInSameStrand: wcd.IsChannelInSameStrand != 0,
						Send:                wcd.IsSend != 0,
					}
					channels = append(channels, channel)
				}
			}
		}
		lhsOp := parseOperand(b, flushIns.LhsOperand)
		thenBBName := cpString(b, flushIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorFlush(pos, channels, lhsOp, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_FLUSH,
	}
}

// parseWorkerReceiveTerminator parses a WorkerReceive terminator
func parseWorkerReceiveTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if wrIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWkReceive); ok && wrIns != nil {
		workerName := Name(cpString(b, wrIns.WorkerNameCpIndex))
		lhsOp := parseOperand(b, wrIns.LhsOperand)
		isSameStrand := wrIns.IsSameStrand != 0
		thenBBName := cpString(b, wrIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorWorkerReceive(pos, workerName, lhsOp, isSameStrand, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WK_RECEIVE,
	}
}

// parseWorkerSendTerminator parses a WorkerSend terminator
func parseWorkerSendTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if wsIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWkSend); ok && wsIns != nil {
		channel := Name(cpString(b, wsIns.ChannelNameCpIndex))
		data := parseOperand(b, wsIns.WorkerDataOperand)
		isSameStrand := wsIns.IsSameStrand != 0
		isSync := wsIns.IsSynch != 0
		var lhsOp BIROperand
		if isSync && wsIns.LhsOperand != nil {
			lhsOp = parseOperand(b, wsIns.LhsOperand)
		}
		thenBBName := cpString(b, wsIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorWorkerSend(pos, channel, data, isSameStrand, isSync, lhsOp, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WK_SEND,
	}
}

// parseWorkerAlternateReceiveTerminator parses a WorkerAlternateReceive terminator
func parseWorkerAlternateReceiveTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if arIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWkAltReceive); ok && arIns != nil {
		channels := make([]string, 0, len(arIns.ChannelNameCpIndex))
		for _, chIdx := range arIns.ChannelNameCpIndex {
			channels = append(channels, cpString(b, chIdx))
		}
		lhsOp := parseOperand(b, arIns.LhsOperand)
		isSameStrand := arIns.IsSameStrand != 0
		thenBBName := cpString(b, arIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorWorkerAlternateReceive(pos, channels, lhsOp, isSameStrand, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WK_ALT_RECEIVE,
	}
}

// parseWorkerMultipleReceiveTerminator parses a WorkerMultipleReceive terminator
func parseWorkerMultipleReceiveTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if mrIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWkMulReceive); ok && mrIns != nil {
		receiveFields := make([]ReceiveField, 0, len(mrIns.ChannelFieldCpIndex))
		for _, rf := range mrIns.ChannelFieldCpIndex {
			if rf != nil {
				key := cpString(b, rf.FieldName)
				workerReceive := cpString(b, rf.ChannelName)
				receiveFields = append(receiveFields, ReceiveField{
					Key:           key,
					WorkerReceive: workerReceive,
				})
			}
		}
		var targetType BType
		if mrIns.TypeCpIndex >= 0 {
			targetType = parseTypeFromCP(b, mrIns.TypeCpIndex)
		}
		lhsOp := parseOperand(b, mrIns.LhsOperand)
		isSameStrand := mrIns.IsSameStrand != 0
		thenBBName := cpString(b, mrIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		mr := NewBIRTerminatorWorkerMultipleReceive(pos, receiveFields, lhsOp, isSameStrand, thenBB)
		mr.TargetType = targetType
		return mr
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE,
	}
}

// parseWaitAllTerminator parses a WaitAll terminator
func parseWaitAllTerminator(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction, bbMap map[string]BIRBasicBlock) BIRTerminator {
	if waIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionWaitAll); ok && waIns != nil {
		keys := make([]string, 0, len(waIns.KeyNameCpIndex))
		for _, keyIdx := range waIns.KeyNameCpIndex {
			keys = append(keys, cpString(b, keyIdx))
		}
		valueExprs := make([]BIROperand, 0, len(waIns.ValueExpression))
		for _, val := range waIns.ValueExpression {
			if val != nil {
				valueExprs = append(valueExprs, parseOperand(b, val))
			}
		}
		lhsOp := parseOperand(b, waIns.LhsOperand)
		thenBBName := cpString(b, waIns.ThenBbIdNameCpIndex)
		var thenBB BIRBasicBlock
		if thenBBName != "" {
			thenBB = bbMap[thenBBName]
		}
		return NewBIRTerminatorWaitAll(pos, lhsOp, keys, valueExprs, thenBB)
	}
	return &terminatorImpl{
		BIRTerminatorBase: NewBIRTerminatorBase(pos),
		kind:              INSTRUCTION_KIND_WAIT_ALL,
	}
}

// createNonTerminator creates a BIRNonTerminator instance from Kaitai instruction data.
func createNonTerminator(b *Bir, kind InstructionKind, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if kaitaiIns == nil || kaitaiIns.InstructionStructure == nil {
		// Fallback to minimal implementation
		return &nonTerminatorImpl{
			BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
			kind:                 kind,
		}
	}

	// Parse based on instruction kind
	switch kind {
	case INSTRUCTION_KIND_MOVE:
		return parseMoveInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_CONST_LOAD:
		return parseConstantLoadInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_NEW_STRUCTURE:
		return parseNewStructureInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_NEW_ARRAY:
		return parseNewArrayInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_MAP_LOAD, INSTRUCTION_KIND_ARRAY_LOAD, INSTRUCTION_KIND_MAP_STORE, INSTRUCTION_KIND_ARRAY_STORE,
		INSTRUCTION_KIND_OBJECT_LOAD, INSTRUCTION_KIND_OBJECT_STORE, INSTRUCTION_KIND_STRING_LOAD,
		INSTRUCTION_KIND_XML_LOAD, INSTRUCTION_KIND_XML_SEQ_LOAD, INSTRUCTION_KIND_XML_ATTRIBUTE_LOAD, INSTRUCTION_KIND_XML_ATTRIBUTE_STORE:
		return parseFieldAccessInstruction(b, pos, kind, kaitaiIns)
	case INSTRUCTION_KIND_NEW_ERROR:
		return parseNewErrorInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_TYPE_CAST:
		return parseTypeCastInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_IS_LIKE:
		return parseIsLikeInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_TYPE_TEST:
		return parseTypeTestInstruction(b, pos, kaitaiIns)
	case INSTRUCTION_KIND_ADD, INSTRUCTION_KIND_SUB, INSTRUCTION_KIND_MUL, INSTRUCTION_KIND_DIV, INSTRUCTION_KIND_MOD,
		INSTRUCTION_KIND_EQUAL, INSTRUCTION_KIND_NOT_EQUAL, INSTRUCTION_KIND_GREATER_THAN, INSTRUCTION_KIND_GREATER_EQUAL,
		INSTRUCTION_KIND_LESS_THAN, INSTRUCTION_KIND_LESS_EQUAL, INSTRUCTION_KIND_AND, INSTRUCTION_KIND_OR,
		INSTRUCTION_KIND_REF_EQUAL, INSTRUCTION_KIND_REF_NOT_EQUAL, INSTRUCTION_KIND_CLOSED_RANGE, INSTRUCTION_KIND_HALF_OPEN_RANGE,
		INSTRUCTION_KIND_ANNOT_ACCESS:
		return parseBinaryOpInstruction(b, pos, kind, kaitaiIns)
	case INSTRUCTION_KIND_TYPEOF, INSTRUCTION_KIND_NOT, INSTRUCTION_KIND_NEGATE:
		return parseUnaryOpInstruction(b, pos, kind, kaitaiIns)
	default:
		// For other instructions, return minimal implementation
		return &nonTerminatorImpl{
			BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
			kind:                 kind,
		}
	}
}

// parseMoveInstruction parses a Move instruction
func parseMoveInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if moveIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionMove); ok && moveIns != nil {
		rhsOp := parseOperand(b, moveIns.RhsOperand)
		lhsOp := parseOperand(b, moveIns.LhsOperand)
		return NewBIRNonTerminatorMove(pos, rhsOp, lhsOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_MOVE,
	}
}

// parseConstantLoadInstruction parses a ConstantLoad instruction
func parseConstantLoadInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if constIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionConstLoad); ok && constIns != nil {
		lhsOp := parseOperand(b, constIns.LhsOperand)
		var constType BType
		if constIns.TypeCpIndex >= 0 {
			constType = parseTypeFromCP(b, constIns.TypeCpIndex)
		}
		// Parse constant value
		var value interface{}
		if constIns.ConstantValueInfo != nil {
			value = parseConstantValueInfo(b, constIns.ConstantValueInfo, constType)
		}
		return NewBIRNonTerminatorConstantLoad(pos, value, constType, lhsOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_CONST_LOAD,
	}
}

// parseBinaryOpInstruction parses a BinaryOp instruction
func parseBinaryOpInstruction(b *Bir, pos diagnostics.Location, kind InstructionKind, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if binOpIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionBinaryOperation); ok && binOpIns != nil {
		rhsOp1 := parseOperand(b, binOpIns.RhsOperandOne)
		rhsOp2 := parseOperand(b, binOpIns.RhsOperandTwo)
		lhsOp := parseOperand(b, binOpIns.LhsOperand)
		return NewBIRNonTerminatorBinaryOp(pos, kind, lhsOp, rhsOp1, rhsOp2)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 kind,
	}
}

// parseUnaryOpInstruction parses a UnaryOp instruction
func parseUnaryOpInstruction(b *Bir, pos diagnostics.Location, kind InstructionKind, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if unaryOpIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionUnaryOperation); ok && unaryOpIns != nil {
		rhsOp := parseOperand(b, unaryOpIns.RhsOperand)
		lhsOp := parseOperand(b, unaryOpIns.LhsOperand)
		return NewBIRNonTerminatorUnaryOP(pos, kind, lhsOp, rhsOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 kind,
	}
}

// parseNewStructureInstruction parses a NewStructure instruction
func parseNewStructureInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if structIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionNewStructure); ok && structIns != nil {
		rhsOp := parseOperand(b, structIns.RhsOperand)
		lhsOp := parseOperand(b, structIns.LhsOperand)
		initialValues := make([]BIRMappingConstructorEntry, 0, len(structIns.InitValues))
		for _, mc := range structIns.InitValues {
			if mc == nil {
				continue
			}
			if mc.MappingConstructorKind == 1 { // key-value
				if kvBody, ok := mc.MappingConstructorBody.(*Bir_MappingConstructorKeyValueBody); ok && kvBody != nil {
					keyOp := parseOperand(b, kvBody.KeyOperand)
					valueOp := parseOperand(b, kvBody.ValueOperand)
					entry := NewBIRMappingConstructorKeyValueEntry(keyOp, valueOp)
					initialValues = append(initialValues, entry)
				}
			} else { // spread field
				if spreadBody, ok := mc.MappingConstructorBody.(*Bir_MappingConstructorSpreadFieldBody); ok && spreadBody != nil {
					exprOp := parseOperand(b, spreadBody.ExprOperand)
					entry := NewBIRMappingConstructorSpreadFieldEntry(exprOp)
					initialValues = append(initialValues, entry)
				}
			}
		}
		return NewBIRNonTerminatorNewStructure(pos, lhsOp, rhsOp, initialValues)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_NEW_STRUCTURE,
	}
}

// parseNewArrayInstruction parses a NewArray instruction
func parseNewArrayInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if arrayIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionNewArray); ok && arrayIns != nil {
		var arrayType BType
		if arrayIns.TypeCpIndex >= 0 {
			arrayType = parseTypeFromCP(b, arrayIns.TypeCpIndex)
		}
		lhsOp := parseOperand(b, arrayIns.LhsOperand)
		sizeOp := parseOperand(b, arrayIns.SizeOperand)
		values := make([]BIRListConstructorEntry, 0, len(arrayIns.InitValues))
		for _, val := range arrayIns.InitValues {
			if val == nil {
				continue
			}
			exprOp := parseOperand(b, val)
			entry := NewBIRListConstructorExprEntry(exprOp)
			values = append(values, entry)
		}
		na := NewBIRNonTerminatorNewArray(pos, arrayType, lhsOp, sizeOp, values)
		if arrayIns.HasTypedescOperand == 1 && arrayIns.TypedescOperand != nil {
			na.TypedescOp = parseOperand(b, arrayIns.TypedescOperand)
		}
		if arrayIns.HasElementTypedescOperand == 1 && arrayIns.ElementTypedescOperand != nil {
			na.ElementTypedescOp = parseOperand(b, arrayIns.ElementTypedescOperand)
		}
		return na
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_NEW_ARRAY,
	}
}

// parseFieldAccessInstruction parses a FieldAccess instruction (MAP_LOAD, ARRAY_LOAD, etc.)
func parseFieldAccessInstruction(b *Bir, pos diagnostics.Location, kind InstructionKind, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	// FieldAccess is used for MAP_LOAD, ARRAY_LOAD, MAP_STORE, ARRAY_STORE, etc.
	// They all use index_access structure
	var keyOp, rhsOp, lhsOp BIROperand
	var optionalFieldAccess, fillingRead bool

	switch ins := kaitaiIns.InstructionStructure.(type) {
	case *Bir_InstructionMapLoad:
		if ins != nil {
			optionalFieldAccess = ins.IsOptionalFieldAccess != 0
			fillingRead = ins.IsFillingRead != 0
			if ins.MapLoad != nil {
				lhsOp = parseOperand(b, ins.MapLoad.LhsOperand)
				keyOp = parseOperand(b, ins.MapLoad.KeyOperand)
				rhsOp = parseOperand(b, ins.MapLoad.RhsOperand)
			}
		}
	case *Bir_InstructionArrayLoad:
		if ins != nil {
			optionalFieldAccess = ins.IsOptionalFieldAccess != 0
			fillingRead = ins.IsFillingRead != 0
			if ins.ArrayLoad != nil {
				lhsOp = parseOperand(b, ins.ArrayLoad.LhsOperand)
				keyOp = parseOperand(b, ins.ArrayLoad.KeyOperand)
				rhsOp = parseOperand(b, ins.ArrayLoad.RhsOperand)
			}
		}
	case *Bir_InstructionMapStore:
		if ins != nil && ins.MapStore != nil {
			lhsOp = parseOperand(b, ins.MapStore.LhsOperand)
			keyOp = parseOperand(b, ins.MapStore.KeyOperand)
			rhsOp = parseOperand(b, ins.MapStore.RhsOperand)
		}
	case *Bir_InstructionArrayStore:
		if ins != nil && ins.ArrayStore != nil {
			lhsOp = parseOperand(b, ins.ArrayStore.LhsOperand)
			keyOp = parseOperand(b, ins.ArrayStore.KeyOperand)
			rhsOp = parseOperand(b, ins.ArrayStore.RhsOperand)
		}
	}

	if lhsOp != nil {
		fa := NewBIRNonTerminatorFieldAccess(pos, kind, lhsOp, keyOp, rhsOp)
		fa.OptionalFieldAccess = optionalFieldAccess
		fa.FillingRead = fillingRead
		return fa
	}

	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 kind,
	}
}

// parseNewErrorInstruction parses a NewError instruction
func parseNewErrorInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if errorIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionNewError); ok && errorIns != nil {
		var errorType BType
		if errorIns.ErrorTypeCpIndex >= 0 {
			errorType = parseTypeFromCP(b, errorIns.ErrorTypeCpIndex)
		}
		lhsOp := parseOperand(b, errorIns.LhsOperand)
		messageOp := parseOperand(b, errorIns.MessageOperand)
		causeOp := parseOperand(b, errorIns.CauseOperand)
		detailOp := parseOperand(b, errorIns.DetailOperand)
		return NewBIRNonTerminatorNewError(pos, errorType, lhsOp, messageOp, causeOp, detailOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_NEW_ERROR,
	}
}

// parseTypeCastInstruction parses a TypeCast instruction
func parseTypeCastInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if castIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionTypeCast); ok && castIns != nil {
		lhsOp := parseOperand(b, castIns.LhsOperand)
		rhsOp := parseOperand(b, castIns.RhsOperand)
		var castType BType
		if castIns.TypeCpIndex >= 0 {
			castType = parseTypeFromCP(b, castIns.TypeCpIndex)
		}
		checkTypes := castIns.IsCheckTypes != 0
		return NewBIRNonTerminatorTypeCast(pos, lhsOp, rhsOp, castType, checkTypes)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_TYPE_CAST,
	}
}

// parseIsLikeInstruction parses an IsLike instruction
func parseIsLikeInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if isLikeIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionIsLike); ok && isLikeIns != nil {
		var likeType BType
		if isLikeIns.TypeCpIndex >= 0 {
			likeType = parseTypeFromCP(b, isLikeIns.TypeCpIndex)
		}
		lhsOp := parseOperand(b, isLikeIns.LhsOperand)
		rhsOp := parseOperand(b, isLikeIns.RhsOperand)
		return NewBIRNonTerminatorIsLike(pos, likeType, lhsOp, rhsOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_IS_LIKE,
	}
}

// parseTypeTestInstruction parses a TypeTest instruction
func parseTypeTestInstruction(b *Bir, pos diagnostics.Location, kaitaiIns *Bir_Instruction) BIRNonTerminator {
	if testIns, ok := kaitaiIns.InstructionStructure.(*Bir_InstructionTypeTest); ok && testIns != nil {
		var testType BType
		if testIns.TypeCpIndex >= 0 {
			testType = parseTypeFromCP(b, testIns.TypeCpIndex)
		}
		lhsOp := parseOperand(b, testIns.LhsOperand)
		rhsOp := parseOperand(b, testIns.RhsOperand)
		return NewBIRNonTerminatorTypeTest(pos, testType, lhsOp, rhsOp)
	}
	return &nonTerminatorImpl{
		BIRNonTerminatorBase: NewBIRNonTerminatorBase(pos),
		kind:                 INSTRUCTION_KIND_TYPE_TEST,
	}
}

// positionToLocation converts a Bir_Position to diagnostics.Location.
// For now, returns nil as we don't have a full Location implementation.
func positionToLocation(b *Bir, pos *Bir_Position) diagnostics.Location {
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
	if t.GetThenBB() != nil {
		return []BIRBasicBlock{t.GetThenBB()}
	}
	return []BIRBasicBlock{}
}

func (t *terminatorImpl) Accept(visitor BIRVisitor) {
	// Dispatch to appropriate visitor method based on kind
	switch t.kind {
	case INSTRUCTION_KIND_GOTO:
		visitor.VisitBIRTerminatorGOTO(t)
	case INSTRUCTION_KIND_RETURN:
		visitor.VisitBIRTerminatorReturn(t)
	case INSTRUCTION_KIND_BRANCH:
		visitor.VisitBIRTerminatorBranch(t)
	default:
		// Default to generic terminator visit
		visitor.VisitBIRTerminatorGOTO(t)
	}
}

type nonTerminatorImpl struct {
	BIRNonTerminatorBase
	kind InstructionKind
}

func (n *nonTerminatorImpl) Accept(visitor BIRVisitor) {
	// Dispatch to appropriate visitor method based on kind
	// For now, use a generic non-terminator visit
	visitor.VisitBIRNonTerminatorMove(n)
}

// Helper functions for parsing various structures

// populatePathParameters populates path parameters for resource functions.
func populatePathParameters(b *Bir, fn BIRFunction, rfc *Bir_ResourceFunctionContent) error {
	if rfc == nil {
		return nil
	}

	// Path parameters
	if rfc.PathParamsCount > 0 && len(rfc.PathParams) > 0 {
		pathParams := make([]BIRVariableDcl, 0, len(rfc.PathParams))
		for _, pp := range rfc.PathParams {
			if pp == nil {
				continue
			}
			metaVarName := cpString(b, pp.PathParamNameCpIndex)
			var t BType
			if pp.PathParamTypeCpIndex >= 0 {
				t = parseTypeFromCP(b, pp.PathParamTypeCpIndex)
			}
			name := Name(metaVarName)
			pathVar := NewBIRVariableDclSimple(
				t,
				name,
				VAR_SCOPE_FUNCTION,
				VAR_KIND_LOCAL,
			)
			pathParams = append(pathParams, pathVar)
		}
		if len(pathParams) > 0 {
			fn.SetPathParams(&pathParams)
		}
	}

	// Rest path parameter
	if rfc.HasRestPathParam != 0 && rfc.RestPathParam != nil {
		metaVarName := cpString(b, rfc.RestPathParam.PathParamNameCpIndex)
		name := Name(metaVarName)
		var t BType
		if rfc.RestPathParam.PathParamTypeCpIndex >= 0 {
			t = parseTypeFromCP(b, rfc.RestPathParam.PathParamTypeCpIndex)
		}
		restPathVar := NewBIRVariableDclSimple(
			t,
			name,
			VAR_SCOPE_FUNCTION,
			VAR_KIND_LOCAL,
		)
		fn.SetRestPathParam(restPathVar)
	}

	// Resource path segments
	if rfc.ResourcePathSegmentCount > 0 && len(rfc.ResourcePathSegments) > 0 {
		resourcePath := make([]Name, 0, len(rfc.ResourcePathSegments))
		pathSegmentPosList := make([]diagnostics.Location, 0, len(rfc.ResourcePathSegments))
		pathSegmentTypeList := make([]BType, 0, len(rfc.ResourcePathSegments))

		for _, seg := range rfc.ResourcePathSegments {
			if seg == nil {
				continue
			}
			pathName := Name(cpString(b, seg.ResourcePathSegmentCpIndex))
			resourcePath = append(resourcePath, pathName)
			pathSegmentPosList = append(pathSegmentPosList, positionToLocation(b, seg.ResourcePathSegmentPos))
			var t BType
			if seg.ResourcePathSegmentType >= 0 {
				t = parseTypeFromCP(b, seg.ResourcePathSegmentType)
			}
			pathSegmentTypeList = append(pathSegmentTypeList, t)
		}

		if len(resourcePath) > 0 {
			fn.SetResourcePath(&resourcePath)
			fn.SetResourcePathSegmentPosList(&pathSegmentPosList)
			fn.SetPathSegmentTypeList(&pathSegmentTypeList)
		}
	}

	// Accessor
	if rfc.ResourceAccessor >= 0 {
		// TODO: Map accessor from CP index or value
		accessorName := Name(cpString(b, rfc.ResourceAccessor))
		fn.SetAccessor(accessorName)
	}

	return nil
}

// parseAnnotationAttachments parses annotation attachments.
func parseAnnotationAttachments(b *Bir, aac *Bir_AnnotationAttachmentsContent) []BIRAnnotationAttachment {
	if aac == nil || aac.AttachmentsCount == 0 {
		return []BIRAnnotationAttachment{}
	}

	attachments := make([]BIRAnnotationAttachment, 0, aac.AttachmentsCount)
	for _, aa := range aac.AnnotationAttachments {
		if aa == nil {
			continue
		}

		// Parse PackageID
		var annotPkgID model.PackageID
		if aa.PackageIdCpIndex >= 0 {
			pkgCp, err := cpAsPackage(b, aa.PackageIdCpIndex)
			if err == nil && pkgCp != nil {
				org := Name(cpString(b, pkgCp.OrgIndex))
				pkgName := Name(cpString(b, pkgCp.PackageNameIndex))
				namePart := Name(cpString(b, pkgCp.NameIndex))
				version := Name(cpString(b, pkgCp.VersionIndex))
				annotPkgID = model.NewPackageID(org, []model.Name{pkgName, namePart}, version)
			}
		}

		// Parse tag reference
		tagRef := Name(cpString(b, aa.TagReferenceCpIndex))
		pos := positionToLocation(b, aa.Position)

		// Check if this is a const annotation
		if aa.IsConstAnnot != 0 && aa.ConstantValue != nil {
			// Create const annotation attachment
			constValue := parseConstantValue(b, aa.ConstantValue)
			att := NewBIRConstAnnotationAttachment(pos, annotPkgID, tagRef, constValue)
			attachments = append(attachments, att)
		} else {
			// Regular annotation attachment
			att := NewBIRAnnotationAttachment(pos, annotPkgID, tagRef)
			attachments = append(attachments, att)
		}
	}

	return attachments
}

// parseReturnVar parses a return variable.
func parseReturnVar(b *Bir, rv *Bir_ReturnVar) BIRVariableDcl {
	if rv == nil {
		return nil
	}

	name := Name(cpString(b, rv.NameCpIndex))
	kind := VarKind(rv.Kind)
	var t BType
	if rv.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, rv.TypeCpIndex)
	}

	return NewBIRVariableDclSimple(
		t,
		name,
		VAR_SCOPE_FUNCTION,
		kind,
	)
}

// parseFunctionParameter parses a function parameter (default parameter).
func parseFunctionParameter(b *Bir, dp *Bir_DefaultParameter) BIRFunctionParameter {
	if dp == nil {
		return nil
	}

	name := Name(cpString(b, dp.NameCpIndex))
	kind := VarKind(dp.Kind)
	var t BType
	if dp.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, dp.TypeCpIndex)
	}
	metaVarName := cpString(b, dp.MetaVarNameCpIndex)
	hasDefaultExpr := dp.HasDefaultExpr != 0

	return NewBIRFunctionParameter(
		nil, // pos
		t,
		name,
		VAR_SCOPE_FUNCTION,
		kind,
		metaVarName,
		hasDefaultExpr,
	)
}

// parseLocalVariable parses a local variable.
func parseLocalVariable(b *Bir, lv *Bir_LocalVariable) BIRVariableDcl {
	if lv == nil {
		return nil
	}

	name := Name(cpString(b, lv.NameCpIndex))
	kind := VarKind(lv.Kind)
	var t BType
	if lv.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, lv.TypeCpIndex)
	}

	localVar := NewBIRVariableDclSimple(
		t,
		name,
		VAR_SCOPE_FUNCTION,
		kind,
	)

	// Parse enclosing basic block info if this is a LOCAL variable
	if kind == VAR_KIND_LOCAL && lv.EnclosingBasicBlockId != nil {
		// TODO: Set startBB, endBB, and insOffset if needed
		// This requires access to the basic blocks map
	}

	return localVar
}

// parseErrorEntry parses an error table entry.
func parseErrorEntry(b *Bir, ee *Bir_ErrorEntry) BIRErrorEntry {
	if ee == nil {
		return nil
	}

	// Get basic block names from CP indices
	trapBBName := cpString(b, ee.TrapBbIdCpIndex)
	endBBName := cpString(b, ee.EndBbIdCpIndex)
	targetBBName := cpString(b, ee.TargetBbIdCpIndex)

	// Create basic block references by name (they will be resolved later if needed)
	// For now, we create minimal basic blocks with just the name
	trapBB := NewBIRBasicBlock(Name(trapBBName), 0)
	endBB := NewBIRBasicBlock(Name(endBBName), 0)
	targetBB := NewBIRBasicBlock(Name(targetBBName), 0)

	// Parse error operand
	var errorOp BIROperand
	if ee.ErrorOperand != nil {
		errorOp = parseOperand(b, ee.ErrorOperand)
	}

	return NewBIRErrorEntry(trapBB, endBB, errorOp, targetBB)
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
		name:  Name(cpString(b, ti.NameIndex)),
		flags: ti.TypeFlag,
	}
}

// minimalBType is a minimal implementation of BType for parsing purposes.
type minimalBType struct {
	tag     int
	name    Name
	flags   int64
	semType interface{}
	tsymbol interface{}
}

func (t *minimalBType) GetTag() int                    { return t.tag }
func (t *minimalBType) SetTag(tag int)                 { t.tag = tag }
func (t *minimalBType) GetName() Name                  { return t.name }
func (t *minimalBType) SetName(name Name)              { t.name = name }
func (t *minimalBType) GetFlags() int64                { return t.flags }
func (t *minimalBType) SetFlags(flags int64)           { t.flags = flags }
func (t *minimalBType) GetSemType() interface{}        { return t.semType }
func (t *minimalBType) SetSemType(semType interface{}) { t.semType = semType }
func (t *minimalBType) GetTsymbol() interface{}        { return t.tsymbol }
func (t *minimalBType) SetTsymbol(tsymbol interface{}) { t.tsymbol = tsymbol }
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
func parseMarkdown(b *Bir, md *Bir_Markdown) MarkdownDocAttachment {
	if md == nil || md.HasDoc == 0 || md.MarkdownContent == nil {
		return MarkdownDocAttachment{}
	}

	mc := md.MarkdownContent
	mdDoc := MarkdownDocAttachment{}

	// Parse description
	if mc.DescriptionCpIndex >= 0 {
		desc := cpString(b, mc.DescriptionCpIndex)
		mdDoc.Description = &desc
	}

	// Parse return value description
	if mc.ReturnValueDescriptionCpIndex >= 0 {
		retDesc := cpString(b, mc.ReturnValueDescriptionCpIndex)
		mdDoc.ReturnValueDescription = &retDesc
	}

	// Parse parameters
	if mc.ParametersCount > 0 && len(mc.Parameters) > 0 {
		params := make([]model.Parameters, 0, len(mc.Parameters))
		for _, mp := range mc.Parameters {
			if mp == nil {
				continue
			}
			paramName := cpString(b, mp.NameCpIndex)
			paramDesc := cpString(b, mp.DescriptionCpIndex)
			params = append(params, model.Parameters{
				Name:        &paramName,
				Description: &paramDesc,
			})
		}
		mdDoc.Parameters = params
	}

	// Parse deprecated docs
	if mc.DeprecatedDocsCpIndex >= 0 {
		deprecatedDocs := cpString(b, mc.DeprecatedDocsCpIndex)
		mdDoc.DeprecatedDocumentation = &deprecatedDocs
	}

	// Parse deprecated params
	if mc.DeprecatedParamsCount > 0 && len(mc.DeprecatedParams) > 0 {
		deprecatedParams := make([]model.Parameters, 0, len(mc.DeprecatedParams))
		for _, mp := range mc.DeprecatedParams {
			if mp == nil {
				continue
			}
			paramName := cpString(b, mp.NameCpIndex)
			paramDesc := cpString(b, mp.DescriptionCpIndex)
			deprecatedParams = append(deprecatedParams, model.Parameters{
				Name:        &paramName,
				Description: &paramDesc,
			})
		}
		mdDoc.DeprecatedParameters = deprecatedParams
	}

	return mdDoc
}

// parseConstantValue parses a constant value from Bir_ConstantValue.
func parseConstantValue(b *Bir, cv *Bir_ConstantValue) ConstValue {
	if cv == nil {
		return ConstValue{}
	}

	var constType BType
	if cv.ConstantValueTypeCpIndex >= 0 {
		constType = parseTypeFromCP(b, cv.ConstantValueTypeCpIndex)
	}

	// Parse constant value info based on type
	var value interface{}
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
func parseConstantValueInfo(b *Bir, cvi kaitai.Struct, constType BType) interface{} {
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
			result := make(map[string]interface{})
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
			result := make([]interface{}, 0, len(v.ListMemberValueInfo))
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
func parseReceiver(b *Bir, rec *Bir_Reciever) BIRVariableDcl {
	if rec == nil {
		return nil
	}

	name := Name(cpString(b, rec.NameCpIndex))
	kind := VarKind(rec.Kind)
	var t BType
	if rec.TypeCpIndex >= 0 {
		t = parseTypeFromCP(b, rec.TypeCpIndex)
	}

	return NewBIRVariableDclSimple(
		t,
		name,
		VAR_SCOPE_FUNCTION,
		kind,
	)
}

// parseOperand parses a BIROperand from Bir_Operand.
func parseOperand(b *Bir, op *Bir_Operand) BIROperand {
	if op == nil {
		return nil
	}

	// Check if variable is ignored
	if op.IgnoredVariable != 0 {
		// Ignored variable - create a minimal variable declaration
		var ignoredType BType
		if op.IgnoredTypeCpIndex >= 0 {
			ignoredType = parseTypeFromCP(b, op.IgnoredTypeCpIndex)
		}
		// Create a minimal variable for ignored operands
		ignoredVar := NewBIRVariableDclSimple(
			ignoredType,
			Name("_"),
			VAR_SCOPE_FUNCTION,
			VAR_KIND_LOCAL,
		)
		return NewBIROperand(ignoredVar)
	}

	// Parse variable
	if op.Variable == nil {
		return nil
	}

	var varType BType
	varName := Name(cpString(b, op.Variable.VariableDclNameCpIndex))
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
	varDcl := NewBIRVariableDclSimple(
		varType,
		varName,
		scope,
		kind,
	)

	return NewBIROperand(varDcl)
}
