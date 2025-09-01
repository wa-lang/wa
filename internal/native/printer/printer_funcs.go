// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
)

// func $add(%a:i32, %b:i32, %c:i32) => f64 {
//     local %d: i32 // 局部变量必须先声明, i32 大小的空间
//
//     // 指令
// Loop:
// }

func (p *wsPrinter) printFuncs() error {
	for i, fn := range p.f.Funcs {
		if i > 0 {
			fmt.Fprintln(p.w)
		}

		// 函数参数和返回值定义
		if err := p.printFuncs_sig(fn); err != nil {
			return err
		}

		// 打印函数实现
		if err := p.printFuncs_body(fn); err != nil {
			return err
		}
	}
	return nil
}

// 打印函数签名
// func $add(%a:i32, %b:i32, %c:i32) => f64
func (p *wsPrinter) printFuncs_sig(fn *ast.Func) error {
	fmt.Fprintf(p.w, "func %s(", fn.Name)
	for j, x := range fn.Type.Args {
		if j > 0 {
			fmt.Fprint(p.w, ",")
		}
		fmt.Fprintf(p.w, "%v:%v", x.Name, x.Type)
	}
	fmt.Fprintf(p.w, ")")
	if fn.Type.Return != 0 {
		fmt.Fprintf(p.w, " => %v", fn.Type.Return)
	}
	return nil
}

// 打印函数实现
func (p *wsPrinter) printFuncs_body(fn *ast.Func) error {
	fmt.Fprintf(p.w, "{\n")
	defer fmt.Fprintf(p.w, "}\n")

	// 打印局部变量
	for _, local := range fn.Body.Locals {
		fmt.Fprintf(p.w, "%slocal %v:%v\n", p.indent, local.Name, local.Type)
	}

	// 打印函数指令
	switch p.cpu {
	case abi.RISCV32, abi.RISCV64:
		var pc int64
		for _, inst := range fn.Body.Insts {
			if inst.Label != "" {
				fmt.Fprintf(p.w, "%s:\n", inst.Label)
			}
			fmt.Fprintf(p.w, "%s%s\n", p.indent, riscv.AsmSyntax(pc, inst.As, inst.Arg))
			pc += 4 // TODO: 是否可以省略 pc?
		}
	}

	return nil
}
