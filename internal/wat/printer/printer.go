// 版权 @2024 凹语言 作者。保留所有权利。

// 打印 wat 格式的模块, 不含注释信息

package printer

import (
	"fmt"
	"io"

	"wa-lang.org/wa/internal/wat/ast"
)

func Fprint(output io.Writer, m *ast.Module) error {
	return new(watPrinter).Fprint(output, m)
}

type watPrinter struct {
	m *ast.Module
	w io.Writer

	indent string
}

func (p *watPrinter) Fprint(w io.Writer, m *ast.Module) error {
	p.m = m
	p.w = w

	p.indent = "\t"

	fmt.Fprint(p.w, "(module")
	if p.m.Name != "" {
		fmt.Fprint(p.w, " $", p.m.Name)
	}

	if p.isModuleEmpty() {
		fmt.Fprint(p.w, ")")
		return nil
	} else {
		defer fmt.Fprintln(p.w, ")")
	}

	if err := p.printImport(); err != nil {
		return err
	}
	if err := p.printExport(); err != nil {
		return err
	}

	if err := p.printMemory(); err != nil {
		return err
	}
	if err := p.printTable(); err != nil {
		return err
	}
	if err := p.printTypes(); err != nil {
		return err
	}
	if err := p.printGlobals(); err != nil {
		return err
	}
	if err := p.printFuncs(); err != nil {
		return err
	}
	if err := p.printData(); err != nil {
		return err
	}
	if err := p.printElem(); err != nil {
		return err
	}

	return nil
}
