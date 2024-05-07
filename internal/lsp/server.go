// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"wa-lang.org/wa/internal/lsp/fakenet"
	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	_ "wa-lang.org/wa/internal/lsp/loaderx"
	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/version"
)

var _ protocol.Server = (*LSPServer)(nil)

type Option struct {
	LogFile     string
	SyncFileDir string // 文件快照目录
	WaOS        string // build 的目标系统
	WaRoot      string // 本地的 waroot 目录
}

type LSPServer struct {
	logger     *log.Logger
	conn       jsonrpc2.Conn
	client     protocol.ClientCloser
	clientName string

	waModules map[protocol.DocumentURI]*WaModule
	fileMap   map[string]string // 被编辑和修改的文件, 临时缓存

	syncFile *SyncFile
}

func NewLSPServer(opt *Option) *LSPServer {
	p := &LSPServer{
		waModules: make(map[protocol.DocumentURI]*WaModule),
		fileMap:   make(map[string]string),
	}

	logPrefix := fmt.Sprintf("[PID:%d] ", os.Getpid())

	if opt != nil && opt.LogFile != "" {
		f, err := os.OpenFile(opt.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		p.logger = log.New(f, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	}
	if p.logger == nil {
		p.logger = log.New(io.Discard, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	}

	if opt != nil && opt.SyncFileDir != "" {
		p.syncFile = &SyncFile{
			RootDir: opt.SyncFileDir,
		}
	}
	if p.syncFile == nil {
		p.syncFile = &SyncFile{}
	}

	stream := jsonrpc2.NewHeaderStream(fakenet.NewConn("stdio", os.Stdin, os.Stdout))
	p.conn = jsonrpc2.NewConn(stream)
	p.client = protocol.ClientDispatcher(p.conn)

	return p
}

func (p *LSPServer) Run() error {
	p.logger.Println("LSPServer.Version:", version.Version)

	info, _ := debug.ReadBuildInfo()
	p.logger.Println("LSPServer.BuildInfo:", jsonMarshal(info))

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

	reply := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				Change:    protocol.Incremental,
				OpenClose: true,
				Save: &protocol.SaveOptions{
					IncludeText: false,
				},
			},

			DocumentFormattingProvider: &protocol.Or_ServerCapabilities_documentFormattingProvider{Value: true},
			HoverProvider:              &protocol.Or_ServerCapabilities_hoverProvider{Value: true},
			DefinitionProvider:         &protocol.Or_ServerCapabilities_definitionProvider{Value: true},
			TypeDefinitionProvider:     &protocol.Or_ServerCapabilities_typeDefinitionProvider{Value: true},

			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{"."},
			},

			// TODO(chai): 实现基本能力
			FoldingRangeProvider:   &protocol.Or_ServerCapabilities_foldingRangeProvider{Value: false},
			CodeActionProvider:     false,
			RenameProvider:         false,
			ExecuteCommandProvider: nil,
			CodeLensProvider:       nil,
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

func (s *LSPServer) Shutdown(context.Context) error {
	return nil
}

func (s *LSPServer) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	s.logger.Println("DidChangeConfiguration:", jsonMarshal(params))
	return nil
}

func (s *LSPServer) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	s.logger.Println("DidChangeConfiguration:", jsonMarshal(params))
	return nil
}
