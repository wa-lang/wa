// 版权 @2022 凹语言 作者。保留所有权利。

package waruntime

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func WasiInstantiate(ctx context.Context, rt wazero.Runtime) {
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
}
