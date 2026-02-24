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
	"fmt"
	"sort"
	"strings"
)

// CFGPrettyPrinter prints a PackageCFG in a human-readable format
type CFGPrettyPrinter struct {
	ctx    *context.CompilerContext
	buffer strings.Builder
}

// NewCFGPrettyPrinter creates a new CFG pretty printer
func NewCFGPrettyPrinter(ctx *context.CompilerContext) *CFGPrettyPrinter {
	return &CFGPrettyPrinter{
		ctx: ctx,
	}
}

// Print generates a string representation of the CFG
func (p *CFGPrettyPrinter) Print(cfg *PackageCFG) string {
	p.buffer.Reset()

	// Sort functions by name for deterministic output
	type fnEntry struct {
		ref  model.SymbolRef
		name string
		cfg  functionCFG
	}

	entries := make([]fnEntry, 0, len(cfg.funcCfgs))
	for ref, fnCfg := range cfg.funcCfgs {
		name := p.ctx.SymbolName(ref)
		entries = append(entries, fnEntry{ref: ref, name: name, cfg: fnCfg})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})

	// Print each function's CFG
	for i, entry := range entries {
		if i > 0 {
			p.buffer.WriteString("\n")
		}
		p.printFunctionCFG(entry.name, entry.cfg)
	}

	return p.buffer.String()
}

// printFunctionCFG prints the CFG for a single function
func (p *CFGPrettyPrinter) printFunctionCFG(funcName string, cfg functionCFG) {
	// Print function name
	p.buffer.WriteString("(")
	p.buffer.WriteString(funcName)
	p.buffer.WriteString("\n")

	// Print each basic block
	for _, bb := range cfg.bbs {
		p.printBasicBlock(&bb)
	}

	p.buffer.WriteString(")")
}

// printBasicBlock prints a single basic block
func (p *CFGPrettyPrinter) printBasicBlock(bb *basicBlock) {
	// Print basic block header: (bb<id> (<parents>) (<children>)
	p.buffer.WriteString("  (bb")
	p.buffer.WriteString(fmt.Sprintf("%d", bb.id))
	p.buffer.WriteString(" ")

	// Print parents
	p.buffer.WriteString("(")
	for i, parent := range bb.parents {
		if i > 0 {
			p.buffer.WriteString(" ")
		}
		p.buffer.WriteString(fmt.Sprintf("bb%d", parent))
	}
	p.buffer.WriteString(")")
	p.buffer.WriteString(" ")

	// Print children
	p.buffer.WriteString("(")
	for i, child := range bb.children {
		if i > 0 {
			p.buffer.WriteString(" ")
		}
		p.buffer.WriteString(fmt.Sprintf("bb%d", child))
	}
	p.buffer.WriteString(")")

	// Print nodes if any
	if len(bb.nodes) > 0 {
		p.buffer.WriteString("\n")
		p.printNodes(bb.nodes)
		p.buffer.WriteString("  ")
	}

	p.buffer.WriteString(")\n")
}

// printNodes prints the AST nodes within a basic block
func (p *CFGPrettyPrinter) printNodes(nodes []model.Node) {
	// Create a new PrettyPrinter for each basic block

	for _, node := range nodes {
		// Cast model.Node to ast.BLangNode
		if blangNode, ok := node.(ast.BLangNode); ok {
			printer := &ast.PrettyPrinter{}
			nodeStr := printer.Print(blangNode)

			// Indent each line of the output
			lines := strings.SplitSeq(nodeStr, "\n")
			for line := range lines {
				if line != "" {
					p.buffer.WriteString("    ")
					p.buffer.WriteString(line)
					p.buffer.WriteString("\n")
				}
			}
		}
	}
}
