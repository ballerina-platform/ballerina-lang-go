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

// Package langlib provides a test-only helper that compiles the migrated lang
// library bundles into a caller-supplied compiler context, mirroring what the
// project API does via package resolution. Hand-rolled compile drivers (which
// bypass the project API) use it to obtain the implicit-imports map.
package langlib

import (
	"fmt"
	"io/fs"
	"strings"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lib/langlibs"
	"ballerina-lang-go/lib/stdlibs"
	"ballerina-lang-go/model"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/tools/text"
)

type langLib struct {
	org       string
	nameComps []string
	// implicitID, when set, is the key under which the lib is exposed in the
	// implicit-imports map (used without an import statement). When empty the
	// lib requires an explicit import and is exposed via publicSymbols instead.
	implicitID string
	srcFS      fs.FS  // embedded bundle filesystem the source is read from
	balPath    string // path within srcFS
	version    string
}

var migratedLangLibs = []langLib{
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "int"},
		implicitID: "lang.int",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.int/0.0.1/any/lang.int.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "error"},
		implicitID: "lang.error",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.error/0.0.1/any/lang.error.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "string"},
		implicitID: "lang.string",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.string/0.0.1/any/lang.string.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "value"},
		implicitID: "lang.value",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.value/0.0.1/any/lang.value.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "xml"},
		implicitID: "lang.xml",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.xml/0.0.1/any/lang.xml.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "array"},
		implicitID: "lang.array",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.array/0.0.1/any/lang.array.bal",
		version:    "0.0.1",
	},
	{
		org:        "ballerina",
		nameComps:  []string{"lang", "map"},
		implicitID: "lang.map",
		srcFS:      langlibs.FS,
		balPath:    "ballerina/lang.map/0.0.1/any/lang.map.bal",
		version:    "0.0.1",
	},
	{
		org:       "ballerina",
		nameComps: []string{"io"},
		srcFS:     stdlibs.FS,
		balPath:   "ballerina/io/0.0.1/any/io.bal",
		version:   "0.0.1",
	},
}

// ImplicitImports returns the implicit-imports map for a hand-rolled compile
// driver: the still-intrinsic langlibs from semantics.GetImplicitImports plus
// the migrated lang libraries compiled into cx. Compilation happens in cx's
// env so the returned symbol spaces resolve when the driver compiles user code
// in the same context.
func ImplicitImports(cx *context.CompilerContext) (map[string]model.ExportedSymbolSpace, error) {
	result := semantics.GetImplicitImports(cx)
	for _, lib := range migratedLangLibs {
		if lib.implicitID == "" {
			continue
		}
		space, err := compileLangLib(cx, lib)
		if err != nil {
			return nil, err
		}
		result[lib.implicitID] = space
	}
	return result, nil
}

// SeedPublicSymbols compiles the migrated lang libraries into cx and registers
// them in publicSymbols keyed by package identifier, so a hand-rolled driver
// resolves them like any other dependency when the user code imports them.
// This includes the implicitly-used libs (e.g. lang.array) so that user code
// which also imports them explicitly resolves, mirroring the project pipeline
// where bundled lang libs are published to publicSymbols. A nil publicSymbols
// map is initialized and returned.
func SeedPublicSymbols(cx *context.CompilerContext, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace) (map[semantics.PackageIdentifier]model.ExportedSymbolSpace, error) {
	if publicSymbols == nil {
		publicSymbols = make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace)
	}
	for _, lib := range migratedLangLibs {
		space, err := compileLangLib(cx, lib)
		if err != nil {
			return nil, err
		}
		publicSymbols[semantics.PackageIdentifier{
			OrgName:    lib.org,
			ModuleName: strings.Join(lib.nameComps, "."),
		}] = space
	}
	return publicSymbols, nil
}

// libCacheKey identifies a compiled lang library within a single compiler
// context. A library may be requested by both ImplicitImports and
// SeedPublicSymbols; caching ensures it is compiled (and its source file
// registered) exactly once per context.
type libCacheKey struct {
	cx   *context.CompilerContext
	path string
}

var libCache sync.Map // libCacheKey -> model.ExportedSymbolSpace

// compileLangLib compiles a single bundled lang library's source into cx and
// returns its exported symbol space, reusing a previous compilation in the same
// context if present.
func compileLangLib(cx *context.CompilerContext, lib langLib) (model.ExportedSymbolSpace, error) {
	key := libCacheKey{cx: cx, path: lib.balPath}
	if cached, ok := libCache.Load(key); ok {
		return cached.(model.ExportedSymbolSpace), nil
	}
	content, err := fs.ReadFile(lib.srcFS, lib.balPath)
	if err != nil {
		return model.ExportedSymbolSpace{}, fmt.Errorf("langlib: read %s: %w", lib.balPath, err)
	}

	cx.DiagnosticEnv().RegisterFile(lib.balPath, text.NewStringTextDocument(string(content)))
	syntaxTree, err := parser.GetSyntaxTree(cx, lib.balPath, string(content))
	if err != nil {
		return model.ExportedSymbolSpace{}, fmt.Errorf("langlib: parse %s: %w", lib.implicitID, err)
	}
	cu := ast.GetCompilationUnit(cx, syntaxTree)
	if cu == nil {
		return model.ExportedSymbolSpace{}, fmt.Errorf("langlib: AST generation failed for %s", lib.implicitID)
	}
	pkg := ast.ToPackage(cu)

	nameComps := make([]model.Name, len(lib.nameComps))
	for i, c := range lib.nameComps {
		nameComps[i] = model.Name(c)
	}
	pkg.PackageID = cx.NewPackageID(model.Name(lib.org), nameComps, model.Name(lib.version))

	// lang libraries do not themselves import migrated libs, so the
	// still-intrinsic implicit imports are sufficient here.
	imported := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx),
		make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace), lib.org)
	exported := semantics.ResolveSymbols(cx, pkg, imported)
	semantics.ResolveTopLevelNodes(cx, pkg, imported)
	libCache.Store(key, exported)
	return exported, nil
}
