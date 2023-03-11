// 版权 @2023 凹语言 作者。保留所有权利。

// https://microsoft.github.io/language-server-protocol/

package lsp

type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	RootURI               string             `json:"rootUri,omitempty"`
	InitializationOptions InitializeOptions  `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
}

type InitializeOptions struct {
	DocumentFormatting bool `json:"documentFormatting"`
	Hover              bool `json:"hover"`
	DocumentSymbol     bool `json:"documentSymbol"`
	CodeAction         bool `json:"codeAction"`
	Completion         bool `json:"completion"`
}

type ClientCapabilities struct {
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

type MessageType int

const (
	LogError   MessageType = 1
	LogWarning MessageType = 2
	LogInfo    MessageType = 3
	LogLog     MessageType = 4
)

type TextDocumentSyncKind int

const (
	TDSKNone        TextDocumentSyncKind = 0
	TDSKFull        TextDocumentSyncKind = 1
	TDSKIncremental TextDocumentSyncKind = 2
)

type CompletionProvider struct {
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
	TriggerCharacters []string `json:"triggerCharacters"`
}

type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported"`
	ChangeNotifications bool `json:"changeNotifications"`
}

type ServerCapabilitiesWorkspace struct {
	WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders"`
}

type ServerCapabilities struct {
	TextDocumentSync           TextDocumentSyncKind         `json:"textDocumentSync,omitempty"`
	DocumentSymbolProvider     bool                         `json:"documentSymbolProvider,omitempty"`
	CompletionProvider         *CompletionProvider          `json:"completionProvider,omitempty"`
	DefinitionProvider         bool                         `json:"definitionProvider,omitempty"`
	ReferencesProvider         bool                         `json:"referencesProvider,omitempty"`
	DocumentFormattingProvider bool                         `json:"documentFormattingProvider,omitempty"`
	HoverProvider              bool                         `json:"hoverProvider,omitempty"`
	CodeActionProvider         bool                         `json:"codeActionProvider,omitempty"`
	Workspace                  *ServerCapabilitiesWorkspace `json:"workspace,omitempty"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range"`
	RangeLength int    `json:"rangeLength"`
	Text        string `json:"text"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type DidSaveTextDocumentParams struct {
	Text         *string                `json:"text"`
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	CompletionContext CompletionContext `json:"contentChanges"`
}

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter *string               `json:"triggerCharacter"`
}

type CompletionTriggerKind int

const (
	Invoked                         CompletionTriggerKind = 1
	TriggerCharacter                CompletionTriggerKind = 2
	TriggerForIncompleteCompletions CompletionTriggerKind = 3
)

type HoverParams struct {
	TextDocumentPositionParams
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           int                            `json:"severity,omitempty"`
	Code               *string                        `json:"code,omitempty"`
	Source             *string                        `json:"source,omitempty"`
	Message            string                         `json:"message"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Version     int          `json:"version"`
}

