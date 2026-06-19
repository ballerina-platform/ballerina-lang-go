// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"ballerina-lang-go/lsp/protocol"
)

const (
	parseError     = -32700
	invalidRequest = -32600
	methodNotFound = -32601
	invalidParams  = -32602
	internalError  = -32603
)

type Server struct {
	in            io.Reader
	out           io.Writer
	writeMu       sync.Mutex
	snapshots     *SnapshotManager
	notifications chan notification
	shutdown      atomic.Bool
}

type notification struct {
	method string
	params json.RawMessage
}

func NewServer(in io.Reader, out io.Writer) *Server {
	return &Server{
		in:            in,
		out:           out,
		snapshots:     NewSnapshotManager(),
		notifications: make(chan notification, 64),
	}
}

func (s *Server) Run() error {
	go s.runNotifications()
	reader := bufio.NewReader(s.in)
	for {
		payload, err := readMessage(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		var msg protocol.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			s.writeError(nil, parseError, "parse error")
			continue
		}
		if msg.Method == "" {
			continue
		}
		if len(msg.ID) == 0 {
			s.notifications <- notification{method: msg.Method, params: msg.Params}
			continue
		}
		go s.handleRequest(msg)
	}
}

func (s *Server) handleRequest(msg protocol.Message) {
	result, errCode, errMessage := s.dispatchRequest(msg.Method, msg.Params)
	if errCode != 0 {
		s.writeError(msg.ID, errCode, errMessage)
		return
	}
	s.writeResponse(msg.ID, result)
}

func (s *Server) dispatchRequest(method string, params json.RawMessage) (any, int, string) {
	snapshot := s.snapshots.Current()
	_ = snapshot
	switch method {
	case "initialize":
		return protocol.InitializeResult{Capabilities: protocol.ServerCapabilities{TextDocumentSync: protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    1,
			Save:      protocol.SaveOptions{IncludeText: true},
		}}}, 0, ""
	case "shutdown":
		s.shutdown.Store(true)
		return nil, 0, ""
	default:
		return nil, methodNotFound, "method not found: " + method
	}
}

func (s *Server) runNotifications() {
	for notification := range s.notifications {
		s.handleNotification(notification.method, notification.params)
	}
}

func (s *Server) handleNotification(method string, params json.RawMessage) {
	switch method {
	case "initialized":
		return
	case "exit":
		if s.shutdown.Load() {
			os.Exit(0)
		}
		os.Exit(1)
	case "textDocument/didOpen":
		var p protocol.DidOpenTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.updateSnapshot(func(files map[protocol.DocumentURI]SourceFile) (SourceFile, bool) {
			uri := p.TextDocument.URI
			path := pathFromURI(uri)
			file := SourceFile{URI: uri, Path: path, File: path, Version: p.TextDocument.Version, Content: p.TextDocument.Text, Open: true}
			files[uri] = file
			return file, true
		})
	case "textDocument/didChange":
		var p protocol.DidChangeTextDocumentParams
		if decodeParams(params, &p) != nil || len(p.ContentChanges) == 0 {
			return
		}
		s.updateSnapshot(func(files map[protocol.DocumentURI]SourceFile) (SourceFile, bool) {
			uri := p.TextDocument.URI
			file := files[uri]
			if file.URI == "" {
				path := pathFromURI(uri)
				file = SourceFile{URI: uri, Path: path, File: path, Open: true}
			}
			file.Version = p.TextDocument.Version
			file.Content = p.ContentChanges[len(p.ContentChanges)-1].Text
			file.Open = true
			files[uri] = file
			return file, true
		})
	case "textDocument/didClose":
		var p protocol.DidCloseTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.updateSnapshot(func(files map[protocol.DocumentURI]SourceFile) (SourceFile, bool) {
			file, ok := files[p.TextDocument.URI]
			if !ok {
				return SourceFile{}, false
			}
			file.Open = false
			files[p.TextDocument.URI] = file
			return file, true
		})
	case "textDocument/didSave":
		var p protocol.DidSaveTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.updateSnapshot(func(files map[protocol.DocumentURI]SourceFile) (SourceFile, bool) {
			file, ok := files[p.TextDocument.URI]
			if !ok {
				return SourceFile{}, false
			}
			if p.Text != nil {
				file.Content = *p.Text
			}
			files[p.TextDocument.URI] = file
			return file, true
		})
	}
}

