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
	"ballerina-lang-go/lib"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type symbolResolver interface {
	GetSymbol(name string) (model.Symbol, bool)
	GetPrefixedSymbol(prefix, name string) (model.Symbol, bool)
	AddSymbol(name string, symbol model.Symbol)
	GetPkgID() model.PackageID
	GetScope() model.Scope
	GetCtx() *context.CompilerContext
}

type (
	moduleSymbolResolver struct {
		ctx *context.CompilerContext
		scope *model.ModuleScope
		pkgID model.PackageID
	}

	blockSymbolResolver struct {
		parent symbolResolver
		scope model.Scope
		node ast.BLangNode
	}
)

var _ symbolResolver = &moduleSymbolResolver{}
var _ symbolResolver = &blockSymbolResolver{}

func newModuleSymbolResolver(ctx *context.CompilerContext, pkgID model.PackageID, importedSymbols map[string]model.ExportedSymbolSpace) *moduleSymbolResolver {
	if importedSymbols == nil {
		importedSymbols = make(map[string]model.ExportedSymbolSpace)
	}
	scope := &model.ModuleScope{
		Main:       *model.NewSymbolSpace(pkgID),
		Prefix:     importedSymbols,
		Annotation: *model.NewSymbolSpace(pkgID),
	}
	return &moduleSymbolResolver{
		ctx:   ctx,
		scope: scope,
		pkgID: pkgID,
	}
}

func newFunctionResolver(parent symbolResolver, node ast.BLangNode) *blockSymbolResolver {
	pkgID := parent.GetPkgID()
	parentScope := parent.GetScope()
	scope := model.NewFunctionScope(parentScope, pkgID)
	return &blockSymbolResolver{
		parent: parent,
		scope:  scope,
		node:   node,
	}
}

func newBlockSymbolResolverWithBlockScope(parent symbolResolver, node ast.BLangNode) *blockSymbolResolver {
	pkgID := parent.GetPkgID()
	parentScope := parent.GetScope()
	scope := model.NewBlockScope(parentScope, pkgID)
	return &blockSymbolResolver{
		parent: parent,
		scope:  scope,
		node:   node,
	}
}

func (ms *moduleSymbolResolver) GetSymbol(name string) (model.Symbol, bool) {
	return ms.scope.GetSymbol(name)
}

func (ms *moduleSymbolResolver) GetPkgID() model.PackageID {
	return ms.pkgID
}

func (ms *moduleSymbolResolver) GetScope() model.Scope {
	return ms.scope
}

func (ms *moduleSymbolResolver) GetPrefixedSymbol(prefix, name string) (model.Symbol, bool) {
	return ms.scope.GetPrefixedSymbol(prefix, name)
}

func (ms *moduleSymbolResolver) AddSymbol(name string, symbol model.Symbol) {
	ms.scope.AddSymbol(name, symbol)
}

func (ms *moduleSymbolResolver) GetCtx() *context.CompilerContext {
	return ms.ctx
}

func (bs *blockSymbolResolver) GetSymbol(name string) (model.Symbol, bool) {
	return bs.scope.GetSymbol(name)
}

func (bs *blockSymbolResolver) GetPrefixedSymbol(prefix, name string) (model.Symbol, bool) {
	return bs.scope.GetPrefixedSymbol(prefix, name)
}

func (bs *blockSymbolResolver) AddSymbol(name string, symbol model.Symbol) {
	bs.scope.AddSymbol(name, symbol)
}

func (bs *blockSymbolResolver) GetPkgID() model.PackageID {
	return bs.parent.GetPkgID()
}

func (bs *blockSymbolResolver) GetScope() model.Scope {
	return bs.scope
}

func (bs *blockSymbolResolver) GetCtx() *context.CompilerContext {
	return bs.parent.GetCtx()
}

