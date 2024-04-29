// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) CodeAction(context.Context, *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) CodeLens(context.Context, *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, fmt.Errorf("TODO")
}
