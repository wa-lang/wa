// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) Definition(context.Context, *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, fmt.Errorf("TODO")
}
