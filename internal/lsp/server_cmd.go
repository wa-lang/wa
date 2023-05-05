// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"encoding/json"
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/jsonrpc2"
)

func (h *LSPServer) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("TODO")
}

func (h *LSPServer) handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	if len(params.ContentChanges) == 1 {
		// todo
	}
	return nil, fmt.Errorf("TODO")
}

func (h *LSPServer) handleTextDocumentDidSave(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidSaveTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	if params.Text != nil {
		// update
	} else {
		// save
	}

	return nil, fmt.Errorf("TODO")
}

func (h *LSPServer) handleTextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("TODO")
}

func (h *LSPServer) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}
	h.conn = conn

	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// Only try to parse the workspace root if its not null. Otherwise initialize will fail
	if params.RootURI != "" {
		//
	}

	var completion = &CompletionProvider{
		TriggerCharacters: []string{"*"},
	}

	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync:           TDSKFull,
			DocumentFormattingProvider: true,
			DocumentSymbolProvider:     true,
			DefinitionProvider:         true,
			ReferencesProvider:         true,
			CompletionProvider:         completion,
			HoverProvider:              true,
			CodeActionProvider:         true,
			Workspace: &ServerCapabilitiesWorkspace{
				WorkspaceFolders: WorkspaceFoldersServerCapabilities{
					Supported:           true,
					ChangeNotifications: true,
				},
			},
		},
	}, nil
}

func (h *LSPServer) logMessage(typ MessageType, message string) {
	h.conn.Notify(context.Background(),
		k_window_logMessage, &LogMessageParams{
			Type: typ, Message: message,
		},
	)
}

func (p *LSPServer) executeCommand(params *ExecuteCommandParams) (interface{}, error) {
	if len(params.Arguments) != 1 {
		return nil, fmt.Errorf("invalid command")
	}
	return nil, nil
}

func (p *LSPServer) codeAction(uri string, params *CodeActionParams) ([]Command, error) {
	return nil, nil
}

func (p *LSPServer) completion(uri string, params *CompletionParams) ([]CompletionItem, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) definition(uri string, params *TextDocumentPositionParams) ([]Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) formatting(uri string, options FormattingOptions) ([]TextEdit, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) hover(uri string, params *HoverParams) (*Hover, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *LSPServer) findRefs(uri string, params *TextDocumentPositionParams) ([]Location, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) symbol(uri string) ([]DocumentSymbol, error) {
	return nil, fmt.Errorf("TODO")
}
