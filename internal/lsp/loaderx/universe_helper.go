// 版权 @2024 凹语言 作者。保留所有权利。

package loaderx

import (
	"bytes"
	"fmt"

	"wa-lang.org/wa/internal/lsp/jsonrpc2"
	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *Universe) addWorkspaceFolders(uri protocol.URI) {
	for _, x := range s.WorkspaceFolder {
		if x == uri {
			return
		}
	}
	s.WorkspaceFolder = append(s.WorkspaceFolder, uri)
}

func (s *Universe) removeWorkspaceFolders(uri protocol.URI) {
	folders := s.WorkspaceFolder[:0]
	for _, x := range s.WorkspaceFolder {
		if x != uri {
			folders = append(folders, x)
		}
	}
	s.WorkspaceFolder = folders
}

//WorkspaceFolders

func (s *Universe) changedText(uri protocol.DocumentURI, changes []protocol.TextDocumentContentChangeEvent) ([]byte, error) {
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
func (s *Universe) applyIncrementalChanges(uri protocol.DocumentURI, changes []protocol.TextDocumentContentChangeEvent) ([]byte, error) {
	content := s.Files[uri].Data

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
