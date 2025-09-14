// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package printer

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/native/abi"
	"wa-lang.org/wa/internal/native/ast"
	"wa-lang.org/wa/internal/native/riscv"
)

// 打印指令
func PrintInstruction(w io.Writer, inst ast.Instruction) {
	fmt.Fprint(w, riscv.AsmSyntax(inst.As, inst.Arg))
}

// 打印汇编格式
func Fprint(w io.Writer, f *ast.File) error {
	return new(wsPrinter).Fprint(w, f)
}

type wsPrinter struct {
	cpu abi.CPUType

	f *ast.File
	w io.Writer

	indent string
}

func (p *wsPrinter) Fprint(w io.Writer, f *ast.File) error {
	p.f = f
	p.w = w
	p.indent = "\t"

	if err := p.printConsts(); err != nil {
		return err
	}
	if err := p.printGlobals(); err != nil {
		return err
	}
	if err := p.printFuncs(); err != nil {
		return err
	}

	return nil
}
