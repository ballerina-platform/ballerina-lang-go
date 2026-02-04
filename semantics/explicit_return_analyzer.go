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

func AnalyzeExplicitReturn(ctx *context.CompilerContext, pkg *ast.BLangPackage, cfg *PackageCFG) {
	var wg sync.WaitGroup
	panicChan := make(chan interface{}, len(pkg.Functions))
	for i := range pkg.Functions {
		fn := &pkg.Functions[i]
		wg.Add(1)
		go func(f *ast.BLangFunction) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicChan <- r
				}
			}()
			analyzeFunctionExplicitReturn(ctx, f, cfg)
		}(fn)
	}
	wg.Wait()

	close(panicChan)

	for p := range panicChan {
		panic(p)
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
