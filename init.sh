#!/bin/bash

# Build script for compiler-tools
# Builds all compiler-tools and places them in the root directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Building compiler-tools..."

# Build tree-gen
echo "Building tree-gen..."
cd compiler-tools/tree-gen
go build -o ../../tree-gen
if [ $? -ne 0 ]; then
    echo "Error: Failed to build tree-gen"
    exit 1
fi
cd "$SCRIPT_DIR"
echo "✓ tree-gen built successfully"

# Build update-corpus
echo "Building update-corpus..."
cd compiler-tools/update-corpus
go build -o ../../update-corpus
if [ $? -ne 0 ]; then
    echo "Error: Failed to build update-corpus"
    exit 1
fi
cd "$SCRIPT_DIR"
echo "✓ update-corpus built successfully"

echo ""
echo "All compiler-tools built successfully!"
echo "Executables are available in the root directory:"
echo "  - ./tree-gen"
echo "  - ./update-corpus"

