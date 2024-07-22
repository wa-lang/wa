// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
)

const unknownModuleName = "unknown"

func UnknownInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(unknownModuleName).
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
