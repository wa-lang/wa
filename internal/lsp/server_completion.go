// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *LSPServer) Completion(context.Context, *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *LSPServer) ResolveCompletionItem(context.Context, *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, fmt.Errorf("TODO")
}
