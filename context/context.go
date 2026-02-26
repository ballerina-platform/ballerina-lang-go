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

package context

import (
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type CompilerContext struct {
	env         *CompilerEnvironment
	diagnostics []diagnostics.Diagnostic
}

func (this *CompilerContext) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	return this.env.NewSymbolSpace(packageId)
}

func (this *CompilerContext) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return this.env.NewFunctionScope(parent, pkg)
}

func (this *CompilerContext) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return this.env.NewBlockScope(parent, pkg)
}

func (this *CompilerContext) GetSymbol(symbol model.SymbolRef) model.Symbol {
	return this.env.GetSymbol(symbol)
}

// CreateNarrowedSymbol create a narrowed symbol for the given baseRef symbol. IMPORTANT: baseRef must be the actual symbol
// not a narrowed symbol.
func (this *CompilerContext) CreateNarrowedSymbol(baseRef model.SymbolRef) model.SymbolRef {
	return this.env.CreateNarrowedSymbol(baseRef)
}

func (this *CompilerContext) UnnarrowedSymbol(symbol model.SymbolRef) model.SymbolRef {
	return this.env.UnnarrowedSymbol(symbol)
}

func (this *CompilerContext) SymbolName(symbol model.SymbolRef) string {
	return this.env.GetSymbol(symbol).Name()
}

func (this *CompilerContext) SymbolType(symbol model.SymbolRef) semtypes.SemType {
	return this.env.GetSymbol(symbol).Type()
}

func (this *CompilerContext) SymbolKind(symbol model.SymbolRef) model.SymbolKind {
	return this.env.GetSymbol(symbol).Kind()
}

func (this *CompilerContext) SymbolIsPublic(symbol model.SymbolRef) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this *CompilerContext) SetSymbolType(symbol model.SymbolRef, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this *CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.env.GetDefaultPackage()
}

func (this *CompilerContext) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return this.env.NewPackageID(orgName, nameComps, version)
}

func (this *CompilerContext) Unimplemented(message string, pos diagnostics.Location) {
	this.addDiagnostic("UNIMPLEMENTED_ERROR", diagnostics.Fatal, message, pos)
}

func (this *CompilerContext) InternalError(message string, pos diagnostics.Location) {
	this.addDiagnostic("INTERNAL_ERROR", diagnostics.Fatal, message, pos)
}

func (this *CompilerContext) SyntaxError(message string, pos diagnostics.Location) {
	this.addDiagnostic("SYNTAX_ERROR", diagnostics.Error, message, pos)
}

func (this *CompilerContext) SemanticError(message string, pos diagnostics.Location) {
	this.addDiagnostic("SEMANTIC_ERROR", diagnostics.Error, message, pos)
}

func (this *CompilerContext) addDiagnostic(code string, severity diagnostics.DiagnosticSeverity, message string, pos diagnostics.Location) {
	diagnostic := diagnostics.CreateDiagnostic(diagnostics.NewDiagnosticInfo(&code, message, severity), pos)
	this.diagnostics = append(this.diagnostics, diagnostic)
}

func (this *CompilerContext) HasDiagnostics() bool {
	return len(this.diagnostics) > 0
}

func (this *CompilerContext) HasErrors() bool {
	for _, diag := range this.diagnostics {
		if diag.DiagnosticInfo().Severity() == diagnostics.Error {
			return true
		}
	}
	return false
}

func (this *CompilerContext) Diagnostics() []diagnostics.Diagnostic {
	return this.diagnostics
}

func NewCompilerContext(env *CompilerEnvironment) *CompilerContext {
	return &CompilerContext{
		env: env,
	}
}

// GetTypeEnv returns the type environment for this context
func (this *CompilerContext) GetTypeEnv() semtypes.Env {
	return this.env.GetTypeEnv()
}

func (this *CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	return this.env.GetNextAnonymousTypeKey(packageID)
}
