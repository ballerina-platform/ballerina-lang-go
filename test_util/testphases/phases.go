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

// Package testphases provide utilities to run frontend upto a certain point so that a given
// frontend phase can be validated after that point
package testphases

import (
	"fmt"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/semantics"
)

// Phase represents a frontend compilation phase
type Phase int

const (
	// PhaseParse runs only parsing (syntax tree generation)
	PhaseParse Phase = iota
	// PhaseAST runs parsing + AST generation
	PhaseAST
	// PhaseSymbolResolution runs through symbol resolution
	PhaseSymbolResolution
	// PhaseTypeResolution runs through type resolution
	PhaseTypeResolution
	// PhaseTypeNarrowing runs through type narrowing
	PhaseTypeNarrowing
	// PhaseSemanticAnalysis runs through semantic analysis
	PhaseSemanticAnalysis
	// PhaseCFG runs through CFG generation
	PhaseCFG
	// PhaseCFGAnalysis runs through CFG analysis (reachability, explicit return)
	PhaseCFGAnalysis
	// PhaseDesugar runs through desugaring
	PhaseDesugar
	// PhaseBIR runs through BIR generation
	PhaseBIR
)

// PipelineResult holds the results from running the frontend pipeline
type PipelineResult struct {
	CompilationUnit *ast.BLangCompilationUnit
	Package         *ast.BLangPackage
	CFG             *semantics.PackageCFG
	BIRPackage      *bir.BIRPackage
}

// RunPipeline runs the frontend compilation pipeline up to the specified phase.
// It returns a PipelineResult containing the outputs relevant to that phase.
func RunPipeline(cx *context.CompilerContext, phase Phase, inputPath string) (*PipelineResult, error) {
	result := &PipelineResult{}

	// Create debug context with channel
	debugCtx := &debugcommon.DebugContext{
		Channel: make(chan string),
	}
	// Drain channel in background to prevent blocking
	go func() {
		for range debugCtx.Channel {
			// Discard debug messages
		}
	}()
	defer close(debugCtx.Channel)

	// Phase 1: Parse
	syntaxTree, err := parser.GetSyntaxTree(cx, debugCtx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}
	if phase == PhaseParse {
		return result, nil
	}

	// Phase 2: AST
	result.CompilationUnit = ast.GetCompilationUnit(cx, syntaxTree)
	if result.CompilationUnit == nil {
		return nil, fmt.Errorf("AST generation failed: compilation unit is nil")
	}
	result.Package = ast.ToPackage(result.CompilationUnit)
	if phase == PhaseAST {
		return result, nil
	}

	// Phase 3: Symbol Resolution
	importedSymbols := semantics.ResolveImports(cx, result.Package, semantics.GetImplicitImports(cx))
	semantics.ResolveSymbols(cx, result.Package, importedSymbols)
	if phase == PhaseSymbolResolution {
		return result, nil
	}

	// Phase 4: Type Resolution
	typeResolver := semantics.NewTypeResolver(cx, importedSymbols)
	typeResolver.ResolveTypes(cx, result.Package)
	if phase == PhaseTypeResolution {
		return result, nil
	}

	// Phase 5: Type Narrowing
	semantics.NarrowTypes(cx, result.Package)
	if phase == PhaseTypeNarrowing {
		return result, nil
	}

	// Phase 6: Semantic Analysis
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(result.Package)
	if phase == PhaseSemanticAnalysis {
		return result, nil
	}

	// Phase 7: CFG Generation
	result.CFG = semantics.CreateControlFlowGraph(cx, result.Package)
	if phase == PhaseCFG {
		return result, nil
	}

	// Phase 8: CFG Analysis
	semantics.AnalyzeCFG(cx, result.Package, result.CFG)
	if phase == PhaseCFGAnalysis {
		return result, nil
	}

	// Phase 9: Desugar
	result.Package = desugar.DesugarPackage(cx, result.Package, importedSymbols)
	if phase == PhaseDesugar {
		return result, nil
	}

	// Phase 10: BIR Generation
	result.BIRPackage = bir.GenBir(cx, result.Package)
	return result, nil
}
