// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type scriptStep struct {
	Message json.RawMessage
	Wait    time.Duration
}

type document struct {
	URI     string
	Path    string
	Version int32
	Content string
}

type client struct {
	stdin  io.WriteCloser
	nextID int64
	docs   map[string]*document
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var scriptPath string
	var serverCommand string
	var drainDuration time.Duration
	var timeout time.Duration
	var pretty bool
	flag.StringVar(&scriptPath, "script", "", "run a JSON script instead of starting the REPL")
	flag.StringVar(&serverCommand, "server", "go run ./cli/cmd lsp", "shell command used to start the LSP server")
	flag.DurationVar(&drainDuration, "drain", 500*time.Millisecond, "time to keep reading server messages after script/REPL completion")
	flag.DurationVar(&timeout, "timeout", 0, "maximum client runtime; 0 disables the timeout")
	flag.BoolVar(&pretty, "pretty", false, "pretty print received JSON messages")
	flag.Parse()

	ctx := context.Background()
	var cancel context.CancelFunc = func() {}
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", serverCommand)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go copyPrefixed(os.Stderr, stderr, "server stderr: ")

	messages := make(chan json.RawMessage, 16)
	errs := make(chan error, 1)
	go readServerMessages(stdout, messages, errs)

	printed := make(chan struct{})
	go func() {
		defer close(printed)
		for msg := range messages {
			writeReceivedMessage(os.Stdout, msg, pretty)
		}
	}()

	c := &client{stdin: stdin, docs: make(map[string]*document)}
	if scriptPath != "" {
		if err := c.runScript(ctx, scriptPath); err != nil {
			return err
		}
	} else if err := c.runREPL(ctx); err != nil {
		return err
	}

	if err := sleep(ctx, drainDuration); err != nil {
		return err
	}
	_ = stdin.Close()

	waitErr := cmd.Wait()
	<-printed

	select {
	case err := <-errs:
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	default:
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return waitErr
}

func (c *client) runScript(ctx context.Context, scriptPath string) error {
	steps, err := readScript(scriptPath)
	if err != nil {
		return err
	}
	for _, step := range steps {
		if step.Wait > 0 {
			if err := sleep(ctx, step.Wait); err != nil {
				return err
			}
			continue
		}
		if err := writeClientMessage(c.stdin, step.Message); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) runREPL(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "LSP client REPL. Type 'help' for commands. Positions are 1-based line/column.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "lsp> ")
		if !scanner.Scan() {
			return scanner.Err()
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "quit" || line == "q" {
			return nil
		}
		if err := c.handleCommand(ctx, line); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}

func (c *client) handleCommand(ctx context.Context, line string) error {
	fields, err := splitCommand(line)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}
	switch fields[0] {
	case "help", "h":
		printHelp()
	case "init", "initialize":
		root := "."
		if len(fields) > 1 {
			root = fields[1]
		}
		return c.request("initialize", map[string]any{"rootUri": uriFromPath(root), "capabilities": map[string]any{}})
	case "initialized":
		return c.notify("initialized", map[string]any{})
	case "open":
		if len(fields) != 2 {
			return fmt.Errorf("usage: open <file>")
		}
		return c.open(fields[1])
	case "reload", "change":
		if len(fields) != 2 {
			return fmt.Errorf("usage: reload <file>")
		}
		return c.reload(fields[1])
	case "insert":
		if len(fields) < 5 {
			return fmt.Errorf("usage: insert <file> <line> <column> <text>")
		}
		lineNo, col, err := parsePosition(fields[2], fields[3])
		if err != nil {
			return err
		}
		return c.replace(fields[1], lineNo, col, lineNo, col, fields[4])
	case "replace":
		if len(fields) < 7 {
			return fmt.Errorf("usage: replace <file> <start-line> <start-column> <end-line> <end-column> <text>")
		}
		startLine, startCol, err := parsePosition(fields[2], fields[3])
		if err != nil {
			return err
		}
		endLine, endCol, err := parsePosition(fields[4], fields[5])
		if err != nil {
			return err
		}
		return c.replace(fields[1], startLine, startCol, endLine, endCol, fields[6])
	case "completion", "complete":
		if len(fields) != 4 {
			return fmt.Errorf("usage: completion <file> <line> <column>")
		}
		return c.positionRequest("textDocument/completion", fields[1], fields[2], fields[3])
	case "definition", "def":
		if len(fields) != 4 {
			return fmt.Errorf("usage: definition <file> <line> <column>")
		}
		return c.positionRequest("textDocument/definition", fields[1], fields[2], fields[3])
	case "code-action", "codeAction":
		if len(fields) != 6 {
			return fmt.Errorf("usage: code-action <file> <start-line> <start-column> <end-line> <end-column>")
		}
		return c.codeAction(fields)
	case "save":
		if len(fields) != 2 {
			return fmt.Errorf("usage: save <file>")
		}
		return c.save(fields[1])
	case "close":
		if len(fields) != 2 {
			return fmt.Errorf("usage: close <file>")
		}
		return c.close(fields[1])
	case "wait":
		if len(fields) != 2 {
			return fmt.Errorf("usage: wait <duration>")
		}
		duration, err := time.ParseDuration(fields[1])
		if err != nil {
			return err
		}
		return sleep(ctx, duration)
	case "raw":
		jsonText := strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
		return writeClientMessage(c.stdin, normalizeClientMessage(json.RawMessage(jsonText)))
	case "shutdown":
		return c.request("shutdown", nil)
	case "exit":
		return c.notify("exit", nil)
	default:
		return fmt.Errorf("unknown command %q", fields[0])
	}
	return nil
}

func (c *client) open(path string) error {
	doc, err := c.load(path)
	if err != nil {
		return err
	}
	return c.notify("textDocument/didOpen", map[string]any{"textDocument": map[string]any{
		"uri": doc.URI, "languageId": "ballerina", "version": doc.Version, "text": doc.Content,
	}})
}

func (c *client) reload(path string) error {
	doc, err := c.document(path)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(doc.Path)
	if err != nil {
		return err
	}
	doc.Content = string(content)
	doc.Version++
	return c.didChange(doc)
}

func (c *client) replace(path string, startLine, startCol, endLine, endCol int, text string) error {
	doc, err := c.document(path)
	if err != nil {
		return err
	}
	start, err := offsetAt(doc.Content, startLine, startCol)
	if err != nil {
		return err
	}
	end, err := offsetAt(doc.Content, endLine, endCol)
	if err != nil {
		return err
	}
	if end < start {
		return fmt.Errorf("end position is before start position")
	}
	doc.Content = doc.Content[:start] + text + doc.Content[end:]
	doc.Version++
	return c.didChange(doc)
}

func (c *client) positionRequest(method, path, lineText, colText string) error {
	doc, err := c.document(path)
	if err != nil {
		return err
	}
	lineNo, col, err := parsePosition(lineText, colText)
	if err != nil {
		return err
	}
	return c.request(method, map[string]any{
		"textDocument": map[string]any{"uri": doc.URI},
		"position":     map[string]any{"line": lineNo, "character": col},
	})
}

func (c *client) codeAction(fields []string) error {
	doc, err := c.document(fields[1])
	if err != nil {
		return err
	}
	startLine, startCol, err := parsePosition(fields[2], fields[3])
	if err != nil {
		return err
	}
	endLine, endCol, err := parsePosition(fields[4], fields[5])
	if err != nil {
		return err
	}
	return c.request("textDocument/codeAction", map[string]any{
		"textDocument": map[string]any{"uri": doc.URI},
		"range": map[string]any{
			"start": map[string]any{"line": startLine, "character": startCol},
			"end":   map[string]any{"line": endLine, "character": endCol},
		},
		"context": map[string]any{"diagnostics": []any{}},
	})
}

func (c *client) save(path string) error {
	doc, err := c.document(path)
	if err != nil {
		return err
	}
	return c.notify("textDocument/didSave", map[string]any{
		"textDocument": map[string]any{"uri": doc.URI}, "text": doc.Content,
	})
}

func (c *client) close(path string) error {
	doc, err := c.document(path)
	if err != nil {
		return err
	}
	delete(c.docs, doc.Path)
	return c.notify("textDocument/didClose", map[string]any{"textDocument": map[string]any{"uri": doc.URI}})
}

func (c *client) load(path string) (*document, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	doc := &document{URI: uriFromPath(abs), Path: abs, Version: 1, Content: string(content)}
	c.docs[abs] = doc
	return doc, nil
}

func (c *client) document(path string) (*document, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	doc := c.docs[abs]
	if doc == nil {
		return nil, fmt.Errorf("document is not open: %s", abs)
	}
	return doc, nil
}

func (c *client) didChange(doc *document) error {
	return c.notify("textDocument/didChange", map[string]any{
		"textDocument":   map[string]any{"uri": doc.URI, "version": doc.Version},
		"contentChanges": []any{map[string]any{"text": doc.Content}},
	})
}

func (c *client) request(method string, params any) error {
	c.nextID++
	return c.send(map[string]any{"jsonrpc": "2.0", "id": c.nextID, "method": method, "params": params})
}

func (c *client) notify(method string, params any) error {
	return c.send(map[string]any{"jsonrpc": "2.0", "method": method, "params": params})
}

func (c *client) send(message any) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return writeClientMessage(c.stdin, payload)
}

func printHelp() {
	fmt.Fprintln(os.Stderr, `Commands:
  init [root]                                send initialize
  initialized                                send initialized notification
  open <file>                                read file and send didOpen
  reload <file>                              reread file from disk and send full didChange
  insert <file> <line> <col> <text>          insert text in the open document
  replace <file> <sl> <sc> <el> <ec> <text>  replace text in the open document
  completion <file> <line> <col>             request completions
  definition <file> <line> <col>             request definition
  code-action <file> <sl> <sc> <el> <ec>     request code actions with empty diagnostics
  save <file>                                send didSave with current in-memory content
  close <file>                               send didClose
  wait <duration>                            wait, e.g. 200ms or 1s
  raw <json>                                 send raw JSON-RPC message
  shutdown                                   request shutdown
  exit                                       send exit notification
  quit                                       stop the client

Text arguments support Go-style quoted strings, so use "text with spaces" or "line1\nline2".
Positions are 1-based in the REPL and converted to LSP's 0-based positions.`)
}

func readScript(path string) ([]scriptStep, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseScript(data)
}

func parseScript(data []byte) ([]scriptStep, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, fmt.Errorf("empty script")
	}

	var rawSteps []json.RawMessage
	if data[0] == '[' {
		if err := json.Unmarshal(data, &rawSteps); err != nil {
			return nil, err
		}
	} else {
		var envelope struct {
			Steps    []json.RawMessage `json:"steps"`
			Messages []json.RawMessage `json:"messages"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			return nil, err
		}
		rawSteps = envelope.Steps
		if len(rawSteps) == 0 {
			rawSteps = envelope.Messages
		}
	}

	steps := make([]scriptStep, 0, len(rawSteps))
	for i, raw := range rawSteps {
		step, err := parseStep(raw)
		if err != nil {
			return nil, fmt.Errorf("step %d: %w", i, err)
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func parseStep(raw json.RawMessage) (scriptStep, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return scriptStep{}, err
	}
	if message, ok := fields["message"]; ok {
		return scriptStep{Message: normalizeClientMessage(message)}, nil
	}
	for _, name := range []string{"wait", "waitMs", "waitMillis"} {
		if value, ok := fields[name]; ok {
			duration, err := parseWait(name, value)
			if err != nil {
				return scriptStep{}, err
			}
			return scriptStep{Wait: duration}, nil
		}
	}
	if _, ok := fields["method"]; ok {
		return scriptStep{Message: normalizeClientMessage(raw)}, nil
	}
	return scriptStep{}, fmt.Errorf("expected a message or wait step")
}

func parseWait(name string, raw json.RawMessage) (time.Duration, error) {
	if name == "wait" {
		var text string
		if err := json.Unmarshal(raw, &text); err == nil {
			return time.ParseDuration(text)
		}
	}
	var millis int64
	if err := json.Unmarshal(raw, &millis); err != nil {
		return 0, err
	}
	return time.Duration(millis) * time.Millisecond, nil
}

func normalizeClientMessage(raw json.RawMessage) json.RawMessage {
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return raw
	}
	if _, ok := object["jsonrpc"]; !ok {
		object["jsonrpc"] = "2.0"
	}
	payload, err := json.Marshal(object)
	if err != nil {
		return raw
	}
	return payload
}

func writeClientMessage(writer io.Writer, message json.RawMessage) error {
	_, err := fmt.Fprintf(writer, "Content-Length: %d\r\n\r\n%s", len(message), message)
	return err
}

func readServerMessages(reader io.Reader, messages chan<- json.RawMessage, errs chan<- error) {
	defer close(messages)
	buf := bufio.NewReader(reader)
	for {
		message, err := readFramedMessage(buf)
		if err != nil {
			errs <- err
			return
		}
		messages <- message
	}
}

func readFramedMessage(reader *bufio.Reader) (json.RawMessage, error) {
	contentLength := -1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok || !strings.EqualFold(strings.TrimSpace(name), "Content-Length") {
			continue
		}
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return nil, err
		}
		contentLength = parsed
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}
	payload := make([]byte, contentLength)
	_, err := io.ReadFull(reader, payload)
	return payload, err
}

func writeReceivedMessage(writer io.Writer, message json.RawMessage, pretty bool) {
	if pretty {
		var value any
		if json.Unmarshal(message, &value) == nil {
			formatted, err := json.MarshalIndent(value, "", "  ")
			if err == nil {
				_, _ = writer.Write(formatted)
				_, _ = writer.Write([]byte("\n"))
				return
			}
		}
	}
	_, _ = writer.Write(message)
	_, _ = writer.Write([]byte("\n"))
}

func copyPrefixed(writer io.Writer, reader io.Reader, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		_, _ = fmt.Fprintln(writer, prefix+scanner.Text())
	}
}

func sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func uriFromPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return (&url.URL{Scheme: "file", Path: filepath.ToSlash(abs)}).String()
}

func parsePosition(lineText string, colText string) (int, int, error) {
	lineNo, err := strconv.Atoi(lineText)
	if err != nil {
		return 0, 0, err
	}
	col, err := strconv.Atoi(colText)
	if err != nil {
		return 0, 0, err
	}
	if lineNo < 1 || col < 1 {
		return 0, 0, fmt.Errorf("line and column must be >= 1")
	}
	return lineNo - 1, col - 1, nil
}

func offsetAt(content string, lineNo int, col int) (int, error) {
	line, column := 0, 0
	for offset, r := range content {
		if line == lineNo && column == col {
			return offset, nil
		}
		if r == '\n' {
			line++
			column = 0
			continue
		}
		column++
	}
	if line == lineNo && column == col {
		return len(content), nil
	}
	return 0, fmt.Errorf("position is outside document")
}

func splitCommand(line string) ([]string, error) {
	var fields []string
	for len(strings.TrimSpace(line)) > 0 {
		line = strings.TrimLeft(line, " \t")
		if strings.HasPrefix(line, "\"") || strings.HasPrefix(line, "`") {
			value, rest, err := readQuoted(line)
			if err != nil {
				return nil, err
			}
			fields = append(fields, value)
			line = rest
			continue
		}
		idx := strings.IndexAny(line, " \t")
		if idx < 0 {
			fields = append(fields, line)
			break
		}
		fields = append(fields, line[:idx])
		line = line[idx+1:]
	}
	return fields, nil
}

func readQuoted(input string) (string, string, error) {
	quote := input[0]
	for i := 1; i < len(input); i++ {
		if input[i] == '\\' && quote == '"' {
			i++
			continue
		}
		if input[i] == quote {
			value, err := strconv.Unquote(input[:i+1])
			if err != nil {
				return "", "", err
			}
			return value, input[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("unterminated quoted string")
}
