package experimental

import (
	"context"
	"io/fs"

	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	internalfs "wa-lang.org/wa/internal/3rdparty/wazero/internalx/sys"
)

// WithFS overrides fs.FS in the context-based manner. Caller needs to take
// responsibility for closing the filesystem.
//
// Note: This has the same effect as the same function on wazero.ModuleConfig.
func WithFS(ctx context.Context, fs fs.FS) (context.Context, api.Closer) {
	if fs == nil {
		fs = internalfs.EmptyFS
	}
	fsCtx := internalfs.NewFSContext(fs)
	return context.WithValue(ctx, internalfs.FSKey{}, fsCtx), fsCtx
}
