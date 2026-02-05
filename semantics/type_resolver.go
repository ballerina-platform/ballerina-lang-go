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

package semantics

import (
	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
	"math/big"
	"strconv"
)

type (
	TypeResolver struct {
		ctx       *context.CompilerContext
		typeEnv   semtypes.Env
		typeDefns map[model.SymbolRef]*ast.BLangTypeDefinition
	}
)

// symbolTypeSetter is redefined here to allow setting types on symbols.
// This must match the private interface in model/symbol.go.
type symbolTypeSetter interface {
	SetType(semtypes.SemType)
}

var _ ast.Visitor = &TypeResolver{}

func NewTypeResolver(ctx *context.CompilerContext) *TypeResolver {
	return &TypeResolver{
		typeEnv:   semtypes.GetTypeEnv(),
		ctx:       ctx,
		typeDefns: make(map[model.SymbolRef]*ast.BLangTypeDefinition),
	}
}

// NewIsolatedTypeResolver is meant for testing so that we can run each test in parallel
func NewIsolatedTypeResolver(ctx *context.CompilerContext) *TypeResolver {
	return &TypeResolver{
		typeEnv:   semtypes.GetIsolatedTypeEnv(),
		ctx:       ctx,
		typeDefns: make(map[model.SymbolRef]*ast.BLangTypeDefinition),
	}
}

// ResolveTypes resolves all the type definitions and return a map of all the types of symbols exported by the package.
// After this (for the given package) all the semtypes are known. Semantic analysis will validate and propagate these
// types to the rest of nodes based on semantic information. This means after Resolving types of all the packages
// it is safe use the closed world assumption to optimize type checks.
func (t *TypeResolver) ResolveTypes(ctx *context.CompilerContext, pkg *ast.BLangPackage) {
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		symbol := defn.Symbol().(*model.SymbolRef)
		t.typeDefns[*symbol] = defn
	}
	for i := range pkg.TypeDefinitions {
		defn := &pkg.TypeDefinitions[i]
		t.resolveTypeDefinition(defn, 0)
	}
	ast.Walk(t, pkg)
	tctx := semtypes.ContextFrom(t.typeEnv)
	for _, defn := range pkg.TypeDefinitions {
		if semtypes.IsEmpty(tctx, defn.DeterminedType) {
			t.ctx.SemanticError(fmt.Sprintf("type definition %s is empty", defn.Name.GetValue()), defn.GetPosition())
		}
	}
	for _, fn := range pkg.Functions {
		t.resolveFunction(ctx, &fn)
	}
}

func (t *TypeResolver) resolveFunction(ctx *context.CompilerContext, fn *ast.BLangFunction) semtypes.SemType {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i, param := range fn.RequiredParams {
		typeData := param.GetTypeData()
		if typeData.Type != nil {
			// Already resolved
			paramTypes[i] = typeData.Type
		} else {
			// This should be resolved already
			paramTypes[i] = typeData.Type
		}
	}
	var restTy semtypes.SemType
	if fn.RestParam != nil {
		t.ctx.Unimplemented("var args not supported", fn.RestParam.GetPosition())
	} else {
		restTy = &semtypes.NEVER
	}
	paramListDefn := semtypes.NewListDefinition()
	paramListTy := paramListDefn.DefineListTypeWrapped(t.typeEnv, paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	returnTypeData := fn.GetReturnTypeData()
	if returnTypeData.TypeDescriptor != nil {
		// Already resolved
		returnTy = returnTypeData.Type
	} else {
		returnTy = &semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	fnType := functionDefn.Define(t.typeEnv, paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.typeEnv, false, false))

	// Update symbol type for the function
	updateSymbolType(t.ctx, fn, fnType)
	fnSymbol := ctx.GetSymbol(fn.Symbol()).(*model.FunctionSymbol)
	fnSymbol.Signature.ParamTypes = paramTypes
	fnSymbol.Signature.ReturnType = returnTy
	fnSymbol.Signature.RestParamType = restTy

	return fnType
}

func (t *TypeResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return t
	}
	ty := t.resolveBType(typeData.TypeDescriptor.(ast.BType), 0)
	typeData.Type = ty

	// Update symbol type if the type descriptor has a symbol
	if tdNode, ok := typeData.TypeDescriptor.(ast.BLangNode); ok {
		updateSymbolType(t.ctx, tdNode, ty)
	}

	return t
}

func (t *TypeResolver) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangFunction, *ast.BLangConstant:
		return t
	case *ast.BLangSimpleVariable:
		t.resolveSimpleVariable(node.(*ast.BLangSimpleVariable))
		return t
	case *ast.BLangArrayType, *ast.BLangBuiltInRefTypeNode, *ast.BLangValueType, *ast.BLangUserDefinedType, *ast.BLangFiniteTypeNode, *ast.BLangUnionTypeNode, *ast.BLangErrorTypeNode:
		t.resolveBType(node.(ast.BType), 0)
		return t
	case *ast.BLangLiteral:
		t.resolveLiteral(n)
		return nil
	case *ast.BLangNumericLiteral:
		t.resolveNumericLiteral(n)
		return nil
	case *ast.BLangTypeDefinition:
		t.resolveTypeDefinition(n, 0)
		return nil
	default:
		return t
	}
}