func addTopLevelSymbol(resolver *moduleSymbolResolver, name string, symbol model.Symbol, pos diagnostics.Location) {
	if _, exists := resolver.GetSymbol(name); exists {
		semanticError(resolver, "redeclared symbol '"+name+"'", pos)
		return
	}
	resolver.AddSymbol(name, symbol)
}

func ResolveSymbols(cx *context.CompilerContext, pkg *ast.BLangPackage, importedSymbols map[string]model.ExportedSymbolSpace) model.ExportedSymbolSpace {
	moduleResolver := newModuleSymbolResolver(cx, *pkg.PackageID, importedSymbols)
	// First add all the top level symbols they can be referred from anywhere
	for _, fn := range pkg.Functions {
		name := fn.Name.Value
		isPublic := fn.FlagSet.Contains(model.Flag_PUBLIC)
		// We are going to fill this in type resolver
		signature := model.FunctionSignature{}
		symbol := model.NewFunctionSymbol(name, signature, isPublic)
		addTopLevelSymbol(moduleResolver, name, &symbol, fn.Name.GetPosition())
	}
	for _, constDef := range pkg.Constants {
		name := constDef.Name.Value
		isPublic := constDef.FlagSet.Contains(model.Flag_PUBLIC)
		symbol := model.NewValueSymbol(name, isPublic, true, false)
		addTopLevelSymbol(moduleResolver, name, &symbol, constDef.Name.GetPosition())
	}
	for _, typeDef := range pkg.TypeDefinitions {
		name := typeDef.Name.Value
		isPublic := typeDef.FlagSet.Contains(model.Flag_PUBLIC)
		symbol := model.NewTypeSymbol(name, isPublic)
		addTopLevelSymbol(moduleResolver, name, &symbol, typeDef.Name.GetPosition())
	}
	// Now properly resolve top level nodes
	for _, fn := range pkg.Functions {
		functionResolver := newFunctionResolver(moduleResolver, &fn)
		resolveFunction(functionResolver, &fn)
	}
	return moduleResolver.scope.Exports()
}

func resolveFunction(functionResolver *blockSymbolResolver, function *ast.BLangFunction) {
	// First add all the parameters to the functionResolver scope
	for _, param := range function.RequiredParams {
		name := param.Name.Value
		symbol := model.NewValueSymbol(name, false, false, true)
		functionResolver.AddSymbol(name, &symbol)
	}

	if function.RestParam != nil {
		if restParam, ok := function.RestParam.(*ast.BLangSimpleVariable); ok {
			name := restParam.Name.Value
			symbol := model.NewValueSymbol(name, false, false, true)
			functionResolver.AddSymbol(name, &symbol)
		}
	}

	ast.Walk(functionResolver, function)
}

// This is a tempary hack since we can only have one import io
func ResolveImports(env semtypes.Env, pkg *ast.BLangPackage) map[string]model.ExportedSymbolSpace {
	result := make(map[string]model.ExportedSymbolSpace)

	for _, imp := range pkg.Imports {
		// Check if this is ballerina/io import
		if imp.OrgName != nil && imp.OrgName.Value == "ballerina" {
			if len(imp.PkgNameComps) == 1 && imp.PkgNameComps[0].Value == "io" {
				// Use alias if available, otherwise use package name
				key := "io"
				if imp.Alias != nil {
					key = imp.Alias.Value
				}
				result[key] = lib.GetIoSymbols(env)
			}
		}
	}

	return result
}

func (bs *blockSymbolResolver) Visit(node ast.BLangNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.BLangFunction:
		// This happens because we visit from the top in [resolveFunction]
		if n == bs.node {
			return bs
		}
		functionResolver := newFunctionResolver(bs, n)
		resolveFunction(functionResolver, n)
		return nil
	case *ast.BLangIf, *ast.BLangWhile, *ast.BLangBlockStmt, *ast.BLangDo:
		return newBlockSymbolResolverWithBlockScope(bs, n)
	case *ast.BLangSimpleVariableDef:
		defineVariable(bs, n.GetVariable())
	case model.InvocationNode:
		resolveFunctionRef(bs, n.(functionRefNode))
	case model.VariableNode:
		referVariable(bs, n.(variableNode))
	}
	return bs
}

