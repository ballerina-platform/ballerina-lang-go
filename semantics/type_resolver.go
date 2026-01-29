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

// UniformRef is a reference to a symbol that is same whether we are referring to it from different package or same
// so we can efficiently implement lookup logic.
// TODO: string here is just a placeholder
type UniformRef string

func refInPackage(pkg *ast.BLangPackage, name string) UniformRef {
	return UniformRef(name)
}

type (
	TypeResolver struct {
		env *symbolEnv
		ctx *context.CompilerContext
	}

	symbolEnv struct {
		// TODO: keep a map of ast nodes by identifier. When we run into a an indentifier in different package we
		// should be able to lookup that from here and get the semtype there
		typeEnv semtypes.Env
	}

	TypeResolutionResult struct {
		functions map[UniformRef]semtypes.SemType
		// We can't resolve constants fully here because they can have type descriptors so they'll be resolved at semantic analysis
	}

)

var _ ast.Visitor = &TypeResolver{}

func NewTypeResolver(ctx *context.CompilerContext) *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetTypeEnv()}, ctx: ctx}
}

// NewIsolatedTypeResolver is meant for testing so that we can run each test in parallel
func NewIsolatedTypeResolver(ctx *context.CompilerContext) *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetIsolatedTypeEnv()}, ctx: ctx}
}

// ResolveTypes resolves all the type definitions and return a map of all the types of symbols exported by the package.
// After this (for the given package) all the semtypes are known. Semantic analysis will validate and propagate these
// types to the rest of nodes based on semantic information. This means after Resolving types of all the packages
// it is safe use the closed world assumption to optimize type checks.
func (t *TypeResolver) ResolveTypes(pkg *ast.BLangPackage) TypeResolutionResult {
	ast.Walk(t, pkg)
	// TODO: We need to build symbol for function types here (and in the future type decl)
	functions := make(map[UniformRef]semtypes.SemType)
	for _, fn := range pkg.Functions {
		ty := t.resolveFunction(&fn)
		functions[refInPackage(pkg, fn.Name.Value)] = ty
	}
	return TypeResolutionResult{functions: functions}
}

func (t *TypeResolver) resolveFunction(fn *ast.BLangFunction) semtypes.SemType {
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
	paramListTy := paramListDefn.DefineListTypeWrapped(t.env.typeEnv, paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	returnTypeData := fn.GetReturnTypeData()
	if returnTypeData.TypeDescriptor != nil {
		// Already resolved
		returnTy = returnTypeData.Type
	} else {
		returnTy = &semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	return functionDefn.Define(t.env.typeEnv, paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.env.typeEnv, false, false))
}

func (t *TypeResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return t
	}
	ty := t.resolveBType(typeData.TypeDescriptor.(ast.BType))
	typeData.Type = ty
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
	case *ast.BLangArrayType, *ast.BLangBuiltInRefTypeNode, *ast.BLangValueType, *ast.BLangUserDefinedType, *ast.BLangFiniteTypeNode:
		t.ctx.InternalError("unexpected type definition node", n.GetPosition())
		return nil
	case *ast.BLangLiteral:
		t.resolveLiteral(n)
		return nil
	case *ast.BLangNumericLiteral:
		t.resolveNumericLiteral(n)
		return nil
	case *ast.BLangTypeDefinition:
		t.ctx.Unimplemented("type definitions not supported", n.GetPosition())
		return nil
	default:
		return t
	}
}

func (t *TypeResolver) resolveLiteral(n *ast.BLangLiteral) {
	typeData := n.GetBType()
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
	n.SetBType(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)
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
	typeData := n.GetBType()
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
	n.SetBType(typeData)

	// Set on determinedType
	n.SetDeterminedType(ty)
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

func (t *TypeResolver) resolveSimpleVariable(node *ast.BLangSimpleVariable) {
	typeData := node.GetTypeData()
	if typeData.TypeDescriptor == nil {
		return
	}

	// Resolve the type descriptor and get the semtype
	semType := t.resolveBType(typeData.TypeDescriptor.(ast.BType))

	// Set on TypeData
	typeData.Type = semType
	node.SetTypeData(typeData)

	// Set on determinedType
	node.SetDeterminedType(semType)
}

// TODO: do we need to track depth (similar to nBallerina)?
func (tr *TypeResolver) resolveBType(btype ast.BType) semtypes.SemType {
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
			memberTy := tr.resolveBType(elemTypeData.TypeDescriptor.(ast.BType))
			elemTypeData.Type = memberTy
			ty.Elemtype = elemTypeData

			if ty.IsOpenArray() {
				semTy = d.DefineListTypeWrappedWithEnvSemType(tr.env.typeEnv, memberTy)
			} else {
				length := ty.Sizes[0].(*ast.BLangLiteral).Value.(int)
				semTy = d.DefineListTypeWrappedWithEnvSemTypesInt(tr.env.typeEnv, []semtypes.SemType{memberTy}, length)
			}
		} else {
			semTy = defn.GetSemType(tr.env.typeEnv)
		}
		return semTy
	default:
		// TODO: here we need to implement type resolution logic for each type
		tr.ctx.Unimplemented("unsupported type", nil)
		return nil
	}
}
