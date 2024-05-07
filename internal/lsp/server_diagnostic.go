// 版权 @2024 凹语言 作者。保留所有权利。

package lsp

import (
	"context"

	"wa-lang.org/wa/internal/lsp/protocol"
)

func (p *LSPServer) publishFileDiagnosticsLocked(ctx context.Context, uri protocol.DocumentURI, version int32) error {
	// TODO: 服务器异步通知IDE诊断结果, 有文件版本号信息, 后台进行
	p.client.PublishDiagnostics(ctx, nil)
	return nil
}
