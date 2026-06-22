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
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"ballerina-lang-go/lsp/protocol"
	"ballerina-lang-go/model"
)

const (
	parseError     = -32700
	invalidRequest = -32600
	methodNotFound = -32601
	invalidParams  = -32602
	internalError  = -32603

	indexingProgressToken protocol.ProgressToken = "indexing"
)

type Server struct {
	in                              io.Reader
	out                             io.Writer
	snapshots                       map[string]*SnapshotManager
	root                            string
	shutdown                        bool
	workDoneProgress                bool
	watchedFilesDynamicRegistration bool
	progressCreated                 bool
	nextServerRequestID             int64
}

func NewServer(in io.Reader, out io.Writer) *Server {
	return &Server{
		in:        in,
		out:       out,
		snapshots: make(map[string]*SnapshotManager),
	}
}

func (s *Server) Run() error {
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
			s.handleNotification(msg.Method, msg.Params)
			continue
		}
		s.handleRequest(msg)
	}
}

func (s *Server) handleRequest(msg protocol.Message) {
	logLS(s.root, "request handling method=%s id=%s", msg.Method, string(msg.ID))
	result, errCode, errMessage := s.dispatchRequest(msg.Method, msg.Params)
	if errCode != 0 {
		s.writeError(msg.ID, errCode, errMessage)
		return
	}
	s.writeResponse(msg.ID, result)
}

func (s *Server) dispatchRequest(method string, params json.RawMessage) (any, int, string) {
	logLS(s.root, "request received method=%s", method)
	switch method {
	case "initialize":
		var p protocol.InitializeParams
		if err := decodeParams(params, &p); err != nil {
			return nil, invalidParams, "invalid initialize params"
		}
		s.initializeSnapshots(p)
		return protocol.InitializeResult{Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    1,
				Save:      protocol.SaveOptions{IncludeText: true},
			},
			CompletionProvider: &protocol.CompletionOptions{TriggerCharacters: []string{":", "."}},
			DefinitionProvider: true,
			ReferenceProvider:  true,
			CodeActionProvider: true,
		}}, 0, ""
	case "textDocument/completion":
		var p protocol.CompletionParams
		if err := decodeParams(params, &p); err != nil {
			return nil, invalidParams, "invalid completion params"
		}
		return s.completion(p), 0, ""
	case "textDocument/definition":
		var p protocol.DefinitionParams
		if err := decodeParams(params, &p); err != nil {
			return nil, invalidParams, "invalid definition params"
		}
		return s.definition(p), 0, ""
	case "textDocument/references":
		var p protocol.ReferenceParams
		if err := decodeParams(params, &p); err != nil {
			return nil, invalidParams, "invalid references params"
		}
		return s.references(p), 0, ""
	case "textDocument/codeAction":
		var p protocol.CodeActionParams
		if err := decodeParams(params, &p); err != nil {
			return nil, invalidParams, "invalid code action params"
		}
		return s.codeActions(p), 0, ""
	case "shutdown":
		s.shutdown = true
		return nil, 0, ""
	default:
		return nil, methodNotFound, "method not found: " + method
	}
}

