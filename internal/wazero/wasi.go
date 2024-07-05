// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"io"
	"os"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	wasi "wa-lang.org/wa/internal/3rdparty/wazero/imports/wasi_snapshot_preview1"
)

// https://github.com/WebAssembly/WASI/blob/snapshot-01/phases/snapshot/docs.md

const wasiModuleName = "wasi_snapshot_preview1"

func WasiInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return wasi.Instantiate(ctx, rt)
}

func wasi_fdWriteFn(ctx context.Context, mod api.Module, fd, iovs, iovsCount, resultSize uint32) wasi.Errno {
	var err error
	var nwritten uint32
	var writer = os.Stdout

	for i := uint32(0); i < iovsCount; i++ {
		iov := iovs + i*8
		offset, ok := mod.Memory().ReadUint32Le(ctx, iov)
		if !ok {
			return wasi.ErrnoFault
		}
		// Note: emscripten has been known to write zero length iovec. However,
		// it is not common in other compilers, so we don't optimize for it.
		l, ok := mod.Memory().ReadUint32Le(ctx, iov+4)
		if !ok {
			return wasi.ErrnoFault
		}

		var n int
		if writer == io.Discard { // special-case default
			n = int(l)
		} else {
			b, ok := mod.Memory().Read(ctx, offset, l)
			if !ok {
				return wasi.ErrnoFault
			}
			n, err = writer.Write(b)
			if err != nil {
				return wasi.ErrnoIo
			}
		}
		nwritten += uint32(n)
	}
	if !mod.Memory().WriteUint32Le(ctx, resultSize, nwritten) {
		return wasi.ErrnoFault
	}
	return wasi.ErrnoSuccess
}