func (t *TypeResolver) resolveTypeDefinition(defn *ast.BLangTypeDefinition, depth int) semtypes.SemType {
	if defn.DeterminedType != nil {
		return defn.DeterminedType
	}
	if depth == defn.CycleDepth {
		t.ctx.SemanticError(fmt.Sprintf("invalid cycle detected for type definition %s", defn.Name.GetValue()), defn.GetPosition())
	}
	defn.CycleDepth = depth
	semType := t.resolveBType(defn.GetTypeData().TypeDescriptor.(ast.BType), depth)
	if defn.DeterminedType == nil {
		defn.SetDeterminedType(semType)
		updateSymbolType(t.ctx, defn, semType)
		defn.CycleDepth = -1
		typeData := defn.GetTypeData()
		typeData.Type = semType
		defn.SetTypeData(typeData)
		return semType
	} else {
		// This can happen with recursion
		// We use the first definition we produced
		// and throw away the others
		return defn.GetDeterminedType()
	}
}

func (t *TypeResolver) resolveLiteral(n *ast.BLangLiteral) {
	typeData := n.GetTypeData()
	bType := typeData.TypeDescriptor.(ast.BType)
	var ty semtypes.SemType

	switch bType.BTypeGetTag() {
	case model.TypeTags_INT:
		// INT literals are usually handled via BLangNumericLiteral path
		// but we resolve the type here as well for completeness
		value := n.GetValue().(int64)
		ty = semtypes.IntConst(value)
	case model.TypeTags_BYTE:
		value := n.GetValue().(int64)
		ty = semtypes.IntConst(value)
	case model.TypeTags_BOOLEAN:
		value := n.GetValue().(bool)
		ty = semtypes.BooleanConst(value)
	case model.TypeTags_STRING:
		value := n.GetValue().(string)
		ty = semtypes.StringConst(value)
	case model.TypeTags_NIL:
		ty = &semtypes.NIL
	case model.TypeTags_DECIMAL:
		strValue := n.GetValue().(string)
		r := t.parseDecimalValue(strValue, n.GetPosition())
		ty = semtypes.DecimalConst(*r)
	case model.TypeTags_FLOAT:
		strValue := n.GetValue().(string)
		f := t.parseFloatValue(strValue, n.GetPosition())
		ty = semtypes.FloatConst(f)
	default:
		t.ctx.Unimplemented("unsupported literal type", n.GetPosition())
	}

	// Set on TypeData
	typeData.Type = ty
	n.SetTypeData(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)

	// Update symbol type if this literal has a symbol
	updateSymbolType(t.ctx, n, ty)
}

// parseFloatValue parses a string as float64 with error handling
func (t *TypeResolver) parseFloatValue(strValue string, pos diagnostics.Location) float64 {
	f, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		t.ctx.SyntaxError(fmt.Sprintf("invalid float literal: %s", strValue), pos)
		return 0
	}
	return f
}

// parseDecimalValue parses a string as big.Rat with error handling
func (t *TypeResolver) parseDecimalValue(strValue string, pos diagnostics.Location) *big.Rat {
	r := new(big.Rat)
	if _, ok := r.SetString(strValue); !ok {
		t.ctx.SyntaxError(fmt.Sprintf("invalid decimal literal: %s", strValue), pos)
		return big.NewRat(0, 1)
	}
	return r
}

func (t *TypeResolver) resolveNumericLiteral(n *ast.BLangNumericLiteral) {
	typeData := n.GetTypeData()
	bType := typeData.TypeDescriptor.(ast.BType)
	typeTag := bType.BTypeGetTag()

	var ty semtypes.SemType

	switch n.Kind {
	case model.NodeKind_INTEGER_LITERAL:
		ty = t.resolveIntegerLiteral(n, typeTag)
	case model.NodeKind_DECIMAL_FLOATING_POINT_LITERAL:
		ty = t.resolveDecimalFloatingPointLiteral(n, typeTag)
	case model.NodeKind_HEX_FLOATING_POINT_LITERAL:
		ty = t.resolveHexFloatingPointLiteral(n, typeTag)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected numeric literal kind: %v", n.Kind), n.GetPosition())
		return
	}

	// Set on TypeData
	typeData.Type = ty
	n.SetTypeData(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)

	// Update symbol type if this numeric literal has a symbol
	updateSymbolType(t.ctx, n, ty)
}

