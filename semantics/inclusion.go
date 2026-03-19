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
)

type transitiveInclusionResult struct {
	defns   []model.TypeDefinition
	members []includedMember
}

type includedMember struct {
	objectMember model.ObjectMember
	classField   *ast.BLangSimpleVariable
	classMethod  *ast.BLangFunction
}

func collectTransitiveInclusions(ctx *context.CompilerContext, inclusions []model.SymbolRef) transitiveInclusionResult {
	visited := make(map[model.SymbolRef]bool)
	return collectTransitiveInclusionsInner(ctx, inclusions, visited)
}

func collectTransitiveInclusionsInner(ctx *context.CompilerContext, inclusions []model.SymbolRef, visited map[model.SymbolRef]bool) transitiveInclusionResult {
	var result transitiveInclusionResult
	for _, symRef := range inclusions {
		if visited[symRef] {
			tDefn, ok := ctx.GetTypeDefinition(symRef)
			if ok {
				ctx.SemanticError("cyclic type inclusion", tDefn.(model.Node).GetPosition())
			}
			continue
		}
		visited[symRef] = true
		tDefn, ok := ctx.GetTypeDefinition(symRef)
		if !ok {
			ctx.InternalError("type definition not found for inclusion", nil)
			continue
		}
		switch defn := tDefn.(type) {
		case *ast.BLangTypeDefinition:
			objTy := defn.GetTypeData().TypeDescriptor.(*ast.BLangObjectType)
			sub := collectTransitiveInclusionsInner(ctx, objTy.Inclusions, visited)
			result.defns = append(result.defns, sub.defns...)
			result.members = append(result.members, sub.members...)
			result.defns = append(result.defns, defn)
			for m := range objTy.Members() {
				result.members = append(result.members, includedMember{objectMember: m})
			}
		case *ast.BLangClassDefinition:
			sub := collectTransitiveInclusionsInner(ctx, defn.Inclusions, visited)
			result.defns = append(result.defns, sub.defns...)
			result.members = append(result.members, sub.members...)
			result.defns = append(result.defns, defn)
			for _, f := range defn.Fields {
				field := f.(*ast.BLangSimpleVariable)
				result.members = append(result.members, includedMember{classField: field})
			}
			for _, method := range defn.Methods {
				result.members = append(result.members, includedMember{classMethod: method})
			}
		default:
			ctx.InternalError("unexpected type definition kind for inclusion", tDefn.(model.Node).GetPosition())
		}
	}
	return result
}
