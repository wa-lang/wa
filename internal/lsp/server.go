// Package server defines gopls' implementation of the LSP server
// interface, [protocol.Server]. Call [New] to create an instance.
package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

// New creates an LSP server and binds it to handle incoming client
// messages on the supplied stream.
func New(client protocol.ClientCloser) protocol.Server {
	return &server{}
}

// server implements the protocol.server interface.
type server struct {
}

func (s *server) Progress(context.Context, *protocol.ProgressParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) SetTrace(context.Context, *protocol.SetTraceParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) IncomingCalls(context.Context, *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) OutgoingCalls(context.Context, *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ResolveCodeAction(context.Context, *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ResolveCodeLens(context.Context, *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ResolveCompletionItem(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ResolveDocumentLink(context.Context, *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Exit(context.Context) error { return fmt.Errorf("TODO") } // exit

func (s *server) Initialize(context.Context, *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	return nil, fmt.Errorf("TODO: Initialize")
}

func (s *server) Initialized(context.Context, *protocol.InitializedParams) error {
	return fmt.Errorf("TODO: Initialized")
}

func (s *server) Resolve(context.Context, *protocol.InlayHint) (*protocol.InlayHint, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DidChangeNotebookDocument(context.Context, *protocol.DidChangeNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidCloseNotebookDocument(context.Context, *protocol.DidCloseNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidOpenNotebookDocument(context.Context, *protocol.DidOpenNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidSaveNotebookDocument(context.Context, *protocol.DidSaveNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) Shutdown(context.Context) error {
	return fmt.Errorf("TODO")
}

func (s *server) CodeAction(context.Context, *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) CodeLens(context.Context, *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ColorPresentation(context.Context, *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Completion(context.Context, *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Declaration(context.Context, *protocol.DeclarationParams) (*protocol.Or_textDocument_declaration, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Definition(context.Context, *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Diagnostic(context.Context, *string) (*string, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DidChange(context.Context, *protocol.DidChangeTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidClose(context.Context, *protocol.DidCloseTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidOpen(context.Context, *protocol.DidOpenTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidSave(context.Context, *protocol.DidSaveTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DocumentColor(context.Context, *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DocumentHighlight(context.Context, *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DocumentLink(context.Context, *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DocumentSymbol(context.Context, *protocol.DocumentSymbolParams) ([]interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) FoldingRange(context.Context, *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Hover(context.Context, *protocol.HoverParams) (*protocol.Hover, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Implementation(context.Context, *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) InlayHint(context.Context, *protocol.InlayHintParams) ([]protocol.InlayHint, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) InlineCompletion(context.Context, *protocol.InlineCompletionParams) (*protocol.Or_Result_textDocument_inlineCompletion, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) InlineValue(context.Context, *protocol.InlineValueParams) ([]protocol.InlineValue, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) LinkedEditingRange(context.Context, *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Moniker(context.Context, *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) OnTypeFormatting(context.Context, *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) PrepareCallHierarchy(context.Context, *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) PrepareRename(context.Context, *protocol.PrepareRenameParams) (*protocol.PrepareRenameResult, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) PrepareTypeHierarchy(context.Context, *protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) RangeFormatting(context.Context, *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) RangesFormatting(context.Context, *protocol.DocumentRangesFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) References(context.Context, *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Rename(context.Context, *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) SelectionRange(context.Context, *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) SemanticTokensFull(context.Context, *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) SemanticTokensFullDelta(context.Context, *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) SemanticTokensRange(context.Context, *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) SignatureHelp(context.Context, *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) TypeDefinition(context.Context, *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) WillSave(context.Context, *protocol.WillSaveTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) WillSaveWaitUntil(context.Context, *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Subtypes(context.Context, *protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Supertypes(context.Context, *protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) WorkDoneProgressCancel(context.Context, *protocol.WorkDoneProgressCancelParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DiagnosticWorkspace(context.Context, *protocol.WorkspaceDiagnosticParams) (*protocol.WorkspaceDiagnosticReport, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) DidChangeConfiguration(context.Context, *protocol.DidChangeConfigurationParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidChangeWatchedFiles(context.Context, *protocol.DidChangeWatchedFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidChangeWorkspaceFolders(context.Context, *protocol.DidChangeWorkspaceFoldersParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidCreateFiles(context.Context, *protocol.CreateFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidDeleteFiles(context.Context, *protocol.DeleteFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) DidRenameFiles(context.Context, *protocol.RenameFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *server) ExecuteCommand(context.Context, *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) Symbol(context.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) WillCreateFiles(context.Context, *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) WillDeleteFiles(context.Context, *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) WillRenameFiles(context.Context, *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *server) ResolveWorkspaceSymbol(context.Context, *protocol.WorkspaceSymbol) (*protocol.WorkspaceSymbol, error) {
	return nil, fmt.Errorf("TODO")
}
