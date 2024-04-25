// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"wa-lang.org/wa/internal/3rdparty/jsonrpc2"
)

type Option struct {
	LogFile string
}

type LSPServer struct {
	mu sync.Mutex

	conn     *jsonrpc2.Conn
	rootPath string
	shutdown bool

	logger   *log.Logger
	loglevel int
}

func (p *Option) clone() *Option {
	if p == nil {
		return &Option{}
	}
	return &Option{
		LogFile: p.LogFile,
	}
}

func NewLSPServer(opt *Option) *LSPServer {
	return &LSPServer{}
}

func (p *LSPServer) Run() {
	var rwc = struct {
		io.Writer
		io.ReadCloser
	}{
		ReadCloser: io.NopCloser(os.Stdin),
		Writer:     os.Stdout,
	}

	conn := jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(rwc, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(p.handle),
	)
	ch := conn.DisconnectNotify()
	<-ch
}

func (h *LSPServer) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Notif {
		switch req.Method {
		case k_initialized:
			return
		case k_textDocument_didOpen:
			return h.handleTextDocumentDidOpen(ctx, conn, req)
		case k_textDocument_didChange:
			return h.handleTextDocumentDidChange(ctx, conn, req)
		case k_textDocument_didSave:
			return h.handleTextDocumentDidSave(ctx, conn, req)
		case k_textDocument_didClose:
			return h.handleTextDocumentDidClose(ctx, conn, req)
		case k_exit:
			return nil, nil
		}
		return
	}

	switch req.Method {
	case k_initialize:
		return nil, fmt.Errorf("TODO")
	case k_initialized:
		return
	case k_shutdown:
		return nil, fmt.Errorf("TODO")

	case k_workspace_executeCommand:
		return nil, fmt.Errorf("TODO")
	case k_workspace_didChangeConfiguration:
		return nil, fmt.Errorf("TODO")
	case k_workspace_workspaceFolders:
		return nil, fmt.Errorf("TODO")
	case k_workspace_didChangeWorkspaceFolders:
		return nil, fmt.Errorf("TODO")

	case k_textDocument_didOpen:
		//var params DidOpenTextDocumentParams
		//if err := json.Unmarshal(argsData, &params); err != nil {
		//	return nil, err
		//}
		return nil, fmt.Errorf("TODO")
	case k_textDocument_didChange:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_didSave:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_didClose:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_formatting:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_documentSymbol:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_completion:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_definition:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_references:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_hover:
		return nil, fmt.Errorf("TODO")
	case k_textDocument_codeAction:
		return nil, fmt.Errorf("TODO")
	}
	return nil, nil
}
