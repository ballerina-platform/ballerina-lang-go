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
	"maps"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"strings"

	array "ballerina-lang-go/lib/array/compile"
	bInt "ballerina-lang-go/lib/int/compile"
	io "ballerina-lang-go/lib/io/compile"
	bMap "ballerina-lang-go/lib/map/compile"
)

type scopeKind int

const (
	moduleScopeKind scopeKind = iota
	blockScopeKind
)

type symbolResolver interface {
	GetSymbol(name string) (model.SymbolRef, scopeKind, bool)
	ast.Visitor
	GetPrefixedSymbol(prefix, name string) (model.SymbolRef, bool)
	AddSymbol(name string, symbol model.Symbol)
	GetPkgID() model.PackageID
	GetScope() model.Scope
	GetCtx() *context.CompilerContext
}

type (
	moduleSymbolResolver struct {
		ctx   *context.CompilerContext
		scope *model.ModuleScope
		pkgID model.PackageID
	}

	blockSymbolResolver struct {
		parent symbolResolver
		scope  model.BlockLevelScope
		node   ast.BLangNode
	}
)

var (
	_ symbolResolver = &moduleSymbolResolver{}
	_ symbolResolver = &blockSymbolResolver{}
)

func newModuleSymbolResolver(ctx *context.CompilerContext, pkgID model.PackageID, importedSymbols map[string]model.ExportedSymbolSpace) *moduleSymbolResolver {
	if importedSymbols == nil {
		importedSymbols = make(map[string]model.ExportedSymbolSpace)
	}
	scope := &model.ModuleScope{
		Main:       ctx.NewSymbolSpace(pkgID),
		Prefix:     importedSymbols,
		Annotation: ctx.NewSymbolSpace(pkgID),
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
	scope := parent.GetCtx().NewFunctionScope(parentScope, pkgID)
	return &blockSymbolResolver{
		parent: parent,
		scope:  scope,
		node:   node,
	}
}

func newBlockSymbolResolverWithBlockScope(parent symbolResolver, node ast.BLangNode) *blockSymbolResolver {
	pkgID := parent.GetPkgID()
	parentScope := parent.GetScope()
	scope := parent.GetCtx().NewBlockScope(parentScope, pkgID)
	return &blockSymbolResolver{
		parent: parent,
		scope:  scope,
		node:   node,
	}
}

func (ms *moduleSymbolResolver) GetSymbol(name string) (model.SymbolRef, scopeKind, bool) {
	ref, ok := ms.scope.Main.GetSymbol(name)
	return ref, moduleScopeKind, ok
}

func (ms *moduleSymbolResolver) GetPkgID() model.PackageID {
	return ms.pkgID
}

func (ms *moduleSymbolResolver) GetScope() model.Scope {
	return ms.scope
}

func (ms *moduleSymbolResolver) GetPrefixedSymbol(prefix, name string) (model.SymbolRef, bool) {
	return ms.scope.GetPrefixedSymbol(prefix, name)
}

func (ms *moduleSymbolResolver) AddSymbol(name string, symbol model.Symbol) {
	ms.scope.AddSymbol(name, symbol)
}

func (ms *moduleSymbolResolver) GetCtx() *context.CompilerContext {
	return ms.ctx
}

func (bs *blockSymbolResolver) GetSymbol(name string) (model.SymbolRef, scopeKind, bool) {
	ref, ok := bs.scope.MainSpace().GetSymbol(name)
	if ok {
		return ref, blockScopeKind, true
	}
	return bs.parent.GetSymbol(name)
}

func (bs *blockSymbolResolver) GetPrefixedSymbol(prefix, name string) (model.SymbolRef, bool) {
	return bs.parent.GetPrefixedSymbol(prefix, name)
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
	if _, _, exists := resolver.GetSymbol(name); exists {
		semanticError(resolver, "redeclared symbol '"+name+"'", pos)
		return
	}
	resolver.AddSymbol(name, symbol)
}

func addSymbolAndSetOnNode[T symbolResolver](resolver T, name string, symbol model.Symbol, node ast.BNodeWithSymbol) {
	resolver.AddSymbol(name, symbol)
	symRef, _, _ := resolver.GetSymbol(name)
	node.SetSymbol(symRef)
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
		addTopLevelSymbol(moduleResolver, name, symbol, fn.Name.GetPosition())
	}
	for _, constDef := range pkg.Constants {
		name := constDef.Name.Value
		isPublic := constDef.FlagSet.Contains(model.Flag_PUBLIC)
		symbol := model.NewValueSymbol(name, isPublic, true, false)
		addTopLevelSymbol(moduleResolver, name, &symbol, constDef.Name.GetPosition())
	}
	for i := range pkg.TypeDefinitions {
		typeDef := &pkg.TypeDefinitions[i]
		name := typeDef.Name.Value
		isPublic := typeDef.FlagSet.Contains(model.Flag_PUBLIC)
		symbol := model.NewTypeSymbol(name, isPublic)
		addTopLevelSymbol(moduleResolver, name, &symbol, typeDef.Name.GetPosition())
		symRef, _, _ := moduleResolver.GetSymbol(name)
		cx.SetTypeDefinition(symRef, typeDef)
	}
	for i := range pkg.ClassDefinitions {
		classDef := &pkg.ClassDefinitions[i]
		name := classDef.Name.Value
		isPublic := classDef.FlagSet.Contains(model.Flag_PUBLIC)
		symbol := model.NewClassSymbol(name, isPublic)
		addTopLevelSymbol(moduleResolver, name, &symbol, classDef.Name.GetPosition())
		symRef, _, _ := moduleResolver.GetSymbol(name)
		cx.SetTypeDefinition(symRef, classDef)
	}
	ast.Walk(moduleResolver, pkg)
	return moduleResolver.scope.Exports()
}

