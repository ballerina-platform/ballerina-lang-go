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
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
	"sync"
)

// FIXME: we have these panic chaining stuff because we don't handle compiler errors correctly
// need to fix these once we do.

// AnalyzeCFG runs both reachability and explicit return analyses concurrently
// with centralized panic handling.
func AnalyzeCFG(ctx *context.CompilerContext, pkg *ast.BLangPackage, cfg *PackageCFG) {
	var wg sync.WaitGroup
	var panicErr any = nil

	// Run reachability analysis
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()
		analyzeReachability(ctx, cfg)
	}()

	// Run explicit return analysis
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()
		analyzeExplicitReturn(ctx, pkg, cfg)
	}()

	wg.Wait()
	if panicErr != nil {
		panic(panicErr)
	}
}

// analyzeReachability checks for unreachable code in all functions.
// This is now a private function called by AnalyzeCFG.
func analyzeReachability(ctx *context.CompilerContext, cfg *PackageCFG) {
	var wg sync.WaitGroup
	var panicErr any = nil
	for _, fnCfg := range cfg.funcCfgs {
		wg.Add(1)
		go func(fcfg *functionCFG) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			for _, bb := range fcfg.bbs {
				if !bb.isReachable() {
					for _, node := range bb.nodes {
						ctx.SemanticError("unreachable code", node.GetPosition())
					}
				}
			}
		}(&fnCfg)
	}
	wg.Wait()

	if panicErr != nil {
		panic(panicErr)
	}
}

// analyzeExplicitReturn validates that functions with non-nil return types
// have explicit return statements.
// This is now a private function called by AnalyzeCFG.
func analyzeExplicitReturn(ctx *context.CompilerContext, pkg *ast.BLangPackage, cfg *PackageCFG) {
	var wg sync.WaitGroup
	var panicErr any = nil
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		wg.Add(1)
		go func(f *ast.BLangFunction) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicErr = r
				}
			}()
			analyzeFunctionExplicitReturn(ctx, f, cfg)
		}(fn)
	}
	wg.Wait()

	if panicErr != nil {
		panic(panicErr)
	}
}

func analyzeFunctionExplicitReturn(ctx *context.CompilerContext, fn *ast.BLangFunction, cfg *PackageCFG) {
	sym := ctx.GetSymbol(fn.Symbol()).(*model.FunctionSymbol)
	retType := sym.Signature.ReturnType
	if semtypes.IsSubtypeSimple(retType, semtypes.NIL) {
		return
	}

	ref := ctx.RefSymbol(fn.Symbol())
	fnCfg, ok := cfg.funcCfgs[ref]
	if !ok {
		return
	}

	for _, bb := range fnCfg.bbs {
		if !bb.isTerminal() || !bb.isReachable() {
			continue
		}
		if terminalBlockHasReturnOrPanic(bb) {
			continue
		}
		pos := positionForMissingReturn(bb, fn)
		ctx.SemanticError("missing return statement", pos)
	}
}

func terminalBlockHasReturnOrPanic(bb basicBlock) bool {
	if len(bb.nodes) == 0 {
		return false
	}
	last := bb.nodes[len(bb.nodes)-1]
	k := last.GetKind()
	return k == model.NodeKind_RETURN || k == model.NodeKind_PANIC
}

func positionForMissingReturn(bb basicBlock, fn *ast.BLangFunction) diagnostics.Location {
	if len(bb.nodes) > 0 {
		return bb.nodes[len(bb.nodes)-1].GetPosition()
	}
	return fn.GetPosition()
}
