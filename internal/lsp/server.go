// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
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

	clientName       string
	workspaceFolders []string
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

	// 记录客户端名字
	if params != nil && params.ClientInfo != nil {
		s.clientName = params.ClientInfo.Name
	}

	// 记录工作区目录列表
	for _, folder := range params.WorkspaceFolders {
		if folder.URI == "" {
			return nil, fmt.Errorf("empty WorkspaceFolder.URI")
		}
		if _, err := protocol.ParseDocumentURI(folder.URI); err != nil {
			return nil, fmt.Errorf("invalid WorkspaceFolder.URI: %v", err)
		}
		s.workspaceFolders = append(s.workspaceFolders, folder.URI)
	}

	// 如果没有工作区, 且打开的是目录, 则用改目录
	if len(s.workspaceFolders) == 0 && params.RootURI != "" {
		s.workspaceFolders = append(s.workspaceFolders, params.RootURI.Path())
	}

	// TODO(chai): 工作区需要对应到 Wa 模块根目录

	reply := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DocumentFormattingProvider: &protocol.Or_ServerCapabilities_documentFormattingProvider{Value: true},

			// TODO(chai): 实现基本能力
			HoverProvider:          &protocol.Or_ServerCapabilities_hoverProvider{Value: false},
			DefinitionProvider:     &protocol.Or_ServerCapabilities_definitionProvider{Value: false},
			TypeDefinitionProvider: &protocol.Or_ServerCapabilities_typeDefinitionProvider{Value: false},
			FoldingRangeProvider:   &protocol.Or_ServerCapabilities_foldingRangeProvider{Value: false},
			CodeActionProvider:     false,
			RenameProvider:         false,
			ExecuteCommandProvider: nil,
			CodeLensProvider:       nil,

			// TODO: 支持工作区列表变化
			Workspace: &protocol.WorkspaceOptions{
				WorkspaceFolders: &protocol.WorkspaceFolders5Gn{
					Supported:           true,
					ChangeNotifications: "workspace/didChangeWorkspaceFolders",
				},
			},
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
	s.logger.Println("Exit")
	return nil
}

func (s *LSPServer) DidChangeConfiguration(context.Context, *protocol.DidChangeConfigurationParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidChangeWorkspaceFolders(context.Context, *protocol.DidChangeWorkspaceFoldersParams) error {
	return fmt.Errorf("TODO")
}