func (s *Server) handleNotification(method string, params json.RawMessage) {
	logLS(s.root, "notification received method=%s", method)
	switch method {
	case "initialized":
		if s.root == "" {
			return
		}
		s.registerWatchedFiles()
		manager := s.snapshots[s.root]
		if manager != nil && manager.Current().Kind == ProjectKindBuild {
			s.publishDiagnostics(manager, manager.Current(), SourceFile{Path: s.root})
		}
		return
	case "exit":
		if s.shutdown {
			os.Exit(0)
		}
		os.Exit(1)
	case "textDocument/didOpen":
		var p protocol.DidOpenTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		uri := p.TextDocument.URI
		path := pathFromURI(uri)
		file := SourceFile{URI: uri, Path: path, File: path, Version: p.TextDocument.Version, Content: p.TextDocument.Text, Open: true}
		s.updateSnapshot(file, func(files map[protocol.DocumentURI]SourceFile) SourceFile {
			files[uri] = file
			return file
		})
	case "textDocument/didChange":
		var p protocol.DidChangeTextDocumentParams
		if decodeParams(params, &p) != nil || len(p.ContentChanges) == 0 {
			return
		}
		uri := p.TextDocument.URI
		file := s.sourceFile(uri)
		if file.URI == "" {
			return
		}
		file.Version = p.TextDocument.Version
		file.Content = p.ContentChanges[len(p.ContentChanges)-1].Text
		file.File = file.Path
		file.Open = true
		s.updateSnapshot(file, func(files map[protocol.DocumentURI]SourceFile) SourceFile {
			files[uri] = file
			return file
		})
	case "textDocument/didClose":
		var p protocol.DidCloseTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		uri := p.TextDocument.URI
		file := s.sourceFile(uri)
		if file.URI == "" || !file.Open {
			return
		}
		file.Open = false
		if content, err := os.ReadFile(file.Path); err == nil {
			file.Content = string(content)
		}
		s.updateSnapshot(file, func(files map[protocol.DocumentURI]SourceFile) SourceFile {
			delete(files, uri)
			return file
		})
		return
	case "textDocument/didSave":
		var p protocol.DidSaveTextDocumentParams
		if decodeParams(params, &p) != nil {
			return
		}
		uri := p.TextDocument.URI
		file := s.sourceFile(uri)
		if file.URI == "" {
			return
		}
		content, err := os.ReadFile(file.Path)
		if err != nil || string(content) == file.Content {
			return
		}
		file.Content = string(content)
		file.Open = true
		s.updateSnapshot(file, func(files map[protocol.DocumentURI]SourceFile) SourceFile {
			files[uri] = file
			return file
		})
	case "workspace/didChangeWatchedFiles":
		var p protocol.DidChangeWatchedFilesParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.handleWatchedFileChanges(p)
	case "workspace/didRenameFiles":
		var p protocol.RenameFilesParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.handleRenamedFiles(p)
	case "workspace/didCreateFiles":
		var p protocol.CreateFilesParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.handleCreatedFiles(p)
	case "workspace/didDeleteFiles":
		var p protocol.DeleteFilesParams
		if decodeParams(params, &p) != nil {
			return
		}
		s.handleDeletedFiles(p)
	}
}

func (s *Server) updateSnapshot(source SourceFile, update func(map[protocol.DocumentURI]SourceFile) SourceFile) {
	key := s.snapshotKey(source)
	if s.root != "" && s.root != key {
		logLS(s.root, "document mapped source=%s snapshotKey=%s projectLog=%s", source.Path, key, filepath.Join(key, ".bal", "lsp.log"))
	}
	manager := s.snapshotManager(key, source)
	old := manager.Current()
	logLS(key, "snapshot update start key=%s kind=%s oldID=%d source=%s", key, projectKindString(old.Kind), old.ID, source.Path)
	var newSnapshot *Snapshot
	var changed SourceFile
	if old.Kind == ProjectKindBuild {
		newSnapshot = nextBuildSnapshot(old, func(files map[protocol.DocumentURI]SourceFile) {
			changed = update(files)
		})
		invalidateChangedDependents(old, newSnapshot, changed.URI)
	} else {
		changed = update(map[protocol.DocumentURI]SourceFile{})
		newSnapshot = nextSingleFileSnapshot(old, changed)
	}
	manager.Publish(newSnapshot)
	logLS(key, "snapshot update published key=%s kind=%s newID=%d modules=%d files=%d", key, projectKindString(newSnapshot.Kind), newSnapshot.ID, len(newSnapshot.Modules), len(newSnapshot.Files))
	s.publishDiagnostics(manager, newSnapshot, changed)
}

func (s *Server) handleWatchedFileChanges(params protocol.DidChangeWatchedFilesParams) {
	for _, change := range params.Changes {
		path := pathFromURI(change.URI)
		root := s.projectRootForWatchedPath(path)
		if root == "" || !isRelevantWatchedPath(root, path) {
			continue
		}
		if change.Type == protocol.FileChangeTypeChanged && strings.HasSuffix(path, ".bal") {
			s.refreshChangedBuildFile(root, change.URI)
			continue
		}
		s.refreshBuildProject(root, true, SourceFile{URI: change.URI, Path: path, File: path})
	}
}

