// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
)

// 打印指令
func PrintInstruction(w io.Writer, pc int64, inst ast.Instruction) {
	fmt.Fprint(w, riscv.AsmSyntax(pc, inst.As, inst.Arg))
}

// 打印汇编格式
func Fprint(w io.Writer, pc int64, f *ast.File) {
	for _, fn := range f.Funcs {
		fmt.Fprintln(w, f.Name)
		for _, inst := range fn.Body.Insts {
			fmt.Fprint(w, riscv.AsmSyntax(pc, inst.As, inst.Arg))
			pc += 4
		}
	}
}
