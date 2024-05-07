// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	"wa-lang.org/wa/internal/lsp/protocol"
)

func (p *LSPServer) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	p.logger.Println("DidOpen:", jsonMarshal(params))
	p.fileMap[params.TextDocument.URI.Path()] = params.TextDocument.Text

	p.syncFile.SaveFile(
		params.TextDocument.URI.Path()+"."+strconv.Itoa(int(params.TextDocument.Version)),
		params.TextDocument.Text,
	)
	return nil
}

func (s *LSPServer) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	s.logger.Println("DidSave:", jsonMarshal(params))
	return nil
}

func (s *LSPServer) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.logger.Println("DidClose:", jsonMarshal(params))
	return nil
}

func (p *LSPServer) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	p.logger.Println("DidChange:", jsonMarshal(params))

	if !strings.HasSuffix(string(params.TextDocument.URI), ".wa") {
		return nil
	}

	text, err := p.changedText(params.TextDocument.URI, params.ContentChanges)
	if err != nil {
		return err
	}

	p.logger.Println("DidChange.text:", string(text))
	p.syncFile.SaveFile(
		params.TextDocument.URI.Path()+"."+strconv.Itoa(int(params.TextDocument.Version)),
		string(text),
	)

	p.fileMap[params.TextDocument.URI.Path()] = string(text)
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
	content := []byte(s.fileMap[uri.Path()])

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