func resolveFunction(functionResolver *blockSymbolResolver, function *ast.BLangFunction) {
	// First add all the parameters to the functionResolver scope
	for i := range function.RequiredParams {
		param := &function.RequiredParams[i]
		name := param.Name.Value
		symbol := model.NewValueSymbol(name, false, false, true)
		addSymbolAndSetOnNode(functionResolver, name, &symbol, param)
	}

	if function.RestParam != nil {
		if restParam, ok := function.RestParam.(*ast.BLangSimpleVariable); ok {
			name := restParam.Name.Value
			symbol := model.NewValueSymbol(name, false, false, true)
			addSymbolAndSetOnNode(functionResolver, name, &symbol, restParam)
		}
	}

	ast.Walk(functionResolver, function)
}

func resolveLambdaFunction(functionResolver *blockSymbolResolver, parent *blockSymbolResolver, function *ast.BLangFunction) {
	// Check for shadowing on parameters against the enclosing function scope
	for i := range function.RequiredParams {
		param := &function.RequiredParams[i]
		name := param.Name.Value
		if isShadowed(parent, name) {
			semanticError(functionResolver, "Variable already defined: "+name, param.GetPosition())
		}
		symbol := model.NewValueSymbol(name, false, false, true)
		addSymbolAndSetOnNode(functionResolver, name, &symbol, param)
	}

	if function.RestParam != nil {
		if restParam, ok := function.RestParam.(*ast.BLangSimpleVariable); ok {
			name := restParam.Name.Value
			if isShadowed(parent, name) {
				semanticError(functionResolver, "Variable already defined: "+name, restParam.GetPosition())
			}
			symbol := model.NewValueSymbol(name, false, false, true)
			addSymbolAndSetOnNode(functionResolver, name, &symbol, restParam)
		}
	}

	ast.Walk(functionResolver, function)
}

