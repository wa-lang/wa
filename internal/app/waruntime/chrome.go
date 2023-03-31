// 版权 @2022 凹语言 作者。保留所有权利。

package waruntime

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"wa-lang.org/wa/internal/config"
)

func ChromeInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(config.WaOS_chrome).Instantiate(ctx)
}
