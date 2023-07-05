// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/sys"
	"wa-lang.org/wa/internal/config"
)

func MvpInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(config.WaOS_mvp).
		// func waPrintI32(v: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v int32) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintI32").

		// func waPrintU32(v: u32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v uint32) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintU32").

		// func waPrintI64(v: i64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v int64) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintI64").

		// func waPrintU64(v: u64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v uint64) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintU64").

		// func waPrintF32(v: f32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v float32) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintF32").

		// func waPrintF64(v: f64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, v float64) {
			fmt.Print(v)
		}).
		WithParameterNames("v").
		Export("waPrintF64").

		// func waPrintRune(ch: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, ch uint32) {
			fmt.Printf("%c", rune(ch))
		}).
		WithParameterNames("ch").
		Export("waPrintRune").

		// func waPuts(ptr: i32, len: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) {
			bytes, _ := m.Memory().Read(ctx, ptr, len)
			fmt.Print(string(bytes))
		}).
		WithParameterNames("ptr", "len").
		Export("waPuts").

		// func ProcExit(code: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, exitCode uint32) {
			panic(sys.NewExitError(m.Name(), exitCode))
		}).
		WithParameterNames("code").
		Export("proc_exit").

		// Done
		Instantiate(ctx, rt)
}