func ResolveImports(ctx *context.CompilerContext, pkg *ast.BLangPackage, implicitImports map[string]model.ExportedSymbolSpace,
	publicSymbols map[PackageIdentifier]model.ExportedSymbolSpace, defaultOrg string,
) map[string]model.ExportedSymbolSpace {
	result := make(map[string]model.ExportedSymbolSpace)

	for _, imp := range pkg.Imports {
		// Check if this is ballerina import
		if imp.OrgName != nil && imp.OrgName.Value == "ballerina" {
			if isIoImport(&imp) {
				// Use alias if available, otherwise use package name
				key := "io"
				if imp.Alias != nil {
					key = imp.Alias.Value
				}
				result[key] = io.GetIoSymbols(ctx)
			} else if isLangImport(&imp, "array") {
				key := "array"
				if imp.Alias != nil {
					key = imp.Alias.Value
				}
				result[key] = array.GetArraySymbols(ctx)
			} else if isLangImport(&imp, "map") {
				key := "map"
				if imp.Alias != nil {
					key = imp.Alias.Value
				}
				result[key] = bMap.GetMapSymbols(ctx)
			} else {
				ctx.Unimplemented("unsupported ballerina import: "+imp.OrgName.Value+"/"+imp.PkgNameComps[0].Value, imp.GetPosition())
			}
		} else {
			id := resolveImportPackageIdentifier(&imp, defaultOrg)
			if symbols, ok := publicSymbols[id]; ok {
				var key string
				if imp.Alias != nil {
					key = imp.Alias.Value
				} else {
					comps := imp.GetPackageName()
					key = comps[len(comps)-1].GetValue()
				}
				result[key] = symbols
			} else {
				ctx.SemanticError("Unknown import: "+id.OrgName+"/"+id.ModuleName, imp.GetPosition())
			}
		}
	}

	maps.Copy(result, implicitImports)

	return result
}

type PackageIdentifier struct {
	OrgName    string
	ModuleName string
}

func resolveImportPackageIdentifier(imp *ast.BLangImportPackage, defaultOrg string) PackageIdentifier {
	nameComps := imp.GetPackageName()
	nameParts := make([]string, len(nameComps))
	for i, name := range nameComps {
		nameParts[i] = name.GetValue()
	}
	moduleName := strings.Join(nameParts, ".")
	var orgName string
	if imp.OrgName == nil || imp.OrgName.GetValue() == "" {
		orgName = defaultOrg
	} else {
		orgName = imp.OrgName.GetValue()
	}
	return PackageIdentifier{orgName, moduleName}
}

func GetImplicitImports(ctx *context.CompilerContext) map[string]model.ExportedSymbolSpace {
	result := make(map[string]model.ExportedSymbolSpace)
	result[array.PackageName] = array.GetArraySymbols(ctx)
	result[bInt.PackageName] = bInt.GetArraySymbols(ctx)
	result[bMap.PackageName] = bMap.GetMapSymbols(ctx)
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
		n.SetScope(functionResolver.scope)
		resolveFunction(functionResolver, n)
		return nil
	case *ast.BLangIf:
		resolver := newBlockSymbolResolverWithBlockScope(bs, n)
		n.SetScope(resolver.scope)
		return resolver
	case *ast.BLangWhile:
		resolver := newBlockSymbolResolverWithBlockScope(bs, n)
		n.SetScope(resolver.scope)
		return resolver
	case *ast.BLangForeach:
		resolveForeachSymbols(bs, n)
		return nil
	case *ast.BLangBlockStmt, *ast.BLangDo:
		return newBlockSymbolResolverWithBlockScope(bs, n)
	case *ast.BLangSimpleVariableDef:
		defineVariable(bs, n.GetVariable(), n.GetVariable().GetFlags().Contains(model.Flag_FINAL))
	case *ast.BLangLambdaFunction:
		fn := n.Function
		name := fn.Name.Value
		signature := model.FunctionSignature{}
		symbol := model.NewFunctionSymbol(name, signature, false)
		addSymbolAndSetOnNode(bs, name, symbol, fn)
		functionResolver := newFunctionResolver(bs, fn)
		fn.SetScope(functionResolver.scope)
		resolveLambdaFunction(functionResolver, bs, fn)
		return nil
	default:
		return visitInnerSymbolResolver(bs, n)
	}
	return bs
}

func visitInnerSymbolResolver[T symbolResolver](resolver T, node ast.BLangNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.BLangFieldBaseAccess:
		if classDef := getEnclosingClassDef(resolver); isSelfFieldAccess(n) && classDef != nil {
			resolveSelfFieldAccess(resolver, n, classDef)
			return nil
		}
	case *ast.BLangMappingConstructorExpr:
		return resolveMappingConstructor(resolver, n)
	case *ast.BLangQueryExpr:
		return newBlockSymbolResolverWithBlockScope(resolver, n)
	case model.InvocationNode:
		if n.GetExpression() != nil {
			createDeferredMethodSymbol(resolver, n)
		} else {
			resolveFunctionRef(resolver, n.(functionRefNode))
		}
	case model.VariableNode:
		referVariable(resolver, n.(variableNode))
	case model.SimpleVariableReferenceNode:
		referSimpleVariableReference(resolver, n)
	case *ast.BLangUserDefinedType:
		referUserDefinedType(resolver, n)
	case *ast.BLangObjectType:
		n.Inclusions = resolveObjectInclusions(resolver, n.PopUnresolvedInclusions())
	}
	return resolver
}

