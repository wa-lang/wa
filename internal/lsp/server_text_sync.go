// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) DidChangeWatchedFiles(context.Context, *protocol.DidChangeWatchedFilesParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidOpen(context.Context, *protocol.DidOpenTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidSave(context.Context, *protocol.DidSaveTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (s *LSPServer) DidClose(context.Context, *protocol.DidCloseTextDocumentParams) error {
	return fmt.Errorf("TODO")
}

func (p *LSPServer) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	p.logger.Println("DidChange:", jsonMarshal(params))

	if !strings.HasSuffix(string(params.TextDocument.URI), ".wa") {
		return nil
	}

	path := params.TextDocument.URI.Path()

	// todo: 目前只支持全局同步
	for _, x := range params.ContentChanges {
		if x.Range == nil {
			p.fileMap[path] = x.Text
			break
		}
	}

	return nil
}
func (s *LSPServer) changedText(uri protocol.DocumentURI, changes []protocol.TextDocumentContentChangeEvent) ([]byte, error) {
	if len(changes) == 0 {
		return nil, fmt.Errorf("%w: no content changes provided", jsonrpc2.ErrInternal)
	}

	// 全量变更数据
	if len(changes) == 1 && changes[0].Range == nil && changes[0].RangeLength == 0 {
		return []byte(changes[0].Text), nil
	}

	// 增量变更
	return s.applyIncrementalChanges(uri, changes)
}

// 接受增量变更, 用于同步数据内容
func (s *LSPServer) applyIncrementalChanges(uri protocol.DocumentURI, changes []protocol.TextDocumentContentChangeEvent) ([]byte, error) {
	content, err := os.ReadFile(uri.Path())
	if err != nil {
		return nil, err
	}

	for _, change := range changes {
		m := protocol.NewMapper(uri, content)
		if change.Range == nil {
			return nil, fmt.Errorf("%w: unexpected nil range for change", jsonrpc2.ErrInternal)
		}
		start, end, err := m.RangeOffsets(*change.Range)
		if err != nil {
			return nil, err
		}
		if end < start {
			return nil, fmt.Errorf("%w: invalid range for content change", jsonrpc2.ErrInternal)
		}
		var buf bytes.Buffer
		buf.Write(content[:start])
		buf.WriteString(change.Text)
		buf.Write(content[end:])
		content = buf.Bytes()
	}
	return content, nil
}
