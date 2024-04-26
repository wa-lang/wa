// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"io"
	"log"
	"os"

	"wa-lang.org/wa/internal/lsp/fakenet"
	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/version"
)

type Option struct {
	LogFile string
}

type LSPServer struct {
	logger *log.Logger
	conn   jsonrpc2.Conn
	client protocol.ClientCloser
}

func NewLSPServer(opt *Option) *LSPServer {
	p := &LSPServer{}

	if opt != nil && opt.LogFile != "" {
		f, err := os.OpenFile(opt.LogFile, os.O_CREATE|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		p.logger = log.New(f, "", log.Ltime|log.Lshortfile)
	}
	if p.logger == nil {
		p.logger = log.New(io.Discard, "", log.Ltime|log.Lshortfile)
	}

	stream := jsonrpc2.NewHeaderStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	p.conn = jsonrpc2.NewConn(stream)
	p.client = protocol.ClientDispatcher(p.conn)

	return p
}

func (p *LSPServer) Run() error {
	p.logger.Println("LSPServer.Run")

	ctx := protocol.WithClient(context.Background(), p.client)
	p.conn.Go(ctx, protocol.Handlers(
		protocol.ServerHandler(p, jsonrpc2.MethodNotFound),
	))
	<-p.conn.Done()
	return p.conn.Err()
}

func (s *LSPServer) Initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	s.logger.Println("Initialize:", jsonMarshal(params))

	reply := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DocumentFormattingProvider: &protocol.Or_ServerCapabilities_documentFormattingProvider{Value: true},
		},
		ServerInfo: &protocol.ServerInfo{
			Version: version.Version,
			Name:    "wa lsp",
		},
	}
	return reply, nil
}

func (s *LSPServer) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.logger.Println("Initialized:", jsonMarshal(params))
	return nil
}

func (s *LSPServer) Exit(ctx context.Context) error {
	s.logger.Println("Exit:")
	return nil
}