func resolveMappingConstructor[T symbolResolver](resolver T, n *ast.BLangMappingConstructorExpr) ast.Visitor {
	blockResolver := newBlockSymbolResolverWithBlockScope(resolver, n)
	for _, field := range n.Fields {
		if kv, ok := field.(*ast.BLangMappingKeyValueField); ok {
			if !kv.Key.ComputedKey {
				if varRef, ok := kv.Key.Expr.(*ast.BLangSimpleVarRef); ok {
					name := varRef.VariableName.Value
					symbol := model.NewValueSymbol(name, false, false, false)
					addSymbolAndSetOnNode(blockResolver, name, &symbol, varRef)
				}
			}
		}
	}
	return blockResolver
}

// since we don't have type information we can't determine if this is an actual method call or need to be converted
// to a function call.
func createDeferredMethodSymbol[T symbolResolver](resolver T, n model.InvocationNode) {
	invocation := n.(*ast.BLangInvocation)
	name := invocation.Name.GetValue()
	scope := resolver.GetScope().(model.SymbolSpaceProvider)
	invocation.RawSymbol = &deferredMethodSymbol{name: name, space: scope.MainSpace()}
}

type deferredMethodSymbol struct {
	name  string
	space *model.SymbolSpace
}

var _ model.Symbol = &deferredMethodSymbol{}

// IsDeferredMethodSymbol returns true if the symbol is a deferred method symbol
// (a placeholder used during symbol resolution that will be resolved later).
func IsDeferredMethodSymbol(symbol any) bool {
	_, ok := symbol.(*deferredMethodSymbol)
	return ok
}

func (d *deferredMethodSymbol) Name() string {
	panic("method symbol has not been resolved yet")
}

func (d *deferredMethodSymbol) Type() semtypes.SemType {
	panic("method symbol has not been resolved yet")
}

func (d *deferredMethodSymbol) Kind() model.SymbolKind {
	panic("method symbol has not been resolved yet")
}

func (d *deferredMethodSymbol) SetType(semtypes.SemType) {
	panic("method symbol has not been resolved yet")
}

func (d *deferredMethodSymbol) IsPublic() bool {
	panic("method symbol has not been resolved yet")
}

func (d *deferredMethodSymbol) Copy() model.Symbol {
	panic("method symbol has not been resolved yet")
}

func referUserDefinedType[T symbolResolver](resolver T, n *ast.BLangUserDefinedType) {
	name := n.GetTypeName().GetValue()
	var prefix string
	if n.GetPackageAlias() != nil {
		prefix = n.GetPackageAlias().GetValue()
	}
	if prefix != "" {
		symRef, ok := resolver.GetPrefixedSymbol(prefix, name)
		if !ok {
			semanticError(resolver, "Unknown type: "+name, n.GetPosition())
		}
		n.SetSymbol(symRef)
	} else {
		symRef, _, ok := resolver.GetSymbol(name)
		if !ok {
			semanticError(resolver, "Unknown type: "+name, n.GetPosition())
		}
		n.SetSymbol(symRef)
	}
}

type symbolRefNode interface {
	SetSymbol(symbolRef model.SymbolRef)
}

func resolveSymbolRef[T symbolResolver](resolver T, name, prefix string, pos diagnostics.Location, target symbolRefNode) {
	if prefix != "" {
		symRef, ok := resolver.GetPrefixedSymbol(prefix, name)
		if !ok {
			semanticError(resolver, "Unknown symbol: "+name, pos)
		}
		target.SetSymbol(symRef)
	} else {
		symRef, _, ok := resolver.GetSymbol(name)
		if !ok {
			semanticError(resolver, "Unknown symbol: "+name, pos)
		}
		target.SetSymbol(symRef)
	}
}

