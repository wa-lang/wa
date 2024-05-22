// 版权 @2024 凹语言 作者。保留所有权利。

// 针对 LSP 定制的 loader

package loaderx

import (
	"context"
	"io"
	"io/fs"
	"log"

	"wa-lang.org/wa/internal/lsp/protocol"
	"wa-lang.org/wa/internal/token"
)

func NewUniverse(cfg Config) *Universe {
	p := &Universe{
		Version: 1,

		WaOS:            cfg.WaOS,
		WaRoot:          cfg.WaRoot,
		WorkspaceFolder: []URI{},

		Fset:  token.NewFileSet(),
		Files: make(map[protocol.DocumentURI]*File),
		Pkgs:  make(map[PkgpathURI]*Package),

		logger: cfg.logger,
	}

	if p.logger == nil {
		p.logger = log.New(io.Discard, "", 0)
	}

	return p
}

func (p *Universe) VFS() fs.FS {
	return &openedFS{p}
}

// 配置发生变化
func (p *Universe) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	p.logger.Println("Universe.DidChangeConfiguration:", params.Settings)
	return nil
}

// 监视的文件列表发生变化
func (p *Universe) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	for _, x := range params.Changes {
		switch x.Type {
		case protocol.Created:
			p.logger.Println("Universe.DidChangeWatchedFiles:", x.URI, "Created")
		case protocol.Changed:
			p.logger.Println("Universe.DidChangeWatchedFiles:", x.URI, "Changed")
		case protocol.Deleted:
			p.logger.Println("Universe.DidChangeWatchedFiles:", x.URI, "Deleted")
		}
	}

	return nil
}

// 工作区发生变化
func (p *Universe) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	for _, x := range params.Event.Removed {
		p.removeWorkspaceFolders(x.URI)
	}
	for _, x := range params.Event.Added {
		p.addWorkspaceFolders(x.URI)
	}
	return nil
}

// 创建文件
func (p *Universe) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) error {
	for _, f := range params.Files {
		p.logger.Println("Universe.DidCreateFiles:", f.URI)
	}
	return nil
}

// 删除文件
func (p *Universe) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) error {
	for _, f := range params.Files {
		p.logger.Println("Universe.DidDeleteFiles:", f.URI)
		delete(p.Files, protocol.DocumentURI(f.URI))
	}
	return nil
}

// 重新命名文件
func (p *Universe) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) error {
	for _, f := range params.Files {
		p.logger.Println("Universe.DidRenameFiles:", f.OldURI, "=>", f.NewURI)
		delete(p.Files, protocol.DocumentURI(f.OldURI))
		delete(p.Files, protocol.DocumentURI(f.NewURI))
	}
	return nil
}

func (p *Universe) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	p.logger.Println("Universe.DidOpen:", params.TextDocument.URI)
	p.Files[params.TextDocument.URI] = &File{
		Version: 0,
		FileUri: params.TextDocument.URI,
		PkgPath: "",
		Data:    []byte(params.TextDocument.Text),
	}
	return nil
}

func (p *Universe) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	p.logger.Println("Universe.DidSave:", params.TextDocument.URI)
	delete(p.Files, protocol.DocumentURI(params.TextDocument.URI))
	return nil
}

func (p *Universe) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	p.logger.Println("Universe.DidClose:", params.TextDocument.URI)
	delete(p.Files, protocol.DocumentURI(params.TextDocument.URI))
	return nil
}

func (p *Universe) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	text, err := p.changedText(params.TextDocument.URI, params.ContentChanges)
	if err != nil {
		return err
	}
	p.Files[params.TextDocument.URI].Data = text
	return nil
}
