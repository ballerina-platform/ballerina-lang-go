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
// This is a Go port of the Java io.ballerina.projects package. It implements
// the orchestration layer that loads projects from the filesystem, parses
// manifests, resolves dependencies, compiles modules, and generates BIR
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
// Concrete project implementations live in the directory subpackage:
//
//   - directory.BuildProject: A standard project with a Ballerina.toml manifest
//   - directory.SingleFileProject: A standalone .bal file without a manifest
//
// Use [directory.LoadProject] to auto-detect the project type from a path,
// or use [directory.LoadBuildProject] / [directory.LoadSingleFileProject] directly.
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
// Example usage:
//
//	project, err := directory.LoadProject(path)
//	pkg := project.CurrentPackage()
//	compilation := pkg.Compilation()
//	backend := NewBallerinaBackend(compilation)
//	birPkg := backend.BIR()
//
// # Immutability and Modifiers
//
// Package, Module, and Document are immutable after creation. To create modified
// copies, use the modifier pattern:
//
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
//   - projects/directory: Concrete project loaders (BuildProject, SingleFileProject)
//   - projects/internal: Internal helpers (ManifestBuilder, PackageConfigCreator)
package projects
