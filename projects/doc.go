/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package projects provides the Ballerina Project API for loading, compiling,
// and managing Ballerina projects and packages.
//
// It implements the orchestration layer that loads projects from the filesystem,
// parses manifests, resolves dependencies, compiles modules, and generates BIR
// (Ballerina Intermediate Representation) for execution.
//
// # Type Hierarchy
//
// The core type hierarchy mirrors the Ballerina project model:
//
//	Project           Interface representing a Ballerina project (build, single-file, or bala)
//	  Package         A versioned package identified by org/name/version
//	    Module        A single module within a package (one default, zero or more named)
//	      Document    A single .bal source file within a module
//
// Each level has associated types:
//   - Identity:    PackageID, ModuleID, DocumentID (UUID-based)
//   - Descriptor:  PackageDescriptor, ModuleDescriptor (org + name + version)
//   - Config:      PackageConfig, ModuleConfig, DocumentConfig (filesystem layout)
//
// # Project Types
//
// Two concrete project implementations are provided:
//
//   - [BuildProject]: A standard project with a Ballerina.toml manifest
//   - [SingleFileProject]: A standalone .bal file without a manifest
//
// Projects are loaded using [directory.LoadProject] which auto-detects the
// project type based on the path:
//   - Directory with Ballerina.toml -> BuildProject
//   - Single .bal file -> SingleFileProject
//
// # Loading Projects
//
// Use the directory subpackage to load projects from the filesystem:
//
//	import "ballerina-lang-go/projects/directory"
//
//	// Load with default options
//	result, err := directory.LoadProject("./myproject")
//
//	// Load with custom build options
//	buildOpts := projects.NewBuildOptionsBuilder().
//	    WithOffline(true).
//	    WithSkipTests(true).
//	    Build()
//	result, err := directory.LoadProject("./myproject", directory.ProjectLoadConfig{
//	    BuildOptions: &buildOpts,
//	})
//
//	// Load a single .bal file
//	result, err := directory.LoadProject("./main.bal")
//
// The [directory.ProjectLoadConfig] struct allows optional configuration:
//   - BuildOptions: Compilation settings (offline mode, skip tests, etc.)
//   - Future fields can be added without breaking existing callers
//
// # Compilation Pipeline
//
// The compilation pipeline follows a three-phase design:
//
//  1. Load: Project loading creates PackageConfig from the filesystem
//  2. Compile: [PackageCompilation] parses, analyzes, and type-checks all modules
//     in topological order via [PackageResolution]
//  3. CodeGen: [BallerinaBackend] generates BIR from compiled modules
//
// Complete example:
//
//	// Load project
//	result, err := directory.LoadProject("./myproject")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check for loading diagnostics
//	if result.Diagnostics().HasErrors() {
//	    // Handle errors
//	}
//
//	// Get project and package
//	project := result.Project()
//	pkg := project.CurrentPackage()
//
//	// Compile (triggers parsing, type checking, semantic analysis)
//	compilation := pkg.Compilation()
//	if compilation.DiagnosticResult().HasErrors() {
//	    // Handle compilation errors
//	}
//
//	// Generate BIR for execution
//	backend := projects.NewBallerinaBackend(compilation)
//	birPkg := backend.BIR()
//
// # Immutability and Modifiers
//
// Package, Module, and Document are immutable after creation. To create modified
// copies, use the modifier pattern:
//
//	// Modify a document
//	doc := module.Document(docID)
//	updatedDoc := doc.Modify().WithContent(newContent).Apply()
//
//	// Modify a package
//	modifier := pkg.Modify()
//	modifier.AddModule(moduleConfig)
//	newPkg := modifier.Apply()
//
// # Thread Safety
//
// Lazy-initialized fields (compilation, resolution, compiler backends) are
// protected by sync.Once or sync.Mutex for safe concurrent access. Document
// content supports both eager and lazy loading via the DocumentConfig interface.
//
// # Subpackages
//
//   - projects/directory: Project loading from filesystem ([LoadProject], [ProjectLoadConfig])
//   - projects/internal: Internal helpers (ManifestBuilder, PackageConfigCreator)
package projects
