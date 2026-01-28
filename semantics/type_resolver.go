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
)

var (
	_ ast.Visitor = &TypeResolver{}
)

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetTypeEnv()}}
}

// NewIsolatedTypeResolver is meant for testing so that we can run each test in parallel
func NewIsolatedTypeResolver() *TypeResolver {
	return &TypeResolver{env: &symbolEnv{typeEnv: semtypes.GetIsolatedTypeEnv()}}
}

func (t *TypeResolver) ResolveTypes(pkg *ast.BLangPackage) {
	ast.Walk(t, pkg)
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
	case *ast.BLangArrayType:
	case *ast.BLangBuiltInRefTypeNode:
	case *ast.BLangValueType:
	case *ast.BLangUserDefinedType:
	case *ast.BLangFiniteTypeNode:
		resolveBType(t.env, n)
		return nil
	default:
		return t
	}
	panic("unreachable")
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
