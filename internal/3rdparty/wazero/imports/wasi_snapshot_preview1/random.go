package wasi_snapshot_preview1

import (
	"context"
	"io"

	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/wasm"
)

const functionRandomGet = "random_get"

// randomGet is the WASI function named functionRandomGet which writes random
// data to a buffer.
//
// # Parameters
//
//   - buf: api.Memory offset to write random values
//   - bufLen: size of random data in bytes
//
// Result (Errno)
//
// The return value is ErrnoSuccess except the following error conditions:
//   - ErrnoFault: `buf` or `bufLen` point to an offset out of memory
//   - ErrnoIo: a file system error
//
// For example, if underlying random source was seeded like
// `rand.NewSource(42)`, we expect api.Memory to contain:
//
//	                   bufLen (5)
//	          +--------------------------+
//	          |                        	 |
//	[]byte{?, 0x53, 0x8c, 0x7f, 0x96, 0xb1, ?}
//	    buf --^
//
// See https://github.com/WebAssembly/WASI/blob/snapshot-01/phases/snapshot/docs.md#-random_getbuf-pointeru8-bufLen-size---errno
var randomGet = &wasm.HostFunc{
	ExportNames: []string{functionRandomGet},
	Name:        functionRandomGet,
	ParamTypes:  []api.ValueType{i32, i32},
	ParamNames:  []string{"buf", "buf_len"},
	ResultTypes: []api.ValueType{i32},
	Code: &wasm.Code{
		IsHostFunction: true,
		GoFunc:         wasiFunc(randomGetFn),
	},
}

func randomGetFn(ctx context.Context, mod api.Module, params []uint64) Errno {
	sysCtx := mod.(*wasm.CallContext).Sys
	randSource := sysCtx.RandSource()
	buf, bufLen := uint32(params[0]), uint32(params[1])

	randomBytes, ok := mod.Memory().Read(ctx, buf, bufLen)
	if !ok { // out-of-range
		return ErrnoFault
	}

	// We can ignore the returned n as it only != byteCount on error
	if _, err := io.ReadAtLeast(randSource, randomBytes, int(bufLen)); err != nil {
		return ErrnoIo
	}

	return ErrnoSuccess
}
