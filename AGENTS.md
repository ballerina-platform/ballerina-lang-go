## Coding style

- Don't make symbols public unless asked for or needed
- Constructor methods should provide data for all the fields unless there is default initialization
  - Map values should always be initialized to an empty map

- If multiple structs need to hold same set of fields and implement methods on those fields add \*Base struct and use type inclusion on other structs
  - Make this base struct private
  - Implement the relevant methods on the base struct

- Don't add comments explaining each line of code. If you need to add comments to describe a block of statements then you should extract them to
  a function with meaningful name.

## Symbols
- IMPORTANT: never store `model.Symbol` as the key in a map, always use a `model.SymbolRef`
- Don't call operations on symbols directly instead call them via compiler context

## Interpreter stages

1. Generate Syntax Tree
2. Do symbol resolution
3. Do type resolution
4. Generate Control Flow Graph (CFG)
5. Do semantic analysis
6. Analyze CFG
   - Reachability analysis
   - Explicit return analysis
7. Generate BIR
8. Interpret generated BIR

Stages up to 7 are considered front end.

## Tests

### Corpus tests

- We have 3 kinds of tests indicated by file name in `./corpus/bal`
  1. valid tests (`*-v.bal`)
     These are expected to run end to end and generate output (outputs are indicated with `@output` comments)
  2. error tests (`*-e.bal`)
     These have errors that should be detected before interpreter (error lines are marked with `@error` comments)
  3. panic tests (`*-p.bal`)
     These would trigger runtime panics in the interpreter

- For valid tests for each stage we have expected output defined in `./corpus/$stage` directory. We have corpus tests that generate the actual output and compare against them
  - Each corpus test accepts `-update` flag that will update expected output to match actual output
  - Each corpus tests will run the interpreter up to that stage.
- IMPORTANT: This is the preferred way of testing for any interpreter stage.

## Commands

- You can run interpreter as `go run ./cli/cmd run [flags] <path to bal file>`

##  Profiling
- In order to profile a `.bal` file first you need to get a debug build (`go build -tags debug -o bal-debug ./cli/cmd`)
- Then run the debug build against the bal file `./bal-debug [flag] -prof <path to bal file>`.

### Opening interactive profiler on log running processes
- After running the interpreter with profiling flags run `go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30` and open `localhost:8080` in the browser
