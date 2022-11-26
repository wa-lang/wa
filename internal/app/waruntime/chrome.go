// 版权 @2022 凹语言 作者。保留所有权利。

package waruntime

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/wa-lang/wa/internal/config"
)

func ChromeInstantiate(ctx context.Context, rt wazero.Runtime) (api.Module, error) {
	return rt.NewHostModuleBuilder(config.WaOS_Chrome).Instantiate(ctx, rt)
}
