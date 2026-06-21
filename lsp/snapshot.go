// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lsp

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/common/tomlparser"
	"ballerina-lang-go/context"
	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/text"
)

type ProjectKind int

const (
	ProjectKindSingleFile ProjectKind = iota
	ProjectKindBuild
)

const (
	defaultModuleName              = "."
	maxIncrementalSnapshotID       = 100000
	initialSnapshotID        int64 = 0
)

type SourceFile struct {
	URI     protocol.DocumentURI
	Path    string
	File    string
	Version int32
	Content string
	Open    bool
}

type ModuleImport struct {
	Identifier semantics.PackageIdentifier
	ModuleName string
}

type FrontendStage int

const (
	FrontendStageNone FrontendStage = iota
	FrontendStageParsed
	FrontendStageSymbolResolved
	FrontendStageTopLevelTypeResolved
	FrontendStageLocalTypeResolved
	FrontendStageSemanticAnalyzed
	FrontendStageCFGBuilt
	FrontendStageCFGAnalyzed
)

type Module struct {
	Name             string
	Root             string
	PackageID        *model.PackageID
	Files            map[protocol.DocumentURI]SourceFile
	CompilationUnits map[protocol.DocumentURI]*ast.BLangCompilationUnit
	Fingerprint      string
	Stage            FrontendStage
	Imports          []ModuleImport
	ImportedByCU     []semantics.CompilationUnitImports
	ImportedSymbols  map[string]model.ExportedSymbolSpace
	Package          *ast.BLangPackage
	Exported         model.ExportedSymbolSpace
	CFG              *semantics.PackageCFG
}

type Snapshot struct {
	ID        int64
	Kind      ProjectKind
	Root      string
	OrgName   string
	PkgName   string
	Version   string
	Env       *context.CompilerEnvironment
	Files     map[protocol.DocumentURI]SourceFile
	Modules   map[string]*Module
	TopoOrder []string
}

type SnapshotManager struct {
	current *Snapshot
}

func NewSingleFileSnapshotManager(file SourceFile) *SnapshotManager {
	return &SnapshotManager{current: newSingleFileSnapshot(initialSnapshotID, file)}
}

func NewBuildSnapshotManager(root string) *SnapshotManager {
	return &SnapshotManager{current: newBuildSnapshot(initialSnapshotID, nil, root, nil)}
}

func (m *SnapshotManager) Current() *Snapshot {
	return m.current
}

func (m *SnapshotManager) Publish(snapshot *Snapshot) {
	m.current = snapshot
}

func (m *SnapshotManager) IsCurrent(snapshot *Snapshot) bool {
	return m.current == snapshot
}

func nextSingleFileSnapshot(old *Snapshot, file SourceFile) *Snapshot {
	id := nextSnapshotID(old.ID)
	return newSingleFileSnapshot(id, file)
}

func nextBuildSnapshot(old *Snapshot, update func(map[protocol.DocumentURI]SourceFile)) *Snapshot {
	files := make(map[protocol.DocumentURI]SourceFile, len(old.Files))
	for uri, file := range old.Files {
		files[uri] = file
	}
	if update != nil {
		update(files)
	}
	id := nextSnapshotID(old.ID)
	if id == initialSnapshotID {
		return newBuildSnapshot(id, nil, old.Root, files)
	}
	return newBuildSnapshot(id, old, old.Root, files)
}

func newSingleFileSnapshot(id int64, file SourceFile) *Snapshot {
	env := newCompilerEnvironment()
	file.File = file.Path
	files := map[protocol.DocumentURI]SourceFile{file.URI: file}
	module := &Module{
		Name:             defaultModuleName,
		Root:             filepath.Dir(file.Path),
		PackageID:        env.GetDefaultPackage(),
		Files:            files,
		CompilationUnits: make(map[protocol.DocumentURI]*ast.BLangCompilationUnit),
		Fingerprint:      fingerprintFiles(files),
	}
	registerFiles(env, files)
	return &Snapshot{
		ID:        id,
		Kind:      ProjectKindSingleFile,
		Root:      file.Path,
		OrgName:   string(model.ANON_ORG),
		PkgName:   string(model.DEFAULT_PACKAGE),
		Version:   string(model.DEFAULT_VERSION),
		Env:       env,
		Files:     files,
		Modules:   map[string]*Module{defaultModuleName: module},
		TopoOrder: []string{defaultModuleName},
	}
}

func newBuildSnapshot(id int64, old *Snapshot, root string, openFiles map[protocol.DocumentURI]SourceFile) *Snapshot {
	root = normalizePath(root)
	env := newCompilerEnvironment()
	if old != nil && old.Kind == ProjectKindBuild && old.Root == root && old.Env != nil {
		env = old.Env
	}
	orgName, pkgName, version := readPackageDescriptor(root)
	files, modules := scanBuildProject(env, root, orgName, pkgName, version, openFiles)
	if old != nil && old.Env == env {
		for name, module := range modules {
			oldModule := old.Modules[name]
			if oldModule == nil {
				continue
			}
			reuseCompilationUnits(module, oldModule)
			if oldModule.Fingerprint == module.Fingerprint {
				module.Stage = oldModule.Stage
				module.Imports = oldModule.Imports
				module.ImportedByCU = oldModule.ImportedByCU
				module.ImportedSymbols = oldModule.ImportedSymbols
				module.Package = oldModule.Package
				module.Exported = oldModule.Exported
				module.CFG = oldModule.CFG
			}
		}
	}
	registerFiles(env, files)
	return &Snapshot{
		ID:        id,
		Kind:      ProjectKindBuild,
		Root:      root,
		OrgName:   orgName,
		PkgName:   pkgName,
		Version:   version,
		Env:       env,
		Files:     files,
		Modules:   modules,
		TopoOrder: copyStringSlice(oldTopoOrder(old)),
	}
}

