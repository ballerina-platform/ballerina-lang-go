# Tree Gen

A code generation tool for generating syntax tree node structs from JSON configuration.

## Building

```bash
cd compiler-tools/tree-gen && go build -o ../../tree-gen
```

## Usage

The tool is invoked via `go generate` directives and accepts the following flags:

- `-config`: Path to JSON configuration file (required)
- `-type`: Node type ("st-node" or "node") (required)
- `-template`: Path to Go template file for main node definitions (required)
- `-output`: Output file path for main node definitions (required)
- `-util-template`: Path to Go template file for utility functions (optional)
- `-util-output`: Output file path for utility functions (optional)

## Configuration Format

The JSON configuration file should have the following structure:

```json
{
  "nodes": [
    "ModulePart",
    "FunctionDefinition",
    "VariableDeclaration"
  ]
}
```

## Running Code Generation

To generate code, run:

```bash
go generate ./...
```

This will execute all `//go:generate` directives in the codebase.

## Generated Files

The tool generates two types of files:

1. **Main node definitions** (e.g., `st-node-gen.go`): Contains the struct definitions for all syntax tree nodes
2. **Utility functions** (e.g., `st-node-util-gen.go`): Contains helper functions like `replaceInner` for working with nodes

## Example

```bash
tree-gen \
  -config ../nodes.json \
  -type st-node \
  -template ../../compiler-tools/tree-gen/templates/st-node.go.tmpl \
  -output st-node-gen.go \
  -util-template ../../compiler-tools/tree-gen/templates/st-node-util.go.tmpl \
  -util-output st-node-util-gen.go
```