func (s *Server) handleRenamedFiles(params protocol.RenameFilesParams) {
	roots := make(map[string]SourceFile)
	for _, file := range params.Files {
		for _, uri := range []protocol.DocumentURI{file.OldURI, file.NewURI} {
			path := pathFromURI(uri)
			root := s.projectRootForWatchedPath(path)
			if root == "" || !isRelevantWatchedPath(root, path) {
				continue
			}
			roots[root] = SourceFile{URI: uri, Path: path, File: path}
		}
	}
	for root, source := range roots {
		s.refreshBuildProject(root, true, source)
	}
}

func (s *Server) handleCreatedFiles(params protocol.CreateFilesParams) {
	for _, file := range params.Files {
		path := pathFromURI(file.URI)
		root := s.projectRootForWatchedPath(path)
		if root != "" && isRelevantWatchedPath(root, path) {
			s.refreshBuildProject(root, true, SourceFile{URI: file.URI, Path: path, File: path})
		}
	}
}

func (s *Server) handleDeletedFiles(params protocol.DeleteFilesParams) {
	for _, file := range params.Files {
		path := pathFromURI(file.URI)
		root := s.projectRootForWatchedPath(path)
		if root != "" && isRelevantWatchedPath(root, path) {
			s.refreshBuildProject(root, true, SourceFile{URI: file.URI, Path: path, File: path})
		}
	}
}

func (s *Server) refreshChangedBuildFile(root string, uri protocol.DocumentURI) {
	manager := s.snapshots[root]
	if manager == nil || manager.Current().Kind != ProjectKindBuild {
		return
	}
	if file, ok := manager.Current().Files[uri]; ok && file.Open {
		return
	}
	s.refreshBuildProject(root, false, SourceFile{URI: uri, Path: pathFromURI(uri), File: pathFromURI(uri)})
}

func (s *Server) refreshBuildProject(root string, clean bool, source SourceFile) {
	manager := s.snapshots[root]
	if manager == nil {
		manager = NewBuildSnapshotManager(root)
		s.snapshots[root] = manager
	}
	old := manager.Current()
	if old.Kind != ProjectKindBuild {
		return
	}
	openFiles := openSnapshotFiles(old)
	id := nextSnapshotID(old.ID)
	reuseFrom := old
	if clean || id == initialSnapshotID {
		reuseFrom = nil
	}
	newSnapshot := newBuildSnapshot(id, reuseFrom, root, openFiles)
	if !clean {
		invalidateChangedDependents(old, newSnapshot, source.URI)
	}
	manager.Publish(newSnapshot)
	logLS(root, "snapshot refresh published key=%s kind=%s newID=%d modules=%d files=%d clean=%t", root, projectKindString(newSnapshot.Kind), newSnapshot.ID, len(newSnapshot.Modules), len(newSnapshot.Files), clean)
	s.publishDiagnostics(manager, newSnapshot, source)
	s.publishRemovedFileDiagnostics(old, newSnapshot)
}

func (s *Server) publishRemovedFileDiagnostics(old *Snapshot, snapshot *Snapshot) {
	for uri, file := range old.Files {
		if _, ok := snapshot.Files[uri]; ok {
			continue
		}
		version := file.Version
		s.writeNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         uri,
			Version:     &version,
			Diagnostics: []protocol.Diagnostic{},
		})
	}
}

func (s *Server) projectRootForWatchedPath(path string) string {
	path = normalizePath(path)
	if s.root == "" || !isUnder(path, s.root) {
		return ""
	}
	best := ""
	for root, manager := range s.snapshots {
		if manager == nil || manager.Current().Kind != ProjectKindBuild || !isUnder(path, root) {
			continue
		}
		if len(root) > len(best) {
			best = root
		}
	}
	if best != "" {
		return best
	}
	return s.projectRootForFile(path)
}

