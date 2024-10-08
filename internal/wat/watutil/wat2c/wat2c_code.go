// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

//go:embed _math_x.c
var math_x_code string

func (p *wat2cWorker) buildCode(w io.Writer) error {
	fmt.Fprintf(w, "// Auto Generated by https://wa-lang.org/wa/wat2c. DONOT EDIT!!!\n\n")

	if p.m.Name != "" {
		fmt.Fprintf(w, "// module %s\n\n", p.m.Name)
	}

	fmt.Fprintf(w, "#include <stdint.h>\n")
	fmt.Fprintf(w, "#include <string.h>\n")
	fmt.Fprintf(w, "#include <math.h>\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "typedef uint8_t   u8_t;\n")
	fmt.Fprintf(w, "typedef int8_t    i8_t;\n")
	fmt.Fprintf(w, "typedef uint16_t  u16_t;\n")
	fmt.Fprintf(w, "typedef int16_t   i16_t;\n")
	fmt.Fprintf(w, "typedef uint32_t  u32_t;\n")
	fmt.Fprintf(w, "typedef int32_t   i32_t;\n")
	fmt.Fprintf(w, "typedef uint64_t  u64_t;\n")
	fmt.Fprintf(w, "typedef int64_t   i64_t;\n")
	fmt.Fprintf(w, "typedef float     f32_t;\n")
	fmt.Fprintf(w, "typedef double    f64_t;\n")
	fmt.Fprintf(w, "typedef uintptr_t ref_t;\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "typedef union val_t {\n")
	fmt.Fprintf(w, "  i64_t i64;\n")
	fmt.Fprintf(w, "  f64_t f64;\n")
	fmt.Fprintf(w, "  i32_t i32;\n")
	fmt.Fprintf(w, "  f32_t f32;\n")
	fmt.Fprintf(w, "  ref_t ref;\n")
	fmt.Fprintf(w, "} val_t;\n\n")

	if err := p.buildImport(w); err != nil {
		return err
	}

	if err := p.buildMemory(w); err != nil {
		return err
	}
	if err := p.buildTable(w); err != nil {
		return err
	}

	if err := p.buildGlobal(w); err != nil {
		return err
	}
	if err := p.buildFuncs(w); err != nil {
		return err
	}

	return nil
}

func (p *wat2cWorker) buildImport(w io.Writer) error {
	if len(p.m.Imports) == 0 {
		return nil
	}
	for _, importSpec := range p.m.Imports {
		if importSpec.ObjKind != token.FUNC {
			continue
		}

		fnName := importSpec.FuncName
		fnType := importSpec.FuncType

		fmt.Fprintf(w, "// import %s.%s => func $%s",
			importSpec.ObjModule, importSpec.ObjName,
			fnName,
		)

		if len(fnType.Params) > 0 {
			for i, x := range fnType.Params {
				if x.Name != "" {
					fmt.Fprintf(w, " (param $%s %v)", x.Name, x.Type)
				} else {
					fmt.Fprintf(w, " (param $%d %v)", i, x.Type)
				}
			}
		}
		if len(fnType.Results) > 0 {
			fmt.Fprintf(w, " (result")
			for _, x := range fnType.Results {
				fmt.Fprintf(w, " %v", x)
			}
			fmt.Fprint(w, ")")
		}
		fmt.Fprintln(w)

		// 返回值通过栈传递
		fmt.Fprintf(w, "extern int fn_%s(val_t $result[]", toCName(fnName))
		if len(fnType.Params) > 0 {
			for i, x := range fnType.Params {
				if x.Name != "" {
					fmt.Fprintf(w, ", val_t %v", toCName(x.Name))
				} else {
					fmt.Fprintf(w, ", val_t $arg%d", i)
				}
			}
		}
		fmt.Fprintf(w, ");\n")
	}

	fmt.Fprintln(w)
	return nil
}

func (p *wat2cWorker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}
	if p.m.Memory.Name != "" {
		fmt.Fprintf(w, "// memory $%s\n", p.m.Memory.Name)
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		fmt.Fprintf(w, "static uint8_t       wasm_memoy[%d*64*1024];\n", max)
		fmt.Fprintf(w, "static int32_t       wasm_memoy_size = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int32_t wasm_memoy_max_pages = %d;\n", max)
		fmt.Fprintf(w, "static const int32_t wasm_memoy_pages = %d;\n", p.m.Memory.Pages)
	} else {
		fmt.Fprintf(w, "static uint8_t       wasm_memoy[%d*64*1024];\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static int32_t       wasm_memoy_size = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int32_t wasm_memoy_max_pages = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int32_t wasm_memoy_pages = %d;\n", p.m.Memory.Pages)
	}
	fmt.Fprintln(w)
	return nil
}

func (p *wat2cWorker) buildTable(w io.Writer) error {
	if p.m.Table == nil {
		return nil
	}
	if p.m.Table.Type != token.FUNCREF {
		return fmt.Errorf("unsupported table type: %s", p.m.Table.Type)
	}

	if p.m.Table.Name != "" {
		fmt.Fprintf(w, "// table $%s\n", p.m.Table.Name)
	}
	if max := p.m.Table.MaxSize; max > 0 {
		fmt.Fprintf(w, "static ref_t     wasm_table[%d];\n", max)
		fmt.Fprintf(w, "static int32_t   wasm_table_size = %d;\n", p.m.Table.Size)
		fmt.Fprintf(w, "static const int wasm_table_max_size = %d;\n", max)
	} else {
		fmt.Fprintf(w, "static ref_t     wasm_table[%d];\n", p.m.Table.Size)
		fmt.Fprintf(w, "static int32_t   wasm_table_size = %d;\n", p.m.Table.Size)
		fmt.Fprintf(w, "static const int wasm_table_max_size = %d;\n", p.m.Table.Size)
	}
	fmt.Fprintln(w)
	return nil
}

func (p *wat2cWorker) buildGlobal(w io.Writer) error {
	if len(p.m.Globals) == 0 {
		return nil
	}
	for _, g := range p.m.Globals {
		fmt.Fprintf(w, "// global $%s: %v\n", g.Name, g.Type)
		switch g.Type {
		case token.I32:
			if g.Mutable {
				fmt.Fprintf(w, "static val_t var_%s = {.i32=%d};\n", toCName(g.Name), g.I32Value)
			} else {
				fmt.Fprintf(w, "static const val_t var_%s = {.i32=%d};\n", toCName(g.Name), g.I32Value)
			}
		case token.I64:
			if g.Mutable {
				fmt.Fprintf(w, "static val_t var_%s = {.i64=%d};\n", toCName(g.Name), g.I64Value)
			} else {
				fmt.Fprintf(w, "static const val_t var_%s = {.i64=%d};\n", toCName(g.Name), g.I64Value)
			}
		case token.F32:
			if g.Mutable {
				fmt.Fprintf(w, "static val_t var_%s = {.f32=%f|;\n", toCName(g.Name), g.F32Value)
			} else {
				fmt.Fprintf(w, "static const val_t var_%s = {.f32=%f};\n", toCName(g.Name), g.F32Value)
			}
		case token.F64:
			if g.Mutable {
				fmt.Fprintf(w, "static val_t var_%s = {.f64=%f};\n", toCName(g.Name), g.F64Value)
			} else {
				fmt.Fprintf(w, "static const val_t var_%s = {.f64=%f};\n", toCName(g.Name), g.F64Value)
			}
		default:
			return fmt.Errorf("unsupported global type: %s", g.Type)
		}
	}
	fmt.Fprintln(w)
	return nil
}

func (p *wat2cWorker) buildFuncs(w io.Writer) error {
	if len(p.m.Funcs) == 0 {
		return nil
	}

	// 函数声明
	for _, f := range p.m.Funcs {
		fmt.Fprintf(w, "// func $%s", f.Name)
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				if x.Name != "" {
					fmt.Fprintf(w, " (param $%s %v)", x.Name, x.Type)
				} else {
					fmt.Fprintf(w, " (param $%d %v)", i, x.Type)
				}
			}
		}
		if len(f.Type.Results) > 0 {
			fmt.Fprintf(w, " (result")
			for _, x := range f.Type.Results {
				fmt.Fprintf(w, " %v", x)
			}
			fmt.Fprint(w, ")")
		}
		fmt.Fprintln(w)

		// 返回值类型
		fmt.Fprintf(w, "typedef struct {")
		for i := 0; i < len(f.Type.Results); i++ {
			if i == 0 {
				fmt.Fprintf(w, " val_t $R%d", i)
			} else {
				fmt.Fprintf(w, ", $R%d", i)
			}
			if i == len(f.Type.Results)-1 {
				fmt.Fprintf(w, "; ")
			}
		}
		fmt.Fprintf(w, "} fn_%s_ret_t;\n", toCName(f.Name))

		// 返回值通过栈传递, 返回入栈的个数
		fmt.Fprintf(w, "static fn_%s_ret_t fn_%s(", toCName(f.Name), toCName(f.Name))
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				if i > 0 {
					fmt.Fprintf(w, ", ")
				}
				if x.Name != "" {
					fmt.Fprintf(w, "val_t %v", toCName(x.Name))
				} else {
					fmt.Fprintf(w, "val_t $arg%d", i)
				}
			}
		}
		fmt.Fprintf(w, ");\n")
	}
	fmt.Fprintln(w)

	// 函数的实现
	var funcImplBuf bytes.Buffer
	var wBackup = w

	w = &funcImplBuf
	for _, f := range p.m.Funcs {
		p.localNames = nil
		p.localTypes = nil
		p.scopeLabels = nil

		fmt.Fprintf(w, "// func %s", f.Name)
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				if x.Name != "" {
					fmt.Fprintf(w, " (param $%s %v)", x.Name, x.Type)
				} else {
					fmt.Fprintf(w, " (param $%d %v)", i, x.Type)
				}
			}
		}
		if len(f.Type.Results) > 0 {
			fmt.Fprintf(w, " (result")
			for _, x := range f.Type.Results {
				fmt.Fprintf(w, " %v", x)
			}
			fmt.Fprint(w, ")")
		}
		fmt.Fprintln(w)

		// 返回值通过栈传递, 返回入栈的个数
		fmt.Fprintf(w, "static fn_%s_ret_t fn_%s(", toCName(f.Name), toCName(f.Name))
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				var argName string
				if x.Name != "" {
					argName = toCName(x.Name)
				} else {
					argName = fmt.Sprintf("$arg%d", i)
				}

				p.localNames = append(p.localNames, argName)
				p.localTypes = append(p.localTypes, x.Type)

				if i > 0 {
					fmt.Fprint(w, ", ")
				}
				fmt.Fprintf(w, "val_t %v", argName)
			}
		}
		fmt.Fprintf(w, ") {\n")
		if err := p.buildFunc_body(w, f); err != nil {
			return err
		}
		fmt.Fprintf(w, "}\n\n")
	}

	// 恢复输出流
	w = wBackup

	// 扩展数学函数
	if p.useMathX {
		fmt.Fprintln(w, math_x_code)
	}

	// 复制函数实现
	{
		code := bytes.TrimSpace(funcImplBuf.Bytes())
		if _, err := w.Write(code); err != nil {
			return err
		}
	}

	fmt.Fprintln(w)
	return nil
}
