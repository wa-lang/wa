// 版权 @2023 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"log"
	"os"
	"sync"

	"wa-lang.org/wa/internal/lsp/fakenet"
	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	"wa-lang.org/wa/internal/lsp/protocol"
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

func NewLSPServer(opt *Option) *LSPServer {
	return &LSPServer{}
}

func (p *LSPServer) Run() error {
	stream := jsonrpc2.NewHeaderStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	conn := jsonrpc2.NewConn(stream)
	client := protocol.ClientDispatcher(conn)
	svr := New(client)

	ctx := protocol.WithClient(context.Background(), client)
	conn.Go(ctx, protocol.Handlers(
		protocol.ServerHandler(svr, jsonrpc2.MethodNotFound),
	))
	<-conn.Done()
	return conn.Err()
}
