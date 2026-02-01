#!/bin/bash

# Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

# Build script for compiler-tools
# Builds all compiler-tools and places them in the root directory

set -e

ensure_kaitai_compiler() {
    if command -v kaitai-struct-compiler >/dev/null 2>&1; then
        echo "✓ kaitai-struct-compiler is already installed"
        return 0
    fi

    echo "kaitai-struct-compiler not found. Installing..."

    case "$(uname -s)" in
        Darwin)
            if ! command -v brew >/dev/null 2>&1; then
                echo "Error: Homebrew is required to install kaitai-struct-compiler on macOS."
                echo "Install Homebrew first, then re-run this script."
                exit 1
            fi
            brew install kaitai-struct-compiler
            ;;
        Linux)
            curl -fsSLO https://github.com/kaitai-io/kaitai_struct_compiler/releases/download/0.11/kaitai-struct-compiler_0.11_all.deb
            sudo apt-get install ./kaitai-struct-compiler_0.11_all.deb
            ;;
        *)
            echo "Error: Unsupported OS for auto-install of kaitai-struct-compiler: $(uname -s)"
            echo "Please install kaitai-struct-compiler manually and re-run this script."
            exit 1
            ;;
    esac

    if ! command -v kaitai-struct-compiler >/dev/null 2>&1; then
        echo "Error: kaitai-struct-compiler installation did not succeed (still not found on PATH)."
        exit 1
    fi
    echo "✓ kaitai-struct-compiler installed"
}

ensure_kaitai_compiler

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