func referSimpleVariableReference[T symbolResolver](resolver T, n model.SimpleVariableReferenceNode) {
	name := n.GetVariableName().GetValue()
	var prefix string
	if n.GetPackageAlias() != nil {
		prefix = n.GetPackageAlias().GetValue()
	}
	resolveSymbolRef(resolver, name, prefix, n.GetPosition(), n.(ast.BNodeWithSymbol))
}

type functionRefNode interface {
	GetName() model.IdentifierNode
	GetPosition() diagnostics.Location
	GetPackageAlias() model.IdentifierNode
	SetSymbol(symbolRef model.SymbolRef)
}

func resolveFunctionRef[T symbolResolver](resolver T, functionRef functionRefNode) {
	resolveSymbolRef(resolver, functionRef.GetName().GetValue(), functionRef.GetPackageAlias().GetValue(), functionRef.GetPosition(), functionRef)
}

type variableNode interface {
	GetName() model.IdentifierNode
	GetPosition() diagnostics.Location
	SetSymbol(symbolRef model.SymbolRef)
}

func referVariable[T symbolResolver](resolver T, variable variableNode) {
	resolveSymbolRef(resolver, variable.GetName().GetValue(), "", variable.GetPosition(), variable)
}

// isShadowed checks if a name is already defined in an enclosing block scope.
// Mapping constructor scopes contain record keys that are not real variable bindings, so they are skipped.
func isShadowed(resolver *blockSymbolResolver, name string) bool {
	if name == string(model.IGNORE) {
		return false
	}
	current := resolver
	for current != nil {
		// Issue here is mapping constructor treats some of it's keys as simple variable ref; which is wrong but since they are variable they have symbols
		// and we have to resolve them. But they are not real variables
		if _, isMappingScope := current.node.(*ast.BLangMappingConstructorExpr); !isMappingScope {
			if _, ok := current.scope.MainSpace().GetSymbol(name); ok {
				return true
			}
		}
		if next, ok := current.parent.(*blockSymbolResolver); ok {
			current = next
		} else {
			break
		}
	}
	return false
}

func defineVariable(resolver *blockSymbolResolver, variable model.VariableNode, isFinal bool) {
	switch variable := variable.(type) {
	case *ast.BLangSimpleVariable:
		name := variable.Name.Value
		if isShadowed(resolver, name) {
			semanticError(resolver, "Variable already defined: "+name, variable.GetPosition())
		}
		symbol := model.NewValueSymbol(name, false, isFinal, false)
		addSymbolAndSetOnNode(resolver, name, &symbol, variable)
	default:
		internalError(resolver, "Unsupported variable", variable.GetPosition())
		return
	}
}

func resolveForeachSymbols(bs *blockSymbolResolver, n *ast.BLangForeach) {
	resolver := newBlockSymbolResolverWithBlockScope(bs, n)
	n.SetScope(resolver.scope)
	if n.Collection != nil {
		ast.Walk(resolver, n.Collection.(ast.BLangNode))
	}
	if n.VariableDef != nil {
		defineVariable(resolver, n.VariableDef.GetVariable(), true)
		ast.Walk(resolver, n.VariableDef.Var)
	}
	ast.Walk(resolver, &n.Body)
	if n.OnFailClause != nil {
		ast.Walk(resolver, n.OnFailClause)
	}
}

func (bs *blockSymbolResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	if typeData.TypeDescriptor == nil {
		return nil
	}
	td := typeData.TypeDescriptor
	setTypeDescriptorSymbol(bs, td)
	return bs
}

func setTypeDescriptorSymbol[T symbolResolver](resolver T, td model.TypeDescriptor) {
	if bNodeWithSymbol, ok := td.(ast.BNodeWithSymbol); ok {
		if ast.SymbolIsSet(bNodeWithSymbol) {
			return
		}
		switch td := td.(type) {
		case *ast.BLangUserDefinedType:
			pkg := td.GetPackageAlias().GetValue()
			tyName := td.GetTypeName().GetValue()
			var symRef model.SymbolRef
			if pkg != "" {
				symRef, ok = resolver.GetPrefixedSymbol(pkg, tyName)
				if !ok {
					semanticError(resolver, "Unknown type: "+tyName, td.GetPosition())
				}
			} else {
				symRef, _, ok = resolver.GetSymbol(tyName)
				if !ok {
					semanticError(resolver, "Unknown type: "+tyName, td.GetPosition())
				}
			}
			bNodeWithSymbol.SetSymbol(symRef)
		default:
			internalError(resolver, "Unsupported type descriptor", td.GetPosition())
		}
	}
}

