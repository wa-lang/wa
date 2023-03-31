// 版权 @2022 凹语言 作者。保留所有权利。

package waruntime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// https://github.com/WebAssembly/WASI/blob/snapshot-01/phases/snapshot/docs.md

const wasiModuleName = "wasi_snapshot_preview1"

func WasiInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	if true {
		return wasi.Instantiate(ctx, rt)
	}
	return rt.NewHostModuleBuilder(wasiModuleName).
		NewFunctionBuilder().
		WithFunc(wasi_fdWriteFn).
		WithParameterNames("fd", "iovs", "iovsCount", "resultSize").
		Export("fd_write").
		// ----
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, pos, len uint32) {
			bytes, _ := m.Memory().Read(pos, len)
			fmt.Print(string(bytes))
		}).
		WithParameterNames("pos", "len").
		Export("waPuts").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v uint32) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintI32").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, ch uint32) {
			fmt.Printf("%c", rune(ch))
		}).
		WithParameterNames("ch").
		Export("waPrintRune").
		Instantiate(ctx)
}

// errno is neither uint16 nor an alias for parity with wasm.ValueType.
type errno = uint32

const (
	// ErrnoSuccess No error occurred. System call completed successfully.
	errnoSuccess errno = 0

	// errnoFault Bad address.
	errnoFault errno = 21

	// errnoIo I/O error.
	errnoIo errno = 29
)

func wasi_fdWriteFn(ctx context.Context, mod api.Module, fd, iovs, iovsCount, resultSize uint32) errno {
	var err error
	var nwritten uint32
	var writer = os.Stdout

	for i := uint32(0); i < iovsCount; i++ {
		iov := iovs + i*8
		offset, ok := mod.Memory().ReadUint32Le(iov)
		if !ok {
			return errnoFault
		}
		// Note: emscripten has been known to write zero length iovec. However,
		// it is not common in other compilers, so we don't optimize for it.
		l, ok := mod.Memory().ReadUint32Le(iov + 4)
		if !ok {
			return errnoFault
		}

		var n int
		if writer == io.Discard { // special-case default
			n = int(l)
		} else {
			b, ok := mod.Memory().Read(offset, l)
			if !ok {
				return errnoFault
			}
			n, err = writer.Write(b)
			if err != nil {
				return errnoIo
			}
		}
		nwritten += uint32(n)
	}
	if !mod.Memory().WriteUint32Le(resultSize, nwritten) {
		return errnoFault
	}
	return errnoSuccess
}