func isRelevantWatchedPath(root string, path string) bool {
	path = normalizePath(path)
	if !isUnder(path, root) {
		return false
	}
	if strings.HasSuffix(path, ".bal") || filepath.Base(path) == "Ballerina.toml" {
		return true
	}
	modulesDir := filepath.Join(root, "modules")
	return isUnder(path, modulesDir)
}

func invalidateChangedDependents(old *Snapshot, snapshot *Snapshot, changedURI protocol.DocumentURI) {
	if changedURI == "" {
		return
	}
	changedModule := ""
	for name, module := range snapshot.Modules {
		if _, ok := module.Files[changedURI]; ok {
			changedModule = name
			break
		}
	}
	if changedModule == "" {
		return
	}
	for _, name := range dependentModuleClosure(old, changedModule) {
		resetModuleState(snapshot.Modules[name])
	}
	snapshot.TopoOrder = nil
}

func dependentModuleClosure(snapshot *Snapshot, changedModule string) []string {
	if snapshot == nil {
		return []string{changedModule}
	}
	dependents := make(map[string][]string, len(snapshot.Modules))
	for name, module := range snapshot.Modules {
		for _, imp := range module.Imports {
			if imp.ModuleName == name {
				continue
			}
			dependents[imp.ModuleName] = append(dependents[imp.ModuleName], name)
		}
	}
	seen := map[string]bool{changedModule: true}
	queue := []string{changedModule}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		for _, dependent := range dependents[name] {
			if seen[dependent] {
				continue
			}
			seen[dependent] = true
			queue = append(queue, dependent)
		}
	}
	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func resetModuleState(module *Module) {
	if module == nil {
		return
	}
	module.Stage = FrontendStageNone
	module.Imports = nil
	module.ImportedByCU = nil
	module.ImportedSymbols = nil
	module.Package = nil
	module.Exported = model.ExportedSymbolSpace{}
	module.CFG = nil
}

func (s *Server) publishDiagnostics(manager *SnapshotManager, snapshot *Snapshot, source SourceFile) {
	s.beginIndexingProgress()
	defer s.endIndexingProgress()
	logLS(snapshot.Root, "diagnostics start snapshotID=%d kind=%s source=%s", snapshot.ID, projectKindString(snapshot.Kind), source.Path)
	diagnosticsByURI := runDiagnostics(snapshot, source)
	if !manager.IsCurrent(snapshot) {
		return
	}
	logLS(snapshot.Root, "diagnostics complete snapshotID=%d diagnosticFiles=%d", snapshot.ID, len(diagnosticsByURI))
	for uri, file := range snapshot.Files {
		diagnostics := diagnosticsByURI[uri]
		if diagnostics == nil {
			diagnostics = []protocol.Diagnostic{}
		}
		version := file.Version
		s.writeNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         uri,
			Version:     &version,
			Diagnostics: diagnostics,
		})
	}
}

func (s *Server) initializeSnapshots(params protocol.InitializeParams) {
	root := pathFromRoot(params)
	if root == "" {
		return
	}
	s.root = root
	s.workDoneProgress = params.Capabilities.Window != nil && params.Capabilities.Window.WorkDoneProgress
	s.watchedFilesDynamicRegistration = params.Capabilities.Workspace != nil &&
		params.Capabilities.Workspace.DidChangeWatchedFiles != nil &&
		params.Capabilities.Workspace.DidChangeWatchedFiles.DynamicRegistration
	logLS(root, "initialize root=%s build=%t", root, isBuildProjectRoot(root))
	if isBuildProjectRoot(root) {
		s.snapshots[root] = NewBuildSnapshotManager(root)
	}
}

func pathFromRoot(params protocol.InitializeParams) string {
	if params.RootURI != "" {
		return pathFromURI(protocol.DocumentURI(params.RootURI))
	}
	if params.RootPath != "" {
		return normalizePath(params.RootPath)
	}
	return ""
}

func (s *Server) snapshotKey(file SourceFile) string {
	if root := s.projectRootForFile(file.Path); root != "" {
		return root
	}
	return file.Path
}

