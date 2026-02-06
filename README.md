# ballerina-lang-go

## Goals

The primary goal of this effort is to come up with a native Ballerina compiler frontend that is fast, memory-efficient and has a fast startup. Eventually it could replace the current  [jBallerina](https://github.com/ballerina-platform/ballerina-lang) compiler frontend.

## Implementation plan

The implementation strategy involves a one-to-one mapping of the jBallerina compiler.

## Usage

### Dependencies

The project is built using the [Go programming language](https://go.dev/). The following dependencies are required:
- [Go 1.24 or later](https://go.dev/dl/)

### Build the CLI
```bash
go build -o bal ./cli/cmd
```

### Using the CLI

#### CLI Help
```bash
./bal --help
```

```bash
./bal run --help
```

#### Running a bal source file

Currently, only single files are supported.

E.g 
```bash
./bal run --dump-bir corpus/bal/subset1/01-boolean/equal1-v.bal
```

#### Apply shell completions

##### Zsh

Use these commands to generate and install completions for zsh.

```bash
mkdir -p ~/.zsh/completions
bal completion zsh > ~/.zsh/completions/_bal

# Add to ~/.zshrc (if not already present)
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc

# Apply changes
source ~/.zshrc
```

##### Bash

Use these commands to generate and install completions for bash.

```bash
mkdir -p ~/.local/share/bash-completion/completions
bal completion bash > ~/.local/share/bash-completion/completions/bal

# Apply changes (or restart terminal)
source ~/.local/share/bash-completion/completions/bal
```

##### PowerShell

Use these commands to add completions to your PowerShell profile.

```powershell
bal completion powershell >> $PROFILE
```

### Testing

To run the tests, use the following command:

```bash
go test ./...
```
