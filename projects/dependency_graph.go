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

package projects

import (
	"maps"
	"slices"
	"sync"
)

// DependencyGraph represents a directed graph of dependencies between nodes.
// It supports topological sorting, cycle detection, and dependency traversal.
type DependencyGraph[T comparable] struct {
	rootNode            *T
	dependencies        map[T]map[T]struct{}
	topologicallySorted []T
	cyclicDependencies  [][]T
	sortOnce            sync.Once
}

// nodes returns all nodes in the graph.
func (g *DependencyGraph[T]) nodes() []T {
	return slices.Collect(maps.Keys(g.dependencies))
}

// DirectDependencies returns the direct dependencies of the given node.
// Returns nil if the node does not exist in the graph.
func (g *DependencyGraph[T]) DirectDependencies(node T) []T {
	deps, ok := g.dependencies[node]
	if !ok {
		return nil
	}
	return slices.Collect(maps.Keys(deps))
}

// ToTopologicallySortedList returns nodes in dependency order (dependencies first).
// The result is computed lazily and cached for subsequent calls.
func (g *DependencyGraph[T]) ToTopologicallySortedList() []T {
	g.ensureSorted()
	return slices.Clone(g.topologicallySorted)
}

// FindCycles detects and returns any cycles in the graph.
// Each cycle is represented as a slice of nodes forming the cycle.
// Returns nil if no cycles exist.
func (g *DependencyGraph[T]) FindCycles() [][]T {
	g.ensureSorted()
	if len(g.cyclicDependencies) == 0 {
		return nil
	}
	result := make([][]T, len(g.cyclicDependencies))
	for i, cycle := range g.cyclicDependencies {
		result[i] = slices.Clone(cycle)
	}
	return result
}

func (g *DependencyGraph[T]) ensureSorted() {
	g.sortOnce.Do(func() {
		g.topologicallySorted, g.cyclicDependencies = g.computeTopologicalSort()
	})
}

// computeTopologicalSort performs DFS-based topological sort on the graph.
// Returns nodes in dependency order (dependencies before dependents)
// and any cycles detected.
func (g *DependencyGraph[T]) computeTopologicalSort() ([]T, [][]T) {
	nodes := g.nodes()
	visited := make(map[T]bool, len(nodes))
	inStack := make(map[T]bool, len(nodes))
	var stack []T
	sorted := make([]T, 0, len(nodes))
	var cycles [][]T

	var visit func(vertex T)
	visit = func(vertex T) {
		inStack[vertex] = true
		stack = append(stack, vertex)

		for dep := range g.dependencies[vertex] {
			if inStack[dep] {
				// Found a cycle - extract it from the stack
				if startIdx := slices.Index(stack, dep); startIdx >= 0 {
					cycles = append(cycles, slices.Clone(stack[startIdx:]))
				}
			} else if !visited[dep] {
				visit(dep)
			}
		}
		// Post-order: add to sorted list after processing all dependencies
		sorted = append(sorted, vertex)
		visited[vertex] = true
		delete(inStack, vertex)
		stack = stack[:len(stack)-1]
	}

	for _, node := range nodes {
		if !visited[node] {
			visit(node)
		}
	}

	return sorted, cycles
}

type dependencyGraphBuilder[T comparable] struct {
	rootNode     *T
	dependencies map[T]map[T]struct{}
}

func newDependencyGraphBuilder[T comparable]() *dependencyGraphBuilder[T] {
	return &dependencyGraphBuilder[T]{
		dependencies: make(map[T]map[T]struct{}),
	}
}

func (b *dependencyGraphBuilder[T]) addNode(node T) *dependencyGraphBuilder[T] {
	b.ensureNode(node)
	return b
}

func (b *dependencyGraphBuilder[T]) addDependency(from, to T) *dependencyGraphBuilder[T] {
	// Both nodes are added to the graph if they don't exist.
	b.ensureNode(from)
	b.ensureNode(to)
	b.dependencies[from][to] = struct{}{}
	return b
}

// build creates the immutable DependencyGraph from the builder's current state.
// The builder can continue to be used after build is called.
func (b *dependencyGraphBuilder[T]) build() *DependencyGraph[T] {
	cloned := make(map[T]map[T]struct{}, len(b.dependencies))
	for k, v := range b.dependencies {
		cloned[k] = maps.Clone(v)
	}

	var rootCopy *T
	if b.rootNode != nil {
		r := *b.rootNode
		rootCopy = &r
	}

	return &DependencyGraph[T]{
		rootNode:     rootCopy,
		dependencies: cloned,
	}
}

func (b *dependencyGraphBuilder[T]) ensureNode(node T) {
	if _, ok := b.dependencies[node]; !ok {
		b.dependencies[node] = make(map[T]struct{})
	}
}
