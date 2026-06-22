# LSP client

`compiler-tools/lsp-client` is a small stdio LSP client for validating the Ballerina language server.
It starts the server, lets you mimic editor actions from a REPL, and prints every server response/request/notification as JSON.

## REPL usage

```bash
go run ./compiler-tools/lsp-client -pretty
```

Useful flags:

- `-server 'go run ./cli/cmd lsp'`: command used to start the server.
- `-drain 500ms`: how long to keep reading after quitting the REPL.
- `-timeout 30s`: maximum client runtime; `0` disables the timeout.
- `-pretty`: pretty print received JSON messages.

REPL commands use 1-based line/column positions and convert them to LSP's 0-based positions.

```text
init [root]
initialized
open <file>
reload <file>
insert <file> <line> <col> <text>
replace <file> <start-line> <start-col> <end-line> <end-col> <text>
completion <file> <line> <col>
definition <file> <line> <col>
code-action <file> <start-line> <start-col> <end-line> <end-col>
save <file>
close <file>
wait <duration>
raw <json>
shutdown
exit
quit
```

Example:

```text
init /tmp/sample
initialized
open /tmp/sample/main.bal
completion /tmp/sample/main.bal 3 8
insert /tmp/sample/main.bal 3 8 "io:"
completion /tmp/sample/main.bal 3 11
save /tmp/sample/main.bal
shutdown
exit
quit
```

Text arguments support Go-style quoted strings, so use `"text with spaces"` or `"line1\nline2"`.

## Script mode

For reproducible runs, pass a JSON script instead of using the REPL:

```bash
go run ./compiler-tools/lsp-client -script /path/to/script.json -pretty
```

A script is a JSON array. Each entry is either an LSP message, a wrapped message, or a wait step.
`jsonrpc: "2.0"` is added automatically when omitted.

```json
[
  {
    "id": 1,
    "method": "initialize",
    "params": {
      "rootUri": "file:///tmp/sample",
      "capabilities": {}
    }
  },
  { "method": "initialized", "params": {} },
  { "waitMs": 200 },
  { "id": 2, "method": "shutdown", "params": null },
  { "method": "exit" }
]
```
