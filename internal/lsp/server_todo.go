// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) Progress(context.Context, *protocol.ProgressParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) SetTrace(context.Context, *protocol.SetTraceParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) IncomingCalls(context.Context, *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) OutgoingCalls(context.Context, *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveCodeAction(context.Context, *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveCodeLens(context.Context, *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveCompletionItem(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveDocumentLink(context.Context, *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Resolve(context.Context, *protocol.InlayHint) (*protocol.InlayHint, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DidChangeNotebookDocument(context.Context, *protocol.DidChangeNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidCloseNotebookDocument(context.Context, *protocol.DidCloseNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidOpenNotebookDocument(context.Context, *protocol.DidOpenNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidSaveNotebookDocument(context.Context, *protocol.DidSaveNotebookDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) Shutdown(context.Context) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) CodeAction(context.Context, *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) CodeLens(context.Context, *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ColorPresentation(context.Context, *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Completion(context.Context, *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Declaration(context.Context, *protocol.DeclarationParams) (*protocol.Or_textDocument_declaration, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Definition(context.Context, *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Diagnostic(context.Context, *string) (*string, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DidChange(context.Context, *protocol.DidChangeTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidClose(context.Context, *protocol.DidCloseTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidOpen(context.Context, *protocol.DidOpenTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidSave(context.Context, *protocol.DidSaveTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DocumentColor(context.Context, *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DocumentHighlight(context.Context, *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DocumentLink(context.Context, *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DocumentSymbol(context.Context, *protocol.DocumentSymbolParams) ([]interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) FoldingRange(context.Context, *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Hover(context.Context, *protocol.HoverParams) (*protocol.Hover, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Implementation(context.Context, *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) InlayHint(context.Context, *protocol.InlayHintParams) ([]protocol.InlayHint, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) InlineCompletion(context.Context, *protocol.InlineCompletionParams) (*protocol.Or_Result_textDocument_inlineCompletion, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) InlineValue(context.Context, *protocol.InlineValueParams) ([]protocol.InlineValue, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) LinkedEditingRange(context.Context, *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Moniker(context.Context, *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) OnTypeFormatting(context.Context, *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) PrepareCallHierarchy(context.Context, *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) PrepareRename(context.Context, *protocol.PrepareRenameParams) (*protocol.PrepareRenameResult, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) PrepareTypeHierarchy(context.Context, *protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) RangeFormatting(context.Context, *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) RangesFormatting(context.Context, *protocol.DocumentRangesFormattingParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) References(context.Context, *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Rename(context.Context, *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) SelectionRange(context.Context, *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) SemanticTokensFull(context.Context, *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) SemanticTokensFullDelta(context.Context, *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) SemanticTokensRange(context.Context, *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) SignatureHelp(context.Context, *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) TypeDefinition(context.Context, *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) WillSave(context.Context, *protocol.WillSaveTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) WillSaveWaitUntil(context.Context, *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Subtypes(context.Context, *protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Supertypes(context.Context, *protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) WorkDoneProgressCancel(context.Context, *protocol.WorkDoneProgressCancelParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DiagnosticWorkspace(context.Context, *protocol.WorkspaceDiagnosticParams) (*protocol.WorkspaceDiagnosticReport, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) DidChangeConfiguration(context.Context, *protocol.DidChangeConfigurationParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidChangeWatchedFiles(context.Context, *protocol.DidChangeWatchedFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidChangeWorkspaceFolders(context.Context, *protocol.DidChangeWorkspaceFoldersParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidCreateFiles(context.Context, *protocol.CreateFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidDeleteFiles(context.Context, *protocol.DeleteFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidRenameFiles(context.Context, *protocol.RenameFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) ExecuteCommand(context.Context, *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) Symbol(context.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) WillCreateFiles(context.Context, *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) WillDeleteFiles(context.Context, *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) WillRenameFiles(context.Context, *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveWorkspaceSymbol(context.Context, *protocol.WorkspaceSymbol) (*protocol.WorkspaceSymbol, error) {
	return nil, fmt.Errorf("TODO")
}
