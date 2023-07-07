// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"fmt"

	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wazero"
	"wa-lang.org/wazero/api"
	"wa-lang.org/wazero/imports/walang"
	"wa-lang.org/wazero/sys"
)

func MvpInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(config.WaOS_mvp).
		// func waPrintI32(v: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v int32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintI32").

		// func waPrintU32(v: u32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintU32").

		// func waPrintI64(v: i64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v int64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintI64").

		// func waPrintU64(v: u64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintU64").

		// func waPrintF32(v: f32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v float32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintF32").

		// func waPrintF64(v: f64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v float64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("waPrintF64").

		// func waPrintRune(ch: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ch uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprintf(w, "%c", rune(ch))
		}).
		WithParameterNames("ch").
		Export("waPrintRune").

		// func waPuts(ptr: i32, len: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			bytes, _ := m.Memory().Read(ctx, ptr, len)
			fmt.Fprint(w, string(bytes))
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
