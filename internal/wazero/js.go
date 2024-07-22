// 版权 @2022 凹语言 作者。保留所有权利。

package wazero

import (
	"context"
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/wazero"
	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	"wa-lang.org/wa/internal/3rdparty/wazero/imports/walang"
	"wa-lang.org/wa/internal/3rdparty/wazero/sys"
)

const jsModuleName = "syscall_js"

func JsInstantiate(ctx context.Context, rt wazero.Runtime) (api.Closer, error) {
	return rt.NewHostModuleBuilder(jsModuleName).
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

		// 非标准, 仅用于辅助测试
		// func debug_read_file_len(name_ptr, name_len) => i32
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, name_ptr, name_len uint32) uint32 {
			name_bytes, _ := m.Memory().Read(ctx, name_ptr, name_len)
			if fi, err := os.Lstat(string(name_bytes)); err == nil {
				return uint32(fi.Size())
			}
			return 0
		}).
		WithParameterNames("name_ptr", "name_len").
		Export("debug_read_file_len").

		// 非标准, 仅用于辅助测试
		// func debug_read_file_data(name_ptr, name_len, data_ptr)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, name_ptr, name_len, data_ptr, data_len uint32) {
			name_bytes, _ := m.Memory().Read(ctx, name_ptr, name_len)
			data, _ := os.ReadFile(string(name_bytes))
			if len(data) > int(data_len) {
				data = data[:data_len]
			}
			m.Memory().Write(ctx, data_ptr, data)
		}).
		WithParameterNames("name_ptr", "name_len", "data_ptr", "data_len").
		Export("debug_read_file_data").

		// 非标准, 仅用于辅助测试
		// func debug_write_file(name_ptr, name_len, data_ptr, data_len: i32)
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, name_ptr, name_len, data_ptr, data_len uint32) {
			name_bytes, _ := m.Memory().Read(ctx, name_ptr, name_len)
			data_bytes, _ := m.Memory().Read(ctx, data_ptr, data_len)
			os.WriteFile(string(name_bytes), data_bytes, 0666)
		}).
		WithParameterNames("name_ptr", "name_len", "data_ptr", "data_len").
		Export("debug_write_file").

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