func (s *Server) updateSnapshot(update func(map[protocol.DocumentURI]SourceFile) (SourceFile, bool)) {
	old := s.snapshots.Current()
	var changed SourceFile
	var shouldDiagnose bool
	newSnapshot := nextSnapshot(old, func(files map[protocol.DocumentURI]SourceFile) {
		changed, shouldDiagnose = update(files)
	})
	s.snapshots.Publish(newSnapshot)
	if shouldDiagnose {
		go s.publishDiagnostics(newSnapshot, changed)
	}
}

func (s *Server) publishDiagnostics(snapshot *Snapshot, source SourceFile) {
	diagnosticsByURI := runDiagnostics(snapshot, source)
	if !s.snapshots.IsCurrent(snapshot) {
		return
	}
	if _, ok := diagnosticsByURI[source.URI]; !ok {
		diagnosticsByURI[source.URI] = nil
	}
	for uri, diagnostics := range diagnosticsByURI {
		if diagnostics == nil {
			diagnostics = []protocol.Diagnostic{}
		}
		file, ok := snapshot.Files[uri]
		if !ok {
			continue
		}
		version := file.Version
		s.writeNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         uri,
			Version:     &version,
			Diagnostics: diagnostics,
		})
	}
}

func decodeParams(params json.RawMessage, target any) error {
	if len(params) == 0 {
		return nil
	}
	return json.Unmarshal(params, target)
}

func pathFromURI(uri protocol.DocumentURI) string {
	parsed, err := url.Parse(string(uri))
	if err != nil || parsed.Scheme != "file" {
		return normalizePath(string(uri))
	}
	path, err := url.PathUnescape(parsed.Path)
	if err != nil {
		path = parsed.Path
	}
	if runtime.GOOS == "windows" && strings.HasPrefix(path, "/") && len(path) >= 3 && path[2] == ':' {
		path = path[1:]
	}
	return normalizePath(path)
}

func readMessage(reader *bufio.Reader) ([]byte, error) {
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
		if !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(name), "Content-Length") {
			parsed, err := strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return nil, err
			}
			contentLength = parsed
		}
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}
	payload := make([]byte, contentLength)
	_, err := io.ReadFull(reader, payload)
	return payload, err
}

func (s *Server) writeResponse(id json.RawMessage, result any) {
	payload := mustMarshal(result)
	if payload == nil {
		payload = json.RawMessage("null")
	}
	s.writeMessage(struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Result  json.RawMessage `json:"result"`
	}{JSONRPC: "2.0", ID: id, Result: payload})
}

func (s *Server) writeError(id json.RawMessage, code int, message string) {
	if id == nil {
		id = json.RawMessage("null")
	}
	s.writeMessage(struct {
		JSONRPC string                  `json:"jsonrpc"`
		ID      json.RawMessage         `json:"id"`
		Error   *protocol.ResponseError `json:"error"`
	}{JSONRPC: "2.0", ID: id, Error: &protocol.ResponseError{Code: code, Message: message}})
}

func (s *Server) writeNotification(method string, params any) {
	s.writeMessage(struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
	}{JSONRPC: "2.0", Method: method, Params: mustMarshal(params)})
}

func (s *Server) writeMessage(msg any) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	_, _ = fmt.Fprintf(&buffer, "Content-Length: %d\r\n\r\n", len(payload))
	buffer.Write(payload)
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_, _ = s.out.Write(buffer.Bytes())
}

func mustMarshal(value any) json.RawMessage {
	payload, _ := json.Marshal(value)
	return payload
}
