// 版权 @2022 凹语言 作者。保留所有权利。

package waruntime

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const envWalang = "wa_js_env"

func WalangInstantiate(ctx context.Context, rt wazero.Runtime) (api.Module, error) {
	return rt.NewHostModuleBuilder(envWalang).
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, pos, len uint32) {
			bytes, _ := m.Memory().Read(ctx, pos, len)
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
		Instantiate(ctx, rt)
}
