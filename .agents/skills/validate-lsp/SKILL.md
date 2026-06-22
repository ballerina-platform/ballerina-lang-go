# Validate LSP implementation

Use this skill when changing or reviewing code under `lsp/` or the `lsp` CLI command.

## Fast tests

Run focused LSP tests first:

```bash
go test ./lsp
```

Run broader tests when the change may affect compiler behaviour used by LSP:

```bash
go test ./...
```

## Server logging

Enable LSP logs while manually validating behaviour:

```bash
BAL_LSP_LOG=1 go run ./cli/cmd lsp
```

When using an editor or the compiler-tool client, set `BAL_LSP_LOG=1` on the server process. Logs are written to `.bal/lsp.log` under the project/root directory selected by the server. Use these logs to verify request routing, snapshot updates, diagnostics publishing, and response flow.

## REPL client

Use the compiler-tool client to mimic editor actions over stdio:

```bash
BAL_LSP_LOG=1 go run ./compiler-tools/lsp-client -pretty
```

By default it starts:

```bash
go run ./cli/cmd lsp
```

Override the server command when needed:

```bash
BAL_LSP_LOG=1 go run ./compiler-tools/lsp-client -server './bal lsp' -pretty
```

Common REPL flow:

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

Positions in the REPL are 1-based line/column values. The client converts them to LSP's 0-based positions. The client prints every server response/request/notification to stdout as JSON; compare that output with `.bal/lsp.log` to validate both protocol-visible behaviour and internal server decisions.

Useful commands:

- `open <file>`: read file from disk and send `textDocument/didOpen`.
- `reload <file>`: reread file from disk and send full `textDocument/didChange`.
- `insert <file> <line> <col> <text>` / `replace ...`: mutate the in-memory open document and send full `didChange`.
- `completion <file> <line> <col>` / `definition <file> <line> <col>` / `code-action ...`: send editor-like requests.
- `raw <json>`: send a raw JSON-RPC message for cases not covered by a command.

## Script mode

For reproducible non-interactive runs, pass a JSON script:

```bash
BAL_LSP_LOG=1 go run ./compiler-tools/lsp-client -script /path/to/script.json -pretty
```

A script is a JSON array of LSP messages and wait steps, for example `{ "wait": "500ms" }` or `{ "waitMs": 200 }`.
