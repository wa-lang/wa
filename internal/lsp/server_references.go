// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

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
