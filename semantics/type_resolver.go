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
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"math/big"
	"strconv"
)

type (
	TypeResolver struct {
		env *symbolEnv
	}

	symbolEnv struct {
		// TODO: keep a map of ast nodes by identifier. When we run into a an indentifier in different package we
		// should be able to lookup that from here and get the semtype there
		typeEnv semtypes.Env
	}

	TypeResolutionResult struct {
		functions map[string]semtypes.SemType
		// We can't resolve constants fully here because they can have type descriptors so they'll be resolved at semantic analysis
	}
)

var _ ast.Visitor = &TypeResolver{}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetTypeEnv()}}
}

// NewIsolatedTypeResolver is meant for testing so that we can run each test in parallel
func NewIsolatedTypeResolver() *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetIsolatedTypeEnv()}}
}

// ResolveTypes resolves all the type definitions and return a map of all the types of symbols exported by the package.
// After this (for the given package) all the semtypes are known. Semantic analysis will validate and propagate these
// types to the rest of nodes based on semantic information. This means after Resolving types of all the packages
// it is safe use the closed world assumption to optimize type checks.
func (t *TypeResolver) ResolveTypes(pkg *ast.BLangPackage) TypeResolutionResult {
	ast.Walk(t, pkg)
	// TODO: We need to build symbol for function types here (and in the future type decl)
	functions := make(map[string]semtypes.SemType)
	for _, fn := range pkg.Functions {
		ty := t.resolveFunction(&fn)
		functions[fn.Name.Value] = ty
	}
	return TypeResolutionResult{functions: functions}
}

func (t *TypeResolver) resolveFunction(fn *ast.BLangFunction) semtypes.SemType {
	paramTypes := make([]semtypes.SemType, len(fn.RequiredParams))
	for i, param := range fn.RequiredParams {
		paramTypes[i] = param.GetBType().(ast.BType).SemType()
	}
	var restTy semtypes.SemType
	if fn.RestParam != nil {
		panic("unimplemented")
	} else {
		restTy = &semtypes.NEVER
	}
	paramListDefn := semtypes.NewListDefinition()
	paramListTy := paramListDefn.DefineListTypeWrapped(t.env.typeEnv, paramTypes, len(paramTypes), restTy, semtypes.CellMutability_CELL_MUT_NONE)
	var returnTy semtypes.SemType
	if fn.ReturnTypeNode != nil {
		bType := fn.ReturnTypeNode.(ast.BType)
		returnTy = bType.SemType()
	} else {
		returnTy = &semtypes.NIL
	}
	functionDefn := semtypes.NewFunctionDefinition()
	return functionDefn.Define(t.env.typeEnv, paramListTy, returnTy,
		semtypes.FunctionQualifiersFrom(t.env.typeEnv, false, false))
}

func (t *TypeResolver) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		// Done
		return nil
	}
	switch n := node.(type) {
	case *ast.BLangFunction:
		return t
	case *ast.BLangConstant:
		return nil
	case *ast.BLangSimpleVariable:
		t.resolveSimpleVariable(node.(*ast.BLangSimpleVariable))
		return nil
	case *ast.BLangArrayType, *ast.BLangBuiltInRefTypeNode, *ast.BLangValueType, *ast.BLangUserDefinedType, *ast.BLangFiniteTypeNode:
		resolveBType(t.env, n.(ast.BType))
		return nil
	case *ast.BLangLiteral:
		resolveLiteral(t.env, n)
		return nil
	default:
		return t
	}
	panic("unreachable")
}

func resolveLiteral(_ *symbolEnv, n *ast.BLangLiteral) {
	bType := n.GetBType().(ast.BType)
	var ty semtypes.SemType
	switch bType.BTypeGetTag() {
	case model.TypeTags_INT:
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
	// Get value from string
	case model.TypeTags_DECIMAL:
		strValue := n.GetValue().(string)
		r := new(big.Rat)
		if _, ok := r.SetString(strValue); !ok {
			panic("unimplemented")
		}
		ty = semtypes.DecimalConst(*r)
	case model.TypeTags_FLOAT:
		strValue := n.GetValue().(string)
		f, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			panic("unimplemented")
		}
		ty = semtypes.FloatConst(f)
	default:
		panic("unimplemented")
	}
	bType.SetSemType(ty)
	setSemType(n.GetDeterminedType(), ty)
}

func (t *TypeResolver) resolveSimpleVariable(node *ast.BLangSimpleVariable) {
	ty := node.GetBType()
	bType := ty.(ast.BType)
	resolveBType(t.env, bType)
	semType := bType.SemType()
	setSemType(node.GetBType(), semType)
	setSemType(node.GetDeterminedType(), semType)
}

func setSemType(node model.TypeNode, ty semtypes.SemType) {
	if node == nil {
		return
	}
	bType := node.(ast.BType)
	bType.SetSemType(ty)
}

// TODO: do we need to track depth (similar to nBallerina)?
func resolveBType(env *symbolEnv, btype ast.BType) {
	if btype.SemType() != nil {
		// already resolved
		return
	}
	switch ty := btype.(type) {
	case *ast.BLangValueType:
		switch ty.TypeKind {
		case model.TypeKind_BOOLEAN:
			btype.SetSemType(&semtypes.BOOLEAN)
		case model.TypeKind_INT:
			btype.SetSemType(&semtypes.INT)
		case model.TypeKind_FLOAT:
			btype.SetSemType(&semtypes.FLOAT)
		case model.TypeKind_STRING:
			btype.SetSemType(&semtypes.STRING)
		case model.TypeKind_NIL:
			btype.SetSemType(&semtypes.NIL)
		case model.TypeKind_ANY:
			btype.SetSemType(&semtypes.ANY)
		default:
			panic("unexpected")
		}
	case *ast.BLangArrayType:
		defn := ty.Definition
		var t semtypes.SemType
		if defn == nil {
			d := semtypes.NewListDefinition()
			ty.Definition = &d
			resolveBType(env, ty.Elemtype.(ast.BType))
			memberTy := ty.Elemtype.(ast.BType).SemType()
			if ty.IsOpenArray() {
				t = d.DefineListTypeWrappedWithEnvSemType(env.typeEnv, memberTy)
			} else {
				length := ty.Sizes[0].(*ast.BLangLiteral).Value.(int)
				t = d.DefineListTypeWrappedWithEnvSemTypesInt(env.typeEnv, []semtypes.SemType{memberTy}, length)
			}
		} else {
			t = defn.GetSemType(env.typeEnv)
		}

		ty.SetSemType(t)
	default:
		// TODO: here we need to implement type resolution logic for each type
		panic("not implemented")
	}
}
