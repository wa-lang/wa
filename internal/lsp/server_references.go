// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (p *LSPServer) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	p.logger.Println("References:", jsonMarshal(params))

	return nil, fmt.Errorf("TODO")
}
