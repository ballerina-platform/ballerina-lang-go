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
)

type fieldInitState struct {
	initFields map[string]bool
}

func newFieldInitState() *fieldInitState {
	return &fieldInitState{initFields: make(map[string]bool)}
}

func (s *fieldInitState) clone() *fieldInitState {
	c := newFieldInitState()
	for k, v := range s.initFields {
		c.initFields[k] = v
	}
	return c
}

func (s *fieldInitState) markInitialized(name string) {
	s.initFields[name] = true
}

func (s *fieldInitState) isInitialized(name string) bool {
	return s.initFields[name]
}

func mergeFieldStates(s1, s2 *fieldInitState) *fieldInitState {
	if len(s1.initFields) == 0 {
		return s2.clone()
	}
	if len(s2.initFields) == 0 {
		return s1.clone()
	}
	result := newFieldInitState()
	allFields := make(map[string]bool)
	for name := range s1.initFields {
		allFields[name] = true
	}
	for name := range s2.initFields {
		allFields[name] = true
	}
	for name := range allFields {
		init1 := s1.initFields[name]
		init2 := s2.initFields[name]
		result.initFields[name] = init1 && init2
	}
	return result
}

type fieldBlockState struct {
	entry *fieldInitState
	exit  *fieldInitState
}

type uninitFieldAnalyzer struct {
	ctx      *context.CompilerContext
	classDef *ast.BLangClassDefinition
	fcfg     *functionCFG
	states   map[int]*fieldBlockState
	fields   []string // field names that need initialization
}

func newUninitFieldAnalyzer(ctx *context.CompilerContext, classDef *ast.BLangClassDefinition, fcfg *functionCFG) *uninitFieldAnalyzer {
	var fields []string
	for _, field := range classDef.Fields {
		if field.GetInitialExpression() == nil {
			fields = append(fields, field.GetName().GetValue())
		}
	}
	a := &uninitFieldAnalyzer{
		ctx:      ctx,
		classDef: classDef,
		fcfg:     fcfg,
		states:   make(map[int]*fieldBlockState),
		fields:   fields,
	}
	for i := range fcfg.bbs {
		a.states[i] = &fieldBlockState{
			entry: newFieldInitState(),
			exit:  newFieldInitState(),
		}
	}
	return a
}

func (a *uninitFieldAnalyzer) analyze() {
	if len(a.fcfg.bbs) == 0 {
		return
	}
	for _, i := range a.fcfg.topoOrder {
		bb := &a.fcfg.bbs[i]
		entry := a.mergePredecessors(bb)
		if i == 0 {
			for _, name := range a.fields {
				if !entry.isInitialized(name) {
					entry.initFields[name] = false
				}
			}
		}
		a.states[i].entry = entry
		exit := a.analyzeBlock(bb, entry.clone())
		a.states[i].exit = exit
	}
}

func (a *uninitFieldAnalyzer) mergePredecessors(bb *basicBlock) *fieldInitState {
	backedgeSet := make(map[int]bool, len(bb.backedgeParents))
	for _, p := range bb.backedgeParents {
		backedgeSet[p] = true
	}
	var result *fieldInitState
	for _, parentID := range bb.parents {
		if backedgeSet[parentID] {
			continue
		}
		if result == nil {
			result = a.states[parentID].exit.clone()
		} else {
			result = mergeFieldStates(result, a.states[parentID].exit)
		}
	}
	if result == nil {
		result = newFieldInitState()
	}
	return result
}

func (a *uninitFieldAnalyzer) analyzeBlock(bb *basicBlock, state *fieldInitState) *fieldInitState {
	for _, node := range bb.nodes {
		a.analyzeNode(node, state)
	}
	return state
}

func (a *uninitFieldAnalyzer) analyzeNode(node model.Node, state *fieldInitState) {
	if assignment, ok := node.(*ast.BLangAssignment); ok {
		if fieldAccess, ok := assignment.VarRef.(*ast.BLangFieldBaseAccess); ok {
			if isSelfFieldAccess(fieldAccess) {
				state.markInitialized(fieldAccess.Field.Value)
			}
		}
	}
}

func (a *uninitFieldAnalyzer) checkResult() {
	for _, name := range a.fields {
		if !a.isInitializedAtAllTerminals(name) {
			a.reportError(name)
		}
	}
}

func (a *uninitFieldAnalyzer) isInitializedAtAllTerminals(name string) bool {
	for _, bb := range a.fcfg.bbs {
		if !bb.isTerminal() || !bb.isReachable() {
			continue
		}
		if !a.states[bb.id].exit.isInitialized(name) {
			return false
		}
	}
	return true
}

func (a *uninitFieldAnalyzer) reportError(name string) {
	for _, field := range a.classDef.Fields {
		if field.GetName().GetValue() == name {
			a.ctx.SemanticError("field '"+name+"' is not initialized", field.GetPosition())
			return
		}
	}
}

func analyzeUninitializedFields(ctx *context.CompilerContext, pkg *ast.BLangPackage, cfg *PackageCFG) {
	for i := range pkg.ClassDefinitions {
		classDef := &pkg.ClassDefinitions[i]
		var fieldsNeedingInit []string
		for _, field := range classDef.Fields {
			if field.GetInitialExpression() == nil {
				fieldsNeedingInit = append(fieldsNeedingInit, field.GetName().GetValue())
			}
		}
		if len(fieldsNeedingInit) == 0 {
			continue
		}
		if classDef.InitFunction == nil {
			for _, name := range fieldsNeedingInit {
				for _, field := range classDef.Fields {
					if field.GetName().GetValue() == name {
						ctx.SemanticError("field '"+name+"' is not initialized", field.GetPosition())
						break
					}
				}
			}
			continue
		}
		fnCfg, ok := cfg.lookupFunctionCfg(classDef.InitFunction.Symbol())
		if !ok {
			ctx.InternalError("init function CFG not found", classDef.InitFunction.GetPosition())
			continue
		}
		analyzer := newUninitFieldAnalyzer(ctx, classDef, &fnCfg)
		analyzer.analyze()
		analyzer.checkResult()
	}
}