type functionRefNode interface {
	GetName() model.IdentifierNode
	GetPosition() diagnostics.Location
	GetPackageAlias() model.IdentifierNode
	SetSymbol(symbol model.Symbol)
}

func resolveFunctionRef[T symbolResolver](resolver T, functionRef functionRefNode) {
	name := functionRef.GetName().GetValue()
	prefix := functionRef.GetPackageAlias().GetValue()
	if prefix != "" {
		symbol, ok := resolver.GetPrefixedSymbol(prefix, name)
		if !ok {
			syntaxError(resolver, "Unknown function: "+name, functionRef.GetPosition())
		}
		functionRef.SetSymbol(symbol)
	} else {
		symbol, ok := resolver.GetSymbol(name)
		if !ok {
			syntaxError(resolver, "Unknown function: "+name, functionRef.GetPosition())
		}
		functionRef.SetSymbol(symbol)
	}
}

type variableNode interface {
	GetName() model.IdentifierNode
	GetPosition() diagnostics.Location
	SetSymbol(symbol model.Symbol)
}

func referVariable[T symbolResolver](resolver T, variable variableNode) {
	name := variable.GetName().GetValue()
	symbol, ok := resolver.GetSymbol(name)
	if !ok {
		syntaxError(resolver, "Unknown variable: "+name, variable.GetPosition())
	}
	variable.SetSymbol(symbol)
}

func defineVariable[T symbolResolver](resolver T, variable model.VariableNode) {
	switch variable := variable.(type) {
	case *ast.BLangSimpleVariable:
		name := variable.Name.Value
		_, ok := resolver.GetSymbol(name)
		if ok {
			syntaxError(resolver, "Variable already defined: "+name, variable.GetPosition())
		}
		symbol := model.NewValueSymbol(name, false, false, true)
		resolver.AddSymbol(name, &symbol)
	default:
		internalError(resolver, "Unsupported variable", variable.GetPosition())
		return
	}
}

func (bs *blockSymbolResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return nil
	}
	td := typeData.TypeDescriptor
	setTypeDescriptorSymbol(bs, td)
	return nil
}

func setTypeDescriptorSymbol[T symbolResolver](resolver T, td model.TypeDescriptor) {
	if bNodeWithSymbol, ok := td.(ast.BNodeWithSymbol); ok {
		symbol := bNodeWithSymbol.Symbol()
		if symbol != nil {
			return
		}
		switch td := td.(type) {
		case *ast.BLangUserDefinedType:
			pkg := td.GetPackageAlias().GetValue()
			tyName := td.GetTypeName().GetValue()
			var symbol model.Symbol
			if pkg != "" {
				symbol, ok = resolver.GetPrefixedSymbol(pkg, tyName)
				if !ok {
					syntaxError(resolver, "Unknown type: "+tyName, td.GetPosition())
				}
			} else {
				symbol, ok = resolver.GetSymbol(tyName)
				if !ok {
					syntaxError(resolver, "Unknown type: "+tyName, td.GetPosition())
				}
			}
			bNodeWithSymbol.SetSymbol(symbol)
		default:
			internalError(resolver, "Unsupported type descriptor", td.GetPosition())
		}
	}
	return
}

func internalError[T symbolResolver](resolver T, message string, pos diagnostics.Location) {
	resolver.GetCtx().InternalError(message, pos)
}

func syntaxError[T symbolResolver](resolver T, message string, pos diagnostics.Location) {
	resolver.GetCtx().SyntaxError(message, pos)
}

func semanticError[T symbolResolver](resolver T, message string, pos diagnostics.Location) {
	resolver.GetCtx().SemanticError(message, pos)
}
