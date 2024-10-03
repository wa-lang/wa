// 版权 @2024 凹语言 作者。保留所有权利。

package wat2c

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/token"
)

func (p *wat2cWorker) buildCode(w io.Writer) error {
	fmt.Fprintf(w, "// Auto Generated by http://wa-lang.org/wat2c. DONOT EDIT!!!\n\n")

	if p.m.Name != "" {
		fmt.Fprintf(w, "// module %s\n\n", p.m.Name)
	}

	fmt.Fprintf(w, "#include <stdint.h>\n\n")

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

func (p *wat2cWorker) buildMemory(w io.Writer) error {
	if p.m.Memory == nil {
		return nil
	}
	if p.m.Memory.Name != "" {
		fmt.Fprintf(w, "// memory $%s\n", p.m.Memory.Name)
	}
	if max := p.m.Memory.MaxPages; max > 0 {
		fmt.Fprintf(w, "static uint8_t   wasm_memoy[%d*64*1024];\n", max)
		fmt.Fprintf(w, "static int       wasm_memoy_size = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int wasm_memoy_max_pages = %d;\n", max)
		fmt.Fprintf(w, "static const int wasm_memoy_pages = %d;\n", p.m.Memory.Pages)
	} else {
		fmt.Fprintf(w, "static uint8_t   wasm_memoy[%d*64*1024];\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static int       wasm_memoy_size = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int wasm_memoy_max_pages = %d;\n", p.m.Memory.Pages)
		fmt.Fprintf(w, "static const int wasm_memoy_pages = %d;\n", p.m.Memory.Pages)
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
		fmt.Fprintf(w, "static uintptr_t wasm_table[%d];\n", max)
		fmt.Fprintf(w, "static int32_t   wasm_table_size = %d;\n", p.m.Table.Size)
		fmt.Fprintf(w, "static const int wasm_table_max_size = %d;\n", max)
	} else {
		fmt.Fprintf(w, "static uintptr_t wasm_table[%d];\n", p.m.Table.Size)
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
				fmt.Fprintf(w, "static int32_t var_%s = %d;\n", toCName(g.Name), g.I32Value)
			} else {
				fmt.Fprintf(w, "static const int32_t var_%s = %d;\n", toCName(g.Name), g.I32Value)
			}
		case token.I64:
			if g.Mutable {
				fmt.Fprintf(w, "static int64_t var_%s = %d;\n", toCName(g.Name), g.I64Value)
			} else {
				fmt.Fprintf(w, "static const int64_t var_%s = %d;\n", toCName(g.Name), g.I64Value)
			}
		case token.F32:
			if g.Mutable {
				fmt.Fprintf(w, "static float var_%s = %f;\n", toCName(g.Name), g.F32Value)
			} else {
				fmt.Fprintf(w, "static const float var_%s = %f;\n", toCName(g.Name), g.F32Value)
			}
		case token.F64:
			if g.Mutable {
				fmt.Fprintf(w, "static double var_%s = %f;\n", toCName(g.Name), g.F64Value)
			} else {
				fmt.Fprintf(w, "static const double var_%s = %f;\n", toCName(g.Name), g.F64Value)
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

		// 返回值通过栈传递, 返回入栈的个数
		fmt.Fprintf(w, "static int fn_%s(int64_t $result[]", toCName(f.Name))
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				if x.Name != "" {
					fmt.Fprintf(w, ", %v %v", toCType(x.Type), toCName(x.Name))
				} else {
					fmt.Fprintf(w, ", %v $arg%d", toCType(x.Type), i)
				}
			}
		}
		fmt.Fprintf(w, ");\n")
	}
	fmt.Fprintln(w)

	// 函数的实现
	for _, f := range p.m.Funcs {
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
		fmt.Fprintf(w, "static int fn_%s(int64_t $result[]", toCName(f.Name))
		if len(f.Type.Params) > 0 {
			for i, x := range f.Type.Params {
				if x.Name != "" {
					fmt.Fprintf(w, ", %v %v", toCType(x.Type), toCName(x.Name))
				} else {
					fmt.Fprintf(w, ", %v $arg%d", toCType(x.Type), i)
				}
			}
		}
		fmt.Fprintf(w, ") {\n")
		if err := p.buildFunc_body(w, f); err != nil {
			return err
		}
		fmt.Fprintf(w, "}\n\n")
	}

	fmt.Fprintln(w)
	return nil
}