type FormattingOptions map[string]interface{}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Kind           int64            `json:"kind"`
	Deprecated     bool             `json:"deprecated"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
	Detail         string           `json:"detail,omitempty"`
}

type SymbolInformation struct {
	Name          string   `json:"name"`
	Kind          int64    `json:"kind"`
	Deprecated    bool     `json:"deprecated"`
	Location      Location `json:"location"`
	ContainerName *string  `json:"containerName"`
}

type CompletionItemKind int

const (
	TextCompletion          CompletionItemKind = 1
	MethodCompletion        CompletionItemKind = 2
	FunctionCompletion      CompletionItemKind = 3
	ConstructorCompletion   CompletionItemKind = 4
	FieldCompletion         CompletionItemKind = 5
	VariableCompletion      CompletionItemKind = 6
	ClassCompletion         CompletionItemKind = 7
	InterfaceCompletion     CompletionItemKind = 8
	ModuleCompletion        CompletionItemKind = 9
	PropertyCompletion      CompletionItemKind = 10
	UnitCompletion          CompletionItemKind = 11
	ValueCompletion         CompletionItemKind = 12
	EnumCompletion          CompletionItemKind = 13
	KeywordCompletion       CompletionItemKind = 14
	SnippetCompletion       CompletionItemKind = 15
	ColorCompletion         CompletionItemKind = 16
	FileCompletion          CompletionItemKind = 17
	ReferenceCompletion     CompletionItemKind = 18
	FolderCompletion        CompletionItemKind = 19
	EnumMemberCompletion    CompletionItemKind = 20
	ConstantCompletion      CompletionItemKind = 21
	StructCompletion        CompletionItemKind = 22
	EventCompletion         CompletionItemKind = 23
	OperatorCompletion      CompletionItemKind = 24
	TypeParameterCompletion CompletionItemKind = 25
)

type CompletionItemTag int

type InsertTextFormat int

const (
	PlainTextTextFormat InsertTextFormat = 1
	SnippetTextFormat   InsertTextFormat = 2
)

type Command struct {
	Title     string        `json:"title" yaml:"title"`
	Command   string        `json:"command" yaml:"command"`
	Arguments []interface{} `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	OS        string        `json:"-" yaml:"os,omitempty"`
}

type WorkspaceEdit struct {
	Changes         interface{} `json:"changes"`
	DocumentChanges interface{} `json:"documentChanges"`
}

type CodeAction struct {
	Title       string         `json:"title"`
	Diagnostics []Diagnostic   `json:"diagnostics"`
	IsPreferred bool           `json:"isPreferred"`
	Edit        *WorkspaceEdit `json:"edit"`
	Command     *Command       `json:"command"`
}

type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                CompletionItemKind  `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              string              `json:"detail,omitempty"`
	Documentation       string              `json:"documentation,omitempty"` // string | MarkupContent
	Deprecated          bool                `json:"deprecated,omitempty"`
	Preselect           bool                `json:"preselect,omitempty"`
	SortText            string              `json:"sortText,omitempty"`
	FilterText          string              `json:"filterText,omitempty"`
	InsertText          string              `json:"insertText,omitempty"`
	InsertTextFormat    InsertTextFormat    `json:"insertTextFormat,omitempty"`
	TextEdit            *TextEdit           `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit          `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string            `json:"commitCharacters,omitempty"`
	Command             *Command            `json:"command,omitempty"`
	Data                interface{}         `json:"data,omitempty"`
}

type Hover struct {
	Contents interface{} `json:"contents"`
	Range    *Range      `json:"range"`
}

type MarkedString struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type WorkDoneProgressParams struct {
	WorkDoneToken interface{} `json:"workDoneToken"`
}

type ExecuteCommandParams struct {
	WorkDoneProgressParams

	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

type CodeActionKind string

const (
	Empty                 CodeActionKind = ""
	QuickFix              CodeActionKind = "quickfix"
	Refactor              CodeActionKind = "refactor"
	RefactorExtract       CodeActionKind = "refactor.extract"
	RefactorInline        CodeActionKind = "refactor.inline"
	RefactorRewrite       CodeActionKind = "refactor.rewrite"
	Source                CodeActionKind = "source"
	SourceOrganizeImports CodeActionKind = "source.organizeImports"
)

type CodeActionContext struct {
	Diagnostics []Diagnostic     `json:"diagnostics"`
	Only        []CodeActionKind `json:"only,omitempty"`
}

type PartialResultParams struct {
	PartialResultToken interface{} `json:"partialResultToken"`
}

type CodeActionParams struct {
	WorkDoneProgressParams
	PartialResultParams

	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type DidChangeConfigurationParams struct {
	Settings struct{} `json:"settings"`
}

type NotificationMessage struct {
	Method string      `json:"message"`
	Params interface{} `json:"params"`
}

type DocumentDefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

type ReferenceParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

type LogMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

type DidChangeWorkspaceFoldersParams struct {
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

type WorkspaceFoldersChangeEvent struct {
	Added   []WorkspaceFolder `json:"added,omitempty"`
	Removed []WorkspaceFolder `json:"removed,omitempty"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}
