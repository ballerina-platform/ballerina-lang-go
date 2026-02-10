# ballerina-lang-go

## Goals

The primary goal of this effort is to come up with a native Ballerina compiler frontend that is fast, memory-efficient and has a fast startup. Eventually it could replace the current [jBallerina](https://github.com/ballerina-platform/ballerina-lang) compiler frontend.

## Implementation plan

The implementation strategy involves a one-to-one mapping of the jBallerina compiler.

## Usage

### Dependencies

The project is built using the [Go programming language](https://go.dev/). The following dependencies are required:

- [Go 1.24 or later](https://go.dev/dl/)

### Build the CLI

#### Production Build (default)

```bash
go build -o bal ./cli/cmd
```

#### Debug Build

```bash
go build -tags debug -o bal-debug ./cli/cmd
```

### Using Profiling

Profiling is only available in debug builds (compiled with `-tags debug`).

#### Enable Profiling

```bash
# Default profiling port (:6060)
./bal-debug run --prof corpus/bal/subset1/01-boolean/equal1-v.bal

# Custom port
./bal-debug run --prof --prof-addr=:8080 corpus/bal/subset1/01-boolean/equal1-v.bal
```

#### Access Profiling Data

- Web UI: http://localhost:6060/debug/pprof/
- CPU Profile: http://localhost:6060/debug/pprof/profile?seconds=30
- Heap Profile: http://localhost:6060/debug/pprof/heap
- Goroutines: http://localhost:6060/debug/pprof/goroutine

#### Analyze with pprof Tool

```bash
# CPU profiling (30 second sample)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profiling
go tool pprof http://localhost:6060/debug/pprof/heap

# Interactive web UI
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
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

Currently, only single files are supported
E.g

```bash
./bal run --dump-bir corpus/bal/subset1/01-boolean/equal1-v.bal
```

### Testing

To run the tests, use the following command:

```bash
go test ./...
```
