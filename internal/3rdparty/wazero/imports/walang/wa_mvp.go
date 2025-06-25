// 版权 @2023 凹语言 作者。保留所有权利。

// github.com/tetratelabs/wazero/imports/walang

package walang

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/sys"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/wasm"
)

func ModCallContextSys(m api.Module) *sys.Context {
	return m.(*wasm.CallContext).Sys
}
