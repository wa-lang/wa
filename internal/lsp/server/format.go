// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/lsp/event"
	"wa-lang.org/wa/internal/lsp/event/tag"
	"wa-lang.org/wa/internal/lsp/protocol"
)

func (s *server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	ctx, done := event.Start(ctx, "lsp.Server.formatting", tag.URI.Of(params.TextDocument.URI))
	defer done()

	return nil, fmt.Errorf("TODO: Formatting")
}
