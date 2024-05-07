// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (p *LSPServer) CodeAction(ctx context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	p.logger.Println("CodeAction:", jsonMarshal(params))
	return nil, fmt.Errorf("TODO")
}

func (p *LSPServer) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	p.logger.Println("CodeAction:", jsonMarshal(params))
	return nil, fmt.Errorf("TODO")
}