func (ms *moduleSymbolResolver) Visit(node ast.BLangNode) ast.Visitor {
	switch n := node.(type) {
	case *ast.BLangFunction:
		name := n.Name.Value
		symRef, _, ok := ms.GetSymbol(name)
		if !ok {
			internalError(ms, "Module level function symbol not found: "+name, n.Name.GetPosition())
		}
		n.SetSymbol(symRef)
		functionResolver := newFunctionResolver(ms, n)
		n.SetScope(functionResolver.scope)
		resolveFunction(functionResolver, n)
		return nil
	case *ast.BLangConstant:
		name := n.Name.Value
		symRef, _, ok := ms.GetSymbol(name)
		if !ok {
			internalError(ms, "Module level constant symbol not found: "+name, n.Name.GetPosition())
		}
		n.SetSymbol(symRef)
		// TODO: create a local scope and resolve the body?
		return ms
	case *ast.BLangTypeDefinition:
		name := n.Name.Value
		symRef, _, ok := ms.GetSymbol(name)
		if !ok {
			internalError(ms, "Module level type symbol not found: "+name, n.Name.GetPosition())
		}
		n.SetSymbol(symRef)
		return ms
	case *ast.BLangClassDefinition:
		name := n.Name.Value
		symRef, _, ok := ms.GetSymbol(name)
		if !ok {
			internalError(ms, "Module level class symbol not found: "+name, n.Name.GetPosition())
		}
		n.SetSymbol(symRef)
		resolveClassDefinition(ms, n)
		return nil
	default:
		return visitInnerSymbolResolver(ms, n)
	}
}

func (ms *moduleSymbolResolver) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return ms
}

// resolveObjectInclusions update the AST node references with correct symbol references. Will add semantic errors if the type
// reference is for something that can't be included. This means after this stage we have the gurantee symbol ref always refer
// to a valid AST node.
func resolveObjectInclusions[T symbolResolver](resolver T, unresolvedInclusions []*ast.BLangUserDefinedType) []model.SymbolRef {
	ctx := resolver.GetCtx()
	inclusions := make([]model.SymbolRef, 0, len(unresolvedInclusions))
	for _, inc := range unresolvedInclusions {
		ast.Walk(resolver, inc)
		symRef := inc.Symbol()
		tDefn, ok := ctx.GetTypeDefinition(symRef)
		if !ok {
			ctx.InternalError("type definition not found for inclusion", inc.GetPosition())
			continue
		}
		switch defn := tDefn.(type) {
		case *ast.BLangTypeDefinition:
			if _, ok := defn.GetTypeData().TypeDescriptor.(*ast.BLangObjectType); !ok {
				ctx.SemanticError("type inclusion must be an object type or class", inc.GetPosition())
				continue
			}
		case *ast.BLangClassDefinition:
		default:
			ctx.InternalError("unexpected type definition kind for inclusion", inc.GetPosition())
			continue
		}
		inclusions = append(inclusions, symRef)
	}
	return inclusions
}

