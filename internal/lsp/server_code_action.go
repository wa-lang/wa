// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