func scanBuildProject(env *context.CompilerEnvironment, root, orgName, pkgName, version string, openFiles map[protocol.DocumentURI]SourceFile) (map[protocol.DocumentURI]SourceFile, map[string]*Module) {
	files := make(map[protocol.DocumentURI]SourceFile)
	modules := make(map[string]*Module)
	addModule := func(name, moduleRoot string, packageID *model.PackageID) {
		moduleFiles := scanModuleFiles(moduleRoot, openFiles)
		for uri, file := range moduleFiles {
			files[uri] = file
		}
		modules[name] = &Module{
			Name:             name,
			Root:             moduleRoot,
			PackageID:        packageID,
			Files:            moduleFiles,
			CompilationUnits: make(map[protocol.DocumentURI]*ast.BLangCompilationUnit),
			Fingerprint:      fingerprintFiles(moduleFiles),
		}
	}

	addModule(defaultModuleName, root, modulePackageID(env, orgName, pkgName, version))
	modulesDir := filepath.Join(root, "modules")
	_ = filepath.WalkDir(modulesDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil || !entry.IsDir() || path == modulesDir {
			return nil
		}
		if !hasBalFiles(path) {
			return nil
		}
		rel, err := filepath.Rel(modulesDir, path)
		if err != nil {
			return nil
		}
		modulePart := strings.ReplaceAll(filepath.ToSlash(rel), "/", ".")
		addModule(modulePart, path, modulePackageID(env, orgName, pkgName+"."+modulePart, version))
		return nil
	})
	return files, modules
}

func hasBalFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".bal") {
			return true
		}
	}
	return false
}

func scanModuleFiles(moduleRoot string, openFiles map[protocol.DocumentURI]SourceFile) map[protocol.DocumentURI]SourceFile {
	result := make(map[protocol.DocumentURI]SourceFile)
	entries, err := os.ReadDir(moduleRoot)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".bal") {
				continue
			}
			path := normalizePath(filepath.Join(moduleRoot, entry.Name()))
			uri := uriFromPath(path)
			content, _ := os.ReadFile(path)
			result[uri] = SourceFile{URI: uri, Path: path, File: path, Content: string(content)}
		}
	}
	for uri, file := range openFiles {
		if filepath.Dir(file.Path) != moduleRoot {
			continue
		}
		file.File = file.Path
		result[uri] = file
	}
	return result
}

func reuseCompilationUnits(module *Module, oldModule *Module) {
	if oldModule.CompilationUnits == nil {
		return
	}
	for uri, file := range module.Files {
		oldFile, ok := oldModule.Files[uri]
		if !ok || sourceFileFingerprint(file) != sourceFileFingerprint(oldFile) {
			continue
		}
		if unit := oldModule.CompilationUnits[uri]; unit != nil {
			module.CompilationUnits[uri] = unit
		}
	}
}

func readPackageDescriptor(root string) (string, string, string) {
	orgName := string(model.ANON_ORG)
	pkgName := filepath.Base(root)
	version := string(model.DEFAULT_VERSION)
	toml, err := tomlparser.Read(os.DirFS(root), "Ballerina.toml")
	if err != nil {
		return orgName, pkgName, version
	}
	if value, ok := toml.GetString("package.org"); ok && value != "" {
		orgName = value
	}
	if value, ok := toml.GetString("package.name"); ok && value != "" {
		pkgName = value
	}
	if value, ok := toml.GetString("package.version"); ok && value != "" {
		version = value
	}
	return orgName, pkgName, version
}

func modulePackageID(env *context.CompilerEnvironment, orgName string, moduleName string, version string) *model.PackageID {
	nameParts := strings.Split(moduleName, ".")
	comps := make([]model.Name, len(nameParts))
	for i, part := range nameParts {
		comps[i] = model.Name(part)
	}
	return env.NewPackageID(model.Name(orgName), comps, model.Name(version))
}

func registerFiles(env *context.CompilerEnvironment, files map[protocol.DocumentURI]SourceFile) {
	for _, file := range files {
		env.DiagnosticEnv().RegisterFile(file.File, text.NewStringTextDocument(file.Content))
	}
}

func newCompilerEnvironment() *context.CompilerEnvironment {
	return context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
}

func nextSnapshotID(id int64) int64 {
	if id >= maxIncrementalSnapshotID {
		return initialSnapshotID
	}
	return id + 1
}

func sourceFileFingerprint(file SourceFile) string {
	return fmt.Sprintf("%s\x00%d\x00%s", file.Path, file.Version, file.Content)
}

func fingerprintFiles(files map[protocol.DocumentURI]SourceFile) string {
	uris := make([]string, 0, len(files))
	for uri := range files {
		uris = append(uris, string(uri))
	}
	sort.Strings(uris)
	h := sha256.New()
	for _, rawURI := range uris {
		file := files[protocol.DocumentURI(rawURI)]
		_, _ = fmt.Fprintf(h, "%s\x00", sourceFileFingerprint(file))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func oldTopoOrder(snapshot *Snapshot) []string {
	if snapshot == nil {
		return nil
	}
	return snapshot.TopoOrder
}

func copyStringSlice(values []string) []string {
	if values == nil {
		return nil
	}
	result := make([]string, len(values))
	copy(result, values)
	return result
}

func uriFromPath(path string) protocol.DocumentURI {
	return protocol.DocumentURI((&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String())
}

func isBuildProjectRoot(root string) bool {
	info, err := os.Stat(filepath.Join(root, "Ballerina.toml"))
	return err == nil && !info.IsDir()
}

func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(abs)
}