func resolveClassDefinition(ms *moduleSymbolResolver, classDef *ast.BLangClassDefinition) {
	classResolver := newBlockSymbolResolverWithBlockScope(ms, classDef)
	classDef.SetScope(classResolver.scope)

	classDef.Inclusions = resolveObjectInclusions(ms, classDef.PopUnresolvedInclusions())

	for _, field := range classDef.Fields {
		name := field.GetName().GetValue()
		if _, sk, exists := classResolver.GetSymbol(name); exists && sk == blockScopeKind {
			semanticError(classResolver, "redeclared symbol '"+name+"'", field.GetPosition())
			continue
		}
		isPublic := field.GetFlags().Contains(model.Flag_PUBLIC)
		symbol := model.NewValueSymbol(name, isPublic, false, false)
		classResolver.AddSymbol(name, &symbol)
	}

	for methodName := range classDef.Methods {
		method := classDef.Methods[methodName]
		if _, sk, exists := classResolver.GetSymbol(methodName); exists && sk == blockScopeKind {
			semanticError(classResolver, "redeclared symbol '"+methodName+"'", method.Name.GetPosition())
			continue
		}
		isPublic := method.FlagSet.Contains(model.Flag_PUBLIC)
		signature := model.FunctionSignature{}
		symbol := model.NewFunctionSymbol(methodName, signature, isPublic)
		addSymbolAndSetOnNode(classResolver, methodName, symbol, method)
	}

	inc := collectTransitiveInclusions(ms.ctx, classDef.Inclusions)
	for _, m := range inc.members {
		switch {
		case m.objectMember != nil:
			if m.objectMember.MemberKind() != model.ObjectMemberKindField {
				continue
			}
			name := m.objectMember.Name()
			if _, _, exists := classResolver.GetSymbol(name); exists {
				continue
			}
			isPublic := m.objectMember.Visibility() == model.VisibilityPublic
			symbol := model.NewValueSymbol(name, isPublic, false, false)
			classResolver.AddSymbol(name, &symbol)
		case m.classField != nil:
			name := m.classField.Name.Value
			if _, _, exists := classResolver.GetSymbol(name); exists {
				continue
			}
			isPublic := m.classField.FlagSet.Contains(model.Flag_PUBLIC)
			symbol := model.NewValueSymbol(name, isPublic, false, false)
			classResolver.AddSymbol(name, &symbol)
		}
	}

	if classDef.InitFunction != nil {
		signature := model.FunctionSignature{}
		symbol := model.NewFunctionSymbol("init", signature, false)
		addSymbolAndSetOnNode(classResolver, "init", symbol, classDef.InitFunction)
	}

	selfSymbol := model.NewValueSymbol("self", false, false, false)
	classResolver.AddSymbol("self", &selfSymbol)

	for _, field := range classDef.Fields {
		ast.Walk(classResolver, field.(ast.BLangNode))
	}

	if classDef.InitFunction != nil {
		initResolver := newFunctionResolver(classResolver, classDef.InitFunction)
		classDef.InitFunction.SetScope(initResolver.scope)
		resolveFunction(initResolver, classDef.InitFunction)
	}

	for _, method := range classDef.Methods {
		methodResolver := newFunctionResolver(classResolver, method)
		method.SetScope(methodResolver.scope)
		resolveFunction(methodResolver, method)
	}

}

func getEnclosingClassDef(resolver symbolResolver) *ast.BLangClassDefinition {
	for {
		bs, ok := resolver.(*blockSymbolResolver)
		if !ok {
			return nil
		}
		if classDef, ok := bs.node.(*ast.BLangClassDefinition); ok {
			return classDef
		}
		resolver = bs.parent
	}
}

func isSelfFieldAccess(n *ast.BLangFieldBaseAccess) bool {
	varRef, ok := n.Expr.(*ast.BLangSimpleVarRef)
	if !ok {
		return false
	}
	return varRef.VariableName.Value == "self"
}

func resolveSelfFieldAccess[T symbolResolver](resolver T, n *ast.BLangFieldBaseAccess, classDef *ast.BLangClassDefinition) {
	varRef := n.Expr.(*ast.BLangSimpleVarRef)
	referSimpleVariableReference(resolver, varRef)
	fieldName := n.Field.Value
	classScope := classDef.Scope().(model.BlockLevelScope)
	if _, ok := classScope.MainSpace().GetSymbol(fieldName); !ok {
		semanticError(resolver, "undefined member '"+fieldName+"'", n.Field.GetPosition())
	}
}

func internalError[T symbolResolver](resolver T, message string, pos diagnostics.Location) {
	resolver.GetCtx().InternalError(message, pos)
}

func semanticError[T symbolResolver](resolver T, message string, pos diagnostics.Location) {
	resolver.GetCtx().SemanticError(message, pos)
}

// We can't determine if a symbol is actually a method or not without resolivng the expression
// Also we can't really resolve the actual method until we know the type of reciever
// Thus we need to defer the resolution of the method until type resolution
type defferedMethodSymbol struct{}
