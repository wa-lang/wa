// 版权 @2024 凹语言 作者。保留所有权利。

package printer

import (
	"fmt"

	"wa-lang.org/wa/internal/wat/ast"
	"wa-lang.org/wa/internal/wat/token"
)

func (p *watPrinter) printImport() error {
	if len(p.m.Imports) > 0 {
		fmt.Fprintln(p.w)
	}
	for _, importSpec := range p.m.Imports {
		switch importSpec.ObjKind {
		case token.GLOBAL:
			panic("TODO")

		case token.FUNC:
			fmt.Fprint(p.w, p.indent)
			fmt.Fprintf(p.w, "(import %q %q", importSpec.ObjModule, importSpec.ObjName)
			p.printImport_func(importSpec)
			fmt.Fprint(p.w, ")\n")

		case token.MEMORY:
			fmt.Fprint(p.w, p.indent)
			fmt.Fprintf(p.w, "(import %q %q", importSpec.ObjModule, importSpec.ObjName)
			fmt.Fprintf(p.w, " (memory")
			if s := importSpec.Memory.Name; s != "" {
				fmt.Fprintf(p.w, " $"+s)
			}
			fmt.Fprintf(p.w, " %d", importSpec.Memory.Pages)
			if importSpec.Memory.MaxPages != 0 {
				fmt.Fprintf(p.w, " %d", importSpec.Memory.MaxPages)
			}
			fmt.Fprint(p.w, "))\n")

		case token.TABLE:
			panic("TODO")

		default:
			panic("unreachable")
		}
	}
	return nil
}

func (p *watPrinter) printImport_func(importSpec *ast.ImportSpec) {
	fmt.Fprintf(p.w, " (func %s", p.identOrIndex(importSpec.FuncName))

	fnType := importSpec.FuncType
	if len(fnType.Params) > 0 {
		for _, x := range fnType.Params {
			fmt.Fprintf(p.w, " (param %v)", x.Type)
		}
	}
	if len(fnType.Results) > 0 {
		fmt.Fprintf(p.w, " (result")
		for _, x := range fnType.Results {
			fmt.Fprintf(p.w, " %v", x)
		}
		fmt.Fprint(p.w, ")")
	}

	fmt.Fprint(p.w, ")")
	return
}
