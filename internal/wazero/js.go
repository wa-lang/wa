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

func JsInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(config.WaOS_js).
		// func print_bool(v: bool)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			b := v != 0
			fmt.Fprint(w, b)
		}).
		WithParameterNames("v").
		Export("print_bool").

		// func print_i32(v: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v int32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_i32").

		// func print_u32(v: u32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_u32").

		// func print_ptr(v: u32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprintf(w, "0x%x", v)
		}).
		WithParameterNames("v").
		Export("print_ptr").

		// func print_i64(v: i64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v int64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_i64").

		// func print_u64(v: u64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v uint64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_u64").

		// func print_f32(v: f32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v float32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_f32").

		// func print_f64(v: f64)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, v float64) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprint(w, v)
		}).
		WithParameterNames("v").
		Export("print_f64").

		// func print_rune(ch: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ch uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			fmt.Fprintf(w, "%c", rune(ch))
		}).
		WithParameterNames("ch").
		Export("print_rune").

		// func print_str(ptr: i32, len: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) {
			w := walang.ModCallContextSys(m).Stdout()
			bytes, _ := m.Memory().Read(ctx, ptr, len)
			fmt.Fprint(w, string(bytes))
		}).
		WithParameterNames("ptr", "len").
		Export("print_str").

		// func proc_exit(code: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, exitCode uint32) {
			panic(sys.NewExitError(m.Name(), exitCode))
		}).
		WithParameterNames("code").
		Export("proc_exit").

		// Done
		Instantiate(ctx, rt)
}