func (t *TypeResolver) resolveIntegerLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	value := n.GetValue().(int64)

	switch typeTag {
	case model.TypeTags_INT, model.TypeTags_BYTE:
		return semtypes.IntConst(value)
	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for integer literal", typeTag), n.GetPosition())
		return nil
	}
}

func (t *TypeResolver) resolveDecimalFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	strValue := n.GetValue().(string)

	switch typeTag {
	case model.TypeTags_FLOAT:
		f := t.parseFloatValue(strValue, n.GetPosition())
		return semtypes.FloatConst(f)

	case model.TypeTags_DECIMAL:
		r := t.parseDecimalValue(strValue, n.GetPosition())
		return semtypes.DecimalConst(*r)

	default:
		t.ctx.InternalError(fmt.Sprintf("unexpected type tag %v for decimal floating point literal", typeTag), n.GetPosition())
		return nil
	}
}

func (t *TypeResolver) resolveHexFloatingPointLiteral(n *ast.BLangNumericLiteral, typeTag model.TypeTags) semtypes.SemType {
	t.ctx.Unimplemented("hex floating point literals not supported", n.GetPosition())
	return nil
}

// updateSymbolType updates the symbol's type if the node has an associated symbol.
// This synchronizes the symbol's type with the node's resolved type.
func updateSymbolType(ctx *context.CompilerContext, node ast.BLangNode, ty semtypes.SemType) {
	if nodeWithSymbol, ok := node.(ast.BNodeWithSymbol); ok {
		symbol := nodeWithSymbol.Symbol()
		// symbol resolver should initialize the symbol
		ctx.SetSymbolType(symbol, ty)
	}
}

func (t *TypeResolver) resolveSimpleVariable(node *ast.BLangSimpleVariable) {
	typeData := node.GetTypeData()
	if typeData.TypeDescriptor == nil {
		return
	}

	// Resolve the type descriptor and get the semtype
	semType := t.resolveBType(typeData.TypeDescriptor.(ast.BType), 0)

	// Set on TypeData
	typeData.Type = semType
	node.SetTypeData(typeData)

	// Set on determinedType
	node.SetDeterminedType(semType)

	// Update symbol type
	updateSymbolType(t.ctx, node, semType)
}

func (tr *TypeResolver) resolveBType(btype ast.BType, depth int) semtypes.SemType {
	switch ty := btype.(type) {
	case *ast.BLangValueType:
		switch ty.TypeKind {
		case model.TypeKind_BOOLEAN:
			return &semtypes.BOOLEAN
		case model.TypeKind_INT:
			return &semtypes.INT
		case model.TypeKind_FLOAT:
			return &semtypes.FLOAT
		case model.TypeKind_STRING:
			return &semtypes.STRING
		case model.TypeKind_NIL:
			return &semtypes.NIL
		case model.TypeKind_ANY:
			return &semtypes.ANY
		default:
			tr.ctx.InternalError("unexpected type kind", nil)
			return nil
		}
	case *ast.BLangArrayType:
		defn := ty.Definition
		var semTy semtypes.SemType
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			// Resolve element type and update its TypeData
			elemTypeData := ty.Elemtype
			memberTy := tr.resolveBType(elemTypeData.TypeDescriptor.(ast.BType), depth+1)
			elemTypeData.Type = memberTy
			ty.Elemtype = elemTypeData

			if ty.IsOpenArray() {
				semTy = d.DefineListTypeWrappedWithEnvSemType(tr.typeEnv, memberTy)
			} else {
				length := ty.Sizes[0].(*ast.BLangLiteral).Value.(int)
				semTy = d.DefineListTypeWrappedWithEnvSemTypesInt(tr.typeEnv, []semtypes.SemType{memberTy}, length)
			}
		} else {
			semTy = defn.GetSemType(tr.typeEnv)
		}
		return semTy
	case *ast.BLangUnionTypeNode:
		lhs := tr.resolveBType(ty.Lhs().TypeDescriptor.(ast.BType), depth+1)
		rhs := tr.resolveBType(ty.Rhs().TypeDescriptor.(ast.BType), depth+1)
		return semtypes.Union(lhs, rhs)
	case *ast.BLangErrorTypeNode:
		if ty.IsDistinct() {
			panic("distinct error types not supported")
		}
		if ty.IsTop() {
			return &semtypes.ERROR
		} else {
			detailTy := tr.resolveBType(ty.GetDetailType().TypeDescriptor.(ast.BType), depth+1)
			return semtypes.ErrorDetail(detailTy)
		}
	case *ast.BLangUserDefinedType:
		symbol := ty.Symbol().(*model.SymbolRef)
		defn, ok := tr.typeDefns[*symbol]
		if !ok {
			// This should have been detected by the symbol resolver
			tr.ctx.InternalError("type definition not found", nil)
			return nil
		}
		return tr.resolveTypeDefinition(defn, depth)
	default:
		// TODO: here we need to implement type resolution logic for each type
		tr.ctx.Unimplemented("unsupported type", nil)
		return nil
	}
}
