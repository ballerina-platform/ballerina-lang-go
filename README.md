# ballerina-lang-go

## Goals

The primary goal of this effort is to come up with a native Ballerina compiler frontend that is fast, memory-efficient and has a fast startup. Eventually it could replace the current  [jBallerina](https://github.com/ballerina-platform/ballerina-lang) compiler frontend.

## Implementation plan

The implementation strategy involves a one-to-one mapping of the jBallerina compiler.

## Usage

### Dependencies

The project is built using the [Go programming language](https://go.dev/). The following dependencies are required:
- [Go 1.24 or later](https://go.dev/dl/)

### Testing

To run the tests, use the following command:

```bash
go test ./...
```
