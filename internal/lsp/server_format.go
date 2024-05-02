// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"os"

	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	s.logger.Println("Formatting:", jsonMarshal(params))

	path := params.TextDocument.URI.Path()
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	code, changed, err := format.File(nil, path, content)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, nil
	}

	m := protocol.NewMapper(params.TextDocument.URI, content)
	rng, err := m.OffsetRange(0, len(content))
	if err != nil {
		return nil, err
	}

	return []protocol.TextEdit{
		{Range: rng, NewText: string(code)},
	}, nil
}