func (s *Server) snapshotManager(key string, file SourceFile) *SnapshotManager {
	manager := s.snapshots[key]
	if manager != nil {
		return manager
	}
	if isBuildProjectRoot(key) {
		manager = NewBuildSnapshotManager(key)
	} else {
		manager = NewSingleFileSnapshotManager(file)
	}
	s.snapshots[key] = manager
	return manager
}

func (s *Server) sourceFile(uri protocol.DocumentURI) SourceFile {
	path := pathFromURI(uri)
	key := s.snapshotKey(SourceFile{URI: uri, Path: path, File: path})
	manager := s.snapshots[key]
	if manager == nil {
		return SourceFile{}
	}
	file := manager.Current().Files[uri]
	return file
}

func (s *Server) projectRootForFile(path string) string {
	path = normalizePath(path)
	if s.root != "" && !isUnder(path, s.root) {
		return ""
	}
	boundary := normalizePath(s.root)
	dir := path
	if filepath.Ext(path) != "" {
		dir = filepath.Dir(path)
	}
	for {
		if isBuildProjectRoot(dir) {
			return dir
		}
		if dir == boundary {
			return ""
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func isUnder(path, root string) bool {
	path = normalizePath(path)
	root = normalizePath(root)
	return path == root || strings.HasPrefix(path, root+string(filepath.Separator))
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
	logLS(s.root, "response sent id=%s", string(id))
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
	logLS(s.root, "response error id=%s code=%d message=%s", string(id), code, message)
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
	logLS(s.root, "notification sent method=%s", method)
	s.writeMessage(struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
	}{JSONRPC: "2.0", Method: method, Params: mustMarshal(params)})
}

func (s *Server) writeRequest(method string, params any) {
	s.nextServerRequestID++
	id := s.nextServerRequestID
	logLS(s.root, "request sent method=%s id=%d", method, id)
	s.writeMessage(struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      int64           `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
	}{JSONRPC: "2.0", ID: id, Method: method, Params: mustMarshal(params)})
}

func (s *Server) registerWatchedFiles() {
	if !s.watchedFilesDynamicRegistration {
		return
	}
	s.writeRequest("client/registerCapability", protocol.RegistrationParams{
		Registrations: []protocol.Registration{{
			ID:     "bal-workspace-file-watchers",
			Method: "workspace/didChangeWatchedFiles",
			RegisterOptions: protocol.DidChangeWatchedFilesRegistrationOptions{
				Watchers: []protocol.FileSystemWatcher{
					{GlobPattern: "**/*.bal", Kind: protocol.WatchKindCreate | protocol.WatchKindChange | protocol.WatchKindDelete},
					{GlobPattern: "**/Ballerina.toml", Kind: protocol.WatchKindCreate | protocol.WatchKindChange | protocol.WatchKindDelete},
					{GlobPattern: "**/modules/**", Kind: protocol.WatchKindCreate | protocol.WatchKindDelete},
				},
			},
		}},
	})
}

func (s *Server) beginIndexingProgress() {
	if !s.workDoneProgress {
		return
	}
	if !s.progressCreated {
		s.writeRequest("window/workDoneProgress/create", protocol.WorkDoneProgressCreateParams{Token: indexingProgressToken})
		s.progressCreated = true
	}
	s.writeNotification("$/progress", protocol.ProgressParams{
		Token: indexingProgressToken,
		Value: protocol.WorkDoneProgressBegin{Kind: "begin", Title: "indexing"},
	})
}

func (s *Server) endIndexingProgress() {
	if !s.workDoneProgress {
		return
	}
	s.writeNotification("$/progress", protocol.ProgressParams{
		Token: indexingProgressToken,
		Value: protocol.WorkDoneProgressEnd{Kind: "end"},
	})
}

func (s *Server) writeMessage(msg any) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	_, _ = fmt.Fprintf(&buffer, "Content-Length: %d\r\n\r\n", len(payload))
	buffer.Write(payload)
	_, _ = s.out.Write(buffer.Bytes())
}

func mustMarshal(value any) json.RawMessage {
	payload, _ := json.Marshal(value)
	return payload
}
