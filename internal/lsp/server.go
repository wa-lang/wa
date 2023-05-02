// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type LSPServer struct {
	conn *lspChannel
}

func NewLSPServer() *LSPServer {
	return &LSPServer{
		conn: &lspChannel{r: os.Stdin, w: os.Stdout},
	}
}

func (p *LSPServer) Run() error {
	ctx := context.Background()
	for {
		req, err := p.conn.RecvRequest()
		if err != nil {
			return err
		}

		var resp = &lspResponse{
			ID:      req.ID,
			JSONRPC: req.JSONRPC,
		}
		resp.Result, err = p.handleMethod(ctx, req)
		if err != nil {
			resp.Error = &lspError{
				Message: err.Error(),
			}
		}

		if err = p.conn.SendRespose(resp); err != nil {
			// log
		}
	}
}

func (p *LSPServer) handleMethod(ctx context.Context, req *lspRequest) (result interface{}, err error) {
	var argsData []byte
	if len(req.Params) > 0 {
		argsData = []byte(req.Params[0])
	}

	switch req.Method {
	case "initialized":
		return
	case "textDocument/didOpen":
		var params DidOpenTextDocumentParams
		if err := json.Unmarshal(argsData, &params); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("TODO")
	case "textDocument/didChange":
		return nil, fmt.Errorf("TODO")
	case "textDocument/didSave":
		return nil, fmt.Errorf("TODO")
	case "textDocument/didClose":
		return nil, fmt.Errorf("TODO")
	case "exit":
		return nil, fmt.Errorf("TODO")

	case "initialize":
		return nil, fmt.Errorf("TODO")
	case "shutdown":
		return nil, fmt.Errorf("TODO")
	case "textDocument/formatting":
		return nil, fmt.Errorf("TODO")
	case "textDocument/documentSymbol":
		return nil, fmt.Errorf("TODO")
	case "textDocument/completion":
		return nil, fmt.Errorf("TODO")
	case "textDocument/definition":
		return nil, fmt.Errorf("TODO")
	case "textDocument/references":
		return nil, fmt.Errorf("TODO")
	case "textDocument/hover":
		return nil, fmt.Errorf("TODO")
	case "textDocument/codeAction":
		return nil, fmt.Errorf("TODO")
	case "workspace/executeCommand":
		return nil, fmt.Errorf("TODO")
	case "workspace/didChangeConfiguration":
		return nil, fmt.Errorf("TODO")
	case "workspace/workspaceFolders":
		return nil, fmt.Errorf("TODO")
	case "workspace/didChangeWorkspaceFolders":
		return nil, fmt.Errorf("TODO")
	}
	return nil, nil
}
