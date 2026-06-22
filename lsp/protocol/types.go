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

package protocol

import "encoding/json"

type DocumentURI string

type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeParams struct {
	RootURI      string             `json:"rootUri,omitempty"`
	RootPath     string             `json:"rootPath,omitempty"`
	Capabilities ClientCapabilities `json:"capabilities,omitempty"`
}

type ClientCapabilities struct {
	Window    *WindowClientCapabilities    `json:"window,omitempty"`
	Workspace *WorkspaceClientCapabilities `json:"workspace,omitempty"`
}

type WorkspaceClientCapabilities struct {
	DidChangeWatchedFiles *DidChangeWatchedFilesClientCapabilities `json:"didChangeWatchedFiles,omitempty"`
}

type DidChangeWatchedFilesClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type WindowClientCapabilities struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync        TextDocumentSyncOptions `json:"textDocumentSync"`
	CompletionProvider      *CompletionOptions      `json:"completionProvider,omitempty"`
	DefinitionProvider      bool                    `json:"definitionProvider,omitempty"`
	ReferenceProvider       bool                    `json:"referencesProvider,omitempty"`
	CodeActionProvider      bool                    `json:"codeActionProvider,omitempty"`
	DocumentSymbolProvider  bool                    `json:"documentSymbolProvider,omitempty"`
	WorkspaceSymbolProvider bool                    `json:"workspaceSymbolProvider,omitempty"`
}

type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool        `json:"openClose"`
	Change    int         `json:"change"`
	Save      SaveOptions `json:"save"`
}

type SaveOptions struct {
	IncludeText bool `json:"includeText"`
}

type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int32       `json:"version"`
	Text       string      `json:"text"`
}

type VersionedTextDocumentIdentifier struct {
	URI     DocumentURI `json:"uri"`
	Version int32       `json:"version"`
}

type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         *string                `json:"text,omitempty"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Kind           SymbolKind       `json:"kind"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type WorkspaceSymbolParams struct {
	Query string `json:"query"`
}

type SymbolInformation struct {
	Name     string     `json:"name"`
	Kind     SymbolKind `json:"kind"`
	Location Location   `json:"location"`
}

type SymbolKind int

const (
	SymbolKindFile SymbolKind = iota + 1
	SymbolKindModule
	SymbolKindNamespace
	SymbolKindPackage
	SymbolKindClass
	SymbolKindMethod
	SymbolKindProperty
	SymbolKindField
	SymbolKindConstructor
	SymbolKindEnum
	SymbolKindInterface
	SymbolKindFunction
	SymbolKindVariable
	SymbolKindConstant
	SymbolKindString
	SymbolKindNumber
	SymbolKindBoolean
	SymbolKindArray
	SymbolKindObject
	SymbolKindKey
	SymbolKindNull
	SymbolKindEnumMember
	SymbolKindStruct
	SymbolKindEvent
	SymbolKindOperator
	SymbolKindTypeParameter
)

type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

type FileEvent struct {
	URI  DocumentURI    `json:"uri"`
	Type FileChangeType `json:"type"`
}

type FileChangeType int

const (
	FileChangeTypeCreated FileChangeType = 1
	FileChangeTypeChanged FileChangeType = 2
	FileChangeTypeDeleted FileChangeType = 3
)

type RenameFilesParams struct {
	Files []FileRename `json:"files"`
}

type FileRename struct {
	OldURI DocumentURI `json:"oldUri"`
	NewURI DocumentURI `json:"newUri"`
}

type CreateFilesParams struct {
	Files []FileCreate `json:"files"`
}

type FileCreate struct {
	URI DocumentURI `json:"uri"`
}

type DeleteFilesParams struct {
	Files []FileDelete `json:"files"`
}

type FileDelete struct {
	URI DocumentURI `json:"uri"`
}

type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type DefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type ReferenceParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	Context      ReferenceContext       `json:"context"`
}

type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type CodeActionContext struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type CodeAction struct {
	Title       string        `json:"title"`
	Kind        string        `json:"kind,omitempty"`
	Diagnostics []Diagnostic  `json:"diagnostics,omitempty"`
	Edit        WorkspaceEdit `json:"edit,omitempty"`
}

type WorkspaceEdit struct {
	Changes map[DocumentURI][]TextEdit `json:"changes,omitempty"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label               string     `json:"label"`
	Kind                int        `json:"kind,omitempty"`
	Detail              string     `json:"detail,omitempty"`
	InsertText          string     `json:"insertText,omitempty"`
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
}

const (
	CompletionItemKindFunction = 3
	CompletionItemKindVariable = 6
	CompletionItemKindClass    = 7
	CompletionItemKindModule   = 9
	CompletionItemKindConstant = 21
)

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity,omitempty"`
	Code     string `json:"code,omitempty"`
	Source   string `json:"source,omitempty"`
	Message  string `json:"message"`
}

type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Version     *int32       `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type ProgressToken string

type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

type Registration struct {
	ID              string `json:"id"`
	Method          string `json:"method"`
	RegisterOptions any    `json:"registerOptions,omitempty"`
}

type DidChangeWatchedFilesRegistrationOptions struct {
	Watchers []FileSystemWatcher `json:"watchers"`
}

type FileSystemWatcher struct {
	GlobPattern string    `json:"globPattern"`
	Kind        WatchKind `json:"kind,omitempty"`
}

type WatchKind int

const (
	WatchKindCreate WatchKind = 1
	WatchKindChange WatchKind = 2
	WatchKindDelete WatchKind = 4
)

type WorkDoneProgressCreateParams struct {
	Token ProgressToken `json:"token"`
}

type ProgressParams struct {
	Token ProgressToken `json:"token"`
	Value any           `json:"value"`
}

type WorkDoneProgressBegin struct {
	Kind  string `json:"kind"`
	Title string `json:"title"`
}

type WorkDoneProgressEnd struct {
	Kind string `json:"kind"`
}
